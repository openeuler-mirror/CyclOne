package mysqlrepo

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/voidint/page"

	"idcos.io/cloudboot/model"
)

// GetTORs 返回所有的TOR名称列表
func (repo *MySQLRepo) GetTORs() (items []string, err error) {
	if err = repo.db.Raw(`SELECT DISTINCT(tor) FROM network_device WHERE tor IS NOT NULL AND tor != ""`).Pluck("tor", &items).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return nil, nil
}

// GetNetworkDeviceByTORS 根据指定的tor返回网络设备
func (repo *MySQLRepo) GetNetworkDeviceByTORS(tor ...string) (nd []*model.NetworkDevice, err error) {
	db := repo.db.Where("tor in (?)", tor).Find(&nd)
	if err = db.Error; err != nil {
		return nil, err
	}
	return nd, err
}

// GetTORBySN 返回目标设备所属的TOR名称
func (repo *MySQLRepo) GetTORBySN(sn string) (tor string, err error) {
	switcher, err := repo.GetIntranetSwitchBySN(sn)
	if err != nil {
		return "", err
	}
	return switcher.TOR, nil
}

// GetIntranetSwitchBySN 查询设备所在机位的内网交换机
func (repo *MySQLRepo) GetIntranetSwitchBySN(sn string) (*model.NetworkDevice, error) {
	return repo.getSwitchBySN(sn, intranetField)
}

// GetExtranetSwitchBySN 查询设备所在机位的外网交换机
func (repo *MySQLRepo) GetExtranetSwitchBySN(sn string) (*model.NetworkDevice, error) {
	return repo.getSwitchBySN(sn, extranetField)
}

// SaveNetworkDevice 保存网络设备
func (repo *MySQLRepo) SaveNetworkDevice(na *model.NetworkDevice) (networkDevice *model.NetworkDevice, err error) {
	if na.ID > 0 {
		// 更新
		db := repo.db.Model(&model.NetworkDevice{}).Where("id = ?", na.ID).Updates(na)
		if err = db.Error; err != nil {
			repo.log.Error(err)
			return nil, err
		}
		return na, nil
	}
	// 新增
	db := repo.db.Create(na)
	if err = db.Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return na, nil
}

// GetPeerNetworkDeviceByCabinetID 返回指定ID的机架(柜)
func (repo *MySQLRepo) GetPeerNetworkDeviceByCabinetID(id uint) (*model.NetworkDevice, error) {
	var netDev = make([]*model.NetworkDevice, 0)
	rawSQL := fmt.Sprintf(`select * from network_device where tor = (select tor from network_device where server_cabinet_id = %d limit 1) and server_cabinet_id != %d`, id, id)
	if err := repo.db.Raw(rawSQL).Scan(&netDev).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	if len(netDev) > 0 {
		return netDev[0], nil
	}
	return nil, nil
}

const (
	oobnetField   = "oobnet_switches"
	intranetField = "intranet_switches"
	extranetField = "extranet_switches"
)

// getSwitchBySN 查询设备所在机位的管理网/内网/外网交换机。
// netField可选值：oobnet_switch_name、intranet_switch_name、extranet_switch_name
func (repo *MySQLRepo) getSwitchBySN(sn string, netField string) (*model.NetworkDevice, error) {
	sql := `
select * from network_device where type='switch' and name = (
	select JSON_UNQUOTE(%s->'$[0].name') from server_usite where id = (
		select server_usite_id from device where sn = ?
	)
);`

	var netDev model.NetworkDevice
	if err := repo.db.Raw(fmt.Sprintf(sql, netField), sn).Scan(&netDev).Error; err != nil {
		repo.log.Errorf("device(sn:%s) has not found %s net-device, err: %s", sn, netField, err)
		return nil, fmt.Errorf("机位没有关联内/外网交换机")
	}
	return &netDev, nil
}

// GetNetworkDevicesByCond 返回满足过滤条件的网络设备(不支持模糊查找)
func (repo *MySQLRepo) GetNetworkDevicesByCond(cond *model.NetworkDeviceCond, orderby model.OrderBy, limiter *page.Limiter) (items []*model.NetworkDevice, err error) {
	db := repo.db.Model(&model.NetworkDevice{})
	db = addCond(db, cond)

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

// CountNetworkDevices 统计满足过滤条件的网络设备数量
func (repo *MySQLRepo) CountNetworkDevices(cond *model.NetworkDeviceCond) (count int64, err error) {
	db := repo.db.Model(&model.NetworkDevice{})
	db = addCond(db, cond)

	if err = db.Select("id").Count(&count).Error; err != nil {
		repo.log.Error(err)
		return 0, err
	}
	return count, nil
}

// RemoveNetworkDeviceByID 删除指定ID的网络设备
func (repo *MySQLRepo) RemoveNetworkDeviceByID(id uint) (err error) {
	db := repo.db.Unscoped().Where("id =?", id).Delete(model.NetworkDevice{})

	if err = db.Error; err != nil {
		return err
	}
	return
}

// GetNetworkDeviceByID 查询指定ID的网络设备
func (repo *MySQLRepo) GetNetworkDeviceByID(id uint) (network *model.NetworkDevice, err error) {
	var net model.NetworkDevice
	db := repo.db.Where("id =?", id).Find(&net)

	if err = db.Error; err != nil {
		return nil, err
	}
	return &net, nil
}

// GetNetworkDeviceBySN 查询指定sn的网络设备
func (repo *MySQLRepo) GetNetworkDeviceBySN(sn string) (network []*model.NetworkDevice, err error) {
	var net []*model.NetworkDevice
	db := repo.db.Where("sn =?", sn).Find(&net)

	if err = db.Error; err != nil {
		return nil, err
	}
	return net, nil
}

// GetNetworkDeviceByFixedAssetNumber 查询指定FixedAssetNumber的网络设备
func (repo *MySQLRepo) GetNetworkDeviceByFixedAssetNumber(FixedAssetNumber string) (network []*model.NetworkDevice, err error) {
	var net []*model.NetworkDevice
	db := repo.db.Where("fixed_asset_number =?", FixedAssetNumber).Find(&net)

	if err = db.Error; err != nil {
		return nil, err
	}
	return net, nil
}

// addCond 添加查询条件
func addCond(db *gorm.DB, cond *model.NetworkDeviceCond) *gorm.DB {
	if cond != nil {
		db = MultiNumQuery(db, "network_device.idc_id", cond.IDCID)
		db = MultiNumQuery(db, "network_device.server_room_id", cond.ServerRoomID)
		db = MultiNumQuery(db, "network_device.server_cabinet_id", cond.ServerCabinetID)
		if cond.ServerRoomName != "" {
			db = db.Joins("left join server_room on server_room.id = network_device.server_room_id")
			db = MultiQuery(db, "server_room.name", cond.ServerRoomName)
		}

		if cond.ServerCabinetNumber != "" {
			db = db.Joins("left join server_cabinet on server_cabinet.id = network_device.server_cabinet_id")
			db = MultiQuery(db, "server_cabinet.number", cond.ServerCabinetNumber)
		}
		db = MultiQuery(db, "network_device.name", cond.Name)
		db = MultiQuery(db, "network_device.fixed_asset_number", cond.FixedAssetNumber)
		db = MultiQuery(db, "network_device.sn", cond.SN)
		db = MultiQuery(db, "network_device.tor", cond.TOR)
		db = MultiQuery(db, "network_device.type", cond.Type)
		db = MultiQuery(db, "network_device.model", cond.ModelNumber)
		db = MultiQuery(db, "network_device.vendor", cond.Vendor)
		db = MultiQuery(db, "network_device.os", cond.OS)
		db = MultiQuery(db, "network_device.usage", cond.Usage)
		db = MultiQuery(db, "network_device.status", cond.Status)
	}
	return db
}
