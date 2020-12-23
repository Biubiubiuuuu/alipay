package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Biubiubiuuuu/go-pay/alipay"
	"github.com/Biubiubiuuuu/go-pay/config"
	"github.com/Biubiubiuuuu/go-pay/wxpay"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	uuid "github.com/satori/go.uuid"
)

func main() {
	r := gin.Default()
	r.POST("alipay/login", AlipayLogin)
	r.POST("alipay/pay", Pay)
	r.POST("alipay/refund", Refund)
	r.POST("alipay/notify", Notify)
	r.POST("alipay/freeze/order", AlipayFreezeCreateOrder)
	r.POST("alipay/freeze/notify", FreezeNotify)
	r.POST("alipay/freeze/unfreeze/order", UnFreeze)
	r.POST("wxpay/score/order/pay", WxPayScoreV3CreateOrder)
	// 其他待补充
	r.Run("127.0.0.1:8030")
}

// AlipayLoginReq 支付宝小程序授权请求参数
type AlipayLoginReq struct {
	AuthCode string `json:"auth_code"`
}

// AlipayLoginResp 支付宝小程序授权返回参数
type AlipayLoginResp struct {
	UserID       string `json:"user_id"`       // 支付宝用户的唯一userid
	AccessToken  string `json:"access_token"`  // 访问令牌
	RefreshToken string `json:"refresh_token"` // 刷新令牌
}

// CreatePayReq 支付宝创建订单请求参数
type CreatePayReq struct {
	TotalAmount float64 `json:"total_amount,omitempy"`
	Subject     string  `json:"subject,omitempy"`
	Body        string  `json:"body,omitempy"`
	BuyerID     string  `json:"buyer_id,omitempy"`
}

// RefundPayReq 支付宝退款请求参数
type RefundPayReq struct {
	OutTradeNo   string  `json:"out_trade_no"`  // 商户订单号,和trade_no不能同时为空
	TradeNo      string  `json:"trade_no"`      // 支付宝交易号,和trade_no不能同时为空
	RefundReason string  `json:"refund_reason"` // 退款的原因说明
	RefundAmount float64 `json:"refund_amount"` // 需要退款的金额
}

// AlipayLogin 支付宝小程序用户授权登录
func AlipayLogin(c *gin.Context) {
	var resp AlipayLoginResp
	var req AlipayLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}
	authResp, err := alipay.AuthToken(req.AuthCode)
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}
	if authResp.OauthTokenResp.UserID == "" {
		c.JSON(http.StatusOK, authResp)
		return
	}
	resp.UserID = authResp.OauthTokenResp.UserID
	resp.AccessToken = authResp.OauthTokenResp.AccessToken
	resp.RefreshToken = authResp.OauthTokenResp.RefreshToken
	c.JSON(http.StatusOK, resp)
}

func Pay(c *gin.Context) {
	var req CreatePayReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}
	tradeCreateResp, err := alipay.CreatePay(req.BuyerID, req.Subject, req.Body, req.TotalAmount)
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}
	// 统一下单失败
	if tradeCreateResp.CreatedResp.Code != "10000" {
		c.JSON(http.StatusOK, tradeCreateResp)
		return
	}
	// 创建成功
	// 处理业务逻辑
	c.JSON(http.StatusOK, tradeCreateResp)
}

func Refund(c *gin.Context) {
	var req RefundPayReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}
	refundResp, err := alipay.RefundPay(req.OutTradeNo, req.TradeNo, req.RefundReason, req.RefundAmount)
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}
	// 退款失败
	if refundResp.TradeRefundResp.Code != "10000" {
		c.JSON(http.StatusOK, refundResp)
		return
	}
	// 退款成功
	// 处理业务逻辑
	c.JSON(http.StatusOK, refundResp)
	return
}

func Notify(c *gin.Context) {
	var req alipay.AlipayNotifyPayOrRefund
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}
	// 验签失败
	if !alipay.NotifyPay(req) {
		c.JSON(http.StatusOK, "签名失败")
		return
	}
	// 业务流水数据判断 校验数据正确性
	// out_trade_no
	// total_amount
	// seller_id
	// app_id
	fmt.Println("SUCCESS")
	// 验签成功 写入SUCCESS
	c.JSON(http.StatusOK, "SUCCESS")
	return
}

type AlipayFreezeCreateOrderResponse struct {
	OutOrderNo string `json:"out_order_no"`
	PostData   string `json:"post_data"` // 预授权请求参数
}

// alipay.fund.auth.order.app.freeze(线上资金授权冻结接口)
// https://opendocs.alipay.com/apis/api_28/alipay.fund.auth.order.app.freeze
func AlipayFreezeCreateOrder(c *gin.Context) {
	var resp AlipayFreezeCreateOrderResponse
	var desposit float64
	out_order_no, returnStr, err := alipay.FreezeCreate("设备码", desposit)
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}
	resp.OutOrderNo = out_order_no
	resp.PostData = returnStr
	c.JSON(http.StatusOK, resp)
	// 业务处理
	return
}

// 资金冻结授权回调
func FreezeNotify(c *gin.Context) {
	req := alipay.FreezeNotifyReq{
		AppID:                     c.PostForm("app_id"),
		Version:                   c.PostForm("version"),
		AuthAppID:                 c.PostForm("auth_app_id"),
		GmtCreate:                 c.PostForm("gmt_create"),
		Charset:                   c.PostForm("charset"),
		RestCreditAmount:          c.PostForm("rest_credit_amount"),
		NotifyType:                c.PostForm("notify_type"),
		NotifyID:                  c.PostForm("notify_id"),
		NotifyTime:                c.PostForm("notify_time"),
		SignType:                  c.PostForm("sign_type"),
		Sign:                      c.PostForm("sign"),
		AuthNo:                    c.PostForm("auth_no"),
		OutOrderNo:                c.PostForm("out_order_no"),
		OperationID:               c.PostForm("operation_id"),
		OutRequestNo:              c.PostForm("out_request_no"),
		OperationType:             c.PostForm("operation_type"),
		Amount:                    c.PostForm("amount"),
		Status:                    c.PostForm("status"),
		GmtTrans:                  c.PostForm("gmt_trans"),
		PayerLogonID:              c.PostForm("payer_logon_id"),
		PayerUserID:               c.PostForm("payer_user_id"),
		TotalFreezeAmount:         c.PostForm("total_freeze_amount"),
		TotalUnfreezeAmount:       c.PostForm("total_unfreeze_amount"),
		TotalPayAmount:            c.PostForm("total_pay_amount"),
		RestAmount:                c.PostForm("rest_amount"),
		CreditAmount:              c.PostForm("credit_amount"),
		FundAmount:                c.PostForm("fund_amount"),
		TotalFreezeCreditAmount:   c.PostForm("total_freeze_credit_amount"),
		TotalFreezeFundAmount:     c.PostForm("total_freeze_fund_amount"),
		TotalPayCreditAmount:      c.PostForm("total_pay_credit_amount"),
		TotalPayFundAmount:        c.PostForm("total_pay_fund_amount"),
		TotalUnfreezeCreditAmount: c.PostForm("total_unfreeze_credit_amount"),
		TotalUnfreezeFundAmount:   c.PostForm("total_unfreeze_fund_amount"),
		RestFundAmount:            c.PostForm("rest_fund_amount"),
		PreAuthType:               c.PostForm("pre_auth_type"),
	}

	// 验签失败
	if !alipay.FreezeNotify(req) {
		fmt.Println("ERROR")
		return
	}
	// 回复处理成功
	fmt.Println("success")
	c.JSON(http.StatusOK, "SUCCESS")
	return
}

type AlipayFreezeOrderQueryResponse struct {
	Status string `json:"status"`
}

func UnFreeze(c *gin.Context) {
	var resp AlipayFreezeOrderQueryResponse
	out_order_no := c.Query("out_order_no")
	respQuy, _ := alipay.FreezeOrderQuery(out_order_no)
	respData, err := alipay.UnFreeze(respQuy.Data.AuthNo, respQuy.Data.OutRequestNo, 99)
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}
	if respData.Data.Code != "10000" {
		c.JSON(http.StatusOK, respData.Data.Msg)
		return
	}
	resp.Status = respData.Data.Status
	c.JSON(http.StatusOK, resp)
	return
}

// 微信信用分创建订单返回参数
type WxPayScoreV3CreateOrderResponse struct {
	OrderID   string `json:"order_id"`  // 服务单号
	MchID     string `json:"mch_id"`    // 商户号
	Package   string `json:"package"`   // 扩展字符串 创建订单中返回
	Timestamp string `json:"timestamp"` // 时间戳
	NonceStr  string `json:"nonce_str"` // 随机字符串
	SignType  string `json:"sign_type"` // 签名方式
	Sign      string `json:"sign"`      // 签名
}

// 微信信用分创建订单
func WxPayScoreV3CreateOrder(c *gin.Context) {
	var sendResp WxPayScoreV3CreateOrderResponse
	var desposit float64
	resp, err := wxpay.CreateScoreOrder("微信用户唯一标识", "商品描述", "位置", desposit)
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}
	log.Infof("weixin zhifufen create order response struct:[%+v].", resp)
	if resp.Code != "" && resp.Code != "ORDER_DONE" {
		c.JSON(http.StatusOK, resp.Message)
		return
	}
	if resp.State != "CREATED" {
		c.JSON(http.StatusOK, resp.Message)
		return
	}
	u4 := uuid.NewV4()
	sendResp = WxPayScoreV3CreateOrderResponse{
		MchID:     config.WxpayMchID,
		Package:   resp.Package,
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
		NonceStr:  u4.String()[1:32],
		SignType:  "HMAC-SHA256",
		Sign:      "",
	}
	urlValue := alipay.StructToURLVal(sendResp)
	signStr := alipay.URLValues(urlValue, false)
	signStr = fmt.Sprintf("%s&key=%s", signStr, config.WxpayMchKey)
	sign := wxpay.GetHmacSha256Encoding(signStr, config.WxpayMchKey)
	sendResp.Sign = sign
	sendResp.OrderID = resp.OrderID
	c.JSON(http.StatusOK, sendResp)
	// 业务处理
	return
}
