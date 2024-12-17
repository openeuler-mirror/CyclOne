package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"
	"sort"

	"github.com/voidint/binding"
	"github.com/voidint/page"

	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils/times"
)

// InspectionStatistics 硬件巡检结果统计
type InspectionStatistics struct {
	//日期
	Date string `json:"date"`
	//正常设备数量
	NominalCount int `json:"nominal_count"`
	//警告设备数量
	WarningCount int `json:"warning_count"`
	//异常设备数量
	CriticalCount int `json:"critical_count"`
}

// GetGetInspectionStatisticsReq 查询硬件巡检结果状态(正常、警告、异常)统计请求结构体
type GetGetInspectionStatisticsReq struct {
	Period    string `json:"period"`
	Direction string `json:"direction"`
}

// FieldMap 字段映射
func (reqData *GetGetInspectionStatisticsReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑
	return binding.FieldMap{
		&reqData.Period:    "period",
		&reqData.Direction: "direction",
	}
}

// Validate 结构体数据校验
func (reqData *GetGetInspectionStatisticsReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {

	if reqData.Period != "" && reqData.Period != StatiPeriodLatestWeek && reqData.Period != StatiPeriodLatestMonth {
		errs.Add([]string{"period"}, binding.BusinessError, fmt.Sprintf("基本参数错误:%s=%s", "period", reqData.Period))
		return errs
	}
	if reqData.Direction != "" && reqData.Direction != DirectionASC && reqData.Direction != DirectionDESC {
		errs.Add([]string{"direction"}, binding.BusinessError, fmt.Sprintf("基本参数错误:%s=%s", "direction", reqData.Direction))
		return errs
	}
	return errs
}

const (
	// StatiPeriodLatestWeek 统计周期-近一周
	StatiPeriodLatestWeek = "latest_week"
	// StatiPeriodLatestMonth 统计周期-近一月
	StatiPeriodLatestMonth = "latest_month"
)

const (
	// DirectionASC 排序方向-递增
	DirectionASC = "asc"
	// DirectionDESC 排序方向-递减
	DirectionDESC = "desc"
)

const (
	dateLayout = "2006-01-02"
	selTimeLayout = "Jan-02-2006 15:04:05"

)

// GetGetInspectionStatistics 查询硬件巡检结果状态(正常、警告、异常)统计
func GetGetInspectionStatistics(log logger.Logger, repo model.Repo, reqData *GetGetInspectionStatisticsReq) (items []*InspectionStatistics, err error) {
	// 1、生成日期列表
	if reqData.Period == "" {
		reqData.Period = StatiPeriodLatestWeek
	}
	isASC := true // 默认升序排列
	if reqData.Direction == DirectionDESC {
		isASC = false
	}

	var dates []time.Time
	switch reqData.Period {
	case StatiPeriodLatestWeek:
		dates = times.LatestWeek(time.Now(), isASC)
	case StatiPeriodLatestMonth:
		dates = times.LatestMonth(time.Now(), isASC)
	}

	items = make([]*InspectionStatistics, 0, len(dates))

	// 2、遍历日期数组并逐天统计硬件正常、警告、异常设备数量
	for i := range dates {
		// 查询指定日期内已完成的设备巡检记录
		insps, err := repo.GetInspectionStatisticsGroupBySN(&model.Inspection{
			StartTime:     &dates[i],
			RunningStatus: model.RunningStatusDone,
		}, model.OneOrderBy("start_time", model.ASC), nil)
		if err != nil {
			return nil, err
		}
		// 统计每台设备当天最新一次巡检的结果
		var critical, warning, nominal int
		for i := range insps {
			if insps[i] == nil || insps[i].HealthStatus == model.HealthStatusUnknown {
				continue
			}
			switch insps[i].HealthStatus {
			case model.HealthStatusCritical:
				critical++

			case model.HealthStatusWarning:
				warning++

			case model.HealthStatusNominal:
				nominal++
			}			
		}

		items = append(items, &InspectionStatistics{
			Date:          dates[i].Format(dateLayout),
			NominalCount:  nominal,
			WarningCount:  warning,
			CriticalCount: critical,
		})
	}
	return items, nil
}

// GetInspectionsPageReq 查询巡检分页请求结构体
type GetInspectionsPageReq struct {
	// 关键词
	Keyword string `json:"keyword"`
	// 设备固资号
	FixedAssetNumber string `json:"fixed_asset_number"`
	// 设备序列号，支持英文逗号分隔多个序列号
	SN string `json:"sn"`
    // 内网IP
	IntranetIP  string  `json:"intranet_ip"`	
	// 巡检开始时间
	StartTime string `json:"start_time"`
	// 巡检结束时间
	EndTime string `json:"end_time"`
	// 带外IP
	//OOBIP string `json:"oob_ip"`
	// 巡检执行状态
	RuningStatus string `json:"running_status"`
	// 机器的巡检结果（包含正常，告警，致命，未知）
	HealthStatus string `json:"health_status"`
	// 分页页号
	Page int64 `json:"page"`
	// 分页大小
	PageSize int64 `json:"page_size"`
}

// FieldMap 字段映射
func (reqData *GetInspectionsPageReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.FixedAssetNumber:        "fixed_asset_number",
		&reqData.IntranetIP:        "intranet_ip",
		&reqData.SN:        "sn",
		&reqData.StartTime: "start_time",
		&reqData.EndTime:   "end_time",
		//&reqData.OOBIP:        "oob_ip",
		&reqData.RuningStatus: "running_status",
		&reqData.HealthStatus: "health_status",
		&reqData.Page:         "page",
		&reqData.PageSize:     "page_size",
	}
}

// GetInspectionsPage 返回满足条件的物理机设备巡检结果的分页列表
func GetInspectionsPage(log logger.Logger, repo model.Repo, reqData *GetInspectionsPageReq) (pg *page.Page, err error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 100 {
		reqData.PageSize = 20
	}
	if reqData.Page < 0 {
		reqData.Page = 0
	}

	cond := model.InspectionCond{
		FixedAssetNumber: reqData.FixedAssetNumber,
		SN:               reqData.SN,
		IntranetIP:       reqData.IntranetIP,
		StartTime:        reqData.StartTime,
		EndTime:          reqData.EndTime,
		//OOBIP:          reqData.OOBIP,
		HealthStatus:     reqData.HealthStatus,
		RuningStatus:     reqData.RuningStatus,
	}

	totalRecords, err := repo.CountInspctionsByCond(&cond)
	if err != nil {
		return nil, err
	}

	pager := page.NewPager(reflect.TypeOf(&InspectionFullWithPage{}), reqData.Page, reqData.PageSize, totalRecords)
	items, err := repo.GetInspectionListWithPageNew(&cond, pager.BuildLimiter())
	if err != nil {
		return nil, err
	}
	for i := range items {
		item := convert2IPMIResult(log, &items[i])
		pager.AddRecords(item)
	}
	return pager.BuildPage(), nil
}

/*
temperature: 温度。可选值:nominal, warning, critical、unknown
voltage: 电压。可选值:nominal, warning, critical、unknown
fan: 风扇。可选值:nominal, warning, critical、unknown
memory: 内存。可选值:nominal, warning, critical、unknown
power_supply: 电源。可选值:nominal, warning, critical、unknown
*/

//InspectionResultItem 单项巡检结果（包括temperature，voltage，fan，memory，power_supply）
type InspectionResultItem struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}

const (
	temperature = iota
	voltage
	fan
	memory
	powerSupply
)

//InspectionFullWithPage 硬件巡检分页查询
type InspectionFullWithPage struct {
	//索引
	ID uint `json:"id"`
	// 设备固资号
	FixedAssetNumber string `json:"fixed_asset_number"`
	//设备序列号
	SN string `json:"sn"`
    // 内网IP
	IntranetIP            string  `json:"intranet_ip"`
	//巡检开始时间
	StartTime times.ISOTime `json:"start_time"`
	//巡检结束时间
	EndTime times.ISOTime `json:"end_time"`
	//运行状态
	RuningStatus string `json:"running_status"`
	//带外IP
	//OOBIP string `json:"oob_ip"`
	//巡检错误信息
	Error string `json:"error"`
	//健康状态
	HealthStatus string `json:"health_status"`
	//巡检结果
	Result []InspectionResultItem `json:"result"`
	//巡检日志创建时间
	CreatedAt times.ISOTime `json:"created_at"`
	//巡检日志更新时间
	UpdatedAt times.ISOTime `json:"updated_at"`
}

//convert2IPMIResult 对IPMIResult结果做进一步处理， 有数据之后，这块结果放到repo做掉
func convert2IPMIResult(log logger.Logger, item *model.InspectionFullWithPage) *InspectionFullWithPage {
	var result []*model.SensorData
	var ifwp InspectionFullWithPage
	ifwp.ID = item.ID
	ifwp.FixedAssetNumber = item.FixedAssetNumber
	ifwp.SN = item.SN
	ifwp.IntranetIP = item.IntranetIP
	ifwp.StartTime = item.StartTime
	ifwp.EndTime = item.EndTime
	ifwp.RuningStatus = item.RuningStatus
	ifwp.HealthStatus = item.HealthStatus
	//ifwp.OOBIP = item.OOBIP
	ifwp.Error = item.Error
	ifwp.CreatedAt = item.CreatedAt
	ifwp.UpdatedAt = item.UpdatedAt

	if ifwp.Error != "" {
		ifwp.Result = []InspectionResultItem{}
		return &ifwp
	}

	if err := json.Unmarshal([]byte(item.Result), &result); err != nil {
		log.Info(err)
	}

	isr := []InspectionResultItem{
		InspectionResultItem{
			Type:   "temperature",
			Status: "unknown",
		},
		InspectionResultItem{
			Type:   "voltage",
			Status: "unknown",
		},
		InspectionResultItem{
			Type:   "fan",
			Status: "unknown",
		},
		InspectionResultItem{
			Type:   "memory",
			Status: "unknown",
		},
		InspectionResultItem{
			Type:   "power_supply",
			Status: "unknown",
		},
	}
	//"nominal", "warning", "critical", "unknown"
	//    1    ,    2     ,     3     ,    0
	status := map[string]int{"unknown": 0, "nominal": 1, "warning": 2, "critical": 3}
	for _, v := range result {
		if strings.Contains(v.Type, "Temperature") {
			if isr[temperature].Status == "critical" {
				continue
			}

			if status[strings.ToLower(v.State)] > status[isr[temperature].Status] {
				isr[temperature].Status = strings.ToLower(v.State)
			}
		}
		if strings.Contains(v.Type, "Voltage") {
			if isr[voltage].Status == "critical" {
				continue
			}

			if status[strings.ToLower(v.State)] > status[isr[voltage].Status] {
				isr[voltage].Status = strings.ToLower(v.State)
			}
		}
		if strings.Contains(v.Type, "Fan") {
			if isr[fan].Status == "critical" {
				continue
			}

			if status[strings.ToLower(v.State)] > status[isr[fan].Status] {
				isr[fan].Status = strings.ToLower(v.State)
			}
		}
		if strings.Contains(v.Type, "Memory") {
			if isr[memory].Status == "critical" {
				continue
			}

			if status[strings.ToLower(v.State)] > status[isr[memory].Status] {
				isr[memory].Status = strings.ToLower(v.State)
			}
		}
		if strings.Contains(v.Type, "Power Supply") {
			if isr[powerSupply].Status == "critical" {
				continue
			}

			if status[strings.ToLower(v.State)] > status[isr[powerSupply].Status] {
				isr[powerSupply].Status = strings.ToLower(v.State)
			}
		}
	}
	ifwp.Result = isr
	return &ifwp
}

// GetInspectionBySNReq 查询巡检明细请求结构体
type GetInspectionBySNReq struct {
	// 硬件巡检目标设备
	SN string `json:"sn"`
	// 硬件巡检开始时间
	StartTime string `json:"start_time"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *GetInspectionBySNReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.SN:        "sn",
		&reqData.StartTime: "start_time",
	}
}

// GetInspectionBySN 查询指定设备某一次硬件巡检明细
func GetInspectionBySN(repo model.Repo, req *GetInspectionBySNReq) (map[string]interface{}, error) {
	mod, err := repo.GetInspectionDetail(req.SN, req.StartTime)
	if err != nil {
		return nil, err
	}
	var ipmiSensor []model.SensorData
	var ipmiSel []model.SelData
	_ = json.Unmarshal([]byte(mod.IPMIResult), &ipmiSensor)
    _ = json.Unmarshal([]byte(mod.IPMISELResult), &ipmiSel)
	
	return map[string]interface{}{
		"id":             mod.ID,
		"job_id":         mod.JobID,
		"start_time":     times.ISOTime(*mod.StartTime),
		"end_time":       times.ISOTime(*mod.EndTime),
		"running_status": mod.RunningStatus,
		"error":          mod.Error,
		"result":         sensorDataFilter(ipmiSensor, model.HealthStatusNominal),
		"selresult":      selDataFilter(ipmiSel, model.HealthStatusNominal),
		"created_at":     times.ISOTime(mod.CreatedAt),
		"updated_at":     times.ISOTime(mod.UpdatedAt),
	}, nil
}

func sensorDataFilter(ipmiSensor []model.SensorData, state string) []model.SensorData {
	var filtered []model.SensorData

	for _, it := range ipmiSensor {
		if strings.ToLower(it.State) == state {
			continue
		}

		filtered = append(filtered, it)
	}

	return filtered
}

// 实现系统事件日志数据按时间最新排序
type SelDataCollect []model.SelData
func (sdc SelDataCollect) Len() int {
	return len(sdc)
}

// 时间比较
func (sdc SelDataCollect) Less(i, j int) bool {
    // 将字符串的日期时间转换time.Time 再通过Before or After 进行比较
	formatTimeI, _ := time.Parse(selTimeLayout, fmt.Sprintf("%s %s", sdc[i].Date, sdc[i].Time))
	formatTimeJ, _ := time.Parse(selTimeLayout, fmt.Sprintf("%s %s", sdc[j].Date, sdc[j].Time))

    if formatTimeI.After(formatTimeJ) {
		return true
	}
	return false
}

func (sdc SelDataCollect) Swap(i, j int) {
	sdc[i], sdc[j] = sdc[j], sdc[i]
}

// 过滤Nominal的日志并重新排序
func selDataFilter(ipmiSel []model.SelData, state string) []model.SelData {
	var filtered SelDataCollect
	for _, it := range ipmiSel {
		if strings.ToLower(it.State) == state {
			continue
		}

		filtered = append(filtered, it)
	}
	sort.Sort(filtered)
	return filtered 
}

// GetInspectionStartTimesBySNReq 查询巡检明细请求结构体
type GetInspectionStartTimesBySNReq struct {
	// 硬件巡检目标设备
	SN string
	// 硬件巡检执行状态
	RunningStatus string
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *GetInspectionStartTimesBySNReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.SN:            "sn",
		&reqData.RunningStatus: "running_status",
	}
}

// GetInspectionStartTimesBySN 查询指定设备的历史硬件巡检开始时间列表
func GetInspectionStartTimesBySN(repo model.Repo, req *GetInspectionStartTimesBySNReq) ([]string, error) {
	mod, err := repo.GetInspectionStartTimeBySN(req.SN, req.RunningStatus)
	if err != nil {
		return nil, err
	}
	var starts []string
	for _, m := range mod {
		starts = append(starts, times.ISOTime(m).String())
	}
	return starts, nil
}

// InspectionPageReq 物理机巡检记录查询分页请求
type InspectionRecordsPageReq struct {
	// 设备序列号，支持英文逗号分隔多个序列号
	SN string `json:"sn"`
	// 巡检开始时间
	StartTime string `json:"start_time"`
	// 巡检结束时间
	EndTime string `json:"end_time"`
	// 机器的巡检结果（包含正常，告警，致命，未知） enum('nominal','warning','critical','unknown') 
	HealthStatus string `json:"health_status"`
	// 分页页号
	Page int64 `json:"page"`
	// 分页大小 限制1000
	PageSize int64 `json:"page_size"`	
}

// FieldMap 字段映射
func (reqData *InspectionRecordsPageReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.SN:        "sn",
		&reqData.StartTime: "start_time",
		&reqData.EndTime:   "end_time",
		&reqData.HealthStatus: "health_status",
		&reqData.Page:         "page",
		&reqData.PageSize:     "page_size",
	}
}

//InspectionRecordsPage 物理机巡检记录查询分页结果
type InspectionRecordsPage struct {
	//索引
	ID uint `json:"id"`
	//设备序列号
	SN string `json:"sn"`
	//巡检开始时间
	StartTime times.ISOTime `json:"start_time"`
	//巡检结束时间
	EndTime times.ISOTime `json:"end_time"`
	//运行状态
	RuningStatus string `json:"running_status"`
	//巡检错误信息
	Error string `json:"error"`
	//健康状态
	HealthStatus string `json:"health_status"`
	//巡检结果
	Result []model.SensorData `json:"result"`
	SelResult []model.SelData `json:"selresult"`
	//巡检日志创建时间
	CreatedAt times.ISOTime `json:"created_at"`
	//巡检日志更新时间
	UpdatedAt times.ISOTime `json:"updated_at"`
}

// GetInspectionRecordsPage 查询物理机巡检记录分页列表
func GetInspectionRecordsPage(log logger.Logger, repo model.Repo, reqData *InspectionRecordsPageReq) (pg *page.Page, err error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 1000 {
		reqData.PageSize = 20
	}
	if reqData.Page < 0 {
		reqData.Page = 1
	}

	cond := model.InspectionCond{
		SN:        reqData.SN,
		StartTime: reqData.StartTime,
		EndTime:   reqData.EndTime,
		HealthStatus: reqData.HealthStatus,
	}

	totalRecords, err := repo.CountInspectionRecordsPage(&cond)
	if err != nil {
		return nil, err
	}
    
	pager := page.NewPager(reflect.TypeOf(&InspectionRecordsPage{}), reqData.Page, reqData.PageSize, totalRecords)
	items, err := repo.GetInspectionRecordsPage(&cond, pager.BuildLimiter())
	if err != nil {
		return nil, err
	}
	for i := range items {
		item := FilterIPMIResult(log, items[i])
		pager.AddRecords(item)
	}
	return pager.BuildPage(), nil
}

//convert2IPMIResult 对IPMIResult结果做进一步处理， 有数据之后，这块结果放到repo做掉
func FilterIPMIResult(log logger.Logger, item *model.InspectionRecordsPage) *InspectionRecordsPage {
	var result []*model.SensorData
	var selresult []*model.SelData
	var IRP InspectionRecordsPage
	IRP.ID = item.ID
	IRP.SN = item.SN
	IRP.StartTime = item.StartTime
	IRP.EndTime = item.EndTime
	IRP.RuningStatus = item.RuningStatus
	IRP.HealthStatus = item.HealthStatus
	IRP.Error = item.Error
	IRP.CreatedAt = item.CreatedAt
	IRP.UpdatedAt = item.UpdatedAt

	// 过滤传感器中读数N/A的数据
	var filteredReading = "N/A"
	var filteredresult []model.SensorData
	if err := json.Unmarshal([]byte(item.Result), &result); err != nil {
		log.Info(err)
	}
	for _, v := range result {
		if strings.Contains(v.Reading, filteredReading) {
			continue
		}
		filteredresult = append(filteredresult, *v)
	}
	IRP.Result = filteredresult

	// 系统事件日志暂时不过滤
	var filteredselresult []model.SelData
	if err := json.Unmarshal([]byte(item.SelResult), &selresult); err != nil {
		log.Info(err)
	}
	for _, v := range selresult {
		filteredselresult = append(filteredselresult, *v)
	}
	IRP.SelResult = filteredselresult
	return &IRP
}