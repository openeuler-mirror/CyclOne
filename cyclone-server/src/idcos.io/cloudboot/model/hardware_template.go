package model

import (
	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
)

// HardwareTplCond 硬件模板查询条件
type HardwareTplCond struct {
	Builtin   string
	Name      string
	Vendor    string
	ModelName string
}

// HardwareTemplate 硬件配置模板
type HardwareTemplate struct {
	gorm.Model
	Name      string `gorm:"column:name"`
	Vendor    string `gorm:"column:vendor"`
	ModelName string `gorm:"column:model"`
	Builtin   string `gorm:"column:builtin"`
	Data      string `gorm:"column:data"`
}

// TableName 指定数据库表名
func (HardwareTemplate) TableName() string {
	return "hardware_template"
}

// IHardwareTemplate 硬件配置操作接口
type IHardwareTemplate interface {
	// GetHardwareTemplateByName 返回指定名称的硬件配置模板。
	GetHardwareTemplateByName(name string) (template *HardwareTemplate, err error)
	// CountHardwareByCond 统计查询硬件模板数量
	CountHardwareByCond(cond *HardwareTplCond) (int64, error)
	// GetHardwaresByCond 分页查询硬件模板
	GetHardwaresByCond(cond *HardwareTplCond, limiter *page.Limiter) ([]*HardwareTemplate, error)
	// GetHardwareTemplateByID 返回指定ID的镜像安装模板
	GetHardwareTemplateByID(id uint) (*HardwareTemplate, error)
	// SaveHardwareTemplate 保存、修改硬件模板
	SaveHardwareTemplate(na *HardwareTemplate) (id uint, err error)
	//RemoveHardwareTemplateByID 删除硬件模板配置
	RemoveHardwareTemplateByID(id uint) (affected int64, err error)
}
