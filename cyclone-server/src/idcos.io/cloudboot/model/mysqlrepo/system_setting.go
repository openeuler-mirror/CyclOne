package mysqlrepo

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"

	"idcos.io/cloudboot/model"
)

// GetSystemSetting 查询指定名称的系统配置
func (repo *MySQLRepo) GetSystemSetting(key string) (*model.SystemSetting, error) {
	var sett model.SystemSetting
	if err := repo.db.Where("`key` = ?", key).Find(&sett).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &sett, nil
}

// GetSystemSetting4InstallatonTimeout 查询安装超时时间的系统设置值。若发生错误或者值不存在，则返回默认值。
func (repo *MySQLRepo) GetSystemSetting4InstallatonTimeout(defValue int64) (sec int64) {
	sett, err := repo.GetSystemSetting(model.SysSettingInstallationTimeout)
	if err != nil || sett == nil || sett.Value == "" {
		return defValue
	}
	sec, err = strconv.ParseInt(strings.TrimSpace(sett.Value), 10, 64)
	if err != nil {
		return defValue
	}
	return sec
}

// GetSystemSetting4AuthorizationAPIs 查询授权API配置
func (repo *MySQLRepo) GetSystemSetting4AuthorizationAPIs() (items []*model.AuthorizationAPI, err error) {
	value, err := repo.GetSystemSetting(model.SysSettingAuthorization)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	if err = json.Unmarshal([]byte(value.Value), &items); err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return items, nil
}
