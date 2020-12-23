package wxpay

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Biubiubiuuuu/go-pay/alipay"
	"github.com/Biubiubiuuuu/go-pay/config"
	"github.com/lexkong/log"
	uuid "github.com/satori/go.uuid"
)

// 微信支付回调通知应答参数
type WxpayV3NotifyResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// V3微信支付回调通知参数
type WxPayV3NotifyReq struct {
	ID           string `json:"id"`
	CreateTime   string `json:"create_time"`
	ResourceType string `json:"resource_type"`
	EventType    string `json:"event_type"`
	Resource     struct {
		Algorithm      string `json:"algorithm"`
		Ciphertext     string `json:"ciphertext"`
		Nonce          string `json:"nonce"`
		AssociatedData string `json:"associated_data"`
	} `json:"resource"`
	Summary string `json:"summary"`
}

//利用反射通用生成数字签名函数。
func Sign(obj interface{}) string {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	sign := url.Values{}
	for i := 0; i < t.NumField(); i++ {
		//如果为空值不参与签名
		xmlname := t.Field(i).Tag.Get("xml")
		value := v.Field(i).Interface()
		var tmp_value string
		//只能添加string类型其他类型一律转为string
		switch value.(type) {
		case int:
			tmp_value = strconv.Itoa(value.(int))
		case string:
			tmp_value, _ = value.(string)
		default:
			log.Errorf(nil, "sign unknown type.")
		}
		if tmp_value == "" {
			continue
		}
		sign.Add(xmlname, tmp_value)
	}
	r, _ := url.QueryUnescape(sign.Encode())
	r += "&key="
	r += config.WxpayMchKey
	log.Infof("sign:%s", r)
	return GetHmacSha256Encoding(r, config.WxpayMchKey)
}

//利用反射通用生成数字签名函数。
func SignMd5(obj interface{}) string {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	sign := url.Values{}
	for i := 0; i < t.NumField(); i++ {
		//如果为空值不参与签名
		xmlname := t.Field(i).Tag.Get("xml")
		value := v.Field(i).Interface()
		var tmp_value string
		//只能添加string类型其他类型一律转为string
		switch value.(type) {
		case int:
			tmp_value = strconv.Itoa(value.(int))
		case string:
			tmp_value, _ = value.(string)
		default:
			log.Errorf(nil, "sign unknown type.")
		}
		if tmp_value == "" {
			continue
		}
		sign.Add(xmlname, tmp_value)
	}
	r, _ := url.QueryUnescape(sign.Encode())
	r += "&key="
	r += config.WxpayMchKey
	log.Infof("sign:%s", r)
	return GetMd5Encoding(r)
}

//利用hmac-sha256对字符串进行加密
func GetHmacSha256Encoding(msg, key string) string {
	ret := hmac.New(sha256.New, []byte(key))
	ret.Write([]byte(msg))
	return strings.ToUpper(hex.EncodeToString(ret.Sum(nil)))
}

//利用hmac-sha1进行加密,并进行base64编码
func GetHmacSha1Encoding(msg, key string) string {
	ret := hmac.New(sha1.New, []byte(key))
	ret.Write([]byte(msg))
	return base64.StdEncoding.EncodeToString(ret.Sum(nil))
}

//利用md5加密
func GetMd5Encoding(msg string) string {
	m := md5.New()
	m.Write([]byte(msg))
	return strings.ToUpper(hex.EncodeToString(m.Sum(nil)))
}

//发起携带微信证书的安全请求
func SecurePost(url string, req []byte) (*http.Response, error) {
	wechat_pay_cert := config.WechatPayCert
	wechat_pay_key := config.WechatPayKey
	wechat_root_ca := config.WechatRootCa
	//读取证书对
	certs, err := tls.LoadX509KeyPair(wechat_pay_cert, wechat_pay_key)
	if err != nil {
		log.Errorf(err, "load wechat cert failed.")
		return nil, err
	}
	root_ca, err := ioutil.ReadFile(wechat_root_ca)
	if err != nil {
		log.Errorf(err, "load wechat rootca failed.")
		return nil, err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(root_ca)
	tls_config := tls.Config{
		Certificates: []tls.Certificate{certs},
		RootCAs:      pool,
	}
	transport := &http.Transport{TLSClientConfig: &tls_config}
	client := &http.Client{Transport: transport}
	return client.Post(url, "application/xml", strings.NewReader(string(req)))
}

//发起微信post请求
func SimplePost(url string, req []byte) (*http.Response, error) {
	client := &http.Client{}
	return client.Post(url, "application/xml", strings.NewReader(string(req)))
}

//生成微信商户订单ID,微信要求32位以内
func CreateWxOrderId() string {
	order_id := "wx"
	now := time.Now()
	order_id += now.Format("060102150405")
	order_id += strconv.FormatInt(now.UnixNano()/1e3, 10)
	return order_id
}

//生成微信商户退款订单ID,微信要求32位以内
func CreateWxRefundId() string {
	order_id := "wxrf"
	now := time.Now()
	order_id += now.Format("060102150405")
	order_id += strconv.FormatInt(now.UnixNano()/1e3, 10)
	return order_id
}

//生成微信商户提现订单ID,微信要求32位以内
func CreateWxWithdrawId() string {
	order_id := "wxwd"
	now := time.Now()
	order_id += now.Format("060102150405")
	order_id += strconv.FormatInt(now.UnixNano()/1e3, 10)
	return order_id
}

//判断二进制串是否是utf编码
func IsUTF8(buf []byte) bool {
	nBytes := 0
	for i := 0; i < len(buf); i++ {
		if nBytes == 0 {
			if (buf[i] & 0x80) != 0 { //与操作之后不为0，说明首位为1
				for (buf[i] & 0x80) != 0 {
					buf[i] <<= 1 //左移一位
					nBytes++     //记录字符共占几个字节
				}
				if nBytes < 2 || nBytes > 6 { //因为UTF8编码单字符最多不超过6个字节
					return false
				}
				nBytes-- //减掉首字节的一个计数
			}
		} else { //处理多字节字符
			if buf[i]&0xc0 != 0x80 { //判断多字节后面的字节是否是10开头
				return false
			}
			nBytes--
		}
	}
	return nBytes == 0
}

//本系统统一使用元作为金额的计数单位，微信中使用的是分，因为这里有个转换
func ToYuan(fen int) float64 {
	return float64(fen) / 100.00
}

func ToFen(yuan float64) int64 {
	return int64(yuan * 100)
}

// go AEAD_AES_256_GCM 解密
func AESDecrypt(keyStr, body, nonceStr, additionalDataStr string) []byte {
	ciphertext, _ := hex.DecodeString(body)
	key := []byte(keyStr)
	nonce, _ := hex.DecodeString(nonceStr)
	var additionalData []byte
	if additionalDataStr != "" {
		additionalDataByte, _ := hex.DecodeString(additionalDataStr)
		additionalData = additionalDataByte
	} else {
		additionalData = nil
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Errorf(err, "AES decrypt failed.")
		return nil
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Errorf(err, "AES decrypt failed.")
		return nil
	}
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, additionalData)
	if err != nil {
		log.Errorf(err, "AES decrypt failed.")
		return nil
	}
	return plaintext
}

// http请求
func HttpRequestFunc(url, method, postData string, headers map[string]interface{}) (*http.Response, error) {
	log.Infof("request wxpay score order url:%s,method:%s,header:%s,postData:%+v", url, method, headers, postData)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(postData))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v.(string))
	}
	return client.Do(req)
}

// 	获取authorization
func GetAuthorization(method, url, requestbody string) string {
	// 组装待签名字符串
	var signBur strings.Builder
	// HTTP请求方法\n
	signBur.WriteString(fmt.Sprintf("%s\n", method))
	// URL\n
	signBur.WriteString(fmt.Sprintf("%s\n", url))
	// 请求时间戳\n
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signBur.WriteString(fmt.Sprintf("%v\n", timestamp))
	// 请求随机串\n
	u4 := uuid.NewV4()
	signBur.WriteString(fmt.Sprintf("%s\n", u4.String()[1:32]))
	// 请求报文主体\n
	signBur.WriteString(fmt.Sprintf("%s\n", requestbody))
	sign := alipay.Rsa2Sign(signBur.String(), config.WxpayPrivateKey, crypto.SHA256)
	return fmt.Sprintf("WECHATPAY2-SHA256-RSA2048 mchid=\"%s\",nonce_str=\"%s\",signature=\"%s\",timestamp=\"%s\",serial_no=\"%s\"", config.WxpayMchID, u4.String()[1:32], sign, timestamp, config.WxpaySerialNo)
}

// AES-256-GCM go 解密
func AES256GCMDecrypt(ciphertext, nonce2, associatedData2 string) (plaindata []byte, err error) {
	key := []byte(config.WXpayAppSecretV3)
	additionalData := []byte(associatedData2)
	nonce := []byte(nonce2)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCMWithNonceSize(block, len(nonce))
	if err != nil {
		return nil, err
	}
	cipherdata, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}
	plaindata, err = aesgcm.Open(nil, nonce, cipherdata, additionalData)
	return plaindata, err
}

// 	V3 小程序调起支付签名计算
func GetV3Sign(appid, timeStamp, nonceStr, packageStr string) string {
	// 组装待签名字符串
	var signBur strings.Builder
	// 公众号id\n
	signBur.WriteString(fmt.Sprintf("%s\n", appid))
	// 时间戳\n
	signBur.WriteString(fmt.Sprintf("%v\n", timeStamp))
	// 随机字符串\n
	signBur.WriteString(fmt.Sprintf("%s\n", nonceStr))
	// 订单详情扩展字符串\n
	signBur.WriteString(fmt.Sprintf("%s\n", packageStr))
	return alipay.Rsa2Sign(signBur.String(), config.WxpayPrivateKey, crypto.SHA256)
}
