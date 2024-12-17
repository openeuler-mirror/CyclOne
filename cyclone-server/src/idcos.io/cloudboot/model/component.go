package model

import "github.com/jinzhu/gorm"

const (
	// ComponentAgent 组件-cloudboot-agent
	ComponentAgent = "agent"
	// ComponentHWServer 组件-hw-server
	ComponentHWServer = "hw-server"
	// ComponentPEConfig 组件-peconfig
	ComponentPEConfig = "peconfig"
	// ComponentWINConfig 组件-winconfig
	ComponentWINConfig = "winconfig"
	// OSConfigLog 系统配置日志
	OSConfigLog = "os-config"
	// ImageCloneLog 系统配置日志
	ImageCloneLog = "image-clone"
)

// ComponentLog 组件日志
type ComponentLog struct {
	gorm.Model
	SN        string `gorm:"column:sn"`
	Component string `gorm:"column:component"`
	LogData   string `gorm:"column:log"`
}

// TableName 指定数据库表名
func (ComponentLog) TableName() string {
	return "component_log"
}

// IComponentLog 组件日志操作接口
type IComponentLog interface {
	// SaveComponentLogBySN 保存组件日志
	SaveComponentLogBySN(*ComponentLog) error
	// GetComponentLog 查询指定设备的指定组件日志
	GetComponentLog(sn, component string) (*ComponentLog, error)
}
