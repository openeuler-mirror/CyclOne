package mysqlrepo

import "idcos.io/cloudboot/model"

// GetDeviceLogsByDeviceSettingID 根据装机参数id获取装机日志信息
func (repo *MySQLRepo) GetDeviceLogsByDeviceSettingID(id uint) (deviceLogs []model.DeviceLog, err error) {
	deviceLogs = []model.DeviceLog{}
	db := repo.db.Model(&model.DeviceLog{}).Where("device_setting_id = ? and type = ?", id, "install").Find(&deviceLogs)
	if db.Error != nil {
		repo.log.Error(err)
		return nil, db.Error
	}
	return deviceLogs, nil
}

// UpdateDeviceLogType 修改操作系统安装进度记录
func (repo *MySQLRepo) UpdateDeviceLogType(settingID uint, fromLogType, toLogType string) (affected int64, err error) {
	db := repo.db.Model(&model.DeviceLog{}).Where("device_setting_id = ? and type = ?", settingID, fromLogType).Update("type", toLogType)
	if db.Error != nil {
		repo.log.Error(err)
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// SaveDeviceLog 新增操作系统安装进度记录
func (repo *MySQLRepo) SaveDeviceLog(log *model.DeviceLog) (affected int64, err error) {
	if log.ID > 0 {
		// 更新
		db := repo.db.Model(&model.DeviceLog{}).Where("id = ?", log.ID).Updates(log)
		if err = db.Error; err != nil {
			repo.log.Error(err)
			return db.RowsAffected, err
		}
		return db.RowsAffected, nil
	}
	// 新增
	db := repo.db.Create(log)
	if err = db.Error; err != nil {
		repo.log.Error(err)
		return db.RowsAffected, err
	}
	return db.RowsAffected, nil
}
