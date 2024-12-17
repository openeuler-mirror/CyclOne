package mysqlrepo

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/voidint/page"

	"idcos.io/cloudboot/model"
)

// GetDeviceSettingsByInstallationTimeout 查询安装超时的设备装机记录列表。timeout-超时时间，单位秒。
func (repo *MySQLRepo) GetDeviceSettingsByInstallationTimeout(timeout int64) (items []*model.DeviceSetting, err error) {
	if timeout <= 0 {
		return make([]*model.DeviceSetting, 0, 0), nil
	}
	if err = repo.db.Model(&model.DeviceSetting{}).Where("status = ? AND UNIX_TIMESTAMP(NOW())-UNIX_TIMESTAMP(installation_start_time) >= ?", model.InstallStatusIng, timeout).Find(&items).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return items, nil
}

// SetInstallationTimeout 为指定设备序列号的装机参数进行'安装超时'处理。
func (repo *MySQLRepo) SetInstallationTimeout(sns ...string) (affected int64, err error) {
	db := repo.db.Model(&model.DeviceSetting{}).Where("sn IN (?) AND status = ?", sns, model.InstallStatusIng).Updates(map[string]interface{}{
		"install_progress":      0,
		"status":                model.InstallStatusFail,
		"installation_end_time": time.Now(),
	})
	if err = db.Error; err != nil {
		repo.log.Error(err)
		return 0, err
	}
	return db.RowsAffected, nil
}

//UpdateInstallStatusAndProgressByID 更新安装状态和进度
func (repo *MySQLRepo) UpdateInstallStatusAndProgressByID(id uint, status string, progress float64) (affected int64, err error) {
	kv := make(map[string]interface{}, 3)
	kv["status"] = status
	kv["install_progress"] = progress

	if status == model.InstallStatusIng && progress == 0 {
		now := time.Now()
		kv["installation_start_time"] = &now
		kv["installation_end_time"] = nil
	}
	if status == model.InstallStatusFail || status == model.InstallStatusSucc {
		now := time.Now()
		kv["installation_end_time"] = &now
	}

	db := repo.db.Model(&model.DeviceSetting{}).Where("id = ?", id).Updates(kv)
	if err = db.Error; err != nil {
		repo.log.Error(err)
		return 0, err
	}
	return db.RowsAffected, nil
}

//DeleteDeviceSettingByID 删除指定ID的装机参数
func (repo *MySQLRepo) DeleteDeviceSettingByID(id uint) (*model.DeviceSetting, error) {
	mod := model.DeviceSetting{}
	err := repo.db.Unscoped().Where("id = ?", id).Delete(&mod).Error
	if err != nil {
		repo.log.Errorf("DeleteDeviceSettingByID failure,%s", err.Error())
	}
	return &mod, err
}

//DeleteDeviceSettingBySN 删除指定SN的装机参数
func (repo *MySQLRepo) DeleteDeviceSettingBySN(sn string) (*model.DeviceSetting, error) {
	mod := model.DeviceSetting{}
	err := repo.db.Unscoped().Where("sn = ?", sn).Delete(&mod).Error
	if err != nil {
		repo.log.Errorf("DeleteDeviceSettingBySN failure,%s", err.Error())
	}
	return &mod, err
}

// AddDeviceSettings 批量添加设备装机参数
func (repo *MySQLRepo) AddDeviceSettings(items ...*model.DeviceSetting) (err error) {
	tx := repo.db.Begin()
	for i := range items {
		if err = tx.Create(items[i]).Error; err != nil {
			tx.Rollback()
			repo.log.Error(err)
			return err
		}
	}
	return tx.Commit().Error
}

// UpdateDeviceSettingBySN 根据SN更新设备装机参数
func (repo *MySQLRepo) UpdateDeviceSettingBySN(sett *model.DeviceSetting) (affected int64, err error) {
	db := repo.db.Model(&model.DeviceSetting{}).Where("sn = ?", sett.SN).Updates(sett) // Updates仅更新非0值
	if db.Error != nil {
		repo.log.Error(err)
		return db.RowsAffected, db.Error
	}
	return db.RowsAffected, nil
}

// UpdateDeviceSettingIPConfigBySN 根据SN更新设备装机参数
//func (repo *MySQLRepo) UpdateDeviceSettingIPConfigBySN(sn string, intranet bool) (affected int64, err error) {
//	db := repo.db.Model(&model.DeviceSetting{}).Where("sn = ?", sn)
//	if intranet {
//		db.Updates(map[string]interface{}{
//			"intranet_ip_network_id": 0,
//			"intranet_ip":            "",
//		})
//	} else {
//		db.Updates(map[string]interface{}{
//			"extranet_ip_network_id": 0,
//			"extranet_ip":            "",
//		})
//	}
//	if db.Error != nil {
//		repo.log.Error(err)
//		return db.RowsAffected, db.Error
//	}
//	return db.RowsAffected, nil
//}

// SaveDeviceSetting 保存设备装机参数。若入参包含主键ID，则进行更新操作，否则进行新增操作。
func (repo *MySQLRepo) SaveDeviceSetting(sett *model.DeviceSetting) (err error) {
	db := repo.db.Save(sett)
	if err = db.Error; err != nil {
		repo.log.Error(err)
		return err
	}
	return nil
}

// GetDeviceSettingBySN 根据sn查询设备装机参数
func (repo *MySQLRepo) GetDeviceSettingBySN(sn string) (devSetting *model.DeviceSetting, err error) {
	var row model.DeviceSetting
	if err = repo.db.Where("sn = ?", sn).Find(&row).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &row, nil
}

// GetDeviceSettingByID 根据id查询设备装机参数
func (repo *MySQLRepo) GetDeviceSettingByID(id uint) (devSetting *model.DeviceSetting, err error) {
	var row model.DeviceSetting
	if err = repo.db.Where("id = ?", id).Find(&row).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &row, nil
}

// CountDeviceSettingCombines 统计满足过滤条件的装机参数数量
func (repo *MySQLRepo) CountDeviceSettingCombines(cond *model.CombineDeviceSetting) (count int64, err error) {
	db := repo.db.Table("device_setting ").
		Joins("inner join device ON  device_setting.sn = device.sn")

	db = addDeviceSettingCond(db, cond)

	if err = db.Select("device_setting.id").Count(&count).Error; err != nil {
		repo.log.Error(err)
		return 0, err
	}

	return count, nil

}

// CountDeviceSettingByStatus 统计对应设备安装状态的数量
func (repo *MySQLRepo) CountDeviceSettingByStatus(status string) (count int64, err error) {
	db := repo.db.Table("device_setting  as setting").
		Joins("inner join device on device.sn = setting.sn").Where("setting.status = ?", status).Count(&count)

	if db.Error != nil {
		return 0, db.Error
	}

	return count, nil
}

// GetDeviceSettingCombinesByCond 返回满足过滤条件的装机参数
func (repo *MySQLRepo) GetDeviceSettingCombinesByCond(cond *model.CombineDeviceSetting, orderby model.OrderBy, limiter *page.Limiter) (item []*model.DeviceSetting, err error) {
	db := repo.db.Table("device_setting ").
		Joins("inner join device ON  device_setting.sn = device.sn")

	db = addDeviceSettingCond(db, cond)

	for i := range orderby {
		db = db.Order(orderby[i].String())
	}

	if limiter != nil {
		db = db.Limit(limiter.Limit).Offset(limiter.Offset)
	}

	if err = db.Select("device_setting.*").Find(&item).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return item, nil
}

// addDeviceSettingCond 添加物理设备参数信息
func addDeviceSettingCond(db *gorm.DB, cond *model.CombineDeviceSetting) *gorm.DB {
	if cond.ServerRoomName != "" {
		db = db.Joins("left join server_room ON  device.server_room_id = server_room.id")
	}
	if cond.ServerCabinetNumber != "" {
		db = db.Joins("left join server_cabinet ON  device.server_cabinet_id = server_cabinet.id")
	}
	if cond != nil {
		db = MultiNumQuery(db, "device.idc_id", cond.IDCID)
		db = MultiNumQuery(db, "device.server_room_id", cond.ServerRoomID)
		db = MultiNumQuery(db, "device.server_cabinet_id", cond.ServerCabinetID)
		db = MultiNumQuery(db, "device.server_usite_id", cond.ServerUsiteID)
		db = MultiQuery(db, "server_room.name", cond.ServerRoomName)
		db = MultiQuery(db, "server_cabinet.number", cond.ServerCabinetNumber)
		db = MultiQuery(db, "device.sn", cond.Sn)
		db = MultiQuery(db, "device.fixed_asset_number", cond.FN)
		db = MultiQuery(db, "device.category", cond.Category)
		db = MultiQuery(db, "device_setting.extranet_ip", cond.ExtranetIP)
		db = MultiQuery(db, "device_setting.intranet_ip", cond.IntranetIP)
		if cond.HardwareTemplateID > 0 {
			db = db.Where("device_setting.hardware_template_id = ?", cond.HardwareTemplateID)
		}
		if cond.ImageTemplateID > 0 {
			db = db.Where("device_setting.image_template_id = ?", cond.ImageTemplateID)
		}
		if cond.Status != "" {
			db = db.Where("device_setting.Status = ?", cond.Status)
		}

	}
	return db
}
