package model

import (
	"time"

	"github.com/voidint/page"
)

// Job 任务
type Job struct {
	ID          string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
	Builtin     string     `gorm:"column:builtin"`
	Title       string     `gorm:"column:title"`
	Category    string     `gorm:"column:category"`
	Rate        string     `gorm:"column:rate"`
	Cron        string     `gorm:"column:cron"`
	CronRender  string     `gorm:"column:cron_render"`
	NextRunTime *time.Time `gorm:"column:next_run_time"`
	Target      string     `gorm:"column:target"`
	Status      string     `gorm:"column:status"`
	Creator     string     `gorm:"column:creator"`
}

// TableName 指定数据库表名
func (Job) TableName() string {
	return "job"
}

// IJob 任务操作接口
type IJob interface {
	// SaveJob 保存(新增/更新)任务
	SaveJob(*Job) error
	// RemoveJob 移除指定ID的任务
	RemoveJob(id string) (affected int64, err error)
	// GetJobByID 查询指定ID的任务
	GetJobByID(id string) (job *Job, err error)
	// CountJobs 统计满足条件的任务数量
	CountJobs(cond *Job) (count int64, err error)
	// GetJobs 返回满足过滤条件的任务
	GetJobs(cond *Job, orderby OrderBy, limiter *page.Limiter) (items []*Job, err error)
}
