package alipay

import (
	"fmt"
	"net/url"
)

// GetAuth 应用授权并返回app_auth_code
//  9d5e7b16a1cb4b788bad56c853508X30
func GetAuth() (string, error) {
	reqUrl := fmt.Sprintf("%s?app_id=%s&application_type=TINYAPP,WEBAPP&redirect_uri=%s", AuthUrl, AppID, RedirectUrl)
	fmt.Println(reqUrl)
	resp := HttpGet(reqUrl)
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

// GetToken 获取授权token
func GetToken() {

}
