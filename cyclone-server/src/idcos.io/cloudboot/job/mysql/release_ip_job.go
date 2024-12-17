package mysql

import (
	"idcos.io/cloudboot/config"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
)

// ReleaseIPJob 释放IP任务
type ReleaseIPJob struct {
	id   string // 任务ID(全局唯一)
	log  logger.Logger
	repo model.Repo
	conf *config.Config
}

// NewReleaseIPJob 实例化任务管理器
func NewReleaseIPJob(log logger.Logger, repo model.Repo, conf *config.Config, jobid string) *ReleaseIPJob {
	return &ReleaseIPJob{
		log:  log,
		repo: repo,
		conf: conf,
		id:   jobid,
	}
}

// Run 运行IP释放任务
// 内置任务，根据release_date 将 disabled的IP 置为空闲
func (j *ReleaseIPJob) Run() {
	j.log.Debugf("Start checking the IP release date")
	defer j.log.Debugf("IP release date checking is completed")

	defer func() {
		if err := recover(); err != nil {
			j.log.Errorf("IP release job panic: \n%s", err)
		}
	}()

	items, _ := j.repo.GetReleasableIP()
	if len(items) <= 0 {
		j.log.Debugf("The releasable IP list is empty")
		return
	}

	for i := range items {
		items[i].IsUsed = model.IPNotUsed
		_, err := j.repo.SaveIP(items[i])
		if err != nil {
			j.log.Errorf("IP %s  release fail %s .", items[i].IP, err)
			continue
		}
		j.log.Infof("IP %s has been released .", items[i].IP)
	}
}
