package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"matching/engine"
	"matching/model"
	"matching/utils/cache"
	"matching/utils/code"
	"matching/utils/enum"
	"matching/utils/redis"
	"net/http"
)

type handleOrder struct {
	Accid int `json:"accid" form:"accid" binding:"require"`
	Action enum.OrderAction `json:"action" form:"action" binding:"require"`
	Symbol string `json:"symbol" form:"symbol" binding:"require"`
	OrderId string `json:"orderid" form:"orderid" binding:"require"`
	Side enum.OrderSide `json:"side" form:"side" binding:"require"`
	Type enum.OrderType `json:"type" form:"type" binding:"require"`
	Amount decimal.Decimal `json:"amount" form:"amount" binding:"require"`
	Price decimal.Decimal `json:"price" form:"price" binding:"require"`
}

func HandleOrder(c *gin.Context) {
	// 绑定参数
	var hOrder handleOrder
	if c.ShouldBind(&hOrder) != nil {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_PARAMS_NOTEXISTS, "msg": "参数缺失"})
		return
	}

	// 判断参数
	if hOrder.Symbol == "" {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_PARAMS_ERROR, "msg": "交易标参数不能为空"})
		return
	}

	// 判断该交易标引擎是否开启，从redis缓存中查询
	if redis.HasSymbol(hOrder.Symbol) {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_SYMBOL_MATCHINIG_OPEN_REPEAT, "msg": "交易标引擎重复开启"})
		return
	}

	// 写入redis并发给通道
	order := model.Order{
		Accid: hOrder.Accid,
		Action: hOrder.Action,
		Symbol: hOrder.Symbol,
		OrderId: hOrder.OrderId,
		Side: hOrder.Side,
		Type: hOrder.Type,
		Amount: hOrder.Amount,
		Price: hOrder.Price,
	}
	cache.SaveOrder(order)
	engine.OrderChanMap[order.Symbol] <- order

	c.JSON(http.StatusOK, gin.H{"code": code.HTTP_OK, "msg": "success"})
}