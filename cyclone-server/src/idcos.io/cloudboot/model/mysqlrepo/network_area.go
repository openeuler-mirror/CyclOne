package mysqlrepo

import (
	"github.com/jinzhu/gorm"
	"github.com/voidint/page"

	"idcos.io/cloudboot/model"
)

// RemoveNetworkAreaByID 删除指定ID的网络区域
func (repo *MySQLRepo) RemoveNetworkAreaByID(id uint) (affected int64, err error) {
	db := repo.db.Unscoped().Delete(&model.NetworkArea{}, "id = ?", id)
	if err = db.Error; err != nil {
		repo.log.Error(err)
		return 0, err
	}
	return db.RowsAffected, nil
}

// SaveNetworkArea 保存网络区域
func (repo *MySQLRepo) SaveNetworkArea(na *model.NetworkArea) (affected int64, err error) {
	if na.ID > 0 {
		// 更新
		db := repo.db.Model(&model.NetworkArea{}).Where("id = ?", na.ID).Updates(na)
		if err = db.Error; err != nil {
			repo.log.Error(err)
			return db.RowsAffected, err
		}
		return db.RowsAffected, nil
	}
	// 新增
	db := repo.db.Create(na)
	if err = db.Error; err != nil {
		repo.log.Error(err)
		return db.RowsAffected, err
	}
	return db.RowsAffected, nil
}

// UpdateNetworkAreaStatus 批量更新网络区域状态
func (repo *MySQLRepo) UpdateNetworkAreaStatus(status string, ids ...uint) (affected int64, err error) {
	db := repo.db.Model(&model.NetworkArea{}).Where("id IN(?)", ids).Update("status", status)
	if err = db.Error; err != nil {
		repo.log.Error(err)
		return 0, err
	}
	return db.RowsAffected, nil
}

// GetNetworkAreaByID 返回指定ID的网络区域
func (repo *MySQLRepo) GetNetworkAreaByID(id uint) (*model.NetworkArea, error) {
	var na model.NetworkArea
	if err := repo.db.Where("id = ?", id).Find(&na).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &na, nil
}

// CountNetworkAreas 统计满足过滤条件的网络区域数量
func (repo *MySQLRepo) CountNetworkAreas(cond *model.NetworkAreaCond) (count int64, err error) {
	db := repo.db.Model(&model.NetworkArea{})
	if cond != nil {
		db = MultiNumQuery(db, "idc_id", cond.IDCID)
		db = MultiNumQuery(db, "server_room_id", cond.ServerRoomID)
		db = MultiQuery(db, "name", cond.Name)
		db = MultiQuery(db, "physical_area", cond.PhysicalArea)
		db = MultiEnumQuery(db, "status", cond.Status)
	}
	if err = db.Select("id").Count(&count).Error; err != nil {
		repo.log.Error(err)
		return count, err
	}
	return count, nil
}

// GetNetworkAreasByCond 返回满足过滤条件的网络区域(不支持模糊查找)
func (repo *MySQLRepo) GetNetworkAreasByCond(cond *model.NetworkArea) (item []*model.NetworkArea, err error) {
	db := repo.db.Model(&model.NetworkArea{})
	if cond != nil {
		if cond.IDCID > 0 {
			db = db.Where("idc_id = ?", cond.IDCID)
		}
		if cond.ServerRoomID > 0 {
			db = db.Where("server_room_id = ?", cond.ServerRoomID)
		}
		if cond.Name != "" {
			db = db.Where("name = ?", cond.Name)
		}
	}

	if err = db.Find(&item).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return item, nil
}

func addNetAreaCond(db *gorm.DB, cond *model.NetworkAreaCond) *gorm.DB {
	if cond.ServerRoomName != "" {
		db = db.Joins("left join server_room on server_room.id = network_area.server_room_id")
	}
	if cond != nil {
		db = MultiNumQuery(db, "idc_id", cond.IDCID)
		db = MultiNumQuery(db, "server_room_id", cond.ServerRoomID)
		db = MultiQuery(db, "name", cond.Name)
		db = MultiQuery(db, "server_room.name", cond.ServerRoomName)
		db = MultiQuery(db, "physical_area", cond.PhysicalArea)
		db = MultiEnumQuery(db, "status", cond.Status)
	}
	return db
}

// GetNetworkAreas 返回满足过滤条件的网络区域
func (repo *MySQLRepo) GetNetworkAreas(cond *model.NetworkAreaCond, orderby model.OrderBy, limiter *page.Limiter) (items []*model.NetworkArea, err error) {
	db := repo.db.Model(&model.NetworkArea{})

	db = addNetAreaCond(db, cond)

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
