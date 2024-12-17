package mysqlrepo

import (
	"github.com/jinzhu/gorm"

	"idcos.io/cloudboot/model"
)

// OverwriteHardwareSettings 覆写指定设备的硬件配置参数
func (repo *MySQLRepo) OverwriteHardwareSettings(sn string, items ...*model.DeviceHardwareSetting) (err error) {
	tx := repo.db.Begin()

	// 删除设备所有记录
	if err = tx.Where("sn = ?", sn).Delete(model.DeviceHardwareSetting{}).Error; err != nil {
		tx.Rollback()
		repo.log.Error(err)
		return err
	}

	for i := range items {
		if items[i] == nil {
			continue
		}
		// 新增设备记录
		if err = tx.Create(items[i]).Error; err != nil {
			tx.Rollback()
			repo.log.Error(err)
			return err
		}
	}
	return tx.Commit().Error
}

// GetHardwareSettingsBySN 返回指定设备的硬件配置配置参数
func (repo *MySQLRepo) GetHardwareSettingsBySN(sn string) (items []*model.DeviceHardwareSetting, err error) {
	if err := repo.db.Where("sn = ?", sn).Find(&items).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return items, nil
}

// UpdateHardwareSettingsApplied 更新满足过滤条件的硬件配置项实施状态
func (repo *MySQLRepo) UpdateHardwareSettingsApplied(cond *model.DeviceHardwareSetting, applied string) (affected int64, err error) {
	db := repo.db.Model(&model.DeviceHardwareSetting{})
	if cond != nil {
		if cond.SN != "" {
			db = db.Where("sn = ?", cond.SN)
		}
		if cond.Index >= 0 {
			db = db.Where("`index` = ?", cond.Index)
		}
		if cond.Category != "" {
			db = db.Where("category = ?", cond.Category)
		}
		if cond.Action != "" {
			db = db.Where("action = ?", cond.Action)
		}
	}

	db = db.Update("applied", applied)
	if err = db.Error; err != nil {
		repo.log.Error(err)
	}
	return db.RowsAffected, err
}

// RedoHardwareSettings 将指定设备的硬件配置项设置为'未实施'状态。
func (repo *MySQLRepo) RedoHardwareSettings(sn string) (affected int64, err error) {
	db := repo.db.Model(&model.DeviceHardwareSetting{}).Where("sn = ?", sn).Update("applied", model.NO)
	if err = db.Error; err != nil {
		repo.log.Error(err)
	}
	return db.RowsAffected, err
}
