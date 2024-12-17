package mysqlrepo

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/voidint/page"

	"fmt"

	"idcos.io/cloudboot/model"
)

// RemoveServerCabinetByID 删除指定ID的机架(柜)
func (repo *MySQLRepo) RemoveServerCabinetByID(id uint) (affected int64, err error) {
	db := repo.db.Unscoped().Delete(&model.ServerCabinet{}, "id = ?", id)
	if err = db.Error; err != nil {
		repo.log.Error(err)
		return 0, err
	}
	return db.RowsAffected, nil
}

//GetServerCabinetCountByServerRoomID 根据给定的机房获取机房内的机架(柜)数
func (repo *MySQLRepo) GetServerCabinetCountByServerRoomID(id uint) (count int64, err error) {

	if err := repo.db.Model(&model.ServerCabinet{}).Where("server_room_id = ?", id).Select("id").Count(&count).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return 0, err
	}
	return count, nil
}

// SaveServerCabinet 保存机架(柜)
func (repo *MySQLRepo) SaveServerCabinet(cabinet *model.ServerCabinet) (affected int64, err error) {
	if cabinet.ID > 0 {
		// 更新
		db := repo.db.Model(&model.ServerCabinet{}).Where("id = ?", cabinet.ID).Updates(cabinet)
		if err = db.Error; err != nil {
			repo.log.Error(err)
			return db.RowsAffected, err
		}
		return db.RowsAffected, nil
	}
	// 新增
	db := repo.db.Create(cabinet)
	if err = db.Error; err != nil {
		repo.log.Error(err)
		return db.RowsAffected, err
	}
	return db.RowsAffected, nil
}

// GetServerCabinetByID 返回指定ID的机架(柜)
func (repo *MySQLRepo) GetServerCabinetByID(id uint) (*model.ServerCabinet, error) {
	var cabinet model.ServerCabinet
	if err := repo.db.Where("id = ?", id).Find(&cabinet).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &cabinet, nil
}

// GetServerCabinetByNumber 返回指定编号的机架(柜)
func (repo *MySQLRepo) GetServerCabinetByNumber(serverRoomID uint, number string) (*model.ServerCabinet, error) {
	var cabinet model.ServerCabinet
	if err := repo.db.Where("server_room_id = ? AND number = ? ", serverRoomID, number).Find(&cabinet).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &cabinet, nil
}

//UpdateServerCabinetStatus 批量修改机架(柜)状态
func (repo *MySQLRepo) UpdateServerCabinetStatus(ids []uint, status string) (affected int64, err error) {
	tx := repo.db.Begin()
	time := time.Now()
	//if err = tx.Model(&model.ServerCabinet{}).Where("id in (?)", ids).Update(&model.ServerCabinet{
	//	EnableTime: &time,
	//	Status:     status,
	//	IsEnabled:  isEnabled,
	//}).Error; err != nil {
	//	tx.Rollback()
	//	repo.log.Error(err)
	//	return affected, err
	//}
	// 针对已锁定的机柜，将其关联的非已使用的机位状态修改为不可用
	if status == model.CabinetStatLocked {
		if err = tx.Model(&model.ServerCabinet{}).Where("id in (?)", ids).Update(&model.ServerCabinet{
			Status:     status,
		}).Error; err != nil {
			tx.Rollback()
			repo.log.Error(err)
			return affected, err
		}
		db := repo.db.Model(&model.ServerUSite{}).Where("server_cabinet_id in (?)", ids).Where("status != (?)", model.USiteStatUsed).Update("status", model.USiteStatDisabled)
		if err = db.Error; err != nil {
			repo.log.Error(err)
			tx.Rollback()
			return affected, err
		}
	} else {
		isEnabled := model.NO
		switch status {
		case model.CabinetStatEnabled:
			isEnabled = model.YES
		case model.CabinetStatNotEnabled:
			isEnabled = model.NO
		}
		if err = tx.Model(&model.ServerCabinet{}).Where("id in (?)", ids).Update(&model.ServerCabinet{
			EnableTime: &time,
			Status:     status,
			IsEnabled:  isEnabled,
		}).Error; err != nil {
			tx.Rollback()
			repo.log.Error(err)
			return affected, err
		}
	}
	affected = tx.RowsAffected
	return affected, tx.Commit().Error
}

// GetServerCabinetID 根据条件返回机架(柜)ID,如果不存在则返回0
func (repo *MySQLRepo) GetServerCabinetID(cond *model.ServerCabinet) (id []uint, err error) {
	db := repo.db.Model(&model.ServerCabinet{})
	if cond != nil {
		if cond.IDCID > 0 {
			db = db.Where("idc_id = ?", cond.IDCID)
		}
		if cond.ServerRoomID > 0 {
			db = db.Where("server_room_id = ?", cond.ServerRoomID)
		}
		if cond.NetworkAreaID > 0 {
			db = db.Where("network_area_id = ?", cond.NetworkAreaID)
		}
		//if cond.Number != "" {
		//	db = db.Where("number = ?", cond.Number)
		//}
		//if cond.Type != "" {
		//	db = db.Where("type = ?", cond.Type)
		//}
		//if cond.Status != "" {
		//	db = db.Where("status = ?", cond.Status)
		//}
		//if cond.IsEnabled != "" {
		//	db = db.Where("is_enabled = ?", cond.IsEnabled)
		//}
		//if cond.IsPowered != "" {
		//	db = db.Where("is_powered = ?", cond.IsPowered)
		//}
		db = MultiQuery(db, "number", cond.Number)
		db = MultiEnumQuery(db, "type", cond.Type)
		db = MultiEnumQuery(db, "status", cond.Status)
		db = MultiEnumQuery(db, "is_enabled", cond.IsEnabled)
		db = MultiEnumQuery(db, "is_powered", cond.IsPowered)
	}
	if err = db.Pluck("id", &id).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return id, nil
}

// CountServerCabinets 统计满足过滤条件的机架(柜)数量
func (repo *MySQLRepo) CountServerCabinets(cond *model.ServerCabinetCond) (count int64, err error) {
	db := repo.db.Model(&model.ServerCabinet{})
	db = addServerCabinetCond(db, cond)
	if err = db.Select("id").Count(&count).Error; err != nil {
		repo.log.Error(err)
		return count, err
	}
	return count, nil
}

func addServerCabinetCond(db *gorm.DB, cond *model.ServerCabinetCond) *gorm.DB {
	if cond.ServerRoomName != "" {
		db = db.Joins("left join server_room on server_room.id = server_cabinet.server_room_id")
	}
	if cond.NetworkAreaName != "" {
		db = db.Joins("left join network_area on network_area.id = server_cabinet.network_area_id")
	}
	if cond != nil {
		db = MultiNumQuery(db, "server_cabinet.idc_id", cond.IDCID)
		db = MultiNumQuery(db, "server_cabinet.server_room_id", cond.ServerRoomID)
		db = MultiNumQuery(db, "server_cabinet.id", cond.ServerCabinetID)
		db = MultiQuery(db, "server_room.name", cond.ServerRoomName)
		db = MultiNumQuery(db, "network_area_id", cond.NetworkAreaID)
		db = MultiQuery(db, "network_area.name", cond.NetworkAreaName)
		db = MultiQuery(db, "number", cond.Number)
		db = MultiEnumQuery(db, "type", cond.Type)
		db = MultiEnumQuery(db, "status", cond.Status)
		db = MultiEnumQuery(db, "is_enabled", cond.IsEnabled)
		db = MultiEnumQuery(db, "is_powered", cond.IsPowered)
	}
	return db
}

// GetServerCabinets 返回满足过滤条件的机架(柜)
func (repo *MySQLRepo) GetServerCabinets(cond *model.ServerCabinetCond, orderby model.OrderBy, limiter *page.Limiter) (items []*model.ServerCabinet, err error) {
	db := repo.db.Model(&model.ServerCabinet{})
	db = addServerCabinetCond(db, cond)
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

//按照可用机位数倒序的方式查询指定机房下的机架列表
func (repo *MySQLRepo) GetCabinetOrderByFreeUsites(req *model.ServerCabinet, physicalArea string) (mod []*model.OrderedCabientResp, err error) {
	mod = make([]*model.OrderedCabientResp, 0)
	rawSQL := fmt.Sprintf(`select server_cabinet.*,count(server_usite.number) as available_usite_count 
	from server_cabinet right join server_usite on server_cabinet.id = server_usite.server_cabinet_id 
	where server_usite.status = 'free' and server_cabinet.idc_id = %d 
	and server_cabinet.server_room_id = %d 
	and server_usite.physical_area = '%s' 
	group by server_cabinet.id order by available_usite_count desc;`,
		req.IDCID,
		req.ServerRoomID,
		physicalArea)

	err = repo.db.Raw(rawSQL).Scan(&mod).Error
	return
}

//PowerOnServerCabinetByID 根据ID将机架(柜)上电
func (repo *MySQLRepo) PowerOnServerCabinetByID(ids []uint) (affected int64, err error) {
	time := time.Now()
	db := repo.db.Model(&model.ServerCabinet{}).Where("id in (?)", ids).Update(&model.ServerCabinet{
		PowerOnTime: &time,
		IsPowered:   model.CabinetPowerOn,
	})
	if db.Error != nil {
		repo.log.Error(db.Error)
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

//PowerOffServerCabinetByID 根据ID将机架(柜)下电
func (repo *MySQLRepo) PowerOffServerCabinetByID(id uint) (affected int64, err error) {
	time := time.Now()
	db := repo.db.Model(&model.ServerCabinet{}).Where("id = ?", id).Update(&model.ServerCabinet{
		PowerOffTime: &time,
		IsPowered:    model.CabinetPowerOff,
	})
	if db.Error != nil {
		repo.log.Error(db.Error)
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

//UpdateServerCabinetType 批量修改机架(柜)类型
func (repo *MySQLRepo) UpdateServerCabinetType(ids []uint, typ string) (affected int64, err error) {
	tx := repo.db.Begin()
	if err = tx.Model(&model.ServerCabinet{}).Where("id in (?)", ids).Update(&model.ServerCabinet{
		Type:     typ,
	}).Error; err != nil {
		tx.Rollback()
		repo.log.Error(err)
		return affected, err
	}
	affected = tx.RowsAffected
	return affected, tx.Commit().Error
}

//UpdateServerCabinetRemark 批量修改机架(柜)备注信息
func (repo *MySQLRepo) UpdateServerCabinetRemark(ids []uint, remark string) (affected int64, err error) {
	tx := repo.db.Begin()
	if err = tx.Model(&model.ServerCabinet{}).Where("id in (?)", ids).Update(&model.ServerCabinet{
		Remark:     remark,
	}).Error; err != nil {
		tx.Rollback()
		repo.log.Error(err)
		return affected, err
	}
	affected = tx.RowsAffected
	return affected, tx.Commit().Error
}