package hpsmartarray

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"idcos.io/cloudboot/hardware"
	"idcos.io/cloudboot/hardware/raid"
	"idcos.io/cloudboot/utils/collection"
)

const (
	// eqSeparator 等号分隔符
	eqSeparator = "="
	// colonSeparator 冒号分隔符
	colonSeparator = ":"
)

// extractValue 截取kv对中v的内容。假设，kv内容为"name : voidint"，那么将返回"voidint"。
func (worker *worker) extractValue(kv, sep string) (value string) {
	if !strings.Contains(kv, sep) {
		return kv
	}
	return strings.TrimSpace(strings.SplitN(kv, sep, 2)[1])
}

// validateLevel 校验RAID级别
func (worker *worker) validateLevel(level raid.Level) error {
	if level != raid.RAID0 &&
		level != raid.RAID1 &&
		level != raid.RAID10 &&
		level != raid.RAID5 &&
		level != raid.RAID50 &&
		level != raid.RAID6 &&
		level != raid.RAID60 {
		return raid.ErrUnsupportedLevel
	}
	return nil
}

// strLevel 返回字符串形式的RAID级别
func (worker *worker) strLevel(level raid.Level) (l string) {
	switch level {
	case raid.RAID0:
		return "0"
	case raid.RAID1:
		return "1"
	case raid.RAID10:
		return "1+0"
	case raid.RAID5:
		return "5"
	case raid.RAID50:
		return "50"
	case raid.RAID6:
		return "6"
	case raid.RAID60:
		return "60"
	}
	panic("unreachable")
}

const (
	raidmode = "raidmode"
	hbamode  = "hbamode"
)

func (worker *worker) setManagingMode(ctrlID string, mode string) (err error) {
	if mode != raidmode && mode != hbamode {
		return raid.ErrUnsupportedRAIDManagingMode
	}
	ctrl, err := worker.findController(ctrlID)
	if err != nil {
		return err
	}
	if (mode == raidmode && ctrl.Mode == "RAID") || (mode == hbamode && ctrl.Mode != "RAID") {
		return nil
	}
	// hpssacli ctrl slot=%s modify raidmode=on forced
	_, err = worker.ExecByShell(tool, "ctrl", fmt.Sprintf("slot=%s", ctrlID), "modify", fmt.Sprintf("%s=on", mode), "forced")
	return err
}

// hotspare 添加热备盘
func (worker *worker) hotspare(target string, drives []string) (err error) {
	// hpssacli ctrl slot=%s array all add spares=%s
	_, err = worker.Base.ExecByShell(tool, target, "add", fmt.Sprintf("spares=%s", strings.Join(drives, ",")))
	return err
}

// controllerIDs 返回当前设备所有RAID卡的ID列表
func (worker *worker) controllerIDs() (items []string, err error) {
	// hpssacli ctrl all show
	output, err := worker.Base.ExecByShell(tool, "ctrl", "all", "show")
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		begin := strings.Index(line, "Slot")
		end := strings.Index(line, "(")
		if !(strings.HasPrefix(line, "Smart Array") && begin > 0 && end > 0 && begin < end) {
			continue
		}
		items = append(items, strings.TrimSpace(line[begin+4:end]))
	}
	return items, scanner.Err()
}

// findController 返回指定编号的Controller
func (worker *worker) findController(ctrlID string) (*raid.Controller, error) {
	// hpssacli ctrl slot=1 show
	output, err := worker.Base.ExecByShell(tool, "ctrl", fmt.Sprintf("slot=%s", ctrlID), "show")
	if err != nil {
		return nil, err
	}

	ctrl := raid.Controller{
		ID: ctrlID,
	}

	keyword := fmt.Sprintf("in Slot %s", ctrlID)
	scanner := bufio.NewScanner(bytes.NewBuffer(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if idx := strings.Index(line, keyword); idx > 0 {
			ctrl.ModelName = fmt.Sprintf("HP %s", strings.TrimSpace(line[:idx]))
		}
		if strings.HasPrefix(line, "Firmware Version:") {
			ctrl.FirmwareVersion = worker.extractValue(line, colonSeparator)
		}
		if strings.HasPrefix(line, "Controller Mode:") {
			ctrl.Mode = worker.extractValue(line, colonSeparator)
		}
	}
	return &ctrl, scanner.Err()
}

var (
	arrLineReg        = regexp.MustCompile(`^array [A-Z]$`)
	unassignedLineReg = regexp.MustCompile(`^unassigned$`)
)

type physicalDrive struct {
	Array         string // 所属的阵列
	Port          string
	Box           string
	Bay           string
	InterfaceType string
	Size          string
	Spare         bool // 是否是备份盘
}

func (pd *physicalDrive) identity() string {
	return fmt.Sprintf("%s:%s:%s", pd.Port, pd.Box, pd.Bay)
}

// findPhysicalDriveIDByIndex 通过切片索引查找物理驱动器ID
func findPhysicalDriveIDByIndex(items []physicalDrive, index int) (id string, err error) {
	if index < 0 || index > len(items)-1 {
		return "", raid.ErrInvalidDiskIdentity
	}
	return items[index].identity(), nil
}

// findPhysicalDriveIDsByIndexes 通过切片索引列表查找物理驱动器ID
func findPhysicalDriveIDsByIndexes(items []physicalDrive, indexes ...int) (ids []string, err error) {
	if len(indexes) <= 0 {
		for i := range items {
			ids = append(ids, items[i].identity())
		}
		return ids, nil
	}
	ids = make([]string, 0, len(indexes))
	for i := range indexes {
		id, err := findPhysicalDriveIDByIndex(items, indexes[i])
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func fromIdentity(id string) (port, box, bay string) {
	if fields := strings.Split(id, ":"); len(fields) == 3 {
		return fields[0], fields[1], fields[2]
	}
	return "", "", ""
}

// physicalDrives 返回指定slot下的物理驱动器列表
func (worker *worker) physicalDrives(ctrlID string) (items []physicalDrive, err error) {
	// hpssacli ctrl slot=1 pd all show
	output, err := worker.Base.ExecByShell(tool, "ctrl", fmt.Sprintf("slot=%s", ctrlID), "pd", "all", "show")
	if err != nil {
		return nil, err
	}

	var started bool
	var array string

	scanner := bufio.NewScanner(bytes.NewBuffer(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !started && (arrLineReg.MatchString(line) || unassignedLineReg.MatchString(line)) {
			started = true
		}
		if arrLineReg.MatchString(line) || unassignedLineReg.MatchString(line) {
			array = line
		}

		begin := strings.Index(line, "(")
		end := strings.Index(line, ")")
		if !(strings.HasPrefix(line, "physicaldrive") && begin > 0 && end > 0 && begin < end) {
			continue
		}

		arr := strings.Split(line[begin+1:end], ",")
		if len(arr) < 3 {
			continue
		}

		id := strings.TrimSpace(strings.TrimPrefix(line[:begin], "physicaldrive"))

		var pd physicalDrive
		pd.Array = array
		pd.Port, pd.Box, pd.Bay = fromIdentity(id)
		pd.InterfaceType, pd.Size = strings.TrimSpace(arr[1]), strings.TrimSpace(arr[2])
		pd.Spare = len(arr) >= 5 && strings.TrimSpace(arr[4]) == "spare"

		items = append(items, pd)
	}
	return items, scanner.Err()
}

// physicalDriveIDs 根据入参的物理驱动器，返回由物理驱动器ID(port:box:bay)组成的字符串切片。
func physicalDriveIDs(items []physicalDrive) (ids []string) {
	ids = make([]string, 0, len(items))
	for i := range items {
		ids = append(ids, items[i].identity())
	}
	return ids
}

type pdScanFunc func(item *physicalDrive) bool

// scanPDFunc 逐个扫描物理驱动器列表并返回满足所有过滤条件的物理驱动器
func scanPDFunc(pds []physicalDrive, filters ...pdScanFunc) (items []physicalDrive) {
	for i := range pds {
		var unmatched bool
		for j := range filters {
			if !filters[j](&pds[i]) {
				unmatched = true
				break
			}
		}
		if !unmatched {
			items = append(items, pds[i])
		}
	}
	return items
}

// sparePDs 若当前物理驱动器是用于热备，则返回true，反之返回false。
func sparePDs(pd *physicalDrive) bool {
	return pd != nil && pd.Spare
}

// notSparePDs 若当前物理驱动器是用于热备，则返回false，反之返回true。
func notSparePDs(pd *physicalDrive) bool {
	return pd != nil && !pd.Spare
}

// arrayedPDs 若当前物理驱动器已经分配给某个阵列，则返回true，反之返回false。
func arrayedPDs(pd *physicalDrive) bool {
	return pd != nil && pd.Array != "" && pd.Array != "unassigned"
}

// pdGroupByArray 将物理驱动器以阵列进行分组
// func pdGroupByArray(items []physicalDrive) map[string][]physicalDrive {
// 	store := make(map[string][]physicalDrive)
// 	for i := range items {
// 		if _, exist := store[items[i].Array]; exist {
// 			store[items[i].Array] = append(store[items[i].Array], items[i])
// 		} else {
// 			store[items[i].Array] = []physicalDrive{items[i]}
// 		}
// 	}
// 	return store
// }

type logicalDrive struct {
	Array     string // 所属的阵列
	Title     string
	Size      string
	RAIDLevel string
}

// logicalDrives 返回指定RAID卡控制器下的逻辑驱动器列表
func (worker *worker) logicalDrives(ctrlID string) (items []logicalDrive, err error) {
	// hpssacli ctrl slot=1 ld all show
	output, err := worker.Base.ExecByShell(tool, "ctrl", fmt.Sprintf("slot=%s", ctrlID), "ld", "all", "show")
	if err != nil {
		return nil, nil // 若实际不存在逻辑驱动器，执行以上命令也会得到非0的退出码。
	}

	var started bool
	var array string

	scanner := bufio.NewScanner(bytes.NewBuffer(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !started && arrLineReg.MatchString(line) {
			started = true
		}
		if arrLineReg.MatchString(line) {
			array = line
		}
		begin := strings.Index(line, "(")
		end := strings.Index(line, ")")
		if !(strings.HasPrefix(line, "logicaldrive") && begin > 0 && end > 0 && begin < end) {
			continue
		}
		arr := strings.Split(line[begin+1:end], ",")
		if len(arr) < 2 {
			continue
		}
		items = append(items, logicalDrive{
			Array:     array,
			Title:     strings.TrimSpace(line[:begin-1]),
			Size:      strings.TrimSpace(arr[0]),
			RAIDLevel: strings.TrimSpace(arr[1]),
		})
	}
	return items, scanner.Err()
}

// logicalDrivesExist 阵列(逻辑磁盘)是否存在
func (worker *worker) logicalDrivesExist(ctrlID string) bool {
	lds, _ := worker.logicalDrives(ctrlID)
	return len(lds) > 0
}

func (worker *worker) findArraysByLevel(ctrlID string, level raid.Level) (arrays []string, err error) {
	lds, err := worker.logicalDrives(ctrlID)
	if err != nil {
		return nil, err
	}
	set := collection.NewSSet(1)
	for i := range lds {
		if lds[i].RAIDLevel == fmt.Sprintf("RAID %s", worker.strLevel(level)) {
			set.Add(lds[i].Array)
		}
	}
	return set.Elements(), nil
}

// checkArrays 检查实际的RAID阵列是否与预期配置一致
func (worker *worker) checkArrays(ctrlID string, arraySetts []raid.ArraySetting) (items []hardware.CheckingItem) {
	if ctrlID == "" || len(arraySetts) == 0 {
		return nil
	}
	for i := range arraySetts {
		if item := worker.checkArray(ctrlID, &arraySetts[i]); item != nil {
			items = append(items, *item)
		}
	}
	return items
}

// checkArray 检查指定的单个RAID阵列配置与实际是否一致
func (worker *worker) checkArray(ctrlID string, sett *raid.ArraySetting) *hardware.CheckingItem {
	if sett == nil {
		return nil
	}
	item := hardware.CheckingItem{
		Title:    fmt.Sprintf("[Controller%s] Create array", ctrlID),
		Expected: fmt.Sprintf("%s@%s", strings.ToUpper(sett.Level), sett.Drives),
		Matched:  hardware.MatchedUnknown,
	}

	allpds, err := worker.physicalDrives(ctrlID)
	if err != nil {
		item.Error = err.Error()
		return &item
	}

	level, err := raid.SelectLevel(sett.Level)
	if err != nil {
		item.Error = err.Error()
		return &item
	}

	arrays, err := worker.findArraysByLevel(ctrlID, level)
	if err != nil {
		item.Error = err.Error()
		return &item
	}

	drives := strings.Split(sett.Drives, raid.Sep) // 预期的物理驱动器列表

	switch len(arrays) {
	case 0: // 不存在与预期级别相同的阵列
		item.Matched = hardware.MatchedNO
		item.Actual = "missing"
		return &item

	case 1: // 存在一组与预期级别相同的阵列
		pds := scanPDFunc(allpds,
			func(item *physicalDrive) bool {
				return item != nil && item.Array == arrays[0]
			},
			notSparePDs,
		) // 过滤获得指定阵列中实际的非热备物理驱动器列表

		item.Actual = fmt.Sprintf("%s@%s",
			strings.ToUpper(sett.Level),
			strings.Join(physicalDriveIDs(pds), raid.Sep),
		)

		if len(pds) != len(drives) {
			item.Matched = hardware.MatchedNO
			return &item
		}

		for i := range pds {
			if !collection.InSlice(pds[i].identity(), drives) {
				item.Matched = hardware.MatchedNO
				return &item
			}
		}
		item.Matched = hardware.MatchedYES
		return &item

	default: // 存在多组与预期级别相同的阵列
		// TODO 待实现

		return &item
	}

}

// checkGlobalHotSpares 检查实际的全局热备盘设置是否与预期配置一致
func (worker *worker) checkGlobalHotSpares(ctrlID string, sett *raid.ControllerSetting) *hardware.CheckingItem {
	if ctrlID == "" || sett == nil || sett.Hotspares == "" {
		return nil
	}
	item := hardware.CheckingItem{
		Title:    fmt.Sprintf("[Controller%s] Set global hot spare(s)", ctrlID),
		Expected: sett.Hotspares,
		Matched:  hardware.MatchedUnknown,
	}

	pds, err := worker.physicalDrives(ctrlID)
	if err != nil {
		item.Error = err.Error()
		return &item
	}

	// 扫描查找所有已经分配作为热备的物理驱动器
	sparePDs := scanPDFunc(pds, sparePDs)
	item.Actual = strings.Join(physicalDriveIDs(sparePDs), raid.Sep)

	spares := strings.Split(sett.Hotspares, raid.Sep) // 预期的用于备份的物理驱动器列表
	if len(sparePDs) != len(spares) {
		item.Matched = hardware.MatchedNO
		return &item
	}

	for i := range sparePDs {
		if !collection.InSlice(sparePDs[i].identity(), spares) {
			item.Matched = hardware.MatchedNO
			return &item
		}
	}

	item.Matched = hardware.MatchedYES
	return &item
}
