package alipay

import (
	"crypto"
	"encoding/json"
	"fmt"
	"time"
)

type PayReq struct {
	OutTradeNo         string              `json:"out_trade_no"`        //商户订单号,64个字符以内、只能包含字母、数字、下划线；需保证在商户端不重复
	SellerID           string              `json:"seller_id"`           //卖家支付宝用户ID
	TotalAmount        float64             `json:"total_amount"`        //订单总金额，单位为元，精确到小数点后两位
	DiscountableAmount float64             `json:"discountable_amount"` //可打折金额.参与优惠计算的金额，单位为元，精确到小数点后两位
	Subject            string              `json:"subject"`             //订单标题
	Body               string              `json:"body"`                //对交易或商品的描述
	BuyerID            string              `json:"buyer_id"`            //买家的支付宝唯一用户号（2088开头的16位纯数字）
	GoodsDetail        []PayReqGoodsDetail `json:"goods_detail"`        //订单包含的商品列表信息
}

type PayReqGoodsDetail struct {
	GoodsID        string  `json:"goods_id"`        //商品的编号
	GoodsName      string  `json:"goods_name"`      //商品名称
	Quantity       int64   `json:"quantity"`        //商品数量
	Price          float64 `json:"price"`           //商品单价，单位为元
	GoodsCategory  string  `json:"goods_category"`  //商品类目
	CategoriesTree string  `json:"categories_tree"` //商品类目树，从商品类目根节点到叶子节点的类目id组成，类目id值使用|分割
	Body           string  `json:"body"`            //商品描述信息
	ShowURL        string  `json:"show_url"`        //商品的展示地址
}

type PaySettleReq struct {
	OutRequestNo      string                          `json:"out_request_no"`     //结算请求流水号 开发者自行生成并保证唯一性
	TradeNo           string                          `json:"trade_no"`           //支付宝订单号
	RoyaltyParameters []PaySettleReqRoyaltyParameters `json:"royalty_parameters"` //分账明细信息
	OperatorID        string                          `json:"operator_id"`        //操作员id
}

type PaySettleReqRoyaltyParameters struct {
	RoyaltyType  string `json:"royalty_type"`   //分账类型.普通分账为：transfer;补差为：replenish;为空默认为分账transfer;
	TransOut     string `json:"trans_out"`      //支出方账户
	TransOutType string `json:"trans_out_type"` //支出方账户类型
	TransInType  string `json:"trans_in_type"`  //收入方账户类型
	TransIn      string `json:"trans_in"`       //收入方账户
	Amount       int    `json:"amount"`         //分账的金额，单位为元
	Desc         string `json:"desc"`           //分账描述
}

func CreatePay() {
	payReqGoodsDetail := PayReqGoodsDetail{
		GoodsID:   "apple-01",
		GoodsName: "测试",
		Quantity:  1,
		Price:     11111,
		Body:      "测试apple-01",
	}
	var goodsDetail []PayReqGoodsDetail
	goodsDetail = append(goodsDetail, payReqGoodsDetail)
	req := PayReq{
		OutTradeNo:  "202009010000000002",
		TotalAmount: 112220,
		Subject:     "测试下单",
		Body:        "测试下单",
		BuyerID:     "2088102176090402",
		GoodsDetail: goodsDetail,
	}
	bizContent, _ := json.Marshal(req)
	commonReq := CommonReq{
		AppID:        AppID,
		Method:       "alipay.trade.create",
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

func PaySettle() {
	var royaltyParameters []PaySettleReqRoyaltyParameters
	paySettleReqRoyaltyParameters := PaySettleReqRoyaltyParameters{
		TransOut:     "2088102176090402",
		TransOutType: "userId",
		TransIn:      "bpfcxt6304@sandbox.com",
		TransInType:  "loginName",
	}
	royaltyParameters = append(royaltyParameters, paySettleReqRoyaltyParameters)
	req := PaySettleReq{
		OutRequestNo:      "202009010000000002",
		TradeNo:           "2020090222001490400500708455",
		RoyaltyParameters: royaltyParameters,
	}
	bizContent, _ := json.Marshal(req)
	commonReq := CommonReq{
		AppID:        AppID,
		Method:       "alipay.trade.order.settle",
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
