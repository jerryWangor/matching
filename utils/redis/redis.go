package redis

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"matching/model"
	"matching/utils/common"
	"matching/utils/enum"
	"matching/utils/log"
	"strconv"
	"strings"
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
	action := enum.OrderAction(order["action"].(int)).String()
	timestamp := order["timestamp"].(float64) // time.Now().UnixMicro() 16位

	// 数据转换才能存入redis   redis: can't marshal enum.OrderAction (implement encoding.BinaryMarshaler)
	key := "matching:h:order:" + symbol + ":" + orderId + ":" + action
	redisClient.HMSet(key, order)

	// Zset(sorted_set类型) 创建以timestamp排序的数据
	key = "matching:z:orderids:" + symbol
	z := &redis.Z{
		Score:  timestamp,
		Member: orderId + ":" + action,
	}
	redisClient.ZAdd(key, *z)
}

func GetOrder(symbol, orderId, action string) map[string]interface{} {
	// 从redis中查询订单
	var orderMap = make(map[string]interface{})
	key := "matching:h:order:" + symbol + ":" + orderId + ":" + action
	result, _ := redisClient.HGetAll(key).Result()
	for k, v := range result {
		orderMap[k] = v
	}
	return orderMap
}

func UpdateOrder(order map[string]interface{}) {
	symbol := order["symbol"].(string)
	orderId := order["orderId"].(string)
	action := enum.OrderAction(order["action"].(int)).String()
	timestamp, _ := order["timestamp"].(float64) // time.Now().UnixMicro() 16位
	key := "matching:h:order:" + symbol + ":" + orderId + ":" + action
	redisClient.HMSet(key, order)

	// Zset(sorted_set类型) 创建以timestamp排序的数据
	key = "matching:z:orderids:" + symbol
	z := &redis.Z{
		Score:  timestamp,
		Member: orderId + ":" + action,
	}
	redisClient.ZAdd(key, *z)
}

func RemoveOrder(order map[string]interface{}) error {
	symbol := order["symbol"].(string)
	orderId := order["orderId"].(string)
	action := enum.OrderAction(order["action"].(int)).String()
	//action := enum.ActionCreate.String()
	// 删除hash
	key := "matching:h:order:" + symbol + ":" + orderId + ":" + action
	result := redisClient.Del(key)
	if result.Err() != nil {
		return result.Err()
	}
	// 删除zset
	key = "matching:z:orderids:" + symbol
	result = redisClient.ZRem(key, orderId + ":" + action)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func OrderExist(symbol, orderId, action string) bool {
	key := "matching:h:order:" + symbol + ":" + orderId + ":" + action
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
	return redisClient.ZRange(key, 0, -1).Val()
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
		MaxLenApprox: 10000,
		Values:       values,
	}
	redisClient.XAdd(a)
}

func SendTrade(symbol string, trade map[string]interface{}) {
	a := &redis.XAddArgs{
		Stream:       "matching:trades:" + symbol,
		MaxLenApprox: 10000,
		Values:       trade,
	}
	redisClient.XAdd(a)
}

func GetCancelResult(symbol string) map[string]string {
	strMap := make(map[string]string)
	stream := "matching:cancelresults:" + symbol
	// XRange正向查看 XRevRange反向查看 所以后面的+-要反着来
	result := redisClient.XRevRange(stream, "+", "-")
	val, _ := result.Result()
	for _, v := range val {
		jsons, _ := json.Marshal(v.Values)
		strMap[v.ID] = string(jsons)
	}
	return strMap
}

func GetTradeResult(symbol string) map[string]string {
	strMap := make(map[string]string)
	stream := "matching:trades:" + symbol
	result := redisClient.XRevRange(stream, "+", "-")
	val, _ := result.Result()
	for _, v := range val {
		jsons, _ := json.Marshal(v.Values)
		strMap[v.ID] = string(jsons)
	}
	return strMap
}

//  TopN K线图
func SetTopN(symbol string, num int, data map[string]interface{}) {
	// topN打算用有序集合进行保存
	key := "data:top" + strconv.Itoa(num) + ":s:symbol:" + symbol
	redisClient.Set(key, common.ToJson(data), 0)
}

func GetTopN(symbol string, num int) map[string]interface{} {
	var orderMap = make(map[string]interface{})
	key := "data:top" + strconv.Itoa(num) + ":s:symbol:" + symbol
	result, er := redisClient.Get(key).Result()
	if er != nil {
		log.Error("TopN 缓存获取失败，" + er.Error())
		return orderMap
	}
	err := json.Unmarshal([]byte(result), &orderMap)
	if err != nil {
		log.Error("TopN 数据解析失败")
	}
	return orderMap
}

// SetKData kdata把数据存成json，以timestamp排序
func SetKData(symbol string, timestamp int64, data string) {
	// Zset(sorted_set类型) 创建以timestamp排序的数据
	key := "data:kdata:z:" + symbol
	z := &redis.Z{
		Score:  float64(timestamp),
		Member: data,
	}
	redisClient.ZAdd(key, *z)
}

// GetKData 查询一段时间内的所有kdata
func GetKData(symbol string, time1, time2 int64) map[int64]model.KData {
	var orderMap = make(map[int64]model.KData)
	key := "data:kdata:z:" + symbol
	stime1 := strconv.FormatInt(time1, 10)
	stime2 := strconv.FormatInt(time2, 10)
	z := &redis.ZRangeBy{
		Min: stime1,
		Max: stime2,
	}
	result := redisClient.ZRangeByScore(key, *z)
	value, err := result.Result()
	if err != nil {
		log.Error("获取KData失败：" + err.Error())
	}
	for _, v := range value {
		var kdata model.KData
		dec := json.NewDecoder(strings.NewReader(v))
		dec.Decode(&kdata)
		orderMap[kdata.Timestamp] = kdata
	}
	return orderMap
}