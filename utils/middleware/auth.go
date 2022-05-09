package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"matching/model"
	"matching/utils/common"
	"net/http"
	"strconv"
)

func AuthSign() gin.HandlerFunc {
	return func (c *gin.Context) {

		var flag = true
		var msg string
		var timesign model.TimeSign
		// 如果是 `GET` 请求，只使用 `Form` 绑定引擎（`query`）。
		// 如果是 `POST` 请求，首先检查 `content-type` 是否为 `JSON` 或 `XML`，
		// 然后再使用 `Form`（`form-data`）。
		if c.ShouldBind(&timesign) != nil {
			flag = false
			msg = "参数错误"
		}

		// 判断时间过期没
		nowtime := common.GetNowTimeStamp()
		log.Println("time", timesign.Time)
		log.Println("ctime", common.GetNowTimeStamp())
		log.Println("nowtime", nowtime)
		if timesign.Time + viper.GetInt64("http.timeout") < nowtime {
			flag = false
			msg = "请求超时"
		}
		// 判断sign
		lsign := common.GetMd5String(strconv.FormatInt(timesign.Time, 10) + viper.GetString("http.sign"))
		if lsign != timesign.Sign {
			flag = false
			msg = "sign验证失败"
		}

		if flag == false {
			c.JSON(http.StatusOK, gin.H{"code": msg})
			//终止后续接口调用，不加的话recover到异常后，还会继续执行接口里后续代码
			c.Abort()
		}

		c.Next()
	}
}