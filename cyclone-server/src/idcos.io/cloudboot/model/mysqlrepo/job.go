package mysqlrepo

import (
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/voidint/page"

	"idcos.io/cloudboot/model"
)

// SaveJob 保存(新增/更新)任务
func (repo *MySQLRepo) SaveJob(job *model.Job) (err error) {
	var affected int64
	if job.ID != "" {
		db := repo.db.Model(&model.Job{}).Updates(job) // 根据主键ID更新任务
		if err = db.Error; err != nil {
			repo.log.Error(err)
			return err
		}
		if affected = db.RowsAffected; affected > 0 {
			return nil
		}
	}

	if err = repo.db.Create(job).Error; err != nil {
		repo.log.Error(err)
		return err
	}
	return nil
}

// RemoveJob 移除指定ID的任务
func (repo *MySQLRepo) RemoveJob(id string) (affected int64, err error) {
	db := repo.db.Delete(&model.Job{}, "id = ?", id)
	if err = db.Error; err != nil {
		repo.log.Error(err)
	}
	return db.RowsAffected, db.Error
}

// GetJobByID 查询指定ID的任务
func (repo *MySQLRepo) GetJobByID(id string) (job *model.Job, err error) {
	var one model.Job
	if err := repo.db.Where("id = ?", id).Find(&one).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &one, nil
}

// CountJobs 统计满足条件的任务数量
func (repo *MySQLRepo) CountJobs(cond *model.Job) (count int64, err error) {
	db := repo.db.Model(&model.Job{})
	db = repo.setWhereSQL4Job(db, cond)

	if err = db.Select("id").Count(&count).Error; err != nil {
		repo.log.Error(err)
		return count, err
	}
	return count, nil
}

// GetJobs 返回满足过滤条件的任务
func (repo *MySQLRepo) GetJobs(cond *model.Job, orderby model.OrderBy, limiter *page.Limiter) (items []*model.Job, err error) {
	db := repo.db.Model(&model.Job{})
	db = repo.setWhereSQL4Job(db, cond)

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

func (repo *MySQLRepo) setWhereSQL4Job(db *gorm.DB, cond *model.Job) *gorm.DB {
	if db == nil || cond == nil {
		return db
	}
	if cond.Title != "" {
		db = db.Where("title LIKE ?", "%"+cond.Title+"%")
	}
	if cond.Builtin != "" {
		db = db.Where("builtin = ?", cond.Builtin)
	}
	if cond.Category != "" {
		db = db.Where("category = ?", cond.Category)
	}
	if cond.Rate != "" {
		db = db.Where("rate = ?", cond.Rate)
	}
	if cond.Status != "" {
		db = db.Where("`status` IN(?)", strings.Split(cond.Status, ","))
	}
	if cond.Creator != "" {
		db = db.Where("creator = ?", cond.Creator)
	}
	return db
}
