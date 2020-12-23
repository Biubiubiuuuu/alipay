package alipay

import (
	"crypto"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Biubiubiuuuu/go-pay/config"
)

type FreezeCreateRequst struct {
	OutOrderNo   string  `json:"out_order_no"`
	OutRequestNo string  `json:"out_request_no"`
	OrderTitle   string  `json:"order_title"`
	Amount       float64 `json:"amount"`
	ProductCode  string  `json:"product_code"`
	PayeeUserID  string  `json:"payee_user_id"`
	ExtraParam   string  `json:"extra_param"`
	//EnablePayChannels string  `json:"enable_pay_channels"`
}

// alipay.fund.auth.order.app.freeze(线上资金授权冻结接口)
// https://opendocs.alipay.com/apis/api_28/alipay.fund.auth.order.app.freeze
func FreezeCreate(sn string, amount float64) (out_order_no, postData string, err error) {
	number := CreateAlipayOrderID()
	req := FreezeCreateRequst{
		OutOrderNo:   number,
		OutRequestNo: number,
		OrderTitle:   "预授权冻结",
		Amount:       amount,
		ProductCode:  "PRE_AUTH_ONLINE",
		ExtraParam:   fmt.Sprintf("{\"category\":\"RENT_SHARABLE_CHARGERS\",\"outStoreCode\":\"%s\",\"outStoreAlias\":\"租借设备\"}", sn),
		PayeeUserID:  "",
		//EnablePayChannels: "[{\"payChannelType\":\"PCREDIT_PAY\"},{\"payChannelType\":\"MONEY_FUND\"},{\"payChannelType\":\"CREDITZHIMA\"}]",
	}
	bizContent, err := json.Marshal(req)
	if err != nil {
		return number, postData, errors.New("rsp body unmarshal failed.")
	}
	commonReq := CommonReq{
		AppID:        config.AlipayAppID,
		Method:       "alipay.fund.auth.order.app.freeze",
		Format:       "JSON",
		Charset:      "utf-8",
		SignType:     "RSA2",
		Timestamp:    GetTimestamp(),
		Version:      "1.0",
		AppAuthToken: "",
		BizContent:   string(bizContent),
		NotifyURL:    config.AlipayNotifyFreezeURL,
	}
	signURLValues := GetSignStr(commonReq)
	signStr := URLValues(signURLValues, false)
	sign := Rsa2Sign(signStr, config.AlipayPrivateKey, crypto.SHA256)
	commonReq.Sign = sign
	postDataURLValues := GetSignStr(commonReq)
	postData = URLValues(postDataURLValues, true)
	return number, postData, nil
}
