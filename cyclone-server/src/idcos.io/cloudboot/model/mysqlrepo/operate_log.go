package mysqlrepo

import (
	"github.com/jinzhu/gorm"
	"github.com/voidint/page"

	"idcos.io/cloudboot/model"
)

// SaveOperateLog 保存操作记录
func (repo *MySQLRepo) SaveOperateLog(operate *model.OperateLog) (id uint, err error) {
	if err = repo.db.Model(&model.OperateLog{}).Save(operate).Error; err != nil {
		return 0, err
	}
	return operate.ID, nil
}

// CountOperateLog 统计满足过滤条件的操作记录数量
func (repo *MySQLRepo) CountOperateLog(cond *model.OperateLog) (count int64, err error) {
	db := repo.db.Model(&model.OperateLog{})
	if cond != nil {
		db = addOperatesCond(db, cond)
	}
	err = db.Count(&count).Error
	return count, err
}

// GetOperateLogByCond 返回满足过滤条件的操作记录
func (repo *MySQLRepo) GetOperateLogByCond(cond *model.OperateLog, orderby model.OrderBy, limiter *page.Limiter) (items []*model.OperateLog, err error) {
	db := repo.db.Model(&model.OperateLog{})
	db = addOperatesCond(db, cond)

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

// addOperatesCond 添加查询条件
func addOperatesCond(db *gorm.DB, cond *model.OperateLog) *gorm.DB {
	if cond != nil {
		if cond.HTTPMethod != "" {
			db = db.Where("http_method = ?", cond.HTTPMethod)
		}
		if cond.URL != "" {
			db = db.Where("url like ?", "%"+cond.URL+"%")
		}
		if cond.Source != "" {
			db = db.Where("source like ?", "%"+cond.Source+"%")
		}
		if cond.Destination != "" {
			db = db.Where("destination like ?", "%"+cond.Destination+"%")
		}
		if cond.CategoryName !=""{
			db = db.Where("category_name like ?", "%"+cond.CategoryName+"%")
		}
	}
	return db
}
