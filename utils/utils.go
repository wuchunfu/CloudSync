package utils

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// IsDir 判断是否是一个目录
func IsDir(filePath string) bool {
	file, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return file.IsDir()
}

// GetFullPath 获取绝对路径
func GetFullPath(path string) string {
	abPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}
	return abPath
}

// GetFileMd5 获取文件的md5码
func GetFileMd5(filePath string) string {
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	md5hash := md5.New()
	sum := md5hash.Sum(file)
	return hex.EncodeToString(sum)
}

var IgnoreFiles = []string{".git", ".idea", ".swp"}

// IgnoreFile Check if the file is contains the ignore file
func IgnoreFile(filename string) bool {
	for _, ignoreFile := range IgnoreFiles {
		if strings.Contains(filename, ignoreFile) {
			return true
		}
	}
	return false
}
