package model

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"matching/utils"
	"matching/utils/common"
	"matching/utils/enum"
	"reflect"
)

type Order struct {
	Accid int `json:"accid"`
	Action enum.OrderAction `json:"action"`
	Symbol string `json:"symbol"`
	OrderId string `json:"orderid"`
	Type enum.OrderType `json:"type"`
	Side enum.OrderSide `json:"side"`
	Amount decimal.Decimal `json:"amount"`
	Price decimal.Decimal `json:"price"`
	Timestamp int64 `json:"timestamp"`
}

func (o Order) ToJson() string {
	data, err := json.Marshal(o)
	if err != nil {
		utils.LogError("json marshal error", err, o.OrderId)
	}
	return string(data)
}

// 把缓存中的订单转换成Order结构
func (o Order) FromMap(ordermap map[string]interface{}) (Order, error) {
	var order Order
	if err := mapstructure.Decode(ordermap, &order); err != nil {
		return order, common.Errors(fmt.Sprintf("map decode error"))
	}

	return order, nil
}

// 把Order结构转换成缓存中的订单
func (o Order) ToMap() (map[string]interface{}, error) {

	out := make(map[string]interface{})

	// 通过反射获取信息
	v := reflect.ValueOf(o)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// 判断是不是结构体
	if v.Kind() != reflect.Struct {  // 非结构体返回错误提示
		return out, common.Errors(fmt.Sprintf("ToMap only accepts struct or struct pointer; got %T", v))
	}

	t := v.Type()
	// 遍历结构体字段
	// 指定tagName值为map中key;字段值为map中value
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get("json"); tagValue != "" {
			out[tagValue] = v.Field(i).Interface()
		}
	}
	return out, nil

}