package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/voidint/page"

	"idcos.io/cloudboot/utils/times"
)

const (
	// RunningStatusRunning 运行状态-运行中
	RunningStatusRunning = "running"
	// RunningStatusDone 运行状态-完成
	RunningStatusDone = "done"
)

const (
	// HealthStatusNominal 设备健康状况-正常
	HealthStatusNominal = "nominal"
	// HealthStatusWarning 设备健康状况-警告
	HealthStatusWarning = "warning"
	// HealthStatusCritical 设备健康状况-异常
	HealthStatusCritical = "critical"
	// HealthStatusUnknown 设备健康状况-未知
	HealthStatusUnknown = "unknown"
)

// SensorData IPMI传感器数据
type SensorData struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	State   string `json:"state"`
	Reading string `json:"reading"`
	Units   string `json:"units"`
	Event   string `json:"event"`
}

// SelData IPMI系统事件数据 ID,Date,Time,Name,Type,State,Event
type SelData struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Time    string `json:"time"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	State   string `json:"state"`
	Event   string `json:"event"`
}

// Inspection 硬件巡检
type Inspection struct {
	gorm.Model
	JobID         string     `gorm:"column:job_id"`
	StartTime     *time.Time `gorm:"column:start_time"`
	EndTime       *time.Time `gorm:"column:end_time"`
	OriginNode    string     `gorm:"column:origin_node"`
	SN            string     `gorm:"column:sn"`
	RunningStatus string     `gorm:"column:running_status"`
	HealthStatus  string     `gorm:"column:health_status"`
	Error         string     `gorm:"column:error"`
	IPMIResult    string     `gorm:"column:ipmi_result"` // 通过IPMI获取的硬件巡检结果。JSON结构同[]*SensorData
	IPMISELResult string     `gorm:"column:ipmisel_result"` // 通过ipmi-sel获取的硬件系统时间日志。JSON结构同[]*SelData
}

// TableName 指定数据库表名
func (Inspection) TableName() string {
	return "inspection"
}

//InspectionCond 分页查询硬件巡检结果条件
type InspectionCond struct {
	FixedAssetNumber  string
	SN                string
	IntranetIP        string 
	StartTime         string
	EndTime           string
	OOBIP             string
	RuningStatus      string
	HealthStatus      string
	Order             string
	KeyWord           string
}

//InspectionFullWithPage 分页获取硬件巡检结果，每一个SN仅仅展示最新的一条记录
type InspectionFullWithPage struct {
	FixedAssetNumber string    `json:"fixed_asset_number" gorm:"column:fixed_asset_number"`
	ID           uint          `json:"id" gorm:"column:id"`
	SN           string        `json:"sn" gorm:"column:sn"`
	IntranetIP   string        `json:"intranet_ip" gorm:"column:intranet_ip"`
	StartTime    times.ISOTime `json:"start_time" gorm:"column:start_time"`
	EndTime      times.ISOTime `json:"end_time" gorm:"column:end_time"`
	RuningStatus string        `json:"running_status" gorm:"column:running_status"`
	HealthStatus string        `json:"health_status" gorm:"column:health_status"` // 机器的巡检结果（包含正常，告警，致命，未知）
	//OOBIP        string        `json:"oob_ip" gorm:"column:oob_ip"`
	Error     string        `json:"error" gorm:"column:error"`
	Result    string        `json:"result" gorm:"column:ipmi_result"`
	CreatedAt times.ISOTime `json:"created_at" gorm:"column:created_at"`
	UpdatedAt times.ISOTime `json:"updated_at" gorm:"column:updated_at"`
}

//InspectionRecordsPage 物理机巡检记录查询分页结果
type InspectionRecordsPage struct {
	ID           uint          `json:"id" gorm:"column:id"`
	SN           string        `json:"sn" gorm:"column:sn"`
	StartTime    times.ISOTime `json:"start_time" gorm:"column:start_time"`
	EndTime      times.ISOTime `json:"end_time" gorm:"column:end_time"`
	RuningStatus string        `json:"running_status" gorm:"column:running_status"`
	HealthStatus string        `json:"health_status" gorm:"column:health_status"` // 机器的巡检结果（包含正常，告警，致命，未知）
	Error     string        `json:"error" gorm:"column:error"`
	Result    string        `json:"result" gorm:"column:ipmi_result"`
	SelResult    string        `json:"selresult" gorm:"column:ipmisel_result"`
	CreatedAt times.ISOTime `json:"created_at" gorm:"column:created_at"`
	UpdatedAt times.ISOTime `json:"updated_at" gorm:"column:updated_at"`
}

// TableName 指定数据库表名
func (InspectionRecordsPage) TableName() string {
	return "inspection"
}

// Inspection 硬件巡检记录健康状态统计
type InspectionStatistics struct {
	gorm.Model
	JobID         string     `gorm:"column:job_id"`
	StartTime     *time.Time `gorm:"column:start_time"`
	EndTime       *time.Time `gorm:"column:end_time"`
	SN            string     `gorm:"column:sn"`
	RunningStatus string     `gorm:"column:running_status"`
	HealthStatus  string     `gorm:"column:health_status"`
}

// TableName 指定数据库表名
func (InspectionStatistics) TableName() string {
	return "inspection"
}

// IInspection 硬件巡检操作接口
type IInspection interface {
	// AddInspections 批量新增硬件巡检记录
	AddInspections(...*Inspection) error
	// UpdateInspectionByID 根据ID更新硬件巡检记录
	UpdateInspectionByID(*Inspection) (affected int64, err error)
	// GetInspectedSN 查询已经执行过硬件巡检的设备SN列表
	GetInspectedSN() (items []string, err error)
	//CountInspctionsByCond 根据条件统计
	CountInspctionsByCond(cond *InspectionCond) (count int64, err error)
	//GetInspectionListWithPageNew 分页查询硬件巡检结果
	GetInspectionListWithPageNew(cond *InspectionCond, limiter *page.Limiter) ([]InspectionFullWithPage, error)
	// GetInspectionDetail 查询指定设备的某次巡检详情
	GetInspectionDetail(SN, startTime string) (i *Inspection, err error)
	//GetInspectionStartTimeBySN 查询指定设备的历史硬件巡检开始时间列表
	GetInspectionStartTimeBySN(SN, Where string) ([]time.Time, error)
	// GetInspections 查询满足过滤条件的硬件巡检列表
	GetInspections(cond *Inspection, orderby OrderBy, limiter *page.Limiter) (items []*Inspection, err error)
	// GetInspectionStatistics 查询满足过滤条件的硬件巡检列表用于健康状态统计
	GetInspectionStatistics(cond *Inspection, orderby OrderBy, limiter *page.Limiter) (items []*InspectionStatistics, err error)
	// GetInspectionStatisticsGroupBySN 查询满足过滤条件的硬件巡检列表用于健康状态统计
	GetInspectionStatisticsGroupBySN(cond *Inspection, orderby OrderBy, limiter *page.Limiter) (items []*InspectionStatistics, err error)	
	//RemoveInspectionOnStartTimeBySN 根据设备SN删除巡检记录，按时间排序保留一定数量的记录
	RemoveInspectionOnStartTimeBySN(SN string) (err error)
	// CountInspectionPageRecords 获取巡检记录条数
	CountInspectionRecordsPage(cond *InspectionCond) (count int64, err error)
	// GetInspectionPage 获取巡检记录分页
	GetInspectionRecordsPage(cond *InspectionCond, limiter *page.Limiter) (items []*InspectionRecordsPage, err error)
}

const (
	// EmptyJSONObject 空JSON对象字符串标识
	EmptyJSONObject = "{}"
	// EmptyJSONArray 空JSON数组字符串标识
	EmptyJSONArray = "[]"
)

func replaceIfBlank(field *string, newS string) {
	if *field == "" {
		*field = newS
	}
}

// BeforeSave 保存硬件巡检对象前的钩子方法。
// 防止将空字符串写入类型为JSON的数据库字段中引发报错。
func (insp *Inspection) BeforeSave() (err error) {
	replaceIfBlank(&insp.IPMIResult, EmptyJSONArray)
	replaceIfBlank(&insp.IPMISELResult, EmptyJSONArray)
	return
}
