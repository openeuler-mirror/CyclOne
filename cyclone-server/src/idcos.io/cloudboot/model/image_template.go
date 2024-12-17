package model

import (
	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
)

// ImageTemplate 镜像安装模板
type ImageTemplate struct {
	gorm.Model
	Family          string `gorm:"column:family"`
	BootMode        string `gorm:"column:boot_mode"`
	Name            string `gorm:"column:name"`
	ImageURL        string `gorm:"column:image_url"`
	Username        string `gorm:"column:username"`
	Password        string `gorm:"column:password"`
	Partition       string `gorm:"column:partition"`
	PreScript       string `gorm:"column:pre_script"`
	PostScript      string `gorm:"column:post_script"`
	OSLifecycle     string `gorm:"column:os_lifecycle"` // OS生命周期：Testing|Active(Default)|Active|Containment|EOL
	Arch            string `gorm:"column:arch"`  //  OS架构平台：x86_64|aarch64
}

// TableName 指定数据库表名
func (ImageTemplate) TableName() string {
	return "image_template"
}

// IImageTemplate 镜像安装模板持久化接口
type IImageTemplate interface {
	// RemoveImageTemplate 删除指定ID的镜像安装模板
	RemoveImageTemplate(id uint) (affected int64, err error)
	// SaveImageTemplate 保存镜像安装模板
	SaveImageTemplate(*ImageTemplate) (id uint, err error)
	// GetImageTemplateByID 返回指定ID的镜像安装模板
	GetImageTemplateByID(id uint) (*ImageTemplate, error)
	// CountImageTemplates 统计满足过滤条件的镜像安装模板数量
	CountImageTemplates(cond *ImageTemplate) (count int64, err error)
	// GetImageTemplates 返回满足过滤条件的镜像安装模板
	GetImageTemplates(cond *ImageTemplate, orderby OrderBy, limiter *page.Limiter) (items []*ImageTemplate, err error)
	// 查询指定设备关联的镜像安装模板
	GetImageTemplateBySN(sn string) (*ImageTemplate, error)

	// ========== deprecated methods ==========
	CountImageTemplateByName(name string) (uint, error)
	CountImageTemplateByNameAndID(name string, ID uint) (uint, error)
	GetImageTemplateIDByName(name string) (uint, error)
	CountImageTemplate() (uint, error)
	GetImageTemplateListWithPage(limit uint, offset uint) ([]ImageTemplate, error)
	GetImageTemplateByName(name string) (*ImageTemplate, error)
}
