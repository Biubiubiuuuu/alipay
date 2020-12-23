package alipay

import (
	"crypto"

	"github.com/Biubiubiuuuu/go-pay/config"
)

// NotifyPay 支付宝统一下单回调
func NotifyPay(req AlipayNotifyPayOrRefund) bool {
	sign := req.Sign
	req.Sign = ""
	req.SignType = ""
	signURLValues := StructToURLVal(req)
	signStr := URLValues(signURLValues, false)
	return Rsa2PubSign(signStr, sign, config.AlipayPublicKey, crypto.SHA256)
}
