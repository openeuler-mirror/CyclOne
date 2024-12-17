package model

import (
	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
)

// DeviceCategory 设备类型
type DeviceCategory struct {
	gorm.Model
	Category 						string 	`gorm:"column:category"`
	Hardware 						string 	`gorm:"column:hardware"` //硬件配置
	CentralProcessorManufacturer	string  `gorm:"column:central_processor_manufacture"` // 处理器生产商
	CentralProcessorArch			string  `gorm:"column:central_processor_arch"` // 处理器架构
	Power    						string 	`gorm:"column:power"`    //功率
	Unit 							uint	`gorm:"column:unit"`    //设备所占用的U数 机柜参数Unit 1U = 44.45mm
	IsFITIEcoProduct				string	`gorm:"column:is_fiti_eco_product"`	// Financial Information Technology Innovation Ecological Product 金融信创生态产品
	Remark   						string 	`gorm:"column:remark"`
	Creator  						string 	`gorm:"column:creator"`
}

// TableName 指定数据库表名
func (DeviceCategory) TableName() string {
	return "device_category"
}

// IDeviceCategory 设备类型持久化接口
type IDeviceCategory interface {
	// RemoveDeviceCategoryByID 删除指定ID的设备类型
	RemoveDeviceCategoryByID(id uint) (affected int64, err error)
	// SaveDeviceCategory 保存设备类型
	SaveDeviceCategory(*DeviceCategory) (affected int64, err error)
	// GetDeviceCategoryByID 返回指定ID的设备类型
	GetDeviceCategoryByID(id uint) (*DeviceCategory, error)
	// CountServerRooms 统计满足过滤条件的设备类型数量
	CountDeviceCategorys(cond *DeviceCategory) (count int64, err error)
	// GetDeviceCategorys 返回满足过滤条件的设备类型
	GetDeviceCategorys(cond *DeviceCategory, orderby OrderBy, limiter *page.Limiter) (items []*DeviceCategory, err error)
	GetDeviceCategoryByName(name string) (dc *DeviceCategory, err error)
	GetDeviceCategoryQuerys(param string) (*DeviceQueryParamResp, error)
}
