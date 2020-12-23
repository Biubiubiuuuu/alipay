package wxpay

import (
	"encoding/xml"
	"io/ioutil"

	"github.com/Biubiubiuuuu/go-pay/config"
	"github.com/lexkong/log"
	uuid "github.com/satori/go.uuid"
)

//微信企业付款到零钱
type WxPayWithdrawReq struct {
	MchAppId       string `xml:"mch_appid"`
	MchId          string `xml:"mchid"`
	DeviceInfo     string `xml:"device_info"`
	NonceStr       string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	SignType       string `xml:"sign_type"`
	PartnerTradeNo string `xml:"partner_trade_no"`
	Openid         string `xml:"openid"`
	CheckName      string `xml:"check_name"`
	ReUserName     string `xml:"re_user_name"`
	Amount         int    `xml:"amount"`
	Desc           string `xml:"desc"`
	SpbillCreateIp string `xml:"spbill_create_ip"`
}

type WxPayWithdrawRsp struct {
	ReturnCode     string `xml:"return_code"`
	ReturnMsg      string `xml:"return_msg"`
	MchAppid       string `xml:"mch_appid"`
	MchId          string `xml:"mch_id"`
	DeviceInfo     string `xml:"device_info"`
	NonceStr       string `xml:"nonce_str"`
	ResultCode     string `xml:"result_code"`
	ErrCode        string `xml:"err_code"`
	ErrCodeDes     string `xml:"err_code_des"`
	PartnerTradeNo string `xml:"partner_trade_no"`
	PaymentNo      string `xml:"payment_no"`
	PaymentTime    string `xml:"payment_time"`
}

//微信提现
func WxPayWithdraw(amount float64, openid string, seqid string) error {
	host := config.WxpayMchURL
	url := host + "/mmpaymkttransfers/promotion/transfers"
	var req WxPayWithdrawReq
	req.Openid = openid
	//封装请求包
	req.MchAppId = config.WxpayAppID
	req.MchId = config.WxpayMchID
	u4 := uuid.NewV4()
	req.NonceStr = u4.String()[1:32]
	req.PartnerTradeNo = seqid
	req.Openid = openid
	req.CheckName = "NO_CHECK"
	//元转分
	req.Amount = int(amount * 100)
	req.Desc = "提现"
	req.Sign = SignMd5(req)
	xml_req, _ := xml.MarshalIndent(req, "", "")
	log.Infof("create withdraw url:%s,req:%+v", url, req)
	rsp, err := SecurePost(url, xml_req)
	defer rsp.Body.Close()
	if err != nil {
		log.Errorf(err, "request wxpay withdraw failed.")
		return err
	}
	body, err := ioutil.ReadAll(rsp.Body)
	var rsp_body WxPayWithdrawRsp
	err = xml.Unmarshal(body, &rsp_body)
	if err != nil {
		log.Errorf(err, "rsp body unmarshal failed.")
		return err
	}
	log.Infof("wxpay withdraw rsp body:%+v", rsp_body)
	return nil
}
