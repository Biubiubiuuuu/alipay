package alipay

import (
	"crypto"
	"encoding/json"
	"io/ioutil"

	"github.com/Biubiubiuuuu/go-pay/config"
	"github.com/lexkong/log"
)

// 统一订单查询请求参数
type TradeQueryRequest struct {
	OutTradeNo string `json:"out_trade_no"` // 商户订单号
	TradeNo    string `json:"trade_no"`     // 支付宝交易订单号
}

// 统一订单查询返回参数
type TradeQueryResponse struct {
	CommonRspSign
	Data TradeQueryResponseData `json:"alipay_trade_query_response"`
}

// 统一订单查询响应返回参数
type TradeQueryResponseData struct {
	CommonRsp
	OutTradeNo  string `json:"out_trade_no"`
	TradeNo     string `json:"trade_no"`
	TradeStatus string `json:"trade_status"` // 交易状态：WAIT_BUYER_PAY（交易创建，等待买家付款）、TRADE_CLOSED（未付款交易超时关闭，或支付完成后全额退款）、TRADE_SUCCESS（交易支付成功）、TRADE_FINISHED（交易结束，不可退款）
}

// 订单交易结果查询
//  alipay.trade.query(统一收单线下交易查询)
func TradeQuery(trade_no string) (TradeQueryResponse, error) {
	req := TradeQueryRequest{
		TradeNo: trade_no,
	}
	var resp TradeQueryResponse
	bizContent, err := json.Marshal(req)
	if err != nil {
		log.Errorf(err, "rsp body unmarshal failed.")
		return resp, err
	}
	commonReq := CommonReq{
		AppID:        config.AlipayAppID,
		Method:       "alipay.trade.query",
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
	log.Infof("request alipay order query url:%s,req:%+v", url, postData)
	rsp, err := HTTPPost(url, postData, "")
	defer rsp.Body.Close()
	if err != nil {
		log.Errorf(err, "request alipay order query failed.")
		return resp, err
	}
	rspBody, err := ioutil.ReadAll(rsp.Body)
	err = json.Unmarshal(rspBody, &resp)
	if err != nil {
		log.Errorf(err, "rsp body unmarshal failed.")
		return resp, err
	}
	log.Infof("alipay order query resp body:%+v", resp)
	return resp, nil
}
