package mysqlrepo

import (
	"idcos.io/cloudboot/model"
	"github.com/jinzhu/gorm"
)



// SaveDeviceLifecycle 保存设备生命周期记录
func (repo *MySQLRepo) SaveDeviceLifecycle(mod *model.DeviceLifecycle) (affected int64, err error) {
	if mod.ID != 0 {
		return repo.UpdateDeviceLifecycle(mod)
	} else {
		affected, _, err = repo.AddDeviceLifecycle(mod)
	}
	return affected, err
}

// AddDeviceLifecycle 增加设备生命周期记录
func (repo *MySQLRepo) AddDeviceLifecycle(mod *model.DeviceLifecycle) (affected int64, DeviceLifecycle *model.DeviceLifecycle, err error) {
	err = repo.db.Create(mod).Error
	if err != nil {
		return 0, nil, err
	}
	return 1, mod, nil
}

// UpdateDeviceLifecycle 修改设备生命周期记录
func (repo *MySQLRepo) UpdateDeviceLifecycle(mod *model.DeviceLifecycle) (affected int64, err error) {
	err = repo.db.Model(&model.DeviceLifecycle{}).Where("id = ?", mod.ID).Update(mod).Error
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// RemoveDeviceLifecycleByID 删除指定ID的设备生命周期记录
func (repo *MySQLRepo) RemoveDeviceLifecycleByID(id uint) (affected int64, err error) {
	mod := model.DeviceLifecycle{Model: gorm.Model{ID: id}}
	err = repo.db.Unscoped().Delete(mod).Error
	if err != nil {
		return 0, err
	}
	return 1, nil
}

//RemoveDeviceLifecycleBySN 删除指定SN的设备生命周期记录
func (repo *MySQLRepo) RemoveDeviceLifecycleBySN(sn string) (affected int64, err error) {
	mod := model.DeviceLifecycle{}
	db := repo.db.Unscoped().Where("sn = ?", sn).Delete(&mod)
	if err = db.Error; err != nil {
		repo.log.Errorf("RemoveDeviceLifecycleBySN failure:%s", err.Error())
	}
	return db.RowsAffected, nil
}

// GetDeviceLifecycleBySN 返回指定SN的设备生命周期
func (repo *MySQLRepo) GetDeviceLifecycleBySN(sn string) (*model.DeviceLifecycle, error) {
	var devLifecycle model.DeviceLifecycle
	if err := repo.db.Where("sn = ?", sn).Find(&devLifecycle).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &devLifecycle, nil
}


// UpdateDeviceLifecycleBySN 根据SN更新设备生命周期变更记录
func (repo *MySQLRepo) UpdateDeviceLifecycleBySN(deviceLifecycle *model.DeviceLifecycle) (err error) {
	db := repo.db.Model(&model.DeviceLifecycle{}).Where("sn = ?", deviceLifecycle.SN).Updates(deviceLifecycle)
	if db.Error != nil {
		repo.log.Error(err)
		return db.Error
	}
	return nil
}