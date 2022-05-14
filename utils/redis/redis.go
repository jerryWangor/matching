package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"matching/utils"
)

var RedisClient *redis.Client

func InitRedis() {

	defer func() {
		fuckedUp := recover() //recover() 捕获错误保存到变量中
		if fuckedUp != nil {
			fmt.Println(fuckedUp)
		}
	}()

	addr := viper.GetString("redis.addr")
	password := viper.GetString("redis.password")
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       0,  // use default DB
	})

	_, err := RedisClient.Ping().Result()
	if err != nil {
		panic(err)
	} else {
		utils.LogInfo(fmt.Sprintf("Connected to redis: %s", addr))
	}
}

func SaveSymbol(symbol string) {
	key := "matching:symbols"
	RedisClient.SAdd(key, symbol)
}

func RemoveSymbol(symbol string) {
	key := "matching:symbols"
	RedisClient.SRem(key, symbol)
}

func HasSymbol(symbol string) bool {
	key := "matching:symbols"
	symbols := RedisClient.SMembers(key).Val()
	for _, v := range symbols {
		if v == symbol {
			return true
			break
		}
	}
	return false
}

func GetSymbols() []string {
	key := "matching:symbols"
	return RedisClient.SMembers(key).Val()
}

func SavePrice(symbol string, price decimal.Decimal) {
	key := "matching:price:" + symbol
	RedisClient.Set(key, price.String(), 0)
}

func GetPrice(symbol string) decimal.Decimal {
	key := "matching:price:" + symbol
	priceStr := RedisClient.Get(key).Val()
	result, err := decimal.NewFromString(priceStr)
	if err != nil {
		result = decimal.Zero
	}
	return result
}

func RemovePrice(symbol string) {
	key := "matching:price:" + symbol
	RedisClient.Del(key)
}

// SaveOrder 保存订单
func SaveOrder(order map[string]interface{}) {
	symbol := order["symbol"].(string)
	orderId := order["orderId"].(string)
	timestamp := order["timestamp"].(float64) // time.Now().UnixMicro() 16位
	action := order["action"].(string)

	//key := "matching:order:" + symbol + ":" + orderId + ":" + action
	key := "matching:order:" + symbol + ":" + orderId
	RedisClient.HMSet(key, order)

	// Zset(sorted_set类型) 创建以timestamp排序的数据
	key = "matching:orderids:" + symbol
	z := &redis.Z{
		Score:  timestamp,
		Member: orderId + ":" + action,
	}
	RedisClient.ZAdd(key, *z)
}

func GetOrder(symbol string, orderId string) map[string]interface{} {
	// 从redis中查询订单
	var maporder map[string]interface{}
	maporder = make(map[string]interface{})
	key := "matching:order:" + symbol + ":" + orderId
	result := RedisClient.HMGet(key)
	for _, v := range result.Val() {
		utils.LogDebug(fmt.Sprintf("getorder 的值是：%s", v))
	}
	return maporder
}

func UpdateOrder(order map[string]interface{}) {
	symbol := order["symbol"].(string)
	orderId := order["orderId"].(string)
	timestamp := order["timestamp"].(float64) // time.Now().UnixMicro() 16位
	action := order["action"].(string)
	key := "matching:order:" + symbol + ":" + orderId
	RedisClient.HMSet(key, order)

	// Zset(sorted_set类型) 创建以timestamp排序的数据
	key = "matching:orderids:" + symbol
	z := &redis.Z{
		Score:  timestamp,
		Member: orderId + ":" + action,
	}
	RedisClient.ZAdd(key, *z)
}

func RemoveOrder(order map[string]interface{}) {
	symbol := order["symbol"].(string)
	orderId := order["orderId"].(string)
	timestamp := order["timestamp"].(float64) // time.Now().UnixMicro() 16位
	action := order["action"].(string)
	key := "matching:order:" + symbol + ":" + orderId
	RedisClient.HDel(key)

	// Zset(sorted_set类型) 创建以timestamp排序的数据
	key = "matching:orderids:" + symbol
	z := &redis.Z{
		Score:  timestamp,
		Member: orderId + ":" + action,
	}

	RedisClient.ZRem(key, *z)
}

func OrderExist(symbol string, orderId string) bool {
	key := "matching:order:" + symbol + ":" + orderId
	result := RedisClient.HMGet(key)
	if result != nil {
		return true
	}
	return false
}

func GetOrderIdsWithSymbol(symbol string) []string {
	key := "matching:orderids:" + symbol
	return RedisClient.ZRange(key, 0, -1).Val()
}


// SendCancelResult 队列操作
/**
其中，matching:cancelresults:{symbol} 就是撤单结果的 MQ 所属的 Key，
matching:trades:{symbol} 则是成交记录的 MQ 所属的 Key。
可以看到，我们还根据不同 symbol 分不同 MQ，这样还方便下游服务可以根据需要实现分布式订阅不同 symbol 的 MQ。
*/
func SendCancelResult(symbol, orderId string, ok bool) {
	values := map[string]interface{}{"orderId": orderId, "ok": ok}
	a := &redis.XAddArgs{
		Stream:       "matching:cancelresults:" + symbol,
		MaxLenApprox: 1000,
		Values:       values,
	}
	RedisClient.XAdd(a)
}

func SendTrade(symbol string, trade map[string]interface{}) {
	a := &redis.XAddArgs{
		Stream:       "matching:trades:" + symbol,
		MaxLenApprox: 1000,
		Values:       trade,
	}
	RedisClient.XAdd(a)
}
