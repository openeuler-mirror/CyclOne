package adaptecsmartraid

import (
	"bufio"
	"bytes"
	"fmt"
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
	// commaSeparator 逗号分隔符
	commaSeparator = ","
)

// extractValue 截取kv对中v的内容。假设，kv内容为"name : voidint"，那么将返回"voidint"。
func (worker *worker) extractValue(kv, sep string) (value string) {
	if !strings.Contains(kv, sep) {
		return kv
	}
	return strings.TrimSpace(strings.SplitN(kv, sep, 2)[1])
}

const (
	// RAID模式仅发送LD至OS
	modeRAID = "3"
	// Mixed模式发送PD\LD至OS
	modeMIXED = "5"
)

// switchControllerMode
func (worker *worker) switchControllerMode(ctrlID, modeSet string) (err error) {
	_, err = worker.Base.ExecByShell(tool, "setcontrollermode", ctrlID, modeSet, "noprompt")
	return nil //已设置为RAID模式时，command aborted.
}

// validateLevel 校验RAID级别
func (worker *worker) validateLevel(level raid.Level) error {
	if level != raid.RAID0 && level != raid.RAID1 && level != raid.RAID10 &&
		level != raid.RAID5 && level != raid.RAID50 && level != raid.RAID6 && level != raid.RAID60 {
		return raid.ErrUnsupportedLevel
	}
	return nil
}


// controllerIDs 返回当前设备所有RAID卡的ID列表
func (worker *worker) controllerIDs() (items []string, err error) {
	// arcconf list
	output, err := worker.Base.ExecByShell(tool, "list")
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.Contains(line, "Controller") && strings.Contains(line, "Optimal") && strings.Contains(line, "Slot"){
			fields := strings.Fields(line)
			// Controller 1:
			if strings.Contains(fields[0], "Controller") {
				ctlid := strings.TrimRight(fields[1], colonSeparator)
				items = append(items, ctlid)
			}
		}
	}
	return items, scanner.Err()
}

// findController 返回指定编号的Controller
func (worker *worker) findController(ctrlID string) (*raid.Controller, error) {
	// arcconf getconfig 1 AD 获取整个适配器信息
	output, err := worker.Base.ExecByShell(tool, "getconfig", fmt.Sprintf("%s", ctrlID), "AD")
	if err != nil {
		return nil, err
	}

	ctl := raid.Controller{
		ID: ctrlID,
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "Firmware") && strings.Contains(line, colonSeparator) {
			ctl.FirmwareVersion = worker.extractValue(line, colonSeparator)
			break
		}
	}
	ctl.ModelName, _ = worker.findModelName(ctrlID)
	return &ctl, scanner.Err()
}

// findModelName 查询指定RAID卡控制器的型号
func (worker *worker) findModelName(ctrlID string) (string, error) {
	// arcconf getconfig 1 AD
	output, err := worker.Base.ExecByShell(tool, "getconfig", fmt.Sprintf("%s", ctrlID), "AD")
	if err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "Controller Model") && strings.Contains(line, colonSeparator) {
			return worker.extractValue(line, colonSeparator), nil
		}
	}
	return "", scanner.Err()
}

// drivesEq 对比两个驱动器列表是否相等(顺序无关)
func (worker *worker) drivesEq(drives0, drives1 []string) bool {
	if len(drives0) != len(drives1) {
		return false
	}
	for i := range drives0 {
		if !collection.InSlice(drives0[i], drives1) {
			return false
		}
	}
	return true
}

// physicalDrive 物理驱动器
type physicalDrive struct {
	PhysicalID  string // Channel Device
	State 		string
	SizeMB      string // MB
	Type        string // SSD
}

// findPhysicalDriveIDByIndex 通过切片索引查找物理驱动器ID
func findPhysicalDriveIDByIndex(items []physicalDrive, index int) (id string, err error) {
	if index < 0 || index > len(items)-1 {
		return "", raid.ErrInvalidDiskIdentity
	}
	return items[index].PhysicalID, nil
}

// findPhysicalDriveIDsByIndexes 通过切片索引列表查找物理驱动器ID
func findPhysicalDriveIDsByIndexes(items []physicalDrive, indexes ...int) (ids []string, err error) {
	if len(indexes) <= 0 {
		for i := range items {
			ids = append(ids, items[i].PhysicalID)
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

// physicalDriveIDs 根据入参的物理驱动器，返回由物理驱动器ID(Enclosure:Bay)组成的字符串切片。
func physicalDriveIDs(items []physicalDrive) (ids []string) {
	ids = make([]string, 0, len(items))
	for i := range items {
		ids = append(ids, items[i].PhysicalID)
	}
	return ids
}

type pdScanFunc func(item *physicalDrive) bool

// scanPDs 逐个扫描物理驱动器列表并返回满足所有过滤条件的物理驱动器
//func scanPDs(pds []physicalDrive, filters ...pdScanFunc) (items []physicalDrive) {
//	for i := range pds {
//		var unmatched bool
//		for j := range filters {
//			if !filters[j](&pds[i]) {
//				unmatched = true
//				break
//			}
//		}
//		if !unmatched {
//			items = append(items, pds[i])
//		}
//	}
//	return items
//}

// sparePDs 若当前物理驱动器是用于热备，则返回true，反之返回false。
//func sparePDs(pd *physicalDrive) bool {
//	return pd != nil && strings.Contains(pd.State, "GHS")
//}

// parsePDs 解析内容中包含的物理驱动器列表
func parsePDs(output []byte) (list []physicalDrive, err error) {
	var started bool
	var hd physicalDrive
	scanner := bufio.NewScanner(bytes.NewBuffer(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, "Physical") && strings.Contains(line, "ID") {
			started = true
			continue
		}
		if !started || line == "" {
			continue
		}
		if strings.Contains(line, "Physical") && strings.Contains(line, "Enclosure") && strings.Contains(line, "Slot") {
			fields := strings.Fields(line)
			if  len(fields) > 11 {
				hd.PhysicalID = strings.ReplaceAll(fields[1], commaSeparator, " ") // 0,8 (Channel,Device) -> 0 8
				hd.State = fields[3] //Ready
				hd.SizeMB = strings.TrimRight(fields[7], commaSeparator) // 1000MB
				switch fields[10] {
				case "Hard":
					hd.Type = "Hard Drive"
				case "Solid":
					hd.Type = "Solid State Drive"
				default:
					hd.Type = "Unknown"
				}
				list = append(list, hd)
			}
		}
	}
	return
}

// physicalDrives 返回指定controller下的物理驱动器列表
func (worker *worker) physicalDrives(ctrlID string) (items []physicalDrive, err error) {
	//arcconf list 1
	output, err := worker.Base.ExecByShell(tool, "list", fmt.Sprintf("%s", ctrlID))
	if err != nil {
		return nil, err
	}
	return parsePDs(output)
}

// logicalDrive 逻辑驱动器(阵列)
type logicalDrive struct {
	ID             string
	Name           string
	Type           string  //RAID0
	State          string
	Size           string
	PhysicalDrives []physicalDrive // 逻辑驱动器底层所依赖的物理驱动器
}

type ldScanFunc func(item *logicalDrive) bool

// scanLDs 逐个扫描逻辑驱动器列表并返回满足所有过滤条件的逻辑驱动器
func scanLDs(lds []logicalDrive, filters ...ldScanFunc) (items []logicalDrive) {
	for i := range lds {
		var unmatched bool
		for j := range filters {
			if !filters[j](&lds[i]) {
				unmatched = true
				break
			}
		}
		if !unmatched {
			items = append(items, lds[i])
		}
	}
	return items
}

// logicalDrives 返回指定controller下的LD列表
func (worker *worker) logicalDrives(ctrlID string) (items []logicalDrive, err error) {
	// arcconf getconfig 1 LD
	output, err := worker.Base.ExecByShell(tool, "getconfig", ctrlID, "LD")
	if err != nil {
		return nil, err
	}
	var ld logicalDrive
	scanner := bufio.NewScanner(bytes.NewBuffer(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.Contains(line, "Logical Device number") {
			if ld.ID != "" {
				items = append(items, ld)
			}
			ld = logicalDrive{}
			ld.ID = strings.TrimRight(line, "\n")
			continue
		}
		if strings.Contains(line, "Disk Name") && strings.Contains(line, colonSeparator){
			ld.Name = worker.extractValue(line, eqSeparator) // 形如 /dev/sda
		} else if strings.HasPrefix(line, "RAID level") && strings.Contains(line, colonSeparator){
			ld.Type = "RAID" + worker.extractValue(line, eqSeparator)  // 形如 0、1、5、50
		} else if strings.Contains(line, "Status of Logical Device") && strings.Contains(line, colonSeparator){
			ld.State = worker.extractValue(line, eqSeparator) // 形如 Optimal
		} else if strings.HasPrefix(line, "Size") && strings.Contains(line, colonSeparator){
			ld.Size = worker.extractValue(line, eqSeparator) // 形如 300 MB
		} else if strings.Contains(line, "Enclosure") {
			fields := strings.Fields(line)
			pd := physicalDrive {
				PhysicalID: "0 " + fields[1],
				State: "Online",
				SizeMB: strings.TrimRight(strings.TrimPrefix(fields[4], "("), ","),
				Type: strings.TrimRight(fields[6], ","),
			}
			ld.PhysicalDrives = append(ld.PhysicalDrives, pd)
		}
	}

	return items, nil
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
	if sett == nil || (sett.Drives == "" && sett.Level == "") {
		return nil
	}
	item := hardware.CheckingItem{
		Title:    fmt.Sprintf("[Controller%s] 创建阵列", ctrlID),
		Expected: fmt.Sprintf("%s@%s", strings.ToUpper(sett.Level), sett.Drives),
		Matched:  hardware.MatchedUnknown,
	}

	level, err := raid.SelectLevel(sett.Level)
	if err != nil {
		item.Error = err.Error()
		return &item
	}

	// 获取该RAID卡下已有的阵列列表
	allLDs, err := worker.logicalDrives(ctrlID)
	if err != nil {
		item.Error = err.Error()
		return &item
	}

	// 过滤获得指定级别的阵列列表
	lds := scanLDs(allLDs, func(item *logicalDrive) bool {
		return item != nil && strings.ToLower(item.Type) == strings.ToLower(string(level))
	})

	if len(lds) <= 0 {
		// 不存在与预期级别相同的阵列
		item.Matched = hardware.MatchedNO
		item.Actual = fmt.Sprintf("缺失%s", strings.ToUpper(sett.Level))
		return &item
	}

	// 预期的物理驱动器列表
	drives := strings.Split(sett.Drives, raid.Sep)

	var pds4ld []string
	// 遍历该RAID级别的阵列，查找是否存在一个阵列所使用的物理驱动器与预期物理驱动器一致。
	for i := range lds {
		// 当前阵列实际所使用的物理驱动器列表
		pdids := physicalDriveIDs(lds[i].PhysicalDrives)
		pds4ld = append(pds4ld, strings.Join(pdids, raid.Sep))

		// 比较实际使用的物理驱动器列表与预期物理驱动器列表是否一致
		if worker.drivesEq(pdids, drives) {
			item.Actual = fmt.Sprintf("%s@%s", strings.ToUpper(sett.Level), sett.Drives)
			item.Matched = hardware.MatchedYES
			return &item
		}
	}
	item.Actual = fmt.Sprintf("%s@%s", strings.ToUpper(sett.Level), strings.Join(pds4ld, "|"))
	item.Matched = hardware.MatchedNO
	return &item
}

// checkGlobalHotSpares 检查单个RAID卡控制器的实际全局热备盘设置是否与预期配置一致
//func (worker *worker) checkGlobalHotSpares(ctrlID string, sett *raid.ControllerSetting) *hardware.CheckingItem {
//	if sett == nil || sett.Hotspares == "" {
//		return nil
//	}
//	item := hardware.CheckingItem{
//		Title:    fmt.Sprintf("[Controller%s] 设置全局热备盘", ctrlID),
//		Expected: sett.Hotspares,
//		Matched:  hardware.MatchedUnknown,
//	}
//
//	pds, err := worker.physicalDrives(ctrlID)
//	if err != nil {
//		item.Error = err.Error()
//		return &item
//	}
//
//	// 扫描查找所有已经分配作为热备的物理驱动器
//	sparePDs := scanPDs(pds, sparePDs)
//	item.Actual = strings.Join(physicalDriveIDs(sparePDs), raid.Sep)
//
//	// 预期的用于备份的物理驱动器列表
//	spares := strings.Split(sett.Hotspares, raid.Sep)
//
//	// 对比实际与预期的热备盘列表是否一致
//	if worker.drivesEq(physicalDriveIDs(sparePDs), spares) {
//		item.Matched = hardware.MatchedYES
//	} else {
//		item.Matched = hardware.MatchedNO
//	}
//	return &item
//}
