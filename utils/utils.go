package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/wuchunfu/CloudSync/config"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// IsDir 判断是否是一个目录
func IsDir(filePath string) bool {
	file, err := os.Stat(filePath)
	return err == nil && file.IsDir()
}

// GetFullPath 获取绝对路径
func GetFullPath(path string) string {
	abPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}
	return abPath
}

// GetFileMd5 获取文件的 md5 值
func GetFileMd5(filePath string) string {
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("Generate md5 fail: ", filePath, err)
		return ""
	}
	md5hash := md5.New()
	sum := md5hash.Sum(file)
	return hex.EncodeToString(sum)
}

// GenerateMd5 生成 md5 值
func GenerateMd5(filePath string) string {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal("Generate md5 fail: ", filePath, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)

	md5hash := md5.New()
	_, copyErr := io.Copy(md5hash, file)
	if copyErr != nil {
		return ""
	}
	return fmt.Sprintf("%x", md5hash.Sum([]byte("")))
}

//var IgnoreFiles = []string{".git", ".idea", ".swp", ".swx"}

// IgnoreFile Check if the file is contains the ignore file
func IgnoreFile(filename string) bool {
	for _, ignoreFile := range config.GlobalObject.IgnoreFiles {
		if strings.Contains(filename, ignoreFile) {
			return true
		}
	}
	return false
}
