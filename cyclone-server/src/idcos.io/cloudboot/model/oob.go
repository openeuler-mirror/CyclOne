package model

import "github.com/jinzhu/gorm"

type (
	// OOBHistory 带外修改历史
	OOBHistory struct {
		gorm.Model
		SN          string `gorm:"column:sn"`
		UsernameOld string `gorm:"column:username_old"`
		UsernameNew string `gorm:"column:username_new"`
		PasswordOld string `gorm:"column:password_old"`
		PasswordNew string `gorm:"column:password_new"`
		Remark      string `gorm:"column:remark"`
		Creator     string `gorm:"column:creator"`
	}
)

// TableName 指定数据库表名
func (OOBHistory) TableName() string {
	return "device_oob_history"
}

// IOOBHistory 带外修改历史接口
type IOOBHistory interface {
	// AddOOBHistory 新增带外修改历史-->改由触发器实现
	// AddOOBHistory(*OOBHistory) (affected int64, err error)
	// GetLastOOBHistoryBySN 查询指定SN最近的一条修改记录，用以找到/确认当前的用户密码
	GetLastOOBHistoryBySN(sn string) (mod OOBHistory, err error)
}

//const (
//	OOBHistoryRemarkManu = "手动修改带外密码"
//	OOBHistoryRemarkAuto = "自动批量初始化带外密码"
//)
