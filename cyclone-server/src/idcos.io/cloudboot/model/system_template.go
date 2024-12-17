package model

import (
	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
)

const (
	// BootModeUEFI UEFI启动模式
	BootModeUEFI = "uefi"
	// BootModeBIOS Legacy BIOS启动模式
	BootModeBIOS = "legacy_bios"
)

// SystemTemplate 系统模板
type SystemTemplate struct {
	gorm.Model
	Family          string
	BootMode        string
	Name            string
	PXE             string `gorm:"column:pxe"`
	Username        string
	Password        string
	Content         string
	OSLifecycle     string `gorm:"column:os_lifecycle"` // OS生命周期：Testing|Active(Default)|Active|Containment|EOL
	Arch            string `gorm:"column:arch"`  //  OS架构平台：x86_64|aarch64
}

// TableName 指定数据库表名
func (SystemTemplate) TableName() string {
	return "system_template"
}

// ISystemTemplate 系统安装模板持久化接口
type ISystemTemplate interface {
	// RemoveSystemTemplate 删除指定ID的系统安装模板
	RemoveSystemTemplate(id uint) (affected int64, err error)
	// SaveSystemTemplate 保存系统安装模板
	SaveSystemTemplate(*SystemTemplate) (id uint, err error)
	// GetSystemTemplateByID 返回指定ID的系统安装模板
	GetSystemTemplateByID(id uint) (*SystemTemplate, error)
	// CountSystemTemplates 统计满足过滤条件的系统安装模板数量
	CountSystemTemplates(cond *SystemTemplate) (count int64, err error)
	// GetSystemTemplates 返回满足过滤条件的系统安装模板
	GetSystemTemplates(cond *SystemTemplate, orderby OrderBy, limiter *page.Limiter) (items []*SystemTemplate, err error)

	// ========== deprecated methods ==========
	CountSystemTemplate() (uint, error)
	CountSystemTemplateByName(name string) (uint, error)
	CountSystemTemplateByNameAndID(name string, id uint) (uint, error)
	CountSystemTemplateByShield(cond *SystemTemplate) (uint, error)
	GetSystemTemplateIDByName(name string) (uint, error)
	GetSystemTemplateListWithPage(Limit uint, Offset uint) ([]SystemTemplate, error)
	GetSystemTemplateListWithPageAndShield(Limit uint, Offset uint, cond *SystemTemplate) ([]SystemTemplate, error)
	GetSystemTemplateByName(n string) (*SystemTemplate, error)
}
