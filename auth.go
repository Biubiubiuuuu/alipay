package alipay

import (
	"crypto"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// GetAuth 应用授权并返回app_auth_code
//  353f0f6d3c0a4eeeaef96684016e6X40
func GetAuth() (string, error) {
	reqURL := fmt.Sprintf("%s?app_id=%s&application_type=TINYAPP,WEBAPP&redirect_uri=%s", AuthURL, AppID, RedirectRUL)
	fmt.Println(reqURL)
	resp := HTTPGet(reqURL)
	u, err := url.Parse(resp)
	if err != nil {
		return "", err
	}
	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", err
	}
	return m["app_auth_code"][0], nil
}

// CommonReq 公共请求参数
// ced30c44a01f4b47b3ee4a873d7dTD04
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

// GetAccessTokenReq 换取授权访问令牌请求参数
type GetAccessTokenReq struct {
	GrantType    string `json:"grant_type"`    // 值为authorization_code时，代表用code换取；值为refresh_token时，代表用refresh_token换取
	Code         string `json:"code"`          // 授权码
	RefreshToken string `json:"refresh_token"` // 刷新令牌
}

// GetAccessToken 获取授权token
func GetAccessToken() {
	req := GetAccessTokenReq{
		GrantType:    "authorization_code",
		Code:         "7bb7cbe5ebcc4c7ba75b390b90e0fX40",
		RefreshToken: "",
	}
	bizContent, _ := json.Marshal(req)
	commonReq := CommonReq{
		AppID:        AppID,
		Method:       "alipay.open.auth.token.app",
		Format:       "JSON",
		Charset:      "UTF-8",
		SignType:     "RSA2",
		Timestamp:    time.Now().Format("2006-01-02 15:04:05"),
		Version:      "1.0",
		AppAuthToken: "",
		BizContent:   string(bizContent),
	}
	signStr, _ := SignStr(commonReq, GetAccessTokenReq{}, false)
	fmt.Println("signStr+++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println(signStr)
	sign := RsaSign(signStr, AppPrivateKey, crypto.SHA256)
	fmt.Println("sign+++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println(sign)
	commonReq.Sign = sign
	fmt.Println("commonReq++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println(commonReq)
	postData, _ := SignStr(commonReq, GetAccessTokenReq{}, true)
	fmt.Println("postData+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println(postData)
	resp := HTTPPost(APIUrl, "", postData)
	fmt.Println(resp)
}
