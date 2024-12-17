package mysql

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"runtime/debug"
	"sync"

	"idcos.io/cloudboot/config"
	"idcos.io/cloudboot/job"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils"
	"idcos.io/cloudboot/utils/oob"
	"idcos.io/cloudboot/utils/sh"
)

// InspectionJob 硬件巡检任务
type InspectionJob struct {
	id      string // 任务ID(全局唯一)
	log     logger.Logger
	repo    model.Repo
	conf    *config.Config
	retries int
}

// IPMIData ipmi工具巡检数据
type IPMIData struct {
	IPMISensorData    []*model.SensorData
	IPMISelData       []*model.SelData
}

// NewInspectionJob 实例化任务管理器
func NewInspectionJob(log logger.Logger, repo model.Repo, conf *config.Config, jobid string) *InspectionJob {
	return &InspectionJob{
		log:     log,
		repo:    repo,
		conf:    conf,
		id:      jobid,
		retries: 5,
	}
}

// Run 运行硬件巡检任务
func (j *InspectionJob) Run() {
	j.log.Infof("[%s]Inspection job starts running", j.id)
	defer j.log.Infof("[%s]Inspection job completed", j.id)

	defer func() {
		if err := recover(); err != nil {
			j.log.Errorf("[%s]Inspection job panic: \n%s\n%s", j.id, err, debug.Stack())
		}
	}()

	// 1、根据任务ID获取任务明细
	mjob, err := j.repo.GetJobByID(j.id)
	if err != nil {
		j.log.Error(err)
		return
	}
	j.log.Debugf("[%s]Inspection job: %+v", j.id, mjob)

	// 2、提取出任务的目标SN列表
	var target map[string][]string
	if err = json.Unmarshal([]byte(mjob.Target), &target); err != nil {
		j.log.Error(err)
		return
	}

	// 若未指定具体的巡检设备，则巡检所有设备，过滤待退役、已退役设备
	if len(target) <= 0 {
		target, _ = j.getOriginNodeDevicesPairs()
		if len(target) <= 0 {
			j.log.Warnf("[%s]The device to be inspected cannot be obtained. The inspection is terminated.", j.id)
			return
		}
	}

	startTime := time.Now()

	// 3、将整体任务分解成针对各个SN的巡检子任务并持久化
	var insps []*model.Inspection
	for node, sns := range target {
		for _, sn := range sns {
			insps = append(insps, &model.Inspection{
				JobID:         j.id,
				StartTime:     &startTime,
				OriginNode:    node,
				SN:            sn,
				RunningStatus: model.RunningStatusRunning,
				HealthStatus:  model.HealthStatusUnknown,
			})
		}
	}
	if err = j.repo.AddInspections(insps...); err != nil {
		j.log.Error(err)
		return
	}
	j.log.Debugf("[%s]Save inspection records to MySQL", j.id)

	// 4、遍历巡检子任务并依次采集ipmi传感器、事件数据等后将结果持久化
	var wg sync.WaitGroup
	sem := make(chan struct{}, 50) // 最多允许50个并发同时执行

	for i := range insps {
		wg.Add(1)
		go func(k int) {
			sem <- struct{}{}        // 获取信号
			defer func() { <-sem }() // 释放信号
			j.collectThenRestore(insps[k])
			wg.Done()
		}(i)		
	}
	wg.Wait()

	// 5、若当前任务为一次性执行的任务，则将任务状态修改为'已停止'。
	if mjob.Rate == job.RateImmediately {
		_ = j.repo.SaveJob(&model.Job{
			ID:     j.id,
			Status: job.Stoped,
		})
	}

	// 6. 发送告警邮件
	//go SendNotifyByJobID(j)
}

// translateErrMsg 错误信息翻译转换
func (j *InspectionJob) translateErrMsg(err error) string {
	switch err {
	case ErrServer:
		return "服务器内部错误"
	case ErrDeviceSN:
		return "该序列号设备不存在"
	case ErrProxyIPUnreachable:
		return "代理节点IP不可达"
	case ErrFetchOOB:
		return "带外信息缺失"
	case ErrIPMIDataCollection:
		return "远程采集IPMI数据失败(请检查带外IP、用户名、密码)"
	case ErrCollectDataByProxy:
		return "通过代理节点采集数据失败"
	case ErrInvalidUAMToken:
		return "请检查配置的UAM令牌是否有效"
	case oob.ErrOOBIPUnreachable:
		return "带外IP不可达"
	case oob.ErrMissingOOBInfo:
		return "带外信息(IP、用户名、密码)不完整"
	case oob.ErrUsernamePassword:
		return "带外用户名、密码不匹配"
	case oob.ErrMissingIPMImonitoring:
		return "ipmimonitoring工具未安装在服务端PATH环境变量目录下"
	}
	return err.Error()
}

var (
	// ErrServer 服务器错误
	ErrServer = errors.New("internal server error")
	// ErrDeviceSN 无效的设备SN
	ErrDeviceSN = errors.New("invalid device SN")
	// ErrProxyIPUnreachable 代理节点IP不可达
	ErrProxyIPUnreachable = errors.New("proxy ip is unreachable")
	// ErrFetchOOB 无法获取设备OOB信息
	ErrFetchOOB = errors.New("oob information is not available")
	// ErrIPMIDataCollection 远程采集IPMI数据失败(请检查带外用户名和密码)
	ErrIPMIDataCollection = errors.New("ipmi data collection failed")
	// ErrCollectDataByProxy 通过代理节点采集数据失败
	ErrCollectDataByProxy = errors.New("failed to collect data through proxy node")
	// ErrInvalidUAMToken 无效的UAM令牌
	ErrInvalidUAMToken = errors.New("invalid UAM token")
)

// collectThenRestore 采集目标设备的ipmi传感器、事件数据并持久化
func (j *InspectionJob) collectThenRestore(insp *model.Inspection) {
	defer func(insp *model.Inspection) {
		// 4、巡检结果持久化
		endTime := time.Now()
		insp.EndTime = &endTime
		insp.RunningStatus = model.RunningStatusDone

		// 将结果保存至MySQL
		j.log.Debugf("[%s]Start saving inspection results: %+v", j.id, *insp)
		if affected, _ := j.repo.UpdateInspectionByID(insp); affected > 0 {
			j.log.Debugf("[%s]The inspection result is saved", j.id)
		}
		
		// 清理较旧的巡检记录
		j.log.Debugf("[%s]Start deleting old inspection results", j.id)
		if err := j.repo.RemoveInspectionOnStartTimeBySN(insp.SN); err != nil {
			j.log.Debugf("Old inspection result of %s is not deleted, ERR: %s ", insp.SN, err)
		} else {
			j.log.Debugf("Only keep few latest inspection results of %s , old record had been deleted.", insp.SN)
		}

	}(insp)

	insp.HealthStatus = model.HealthStatusUnknown

	// 1、查询设备信息
	dev, err := j.repo.GetDeviceBySN(insp.SN)
	if err != nil {
		err = ErrDeviceSN
		insp.Error = j.translateErrMsg(err)
		return
	}

	// 3、采集
	switch dev.Arch {
	case model.DevArchX8664, model.DevArchAarch64, "": // 暂时仅支持x86_64、arm服务器
		// 获取数据
		ipmidata, err := j.fetchIPMIData(dev)
		if err != nil {
			insp.Error = j.translateErrMsg(err)
			return
		}

        // 处理传感器数据
		jsonSensorData, err := json.Marshal(ipmidata.IPMISensorData)
		if err != nil {
			j.log.Error(err)
			insp.Error = j.translateErrMsg(ErrServer)
			return
		}
		// 根据传感器安状态判断整体健康状况
		insp.HealthStatus = j.healthStatus(ipmidata.IPMISensorData)
		insp.IPMIResult = string(jsonSensorData)
		
		// 处理系统事件日志
		jsonSelData, err := json.Marshal(ipmidata.IPMISelData)
		if err != nil {
			j.log.Error(err)
			insp.Error = j.translateErrMsg(ErrServer)
			return
		}
		insp.IPMISELResult = string(jsonSelData)
		
	default:
		j.log.Warnf("Unsupported arch: %s", dev.Arch)
	}
}

// fetchIPMIData 获取指定设备巡检项的ipmi数据(传感器、系统事件))
func (j *InspectionJob) fetchIPMIData(dev *model.Device) (*IPMIData, error) {
	j.log.Debugf("[%s]Start collecting ipmi data for device(%s)", j.id, dev.SN)

	if dev.OOBUser == "" || dev.OOBPassword == "" {
		j.log.Warnf("设备带外用户或密码为空，尝试找回，[SN:%s]", dev.SN)
		history, err := j.repo.GetLastOOBHistoryBySN(dev.SN)
		if err != nil {
			j.log.Errorf("find back oob history by sn:%s fail,%s", dev.SN, err.Error())
			return nil, err
		}
		dev.OOBUser = history.UsernameNew
		dev.OOBPassword = history.PasswordNew
	}
	if dev.SN == "" || dev.OOBUser == "" || dev.OOBPassword == "" {
		j.log.Warnf("[%s]OOB information for the device(%s) is incomplete", j.id, dev.SN)
		return nil, oob.ErrMissingOOBInfo
	}

	oobHost := utils.GetOOBHost(dev.SN, dev.Vendor, j.conf.Server.OOBDomain)
	// dns dig name后若解析到的ip与device表oob_ip字段不一致则更新入库
	oobIP := oob.TransferHostname2IP(j.log, j.repo, dev.SN, oobHost)
    if oobIP == "" {
		return nil, oob.ErrMissingOOBInfo
	}

	pwd, err := utils.AESDecrypt(dev.OOBPassword, []byte(j.conf.Crypto.Key))
	if err != nil {
		return nil, oob.ErrMissingOOBInfo
	}
	password := string(pwd)

	var ipmidata IPMIData
	// 调用ipmimonitoring工具采集传感器数据
	sensors, err := j.ipmimonitoring(oobIP, dev.OOBUser, password)
	if err != nil {
		j.log.Warnf("[%s]IPMI-Sensor data collection failed", j.id)
		return nil, ErrIPMIDataCollection
	}
	if sensors != nil {
		ipmidata.IPMISensorData = sensors
	}
	// 调用ipmi-sel工具采集系统事件日志
	sel, err := j.ipmisel(oobIP, dev.OOBUser, password)
	if err != nil {
		j.log.Warnf("[%s]IPMI-Sel data collection failed", j.id)
		return nil, ErrIPMIDataCollection
	}
	if sel != nil {
		ipmidata.IPMISelData = sel
	}
	return &ipmidata, nil
}

// ipmimonitoring 通过ipmimonitoring工具获取指定带外的传感器数据
func (j *InspectionJob) ipmimonitoring(ip, username, password string) (items []*model.SensorData, err error) {
	// -W authcap 各厂商实现IPMI协议存在不同，部分型号需要指定 workaround
	cmd := fmt.Sprintf("ipmimonitoring --interpret-oem-data --ignore-not-available-sensors --ignore-unrecognized-events --output-sensor-state --entity-sensor-names --comma-separated-output --no-header-output -D LAN_2_0 -W authcap -h %s -u %s -p %s", ip, username, password)
	out, err := sh.ExecOutputWithLog(j.log, cmd)
	if err != nil {
		j.log.Debugf("error of ipmimonitoring: %s", err)
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Split(line, ",")
		if len(fields) != 7 {
			continue
		}
		items = append(items, &model.SensorData{
			ID:      fields[0],
			Name:    fields[1],
			Type:    fields[2],
			State:   fields[3],
			Reading: fields[4],
			Units:   fields[5],
			Event:   strings.Replace(fields[6], "'", "", -1),
		})
	}
	return items, scanner.Err()
}

// ipmi-sel 通过ipmi-sel工具获取指定带外的系统事件日志
func (j *InspectionJob) ipmisel(ip, username, password string) (items []*model.SelData, err error) {
	// -W authcap 各厂商实现IPMI协议存在不同，部分型号需要指定 workaround
	cmd := fmt.Sprintf("ipmi-sel --tail=50 --output-event-state --non-abbreviated-units --interpret-oem-data --system-event-only --comma-separated-output --no-header-output -D LAN_2_0 -W authcap -h %s -u %s -p %s", ip, username, password)
	out, err := sh.ExecOutputWithLog(j.log, cmd)
	if err != nil {
		j.log.Debugf("error of ipmi-sel: %s", err)
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Split(line, ",")
		if len(fields) != 7 {
			continue
		}
		items = append(items, &model.SelData{
			ID:      fields[0],
			Date:    fields[1],
			Time:    fields[2],
			Name:    fields[3],
			Type:    fields[4],
			State:   fields[5],
			Event:   fields[6],
		})
	}
	return items, scanner.Err()
}

// healthStatus 返回设备的健康状况。判定规则：
// 1、若未实际获取ipmi传感器数据，则该设备健康状况为未知(unknown)。
// 2、若任意传感器包含异常(critical)，则该设备健康状况为异常(critical)。
// 3、若任意传感器都未包含异常，但存在一个或者多个传感器存在警告(warning)，则该设备健康状况为警告(warning)。
// 4、除了以上情况外，该设备健康状况为正常(nominal)。
func (j *InspectionJob) healthStatus(items []*model.SensorData) string {
	if len(items) <= 0 {
		return model.HealthStatusUnknown
	}
	var hasWarnings bool
	for _, sd := range items {
		switch strings.ToLower(sd.State) {
		case model.HealthStatusCritical:
			return model.HealthStatusCritical

		case model.HealthStatusWarning:
			hasWarnings = true
		}
	}
	if hasWarnings {
		return model.HealthStatusWarning
	}
	return model.HealthStatusNominal
}

// getOriginNodeDevicesPairs 根据config文件 NODE IP 构成的键值对
// 实现各机房管理单元NODE节点对应自身机房的SN，多个NODE节点时轮询分配
func (j *InspectionJob) getOriginNodeDevicesPairs() (pairs map[string][]string, err error) {
	devices, err := j.repo.GetCombinedDevices(nil, nil, nil)
	if err != nil {
		return nil, err
	}
	pairs = make(map[string][]string, 1)
	var nodeIPs []string
	var nodeIPLength int
	for i := range devices {
		// 忽略待退役、已退役设备
		if devices[i].OperationStatus == model.DevOperStatPreRetire || devices[i].OperationStatus == model.DevOperStatRetiring || devices[i].OperationStatus == model.DevOperStateRetired {
			j.log.Debugf("Device sn %s is %s  ,ignore it.", devices[i].SN, devices[i].OperationStatus)
			continue
		}

		if devices[i].ServerRoomID != 0 {
			if nodeIPLength = len(middleware.MapDistributeNode.MDistribute[devices[i].ServerRoomID]); nodeIPLength != 0 {
				nodeIPs = middleware.MapDistributeNode.MDistribute[devices[i].ServerRoomID]
				_, ok := pairs[nodeIPs[i%nodeIPLength]]
				if !ok {
					pairs[nodeIPs[i%nodeIPLength]] = []string{}
				}
				pairs[nodeIPs[i%nodeIPLength]] = append(pairs[nodeIPs[i%nodeIPLength]], devices[i].SN)
			} else {
				_, ok := pairs["master"]
				if !ok {
					pairs["master"] = []string{}
				}
				pairs["master"] = append(pairs["master"], devices[i].SN)				
			}
		}
	}
	return pairs, nil
}


// fetchSensors 获取指定设备巡检项的ipmi传感器数据
func (j *InspectionJob) fetchSensors(dev *model.Device) (items []*model.SensorData, err error) {
	j.log.Debugf("[%s]Start collecting sensor data for device(%s)", j.id, dev.SN)

	if dev.OOBUser == "" || dev.OOBPassword == "" {
		j.log.Warnf("设备带外用户或密码为空，尝试找回，[SN:%s]", dev.SN)
		history, err := j.repo.GetLastOOBHistoryBySN(dev.SN)
		if err != nil {
			j.log.Errorf("find back oob history by sn:%s fail,%s", dev.SN, err.Error())
			return nil, err
		}
		dev.OOBUser = history.UsernameNew
		dev.OOBPassword = history.PasswordNew
	}
	if dev.SN == "" || dev.OOBUser == "" || dev.OOBPassword == "" {
		j.log.Warnf("[%s]OOB information for the device(%s) is incomplete", j.id, dev.SN)
		return nil, oob.ErrMissingOOBInfo
	}

	oobHost := utils.GetOOBHost(dev.SN, dev.Vendor, j.conf.Server.OOBDomain)
	oobIP := oob.TransferHostname2IP(j.log, j.repo, dev.SN, oobHost)

	// 调用ipmimonitoring工具采集传感器数据
	pwd, err := utils.AESDecrypt(dev.OOBPassword, []byte(j.conf.Crypto.Key))
	if err != nil {
		return nil, oob.ErrMissingOOBInfo
	}
	password := string(pwd)
	items, err = j.ipmimonitoring(oobIP, dev.OOBUser, password)
	if err != nil {
		j.log.Warnf("[%s]IPMI data collection failed", j.id)
		return nil, ErrIPMIDataCollection
	}
	return items, nil
}