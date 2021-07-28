package utils

import (
	"math/rand"
	"strings"
	"time"
)

// IgnoreFile Check if the file is contains the ignore file
func IgnoreFile(ignoreFiles []string, filename string) bool {
	for _, ignoreFile := range ignoreFiles {
		if strings.Contains(filename, ignoreFile) {
			return true
		}
	}
	return false
}

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
