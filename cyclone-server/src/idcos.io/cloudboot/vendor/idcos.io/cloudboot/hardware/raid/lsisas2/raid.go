package lsisas2

import (
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
	name = raid.LSISAS2
	// tool 硬件配置工具
	tool = "sas2ircu"
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

// Controllers 返回设备的RAID卡控制器列表。
func (worker *worker) Controllers() (ctrls []raid.Controller, err error) {
	ids, err := worker.controllerIDs()
	if err != nil {
		return nil, err
	}
	ctrls = make([]raid.Controller, 0, len(ids))
	for i := range ids {
		ctrl, err := worker.findController(ids[i])
		if err != nil {
			return nil, err
		}
		ctrls = append(ctrls, *ctrl)
	}
	return ctrls, nil
}

// PhysicalDrives 返回指定RAID卡控制器的物理驱动器(物理磁盘)列表。
// ctrlID 表示RAID卡控制器ID，通过Convert2ControllerID方法可获取到第N块RAID卡的控制器ID。
func (worker *worker) PhysicalDrives(ctrlID string) (pds []raid.PhysicalDrive, err error) {
	items, err := worker.physicalDrives(ctrlID)
	if err != nil {
		return nil, err
	}
	pds = make([]raid.PhysicalDrive, 0, len(items))
	for i := range items {
		pds = append(pds, raid.PhysicalDrive{
			ID:           strconv.Itoa(items[i].Slot),
			Name:         items[i].identity(),
			RawSize:      items[i].Size,
			MediaType:    items[i].DriveType,
			ControllerID: ctrlID,
		})
	}
	return pds, nil
}

// Clear 擦除指定RAID卡控制器的配置。
// ctrlID 表示RAID卡控制器ID，通过Convert2ControllerID方法可获取到第N块RAID卡的控制器ID。
func (worker *worker) Clear(ctrlID string) (err error) {
	_, err = worker.ExecByShell(tool, ctrlID, "delete", "noprompt")
	worker.Sleep()
	return err
}

// InitLogicalDrives 初始化指定RAID卡控制器下的逻辑驱动器(逻辑磁盘)。
func (worker *worker) InitLogicalDrives(ctrlID string) error {
	return raid.ErrInitLogicalDrivesNotSupport
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

// CreateArray 在指定RAID卡控制器下的若干物理驱动器上创建指定级别的阵列(逻辑磁盘)。
// ctrlID 表示RAID卡控制器ID，通过Convert2ControllerID方法可获取到第N块RAID卡的控制器ID。
// level RAID级别，不同的RAID卡所支持的级别有所差异。
// drives 物理驱动器列表。PhysicalDrives方法可获得的物理驱动器信息列表，其中'Name'属性即为某个物理驱动器标识。
func (worker *worker) CreateArray(ctrlID string, level raid.Level, drives []string) (err error) {
	if level != raid.RAID0 && level != raid.RAID1 && level != raid.RAID1E && level != raid.RAID10 {
		return raid.ErrUnsupportedLevel
	}
	args := []string{ctrlID, "create", string(level), "max"}
	args = append(args, drives...)
	args = append(args, "noprompt")
	_, err = worker.ExecByShell(tool, args...)
	worker.Sleep()
	return err
}

// SetGlobalHotspares 将指定RAID卡下的物理驱动器设置为全局热备盘。
// ctrlID 表示RAID卡控制器ID，通过Convert2ControllerID方法可获取到第N块RAID卡的控制器ID。
func (worker *worker) SetGlobalHotspares(ctrlID string, drives []string) (err error) {
	for i := range drives {
		if err = worker.hotspare(ctrlID, drives[i], false); err != nil {
			return err
		}
	}
	worker.Sleep()
	return nil
}

// SetJBODs 将指定RAID卡控制器设置为直通模式，或者将RAID卡控制器下部分物理驱动器设置为直通模式。
// ctrlID 表示RAID卡控制器ID，通过Convert2ControllerID方法可获取到第N块RAID卡的控制器ID。
// drives 物理驱动器列表。该RAID卡不支持将部分物理驱动器设置为直通模式，因此该参数并不会有实际作用。
func (worker *worker) SetJBODs(ctrlID string, drives []string) (err error) {
	return raid.ErrJBODModeNotSupport
}

// PostCheck RAID配置实施后置检查
func (worker *worker) PostCheck(sett *raid.Setting) (items []hardware.CheckingItem) {
	for i := range sett.Controllers {
		ctrlID, err := worker.Convert2ControllerID(sett.Controllers[i].Index)
		if err != nil {
			items = append(items, hardware.CheckingItem{
				Title:   "Get controller ID",
				Matched: hardware.MatchedUnknown,
				Error:   err.Error(),
			})
			continue
		}
		items = append(items, worker.checkArrays(ctrlID, sett.Controllers[i].Arrays)...)

		if item := worker.checkGlobalHotSpares(ctrlID, &sett.Controllers[i]); item != nil {
			items = append(items, *item)
		}
	}
	return items
}

// 返回 worker 所使用的cmdline
func (worker *worker) GetCMDLine() string {
	return tool
}