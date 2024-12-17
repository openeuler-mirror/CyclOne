package mysqlrepo

import (
	"github.com/voidint/page"
	"idcos.io/cloudboot/model"
)

// GetPermissionCodes 返回满足过滤条件的权限码
func (repo *MySQLRepo) GetPermissionCodes(cond *model.PermissionCode, orderby model.OrderBy, limiter *page.Limiter) (items []*model.PermissionCode, err error) {
	db := repo.db.Model(&model.PermissionCode{})
	if cond != nil {
		if cond.PID >= 0 {
			db = db.Where("pid = ?", cond.PID)
		}
		if cond.Code != "" {
			db = db.Where("code = ?", cond.Code)
		}
		if cond.Title != "" {
			db = db.Where("title LIKE ?", "%"+cond.Title+"%")
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
