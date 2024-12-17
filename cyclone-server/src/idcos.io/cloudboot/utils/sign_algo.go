package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"
	"time"
)

const secretKey = "idcos_cloudboot$2017"

// GetSign 获取URL签名
func GetSign(salt string) string {
	signStr := "salt" + salt + "secret_key" + secretKey
	h := md5.New()
	h.Write([]byte(signStr))
	return strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil)))
}

// GetRespSign 获取返回值签名
func GetRespSign(stat, msg interface{}) string {
	signStr := "status" + stat.(string) + "message" + msg.(string) + "secret_key" + secretKey
	h := md5.New()
	h.Write([]byte(signStr))
	return strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil)))
}

// GenSalt 生成salt(随机8位字符串)
func GenSalt() string {
	t := time.Now()
	h := md5.New()
	io.WriteString(h, "crazyof.me")
	io.WriteString(h, t.String())
	salt := fmt.Sprintf("%x", h.Sum(nil))
	return salt[0:8]
}
