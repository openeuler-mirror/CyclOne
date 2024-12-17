package mysqlrepo

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
	"idcos.io/cloudboot/utils"

	"idcos.io/cloudboot/model"
)

//GetServerUSiteCountByServerCabinetID 根据给定的机架获取机架上的机位数
func (repo *MySQLRepo) GetServerUSiteCountByServerCabinetID(id uint) (count int64, err error) {

	if err := repo.db.Model(&model.ServerUSite{}).Where("server_cabinet_id = ?", id).Select("id").Count(&count).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return 0, err
	}
	return count, nil
}

// SaveServerUSite 保存机位(U位)
func (repo *MySQLRepo) SaveServerUSite(usite *model.ServerUSite) (affected int64, err error) {
	if usite.ID > 0 {
		// 更新
		db := repo.db.Model(&model.ServerUSite{}).Where("id = ?", usite.ID).Updates(usite)
		if err = db.Error; err != nil {
			repo.log.Error(err)
			return db.RowsAffected, err
		}
		return db.RowsAffected, nil
	}
	// 新增
	db := repo.db.Create(usite)
	if err = db.Error; err != nil {
		repo.log.Error(err)
		return db.RowsAffected, err
	}
	return db.RowsAffected, nil
}

// GetServerUSiteID 返回满足过滤条件的机位(U位)ID  不支持模糊查找
func (repo *MySQLRepo) GetServerUSiteID(cond *model.ServerUSite) (id []uint, err error) {
	db := repo.db.Model(&model.ServerUSite{})
	if cond != nil {
		if cond.IDCID > 0 {
			db = db.Where("idc_id = ?", cond.IDCID)
		}
		if cond.ServerRoomID > 0 {
			db = db.Where("server_room_id = ?", cond.ServerRoomID)
		}
		if cond.ServerCabinetID > 0 {
			db = db.Where("server_cabinet_id = ?", cond.ServerCabinetID)
		}
		if cond.Number != "" {
			db = db.Where("number = ?", cond.Number)
		}
	}

	if err = db.Pluck("id", &id).Error; err != nil {
		repo.log.Error(err)
		return id, err
	}
	return id, nil
}

// GetServerUSiteByID 返回指定ID的机位(U位)
func (repo *MySQLRepo) GetServerUSiteByID(id uint) (*model.ServerUSite, error) {
	var usite model.ServerUSite
	if err := repo.db.Where("id = ?", id).Find(&usite).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &usite, nil
}

// GetServerUSiteByNumber 返回指定机架及编号的机位(U位)
func (repo *MySQLRepo) GetServerUSiteByNumber(cabinetID uint, number string) (*model.ServerUSite, error) {
	var usite model.ServerUSite
	if err := repo.db.Where("server_cabinet_id = ? AND number = ?", cabinetID, number).Find(&usite).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &usite, nil
}

// BatchUpdateServerUSitesStatus 批量更新机位状态信息
func (repo *MySQLRepo) BatchUpdateServerUSitesStatus(ids []uint, status string) (affected int64, err error) {
	db := repo.db.Model(&model.ServerUSite{}).Where("id in (?)", ids).Update("status", status)
	if err = db.Error; err != nil {
		repo.log.Error(err)
		return 0, err
	}
	return db.RowsAffected, nil
}

// DeleteServerUSitePort 删除机位端口号
func (repo *MySQLRepo) DeleteServerUSitePort(id uint) (affected int64, err error) {
	var nullArray []*model.SwitchInfo

	db := repo.db.Model(&model.ServerUSite{}).Where("id = ?", id).Updates(map[string]interface{}{
		"oobnet_switches":   utils.ToJsonString(nullArray),
		"intranet_switches": utils.ToJsonString(nullArray),
		"extranet_switches": utils.ToJsonString(nullArray),
	})
	if err = db.Error; err != nil {
		repo.log.Error(err)
		return 0, err
	}
	return db.RowsAffected, nil
}

// RemoveServerUSiteByID 删除机位
func (repo *MySQLRepo) RemoveServerUSiteByID(id uint) (affected int64, err error) {
	db := repo.db.Unscoped().Delete(&model.ServerUSite{}, "id = ? ", id)
	if err = db.Error; err != nil {
		repo.log.Error(err)
		return 0, err
	}
	return db.RowsAffected, nil
}

// CountServerUSite 统计机位数量
func (repo *MySQLRepo) CountServerUSite(cond *model.CombinedServerUSite) (count int64, err error) {
	db := repo.db.Table("server_usite").
		Joins("left join server_cabinet cabinet on cabinet.id = server_usite.server_cabinet_id")

	db = addUsiteCond(db, cond)

	if err = db.Select("server_usite.id").Count(&count).Error; err != nil {
		repo.log.Error(err)
		return count, err
	}

	return count, nil
}

// GetServerUSiteByCond 查询机位信息
func (repo *MySQLRepo) GetServerUSiteByCond(cond *model.CombinedServerUSite, orderby model.OrderBy, limiter *page.Limiter) (items []*model.ServerUSite, err error) {
	db := repo.db.Table("server_usite").
		Joins("left join server_cabinet cabinet on cabinet.id = server_usite.server_cabinet_id")

	db = addUsiteCond(db, cond)

	for i := range orderby {
		db = db.Order(orderby[i].String())
	}
	if limiter != nil {
		db = db.Limit(limiter.Limit).Offset(limiter.Offset)
	}

	if err = db.Select("server_usite.*").Find(&items).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}

	return items, nil
}

// addUsiteCond 添加机位查询条件
func addUsiteCond(db *gorm.DB, cond *model.CombinedServerUSite) *gorm.DB {
	if cond.ServerRoomName != "" {
		db = db.Joins("left join server_room on server_room.id = server_usite.server_room_id")
	}
	if cond.ServerRoomNameCabinetNumUSiteNumSlice != nil || cond.ServerRoomNameCabinetNumUSiteNum != "" {
		db = db.Joins("LEFT JOIN server_room room ON room.id = server_usite.server_room_id")
	}
	if cond != nil {
		db = MultiNumQuery(db, "server_usite.idc_id", cond.IDCID)
		db = MultiNumQuery(db, "server_usite.server_room_id", cond.ServerRoomID)
		db = MultiQuery(db, "server_room.name", cond.ServerRoomName)
		db = ConcatColumnStringQuery(db, "room.name, cabinet.number, server_usite.number", cond.ServerRoomNameCabinetNumUSiteNum)
		db = ConcatColumnSliceStringQuery(db, "room.name, cabinet.number, server_usite.number", cond.ServerRoomNameCabinetNumUSiteNumSlice)
		db = MultiNumQuery(db, "server_usite.server_cabinet_id", cond.ServerCabinetID)
		db = MultiNumQuery(db, "server_usite.network_area_id", cond.NetAreaID)
		db = MultiQuery(db, "cabinet.number", cond.CabinetNumber)
		db = MultiEnumQuery(db, "server_usite.status", cond.Status)
		db = MultiEnumQuery(db, "server_usite.la_wa_port_rate", cond.LAWAPortRate)
		db = MultiMatchWithSpaceQuery(db, "server_usite.physical_area", cond.PhysicalArea)
		db = MultiEnumQuery(db, "server_usite.height", cond.Height)
		db = MultiMatchQuery(db, "server_usite.number", cond.USiteNumber)
	}
	db = db.Where("server_usite.deleted_at is null ")

	return db
}

// GetServerUsiteByNetworkDeviceName 根据指定的网络设备返回机位ID
func (repo *MySQLRepo) GetServerUsiteByNetworkDeviceName(name []string) (nd []uint, err error) {
	nd = make([]uint, 0)
	sql := `select DISTINCT(id) from (select TRIM(BOTH '"' FROM name->'$[0]') as name, id from (select intranet_switches->'$[*].name' as name, id  from server_usite) as ali) as alii where `
	for i := range name {
		name[i] = fmt.Sprintf("'%s'", name[i])
	}
	rawsql := fmt.Sprintf("%s name IN(%s)", sql, strings.Join(name, ","))
	repo.log.Debugf(rawsql)
	if err = repo.db.Raw(rawsql).Pluck("id", &nd).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return nd, err
}

func (repo *MySQLRepo) GetPhysicalAreas() (*model.DeviceQueryParamResp, error) {
	p := model.DeviceQueryParamResp{}
	p.ParamName = "physical_area"
	rawSQL := ""
	mods := make([]model.ParamList, 0)
	rawSQL = "select DISTINCT 0 AS id, physical_area AS name from server_usite"
	_ = repo.db.Raw(rawSQL).Scan(&mods).Error
	if len(mods) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	p.List = mods
	return &p, nil
}

// BatchUpdateServerUSitesRemark 批量更新机位备注信息
func (repo *MySQLRepo) BatchUpdateServerUSitesRemark(ids []uint, remark string) (affected int64, err error) {
	db := repo.db.Model(&model.ServerUSite{}).Where("id in (?)", ids).Update("remark", remark)
	if err = db.Error; err != nil {
		repo.log.Error(err)
		return 0, err
	}
	return db.RowsAffected, nil
}