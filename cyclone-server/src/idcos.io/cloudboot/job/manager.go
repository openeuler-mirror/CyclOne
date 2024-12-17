package job

import (
	"errors"
	"time"

	"github.com/jakecoffman/cron"
)

const (
	// CategoryInspection 任务类别-硬件巡检
	CategoryInspection = "inspection"
	// CategoryInstallationTime 任务类别-安装超时
	CategoryInstallationTime = "installation_timeout"
	// CategoryReleaseIP 任务类别-释放IP
	CategoryReleaseIP = "release_ip"
	// CategoryAutoDeploy 任务类别-自动部署
	CategoryAutoDeploy = "auto_deploy"
	// CategoryMailSend 任务类别-邮件推送
	CategoryMailSend = "mailsend"
	// CategoryUpdateDeviceLifecycle 任务类别-更新设备生命周期（维保状态）
	CategoryUpdateDeviceLifecycle = "update_device_lifecycle"	
)

const (
	// RateImmediately 任务执行频率-立即执行
	RateImmediately = "immediately"
	// RateFixedRate 任务执行频率-固定频率(周期性)执行
	RateFixedRate = "fixed_rate"
)

const (
	// Builtin 内建任务
	Builtin = "yes"
	// NoBuiltin 非内建任务
	NoBuiltin = "no"
)

const (
	// Running 任务状态-运行中
	Running = "running"
	// Paused 任务状态-已暂停
	Paused = "paused"
	// Stoped 任务状态-已停止
	Stoped = "stoped"
	// Deleted 任务状态-已删除
	Deleted = "deleted"
)

// Job 任务
type Job struct {
	ID             string              // 任务ID（保证全局唯一）
	Creator        string              // 任务创建者
	Builtin        string              // 是否是内建任务
	Title          string              // 任务标题
	Category       string              // 任务类型。可选值:inspection、installation_timeout、release_ip
	Rate           string              // 任务执行频率。可选值:immediately-立刻执行; fixed_rate-固定频率(周期性)执行;
	CronExpression string              // cron表达式。若为一次性任务，则该值为空。
	CronRender     string              // cron表达式UI渲染信息
	Target         map[string][]string // 任务作用目标。map中的value暂定为[]string，可能发生变化。
	Status         string              // 任务状态。
	CronJob        cron.Job            // 任务的具体实现步骤
	CreatedAt      time.Time           // 创建时间
	UpdatedAt      time.Time           // 更新时间
}

var (
	// ErrInvalidJobID 无效的任务ID
	ErrInvalidJobID = errors.New("invalid job id")
	// ErrPauseJob 无法暂停此任务(非运行中的定时任务无法暂停)
	ErrPauseJob = errors.New("unable to pause this job")
	// ErrUnpauseJob 无法继续此任务(非暂停的定时任务无法继续)
	ErrUnpauseJob = errors.New("unable to continue this job")
	// ErrRemoveBuiltinJob 无法删除系统内建任务
	ErrRemoveBuiltinJob = errors.New("unable to remove builtin job")
)

// Manager 任务管理器
type Manager interface {
	// Rebuild 重建/恢复任务（用于系统重启场景）
	Rebuild() error
	// Submit 向任务管理器提交任务
	Submit(cjob *Job) error
	// Remove 移除指定ID的任务。
	// 若指定的任务ID无效，则返回ErrInvalidJobID错误。
	// 若该任务已经开始运行，则此操作不会有任何效果。
	Remove(jobid string) error
	// Pause 暂停指定ID的任务。
	// 若指定的任务ID无效，则返回ErrInvalidJobID错误。
	// 若该任务无法暂停，则返回ErrPauseJob错误。
	Pause(jobid string) error
	// Unpause 继续运行指定ID的任务。
	// 若指定的任务ID无效，则返回ErrInvalidJobID错误。
	// 若该任务无法继续，则返回ErrUnpauseJob错误。
	Unpause(jobid string) error
	// GetJobByID 返回指定ID的任务的一份拷贝
	GetJobByID(jobid string) (*Job, error)
	// GetJobs 返回指定分页的任务列表
	GetJobs(cond *Job, pageNo, pageSize int64) (totalRecords int64, items []*Job, err error)
	// Stop 停止任务管理器
	Stop() error
}
