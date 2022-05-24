package model

type TimeSign struct {
	Time int64 `form:"time" json:"time" binding:"required" comment:"时间戳"`
	Sign string `form:"sign" json:"sign" binding:"required" comment:"加密串，规则：md5(Time+token)"`
}

type PriceTopN struct {
	Price float64 `form:"price" json:"price" comment:"价格"`
	Amount float64 `form:"amount" json:"amount" comment:"数量"`
}