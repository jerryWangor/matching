package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/viper"
	"matching/model"
	"matching/utils/code"
	"matching/utils/common"
	"net/http"
	"strconv"
)

func AuthSign() gin.HandlerFunc {
	return func (c *gin.Context) {

		var timeSign model.TimeSign
		// 如果是 `GET` 请求，只使用 `Form` 绑定引擎（`query`）。
		// 如果是 `POST` 请求，首先检查 `content-type` 是否为 `JSON` 或 `XML`，
		// 然后再使用 `Form`（`form-data`）。
		//fmt.Println("request", c.Request)
		//c.Request.ParseForm()
		//for k, v := range c.Request.PostForm {
		//	fmt.Printf("k:%v - v:%v", k, v)
		//}
		//fmt.Println(c.Params)

		//d :=json.NewDecoder(c.Request.Body)
		//var aa map[string] interface{}
		//d.Decode(&aa)
		//fmt.Printf("%+v",aa)

		if c.ContentType() == "multipart/form-data" {
			if result := c.ShouldBind(&timeSign); result != nil {
				returnJson(c, code.HTTP_PARAMS_NOTEXISTS, "参数缺失：" + result.Error())
				return
			}
		} else if c.ContentType() == "application/json" {
			if result := c.ShouldBindBodyWith(&timeSign, binding.JSON); result != nil {
				returnJson(c, code.HTTP_PARAMS_NOTEXISTS, "参数缺失：" + result.Error())
				return
			}
		}

		// 判断时间过期没
		nowTime := common.GetNowTimeStamp()
		//time, _ := strconv.Atoi(timeSign.Time)
		if timeSign.Time + viper.GetInt64("http.timeout") < nowTime {
			returnJson(c, code.HTTP_PARAMS_ERROR, "请求超时")
			return
		}
		// 判断sign
		sign := common.GetMd5String(strconv.FormatInt(timeSign.Time, 10) + viper.GetString("http.sign"))
		if sign != timeSign.Sign {
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