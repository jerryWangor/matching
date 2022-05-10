package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"matching/model"
	"matching/utils/code"
	"matching/utils/common"
	"net/http"
	"strconv"
)

func AuthSign() gin.HandlerFunc {
	return func (c *gin.Context) {

		var timesign model.TimeSign
		// 如果是 `GET` 请求，只使用 `Form` 绑定引擎（`query`）。
		// 如果是 `POST` 请求，首先检查 `content-type` 是否为 `JSON` 或 `XML`，
		// 然后再使用 `Form`（`form-data`）。
		if c.ShouldBind(&timesign) != nil {
			returnJson(c, code.HTTP_PARAMS_NOTEXISTS, "参数缺失")
			return
		}
		// 判断时间过期没
		nowtime := common.GetNowTimeStamp()
		if timesign.Time + viper.GetInt64("http.timeout") < nowtime {
			returnJson(c, code.HTTP_PARAMS_ERROR, "请求超时")
			return
		}
		// 判断sign
		lsign := common.GetMd5String(strconv.FormatInt(timesign.Time, 10) + viper.GetString("http.sign"))
		if lsign != timesign.Sign {
			returnJson(c, code.HTTP_PARAMS_ERROR, "验证失败")
			return
		}
		// 继续执行后面的函数
		c.Next()
	}
}

// 中间件退出，不执行后面的函数
func returnJson(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, gin.H{"code": code, "msg": msg})
	//终止后续接口调用，不加的话recover到异常后，还会继续执行接口里后续代码
	c.Abort()
}