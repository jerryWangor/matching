package middleware

import (

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"log"
)

var RedisClient *redis.Client

func Init() {
	addr := viper.GetString("redis.addr")
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := RedisClient.Ping().Result()
	if err != nil {
		//panic(err)
	} else {
		log.Printf("Connected to redis: %s", addr)
	}
}