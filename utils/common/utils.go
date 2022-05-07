package common

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

// 获取当前时间
func GetNowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// 获取当前日期
func GetNowDate() string {
	return time.Now().Format("2006-01-02")
}

// 获取where条件
func GetWheres(where []string) string {
	var wheres = strings.Join(where, " and ")
	return wheres
}

// md5加密
func GetMd5String(b []byte) string {
	return fmt.Sprintf("%x", md5.Sum(b))
}

// int转字节
func IntToBytes(n int) []byte {
	data := int64(n)
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}

// 字节转int
func BytesToInt(bys []byte) int {
	bytebuff := bytes.NewBuffer(bys)
	var data int64
	binary.Read(bytebuff, binary.BigEndian, &data)
	return int(data)
}
