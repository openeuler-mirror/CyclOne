package service

import (
	"encoding/json"
	"fmt"

	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
)

// HardwareSettingResp 硬件配置模板返回结构体
type HardwareSettingResp struct {
	// 配置信息
	Settings []HardwareSetting `json:"settings"`
}

// HardwareSetting 硬件配置参数
type HardwareSetting struct {
	// 硬件配置动作类别 raid|oob|bios|firmware|reboot
	Category string `json:"category"`
	// 硬件配置动作 reboot,add_user,reset_bmc,set_ip,set_bios,clear_settings,create_array,set_global_hotspare,init_disk,set_jbod,update_package
	Action string `json:"action"`
	// 配置元数据
	Metadata map[string]interface{} `json:"metadata"`
}

// HardwareTemplateResp 硬件模板返回体
type HardwareTemplateResp struct {
	// 是否内置
	Builtin string `json:"builtin"`
	// 配置数据
	Data Data `json:"data"`
	// 型号
	ModelName string `json:"model_name"`
	// 名称
	Name string `json:"name"`
	// 厂商
	Vendor string `json:"vendor"`
}

// HardwareInfo 硬件相关信息
type HardwareInfo struct {
	//序列号
	SN               string    `json:"sn"`
	//厂商
	Vendor           string    `json:"vendor"`
	//厂商型号
	DevModel         string    `json:"model"`
	//设备类型（自定义型号，如M10）
	Category         string    `json:"category"`
	//硬件备注（CPU\MEM\NET\DISK）
	HardwareRemark   string    `json:"hardware_remark"`
}

// GetHardwareSettingsBySN 查询设备硬件配置参数信息
func GetHardwareSettingsBySN(log logger.Logger, repo model.Repo, sn string) (item *HardwareSettingResp, err error) {
	item = &HardwareSettingResp{
		Settings: []HardwareSetting{},
	}

	devSetting, err := repo.GetDeviceSettingBySN(sn)
	if err != nil {
		return item, err
	}

	if devSetting.ID <= 0 {
		return item, nil
	}

	hardwareTemp, err := repo.GetHardwareTemplateByID(devSetting.HardwareTemplateID)
	if err != nil {
		return item, err
	}

	if hardwareTemp.ID <= 0 {
		return item, nil
	}

	settings := item.Settings

	err = json.Unmarshal([]byte(hardwareTemp.Data), &settings)
	if err != nil {
		return item, err
	}

	item.Settings = settings

	return item, nil
}

// GetHardwareSettingsByID 查询设备硬件配置参数信息
func GetHardwareSettingsByID(repo model.Repo, sn uint) (resp *HardwareTemplateResp, err error) {
	item, err := repo.GetHardwareTemplateByID(sn)
	if err != nil {
		return nil, err
	}

	fmt.Println(item.Data)

	var datas Data

	err = json.Unmarshal([]byte(item.Data), &datas)
	if err != nil {
		return nil, err
	}

	resp = &HardwareTemplateResp{
		Name:      item.Name,
		Vendor:    item.Vendor,
		ModelName: item.ModelName,
		Builtin:   item.Builtin,
		Data:      datas,
	}

	return
}

// GetHardwareInfoBySN 根据SN获取设备硬件相关信息
func GetHardwareInfoBySN(log logger.Logger, repo model.Repo, sn string) (resp *HardwareInfo, err error) {
	devInfo, err := repo.GetDeviceBySN(sn)
	if err != nil {
		return nil, err
	}

	resp = &HardwareInfo{
		SN: devInfo.SN,
		Vendor: devInfo.Vendor,
		DevModel: devInfo.DevModel,
		Category: devInfo.Category,
		HardwareRemark: devInfo.HardwareRemark,
	}
	return resp,nil
}