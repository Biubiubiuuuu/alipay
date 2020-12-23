package alipay

import (
	"crypto"

	"github.com/Biubiubiuuuu/go-pay/config"
)

// 验证授权冻结回调
func FreezeNotify(req FreezeNotifyReq) bool {
	sign := req.Sign
	req.Sign = ""
	req.SignType = ""
	signURLValues := StructToURLVal(req)
	signStr := URLValues(signURLValues, false)
	return Rsa2PubSign(signStr, sign, config.AlipayPublicKey, crypto.SHA256)
}
