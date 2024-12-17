package mysqlrepo

import "idcos.io/cloudboot/model"

// AddDataDicts 新增数据字典
func (repo *MySQLRepo) AddDataDicts(mods []*model.DataDict) error {
	tx := repo.db.Begin()
	for _, m := range mods {
		err := repo.db.Create(m).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return nil
}

//DelDataDicts 删除数据字典
func (repo *MySQLRepo) DelDataDicts(mods []*model.DataDict) error {
	tx := repo.db.Begin()
	for _, m := range mods {
		err := repo.db.Delete(m).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return nil
}

//UpdateDataDicts 修改数据字典
func (repo *MySQLRepo) UpdateDataDicts(mods []*model.DataDict) error {
	tx := repo.db.Begin()
	for _, m := range mods {
		err := repo.db.Model(&model.DataDict{}).Update(m).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return nil
}

// GetDataDictsByType 查询指定类型的数据字典列表
func (repo *MySQLRepo) GetDataDictsByType(typ string) (mods []*model.DataDict, err error) {
	if typ != "" {
		err = repo.db.Where("type = ?", typ).Find(&mods).Error
		if err != nil {
			repo.log.Errorf("查询数据字典异常， %s", err.Error())
			return
		}
		return
	}
	repo.db.Find(&mods)
	return
}
