package alipay

import (
	"crypto"
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/Biubiubiuuuu/go-pay/config"

	"github.com/lexkong/log"
)

type TradePayRequest struct {
	OutTradeNo  string `json:"out_trade_no"`
	ProductCode string `json:"product_code"`
	AuthNo      string `json:"auth_no"`
	Subject     string `json:"subject"`
	TotalAmount string `json:"total_amount"`
	BuyerID     string `json:"buyer_id"`
	SellerID    string `json:"seller_id"`
	Body        string `json:"body"`
	//StoreID         string `json:"store_id"`
	AuthConfirmMode string `json:"auth_confirm_mode"`
}

type TradePayResponse struct {
	AlipayTradePayResponse struct {
		Code            string `json:"code"`
		Msg             string `json:"msg"`
		TradeNo         string `json:"trade_no"`
		OutTradeNo      string `json:"out_trade_no"`
		BuyerLogonID    string `json:"buyer_logon_id"`
		BuyerPayAmount  string `json:"buyer_pay_amount"`
		SettleAmount    string `json:"settle_amount"`
		PayCurrency     string `json:"pay_currency"`
		PayAmount       string `json:"pay_amount"`
		SettleTransRate string `json:"settle_trans_rate"`
		TransPayRate    string `json:"trans_pay_rate"`
		TotalAmount     string `json:"total_amount"`
		TransCurrency   string `json:"trans_currency"`
		SettleCurrency  string `json:"settle_currency"`
		ReceiptAmount   string `json:"receipt_amount"`
		PointAmount     string `json:"point_amount"`
		InvoiceAmount   string `json:"invoice_amount"`
		GmtPayment      string `json:"gmt_payment"`
		FundBillList    []struct {
			FundChannel string `json:"fund_channel"`
			Amount      string `json:"amount"`
			RealAmount  string `json:"real_amount"`
		} `json:"fund_bill_list"`
		CardBalance         string `json:"card_balance"`
		StoreName           string `json:"store_name"`
		BuyerUserID         string `json:"buyer_user_id"`
		DiscountGoodsDetail string `json:"discount_goods_detail"`
		AdvanceAmount       string `json:"advance_amount"`
		AuthTradePayMode    string `json:"auth_trade_pay_mode"`
		ChargeAmount        string `json:"charge_amount"`
		ChargeFlags         string `json:"charge_flags"`
		SettlementID        string `json:"settlement_id"`
		BusinessParams      string `json:"business_params"`
		BuyerUserType       string `json:"buyer_user_type"`
		MdiscountAmount     string `json:"mdiscount_amount"`
		DiscountAmount      string `json:"discount_amount"`
		BuyerUserName       string `json:"buyer_user_name"`
	} `json:"alipay_trade_pay_response"`
	Sign string `json:"sign"`
}

// alipay.trade.pay(统一收单交易支付接口)
// https://opendocs.alipay.com/apis/api_1/alipay.trade.pay
func TradePay(auth_no, out_trade_no, subject, body, buyer_id, store_id string, amount float64) (resp TradePayResponse, err error) {
	req := TradePayRequest{
		OutTradeNo:      out_trade_no,
		ProductCode:     "PRE_AUTH_ONLINE",
		AuthNo:          auth_no,
		Subject:         subject,
		TotalAmount:     strconv.FormatFloat(amount, 'f', 2, 64),
		SellerID:        config.AlipaySellerID,
		BuyerID:         buyer_id,
		Body:            body,
		AuthConfirmMode: "COMPLETE",
	}
	bizContent, err := json.Marshal(req)
	if err != nil {
		log.Errorf(err, "rsp body unmarshal failed.")
		return resp, err
	}
	commonReq := CommonReq{
		AppID:        config.AlipayAppID,
		Method:       "alipay.trade.pay",
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
	log.Infof("request alipay order pay url:%s,req:%+v", url, postData)
	rsp, err := HTTPPost(url, postData, "application/x-www-form-urlencoded; charset=utf-8")
	defer rsp.Body.Close()
	if err != nil {
		log.Errorf(err, "request alipay pay failed.")
		return resp, err
	}
	rspBody, err := ioutil.ReadAll(rsp.Body)
	log.Infof("traderpay response body:%+v", string(rspBody))
	err = json.Unmarshal(rspBody, &resp)
	if err != nil {
		log.Errorf(err, "rsp body unmarshal failed.")
		return resp, err
	}
	return resp, nil
}
