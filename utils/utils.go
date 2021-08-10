package utils

import (
	"math/rand"
	"os"
	"os/user"
	"strings"
	"time"
)

// CurrentUser 获取当前SSH连接的用户
func CurrentUser() string {
	currentUser, _ := user.Current()
	return currentUser.Username
}

// UserHome 获取用户的家目录
func UserHome() (string, error) {
	currentUser, err := user.Current()
	if err == nil {
		return currentUser.HomeDir, nil
	}
	return os.Getenv("HOME"), nil
}

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

// GenRandomString 生成随机字符串
// length 生成长度
// specialChar 是否生成特殊字符
func GenRandomString(length int, specialChar bool) string {
	letterBytes := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	special := "!@#%$*.="

	if specialChar {
		letterBytes = letterBytes + special
	}
	chars := []byte(letterBytes)
	if length == 0 {
		return ""
	}

	clen := len(chars)
	maxChar := 255 - (256 % clen)
	b := make([]byte, length)
	// storage for random bytes.
	r := make([]byte, length+(length/4))
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			return ""
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxChar {
				// Skip this number to avoid modulo bias.
				continue
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}
