package alipay

import (
	"crypto"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/Biubiubiuuuu/go-pay/config"
	"github.com/lexkong/log"
)

type UnFreezeRequest struct {
	AuthNo       string  `json:"auth_no"`
	OutRequestNo string  `json:"out_request_no"`
	Amount       float64 `json:"amount"`
	Remark       string  `json:"remark"`
	ExtraParam   string  `json:"extra_param"`
}

type UnFreezeResponse struct {
	CommonRspSign
	Data UnFreezeResponseData `json:"alipay_fund_auth_order_unfreeze_response"`
}

type UnFreezeResponseData struct {
	CommonRsp
	AuthNo       string  `json:"auth_no"`
	OutOrderNo   string  `json:"out_order_no"`
	OperationID  string  `json:"operation_id"`
	OutRequestNo string  `json:"out_request_no"`
	Amount       float64 `json:"amount"`
	Status       string  `json:"status"`
	GmtTrans     string  `json:"gmt_trans"`
	CreditAmount float64 `json:"credit_amount"`
	FundAmount   float64 `json:"fund_amount"`
}

// alipay.fund.auth.order.unfreeze(资金授权解冻接口)
// https://opendocs.alipay.com/apis/api_28/alipay.fund.auth.order.unfreeze
func UnFreeze(auth_no, out_request_no string, amount float64) (resp UnFreezeResponse, err error) {
	req := UnFreezeRequest{
		AuthNo:       auth_no,
		OutRequestNo: out_request_no,
		Amount:       amount,
		Remark:       fmt.Sprintf("%s解冻%v元", time.Now().Format("2006/01/02/ 15:04:05"), amount),
		ExtraParam:   "{\"unfreezeBizInfo\":{\"bizComplete\":\"true\"}}", // 履约完成
	}
	bizContent, err := json.Marshal(req)
	if err != nil {
		return resp, errors.New("rsp body unmarshal failed.")
	}
	commonReq := CommonReq{
		AppID:        config.AlipayAppID,
		Method:       "alipay.fund.auth.order.unfreeze",
		Format:       "JSON",
		Charset:      "utf-8",
		SignType:     "RSA2",
		Timestamp:    GetTimestamp(),
		Version:      "1.0",
		AppAuthToken: "",
		BizContent:   string(bizContent),
		NotifyURL:    "",
	}
	signURLValues := GetSignStr(commonReq)
	signStr := URLValues(signURLValues, false)
	sign := Rsa2Sign(signStr, config.AlipayPrivateKey, crypto.SHA256)
	commonReq.Sign = sign
	postDataURLValues := GetSignStr(commonReq)
	postData := URLValues(postDataURLValues, true)
	url := config.AlipayURL
	log.Infof("request alipay freeze order unfreeze url:%s,req:%+v", url, postData)
	rsp, err := HTTPPost(url, postData, "")
	defer rsp.Body.Close()
	if err != nil {
		return resp, errors.New("request alipay unfreeze order failed.")
	}
	rspBody, err := ioutil.ReadAll(rsp.Body)
	err = json.Unmarshal(rspBody, &resp)
	if err != nil {
		return resp, errors.New("rsp body unmarshal failed.")
	}
	log.Infof("alipay order unfreeze resp body:%+v", resp)
	return resp, nil
}
