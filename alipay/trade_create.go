package alipay

import (
	"crypto"
	"encoding/json"
	"io/ioutil"

	"github.com/Biubiubiuuuu/go-pay/config"
	"github.com/lexkong/log"
)

// TradeCreateReq 统一订单创建请求参数
type TradeCreateReq struct {
	OutTradeNo         string        `json:"out_trade_no,omitempy"`
	SellerID           string        `json:"seller_id,omitempy"`
	TotalAmount        float64       `json:"total_amount,omitempy"`
	DiscountableAmount string        `json:"discountable_amount,omitempy"`
	Subject            string        `json:"subject,omitempy"`
	Body               string        `json:"body,omitempy"`
	BuyerID            string        `json:"buyer_id,omitempy"`
	GoodsDetail        []GoodsDetail `json:"goods_detail,omitempy"`
	ProductCode        string        `json:"product_code,omitempy"`
	OperatorID         string        `json:"operator_id,omitempy"`
	StoreID            string        `json:"store_id,omitempy"`
	TerminalID         string        `json:"terminal_id,omitempy"`
	ExtendParams       ExtendParams  `json:"extend_params,omitempy"`
	TimeoutExpress     string        `json:"timeout_express,omitempy"`
}

// TradeCreateResp 统一订单创建公共响应返回参数
type TradeCreateResp struct {
	CreatedResp TradeCreateRespCreatedResp `json:"alipay_trade_create_response"`
	CommonRspSign
}

// TradeCreateRespCreatedResp 统一订单创建响应返回参数
type TradeCreateRespCreatedResp struct {
	CommonRsp
	OutTradeNo string `json:"out_trade_no"`
	TradeNo    string `json:"trade_no"`
}

// CreatePay 支付宝统一下单
//  alipay.trade.create(统一收单交易创建接口)
//  user_id 买家支付宝账号ID
//  subject 订单标题
//  body    订单描述内容
//  amount  订单金额
func CreatePay(userID, subject, body string, amount float64) (TradeCreateResp, error) {
	req := TradeCreateReq{
		OutTradeNo:  CreateAlipayOrderID(),
		TotalAmount: amount,
		Subject:     subject,
		Body:        body,
		BuyerID:     userID,
	}
	var resp TradeCreateResp
	bizContent, err := json.Marshal(req)
	if err != nil {
		log.Errorf(err, "rsp body unmarshal failed.")
		return resp, err
	}
	commonReq := CommonReq{
		AppID:        config.AlipayAppID,
		Method:       "alipay.trade.create",
		Format:       "JSON",
		Charset:      "utf-8",
		SignType:     "RSA2",
		Timestamp:    GetTimestamp(),
		Version:      "1.0",
		AppAuthToken: "",
		BizContent:   string(bizContent),
		NotifyURL:    config.AlipayNotifyURL,
	}
	signURLValues := GetSignStr(commonReq)
	signStr := URLValues(signURLValues, false)
	sign := Rsa2Sign(signStr, config.AlipayPrivateKey, crypto.SHA256)
	commonReq.Sign = sign
	postDataURLValues := GetSignStr(commonReq)
	postData := URLValues(postDataURLValues, true)
	url := config.AlipayURL
	log.Infof("request alipay order create url:%s,req:%+v", url, postData)
	rsp, err := HTTPPost(url, postData, "application/x-www-form-urlencoded; charset=utf-8")
	defer rsp.Body.Close()
	if err != nil {
		log.Errorf(err, "request alipay order create failed.")
		return resp, err
	}
	rspBody, err := ioutil.ReadAll(rsp.Body)
	err = json.Unmarshal(rspBody, &resp)
	if err != nil {
		log.Errorf(err, "rsp body unmarshal failed.")
		return resp, err
	}
	log.Infof("alipay order create resp body:%+v", resp)
	return resp, nil
}
