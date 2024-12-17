package model

import "github.com/jinzhu/gorm"

var (
	// SupportedRAIDLevels 支持的raid级别
	SupportedRAIDLevels = []string{"raid0", "raid1", "raid5", "raid6", "raid10", "raid1e"}
	// SupportedControllerIndexes RAID卡控制器索引(序号)。0-设备第一块RAID卡；1-设备第二块RAID卡；以此类推
	SupportedControllerIndexes = []string{"0", "1", "2", "3", "4"}
	// SupportedOOBUserLevels 带外用户特权级别
	SupportedOOBUserLevels = []string{"1", "2", "3", "4", "5"}
)

const (
	// CategoryRAID 硬件配置动作类别-RAID
	CategoryRAID = "raid"
	// CategoryBIOS 硬件配置动作类别-BIOS
	CategoryBIOS = "bios"
	// CategoryOOB 硬件配置动作类别-OOB
	CategoryOOB = "oob"
	// CategoryFW 硬件配置动作类别-固件升级
	CategoryFW = "firmware"
	// CategoryReboot 硬件配置动作类别-重启
	CategoryReboot = "reboot"
)

const (
	// ON 打开
	ON = "ON"
	// OFF 关闭
	OFF = "OFF"
)

const (
	// ActionRAIDClear 硬件配置动作-清除RAID配置
	ActionRAIDClear = "clear_settings"
	// ActionRAIDCreate 硬件配置动作-创建RAID阵列
	ActionRAIDCreate = "create_array"
	// ActionRAIDSetGlobalHotspare 硬件配置动作-设置全局热备盘
	ActionRAIDSetGlobalHotspare = "set_global_hotspare"
	// ActionRAIDInitDisk 硬件配置动作-初始化逻辑磁盘
	ActionRAIDInitDisk = "init_disk"
	// ActionRAIDSetJBOD 硬件配置动作-设置直通盘
	ActionRAIDSetJBOD = "set_jbod"

	// ActionReboot 硬件配置动作-重启
	ActionReboot = "reboot"
	// ActionOOBAddUser 硬件配置动作-添加OOB用户
	ActionOOBAddUser = "add_user"
	// ActionOOBResetBMC 硬件配置动作-重启BMC
	ActionOOBResetBMC = "reset_bmc"
	// ActionOOBSetIP 硬件配置动作-设置OOB BootOSIP
	ActionOOBSetIP = "set_ip"
	// ActionBIOSSet 硬件配置动作-设置BIOS项
	ActionBIOSSet = "set_bios"
	// ActionFWUpdatePkg 硬件配置动作-更新固件包
	ActionFWUpdatePkg = "update_package"
)

// DeviceHardwareSetting 装机参数-硬件配置参数
type DeviceHardwareSetting struct {
	gorm.Model
	// 设备序列号
	SN string
	// 当前配置项在配置项序列中的索引号（自0始）
	Index uint
	// 硬件配置动作类别 raid,oob,bios,firmware,reboot
	Category string
	// 硬件配置动作 reboot,add_user,reset_bmc,set_ip,set_bios,clear_settings,create_array,set_global_hotspare,init_disk,set_jbod,update_package
	Action string
	// 配置元数据
	Metadata map[string]string
	// 是否已实施 yes/no
	Applied string
}

// TableName 指定数据库表名
func (DeviceHardwareSetting) TableName() string {
	return "device_hardware_setting"
}

// IHardwareSetting 硬件配置操作接口
type IHardwareSetting interface {
	// OverwriteHardwareSettings 覆写指定设备的硬件配置参数
	OverwriteHardwareSettings(sn string, items ...*DeviceHardwareSetting) (err error)
	// GetHardwareSettingsBySN 返回指定设备的硬件配置配置参数
	GetHardwareSettingsBySN(sn string) (items []*DeviceHardwareSetting, err error)
	// UpdateHardwareSettingsApplied 更新满足过滤条件的硬件配置项实施状态
	UpdateHardwareSettingsApplied(cond *DeviceHardwareSetting, applied string) (affected int64, err error)
	// RedoHardwareSettings 将指定设备的硬件配置项设置为'未实施'状态。
	RedoHardwareSettings(sn string) (affected int64, err error)
}
