package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"matching/config"
	"matching/engine"
	"matching/utils/code"
	"matching/utils/common"
	"matching/utils/redis"
	"net/http"
)

type CloseMatch struct {
	Symbol string `form:"symbol" binding:"required" comment:"交易标"`
}

func CloseMatching(c *gin.Context) {

	// 绑定参数
	var closeMatch CloseMatch
	if c.ContentType() == config.HttpContentFormData {
		if result := c.ShouldBind(&closeMatch); result != nil {
			c.JSON(http.StatusOK, gin.H{"code": code.HTTP_PARAMS_NOTEXISTS, "msg": "参数缺失：" + result.Error()})
			return
		}
	} else if c.ContentType() == config.HttpContentJson {
		if result := c.ShouldBindBodyWith(&closeMatch, binding.JSON); result != nil {
			c.JSON(http.StatusOK, gin.H{"code": code.HTTP_PARAMS_NOTEXISTS, "msg": "参数缺失：" + result.Error()})
			return
		}
	}

	// 判断是否为空
	if closeMatch.Symbol == "" {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_PARAMS_ERROR, "msg": "交易标参数不能为空"})
		return
	}

	// 判断是否在allow里面
	if !common.InArray(closeMatch.Symbol, AllowSymbol) {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_SYMBOL_NOTIN_ALLOWLIST, "msg": "交易标不在允许列表中"})
		return
	}

	// 判断该交易标引擎是否开启，从redis缓存中查询
	if !redis.HasSymbol(closeMatch.Symbol) {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_SYMBOL_MATCHINIG_OPEN_REPEAT, "msg": "交易标引擎未开启，无法关闭"})
		return
	}

	if err := engine.CloseEngine(closeMatch.Symbol); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": code.HTTP_SYMBOL_MATCHINIG_CLOSE_REPEAT, "msg": "交易标引擎关闭失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": code.HTTP_OK, "msg": "交易标引擎关闭成功"})
}