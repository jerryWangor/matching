package model

type TimeSign struct {
	Time int64 `form:"time" json:"time" binding:"required" comment:"时间戳"`
	Sign string `form:"sign" json:"sign" binding:"required" comment:"加密串，规则：md5(Time+token)"`
}