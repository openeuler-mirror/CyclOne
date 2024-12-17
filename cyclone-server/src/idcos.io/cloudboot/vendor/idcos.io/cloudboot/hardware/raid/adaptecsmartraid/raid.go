package adaptecsmartraid

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"idcos.io/cloudboot/hardware"
	"idcos.io/cloudboot/hardware/raid"
	"idcos.io/cloudboot/logger"
)

const (
	// name RAID处理器名称
	name = raid.AdaptecSmartRAID
	// tool 硬件配置工具
	tool = "/usr/sbin/arcconf"
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
func (worker *worker) Clear(ctrlID string) (err error) {
	// arcconf delete 1 logicaldrive all noprompt
	_, err = worker.Base.ExecByShell(tool, "delete", ctrlID, "logicaldrive", "all", "noprompt")
	worker.Base.Sleep()
	return nil // 当不存在任何LD时,command aborted,属于err，故不作异常处理
}

// InitLogicalDrives 初始化指定RAID卡控制器下的逻辑驱动器(逻辑磁盘)。
func (worker *worker) InitLogicalDrives(ctrlID string) error {
	// archonf create时默认自动格式化drvie，此处仅设置启用write cache
	// arcconf setcache 1 drivewritecachepolicy Configured 1
	// 将state=Configured的硬盘Enable Write Cache功能
	_, _ = worker.Base.ExecByShell(tool, "setcache", ctrlID, "drivewritecachepolicy", "Configured", "1")
	worker.Base.Sleep()
	return nil // when already set, command aborted,属于err，故不作异常处理
}

// PhysicalDrives 返回物理驱动器(物理磁盘)列表。
func (worker *worker) PhysicalDrives(ctrlID string) (pds []raid.PhysicalDrive, err error) {
	items, err := worker.physicalDrives(ctrlID)
	if err != nil {
		return nil, err
	}
	for i := range items {
		pds = append(pds, raid.PhysicalDrive{
			ID:           worker.extractValue(items[i].PhysicalID, commaSeparator),
			Name:         items[i].PhysicalID,
			RawSize:      items[i].SizeMB,
			MediaType:    items[i].Type,
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
func (worker *worker) CreateArray(ctrlID string, level raid.Level, drives []string) (err error) {
	if err = worker.validateLevel(level); err != nil {
		return err
	}

	// 设置整个RAID卡控制器的模式（开启RAID模式）
	// 当前已是RAID模式时，arcconf命令返回command aborted,属于err，故不作判断
	_ = worker.switchControllerMode(ctrlID, modeRAID)

	_, _ = worker.Base.ExecByShell(tool, "list", ctrlID)

	if level == raid.RAID0 {
		// arcconf create 1 logicaldrive method default ldcache lon max raidlevel physical_id
		_, err = worker.Base.ExecByShell(tool, "create", ctrlID, "logicaldrive",
			"method",
			"default",
			"ldcache",
			"lon",
			"max", // size
			"0", //raidlevel
			fmt.Sprintf("%s", strings.Join(drives, " ")),
			"noprompt",
		)
	} else if level == raid.RAID1 {
		_, err = worker.Base.ExecByShell(tool, "create", ctrlID, "logicaldrive",
			"method",
			"default",
			"ldcache",
			"lon",
			"max", // size
			"1", //raidlevel
			fmt.Sprintf("%s", strings.Join(drives, " ")),
			"noprompt",
		)
	} else {
		err = fmt.Errorf("RAID LEVEL %s Not Supported yet.", level)
	}
	worker.Base.Sleep()
	_, _ = worker.Base.ExecByShell(tool, "list", ctrlID)
	_, _ = worker.Base.ExecByShell(tool, "getconfig", ctrlID, "LD")
	return err
}

// SetGlobalHotspares 添加热备盘
func (worker *worker) SetGlobalHotspares(ctrlID string, drives []string) (err error) {
	// arcconf /cx[/ex]/sx add hotsparedrive [{dgs=<n|0,1,2...>}] [enclaffinity][nonrevertible]
	for i := range drives {
		if _, err = worker.Base.ExecByShell(tool, "setstate", ctrlID, "device", drives[i], "hsp", "sparetype", "2", "noprompt"); err != nil {
			return err
		}
	}
	worker.Base.Sleep()
	return nil
}

// SetJBODs 将指定RAID卡控制器设置为直通模式，或者将RAID卡控制器下部分物理驱动器设置为直通模式。
// ctrlID 表示RAID卡控制器ID，通过Convert2ControllerID方法可获取到第N块RAID卡的控制器ID。
// drives 物理驱动器列表。若物理驱动器列表为空，则意味着将指定的RAID卡下所有的物理驱动器都设置为直通模式。
func (worker *worker) SetJBODs(ctrlID string, drives []string) (err error) {
	defer worker.Base.Sleep()

	if len(drives) <= 0 {
		// 为整个RAID卡控制器下的物理驱动器开启直通模式
		if err = worker.switchControllerMode(ctrlID, modeMIXED); err != nil {
			return err
		}
	}
	// 为指定RAID卡控制器下的部分物理驱动器开启直通模式
	// arcconf /cx[/ex]/sx set jbod
	//for i := range drives {
	//	if _, err = worker.Base.ExecByShell(tool, worker.genPDIdentity(ctrlID, drives[i]), "set", "jbod"); err != nil {
	//		return err
	//	}
	//}
	return nil
}

// PostCheck RAID配置实施后置检查
func (worker *worker) PostCheck(sett *raid.Setting) (items []hardware.CheckingItem) {
	for i := range sett.Controllers {
		ctrlID, err := worker.Convert2ControllerID(sett.Controllers[i].Index)
		if err != nil {
			items = append(items, hardware.CheckingItem{
				Title:   "获取RAID卡控制器ID",
				Matched: hardware.MatchedUnknown,
				Error:   err.Error(),
			})
			continue
		}
		items = append(items, worker.checkArrays(ctrlID, sett.Controllers[i].Arrays)...)

		//if item := worker.checkGlobalHotSpares(ctrlID, &sett.Controllers[i]); item != nil {
		//	items = append(items, *item)
		//}
	}
	return items
}

// 返回 worker 所使用的cmdline
func (worker *worker) GetCMDLine() string {
	return tool
}
