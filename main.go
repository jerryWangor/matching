package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "matching/config"
	"matching/handler"
	"matching/utils"
)

func main() {

	// 日志初始化
	utils.InitLog()
	// redis初始化
	utils.InitRedis()

	engine := gin.Default()
	// 交易标路由分组
	symbolGroup := engine.Group("/symbol")
	// open
	symbolGroup.GET("/openMatching", handler.OpenMatching)
	// close
	symbolGroup.GET("/closeMatching", handler.CloseMatching)

	// 订单处理分组
	orderGroup := engine.Group("/order")
	// 订单处理
	orderGroup.GET("/handleOrder", handler.HandleOrder)

	engine.Run(":8080")

}
