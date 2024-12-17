package service

import (
	"strings"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/voidint/binding"

	"time"

	"idcos.io/cloudboot/config"
	"idcos.io/cloudboot/limiter"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils/centos6"
	"idcos.io/cloudboot/utils/ping"
)

// InstallProgressReq 安装进度上报结构体
type InstallProgressReq struct {
	DeviceSettingID uint    `json:"-"`
	SN              string  `json:"-"`
	Title           string  `json:"title"`
	Progress        float64 `json:"progress"`
	Log             string  `json:"log"`
}

// FieldMap 请求字段映射
func (reqData *InstallProgressReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.Title:    "title",
		&reqData.Progress: "progress",
		&reqData.Log:      "log",
	}
}

// Validate 装机参数校验
func (reqData *InstallProgressReq) Validate(r *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(r.Context())

	if reqData.SN == "" {
		errs.Add([]string{"sn"}, binding.RequiredError, "sn不能为空")
	}

	if reqData.Progress > 1 || (reqData.Progress < 0 && reqData.Progress != -1) {
		errs.Add([]string{"progress"}, binding.BusinessError, "无效的进度值")
		return errs
	}
	sett, err := repo.GetDeviceSettingBySN(reqData.SN)
	if err == gorm.ErrRecordNotFound {
		errs.Add([]string{"sn"}, binding.BusinessError, "未发现设备装机参数")
		return errs
	}
	if err != nil {
		errs.Add([]string{"sn"}, binding.SystemError, err.Error())
		return errs
	}
	reqData.DeviceSettingID = sett.ID

	//有个需求是即使安装失败的状态，也可以通过手动干预的方式在确认OK（比如检查脚本）后上报进度，故不再限制死此状态
	//已经安装成功或者已经取消安装(status==failure,progress==0.0)
	if (sett.Status == model.InstallStatusSucc) ||
		(sett.Status == model.InstallStatusFail && sett.InstallProgress == 0.0) {
		errs.Add([]string{"sn"}, binding.BusinessError, "当前设备并不在安装列表中")
		return errs
	}
	return errs
}

// ReportInstallProgress 处理安装进度上报
func ReportInstallProgress(logger logger.Logger, repo model.Repo, conf *config.Config, lim limiter.Limiter, reqData *InstallProgressReq) (err error) {
	logger.Infof("title: %v, progress: %v, log: %v\n", reqData.Title, reqData.Progress, reqData.Log)

	var status, title, installLog string
	if byteDecode, err := base64.StdEncoding.DecodeString(reqData.Log); err == nil {
		installLog = string(byteDecode)
	}

	if (reqData.Progress == 1 || reqData.Progress == -1) && centos6.IsPXEUEFI(logger, repo, reqData.SN) {
		_ = centos6.DropConfigurations(logger, repo, reqData.SN) // TODO 为支持centos6的UEFI方式安装而临时增加的逻辑，后续会删除。
	}

	go func() {
		switch reqData.Progress {
		case 1: // 安装成功
			devSett, err := GetDeviceSettingBySN(logger, repo, reqData.SN)
			if err != nil {
				logger.Error("get device setting by sn:%s fail,%v", reqData.SN, err)
				//return err
			}
			// 可能有多个IP逗号分隔的场景，默认取第一个
			intranetIPList := strings.Split(devSett.IntranetIP.IP, ",")
			//异步ping 内网IP，直到ping通或者超时
			timeout := 30 * time.Minute //30 min
			ok := make(chan bool, 0)
			go func() {
				for i := 0; i < int(timeout/(10*time.Second)); i++ {
					time.Sleep(10 * time.Second)
					if err := ping.PingTest(intranetIPList[0]); err == nil {
						ok <- true
						return
					}
				}
			}()

			select {
			case <-time.After(timeout):
				// 临门一脚失败了，设备ping超时
				logger.Error("device:%s install complete but intranet ip not available after reboot", reqData.SN)
				status, title, installLog = model.InstallStatusFail, reqData.Title, installLog+"(but ping failed at last)"
				reqData.Progress = -1.0		
			case <-ok:
				// 千辛万苦，部署成功啦！
				logger.Infof("Congratulations!!! device:%s installed success", reqData.SN)
				status = model.InstallStatusSucc
				title = fmt.Sprintf("%s(%.1f%%)", reqData.Title, reqData.Progress*100)
				installLog += "(and ping client successful!)"
				//加个特殊的流程begin
				//如果备注字段有运行状态，则说明是回收重装的
				d, err := repo.GetDeviceBySN(devSett.SN)
				if d == nil {
					logger.Error("device:%s not exist,%v", devSett.SN, err)
					//return err
				} else if d != nil && d.Remark != "" && validOperationStatus(d.Remark) {
					d.OperationStatus = d.Remark
					d.Remark = "" //置空,这里会比较暴力
				} else {
					d.OperationStatus = model.DevOperStatOnShelve // 系统安装完成后运营状态改为'已上架'
				}
				//end
				d.RAIDRemark = devSett.HardwareTemplate.Name //将RAID模板名称同步写到device表的RAID描述字段

				if _, err := repo.SaveDevice(d); err != nil {
					logger.Error("save device:%d fail,%v", d, err)
				}
				// 归还令牌
				if conf.DHCPLimiter.Enable {
					if bucket, _ := lim.Route(reqData.SN); bucket != nil {
						if token, _ := repo.GetTokenBySN(reqData.SN); token != "" {
							_ = bucket.Return(reqData.SN, limiter.Token(token))
						}
					}
				}
			}
		case -1: // 安装失败
			status, title = model.InstallStatusFail, reqData.Title
		default: // 安装中
			status, title = model.InstallStatusIng, fmt.Sprintf("%s(%.1f%%)", reqData.Title, reqData.Progress*100)
		}

		if _, err := repo.UpdateInstallStatusAndProgressByID(reqData.DeviceSettingID, status, reqData.Progress); err != nil {
			logger.Errorf("update install status by id:%d err:%v", reqData.DeviceSettingID, err)
			//return err
		}

		_, err = repo.SaveDeviceLog(&model.DeviceLog{
			DeviceSettingID: reqData.DeviceSettingID,
			Title:           title,
			LogType:         model.DeviceLogInstallType,
			Content:         installLog,
			SN:              reqData.SN,
		})
		if err != nil {
			logger.Errorf("save device log err:%v", err)
		}
	}()
	return err
}

// InstallationStatus 设备操作系统安装状态
type InstallationStatus struct {
	Type     string  `json:"type"`
	Status   string  `json:"status"`
	Progress float64 `json:"progress"`
}

// GetInstallationStatus 返回指定设备的OS安装状态
func GetInstallationStatus(log logger.Logger, repo model.Repo, sn string) (*InstallationStatus, error) {
	sett, err := repo.GetDeviceSettingBySN(sn)
	if err != nil {
		return nil, err
	}
	return &InstallationStatus{
		Type:     sett.InstallType,
		Status:   sett.Status,
		Progress: sett.InstallProgress,
	}, nil
}

// GetIsInInstallListBySN 返回指定设备是否处于装机队列的布尔值
func GetIsInInstallListBySN(log logger.Logger, repo model.Repo, sn string) (isInList bool, err error) {
	// var affected int64
	// if affected, err = repo.UpdateDeviceBySN(&model.Device{
	// 	SN:                   sn,
	// 	BootosLastActiveTime: time.Now().Format("2006-01-02 15:04:05"),
	// }, true); err != nil {
	// 	log.Errorf("update bootos last active time error:%s", err.Error())
	// }
	// if affected <= 0 {
	// 	return false, nil
	// }

	sett, err := repo.GetDeviceSettingBySN(sn)
	if err != nil {
		return false, err
	}
	return sett.Status == model.InstallStatusPre || sett.Status == model.InstallStatusIng, nil
}
