package fileUtils

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/wuchunfu/CloudSync/middleware/config"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// IsDir 判断所给路径是否为文件夹
// IsDir returns true if given path is a dir,
// or returns false when it's a directory or does not exist.
func IsDir(filePath string) bool {
	file, err := os.Stat(filePath)
	return err == nil && file.IsDir()
}

// GetFullPath 获取绝对路径
func GetFullPath(path string) string {
	abPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Println(err)
	}
	return abPath
}

// IsFile 判断所给路径是否为文件
// IsFile returns true if given path is a file,
// or returns false when it's a directory or does not exist.
func IsFile(filePath string) bool {
	return !IsDir(filePath)
}

// FileExist 判断所给路径文件/文件夹是否存在
func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

// PathExists return true if given path exist.
func PathExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// Sha1f return file sha1 encode
func Sha1f(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	h := sha1.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// ReadFile 读取文件
func ReadFile(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer func(fi *os.File) {
		err := fi.Close()
		if err != nil {
			fmt.Println("Close error: ", err)
		}
	}(fi)
	fd, err := io.ReadAll(fi)
	return string(fd)
}

// GetFileSize 获取文件路径
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
