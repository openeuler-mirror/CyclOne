package mysqlrepo

import (
	"encoding/json"
	"strings"

	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
	"idcos.io/cloudboot/model"
	nw "idcos.io/cloudboot/utils/network"
)

// GetIntranetIPNetworksBySN 查询指定物理机的内网IP所属网段集合
func (repo *MySQLRepo) GetIntranetIPNetworksBySN(sn string) (ipnet []*model.IPNetwork, err error) {
	return repo.getIPNetworksBySN(sn, true)
}

// GetExtranetIPNetworksBySN 查询指定物理机的外网IP所属网段集合
func (repo *MySQLRepo) GetExtranetIPNetworksBySN(sn string) (ipnet []*model.IPNetwork, err error) {
	return repo.getIPNetworksBySN(sn, false)
}

const (
	// sep 分隔符
	sep = ","
)

// 查询指定物理机的内/外网IP所属网段
func (repo *MySQLRepo) getIPNetworksBySN(sn string, intranet bool) (ipnet []*model.IPNetwork, err error) {
	var switcher *model.NetworkDevice
	if intranet {
		switcher, err = repo.GetIntranetSwitchBySN(sn)
	} else {
		switcher, err = repo.GetExtranetSwitchBySN(sn)
	}
	if err != nil || switcher == nil {
		return nil, fmt.Errorf("查询设备(sn:%s)所在机位关联的交换机失败", sn)
	}

	var items []*model.IPNetwork
	if err = repo.db.Model(&model.IPNetwork{}).Where("switches LIKE ?", "%"+switcher.FixedAssetNumber+"%").Find(&items).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	for i := range items {
		var switches []string
		if err = json.Unmarshal([]byte(items[i].Switches), &switches); err != nil {
			repo.log.Error(err)
			continue
		}
		for _, s := range switches {
			if strings.TrimSpace(s) == strings.TrimSpace(switcher.FixedAssetNumber) {
				ipnet = append(ipnet, items[i])
			}
		}
	}
	if len(ipnet) > 0 {
		return ipnet, nil
	}
	return nil, fmt.Errorf("没有找到覆盖交换机(固资号:%s)的网段", switcher.FixedAssetNumber)
}

// RemoveIPNetworkByID 删除指定ID的网段
func (repo *MySQLRepo) RemoveIPNetworkByID(id uint) (affected int64, err error) {
	tx := repo.db.Begin()

	// 删除网段
	tx = tx.Unscoped().Delete(&model.IPNetwork{}, "id = ?", id)
	if err = tx.Error; err != nil {
		tx.Rollback()
		repo.log.Error(err)
		return 0, err
	}

	if affected = tx.RowsAffected; affected <= 0 {
		return affected, tx.Commit().Error
	}

	// 删除网段内IP
	tx = tx.Unscoped().Delete(&model.IP{}, "ip_network_id = ?", id)
	if err = tx.Error; err != nil {
		tx.Rollback()
		repo.log.Error(err)
		return 0, err
	}
	return affected, tx.Commit().Error
}

// SaveIPNetwork 新增/更新网段
func (repo *MySQLRepo) SaveIPNetwork(ipn *model.IPNetwork) (affected int64, err error) {

	tx := repo.db.Begin()
	if ipn.ID > 0 {
		_, err := repo.GetIPNetworkByID(ipn.ID)
		if err != nil {
			return 0, err
		}
		// 更新网段信息，为防止IP覆盖式写入，只允许修改覆盖交换机、网关、数据中心、机房管理单元、VLAN
		if err = tx.Model(&model.IPNetwork{}).Where("id = ?", ipn.ID).
			Update("switches", ipn.Switches).
			Update("gateway", ipn.Gateway).
			Update("idc_id", ipn.IDCID).
			Update("server_room_id", ipn.ServerRoomID).
			Update("vlan", ipn.Vlan).Error; err != nil {
			tx.Rollback()
			repo.log.Error(err)
			return 0, err
		}
	} else {
		switch ipn.Version { // 针对IPv4 IPv6分别处理
		case model.IPv6:
			if err = tx.Create(ipn).Error; err != nil {
				tx.Rollback()
				repo.log.Error(err)
				return affected, err
			}

		default: //新增网段同时裂解IP资源池 -- 默认仅针对IPv4作处理
			//获取全部IP的列表
			ips, err := nw.GetIPListByMinAndMaxIP(ipn.IPPool)
			min, max := nw.GetCidrIPFirstAndLast(ipn.CIDR)
			if err != nil {
				repo.log.Error(err)
				return 0, err
			}
			//获取pxeIP的列表
			var pxeIps []string
			if ipn.PXEPool != "" {
				pxeIps, err = nw.GetIPListByMinAndMaxIP(ipn.PXEPool)
				if err != nil {
					repo.log.Error(err)
					return 0, err
				}
			}
			if err = tx.Create(ipn).Error; err != nil {
				tx.Rollback()
				repo.log.Error(err)
				return affected, err
			}
			affected = tx.RowsAffected
			// 虚拟化网段、ILO网段不做IP分配使用
			if ipn.Category == model.VIntranet || ipn.Category == model.VExtranet || ipn.Category == model.ILO {
				return affected, tx.Commit().Error
			}
	
			var scope string
			if strings.Contains(ipn.Category, "intranet") {
				scope = model.IPScopeIntranet
			} else if strings.Contains(ipn.Category, "extranet") {
				scope = model.IPScopeExtranet
			}
	
			for i := range ips {
				if ips[i] == min || ips[i] == max {
					continue
				}
	
				if err = tx.Create(&model.IP{
					IPNetworkID: ipn.ID,
					IP:          ips[i],
					Scope:       &scope,
					Category:    model.BusinessIP,
					IsUsed:      model.IPNotUsed,
					ReleaseDate: time.Now(),
				}).Error; err != nil {
					tx.Rollback()
					repo.log.Error(err)
					return 0, err
				}
			}
			//对于PXE网段，更新IP类别
			for i := range pxeIps {
				if ips[i] == min || ips[i] == max {
					continue
				}
	
				if err = tx.Model(&model.IP{}).Where("ip_network_id = ? and ip = ?", ipn.ID, pxeIps[i]).Updates(model.IP{
					Category: model.PXEIP,
				}).Error; err != nil {
					tx.Rollback()
					repo.log.Error(err)
					return 0, err
				}
			}
		}
	}
	return affected, tx.Commit().Error
}

// GetIPNetworkByID 返回指定ID的网段
func (repo *MySQLRepo) GetIPNetworkByID(id uint) (*model.IPNetwork, error) {
	var ipn model.IPNetwork
	if err := repo.db.Where("id = ?", id).Find(&ipn).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &ipn, nil
}

func addIPNetworkCond(db *gorm.DB, cond *model.IPNetworkCond) *gorm.DB {
	if cond != nil {
		if cond.ServerRoomName != "" {
	        db = db.Joins("left join server_room on server_room.id = ip_network.server_room_id")
	    }
		db = MultiNumQuery(db, "server_room_id", cond.ServerRoomID)
		db = MultiEnumQuery(db, "ip_network.category", cond.Category)
		db = MultiQuery(db, "ip_network.cidr", cond.CIDR)
		db = MultiQuery(db, "server_room.name", cond.ServerRoomName)
		if cond.Switches != "" {
			db = MultiQuery(db, "ip_network.switches", cond.Switches)
		}
	}
	return db
}

// CountIPNetworks 统计满足过滤条件的网段数量
func (repo *MySQLRepo) CountIPNetworks(cond *model.IPNetworkCond) (count int64, err error) {
	db := repo.db.Model(&model.IPNetwork{})
	db = addIPNetworkCond(db, cond)
	if err = db.Select("id").Count(&count).Error; err != nil {
		repo.log.Error(err)
		return count, err
	}
	return count, nil
}

// GetIPNetworks 返回满足过滤条件的网段
func (repo *MySQLRepo) GetIPNetworks(cond *model.IPNetworkCond, orderby model.OrderBy, limiter *page.Limiter) (items []*model.IPNetwork, err error) {
	db := repo.db.Model(&model.IPNetwork{}) //.Joins("left join network_device on ip_network.switches LIKE concat(%,network_device.name,%)")
	db = addIPNetworkCond(db, cond)
	for i := range orderby {
		db = db.Order(orderby[i].String())
	}
	if limiter != nil {
		db = db.Limit(limiter.Limit).Offset(limiter.Offset)
	}

	if err = db.Select("ip_network.*").Find(&items).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return items, nil
}

// GetIPNetworksBySwitchNumber 根据设备编号查询网段信息
func (repo *MySQLRepo) GetIPNetworksBySwitchNumber(switchNum string) (items []*model.IPNetwork, err error) {
	db := repo.db.Model(&model.IPNetwork{}).Where("switches like ?", "%"+switchNum+"%")

	if err = db.Find(&items).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return items, nil
}


// 查询指定物理机的内/外网IPv6所属网段
func (repo *MySQLRepo) GetIPv6NetworkBySN(sn string, category string) (ipnet *model.IPNetwork, err error) {
	var switcher *model.NetworkDevice
	if category == model.Intranet {
		switcher, err = repo.GetIntranetSwitchBySN(sn)
	} else if category == model.Extranet {
		switcher, err = repo.GetExtranetSwitchBySN(sn)
	}
	if err != nil || switcher == nil {
		return nil, fmt.Errorf("查询设备(sn:%s)所在机位关联的交换机失败", sn)
	}

	var items []*model.IPNetwork
	if err = repo.db.Model(&model.IPNetwork{}).Where("version = ? and switches LIKE ?", model.IPv6, "%"+switcher.FixedAssetNumber+"%").Find(&items).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	for i := range items {
		var switches []string
		if err = json.Unmarshal([]byte(items[i].Switches), &switches); err != nil {
			repo.log.Error(err)
			continue
		}
		for _, s := range switches {
			if strings.TrimSpace(s) == strings.TrimSpace(switcher.FixedAssetNumber) {
				// 仅处理一个网段
				ipnet = items[i]
				break
			}
		}
	}
	if ipnet != nil {
		return ipnet, nil
	}
	return nil, fmt.Errorf("没有找到覆盖交换机(固资号:%s)的IPv6网段", switcher.FixedAssetNumber)
}