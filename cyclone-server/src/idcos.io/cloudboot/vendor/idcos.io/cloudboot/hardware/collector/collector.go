package collector

import (
	"encoding/json"
	"errors"

	"idcos.io/cloudboot/hardware/raid"
	"idcos.io/cloudboot/logger"
)

const (
	// HP 厂商-惠普
	HP = "HP"
	// Dell 厂商-戴尔
	Dell = "Dell"
	// Huawei 厂商-华为
	Huawei = "Huawei"
	// xFusion 厂商-超聚变
	XFusion = "xFusion"	
	// Lenovo 厂商-联想
	Lenovo = "Lenovo"
	// H3C 厂商-华三
	H3C = "H3C"
	// Inspur 厂商-浪潮
	Inspur = "Inspur"
	// Sugon 厂商-曙光
	Sugon = "Sugon"
	// UNIS 厂商-紫光
	UNIS = "UNIS"
	// suma 厂商-曙光
	Suma = "suma"		
	// Supermicro 厂商-超微
	Supermicro = "Supermicro"
	// Greatwall 厂商-长城
	Greatwall = "Greatwall"

	// DefaultCollector 默认采集器名称
	DefaultCollector = "BOOTOS"
)

var (
	// ErrUnregisteredCollector 未注册的采集器
	ErrUnregisteredCollector = errors.New("unregistered collector")
)

// Collector 设备信息采集器
type Collector interface {
	// SetDebug 设置是否开启debug。若开启debug，会将关键日志信息写入console。
	SetDebug(debug bool)
	// SetLog 更换日志实现。默认情况下内部无日志实现。
	SetLog(log logger.Logger)
	// 采集并返回当前设备基本信息
	BASE() (isVM bool, sn, vendor, model, arch string, err error)
	// CPU 采集并返回当前设备的CPU信息
	CPU() (*CPU, error)
	// Memory 采集并返回当前设备内存信息
	Memory() (*Memory, error)
	// Motherboard 采集并返回当前设备的主板信息
	Motherboard() (*Motherboard, error)
	// Disk 采集并返回当前设备的逻辑磁盘信息
	Disk() (*Disk, error)
	// DiskSlot 采集并返回当前设备的磁盘槽位(物理驱动器)信息
	DiskSlot() (*DiskSlot, error)
	// NIC 采集并返回当前设备的网卡信息
	NIC() (*NIC, error)
	// HBA 采集并返回当前设备的HBA卡信息。若当前设备不包含HBA卡，则返回值都为nil。
	HBA() (*HBA, error)
	// OOB 采集并返回当前设备的OOB信息。
	OOB() (*OOB, error)
	// BIOS 采集并返回当前设备的BIOS信息。
	BIOS() (*BIOS, error)
	// RAID 采集并返回当前设备的RAID信息。
	RAID() (*RAID, error)
	// PCI 采集并返回当前设备的所有PCI插槽信息。
	PCI() (*PCI, error)
	// Fan 采集并返回当前设备的所有风扇信息。
	Fan() (*Fan, error)
	// Power 采集并返回当前设备的电源信息。
	Power() (*Power, error)
	// LLDP 采集并返回当前设备LLDP信息。
	LLDP() (*LLDP, error)
	// IDRAC 采集并返回Dell iDRAC信息。
	IDRAC() (*IDRAC, error)
	// ILO 采集并返回HP iLO信息。
	ILO() (*ILO, error)
	// Backplane 采集并返回背板信息。
	Backplane() (*Backplane, error)
	// Extra 执行采集脚本并返回采集到的信息。若执行采集脚本时发生错误，则丢弃该错误，继续执行后续脚本。
	// 采集脚本需满足以下条件：
	// 1、脚本的开始位置通过shebang指定执行该脚本的程序，如'#! /usr/bin/env python'。
	// 2、脚本执行完毕后，需要将采集到的设备信息以JSON Object格式字符串写入stdout。
	// 3、多个脚本输出的JSON Object属性名不能重复，否则在合并多个JSON Object过程中可能导致数据丢失。
	Extra(scripts [][]byte) *Extra
}

// Backplane 背板
type Backplane struct {
	Items []BackplaneItem `json:"items"`
}

// BackplaneItem 背板条目
type BackplaneItem struct {
	Name            string `json:"name"`             // 背板名称
	FirmwareVersion string `json:"firmware_version"` // 背板固件版本号
}

// ILO HP iLO
type ILO struct {
	FirmwareDate    string `json:"firmware_date"`
	FirmwareVersion string `json:"firmware_version"`
}

// IDRAC Dell iDRAC
type IDRAC struct {
	RACDateTime        string `json:"rac_date_time"`
	FirmwareVersion    string `json:"firmware_version"`
	FirmwareBuild      string `json:"firmware_build"`
	LastFirmwareUpdate string `json:"last_firmware_update"`
}

// NIC 网卡
type NIC struct {
	Items []NICItem `json:"items"`
}

// NICItem 网卡信息
type NICItem struct {
	Name            string `json:"name"`             // 网卡逻辑名称
	Mac             string `json:"mac"`              // mac地址
	IP              string `json:"ip"`               // BootOS BootOSIP
	Businfo         string `json:"businfo"`          // 网口
	Designation     string `json:"designation"`      // 槽位
	Side            string `json:"side"`             // 内/外置
	Speed           string `json:"speed"`            // 速率
	Company         string `json:"company"`          // 厂商
	ModelName       string `json:"model_name"`       // 型号
	FirmwareVersion string `json:"firmware_version"` // 固件版本
	Link            string `json:"link"`             //网卡链接状态
}

// ToJSON 序列化为JSON
func (nic NIC) ToJSON() []byte {
	b, _ := json.Marshal(nic)
	return b
}

// CPU CPU信息。
// 总核数 = 物理CPU个数 X 每颗物理CPU的核数
// 总逻辑CPU数 = 物理CPU个数 X 每颗物理CPU的核数 X 超线程数
type CPU struct {
	TotalCores     int           `json:"total_cores"`     // 物理CPU总核心数
	Threads        int           `json:"threads"`         // 超线程数
	TotalLogicals  int           `json:"total_logicals"`  // 逻辑CPU数量
	TotalPhysicals int           `json:"total_physicals"` // 物理CPU数量
	Physicals      []PhysicalCPU `json:"physicals"`       // 物理CPU列表
}

// PhysicalCPU 物理CPU
type PhysicalCPU struct {
	ModelName  string `json:"model_name"`  // 型号
	ClockSpeed string `json:"clock_speed"` // 主频
	Cores      int    `json:"cores"`       // 单颗物理CPU的核心数
}

// ToJSON 序列化为JSON
func (cpu CPU) ToJSON() []byte {
	b, _ := json.Marshal(cpu)
	return b
}

// Memory 内存
type Memory struct {
	TotalSizeMB int            `json:"total_size_mb"` // 当前内存总容量，单位MB。
	TotalSize   int64          `json:"total_size"`    // 当前内存总容量，单位Byte。
	Items       []MemoryDevice `json:"items"`         // 物理内存条列表
}

// MemoryDevice 插在内存插槽上的内存条
type MemoryDevice struct {
	Locator string `json:"locator"`
	Size    string `json:"size"`  // 容量，8192 MB
	Type    string `json:"type"`  // 类型，如DDR3
	Speed   string `json:"speed"` // 速率
}

// ToJSON 序列化为JSON
func (m Memory) ToJSON() []byte {
	b, _ := json.Marshal(m)
	return b
}

// Disk 逻辑磁盘
type Disk struct {
	TotalSizeGB int        `json:"total_size_gb"` // 取整后的磁盘总容量，单位GB。
	TotalSize   int64      `json:"total_size"`    // 磁盘总容量，单位Byte。
	Items       []DiskItem `json:"items"`
}

// DiskItem 逻辑磁盘条目
type DiskItem struct {
	Name string `json:"name"` // 名称，如'/dev/sda'
	Size string `json:"size"` // 容量，如'599.6 GB, 599550590976 bytes'
}

// ToJSON 序列化为JSON
func (d Disk) ToJSON() []byte {
	b, _ := json.Marshal(d)
	return b
}

// DiskSlot 磁盘槽位(物理驱动器)
type DiskSlot struct {
	Items []raid.PhysicalDrive `json:"items"`
}

// ToJSON 序列化为JSON
func (d DiskSlot) ToJSON() []byte {
	b, _ := json.Marshal(d)
	return b
}

// Motherboard 主板
type Motherboard struct {
	Manufacturer string `json:"manufacturer"`  // 厂商名
	ProductName  string `json:"product_name"`  // 产品名
	SerialNumber string `json:"serial_number"` // 序列号
}

// ToJSON 序列化为JSON
func (board Motherboard) ToJSON() []byte {
	b, _ := json.Marshal(board)
	return b
}

// RAID RAID
type RAID struct {
	Items []raid.Controller `json:"items"`
}

// RAIDArray RAID阵列
type RAIDArray struct {
	Level string `json:"level"`
}

// ToJSON 序列化为JSON
func (raid RAID) ToJSON() []byte {
	b, _ := json.Marshal(raid)
	return b
}

// OOB OOB
type OOB struct {
	Network  *OOBNetwork `json:"network,omitempty"`
	User     []OOBUser   `json:"user,omitempty"`
	Firmware string      `json:"firmware,omitempty"` // 固件版本
}

// OOBNetwork 带外网络信息
type OOBNetwork struct {
	IPSrc   string `json:"ip_src,omitempty"` // IP来源。可选值: static|dhcp
	IP      string `json:"ip,omitempty"`
	Netmask string `json:"netmask,omitempty"`
	Gateway string `json:"gateway,omitempty"`
}

// OOBUser 带外用户信息
type OOBUser struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Password       string `json:"password"` //这个值采集不到
	PrivilegeLevel string `json:"privilege_level"`
}

// ToJSON 序列化为JSON
func (oob OOB) ToJSON() []byte {
	b, _ := json.Marshal(oob)
	return b
}

// BIOS BIOS
type BIOS struct {
	Vendor  string `json:"vendor,omitempty"`  // 厂商
	Version string `json:"version,omitempty"` // 固件版本号
}

// ToJSON 序列化为JSON
func (bios BIOS) ToJSON() []byte {
	b, _ := json.Marshal(bios)
	return b
}

// Power 电源
type Power struct {
	// TODO 定义字段
}

// ToJSON 序列化为JSON
func (p Power) ToJSON() []byte {
	b, _ := json.Marshal(p)
	return b
}

// Fan 风扇
type Fan struct {
	Items []FanItem `json:"items,omitempty"`
}

// FanItem 风扇
type FanItem struct {
	ID    string `json:"id"`    // 风扇传感器ID
	Speed string `json:"speed"` // 转速
}

// ToJSON 序列化为JSON
func (fan Fan) ToJSON() []byte {
	b, _ := json.Marshal(fan)
	return b
}

// PCI PCI
type PCI struct {
	TotalSlots int        `json:"total_slots"`
	Items      []SlotItem `json:"slots,omitempty"`
}

// SlotItem PCI插槽
type SlotItem struct {
	Designation  string `json:"designation"`   // 槽位
	Type         string `json:"type"`          // 设备类型
	CurrentUsage string `json:"current_usage"` // 当前使用情况
}

// ToJSON 序列化为JSON
func (pci PCI) ToJSON() []byte {
	b, _ := json.Marshal(pci)
	return b
}

// HBA HBA卡
type HBA struct {
	Items []HBAItem `json:"items,omitempty"`
}

// HBAItem 单块HBA卡
type HBAItem struct {
	Host            string   `json:"host"`
	WWPNs           []string `json:"wwpns"`
	WWNNs           []string `json:"wwnns"`
	FirmwareVersion string   `json:"firmware_version"`
}

// ToJSON 序列化为JSON
func (hba HBA) ToJSON() []byte {
	b, _ := json.Marshal(hba)
	return b
}

// LLDP LLDP采集到的交换机信息
type LLDP map[string]interface{}

// ToJSON 序列化为JSON
func (lldp LLDP) ToJSON() []byte {
	if lldp == nil {
		return nil
	}
	b, _ := json.Marshal(lldp)
	return b
}

// Extra 自定义扩展属性
type Extra map[string]interface{}

// ToJSON 序列化为JSON
func (extra Extra) ToJSON() []byte {
	b, _ := json.Marshal(extra)
	return b
}
