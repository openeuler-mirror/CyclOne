package mysqlrepo

import (
	"github.com/jinzhu/gorm"
	"github.com/voidint/page"

	"idcos.io/cloudboot/model"
)

// RemoveServerRoomByID 删除指定ID的机房
func (repo *MySQLRepo) RemoveServerRoomByID(id uint) (affected int64, err error) {
	// TODO 待实现
	tx := repo.db.Begin()
	if err = tx.Unscoped().Where("id = ?", id).Delete(&model.ServerRoom{}).Error; err != nil {
		tx.Rollback()
		repo.log.Error(err)
		return affected, err
	}
	affected = tx.RowsAffected
	return affected, tx.Commit().Error
}

//UpdateServerRoom 更新机房
func (repo *MySQLRepo) UpdateServerRoom(srs []*model.ServerRoom) (affected int64, err error) {
	tx := repo.db.Begin()
	for k := range srs {
		if err = tx.Model(&model.ServerRoom{}).Where("id = ?", srs[k].ID).Updates(srs[k]).Error; err != nil {
			tx.Rollback()
			repo.log.Error(err)
			return affected, err
		}
	}
	affected = tx.RowsAffected
	return affected, tx.Commit().Error
}

// SaveServerRoom 保存机房
func (repo *MySQLRepo) SaveServerRoom(sr *model.ServerRoom) (affected int64, err error) {
	tx := repo.db.Begin()
	if err = tx.Create(sr).Error; err != nil {
		tx.Rollback()
		repo.log.Error(err)
		return affected, err
	}
	affected = tx.RowsAffected
	return affected, tx.Commit().Error
}

// UpdateServerRoomStatus 批量更新机房状态
func (repo *MySQLRepo) UpdateServerRoomStatus(status string, ids ...uint) (affected int64, err error) {
	tx := repo.db.Begin()
	if err = tx.Model(&model.ServerRoom{}).Where("id in (?)", ids).Update("status", status).Error; err != nil {
		tx.Rollback()
		repo.log.Error(err)
		return affected, err
	}
	affected = tx.RowsAffected
	return affected, tx.Commit().Error
}

// GetServerRoomByID 返回指定ID的机房
func (repo *MySQLRepo) GetServerRoomByID(id uint) (*model.ServerRoom, error) {
	var room model.ServerRoom
	if err := repo.db.Where("id = ?", id).Find(&room).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &room, nil
}

// GetServerRoomByName 返回指定Name的机房
func (repo *MySQLRepo) GetServerRoomByName(Name string) (*model.ServerRoom, error) {
	var room model.ServerRoom
	if err := repo.db.Where("name = ?", Name).Find(&room).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &room, nil
}

// CountServerRooms 统计满足过滤条件的机房数量
func (repo *MySQLRepo) CountServerRooms(cond *model.ServerRoomCond) (count int64, err error) {

	db := repo.db.Model(&model.ServerRoom{})

	if cond != nil {
		db = MultiNumQuery(db, "id", cond.ID)
		db = MultiNumQuery(db, "idc_id", cond.IDCID)
		db = MultiNumQuery(db, "first_server_room", cond.FirstServerRoom)
		db = MultiQuery(db, "name", cond.Name)
		db = MultiQuery(db, "city", cond.City)
		db = MultiQuery(db, "address", cond.Address)
		db = MultiQuery(db, "server_room_manager", cond.ServerRoomManager)
		db = MultiQuery(db, "vendor_manager", cond.VendorManager)
		db = MultiQuery(db, "network_asset_manager", cond.NetworkAssetManager)
		db = MultiQuery(db, "support_phone_number", cond.SupportPhoneNumber)
		db = MultiEnumQuery(db, "status", cond.Status)
		db = MultiQuery(db, "creator", cond.Creator)
	}

	if err = db.Select("id").Count(&count).Error; err != nil {
		repo.log.Error(err)
		return count, err
	}
	return count, nil
}

// GetServerRooms 返回满足过滤条件的机房
func (repo *MySQLRepo) GetServerRooms(cond *model.ServerRoomCond, orderby model.OrderBy, limiter *page.Limiter) (items []*model.ServerRoom, err error) {
	db := repo.db.Model(&model.ServerRoom{})

	if cond != nil {
		db = MultiNumQuery(db, "id", cond.ID)
		db = MultiNumQuery(db, "idc_id", cond.IDCID)
		db = MultiNumQuery(db, "first_server_room", cond.FirstServerRoom)
		db = MultiQuery(db, "name", cond.Name)
		db = MultiQuery(db, "city", cond.City)
		db = MultiQuery(db, "address", cond.Address)
		db = MultiQuery(db, "server_room_manager", cond.ServerRoomManager)
		db = MultiQuery(db, "vendor_manager", cond.VendorManager)
		db = MultiQuery(db, "network_asset_manager", cond.NetworkAssetManager)
		db = MultiQuery(db, "support_phone_number", cond.SupportPhoneNumber)
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
