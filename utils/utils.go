package utils

import (
	"math/rand"
	"time"
)

// RandomString 导出随机字符串
func RandomString(n int) string {
	var byteStr = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = byteStr[rand.Intn(len(byteStr))]
	}
	return string(result)
}
