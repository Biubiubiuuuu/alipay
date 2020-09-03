package pay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	// PEMBEGIN 开头
	PEMBEGIN = "-----BEGIN RSA PRIVATE KEY-----\n"
	// PEMEND 结尾
	PEMEND = "\n-----END RSA PRIVATE KEY-----"
)

/*************************************************************************
	> File Name: trade_create.go
	> Description: trade_create.go
	> Author: heohe
	> Mail: yyhejun@163.com
	> Created Time: 2019-12-14 11:17:34
	> Last modified: 2019-12-14 11:17:34
 ************************************************************************/

// CommonReq 公共请求参数结构体
type CommonReq struct {
	AppID        string `json:"app_id"`         // 开发者应用ID
	Method       string `json:"method"`         // 接口名称
	Format       string `json:"format"`         // 请求结构数据类型 仅支持JSON
	Charset      string `json:"charset"`        // 请求使用的编码格式，如utf-8,gbk,gb2312等 默认utf-8
	SignType     string `json:"sign_type"`      // 签名算法类型 RSA2和RSA 默认RSA2
	Sign         string `json:"sign"`           // 签名串
	Timestamp    string `json:"timestamp"`      // 请求时间 格式 yyyy-MM-dd HH:mm:ss
	Version      string `json:"version"`        // 接口版本 1.0
	NotifyURL    string `json:"notify_url"`     // 支付宝服务器主动通知商户服务器里指定的页面http/https路径
	AppAuthToken string `json:"app_auth_token"` // 应用授权app_auth_token
	BizContent   string `json:"biz_content"`    // 请求参数的集合
}

// CommonRsp 公共回包结构体
type CommonRsp struct {
	Code    string `json:"code,omitempy"`
	Msg     string `json:"msg,omitempy"`
	SubCode string `json:"sub_code,omitempy"`
	SubMsg  string `json:"sub_msg,omitempy"`
}

// CommonRspSign 公共回包结构体sign
type CommonRspSign struct {
	Sign string `json:"sign"`
}

// GoodsDetail 公共商品参数结构体
type GoodsDetail struct {
	GoodsID        string  `json:"goods_id,omitempy"`        // 商品的编号
	GoodsName      string  `json:"goods_name,omitempy"`      // 商品名称
	Quantity       int64   `json:"quantity,omitempy"`        // 商品数量
	Price          float64 `json:"price,omitempy"`           // 商品单价，单位为元
	GoodsCategory  string  `json:"goods_category,omitempy"`  // 商品类目
	CategoriesTree string  `json:"categories_tree,omitempy"` // 商品类目树，从商品类目根节点到叶子节点的类目id组成，类目id值使用|分割
	Body           string  `json:"body,omitempy"`            // 商品描述信息
	ShowURL        string  `json:"show_url,omitempy"`        // 商品的展示地址
}

// ExtendParams 公共扩展参数结构体
type ExtendParams struct {
	SysServiceProviderID string `json:"sys_service_provider_id,omitempy"`
	IndustryRefluxInfo   string `json:"industry_reflux_info,omitempy"`
	CardType             string `json:"card_type,omitempy"`
}

// GetAuthSignStr 支付宝授权接口alipay.system.oauth.token待签名数据
func GetAuthSignStr(req AuthTokenReq) url.Values {
	var p = url.Values{}
	p.Add("app_id", req.AppID)
	p.Add("method", req.Method)
	p.Add("format", req.Format)
	p.Add("charset", req.Charset)
	p.Add("sign_type", req.SignType)
	p.Add("timestamp", req.Timestamp)
	p.Add("version", req.Version)
	p.Add("notify_url", req.NotifyURL)
	p.Add("biz_content", req.BizContent)
	p.Add("app_auth_token", req.AppAuthToken)
	p.Add("sign", req.Sign)
	p.Add("grant_type", req.GrantType)
	p.Add("code", req.Code)
	p.Add("refresh_token", req.RefreshToken)
	return p
}

// GetSignStr 支付宝接口待签名数据
func GetSignStr(req CommonReq) url.Values {
	var p = url.Values{}
	p.Add("app_id", req.AppID)
	p.Add("method", req.Method)
	p.Add("format", req.Format)
	p.Add("charset", req.Charset)
	p.Add("sign_type", req.SignType)
	p.Add("timestamp", req.Timestamp)
	p.Add("version", req.Version)
	p.Add("notify_url", req.NotifyURL)
	p.Add("biz_content", req.BizContent)
	p.Add("app_auth_token", req.AppAuthToken)
	p.Add("sign", req.Sign)
	return p
}

// StructToURLVal 统一下单接口结构体转URL Val
func StructToURLVal(param interface{}) url.Values {
	var p = url.Values{}
	typ := reflect.TypeOf(param)
	val := reflect.ValueOf(param)
	kd := val.Kind()
	if kd != reflect.Struct {
		fmt.Println("expect struct")
		return nil
	}
	num := val.NumField()
	//遍历结构体的所有字段
	for i := 0; i < num; i++ {
		tagVal := typ.Field(i).Tag.Get("json")
		if tagVal != "" {
			p.Add(tagVal, val.Field(i).String())
		}
	}
	return p
}

// URLValues 待签名字符串 备注：去除值为空的，按键值ASCII码递增排序
//  param 代签名数据
//  urlencode 是否对数据值URL编码
func URLValues(param url.Values, urlencode bool) string {
	if param == nil {
		param = make(url.Values, 0)
	}
	var pList = make([]string, 0, 0)
	for key := range param {
		var value = strings.TrimSpace(param.Get(key))
		if len(value) > 0 {
			if urlencode {
				value = url.QueryEscape(value)
			}
			pList = append(pList, fmt.Sprintf("%s=%s", key, value))
		}
	}
	sort.Strings(pList)
	return strings.Join(pList, "&")
}

// GetTimestamp 获取时间  格式yyyy-MM-dd HH:mm:ss
func GetTimestamp() string {
	return time.Now().In(time.Local).Format("2006-01-02 15:04:05")
}

// Rsa2Sign RSA2私钥签名
func Rsa2Sign(signContent string, privateKey string, hash crypto.Hash) string {
	shaNew := hash.New()
	shaNew.Write([]byte(signContent))
	hashed := shaNew.Sum(nil)
	priKey, err := ParsePrivateKey(privateKey)
	if err != nil {
		return ""
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, priKey, hash, hashed)
	if err != nil {
		return ""
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

// Rsa2PubSign RSA2公钥验证签名
func Rsa2PubSign(signContent, sign, publicKey string, hash crypto.Hash) bool {
	shaNew := hash.New()
	shaNew.Write([]byte(signContent))
	hashed := shaNew.Sum(nil)
	pubKey, err := ParsePublicKey(publicKey)
	if err != nil {
		return false
	}
	sig, _ := base64.StdEncoding.DecodeString(string(sign))
	err = rsa.VerifyPKCS1v15(pubKey, hash, hashed, sig)
	if err != nil {
		return false
	}
	return true
}

// ParsePublicKey 公钥验证
func ParsePublicKey(publicKey string) (*rsa.PublicKey, error) {
	publicKey = FormatPublicKey(publicKey)
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, errors.New("公钥信息错误！")
	}
	pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}

// FormatPublicKey 组装公钥
func FormatPublicKey(publicKey string) string {
	if !strings.HasPrefix(publicKey, PEMBEGIN) {
		publicKey = PEMBEGIN + publicKey
	}
	if !strings.HasSuffix(publicKey, PEMEND) {
		publicKey = publicKey + PEMEND
	}
	return publicKey
}

// HTTPPost 发送POST请求
func HTTPPost(url, postData, contentType string) (*http.Response, error) {
	if contentType == "" {
		contentType = "application/x-www-form-urlencoded"
	}
	client := &http.Client{}
	return client.Post(url, contentType, strings.NewReader(postData))
}

// CreateAlipayOrderID 生成支付宝商户订单ID,支付宝要求64位以内
func CreateAlipayOrderID() string {
	orderID := "zfb"
	now := time.Now()
	orderID += now.Format("060102150405")
	orderID += strconv.FormatInt(now.UnixNano()/1e3, 10)
	return orderID
}
