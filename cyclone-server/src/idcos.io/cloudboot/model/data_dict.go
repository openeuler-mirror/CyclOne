package model

import "github.com/jinzhu/gorm"

//DataDict 数据字典
type DataDict struct {
	gorm.Model
	Type   string
	Name   string
	Value  string
	Remark string
}

//TableName 表名称
func (DataDict) TableName() string {
	return "data_dict"
}

//IDataDict 数据字典接口
type IDataDict interface {
	AddDataDicts(items []*DataDict) error
	DelDataDicts(items []*DataDict) error
	UpdateDataDicts([]*DataDict) error
	GetDataDictsByType(typ string) ([]*DataDict, error)
}
