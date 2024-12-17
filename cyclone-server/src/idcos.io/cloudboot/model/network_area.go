package model

import (
	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
)

const (
	// NetworkAreaStatNonProduction 网络区域状态-未投产
	NetworkAreaStatNonProduction = "nonproduction"
	// NetworkAreaStatProduction 网络区域状态-已投产
	NetworkAreaStatProduction = "production"
	// NetworkAreaStatOffline 网络区域状态-已经下线(回收)
	NetworkAreaStatOffline = "offline"
)

// NetworkArea 网络区域
type NetworkArea struct {
	gorm.Model
	IDCID        uint   `gorm:"column:idc_id"`
	ServerRoomID uint   `gorm:"column:server_room_id"`
	Name         string `gorm:"column:name"`
	PhysicalArea string `gorm:"column:physical_area"`
	Status       string `gorm:"column:status"`
	Creator      string `gorm:"column:creator"`
}

// NetworkArea 网络区域搜索条件
type NetworkAreaCond struct {
	IDCID          []uint `gorm:"column:idc_id"`
	ServerRoomID   []uint `gorm:"column:server_room_id"`
	Name           string `gorm:"column:name"`
	PhysicalArea   string `gorm:"column:physical_area"`
	Status         string `gorm:"column:status"`
	Creator        string `gorm:"column:creator"`
	ServerRoomName string
}

// TableName 指定数据库表名
func (NetworkArea) TableName() string {
	return "network_area"
}

// INetworkArea 网络区域持久化接口
type INetworkArea interface {
	// RemoveNetworkAreaByID 删除指定ID的网络区域
	RemoveNetworkAreaByID(id uint) (affected int64, err error)
	// SaveNetworkArea 保存网络区域
	SaveNetworkArea(*NetworkArea) (affected int64, err error)
	// UpdateNetworkAreaStatus 批量更新网络区域状态
	UpdateNetworkAreaStatus(status string, ids ...uint) (affected int64, err error)
	// GetNetworkAreaByID 返回指定ID的网络区域
	GetNetworkAreaByID(id uint) (*NetworkArea, error)
	// GetNetworkAreasByCond 返回满足过滤条件的网络区域(不支持模糊查找)
	GetNetworkAreasByCond(cond *NetworkArea) (item []*NetworkArea, err error)
	// CountNetworkAreas 统计满足过滤条件的网络区域数量
	CountNetworkAreas(cond *NetworkAreaCond) (count int64, err error)
	// GetNetworkAreas 返回满足过滤条件的网络区域
	GetNetworkAreas(cond *NetworkAreaCond, orderby OrderBy, limiter *page.Limiter) (items []*NetworkArea, err error)
}
