package mysqlrepo

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
	"idcos.io/cloudboot/model"
	mystrings "idcos.io/cloudboot/utils/strings"
)

// SaveCollectedDeviceBySN 保存(更新/新增)采集到的设备信息
func (repo *MySQLRepo) SaveCollectedDeviceBySN(dev *model.CollectedDevice) error {
	db := repo.db.Model(&model.CollectedDevice{}).Where("sn = ?", dev.SN).Updates(dev)
	if err := db.Error; err != nil {
		repo.log.Error(err)
		return err
	}
	if db.RowsAffected > 0 {
		return nil
	}

	//if err := repo.db.Create(dev).Error; err != nil {
	//	repo.log.Error(err)
	//	return err
	//}
	return nil
}

//GetDeviceBySN 根据SN查找设备
func (repo *MySQLRepo) GetDeviceBySN(SN string) (*model.Device, error) {
	var dev model.Device
	if err := repo.db.Unscoped().Where("sn = ?", SN).Find(&dev).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &dev, nil
}

// GetDeviceBySNOrMAC 根据SN或者mac地址查询设备
func (repo *MySQLRepo) GetDeviceBySNOrMAC(snOrMAC string) (*model.Device, error) {
	var dev model.Device
	if err := repo.db.
		Where(
			"sn = ? OR JSON_CONTAINS(nic->'$.items',JSON_OBJECT('mac',?)) OR JSON_CONTAINS(nic->'$.items',JSON_OBJECT('mac',?))",
			snOrMAC,
			strings.ToLower(snOrMAC),
			strings.ToUpper(snOrMAC),
		).
		Find(&dev).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &dev, nil
}

//GetDeviceByFixAssetNumber 根据固资查找设备
func (repo *MySQLRepo) GetDeviceByFixAssetNumber(fixAssetNum string) (*model.Device, error) {
	var dev model.Device
	if err := repo.db.Unscoped().Where("fixed_asset_number = ?", fixAssetNum).Find(&dev).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &dev, nil
}

// GetDeviceByID 返回指定ID的物理机
func (repo *MySQLRepo) GetDeviceByID(id uint) (*model.Device, error) {
	var dev model.Device
	if err := repo.db.Unscoped().Where("id = ?", id).Find(&dev).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &dev, nil
}

//SaveDevice 新增/修改设备
func (repo *MySQLRepo) SaveDevice(mod *model.Device) (affected int64, err error) {
	db := repo.db.Save(mod)
	if err = db.Error; err != nil {
		repo.log.Error(err)
		return db.RowsAffected, err
	}
	return db.RowsAffected, nil
}

//UpdateDevice 修改物理机
func (repo *MySQLRepo) UpdateDevice(mod *model.Device) (affected int64, err error) {
	db := repo.db.Model(model.Device{}).Where("id = ?", mod.ID).Update(mod)
	if err = db.Error; err != nil {
		repo.log.Error(err)
		return db.RowsAffected, err
	}
	return db.RowsAffected, nil
}

// UpdateDeviceBySN 更新目标设备的
func (repo *MySQLRepo) UpdateDeviceBySN(dev *model.Device) (affected int64, err error) {
	db := repo.db.Model(model.Device{}).Where("sn = ?", dev.SN).Updates(dev)
	//下面这几个字段可以置0
	//Update("server_room_id", dev.ServerRoomID).
	//Update("server_cabinet_id", dev.CabinetID).
	//Update("server_usite_id", dev.USiteID).
	//Update("store_room_id", dev.StoreRoomID).
	//Update("virtual_cabinet_id", dev.VCabinetID)
	if err = db.Error; err != nil {
		repo.log.Error(err)
		return db.RowsAffected, err
	}
	return db.RowsAffected, nil
}

//RemoveDeviceByID 删除物理机
func (repo *MySQLRepo) RemoveDeviceByID(id uint) (affected int64, err error) {
	mod := model.Device{Model: gorm.Model{ID: id}}
	db := repo.db.Unscoped().Delete(mod)
	if err = db.Error; err != nil {
		repo.log.Error(err)
		return db.RowsAffected, err
	}
	return db.RowsAffected, nil
}

//RemoveDeviceBySN 删除物理机
func (repo *MySQLRepo) RemoveDeviceBySN(sn string) (affected int64, err error) {
	mod := model.Device{}
	db := repo.db.Unscoped().Where("sn = ?", sn).Delete(&mod)
	if err = db.Error; err != nil {
		repo.log.Errorf("RemoveDeviceBySN failure:%s", err.Error())
	}
	return db.RowsAffected, nil
}

//GetDevicesByUSiteID 根据机位信息查询设备
func (repo *MySQLRepo) GetDevicesByUSiteID(id uint) (items []*model.Device, err error) {
	db := repo.db.Model(&model.Device{}).Where("server_usite_id =?", id)

	if err = db.Find(&items).Error; err != nil {
		repo.log.Error(err)
		return
	}
	return
}

// GetDevicesByUSiteIDS 根据机位信息查询设备
func (repo *MySQLRepo) GetDevicesByUSiteIDS(ids []uint, usage string) (items []*model.Device, err error) {
	db := repo.db.Model(&model.Device{}).Where("server_usite_id in (?)", ids).Where("`usage` = ?", usage)

	if err = db.Find(&items).Error; err != nil {
		repo.log.Error(err)
		return
	}
	return
}

//CountDevices 统计满足过滤条件的物理机数量
func (repo *MySQLRepo) CountDevices(cond *model.Device) (count int64, err error) {
	db := repo.db.Model(&model.Device{})
	if cond != nil {
		if cond.FixedAssetNumber != "" {
			db = db.Where("`fixed_asset_number` LIKE ?", fmt.Sprintf("%%%s%%", cond.FixedAssetNumber))
		}
		if cond.SN != "" {
			//这里的SN可能是换行符分隔的多个SN构成的数组
			sns := mystrings.MultiLines2Slice(cond.SN)
			for i, sn := range sns {
				if i == 0 {
					db = db.Where("`sn` LIKE ?", fmt.Sprintf("%%%s%%", sn))
				} else {
					db = db.Or("`sn` LIKE ?", fmt.Sprintf("%%%s%%", sn))
				}
			}
		}
		if cond.Vendor != "" {
			db = db.Where("`vendor` LIKE ?", fmt.Sprintf("%%%s%%", cond.Vendor))
		}
		if cond.DevModel != "" {
			db = db.Where("`model` LIKE ?", fmt.Sprintf("%%%s%%", cond.DevModel))
		}
		if cond.Usage != "" {
			db = db.Where("`usage` LIKE ?", fmt.Sprintf("%%%s%%", cond.Usage))
		}
		if cond.Category != "" {
			db = db.Where("`category` LIKE ?", fmt.Sprintf("%%%s%%", cond.Category))
		}
		if cond.OperationStatus != "" {
			db = db.Where("`operation_status` LIKE ?", fmt.Sprintf("%%%s%%", cond.OperationStatus))
		}
	}
	err = db.Count(&count).Error
	if err != nil {
		repo.log.Errorf("统计物理机错误，%s", err.Error())
	}
	return count, err
}

// GetDevices 返回满足过滤条件的物理机
func (repo *MySQLRepo) GetDevices(cond *model.Device, orderby model.OrderBy, limiter *page.Limiter) (items []*model.Device, err error) {
	items = make([]*model.Device, 0)
	db := repo.db.Model(&model.Device{})
	if cond != nil {
		if cond.FixedAssetNumber != "" {
			//这里的固资编号可能是换行符分隔的多个构成的数组
			fns := mystrings.MultiLines2Slice(cond.FixedAssetNumber)
			db = db.Where("`fixed_asset_number` IN (?)", fns)
			//for i, fn := range fns {
			//	if i == 0 {
			//		db = db.Where("`fixed_asset_number` LIKE ?", fmt.Sprintf("%%%s%%", fn))
			//	} else {
			//		db = db.Or("`fixed_asset_number` LIKE ?", fmt.Sprintf("%%%s%%", fn))
			//	}
			//}
		}
		if cond.SN != "" {
			//这里的SN可能是换行符分隔的多个SN构成的数组
			sns := mystrings.MultiLines2Slice(cond.SN)
			db = db.Where("`sn` IN (?)", sns)
			//for i, sn := range sns {
			//	if i == 0 {
			//		db = db.Where("`sn` LIKE ?", fmt.Sprintf("%%%s%%", sn))
			//	} else {
			//		db = db.Or("`sn` LIKE ?", fmt.Sprintf("%%%s%%", sn))
			//	}
			//}
		}
		if cond.Vendor != "" {
			db = db.Where("`vendor` LIKE ?", fmt.Sprintf("%%%s%%", cond.Vendor))
		}
		if cond.DevModel != "" {
			db = db.Where("`model` LIKE ?", fmt.Sprintf("%%%s%%", cond.DevModel))
		}
		if cond.Usage != "" {
			db = db.Where("`usage` LIKE ?", fmt.Sprintf("%%%s%%", cond.Usage))
		}
		if cond.Category != "" {
			db = db.Where("`category` LIKE ?", fmt.Sprintf("%%%s%%", cond.Category))
		}
		if cond.OperationStatus != "" {
			db = db.Where("`operation_status` LIKE ?", fmt.Sprintf("%%%s%%", cond.OperationStatus))
		}
		if cond.OOBAccessible != nil {
			if *cond.OOBAccessible == model.YES {
				db = db.Where("`oob_accessible` = ?", *cond.OOBAccessible)
			} else {
				db = db.Where("`oob_accessible` = ? OR `oob_accessible` is null", model.NO)
			}
		}
		if cond.VCabinetID != 0 {
			db = db.Where("`virtual_cabinet_id` = ?", cond.VCabinetID)
		}

	}
	for i := range orderby {
		db = db.Order(orderby[i].String())
	}
	if limiter != nil {
		db = db.Offset(limiter.Offset).Limit(limiter.Limit)
	}
	err = db.Find(&items).Error
	if err != nil {
		repo.log.Errorf("多条件查询物理机错误，%s", err.Error())
		return nil, err
	}
	return items, nil
}

const condSep = ","

func (repo *MySQLRepo) setWhereSQL4CombinedDevice(db *gorm.DB, cond *model.CombinedDeviceCond) *gorm.DB {
	if db == nil || cond == nil {
		return db
	}
	var sb strings.Builder
	if cond.SN != "" {
		sb.Reset()
		for i, sn := range mystrings.MultiLines2Slice(cond.SN) {
			if i > 0 {
				sb.WriteString(" OR ")
			}
			sb.WriteString(fmt.Sprintf("device.sn LIKE '%%%s%%'", sn))
		}
		db = db.Where(sb.String())
	}
	if cond.FixedAssetNumber != "" {
		sb.Reset()
		for i, number := range mystrings.MultiLines2Slice(cond.FixedAssetNumber) {
			if i > 0 {
				sb.WriteString(" OR ")
			}
			sb.WriteString(fmt.Sprintf("device.fixed_asset_number LIKE '%%%s%%'", number))
		}
		db = db.Where(sb.String())
	}
	if cond.ID != nil {
		db = db.Where("device.id IN (?)", cond.ID)
	}
	if cond.IDCID != nil {
		db = db.Where("device.idc_id IN (?)", cond.IDCID)
	}
	if cond.ServerRoomID != nil {
		db = db.Where("device.server_room_id IN (?)", cond.ServerRoomID)
	}
	if cond.ServerCabinet != "" {
		sb.Reset()
		for i, number := range mystrings.MultiLines2Slice(cond.ServerCabinet) {
			if i > 0 {
				sb.WriteString(" OR ")
			}
			sb.WriteString(fmt.Sprintf("server_cabinet.number LIKE '%%%s%%'", number))
		}
		db = db.Where(sb.String())
	}
	if cond.ServerRoomName != "" {
		sb.Reset()
		for i, roomName := range mystrings.MultiLines2Slice(cond.ServerRoomName) {
			if i > 0 {
				sb.WriteString(" OR ")
			}
			sb.WriteString(fmt.Sprintf("server_room.name LIKE '%%%s%%'", roomName))
		}
		db = db.Where(sb.String())
	}
	if cond.ServerUsiteNumber != "" {
		sb.Reset()
		for i, usiteNum := range mystrings.MultiLines2Slice(cond.ServerUsiteNumber) {
			if i > 0 {
				sb.WriteString(" OR ")
			}
			sb.WriteString(fmt.Sprintf("server_usite.number LIKE '%%%s%%'", usiteNum))
		}
		db = db.Where(sb.String())
	}
	if cond.USiteID != nil {
		db = db.Where("device.server_cabinet_id IN (?)", cond.USiteID)
	}
	if cond.PhysicalArea != "" {
		// PhysicalArea: 'Management bonding区1' contains space 
		db = db.Where("server_usite.physical_area IN (?)", mystrings.MultiLines2SliceWithSpace(cond.PhysicalArea))
		//sb.Reset()
		//for i, phyArea := range mystrings.MultiLines2Slice(cond.PhysicalArea) {
		//	if i > 0 {
		//		sb.WriteString(" OR ")
		//	}
		//	sb.WriteString(fmt.Sprintf("device.operation_status <> 'in_store' AND server_usite.physical_area LIKE '%%%s%%'", phyArea))
		//}
		//db = db.Where(sb.String())
	}
	if cond.Vendor != "" {
		db = db.Where("device.vendor IN (?)", mystrings.MultiLines2Slice(cond.Vendor))
	}
	if cond.DevModel != "" {
		sb.Reset()
		for i, devModel := range mystrings.MultiLines2Slice(cond.DevModel) {
			if i > 0 {
				sb.WriteString(" OR ")
			}
			sb.WriteString(fmt.Sprintf("device.model LIKE '%%%s%%'", devModel))
		}
		db = db.Where(sb.String())
	}
	if cond.HardwareRemark != "" {
		db = db.Where(fmt.Sprintf("device.hardware_remark LIKE '%%%s%%'", cond.HardwareRemark))
	}
	if cond.Usage != "" {
		db = db.Where("device.`usage` IN (?)", mystrings.MultiLines2Slice(cond.Usage))
	}
	if len(cond.USiteID) != 0 {
		db = db.Where("device.server_usite_id IN (?)", cond.USiteID)
	}
	if cond.Category != "" {
		db = db.Where("device.`category` IN (?)", mystrings.MultiLines2Slice(cond.Category))
	}
	if cond.OperationStatus != "" {
		db = db.Where("device.`operation_status` IN (?)", strings.Split(cond.OperationStatus, condSep))
	}
	if cond.PreDeployed { // '预部署'状态，返回不在装机列表或者装机状态是失败的设备。
		db = db.Where("device_setting.sn IS NULL OR device_setting.`status` = ?", model.InstallStatusFail)
	} else if cond.DeployStatus != "" {
		db = db.Where("device_setting.`status` IN (?)", strings.Split(cond.DeployStatus, condSep))
	}
	if cond.IntranetIP != "" {
		//db = db.Where("device_setting.intranet_ip IN (?)", cond.IntranetIP)
		db = MultiMatchQuery(db, "device_setting.intranet_ip", cond.IntranetIP)
	}
	if cond.ExtranetIP != "" {
		//db = db.Where(" IN (?)", cond.ExtranetIP)
		db = MultiMatchQuery(db, "device_setting.extranet_ip", cond.ExtranetIP)
	}
	if cond.IP != "" {
		sb.Reset()
		for i, ip := range mystrings.MultiLines2Slice(cond.IP) {
			if i > 0 {
				sb.WriteString(" OR ")
			}
			sb.WriteString(fmt.Sprintf("device_setting.intranet_ip LIKE '%%%s%%' OR device_setting.extranet_ip LIKE '%%%s%%' OR device.oob_ip LIKE '%%%s%%'", ip, ip, ip))
		}
		db = db.Where(sb.String())
	}
	if cond.OOBAccessible != "" {
		if strings.Contains(cond.OOBAccessible, model.Unknown) {
			db = db.Where("device.`oob_accessible` IN (?) OR device.`oob_accessible` IS NULL", mystrings.MultiLines2Slice(cond.OOBAccessible))
		} else {
			db = db.Where("device.`oob_accessible` IN (?)", mystrings.MultiLines2Slice(cond.OOBAccessible))
		}
	}
	return db
}

// CountCombinedDevices 统计满足过滤条件的记录数量
func (repo *MySQLRepo) CountCombinedDevices(cond *model.CombinedDeviceCond) (count int64, err error) {
	db := repo.db.Table("device")
	db = repo.setWhereSQL4CombinedDevice(db, cond)
	if err = db.Select("COUNT(device.id)").
		Joins("LEFT JOIN device_setting device_setting ON device.sn = device_setting.sn").
		Joins("LEFT JOIN hardware_template hardware_template ON device_setting.hardware_template_id = hardware_template.id").
		Joins("LEFT JOIN system_template system_template ON device_setting.system_template_id = system_template.id").
		Joins("LEFT JOIN image_template image_template ON device_setting.image_template_id = image_template.id").
		Joins("LEFT JOIN server_cabinet server_cabinet ON device.server_cabinet_id = server_cabinet.id").
		Joins("LEFT JOIN server_room server_room ON device.server_room_id = server_room.id").
		Joins("LEFT JOIN server_usite server_usite ON device.server_usite_id = server_usite.id").Count(&count).Error; err != nil {
		repo.log.Error(err)
		return 0, err
	}

	return count, nil
}

// GetCombinedDevices 返回满足过滤条件的设备及其装机参数列表
func (repo *MySQLRepo) GetCombinedDevices(cond *model.CombinedDeviceCond, orderby model.OrderBy, limiter *page.Limiter) (items []*model.CombinedDevice, err error) {
	db := repo.db.Table("device")
	db = repo.setWhereSQL4CombinedDevice(db, cond)
	db = db.Select(`
	device.id,
	device.fixed_asset_number,
	device.sn,
	device.vendor,
	device.model,
	device.power_status,
	device.oob_accessible,
	device.arch,
	device.usage,
	device.category,
	device.idc_id,
	device.server_room_id,
	device.server_cabinet_id,
	device.server_usite_id,
	device.store_room_id,
	device.virtual_cabinet_id,
	device.hardware_remark,
	device.raid_remark,
	device.started_at,
	device.onshelve_at,
	device.operation_status,
	device.oob_ip,
	device.oob_user,
	device.oob_password,
	device.cpu_sum,
	device.cpu,
	device.memory_sum,
	device.memory,
	device.disk_sum,
	device.disk,
	device.disk_slot,
	device.nic,
	device.nic_device,
	device.bootos_ip,
	device.bootos_mac,
	device.motherboard,
	device.raid,
	device.oob,
	device.bios,
	device.fan,
	device.power,
	device.hba,
	device.pci,
	device.switch,
	device.lldp,
	device.extra,
	device.origin_node,
	device.origin_node_ip,
	device.operation_user_id,
	device.creator,
	device.updater,
	device.remark,
	device.order_number,
	device_setting.status AS deploy_status,
	device_setting.install_progress,
	IF(system_template.name != '', system_template.name, image_template.name) AS os,
	device_setting.system_template_id,
	system_template.name AS system_template_name,
	device_setting.hardware_template_id,
	concat(hardware_template.vendor,'-',hardware_template.model) AS hardware_name,
	device_setting.intranet_ip,
	device_setting.extranet_ip,
	device_setting.intranet_ipv6,
	device_setting.extranet_ipv6,	
	device_setting.image_template_id AS image_template_id,
	image_template.name AS image_template_name,
	device_lifecycle.asset_belongs,
	device_lifecycle.owner,
	device_lifecycle.is_rental,
	device_lifecycle.maintenance_service_provider,
	device_lifecycle.maintenance_service,
	device_lifecycle.logistics_service,
	device_lifecycle.maintenance_service_date_begin,
	device_lifecycle.maintenance_service_date_end,
	device_lifecycle.maintenance_service_status,
	device_lifecycle.device_retired_date,
	device_lifecycle.lifecycle_log,
	device.created_at,
	IF(device_setting.updated_at IS NULL, device.updated_at, IF(device.updated_at > device_setting.updated_at, device.updated_at, device_setting.updated_at)) AS updated_at
	`).
		Joins("LEFT JOIN device_setting device_setting ON device.sn = device_setting.sn").
		Joins("LEFT JOIN hardware_template hardware_template ON device_setting.hardware_template_id = hardware_template.id").
		Joins("LEFT JOIN system_template system_template ON device_setting.system_template_id = system_template.id").
		Joins("LEFT JOIN image_template image_template ON device_setting.image_template_id = image_template.id").
		Joins("LEFT JOIN server_cabinet server_cabinet ON device.server_cabinet_id = server_cabinet.id").
		Joins("LEFT JOIN server_room server_room ON device.server_room_id = server_room.id").
		Joins("LEFT JOIN server_usite server_usite ON device.server_usite_id = server_usite.id").
		Joins("LEFT JOIN device_lifecycle device_lifecycle ON device.sn = device_lifecycle.sn")

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

// 查询物理机查询条件的下拉源数据
func (repo *MySQLRepo) GetDeviceQuerys(param string) (*model.DeviceQueryParamResp, error) {
	p := model.DeviceQueryParamResp{}
	p.ParamName = param
	rawSQL := ""
	mods := make([]struct {
		ID   string
		Name string
	}, 0)
	switch param {
	case "idc":
		rawSQL = "select DISTINCT device.idc_id AS id,idc.name from device left join idc ON device.idc_id = idc.id;"
	case "server_room":
		rawSQL = "select DISTINCT device.server_room_id AS id,server_room.name from device left join server_room ON device.server_room_id = server_room.id;"
	case "server_cabinet":
		rawSQL = "select DISTINCT device.server_cabinet_id AS id,server_cabinet.number AS name from device left join server_cabinet ON device.server_cabinet_id = server_cabinet.id;"
	case "physical_area":
		rawSQL = "select DISTINCT 0 AS id, server_usite.physical_area AS name from server_usite where server_usite.physical_area is not null;"
	case "op_status":
		rawSQL = "select DISTINCT 0 AS id, operation_status AS name from device;"
	case "usage":
		rawSQL = "select DISTINCT 0 AS id, `usage` AS name from device;"
	case "category":
		rawSQL = "select DISTINCT 0 AS id, `category` AS name from device;"
	case "category_pre_deploy":
		rawSQL = "select DISTINCT 0 AS id, `category` AS name from device where operation_status='pre_deploy';"		
	case "vendor":
		rawSQL = "select DISTINCT 0 AS id, `vendor` AS name from device;"
	case "switches":
		rawSQL = "select DISTINCT n.fixed_asset_number AS id ,concat(n.name,'(',n.fixed_asset_number,')') AS name from network_device AS n"
	default:
		return nil, errors.New("not supported yet")
	}

	_ = repo.db.Raw(rawSQL).Scan(&mods).Error
	if len(mods) == 0 {
		p.List = []model.ParamList{}
	}
	for _, mod := range mods {
		p.List = append(p.List, model.ParamList{
			ID:   mod.ID,
			Name: mod.Name,
		})
	}
	return &p, nil
}

// 按月份查询最大编号的固资编号
func (repo *MySQLRepo) GetMaxFixedAssetNumber(month string) (fixedAssetNumber string, err error) {
	rawSQL := fmt.Sprintf("select max(fixed_asset_number) As fixed_asset_number from device where fixed_asset_number LIKE 'WDEV%s%%';", month)
	out := struct {
		FixedAssetNumber string
	}{}
	err = repo.db.Raw(rawSQL).Find(&out).Error
	if err == nil {
		return out.FixedAssetNumber, nil
	}
	return "", err
}

// GetDeviceByStartedAt 根据启用日期查询在该日期之前的设备
func (repo *MySQLRepo) GetDeviceByStartedAt(started_date string) (items []*model.Device, err error) {
	items = make([]*model.Device, 0)
	db := repo.db.Model(&model.Device{})
	if started_date != "" {
		db = db.Where("`started_at` < (?)", started_date)
	}
	err = db.Find(&items).Error
	if err != nil {
		repo.log.Errorf("根据启用日期查询在该日期之前的设备，%s", err.Error())
		return nil, err
	}
	return items, nil
}

//GetDeviceByOperationStatus 根据运营状态查找设备
func (repo *MySQLRepo) GetDeviceByOperationStatus(status string) (items []*model.Device, err error) {
	items = make([]*model.Device, 0)
	db := repo.db.Model(&model.Device{})
	if status != "" {
		db = db.Where("`operation_status` = ?", status)
	}
	err = db.Find(&items).Error
	if err != nil {
		repo.log.Errorf("根据运营状态(%s)查找设备失败，Error:%s", status, err.Error())
		return nil, err
	}
	return items, nil
}