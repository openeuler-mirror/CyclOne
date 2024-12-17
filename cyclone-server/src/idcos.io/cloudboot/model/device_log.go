package model

import "github.com/jinzhu/gorm"

const (
	// DeviceLogInstallType 安装类型的主机日志
	DeviceLogInstallType = "install"
	// DeviceLogHistoryType 安装历史类型的主机日志
	DeviceLogHistoryType = "install_history"
)

// DeviceLog 系统安装进度日志
type DeviceLog struct {
	gorm.Model
	SN              string `gorm:"column:sn"`                // 设备序列号
	DeviceSettingID uint   `gorm:"column:device_setting_id"` // 设备装机参数ID
	LogType         string `gorm:"column:type"`              // 进度日志类型
	Title           string `gorm:"column:title"`             // 进度日志标题
	Content         string `gorm:"column:content"`           // 进度日志内容
}

// TableName 表名信息
func (DeviceLog) TableName() string {
	return "device_log"

}

// IDeviceLog 系统安装进度日志接口
type IDeviceLog interface {
	// GetDeviceLogsByDeviceSettingID 根据装机参数id获取装机日志信息
	GetDeviceLogsByDeviceSettingID(id uint) (deviceLogs []DeviceLog, err error)
	// UpdateDeviceLogType 修改操作系统安装进度记录
	UpdateDeviceLogType(settingID uint, fromLogType, toLogType string) (affected int64, err error)
	// SaveDeviceLog 新增操作系统安装进度记录
	SaveDeviceLog(log *DeviceLog) (affected int64, err error)
}
