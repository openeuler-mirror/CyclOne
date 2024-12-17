package model

import (
	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
)

const (
	// RoomStatUnderConstruction 机房状态-建设中
	RoomStatUnderConstruction = "under_construction"
	// RoomStatAccepted 机房状态-已验收
	RoomStatAccepted = "accepted"
	// RoomStatProduction 机房状态-已投产
	RoomStatProduction = "production"
	// RoomStatAbolished 机房状态-已裁撤
	RoomStatAbolished = "abolished"
)

// ServerRoom 机房
type ServerRoom struct {
	gorm.Model
	IDCID               uint   `gorm:"column:idc_id"`
	Name                string `gorm:"column:name"`
	FirstServerRoom     uint   `gorm:"column:first_server_room"`
	City                string `gorm:"column:city"`
	Address             string `gorm:"column:address"`
	ServerRoomManager   string `gorm:"column:server_room_manager"`
	VendorManager       string `gorm:"column:vendor_manager"`
	NetworkAssetManager string `gorm:"column:network_asset_manager"`
	SupportPhoneNumber  string `gorm:"column:support_phone_number"`
	Status              string `gorm:"column:status"`
	Creator             string `gorm:"column:creator"`
}

// ServerRoomCond 机房查询条件
type ServerRoomCond struct {
	ID                  []uint
	IDCID               []uint
	Name                string
	FirstServerRoom     []uint
	City                string
	Address             string
	ServerRoomManager   string
	VendorManager       string
	NetworkAssetManager string
	SupportPhoneNumber  string
	Status              string
	Creator             string
}

// TableName 指定数据库表名
func (ServerRoom) TableName() string {
	return "server_room"
}

// IServerRoom 机房持久化接口
type IServerRoom interface {
	// RemoveServerRoomByID 删除指定ID的机房
	RemoveServerRoomByID(id uint) (affected int64, err error)
	// SaveServerRoom 保存机房
	SaveServerRoom(*ServerRoom) (affected int64, err error)
	// UpdateServerRoom更新机房信息
	UpdateServerRoom(srs []*ServerRoom) (affected int64, err error)
	// UpdateServerRoomStatus 批量更新机房状态
	UpdateServerRoomStatus(status string, ids ...uint) (affected int64, err error)
	// GetServerRoomByID 返回指定ID的机房
	GetServerRoomByID(id uint) (*ServerRoom, error)
	// GetServerRoomByName 返回指Name的机房
	GetServerRoomByName(Name string) (*ServerRoom, error)
	// CountServerRooms 统计满足过滤条件的机房数量
	CountServerRooms(cond *ServerRoomCond) (count int64, err error)
	// GetServerRooms 返回满足过滤条件的机房
	GetServerRooms(cond *ServerRoomCond, orderby OrderBy, limiter *page.Limiter) (items []*ServerRoom, err error)
}
