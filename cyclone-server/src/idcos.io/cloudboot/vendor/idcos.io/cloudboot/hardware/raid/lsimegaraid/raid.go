package lsimegaraid

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"idcos.io/cloudboot/hardware"
	"idcos.io/cloudboot/hardware/raid"
	"idcos.io/cloudboot/logger"
)

const (
	// name RAID处理器名称
	name = raid.LSIMegaRAID
	// tool 硬件配置工具
	tool = "/opt/MegaRAID/MegaCli/MegaCli64"
)

func init() {
	raid.Register(name, new(worker))
}

type worker struct {
	hardware.Base
	mux sync.Mutex
}

// SetDebug 设置是否开启debug。若开启debug，会将关键日志信息写入console。
func (worker *worker) SetDebug(debug bool) {
	worker.Base.SetDebug(debug)
}

// SetLog 更换日志实现。默认情况下内部无日志实现。
func (worker *worker) SetLog(log logger.Logger) {
	worker.mux.Lock()
	defer worker.mux.Unlock()

	worker.Base.SetLog(log)
}

// GetLog 返回日志实例
func (worker *worker) GetLog() logger.Logger {
	worker.mux.Lock()
	defer worker.mux.Unlock()

	return worker.Base.GetLog()
}

// Convert2ControllerID RAID卡控制器索引号转换成RAID卡控制器ID。
// ctrlIndex RAID卡控制器索引号。0表示首块RAID卡，1表示第二块RAID卡，以此类推。
// 如'HP SmartArray' RAID卡，若索引号为'0'，那么可能转换获得的RAID卡控制器ID为'1'，原因是这款RAID卡ID是从'1'开始的。
// 如'LSI SAS3' RAID卡，若索引号为'0'，那么可能转换获得的RAID卡控制器ID依然为'0'，原因是这款RAID卡ID就是从'0'开始的。
func (worker *worker) Convert2ControllerID(ctrlIndex uint) (ctrlID string, err error) {
	ids, err := worker.controllerIDs()
	if err != nil {
		return "", err
	}
	if int(ctrlIndex) > len(ids)-1 {
		return "", raid.ErrControllerNotFound
	}
	return ids[ctrlIndex], nil
}

// Controllers 返回设备的RAID卡列表
func (worker *worker) Controllers() (ctrls []raid.Controller, err error) {
	ids, err := worker.controllerIDs()
	if err != nil {
		return nil, err
	}

	for i := range ids {
		ctrl, err := worker.findController(ids[i])
		if err != nil {
			return nil, err
		}
		ctrls = append(ctrls, *ctrl)
	}
	return ctrls, nil
}

// Clear 擦除raid
func (worker *worker) Clear(ctrlID string) error {
	//%s -CfgForeign -Clear -aALL -NoLog; %s -CfgClr -aALL -NoLog
	_, err := worker.ExecByShell(tool, "-CfgForeign", "-Clear", fmt.Sprintf("-a%s", ctrlID), "-NoLog;")
	if err != nil {
		return err
	}
	worker.Sleep()
	_, err = worker.ExecByShell(tool, "-CfgClr", fmt.Sprintf("-a%s", ctrlID), "-NoLog")
	worker.Sleep()
	return err
}

// InitLogicalDrives 初始化指定RAID卡控制器下的逻辑驱动器(逻辑磁盘)。
func (worker *worker) InitLogicalDrives(ctrlID string) error {
	// MegaCli64 -LDInit -Start -LALL -aALL -NoLog
	_, err := worker.ExecByShell(tool, "-LDInit", "-Start", "-LALL", fmt.Sprintf("-a%s", ctrlID), "-NoLog")
	worker.Sleep()
	return err
}

// PhysicalDrives 返回物理驱动器(物理磁盘)列表。
func (worker *worker) PhysicalDrives(ctrlID string) (pds []raid.PhysicalDrive, err error) {
	items, err := worker.physicalDrives(ctrlID)
	if err != nil {
		return nil, err
	}
	// 排序
	sort.Sort(items)
	pds = make([]raid.PhysicalDrive, 0, len(items))
	for i := range items {
		pds = append(pds, raid.PhysicalDrive{
			ID:           strconv.Itoa(items[i].Slot),
			Name:         items[i].identity(),
			RawSize:      items[i].RawSize,
			MediaType:    items[i].PDType,
			ControllerID: ctrlID,
		})
	}
	return pds, nil
}

// TranslateLegacyDisks 将目标RAID卡控制器下传统"1|2,3|4-6|7-|all"形式的硬盘标识符转换成实际的物理驱动器。
func (worker *worker) TranslateLegacyDisks(ctrlID string, legacyDisks string) (drives []string, err error) {
	pds, err := worker.physicalDrives(ctrlID)
	if err != nil {
		return nil, err
	}
	if regexp.MustCompile("^[[:digit:]]+$").MatchString(legacyDisks) {
		diskNo, _ := strconv.Atoi(legacyDisks)
		return findPhysicalDriveIDsByIndexes(pds, diskNo-1)

	} else if regexp.MustCompile("^[[:digit:],]+$").MatchString(legacyDisks) {
		diskNos := strings.Split(legacyDisks, ",")
		indexes := make([]int, 0, len(diskNos))

		for _, fields := range diskNos {
			diskNo, _ := strconv.Atoi(fields)
			indexes = append(indexes, diskNo-1)
		}
		return findPhysicalDriveIDsByIndexes(pds, indexes...)

	} else if regexp.MustCompile("^[[:digit:]-]+[[:digit:]]$").MatchString(legacyDisks) {
		arr := strings.Split(legacyDisks, "-")
		if len(arr) != 2 {
			return nil, raid.ErrInvalidDiskIdentity
		}
		begin, _ := strconv.Atoi(arr[0])
		end, _ := strconv.Atoi(arr[1])
		if begin >= end {
			return nil, raid.ErrInvalidDiskIdentity
		}
		if begin > len(pds) || end > len(pds) {
			return nil, raid.ErrInvalidDiskIdentity
		}
		var indexes []int
		for diskNo := begin; diskNo <= end; diskNo++ {
			indexes = append(indexes, diskNo-1)
		}
		return findPhysicalDriveIDsByIndexes(pds, indexes...)

	} else if regexp.MustCompile("^[[:digit:]]+-$").MatchString(legacyDisks) {
		arr := strings.Split(legacyDisks, "-")
		begin, _ := strconv.Atoi(arr[0])
		if begin > len(pds) {
			return nil, raid.ErrInvalidDiskIdentity
		}
		var indexes []int
		for diskNo := begin; diskNo <= len(pds); diskNo++ {
			indexes = append(indexes, diskNo-1)
		}
		return findPhysicalDriveIDsByIndexes(pds, indexes...)

	} else if legacyDisks == "all" {
		return findPhysicalDriveIDsByIndexes(pds)
	}
	return nil, raid.ErrInvalidDiskIdentity
}

// CreateArray 用指定的物理磁盘和RAID级别创建阵列(逻辑磁盘)
func (worker *worker) CreateArray(ctrlID string, level raid.Level, drives []string) error {
	// drives数据校验
	err := worker.checkRAID(level, drives)
	if err != nil {
		return err
	}

	// 关闭整个RAID卡控制器的JBOD模式（开启RAID模式）
	if err = worker.switchJBODMode(ctrlID, false); err != nil {
		return err
	}
	worker.Sleep()

	var args []string
	if level == raid.RAID10 {
		// MegaCli64 -CfgSpanAdd -r10 Array0[32:0,N/A:1] Array1[32:2,N/A:3] -a0 -NoLog
		args = append(args, "-CfgSpanAdd", "-r10")
		for i := 0; i < len(drives); i += 2 {
			args = append(args, fmt.Sprintf("Array%d[%s,%s]", i/2, drives[i], drives[i+1]))
		}
	} else if level == raid.RAID50 {
		// MegaCli64 -CfgSpanAdd -r50 Array0[32:0,N/A:1] Array1[32:2,N/A:3] -a0 -NoLog
		//为了方便实现，只分2组，而raid50理论上可以将参数分成多组
		args = append(args, "-CfgSpanAdd", "-r50")
		args = append(args, fmt.Sprintf("Array%d[%s]", 0, strings.Join(drives[:(len(drives)/2)], ",")))
		args = append(args, fmt.Sprintf("Array%d[%s]", 1, strings.Join(drives[(len(drives)/2):], ",")))
	} else {
		//%s -CfgLdAdd -r%s [%s] -a0 -NoLog
		args = append(args, "-CfgLdAdd", fmt.Sprintf("-r%s [%s]", strings.TrimPrefix(string(level), "raid"), strings.Join(drives, ",")))
	}
	args = append(args, "-WB", "-Direct", "-strpsz64", fmt.Sprintf("-a%s", ctrlID), "-NoLog")

	_, err = worker.ExecByShell(tool, args...)
	worker.Sleep()
	return err
}

// SetGlobalHotspares 添加热备盘
func (worker *worker) SetGlobalHotspares(ctrlID string, drives []string) (err error) {
	// -PDHSP -Set -PhysDrv [%s:%s] -a0 -NoLog
	_, err = worker.ExecByShell(tool, "-PDHSP", "-Set", "-PhysDrv", fmt.Sprintf("[%s]", strings.Join(drives, ",")), "-a0", "-NoLog")
	worker.Sleep()
	return err
}

// SetJBODs 将指定RAID卡控制器设置为直通模式，或者将RAID卡控制器下部分物理驱动器设置为直通模式。
// ctrlID 表示RAID卡控制器ID，通过Convert2ControllerID方法可获取到第N块RAID卡的控制器ID。
// drives 物理驱动器列表。若物理驱动器列表为空，则意味着将指定的RAID卡下所有的物理驱动器都设置为直通模式。
func (worker *worker) SetJBODs(ctrlID string, drives []string) (err error) {
	if len(drives) <= 0 {
		// 为整个RAID卡控制器下的物理驱动器开启直通模式
		return worker.switchJBODMode(ctrlID, true)
	}
	// 为指定RAID卡控制器下的部分物理驱动器开启直通模式
	// megacli -PDMakeJBOD -PhysDrv[E0:S0,E1:S1,...] -a0 -NoLog
	_, err = worker.ExecByShell(tool, "-PDMakeJBOD", fmt.Sprintf("-PhysDrv[%s]", strings.Join(drives, ",")), fmt.Sprintf("-a%s", ctrlID), "-NoLog")
	worker.Sleep()
	return err
}

// PostCheck RAID配置实施后置检查
func (worker *worker) PostCheck(sett *raid.Setting) (items []hardware.CheckingItem) {
	items = append(items, worker.checkArrays(sett)...)
	return items
}
