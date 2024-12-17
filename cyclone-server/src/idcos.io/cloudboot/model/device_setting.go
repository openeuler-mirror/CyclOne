package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
)

const (
	// InstallStatusPre 设备安装状态-等待安装
	InstallStatusPre = "pre_install"
	// InstallStatusIng 设备安装状态-正在安装
	InstallStatusIng = "installing"
	// InstallStatusFail 设备安装状态-安装失败
	InstallStatusFail = "failure"
	// InstallStatusSucc 设备安装状态-安装成功
	InstallStatusSucc = "success"
)

const (
	// InstallationPXE pxe安装方式
	InstallationPXE = "pxe"
	// InstallationImage 镜像安装方式
	InstallationImage = "image"
)

const (
	// IPSourceStatic IP来源-static
	IPSourceStatic = "static"
	// IPSourceDHCP IP来源-dhcp
	IPSourceDHCP = "dhcp"
)

// DeviceSetting 设备装机参数
type DeviceSetting struct {
	gorm.Model
	SN                    	string     `gorm:"column:sn"`                      		// 序列号
	HardwareTemplateID    	uint       `gorm:"column:hardware_template_id"`    		// 硬件配置模板ID
	SystemTemplateID      	uint       `gorm:"column:system_template_id"`      		// 系统安装模板
	ImageTemplateID       	uint       `gorm:"column:image_template_id"`       		// 镜像安装模板ID
	NeedExtranetIP        	string     `gorm:"column:need_extranet_ip"`        		// 是否需要配置外网IP
	ExtranetIPNetworkID   	uint       `gorm:"column:extranet_ip_network_id"`  		// 外网IP所属网段ID
	ExtranetIP            	string     `gorm:"column:extranet_ip"`             		// 外网IP
	IntranetIPNetworkID   	uint       `gorm:"column:intranet_ip_network_id"`  		// 内网IP所属网段ID
	IntranetIP            	string     `gorm:"column:intranet_ip"`             		// 内网IP
	InstallType           	string     `gorm:"column:install_type"`            		// 安装方式
	Status                	string     `gorm:"column:status"`                  		// 安装状态
	InstallProgress       	float64    `gorm:"column:install_progress"`        		// 安装进度
	InstallationStartTime 	*time.Time `gorm:"column:installation_start_time"` 		// 安装开始时间
	InstallationEndTime   	*time.Time `gorm:"column:installation_end_time"`   		// 安装结束时间
	Creator               	string     `gorm:"column:creator"`                 		// 创建人ID
	Updater               	string     `gorm:"column:updater"`                 		// 创建人ID
	NeedExtranetIPv6        string     `gorm:"column:need_extranet_ipv6"`        	// 是否需要配置外网IPv6
	ExtranetIPv6NetworkID   uint       `gorm:"column:extranet_ipv6_network_id"`  	// 外网IPv6所属网段IDv6
	ExtranetIPv6            string     `gorm:"column:extranet_ipv6"`             	// 外网IPv6
	NeedIntranetIPv6        string     `gorm:"column:need_intranet_ipv6"`        	// 是否需要配置外网IPv6
	IntranetIPv6NetworkID   uint       `gorm:"column:intranet_ipv6_network_id"`  	// 内网IPv6所属网段ID
	IntranetIPv6            string     `gorm:"column:intranet_ipv6"`             	// 内网IPv6
}

// TableName 指定数据库表名
func (DeviceSetting) TableName() string {
	return "device_setting"
}

// CombineDeviceSetting 装机设备查询条件列表
type CombineDeviceSetting struct {
	// 源节点名
	OriginNode string `json:"origin_node"`
	// 所属数据中心ID
	IDCID []uint `json:"idc_id"`
	// 所属数据中心ID
	ServerRoomID []uint `json:"server_room_id"`
	// 所属机架ID
	ServerCabinetID []uint `json:"server_cabinet_id"`
	// 所属机位ID
	ServerUsiteID []uint `json:"server_usite_id"`
	// 设备序列号(多个SN用英文逗号分隔)
	Sn string `json:"sn"`
	// 硬件配置模板ID
	HardwareTemplateID uint `json:"hardware_template_id"`
	// 镜像配置模板ID
	ImageTemplateID uint `json:"image_template_id"`
	// 装机状态。可选值  pre_install-等待安装; installing-正在安装; failure-安装失败; success-安装成功;
	Status string `json:"status"`

	FN                  string
	Category            string
	ServerCabinetNumber string
	ServerRoomName      string
	IntranetIP          string
	ExtranetIP          string
}

// IDeviceSetting 装机设备参数持久化接口
type IDeviceSetting interface {
	// SaveDeviceSetting 保存设备装机参数。若入参包含主键ID，则进行更新操作，否则进行新增操作。
	SaveDeviceSetting(*DeviceSetting) (err error)
	// CountDeviceSettingCombines 返回满足过滤条件的装机设备参数(不支持模糊查找)
	CountDeviceSettingCombines(cond *CombineDeviceSetting) (count int64, err error)
	// CountDeviceSettingByStatus 统计对应设备安装状态的数量
	CountDeviceSettingByStatus(status string) (count int64, err error)
	// GetDeviceSettingCombinesByCond 统计满足过滤条件的装机设备参数数量
	GetDeviceSettingCombinesByCond(cond *CombineDeviceSetting, orderby OrderBy, limiter *page.Limiter) (items []*DeviceSetting, err error)
	// AddDeviceSettings 批量添加设备装机参数
	AddDeviceSettings(...*DeviceSetting) error
	// UpdateDeviceSettingBySN 根据SN更新设备装机参数(仅会更新非零值的字段)
	UpdateDeviceSettingBySN(*DeviceSetting) (affected int64, err error)
	// GetDeviceSettingByDeviceSettingID 根据SN查询设备装机参数
	GetDeviceSettingBySN(sn string) (*DeviceSetting, error)
	// GetDeviceSettingByID 根据id查询设备装机参数
	GetDeviceSettingByID(id uint) (devSetting *DeviceSetting, err error)
	// DeleteDeviceSettingByID 删除指定ID的装机参数
	DeleteDeviceSettingByID(id uint) (ds *DeviceSetting, err error)
	// DeleteDeviceSettingBySN 删除指定SN的装机参数
	DeleteDeviceSettingBySN(sn string) (ds *DeviceSetting, err error)
	//UpdateInstallStatusAndProgressByID 更新安装状态和进度
	UpdateInstallStatusAndProgressByID(id uint, status string, inp float64) (affected int64, err error)
	// GetDeviceSettingsByInstallationTimeout 查询安装超时的设备装机记录列表。timeout-超时时间，单位秒。
	GetDeviceSettingsByInstallationTimeout(timeout int64) (items []*DeviceSetting, err error)
	// SetInstallationTimeout 为指定设备序列号的装机参数进行'安装超时'处理。
	SetInstallationTimeout(sns ...string) (affected int64, err error)
	// UpdateDeviceSettingIPConfigBySN 更新设备IP配置
	//UpdateDeviceSettingIPConfigBySN(sn string, extranet bool) (affected int64, err error)
}
