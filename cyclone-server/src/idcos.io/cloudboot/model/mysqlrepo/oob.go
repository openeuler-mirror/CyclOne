package mysqlrepo

import "idcos.io/cloudboot/model"

// AddOOBHistory 增加带外修改历史记录
//func (repo *MySQLRepo) AddOOBHistory(mod *model.OOBHistory) (affected int64, err error) {
//	err = repo.db.Create(mod).Error
//	if err != nil {
//		return 0, err
//	}
//	return 1, nil
//}

// GetLastOOBHistoryBySN 查询指定SN最近的一条修改记录，用以找到/确认当前的用户密码
func (repo *MySQLRepo) GetLastOOBHistoryBySN(sn string) (mod model.OOBHistory, err error) {
	err = repo.db.Model(model.OOBHistory{}).Where("sn  = ?", sn).Last(&mod).Error
	if err != nil {
		return mod, err
	}
	return mod, nil
}
