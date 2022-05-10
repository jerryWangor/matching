package model

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"matching/utils"
	"matching/utils/enum"
)

type Order struct {
	Accid int `json:"accid"`
	Action enum.OrderAction `json:"action"`
	Symbol string `json:"symbol"`
	OrderId string `json:"orderid"`
	Amount int `json:"amount"`
	Price decimal.Decimal `json:"price"`
	Timestamp float64 `json:"timestamp"`
}

func (o Order) toJson() string {
	data, err := json.Marshal(o)
	if err != nil {
		utils.LogError("json marshal error", err, o.OrderId)
	}
	return string(data)
}

// 把缓存中的订单转换成Order结构
func (o Order) FromMap(omap Order) {

}