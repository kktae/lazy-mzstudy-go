package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

func EncPW(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	md := hash.Sum(nil)
	// retStr := hex.EncodeToString(md)
	retStr := base64.StdEncoding.EncodeToString(md) // URLEncoding is not at go-resty
	return retStr
}

func GetJSESSIONId() string {
	md5Hash := GetMD5Hash(time.Now().Format("2006-01-02 15:04:05"))
	return fmt.Sprintf("%s.mz_was1", md5Hash)
}
