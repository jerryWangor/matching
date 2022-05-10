package model

type TimeSign struct {
	Time int64 `form:"time" binding:"required"`
	Sign string `form:"sign" binding:"required"`
}