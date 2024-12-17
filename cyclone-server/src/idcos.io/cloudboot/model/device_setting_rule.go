package model

import (
	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
)

const (
	// 规则事实 fact 的属性
	AttributeCategory = "category"
	AttributeArch = "arch"
	AttributePhysicalArea = "physical_area"
	AttributeVendor = "vendor"
	AttributeLogicalOperator = "logical_operator"
	AttributeIsFITIEcoProduct = "is_fiti_eco_product"
	// 规则事实 fact 的属性与值的关系
	OperatorIN = "in"
	OperatorContains = "contains"
	OperatorEqual = "equal"
	// 多个事实 fact 之间的关系 OR 优先于 AND ,如：（F1 OR F2） AND (F3 OR F4)
	OperatorOR = "or"
	OperatorAND = "and"
)

const (
	DeviceSettingRuleOS = "os"
	DeviceSettingRuleRaid = "raid"
	DeviceSettingRuleNetwork = "network"
)

// 规则前件结构体
type RuleP struct {
	Attribute    string      `json:"attribute"`  // 规则事实属性category physical_area arch vendor 以及逻辑运算符 logical_operator
	Operator     string      `json:"operator"`   // 属性与值的关系
    Value        []string    `json:"value"`      // 属性的值
}

// DeviceSettingRule 产生式设备装机参数规则库
type DeviceSettingRule struct {
	gorm.Model
	Condition                string    `gorm:"column:condition"`      // 前件
	Action                   string    `gorm:"column:action"`         // 结论
	RuleCategory             string    `gorm:"column:rule_category"`  // 规则分类 enum('os','network','raid')
}

// TableName 指定数据库表名
func (DeviceSettingRule) TableName() string {
	return "device_setting_rule"
}

// IDeviceSettingRule 持久化接口
type IDeviceSettingRule interface {
	// SaveDeviceSettingRule 保存规则记录
	SaveDeviceSettingRule(*DeviceSettingRule) (affected int64, err error)
	// RemoveDeviceSettingRuleByID 删除指定ID的规则记录
	RemoveDeviceSettingRuleByID(id uint) (affected int64, err error)	
	// 根据规则分类查询获取所有规则
	GetDeviceSettingRulesByType(queryType string) (items []*DeviceSettingRule, err error)
	// GetDeviceSettingRuleByID 返回指定ID的规则
	GetDeviceSettingRuleByID(id uint) (*DeviceSettingRule, error)
	CountDeviceSettingRules(cond *DeviceSettingRule) (count int64, err error)
	GetDeviceSettingRules(cond *DeviceSettingRule, orderby OrderBy, limiter *page.Limiter) (items []*DeviceSettingRule, err error)	
}