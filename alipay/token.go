package alipay

import (
	"crypto"
	"encoding/json"
	"io/ioutil"

	"github.com/Biubiubiuuuu/go-pay/config"
	"github.com/lexkong/log"
)

// AuthTokenReq 授权请求参数
type AuthTokenReq struct {
	CommonReq
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	RefreshToken string `json:"refresh_token"`
}

// AuthTokenRsp 授权返回参数
type AuthTokenRsp struct {
	ErrCommonRsp   CommonRsp                  `json:"error_response"` // 授权失败返回参数
	OauthTokenResp AuthTokenRspOauthTokenResp `json:"alipay_system_oauth_token_response"`
	CommonRspSign
}

// AuthTokenRspOauthTokenResp 授权成功返回参数
type AuthTokenRspOauthTokenResp struct {
	CommonRsp
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    uint64 `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	ReExpiresIn  uint64 `json:"re_expires_in"`
}

// AutErrhToken 获取支付宝授权
//  alipay.system.oauth.token(换取授权访问令牌)
func AuthToken(code string) (AuthTokenRsp, error) {
	commom := CommonReq{
		AppID:        config.AlipayAppID,
		Method:       "alipay.system.oauth.token",
		Format:       "JSON",
		Charset:      "utf-8",
		SignType:     "RSA2",
		Timestamp:    GetTimestamp(),
		Version:      "1.0",
		AppAuthToken: "",
		BizContent:   "",
	}
	req := AuthTokenReq{
		GrantType:    "authorization_code",
		Code:         code,
		RefreshToken: "",
		CommonReq:    commom,
	}
	signURLValues := GetAuthSignStr(req)
	signStr := URLValues(signURLValues, false)
	sign := Rsa2Sign(signStr, config.AlipayPrivateKey, crypto.SHA256)
	req.Sign = sign
	postDataURLValues := GetAuthSignStr(req)
	postData := URLValues(postDataURLValues, true)
	url := config.AlipayURL
	log.Infof("request alipay oauth url:%s,req:%+v", url, postData)
	rsp, err := HTTPPost(url, postData, "")
	defer rsp.Body.Close()
	var resp AuthTokenRsp
	if err != nil {
		log.Errorf(err, "request alipay oauth token failed.")
		return resp, err
	}
	body, err := ioutil.ReadAll(rsp.Body)
	log.Infof("alipay oauth rsp body:%+v", string(body))
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Errorf(err, "rsp body unmarshal failed.")
		return resp, err
	}
	log.Infof("alipay oauth rsp body:%+v", resp)
	return resp, nil
}
