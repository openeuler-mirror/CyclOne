package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
)

const (
	//DevOperStatRunWithAlarm 设备运行状态：运行中(需告警)
	DevOperStatRunWithAlarm = "run_with_alarm"
	//DevOperStatRunWithoutAlarm 运行中(无需告警)
	DevOperStatRunWithoutAlarm = "run_without_alarm"
	//DevOperStatReinstalling 重装中
	DevOperStatReinstalling = "reinstalling"
	//DevOperStatMoving 搬迁中
	DevOperStatMoving = "moving"
	//DevOperStatPreRetire 待退役
	DevOperStatPreRetire = "pre_retire"
	//DevOperStatRetiring 退役中
	DevOperStatRetiring = "retiring"
	//DevOperStateRetired 已退役
	DevOperStateRetired = "retired"
	//DevOperStatPreDeploy 待部署
	DevOperStatPreDeploy = "pre_deploy"
	//DevOperStatOnShelve 已上架
	DevOperStatOnShelve = "on_shelve"
	//DevOperStatRecycling 回收中
	DevOperStatRecycling = "recycling"
	//DevOperStatMaintaining 维护中
	DevOperStatMaintaining = "maintaining"
	//DevOperStatPreMove 待搬迁
	DevOperStatPreMove = "pre_move"
	//DevOperStatInStore 库房中
	DevOperStatInStore = "in_store"

	//DevUsageSpecialDev 特殊设备
	DevUsageSpecialDev = "特殊设备"

	// PowerStatusOn 电源状态-开
	PowerStatusOn = "power_on"
	// PowerStatusOff 电源状态-关
	PowerStatusOff = "power_off"
)

const (
	// DevArchX8664 x86_64架构
	DevArchX8664 = "x86_64"
	// DevArchAarch64 aarch64架构
	DevArchAarch64 = "aarch64"
	// DevArchPPC64 ppc64架构
	DevArchPPC64 = "ppc64"
)

// Device 设备
type Device struct {
	gorm.Model
	FixedAssetNumber string    `gorm:"column:fixed_asset_number"`
	SN               string    `gorm:"column:sn"`
	Vendor           string    `gorm:"column:vendor"`
	DevModel         string    `gorm:"column:model"`
	Arch             string    `gorm:"column:arch"`
	Usage            string    `gorm:"column:usage"`
	Category         string    `gorm:"column:category"`
	IDCID            uint      `gorm:"column:idc_id"`
	ServerRoomID     uint      `gorm:"column:server_room_id"`
	CabinetID        uint      `gorm:"column:server_cabinet_id"`
	USiteID          *uint     `gorm:"column:server_usite_id"`
	StoreRoomID      uint      `gorm:"column:store_room_id"`
	VCabinetID       uint      `gorm:"column:virtual_cabinet_id"`
	HardwareRemark   string    `gorm:"column:hardware_remark"`
	RAIDRemark       string    `gorm:"column:raid_remark"`
	OOBInit          string    `gorm:"column:oob_init"`
	StartedAt        time.Time `gorm:"column:started_at"`
	OnShelveAt       time.Time `gorm:"column:onshelve_at"`
	OperationStatus  string    `gorm:"column:operation_status"`
	PowerStatus      string    `gorm:"power_status"` // 电源状态
	OOBIP            string    `gorm:"column:oob_ip"`
	OOBUser          string    `gorm:"column:oob_user"`
	OOBPassword      string    `gorm:"column:oob_password"`
	OOBAccessible    *string   `gorm:"column:oob_accessible"` //带外纳管状态:yes|no
	CPUSum           uint      `gorm:"column:cpu_sum"`
	CPU              string    `gorm:"column:cpu"`
	MemorySum        uint      `gorm:"column:memory_sum"`
	Memory           string    `gorm:"column:memory"`
	DiskSum          uint      `gorm:"column:disk_sum"`
	Disk             string    `gorm:"column:disk"`
	DiskSlot         string    `gorm:"column:disk_slot"`
	NIC              string    `gorm:"column:nic"`
	NICDevice        string    `gorm:"column:nic_device"`
	BootOSIP         string    `gorm:"column:bootos_ip"`
	BootOSMac        string    `gorm:"column:bootos_mac"`
	Motherboard      string    `gorm:"column:motherboard"`
	RAID             string    `gorm:"column:raid"`
	OOB              string    `gorm:"column:oob"`
	BIOS             string    `gorm:"column:bios"`
	Fan              string    `gorm:"column:fan"`
	Power            string    `gorm:"column:power"`
	HBA              string    `gorm:"column:hba"`
	PCI              string    `gorm:"column:pci"`
	Switch           string    `gorm:"column:switch"`
	LLDP             string    `gorm:"column:lldp"`
	Extra            string    `gorm:"column:extra"`
	OriginNode       string    `gorm:"column:origin_node"`    		
	OriginNodeIP     string    `gorm:"column:origin_node_ip"`		
	OperUserID       string    `gorm:"column:operation_user_id"`	//未使用
	Creator          string    `gorm:"column:creator"`
	Updater          string    `gorm:"column:updater"`
	Remark           string    `gorm:"column:remark"`
	OrderNumber      string    `gorm:"column:order_number"`
}

// BeforeSave 保存设备信息前的钩子方法。
// 防止将空字符串写入类型为JSON的数据库字段中引发报错。
func (dev *Device) BeforeSave() (err error) {
	replaceIfBlank(&dev.CPU, EmptyJSONObject)
	replaceIfBlank(&dev.Memory, EmptyJSONObject)
	replaceIfBlank(&dev.Disk, EmptyJSONObject)
	replaceIfBlank(&dev.DiskSlot, EmptyJSONObject)
	replaceIfBlank(&dev.NIC, EmptyJSONObject)
	replaceIfBlank(&dev.Motherboard, EmptyJSONObject)
	replaceIfBlank(&dev.RAID, EmptyJSONObject)
	replaceIfBlank(&dev.OOB, EmptyJSONObject)
	replaceIfBlank(&dev.BIOS, EmptyJSONObject)
	replaceIfBlank(&dev.Fan, EmptyJSONObject)
	replaceIfBlank(&dev.Power, EmptyJSONObject)
	replaceIfBlank(&dev.HBA, EmptyJSONObject)
	replaceIfBlank(&dev.PCI, EmptyJSONObject)
	replaceIfBlank(&dev.Switch, EmptyJSONObject)
	replaceIfBlank(&dev.LLDP, EmptyJSONObject)
	replaceIfBlank(&dev.Extra, EmptyJSONObject)
	replaceIfBlank(&dev.OOBInit, EmptyJSONObject)
	//replaceIfBlank(dev.OOBAccessible, Unknown)
	return
}

// CombinedDevice 设备信息及设备装机参数联合结构体
type CombinedDevice struct {
	Device
	DeployStatus    				string  	`gorm:"column:deploy_status"`
	InstallProgress 				float64 	`gorm:"column:install_progress"`
	IntranetIP           			string 		`gorm:"column:intranet_ip"`
	ExtranetIP           			string 		`gorm:"column:extranet_ip"`
	IntranetIPv6           			string 		`gorm:"column:intranet_ipv6"`
	ExtranetIPv6           			string 		`gorm:"column:extranet_ipv6"`	
	OS                   			string 		`gorm:"column:os"`
	ImageTemplateID      			int    		`gorm:"column:image_tpl_id"`
	ImageTemplateName    			string 		`gorm:"column:image_tpl_name"`
	HardwareTemplateID   			int    		`gorm:"column:hardware_tpl_id"`
	HardwareTemplateName 			string 		`gorm:"column:hardware_name"`
	Hostname             			string 		`gorm:"column:hostname"`
	InspectionRunStatus  			string 		`gorm:"column:inspection_run_status"`
	InspectionResult     			string 		`gorm:"column:inspection_result"`
	InspectionRemark     			string 		`gorm:"column:inspection_remark"`
	DeviceLifecycleDeatail
}

// CombinedDeviceCond 复杂的设备信息查询条件
type CombinedDeviceCond struct {
	//Device
	IDCID             []uint
	ServerRoomID      []uint
	ServerCabinet     string
	ServerRoomName    string
	ServerUsiteNumber string
	USiteID           []uint
	PhysicalArea      string
	FixedAssetNumber  string
	SN                string
	DevModel          string
	Vendor            string
	Usage             string
	Category          string
	OperationStatus   string
	IntranetIP        string
	ExtranetIP        string
	HardwareRemark    string
	IP                string //模糊搜索内外网
	OOBAccessible     string
	// 部署状态，即OS安装状态
	DeployStatus string
	PreDeployed  bool // 标识是否是预部署状态物理机（没有任何安装记录的物理机）
	ID           []uint
}

// TableName 指定数据库表名
func (Device) TableName() string {
	return "device"
}

// CollectedDevice 采集到的设备信息结构体
type CollectedDevice struct {
	gorm.Model
	OriginNode   string `gorm:"column:origin_node"`
	OriginNodeIP string `gorm:"column:origin_node_ip"`
	SN           string `gorm:"column:sn"`
	Vendor       string `gorm:"column:vendor"`
	ModelName    string `gorm:"column:model"`
	Arch         string `gorm:"column:arch"`
	CPUSum       uint   `gorm:"column:cpu_sum"`
	CPU          string `gorm:"column:cpu"`
	MemorySum    uint   `gorm:"column:memory_sum"`
	Memory       string `gorm:"column:memory"`
	DiskSum      uint   `gorm:"column:disk_sum"`
	Disk         string `gorm:"column:disk"`
	DiskSlot     string `gorm:"column:disk_slot"`
	NIC          string `gorm:"column:nic"`
	NICDevice    string `gorm:"column:nic_device"`
	BootOSIP     string `gorm:"column:bootos_ip"`
	BootOSMac    string `gorm:"column:bootos_mac"`
	Motherboard  string `gorm:"column:motherboard"`
	RAID         string `gorm:"column:raid"`
	OOB          string `gorm:"column:oob"`
	BIOS         string `gorm:"column:bios"`
	Fan          string `gorm:"column:fan"`
	Power        string `gorm:"column:power"`
	HBA          string `gorm:"column:hba"`
	PCI          string `gorm:"column:pci"`
	Switch       string `gorm:"column:switch"`
	LLDP         string `gorm:"column:lldp"`
	Extra        string `gorm:"column:extra"`
}

// TableName 指定数据库表名
func (CollectedDevice) TableName() string {
	return "device"
}


// BeforeSave 保存设备信息前的钩子方法。
// 防止将空字符串写入类型为JSON的数据库字段中引发报错。
func (dev *CollectedDevice) BeforeSave() (err error) {
	replaceIfBlank(&dev.CPU, EmptyJSONObject)
	replaceIfBlank(&dev.Memory, EmptyJSONObject)
	replaceIfBlank(&dev.Disk, EmptyJSONObject)
	replaceIfBlank(&dev.DiskSlot, EmptyJSONObject)
	replaceIfBlank(&dev.NIC, EmptyJSONObject)
	replaceIfBlank(&dev.Motherboard, EmptyJSONObject)
	replaceIfBlank(&dev.RAID, EmptyJSONObject)
	replaceIfBlank(&dev.OOB, EmptyJSONObject)
	replaceIfBlank(&dev.BIOS, EmptyJSONObject)
	replaceIfBlank(&dev.Fan, EmptyJSONObject)
	replaceIfBlank(&dev.Power, EmptyJSONObject)
	replaceIfBlank(&dev.HBA, EmptyJSONObject)
	replaceIfBlank(&dev.PCI, EmptyJSONObject)
	replaceIfBlank(&dev.Switch, EmptyJSONObject)
	replaceIfBlank(&dev.LLDP, EmptyJSONObject)
	replaceIfBlank(&dev.Extra, EmptyJSONObject)
	return
}

// IDevice 物理机设备
type IDevice interface {
	// RemoveDeviceByID 删除指定ID的物理机
	RemoveDeviceByID(id uint) (affected int64, err error)
	// RemoveDeviceBySN 删除指定SN的物理机
	RemoveDeviceBySN(sn string) (affected int64, err error)	
	// SaveDevice 保存物理机
	SaveDevice(*Device) (affected int64, err error)
	// UpdateDevice 修改物理机
	UpdateDevice(*Device) (affected int64, err error)
	// UpdateDeviceBySN 更新目标设备的
	UpdateDeviceBySN(*Device) (affected int64, err error)
	// GetDeviceByID 返回指定ID的物理机
	GetDeviceByID(id uint) (*Device, error)
	// CountDevices 统计满足过滤条件的物理机数量
	CountDevices(cond *Device) (count int64, err error)
	// GetDevices 返回满足过滤条件的物理机
	GetDevices(cond *Device, orderby OrderBy, limiter *page.Limiter) (items []*Device, err error)
	//GetDeviceBySN 根据SN查找设备
	GetDeviceBySN(SN string) (*Device, error)
	// GetDeviceBySNOrMAC 根据SN或者mac地址查询设备
	GetDeviceBySNOrMAC(snOrMAC string) (*Device, error)
	//GetDeviceBySN 根据SN查找设备
	GetDeviceByFixAssetNumber(fixAssetNum string) (*Device, error)
	// CountCombinedDevices 统计满足过滤条件的记录数量
	CountCombinedDevices(cond *CombinedDeviceCond) (count int64, err error)
	// GetCombinedDevices 返回满足过滤条件的设备及其装机参数列表
	GetCombinedDevices(cond *CombinedDeviceCond, orderby OrderBy, limiter *page.Limiter) (items []*CombinedDevice, err error)
	// SaveCollectedDeviceBySN 保存(更新/新增)采集到的设备信息
	SaveCollectedDeviceBySN(*CollectedDevice) error
	//GetDevicesByUSiteID 根据机位信息查询设备
	GetDevicesByUSiteID(id uint) (items []*Device, err error)
	//GetDevicesByUSiteIDS 根据机位信息查询设备
	GetDevicesByUSiteIDS(id []uint, usage string) (items []*Device, err error)
	// GetDeviceQuerys(param string)
	GetDeviceQuerys(param string) (*DeviceQueryParamResp, error)
	GetMaxFixedAssetNumber(month string) (fixedAssetNumber string, err error)
	// GetDeviceByStartedAt 根据启用日期查询在该日期之前的设备
	GetDeviceByStartedAt(started_date string) (items []*Device, err error)
	// GetDeviceByOperationStatus 根据运营状态查找设备
	GetDeviceByOperationStatus(started_date string) (items []*Device, err error)	
}

type DeviceQueryParamResp struct {
	ParamName string      `json:"param_name"`
	List      []ParamList `json:"list"`
}

type ParamList struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ParamListINT struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
