package alipay

import (
	"net/url"
	"sort"
	"strings"
	"time"
)

// SignStr 待签名字符串
func SignStr(req CommonReq, auth GetAccessTokenReq, urlencode bool) (string, error) {
	var p = url.Values{}
	p.Add("app_id", req.AppID)
	p.Add("method", req.Method)
	p.Add("format", req.Format)
	p.Add("charset", req.Charset)
	p.Add("sign_type", req.SignType)
	p.Add("timestamp", GetTimestamp())
	p.Add("version", req.Version)
	p.Add("notify_url", req.NotifyURL)
	p.Add("biz_content", req.BizContent)
	p.Add("app_auth_token", req.AppAuthToken)
	p.Add("sign", req.Sign)
	p.Add("grant_type", auth.GrantType)
	p.Add("code", auth.Code)
	p.Add("refresh_token", auth.RefreshToken)
	str := URLValues(p, urlencode)
	return str, nil
}

// URLValues 待签名字符串 去除值为空的
//  键值 ASCII 码递增排序
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
			pList = append(pList, key+"="+value)
		}
	}
	sort.Strings(pList)
	return strings.Join(pList, "&")
}

// GetTimestamp 获取时间  格式yyyy-MM-dd HH:mm:ss
func GetTimestamp() string {
	return time.Now().In(time.Local).Format("2006-01-02 15:04:05")
}
