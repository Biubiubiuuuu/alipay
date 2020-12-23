package wxpay

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/Biubiubiuuuu/go-pay/config"
	"github.com/lexkong/log"
)

// 微信支付分完结订单请求参数
type WxPayScoreCompleteOrderRequest struct {
	AppID         string                `json:"appid"`          // 公众账号ID
	ServiceID     string                `json:"service_id"`     // 服务ID
	PostPayments  []PostPaymentsData    `json:"post_payments"`  // 后付费项目
	PostDiscounts []PostDiscountsData   `json:"post_discounts"` // 后付费商户优惠
	TimeRange     CompleteTimeRangeData `json:"time_range"`     // 服务时间段
	Location      LocationData          `json:"location"`       // 服务位置
	TotalAmount   int64                 `json:"total_amount"`   // 总金额
	ProfitSharing bool                  `json:"profit_sharing"` // 微信支付服务分账标记
	Goodstag      string                `json:"goods_tag"`      // 订单优惠标记
}

// 服务时间段
type CompleteTimeRangeData struct {
	EndTime       string `json:"end_time"`        // 预计服务结束时间
	EndTimeRemark string `json:"end_time_remark"` // 预计服务结束时间备注
}

// 微信支付分订单完结返回参数
type WxPayScoreCompleteOrderResponse struct {
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
	TotalAmount         int    `json:"total_amount"`
	PostPayments        []struct {
		Name        string `json:"name"`
		Amount      int    `json:"amount"`
		Description string `json:"description"`
		Count       int    `json:"count"`
	} `json:"post_payments"`
	PostDiscounts []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Amount      int    `json:"amount"`
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
	OrderID        string `json:"order_id"`
	NeedCollection bool   `json:"need_collection"`
}

// 微信支付分完结订单
func CompleteScoreOrder(out_order_no, payment_description, location, end_time string, total_amount float64) (resp WxPayScoreCompleteOrderResponse, err error) {
	req := WxPayScoreCompleteOrderRequest{
		AppID:     config.WxpayAppID,
		ServiceID: config.WxpayServiceID,
		PostPayments: []PostPaymentsData{{
			Name:        "租借使用费用",
			Amount:      ToFen(total_amount),
			Description: payment_description,
			Count:       1,
		}},
		PostDiscounts: []PostDiscountsData{
			{
				Name:        "租借使用费用",
				Amount:      0,
				Description: payment_description,
			},
		},
		TimeRange: CompleteTimeRangeData{
			EndTime:       end_time,
			EndTimeRemark: "租借结束时间",
		},
		Location: LocationData{
			StartLocation: location,
			EndLocation:   location,
		},
		TotalAmount: ToFen(total_amount),
	}
	// 请求报文主体
	jsonByte, err := json.Marshal(&req)
	if err != nil {
		return resp, errors.New("request body unmarshal failed.")
	}
	urlStr := fmt.Sprintf("/v3/payscore/serviceorder/%s/complete", out_order_no)
	authorization := GetAuthorization("POST", urlStr, string(jsonByte))
	headers := map[string]interface{}{
		"Authorization": authorization,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}
	httpResp, err := HttpRequestFunc(fmt.Sprintf("%s%s", config.WxpayMchURL, urlStr), "POST", string(jsonByte), headers)
	defer httpResp.Body.Close()
	if err != nil {
		return resp, errors.New("request weixin zhifufen order complete failed.")
	}
	rspBody, err := ioutil.ReadAll(httpResp.Body)
	err = json.Unmarshal(rspBody, &resp)
	if err != nil {
		return resp, errors.New("response body unmarshal failed.")
	}
	log.Infof("weixin score pay complete rsp body:%+v", resp)
	return resp, nil
}
