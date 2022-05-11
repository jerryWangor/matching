package model

import (
	"encoding/json"
	"matching/utils"
	"matching/utils/enum"
)

// 交易记录
type Trade struct {
	MakerId string `json:"makerid"`
	TakerId string `json:"takerid"`
	TakerSide *enum.OrderSide `json:"takerside"`
	Amount int `json:"amount"`
	Price string `json:"price"`
	Timestamp string `json:"timestamp"`
}

func (o Trade) ToJson() string {
	data, err := json.Marshal(o)
	if err != nil {
		utils.LogError("json marshal error", err, o.MakerId, o.TakerId)
	}
	return string(data)
}