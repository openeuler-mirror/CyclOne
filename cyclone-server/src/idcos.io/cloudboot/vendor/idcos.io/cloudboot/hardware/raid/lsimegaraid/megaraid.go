package lsimegaraid

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"idcos.io/cloudboot/hardware"
	"idcos.io/cloudboot/hardware/raid"
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

const (
	// disableJBOD 关闭JBOD模式
	disableJBOD = 0
	// enableJBOD 开启JBOD模式
	enableJBOD = 1
)

// switchJBODMode 打开或者关闭目标RAID卡下物理硬盘的JBOD模式。
// enable为true，开启JBOD模式。
// enable为false，关闭JBOD模式。
func (worker *worker) switchJBODMode(ctrlID string, enable bool) (err error) {
	v := enableJBOD
	if !enable {
		v = disableJBOD
	}
	// megacli -AdpSetProp -EnableJBOD -1 -a0
	_, err = worker.ExecByShell(tool, "-AdpSetProp", "-EnableJBOD", fmt.Sprintf("-%d", v), fmt.Sprintf("-a%s", ctrlID))
	return err
}

// controllerIDs 返回当前设备所有RAID卡的ID列表
func (worker *worker) controllerIDs() (items []string, err error) {
	// megacli64 -AdpAllInfo -aAll | grep "Adapter" | grep "#"
	output, err := worker.ExecByShell(tool, "-AdpAllInfo", "-aAll", "|", "grep", `"Adapter"`, "|", "grep", `"#"`)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "Adapter") {
			continue
		}

		fields := strings.SplitN(line, "#", 2)
		if len(fields) <= 1 {
			continue
		}
		items = append(items, fields[1])
	}
	return items, scanner.Err()
}

// findController 返回指定编号的Controller
func (worker *worker) findController(ctrlID string) (*raid.Controller, error) {
	// megacli64 -AdpAllInfo -a0
	output, err := worker.ExecByShell(tool, "-AdpAllInfo", fmt.Sprintf("-a%s", ctrlID))
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
		if strings.HasPrefix(line, "Product Name") && strings.Contains(line, colonSeparator) {
			ctl.ModelName = worker.extractValue(line, colonSeparator)
		}
		if strings.HasPrefix(line, "FW Version") {
			ctl.FirmwareVersion = worker.extractValue(line, colonSeparator)
			break
		}
	}
	return &ctl, scanner.Err()
}

// checkRAID 检查RAID参数
func (worker *worker) checkRAID(level raid.Level, drives []string) (err error) {
	if level != raid.RAID0 && level != raid.RAID1 && level != raid.RAID10 &&
		level != raid.RAID5 && level != raid.RAID50 && level != raid.RAID6 && level != raid.RAID60 {
		return raid.ErrUnsupportedLevel
	}
	return nil
}

type physicalDrive struct {
	Enclosure int
	Slot      int
	RawSize   string // 包含单位
	PDType    string
}

func (pd physicalDrive) identity() string {
	return fmt.Sprintf("%d:%d", pd.Enclosure, pd.Slot)
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

type pds []physicalDrive

//以下三个函数为了实现磁盘容量的排序规则
func (items pds) Len() int {
	return len(items)
}

func (items pds) Less(i, j int) bool {
	if items[i].Enclosure < items[j].Enclosure {
		return true
	} else if items[i].Slot < items[j].Slot {
		return true
	}
	return false
}

func (items pds) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}

func (worker *worker) physicalDrives(ctrlID string) (pds, error) {
	output, err := worker.ExecByShell(tool, "-PDlist", fmt.Sprintf("-a%s", ctrlID), "-NoLog")
	if err != nil {
		return nil, err
	}
	return parsePDs(output)
}

// parsePDs 解析内容中包含的物理驱动器列表
func parsePDs(output []byte) (list pds, err error) {
	arr := strings.Split(string(output), "\n\n\n")
	for i := 0; i < len(arr)-1; i++ {
		var hd physicalDrive
		lines := strings.Split(arr[i], "\n")
		for _, l := range lines {
			if strings.Contains(l, "Enclosure Device") {
				items := strings.Fields(l)
				if hd.Enclosure, err = strconv.Atoi(strings.TrimSpace(items[len(items)-1])); err != nil {
					hd.Enclosure = -1
					//return
				}
			} else if strings.Contains(l, "Slot Number:") {
				items := strings.Fields(l)
				if hd.Slot, err = strconv.Atoi(strings.TrimSpace(items[len(items)-1])); err != nil {
					hd.Slot = -1
					//return
				}
			} else if strings.Contains(l, "Raw Size:") {
				items := strings.Fields(l)
				hd.RawSize = items[2] + " " + items[3]
			} else if strings.Contains(l, "PD Type:") {
				items := strings.Fields(l)
				hd.PDType = strings.TrimSpace(items[len(items)-1])
			}
		}
		list = append(list, hd)
	}
	return
}

// checkArrays 检查实际的RAID阵列是否与预期配置一致
func (worker *worker) checkArrays(sett *raid.Setting) (items []hardware.CheckingItem) {
	log := worker.GetLog()
	if sett == nil || len(sett.Controllers) == 0 || len(sett.Controllers[0].Arrays) == 0 {
		return nil
	}
	realRAID, _ := worker.getRealRAID()
	if realRAID == nil {
		items = append(items, hardware.CheckingItem{
			Title:   "Check RAID",
			Matched: hardware.MatchedUnknown,
			Error:   "Collect actual RAID info failed",
		})
		return
	}
	if log != nil {
		log.Debugf("actual RAID:%v", *realRAID)
	}
	for i := range sett.Controllers[0].Arrays {
		if item := worker.checkArray(&sett.Controllers[0].Arrays[i], realRAID); item != nil {
			items = append(items, *item)
		}
	}
	// hotspare
	if sett.Controllers[0].Hotspares != "" {
		if item := worker.checkGlobalHotSpares(sett, realRAID); item != nil {
			items = append(items, *item)
		}
	}
	return items
}

// checkArray 检查指定的单个RAID阵列配置与实际是否一致
func (worker *worker) checkArray(sett *raid.ArraySetting, realRAID *raid.Setting) *hardware.CheckingItem {
	if sett == nil {
		return nil
	}
	//TODO 需要增加NonRAID的校验
	item := hardware.CheckingItem{
		Title:    fmt.Sprintf("创建%s", strings.ToUpper(sett.Level)),
		Expected: fmt.Sprintf("%s@%s", strings.ToUpper(sett.Level), sett.Drives),
		Matched:  hardware.MatchedUnknown,
	}

	for _, i := range realRAID.Controllers[0].Arrays {
		if strings.ToUpper(sett.Level) == i.Level {
			item.Actual = fmt.Sprintf("%s@%s", strings.ToUpper(i.Level), i.Drives)
			if sett.Drives == i.Drives {
				item.Matched = hardware.MatchedYES
				return &item
			} else {
				//可能存在多组相同级别的RAID，这里只是置可能值，还要继续寻找
				item.Matched = hardware.MatchedNO
			}
		}
	}

	return &item
}

// getRealRAID 通过megaclisas-status.py查询实际的阵列
func (worker *worker) getRealRAID() (setting *raid.Setting, err error) {
	out, err := worker.Exec("python", "/usr/local/bin/megaraid-status.py") //这个路径需要约定好
	if err != nil {
		worker.GetLog().Error(err)
		return nil, err
	}
	return parseMegacliStatus(out)
}

func parseMegacliStatus(out []byte) (setting *raid.Setting, err error) {
	setting = new(raid.Setting)
	setting.Controllers = make([]raid.ControllerSetting, 1, 5) //TODO 暂时只支持单raid卡场景检查
	sections := strings.Split(string(out), "\n\n")
	m := make(map[string]*raid.ArraySetting, 0) //临时检索容器
	if len(sections) >= 3 {
		//第二自然段，第三行开始读取各RAID
		lines := strings.Split(sections[1], "\n")
		if len(lines) >= 2 {
			for i := 2; i < len(lines); i++ {
				fileds := strings.Split(lines[i], "|")
				id := strings.TrimSpace(fileds[0])
				typ := strings.TrimSpace(fileds[1])
				raidItem := &raid.ArraySetting{
					Level: strings.Replace(typ, "-", "", -1), //RAID-1搞成RAID1
				}
				m[id] = raidItem
			}
		}
		//第三自然段开始解析具体的磁盘IDs
		lines = strings.Split(sections[2], "\n")
		if len(lines) >= 2 {
			drives := make([]string, 0)
			idPre := ""
			for i := 2; i < len(lines); i++ {
				fileds := strings.Split(lines[i], "|")
				id := strings.TrimSpace(fileds[0])[0:4] //只取前4位，和第二段中解析出来的id值就一致了
				slotID := strings.TrimSpace(fileds[7])

				if i != 2 && id != idPre { //这里就是想分组归类磁盘id
					(m[idPre]).Drives = strings.Join(drives, ",")
					drives = drives[:0]
				}
				drives = append(drives, slotID[1:len(slotID)-1]) //[32:1] => 32:1
				if i == len(lines)-1 {
					(m[idPre]).Drives = strings.Join(drives, ",")
				}
				idPre = id
			}
		}

		//热备盘解析
		if len(sections) >= 4 {
			lines = strings.Split(sections[3], "\n")
			if len(lines) >= 2 {
				drives := make([]string, 0)
				for i := 2; i < len(lines); i++ {
					fileds := strings.Split(lines[i], "|")
					if len(fileds) > 8 {
						status := strings.TrimSpace(fileds[4]) //Hotspare, Spun Up|Unconfigured(good), Spun down
						if strings.Contains(status, "Hotspare") {
							slotID := strings.TrimSpace(fileds[7])
							drives = append(drives, slotID[1:len(slotID)-1]) //[32:1] => 32:1
						}
					}
				}
				setting.Controllers[0].Hotspares = strings.Join(drives, ",")
			}
		}
	}
	for _, v := range m {
		setting.Controllers[0].Arrays = append(setting.Controllers[0].Arrays, *v)
	}
	return setting, nil
}

// checkGlobalHotSpares 检查实际的全局热备盘设置是否与预期配置一致
func (worker *worker) checkGlobalHotSpares(sett *raid.Setting, realRAID *raid.Setting) *hardware.CheckingItem {
	if sett == nil || len(sett.Controllers) == 0 || sett.Controllers[0].Hotspares == "" {
		return nil
	}
	item := hardware.CheckingItem{
		Title:    "创建全局热备",
		Expected: sett.Controllers[0].Hotspares,
		Actual:   realRAID.Controllers[0].Hotspares,
		Matched:  hardware.MatchedUnknown,
	}
	if realRAID.Controllers[0].Hotspares != strings.TrimSpace(sett.Controllers[0].Hotspares) {
		item.Matched = hardware.MatchedNO
	} else {
		item.Matched = hardware.MatchedYES
	}
	return &item
}

// 返回 worker 所使用的cmdline
func (worker *worker) GetCMDLine() string {
	return tool
}