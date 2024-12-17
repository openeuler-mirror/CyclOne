package setting

import (
	"strconv"

	"idcos.io/cloudboot/hardware/raid"
	"idcos.io/cloudboot/model"
)

const (
	// ON 打开
	ON = "ON"
	// OFF 关闭
	OFF = "OFF"
)

// HardwareSetting 硬件配置项。
// 详见http://gitlab.idcos.io/cloudboot/cloudboot/issues/344
type HardwareSetting struct {
	Index    uint              `json:"index"`    // 当前配置项在配置项序列中的索引号（自0始）
	Category string            `json:"category"` // 硬件配置动作类别，可选值: raid|oob|bios|firmware|reboot
	Action   string            `json:"action"`   // 配置动作，如add_user、reset_bmc、set_ip、set_bios等
	Metadata map[string]string `json:"metadata"` // 配置项元信息
	Applied  string            `json:"applied"`  // 是否已实施，可选值: Yes|No
	// Metadata HardwareSettingMetadata `json:"metadata"`
}

// HardwareSettings 硬件配置参数集合
type HardwareSettings []*HardwareSetting

// ExtractRAID 提取其中的RAID配置
func (setts HardwareSettings) ExtractRAID() *raid.Setting {
	items := setts.FindByCategory(model.CategoryRAID)
	if len(items) <= 0 {
		return nil
	}
	ctrlMap := make(map[string]*raid.ControllerSetting)
	for i := range items {
		if items[i] == nil {
			continue
		}
		ctrlIdx, ok := items[i].Metadata["controller_index"]
		if !ok || ctrlIdx == "" {
			continue
		}

		ctrl, ok := ctrlMap[ctrlIdx]
		if !ok || ctrl == nil {
			ctrl = new(raid.ControllerSetting)
			ctrlMap[ctrlIdx] = ctrl
		}
		idx, _ := strconv.Atoi(ctrlIdx)
		ctrl.Index = uint(idx)

		switch items[i].Action {
		case model.ActionRAIDCreate:
			ctrl.Arrays = append(ctrl.Arrays, raid.ArraySetting{
				Level:  items[i].Metadata["level"],
				Drives: items[i].Metadata["drives"],
			})
		case model.ActionRAIDClear:
			ctrl.Clear = items[i].Metadata["clear"]
		case model.ActionRAIDSetJBOD:
			ctrl.JBODs = items[i].Metadata["drives"]
		case model.ActionRAIDSetGlobalHotspare:
			ctrl.Hotspares = items[i].Metadata["drives"]
		case model.ActionRAIDInitDisk:
			ctrl.DiskInit = items[i].Metadata["init"]
		}
	}

	ctrls := make([]raid.ControllerSetting, 0, len(ctrlMap))
	for _, ctrl := range ctrlMap {
		ctrls = append(ctrls, *ctrl)
	}

	return &raid.Setting{
		Controllers: ctrls,
	}
}

// Find 返回满足条件的所有元素集合
func (setts HardwareSettings) Find(category, action string) (items []*HardwareSetting) {
	for i := range setts {
		if setts[i] == nil {
			continue
		}
		if setts[i].Category == category && setts[i].Action == action {
			items = append(items, setts[i])
		}
	}
	return items
}

// FindByCategory 返回满足条件的所有元素集合
func (setts HardwareSettings) FindByCategory(category string) (items []*HardwareSetting) {
	for i := range setts {
		if setts[i] == nil {
			continue
		}
		if setts[i].Category == category {
			items = append(items, setts[i])
		}
	}
	return items
}

// HardwareTemplateSetting 硬件模板配置参数
type HardwareTemplateSetting struct {
	ID   int    `json:"id"`   // 硬件模板ID
	Name string `json:"name"` // 硬件模板名称
}
