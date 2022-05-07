package model

import "github.com/shopspring/decimal"

type Order struct {
	Action string
	Symbol string
	OrderId string
	Price decimal.Decimal
}