package model

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"matching/utils/common"
	"matching/utils/enum"
)

type Order struct {
	//Accid int `json:"accid" comment:"账号ID"`
	Symbol string `json:"symbol" comment:"交易标"`
	OrderId string `json:"orderId" comment:"订单ID"`
	Action enum.OrderAction `json:"action" comment:"挂单还是撤单"`
	Type enum.OrderType `json:"type" comment:"竞价类型"`
	Side enum.OrderSide `json:"side" comment:"买/卖"`
	Amount decimal.Decimal `json:"amount" comment:"数量"`
	Price decimal.Decimal `json:"price" comment:"价格"`
	Timestamp int64 `json:"timestamp" comment:"时间"`
}

// FromMap Map转结构体
func (o *Order) FromMap(orderMap map[string]interface{}) (*Order, error) {
	if err := mapstructure.Decode(orderMap, &o); err != nil {
		return o, common.Errors(fmt.Sprintf("map decode error"))
	}

	return o, nil
}