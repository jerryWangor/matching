package model

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"matching/utils"
)

type Order struct {
	Accid int `json:"accid"`
	Action string `json:"action"`
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