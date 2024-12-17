package model

import "github.com/jinzhu/gorm"

// PlatformConfig 平台配置
type PlatformConfig struct {
	gorm.Model
	Name    string `sql:"not null;unique;"`
	Content string `sql:"type:longtext;"`
}

// IPlatformConfig 平台配置
type IPlatformConfig interface {
	// GetPlatformConfigByName 查询指定名称的平台配置项
	GetPlatformConfigByName(name string) (*PlatformConfig, error)
}
