package alipay

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

const (
	// PEMBEGIN 开头
	PEMBEGIN = "-----BEGIN RSA PRIVATE KEY-----\n"
	// PEMEND 结尾
	PEMEND = "\n-----END RSA PRIVATE KEY-----"
)

// RsaSign 签名
func RsaSign(signContent string, privateKey string, hash crypto.Hash) string {
	shaNew := hash.New()
	shaNew.Write([]byte(signContent))
	hashed := shaNew.Sum(nil)
	priKey, err := ParsePrivateKey(privateKey)
	if err != nil {
		panic(err)
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, priKey, hash, hashed)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(signature)
}

// ParsePrivateKey 私钥验证
func ParsePrivateKey(privateKey string) (*rsa.PrivateKey, error) {
	privateKey = FormatPrivateKey(privateKey)
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("私钥信息错误！")
	}
	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return priKey, nil
}

// FormatPrivateKey 组装私钥
func FormatPrivateKey(privateKey string) string {
	if !strings.HasPrefix(privateKey, PEMBEGIN) {
		privateKey = PEMBEGIN + privateKey
	}
	if !strings.HasSuffix(privateKey, PEMEND) {
		privateKey = privateKey + PEMEND
	}
	return privateKey
}

// GetCertRootSn  golang解析支付宝根证书
func GetCertRootSn(certPath string) (string, error) {
	certData, err := ioutil.ReadFile(certPath)
	if err != nil {
		return "", err
	}
	strs := strings.Split(string(certData), "-----END CERTIFICATE-----")

	var cert bytes.Buffer
	for i := 0; i < len(strs); i++ {
		if strs[i] == "" {
			continue
		}
		if blo, _ := pem.Decode([]byte(strs[i] + "-----END CERTIFICATE-----")); blo != nil {
			c, err := x509.ParseCertificate(blo.Bytes)
			if err != nil {
				continue
			}
			if _, ok := alog[c.SignatureAlgorithm.String()]; !ok {
				continue
			}
			si := c.Issuer.String() + c.SerialNumber.String()
			if cert.String() == "" {
				cert.WriteString(md5V(si))
			} else {
				cert.WriteString("_")
				cert.WriteString(md5V(si))
			}
		}

	}
	return cert.String(), nil
}

var alog map[string]string = map[string]string{
	"MD2-RSA":       "MD2WithRSA",
	"MD5-RSA":       "MD5WithRSA",
	"SHA1-RSA":      "SHA1WithRSA",
	"SHA256-RSA":    "SHA256WithRSA",
	"SHA384-RSA":    "SHA384WithRSA",
	"SHA512-RSA":    "SHA512WithRSA",
	"SHA256-RSAPSS": "SHA256WithRSAPSS",
	"SHA384-RSAPSS": "SHA384WithRSAPSS",
	"SHA512-RSAPSS": "SHA512WithRSAPSS",
}

// md5V md5V
func md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// GetAlipayPublicKey 解析应用公钥证书
func GetAlipayPublicKey(certPath string) (publicKey *rsa.PublicKey) {
	certContent, _ := ioutil.ReadFile(certPath)
	certDecode, _ := pem.Decode(certContent)
	x509Cert, err := x509.ParseCertificate(certDecode.Bytes)
	if err != nil {
		return nil
	}
	//解析支付宝公钥证书内容
	if pub, ok := x509Cert.PublicKey.(*rsa.PublicKey); ok {
		publicKeybyte := x509.MarshalPKCS1PublicKey(pub)
		publicKeyString := base64.StdEncoding.EncodeToString(publicKeybyte)
		fmt.Println(publicKeyString)
	}

	//通过rsa.VerifyPKCS1v15验签需要此参数
	return x509Cert.PublicKey.(*rsa.PublicKey)
}
