package alipay

import (
	"crypto"

	"github.com/Biubiubiuuuu/go-pay/config"
)

// alipay.trade.pay(统一收单交易支付接口) 回调验签
func TradePayNotyfy(req TradePayNotyfyReq) bool {
	sign := req.Sign
	req.Sign = ""
	req.SignType = ""
	signURLValues := StructToURLVal(req)
	signStr := URLValues(signURLValues, false)
	return Rsa2PubSign(signStr, sign, config.AlipayPublicKey, crypto.SHA256)
}
