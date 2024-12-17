package lsisas3

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"idcos.io/cloudboot/hardware"
	"idcos.io/cloudboot/hardware/raid"
	"idcos.io/cloudboot/utils/collection"
)

const (
	// colonSeparator 分隔符
	colonSeparator = ":"
)

// extractValue 截取kv对中v的内容。假设，kv内容为"name : voidint"，那么将返回"voidint"。
func (worker *worker) extractValue(kv, sep string) (value string) {
	if !strings.Contains(kv, sep) {
		return kv
	}
	return strings.TrimSpace(strings.SplitN(kv, sep, 2)[1])
}

// hotspare 添加/删除热备盘
func (worker *worker) hotspare(ctrlID, drive string, delete bool) error {
	args := []string{ctrlID, "hotspare"}
	if delete {
		args = append(args, "delete")
	}
	args = append(args, drive)
	_, err := worker.Base.ExecWithStdinPipe([]string{"YES", "NO"}, tool, args...)
	return err
}

var (
	numReg = regexp.MustCompile(`^\d+$`)
)

// controllerIDs 返回当前设备所有RAID卡的ID列表
func (worker *worker) controllerIDs() (items []string, err error) {
	// sas3ircu list
	output, err := worker.ExecByShell(tool, "list")
	if err != nil {
		return nil, err
	}

	var started bool
	scanner := bufio.NewScanner(bytes.NewBuffer(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "---") {
			started = true
			continue
		}
		if !started || line == "" {
			continue
		}

		arr := strings.Fields(line)
		if !numReg.MatchString(arr[0]) {
			continue
		}
		items = append(items, arr[0])
	}
	return items, scanner.Err()
}

// findController 返回指定索引号的Controller
func (worker *worker) findController(ctrlID string) (*raid.Controller, error) {
	output, err := worker.ExecByShell(tool, ctrlID, "display")
	if err != nil {
		return nil, err
	}

	ctl := raid.Controller{
		ID: ctrlID,
	}

	var started bool
	scanner := bufio.NewScanner(bytes.NewBuffer(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "Controller information") {
			started = true
			continue
		}
		if !started {
			continue
		}

		if strings.HasPrefix(line, "Controller type") && strings.Contains(line, ":") {
			ctl.ModelName = worker.extractValue(line, colonSeparator)
		}
		if strings.HasPrefix(line, "Firmware version") && strings.Contains(line, ":") {
			ctl.FirmwareVersion = worker.extractValue(line, colonSeparator)
		}
		if strings.HasPrefix(line, "IR Volume information") || strings.HasPrefix(line, "Physical device information") {
			break
		}
	}
	return &ctl, nil
}

type physicalDrive struct {
	Enclosure int    // 硬盘盒ID
	Slot      int    // 硬盘盒中的插槽ID
	State     string // 物理驱动器状态
	Size      string // 包含单位
	DriveType string // 物理驱动器类型
}

func (hd physicalDrive) identity() string {
	return fmt.Sprintf("%d:%d", hd.Enclosure, hd.Slot)
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

func (worker *worker) physicalDrives(ctrlID string) (pds []physicalDrive, err error) {
	output, err := worker.Exec(tool, ctrlID, "display")
	if err != nil {
		return nil, err
	}
	return worker.parsePDs(output)
}

var (
	enclReg      = regexp.MustCompile(`Enclosure #\s+:\s+(\d+)`)
	slotReg      = regexp.MustCompile(`Slot #\s+:\s+(\d+)`)
	sizeReg      = regexp.MustCompile(`Size.+:\s+(\d+)/(\d+)`)
	driveTypeReg = regexp.MustCompile(`Drive Type\s+:\s+(\S+)`)
	unitReg      = regexp.MustCompile(`(in [a-zA-z]+)`)
	stateReg     = regexp.MustCompile(`State\s+:\s+(\S+\s+\S+)`)
)

// parsePDs 解析内容中包含的物理驱动器列表
func (worker *worker) parsePDs(output []byte) (items []physicalDrive, err error) {
	arr := strings.Split(string(output), "Device is a")
	for i := 1; i < len(arr); i++ {
		var pd physicalDrive
		if pair := strings.Split(stateReg.FindString(arr[i]), ":"); len(pair) == 2 {
			pd.State = strings.TrimSpace(pair[1])
			if strings.Contains(pd.State, "Standby") { // 非硬盘设备
				continue
			}
		}
		if pair := strings.Split(enclReg.FindString(arr[i]), ":"); len(pair) == 2 {
			if pd.Enclosure, err = strconv.Atoi(strings.TrimSpace(pair[1])); err != nil {
				return nil, err
			}
		}
		if pair := strings.Split(slotReg.FindString(arr[i]), ":"); len(pair) == 2 {
			if pd.Slot, err = strconv.Atoi(strings.TrimSpace(pair[1])); err != nil {
				return nil, err
			}
		}
		if pair := strings.Split(sizeReg.FindString(arr[i]), ":"); len(pair) == 2 {
			if idx := strings.Index(pair[1], "/"); idx > 0 {
				pd.Size = fmt.Sprintf("%s%s",
					strings.TrimSpace(pair[1][:idx]),
					strings.TrimPrefix(unitReg.FindString(pair[0]), "in"),
				)
			} else {
				pd.Size = strings.TrimSpace(pair[1])
			}
		}
		if pair := strings.Split(driveTypeReg.FindString(arr[i]), ":"); len(pair) == 2 {
			pd.DriveType = strings.TrimSpace(pair[1])
		}

		items = append(items, pd)
	}
	return items, nil
}

// physicalDriveIDs 根据入参的物理驱动器，返回由物理驱动器ID(Enclosure:Bay)组成的字符串切片。
func physicalDriveIDs(items []physicalDrive) (ids []string) {
	ids = make([]string, 0, len(items))
	for i := range items {
		ids = append(ids, items[i].identity())
	}
	return ids
}

type pdScanFunc func(item *physicalDrive) bool

// scanPDs 逐个扫描物理驱动器列表并返回满足所有过滤条件的物理驱动器
func scanPDs(pds []physicalDrive, filters ...pdScanFunc) (items []physicalDrive) {
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
	return pd != nil && strings.Contains(pd.State, "HotSpare")
}

type volumeScanFunc func(item *volume) bool

func scanVolumes(volumes []*volume, filters ...volumeScanFunc) (items []*volume) {
	for i := range volumes {
		var unmatched bool
		for j := range filters {
			if !filters[j](volumes[i]) {
				unmatched = true
				break
			}
		}
		if !unmatched {
			items = append(items, volumes[i])
		}
	}
	return items
}

// findVolumesByLevel 查找指定RAID卡下指定级别的卷(RAID阵列/逻辑磁盘)列表。
func (worker *worker) findVolumesByLevel(ctrlID string, level raid.Level) (items []*volume, err error) {
	volumes, err := worker.volumes(ctrlID)
	if err != nil {
		return nil, err
	}
	for i := range volumes {
		if volumes[i].RAIDLevel == strings.ToUpper(string(level)) {
			items = append(items, volumes[i])
		}
	}
	return items, nil
}

// volume 逻辑驱动器(卷/RAID阵列)。
type volume struct {
	VolumeID       string
	Status         string
	Size           string
	RAIDLevel      string
	Boot           string
	PhysicalDrives []string
}

// volumes 返回指定controller索引下的逻辑驱动器列表
func (worker *worker) volumes(ctrlID string) (items []*volume, err error) {
	// sas3ircu 0 display
	output, err := worker.Base.ExecByShell(tool, ctrlID, "display")
	if err != nil {
		return nil, nil
	}
	source := string(output)
	begin := strings.Index(source, "IR Volume information")
	end := strings.Index(source, "Physical device information")
	if begin < 0 {
		return nil, nil
	}
	return worker.parseVolumes(source[begin:end])
}

// parseVolumes 解析命令行输出中包含的卷信息
func (worker *worker) parseVolumes(output string) (items []*volume, err error) {
	blocks := strings.Split(output, "IR volume")
	for i := range blocks {
		if blocks[i] = strings.TrimSpace(blocks[i]); blocks[i] == "" {
			continue
		}
		var item volume

		scanner := bufio.NewScanner(strings.NewReader(blocks[i]))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, "Volume ID") && strings.Contains(line, ":") {
				item.VolumeID = worker.extractValue(line, colonSeparator)

			} else if strings.HasPrefix(line, "Status of volume") && strings.Contains(line, ":") {
				item.Status = worker.extractValue(line, colonSeparator)

			} else if strings.HasPrefix(line, "RAID level") && strings.Contains(line, ":") {
				item.RAIDLevel = worker.extractValue(line, colonSeparator)

			} else if strings.HasPrefix(line, "Size (in MB)") && strings.Contains(line, ":") {
				item.Size = fmt.Sprintf("%s MB", worker.extractValue(line, colonSeparator))

			} else if strings.HasPrefix(line, "Boot") && strings.Contains(line, ":") {
				item.Boot = worker.extractValue(line, colonSeparator)

			} else if strings.Contains(line, "Enclosure#/Slot#") && strings.Contains(line, ":") {
				item.PhysicalDrives = append(item.PhysicalDrives, worker.extractValue(line, colonSeparator))
			}
		}
		if item.VolumeID != "" {
			items = append(items, &item)
		}
		if err = scanner.Err(); err != nil {
			if log := worker.Base.GetLog(); log != nil {
				log.Error(err)
			}
			return nil, err
		}
	}
	return items, nil
}

// checkGlobalHotSpares 检查单个RAID卡控制器的实际全局热备盘设置是否与预期配置一致
func (worker *worker) checkGlobalHotSpares(ctrlID string, sett *raid.ControllerSetting) *hardware.CheckingItem {
	if sett == nil || sett.Hotspares == "" {
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
	sparePDs := scanPDs(pds, sparePDs)
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

	level, err := raid.SelectLevel(sett.Level)
	if err != nil {
		item.Error = err.Error()
		return &item
	}

	arrays, err := worker.findVolumesByLevel(ctrlID, level)
	if err != nil {
		item.Error = err.Error()
		return &item
	}

	drives := strings.Split(sett.Drives, raid.Sep) // 预期的物理驱动器列表

	switch len(arrays) {
	case 0:
		// 1、不存在与预期级别相同的阵列
		item.Matched = hardware.MatchedNO
		item.Actual = fmt.Sprintf("missing %s", strings.ToUpper(sett.Level))
		return &item

	case 1:
		// 2、存在一组与预期级别相同的阵列
		item.Actual = fmt.Sprintf("%s@%s",
			strings.ToUpper(sett.Level),
			strings.Join(arrays[0].PhysicalDrives, raid.Sep),
		)

		if len(arrays[0].PhysicalDrives) != len(drives) {
			item.Matched = hardware.MatchedNO
			return &item
		}

		for _, pd := range arrays[0].PhysicalDrives {
			if !collection.InSlice(pd, drives) {
				item.Matched = hardware.MatchedNO
				return &item
			}
		}
		item.Matched = hardware.MatchedYES
		return &item
	}

	// 3、存在多组与预期级别相同的阵列(该RAID卡存在最多创建2组RAID阵列的限制)
	item.Actual = fmt.Sprintf("%s@%s",
		strings.ToUpper(sett.Level),
		fmt.Sprintf("%s|%s", strings.Join(arrays[0].PhysicalDrives, raid.Sep), strings.Join(arrays[1].PhysicalDrives, raid.Sep)),
	)

	// 将首个预期的物理驱动器当作判断基准
	drive := drives[0]

	// 找出使用了当前物理驱动器的卷
	vs := scanVolumes(arrays, func(item *volume) bool {
		return collection.InSlice(drive, item.PhysicalDrives)
	})
	switch len(vs) {
	case 0: // 未找到使用了当前物理驱动器的卷
		item.Matched = hardware.MatchedNO
		return &item

	case 1: // 找到了使用当前物理驱动器的卷
		// 对比预期的物理驱动器数量与此卷所占用的物理驱动器的数量是否一致
		if len(vs[0].PhysicalDrives) != len(drives) {
			item.Matched = hardware.MatchedNO
			return &item
		}
		// 对比预期的物理驱动器列表与此卷所占用的物理驱动器是否完全一致
		for x := range drives {
			if !collection.InSlice(drives[x], vs[0].PhysicalDrives) {
				item.Matched = hardware.MatchedNO
				return &item
			}
		}
		item.Actual = fmt.Sprintf("%s@%s",
			strings.ToUpper(sett.Level),
			strings.Join(vs[0].PhysicalDrives, raid.Sep),
		)
		item.Matched = hardware.MatchedYES
		return &item
	default:
		// 理论上一块物理驱动器不可能存在于多个RAID阵列当中
		panic("unreachable")
	}
}
