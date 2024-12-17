package mysqlrepo

import (
	"github.com/voidint/page"

	"fmt"

	"strconv"

	"github.com/jinzhu/gorm"
	"idcos.io/cloudboot/model"
)

// RemoveOrderByID 删除指定ID的订单
func (repo *MySQLRepo) RemoveOrderByID(id uint) (affected int64, err error) {
	mod := model.Order{Model: gorm.Model{ID: id}}
	err = repo.db.Unscoped().Delete(mod).Error
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// SaveOrder 保存订单
func (repo *MySQLRepo) SaveOrder(mod *model.Order) (affected int64, err error) {
	//if mod.ID != 0 {
	//	return repo.UpdateOrder(mod)
	//} else {
	//	affected, _, err = repo.AddOrder(mod)
	//}
	err = repo.db.Save(mod).Error
	return
}

// AddOrder 增加订单
func (repo *MySQLRepo) AddOrder(mod *model.Order) (affected int64, Order *model.Order, err error) {
	err = repo.db.Create(mod).Error
	if err != nil {
		return 0, nil, err
	}
	return 1, mod, nil
}

// UpdateOrder 修改订单
func (repo *MySQLRepo) UpdateOrder(mod *model.Order) (affected int64, err error) {
	err = repo.db.Model(&model.Order{}).Where("id = ?", mod.ID).Update(mod).Error
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// UpdateOrderStatus 批量更新订单状态
func (repo *MySQLRepo) UpdateOrderStatus(status string, ids ...uint) (affected int64, err error) {
	tx := repo.db.Begin()
	for _, id := range ids {
		err = tx.Model(&model.Order{}).Where("id = ?", id).Update("status", status).Error
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		affected++
	}
	tx.Commit()
	return
}

// GetOrderByID 返回指定ID的订单
func (repo *MySQLRepo) GetOrderByID(id uint) (*model.Order, error) {
	var Order model.Order
	if err := repo.db.Where("id = ?", id).Find(&Order).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &Order, nil
}

// GetOrderByNumber 返回指定订单号的的订单
func (repo *MySQLRepo) GetOrderByNumber(n string) (*model.Order, error) {
	var Order model.Order
	if err := repo.db.Where("number = ?", n).Find(&Order).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &Order, nil
}

// CountOrders 统计满足过滤条件的订单数量
func (repo *MySQLRepo) CountOrders(cond *model.OrderCond) (count int64, err error) {
	db := repo.db.Model(&model.Order{})
	if cond != nil {
		db = MultiNumQuery(db, "id", cond.ID)
		db = MultiQuery(db, "number", cond.Number)
		db = MultiMatchWithSpaceQuery(db, "physical_area", cond.PhysicalArea)
		db = MultiQuery(db, "usage", cond.Usage)
		db = MultiEnumQuery(db, "status", cond.Status)
	}
	if err = db.Count(&count).Error; err != nil {
		repo.log.Error(err.Error())
		return
	}
	return count, err
}

// GetOrders 返回满足过滤条件的订单
func (repo *MySQLRepo) GetOrders(cond *model.OrderCond, orderby model.OrderBy, limiter *page.Limiter) (items []*model.Order, err error) {
	items = make([]*model.Order, 0)
	db := repo.db.Model(&model.Order{})
	if cond != nil {
		db = MultiNumQuery(db, "id", cond.ID)
		db = MultiQuery(db, "number", cond.Number)
		db = MultiMatchWithSpaceQuery(db, "physical_area", cond.PhysicalArea)
		db = MultiQuery(db, "usage", cond.Usage)
		db = MultiEnumQuery(db, "status", cond.Status)
	}
	for i := range orderby {
		db = db.Order(orderby[i].String())
	}
	if limiter != nil {
		db = db.Offset(limiter.Offset).Limit(limiter.Limit)
	}
	err = db.Find(&items).Error
	if err != nil {
		repo.log.Error(err.Error())
		return nil, err
	}
	return items, nil
}

//GetMaxOrderNumber 如果某日最大订单号为IDC120190617009，则返回最后的自增值9
func (repo *MySQLRepo) GetMaxOrderNumber(date string) (orderNumber int, err error) {
	items := make([]*model.Order, 0)
	db := repo.db.Model(&model.Order{})
	err = db.Where("number LIKE ? ", fmt.Sprintf("%%%s%%", date)).Find(&items).Error

	//自增号取最后3位数
	for _, item := range items {
		seqStr := item.Number[len(item.Number)-3:]
		seqInt, err := strconv.Atoi(seqStr)
		if err != nil {
			repo.log.Errorf("strconv:%s fail, %v", seqStr, err)
			return orderNumber, err
		}
		if seqInt > orderNumber {
			orderNumber = seqInt
		}
	}
	return
}
