package main

import "idcos.io/cloudboot/logger"
import "io/ioutil"
import (
	"encoding/json"
)

// User 用户结构体
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoadUserFromFile 从本地文件加载用户/密码信息
func LoadUserFromFile(log logger.Logger, srcFile string) (user []*User, err error) {
	b, err := ioutil.ReadFile(srcFile)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var data struct {
		Content map[string][]*User `json:"content"`
	}
	if err = json.Unmarshal(b, &data); err != nil {
		return nil, err
	}
	if users, ok := data.Content["items"]; ok {
		return users, nil
	}
	return nil, nil
}
