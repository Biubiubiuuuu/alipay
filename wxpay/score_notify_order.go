package wxpay

import (
	"encoding/json"
)

// 微信支付分订单回调解密资源对象参数
type ScoreNotifyOrderResponse struct {
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
	} `json:"post_discounts"`
	RiskFund struct {
		Name        string `json:"name"`
		Amount      int    `json:"amount"`
		Description string `json:"description"`
	} `json:"risk_fund"`
	TimeRange struct {
		StartTime string `json:"start_time"`
	} `json:"time_range"`
	Location struct {
		StartLocation string `json:"start_location"`
		EndLocation   string `json:"end_location"`
	} `json:"location"`
	Attach         string `json:"attach"`
	NotifyURL      string `json:"notify_url"`
	OrderID        string `json:"order_id"`
	NeedCollection bool   `json:"need_collection"`
	Openid         string `json:"openid"`
}

// 微信支付分订单回调验证签名
func ScoreNotifyOrderCheckSign(req WxPayV3NotifyReq) (resp ScoreNotifyOrderResponse, err error) {
	signDecrypt, err := AES256GCMDecrypt(req.Resource.Ciphertext, req.Resource.Nonce, req.Resource.AssociatedData)
	if err != nil {
		return resp, err
	}
	err = json.Unmarshal(signDecrypt, &resp)
	return
}
