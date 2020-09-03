package main

import (
	"fmt"
	"net/http"

	"github.com/Biubiubiuuuu/alipay/pay"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// 登录
	r.POST("alipay/login", AlipayLogin)
	// 支付
	r.POST("alipay/pay", Pay)
	// 退款
	r.POST("alipay/refund", Refund)
	// 支付回调 or 退款回调
	r.GET("alipay/notify", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})
	r.POST("alipay/notify", Notify)
	r.Run("127.0.0.1:8030")
	//r.Run("172.17.48.197:8030")
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
	authResp, err := pay.AuthToken(req.AuthCode)
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
	tradeCreateResp, err := pay.CreatePay(req.BuyerID, req.Subject, req.Body, req.TotalAmount)
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
	refundResp, err := pay.RefundPay(req.OutTradeNo, req.TradeNo, req.RefundReason, req.RefundAmount)
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
	var req pay.NotifyPayReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}
	// 验签失败
	if !pay.NotifyPay(req) {
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
