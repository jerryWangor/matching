package model

import "github.com/shopspring/decimal"

type KData struct {
	TopPrice decimal.Decimal `json:"topPrice" comment:"最高价格"`
	BottomPrice decimal.Decimal `json:"bottomPrice" comment:"最低价格"`
	NowPrice decimal.Decimal `json:"NowPrice" comment:"当前价格"`
	Timestamp int64 `json:"timestamp" comment:"时间点，往后顺移，比如2022-05-24 12:31:00就是到32分的"`
}

type KDataPrice struct {
	TopPrice decimal.Decimal `json:"topPrice" comment:"最高价格"`
	BottomPrice decimal.Decimal `json:"bottomPrice" comment:"最低价格"`
	NowPrice decimal.Decimal `json:"NowPrice" comment:"当前价格"`
}