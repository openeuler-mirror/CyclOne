package model

import "github.com/jinzhu/gorm"

const (
	// SysSettingInstallationTimeout 系统设置-安装超时时间
	SysSettingInstallationTimeout = "installation_timeout"
	// SysSettingAuthorization 系统设置键名-授权API
	SysSettingAuthorization = "authorization"
)

const (
	// DefInstallationTimeout 默认安装超时时间(单位秒)
	DefInstallationTimeout = 3600
)

// SystemSetting 系统设置
type SystemSetting struct {
	gorm.Model
	Key     string `gorm:"column:key"`
	Value   string `gorm:"column:value"`
	Desc    string `gorm:"column:desc"`
	Updater string `gorm:"column:updater"`
}

// TableName 指定数据库表名
func (SystemSetting) TableName() string {
	return "system_setting"
}

// ISystemSetting 系统设置数据库操作接口
type ISystemSetting interface {
	// GetSystemSetting 查询指定名称的系统配置
	GetSystemSetting(key string) (*SystemSetting, error)
	// GetSystemSetting4InstallatonTimeout 查询安装超时时间的系统设置值。若发生错误或者值不存在，则返回默认值。
	GetSystemSetting4InstallatonTimeout(defValue int64) (sec int64)
	// GetSystemSetting4AuthorizationAPIs 查询授权API配置
	GetSystemSetting4AuthorizationAPIs() ([]*AuthorizationAPI, error)
}

// AuthorizationAPI 授权API配置
type AuthorizationAPI struct {
	API   *APIMeta `json:"api"`
	Codes []string `json:"codes"`
}

// APIMeta API元信息
type APIMeta struct {
	Method string `json:"method"`
	URIReg string `json:"uri_regexp"`
	Desc   string `json:"desc"`
}
