package handler

import (
	"github.com/gin-gonic/gin"
	"matching/engine"
	"matching/utils/cache"
	"matching/utils/code"
	"matching/utils/redis"
	"net/http"
)

type OpenMatch struct {
	Symbol string `form:"symbol" binding:"require"`
}

func OpenMatching(c *gin.Context) {
	// 绑定参数
	var openmatch OpenMatch
	if c.ShouldBind(&openmatch) != nil {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_PARAMS_NOTEXISTS, "msg": "参数缺失"})
		return
	}

	// 判断是否为空
	if openmatch.Symbol == "" {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_PARAMS_ERROR, "msg": "交易标参数不能为空"})
		return
	}

	// 判断该交易标引擎是否开启，从redis缓存中查询
	if redis.HasSymbol(openmatch.Symbol) {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_SYMBOL_MATCHINIG_OPEN_REPEAT, "msg": "交易标引擎重复开启"})
		return
	}

	// 开启交易标撮合引擎
	// 从redis缓存里面查询价格
	price := cache.GetPrice(openmatch.Symbol)
	// 开启某个交易标的撮合引擎
	if e := engine.NewEngine(openmatch.Symbol, price); e != nil {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_SYMBOL_MATCHINIG_OPEN_ERROR, "msg": "交易标引擎开启失败：" + e.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": code.HTTP_OK, "msg": "交易标引擎启动成功"})
}
