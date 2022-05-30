package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/shopspring/decimal"
	"matching/config"
	"matching/engine"
	"matching/model"
	"matching/utils/cache"
	"matching/utils/code"
	"matching/utils/common"
	"matching/utils/enum"
	"net/http"
)

// 这里有个坑，表单不能接收0值，只需要把enum.OrderType前面加个*号就可以了
type handleOrder struct {
	//Accid int `json:"accid" form:"accid" binding:"required" comment:"账号ID"`
	Action *enum.OrderAction `json:"action" form:"action" binding:"required" comment:"0 挂单 1 撤单"`
	Symbol string `json:"symbol" form:"symbol" binding:"required" comment:"交易标"`
	OrderId string `json:"orderId" form:"orderId" binding:"required" comment:"订单ID"`
	Type *enum.OrderType `json:"type" form:"type" binding:"required" comment:"竞价类型：0 普通交易"`
	Side *enum.OrderSide `json:"side" form:"side" binding:"required" comment:"0 买 1 卖"`
	Amount decimal.Decimal `json:"amount" form:"amount" binding:"required" comment:"数量"`
	Price decimal.Decimal `json:"price" form:"price" binding:"required" comment:"价格"`
}

func HandleOrder(c *gin.Context) {
	// 绑定参数
	var hOrder handleOrder
	if c.ContentType() == config.HttpContentFormData {
		if result := c.ShouldBind(&hOrder); result != nil {
			c.JSON(http.StatusOK, gin.H{"code": code.HTTP_PARAMS_NOTEXISTS, "msg": "参数缺失：" + result.Error()})
			return
		}
	} else if c.ContentType() == config.HttpContentJson {
		if result := c.ShouldBindBodyWith(&hOrder, binding.JSON); result != nil {
			c.JSON(http.StatusOK, gin.H{"code": code.HTTP_PARAMS_NOTEXISTS, "msg": "参数缺失：" + result.Error()})
			return
		}
	}

	// 判断参数
	if hOrder.Symbol == "" {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_PARAMS_ERROR, "msg": "交易标参数不能为空"})
		return
	}

	// 判断是否在allow里面
	if !common.InArray(hOrder.Symbol, AllowSymbol) {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_SYMBOL_NOTIN_ALLOWLIST, "msg": "交易标不在允许列表中"})
		return
	}

	// 判断Action
	if !hOrder.Action.Valid() {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_ORDER_HANDLE_ACTION_ERROR, "msg": "action error"})
		return
	}

	// 判断type
	if !hOrder.Type.Valid() {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_ORDER_HANDLE_TYPE_ERROR, "msg": "type error"})
		return
	}

	// 判断side
	if !hOrder.Side.Valid() {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_ORDER_HANDLE_SIDE_ERROR, "msg": "side error"})
		return
	}

	// 挂单判断
	if hOrder.Action.String() == enum.ActionCreate.String() {
		// 买单和卖单的价格不能为0
		if hOrder.Price.IsZero() {
			c.JSON(http.StatusOK, gin.H{"code": code.HTTP_ORDER_HANDLE_PRICE_ISZERO, "msg": "价格不能为0"})
			return
		}

		// 买单和卖单的数量不能为0
		if hOrder.Amount.IsZero() {
			c.JSON(http.StatusOK, gin.H{"code": code.HTTP_ORDER_HANDLE_AMOUNT_ISZERO, "msg": "数量不能为0"})
			return
		}
	}

	// 判断该交易标引擎是否开启，从redis缓存中查询
	if !cache.HasSymbol(hOrder.Symbol) {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_SYMBOL_MATCHINIG_OPEN_REPEAT, "msg": "交易标引擎未启动"})
		return
	}

	// 写入redis并发给通道
	order := model.Order{
		//Accid: hOrder.Accid,
		Action: *hOrder.Action,
		Symbol: hOrder.Symbol,
		OrderId: hOrder.OrderId,
		Side: *hOrder.Side,
		Type: *hOrder.Type,
		Amount: hOrder.Amount,
		Price: hOrder.Price,
	}
	// 调用分发订单
	err := engine.Dispatch(order)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_ORDER_HANDLE_ERROR, "msg": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_OK, "msg": "success"})
	}


}