package avagomegaraid

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
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

// genPDIdentity 返回形如'/cx/ex/sx'或者'/cx/sx'的硬盘标识
func (worker *worker) genPDIdentity(ctrlID, esid string) string {
	if strings.Contains(esid, colonSeparator) {
		fields := strings.SplitN(esid, colonSeparator, 2)
		return fmt.Sprintf("/c%s/e%s/s%s", ctrlID, fields[0], fields[1])
	}
	return fmt.Sprintf("/c%s/s%s", ctrlID, esid)
}

const (
	// disableJBOD 关闭JBOD模式
	disableJBOD = "off"
	// enableJBOD 开启JBOD模式
	enableJBOD = "on"
)

// switchJBODMode 打开或者关闭目标RAID卡下物理硬盘的JBOD模式。
// enable为true，开启JBOD模式。
// enable为false，关闭JBOD模式。
func (worker *worker) switchJBODMode(ctrlID string, enable bool) (err error) {
	v := enableJBOD
	if !enable {
		v = disableJBOD
	}
	_, err = worker.ExecByShell(tool, fmt.Sprintf("/c%s", ctrlID), "set", fmt.Sprintf("jbod=%s", v), "force")

	return err
}

// 校验是否支持 JBOD模式
func (worker *worker) checkSupportJBOD(ctrlID string) (isSupport bool, err error) {
	output, err := worker.ExecByShell(tool, fmt.Sprintf("/c%s", ctrlID), "show", "all", "|", "grep", "'Support JBOD ='")
	if err != nil {
		return false, err
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.Contains(line, "=Yes") {
			isSupport = true
		} else {
			isSupport = false
		}
	}
	return isSupport, scanner.Err()
}


// validateLevel 校验RAID级别
func (worker *worker) validateLevel(level raid.Level) error {
	if level != raid.RAID0 && level != raid.RAID1 && level != raid.RAID10 &&
		level != raid.RAID5 && level != raid.RAID50 && level != raid.RAID6 && level != raid.RAID60 {
		return raid.ErrUnsupportedLevel
	}
	return nil
}

var numReg = regexp.MustCompile(`\d+`)

// controllerIDs 返回当前设备所有RAID卡的ID列表
func (worker *worker) controllerIDs() (items []string, err error) {
	// storcli64 show
	output, err := worker.ExecByShell(tool, "show")
	if err != nil {
		return nil, err
	}

	var started bool
	scanner := bufio.NewScanner(bytes.NewBuffer(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.Contains(line, "Ctl") && strings.Contains(line, "Model") && strings.Contains(line, "Ports") {
			started = true
			continue
		}
		if !started || line == "" {
			continue
		}

		fields := strings.Fields(line)
		if !numReg.MatchString(fields[0]) {
			continue
		}
		items = append(items, fields[0])
	}
	return items, scanner.Err()
}

// findController 返回指定编号的Controller
func (worker *worker) findController(ctrlID string) (*raid.Controller, error) {
	// storcli64 /c0 show
	output, err := worker.Base.ExecByShell(tool, fmt.Sprintf("/c%s", ctrlID), "show")
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
		if strings.HasPrefix(line, "FW Version") && strings.Contains(line, eqSeparator) {
			ctl.FirmwareVersion = worker.extractValue(line, eqSeparator)
			break
		}
	}
	ctl.ModelName, _ = worker.findModelName(ctrlID)
	return &ctl, scanner.Err()
}

// findModelName 查询指定RAID卡控制器的型号
func (worker *worker) findModelName(ctrlID string) (string, error) {
	// storcli64 show
	output, err := worker.ExecByShell(tool, "show")
	if err != nil {
		return "", err
	}

	var started bool
	scanner := bufio.NewScanner(bytes.NewBuffer(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.Contains(line, "Ctl") && strings.Contains(line, "Model") && strings.Contains(line, "Ports") {
			started = true
			continue
		}
		if !started || line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 || fields[0] != ctrlID {
			continue
		}
		return fields[1], nil
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
	ESID  string `json:"EID:Slt"`
	State string `json:"State"`
	Size  string `json:"Size"`
	Med   string `json:"Med"`
}

// findPhysicalDriveIDByIndex 通过切片索引查找物理驱动器ID
func findPhysicalDriveIDByIndex(items []physicalDrive, index int) (id string, err error) {
	if index < 0 || index > len(items)-1 {
		return "", raid.ErrInvalidDiskIdentity
	}
	return items[index].ESID, nil
}

// findPhysicalDriveIDsByIndexes 通过切片索引列表查找物理驱动器ID
func findPhysicalDriveIDsByIndexes(items []physicalDrive, indexes ...int) (ids []string, err error) {
	if len(indexes) <= 0 {
		for i := range items {
			ids = append(ids, items[i].ESID)
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
		ids = append(ids, items[i].ESID)
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
	return pd != nil && strings.Contains(pd.State, "GHS")
}

// physicalDrives 返回指定controller下的物理驱动器列表
func (worker *worker) physicalDrives(ctrlID string) (items []physicalDrive, err error) {
	//storcli /c0/eall/sall show J
	output, err := worker.ExecByShell(tool, fmt.Sprintf("/c%s/eall/sall", ctrlID), "show", "J") // TODO 不再JSON方式，因为会在EFI下失效。
	if err != nil {
		return nil, err
	}
	var ctrls struct {
		Ctrls []struct {
			RespData struct {
				Drives []physicalDrive `json:"Drive Information"`
			} `json:"Response Data"`
		} `json:"Controllers"`
	}
	if err = json.Unmarshal(output, &ctrls); err != nil {
		if log := worker.GetLog(); log != nil {
			log.Error(err)
		}
		return nil, err
	}
	if len(ctrls.Ctrls) <= 0 {
		return nil, nil
	}
	return ctrls.Ctrls[0].RespData.Drives, nil
}

// logicalDrive 逻辑驱动器(阵列)
type logicalDrive struct {
	ID             string // 形如'/c0/v0'
	Name           string
	Type           string // 形如'RAID1'
	State          string
	Size           string
	PhysicalDrives []physicalDrive // 逻辑驱动器底层所依赖的物理驱动器
}

type ldScanFunc func(item *logicalDrive) bool

// scanPDs 逐个扫描逻辑驱动器列表并返回满足所有过滤条件的逻辑驱动器
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

// physicalDrives 返回指定controller下的物理驱动器列表
func (worker *worker) logicalDrives(ctrlID string) (items []logicalDrive, err error) {
	// storcli64 /c0/vall show all J
	output, err := worker.ExecByShell(tool, fmt.Sprintf("/c%s/vall", ctrlID), "show", "all", "J") // TODO 不再JSON方式，因为会在EFI下失效。
	if err != nil {
		return nil, err
	}

	// gjson语法: https://github.com/tidwall/gjson#path-syntax
	// 示例: ./testdata/storcli64_c0_vall_show_all_j.json
	gjson.GetBytes(output, `Controllers.0.Response Data`).ForEach(func(key, value gjson.Result) bool {
		k := key.String()
		prefix := fmt.Sprintf("/c%s/v", ctrlID)
		vid := strings.TrimPrefix(k, prefix)

		if !strings.HasPrefix(k, prefix) {
			return true
		}

		if !value.IsArray() {
			return true
		}

		vPDs := gjson.GetBytes(output, fmt.Sprintf(`Controllers.0.Response Data.PDs for VD %s`, vid)).Array()
		pds := make([]physicalDrive, 0, len(vPDs))
		for i := range vPDs {
			pds = append(pds, physicalDrive{
				ESID:  vPDs[i].Get("EID:Slt").String(),
				State: vPDs[i].Get("State").String(),
				Size:  vPDs[i].Get("Size").String(),
				Med:   vPDs[i].Get("Med").String(),
			})
		}

		vVD := value.Array()[0]
		items = append(items, logicalDrive{
			ID:             k,
			Name:           vVD.Get("Name").String(),
			Type:           vVD.Get("TYPE").String(),
			State:          vVD.Get("State").String(),
			Size:           vVD.Get("Size").String(),
			PhysicalDrives: pds,
		})
		return true
	})

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
func (worker *worker) checkGlobalHotSpares(ctrlID string, sett *raid.ControllerSetting) *hardware.CheckingItem {
	if sett == nil || sett.Hotspares == "" {
		return nil
	}
	item := hardware.CheckingItem{
		Title:    fmt.Sprintf("[Controller%s] 设置全局热备盘", ctrlID),
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

	// 预期的用于备份的物理驱动器列表
	spares := strings.Split(sett.Hotspares, raid.Sep)

	// 对比实际与预期的热备盘列表是否一致
	if worker.drivesEq(physicalDriveIDs(sparePDs), spares) {
		item.Matched = hardware.MatchedYES
	} else {
		item.Matched = hardware.MatchedNO
	}
	return &item
}
