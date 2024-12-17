package mysqlrepo

import (
	"github.com/jinzhu/gorm"

	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils/collection"
)

// AddTokenBuckets 新增令牌记录
func (repo *MySQLRepo) AddTokenBuckets(items ...*model.DHCPTokenBucket) (err error) {
	tx := repo.db.Begin()
	// 新增记录
	for i := range items {
		if items[i] == nil {
			continue
		}
		if err = tx.Create(items[i]).Error; err != nil {
			tx.Rollback()
			repo.log.Error(err)
			return err
		}
	}
	return tx.Commit().Error
}

// OverwriteTokenBuckets 覆写令牌桶及其令牌
func (repo *MySQLRepo) OverwriteTokenBuckets(items ...*model.DHCPTokenBucket) (err error) {
	set := collection.NewSSet(1)
	for i := range items {
		set.Add(items[i].Bucket)
	}

	if set.IsEmpty() {
		return nil
	}

	tx := repo.db.Begin()
	// 删除所有记录
	if err = tx.Unscoped().Delete(model.DHCPTokenBucket{}, "bucket IN (?)", set.Elements()).Error; err != nil {
		tx.Rollback()
		repo.log.Error(err)
		return err
	}

	// 新增记录
	for i := range items {
		if items[i] == nil {
			continue
		}
		if err = tx.Create(items[i]).Error; err != nil {
			tx.Rollback()
			repo.log.Error(err)
			return err
		}
	}
	return tx.Commit().Error
}

// BindSNByTokenBucket 绑定SN与令牌桶中令牌
func (repo *MySQLRepo) BindSNByTokenBucket(sn, token, bucket string) (affected int64, err error) {
	db := repo.db.Model(&model.DHCPTokenBucket{}).Where("token = ? AND bucket = ? AND (sn IS NULL OR sn = '')", token, bucket).Update("sn", sn)
	if err = db.Error; err != nil {
		repo.log.Error(err)
		return 0, err
	}
	return db.RowsAffected, nil

}

// UnbindSNByTokenBucket 解绑SN与令牌桶中令牌
func (repo *MySQLRepo) UnbindSNByTokenBucket(sn, token, bucket string) (affected int64, err error) {
	db := repo.db.Model(&model.DHCPTokenBucket{}).Where("token = ? AND bucket = ? AND sn = ?", token, bucket, sn).Update("sn", nil)
	if err = db.Error; err != nil {
		repo.log.Error(err)
		return 0, err
	}
	return db.RowsAffected, nil
}

// GetTokenBuckets 返回满足过滤条件的令牌桶及其令牌
func (repo *MySQLRepo) GetTokenBuckets(cond *model.DHCPTokenBucket) (items []*model.DHCPTokenBucket, err error) {
	db := repo.db.Model(&model.DHCPTokenBucket{})

	if cond != nil {
		if cond.SN != nil {
			db = db.Where("sn = ?", cond.SN)
		}
		if cond.Token != "" {
			db = db.Where("token = ?", cond.Token)
		}
		if cond.Bucket != "" {
			db = db.Where("bucket = ?", cond.Bucket)
		}
	}

	if err = db.Find(&items).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return items, nil
}

// GetUnbindingTokensByBucket 返回未绑定的SN的令牌列表
func (repo *MySQLRepo) GetUnbindingTokensByBucket(bucket string) (tokens []string, err error) {
	if err = repo.db.Model(&model.DHCPTokenBucket{}).Select("token").Where("sn IS NULL OR sn = ''").Pluck("token", &tokens).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return tokens, nil
}

// GetBuckets 返回当前所有令牌桶
func (repo *MySQLRepo) GetBuckets() (buckets []string, err error) {
	if err = repo.db.Raw("SELECT DISTINCT(bucket) FROM dhcp_token_bucket").Pluck("bucket", &buckets).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return buckets, nil
}

// GetTokenBySN 返回目标设备在令牌桶内的令牌
func (repo *MySQLRepo) GetTokenBySN(sn string) (token string, err error) {
	var row model.DHCPTokenBucket
	if err = repo.db.Where("sn = ?", sn).Find(&row).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return "", err
	}
	return row.Token, nil
}

// RemoveTokenBuckets 移除指定名称的令牌桶及其令牌
func (repo *MySQLRepo) RemoveTokenBuckets(buckets ...string) (affected int64, err error) {
	if len(buckets) <= 0 {
		return 0, nil
	}
	db := repo.db.Unscoped().Delete(model.DHCPTokenBucket{}, "bucket IN (?)", buckets)
	if err = db.Error; err != nil {
		repo.log.Error(err)
		return 0, err
	}
	return db.RowsAffected, nil
}
