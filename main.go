package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	_ "matching/config"
	"matching/engine"
	"matching/handler"
	"matching/utils/log"
	"matching/utils/middleware"
	"matching/utils/redis"
)

func main() {

	// 日志初始化
	log.InitLog()
	// redis初始化
	redis.InitRedis()
	// 引擎初始化
	engine.Init()

	httpEngine := gin.Default()
	// 注册中间件
	httpEngine.Use(middleware.AuthSign())
	// 交易标路由分组
	symbolGroup := httpEngine.Group("/symbol")
	// open
	symbolGroup.POST("/openMatching", handler.OpenMatching)
	// close
	symbolGroup.POST("/closeMatching", handler.CloseMatching)

	// 订单处理分组
	orderGroup := httpEngine.Group("/order")
	// 订单处理
	orderGroup.POST("/handleOrder", handler.HandleOrder)

	// 订单处理分组
	logGroup := httpEngine.Group("/log")
	// 订单处理
	logGroup.POST("/showLogs", handler.ShowLogs)

	httpEngine.Run(viper.GetString("server.port"))
}
