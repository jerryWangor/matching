package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"matching/utils/log"
	"strconv"
)

var redisClient *redis.Client

func InitRedis() {

	defer func() {
		fuckedUp := recover() //recover() 捕获错误保存到变量中
		if fuckedUp != nil {
			log.Error("Redis连接异常：", fuckedUp)
		}
	}()

	addr := viper.GetString("redis.addr")
	password := viper.GetString("redis.password")
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       0,  // use default DB
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		panic(err)
	} else {
		log.Info(fmt.Sprintf("Connected to redis: %s", addr))
	}
}

func SaveSymbol(symbol string) {
	key := "matching:s:symbols"
	redisClient.SAdd(key, symbol)
}

func RemoveSymbol(symbol string) {
	key := "matching:s:symbols"
	redisClient.SRem(key, symbol)
}

func HasSymbol(symbol string) bool {
	key := "matching:s:symbols"
	symbols := redisClient.SMembers(key).Val()
	for _, v := range symbols {
		if v == symbol {
			return true
		}
	}
	return false
}

func GetSymbols() []string {
	key := "matching:s:symbols"
	return redisClient.SMembers(key).Val()
}

func SavePrice(symbol string, price decimal.Decimal) {
	key := "matching:s:price:" + symbol
	redisClient.Set(key, price.String(), 0)
}

func GetPrice(symbol string) decimal.Decimal {
	key := "matching:s:price:" + symbol
	priceStr := redisClient.Get(key).Val()
	result, err := decimal.NewFromString(priceStr)
	if err != nil {
		result = decimal.Zero
	}
	return result
}

func RemovePrice(symbol string) {
	key := "matching:s:price:" + symbol
	redisClient.Del(key)
}

// SaveOrder 保存订单
func SaveOrder(order map[string]interface{}) {
	symbol := order["symbol"].(string)
	orderId := order["orderId"].(string)
	timestamp := order["timestamp"].(float64) // time.Now().UnixMicro() 16位

	// 数据转换才能存入redis   redis: can't marshal enum.OrderAction (implement encoding.BinaryMarshaler)
	key := "matching:h:order:" + symbol + ":" + orderId
	redisClient.HMSet(key, order)

	// Zset(sorted_set类型) 创建以timestamp排序的数据
	key = "matching:z:orderids:" + symbol
	z := &redis.Z{
		Score:  timestamp,
		Member: orderId,
	}
	redisClient.ZAdd(key, *z)
}

func GetOrder(symbol string, orderId string) map[string]interface{} {
	// 从redis中查询订单
	var orderMap = make(map[string]interface{})
	key := "matching:h:order:" + symbol + ":" + orderId
	result, _ := redisClient.HGetAll(key).Result()
	for k, v := range result {
		orderMap[k] = v
	}
	return orderMap
}

func UpdateOrder(order map[string]interface{}) {
	symbol := order["symbol"].(string)
	orderId := order["orderId"].(string)
	timestamp, _ := strconv.ParseFloat(order["timestamp"].(string), 64) // time.Now().UnixMicro() 16位
	key := "matching:h:order:" + symbol + ":" + orderId
	redisClient.HMSet(key, order)

	// Zset(sorted_set类型) 创建以timestamp排序的数据
	key = "matching:z:orderids:" + symbol
	z := &redis.Z{
		Score:  timestamp,
		Member: orderId,
	}
	redisClient.ZAdd(key, *z)
}

func RemoveOrder(order map[string]interface{}) {
	symbol := order["symbol"].(string)
	orderId := order["orderId"].(string)
	timestamp, _ := strconv.ParseFloat(order["timestamp"].(string), 64) // time.Now().UnixMicro() 16位
	key := "matching:h:order:" + symbol + ":" + orderId
	redisClient.HDel(key)

	// Zset(sorted_set类型) 创建以timestamp排序的数据
	key = "matching:z:orderids:" + symbol
	z := &redis.Z{
		Score:  timestamp,
		Member: orderId,
	}

	redisClient.ZRem(key, *z)
}

func OrderExist(symbol string, orderId string) bool {
	key := "matching:h:order:" + symbol + ":" + orderId
	result, err := redisClient.HGetAll(key).Result()
	if err != nil {
		return true
	}
	if len(result)>0 {
		return true
	}
	return false
}

func GetOrderIdsWithSymbol(symbol string) []string {
	key := "matching:z:orderids:" + symbol
	return redisClient.ZRange(key, 0, 10).Val()
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
	redisClient.XAdd(a)
}

func SendTrade(symbol string, trade map[string]interface{}) {
	a := &redis.XAddArgs{
		Stream:       "matching:trades:" + symbol,
		MaxLenApprox: 1000,
		Values:       trade,
	}
	redisClient.XAdd(a)
}
