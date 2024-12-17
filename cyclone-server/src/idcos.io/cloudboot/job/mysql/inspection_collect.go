package mysql

import (
	"encoding/json"
	"errors"
	"time"

	"idcos.io/cloudboot/model"
)

// Collect 采集目标设备的ipmi传感器、事件数据
// TODO 这个需求用于实时查询ipmi传感器、事件数据接口，可以更优雅，应该把这部分逻辑提炼到utils/oob里面
func (j *InspectionJob) Collect(sn string) (insp *model.Inspection) {
	insp = &model.Inspection{}
	var err error
	defer func(insp *model.Inspection) {
		// 4、巡检结果
		endTime := time.Now()
		insp.EndTime = &endTime
		insp.RunningStatus = model.RunningStatusDone
		if err != nil {
			j.log.Error(err)
			insp.Error = j.translateErrMsg(ErrServer)
			return
		}
	}(insp)
	insp.SN = sn
	insp.HealthStatus = model.HealthStatusUnknown
	// 1、查询设备信息
	dev, err := j.repo.GetDeviceBySN(insp.SN)
	if err != nil {
		err = ErrDeviceSN
		return
	}

	// 3、采集传感器信息
	switch dev.Arch {
	case model.DevArchX8664, model.DevArchAarch64, "": // 暂时仅支持x86_64、arm服务器
		sensors, err := j.fetchSensors(dev)
		insp.HealthStatus = j.healthStatus(sensors)
		if err != nil {
			return
		}
		b, err := json.Marshal(sensors)
		if err != nil {
			return
		}
		insp.IPMIResult = string(b)
	default:
		err = errors.New("Unsupported arch: " + dev.Arch)
	}
	return
}
