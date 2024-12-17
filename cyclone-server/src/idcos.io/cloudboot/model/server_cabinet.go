package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
)

const (
	// CabinetTypeServer 机架(柜)类型-通用服务器
	CabinetTypeServer = "server"
	// CabinetTypeNetworkDevice 机架(柜)类型-网络设备
	CabinetTypeNetworkDevice = "network_device"
	// CabinetTypeReserved 机架(柜)类型-预留
	CabinetTypeReserved = "reserved"
	// CabinetTypeKvmServer 机架(柜)类型-虚拟化服务器
	CabinetTypeKvmServer = "kvm_server"
)

const (
	// CabinetStatUnderConstruction 机架(柜)状态-建设中
	CabinetStatUnderConstruction = "under_construction"
	// CabinetStatNotEnabled 机架(柜)状态-未启用
	CabinetStatNotEnabled = "not_enabled"
	// CabinetStatEnabled 机架(柜)状态-已启用
	CabinetStatEnabled = "enabled"
	// CabinetStatLocked 机架(柜)状态-已锁定[允许回收不允许新分配资源]
	CabinetStatLocked = "locked"
	// CabinetStatOffline 机架(柜)状态-已下线
	CabinetStatOffline = "offline"
	// CabinetPowerOn 机架(柜)状态-开电
	CabinetPowerOn = "yes"
	// CabinetPowerOff 机架(柜)状态-关电
	CabinetPowerOff = "no"
)

const (
	// YES 是
	YES   = "yes"
	YesCh = "是"
	// NO 否
	NO      = "no"
	NoCh    = "否"
	Unknown = "unknown"
)

// ServerCabinet 机架(柜)
type ServerCabinet struct {
	gorm.Model
	IDCID          uint       `gorm:"column:idc_id"`
	ServerRoomID   uint       `gorm:"column:server_room_id"`
	NetworkAreaID  uint       `gorm:"column:network_area_id"`
	Number         string     `gorm:"column:number"`
	Height         uint       `gorm:"column:height"`
	Type           string     `gorm:"column:type"`
	NetworkRate    string     `gorm:"column:network_rate"`
	Current        string     `gorm:"column:current"`
	AvailablePower string     `gorm:"column:available_power"`
	MaxPower       string     `gorm:"column:max_power"`
	IsEnabled      string     `gorm:"column:is_enabled"`
	EnableTime     *time.Time `gorm:"column:enable_time"`
	IsPowered      string     `gorm:"column:is_powered"`
	PowerOnTime    *time.Time `gorm:"column:power_on_time"`
	PowerOffTime   *time.Time `gorm:"column:power_off_time"`
	Status         string     `gorm:"column:status"`
	Remark         string     `gorm:"column:remark"`
	Creator        string     `gorm:"column:creator"`
}

type ServerCabinetCond struct {
	IDCID           	[]uint `gorm:"column:idc_id"`
	ServerRoomID    	[]uint `gorm:"column:server_room_id"`
	ServerCabinetID    	[]uint `gorm:"column:server_cabinet_id"`
	NetworkAreaID   	[]uint `gorm:"column:network_area_id"`
	Number          	string `gorm:"column:number"`
	Type            	string `gorm:"column:type"`
	IsEnabled       	string `gorm:"column:is_enabled"`
	IsPowered       	string `gorm:"column:is_powered"`
	Status          	string `gorm:"column:status"`
	ServerRoomName  	string
	NetworkAreaName 	string
}

// TableName 指定数据库表名
func (ServerCabinet) TableName() string {
	return "server_cabinet"
}

// IServerCabinet 机架(柜)持久化接口
type IServerCabinet interface {
	// RemoveServerCabinetByID 删除指定ID的机架(柜)
	RemoveServerCabinetByID(id uint) (affected int64, err error)
	// SaveServerCabinet 保存机架(柜)
	SaveServerCabinet(*ServerCabinet) (affected int64, err error)
	// GetServerCabinetByID 返回指定ID的机架(柜)
	GetServerCabinetByID(id uint) (*ServerCabinet, error)
	// GetServerCabinetByNumber 返回指定编号的机架(柜)
	GetServerCabinetByNumber(serverRoomID uint, num string) (*ServerCabinet, error)
	// GetServerCabinetID 根据条件返回ID
	GetServerCabinetID(cond *ServerCabinet) (id []uint, err error)
	// CountServerCabinets 统计满足过滤条件的机架(柜)数量
	CountServerCabinets(cond *ServerCabinetCond) (count int64, err error)
	// GetServerCabinets 返回满足过滤条件的机架(柜)
	GetServerCabinets(cond *ServerCabinetCond, orderby OrderBy, limiter *page.Limiter) (items []*ServerCabinet, err error)
	//PowerOffServerCabinetByID 根据ID对机架(柜)下电
	PowerOffServerCabinetByID(id uint) (affected int64, err error)
	//PowerOnServerCabinetByID 根据ID对机架(柜)上电
	PowerOnServerCabinetByID(ids []uint) (affected int64, err error)
	//UpdateServerCabinetStatus 批量修改机架(柜)状态
	UpdateServerCabinetStatus(ids []uint, status string) (affected int64, err error)
	//UpdateServerCabinetType 批量修改机架(柜)类型
	UpdateServerCabinetType(ids []uint, typ string) (affected int64, err error)
	//UpdateServerCabinetRemakr 批量修改机架(柜)备注信息
	UpdateServerCabinetRemark(ids []uint, remark string) (affected int64, err error)	
	//GetServerCabinetCountByServerRoomID 根据给定的机房获取机房内的机架(柜)数
	GetServerCabinetCountByServerRoomID(id uint) (count int64, err error)
	// GetPeerCabinetByID 根据机架ID查询同组（同属于一个tor）的另一个机架
	//GetPeerCabinetByID(id uint) (*ServerCabinet, error)
	//GetCabinetOrderByFreeUsites 按照可用机位数倒叙顺序返回机架列表
	GetCabinetOrderByFreeUsites(req *ServerCabinet, physicalArea string) (mod []*OrderedCabientResp, err error)
}

type OrderedCabientResp struct {
	ServerCabinet
	AvailableUsiteCount uint `gorm:"column:available_usite_count"`
}
