package mysqlrepo

import (
	"github.com/jinzhu/gorm"

	"idcos.io/cloudboot/model"
)

// GetPlatformConfigByName 返回指定名称的配置对象
func (repo *MySQLRepo) GetPlatformConfigByName(name string) (*model.PlatformConfig, error) {
	var pc model.PlatformConfig
	if err := repo.db.Where("name = ?", name).Find(&pc).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &pc, nil
}
