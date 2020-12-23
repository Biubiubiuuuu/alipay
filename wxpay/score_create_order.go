package wxpay

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/Biubiubiuuuu/go-pay/config"
)

// 微信支付分订单创建请求参数
type WxPayScoreCreateOrderRequest struct {
	OutOrderNo          string              `json:"out_order_no"`         // 商户系统内部服务订单号
	AppID               string              `json:"appid"`                // 公众账号ID
	ServiceID           string              `json:"service_id"`           // 服务ID
	ServiceIntroduction string              `json:"service_introduction"` // 服务信息
	PostPayments        []PostPaymentsData  `json:"post_payments"`        // 后付费项目
	PostDiscounts       []PostDiscountsData `json:"post_discounts"`       // 后付费商户优惠
	TimeRange           TimeRangeData       `json:"time_range"`           // 服务时间段
	Location            LocationData        `json:"location"`             // 服务位置
	RiskFund            RiskFundData        `json:"risk_fund"`            // 订单风险金
	Attach              string              `json:"attach"`               // 商户数据包
	NotifyURL           string              `json:"notify_url"`           // 商户回调地址
	OpenID              string              `json:"openid"`               // 用户标识
	NeedUserConfirm     bool                `json:"need_user_confirm"`    // 是否需要用户确认
}

// 后付费项目
type PostPaymentsData struct {
	Name        string `json:"name"`        // 付费项目名称
	Amount      int64  `json:"amount"`      // 金额
	Description string `json:"description"` // 计费说明
	Count       int64  `json:"count"`       // 付费数量
}

// 后付费商户优惠
type PostDiscountsData struct {
	Name        string `json:"name"`        // 付费项目名称
	Amount      int64  `json:"amount"`      // 金额
	Description string `json:"description"` // 计费说明
}

// 服务时间段
type TimeRangeData struct {
	StartTime       string `json:"start_time"`        // 服务开始时间
	StartTimeRemark string `json:"start_time_remark"` // 服务开始时间备注
	//EndTime         string `json:"end_time"`          // 预计服务结束时间
	//EndTimeRemark   string `json:"end_time_remark"`   // 预计服务结束时间备注
}

// 服务位置
type LocationData struct {
	StartLocation string `json:"start_location"` // 服务开始地点
	EndLocation   string `json:"end_location"`   // 预计服务结束位置
}

type RiskFundData struct {
	Name        string `json:"name"`        // 风险金名称
	Amount      int64  `json:"amount"`      // 风险金额
	Description string `json:"description"` // 风险说明
}

// 微信信用分创建订单返回参数
type WxPayScoreCreateOrderResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  []struct {
		Field string `json:"field"`
		Value string `json:"value"`
		Issue string `json:"issue"`
	} `json:"detail,omitempy"`
	Appid               string `json:"appid"`
	Mchid               string `json:"mchid"`
	OutOrderNo          string `json:"out_order_no"`
	ServiceID           string `json:"service_id"`
	ServiceIntroduction string `json:"service_introduction"`
	State               string `json:"state"`
	StateDescription    string `json:"state_description"`
	PostPayments        []struct {
		Name        string `json:"name"`
		Amount      int    `json:"amount"`
		Description string `json:"description"`
		Count       int    `json:"count"`
	} `json:"post_payments"`
	PostDiscounts []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"post_discounts"`
	RiskFund struct {
		Name        string `json:"name"`
		Amount      int    `json:"amount"`
		Description string `json:"description"`
	} `json:"risk_fund"`
	TimeRange struct {
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	} `json:"time_range"`
	Location struct {
		StartLocation string `json:"start_location"`
		EndLocation   string `json:"end_location"`
	} `json:"location"`
	Attach    string `json:"attach"`
	NotifyURL string `json:"notify_url"`
	OrderID   string `json:"order_id"`
	Package   string `json:"package"`
}

// 创建信用分订单
//  openid 用户标识
//  payment_description 计费说明
//  location 租借地点
//  risk_fund_amount 风险金额（押金）
func CreateScoreOrder(openid, payment_description, location string, risk_fund_amount float64) (resp WxPayScoreCreateOrderResponse, err error) {
	// 组装请求参数
	req := WxPayScoreCreateOrderRequest{
		OutOrderNo:          CreateWxOrderId(),
		AppID:               config.WxpayAppID,
		ServiceID:           config.WxpayServiceID,
		ServiceIntroduction: "共享充电器",
		PostPayments: []PostPaymentsData{
			{
				Name:        "租借使用费用",
				Amount:      0,
				Description: payment_description,
				Count:       1,
			},
		},
		PostDiscounts: []PostDiscountsData{
			{
				Name:        "租借使用费用",
				Amount:      0,
				Description: payment_description,
			},
		},
		TimeRange: TimeRangeData{
			StartTime:       "OnAccept",
			StartTimeRemark: "租借开始时间",
		},
		Location: LocationData{
			StartLocation: location,
			EndLocation:   location,
		},
		RiskFund: RiskFundData{
			Name:        "DEPOSIT",
			Amount:      ToFen(risk_fund_amount),
			Description: "租借押金",
		},
		Attach:          "",
		NotifyURL:       config.WxpayScoreNotifyURL,
		OpenID:          openid,
		NeedUserConfirm: true,
	}
	// 请求报文主体
	jsonByte, err := json.Marshal(&req)
	if err != nil {
		return resp, errors.New("request body unmarshal failed.")
	}
	authorization := GetAuthorization("POST", "/v3/payscore/serviceorder", string(jsonByte))
	headers := map[string]interface{}{
		"Authorization": authorization,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}
	httpResp, err := HttpRequestFunc(fmt.Sprintf("%s/v3/payscore/serviceorder", config.WxpayMchURL), "POST", string(jsonByte), headers)
	defer httpResp.Body.Close()
	if err != nil {
		return resp, errors.New("request weixin zhifufen order create failed.")
	}
	rspBody, err := ioutil.ReadAll(httpResp.Body)
	err = json.Unmarshal(rspBody, &resp)
	if err != nil {
		return resp, errors.New("response body unmarshal failed.")
	}
	return resp, nil
}
