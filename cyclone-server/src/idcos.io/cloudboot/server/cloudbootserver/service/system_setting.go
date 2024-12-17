package service

import (
	"fmt"
	"strings"

	"idcos.io/cloudboot/config"
)

// SystemLoginSetting 系统登录配置
type SystemLoginSetting struct {
	Channel      string `json:"channel"` // 系统登录通道
	URL          string `json:"url"`     // 系统登录URL
	RootEndpoint string `json:"-"`
}

const (
	// LoginChanUAMSSO 系统登录通道-UAM SSO
	LoginChanUAMSSO = "uam_sso"
)

// GetSystemLoginSetting 返回系统登录配置
func GetSystemLoginSetting(conf *config.Config) (*SystemLoginSetting, error) {
	return &SystemLoginSetting{
		Channel:      LoginChanUAMSSO,
		URL:          fmt.Sprintf("%s/login", strings.TrimSuffix(conf.UAM.RootEndpoint, "/")),
		RootEndpoint: conf.UAM.RootEndpoint,
	}, nil
}
