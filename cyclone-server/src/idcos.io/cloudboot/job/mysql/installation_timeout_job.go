package mysql

import (
	"fmt"
	"time"

	"idcos.io/cloudboot/config"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils/centos6"
	"idcos.io/cloudboot/utils/times"
)

// InstallationTimeoutJob 检查安装超时任务
type InstallationTimeoutJob struct {
	id   string // 任务ID(全局唯一)
	log  logger.Logger
	repo model.Repo
	conf *config.Config
}

// NewInstallationTimeoutJob 实例化任务管理器
func NewInstallationTimeoutJob(log logger.Logger, repo model.Repo, conf *config.Config, jobid string) *InstallationTimeoutJob {
	return &InstallationTimeoutJob{
		log:  log,
		repo: repo,
		conf: conf,
		id:   jobid,
	}
}

// Run 运行安装超时处理任务
func (j *InstallationTimeoutJob) Run() {
	j.log.Debugf("Start checking the installation timeout")
	defer j.log.Debugf("Installation timeout checking is completed")

	defer func() {
		if err := recover(); err != nil {
			j.log.Errorf("Installation timeout job panic: \n%s", err)
		}
	}()

	timeout := j.repo.GetSystemSetting4InstallatonTimeout(model.DefInstallationTimeout)
	items, _ := j.repo.GetDeviceSettingsByInstallationTimeout(timeout)
	if len(items) <= 0 {
		j.log.Debugf("The os installation timeout device list is empty")
		return
	}

	now := time.Now().Format(times.DateTimeLayout)

	for i := range items {
		affected, _ := j.repo.SetInstallationTimeout(items[i].SN)
		if affected <= 0 {
			continue
		}

		if centos6.IsPXEUEFI(j.log, j.repo, items[i].SN) {
			_ = centos6.DropConfigurations(j.log, j.repo, items[i].SN) // TODO 为支持centos6的UEFI方式安装而临时增加的逻辑，后续会删除。
		}

		_, _ = j.repo.UpdateDeviceBySN(&model.Device{
			SN:              items[i].SN,
			OperationStatus: model.DevOperStatPreDeploy,
		})

		var startTime string
		if items[i].InstallationStartTime != nil {
			startTime = items[i].InstallationStartTime.Format(times.DateTimeLayout)
		}

		_, _ = j.repo.SaveDeviceLog(&model.DeviceLog{
			SN:              items[i].SN,
			DeviceSettingID: items[i].ID,
			LogType:         model.DeviceLogInstallType,
			Title:           "安装失败(超时)",
			Content:         fmt.Sprintf("安装始于%s，至%s超时。", startTime, now),
		})
		j.log.Infof("The device(%s) has installed the os from %q and has timedout.", items[i].SN, startTime)
	}
}
