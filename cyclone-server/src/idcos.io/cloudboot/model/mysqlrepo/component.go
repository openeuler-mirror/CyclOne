package mysqlrepo

import (
	"github.com/jinzhu/gorm"

	"idcos.io/cloudboot/model"
)

// SaveComponentLogBySN 保存组件日志
func (repo *MySQLRepo) SaveComponentLogBySN(cl *model.ComponentLog) (err error) {
	var count int64
	if err = repo.db.Model(&model.ComponentLog{}).Where("sn = ?", cl.SN).Where("component = ?", cl.Component).Count(&count).Error; err != nil {
		repo.log.Error(err)
		return err
	}

	if count > 0 {
		if err = repo.db.Model(&model.ComponentLog{}).Where("sn = ?", cl.SN).Where("component = ?", cl.Component).Update("log", cl.LogData).Error; err != nil {
			repo.log.Error(err)
			return err
		}
		return nil
	}

	if err = repo.db.Create(cl).Error; err != nil {
		repo.log.Error(err)
		return err
	}
	return nil
}

// GetComponentLog 查询指定设备的指定组件日志
func (repo *MySQLRepo) GetComponentLog(sn, component string) (*model.ComponentLog, error) {
	var cl model.ComponentLog
	if err := repo.db.Model(&model.ComponentLog{}).Where("sn = ?", sn).Where("component = ?", component).Find(&cl).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &cl, nil
}
