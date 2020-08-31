package alipay

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// HTTPGet 发送GET请求
// 	url：         请求地址
// 	response：    请求返回的内容
func HTTPGet(url string) string {
	// 超时时间：20秒
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	return result.String()
}

// HTTPPost 发送POST请求
// url：         请求地址
// contentType： 内容类型（默认：application/x-www-form-urlencoded）
// response：    请求返回的内容
func HTTPPost(url string, contentType string, postData string) string {
	if contentType == "" {
		contentType = "application/x-www-form-urlencoded"
	}
	resp, err := http.Post(url, contentType, strings.NewReader(postData))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(body)
}
