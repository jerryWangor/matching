package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"matching/config"
	"matching/engine"
	"matching/utils/cache"
	"matching/utils/code"
	"matching/utils/redis"
	"net/http"
)

type OpenMatch struct {
	Symbol string `form:"symbol" binding:"required"`
}

func OpenMatching(c *gin.Context) {

	// 绑定参数
	var openMatch OpenMatch
	if c.ContentType() == config.HttpContentFormData {
		if result := c.ShouldBind(&openMatch); result != nil {
			c.JSON(http.StatusOK, gin.H{"code": code.HTTP_PARAMS_NOTEXISTS, "msg": "参数缺失：" + result.Error()})
			return
		}
	} else if c.ContentType() == config.HttpContentJson {
		if result := c.ShouldBindBodyWith(&openMatch, binding.JSON); result != nil {
			c.JSON(http.StatusOK, gin.H{"code": code.HTTP_PARAMS_NOTEXISTS, "msg": "参数缺失：" + result.Error()})
			return
		}
	}

	// 判断是否为空
	if openMatch.Symbol == "" {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_PARAMS_ERROR, "msg": "交易标参数不能为空"})
		return
	}

	// 判断该交易标引擎是否开启，从redis缓存中查询
	if redis.HasSymbol(openMatch.Symbol) {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_SYMBOL_MATCHINIG_OPEN_REPEAT, "msg": "交易标引擎重复开启"})
		return
	}

	// 开启交易标撮合引擎
	// 从redis缓存里面查询价格
	price := cache.GetPrice(openMatch.Symbol)
	// 开启某个交易标的撮合引擎
	if e := engine.NewEngine(openMatch.Symbol, price); e != nil {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_SYMBOL_MATCHINIG_OPEN_ERROR, "msg": "交易标引擎开启失败：" + e.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": code.HTTP_OK, "msg": "交易标引擎启动成功"})
}
