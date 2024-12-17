package model

import (
	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
)

// PermissionCode 权限码
type PermissionCode struct {
	gorm.Model
	PID   uint   `gorm:"column:pid"`
	Code  string `gorm:"column:code"`
	Title string `gorm:"column:title"`
	Note  string `gorm:"column:note"`
}

// TableName 指定数据库表名
func (PermissionCode) TableName() string {
	return "permission_code"
}

// IPermissionCode 权限码数据库操作接口
type IPermissionCode interface {
	// GetPermissionCodes 返回满足过滤条件的权限码
	GetPermissionCodes(cond *PermissionCode, orderby OrderBy, limiter *page.Limiter) (items []*PermissionCode, err error)
}
