package mysqlrepo

import (
	"idcos.io/cloudboot/model"
	"github.com/voidint/page"
	"github.com/jinzhu/gorm"
)



// SaveDeviceSettingRule 保存规则记录
func (repo *MySQLRepo) SaveDeviceSettingRule(mod *model.DeviceSettingRule) (affected int64, err error) {
	if mod.ID != 0 {
		return repo.UpdateDeviceSettingRule(mod)
	} else {
		affected, _, err = repo.AddDeviceSettingRule(mod)
	}
	return affected, err
}

// AddDeviceSettingRule 增加规则记录
func (repo *MySQLRepo) AddDeviceSettingRule(mod *model.DeviceSettingRule) (affected int64, DeviceSettingRule *model.DeviceSettingRule, err error) {
	err = repo.db.Create(mod).Error
	if err != nil {
		return 0, nil, err
	}
	return 1, mod, nil
}

// UpdateDeviceSettingRule 修改规则记录
func (repo *MySQLRepo) UpdateDeviceSettingRule(mod *model.DeviceSettingRule) (affected int64, err error) {
	err = repo.db.Model(&model.DeviceSettingRule{}).Where("id = ?", mod.ID).Update(mod).Error
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// RemoveDeviceSettingRuleByID 删除指定ID的规则记录
func (repo *MySQLRepo) RemoveDeviceSettingRuleByID(id uint) (affected int64, err error) {
	mod := model.DeviceSettingRule{Model: gorm.Model{ID: id}}
	err = repo.db.Unscoped().Delete(mod).Error
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// GetDeviceSettingRulesByType 根据规则分类查询获取所有规则
func (repo *MySQLRepo) GetDeviceSettingRulesByType(queryType string) (items []*model.DeviceSettingRule, err error) {
	if err = repo.db.Model(&model.DeviceSettingRule{}).Where("rule_category = ? ", queryType).Find(&items).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return items, nil
}

// GetDeviceSettingRuleByID 返回指定ID的规则
func (repo *MySQLRepo) GetDeviceSettingRuleByID(id uint) (*model.DeviceSettingRule, error) {
	var DeviceSettingRule model.DeviceSettingRule
	if err := repo.db.Where("id = ?", id).Find(&DeviceSettingRule).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &DeviceSettingRule, nil
}

// CountDeviceSettingRules 统计满足过滤条件的规则数量
func (repo *MySQLRepo) CountDeviceSettingRules(cond *model.DeviceSettingRule) (count int64, err error) {
	db := repo.db.Model(&model.DeviceSettingRule{})
	if cond != nil {
		db = MultiQuery(db, "rule_category", cond.RuleCategory)
	}
	err = db.Count(&count).Error
	return count, err
}

// GetDeviceSettingRules 返回满足过滤条件的规则
func (repo *MySQLRepo) GetDeviceSettingRules(cond *model.DeviceSettingRule, orderby model.OrderBy, limiter *page.Limiter) (items []*model.DeviceSettingRule, err error) {
	items = make([]*model.DeviceSettingRule, 0)
	db := repo.db.Model(&model.DeviceSettingRule{})
	if cond != nil {
		db = MultiQuery(db, "rule_category", cond.RuleCategory)
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