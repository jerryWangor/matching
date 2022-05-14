package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"matching/engine"
	"matching/model"
	"matching/utils/code"
	"matching/utils/enum"
	"matching/utils/redis"
	"net/http"
)

type handleOrder struct {
	Accid int `json:"accid" form:"accid" binding:"require" comment:"账号ID"`
	Action enum.OrderAction `json:"action" form:"action" binding:"require" comment:"0 挂单 1 撤单"`
	Symbol string `json:"symbol" form:"symbol" binding:"require" comment:"交易标"`
	OrderId string `json:"orderId" form:"orderId" binding:"require" comment:"订单ID"`
	Type enum.OrderType `json:"type" form:"type" binding:"require" comment:"竞价类型：0 普通交易"`
	Side enum.OrderSide `json:"side" form:"side" binding:"require" comment:"0 买 1 卖"`
	Amount decimal.Decimal `json:"amount" form:"amount" binding:"require" comment:"数量"`
	Price decimal.Decimal `json:"price" form:"price" binding:"require" comment:"价格"`
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
	// 调用分发订单
	err := engine.Dispatch(order)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_ORDER_HANDLE_ERROR, "msg": err})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_OK, "msg": "success"})
	}


}