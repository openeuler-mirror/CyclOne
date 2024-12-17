package mysqlrepo

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/voidint/page"

	"idcos.io/cloudboot/model"
	mystrings "idcos.io/cloudboot/utils/strings"
)

const dateLayout = "2006-01-02"

// GetInspections 查询满足过滤条件的硬件巡检列表
func (repo *MySQLRepo) GetInspections(cond *model.Inspection, orderby model.OrderBy, limiter *page.Limiter) (items []*model.Inspection, err error) {
	db := repo.db.Model(&model.Inspection{})
	if cond != nil {
		if cond.JobID != "" {
			db = db.Where("job_id = ?", cond.JobID)
		}
		if cond.StartTime != nil {
			db = db.
				Where("start_time >= ?", fmt.Sprintf("%s 00:00:00", cond.StartTime.Format(dateLayout))).
				Where("start_time <= ?", fmt.Sprintf("%s 23:59:59", cond.StartTime.Format(dateLayout)))
		}
		if cond.OriginNode != "" {
			db = db.Where("origin_node = ?", cond.OriginNode)
		}
		if cond.SN != "" {
			db = db.Where("sn IN (?)", strings.Split(cond.SN, ","))
		}
		if cond.RunningStatus != "" {
			db = db.Where("running_status = ?", cond.RunningStatus)
		}
		if cond.HealthStatus != "" {
			db = db.Where("health_status = ?", cond.HealthStatus)
		}
	}

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

// GetInspectionStatistics 查询满足过滤条件的硬件巡检列表用于健康状态统计
func (repo *MySQLRepo) GetInspectionStatistics(cond *model.Inspection, orderby model.OrderBy, limiter *page.Limiter) (items []*model.InspectionStatistics, err error) {
	db := repo.db.Model(&model.Inspection{})
	if cond != nil {
		if cond.JobID != "" {
			db = db.Where("job_id = ?", cond.JobID)
		}
		if cond.StartTime != nil {
			db = db.
				Where("start_time >= ?", fmt.Sprintf("%s 00:00:00", cond.StartTime.Format(dateLayout))).
				Where("start_time <= ?", fmt.Sprintf("%s 23:59:59", cond.StartTime.Format(dateLayout)))
		}
		if cond.OriginNode != "" {
			db = db.Where("origin_node = ?", cond.OriginNode)
		}
		if cond.SN != "" {
			db = db.Where("sn IN (?)", strings.Split(cond.SN, ","))
		}
		if cond.RunningStatus != "" {
			db = db.Where("running_status = ?", cond.RunningStatus)
		}
		if cond.HealthStatus != "" {
			db = db.Where("health_status = ?", cond.HealthStatus)
		}
	}

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


// GetInspectionStatisticsGroupBySN 查询满足过滤条件的硬件巡检列表用于健康状态统计(SN分组返回最新)
func (repo *MySQLRepo) GetInspectionStatisticsGroupBySN(cond *model.Inspection, orderby model.OrderBy, limiter *page.Limiter) (items []*model.InspectionStatistics, err error) {
	db := repo.db.Table("inspection")
	if cond != nil {
		if cond.StartTime != nil && cond.RunningStatus != "" {
			db = db.
				Where("start_time >= ?", fmt.Sprintf("%s 00:00:00", cond.StartTime.Format(dateLayout))).
				Where("start_time <= ?", fmt.Sprintf("%s 23:59:59", cond.StartTime.Format(dateLayout))).
				Where("running_status = ?", cond.RunningStatus)
			// string=(SELECT MAX(id) as id,sn FROM `inspection`  WHERE (start_time >= ?) AND (start_time <= ?) AND (running_status = ?) GROUP BY sn)
			// 通过GROUP BY + MAX(ID) 的子查询方式直接获取各SN最新的巡检结果
			subQuery := db.Select("MAX(id) as id,sn").Group("sn").SubQuery()
			db = db.Select("inspection.id, inspection.start_time, inspection.end_time, inspection.sn, inspection.running_status, inspection.health_status").
				Joins("inner join ? as latest_inspection on inspection.id = latest_inspection.id", subQuery)
		}
	}

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


// AddInspections 批量新增硬件巡检记录
func (repo *MySQLRepo) AddInspections(items ...*model.Inspection) (err error) {
	tx := repo.db.Begin()
	for i := range items {
		if err = tx.Create(items[i]).Error; err != nil {
			tx.Rollback()
			repo.log.Error(err)
			return err
		}
	}
	return tx.Commit().Error
}

// UpdateInspectionByID 根据ID更新硬件巡检记录
func (repo *MySQLRepo) UpdateInspectionByID(one *model.Inspection) (affected int64, err error) {
	db := repo.db.Model(&model.Inspection{}).Where("id = ?", one.ID).Updates(one)
	if err = db.Error; err != nil {
		repo.log.Error(err)
	}
	return db.RowsAffected, db.Error
}

// GetInspectedSN 查询已经执行过硬件巡检的设备SN列表
func (repo *MySQLRepo) GetInspectedSN() (items []string, err error) {
	if err = repo.db.Raw("SELECT DISTINCT(sn) FROM inspection WHERE deleted_at IS NULL").Pluck("sn", &items).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return items, nil
}

//GetInspectionStartTimeBySN 根据设备SN获取巡检开始时间
func (repo *MySQLRepo) GetInspectionStartTimeBySN(SN, runStatus string) (starts []time.Time, err error) {
	sql := fmt.Sprintf("select distinct(start_time) from inspection where sn='%s'",
		SN)
	if runStatus != "" {
		sql += fmt.Sprintf(" and inspection.running_status = '%s'", runStatus)
	}
	sql += " and deleted_at is null order by start_time desc"

	var mods []struct {
		StartTime time.Time `gorn:"column:start_time"`
	}
	err = repo.db.Raw(sql).Scan(&mods).Error
	if len(mods) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	for _, mod := range mods {
		starts = append(starts, mod.StartTime)
	}
	return
}

//GetInspectionDetail 获取巡检详情
func (repo *MySQLRepo) GetInspectionDetail(SN, startTime string) (i *model.Inspection, err error) {
	i = new(model.Inspection)
	db := repo.db.Model(model.Inspection{}).Order("start_time DESC").Where("sn = ?", SN)
	if startTime != "" {
		db = db.Where("start_time = ?", startTime)
		err = db.Find(i).Error
	} else {
		err = db.First(i).Error
	}
	return
}

// CountInspectionRecordsPage 获取巡检记录条数
func (repo *MySQLRepo) CountInspectionRecordsPage(cond *model.InspectionCond) (count int64, err error) {

	db := repo.db.Model(&model.Inspection{}).Order("start_time DESC").Where("running_status = ?", "done")

	var whereSQL strings.Builder
	if cond != nil {
		if cond.SN != "" {
			sns := strings.Split(cond.SN, ",")
			for i := range sns {
				sns[i] = fmt.Sprintf("'%s'", sns[i])
			}
			whereSQL.WriteString(fmt.Sprintf(" sn IN(%s)", strings.Join(sns, ",")))
			db = db.Where(whereSQL.String())
		}
		if cond.StartTime != "" {
			whereSQL.WriteString(fmt.Sprintf(" start_time > '%s'", cond.StartTime))
			db = db.Where(whereSQL.String())
		}
		if cond.EndTime != "" {
			whereSQL.WriteString(fmt.Sprintf(" end_time < '%s'", cond.EndTime))
			db = db.Where(whereSQL.String())
		}
		// enum('nominal','warning','critical','unknown') 
		if cond.HealthStatus != "" {
			whereSQL.WriteString(fmt.Sprintf(" health_status = '%s'", cond.HealthStatus))
			db = db.Where(whereSQL.String())
		}
	}

    // 返回结果
	if err = db.Count(&count).Error; err != nil {
		repo.log.Error(err)
		return 0, err
	}
	return count, nil
}

// GetInspectionPage 获取巡检记录分页
func (repo *MySQLRepo) GetInspectionRecordsPage(cond *model.InspectionCond, limiter *page.Limiter) (items []*model.InspectionRecordsPage, err error) {

	db := repo.db.Model(&model.Inspection{}).Order("start_time DESC").Where("running_status = ?", "done")

	var whereSQL strings.Builder
	if cond != nil {
		if cond.SN != "" {
			sns := strings.Split(cond.SN, ",")
			for i := range sns {
				sns[i] = fmt.Sprintf("'%s'", sns[i])
			}
			whereSQL.WriteString(fmt.Sprintf(" sn IN(%s)", strings.Join(sns, ",")))
			db = db.Where(whereSQL.String())
		}
		if cond.StartTime != "" {
			whereSQL.WriteString(fmt.Sprintf(" start_time > '%s'", cond.StartTime))
			db = db.Where(whereSQL.String())
		}
		if cond.EndTime != "" {
			whereSQL.WriteString(fmt.Sprintf(" end_time < '%s'", cond.EndTime))
			db = db.Where(whereSQL.String())
		}
		// enum('running','done')
		//if cond.RuningStatus != "" {
		//	whereSQL.WriteString(fmt.Sprintf(" AND t10.running_status = '%s'", cond.RuningStatus))
		//}
		// enum('nominal','warning','critical','unknown') 
		if cond.HealthStatus != "" {
			whereSQL.WriteString(fmt.Sprintf(" health_status = '%s'", cond.HealthStatus))
			db = db.Where(whereSQL.String())
		}
	}
	
	// 分页
	if limiter != nil {
		db = db.Limit(limiter.Limit).Offset(limiter.Offset)
	}
    // 返回结果
	if err = db.Find(&items).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return items, nil
}

//GetInspectionListWithPageNew 分页获取硬件巡检结果，每个SN获取最新的一个
func (repo *MySQLRepo) GetInspectionListWithPageNew(cond *model.InspectionCond, limiter *page.Limiter) (result []model.InspectionFullWithPage, err error) {
	//获取硬件巡检中的所有sn（单一）
	//var sn []string
	//var id []uint64
	//err = repo.db.Raw(fmt.Sprintf("select distinct sn from inspection where sn != \"\"")).Pluck("sn", &sn).Error
	//if err != nil {
	//	return nil, err
	//}

	////获取排序后的结果
	//var mod []model.InspectionTemp
	//err = repo.db.Select([]string{"id", "sn"}).Order("start_time desc").Find(&mod).Error
	//if err != nil {
	//	return nil, err
	//}

	//for _, vsn := range sn {
	//	for _, vm := range mod {
	//		if vsn == vm.SN {
	//			id = append(id, vm.ID)
	//			break
	//		}
	//	}
	//}

	//var sql = `select t10.sn as sn, 
	//	t10.start_time, t10.end_time, 
	//	t10.id, t10.running_status, 
	//	t10.error, t10.ipmi_result, 
	//	t10.health_status,
	//	t10.created_at, t10.updated_at 
	//	from device t5 inner join inspection t10 on t5.sn = t10.sn 
	//	where t10.id > 0 %s order by t10.end_time desc %s`
	
	var sql = `select t5.fixed_asset_number,
			   t10.sn as sn, 
			   t20.intranet_ip,
			   t10.start_time,
			   t10.end_time,
			   t10.id,
			   t10.running_status,
			   t10.error,
			   t10.ipmi_result,
			   t10.health_status 
			   from (select max(id) as id,sn from inspection group by sn) as t15 inner join inspection t10 on t15.id = t10.id 
			   inner join device t5 on t5.sn = t10.sn left join device_setting t20 on t5.sn = t20.sn
			   where t10.id > 0 %s %s`

	var whereSQL strings.Builder

	//if id != nil {
	//	var ids []string
	//	for i := range id {
	//		ids = append(ids, fmt.Sprintf("'%d'", id[i]))
	//	}
	//	whereSQL.WriteString(fmt.Sprintf(" AND t10.id IN(%s)", strings.Join(ids, ",")))
	//}

	if cond != nil {
		if cond.SN != "" {
			sns := mystrings.MultiLines2Slice(cond.SN)
			for i := range sns {
				sns[i] = fmt.Sprintf("'%s'", sns[i])
			}
			whereSQL.WriteString(fmt.Sprintf(" AND t5.sn IN(%s)", strings.Join(sns, ",")))
		}
		if cond.FixedAssetNumber != "" {
			fns := mystrings.MultiLines2Slice(cond.FixedAssetNumber)
			for i := range fns {
				fns[i] = fmt.Sprintf("'%s'", fns[i])
			}
			whereSQL.WriteString(fmt.Sprintf(" AND t5.fixed_asset_number IN(%s)", strings.Join(fns, ",")))
		}
		if cond.IntranetIP != "" {
			inips := mystrings.MultiLines2Slice(cond.IntranetIP)
			for i := range inips {
				inips[i] = fmt.Sprintf("'%s'", inips[i])
			}
			whereSQL.WriteString(fmt.Sprintf(" AND t20.intranet_ip IN(%s)", strings.Join(inips, ",")))
		}		
		if cond.StartTime != "" {
			whereSQL.WriteString(fmt.Sprintf(" AND t10.start_time > '%s'", cond.StartTime))
		}
		if cond.EndTime != "" {
			whereSQL.WriteString(fmt.Sprintf(" AND t10.end_time < '%s'", cond.EndTime))
		}
		// if cond.OOBIP != "" {
		// 	whereSQL.WriteString(fmt.Sprintf(" AND t6.manage_ip = '%s'", cond.OOBIP))
		// }
		if cond.RuningStatus != "" {
			whereSQL.WriteString(fmt.Sprintf(" AND t10.running_status = '%s'", cond.RuningStatus))
		}
		if cond.HealthStatus != "" {
			whereSQL.WriteString(fmt.Sprintf(" AND t10.health_status = '%s'", cond.HealthStatus))
		}
	}

	var limitSQL string
	if limiter != nil {
		limitSQL = fmt.Sprintf(" LIMIT %d,%d ", limiter.Offset, limiter.Limit)
	}
	//repo.log.Info(fmt.Sprintf(sql, whereSQL.String(), limitSQL))
	if err := repo.db.Raw(fmt.Sprintf(sql, whereSQL.String(), limitSQL)).Scan(&result).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}

	return result, nil
}

//CountInspctionsByCond 根据条件统计
func (repo *MySQLRepo) CountInspctionsByCond(cond *model.InspectionCond) (count int64, err error) {

	//var sn []string
	//var id []uint64
	//err = repo.db.Raw(fmt.Sprintf("select distinct sn from inspection where sn != \"\"")).Pluck("sn", &sn).Error
	//if err != nil {
	//	return 0, err
	//}
//
	////获取排序后的结果
	//var mod []model.InspectionTemp
	//err = repo.db.Select([]string{"id", "sn"}).Order("start_time desc").Find(&mod).Error
	//if err != nil {
	//	return 0, err
	//}
//
	//for _, vsn := range sn {
	//	for _, vm := range mod {
	//		if vsn == vm.SN {
	//			id = append(id, vm.ID)
	//			break
	//		}
	//	}
	//}
//
	//var sql = `select count(t10.sn) as count from device t5 inner join inspection t10 on t5.sn = t10.sn 
	//where t10.id > 0 %s`

	var sql=`select count(t10.sn) as count from 
			 (select max(id) as id,sn from inspection group by sn) as t15 inner join inspection t10 on t15.id = t10.id 
			 inner join device t5 on t5.sn = t10.sn 
			 left join device_setting t20 on t5.sn = t20.sn
			 where t10.id > 0 %s`

	var whereSQL strings.Builder

	//if id != nil {
	//	var ids []string
	//	for i := range id {
	//		ids = append(ids, fmt.Sprintf("'%d'", id[i]))
	//	}
	//	whereSQL.WriteString(fmt.Sprintf(" AND t10.id IN(%s)", strings.Join(ids, ",")))
	//}

	if cond != nil {
		if cond.SN != "" {
			sns := mystrings.MultiLines2Slice(cond.SN)
			for i := range sns {
				sns[i] = fmt.Sprintf("'%s'", sns[i])
			}
			whereSQL.WriteString(fmt.Sprintf(" AND t5.sn IN(%s)", strings.Join(sns, ",")))
		}
		if cond.FixedAssetNumber != "" {
			fns := mystrings.MultiLines2Slice(cond.FixedAssetNumber)
			for i := range fns {
				fns[i] = fmt.Sprintf("'%s'", fns[i])
			}
			whereSQL.WriteString(fmt.Sprintf(" AND t5.fixed_asset_number IN(%s)", strings.Join(fns, ",")))
		}
		if cond.IntranetIP != "" {
			inips := mystrings.MultiLines2Slice(cond.IntranetIP)
			for i := range inips {
				inips[i] = fmt.Sprintf("'%s'", inips[i])
			}
			whereSQL.WriteString(fmt.Sprintf(" AND t20.intranet_ip IN(%s)", strings.Join(inips, ",")))
		}
		if cond.StartTime != "" {
			whereSQL.WriteString(fmt.Sprintf(" AND t10.start_time > '%s'", cond.StartTime))
		}
		if cond.EndTime != "" {
			whereSQL.WriteString(fmt.Sprintf(" AND t10.end_time < '%s'", cond.EndTime))
		}
		// if cond.OOBIP != "" {
		// 	whereSQL.WriteString(fmt.Sprintf(" AND t6.manage_ip = '%s'", cond.OOBIP))
		// }
		if cond.RuningStatus != "" {
			whereSQL.WriteString(fmt.Sprintf(" AND t10.running_status = '%s'", cond.RuningStatus))
		}
		if cond.HealthStatus != "" {
			whereSQL.WriteString(fmt.Sprintf(" AND t10.health_status = '%s'", cond.HealthStatus))
		}
	}
	//repo.log.Info(fmt.Sprintf(sql, whereSQL.String()))
	if err = repo.db.DB().QueryRow(fmt.Sprintf(sql, whereSQL.String())).Scan(&count); err != nil {
		repo.log.Error(err)
		return 0, err
	}
	return count, nil
}

//RemoveInspectionOnStartTimeBySN 根据设备SN删除巡检记录，按时间排序保留14条记录
func (repo *MySQLRepo) RemoveInspectionOnStartTimeBySN(SN string) (err error) {
	sql := fmt.Sprintf("delete from inspection where sn='%s' and id not in (select t.id from (select id from inspection where sn='%s' order by start_time desc limit 14) as t)",
		SN, SN)
    //repo.log.Info(sql)
    if err := repo.db.Exec(sql).Error; err != nil {
		repo.log.Error(err)
		return err
	}
	return nil
}