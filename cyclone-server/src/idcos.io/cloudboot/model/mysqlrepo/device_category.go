package mysqlrepo

import (
	"errors"

	"github.com/voidint/page"

	"github.com/jinzhu/gorm"
	"idcos.io/cloudboot/model"
)

// RemoveDeviceCategoryByID 删除指定ID的设备类型
func (repo *MySQLRepo) RemoveDeviceCategoryByID(id uint) (affected int64, err error) {
	mod := model.DeviceCategory{Model: gorm.Model{ID: id}}
	err = repo.db.Unscoped().Delete(mod).Error
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// SaveDeviceCategory 保存设备类型
func (repo *MySQLRepo) SaveDeviceCategory(mod *model.DeviceCategory) (affected int64, err error) {
	if mod.ID != 0 {
		return repo.UpdateDeviceCategory(mod)
	} else {
		affected, _, err = repo.AddDeviceCategory(mod)
	}
	return affected, err
}

// AddDeviceCategory 增加设备类型
func (repo *MySQLRepo) AddDeviceCategory(mod *model.DeviceCategory) (affected int64, DeviceCategory *model.DeviceCategory, err error) {
	err = repo.db.Create(mod).Error
	if err != nil {
		return 0, nil, err
	}
	return 1, mod, nil
}

// UpdateDeviceCategory 修改设备类型
func (repo *MySQLRepo) UpdateDeviceCategory(mod *model.DeviceCategory) (affected int64, err error) {
	err = repo.db.Model(&model.DeviceCategory{}).Where("id = ?", mod.ID).Update(mod).Update("remark", mod.Remark).Error
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// GetDeviceCategoryByID 返回指定ID的设备类型
func (repo *MySQLRepo) GetDeviceCategoryByID(id uint) (*model.DeviceCategory, error) {
	var DeviceCategory model.DeviceCategory
	if err := repo.db.Where("id = ?", id).Find(&DeviceCategory).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &DeviceCategory, nil
}

// CountDeviceCategorys 统计满足过滤条件的设备类型数量
func (repo *MySQLRepo) CountDeviceCategorys(cond *model.DeviceCategory) (count int64, err error) {
	db := repo.db.Model(&model.DeviceCategory{})
	if cond != nil {
		db = MultiQuery(db, "category", cond.Category)
		db = MultiQuery(db, "hardware", cond.Hardware)
		db = MultiQuery(db, "power", cond.Power)
		db = MultiQuery(db, "remark", cond.Remark)
	}
	err = db.Count(&count).Error
	return count, err
}

// GetDeviceCategorys 返回满足过滤条件的设备类型
func (repo *MySQLRepo) GetDeviceCategorys(cond *model.DeviceCategory, orderby model.OrderBy, limiter *page.Limiter) (items []*model.DeviceCategory, err error) {
	items = make([]*model.DeviceCategory, 0)
	db := repo.db.Model(&model.DeviceCategory{})
	if cond != nil {
		db = MultiQuery(db, "category", cond.Category)
		db = MultiQuery(db, "hardware", cond.Hardware)
		db = MultiQuery(db, "power", cond.Power)
		db = MultiQuery(db, "remark", cond.Remark)
	}
	for i := range orderby {
		db = db.Order(orderby[i].String())
	}
	if limiter != nil {
		db = db.Offset(limiter.Offset).Limit(limiter.Limit)
	}
	err = db.Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

// 返回指定类型名的设备类型
func (repo *MySQLRepo) GetDeviceCategoryByName(category string) (dc *model.DeviceCategory, err error) {
	var DeviceCategory model.DeviceCategory
	if err := repo.db.Where("category = ?", category).Find(&DeviceCategory).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &DeviceCategory, nil
}

// 查询设备查询条件的下拉源数据
func (repo *MySQLRepo) GetDeviceCategoryQuerys(param string) (*model.DeviceQueryParamResp, error) {
	p := model.DeviceQueryParamResp{}
	p.ParamName = param
	rawSQL := ""
	mods := make([]struct {
		ID   string
		Name string
	}, 0)
	switch param {
	case "category":
		rawSQL = "select DISTINCT id AS id, `category` AS name from device_category;"
	default:
		return nil, errors.New("not supported yet")
	}

	_ = repo.db.Raw(rawSQL).Scan(&mods).Error
	if len(mods) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	for _, mod := range mods {
		p.List = append(p.List, model.ParamList{
			ID:   mod.ID,
			Name: mod.Name,
		})
	}
	return &p, nil
}
