package model

import (
	"github.com/shopspring/decimal"
	"matching/utils/enum"
)

type TimeSign struct {
	Time int64 `form:"time" json:"time" binding:"required" comment:"时间戳"`
	Sign string `form:"sign" json:"sign" binding:"required" comment:"加密串，规则：md5(Time+token)"`
}

type PriceTopN struct {
	Id int `form:"id" json:"id" comment:"排序"`
	Side enum.OrderSide `form:"side" json:"side" comment:"买/卖"`
	Price decimal.Decimal `form:"price" json:"price" comment:"价格"`
	Amount decimal.Decimal `form:"amount" json:"amount" comment:"数量"`
}