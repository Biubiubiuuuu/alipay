package alipay

import (
	"crypto"
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/Biubiubiuuuu/go-pay/config"
	"github.com/lexkong/log"
)

type AlipayFreezeOrderQueryRequest struct {
	AuthNo       string `json:"auth_no"`
	OutOrderNo   string `json:"out_order_no"`
	OperationID  string `json:"operation_id"`
	OutRequestNo string `json:"out_request_no"`
}

type AlipayFreezeOrderQueryResponse struct {
	CommonRspSign
	Data AlipayFreezeOrderQueryResponseData `json:"alipay_fund_auth_operation_detail_query_response"`
}

type AlipayFreezeOrderQueryResponseData struct {
	CommonRsp
	Amount                  string `json:"amount"`
	AuthNo                  string `json:"auth_no"`
	CreditAmount            string `json:"credit_amount"`
	ExtraParam              string `json:"extra_param"`
	FundAmount              string `json:"fund_amount"`
	GmtCreate               string `json:"gmt_create"`
	GmtTrans                string `json:"gmt_trans"`
	OperationID             string `json:"operation_id"`
	OperationType           string `json:"operation_type"`
	OrderTitle              string `json:"order_title"`
	OutOrderNo              string `json:"out_order_no"`
	OutRequestNo            string `json:"out_request_no"`
	PayerLogonID            string `json:"payer_logon_id"`
	PayerUserID             string `json:"payer_user_id"`
	PreAuthType             string `json:"pre_auth_type"`
	Remark                  string `json:"remark"`
	RestAmount              string `json:"rest_amount"`
	RestCreditAmount        string `json:"rest_credit_amount"`
	RestFundAmount          string `json:"rest_fund_amount"`
	Status                  string `json:"status"`
	TotalFreezeAmount       string `json:"total_freeze_amount"`
	TotalFreezeCreditAmount string `json:"total_freeze_credit_amount"`
	TotalFreezeFundAmount   string `json:"total_freeze_fund_amount"`
	TotalPayAmount          string `json:"total_pay_amount"`
	TotalPayCreditAmount    string `json:"total_pay_credit_amount"`
	TotalPayFundAmount      string `json:"total_pay_fund_amount"`
}

// alipay.fund.auth.operation.detail.query(资金授权操作查询接口)
// https://opendocs.alipay.com/apis/api_28/alipay.fund.auth.operation.detail.query
func FreezeOrderQuery(out_order_no string) (resp AlipayFreezeOrderQueryResponse, err error) {
	req := AlipayFreezeOrderQueryRequest{
		OutOrderNo:   out_order_no,
		OutRequestNo: out_order_no,
	}
	bizContent, err := json.Marshal(req)
	if err != nil {
		return resp, errors.New("rsp body unmarshal failed.")
	}
	commonReq := CommonReq{
		AppID:        config.AlipayAppID,
		Method:       "alipay.fund.auth.operation.detail.query",
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
	log.Infof("request alipay freeze order query url:%s,req:%+v", url, postData)
	rsp, err := HTTPPost(url, postData, "")
	defer rsp.Body.Close()
	if err != nil {
		return resp, errors.New("request alipay freeze order query failed.")
	}
	rspBody, err := ioutil.ReadAll(rsp.Body)
	err = json.Unmarshal(rspBody, &resp)
	if err != nil {
		return resp, errors.New("rsp body unmarshal failed.")
	}
	return resp, nil
}
