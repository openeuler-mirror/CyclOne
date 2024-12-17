package mysqlrepo

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"idcos.io/cloudboot/model"
	mystrings "idcos.io/cloudboot/utils/strings"
)

// GetSystemTemplateBySN 返回指定设备所关联的系统模板
func (repo *MySQLRepo) GetSystemTemplateBySN(sn string) (*model.SystemTemplate, error) {
	var tpl model.SystemTemplate
	if err := repo.db.Joins("inner join device_setting on device_setting.system_template_id = system_template.id").Where("device_setting.sn = ?", sn).Find(&tpl).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &tpl, nil
}

// GetSystemTemplatesByCond 根据条件查询系统安装模板
func (repo *MySQLRepo) GetSystemTemplatesByCond(cond *model.CommonTemplateCond) (templates []*model.SystemTemplate, err error) {
	db := repo.db.Model(&model.CommonTemplateCond{})

	if cond.Name != "" {
		db = db.Where("name like ?", fmt.Sprintf("%s%%s%%", cond.Name))
	}

	if cond.BootMode != "" {
		db = db.Where("boot_mode = ?", cond.BootMode)
	}

	if cond.Family != "" {
		db = db.Where("family = ?", cond.Family)
	}
	
	if cond.OSLifecycle != "" {
		db = db.Where("os_lifecycle IN (?)", mystrings.MultiLines2Slice(cond.OSLifecycle))
	}
	if cond.Arch != "" {
		db = db.Where("arch = ?", cond.Arch)
	}

	if err := db.Find(&templates).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}

	return
}

// GetImageTemplatesByCond 根据条件查询镜像安装模板
func (repo *MySQLRepo) GetImageTemplatesByCond(cond *model.CommonTemplateCond) (templates []*model.ImageTemplate, err error) {
	db := repo.db.Model(&model.CommonTemplateCond{})

	if cond.Name != "" {
		db = db.Where("name like ?", fmt.Sprintf("%s%%s%%", cond.Name))
	}

	if cond.BootMode != "" {
		db = db.Where("boot_mode = ?", cond.BootMode)
	}

	if cond.Family != "" {
		db = db.Where("family = ?", cond.Family)
	}

	if cond.OSLifecycle != "" {
		db = db.Where("os_lifecycle IN (?)", mystrings.MultiLines2Slice(cond.OSLifecycle))
	}
	if cond.Arch != "" {
		db = db.Where("arch = ?", cond.Arch)
	}

	if err := db.Find(&templates).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}

	return
}
