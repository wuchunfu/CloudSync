package cryptox

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
)

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
