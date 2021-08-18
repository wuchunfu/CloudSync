package signx

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
)

func GetMd5Sign(params string) (string, error) {
	hash := md5.New()
	_, err := hash.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func GetMd5B64Sign(params string) (string, error) {
	hash := md5.New()
	_, err := hash.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}

func GetHmacMd5Sign(secret string, params string) (string, error) {
	mac := hmac.New(md5.New, []byte(secret))
	_, err := mac.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func GetHmacMd5B64Sign(secret string, params string) (string, error) {
	mac := hmac.New(md5.New, []byte(secret))
	_, err := mac.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

func GetSha1Sign(params string) (string, error) {
	hash := sha1.New()
	_, err := hash.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func GetSha256Sign(params string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func GetSha512Sign(params string) (string, error) {
	hash := sha512.New()
	_, err := hash.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func GetSha1B64Sign(params string) (string, error) {
	hash := sha1.New()
	_, err := hash.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}

func GetSha256B64Sign(params string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}

func GetSha512B64Sign(params string) (string, error) {
	hash := sha512.New()
	_, err := hash.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}

func GetHmacSha1Sign(secret string, params string) (string, error) {
	mac := hmac.New(sha1.New, []byte(secret))
	_, err := mac.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func GetHmacSha256Sign(secret string, params string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func GetHmacSha512Sign(secret string, params string) (string, error) {
	mac := hmac.New(sha512.New, []byte(secret))
	_, err := mac.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func GetHmacSha1B64Sign(secret string, params string) (string, error) {
	mac := hmac.New(sha1.New, []byte(secret))
	_, err := mac.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

func GetHmacSha256B64Sign(secret string, params string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

func GetHmacSha512B64Sign(secret string, params string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write([]byte(params))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}
