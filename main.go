package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "matching/config"
	"matching/handler"
	"matching/utils/middleware"
	"matching/utils/redis"
)

func main() {

	// 日志初始化
	//utils.InitLog()
	// redis初始化
	redis.InitRedis()

	engine := gin.Default()
	// 注册中间件
	engine.Use(middleware.AuthSign())
	// 交易标路由分组
	symbolGroup := engine.Group("/symbol")
	// open
	symbolGroup.POST("/openMatching", handler.OpenMatching)
	// close
	symbolGroup.POST("/closeMatching", handler.CloseMatching)

	// 订单处理分组
	orderGroup := engine.Group("/order")
	// 订单处理
	orderGroup.POST("/handleOrder", handler.HandleOrder)

	engine.Run(":8080")

}
