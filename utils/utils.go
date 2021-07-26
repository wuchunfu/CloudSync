package utils

import (
	"crypto/md5"
	"crypto/sha1"
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

func Sha1(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum([]byte("")))
}

func FileSha1(file *os.File) string {
	_sha1 := sha1.New()
	_, err := io.Copy(_sha1, file)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(_sha1.Sum(nil))
}

func MD5(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum([]byte("")))
}

func FileMD5(file *os.File) string {
	_md5 := md5.New()
	_, err := io.Copy(_md5, file)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(_md5.Sum(nil))
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetFileSize(filename string) int64 {
	var result int64
	err := filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	if err != nil {
		return 0
	}
	return result
}
