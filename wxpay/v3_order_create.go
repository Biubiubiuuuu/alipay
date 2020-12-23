package wxpay

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/Biubiubiuuuu/go-pay/config"
	"github.com/lexkong/log"
)

// 创建订单请求参数
type WxPayV3CreateOrderReq struct {
	AppID       string                      `json:"appid"`        // 直连商户申请的公众号或移动应用appid
	MchID       string                      `json:"mchid"`        // 直连商户的商户号，由微信支付生成并下发
	Description string                      `json:"description"`  // 商品描述
	OutTradeNo  string                      `json:"out_trade_no"` // 商户订单号
	NotifyURL   string                      `json:"notify_url"`   // 通知地址
	Amount      WxPayV3CreateOrderReqAmount `json:"amount"`       // 订单金额
	Payer       WxPayV3CreateOrderReqPayer  `json:"payer"`        // 支付者
}

type WxPayV3CreateOrderReqAmount struct {
	Total    int64  `json:"total"`    // 总金额
	Currency string `json:"currency"` // 货币类型 CNY：人民币，境内商户号仅支持人民币
}

type WxPayV3CreateOrderReqPayer struct {
	OpenID string `json:"openid"` // 用户在直连商户appid下的唯一标识
}

// 创建订单返回参数
type WxPayV3CreateOrderResp struct {
	PrepayID string `json:"prepay_id"`
}

// 新版 V3 微信支付
// JSAPI/小程序下单API
// 文档地址：https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/pay/transactions/chapter3_2.shtml
func V3OrderCreate(openid, out_trade_no string, amount float64) (resp WxPayV3CreateOrderResp, err error) {
	req := WxPayV3CreateOrderReq{
		AppID:       config.WxpayAppID,
		MchID:       config.WxpayMchID,
		Description: "商品描述",
		OutTradeNo:  out_trade_no,
		NotifyURL:   config.WxpayNotifyURL,
		Amount: WxPayV3CreateOrderReqAmount{
			Total:    ToFen(amount),
			Currency: "CNY",
		},
		Payer: WxPayV3CreateOrderReqPayer{
			OpenID: openid,
		},
	}
	log.Infof("WxPay create order request body:%+v", req)
	// 请求报文主体
	jsonByte, err := json.Marshal(&req)
	if err != nil {
		return resp, errors.New("request body unmarshal failed.")
	}
	authorization := GetAuthorization("POST", "/v3/pay/transactions/jsapi", string(jsonByte))
	headers := map[string]interface{}{
		"Authorization": authorization,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}
	httpResp, err := HttpRequestFunc(fmt.Sprintf("%s/v3/pay/transactions/jsapi", config.WxpayMchURL), "POST", string(jsonByte), headers)
	defer httpResp.Body.Close()
	if err != nil {
		return resp, errors.New("request WxPay create order failed.")
	}
	rspBody, err := ioutil.ReadAll(httpResp.Body)
	err = json.Unmarshal(rspBody, &resp)
	if err != nil {
		return resp, errors.New("response body unmarshal failed.")
	}
	log.Infof("WxPay create order response body:%+v", resp)
	return resp, nil
}
