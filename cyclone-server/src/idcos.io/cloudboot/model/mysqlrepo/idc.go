package mysqlrepo

import (
	"github.com/voidint/page"

	"fmt"

	"github.com/jinzhu/gorm"
	"idcos.io/cloudboot/model"
)

// RemoveIDCByID 删除指定ID的数据中心
func (repo *MySQLRepo) RemoveIDCByID(id uint) (affected int64, err error) {
	mod := model.IDC{Model: gorm.Model{ID: id}}
	err = repo.db.Unscoped().Delete(mod).Error
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// SaveIDC 保存数据中心
func (repo *MySQLRepo) SaveIDC(*model.IDC) (affected int64, err error) {
	// TODO 待实现
	return 0, nil
}

// AddIDC 增加数据中心
func (repo *MySQLRepo) AddIDC(mod *model.IDC) (affected int64, idc *model.IDC, err error) {
	err = repo.db.Create(mod).Error
	if err != nil {
		return 0, nil, err
	}
	return 1, mod, nil
}

// UpdateIDC 修改数据中心
func (repo *MySQLRepo) UpdateIDC(mod *model.IDC) (affected int64, err error) {
	err = repo.db.Model(&model.IDC{}).Where("id = ?", mod.ID).Update(mod).Error
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// UpdateIDCStatus 批量更新数据中心状态
func (repo *MySQLRepo) UpdateIDCStatus(status string, ids ...uint) (affected int64, err error) {
	tx := repo.db.Begin()
	for _, id := range ids {
		err = tx.Model(&model.IDC{}).Where("id = ?", id).Update("status", status).Error
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		affected++
	}
	tx.Commit()
	return
}

// GetIDCByName 返回指定Name的数据中心
func (repo *MySQLRepo) GetIDCByName(name string) (*model.IDC, error) {
	var idc model.IDC
	if err := repo.db.Where("name = ?", name).Find(&idc).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &idc, nil
}

// GetIDCByID 返回指定ID的数据中心
func (repo *MySQLRepo) GetIDCByID(id uint) (*model.IDC, error) {
	var idc model.IDC
	if err := repo.db.Where("id = ?", id).Find(&idc).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &idc, nil
}

// CountIDCs 统计满足过滤条件的数据中心数量
func (repo *MySQLRepo) CountIDCs(cond *model.IDC) (count int64, err error) {
	db := repo.db.Model(&model.IDC{})
	if cond != nil {
		db = MultiQuery(db, "name", cond.Name)
		db = MultiQuery(db, "usage", cond.Usage)
		if cond.FirstServerRoom != "" {
			db = db.Where("first_server_room LIKE ?", fmt.Sprintf("%%%s%%", cond.FirstServerRoom))
		}
		db = MultiQuery(db, "first_server_room", cond.FirstServerRoom)
		db = MultiEnumQuery(db, "status", cond.Status)
		db = MultiEnumQuery(db, "vendor", cond.Vendor)
	}
	err = db.Count(&count).Error
	return count, err
}

// GetIDCs 返回满足过滤条件的数据中心
func (repo *MySQLRepo) GetIDCs(cond *model.IDC, orderby model.OrderBy, limiter *page.Limiter) (items []*model.IDC, err error) {
	items = make([]*model.IDC, 0)
	db := repo.db.Model(&model.IDC{})
	if cond != nil {
		db = MultiQuery(db, "name", cond.Name)
		db = MultiEnumQuery(db, "usage", cond.Usage)
		//if cond.FirstServerRoom != "" {
		//	db = db.Where("first_server_room LIKE ?", fmt.Sprintf("%%%s%%", cond.FirstServerRoom))
		//}
		db = MultiQuery(db, "first_server_room", cond.FirstServerRoom)
		db = MultiEnumQuery(db, "status", cond.Status)
		db = MultiQuery(db, "vendor", cond.Vendor)
	}
	for i := range orderby {
		db = db.Order(orderby[i].String())
	}
	db = db.Offset(limiter.Offset).Limit(limiter.Limit)
	err = db.Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}
