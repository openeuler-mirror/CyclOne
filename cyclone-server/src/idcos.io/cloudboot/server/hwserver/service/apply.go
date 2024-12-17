package service

import (
	"sort"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego/httplib"

	"regexp"

	"idcos.io/cloudboot/hardware/raid"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/server/cloudbootserver/types/setting"
	"idcos.io/cloudboot/server/hwserver/config"
)

// SettingWorker 硬件配置实施器
type SettingWorker struct {
	hwSettings []*setting.HardwareSetting
	conf       *config.Configuration
	log        logger.Logger
	sn         string
}

// NewSettingWorker 返回硬件配置实施器实例
func NewSettingWorker(conf *config.Configuration, log logger.Logger, sn string) *SettingWorker {
	return &SettingWorker{
		conf: conf,
		log:  log,
		sn:   sn,
	}
}

// isEmpty 判断硬件配置是否为空
func (worker *SettingWorker) isEmpty() bool {
	return len(worker.hwSettings) <= 0
}

// Apply 实施硬件配置。
func (worker *SettingWorker) Apply() error {
	// 1、拉取硬件配置
	if err := worker.pull(); err != nil {
		_ = worker.postProgress(failureProgress, "拉取硬件配置失败", err.Error())
		return err
	}

	// 硬件配置参数为空，跳过硬件配置。
	if worker.isEmpty() {
		_ = worker.postProgress(endProgress, "跳过硬件配置", "硬件配置参数为空")
		return nil
	}
	// 3、硬件配置实施
	return worker.apply()
}

// apply 配置实施
func (worker *SettingWorker) apply() (err error) {
	if worker.hwSettings == nil {
		return ErrPullSettings
	}

	//由于当前的hw-server完全是模板是怎么样，就实际下发什么配置，包括磁盘顺序
	//但在实际的环境中，常常由于盘的顺序不一致，导致需要定制不一样的模板，比如2RAID1+12RAID0 或者 12RAID0+2RAID1
	//所以在这里增加一个磁盘交换步骤，将RAID1的盘的序号换成小盘的序号
	//if err = worker.drivesSwap(); err != nil {
	//	worker.log.Errorf("swap drives fail, err:%v", err)
	//}
	//worker.log.Debug("after raid1 swap: hw-setting:")
	raidSum := map[string]int {
		"raid1":0,
		"raid0":0,
	}
	for i := range worker.hwSettings {
		worker.log.Debugf("%v", worker.hwSettings[i])
		switch worker.hwSettings[i].Category {
		case model.CategoryRAID:
			if worker.hwSettings[i].Action == model.ActionRAIDCreate && worker.hwSettings[i].Metadata["level"] == "raid1" {
				worker.log.Debugf("raid level is: %v", worker.hwSettings[i].Metadata["level"])
				raidSum["raid1"] += 1
			}
			if worker.hwSettings[i].Action == model.ActionRAIDCreate && worker.hwSettings[i].Metadata["level"] == "raid0" {
				worker.log.Debugf("raid level is: %v", worker.hwSettings[i].Metadata["level"])
				raidSum["raid0"] += 1
			}
		}
		worker.log.Debugf("raid level sum: %v", raidSum)
	}
	if raidSum["raid1"] == 1 && raidSum["raid0"] == 12 {
		if err = worker.drivesSwap(); err != nil {
			worker.log.Errorf("swap drives fail, err:%v", err)
		}
		worker.log.Debug("after raid1 swap: hw-setting:")
	}

	for i := range worker.hwSettings {
		worker.log.Debug(worker.hwSettings[i])
		switch worker.hwSettings[i].Category {
		case model.CategoryRAID:
			if err = worker.doRAID(i, worker.hwSettings[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (worker *SettingWorker) drivesSwap() (err error) {
	var name string
	name, err = raid.Whoami()
	if err != nil {
		return err
	}

	w := raid.SelectWorker(name)
	if w == nil {
		return raid.ErrUnknownHardware
	}

	//先统计一下有几个ctrl,把每个ctl下的磁盘列表（有序的index）拿到
	//mCtrlPhysicalDrvies := make(map[string]raid.PDSlice, 0)
	var items []raid.PhysicalDrive
	currentIdx := -1 //先给个无效值
	for i := range worker.hwSettings {
		ctrlIdx, _ := strconv.Atoi(worker.hwSettings[i].Metadata["controller_index"])
		if currentIdx != ctrlIdx {
			currentIdx = ctrlIdx
			ctrlID, err := w.Convert2ControllerID(uint(ctrlIdx))
			if err != nil {
				return err
			}
			pds, err := w.PhysicalDrives(ctrlID)
			if err != nil {
				return err
			}
			worker.log.Debugf("PhysicalDrives of ctrl:%d  is  %v", ctrlID, pds)
			items = append(items, pds...)
			worker.log.Debugf("ID of ctrl:%d", ctrlID)
			for _, each := range items {
				worker.log.Debugf("original physical drives :%v", each)
			}
		}
	}

	pdsSort := make(raid.PDSlice, len(items), len(items))
	for j, pd := range items {
		pdsSort[j].Index = j
		pdsSort[j].PhysicalDrive = pd
	}
	// raid PDSlice 实现按disk size 由小到到进行sort
	sort.Sort(pdsSort)
	for _, each := range pdsSort {
		worker.log.Debugf("check each sorted physical drives:%v", each)
	}
	//默认处理2RAID1+12RAID0，若发生过排序则代表2RAID1: 1-2 号盘非最小的，直接与末尾对调
	needSwap := false
	if pdsSort[0].Index >1 && pdsSort[1].Index >1 {
		worker.log.Debug("Index of two front sorted physical drives is not 0,1.Going to swap 1st,2nd to last position")
		needSwap = true
	}

	if len(pdsSort) > 2 {
		lastIndex := len(pdsSort) - 1
		if needSwap {
			pd0, pd1 := strconv.Itoa(lastIndex), strconv.Itoa(lastIndex+1) //记录下raid1实际需要使用的序号
			disk1OfRaid1, disk2OfRaid1 := "", ""                           //记录下raid1的原盘序号
			for i, sett := range worker.hwSettings {
				switch worker.hwSettings[i].Category {
				case model.CategoryRAID:
					worker.log.Debugf("hwSettings:%v", sett)
					sLevel, isExist := sett.Metadata["level"]
					if isExist != true {
						continue
					}
					var level raid.Level
					level, err = raid.SelectLevel(sLevel)
					if err != nil {
						e := fmt.Errorf("%s: %s", err.Error(), sLevel)
						return e
					}
					legacyDisks := sett.Metadata["drives"]
					legacyDisksTrans, err := translateLegacyIndex(len(pdsSort), legacyDisks)
					if err != nil {
						return err
					}
					if level == raid.RAID1 {
						strs := strings.Split(legacyDisksTrans, ",")
						if len(strs) == 2 {
							disk1OfRaid1, disk2OfRaid1 = strs[0], strs[1] //准备交换数据
						}
						sett.Metadata["drives"] = fmt.Sprintf("%s,%s", pd0, pd1)
					} else {
						newLegacy := make([]string, 0)
						for _, d := range strings.Split(legacyDisksTrans, ",") {
							if d == pd0 {
								newLegacy = append(newLegacy, disk1OfRaid1) //交换
							} else if d == pd1 {
								newLegacy = append(newLegacy, disk2OfRaid1) //交换
							} else {
								newLegacy = append(newLegacy, d)
							}
						}
						sett.Metadata["drives"] = strings.Join(newLegacy, ",")
					}
				}
			}
		}
	}
	return nil
}

//将目标中定制的1|2,3|4-6|7-|all,统一翻译成逗号分隔的形式。方便做交换，后续还会翻译成实际的磁盘
func translateLegacyIndex(diskCount int, legacyDisks string) (string, error) {
	if regexp.MustCompile("^[[:digit:]]+$").MatchString(legacyDisks) {
		return legacyDisks, nil

	} else if regexp.MustCompile("^[[:digit:],]+$").MatchString(legacyDisks) {
		return legacyDisks, nil

	} else if regexp.MustCompile("^[[:digit:]-]+[[:digit:]]$").MatchString(legacyDisks) {
		arr := strings.Split(legacyDisks, "-")
		if len(arr) != 2 {
			return "", raid.ErrInvalidDiskIdentity
		}
		begin, _ := strconv.Atoi(arr[0])
		end, _ := strconv.Atoi(arr[1])
		if begin >= end {
			return "", raid.ErrInvalidDiskIdentity
		}
		if begin > diskCount || end > diskCount {
			return "", raid.ErrInvalidDiskIdentity
		}
		var indexes []string
		for diskNo := begin; diskNo <= end; diskNo++ {
			indexes = append(indexes, strconv.Itoa(diskNo))
		}
		return strings.Join(indexes, ","), nil

	} else if regexp.MustCompile("^[[:digit:]]+-$").MatchString(legacyDisks) {
		arr := strings.Split(legacyDisks, "-")
		begin, _ := strconv.Atoi(arr[0])
		if begin > diskCount {
			return "", raid.ErrInvalidDiskIdentity
		}
		var indexes []string
		for diskNo := begin; diskNo <= diskCount; diskNo++ {
			indexes = append(indexes, strconv.Itoa(diskNo))
		}
		return strings.Join(indexes, ","), nil

	} else if legacyDisks == "all" {
		var indexes []string
		for diskNo := 0; diskNo <= diskCount; diskNo++ {
			indexes = append(indexes, strconv.Itoa(diskNo))
		}
		return strings.Join(indexes, ","), nil
	}
	return "", raid.ErrInvalidDiskIdentity
}

func (worker *SettingWorker) doRAID(index int, sett *setting.HardwareSetting) (err error) {
	var name string
	name, err = raid.Whoami()
	if err != nil {
		_ = worker.postProgress(worker.calcProgress(index+1), "[RAID] 识别RAID卡失败", err.Error())
		return err
	}

	w := raid.SelectWorker(name)
	if w == nil {
		_ = worker.postProgress(worker.calcProgress(index+1), "[RAID] 识别RAID卡失败", raid.ErrUnknownHardware.Error())
		return raid.ErrUnknownHardware
	}

	ctrlIdx, _ := strconv.Atoi(sett.Metadata["controller_index"])

	ctrlID, err := w.Convert2ControllerID(uint(ctrlIdx))
	if err != nil {
		_ = worker.postProgress(worker.calcProgress(index+1), "[RAID] 获取控制器ID失败", err.Error())
		return err
	}

	switch sett.Action {
	case model.ActionRAIDClear:
		// 清除RAID配置
		if sett.Metadata["clear"] == setting.ON {
			if err = w.Clear(ctrlID); err != nil {
				_ = worker.postProgress(worker.calcProgress(index+1), "[RAID] 清除配置失败", err.Error())
				return err
			}
			_ = worker.postProgress(worker.calcProgress(index+1), "[RAID] 清除配置", msgdone)
		}

	case model.ActionRAIDCreate:
		// 创建阵列
		sLevel := sett.Metadata["level"]
		var level raid.Level
		level, err = raid.SelectLevel(sLevel)
		if err != nil {
			e := fmt.Errorf("%s: %s", err.Error(), sLevel)
			_ = worker.postProgress(worker.calcProgress(index+1), "[RAID] 无效的RAID级别", e.Error())
			return e
		}

		legacyDisks := sett.Metadata["drives"]
		drives, err := w.TranslateLegacyDisks(ctrlID, legacyDisks)
		if err != nil {
			_ = worker.postProgress(worker.calcProgress(index+1), "[RAID] 转换物理驱动器标识符失败", fmt.Sprintf("%s ==> ? : %s", legacyDisks, err.Error()))
			return err
		}

		if err = w.CreateArray(ctrlID, level, drives); err != nil {
			_ = worker.postProgress(worker.calcProgress(index+1), "[RAID] 创建阵列失败", fmt.Sprintf("%s@%s: %s", sLevel, strings.Join(drives, ","), err.Error()))
			return err
		}
		_ = worker.postProgress(worker.calcProgress(index+1), "[RAID] 创建阵列成功", fmt.Sprintf("%s@%s", sLevel, strings.Join(drives, ",")))

	case model.ActionRAIDSetJBOD:
		// 设置直通模式
		legacyDisks := sett.Metadata["drives"]
		drives, err := w.TranslateLegacyDisks(ctrlID, legacyDisks)
		if err != nil {
			_ = worker.postProgress(worker.calcProgress(index+1), "[RAID] 转换物理驱动器标识符失败", fmt.Sprintf("%s ==> ? : %s", legacyDisks, err.Error()))
			return err
		}

		if err = w.SetJBODs(ctrlID, drives); err != nil {
			_ = worker.postProgress(worker.calcProgress(index+1), "[RAID] 配置直通盘失败", fmt.Sprintf("jbod@%s: %s", strings.Join(drives, ","), err.Error()))
			return err
		}
		_ = worker.postProgress(worker.calcProgress(index+1), "[RAID] 配置直通盘成功", fmt.Sprintf("jbod@%s", strings.Join(drives, ",")))

	case model.ActionRAIDSetGlobalHotspare:
		// 设置全局热备盘
		legacyDisks := sett.Metadata["drives"]
		drives, err := w.TranslateLegacyDisks(ctrlID, legacyDisks)
		if err != nil {
			_ = worker.postProgress(worker.calcProgress(index+1), "[RAID] 转换物理驱动器标识符失败", fmt.Sprintf("%s ==> ? : %s", legacyDisks, err.Error()))
			return err
		}

		if err = w.SetGlobalHotspares(ctrlID, drives); err != nil {
			_ = worker.postProgress(worker.calcProgress(index+1), "[RAID] 配置全局热备盘失败", fmt.Sprintf("global hotspare@%s: %s", strings.Join(drives, ","), err.Error()))
			return err
		}
		_ = worker.postProgress(worker.calcProgress(index+1), "[RAID] 配置全局热备盘成功", fmt.Sprintf("global hotspare@%s", strings.Join(drives, ",")))

	case model.ActionRAIDInitDisk:
		// 初始化逻辑磁盘
		if sett.Metadata["init"] == setting.ON {
			err = w.InitLogicalDrives(ctrlID)
			if err == raid.ErrInitLogicalDrivesNotSupport {
				_ = worker.postProgress(worker.calcProgress(index+1), "[RAID] 跳过逻辑磁盘初始化", msgskip)
			} else if err != nil {
				_ = worker.postProgress(worker.calcProgress(index+1), "[RAID] 逻辑磁盘初始化失败", err.Error())
				return err
			}
			_ = worker.postProgress(worker.calcProgress(index+1), "[RAID] 逻辑磁盘初始化成功", msgdone)
		}
	}
	return nil
}

// ErrFetchScript 无法从服务端获取脚本文件错误
var ErrFetchScript = errors.New("unable to get script file from server")

// calcProgress 计算当前步骤下硬件配置进度值
func (worker *SettingWorker) calcProgress(current int) float64 {
	total := len(worker.hwSettings)
	if total == 0 {
		return 0.0
	}
	return startProgress + float64(current)/float64(total*10)
}

var (
	// ErrPullSettings 无法拉取到硬件配置参数
	ErrPullSettings = errors.New("unreachable hardware settings")
)

// pull 拉取硬件配置参数
func (worker *SettingWorker) pull() (err error) {
	var respData struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Content struct {
			Settings []*setting.HardwareSetting `json:"settings"`
		} `json:"content"`
	}
	url := fmt.Sprintf("%s/api/cloudboot/v1/devices/%s/settings/hardwares?lang=en-US", strings.TrimSuffix(worker.conf.HTTPServer.BaseURL, "/"), worker.sn)
	if err = httplib.Get(url).Header("Accept", "application/json").ToJSON(&respData); err != nil {
		worker.log.Errorf("GET %s ==>\n%s", url, err)
		return err
	}

	body, _ := json.Marshal(respData)
	worker.log.Debugf("GET %s ==>\n%s", url, body)

	if respData.Status != "success" {
		return ErrPullSettings
	}
	worker.hwSettings = respData.Content.Settings
	return nil
}

const (
	failureProgress = -1.0
	startProgress   = 0.3
	endProgress     = 0.4
)
const (
	msgdone = "完成"
	msgskip = "跳过"
)

func (worker *SettingWorker) postProgress(progress float64, title, message string) (err error) {
	var reqData struct {
		InstallProgress float64 `json:"progress"`
		InstallLog      string  `json:"log"`
		Title           string  `json:"title"`
	}
	reqData.InstallProgress = progress
	reqData.Title = title
	reqData.InstallLog = base64.StdEncoding.EncodeToString([]byte(message))

	var respData struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	body, _ := json.Marshal(reqData)
	url := fmt.Sprintf("%s/api/cloudboot/v1/devices/%s/installations/progress?lang=en-US", strings.TrimSuffix(worker.conf.HTTPServer.BaseURL, "/"), worker.sn)
	worker.log.Debugf("POST %s ==>\n%s", url, body)

	req, err := httplib.Post(url).Header("Content-Type", "application/json").JSONBody(&reqData)
	if err != nil {
		worker.log.Errorf("%s", err)
		return err
	}
	if err = req.ToJSON(&respData); err != nil {
		worker.log.Errorf("%s", err)
		return err
	}

	if respData.Status != "success" {
		//return errors.New(respData.Status)
		panic("report progress fail, respData.Status: " + respData.Status)
	}
	return nil
}
