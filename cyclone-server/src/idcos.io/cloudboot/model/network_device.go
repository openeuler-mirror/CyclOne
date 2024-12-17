package model

import (
	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
)

// NetworkDevice 网络设备
type NetworkDevice struct {
	gorm.Model
	IDCID            uint   `gorm:"column:idc_id"`
	ServerRoomID     uint   `gorm:"column:server_room_id"`
	ServerCabinetID  uint   `gorm:"column:server_cabinet_id"`
	FixedAssetNumber string `gorm:"column:fixed_asset_number"`
	SN               string `gorm:"column:sn"`
	Name             string `gorm:"column:name"`
	ModelNumber      string `gorm:"column:model"`
	Vendor           string `gorm:"column:vendor"`
	OS               string `gorm:"column:os"`
	Type             string `gorm:"column:type"`
	TOR              string `gorm:"column:tor"`
	Usage            string `gorm:"column:usage"`
	Status           string `gorm:"column:status"` //运营中、待启用
}

// NetworkDeviceCond 网络设备
type NetworkDeviceCond struct {
	gorm.Model
	IDCID               []uint `gorm:"column:idc_id"`
	ServerRoomID        []uint `gorm:"column:server_room_id"`
	ServerRoomName      string `gorm:"column:server_room_name"`
	ServerCabinetID     []uint `gorm:"column:server_cabinet_id"`
	ServerCabinetNumber string
	FixedAssetNumber    string `gorm:"column:fixed_asset_number"`
	SN                  string `gorm:"column:sn"`
	Name                string `gorm:"column:name"`
	ModelNumber         string `gorm:"column:model"`
	Vendor              string `gorm:"column:vendor"`
	OS                  string `gorm:"column:os"`
	Type                string `gorm:"column:type"`
	TOR                 string `gorm:"column:tor"`
	Usage               string `gorm:"column:usage"`
	Status              string `gorm:"column:status"` //运营中、待启用
}

// TableName 指定数据库表名
func (NetworkDevice) TableName() string {
	return "network_device"
}

// INetworkDevice 网络区域持久化接口
type INetworkDevice interface {
	// GetNetworkDevicesByCond 返回满足过滤条件的网络设备(不支持模糊查找)
	GetNetworkDevicesByCond(cond *NetworkDeviceCond, orderby OrderBy, limiter *page.Limiter) (item []*NetworkDevice, err error)
	// CountNetworkDevices 统计满足过滤条件的网络设备数量
	CountNetworkDevices(cond *NetworkDeviceCond) (count int64, err error)
	// RemoveNetworkDeviceByID 删除指定ID的网络设备
	RemoveNetworkDeviceByID(id uint) (err error)
	// GetIntranetSwitchBySN 查询设备所在机位的内网交换机
	GetIntranetSwitchBySN(sn string) (*NetworkDevice, error)
	// GetNetworkDeviceByID 查询指定ID的网络设备
	GetNetworkDeviceByID(id uint) (network *NetworkDevice, err error)
	// GetExtranetSwitchBySN 查询设备所在机位的外网交换机
	GetExtranetSwitchBySN(sn string) (*NetworkDevice, error)
	// SaveNetworkDevice 保存网络设备
	SaveNetworkDevice(na *NetworkDevice) (networkDevice *NetworkDevice, err error)
	// GetNetworkDeviceBySN 查询指定sn的网络设备
	GetNetworkDeviceBySN(sn string) (network []*NetworkDevice, err error)
	// GetNetworkDeviceByFixedAssetNumber 查询指定FixedAssetNumber的网络设备
	GetNetworkDeviceByFixedAssetNumber(FixedAssetNumber string) (network []*NetworkDevice, err error)
	// GetTORs 返回所有的TOR名称列表
	GetTORs() (items []string, err error)
	// GetTORBySN 返回目标设备所属的TOR名称
	GetTORBySN(sn string) (tor string, err error)
	// 根据指定的tor返回对应的网络设备
	GetNetworkDeviceByTORS(tor ...string) (nd []*NetworkDevice, err error)
	// 根据机架id查对端交换机，同时查出对端cabinetID
	GetPeerNetworkDeviceByCabinetID(id uint) (*NetworkDevice, error)
}
