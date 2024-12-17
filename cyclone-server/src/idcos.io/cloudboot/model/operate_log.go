package model

import (
	"time"

	"github.com/voidint/page"
)

type (
	// OperateLog 操作记录
	OperateLog struct {
		ID           uint      `gorm:"column:id"  json:"id"`
		CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
		Operator     string    `gorm:"column:operator" json:"operator"`
		URL          string    `gorm:"column:url" json:"url"`
		HTTPMethod   string    `gorm:"column:http_method" json:"http_method"`
		CategoryCode string    `gorm:"column:category_code" json:"category_code"`
		CategoryName string    `gorm:"column:category_name" json:"category_name"`
		Source       string    `gorm:"column:source" json:"source"`
		Destination  string    `gorm:"column:destination" json:"destination"`
	}
)

// TableName 指定数据库表名
func (OperateLog) TableName() string {
	return "operation_log"
}

// IOperateLog 操作记录持久化接口
type IOperateLog interface {
	// SaveOperateLog 保存操作记录
	SaveOperateLog(*OperateLog) (id uint, err error)
	// CountOperateLog 统计满足过滤条件的操作记录数量
	CountOperateLog(cond *OperateLog) (count int64, err error)
	// GetOperateLogByCond 返回满足过滤条件的操作记录
	GetOperateLogByCond(cond *OperateLog, orderby OrderBy, limiter *page.Limiter) (items []*OperateLog, err error)
}
