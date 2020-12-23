package wxpay

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/Biubiubiuuuu/go-pay/config"
	"github.com/lexkong/log"
)

// 查询订单返回参数
type WxPayV3QueryOrderResp struct {
	TransactionID string `json:"transaction_id"`
	Amount        struct {
		PayerTotal    int    `json:"payer_total"`
		Total         int    `json:"total"`
		Currency      string `json:"currency"`
		PayerCurrency string `json:"payer_currency"`
	} `json:"amount"`
	Mchid      string `json:"mchid"`
	TradeState string `json:"trade_state"`
	Payer      struct {
		Openid string `json:"openid"`
	} `json:"payer"`
	OutTradeNo     string `json:"out_trade_no"`
	Appid          string `json:"appid"`
	TradeStateDesc string `json:"trade_state_desc"`
	TradeType      string `json:"trade_type"`
}

// 新版 V3 微信支付
// 查询订单API
// 文档地址：https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/pay/transactions/chapter3_5.shtml
func V3OrderQuery(transaction_id string) (resp WxPayV3QueryOrderResp, err error) {
	urlStr := fmt.Sprintf("/v3/pay/transactions/id/%s?mchid=%s", transaction_id, config.WxpayMchID)
	authorization := GetAuthorization("GET", urlStr, "")
	headers := map[string]interface{}{
		"Authorization": authorization,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}
	httpResp, err := HttpRequestFunc(fmt.Sprintf("%s%s", config.WxpayMchURL, urlStr), "GET", "", headers)
	defer httpResp.Body.Close()
	if err != nil {
		return resp, errors.New("request weixin order query failed.")
	}
	rspBody, err := ioutil.ReadAll(httpResp.Body)
	err = json.Unmarshal(rspBody, &resp)
	if err != nil {
		return resp, errors.New("response body unmarshal failed.")
	}
	log.Infof("weixin score pay query rsp body:%+v", resp)
	return resp, nil
}
