package common

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"matching/utils/log"
	"math"
	"math/big"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

// GetNowTime 获取当前时间
func GetNowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// GetNowTimeStamp 获取当前时间-时间戳
func GetNowTimeStamp() int64 {
	return time.Now().Unix()
}

// GetNowDate 获取当前日期
func GetNowDate() string {
	return time.Now().Format("2006-01-02")
}

// TimeStampToString 时间戳转年季月时分秒字符串
func TimeStampToString(timestamp int64) string {
	if timestamp > 10000000000 {
		time1 := timestamp / 1000000
		time2 := timestamp % 1000000
		return time.Unix(time1,0).Format("2006-01-02 15:04:05") + "+" + strconv.FormatInt(time2, 10)
	} else {
		return time.Unix(timestamp,0).Format("2006-01-02 15:04:05")
	}
}

// GetWheres 获取where条件
func GetWheres(where []string) string {
	var wheres = strings.Join(where, " and ")
	return wheres
}

// GetMd5String md5加密
func GetMd5String(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

// IntToBytes int转字节
func IntToBytes(n int) []byte {
	data := int64(n)
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}

// BytesToInt 字节转int
func BytesToInt(bys []byte) int {
	bytebuff := bytes.NewBuffer(bys)
	var data int64
	binary.Read(bytebuff, binary.BigEndian, &data)
	return int(data)
}

// Errors 返回error类型错误，并记录error级别日志
func Errors(str string) error {
	log.Error(str)
	return errors.New(str)
}

// Debugs 返回加换行的字符串，并记录debug级别日志
func Debugs(str string) string {
	log.Debug(str)
	return WriteStringLn(str)
}

// ToMap 结构体转Map
func ToMap(o interface{}) (map[string]interface{}, error) {

	out := make(map[string]interface{})

	// 通过反射获取信息
	v := reflect.ValueOf(o)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// 判断是不是结构体
	if v.Kind() != reflect.Struct {  // 非结构体返回错误提示
		return out, Errors(fmt.Sprintf("ToMap only accepts struct or struct pointer; got %T", v))
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

// ToJson 结构体转json
func ToJson(o interface{}) string {
	data, err := json.Marshal(o)
	if err != nil {
		log.Error("json marshal error", err)
	}
	return string(data)
}

// WriteStringLn 字符串自动加上时间和换行符
func WriteStringLn(str string) string {
	return GetNowTime() + str + "\n"
}

// InArray 判断是否在数组中
func InArray(target string, strArray []string) bool {
	sort.Strings(strArray)
	index := sort.SearchStrings(strArray, target)
	if index < len(strArray) && strArray[index] == target {
		return true
	}
	return false
}

// RandInt64 取范围随机数
func RandInt64(min, max int64) int64 {
	if min > max {
		log.Error("the min is greater than max!")
	}

	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(max+1+i64Min))

		return result.Int64() - i64Min
	} else {
		result, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
		return min + result.Int64()
	}
}

// Abs 整数绝对值
func Abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}