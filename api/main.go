package main

import "github.com/Biubiubiuuuu/alipay"

/*
{
	"alipay_open_auth_token_app_response": {
		"code": "10000",
		"msg": "Success",
		"app_auth_token": "202009BB7be34e2075e149bfb1ca8ae3b8416X40",
		"app_refresh_token": "202009BB3432ef94e8cd4757975d584c2fbacX40",
		"auth_app_id": "2021000120622491",
		"expires_in": 31536000,
		"re_expires_in": 32140800,
		"tokens": [{
			"app_auth_token": "202009BB7be34e2075e149bfb1ca8ae3b8416X40",
			"app_refresh_token": "202009BB3432ef94e8cd4757975d584c2fbacX40",
			"auth_app_id": "2021000120622491",
			"expires_in": 31536000,
			"re_expires_in": 32140800,
			"user_id": "2088102176090402"
		}],
		"user_id": "2088102176090402"
	},
	"sign": "SiXSmho2bwj/PIJp2g1wAo+J144Cb9ipJOG+4q8ZOtfXaq12eLnE53L8lpjlK0/MOouonETAtCWbb7kLc4D1q9652J8tT5jAUXP3DAG44sb1WW4+bki6TX620HQxEZRlntxPcmfCXcz8BujyM6Yfb54Bqu0ELg+RnchEQwkaH1d2jelKuBeNvq62jESGHj5kNhBCkSjx+68RyUWGWA1HQ6z8zmsq9Ixb3T62wgFgg/HPAwNJc0XLCaA7x1RvvuhFMmc9Oj1nA11/sVS4+rrDNzUMZHnRJrAzX+XXRHjHI2GLHt6jpIZPyccg3AtotsyZtS5dHc51dpSC7gem6hqEuQ=="
}
{
	"alipay_trade_create_response": {
		"code": "10000",
		"msg": "Success",
		"out_trade_no": "202009010000000002",
		"trade_no": "2020090222001490400500708455"
	},
	"sign": "X/Bg1+mngH1vtN2Dm1SKgtXyLnhJupSsObA21rqPRwOFsBKUTxbEJlwBVFmgHn/PxYJNH0gIo/LkOp4+w5R5wcoW9NgkfoLfoCZt8G9JiWNLbHYkUntpWI97YKfKQszXJSC/7DqeInBqFT4iHIB/eHGwi+fy0rxatTnIteNlr3nYE+HHr6xV85lSqtkV2YF+3p21gkLispYsHkSRyHbrnXbIKvNwDBjJvAhPiOqINPcxFzWZrzS9+dhVuT18jzCaF7k7cEa3iYBRMRQT3svI3Fo05q7O9Sq3bItUiM28bNMj8CYLGulz9pmo8d+HrfV5KnkkAOLb3YfssN/zSGPJBA=="
}
*/
func main() {
	//str, _ := alipay.GetAuth()
	//fmt.Println(str)

	//alipay.GetAccessToken()
	// 统一下单
	//alipay.CreatePay()
	// 结算
	alipay.PaySettle()
}
