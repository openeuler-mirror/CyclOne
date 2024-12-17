package mysqlrepo

import (
	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
	"idcos.io/cloudboot/model"
	mystrings "idcos.io/cloudboot/utils/strings"
)

// RemoveSystemTemplate 删除指定ID的系统安装模板
func (repo *MySQLRepo) RemoveSystemTemplate(id uint) (affected int64, err error) {
	db := repo.db.Unscoped().Where("id = ?", id).Delete(&model.SystemTemplate{})
	if err = db.Error; err != nil {
		repo.log.Error(err)
	}
	return db.RowsAffected, db.Error
}

// SaveSystemTemplate 保存系统安装模板
func (repo *MySQLRepo) SaveSystemTemplate(tpl *model.SystemTemplate) (id uint, err error) {
	// 若ID为0，则新增。
	if tpl.ID == 0 {
		if err = repo.db.Create(tpl).Error; err != nil {
			repo.log.Error(err)
			return 0, err
		}
		return 1, nil
	}
	// 若ID大于0，则更新。
	db := repo.db.Model(&model.SystemTemplate{}).Where("id = ?", tpl.ID).Updates(map[string]interface{}{
		"family":          tpl.Family,
		"name":            tpl.Name,
		"boot_mode":       tpl.BootMode,
		"username":        tpl.Username,
		"password":        tpl.Password,
		"content":         tpl.Content,
		"pxe":             tpl.PXE,
		"os_lifecycle":    tpl.OSLifecycle,
		"arch":            tpl.Arch,
	})
	if err = db.Error; err != nil {
		repo.log.Error(err)
	}
	return tpl.ID, db.Error
}

// GetSystemTemplateByID 返回指定ID的系统安装模板
func (repo *MySQLRepo) GetSystemTemplateByID(id uint) (*model.SystemTemplate, error) {
	var tpl model.SystemTemplate
	if err := repo.db.Where("id = ?", id).Find(&tpl).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &tpl, nil
}

// CountSystemTemplates 统计满足过滤条件的系统安装模板数量
func (repo *MySQLRepo) CountSystemTemplates(cond *model.SystemTemplate) (count int64, err error) {
	db := repo.db.Model(&model.SystemTemplate{})
	if cond != nil {
		if cond.Family != "" {
			db = db.Where("family = ?", cond.Family)
		}
		if cond.Name != "" {
			db = db.Where("name LIKE ?", "%"+cond.Name+"%")
		}
		if cond.BootMode == model.BootModeBIOS || cond.BootMode == model.BootModeUEFI {
			db = db.Where("boot_mode = ?", cond.BootMode)
		}
		if cond.OSLifecycle != "" {
			db = db.Where("os_lifecycle IN (?)", mystrings.MultiLines2Slice(cond.OSLifecycle))
		}
		if cond.Arch != "" {
			db = db.Where("arch = ?", cond.Arch)
		}		
	}

	if err = db.Select("id").Count(&count).Error; err != nil {
		repo.log.Error(err)
		return count, err
	}
	return count, nil
}

// GetSystemTemplates 返回满足过滤条件的系统安装模板
func (repo *MySQLRepo) GetSystemTemplates(cond *model.SystemTemplate, orderby model.OrderBy, limiter *page.Limiter) (items []*model.SystemTemplate, err error) {
	db := repo.db.Model(&model.SystemTemplate{})
	if cond != nil {
		if cond.Family != "" {
			db = db.Where("family = ?", cond.Family)
		}
		if cond.Name != "" {
			db = db.Where("name LIKE ?", "%"+cond.Name+"%")
		}
		if cond.BootMode != "" {
			db = db.Where("boot_mode = ?", cond.BootMode)
		}
		if cond.OSLifecycle != "" {
			db = db.Where("os_lifecycle IN (?)", cond.OSLifecycle)
		}
		if cond.Arch != "" {
			db = db.Where("arch = ?", cond.Arch)
		}	
	}

	for i := range orderby {
		db = db.Order(orderby[i].String())
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

// CountSystemTemplateByName 统计指定名称的系统模板
func (repo *MySQLRepo) CountSystemTemplateByName(name string) (uint, error) {
	mod := model.SystemTemplate{Name: name}
	var count uint
	err := repo.db.Model(mod).Where("name = ?", name).Count(&count).Error
	return count, err
}

// CountSystemTemplateByNameAndID 统计指定名称和ID的系统模板
func (repo *MySQLRepo) CountSystemTemplateByNameAndID(name string, id uint) (uint, error) {
	mod := model.SystemTemplate{}
	var count uint
	err := repo.db.Model(mod).Where("name = ? and id != ?", name, id).Count(&count).Error
	return count, err
}

// CountSystemTemplate 统计系统模板
func (repo *MySQLRepo) CountSystemTemplate() (uint, error) {
	mod := model.SystemTemplate{}
	var count uint
	err := repo.db.Model(mod).Count(&count).Error
	return count, err
}

// CountSystemTemplateByShield 统计指定shield系统模板
func (repo *MySQLRepo) CountSystemTemplateByShield(cond *model.SystemTemplate) (uint, error) {
	mod := model.SystemTemplate{}
	var count uint
	db := repo.db
	if cond.Name != "" {
		db.Where("name NOT like ?", cond.Name)
	}
	err := repo.db.Model(mod).Count(&count).Error
	return count, err
}

// GetSystemTemplateListWithPage 分页查询系统模板
func (repo *MySQLRepo) GetSystemTemplateListWithPage(limit uint, offset uint) ([]model.SystemTemplate, error) {
	var mods []model.SystemTemplate
	err := repo.db.Limit(limit).Offset(offset).Find(&mods).Error
	return mods, err
}

// GetSystemTemplateListWithPageAndShield 分页查询指定shield系统模板
func (repo *MySQLRepo) GetSystemTemplateListWithPageAndShield(limit uint,
	offset uint, cond *model.SystemTemplate) ([]model.SystemTemplate, error) {
	var mods []model.SystemTemplate
	db := repo.db
	if cond.Name != "" {
		db = db.Where("name NOT LIKE ?", "%"+cond.Name+"%")
	}
	err := db.Limit(limit).Offset(offset).Find(&mods).Error
	return mods, err
}

// GetSystemTemplateIDByName 查询指定名称的模板ID
func (repo *MySQLRepo) GetSystemTemplateIDByName(name string) (uint, error) {
	mod := model.SystemTemplate{Name: name}
	err := repo.db.Where("name = ?", name).Find(&mod).Error
	return mod.ID, err
}

// GetSystemTemplateByName 查询指定名称的系统模板
func (repo *MySQLRepo) GetSystemTemplateByName(name string) (*model.SystemTemplate, error) {
	var mod model.SystemTemplate
	err := repo.db.Where("name = ?", name).Find(&mod).Error
	return &mod, err
}
