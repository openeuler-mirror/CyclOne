package model

import (
	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
)

const (
	//ILO 服务器ILO
	ILO = "ilo"
	//TGWIntranet 服务器TGW内网
	TGWIntranet = "tgw_intranet"
	//TGWExtranet 服务器TGW外网
	TGWExtranet = "tgw_extranet"
	//Intranet 服务器普通内网
	Intranet = "intranet"
	//Extranet 服务器普通外网
	Extranet = "extranet"
	//VIntranet 服务器虚拟化内网
	VIntranet = "v_intranet"
	//VExtranet 服务器虚拟化外网
	VExtranet = "v_extranet"
	// IP version
	IPv4 = "ipv4"
	IPv6 = "ipv6"
)

// NetworkTypeMap 网段类型中文和英文的对应表
// 可选值: ilo-服务器ILO; tgw_intranet-服务器TGW内网; tgw_extranet-服务器TGW外网; intranet-服务器普通内网; extranet-服务器普通外网; v_intranet-服务器虚拟化内网"))
var NetworkTypeMap = map[string]string{
	"服务器ILO":   ILO,
	"服务器TGW内网": TGWIntranet,
	"服务器TGW外网": TGWExtranet,
	"服务器普通内网":  Intranet,
	"服务器普通外网":  Extranet,
	"服务器虚拟化内网": VIntranet,
	"服务器虚拟化外网": VExtranet,
}

//IPNetwork 网段表结构
type IPNetwork struct {
	gorm.Model
	IDCID        uint   `gorm:"column:idc_id"`
	ServerRoomID uint   `gorm:"column:server_room_id"`
	Category     string `gorm:"column:category"`
	CIDR         string `gorm:"column:cidr"`
	Netmask      string `gorm:"column:netmask"`
	Gateway      string `gorm:"column:gateway"`
	IPPool       string `gorm:"column:ip_pool"`
	PXEPool      string `gorm:"column:pxe_pool"`
	Switches     string `gorm:"column:switches"`
	Vlan         string `gorm:"column:vlan"`
	Version	     string `gorm:"column:version"`
	Creator      string `gorm:"column:creator"`
}

type IPNetworkCond struct {
	CIDR           string
	Category       string
	ServerRoomID   []uint
	ServerRoomName string
	Switches       string
	NetworkAreaID  []uint
	Status         string
}

// TableName 指定数据库表名
func (IPNetwork) TableName() string {
	return "ip_network"
}

//IIPNetwork 网段持久化接口
type IIPNetwork interface {
	// RemoveIPNetworkByID 删除指定ID的网段
	RemoveIPNetworkByID(id uint) (affected int64, err error)
	// SaveIPNetwork 新增/更新网段
	SaveIPNetwork(*IPNetwork) (affected int64, err error)
	// GetIPNetworkByID 返回指定ID的网段
	GetIPNetworkByID(id uint) (*IPNetwork, error)
	// CountIPNetworks 统计满足过滤条件的网段数量
	CountIPNetworks(cond *IPNetworkCond) (count int64, err error)
	// GetIPNetworks 返回满足过滤条件的网段
	GetIPNetworks(cond *IPNetworkCond, orderby OrderBy, limiter *page.Limiter) (items []*IPNetwork, err error)
	// GetIntranetIPNetworksBySN 查询指定物理机的内网IP所属网段
	GetIntranetIPNetworksBySN(sn string) ([]*IPNetwork, error)
	// GetExtranetIPNetworksBySN 查询指定物理机的外网IP所属网段
	GetExtranetIPNetworksBySN(sn string) ([]*IPNetwork, error)
	// GetIPNetworksBySwitchNumber 根据设备编号查询网段信息
	GetIPNetworksBySwitchNumber(switchNumber string) ([]*IPNetwork, error)
	// GetIPv6NetworksBySN 查询指定物理机SN的IPv6内外网段
	GetIPv6NetworkBySN(sn string, category string) (*IPNetwork, error)
}
