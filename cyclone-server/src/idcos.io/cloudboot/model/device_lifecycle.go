package model

import (
	"time"
	"github.com/jinzhu/gorm"
)

const (
	// 维保状态枚举值:在保-under_warranty;过保-out_of_warranty;未激活-inactive'
	MaintenanceServiceStatusUnderWarranty = "under_warranty"
	MaintenanceServiceStatusOutOfWarranty = "out_of_warranty"
	MaintenanceServiceStatusInactive = "inactive"

	// ChangeLog 变更类型枚举值
	OperationTypeAdd = "设备新增"
	OperationTypeUpdate = "属性更新"
	OperationTypePowerControl = "电源控制"
	OperationTypeOSInstall = "系统部署"
	OperationTypeMove = "设备搬迁"
	OperationTypeRetire = "设备退役"

	// 是否租赁
	YESRental = "yes"
	NORental = "no"
)

// 固资号	WDEV+YYMM+0001，如WDEV21100001代表21年10月录入的设备
// 序列号	设备SN
// 资产归属	腾讯云、华为云
// 负责人	主要应对特殊设备
// 是否租赁	
// 维保服务供应商	合同乙方，huawei\腾讯云l等，或者第三方采购的服务供应商，如神州数码
// 物流服务内容	腾讯云接口人统一安排、华为接口人统一安排
// 维保服务内容	"银牌服务：
// 7x24小时，第二个工作日上门，介质保留服务"
// 维保报障方式	"huawei：
// 800报修，腾讯大客户身份
// lenovo:
// 800报修，腾讯大客户身份"
// 维保起始日期	上架导入日期+设备缓冲期
// 维保截止日期	起始日期+合同维保期
// 维保状态	在保-under_warranty;过保-out_of_warranty;未激活-inactive
// 设备退役日期
// 变更记录	"[
//   {
//     ""操作用户"":""jessewei""，
//     ""操作类型"":""新增"",
//     ""操作内容"":“XXX”,
//     ""操作时间"":""2021-10-14 15:11:16"",
//   },
//   {
//     ""操作用户"":""jessewei""，
//     ""操作类型"":""维修"",
//     ""操作内容"":“XXX”,
//     ""操作时间"":""2021-10-14 15:12:14"",
//   },
// ]"

// LifecycLog JSON结构体
type ChangeLog struct {
	OperationUser		string		`json:"operation_user"`
	OperationType		string		`json:"operation_type"`
	OperationDetail		string		`json:"operation_detail"`
	OperationTime		string		`json:"operation_time"`
}


// DeviceLifecycle 设备生命周期表（维保、资产、变更记录等详细信息）
type DeviceLifecycle struct {
	gorm.Model
	FixedAssetNumber 				string    	`gorm:"column:fixed_asset_number"`
	SN               				string    	`gorm:"column:sn"`
	AssetBelongs	 				string		`gorm:"column:asset_belongs"`
	Owner			 				string		`gorm:"column:owner"`
	IsRental		 				string		`gorm:"column:is_rental"`
	MaintenanceServiceProvider		string		`gorm:"column:maintenance_service_provider"`
	MaintenanceService				string		`gorm:"column:maintenance_service"`
	LogisticsService				string		`gorm:"column:logistics_service"`
	MaintenanceServiceDateBegin     time.Time 	`gorm:"column:maintenance_service_date_begin"`
	MaintenanceServiceDateEnd       time.Time 	`gorm:"column:maintenance_service_date_end"`
	MaintenanceServiceStatus		string		`gorm:"column:maintenance_service_status"`
	DeviceRetiredDate       		time.Time 	`gorm:"column:device_retired_date"`
	LifecycleLog					string		`gorm:"column:lifecycle_log"`
}

// DeviceLifecycleDetail for CombinedDevice 设备信息及设备装机参数联合结构体
type DeviceLifecycleDeatail struct {
	AssetBelongs	 				string		`gorm:"column:asset_belongs"`
	Owner			 				string		`gorm:"column:owner"`
	IsRental		 				string		`gorm:"column:is_rental"`
	MaintenanceServiceProvider		string		`gorm:"column:maintenance_service_provider"`
	MaintenanceService				string		`gorm:"column:maintenance_service"`
	LogisticsService				string		`gorm:"column:logistics_service"`
	MaintenanceServiceDateBegin     time.Time 	`gorm:"column:maintenance_service_date_begin"`
	MaintenanceServiceDateEnd       time.Time 	`gorm:"column:maintenance_service_date_end"`
	MaintenanceServiceStatus		string		`gorm:"column:maintenance_service_status"`
	DeviceRetiredDate       		time.Time 	`gorm:"column:device_retired_date"`
	LifecycleLog					string		`gorm:"column:lifecycle_log"`
}


// TableName 指定数据库表名
func (DeviceLifecycle) TableName() string {
	return "device_lifecycle"
}

// IDeviceLifecycle 持久化接口
type IDeviceLifecycle interface {
	SaveDeviceLifecycle(*DeviceLifecycle) (affected int64, err error)
	RemoveDeviceLifecycleByID(id uint) (affected int64, err error)
	RemoveDeviceLifecycleBySN(sn string) (affected int64, err error)
	GetDeviceLifecycleBySN(sn string) (*DeviceLifecycle, error)
	UpdateDeviceLifecycleBySN(*DeviceLifecycle) (err error)
}