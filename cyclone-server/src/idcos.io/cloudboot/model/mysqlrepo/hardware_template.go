package mysqlrepo

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
	"idcos.io/cloudboot/model"
)

// CountHardwareByCond 统计查询硬件模板数量
func (repo *MySQLRepo) CountHardwareByCond(cond *model.HardwareTplCond) (count int64, err error) {
	db := repo.db.Model(model.HardwareTemplate{})
	if cond.Name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", cond.Name))
	}
	if cond.Builtin != "" {
		db = db.Where("builtin = ?", cond.Builtin)
	}
	if cond.Vendor != "" {
		db = db.Where("vendor LIKE ?", fmt.Sprintf("%%%s%%", cond.Vendor))
	}
	if cond.ModelName != "" {
		db = db.Where("model LIKE ?", fmt.Sprintf("%%%s%%", cond.ModelName))
	}

	if err = db.Select("id").Count(&count).Error; err != nil {
		repo.log.Error(err)
		return 0, err
	}
	return count, nil
}

// GetHardwaresByCond 分页查询硬件模板
func (repo *MySQLRepo) GetHardwaresByCond(cond *model.HardwareTplCond, limiter *page.Limiter) (items []*model.HardwareTemplate, err error) {
	db := repo.db.Model(model.HardwareTemplate{})
	if cond.Name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", cond.Name))
	}
	if cond.Builtin != "" {
		db = db.Where("builtin = ?", cond.Builtin)
	}
	if cond.Vendor != "" {
		db = db.Where("vendor LIKE ?", fmt.Sprintf("%%%s%%", cond.Vendor))
	}
	if cond.ModelName != "" {
		db = db.Where("model LIKE ?", fmt.Sprintf("%%%s%%", cond.ModelName))
	}

	if limiter != nil {
		db = db.Limit(limiter.Limit).Offset(limiter.Offset)
	}

	if err = db.Find(&items).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return items, nil
}

// GetHardwareTemplateByName 返回指定名称的镜像安装模板
func (repo *MySQLRepo) GetHardwareTemplateByName(name string) (*model.HardwareTemplate, error) {
	var tpl model.HardwareTemplate
	if err := repo.db.Where("name = ?", name).Find(&tpl).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &tpl, nil
}

//RemoveHardwareTemplateByID 删除硬件模板配置
func (repo *MySQLRepo) RemoveHardwareTemplateByID(id uint) (affected int64, err error) {
	mod := model.HardwareTemplate{Model: gorm.Model{ID: id}}
	db := repo.db.Unscoped().Delete(mod)
	if err = db.Error; err != nil {
		repo.log.Error(err)
		return db.RowsAffected, err
	}
	return db.RowsAffected, nil
}

// GetHardwareTemplateByID 返回指定ID的镜像安装模板
func (repo *MySQLRepo) GetHardwareTemplateByID(id uint) (*model.HardwareTemplate, error) {
	var tpl model.HardwareTemplate
	if err := repo.db.Where("id = ?", id).Find(&tpl).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &tpl, nil
}

// SaveHardwareTemplate 保存、修改硬件模板
func (repo *MySQLRepo) SaveHardwareTemplate(na *model.HardwareTemplate) (id uint, err error) {
	if na.ID > 0 {
		// 更新
		db := repo.db.Model(&model.HardwareTemplate{}).Where("id = ?", na.ID).Updates(na)
		if err = db.Error; err != nil {
			repo.log.Error(err)
			return na.ID, err
		}
		return na.ID, nil
	}
	// 新增
	db := repo.db.Create(na)
	if err = db.Error; err != nil {
		repo.log.Error(err)
		return 0, err
	}
	return na.ID, nil
}
