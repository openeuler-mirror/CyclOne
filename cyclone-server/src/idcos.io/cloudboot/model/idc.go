package model

import (
	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
)

const (
	// IDCUsageProduction 数据中心用途-生产
	IDCUsageProduction = "production"
	// IDCUsageDisasterRecovery 数据中心用途-容灾
	IDCUsageDisasterRecovery = "disaster_recovery"
	// IDCUsagePreProduction 数据中心用途-准生产
	IDCUsagePreProduction = "pre_production"
	// IDCUsageTesting 数据中心用途-测试
	IDCUsageTesting = "testing"
)

const (
	// IDCStatUnderConstruction 数据中心状态-建设中
	IDCStatUnderConstruction = "under_construction"
	// IDCStatAccepted 数据中心状态-已验收
	IDCStatAccepted = "accepted"
	// IDCStatProduction 数据中心状态-已投产
	IDCStatProduction = "production"
	// IDCStatAbolished 数据中心状态-已裁撤
	IDCStatAbolished = "abolished"
)

type (
	// IDC 数据中心
	IDC struct {
		gorm.Model
		Name            string `gorm:"column:name"`
		Usage           string `gorm:"column:usage"`
		FirstServerRoom string `gorm:"column:first_server_room"`
		Vendor          string `gorm:"column:vendor"`
		Status          string `gorm:"column:status"`
		Creator         string `gorm:"column:creator"`
	}
)

// TableName 指定数据库表名
func (IDC) TableName() string {
	return "idc"
}

// IIDC 数据中心持久化接口
type IIDC interface {
	// RemoveIDCByID 删除指定ID的数据中心
	RemoveIDCByID(id uint) (affected int64, err error)
	// SaveIDC 保存数据中心
	SaveIDC(*IDC) (affected int64, err error)
	// AddIDC 新增数据中心
	AddIDC(*IDC) (affected int64, idc *IDC, err error)
	// UpdateIDC 修改数据中心
	UpdateIDC(*IDC) (affected int64, err error)
	// UpdateIDCStatus 批量更新数据中心状态
	UpdateIDCStatus(status string, ids ...uint) (affected int64, err error)
	// GetIDCByID 返回指定ID的数据中心
	GetIDCByID(id uint) (*IDC, error)
	// GetIDCByName 返回指定Name的数据中心
	GetIDCByName(name string) (*IDC, error)
	// CountIDCs 统计满足过滤条件的数据中心数量
	CountIDCs(cond *IDC) (count int64, err error)
	// GetIDCs 返回满足过滤条件的数据中心
	GetIDCs(cond *IDC, orderby OrderBy, limiter *page.Limiter) (items []*IDC, err error)
}
