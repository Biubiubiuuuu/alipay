package alipay

import (
	"crypto"
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/Biubiubiuuuu/go-pay/config"
	"github.com/lexkong/log"
)

type AlipayFreezeOrderCannelRequest struct {
	OutOrderNo string `json:"out_order_no"`
	Remark     string `json:"remark"`
}

type AlipayFreezeOrderCancelResponse struct {
	CommonRspSign
	Data AlipayFreezeOrderCancelResponseData `json:"alipay_fund_auth_operation_cancel_response"`
}

type AlipayFreezeOrderCancelResponseData struct {
	CommonRsp
	AuthNo       string `json:"auth_no"`
	OutOrderNo   string `json:"out_order_no"`
	OperationID  string `json:"operation_id"`
	OutRequestNo string `json:"out_request_no"`
	Action       string `json:"action"`
}

// alipay.fund.auth.operation.cancel(资金授权撤销接口)
// https://opendocs.alipay.com/apis/api_28/alipay.fund.auth.operation.cancel
func FreezeCancel(out_order_no string) (resp AlipayFreezeOrderCancelResponse, err error) {
	req := AlipayFreezeOrderCannelRequest{
		OutOrderNo: out_order_no,
		Remark:     "取消资金授权冻结",
	}
	bizContent, err := json.Marshal(req)
	if err != nil {
		return resp, errors.New("rsp body unmarshal failed.")
	}
	commonReq := CommonReq{
		AppID:        config.AlipayAppID,
		Method:       "alipay.fund.auth.operation.cancel",
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
	log.Infof("request alipay freeze order cancel url:%s,req:%+v", url, postData)
	rsp, err := HTTPPost(url, postData, "")
	defer rsp.Body.Close()
	if err != nil {
		return resp, errors.New("request alipay freeze order cancel failed.")
	}
	rspBody, err := ioutil.ReadAll(rsp.Body)
	err = json.Unmarshal(rspBody, &resp)
	if err != nil {
		return resp, errors.New("rsp body unmarshal failed.")
	}
	return resp, nil
}
