package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

const (
	LogSwitch bool = true // 日志开关
)

// 配置的一些初始化
func init() {
	initYaml()
}

// yaml配置文件的初始化
func initYaml() {

	// 把config.yaml里面的配置读取到viper里面，后面可以直接用viper.GetSring使用
	path, _ := os.Getwd()
	viper.AddConfigPath(path + "/config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("err:%s\n",err)
	}

	// example
	// addr := viper.GetString("redis.addr")
}