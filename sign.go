package alipay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"hash"

	"github.com/pkg/errors"
)

// RASSign RAS签名
//  SignType：RSA
//  SignType：RSA2
func RASSign(data []byte, SignType, pemPriKey string) (signature string, err error) {
	var h hash.Hash
	var hType crypto.Hash
	switch SignType {
	case "RSA":
		h = sha1.New()
		hType = crypto.SHA1
	case "RSA2":
		h = sha256.New()
		hType = crypto.SHA256
	}
	h.Write(data)
	d := h.Sum(nil)
	pk, err := ParsePrivateKey(pemPriKey)
	if err != nil {
		err = errors.Wrap(err, "私钥错误")
		return
	}
	bs, err := rsa.SignPKCS1v15(rand.Reader, pk, hType, d)

	if err != nil {
		err = errors.Wrap(err, "私钥错误")
		return
	}
	signature = base64.StdEncoding.EncodeToString(bs)
	return
}

// ParsePrivateKey 私钥
func ParsePrivateKey(privateKey string) (pk *rsa.PrivateKey, err error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		err = errors.Errorf("私钥格式错误:%s", privateKey)
		return
	}
	switch block.Type {
	case "RSA PRIVATE KEY":
		rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err == nil {
			pk = rsaPrivateKey
		} else {
			err = errors.Wrap(err, privateKey)
		}
	default:
		err = errors.Errorf("私钥格式错误:%s", privateKey)
	}
	return
}
