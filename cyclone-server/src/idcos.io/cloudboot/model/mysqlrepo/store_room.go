package mysqlrepo

import (
	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
	"idcos.io/cloudboot/model"
)

// RemoveStoreRoomByID 删除指定ID的库房
func (repo *MySQLRepo) RemoveStoreRoomByID(id uint) (affected int64, err error) {
	tx := repo.db.Begin()
	if err = tx.Unscoped().Where("id = ?", id).Delete(&model.StoreRoom{}).Error; err != nil {
		tx.Rollback()
		repo.log.Error(err)
		return affected, err
	}
	affected = tx.RowsAffected
	return affected, tx.Commit().Error
}

//UpdateStoreRoom 更新库房
func (repo *MySQLRepo) UpdateStoreRoom(srs []*model.StoreRoom) (affected int64, err error) {
	tx := repo.db.Begin()
	for k := range srs {
		if err = tx.Model(&model.StoreRoom{}).Where("id = ?", srs[k].ID).Updates(srs[k]).Error; err != nil {
			tx.Rollback()
			repo.log.Error(err)
			return affected, err
		}
	}
	affected = tx.RowsAffected
	return affected, tx.Commit().Error
}

// SaveStoreRoom 保存库房
func (repo *MySQLRepo) SaveStoreRoom(sr *model.StoreRoom) (affected int64, err error) {
	tx := repo.db.Begin()
	if err = tx.Create(sr).Error; err != nil {
		tx.Rollback()
		repo.log.Error(err)
		return affected, err
	}
	affected = tx.RowsAffected
	return affected, tx.Commit().Error
}

// GetStoreRoomByID 返回指定ID的库房
func (repo *MySQLRepo) GetStoreRoomByID(id uint) (*model.StoreRoom, error) {
	var room model.StoreRoom
	if err := repo.db.Where("id = ?", id).Find(&room).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &room, nil
}

// GetStoreRoomByName 返回指定Name的库房
func (repo *MySQLRepo) GetStoreRoomByName(n string) (*model.StoreRoom, error) {
	var room model.StoreRoom
	if err := repo.db.Where("name = ?", n).First(&room).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &room, nil
}

// CountStoreRooms 统计满足过滤条件的库房数量
func (repo *MySQLRepo) CountStoreRooms(cond *model.StoreRoomCond) (count int64, err error) {

	db := repo.db.Model(&model.StoreRoom{})

	if cond != nil {
		db = MultiNumQuery(db, "id", cond.ID)
		db = MultiNumQuery(db, "idc_id", cond.IDCID)
		db = MultiQuery(db, "first_server_room", cond.FirstServerRoom)
		db = MultiQuery(db, "name", cond.Name)
		db = MultiQuery(db, "city", cond.City)
		db = MultiQuery(db, "address", cond.Address)
		db = MultiQuery(db, "store_room_manager", cond.StoreRoomManager)
		db = MultiQuery(db, "vendor_manager", cond.VendorManager)
		db = MultiEnumQuery(db, "status", cond.Status)
		db = MultiQuery(db, "creator", cond.Creator)
	}

	if err = db.Select("id").Count(&count).Error; err != nil {
		repo.log.Error(err)
		return count, err
	}
	return count, nil
}

// GetStoreRooms 返回满足过滤条件的库房
func (repo *MySQLRepo) GetStoreRooms(cond *model.StoreRoomCond, orderby model.OrderBy, limiter *page.Limiter) (items []*model.StoreRoom, err error) {
	db := repo.db.Model(&model.StoreRoom{})

	if cond != nil {
		db = MultiNumQuery(db, "id", cond.ID)
		db = MultiNumQuery(db, "idc_id", cond.IDCID)
		db = MultiQuery(db, "first_server_room", cond.FirstServerRoom)
		db = MultiQuery(db, "name", cond.Name)
		db = MultiQuery(db, "city", cond.City)
		db = MultiQuery(db, "address", cond.Address)
		db = MultiQuery(db, "store_room_manager", cond.StoreRoomManager)
		db = MultiQuery(db, "vendor_manager", cond.VendorManager)
		db = MultiEnumQuery(db, "status", cond.Status)
		db = MultiQuery(db, "creator", cond.Creator)
	}

	for k := range orderby {
		db = db.Order(orderby[k].String())
	}

	if limiter != nil {
		db = db.Limit(limiter.Limit)
		db = db.Offset(limiter.Offset)
	}

	if err = db.Find(&items).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}

	return items, nil
}

//SaveVirtualCabinet 新增虚拟货架
func (repo *MySQLRepo) SaveVirtualCabinet(sr *model.VirtualCabinet) (affected int64, err error) {
	tx := repo.db.Begin()
	if err = tx.Create(sr).Error; err != nil {
		tx.Rollback()
		repo.log.Error(err)
		return affected, err
	}
	affected = tx.RowsAffected
	return affected, tx.Commit().Error
}

// RemoveVirtualCabinetByID 删除指定ID的虚拟货架
func (repo *MySQLRepo) RemoveVirtualCabinetByID(id uint) (affected int64, err error) {
	tx := repo.db.Begin()
	if err = tx.Unscoped().Where("id = ?", id).Delete(&model.VirtualCabinet{}).Error; err != nil {
		tx.Rollback()
		repo.log.Error(err)
		return affected, err
	}
	affected = tx.RowsAffected
	return affected, tx.Commit().Error
}

// CountVirtualCabinets 统计满足过滤条件的虚拟货架数量
func (repo *MySQLRepo) CountVirtualCabinets(cond *model.VirtualCabinet) (count int64, err error) {

	db := repo.db.Model(&model.VirtualCabinet{})

	if cond != nil {
		if cond.StoreRoomID != 0 {
			//待补充
			db = db.Where("store_room_id = ?", cond.StoreRoomID)
		}
	}

	if err = db.Select("id").Count(&count).Error; err != nil {
		repo.log.Error(err)
		return count, err
	}
	return count, nil
}

// GetVirtualCabinets 返回满足过滤条件的库房
func (repo *MySQLRepo) GetVirtualCabinets(cond *model.VirtualCabinet, orderby model.OrderBy, limiter *page.Limiter) (items []*model.VirtualCabinet, err error) {
	db := repo.db.Model(&model.VirtualCabinet{})

	if cond != nil {
		if cond.StoreRoomID != 0 {
			db = db.Where("store_room_id = ?", cond.StoreRoomID)
		}
		if cond.Number != "" {
			db = db.Where("number = ?", cond.Number)
		}
	}

	for k := range orderby {
		db = db.Order(orderby[k].String())
	}

	if limiter != nil {
		db = db.Limit(limiter.Limit)
		db = db.Offset(limiter.Offset)
	}

	if err = db.Find(&items).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}

	return items, nil
}

// GetVirtualCabinetByID 返回指定ID的虚拟货架
func (repo *MySQLRepo) GetVirtualCabinetByID(id uint) (*model.VirtualCabinet, error) {
	var vc model.VirtualCabinet
	if err := repo.db.Where("id = ?", id).Find(&vc).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &vc, nil
}


// GetVirtualCabinetsByCond 返回满足过滤条件的库房
func (repo *MySQLRepo) GetVirtualCabinetsByCond(cond *model.CombinedStoreRoomVirtualCabinet, orderby model.OrderBy, limiter *page.Limiter) (items []*model.VirtualCabinet, err error) {
	db := repo.db.Model(&model.VirtualCabinet{})

	if cond != nil {
		if cond.StoreRoomName != "" {
			db = db.Joins("left join store_room on virtual_cabinet.store_room_id = store_room.id")
			db = db.Where("name = ?", cond.StoreRoomName)
		}
		if cond.StoreRoomID != 0 {
			db = db.Where("store_room_id = ?", cond.StoreRoomID)
		}
		if cond.VirtualCabinetNumber != "" {
			db = db.Where("number = ?", cond.VirtualCabinetNumber)
		}
	}

	for k := range orderby {
		db = db.Order(orderby[k].String())
	}

	if limiter != nil {
		db = db.Limit(limiter.Limit)
		db = db.Offset(limiter.Offset)
	}

	if err = db.Find(&items).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}

	return items, nil
}