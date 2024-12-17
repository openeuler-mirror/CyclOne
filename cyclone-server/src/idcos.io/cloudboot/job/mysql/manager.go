package mysql

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/jakecoffman/cron"
	"github.com/jinzhu/gorm"
	"github.com/voidint/page"

	"idcos.io/cloudboot/config"
	"idcos.io/cloudboot/job"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
)

// JobManager MySQL实现的任务管理器
type JobManager struct {
	log    logger.Logger
	repo   model.Repo
	conf   *config.Config
	runner *cron.Cron
	once   sync.Once
}

// NewJobManager 实例化任务管理器
func NewJobManager(log logger.Logger, repo model.Repo, conf *config.Config) job.Manager {
	runner := cron.New()
	if runner != nil {
		runner.Start()
		return &JobManager{
			log:    log,
			repo:   repo,
			conf:   conf,
			runner: runner,
		}
	} else {
		log.Errorf("Failed to new cron.Cron pointer.")
		return nil
	}
}

// Rebuild 重建/恢复任务（用于系统重启场景）
func (mgr *JobManager) Rebuild() error {
	mgr.log.Info("Start rebuilding jobs")
	defer mgr.log.Info("End job rebuilding")

	cond := model.Job{
		Rate: job.RateFixedRate, // 仅重建周期性任务
	}
	items, err := mgr.repo.GetJobs(&cond, nil, nil)
	if err != nil {
		return err
	}

	for i := range items {
		if items[i] == nil || items[i].Cron == "" || items[i].Status != job.Running {
			continue
		}
		mgr.log.Infof("Job(%s) begins to rebuild", items[i].ID)

		var target map[string][]string
		if items[i].Target != "" {
			if err = json.Unmarshal([]byte(items[i].Target), &target); err != nil {
				mgr.log.Error(err)
				continue
			}
		}
		newJob := job.Job{
			ID:             items[i].ID,
			Builtin:        items[i].Builtin,
			Title:          items[i].Title,
			Category:       items[i].Category,
			Rate:           items[i].Rate,
			CronExpression: items[i].Cron,
			CronRender:     items[i].CronRender,
			Creator:        items[i].Creator,
			Target:         target,
			Status:         job.Running,
		}

		switch items[i].Category {
		case job.CategoryInspection: // 重建硬件巡检任务
			newJob.CronJob = NewInspectionJob(mgr.log, mgr.repo, mgr.conf, items[i].ID)
		case job.CategoryInstallationTime: // 重建安装超时处理任务
			newJob.CronJob = NewInstallationTimeoutJob(mgr.log, mgr.repo, mgr.conf, items[i].ID)
		case job.CategoryReleaseIP: // 重建释放IP任务
			newJob.CronJob = NewReleaseIPJob(mgr.log, mgr.repo, mgr.conf, items[i].ID)
		case job.CategoryAutoDeploy: // 重建自动部署任务
			newJob.CronJob = NewAutoDeployJob(mgr.log, mgr.repo, mgr.conf, items[i].ID)
		case job.CategoryMailSend: // 重建邮件推送任务
			newJob.CronJob = NewMailSendJob(mgr.log, mgr.repo, mgr.conf, items[i].ID)
		case job.CategoryUpdateDeviceLifecycle:
			newJob.CronJob = NewUpdateDeviceLifecycleJob(mgr.log, mgr.repo, mgr.conf, items[i].ID)
		}

		if err = mgr.Submit(&newJob); err != nil {
			mgr.log.Warnf("Job(%s) rebuild failed", items[i].ID)
			continue
		}
		mgr.log.Infof("Job(%s) rebuilt successfully", items[i].ID)
	}
	return nil
}

// Stop 停止任务管理器
func (mgr *JobManager) Stop() error {
	mgr.once.Do(func() {
		mgr.log.Infof("Stop job manager")
		mgr.runner.Stop()
	})
	return nil
}

// Submit 向任务管理器提交任务
func (mgr *JobManager) Submit(cjob *job.Job) (err error) {
	mgr.log.Infof("Submit job(%s)", cjob.ID)

	// 1、将任务写入job表进行持久化
	newJob := model.Job{
		ID:         cjob.ID,
		Title:      cjob.Title,
		Builtin:    cjob.Builtin,
		Category:   cjob.Category,
		Rate:       cjob.Rate,
		Cron:       cjob.CronExpression,
		CronRender: cjob.CronRender,
		Status:     job.Running, // 提交后任务状态改为'运行中'
		Creator:    cjob.Creator,
	}

	if cjob.Rate == job.RateImmediately {
		now := time.Now()
		newJob.NextRunTime = &now
	}
	if btarget, _ := json.Marshal(cjob.Target); btarget != nil {
		newJob.Target = string(btarget)
	}

	if err = mgr.repo.SaveJob(&newJob); err != nil {
		return err
	}

	// 2、将任务提交至任务执行器进行调度执行
	switch cjob.Rate {
	case job.RateImmediately:
		go cjob.CronJob.Run() // TODO 控制协程数量及执行超时时间

	case job.RateFixedRate:
		mgr.runner.AddJob(cjob.CronExpression, cjob.CronJob, cjob.ID)
	}
	return nil
}

// Remove 移除指定ID的任务。若该任务已经开始运行，则此操作不会有任何效果。
func (mgr *JobManager) Remove(jobid string) (err error) {
	mgr.log.Infof("Remove job(id=%s)", jobid)

	// 1、校验任务ID有效性
	mjob, err := mgr.validJob(jobid)
	if err != nil {
		return err
	}

	//  内置任务不能删除
	if mjob.Builtin == job.Builtin {
		return job.ErrRemoveBuiltinJob
	}

	// 2、将指定ID的任务标记为'已删除'
	if err = mgr.repo.SaveJob(&model.Job{
		ID:     jobid,
		Status: job.Deleted,
	}); err != nil {
		return err
	}

	// 3、从任务执行器中移除指定ID的任务
	mgr.runner.RemoveJob(jobid)
	return nil
}

// Pause 暂停指定ID的任务
func (mgr *JobManager) Pause(jobid string) error {
	mgr.log.Infof("Pause job(id=%s)", jobid)

	// 1、校验任务ID有效性
	mjob, err := mgr.validJob(jobid)
	if err != nil {
		return err
	}

	// 非运行中的定时任务都无法暂停
	if mjob.Rate != job.RateFixedRate || mjob.Status != job.Running {
		return job.ErrPauseJob
	}

	// 2、将指定ID的任务标记为'暂停'
	if err = mgr.repo.SaveJob(&model.Job{
		ID:     jobid,
		Status: job.Paused,
	}); err != nil {
		return err
	}

	// 3、从任务执行器中移除指定ID的任务
	mgr.runner.RemoveJob(jobid)
	return nil
}

// Unpause 继续指定ID的任务
func (mgr *JobManager) Unpause(jobid string) error {
	mgr.log.Infof("Unpause job(id=%s)", jobid)

	// 1、校验任务ID有效性
	mjob, err := mgr.validJob(jobid)
	if err != nil {
		return err
	}

	// 非暂停的定时任务都无法继续
	if mjob.Rate != job.RateFixedRate || mjob.Status != job.Paused {
		return job.ErrUnpauseJob
	}

	// 2、将指定ID的任务标记为'运行中'
	if err = mgr.repo.SaveJob(&model.Job{
		ID:     jobid,
		Status: job.Running,
	}); err != nil {
		return err
	}

	// 3、将任务重新提交至任务执行器进行调度执行
	var newjob cron.Job
	switch mjob.Category {
	case job.CategoryInspection: // 重启硬件巡检任务
		newjob = NewInspectionJob(mgr.log, mgr.repo, mgr.conf, mjob.ID)
	case job.CategoryInstallationTime: // 重启安装超时处理任务
		newjob = NewInstallationTimeoutJob(mgr.log, mgr.repo, mgr.conf, mjob.ID)
	case job.CategoryReleaseIP: // 重启释放IP任务
		newjob = NewReleaseIPJob(mgr.log, mgr.repo, mgr.conf, mjob.ID)
	case job.CategoryAutoDeploy: // 重启自动部署任务
		newjob = NewAutoDeployJob(mgr.log, mgr.repo, mgr.conf, mjob.ID)
	case job.CategoryMailSend: // 重启邮件推送任务
		newjob = NewMailSendJob(mgr.log, mgr.repo, mgr.conf, mjob.ID)
	case job.CategoryUpdateDeviceLifecycle:
		newjob = NewUpdateDeviceLifecycleJob(mgr.log, mgr.repo, mgr.conf, mjob.ID)
	}

	mgr.runner.AddJob(mjob.Cron, newjob, mjob.ID)
	return nil
}

// GetJobByID 返回指定ID的任务的一份拷贝
func (mgr *JobManager) GetJobByID(jobid string) (*job.Job, error) {
	mjob, err := mgr.validJob(jobid)
	if err != nil {
		return nil, err
	}
	var target map[string][]string
	if mjob.Target != "" {
		if err = json.Unmarshal([]byte(mjob.Target), &target); err != nil {
			mgr.log.Error(err)
		}
	}
	return &job.Job{
		ID:             mjob.ID,
		Creator:        mjob.Creator,
		Builtin:        mjob.Builtin,
		Title:          mjob.Title,
		Category:       mjob.Category,
		Rate:           mjob.Rate,
		CronExpression: mjob.Cron,
		CronRender:     mjob.CronRender,
		Target:         target,
		Status:         mjob.Status,
		CreatedAt:      mjob.CreatedAt,
		UpdatedAt:      mjob.UpdatedAt,
	}, nil
}

// GetJobs 返回指定分页的任务列表
func (mgr *JobManager) GetJobs(cond *job.Job, pageNo, pageSize int64) (totalRecords int64, jobs []*job.Job, err error) {
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	mcond := model.Job{
		Builtin:  cond.Builtin,
		Title:    cond.Title,
		Category: cond.Category,
		Rate:     cond.Rate,
		Status:   cond.Status,
		Creator:  cond.Creator,
	}
	totalRecords, err = mgr.repo.CountJobs(&mcond)
	if err != nil {
		return 0, nil, err
	}

	items, err := mgr.repo.GetJobs(&mcond, model.OneOrderBy("updated_at", model.DESC), &page.Limiter{
		Limit:  pageSize,
		Offset: (pageNo - 1) * pageSize,
	})
	if err != nil {
		return 0, nil, err
	}
	for _, item := range items {
		var target map[string][]string
		if item.Target != "" {
			if err = json.Unmarshal([]byte(item.Target), &target); err != nil {
				mgr.log.Error(err)
			}
		}
		jobs = append(jobs, &job.Job{
			ID:             item.ID,
			Creator:        item.Creator,
			Builtin:        item.Builtin,
			Title:          item.Title,
			Category:       item.Category,
			Rate:           item.Rate,
			CronExpression: item.Cron,
			CronRender:     item.CronRender,
			Status:         item.Status,
			Target:         target,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		})
	}
	return totalRecords, jobs, nil
}

// validJob 校验任务的有效性，若指定ID是有效的任务ID，则返回该任务对象。
func (mgr *JobManager) validJob(jobid string) (mjob *model.Job, err error) {
	mjob, err = mgr.repo.GetJobByID(jobid)
	if gorm.IsRecordNotFoundError(err) {
		return nil, job.ErrInvalidJobID
	}
	return mjob, err
}
