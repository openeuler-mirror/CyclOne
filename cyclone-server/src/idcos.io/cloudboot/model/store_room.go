package model

import (
	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
)

// StoreRoom 库房
type StoreRoom struct {
	gorm.Model
	IDCID            uint   `gorm:"column:idc_id"`
	Name             string `gorm:"column:name"`
	FirstServerRoom  string `gorm:"column:first_server_room"`
	City             string `gorm:"column:city"`
	Address          string `gorm:"column:address"`
	StoreRoomManager string `gorm:"column:store_room_manager"`
	VendorManager    string `gorm:"column:vendor_manager"`
	Status           string `gorm:"column:status"`
	Creator          string `gorm:"column:creator"`
}

// StoreRoomCond 库房查询条件
type StoreRoomCond struct {
	ID               []uint
	IDCID            []uint
	Name             string
	FirstServerRoom  string
	City             string
	Address          string
	StoreRoomManager string
	VendorManager    string
	Status           string
	Creator          string
}

// TableName 指定数据库表名
func (StoreRoom) TableName() string {
	return "store_room"
}

// IStoreRoom 库房持久化接口
type IStoreRoom interface {
	// RemoveStoreRoomByID 删除指定ID的库房
	RemoveStoreRoomByID(id uint) (affected int64, err error)
	// SaveStoreRoom 保存库房
	SaveStoreRoom(*StoreRoom) (affected int64, err error)
	// UpdateStoreRoom更新库房信息
	UpdateStoreRoom(srs []*StoreRoom) (affected int64, err error)
	// GetStoreRoomByID 返回指定ID的库房
	GetStoreRoomByID(id uint) (*StoreRoom, error)
	// GetStoreRoomByName 返回指定Name的库房
	GetStoreRoomByName(n string) (*StoreRoom, error)
	// CountServerRooms 统计满足过滤条件的库房数量
	CountStoreRooms(cond *StoreRoomCond) (count int64, err error)
	// GetStoreRooms 返回满足过滤条件的库房
	GetStoreRooms(cond *StoreRoomCond, orderby OrderBy, limiter *page.Limiter) (items []*StoreRoom, err error)
}

// VirtualCabinet
type VirtualCabinet struct {
	gorm.Model
	StoreRoomID uint   `gorm:"column:store_room_id"`
	Number      string `gorm:"column:number"`
	Remark      string `gorm:"column:remark"`
	Status      string `gorm:"column:status"`
	Creator     string `gorm:"column:creator"`
}

// TableName 指定数据库表名
func (VirtualCabinet) TableName() string {
	return "virtual_cabinet"
}

// IVirtualCabinet 持久化接口
type IVirtualCabinet interface {
	// RemoveVirtualCabinetByID 删除指定ID的虚拟货架
	RemoveVirtualCabinetByID(id uint) (affected int64, err error)
	// SaveVirtualCabinet 保存库房
	SaveVirtualCabinet(*VirtualCabinet) (affected int64, err error)
	// UpdateStoreRoom更新库房信息
	//UpdateStoreRoom(srs []*VirtualCabinet) (affected int64, err error)
	// GetVirtualCabinetByID 返回指定ID的虚拟货架
	GetVirtualCabinetByID(id uint) (*VirtualCabinet, error)
	// CountServerRooms 统计满足过滤条件的库房数量
	CountVirtualCabinets(cond *VirtualCabinet) (count int64, err error)
	// GetStoreRooms 返回满足过滤条件的库房
	GetVirtualCabinets(cond *VirtualCabinet, orderby OrderBy, limiter *page.Limiter) (items []*VirtualCabinet, err error)
	// GetStoreRooms 返回满足过滤条件的库房
	GetVirtualCabinetsByCond(cond *CombinedStoreRoomVirtualCabinet, orderby OrderBy, limiter *page.Limiter) (items []*VirtualCabinet, err error)	
}

// CombinedStoreRoomVirtualCabinet
type CombinedStoreRoomVirtualCabinet struct {
	StoreRoomName           string
	StoreRoomID 			uint
	VirtualCabinetNumber    string
}