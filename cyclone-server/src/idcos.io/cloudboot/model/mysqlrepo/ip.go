package mysqlrepo

import (
	"github.com/jinzhu/gorm"
	"github.com/voidint/page"

	"bytes"
	"fmt"
	"time"
	strings2 "strings"
	"errors"

	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils/ping"
	"idcos.io/cloudboot/utils/strings"
	"idcos.io/cloudboot/utils/network"
)

// AssignIntranetIP 按照内置规则给设备分配一个内网业务IP
func (repo *MySQLRepo) AssignIntranetIP(sn string) (*model.IP, error) {
	// 1、若设备已存在关联的内网业务IP，则直接返回该IP。
	var intraIP model.IP
	err := repo.db.Model(&model.IP{}).Where("sn = ? AND scope = ? AND category = ? ", sn, model.IPScopeIntranet, model.BusinessIP).First(&intraIP).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		repo.log.Error(err)
		return nil, err
	}
	if intraIP.ID > 0 && intraIP.IP != "" {
		return &intraIP, nil
	}

	// 2、获取待分配IP所属内网网段
	var intranetid uint
	intranet, err := repo.GetIntranetIPNetworksBySN(sn)
	if err != nil || len(intranet) == 0 {
		repo.log.Errorf("intranet segment not find by sn:%s", sn)
		return nil, err
	}

	// 3、尝试释放设备已占用的内网业务IP（可选）
	_, _ = repo.ReleaseIP(sn, model.IPScopeIntranet)

	lenintranet := len(intranet) - 1
	isPrint := false
	tx0 := repo.db.Begin()
	for k := range intranet { //多数情况这里只有一个网段
		if intranet[k].Category == model.VIntranet || intranet[k].Category == model.VExtranet || intranet[k].Category == model.ILO {
			continue
		}
		if lenintranet == k {
			isPrint = true
		}

		intranetid = intranet[k].ID

		//需要加ping检测
		ips := make([]*model.IP, 0)
		// 4、从目标内网IP网段中选择一个空闲业务IP并加排它锁
		if err = tx0.Model(&model.IP{}).Where("ip_network_id = ? AND category = ? AND is_used = ?",
			intranetid, model.BusinessIP, model.IPNotUsed).Order("id ASC").
			Set("gorm:query_option", "FOR UPDATE").Find(&ips).Error; err != nil {
			tx0.Rollback()
			if isPrint {
				repo.log.Error(err)
				return nil, err
			}
			continue
		}

		// 5、将目标业务IP与设备进行关联并释放排它锁
		for _, ip := range ips {
			if ping.PingTest(ip.IP) == nil {
				repo.log.Errorf("ping ip: %s ok, skip this one", ip.IP)
				continue
			}
			if err = tx0.Model(&model.IP{}).Where("id = ?", ip.ID).Updates(map[string]interface{}{
				"is_used": model.IPUsed,
				"sn":      sn,
				//"scope":   model.IPScopeIntranet,
			}).Error; err != nil {
				tx0.Rollback()
				if isPrint {
					repo.log.Error(err)
					return nil, err
				}
				continue
			}
			return ip, tx0.Commit().Error //这里找到就return
		}
	}
	tx0.Commit()
	return nil, errors.New("未找到匹配IP")
}

// AssignExtranetIP 按照内置规则给设备分配一个外网业务IP
func (repo *MySQLRepo) AssignExtranetIP(sn string) (*model.IP, error) {
	// 1、若设备已存在关联的外网业务IP，则直接返回该IP。
	var extraIP model.IP
	err := repo.db.Model(&model.IP{}).Where("sn = ? AND scope = ? AND category = ? ", sn, model.IPScopeExtranet, model.BusinessIP).First(&extraIP).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		repo.log.Error(err)
		return nil, err
	}
	if extraIP.ID > 0 && extraIP.IP != "" {
		return &extraIP, nil
	}

	// 2、获取待分配IP所属外网网段
	var extranetid uint
	extranet, err := repo.GetExtranetIPNetworksBySN(sn)
	if err != nil || len(extranet) == 0 {
		repo.log.Errorf("extranet segment not find by sn:%s", sn)
		return nil, err
	}

	// 3、尝试释放设备已占用的外网业务IP（可选）
	_, _ = repo.ReleaseIP(sn, model.IPScopeExtranet)

	lenextranet := len(extranet) - 1
	isPrint := false
	tx0 := repo.db.Begin()
	for k := range extranet {
		if extranet[k].Category == model.VIntranet || extranet[k].Category == model.VExtranet || extranet[k].Category == model.ILO {
			continue
		}
		if lenextranet == k {
			isPrint = true
		}
		extranetid = extranet[k].ID

		ips := make([]*model.IP, 0)
		// 4、从目标外网IP网段中选择一个空闲业务IP并加排它锁
		if err = tx0.Model(&model.IP{}).Where("ip_network_id = ? AND category = ? AND is_used = ?",
			extranetid, model.BusinessIP, model.IPNotUsed).Order("id ASC").
			Set("gorm:query_option", "FOR UPDATE").Find(&ips).Error; err != nil {
			tx0.Rollback()
			if isPrint {
				repo.log.Error(err)
				return nil, err
			}
			continue
		}

		// 5、将目标业务IP与设备进行关联并释放排它锁
		for _, ip := range ips {
			if ping.PingTest(ip.IP) == nil {
				repo.log.Errorf("ping ip:%s ok, skip this one")
				continue
			}
			if err = tx0.Model(&model.IP{}).Where("id = ?", ip.ID).
				Updates(map[string]interface{}{
					"is_used": model.IPUsed,
					"sn":      sn,
					//"scope":   model.IPScopeExtranet,
				}).Error; err != nil {
				tx0.Rollback()
				if isPrint {
					repo.log.Error(err)
					return nil, err
				}
				continue
			}
			return ip, tx0.Commit().Error //这里找到就return
		}
	}
	return &extraIP, tx0.Commit().Error
}

// ReleaseIP 为目标设备释放内/外网业务IP，有多个IP是一并清空
func (repo *MySQLRepo) ReleaseIP(sn string, scope string) (affected int64, err error) {
	db := repo.db.Model(&model.IP{}).Where("sn = ? AND scope = ? AND category = ?", sn, scope, model.BusinessIP).Updates(map[string]interface{}{
		"is_used": model.NO,
		"sn":      "",
		//"scope":   nil,
	})

	if err = db.Error; err != nil {
		repo.log.Error(err)
		return 0, err
	}

	//同步清理设备参数表的记录
	if scope == model.IPScopeExtranet {
		repo.db.Model(&model.DeviceSetting{}).Where("sn = ?", sn).
			Update("need_extranet_ip", model.NO).
			Update("extranet_ip_network_id", 0).
			Update("extranet_ip", "")
	} else if scope == model.IPScopeIntranet {
		repo.db.Model(&model.DeviceSetting{}).Where("sn = ?", sn).
			Update("intranet_ip_network_id", 0).
			Update("intranet_ip", "")
	}
	return db.RowsAffected, nil
}

// ReserveIP 为目标设备回收内/外网业务IP并保留IP一段时间，设置释放日期
func (repo *MySQLRepo) ReserveIP(sn string, scope string, releasedate time.Time) (affected int64, err error) {
	db := repo.db.Model(&model.IP{}).Where("sn = ? AND scope = ? AND category = ?", sn, scope, model.BusinessIP).Updates(map[string]interface{}{
		"is_used":                 model.IPDisabled,
		"sn":                      "",
		"release_date":            releasedate,
	})

	if err = db.Error; err != nil {
		repo.log.Error(err)
		return 0, err
	}

	//同步清理设备参数表的记录
	if scope == model.IPScopeExtranet {
		repo.db.Model(&model.DeviceSetting{}).Where("sn = ?", sn).
			Update("need_extranet_ip", model.NO).
			Update("extranet_ip_network_id", 0).
			Update("extranet_ip", "")
	} else if scope == model.IPScopeIntranet {
		repo.db.Model(&model.DeviceSetting{}).Where("sn = ?", sn).
			Update("intranet_ip_network_id", 0).
			Update("intranet_ip", "")
	}
	return db.RowsAffected, nil
}

// GetReleasableIP 根据release_date 以及is_used=model.IPDisabled 获取IP记录
func (repo *MySQLRepo) GetReleasableIP() (items []*model.IP, err error) {
	db := repo.db.Model(&model.IP{})
	if err = db.Where("is_used = ? AND category = ? AND datediff(now(),release_date) >= 0", model.IPDisabled, model.BusinessIP).Find(&items).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return items, nil
}

// addIPCond 添加查询条件
func addIPCond(db *gorm.DB, cond *model.IPPageCond) *gorm.DB {
	if cond != nil {
		//网段ID查询
		db = MultiNumQuery(db, "ip.ip_network_id", cond.IPNetworkID)
		//IP ID查询
		db = MultiNumQuery(db, "ip.id", cond.ID)
		//网段类别、网段名称查询
		if cond.Category != "" || cond.CIDR != "" {
			db = db.Joins("left join ip_network on ip.ip_network_id = ip_network.id")
		}
		if cond.Category != "" {
			db = MultiEnumQuery(db, "ip_network.category", cond.Category)
		}
		if cond.CIDR != "" {
			db = MultiQuery(db, "ip_network.cidr", cond.CIDR)
		}
		//IP地址精确查询
		db = MultiMatchQuery(db, "ip.ip", cond.IP)
		if cond.Scope != nil {
			db = MultiQuery(db, "ip.scope", *cond.Scope)
		}
		db = MultiEnumQuery(db, "ip.is_used", cond.IsUsed)
		db = MultiQuery(db, "ip.sn", cond.SN)
		if cond.FixedAssetNumber != "" {
			//db = MultiQuery(db, "device.fixed_asset_number", cond.FixedAssetNumber)
			cs := strings.MultiLines2Slice(cond.FixedAssetNumber)
			var sb bytes.Buffer
			for i, c := range cs {
				cond := fmt.Sprintf("%%%s%%", c)
				if i == 0 {
					sb.WriteString(fmt.Sprintf("device.fixed_asset_number LIKE '%s' OR ip.remark LIKE '%s'", cond, cond))
				} else {
					sb.WriteString(fmt.Sprintf("OR device.fixed_asset_number LIKE '%s' OR ip.remark LIKE '%s'", cond, cond))
				}
			}
			db = db.Where(fmt.Sprintf("(%s)", sb.String()))
		}
	}
	return db
}

// CountIPs 统计满足过滤条件的IP数量
func (repo *MySQLRepo) CountIPs(cond *model.IPPageCond) (count int64, err error) {
	db := repo.db.Model(&model.IP{}).Joins("left join device on ip.sn = device.sn")
	db = addIPCond(db, cond)
	if err = db.Select("id").Count(&count).Error; err != nil {
		repo.log.Error(err)
		return count, err
	}
	return count, nil
}

// GetIPs 返回满足过滤条件的IP
func (repo *MySQLRepo) GetIPs(cond *model.IPPageCond, orderby model.OrderBy, limiter *page.Limiter) (items []*model.IPCombined, err error) {
	db := repo.db.Model(&model.IP{}).Joins("left join device on ip.sn = device.sn")
	db = addIPCond(db, cond)
	for i := range orderby {
		db = db.Order(orderby[i].String())
	}
	if limiter != nil {
		db = db.Limit(limiter.Limit).Offset(limiter.Offset)
	}

	//if err = db.Find(&items).Error; err != nil {
	if err = db.Select("ip.*, device.fixed_asset_number").Find(&items).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return items, nil
}

// GetIPByIP 根据IP地址精确查询IP
func (repo *MySQLRepo) GetIPByIP(ipaddr, scope string) (ip model.IP, err error) {
	db := repo.db.Model(&model.IP{})
	if err = db.Where("ip = ? AND scope = ?", ipaddr, scope).First(&ip).Error; err != nil {
		repo.log.Error(err)
		return
	}
	return
}

//GetNetWorkBySN 根据SN获取网络配置信息
func (repo *MySQLRepo) GetNetWorkBySN(sn string, category string) (ipn []model.IPAndIPNetworkUnion, err error) {
	db := repo.db.Model(&model.IP{})
	db = db.Joins("join ip_network on ip_network.id = ip.ip_network_id")
	if err := db.Select("ip, netmask, gateway, scope, version").Where("ip.sn = ? and ip.is_used = ? and ip.category = ?", sn, model.IPUsed, category).Find(&ipn).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return ipn, nil
}

// AssignIP 分配IP
func (repo *MySQLRepo) AssignIP(sn, scope string, id uint) error {
	if err := repo.db.Model(&model.IP{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_used": model.IPUsed,
		"sn":      sn,
		"scope":   scope,
	}).Error; err != nil {
		repo.log.Error(err)
		return err
	}
	return nil
}

// AssignIPByIP 分配IP
func (repo *MySQLRepo) AssignIPByIP(sn, scope, ip string) error {
	if err := repo.db.Model(&model.IP{}).Where("ip = ?", ip).Updates(map[string]interface{}{
		"is_used": model.IPUsed,
		"sn":      sn,
		"scope":   scope,
	}).Error; err != nil {
		repo.log.Error(err)
		return err
	}
	return nil
}

// UnassignIP 取消分配IP
func (repo *MySQLRepo) UnassignIP(id uint) error {
	if err := repo.db.Model(&model.IP{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_used": model.NO,
		"sn":      "",
		"remark":  nil,
	}).Error; err != nil {
		repo.log.Error(err)
		return err
	}
	return nil
}

// UnassignIPsBySN 释放指定SN的IP资源,所有内外网IP
func (repo *MySQLRepo) UnassignIPsBySN(sn string) error {
	if err := repo.db.Model(&model.IP{}).Where("sn = ?", sn).Updates(map[string]interface{}{
		"is_used": model.NO,
		"sn":      "",
		"remark":  nil,
	}).Error; err != nil {
		repo.log.Error(err)
		return err
	}

	//将装机参数表的数据清空
	repo.db.Model(&model.DeviceSetting{}).Where("sn = ?", sn).
		Update("intranet_ip_network_id", 0).
		Update("intranet_ip", "").
		Update("need_extranet_ip", model.NO).
		Update("extranet_ip_network_id", 0).
		Update("extranet_ip", "")

	return nil
}

// GetIPByID 返回指定ID的IP
func (repo *MySQLRepo) GetIPByID(id uint) (*model.IP, error) {
	var ip model.IP
	if err := repo.db.Unscoped().Where("id = ?", id).Find(&ip).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &ip, nil
}

// GetIPBySNAndScope 返回指定ID和类别的IP,
func (repo *MySQLRepo) GetIPBySNAndScope(sn, scope string) (*model.IP, error) {
	// Fixbug: 现在只能一个外网IP，一个内网IP，如果物理机可以配置多个IP，这里需要改变
	var ip model.IP
	if err := repo.db.Unscoped().Where("sn = ? and scope = ?", sn, scope).Find(&ip).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &ip, nil
}

// 更新
func (repo *MySQLRepo) SaveIP(mod *model.IP) (affected int64, err error) {
	err = repo.db.Model(&model.IP{}).Update(mod).Error
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// 根据网段（ip_network_id）获取一个空闲的业务IP
func (repo *MySQLRepo) GetAvailableIPByIPNetworkID(ipnetworkid uint) (*model.IP, error) {
	var ip model.IP
	if err := repo.db.Model(&model.IP{}).Where("ip_network_id = ? and is_used = ? and category = ?", ipnetworkid, model.IPNotUsed, model.BusinessIP).First(&ip).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return &ip, nil
}

// 根据网段（ip_network_id）获取最后一个业务IPv6
func (repo *MySQLRepo) GetLastIPv6ByIPNetworkID(ipnetworkid uint) (*model.IP, error) {
	var ip model.IP
	if err := repo.db.Model(&model.IP{}).Order("id DESC").Where("ip_network_id = ? and category = ?", ipnetworkid, model.BusinessIP).First(&ip).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return &ip, nil
}

// 创建
func (repo *MySQLRepo) CreateIP(mod *model.IP) error {
	err := repo.db.Create(mod).Error
	if err != nil {
		return err
	}
	return nil
}



// AssignIPv6 给设备分配一个业务IPv6
func (repo *MySQLRepo) AssignIPv6(sn, ipscope string) (*model.IP, error) {
	// 若设备已存在关联的业务IPv6，则直接返回该IP(默认网段唯一，仅返回1个)
	var existIPv6 model.IP
	db := repo.db.Model(&model.IP{}).Joins("left join ip_network on ip.ip_network_id = ip_network.id")
	err := db.Where("ip.sn = ? AND ip.scope = ? AND ip.category = ? AND ip_network.version = ?", sn, ipscope, model.BusinessIP, model.IPv6).First(&existIPv6).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		repo.log.Error(err)
		return nil, err
	}
	if existIPv6.ID > 0 && existIPv6.IP != "" {
		return &existIPv6, nil
	}

	// 计算并新分配一个ipv6
	var ipv6Assign *model.IP
	if ipnetwork, err := repo.GetIPv6NetworkBySN(sn, ipscope); err != nil {
		repo.log.Error(err)
		return nil, err
	} else {
		repo.log.Debugf("GetIPv6NetworkBySN ipv6:%+v", ipnetwork)
		var scope string
		if strings2.Contains(ipnetwork.Category, "intranet") {
			scope = model.IPScopeIntranet
		} else if strings2.Contains(ipnetwork.Category, "extranet") {
			scope = model.IPScopeExtranet
		}

		// 根据网段ID获取可分配的空闲IPv6
		ipv6Assign, err = repo.GetAvailableIPByIPNetworkID(ipnetwork.ID)
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			repo.log.Error(err)
			return nil, err
		}
		repo.log.Debugf("GetAvailableIPByIPNetworkID ipv6:%+v", ipv6Assign)
		// 无空闲IPv6时，获取该网段最后一个IP并计算下一个IPv6地址
		if ipv6Assign == nil {
			repo.log.Debugf("No available IPv6 of CIDR: %s .  Will generate one.", ipnetwork.CIDR)
			if ipv6Last, err := repo.GetLastIPv6ByIPNetworkID(ipnetwork.ID);err != nil && !gorm.IsRecordNotFoundError(err){
				repo.log.Error(err)
				return nil, err
			} else if ipv6Last != nil { // 最后一个IP不为空则获取下一个
				repo.log.Debugf("Last IPv6 of CIDR: %s is %+v", ipnetwork.CIDR, ipv6Last)
				if ipv6Next, err := network.GetNextIPv6OfCIDR(ipv6Last.IP, ipnetwork.CIDR); err != nil {
					repo.log.Error(err)
					return nil, err
				} else if ipv6Next != "" {
					repo.log.Debugf("Next IPv6 of CIDR: %s is %s", ipnetwork.CIDR, ipv6Next)
					ipv6Assign = &model.IP {
						IP: 			ipv6Next,
						SN:				sn,
						IPNetworkID:	ipnetwork.ID,
						Scope:			&scope,
						Category:		model.BusinessIP,
						IsUsed:			model.IPUsed,
						ReleaseDate: 	time.Now(),
					}
					err = repo.CreateIP(ipv6Assign)
					if err != nil {
						repo.log.Error(err)
						return nil, err
					}
				}
			} else {  // 获取不到最后一个IP则当新网段分配处理
				repo.log.Debugf("No IPv6 exist of CIDR: %s .  Will generate one.", ipnetwork.CIDR)
				if ipv6First, err := network.GetFirstIPv6OfCIDR(ipnetwork.CIDR); err != nil {
					repo.log.Error(err)
					return nil, err
				} else if ipv6First != "" {
					repo.log.Debugf("New IPv6 of CIDR: %s is %s", ipnetwork.CIDR, ipv6First)
					ipv6Assign = &model.IP {
						IPNetworkID:	ipnetwork.ID,
						IP: 			ipv6First,
						SN:				sn,
						Scope:			&scope,
						Category:		model.BusinessIP,
						IsUsed:			model.IPUsed,
						ReleaseDate: 	time.Now(),
					}
					err = repo.CreateIP(ipv6Assign)
					if err != nil {
						repo.log.Error(err)
						return nil, err
					}

				}
			 }
		} else { //存在空闲IPv6时
			if err := repo.db.Model(&model.IP{}).Where("id = ?", ipv6Assign.ID).Updates(map[string]interface{}{
				"is_used": model.IPUsed,
				"sn":      sn,
				"scope":   scope,
			}).Error; err != nil {
				repo.log.Error(err)
				return nil,err
			}
		}
	}
	return ipv6Assign, nil
}


// ReleaseIPv6 为目标设备释放内/外网业务IPv6，有多个IP是一并清空
func (repo *MySQLRepo) ReleaseIPv6(sn string, scope string) (affected int64, err error) {
	tx0 := repo.db.Begin()
	ips := make([]*model.IP, 0)
	if err = tx0.Model(&model.IP{}).Joins("left join ip_network on ip.ip_network_id = ip_network.id").
		Where("ip.sn = ? AND ip.scope = ? AND ip.category = ? AND ip_network.version = ?", sn, scope, model.BusinessIP, model.IPv6).
		Set("gorm:query_option", "FOR UPDATE").Find(&ips).Error; err != nil {
			tx0.Rollback()
			repo.log.Error(err)
			return 0, err
	}
	for _, ip := range ips {
		if err = tx0.Model(&model.IP{}).Where("id = ?", ip.ID).
			Updates(map[string]interface{}{
				"is_used": model.IPNotUsed,
				"sn":      "",
			}).Error; err != nil {
				tx0.Rollback()
				repo.log.Error(err)
				return 0, err
		}
	}

	if err = tx0.Commit().Error; err != nil {
		repo.log.Error(err)
		return 0, err
	}

	//同步清理设备参数表的记录
	if scope == model.IPScopeExtranet {
		repo.log.Debugf("device setting reset %s IP of SN:%s ", scope, sn)
		repo.db.Model(&model.DeviceSetting{}).Where("sn = ?", sn).
			Update("need_extranet_ipv6", model.NO).
			Update("extranet_ipv6_network_id", 0).
			Update("extranet_ipv6", "")
	} else if scope == model.IPScopeIntranet {
		repo.log.Debugf("device setting reset %s IP of SN:%s ", scope, sn)
		repo.db.Model(&model.DeviceSetting{}).Where("sn = ?", sn).
			Update("need_intranet_ipv6", model.NO).
			Update("intranet_ipv6_network_id", 0).
			Update("intranet_ipv6", "")
	}
	return tx0.RowsAffected, nil
}