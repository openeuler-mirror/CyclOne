package mysql

import (
	"time"
	"idcos.io/cloudboot/config"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
)

// UpdateDeviceLifecycleJob 更新设备生命周期任务（维保状态）
type UpdateDeviceLifecycleJob struct {
	id   string // 任务ID(全局唯一)
	log  logger.Logger
	repo model.Repo
	conf *config.Config
}

// NewUpdateDeviceLifecycleJob 实例化任务管理器
func NewUpdateDeviceLifecycleJob(log logger.Logger, repo model.Repo, conf *config.Config, jobid string) *UpdateDeviceLifecycleJob {
	return &UpdateDeviceLifecycleJob{
		log:  log,
		repo: repo,
		conf: conf,
		id:   jobid,
	}
}

// Run 运行更新设备生命周期任务（维保状态）
// 内置任务，根据规则更新维保状态
// current_date < maintenance_service_date_begin                                -> maintenance_service_status = inactive
// maintenance_service_date_begin < current_date < maintenance_service_date_end -> maintenance_service_status = under_warranty
// maintenance_service_date_end < current_date                                  -> maintenance_service_status = out_of_warranty
func (j *UpdateDeviceLifecycleJob) Run() {
	j.log.Debugf("[UpdateDeviceLifecycleJob-begin]- compare current_date to maintenance_service_date and update maintenance_service_status")
	defer j.log.Debugf("[UpdateDeviceLifecycleJob-end] - maintenance_service_status update completed")

	defer func() {
		if err := recover(); err != nil {
			j.log.Errorf("UpdateDeviceLifecycle job panic: \n%s", err)
		}
	}()
	
	devices, err := j.repo.GetCombinedDevices(nil, nil, nil)
	if err != nil {
		j.log.Errorf("UpdateDeviceLifecycle job err: \n%s", err)
	}
	for i := range devices {
		// DeviceLifecycle 查询是否已经存在
		devLifecycle, err := j.repo.GetDeviceLifecycleBySN(devices[i].SN)
		if err != nil {
			j.log.Errorf("get deviceLifecycle of %s err: \n%s", devices[i].SN, err)
			continue
		}
		// 过滤 已退役&&已过保
		if devices[i].OperationStatus == model.DevOperStateRetired && devLifecycle.MaintenanceServiceStatus == model.MaintenanceServiceStatusOutOfWarranty {
			j.log.Debugf("Device sn %s is %s  ,ignore it.", devices[i].SN, model.MaintenanceServiceStatusOutOfWarranty)
			continue
		}
		// DeviceLifecycle 结构体
		modDevLifecycle := &model.DeviceLifecycle{
			SN:             				devices[i].SN,
		}
		now := time.Now() 
		currentDate := now.Format("2006-01-02")
		if now.Before(devLifecycle.MaintenanceServiceDateBegin) {
		    if devLifecycle.MaintenanceServiceStatus != model.MaintenanceServiceStatusInactive {
				j.log.Debugf("current_date=%s < maintenance_service_begin_date=%s. Update Device(SN:%s) set status=%s -> %s", 
					currentDate, devLifecycle.MaintenanceServiceDateBegin, devices[i].SN, devLifecycle.MaintenanceServiceStatus, model.MaintenanceServiceStatusInactive)
				modDevLifecycle.MaintenanceServiceStatus = model.MaintenanceServiceStatusInactive
				if err = j.repo.UpdateDeviceLifecycleBySN(modDevLifecycle);err != nil {
					j.log.Errorf("UpdateDeviceLifecycleBySN failed: %s", err.Error())
					continue
				}
			}
		} else if now.After(devLifecycle.MaintenanceServiceDateEnd) {
			if devLifecycle.MaintenanceServiceStatus != model.MaintenanceServiceStatusOutOfWarranty {
				j.log.Debugf("current_date=%s > maintenance_service_end_date=%s. Update Device(%s) set status=%s -> %s", 
					currentDate, devLifecycle.MaintenanceServiceDateEnd, devices[i].SN, devLifecycle.MaintenanceServiceStatus, model.MaintenanceServiceStatusOutOfWarranty)
				modDevLifecycle.MaintenanceServiceStatus = model.MaintenanceServiceStatusOutOfWarranty
				if err = j.repo.UpdateDeviceLifecycleBySN(modDevLifecycle);err != nil {
					j.log.Errorf("UpdateDeviceLifecycleBySN failed: %s", err.Error())
					continue
				}
			}
		} else {
			if devLifecycle.MaintenanceServiceStatus != model.MaintenanceServiceStatusUnderWarranty {
				j.log.Debugf("current_date=%s between  maintenance_service_begin_date=%s and maintenance_service_end_date=%s. Update Device(%s) set status=%s -> %s",
					currentDate, devLifecycle.MaintenanceServiceDateBegin, devLifecycle.MaintenanceServiceDateEnd, devices[i].SN, devLifecycle.MaintenanceServiceStatus, model.MaintenanceServiceStatusUnderWarranty)
				modDevLifecycle.MaintenanceServiceStatus = model.MaintenanceServiceStatusUnderWarranty
				if err = j.repo.UpdateDeviceLifecycleBySN(modDevLifecycle);err != nil {
					j.log.Errorf("UpdateDeviceLifecycleBySN failed: %s", err.Error())
					continue
				}
			}
		}
	}	
}