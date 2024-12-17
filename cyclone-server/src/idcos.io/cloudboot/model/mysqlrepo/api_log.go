package mysqlrepo

import (
	"github.com/jinzhu/gorm"
	"github.com/voidint/page"

	"idcos.io/cloudboot/model"
)

// SaveAPILog 保存API记录
func (repo *MySQLRepo) SaveAPILog(operate *model.APILog) (id uint, err error) {
	if err = repo.db.Model(&model.APILog{}).Save(operate).Error; err != nil {
		repo.log.Error(err.Error())
		return 0, err
	}
	return operate.ID, nil
}

// CountAPILog 统计满足过滤条件的API记录数量
func (repo *MySQLRepo) CountAPILog(cond *model.APILogCond) (count int64, err error) {
	db := repo.db.Model(&model.APILog{})
	if cond != nil {
		db = addAPIsCond(db, cond)
	}
	if err = db.Count(&count).Error; err != nil {
		repo.log.Error(err.Error())
		return 0, err
	}

	return count, err
}

// GetAPILogByCond 返回满足过滤条件的API记录
func (repo *MySQLRepo) GetAPILogByCond(cond *model.APILogCond, orderby model.OrderBy, limiter *page.Limiter) (items []*model.APILog, err error) {
	db := repo.db.Model(&model.APILog{})
	db = addAPIsCond(db, cond)

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

// addAPIsCond 添加查询条件
func addAPIsCond(db *gorm.DB, cond *model.APILogCond) *gorm.DB {
	if cond != nil {
		if cond.Method != "" {
			db = db.Where("method = ?", cond.Method)
		}
		if cond.API != "" {
			db = db.Where("api like ?", "%"+cond.API+"%")
		}
		if cond.Description != "" {
			db = db.Where("description like ?", "%"+cond.Description+"%")
		}
		if cond.Operator != "" {
			db = db.Where("operator like ?", "%"+cond.Operator+"%")
		}
		if cond.Status != "" {
			db = db.Where("status = ?", cond.Status)
		}

		if !cond.CreatedAtStart.IsZero() && !cond.CreatedAtEnd.IsZero() {
			db = db.Where("created_at between ? and ?", cond.CreatedAtStart, cond.CreatedAtEnd)
		}

		if cond.Cost1 > 0.0 && cond.Cost2 > 0.0 {
			db = db.Where("time > ? and time < ?", cond.Cost1, cond.Cost2)
		}

	}
	return db
}
