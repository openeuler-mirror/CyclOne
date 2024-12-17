package model

import (
	"time"

	"github.com/voidint/page"
)

type (
	// APILog 操作记录
	APILog struct {
		ID          uint      `gorm:"column:id"  json:"id"`
		CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
		Operator    string    `gorm:"column:operator" json:"operator"`
		Description string    `gorm:"column:description" json:"description"`
		API         string    `gorm:"column:api" json:"api"`
		ReqBody     string    `gorm:"column:req_body" json:"req_body"`
		Method      string    `gorm:"column:method" json:"method"`
		Status      string    `gorm:"column:status" json:"status"`
		RemoteAddr  string    `gorm:"column:remote_addr" json:"remote_addr"`
		Msg         string    `gorm:"column:msg" json:"msg"`
		Result      string    `gorm:"column:result" json:"result"`
		Time        float64   `gorm:"column:time" json:"time"`
	}

	// APILogCond 操作记录查询条件
	APILogCond struct {
		CreatedAtStart time.Time `json:"created_at_start"`
		CreatedAtEnd   time.Time `json:"created_at_end"`
		Operator       string    `json:"operator"`
		Description    string    `json:"desc"`
		API            string    `json:"api"`
		Method         string    `json:"method"`
		CategoryName   string    `json:"category_name"`
		Status         string    `json:"status"`
		Cost1          float64   `json:"cost1"`
		Cost2          float64   `json:"cost2"`
	}
)

// TableName 指定数据库表名
func (APILog) TableName() string {
	return "api_log"
}

// IAPILog 操作记录持久化接口
type IAPILog interface {
	// SaveAPILog 保存操作记录
	SaveAPILog(*APILog) (id uint, err error)
	// CountAPILog 统计满足过滤条件的操作记录数量
	CountAPILog(cond *APILogCond) (count int64, err error)
	// GetAPILogByCond 返回满足过滤条件的操作记录
	GetAPILogByCond(cond *APILogCond, orderby OrderBy, limiter *page.Limiter) (items []*APILog, err error)
}
