package wxpay

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/Biubiubiuuuu/go-pay/config"
	"github.com/lexkong/log"
)

// 微信支付分查询订单返回参数
type WxpayScoreGetOrderResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  []struct {
		Field string `json:"field"`
		Value string `json:"value"`
		Issue string `json:"issue"`
	} `json:"detail,omitempy"`
	Appid               string `json:"appid"`
	Mchid               string `json:"mchid"`
	ServiceID           string `json:"service_id"`
	OutOrderNo          string `json:"out_order_no"`
	ServiceIntroduction string `json:"service_introduction"`
	State               string `json:"state"`
	StateDescription    string `json:"state_description"`
	TotalAmount         int    `json:"total_amount"`
	PostPayments        []struct {
		Name        string `json:"name"`
		Amount      int    `json:"amount"`
		Description string `json:"description"`
		Count       int    `json:"count"`
	} `json:"post_payments"`
	PostDiscounts []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Amount      int    `json:"amount"`
	} `json:"post_discounts"`
	RiskFund struct {
		Name        string `json:"name"`
		Amount      int    `json:"amount"`
		Description string `json:"description"`
	} `json:"risk_fund"`
	TimeRange struct {
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	} `json:"time_range"`
	Location struct {
		StartLocation string `json:"start_location"`
		EndLocation   string `json:"end_location"`
	} `json:"location"`
	Attach         string `json:"attach"`
	NotifyURL      string `json:"notify_url"`
	OrderID        string `json:"order_id"`
	NeedCollection bool   `json:"need_collection"`
	Collection     struct {
		State        string `json:"state"`
		TotalAmount  int    `json:"total_amount"`
		PayingAmount int    `json:"paying_amount"`
		PaidAmount   int    `json:"paid_amount"`
		Details      []struct {
			Seq           int    `json:"seq"`
			Amount        int    `json:"amount"`
			PaidType      string `json:"paid_type"`
			PaidTime      string `json:"paid_time"`
			TransactionID string `json:"transaction_id"`
		} `json:"details"`
	} `json:"collection"`
}

// 微信支付分订单查询
func WxScorePayGetOrder(out_order_no string) (resp WxpayScoreGetOrderResponse, err error) {
	urlStr := fmt.Sprintf("/v3/payscore/serviceorder?out_order_no=%s&service_id=%s&appid=%s", out_order_no, config.WxpayServiceID, config.WxpayAppID)
	authorization := GetAuthorization("GET", urlStr, "")
	headers := map[string]interface{}{
		"Authorization": authorization,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}
	httpResp, err := HttpRequestFunc(fmt.Sprintf("%s%s", config.WxpayMchURL, urlStr), "GET", "", headers)
	defer httpResp.Body.Close()
	if err != nil {
		return resp, errors.New("request weixin zhifufen order query failed.")
	}
	rspBody, err := ioutil.ReadAll(httpResp.Body)
	err = json.Unmarshal(rspBody, &resp)
	if err != nil {
		return resp, errors.New("response body unmarshal failed.")
	}
	log.Infof("weixin score pay query rsp body:%+v", resp)
	return resp, nil
}
