package service

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/voidint/binding"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
)

//DataDict 数据字典
type DataDict struct {
	// 主键
	ID uint `json:"id"`
	// 类型
	Type string `json:"type"`
	// 名称
	Name string `json:"name"`
	// 值
	Value string `json:"value"`
	// 备注
	Remark string `json:"remark"`
}

// DelDataDictReq 数据字典请求体
type DelDataDictReq struct {
	ID uint `json:"id"`
}

// FieldMap 请求字段映射
func (reqData *DataDict) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.ID:     "id",
		&reqData.Type:   "type",
		&reqData.Name:   "name",
		&reqData.Value:  "value",
		&reqData.Remark: "remark",
	}
}

// GetDataDict 根据type参数查询筛选数据字典信息
func GetDataDict(log logger.Logger, repo model.Repo, typ string) ([]*DataDict, error) {
	dataDicts, err := repo.GetDataDictsByType(typ)
	if err != nil {
		log.Error(err.Error())
	}
	dataDictResp := make([]*DataDict, 0, len(dataDicts))
	for _, d := range dataDicts {
		dataDictResp = append(dataDictResp, &DataDict{
			ID:     d.ID,
			Type:   d.Type,
			Name:   d.Name,
			Value:  d.Value,
			Remark: d.Remark,
		})
	}
	return dataDictResp, nil
}

//AddDataDicts 新增字典
func AddDataDicts(log logger.Logger, repo model.Repo, dataDicts []*DataDict) error {
	mods := make([]*model.DataDict, 0, len(dataDicts))
	for _, d := range dataDicts {
		mods = append(mods, &model.DataDict{
			Type:   d.Type,
			Name:   d.Name,
			Value:  d.Value,
			Remark: d.Remark,
		})
	}
	return repo.AddDataDicts(mods)
}

//DelDataDicts 删除字典
func DelDataDicts(log logger.Logger, repo model.Repo, dataDicts []*DelDataDictReq) error {
	mods := make([]*model.DataDict, 0, len(dataDicts))
	for _, d := range dataDicts {
		mods = append(mods, &model.DataDict{
			Model: gorm.Model{ID: d.ID},
		})
	}
	return repo.DelDataDicts(mods)
}

//UpdateDataDicts 修改字典
func UpdateDataDicts(log logger.Logger, repo model.Repo, dataDicts []*DataDict) error {
	mods := make([]*model.DataDict, 0, len(dataDicts))
	for _, d := range dataDicts {
		mods = append(mods, &model.DataDict{
			Model:  gorm.Model{ID: d.ID},
			Type:   d.Type,
			Name:   d.Name,
			Value:  d.Value,
			Remark: d.Remark,
		})
	}
	return repo.UpdateDataDicts(mods)
}
