package model

type TimeSign struct {
	Time int64 `form:"time" binding:"required" comment:"时间戳"`
	Sign string `form:"sign" binding:"required" comment:"加密串，规则：md5(Time+token)"`
}