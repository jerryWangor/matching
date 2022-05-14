package model

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"matching/utils/common"
	"matching/utils/enum"
)

// Trade 交易记录
type Trade struct {
	MakerId string `json:"makerId" comment:"挂单的订单ID"`
	TakerId string `json:"takerId" comment:"吃单的订单ID"`
	TakerSide enum.OrderSide `json:"takerSide" comment:"买还是卖"`
	Amount decimal.Decimal `json:"amount" comment:"交易成功的数量"`
	Price decimal.Decimal `json:"price" comment:"当前交易价格"`
	Timestamp int64 `json:"timestamp" comment:"交易时间"`
}

// FromMap Map转结构体
func (o *Trade) FromMap(tradeMap map[string]interface{}) (*Trade, error) {
	if err := mapstructure.Decode(tradeMap, &o); err != nil {
		return o, common.Errors(fmt.Sprintf("map decode error"))
	}

	return o, nil
}