package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
)

// Order 订单
type Order struct {
	gorm.Model
	IDCID               uint      `gorm:"column:idc_id"`
	ServerRoomID        uint      `gorm:"column:server_room_id"`
	PhysicalArea        string    `gorm:"column:physical_area"`
	Number              string    `gorm:"column:number"`
	Usage               string    `gorm:"column:usage"`
	Category            string    `gorm:"column:category"`
	Amount              int       `gorm:"column:amount"`
	LeftAmount          int       `gorm:"column:left_amount"`
	ExpectedArrivalDate time.Time `gorm:"column:expected_arrival_date"`
	PreOccupiedCabinets string    `gorm:"column:pre_occupied_cabinets"`
	PreOccupiedUsites   string    `gorm:"column:pre_occupied_usites"`
	Remark              string    `gorm:"column:remark"`
	Status              string    `gorm:"column:status"`
	Creator             string    `gorm:"column:creator"`
	//以下字段参考 DeviceLifecycle
	AssetBelongs	 				string		`gorm:"column:asset_belongs"`
	Owner			 				string		`gorm:"column:owner"`
	IsRental		 				string		`gorm:"column:is_rental"`
	MaintenanceServiceProvider		string		`gorm:"column:maintenance_service_provider"`
	MaintenanceService				string		`gorm:"column:maintenance_service"`
	LogisticsService				string		`gorm:"column:logistics_service"`
	MaintenanceServiceDateBegin     time.Time 	`gorm:"column:maintenance_service_date_begin"`
	MaintenanceServiceDateEnd       time.Time 	`gorm:"column:maintenance_service_date_end"`
}

// OrderCond 订单查询条件
type OrderCond struct {
	Order
	ID []uint
}

// TableName 指定数据库表名
func (Order) TableName() string {
	return "order"
}

// IOrder 订单持久化接口
type IOrder interface {
	// RemoveOrderByID 删除指定ID的订单
	RemoveOrderByID(id uint) (affected int64, err error)
	// SaveOrder 保存订单
	SaveOrder(*Order) (affected int64, err error)
	// UpdateOrder更新订单信息
	UpdateOrder(*Order) (affected int64, err error)
	// GetOrderByID 返回指定ID的订单
	GetOrderByID(id uint) (*Order, error)
	//GetOrderByNumber 返回指定编号的订单
	GetOrderByNumber(n string) (*Order, error)
	// CountServerRooms 统计满足过滤条件的订单数量
	CountOrders(cond *OrderCond) (count int64, err error)
	// GetOrders 返回满足过滤条件的订单
	GetOrders(cond *OrderCond, orderby OrderBy, limiter *page.Limiter) (items []*Order, err error)
	// GetMaxOrderNumber 获取指定日期的最大订单号
	GetMaxOrderNumber(date string) (orderNumber int, err error)
}

const (
	OrderStatusPurchasing    = "purchasing"
	OrderStatusPartlyArrived = "partly_arrived"
	OrderStatusAllArrived    = "all_arrived"
	OrderStatusCanceled      = "canceled"
	OrderStatusfinished      = "finished"
)

// BeforeSave 保存设备信息前的钩子方法。
// 防止将空字符串写入类型为JSON的数据库字段中引发报错。
func (o *Order) BeforeSave() (err error) {
	replaceIfBlank(&o.PreOccupiedUsites, EmptyJSONObject)
	replaceIfBlank(&o.PreOccupiedCabinets, EmptyJSONObject)
	return nil
}
