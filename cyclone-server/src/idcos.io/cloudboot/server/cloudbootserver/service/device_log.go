package service

import (
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
)

// DeviceLogResp 安装日志返回值
type DeviceLogResp struct {
	// 设备序列号
	SN string `json:"sn"`
	// InstallType 系统参数ID
	DeviceSettingID uint `json:"device_setting_id"`
	// 操作系统安装模板
	Logs []DeviceLog `json:"logs"`
}

// DeviceLog 安装日志信息
type DeviceLog struct {
	// 日志记录ID
	ID uint `json:"id"`
	// 日志标题
	Title string `json:"title"`
	// 日志内容
	Content string `json:"content"`
	// 日志记录创建时间
	CreatedAt string `json:"create_at"`
	// 日志记录更新时间
	UpdatedAt string `json:"updated_at"`
}

// GetDeviceLogByDeviceSettingID 返回指定SN返回系统装机日志信息
func GetDeviceLogByDeviceSettingID(log logger.Logger, repo model.Repo, id uint) (deviceLogResp *DeviceLogResp, err error) {
	items, err := repo.GetDeviceLogsByDeviceSettingID(id)

	deviceLogResp = &DeviceLogResp{}
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return deviceLogResp, nil
	}

	deviceLogResp.SN = items[0].SN
	deviceLogResp.DeviceSettingID = items[0].DeviceSettingID

	var deviceLog []DeviceLog
	for _, item := range items {
		deviceLog = append(deviceLog, DeviceLog{
			ID:        item.ID,
			Title:     item.Title,
			Content:   item.Content,
			CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: item.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	deviceLogResp.Logs = deviceLog

	return deviceLogResp, nil
}
