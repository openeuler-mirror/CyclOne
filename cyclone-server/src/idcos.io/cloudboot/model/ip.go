package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
)

const (
	//PXEIP PXE用IP
	PXEIP = "pxe"
	//BusinessIP 业务用IP
	BusinessIP = "business"
)

const (
	//IPUsed IP是否已被占用-是
	IPUsed = "yes"
	//IPNotUsed IP是否已被占用-是
	IPNotUsed = "no"
	//IPDisabled IP是否已被占用-不可用
	IPDisabled = "disabled"
)

const (
	// IPScopeIntranet IP作用范围-内网
	IPScopeIntranet = "intranet"
	// IPScopeExtranet IP作用范围-外网
	IPScopeExtranet = "extranet"
)

// IP BootOSIP
type IP struct {
	gorm.Model
	IPNetworkID uint    `gorm:"column:ip_network_id"`
	Category    string  `gorm:"column:category"`
	IP          string  `gorm:"column:ip"`
	IsUsed      string  `gorm:"column:is_used"`
	SN          string  `gorm:"column:sn"`
	Scope       *string `gorm:"column:scope"`
	Remark      *string `gorm:"column:remark"`
	ReleaseDate time.Time `gorm:"column:release_date"`
}

type IPCombined struct {
	IP
	//固资编号
	FixedAssetNumber string `gorm:"column:fixed_asset_number"`
}

type IPPageCond struct {
	ID               []uint
	IPNetworkID      []uint
	CIDR             string
	Category         string
	IP               string
	IsUsed           string
	SN               string
	FixedAssetNumber string
	Scope            *string
}

// TableName 指定数据库表名
func (IP) TableName() string {
	return "ip"
}

//IIP IP数据操作接口
type IIP interface {
	CountIPs(cond *IPPageCond) (count int64, err error)
	GetIPs(cond *IPPageCond, orderby OrderBy, limiter *page.Limiter) (items []*IPCombined, err error)
	GetIPByIP(ipaddr, scope string) (ip IP, err error)
	GetReleasableIP()(items []*IP, err error)
	GetNetWorkBySN(sn string, category string) ([]IPAndIPNetworkUnion, error)
	// AssignIntranetIP 按照内置规则给设备分配一个内网业务IP
	AssignIntranetIP(sn string) (intraIP *IP, err error)
	// AssignExtranetIP 按照内置规则给设备分配一个外网业务IP
	AssignExtranetIP(sn string) (extraIP *IP, err error)
	// ReleaseIP 为目标设备释放内/外网业务IP
	ReleaseIP(sn string, scope string) (affected int64, err error)
	// ReserveIP 为目标设备回收内/外网业务IP并保留IP一段时间，设置释放日期
	ReserveIP(sn string, scope string, releasedate time.Time) (affected int64, err error)	
	// AssignIP 分配IP
	AssignIP(sn, scope string, id uint) error
	// AssignIP 分配IP
	AssignIPByIP(sn, scope string, ip string) error
	// UnassignIP 取消分配的IP
	UnassignIP(id uint) error
	// UnassignIPsBySN 释放指定SN的IP资源
	UnassignIPsBySN(sn string) error
	// GetIPByID 返回指定ID的IP
	GetIPByID(id uint) (*IP, error)
	// GetIPBySNAndScope 返回被指定SN占用的IP
	GetIPBySNAndScope(sn, scope string) (*IP, error)
	// SaveIP
	SaveIP(ip *IP) (affected int64, err error)
	// CreateIP
	CreateIP(ip *IP) error
	// GetAvailableIPByIPNetworkID
	GetAvailableIPByIPNetworkID(ipnetworkid uint)(*IP, error)
	// GetLastIPv6ByIPNetworkID
	GetLastIPv6ByIPNetworkID(ipnetworkid uint)(*IP, error)
	// AssignIPv6 获取已有业务IPv6记录或按规则新分配一个业务IP
	AssignIPv6(sn, ipscope string) (*IP, error)
	// ReleaseIPv6 为目标设备释放内/外网业务IPv6
	ReleaseIPv6(sn string, scope string) (affected int64, err error)	
}

//IPAndIPNetworkUnion 设备网络配置查询结果
type IPAndIPNetworkUnion struct {
	IP      string  `json:"ip" gorm:"column:ip"`
	Netmask string  `json:"netmask" gorm:"column:netmask"`
	Gateway string  `json:"gateway" gorm:"column:gateway"`
	Scope   *string `json:"scope" gorm:"column:scope"`
	Version	string 	`json:"version" gorm:"column:version"`
}

// TableName 指定数据库表名
func (IPAndIPNetworkUnion) TableName() string {
	return "ip"
}
