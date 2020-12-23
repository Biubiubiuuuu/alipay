package wxpay

import (
	"encoding/json"
	"time"
)

type WxPayV3NotifyOrderResp struct {
	TransactionID string `json:"transaction_id"`
	Amount        struct {
		PayerTotal    int    `json:"payer_total"`
		Total         int    `json:"total"`
		Currency      string `json:"currency"`
		PayerCurrency string `json:"payer_currency"`
	} `json:"amount"`
	Mchid           string `json:"mchid"`
	TradeState      string `json:"trade_state"`
	BankType        string `json:"bank_type"`
	PromotionDetail []struct {
		Amount              int    `json:"amount"`
		WechatpayContribute int    `json:"wechatpay_contribute"`
		CouponID            string `json:"coupon_id"`
		Scope               string `json:"scope"`
		MerchantContribute  int    `json:"merchant_contribute"`
		Name                string `json:"name"`
		OtherContribute     int    `json:"other_contribute"`
		Currency            string `json:"currency"`
		Type                string `json:"type"`
		StockID             string `json:"stock_id"`
		GoodsDetail         []struct {
			GoodsRemark    string `json:"goods_remark"`
			Quantity       int    `json:"quantity"`
			DiscountAmount int    `json:"discount_amount"`
			GoodsID        string `json:"goods_id"`
			UnitPrice      int    `json:"unit_price"`
		} `json:"goods_detail"`
	} `json:"promotion_detail"`
	SuccessTime time.Time `json:"success_time"`
	Payer       struct {
		Openid string `json:"openid"`
	} `json:"payer"`
	OutTradeNo     string `json:"out_trade_no"`
	Appid          string `json:"appid"`
	TradeStateDesc string `json:"trade_state_desc"`
	TradeType      string `json:"trade_type"`
	Attach         string `json:"attach"`
	SceneInfo      struct {
		DeviceID string `json:"device_id"`
	} `json:"scene_info"`
}

// 微信支付订单回调验证签名
// 文档地址：https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/pay/transactions/chapter3_11.shtml
func V3OrderNotify(req WxPayV3NotifyReq) (resp WxPayV3NotifyOrderResp, err error) {
	signDecrypt, err := AES256GCMDecrypt(req.Resource.Ciphertext, req.Resource.Nonce, req.Resource.AssociatedData)
	if err != nil {
		return resp, err
	}
	err = json.Unmarshal(signDecrypt, &resp)
	return
}
