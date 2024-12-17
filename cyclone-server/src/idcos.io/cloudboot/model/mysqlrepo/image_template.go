package mysqlrepo

import (
	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
	"idcos.io/cloudboot/model"
	mystrings "idcos.io/cloudboot/utils/strings"
)

// RemoveImageTemplate 删除指定ID的镜像安装模板
func (repo *MySQLRepo) RemoveImageTemplate(id uint) (affected int64, err error) {
	db := repo.db.Unscoped().Where("id = ?", id).Delete(&model.ImageTemplate{})
	if err = db.Error; err != nil {
		repo.log.Error(err)
	}
	return db.RowsAffected, db.Error
}

// SaveImageTemplate 新增镜像安装模板
func (repo *MySQLRepo) SaveImageTemplate(tpl *model.ImageTemplate) (id uint, err error) {
	// 若ID为0，则新增。
	if tpl.ID == 0 {
		if err = repo.db.Create(tpl).Error; err != nil {
			repo.log.Error(err)
			return 0, err
		}
		return 1, nil
	}
	// 若ID大于0，则更新。
	db := repo.db.Model(&model.ImageTemplate{}).Where("id = ?", tpl.ID).Updates(map[string]interface{}{
		"family":          tpl.Family,
		"name":            tpl.Name,
		"boot_mode":       tpl.BootMode,
		"username":        tpl.Username,
		"password":        tpl.Password,
		"image_url":       tpl.ImageURL,
		"partition":       tpl.Partition,
		"post_script":     tpl.PostScript,
		"pre_script":      tpl.PreScript,
		"os_lifecycle":    tpl.OSLifecycle,
		"arch":            tpl.Arch,
	})
	if err = db.Error; err != nil {
		repo.log.Error(err)
	}
	return tpl.ID, db.Error
}

// GetImageTemplateByID 返回指定id的镜像安装模板
func (repo *MySQLRepo) GetImageTemplateByID(id uint) (*model.ImageTemplate, error) {
	var tpl model.ImageTemplate
	if err := repo.db.Where("id = ?", id).Find(&tpl).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &tpl, nil
}

// CountImageTemplates 统计满足过滤条件的镜像安装模板数量
func (repo *MySQLRepo) CountImageTemplates(cond *model.ImageTemplate) (count int64, err error) {
	db := repo.db.Model(&model.ImageTemplate{})
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

// GetImageTemplates 返回满足过滤条件的镜像安装模板
func (repo *MySQLRepo) GetImageTemplates(cond *model.ImageTemplate, orderby model.OrderBy, limiter *page.Limiter) (items []*model.ImageTemplate, err error) {
	db := repo.db.Model(&model.ImageTemplate{})
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
			db = db.Where("os_lifecycle IN (?)", mystrings.MultiLines2Slice(cond.OSLifecycle))
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

// CountImageTemplateByName 统计指定名称的镜像模板
func (repo *MySQLRepo) CountImageTemplateByName(name string) (uint, error) {
	mod := model.ImageTemplate{Name: name}
	var count uint
	err := repo.db.Model(mod).Where("name = ?", name).Count(&count).Error
	return count, err
}

// CountImageTemplateByNameAndID 统计指定名称和ID的镜像模板
func (repo *MySQLRepo) CountImageTemplateByNameAndID(name string, id uint) (uint, error) {
	mod := model.ImageTemplate{}
	var count uint
	err := repo.db.Model(mod).Where("name = ? and id != ?", name, id).Count(&count).Error
	return count, err
}

// CountImageTemplate 统计镜像模板
func (repo *MySQLRepo) CountImageTemplate() (uint, error) {
	mod := model.ImageTemplate{}
	var count uint
	err := repo.db.Model(mod).Count(&count).Error
	return count, err
}

// GetImageTemplateListWithPage 分页查询镜像模板
func (repo *MySQLRepo) GetImageTemplateListWithPage(limit uint, offset uint) ([]model.ImageTemplate, error) {
	var mods []model.ImageTemplate
	err := repo.db.Limit(limit).Offset(offset).Find(&mods).Error
	return mods, err
}

// GetImageTemplateByName 查询指定名称的镜像模板
func (repo *MySQLRepo) GetImageTemplateByName(name string) (*model.ImageTemplate, error) {
	var mod model.ImageTemplate
	err := repo.db.Where("name = ?", name).Find(&mod).Error
	return &mod, err
}

// GetImageTemplateIDByName 查询指定模板名称的镜像模板ID
func (repo *MySQLRepo) GetImageTemplateIDByName(name string) (uint, error) {
	mod := model.ImageTemplate{Name: name}
	err := repo.db.Where("name = ?", name).Find(&mod).Error
	return mod.ID, err
}

// GetImageTemplateBySN 查询指定SN的镜像模板
func (repo *MySQLRepo) GetImageTemplateBySN(sn string) (*model.ImageTemplate, error) {
	var mod model.ImageTemplate
	if err := repo.db.Joins("inner join device_setting on device_setting.image_template_id = image_template.id").Where("device_setting.sn = ?", sn).Find(&mod).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &mod, nil
}
