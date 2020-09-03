package pay

import (
	"crypto"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Biubiubiuuuu/alipay/config"
)

// RefundReq 支付宝退款请求参数
type RefundReq struct {
	OutTradeNo     string                 `json:"out_trade_no,omitempy"`    // 订单支付时传入的商户订单号,不能和 trade_no同时为空
	TradeNo        string                 `json:"trade_no,omitempy"`        // 支付宝交易号，和商户订单号不能同时为空
	RefundAmount   float64                `json:"refund_amount,omitempy"`   // 需要退款的金额，该金额不能大于订单金额,单位为元，支持两位小数
	RefundCurrency string                 `json:"refund_currency,omitempy"` // 订单退款币种信息
	RefundReason   string                 `json:"refund_reason,omitempy"`   // 退款的原因说明
	OutRequestNo   string                 `json:"out_request_no,omitempy"`  // 标识一次退款请求，同一笔交易多次退款需要保证唯一，如需部分退款，则此参数必传
	OperatorID     string                 `json:"operator_id,omitempy"`     // 商户的操作员编号
	StoreID        string                 `json:"store_id,omitempy"`        // 商户的门店编号
	TerminalID     string                 `json:"terminal_id,omitempy"`     // 商户的终端编号
	DoodsDetail    []RefundReqDoodsDetail `json:"goods_detail,omitempy"`    // 退款包含的商品列表信息
	OrgPid         string                 `json:"org_pid,omitempy"`         // 银行间联模式下有用，其它场景请不要使用；双联通过该参数指定需要退款的交易所属收单机构的pid;
	QueryOptions   []string               `json:"query_options,omitempy"`   // 查询选项，商户通过上送该参数来定制同步需要额外返回的信息字段，数组格式
}

// RefundReqDoodsDetail 退款包含的商品列表信息
type RefundReqDoodsDetail struct {
	GoodsID        string  `json:"goods_id,omitempy"`        // 商品的编号
	AlipayGoodsID  string  `json:"alipay_goods_id,omitempy"` // 支付宝定义的统一商品编号
	GoodsName      string  `json:"goods_name,omitempy"`      // 商品名称
	Quantity       int     `json:"quantity,omitempy"`        // 商品数量
	Price          float64 `json:"price,omitempy"`           // 商品单价，单位为元
	GoodsCategory  string  `json:"goods_category,omitempy"`  // 商品类目
	CategoriesTree string  `json:"categories_tree,omitempy"` // 商品类目树
	Body           string  `json:"body,omitempy"`            // 商品描述信息
	ShowURL        string  `json:"show_url,omitempy"`        // 商品的展示地址
}

// RefundResp 统一收单交易退款公共响应返回参数
type RefundResp struct {
	TradeRefundResp RefundRespTradeRefundResp `json:"alipay_trade_refund_response"`
	CommonRspSign
}

// RefundRespTradeRefundResp 统一收单交易退款响应返回参数
type RefundRespTradeRefundResp struct {
	CommonRsp
	TradeNo                      string                    `json:"trade_no"`
	OutTradeNo                   string                    `json:"out_trade_no"`
	BuyerLogonID                 string                    `json:"buyer_logon_id"`
	FundChange                   string                    `json:"fund_change"`
	RefundFee                    string                    `json:"refund_fee"`
	RefundCurrency               string                    `json:"refund_currency"`
	GmtRefundPay                 string                    `json:"gmt_refund_pay"`
	RefundDetailItemList         []RefundDetailItemList    `json:"refund_detail_item_list"`
	StoreName                    string                    `json:"store_name"`
	BuyerUserID                  string                    `json:"buyer_user_id"`
	RefundSettlementID           string                    `json:"refund_settlement_id"`
	PresentRefundBuyerAmount     string                    `json:"present_refund_buyer_amount"`
	PresentRefundDiscountAmount  string                    `json:"present_refund_discount_amount"`
	PresentRefundMdiscountAmount string                    `json:"present_refund_mdiscount_amount"`
	RefundPresetPaytoolList      []RefundPresetPaytoolList `json:"refund_preset_paytool_list"`
}

// RefundDetailItemList 退款使用的资金渠道
type RefundDetailItemList struct {
	FundChannel string  `json:"fund_channel"`
	BankCode    string  `json:"bank_code"`
	Amount      int     `json:"amount"`
	RealAmount  float64 `json:"real_amount"`
	FundType    string  `json:"fund_type"`
}

// RefundPresetPaytoolList 退回的前置资产列表
type RefundPresetPaytoolList struct {
	Amount         []float64 `json:"amount"`
	AssertTypeCode string    `json:"assert_type_code"`
}

// RefundPay 支付宝退款
//  alipay.trade.refund(统一收单交易退款接口)
func RefundPay(outTradeNo, tradeNo, refundReason string, refundAmount float64) (RefundResp, error) {
	var doodsDetail []RefundReqDoodsDetail
	refundReq := RefundReq{
		OutTradeNo:     outTradeNo,
		TradeNo:        tradeNo,
		RefundAmount:   refundAmount,
		RefundCurrency: "CNY",
		RefundReason:   refundReason,
		DoodsDetail:    doodsDetail,
	}
	var resp RefundResp
	bizContent, err := json.Marshal(refundReq)
	if err != nil {
		return resp, err
	}
	commonReq := CommonReq{
		AppID:        config.AppID,
		Method:       "alipay.trade.refund",
		Format:       "JSON",
		Charset:      "utf-8",
		SignType:     "RSA2",
		Timestamp:    GetTimestamp(),
		Version:      "1.0",
		AppAuthToken: "",
		BizContent:   string(bizContent),
		NotifyURL:    "",
	}
	signURLValues := GetSignStr(commonReq)
	signStr := URLValues(signURLValues, false)
	sign := Rsa2Sign(signStr, config.AppPrivateKey, crypto.SHA256)
	commonReq.Sign = sign
	postDataURLValues := GetSignStr(commonReq)
	postData := URLValues(postDataURLValues, true)
	url := config.APIUrl
	fmt.Println(url)
	fmt.Println(postData)
	rsp, err := HTTPPost(url, postData, "")
	defer rsp.Body.Close()
	if err != nil {
		return resp, err
	}
	rspBody, err := ioutil.ReadAll(rsp.Body)
	fmt.Println(string(rspBody))
	err = json.Unmarshal(rspBody, &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
