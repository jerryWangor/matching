package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "matching/config"
	"matching/handler"
)

func main() {

	//addr := viper.GetString("redis.addr")
	//log.Println(addr)
	handler.Init()
	handler.LogDebug("test")
	//handler.LogInfo("test")
	block2()
}

// 方案2
func block2(){
	select{}
}