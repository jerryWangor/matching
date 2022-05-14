package handler

import (
	"github.com/gin-gonic/gin"
	"matching/engine"
	"matching/utils/code"
	"matching/utils/redis"
	"net/http"
)

type CloseMatch struct {
	Symbol string `form:"symbol" binding:"require" comment:"交易标"`
}

func CloseMatching(c *gin.Context) {

	// 绑定参数
	var closeMatch CloseMatch
	if c.ShouldBind(&closeMatch) != nil {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_PARAMS_NOTEXISTS, "msg": "参数缺失"})
		return
	}

	// 判断是否为空
	if closeMatch.Symbol == "" {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_PARAMS_ERROR, "msg": "交易标参数不能为空"})
		return
	}

	// 判断该交易标引擎是否开启，从redis缓存中查询
	if !redis.HasSymbol(closeMatch.Symbol) {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_SYMBOL_MATCHINIG_OPEN_REPEAT, "msg": "交易标引擎未开启，无法关闭"})
		return
	}

	if err := engine.CloseEngine(closeMatch.Symbol); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_SYMBOL_MATCHINIG_CLOSE_REPEAT, "msg": "交易标引擎关闭失败"})
	}
}