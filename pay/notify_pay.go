package pay

import (
	"crypto"

	"github.com/Biubiubiuuuu/alipay/config"
)

// NotifyPayReq 支付宝统一下单回调请求参数
type NotifyPayReq struct {
	NotifyTime   string `json:"notify_time"`    // 通知的发送时间。格式为yyyy-MM-dd HH:mm:ss
	NotifyType   string `json:"notify_type"`    // 通知的类型
	NotifyID     string `json:"notify_id"`      // 通知校验ID
	AppID        string `json:"app_id"`         // 支付宝分配给开发者的应用Id
	Charset      string `json:"charset"`        // 编码格式，如utf-8、gbk、gb2312等
	Version      string `json:"version"`        // 调用的接口版本，固定为：1.0
	SignType     string `json:"sign_type"`      // 商户生成签名字符串所使用的签名算法类型，目前支持RSA2和RSA，推荐使用RSA2
	Sign         string `json:"sign"`           // 异步返回结果的验签
	TradeNo      string `json:"trade_no"`       // 支付宝交易凭证号
	OutTradeNo   string `json:"out_trade_no"`   // 原支付请求的商户订单号
	OutBizNo     string `json:"out_biz_no"`     // 商户业务ID，主要是退款通知中返回退款申请的流水号
	BuyerID      string `json:"buyer_id"`       // 买家支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字
	BuyerLogonID string `json:"buyer_logon_id"` // 买家支付宝账号
	SellerID     string `json:"seller_id"`      // 卖家支付宝用户号
	SellerEmail  string `json:"seller_email"`   // 卖家支付宝账号
	// WAIT_BUYER_PAY 交易创建，等待买家付款
	// TRADE_CLOSED	  未付款交易超时关闭，或支付完成后全额退款
	// TRADE_SUCCESS  交易支付成功
	// TRADE_FINISHED 交易结束，不可退款
	TradeStatus   string `json:"trade_status"`   // 交易目前所处的状态，见“交易状态说明”
	TotalAmount   string `json:"total_amount"`   // 本次交易支付的订单金额，单位为人民币（元）
	ReceiptAmount string `json:"receipt_amount"` // 商家在交易中实际收到的款项，单位为元
	RefundFee     string `json:"refund_fee"`     // 退款总金额 退款回调才有的参数
}

// NotifyPay 支付宝统一下单回调
func NotifyPay(req NotifyPayReq) bool {
	sign := req.Sign
	req.Sign = ""
	signURLValues := StructToURLVal(req)
	signStr := URLValues(signURLValues, false)
	return Rsa2PubSign(signStr, sign, config.AppPublicKey, crypto.SHA256)
}
