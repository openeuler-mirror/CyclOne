package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/voidint/binding"
	"github.com/voidint/page"

	"strings"

	"os"

	"errors"

	"strconv"

	"idcos.io/cloudboot/config"
	"idcos.io/cloudboot/hardware/collector"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/server/cloudbootserver/types/device"
	"idcos.io/cloudboot/utils"
	"idcos.io/cloudboot/utils/centos6"
	"idcos.io/cloudboot/utils/oob"
	strings2 "idcos.io/cloudboot/utils/strings"
	"idcos.io/cloudboot/utils/times"
	"idcos.io/cloudboot/utils/upload"
)

// ImportPreviewReq 导入参数结构
type ImportPreviewReq struct {
	FileName string `json:"file_name"`
	Limit    uint   `json:"limit"`
	Offset   uint   `json:"offset"`
	// 用户登录名
	LoginName string `json:"-"`
}

type ImportApprovalReq struct {
	ImportPreviewReq
	//审批人
	Approvers []string `json:"approvers"` // 审批人ID构成的JSON数组字符串
}

const commaSep = ","

//ImportDevicesReq 设备导入Excel表对应字段
type ImportDevicesReq struct {
	//固资编号
	FixedAssetNum string `json:"fixed_asset_number"`
	//序列号
	SN string `json:"sn"`
	//厂商
	Vendor string `json:"vendor"`
	//型号
	Model string `json:"model"`
	// CPU架构
	Arch string `json:"arch"`
	//用途 'TDSQL','APP','CVM','TGW','NAS','Other'
	Usage string `json:"usage"`
	//分类
	Category string `json:"category"`
	//机房管理单元名称
	ServerRoomName string `json:"server_room_name"`
	//机架编号
	CabinetNum string `json:"server_cabinet_number"`
	//机位编号
	USiteNum string `json:"server_usite_number"`
	//硬件说明
	HardwareRemark string `json:"hardware_remark"`
	//RAID说明
	RAIDRemark string `json:"raid_remark"`
	// OOB初始用户密码,':'分隔,主要针对的是旧机器导入使用，新机器如果带外是出厂默认，可以缺省
	OOBInit string `json:"oob_init"`
	//OriginNodeIP proxy节点IP
	//OriginNodeIP string `json:"origin_node_ip"` //废弃
	//启用时间
	StartedAt string `json:"started_at"`
	//上架时间
	OnShelveAt string `json:"onshelve_at"`
	//关联订单号(非必填)
	OrderNumber string `json:"order_number"`
	// 数据校验用
	ErrMsgContent string `json:"content"`
	//以上是对应Excel导入字段，以下字段是通过名称关联到
	idcID        uint //数据中心ID
	serverRoomID uint //机房管理单元ID
	cabinetID    uint //机架ID
	uSiteID      uint //机位ID
}

type ImportDevice2StoreReq struct {
	//固资编号-可选
	FixedAssetNum string `json:"fixed_asset_number"`
	//序列号
	SN string `json:"sn"`
	//厂商
	Vendor string `json:"vendor"`
	//型号
	Model string `json:"model"`
	//用途
	Usage string `json:"usage"`
	// 设备类型
	Category string `json:"category"`
	//启用时间
	StartedAt string `json:"started_at"`
	//上架时间
	OnShelveAt string `json:"onshelve_at"`	
	// 参考 DeviceLifecycle
	AssetBelongs	 				string		`json:"asset_belongs"`
	Owner			 				string		`json:"owner"`
	IsRental		 				string		`json:"is_rental"`
	MaintenanceServiceProvider		string		`json:"maintenance_service_provider"`
	MaintenanceService				string		`json:"maintenance_service"`
	LogisticsService				string		`json:"logistics_service"`
	MaintenanceServiceDateBegin     string 		`json:"maintenance_service_date_begin"`
	// 保修期（月数）
	MaintenanceMonths    int `json:"maintenance_months"`

	//库房管理单元名称
	StoreRoomName string `json:"store_room_name"`
	//虚拟货架编号
	VCabinetNum string `json:"virtual_cabinet_number"`
	//机位编号
	//USiteNum string `json:"server_usite_number"`
	//硬件说明
	HardwareRemark string `json:"hardware_remark"`
	//RAID说明
	RAIDRemark string `json:"raid_remark"`
	OOBInit    string `json:"oob_init"`
	//关联订单号(非必填)
	OrderNumber string `json:"order_number"`
	// 数据校验用
	ErrMsgContent string `json:"content"`

	//以上是对应Excel导入字段，以下字段是通过名称关联到
	idcID       uint //数据中心ID
	storeRoomID uint //机房管理单元ID
	vcabinetID  uint //机架ID
	//uSiteID      uint //机位ID
}

// ImportStockDevicesReq 导入存量机器结构
type ImportStockDevicesReq struct {
	ImportDevicesReq
	//内网IP
	IntranetIP string `json:"intranet_ip"`
	//外网IP，可以缺省
	ExtranetIP string `json:"extranet_ip"`
	//操作系统
	OS string `json:"os"`
	//运营状态
	OperationStatus string `json:"operation_status"`
}

//UpdateDevicesReq 设备导入Excel表对应字段
type UpdateDevicesReq struct {
	//ID
	ID uint `json:"id"`
	//固资编号
	FixedAssetNum string `json:"fixed_asset_number"`
	//序列号
	SN string `json:"sn"`
	//带外IP 不支持此处修改
	//OOBIP string `json:"oob_ip"`
	//带外用户 不支持此处修改
	//OOBUser string `json:"oob_user"`
	//带外密码 不支持此处修改
	//OOBPassword string `json:"oob_password"`
	//厂商
	Vendor string `json:"vendor"`
	//型号
	Model string `json:"model"`
	//用途 'TDSQL','APP','CVM','TGW','NAS','Other'
	Usage string `json:"usage"`
	//分类
	Category string `json:"category"`
	//数据中心ID
	IDCID uint `json:"idc_id"`
	//机房管理单元ID
	ServerRoomID uint `json:"server_room_id"`
	//机架ID
	CabinetID uint `json:"server_cabinet_id"`
	//机位ID
	USiteID uint `json:"server_usite_id"`
	//StoreRoomID 库房ID
	StoreRoomID uint `json:"store_room_id"`
	//虚拟货架ID
	VCabinetID uint `json:"virtual_cabinet_id"`
	//硬件说明
	HardwareRemark string `json:"hardware_remark"`
	//RAID说明
	RAIDRemark string `json:"raid_remark"`
	//启用时间
	StartedAt string `json:"started_at"`
	//上架时间
	OnShelveAt string `json:"onshelve_at"`
	//部署状态
	OperationStatus string `json:"operation_status"`
	// 用户登录名
	LoginName string `json:"-"`
}

//UpdateDevicesOperationStatusReq 修改设备部署状态
type UpdateDevicesOperationStatusReq struct {
	//固资编号
	FixedAssetNum string `json:"fixed_asset_number"`
	//序列号
	SN string `json:"sn"`
	//部署状态
	OperationStatus string `json:"operation_status"`
	// 用户登录名
	LoginName string `json:"-"`
}

//BatchUpdateDevicesReq 修改设备运营状态，用途等
type BatchUpdateDevicesReq struct {
	Devices []DeviceUpdateReq `json:"devices"`
	// 用户登录名
	LoginName string `json:"-"`
}

// 批量修改设备运营状态，用途请求结构
type DeviceUpdateReq struct {
	//固资编号
	FixedAssetNum string `json:"fixed_asset_number"`
	//序列号
	SN string `json:"sn"`
	//部署状态
	OperationStatus string `json:"operation_status"`
	//用途
	Usage string `json:"usage"`
	//硬件备注
	HardwareRemark string `json:"hardware_remark"`
}

//DeleteDevicesReq 删除设备请求体
type DeleteDevicesReq struct {
	IDs []uint `json:"ids"`
	SNs []string `json:"sns"`
}

//DevicePageResp 物理机分页
type DevicePageResp struct {
	ID uint `json:"id"`
	//固资编号
	FixedAssetNum string `json:"fixed_asset_number"`
	//序列号
	SN string `json:"sn"`
	//厂商
	Vendor string `json:"vendor"`
	//型号
	Model string `json:"model"`
	// 硬件架构
	Arch string `json:"arch"`
	//用途 'TDSQL','APP','CVM','TGW','NAS','Other'
	Usage string `json:"usage"`
	//分类
	Category string `json:"category"`
	//数据中心
	IDC *IDCSimplify `json:"idc"`
	//机房管理单元
	ServerRoom *ServerRoomSimplify `json:"server_room"`
	//物理区域
	//PhysicalArea string `json:"physical_area"`
	//机架
	ServerCabinet *ServerCabinetSimplify `json:"server_cabinet"`
	//机位
	ServerUSite *ServerUSiteSimplify `json:"server_usite"`
	//库房
	StoreRoom *StoreRoomSimplify `json:"store_room"`
	//虚拟货架
	VCabinet *VCabinetSimplify `json:"virtual_cabinets"`
	//带外IP
	OOBIP string `json:"oob_ip"`
	//带外用户
	OOBUser string `json:"oob_user"`
	//带外密码
	OOBPassword string `json:"oob_password"`
	//电源状态
	PowerStatus string `json:"power_status"`
	//带外纳管状态
	OOBAccessible string `json:"oob_accessible"`
	// TOR
	TOR string `json:"tor"`
	//硬件说明
	HardwareRemark string `json:"hardware_remark"`
	//RAID说明
	RAIDRemark string `json:"raid_remark"`
	//启用时间
	StartedAt string `json:"started_at"`
	//上架时间
	OnShelveAt string `json:"onshelve_at"`
	//运营状态
	OperationStatus string `json:"operation_status"`
	//内网IP
	IntranetIP string `json:"intranet_ip"`
	//外网IP
	ExtranetIP string `json:"extranet_ip"`
	//内网IPv6
	IntranetIPv6 string `json:"intranet_ipv6"`
	//外网IPv6
	ExtranetIPv6 string `json:"extranet_ipv6"`	
	//操作系统
	OS string `json:"os"`
	//OriginNode 代理proxy节点名
	OriginNode string `json:"origin_node"`
	//订单编号
	OrderNumber  string `json:"order_number"`
	//OriginNodeIP proxy代理IP
	OriginNodeIP string `json:"origin_node_ip"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// IDCSimplify 数据中心信息
type IDCSimplify struct {
	//数据中心ID
	ID uint `json:"id"`
	//数据中心名称
	Name string `json:"name"`
}

// ServerRoomSimplify 机房管理单元
type ServerRoomSimplify struct {
	//机房管理单元ID
	ID uint `json:"id"`
	//机房管理单元名称
	Name string `json:"name"`
}

// StoreRoomSimplify 库房管理单元
type StoreRoomSimplify struct {
	//机房管理单元ID
	ID uint `json:"id"`
	//机房管理单元名称
	Name string `json:"name"`
}

// ServerCabinetSimplify 机架
type ServerCabinetSimplify struct {
	//机架ID
	ID uint `json:"id"`
	//机架编号
	Number string `json:"number"`
}

// VCabinetSimplify
type VCabinetSimplify struct {
	//机架ID
	ID uint `json:"id"`
	//机架编号
	Number string `json:"number"`
}

// ServerUSiteSimplify 机位
type ServerUSiteSimplify struct {
	//机位ID
	ID uint `json:"id"`
	//机位编号
	Number string `json:"number"`
	// 物理区域
	PhysicalArea string `json:"physical_area"`
}

//DevicePageReq 物理机分页列表搜索字段
type DevicePageReq struct {
	ID string
	// idc
	IDCID             string `json:"idc_id"`
	ServerRoomID      string `json:"server_room_id"`
	ServerRoomName    string `json:"server_room_name"`
	PhysicalArea      string `json:"physical_area"`
	ServerCabinet     string `json:"server_cabinet_number"`
	ServerUSiteNumber string `json:"server_usite_number"`
	//运营状态
	OperationStatus string `json:"operation_status"`
	//用途 'TDSQL','APP','CVM','TGW','NAS','Other'
	Usage string `json:"usage"`
	//厂商
	Vendor string `json:"vendor"`
	//分类
	Category string `json:"category"`
	CategoryPreDeploy string `json:"category_pre_deploy"`
	//固资编号
	FixedAssetNum string `json:"fixed_asset_number"`
	//序列号
	SN string `json:"sn"`
	// 内网IP
	IntranetIP string `json:"intranet_ip"`
	// 外网IP
	ExtranetIP string `json:"extranet_ip"`
	IP         string `json:"ip"`
	//型号
	Model string `json:"model"`
	// 预部署状态物理机（没有任何安装记录的物理机）。
	PreDeployed bool `json:"pre_deployed"`
	// 硬件备注
	HardwareRemark string `json:"hardware_remark"`
	OOBAccessible  string
	Page           int64 `json:"-"`
	PageSize       int64 `json:"-"`
}

//ExportDevicesReq 物理机分页列表搜索字段
type ExportDevicesReq struct {
	//固资编号
	FixedAssetNum string `json:"fixed_asset_number"`
	//序列号
	SN string `json:"sn"`
	//厂商
	Vendor string `json:"vendor"`
	//型号
	Model string `json:"model"`
	//用途 'TDSQL','APP','CVM','TGW','NAS','Other'
	Usage string `json:"usage"`
	//分类
	Category string `json:"category"`
	//运营状态
	OperationStatus string `json:"operation_status"`
	//Page            int64
	//PageSize        int64
}

// CombinedDevice 设备信息及设备装机参数联合结构体
type CombinedDevice struct {
	DevicePageResp  			DevicePageResp        				`json:"device_page_resp"`
	DeployStatus    			string                				`json:"deploy_status"`
	InstallProgress 			float64               				`json:"install_progress"`
	BootOSIP        			string                				`json:"bootos_ip"`
	BootOSMac       			string                				`json:"bootos_mac"`
	CPU             			collector.CPU         				`json:"cpu"`
	Memory          			collector.Memory      				`json:"memory"`
	LogicDisk       			collector.Disk        				`json:"disk"`
	PhysicalDrive   			collector.DiskSlot    				`json:"disk_slot"`
	NIC             			collector.NIC         				`json:"nic"`
	Motherboard     			collector.Motherboard 				`json:"motherboard"`
	RAID            			collector.RAID        				`json:"raid"`
	OOB             			collector.OOB         				`json:"oob"`
	BIOS            			collector.BIOS        				`json:"bios"`
	Fan             			collector.Fan         				`json:"fan"`
	Power           			collector.Power       				`json:"power"`
	HBA             			collector.HBA         				`json:"hba"`
	PCI             			collector.PCI         				`json:"pci"`
	LLDP            			collector.LLDP        				`json:"lldp"`
	PowerSupplyNum  			int                   				`json:"power_supply_num"`
	Remark          			string                				`json:"remark"`
	IntranetIP       			string                   			`json:"intranet_ip"`
	ExtranetIP       			string                   			`json:"extranet_ip"`
	IntranetIPv6       			string                   			`json:"intranet_ipv6"`
	ExtranetIPv6       			string                   			`json:"extranet_ipv6"`	
	OS               			string                   			`json:"os"`
	ImageTemplate    			ImageTemplateSimplify    			`json:"image_tpl"`
	HardwareTemplate 			HardwareTemplateSimplify 			`json:"hardware_tpl"`
	Hostname         			string                   			`json:"hostname"`
	Inspection       			InspectionSimplify       			`json:"hardware_inspection"`
	DeviceLifecycleDetailPage  	DeviceLifecycleDetailPage        	`json:"device_lifecycle_detail_page"`
}

//ImageTemplateSimplify 系统镜像模板
type ImageTemplateSimplify struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

//HardwareTemplateSimplify 硬件配置模板（RAID，OOB）
type HardwareTemplateSimplify struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

//InspectionSimplify 巡检
type InspectionSimplify struct {
	RunStatus string `json:"run_status"`
	Result    string `json:"result"`
	Remark    string `json:"remark"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *ImportPreviewReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.Limit:    "limit",
		&reqData.Offset:   "offset",
		&reqData.FileName: "file_name",
	}
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *DeleteDevicesReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDs: "ids",
		&reqData.SNs: "sns",
	}
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *BatchUpdateDevicesReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.Devices: "devices",
	}
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *DevicePageReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.ID:                "id",
		&reqData.IDCID:             "idc_id",
		&reqData.ServerRoomID:      "server_room_id",
		&reqData.PhysicalArea:      "physical_area",
		&reqData.ServerRoomName:    "server_room_name",
		&reqData.FixedAssetNum:     "fixed_asset_number",
		&reqData.ServerCabinet:     "server_cabinet_number",
		&reqData.ServerUSiteNumber: "server_usite_number",
		&reqData.SN:                "sn",
		&reqData.IntranetIP:        "intranet_ip",
		&reqData.ExtranetIP:        "extranet_ip",
		&reqData.IP:                "ip", //可以搜索内网，外网
		&reqData.Vendor:            "vendor",
		&reqData.Model:             "model",
		&reqData.Usage:             "usage",
		&reqData.Category:          "category",
		&reqData.CategoryPreDeploy: "category_pre_deploy",
		&reqData.OperationStatus:   "operation_status",
		&reqData.PreDeployed:       "pre_deployed",
		&reqData.HardwareRemark:    "hardware_remark",
		&reqData.OOBAccessible:     "oob_accessible", //yes|no|unknown
		&reqData.Page:              "page",
		&reqData.PageSize:          "page_size",
	}
}

// Validate 结构体数据校验
func (reqData *DevicePageReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	//log, _ := middleware.LoggerFromContext(req.Context())
	// SpecialDev->特殊设备
	if reqData.Usage == "SpecialDev" {
		reqData.Usage = fmt.Sprintf("%s,%s", model.DevUsageSpecialDev, "SpecialDev")
	}
	// 校验设备状态是否正确
	if reqData.OperationStatus != "" {
		ops := strings.Split(reqData.OperationStatus, ",")
		for k := range ops {
			switch ops[k] {
			case model.DevOperStatRunWithAlarm:
			case model.DevOperStatRunWithoutAlarm:
			case model.DevOperStatReinstalling:
			case model.DevOperStatMoving:
			case model.DevOperStatPreRetire:
			case model.DevOperStatRetiring:
			case model.DevOperStateRetired:
			case model.DevOperStatPreDeploy:
			case model.DevOperStatOnShelve:
			case model.DevOperStatRecycling:
			case model.DevOperStatMaintaining:
			case model.DevOperStatPreMove:
			case model.DevOperStatInStore:
			default:
				errs.Add([]string{"sns"}, binding.BusinessError,
					fmt.Sprintf("设备状态：%s 不合法)", ops[k]))
				return errs
			}
		}
	}
	return errs
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *ExportDevicesReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.FixedAssetNum:   "fixed_asset_number",
		&reqData.SN:              "sn",
		&reqData.Vendor:          "vendor",
		&reqData.Model:           "model",
		&reqData.Usage:           "usage",
		&reqData.Category:        "category",
		&reqData.OperationStatus: "operation_status",
		//&reqData.Page:            "page",
		//&reqData.PageSize:        "page_size",
	}
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *UpdateDevicesReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.ID:              "id",
		&reqData.FixedAssetNum:   "fixed_asset_number",
		&reqData.SN:              "sn",
		&reqData.Vendor:          "vendor",
		&reqData.Model:           "model",
		&reqData.Usage:           "usage",
		&reqData.Category:        "category",
		&reqData.IDCID:           "idc_id",
		&reqData.ServerRoomID:    "server_room_id",
		&reqData.CabinetID:       "server_cabinet_id",
		&reqData.USiteID:         "server_usite_id",
		&reqData.HardwareRemark:  "hardware_remark",
		&reqData.RAIDRemark:      "raid_remark",
		&reqData.StartedAt:       "started_at",
		&reqData.OnShelveAt:      "onshelve_at",
		&reqData.OperationStatus: "operation_status",
	}
}

//ImportDevicesPreview 导入预览
func ImportDevicesPreview(log logger.Logger, repo model.Repo, reqData *ImportPreviewReq) (map[string]interface{}, error) {
	ra, err := utils.ParseDataFromXLSX(upload.UploadDir + reqData.FileName)
	if err != nil {
		return nil, err
	}
	length := len(ra)

	var success []*ImportDevicesReq
	var failure []*ImportDevicesReq

	if valid, err := CheckUnique(ra); !valid {
		return nil, err
	}

	//统计一下关联订单号的数量，用于更新订单状态
	var mOrderAmount = make(map[string]int, 0)

	for i := 1; i < length; i++ {
		row := &ImportDevicesReq{}
		if len(ra[i]) < 15 {
			var br string
			if row.ErrMsgContent != "" {
				br = "<br />"
			}
			row.ErrMsgContent += br + "导入文件列长度不对（应为15列）"
			failure = append(failure, row)
			continue
		}
		row.FixedAssetNum = ra[i][0]
		row.SN = ra[i][1]
		row.Vendor = ra[i][10]
		row.Model = ra[i][2]
		row.Usage = ra[i][3]
		row.Category = ra[i][4]
		row.ServerRoomName = ra[i][5]
		row.CabinetNum = ra[i][6]
		row.USiteNum = ra[i][7]
		row.HardwareRemark = ra[i][8]
		row.RAIDRemark = ra[i][9]
		row.StartedAt = ra[i][11]
		row.OnShelveAt = ra[i][12]
		row.OOBInit = ra[i][13]
		//row.OriginNodeIP = ra[i][14]
		row.OrderNumber = ra[i][14]

		utils.StructTrimSpace(row)

		//字段存在性校验
		row.checkLength()

		//以下这段时间转换的代码纯粹是为了转换下Excel中日期格式（如：1990/1/1）
		startedAt, _ := time.Parse(times.DateLayout2, ra[i][11])
		onShelveAt, _ := time.Parse(times.DateLayout2, ra[i][12])
		startedAtStr := startedAt.Format(times.DateLayout)
		if startedAtStr != "0001-01-01" {
			row.StartedAt = startedAtStr
		}
		onShelveAtStr := onShelveAt.Format(times.DateLayout)
		if onShelveAtStr != "0001-01-01" {
			row.OnShelveAt = onShelveAtStr
		}

		//数据有效性校验
		err := row.validate(repo)
		if err != nil {
			return nil, err
		}

		if row.ErrMsgContent != "" {
			failure = append(failure, row)
		} else {
			success = append(success, row)
		}

		dev, _ := repo.GetDeviceBySN(row.SN)
		if dev == nil && row.OrderNumber != "" {
			// 校验订单相关信息
			order, _ := repo.GetOrderByNumber(row.OrderNumber)
			if order != nil {
				switch order.Status {
				case model.OrderStatusCanceled:
					return nil,fmt.Errorf("订单:%s已取消", order.Number)
				case model.OrderStatusfinished:
					return nil,fmt.Errorf("订单:%s已确认完成", order.Number)
				}
				if row.Usage != order.Usage {
					return nil,fmt.Errorf("订单:%s 用途: %s 与导入的设备SN:%s 用途不匹配", order.Number, row.SN, order.Usage)
				}
			}
			mOrderAmount[row.OrderNumber]++ //新增才统计订单到货数量
		}
	}

	for orderNum, arrivalCount := range mOrderAmount {
		order, _ := repo.GetOrderByNumber(orderNum)
		if order != nil {
			if order.LeftAmount < arrivalCount {
				return nil, fmt.Errorf("订单(%s)新增到货数:%d应不大于剩余到货数量数量:%d", orderNum, arrivalCount, order.LeftAmount)
			}
		}
	}

	var data []*ImportDevicesReq
	if len(failure) > 0 {
		data = failure
	} else {
		data = success
	}
	var result []*ImportDevicesReq
	for i := 0; i < len(data); i++ {
		if uint(i) >= reqData.Offset && uint(i) < (reqData.Offset+reqData.Limit) {
			result = append(result, data[i])
		}
	}
	if len(failure) > 0 {
		_ = os.Remove(upload.UploadDir + reqData.FileName)
		return map[string]interface{}{"status": "failure",
			"message":       "导入服务器错误",
			"total_records": len(data),
			"content":       result,
		}, nil
	}
	return map[string]interface{}{"status": "success",
		"message":       "操作成功",
		"import_status": true,
		"total_records": len(data),
		"content":       result,
	}, nil
}

// ImportDevices 将设备放到数据库
func ImportDevices(log logger.Logger, repo model.Repo, conf *config.Config, reqData *ImportPreviewReq) error {
	fileName := upload.UploadDir + reqData.FileName
	log.Debugf("begin to parse data from file(%s)", fileName)
	ra, err := utils.ParseDataFromXLSX(fileName)
	if err != nil {
		log.Debugf("failed to parse data from file(%s)", fileName)
		return err
	}
	//把临时文件删了
	err = os.Remove(fileName)
	if err != nil {
		log.Warnf("remove tmp file: %s fail", fileName)
		return err
	}
	length := len(ra)

	var devices []*model.Device

	if valid, err := CheckUnique(ra); !valid {
		return err
	}

	//统计一下关联订单号的数量，用于更新订单状态
	var mOrderAmount = make(map[string]int, 0)
	for i := 1; i < length; i++ {
		row := &ImportDevicesReq{
			FixedAssetNum:  ra[i][0],
			SN:             ra[i][1],
			Vendor:         ra[i][10],
			Model:          ra[i][2],
			Usage:          ra[i][3],
			Category:       ra[i][4],
			ServerRoomName: ra[i][5],
			CabinetNum:     ra[i][6],
			USiteNum:       ra[i][7],
			HardwareRemark: ra[i][8],
			RAIDRemark:     ra[i][9],
			StartedAt:      ra[i][11],
			OnShelveAt:     ra[i][12],
			OOBInit:        ra[i][13],
			//OriginNodeIP:   ra[i][14],
			OrderNumber: ra[i][14],
		}
		if len(ra[i]) < 15 {
			continue
		}

		utils.StructTrimSpace(row)
		//处理所有字段的多余空格字符

		//必填项校验
		row.checkLength()

		//机房和网络区域校验
		err := row.validate(repo)
		if err != nil {
			return err
		}

		mod := &model.Device{
			SN:             row.SN,
			Vendor:         row.Vendor,
			DevModel:       row.Model,
			Usage:          row.Usage,
			Category:       row.Category,
			IDCID:          row.idcID,
			ServerRoomID:   row.serverRoomID,
			CabinetID:      row.cabinetID,
			USiteID:        &row.uSiteID,
			HardwareRemark: row.HardwareRemark,
			RAIDRemark:     row.RAIDRemark,
			OOBInit:        "{}",
			PowerStatus: model.PowerStatusOff,
			OrderNumber: row.OrderNumber,
			//JSON type的字段需要默认赋空值
			CPU:         "{}",
			Memory:      "{}",
			Disk:        "{}",
			DiskSlot:    "{}",
			NIC:         "{}",
			Motherboard: "{}",
			RAID:        "{}",
			OOB:         "{}",
			BIOS:        "{}",
			Fan:         "{}",
			Power:       "{}",
			HBA:         "{}",
			PCI:         "{}",
			Switch:      "{}",
			LLDP:        "{}",
			Extra:       "{}",
		}
		mod.StartedAt, _ = time.Parse(times.DateLayout2, row.StartedAt)
		mod.OnShelveAt, _ = time.Parse(times.DateLayout2, row.OnShelveAt)
		mod.OperationStatus = model.DevOperStatPreDeploy //default
		if row.OOBInit != "" {
			words := strings.Split(row.OOBInit, ":")
			if len(words) == 2 {
				ou := OOBUser{
					Username: words[0],
					Password: words[1],
				}
				if b, err := json.Marshal(ou); err == nil {
					mod.OOBInit = string(b)
				}
			}
		}

		//查询是否已经存在
		dev, err := repo.GetDeviceBySN(row.SN)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		isExits := false

		if dev != nil {
			mod.FixedAssetNumber = dev.FixedAssetNumber
			mod.Model = dev.Model
			if dev.OperationStatus == "" {
				dev.OperationStatus = model.DevOperStatPreDeploy //default
			}
			mod.OperationStatus = dev.OperationStatus
			mod.OriginNode = dev.OriginNode
			mod.Vendor = dev.Vendor
			mod.Arch = dev.Arch
			mod.CPUSum = dev.CPUSum
			mod.CPU = dev.CPU
			mod.MemorySum = dev.MemorySum
			mod.Memory = dev.Memory
			mod.DiskSum = dev.DiskSum
			mod.Disk = dev.Disk
			mod.DiskSlot = dev.DiskSlot
			mod.NIC = dev.NIC
			mod.NICDevice = dev.NICDevice
			mod.BootOSIP = dev.BootOSIP
			mod.BootOSMac = dev.BootOSMac
			mod.Motherboard = dev.Motherboard
			mod.RAID = dev.RAID
			mod.OOB = dev.OOB
			mod.OOBIP = dev.OOBIP
			mod.OOBUser = dev.OOBUser
			mod.OOBPassword = dev.OOBPassword
			mod.BIOS = dev.BIOS
			mod.Fan = dev.Fan
			mod.Power = dev.Power
			mod.HBA = dev.HBA
			mod.PCI = dev.PCI
			mod.Switch = dev.Switch
			mod.LLDP = dev.LLDP
			mod.Extra = dev.Extra
			mod.HBA = dev.HBA
			mod.PowerStatus = dev.PowerStatus
			mod.Updater = reqData.LoginName
			isExits = true
		} else {
			//自动生成固资编号，新增场景才这么干！
			if row.FixedAssetNum == "" {
				row.FixedAssetNum, err = GenFixedAssetNumber(repo)
				if err != nil {
					log.Errorf("generate fixed_asset number for SN:%s fail", row.SN)
					return fmt.Errorf("自动生成固资号失败：%v", err)
				}
			}
			mod.FixedAssetNumber = row.FixedAssetNum
			mod.Creator = reqData.LoginName
		}
		// 默认电源状态OFF
		if mod.PowerStatus == "" {
			mod.PowerStatus = model.PowerStatusOff
		}
		// 仅记录必要字段到“设备新增”
		optDetail, err := convert2DetailOfOperationTypeAdd(repo, *mod)
		if err != nil {
			log.Errorf("Fail to convert Detail of OperationTypeAdd: %v", err)
		}
		// DeviceLifecycle 变更记录
		deviceLifecycleLog := []model.ChangeLog {
			{
				OperationUser:		reqData.LoginName,
				OperationType:		model.OperationTypeAdd,
				OperationDetail:	optDetail,
				OperationTime:		times.ISOTime(time.Now()).ToTimeStr(),
			},
		}
		b, _ := json.Marshal(deviceLifecycleLog)
		// SaveDeviceLifecycleReq 结构体
		saveDevLifecycleReq := &SaveDeviceLifecycleReq {
			DeviceLifecycleBase: DeviceLifecycleBase{
				FixedAssetNumber: 				mod.FixedAssetNumber,
				SN:             				mod.SN,
				AssetBelongs:					"Undefined",
				Owner:							"Undefined",
				IsRental:						"no",
				MaintenanceServiceProvider:		"Undefined",
				MaintenanceService:				"Undefined",
				LogisticsService:				"Undefined",
				MaintenanceServiceStatus:		model.MaintenanceServiceStatusInactive, //新增场景默认-未激活
				LifecycleLog:					string(b),
			},
		}		
		// 通过订单编号获取资产归属、负责人、维保服务等内容
		// 若无订单编号则以参数输入为准
		if row.OrderNumber != "" {
			order, err := repo.GetOrderByNumber(row.OrderNumber)
			if err != nil {
				log.Errorf("订单(订单号:%s)不存在", row.OrderNumber)
				return err
			}
			if order != nil {
				saveDevLifecycleReq.AssetBelongs = order.AssetBelongs	 			
				saveDevLifecycleReq.Owner = order.Owner
				saveDevLifecycleReq.IsRental = order.IsRental
				saveDevLifecycleReq.MaintenanceServiceProvider = order.MaintenanceServiceProvider
				saveDevLifecycleReq.MaintenanceService = order.MaintenanceService
				saveDevLifecycleReq.LogisticsService = order.LogisticsService
				saveDevLifecycleReq.MaintenanceServiceDateBegin = order.MaintenanceServiceDateBegin
				saveDevLifecycleReq.MaintenanceServiceDateEnd = order.MaintenanceServiceDateEnd
			}
		}
		//插入或者更新
		if _, err = repo.SaveDevice(mod); err != nil {
			return err
		}
		// DeviceLifecycle 查询是否已经存在
		devLifecycle, err := repo.GetDeviceLifecycleBySN(mod.SN)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if devLifecycle != nil {
			log.Debugf("Importing device DeviceLifecycle SN %s already exist.Update it.", mod.SN)
			saveDevLifecycleReq.ID = devLifecycle.ID
		}
		// 若dev已经存在，且密码不为默认出厂密码，则不予修改
		if isExits && checkOriginPassword(log, conf, mod) {
			continue
		}

		devices = append(devices, mod)
		//如果是重复性（更新性质的）导入，不重复更新订单
		if !isExits {
			if row.OrderNumber != "" {
				mOrderAmount[row.OrderNumber]++
			}
		}
		//修改机位占用状态
		if mod.USiteID != nil {
			if _, err = repo.BatchUpdateServerUSitesStatus([]uint{*mod.USiteID}, model.USiteStatUsed); err != nil {
				log.Errorf("update server_usite status failed, usite_num :%s", row.USiteNum)
			}
		}
		// 保存或更新 DeviceLifecycle
		if err = SaveDeviceLifecycle(log, repo, saveDevLifecycleReq); err != nil {
			log.Debug(err)
			return err
		}		
	}

	//更新关联的订单到货数量和订单状态
	for orderNum, arrivalCount := range mOrderAmount {
		if err = UpdateOrderByArrival(log, repo, orderNum, arrivalCount); err != nil {
			return err
		}
	}

	// 导入设备成功后，批量修改密码
	go batchUpdateOOBPassword(log, repo, conf, devices)

	os.Remove(fileName)

	return nil
}

//ImportDevices2StorePreview 导入预览
func ImportDevices2StorePreview(log logger.Logger, repo model.Repo, reqData *ImportPreviewReq) (map[string]interface{}, error) {
	ra, err := utils.ParseDataFromXLSX(upload.UploadDir + reqData.FileName)
	if err != nil {
		return nil, err
	}
	length := len(ra)

	var success []*ImportDevice2StoreReq
	var failure []*ImportDevice2StoreReq

	//if valid, err := CheckUnique(ra); !valid {
	//	return nil, err
	//}
	//统计新增设备。有些重复导入修改的，不统计，用与订单剩余到货数比较
	var mOrderAmount = make(map[string]int, 0)

	for i := 1; i < length; i++ {
		row := &ImportDevice2StoreReq{}
		if len(ra[i]) < 20 {
			var br string
			if row.ErrMsgContent != "" {
				br = "<br />"
			}
			row.ErrMsgContent += br + "导入文件列长度不对（应为21列）"
			failure = append(failure, row)
			continue
		}
		//row.FixedAssetNum = ra[i][0]
		row.SN = ra[i][0]
		row.Vendor = ra[i][1]
		row.Model = ra[i][2]
		row.Usage = ra[i][3]
		row.Category = ra[i][4]
		row.StoreRoomName = ra[i][5]
		row.VCabinetNum = ra[i][6]
		row.HardwareRemark = ra[i][7]
		row.RAIDRemark = ra[i][8]
		row.OOBInit = ra[i][9]
		row.OrderNumber = ra[i][10]
		row.Owner = ra[i][11]

		//转换下Excel中日期格式
		maintenanceDate, _ := time.Parse(times.DateLayout2, ra[i][12])
		maintenanceDateStr := maintenanceDate.Format(times.DateLayout)
		if maintenanceDateStr != "0001-01-01" {
			row.MaintenanceServiceDateBegin = maintenanceDateStr
		} else {
			row.MaintenanceServiceDateBegin = ""
		}
		maintenanceMonths, err := strconv.Atoi(ra[i][13])
		if err != nil {
			return nil, err
		}
		row.MaintenanceMonths = maintenanceMonths
		row.AssetBelongs = ra[i][14]
		row.IsRental = ra[i][15]
		row.MaintenanceServiceProvider = ra[i][16]
		// 允许为空		
		row.MaintenanceService = ra[i][17]
		row.LogisticsService = ra[i][18]
		row.FixedAssetNum = ra[i][19]
		row.StartedAt = ra[i][20]
		row.OnShelveAt = ra[i][21]

		utils.StructTrimSpace(row)

		//字段存在性校验
		row.checkLength()

		//数据有效性校验
		err = row.validate(repo)
		if err != nil {
			return nil, err
		}

		if row.ErrMsgContent != "" {
			failure = append(failure, row)
		} else {
			success = append(success, row)
		}

		dev, _ := repo.GetDeviceBySN(row.SN)
		if dev == nil && row.OrderNumber != "" {
			mOrderAmount[row.OrderNumber]++ //新增才统计订单信息
		}
	}

	for orderNum, arrivalCount := range mOrderAmount {
		order, _ := repo.GetOrderByNumber(orderNum)
		if order != nil {
			switch order.Status {
			case model.OrderStatusCanceled:
				return nil,fmt.Errorf("订单:%s已取消", order.Number)
			case model.OrderStatusfinished:
				return nil,fmt.Errorf("订单:%s已确认完成", order.Number)
			}
			if order.LeftAmount < arrivalCount {
				return nil, fmt.Errorf("订单(%s)新增到货数:%d应不大于剩余到货数量数量:%d", orderNum, arrivalCount, order.LeftAmount)
			}
		}
	}

	var data []*ImportDevice2StoreReq
	if len(failure) > 0 {
		data = failure
	} else {
		data = success
	}
	var result []*ImportDevice2StoreReq
	for i := 0; i < len(data); i++ {
		if uint(i) >= reqData.Offset && uint(i) < (reqData.Offset+reqData.Limit) {
			result = append(result, data[i])
		}
	}
	if len(failure) > 0 {
		_ = os.Remove(upload.UploadDir + reqData.FileName)
		return map[string]interface{}{"status": "failure",
			"message":       "导入服务器错误",
			"total_records": len(data),
			"content":       result,
		}, nil
	}
	return map[string]interface{}{"status": "success",
		"message":       "操作成功",
		"import_status": true,
		"total_records": len(data),
		"content":       result,
	}, nil
}

// ImportDevices2Store 将设备放到导入到库房
func ImportDevices2Store(log logger.Logger, repo model.Repo, conf *config.Config, reqData *ImportPreviewReq) error {
	fileName := upload.UploadDir + reqData.FileName
	log.Debugf("begin to parse data from file(%s)", fileName)
	ra, err := utils.ParseDataFromXLSX(fileName)
	if err != nil {
		log.Debugf("failed to parse data from file(%s)", fileName)
		return err
	}
	//把临时文件删了
	log.Debugf("begin to remove file(%s)", fileName)
	err = os.Remove(fileName)
	if err != nil {
		log.Warnf("remove tmp file: %s fail", fileName)
		return err
	}

	length := len(ra)

	//统计一下关联订单号的数量，用于更新订单状态
	var mOrderAmount = make(map[string]int, 0)

	var mods []*model.Device

	for i := 1; i < length; i++ {
		if len(ra[i]) < 22 {
			continue
		}
		row := &ImportDevice2StoreReq{}
		row.SN = ra[i][0]
		row.Vendor = ra[i][1]
		row.Model = ra[i][2]
		row.Usage = ra[i][3]
		row.Category = ra[i][4]
		row.StoreRoomName = ra[i][5]
		row.VCabinetNum = ra[i][6]
		row.HardwareRemark = ra[i][7]
		row.RAIDRemark = ra[i][8]
		row.OOBInit = ra[i][9]
		row.OrderNumber = ra[i][10]
		row.Owner = ra[i][11]

		//转换下Excel中日期格式
		maintenanceDate, _ := time.Parse(times.DateLayout2, ra[i][12])
		maintenanceDateStr := maintenanceDate.Format(times.DateLayout)
		if maintenanceDateStr != "0001-01-01" {
			row.MaintenanceServiceDateBegin = maintenanceDateStr
		} else {
			row.MaintenanceServiceDateBegin = ""
		}
		maintenanceMonths, err := strconv.Atoi(ra[i][13])
		if err != nil {
			return err
		}
		row.MaintenanceMonths = maintenanceMonths
		row.AssetBelongs = ra[i][14]
		row.IsRental = ra[i][15]
		row.MaintenanceServiceProvider = ra[i][16]
		// 允许为空		
		row.MaintenanceService = ra[i][17]
		row.LogisticsService = ra[i][18]
		row.FixedAssetNum = ra[i][19]
		row.StartedAt = ra[i][20]
		row.OnShelveAt = ra[i][21]

		//处理所有字段的多余空格字符
		utils.StructTrimSpace(row)

		//必填项校验
		row.checkLength()

		//机房和网络区域校验
		err = row.validate(repo)
		if err != nil {
			return err
		}

		mod := &model.Device{
			FixedAssetNumber: row.FixedAssetNum,
			SN:             row.SN,
			Vendor:         row.Vendor,
			DevModel:       row.Model,
			Usage:          row.Usage,
			Category:       row.Category,
			IDCID:          row.idcID,
			StoreRoomID:    row.storeRoomID,
			VCabinetID:     row.vcabinetID,
			HardwareRemark: row.HardwareRemark,
			RAIDRemark:     row.RAIDRemark,
			PowerStatus:    model.PowerStatusOff,
		}
		//mod.StartedAt, _ = time.Parse(times.DateLayout2, "01-01-70") //填充一个非法时间
		//mod.OnShelveAt, _ = time.Parse(times.DateLayout2, "01-01-70")
		mod.PowerStatus = model.PowerStatusOff
		now := time.Now()
		//mod.StartedAt = now
		//mod.OnShelveAt = now
		mod.CreatedAt = now
		mod.StartedAt, _ = time.Parse(times.DateLayout2, row.StartedAt)
		mod.OnShelveAt, _ = time.Parse(times.DateLayout2, row.OnShelveAt)	
		mod.OperationStatus = model.DevOperStatInStore //default
		if row.OOBInit != "" {
			words := strings.Split(row.OOBInit, ":")
			if len(words) == 2 {
				ou := OOBUser{
					Username: words[0],
					Password: words[1],
				}
				if b, err := json.Marshal(ou); err == nil {
					mod.OOBInit = string(b)
				}
			}
		}

		//查询是否已经存在
		dev, err := repo.GetDeviceBySN(row.SN)
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Errorf("check device:%s exist fail,err:%v", row.SN, err)
			return err
		}

		if dev != nil {
			mod.FixedAssetNumber = dev.FixedAssetNumber
			mod.Model = dev.Model
			mod.OriginNode = dev.OriginNode
			mod.Arch = dev.Arch
			mod.CPUSum = dev.CPUSum
			mod.CPU = dev.CPU
			mod.MemorySum = dev.MemorySum
			mod.Memory = dev.Memory
			mod.DiskSum = dev.DiskSum
			mod.Disk = dev.Disk
			mod.DiskSlot = dev.DiskSlot
			mod.NIC = dev.NIC
			mod.NICDevice = dev.NICDevice
			mod.BootOSIP = dev.BootOSIP
			mod.BootOSMac = dev.BootOSMac
			mod.Motherboard = dev.Motherboard
			mod.RAID = dev.RAID
			mod.OOB = dev.OOB
			mod.OOBIP = dev.OOBIP
			mod.OOBUser = dev.OOBUser
			mod.OOBPassword = dev.OOBPassword
			mod.BIOS = dev.BIOS
			mod.Fan = dev.Fan
			mod.Power = dev.Power
			mod.HBA = dev.HBA
			mod.PCI = dev.PCI
			mod.Switch = dev.Switch
			mod.LLDP = dev.LLDP
			mod.Extra = dev.Extra
			mod.HBA = dev.HBA
			mod.PowerStatus = dev.PowerStatus
			mod.Updater = reqData.LoginName
		} else {
			//自动生成固资编号
			if row.FixedAssetNum == "" {
				row.FixedAssetNum, err = GenFixedAssetNumber(repo)
				if err != nil {
					log.Errorf("generate fixed_asset number for SN:%s fail", row.SN)
					return err
				}
			}
			mod.FixedAssetNumber = row.FixedAssetNum
			mod.Creator = reqData.LoginName
			if row.OrderNumber != "" {
				mOrderAmount[row.OrderNumber]++ //新增才统计订单信息
			}
		}
		// 仅记录必要字段到“设备新增”
		optDetail, err := convert2DetailOfOperationTypeAdd(repo, *mod)
		if err != nil {
			log.Errorf("Fail to convert Detail of OperationTypeAdd: %v", err)
		}		
		// DeviceLifecycle 变更记录
		deviceLifecycleLog := []model.ChangeLog {
			{
				OperationUser:		reqData.LoginName,
				OperationType:		model.OperationTypeAdd,
				OperationDetail:	optDetail,
				OperationTime:		times.ISOTime(time.Now()).ToTimeStr(),
			},
		}
		b, _ := json.Marshal(deviceLifecycleLog)

		// SaveDeviceLifecycleReq 结构体
		saveDevLifecycleReq := &SaveDeviceLifecycleReq {
			DeviceLifecycleBase: DeviceLifecycleBase{
				FixedAssetNumber: 				mod.FixedAssetNumber,
				SN:             				mod.SN,
				AssetBelongs:					row.AssetBelongs,
				Owner:							row.Owner,
				IsRental:						row.IsRental,
				MaintenanceServiceProvider:		row.MaintenanceServiceProvider,
				MaintenanceService:				row.MaintenanceService,
				LogisticsService:				row.LogisticsService,
				MaintenanceServiceStatus:		model.MaintenanceServiceStatusInactive, //新增场景默认-未激活
				LifecycleLog:					string(b),
			},
		}

		if row.MaintenanceServiceDateBegin != "" {
			t, err := time.Parse(times.DateLayout, row.MaintenanceServiceDateBegin)
			if err != nil {
				log.Errorf("parse maintenance time %s err:%v , using current time for maintenance-date", row.MaintenanceServiceDateBegin, err)
				saveDevLifecycleReq.MaintenanceServiceDateBegin = now
				saveDevLifecycleReq.MaintenanceServiceDateEnd = now.AddDate(0, row.MaintenanceMonths, 0)
			} else {
				saveDevLifecycleReq.MaintenanceServiceDateBegin = t
				saveDevLifecycleReq.MaintenanceServiceDateEnd = t.AddDate(0, row.MaintenanceMonths, 0)
			}
		} else {
			saveDevLifecycleReq.MaintenanceServiceDateBegin = now
			saveDevLifecycleReq.MaintenanceServiceDateEnd = now.AddDate(0, row.MaintenanceMonths, 0)
		}

	
		// 通过订单编号获取资产归属、负责人、维保服务等内容
		// 若无订单编号则以参数输入为准
		if row.OrderNumber != "" {
			order, err := repo.GetOrderByNumber(row.OrderNumber)
			if err != nil {
				log.Errorf("订单(订单号:%s)不存在", row.OrderNumber)
				return err
			}
			if order != nil {
				mod.OrderNumber = row.OrderNumber
				saveDevLifecycleReq.AssetBelongs = order.AssetBelongs			
				saveDevLifecycleReq.Owner = order.Owner
				saveDevLifecycleReq.IsRental = order.IsRental
				saveDevLifecycleReq.MaintenanceServiceProvider = order.MaintenanceServiceProvider
				saveDevLifecycleReq.MaintenanceService = order.MaintenanceService
				saveDevLifecycleReq.LogisticsService = order.LogisticsService
				saveDevLifecycleReq.MaintenanceServiceDateBegin = order.MaintenanceServiceDateBegin
				saveDevLifecycleReq.MaintenanceServiceDateEnd = order.MaintenanceServiceDateEnd
			}
		}
		// DeviceLifecycle 查询是否已经存在
		devLifecycle, err := repo.GetDeviceLifecycleBySN(mod.SN)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		} 
		if devLifecycle != nil {
			log.Debugf("DeviceLifecycle SN %s already exist.Update it.", mod.SN)
			saveDevLifecycleReq.ID = devLifecycle.ID
		}
		mods = append(mods, mod)
		// 保存或更新 DeviceLifecycle
		if err = SaveDeviceLifecycle(log, repo, saveDevLifecycleReq); err != nil {
			log.Debug(err)
			return err
		}
	}

	//更新关联的订单到货数量和订单状态
	for orderNum, arrivalCount := range mOrderAmount {
		if err = UpdateOrderByArrival(log, repo, orderNum, arrivalCount); err != nil {
			return err
		}
	}

	// 更新完订单再更新设备信息
	for i := range mods {
		if _, err = repo.SaveDevice(mods[i]); err != nil {
			return err
		}
	}
	//TODO 失败回滚
	return nil
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *UpdateDevicesOperationStatusReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.FixedAssetNum:   "fixed_asset_number",
		&reqData.SN:              "sn",
		&reqData.OperationStatus: "operation_status",
	}
}

// Validate 对修改的数据做基本校验
func (reqData *UpdateDevicesOperationStatusReq) Validate(request *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(request.Context())

	if reqData.SN == "" && reqData.FixedAssetNum == "" {
		errs.Add([]string{"设备SN、固资编号"}, binding.RequiredError, fmt.Sprintf("设备SN、固资编号不能同时为空"))
		return errs
	}

	if reqData.OperationStatus != model.DevOperStatRunWithoutAlarm &&
		reqData.OperationStatus != model.DevOperStatRunWithAlarm &&
		reqData.OperationStatus != model.DevOperStateRetired &&
		reqData.OperationStatus != model.DevOperStatMoving &&
		reqData.OperationStatus != model.DevOperStatOnShelve &&
		reqData.OperationStatus != model.DevOperStatPreDeploy &&
		reqData.OperationStatus != model.DevOperStatPreRetire &&
		reqData.OperationStatus != model.DevOperStatRetiring &&
		reqData.OperationStatus != model.DevOperStatRecycling &&
		reqData.OperationStatus != model.DevOperStatPreMove &&
		reqData.OperationStatus != model.DevOperStatMaintaining &&
		reqData.OperationStatus != model.DevOperStatReinstalling {
		errs.Add([]string{"运营状态"}, binding.RequiredError, fmt.Sprintf("运营状态不正确"))
		return errs
	}

	if d, _ := repo.GetDeviceByFixAssetNumber(reqData.FixedAssetNum); d == nil {
		if d, _ = repo.GetDeviceBySN(reqData.SN); d == nil {
			errs.Add([]string{"设备SN"}, binding.RequiredError, fmt.Sprintf("设备信息不存在，SN[%s]，固资号[%s]", reqData.SN, reqData.FixedAssetNum))
			return errs
		}

	}

	return nil
}

// Validate 对修改的数据做基本校验
func (reqData *BatchUpdateDevicesReq) Validate(request *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(request.Context())

	for _, d := range reqData.Devices {
		if d.SN == "" && d.FixedAssetNum == "" {
			errs.Add([]string{"设备SN、固资编号"}, binding.RequiredError, fmt.Sprintf("设备SN、固资编号不能同时为空"))
			return errs
		}

		if d.OperationStatus != "" && !validOperationStatus(d.OperationStatus) {
			errs.Add([]string{"运营状态"}, binding.RequiredError, fmt.Sprintf("运营状态不正确"))
			return errs
		}

		//devUsageEnum := []string{model.DevUsageDB, model.DevUsageAPP, model.DevUsageTGW, model.DevUsageCVM,
		//	model.DevUsageStorage, model.DevUsageContainer, model.DevUsageBDP, model.DevUsageSpecialDev, model.DevUsageOther}
		//if !ValidateEnum(d.Usage, devUsageEnum) {
		//	errs.Add([]string{"用途"}, binding.RequiredError, fmt.Sprintf("用途值(%s)不合法", d.Usage))
		//	return errs
		//}

		if dd, _ := repo.GetDeviceByFixAssetNumber(d.FixedAssetNum); dd == nil {
			if dd, _ = repo.GetDeviceBySN(d.SN); dd == nil {
				errs.Add([]string{"设备SN"}, binding.RequiredError,
					fmt.Sprintf("设备信息不存在，SN[%s]，固资号[%s]", d.SN, d.FixedAssetNum))
				return errs
			}

		}
	}

	return nil
}

//ImportStockDevicesPreview 导入存量设备预览
func ImportStockDevicesPreview(log logger.Logger, repo model.Repo, reqData *ImportPreviewReq) (map[string]interface{}, error) {
	ra, err := utils.ParseDataFromXLSX(upload.UploadDir + reqData.FileName)
	if err != nil {
		return nil, err
	}
	length := len(ra)

	var success []*ImportStockDevicesReq
	var failure []*ImportStockDevicesReq

	//if valid, err := CheckUnique(ra); !valid {
	//	return nil, err
	//}
	for i := 1; i < length; i++ {
		row := &ImportStockDevicesReq{}
		if len(ra[i]) < 19 {
			var br string
			if row.ErrMsgContent != "" {
				br = "<br />"
			}
			row.ErrMsgContent += br + "导入文件列长度不对（应为19列）"
			failure = append(failure, row)
			continue
		}
		row.FixedAssetNum = ra[i][0]
		row.SN = ra[i][1]
		row.Model = ra[i][2]
		row.Arch = ra[i][3]
		row.Usage = ra[i][4]
		row.Category = ra[i][5]
		row.ServerRoomName = ra[i][6]
		row.CabinetNum = ra[i][7]
		row.USiteNum = ra[i][8]
		row.HardwareRemark = ra[i][9]
		row.RAIDRemark = ra[i][10]
		row.Vendor = ra[i][11]
		row.StartedAt = ra[i][12]
		row.OnShelveAt = ra[i][13]
		row.OOBInit = ra[i][14]
		//row.OriginNodeIP = ra[i][15]
		row.IntranetIP = ra[i][15]
		row.ExtranetIP = ra[i][16]
		row.OS = ra[i][17]
		row.OperationStatus = ra[i][18]

		utils.StructTrimSpace(row)
		utils.StructTrimSpace(&row.ImportDevicesReq)

		//字段存在性校验
		row.checkLength()

		//以下这段时间转换的代码纯粹是为了转换下Excel中日期格式
		startedAt, _ := time.Parse(times.DateLayout2, ra[i][11])
		onShelveAt, _ := time.Parse(times.DateLayout2, ra[i][12])
		startedAtStr := startedAt.Format(times.DateLayout)
		if startedAtStr != "0001-01-01" {
			row.StartedAt = startedAtStr
		}
		onShelveAtStr := onShelveAt.Format(times.DateLayout)
		if onShelveAtStr != "0001-01-01" {
			row.OnShelveAt = onShelveAtStr
		}

		//数据有效性校验
		err := row.validate(repo)
		if err != nil {
			return nil, err
		}

		if row.ErrMsgContent != "" {
			failure = append(failure, row)
		} else {
			success = append(success, row)
		}
	}

	var data []*ImportStockDevicesReq
	if len(failure) > 0 {
		data = failure
	} else {
		data = success
	}
	var result []*ImportStockDevicesReq
	for i := 0; i < len(data); i++ {
		if uint(i) >= reqData.Offset && uint(i) < (reqData.Offset+reqData.Limit) {
			result = append(result, data[i])
		}
	}
	if len(failure) > 0 {
		_ = os.Remove(upload.UploadDir + reqData.FileName)
		return map[string]interface{}{"status": "failure",
			"message":       "导入服务器错误",
			"total_records": len(data),
			"content":       result,
		}, nil
	}
	return map[string]interface{}{"status": "success",
		"message":       "操作成功",
		"import_status": true,
		"total_records": len(data),
		"content":       result,
	}, nil
}

// ImportStockDevices 将设备放到数据库
func ImportStockDevices(log logger.Logger, repo model.Repo, conf *config.Config, reqData *ImportPreviewReq) error {
	fileName := upload.UploadDir + reqData.FileName
	ra, err := utils.ParseDataFromXLSX(fileName)
	if err != nil {
		return err
	}
	//把临时文件删了
	err = os.Remove(fileName)
	if err != nil {
		log.Warnf("remove tmp file: %s fail", fileName)
		return err
	}
	length := len(ra)

	var devices []*model.Device

	//if valid, err := CheckUnique(ra); !valid {
	//	return err
	//}

	for i := 1; i < length; i++ {
		row := &ImportStockDevicesReq{
			ImportDevicesReq: ImportDevicesReq{
				FixedAssetNum:  ra[i][0],
				SN:             ra[i][1],
				Vendor:         ra[i][11],
				Model:          ra[i][2],
				Arch:           ra[i][3],
				Usage:          ra[i][4],
				Category:       ra[i][5],
				ServerRoomName: ra[i][6],
				CabinetNum:     ra[i][7],
				USiteNum:       ra[i][8],
				HardwareRemark: ra[i][9],
				RAIDRemark:     ra[i][10],
				StartedAt:      ra[i][12],
				OnShelveAt:     ra[i][13],
				OOBInit:        ra[i][14],
				//OriginNodeIP:   ra[i][15],
			},
			IntranetIP:      ra[i][15],
			ExtranetIP:      ra[i][16],
			OS:              ra[i][17],
			OperationStatus: ra[i][18],
		}
		if len(ra[i]) < 19 {
			continue
		}

		//处理所有字段的多余空格字符
		utils.StructTrimSpace(row)
		utils.StructTrimSpace(&row.ImportDevicesReq)

		//必填项校验
		row.checkLength()

		//机房和网络区域校验
		err := row.validate(repo)
		if err != nil {
			return err
		}

		mod := &model.Device{
			FixedAssetNumber: row.FixedAssetNum,
			SN:               row.SN,
			Vendor:           row.Vendor,
			DevModel:         row.Model,
			Arch:             row.Arch,
			Usage:            row.Usage,
			Category:         row.Category,
			IDCID:            row.idcID,
			ServerRoomID:     row.serverRoomID,
			CabinetID:        row.cabinetID,
			USiteID:          &row.uSiteID,
			HardwareRemark:   row.HardwareRemark,
			RAIDRemark:       row.RAIDRemark,
			OOBInit:          "{}",
			//OriginNodeIP:     row.OriginNodeIP,
			PowerStatus:     model.PowerStatusOn, //默认认为是开的
			OperationStatus: OperationStatusTransfer(row.OperationStatus, false),
			//JSON type的字段需要默认赋空值
			CPU:         "{}",
			Memory:      "{}",
			Disk:        "{}",
			DiskSlot:    "{}",
			NIC:         "{}",
			Motherboard: "{}",
			RAID:        "{}",
			OOB:         "{}",
			BIOS:        "{}",
			Fan:         "{}",
			Power:       "{}",
			HBA:         "{}",
			PCI:         "{}",
			Switch:      "{}",
			LLDP:        "{}",
			Extra:       "{}",
			Remark:      "存量导入",
		}
		mod.StartedAt, _ = time.Parse(times.DateLayout2, row.StartedAt)
		mod.OnShelveAt, _ = time.Parse(times.DateLayout2, row.OnShelveAt)

		if row.OOBInit != "" {
			words := strings.Split(row.OOBInit, ":")
			if len(words) == 2 {
				ou := OOBUser{
					Username: words[0],
					Password: words[1],
				}
				if b, err := json.Marshal(ou); err == nil {
					mod.OOBInit = string(b)
				}
			}
		}

		//查询是否已经存在
		dev, err := repo.GetDeviceBySN(row.SN)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		isExits := false

		if dev != nil {
			mod.Model = dev.Model
			mod.OriginNode = dev.OriginNode
			mod.CPUSum = dev.CPUSum
			mod.CPU = dev.CPU
			mod.MemorySum = dev.MemorySum
			mod.Memory = dev.Memory
			mod.DiskSum = dev.DiskSum
			mod.Disk = dev.Disk
			mod.DiskSlot = dev.DiskSlot
			mod.NIC = dev.NIC
			mod.NICDevice = dev.NICDevice
			mod.BootOSIP = dev.BootOSIP
			mod.BootOSMac = dev.BootOSMac
			mod.Motherboard = dev.Motherboard
			mod.RAID = dev.RAID
			mod.OOB = dev.OOB
			mod.OOBIP = dev.OOBIP
			mod.OOBUser = dev.OOBUser
			mod.OOBPassword = dev.OOBPassword
			mod.BIOS = dev.BIOS
			mod.Fan = dev.Fan
			mod.Power = dev.Power
			mod.HBA = dev.HBA
			mod.PCI = dev.PCI
			mod.Switch = dev.Switch
			mod.LLDP = dev.LLDP
			mod.Extra = dev.Extra
			mod.HBA = dev.HBA
			mod.PowerStatus = dev.PowerStatus
			mod.Updater = reqData.LoginName
			isExits = true
		} else {
			mod.Creator = reqData.LoginName
		}

		//再次用带外远程更新下开关机状态
		if pwrStatus, err := GetDevicePowerStatusBySN(log, repo, conf, row.SN); err != nil {
			log.Errorf("get power status by sn:%s fail:%v", row.SN, err)
		} else if pwrStatus != "" {
			mod.PowerStatus = pwrStatus
		}

		//插入或者更新
		if _, err = repo.SaveDevice(mod); err != nil {
			return err
		}

		//如果设备被导入过，则修改
		if isExits {
			//释放之前占用的机位，因为这次导入的数据可能机位改了，需要释放之前的，后边占用更新后的
			if dev.USiteID != nil {
				if _, err = repo.BatchUpdateServerUSitesStatus([]uint{*dev.USiteID}, model.USiteStatFree); err != nil {
					log.Errorf("update server_usite status failed, usite_id :%s", dev.USiteID)
				}
			}
		}

		//修改机位占用状态
		if mod.USiteID != nil {
			if _, err = repo.BatchUpdateServerUSitesStatus([]uint{*mod.USiteID}, model.USiteStatUsed); err != nil {
				log.Errorf("update server_usite status failed, usite_num :%s", row.USiteNum)
			}
		}
		//模拟一条装机记录，用于保存IP,OS等信息
		ds := model.DeviceSetting{
			SN:              row.SN,
			InstallType:     model.InstallationPXE,
			Status:          model.InstallStatusSucc,
			InstallProgress: 1.0,
			IntranetIP:      row.IntranetIP,
			NeedExtranetIP:  model.NO,
			NeedIntranetIPv6: model.NO,
			NeedExtranetIPv6: model.NO,
		}
		//再查询下是否存在
		dsOrigin, err := repo.GetDeviceSettingBySN(row.SN)
		if err != nil {
			log.Infof("get device setting by sn:%s err%v", row.SN, err)
		} else if dsOrigin != nil {
			ds.Model = dsOrigin.Model //这样就支持更新了
			//释放之前的IP信息
			if _, err = repo.ReleaseIP(row.SN, model.IPScopeIntranet); err != nil {
				log.Errorf("release origin intranet ip by sn:%s fail", row.SN)
			}
			if dsOrigin.ExtranetIP != "" {
				if _, err = repo.ReleaseIP(row.SN, model.IPScopeExtranet); err != nil {
					log.Errorf("release origin intranet ip by sn:%s fail", row.SN)
				}
			}
			ds.Updater = reqData.LoginName
		} else {
			ds.Creator = reqData.LoginName
		}

		if sysTpl, err := repo.GetSystemTemplateByName(row.OS); err == nil {
			ds.SystemTemplateID = sysTpl.ID
		} else {
			log.Warnf("get system template by name:%s err：%v", row.OS, err)
		}
		if row.ExtranetIP != "" {
			ds.ExtranetIP = row.ExtranetIP
			ds.NeedExtranetIP = model.YES
		}

		if err = repo.SaveDeviceSetting(&ds); err != nil {
			log.Errorf("save simulate device setting err,%v", err)
			return err
		}

		//把IP标记为已分配
		//支持多IP场景：如果多IP，则格式为英文逗号（,）分隔
		for _, ip := range strings.Split(row.IntranetIP, commaSep) {
			if err = repo.AssignIPByIP(row.SN, model.Intranet, ip); err != nil {
				log.Errorf("assign intranet ip :%s to sn:%s err :%v", ip, row.SN, err)
				return err
			}
		}
		if row.ExtranetIP != "" {
			for _, ip := range strings.Split(row.ExtranetIP, commaSep) {
				if err = repo.AssignIPByIP(row.SN, model.Extranet, ip); err != nil {
					log.Errorf("assign extranet ip :%s to sn:%s err :%v", ip, row.SN, err)
					return err
				}
			}
		}

		// 若dev已经存在，且密码不为默认出厂密码，则不予修改
		if isExits && checkOriginPassword(log, conf, mod) {
			continue
		}
		devices = append(devices, mod)

	}

	// 导入设备成功后，批量修改密码
	go batchUpdateOOBPassword(log, repo, conf, devices)

	return nil
}

//checkOriginPassword 校验密码是否与出厂密码一致
func checkOriginPassword(log logger.Logger, conf *config.Config, device *model.Device) bool {
	passwd, err := utils.AESDecrypt(device.OOBPassword, []byte(conf.Crypto.Key))
	if err != nil {
		log.Errorf("checkOriginPassword: AESDecrypt password failure, err: %s", err.Error())
		return false
	}

	var mapper map[string][]*OOBUser
	json.Unmarshal([]byte(OOBVendorConfig), &mapper)

	for k, users := range mapper {
		if k != device.Vendor {
			continue
		}

		for _, user := range users {
			if user.Password == string(passwd) {
				return true
			}
		}
	}

	return false
}

// batchUpdateOOBPassword 批量修改带外密码
func batchUpdateOOBPassword(log logger.Logger, repo model.Repo, conf *config.Config, devices []*model.Device) {
	for _, dev := range devices {
		var encryptedPassword string
		bTryImportedOOB := false

		ouInit := OOBUser{}
		_ = json.Unmarshal([]byte(dev.OOBInit), &ouInit)
		if ouInit.Username != "" && ouInit.Password != "" {
			encrypted, err := utils.AESEncrypt(ouInit.Password, []byte(conf.Crypto.Key))
			if err != nil {
				log.Errorf("SN: %s encrypt imported oob password fail", dev.SN)
			} else {
				bTryImportedOOB = true
				encryptedPassword = encrypted
			}
		}

		oobIP := oob.TransferHostname2IP(log, repo, dev.SN, utils.GetOOBHost(dev.SN, dev.Vendor, conf.Server.OOBDomain))
		if oobIP != "" {
			//设备表里有数据, 检查其有效性
			if dev.OOBUser != "" && dev.OOBPassword != "" {
				pwBytes, err := utils.AESDecrypt(dev.OOBPassword, []byte(conf.Crypto.Key))
				if err != nil {
					log.Errorf("SN: %s decrypt old password：%s fail", dev.OOBPassword)
					//continue
				}
				if oobPingTest(log, oobIP, dev.OOBUser, string(pwBytes), dev.OOBPassword) {
					//如果表里的数据有效，则不改。
					continue
				}
			}

			//1->没有则获取一个出厂默认的用户密码
			defaultUser, err := GetDefaultOOBUserPassword(log, repo, conf, dev.SN, dev.Vendor)
			//默认用户不用ping测试了，因为获取的时候已经逐个试过了。返回的都是有效值
			if defaultUser != nil {
				UpdateOOBPasswordBySN(log, repo, &UpdateOOBPasswordReq{
					SN:          dev.SN,
					Username:    defaultUser.Username,
					PasswordOld: defaultUser.Password,
					PasswordNew: GenPassword(),
				}, conf)
				continue //success 1
			}
			log.Infof("SN：%s get default oob user fail,err:%v", dev.SN, err)

			//2. 如果默认出厂用户也无法连通，尝试找回
			oobHis, err := FindOOBByHistory(log, repo, conf, dev.SN)
			if err == nil && oobHis != nil {
				UpdateOOBPasswordBySN(log, repo, &UpdateOOBPasswordReq{
					SN:          dev.SN,
					Username:    oobHis.Username,
					PasswordOld: oobHis.Password,
					PasswordNew: GenPassword(),
				}, conf)
				continue //success 2
			}
			log.Infof("SN:%s find back oob user password fail,err:%v", dev.SN, err)

			//3.如果没有找回，再次尝试从用户导入的字段读取
			if ouInit.Username != "" && ouInit.Password != "" {
				if bTryImportedOOB && oobPingTest(log, oobIP, ouInit.Username, ouInit.Password, encryptedPassword) {
					UpdateOOBPasswordBySN(log, repo, &UpdateOOBPasswordReq{
						SN:          dev.SN,
						Username:    ouInit.Username,
						PasswordOld: ouInit.Password,
						PasswordNew: GenPassword(),
					}, conf)
					continue //success 3
				} else {
					log.Infof("SN：%s oob imported oob user：%s,password: %s connected fail", dev.SN, ouInit.Username, encryptedPassword)
				}
			}
			log.Infof("SN:%s imported oob user empty, all update attempts fail", dev.SN)
		} else {
			log.Errorf("SN：%s get oob ip from dns fail", dev.SN)
		}

		if bTryImportedOOB {
			//客户要求，当所有来源尝试过都失败了，就把用户导入的账户密码保存到设备表里
			dev.OOBUser = ouInit.Username
			dev.OOBPassword = encryptedPassword
			isYes := model.NO
			dev.OOBAccessible = &isYes
			_, err := repo.UpdateDevice(dev)
			if err != nil {
				log.Errorf("SN:%s force write imported oob user fail，Err:%s", dev.SN, err.Error())
				return
			}
		}
	}
}

//checkLength 对导入文件中的数据做字段长度校验
func (impDevReq *ImportDevicesReq) checkLength() {
	leg := len(impDevReq.SN)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:序列号长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.Vendor)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:厂商长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.Model)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:型号长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.Usage)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:用途长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.Category)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:设备类型长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.ServerRoomName)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:机房管理单元名称长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.CabinetNum)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:机架编号长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.USiteNum)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:机位编号长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.StartedAt)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:启用时间长度为(%d)(不能为空，不能大于255)", leg)
	}
	if !strings.Contains(impDevReq.StartedAt, "-") {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + "启用时间格式须为：2019/01/01"
	}
	leg = len(impDevReq.OnShelveAt)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:上架时间长度为(%d)(不能为空，不能大于255)", leg)
	}
	if !strings.Contains(impDevReq.OnShelveAt, "-") {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + "上架时间格式须为：2019/01/01"
	}

	//硬件说明
	//RAID说明
	//带外说明
	//暂时不做检验
}

//checkLength 对导入文件中的数据做字段长度校验
func (impDevReq *ImportDevice2StoreReq) checkLength() {
	//固定资产编号可以由规则生成
	leg := len(impDevReq.SN)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:序列号长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.Vendor)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:厂商长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.Model)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:型号长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.Usage)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:用途长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.Category)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:设备类型长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.StoreRoomName)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:库房管理单元名称长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.VCabinetNum)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:虚拟货架编号长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.Owner)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:负责人长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.MaintenanceServiceDateBegin)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:维保服务起始日期长度为(%d)(不能为空，不能大于255)", leg)
	}
	if impDevReq.MaintenanceMonths == 0 || impDevReq.MaintenanceMonths > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:保修期（月数）长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.IsRental)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验: 是否租赁 长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.MaintenanceServiceProvider)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验: 维保服务供应商 长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.StartedAt)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:启用时间长度为(%d)(不能为空，不能大于255)", leg)
	}
	if !strings.Contains(impDevReq.StartedAt, "-") {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + "启用时间格式须为：2019/01/01"
	}
	leg = len(impDevReq.OnShelveAt)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:上架时间长度为(%d)(不能为空，不能大于255)", leg)
	}
	if !strings.Contains(impDevReq.OnShelveAt, "-") {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + "上架时间格式须为：2019/01/01"
	}
}

//checkLength 对导入文件中的数据做字段长度校验
func (impDevReq *ImportStockDevicesReq) checkLength() {
	leg := len(impDevReq.FixedAssetNum)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:固定资产编号长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.SN)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:序列号长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.Vendor)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:厂商长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.Model)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:型号长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.Usage)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:用途长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.Category)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:设备类型长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.ServerRoomName)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:机房管理单元名称长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.CabinetNum)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:机架编号长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.USiteNum)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:机位编号长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.StartedAt)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:启用时间长度为(%d)(不能为空，不能大于255)", leg)
	}
	if !strings.Contains(impDevReq.StartedAt, "-") {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + "启用时间格式须为：2019/01/01"
	}
	leg = len(impDevReq.OnShelveAt)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:上架时间长度为(%d)(不能为空，不能大于255)", leg)
	}
	if !strings.Contains(impDevReq.OnShelveAt, "-") {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + "上架时间格式须为：2019/01/01"
	}

	//硬件说明
	//RAID说明
	//带外说明
	//暂时不做检验
	leg = len(impDevReq.IntranetIP)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:内网IP长度为(%d)(不能为空，不能大于255)", leg)
	}
	//leg = len(impDevReq.ExtranetIP) 外网可以缺省
	//if leg == 0 || leg > 255 {
	//	var br string
	//	if impDevReq.ErrMsgContent != "" {
	//		br = "<br />"
	//	}
	//	impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:外网IP长度为(%d)(不能为空，不能大于255)", leg)
	//}
	leg = len(impDevReq.OS)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:操作系统长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.OperationStatus)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:运营状态长度为(%d)(不能为空，不能大于255)", leg)
	}
}

//validate 对导入文件中的数据做基本验证
func (impDevReq *ImportDevicesReq) validate(repo model.Repo) error {
	//用途枚举值
	//devUsageEnum := []string{model.DevUsageDB, model.DevUsageAPP, model.DevUsageTGW, model.DevUsageCVM,
	//	model.DevUsageStorage, model.DevUsageContainer, model.DevUsageBDP, model.DevUsageSpecialDev, model.DevUsageOther}
	//if !ValidateEnum(impDevReq.Usage, devUsageEnum) {
	//	var br string
	//	if impDevReq.ErrMsgContent != "" {
	//		br = "<br />"
	//	}
	//	impDevReq.ErrMsgContent += br + fmt.Sprintf("用途字段值:%s非法，必须是:%v", impDevReq.Usage, devUsageEnum)
	//}
	//机房校验
	srs, err := repo.GetServerRoomByName(impDevReq.ServerRoomName)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound || srs == nil {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("机房名(%s)不存在", impDevReq.ServerRoomName)
	} else if srs != nil {
		impDevReq.idcID = srs.IDCID
		impDevReq.serverRoomID = srs.ID
	}

	//机架
	cabinet, err := repo.GetServerCabinetByNumber(impDevReq.serverRoomID, impDevReq.CabinetNum)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound || cabinet == nil {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("机架编号(%s)不存在", impDevReq.CabinetNum)
	} else {
		impDevReq.cabinetID = cabinet.ID
	}
	//机位
	uSite, err := repo.GetServerUSiteByNumber(impDevReq.cabinetID, impDevReq.USiteNum)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound || uSite == nil {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("机位编号(%s)不存在, 机架编号(%s)", impDevReq.USiteNum, impDevReq.CabinetNum)
	} else {
		// 设备已存在，说明是先进入bootos或者数据已经导入过，重新导入刷新，此时的机位如果没有变化，
		// 则不判断机位占用状态，因为就是被自己占着
		dev, _ := repo.GetDeviceBySN(impDevReq.SN)
		if !CheckUSiteFree(repo, uSite.ID, dev) {
			//log.Errorf("机位:%s 被占用或被禁用", row.USiteNum)
			//return errors.New("机位:" + row.USiteNum + "被占用或被禁用")
			var br string
			if impDevReq.ErrMsgContent != "" {
				br = "<br />"
			}
			impDevReq.ErrMsgContent += br + fmt.Sprintf("机位编号(%s), 机架编号(%s)被占用或不可用", impDevReq.USiteNum, impDevReq.CabinetNum)
		} else {
			impDevReq.uSiteID = uSite.ID
		}
	}
	if impDevReq.OOBInit != "" {
		if !strings.Contains(impDevReq.OOBInit, ":") {
			var br string
			if impDevReq.ErrMsgContent != "" {
				br = "<br />"
			}
			impDevReq.ErrMsgContent += br + fmt.Sprintf("带外用户密码初始值:%s格式不正确，应以':'分隔", impDevReq.OOBInit)
		}
	}
	if impDevReq.OrderNumber != "" {
		order, err := repo.GetOrderByNumber(impDevReq.OrderNumber)
		if err != nil {
			var br string
			if impDevReq.ErrMsgContent != "" {
				br = "<br />"
			}
			impDevReq.ErrMsgContent += br + fmt.Sprintf("订单(订单号:%s)不存在", impDevReq.OrderNumber)
		}
		if order != nil {
			if order.IDCID != impDevReq.idcID {
				var br string
				if impDevReq.ErrMsgContent != "" {
					br = "<br />"
				}
				impDevReq.ErrMsgContent += br + fmt.Sprintf("订单(订单号:%s)与设备（SN:%s）不属于同一个数据中心", impDevReq.OrderNumber, impDevReq.SN)
			}
			if order.Category != impDevReq.Category {
				var br string
				if impDevReq.ErrMsgContent != "" {
					br = "<br />"
				}
				impDevReq.ErrMsgContent += br + fmt.Sprintf("订单(订单号:%s)的设备类型（%s） 与 设备（SN:%s）的设备类型（%s）不匹配", impDevReq.OrderNumber, order.Category, impDevReq.SN, impDevReq.Category)
			}
		}
	}

	return nil
}

//validate 对导入文件中的数据做基本验证
func (impDevReq *ImportDevice2StoreReq) validate(repo model.Repo) error {
	//用途枚举值, 放开限制

	//机房校验
	srs, err := repo.GetStoreRoomByName(impDevReq.StoreRoomName)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound || srs == nil {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("库房(%s)不存在", impDevReq.StoreRoomName)
	} else if srs != nil {
		impDevReq.idcID = srs.IDCID
		impDevReq.storeRoomID = srs.ID
	}

	//虚拟货架
	cabinet, err := repo.GetVirtualCabinets(&model.VirtualCabinet{
		StoreRoomID: impDevReq.storeRoomID,
		Number:      impDevReq.VCabinetNum,
	}, nil, nil)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound || len(cabinet) <= 0 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("虚拟货架(%s)不存在", impDevReq.VCabinetNum)
	} else if len(cabinet) > 0 {
		impDevReq.vcabinetID = cabinet[0].ID
	}
	if impDevReq.OOBInit != "" {
		if !strings.Contains(impDevReq.OOBInit, ":") {
			var br string
			if impDevReq.ErrMsgContent != "" {
				br = "<br />"
			}
			impDevReq.ErrMsgContent += br + fmt.Sprintf("带外用户密码初始值:%s格式不正确，应以':'分隔", impDevReq.OOBInit)
		}
	}
	if impDevReq.IsRental != model.NO && impDevReq.IsRental != model.YES {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("是否租赁(%s)不合法，应为 yes or no", impDevReq.IsRental)
	}
	if impDevReq.OrderNumber != "" {
		order, err := repo.GetOrderByNumber(impDevReq.OrderNumber)
		if err != nil || order == nil {
			var br string
			if impDevReq.ErrMsgContent != "" {
				br = "<br />"
			}
			impDevReq.ErrMsgContent += br + fmt.Sprintf("订单(订单号:%s)不存在", impDevReq.OrderNumber)
		} else if order != nil {
			if order.IDCID != impDevReq.idcID {
				var br string
				if impDevReq.ErrMsgContent != "" {
					br = "<br />"
				}
				impDevReq.ErrMsgContent += br + fmt.Sprintf("订单(订单号:%s)与设备（SN:%s）不属于同一个数据中心", impDevReq.OrderNumber, impDevReq.SN)
			}
			if order.Category != impDevReq.Category {
				var br string
				if impDevReq.ErrMsgContent != "" {
					br = "<br />"
				}
				impDevReq.ErrMsgContent += br + fmt.Sprintf("订单(订单号:%s)的设备类型（%s） 与 设备（SN:%s）的设备类型（%s）不匹配", impDevReq.OrderNumber, order.Category, impDevReq.SN, impDevReq.Category)
			}
		}
	}
	return nil
}

//validate 对导入文件中的数据做基本验证
func (impDevReq *ImportStockDevicesReq) validate(repo model.Repo) error {
	//用途枚举值
	//devUsageEnum := []string{model.DevUsageDB, model.DevUsageAPP, model.DevUsageTGW, model.DevUsageCVM,
	//	model.DevUsageStorage, model.DevUsageContainer, model.DevUsageBDP, model.DevUsageSpecialDev, model.DevUsageOther}
	//if !ValidateEnum(impDevReq.Usage, devUsageEnum) {
	//	var br string
	//	if impDevReq.ErrMsgContent != "" {
	//		br = "<br />"
	//	}
	//	impDevReq.ErrMsgContent += br + fmt.Sprintf("用途字段值:%s非法，必须是:%v", impDevReq.Usage, devUsageEnum)
	//}
	//机房校验
	srs, err := repo.GetServerRoomByName(impDevReq.ServerRoomName)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound || srs == nil {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("机房名(%s)不存在", impDevReq.ServerRoomName)
	} else {
		impDevReq.idcID = srs.IDCID
		impDevReq.serverRoomID = srs.ID
	}

	//机架
	cabinet, err := repo.GetServerCabinetByNumber(impDevReq.serverRoomID, impDevReq.CabinetNum)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound || cabinet == nil {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("机架编号(%s)不存在", impDevReq.CabinetNum)
	} else {
		impDevReq.cabinetID = cabinet.ID
	}
	//机位
	uSite, err := repo.GetServerUSiteByNumber(impDevReq.cabinetID, impDevReq.USiteNum)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound || uSite == nil {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("机位编号(%s)不存在, 机架编号(%s)", impDevReq.USiteNum, impDevReq.CabinetNum)
	} else {
		// 设备已存在，说明是先进入bootos或者数据已经导入过，重新导入刷新，此时的机位如果没有变化，
		// 则不判断机位占用状态，因为就是被自己占着
		dev, _ := repo.GetDeviceBySN(impDevReq.SN)
		if !CheckUSiteFree(repo, uSite.ID, dev) {
			//log.Errorf("机位:%s 被占用或被禁用", row.USiteNum)
			//return errors.New("机位:" + row.USiteNum + "被占用或被禁用")
			var br string
			if impDevReq.ErrMsgContent != "" {
				br = "<br />"
			}
			impDevReq.ErrMsgContent += br + fmt.Sprintf("机位编号(%s), 机架编号(%s)被占用或不可用", impDevReq.USiteNum, impDevReq.CabinetNum)
		} else {
			impDevReq.uSiteID = uSite.ID
		}
	}
	if impDevReq.OOBInit != "" {
		if !strings.Contains(impDevReq.OOBInit, ":") {
			var br string
			if impDevReq.ErrMsgContent != "" {
				br = "<br />"
			}
			impDevReq.ErrMsgContent += br + fmt.Sprintf("带外用户密码初始值:%s格式不正确，应以':'分隔", impDevReq.OOBInit)
		}
	}
	if isValidIP(impDevReq.IntranetIP) {
		for _, ip := range strings.Split(impDevReq.IntranetIP, commaSep) {
			if ipDB, err := repo.GetIPs(&model.IPPageCond{IP: ip}, nil, nil); len(ipDB) < 1 || err != nil {
				var br string
				if impDevReq.ErrMsgContent != "" {
					br = "<br />"
				}
				impDevReq.ErrMsgContent += br + fmt.Sprintf("内网IP：%s不存在", ip)
				//return err
			} else if len(ipDB) > 0 {
				if ipDB[0].IsUsed == model.YES && ipDB[0].SN != impDevReq.SN {
					var br string
					if impDevReq.ErrMsgContent != "" {
						br = "<br />"
					}
					impDevReq.ErrMsgContent += br + fmt.Sprintf("内网IP：%s已被占用/禁用，关联SN:%s", ip, ipDB[0].SN)
					//return nil
				}
				if ipDB[0].Scope != nil && *ipDB[0].Scope != model.IPScopeIntranet {
					var br string
					if impDevReq.ErrMsgContent != "" {
						br = "<br />"
					}
					impDevReq.ErrMsgContent += br + fmt.Sprintf("IP：%s非内网IP", ip)
					//return nil
				}
			}
		}
	} else {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("内网IP：%s格式不正确，多IP以英文逗号分隔", impDevReq.IntranetIP)
	}
	if impDevReq.ExtranetIP != "" {
		if isValidIP(impDevReq.ExtranetIP) {
			for _, ip := range strings.Split(impDevReq.ExtranetIP, commaSep) {
				if ipDB, err := repo.GetIPs(&model.IPPageCond{IP: ip}, nil, nil); len(ipDB) < 1 || err != nil {
					var br string
					if impDevReq.ErrMsgContent != "" {
						br = "<br />"
					}
					impDevReq.ErrMsgContent += br + fmt.Sprintf("外网IP：%s不存在", ip)
					//return err
				} else if len(ipDB) > 0 {
					if ipDB[0].IsUsed == model.YES && ipDB[0].SN != impDevReq.SN {
						var br string
						if impDevReq.ErrMsgContent != "" {
							br = "<br />"
						}
						impDevReq.ErrMsgContent += br + fmt.Sprintf("外网IP：%s已被占用/禁用，关联SN:%s", ip, ipDB[0].SN)
						//return nil
					}
					if ipDB[0].Scope != nil && *ipDB[0].Scope != model.IPScopeExtranet {
						var br string
						if impDevReq.ErrMsgContent != "" {
							br = "<br />"
						}
						impDevReq.ErrMsgContent += br + fmt.Sprintf("IP：%s非外网IP", ip)
						//return nil
					}
				}
			}
		} else {
			var br string
			if impDevReq.ErrMsgContent != "" {
				br = "<br />"
			}
			impDevReq.ErrMsgContent += br + fmt.Sprintf("外网IP：%s格式不正确，多IP以英文逗号分隔", impDevReq.ExtranetIP)
		}
	}
	if _, err := repo.GetSystemTemplateByName(impDevReq.OS); err != nil {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("操作系统(对应于系统安装模板名称)：%s不存在", impDevReq.OS)
	}
	//运行状态枚举值
	opStatusEnum := []string{"运行中(需告警)", "运行中(无需告警)", "已上架" /*"重装中", "搬迁中", "待退役", "已退役", "待部署", "回收中"*/}
	if !ValidateEnum(impDevReq.OperationStatus, opStatusEnum) {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("运营状态值:%s非法，必须是:%v", impDevReq.OperationStatus, opStatusEnum)
	}
	return nil
}

//判断是否合规的IP格式，单IP或以英文逗号分隔的多IP
func isValidIP(ipStr string) bool {
	return strings.Count(ipStr, ".") == 3 || (strings.Count(ipStr, ".") > 3 && strings.Contains(ipStr, commaSep))
}

// Validate 对修改的数据做基本校验
func (reqData *UpdateDevicesReq) Validate(request *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(request.Context())
	dev, err := repo.GetDeviceByID(reqData.ID)
	if err != nil {
		errs.Add([]string{"设备ID"}, binding.RequiredError, fmt.Sprintf("设备ID：%d不存在", reqData.ID))
		return errs
	}
	if d, _ := repo.GetDeviceBySN(reqData.SN); d != nil && d.ID != reqData.ID {
		errs.Add([]string{"设备SN"}, binding.RequiredError, fmt.Sprintf("设备SN：%s已存在", reqData.SN))
		return errs
	}
	if d, _ := repo.GetDeviceByFixAssetNumber(reqData.FixedAssetNum); d != nil && d.ID != reqData.ID {
		errs.Add([]string{"设备固资编号"}, binding.RequiredError, fmt.Sprintf("设备固资编号：%s已存在", reqData.FixedAssetNum))
		return errs
	}
	_, err = repo.GetIDCByID(reqData.IDCID)
	if err != nil {
		errs.Add([]string{"数据中心ID"}, binding.RequiredError, fmt.Sprintf("数据中心ID：%d不存在", reqData.IDCID))
		return errs
	}
	serverRoom, err := repo.GetServerRoomByID(reqData.ServerRoomID)
	if err != nil {
		errs.Add([]string{"机房ID"}, binding.RequiredError, fmt.Sprintf("机房ID：%d不存在", reqData.ServerRoomID))
		return errs
	}
	if serverRoom.IDCID != reqData.IDCID {
		errs.Add([]string{"机房ID"}, binding.RequiredError, fmt.Sprintf("机房ID：%d不属于数据中心ID:%d", reqData.ServerRoomID, reqData.IDCID))
		return errs
	}
	cabinet, err := repo.GetServerCabinetByID(reqData.CabinetID)
	if err != nil {
		errs.Add([]string{"机架ID"}, binding.RequiredError, fmt.Sprintf("机架ID：%d不存在", reqData.CabinetID))
		return errs
	}
	if cabinet.ServerRoomID != reqData.ServerRoomID {
		errs.Add([]string{"机架ID"}, binding.RequiredError, fmt.Sprintf("机架ID：%d不属于机房ID:%d", reqData.CabinetID, reqData.ServerRoomID))
		return errs
	}
	usite, err := repo.GetServerUSiteByID(reqData.USiteID)
	if err != nil {
		errs.Add([]string{"机位ID"}, binding.RequiredError, fmt.Sprintf("机位ID：%d不存在", reqData.USiteID))
		return errs
	}
	if usite.ServerCabinetID != reqData.CabinetID {
		errs.Add([]string{"机位ID"}, binding.RequiredError, fmt.Sprintf("机位ID：%d不属于机架ID:%d", reqData.USiteID, reqData.CabinetID))
		return errs
	}
	if !CheckUSiteFree(repo, usite.ID, dev) {
		errs.Add([]string{"机位ID"}, binding.RequiredError, fmt.Sprintf("机位ID：%d已被占用或不可用", reqData.USiteID))
		return errs
	}
	return nil
}

// GetDevicePage 查询物理机分页列表
func GetDevicePage(log logger.Logger, repo model.Repo, conf *config.Config, reqData *DevicePageReq) (pg *page.Page, err error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 100 {
		reqData.PageSize = 20
	}
	if reqData.Page < 0 {
		reqData.Page = 1
	}

	var cond model.CombinedDeviceCond
	cond.IDCID = strings2.Multi2UintSlice(reqData.IDCID)
	cond.ServerRoomID = strings2.Multi2UintSlice(reqData.ServerRoomID)
	cond.PhysicalArea = reqData.PhysicalArea
	cond.ServerCabinet = reqData.ServerCabinet
	cond.ServerRoomName = reqData.ServerRoomName
	cond.ServerUsiteNumber = reqData.ServerUSiteNumber
	cond.FixedAssetNumber = reqData.FixedAssetNum
	cond.SN = reqData.SN
	cond.Vendor = reqData.Vendor
	cond.DevModel = reqData.Model
	cond.Usage = reqData.Usage
	// 物理机-待部署列表仅获取待部署状态的设备类型
	if reqData.CategoryPreDeploy != "" {
		cond.Category = reqData.CategoryPreDeploy
	} else {
		cond.Category = reqData.Category
	}
	cond.OperationStatus = reqData.OperationStatus
	cond.PreDeployed = reqData.PreDeployed
	cond.HardwareRemark = reqData.HardwareRemark
	cond.IntranetIP = reqData.IntranetIP
	cond.ExtranetIP = reqData.ExtranetIP
	cond.IP = reqData.IP
	cond.OOBAccessible = reqData.OOBAccessible

	totalRecords, err := repo.CountCombinedDevices(&cond)
	if err != nil {
		return nil, err
	}

	pager := page.NewPager(reflect.TypeOf(&DevicePageResp{}), reqData.Page, reqData.PageSize, totalRecords)
	orderBy := []model.OrderByPair{
		/*{Name: "device.updated_at", Direction: model.DESC},*/
		{Name: "device.idc_id", Direction: model.ASC},
		{Name: "device.server_room_id", Direction: model.ASC},
		{Name: "device.server_cabinet_id", Direction: model.ASC},
		{Name: "device.id", Direction: model.ASC}} //MYSQL VERSION >=5.6 ，若不增加索引字段作为排序值，可能导致分页（LIMIT&OFFSET）数据重复
	items, err := repo.GetCombinedDevices(&cond, orderBy, pager.BuildLimiter())
	if err != nil {
		return nil, err
	}
	for i := range items {
		item, err := conver2DevicePagesResp(log, repo, conf, &items[i].Device)
		if err != nil {
			log.Error(err)
			//return nil, err
		}

		pager.AddRecords(item)
	}
	return pager.BuildPage(), nil
}

//GetExportDevices 获取导出物理机列表
func GetExportDevices(log logger.Logger, repo model.Repo, conf *config.Config, reqData *DevicePageReq) (devs []*DevicePageResp, err error) {
	var cond model.CombinedDeviceCond
	cond.ID = strings2.Multi2UintSlice(reqData.ID)
	cond.IDCID = strings2.Multi2UintSlice(reqData.IDCID)
	cond.ServerRoomID = strings2.Multi2UintSlice(reqData.ServerRoomID)
	cond.PhysicalArea = reqData.PhysicalArea
	cond.ServerCabinet = reqData.ServerCabinet
	cond.FixedAssetNumber = reqData.FixedAssetNum
	cond.SN = reqData.SN
	cond.Vendor = reqData.Vendor
	cond.DevModel = reqData.Model
	cond.Usage = reqData.Usage
	cond.Category = reqData.Category
	cond.OperationStatus = reqData.OperationStatus
	cond.PreDeployed = reqData.PreDeployed
	cond.IntranetIP = reqData.IntranetIP
	cond.ExtranetIP = reqData.ExtranetIP
	cond.OOBAccessible = reqData.OOBAccessible

	items, err := repo.GetCombinedDevices(&cond, model.OneOrderBy("id", model.DESC), nil)
	if err != nil {
		return nil, err
	}
	devs = make([]*DevicePageResp, 0, len(items))
	for i := range items {
		item, err := conver2DevicePagesResp(log, repo, conf, &items[i].Device)
		if err != nil {
			log.Error(err)
			//return nil, err
		}
		devs = append(devs, item)
	}
	return
}

func ConvertPowerStatus(powerStatus string) string {
	if powerStatus == model.PowerStatusOn {
		return "开电"
	}
	return "关电"
}

//conver2DevicePagesResp 将model层转view层返回结构
func conver2DevicePagesResp(log logger.Logger, repo model.Repo, conf *config.Config, device *model.Device) (resp *DevicePageResp, err error) {
	resp = &DevicePageResp{
		ID:            device.ID,
		FixedAssetNum: device.FixedAssetNumber,
		SN:            device.SN,
		Vendor:        device.Vendor,
		Model:         device.DevModel,
		Arch:          device.Arch,
		Usage:         device.Usage,
		Category:      device.Category,
		OOBIP:         device.OOBIP,
		OOBUser:       device.OOBUser,
		PowerStatus:   device.PowerStatus,
		IDC:           &IDCSimplify{ID: device.IDCID},
		ServerRoom:    &ServerRoomSimplify{ID: device.ServerRoomID},
		ServerCabinet: &ServerCabinetSimplify{ID: device.CabinetID},
		StoreRoom:     &StoreRoomSimplify{ID: device.StoreRoomID},
		VCabinet:      &VCabinetSimplify{ID: device.VCabinetID},		
		//ServerUSite:     &ServerUSiteSimplify{ID: *device.USiteID},
		HardwareRemark:  device.HardwareRemark,
		RAIDRemark:      device.RAIDRemark,
		OperationStatus: device.OperationStatus,
		StartedAt:       times.ISOTime(device.StartedAt).ToDateStr(),
		OnShelveAt:      times.ISOTime(device.OnShelveAt).ToDateStr(),
		CreatedAt:       times.ISOTime(device.CreatedAt).ToDateStr(),
		UpdatedAt:       times.ISOTime(device.UpdatedAt).ToDateStr(),
		OriginNode:      device.OriginNode,
		OriginNodeIP:    device.OriginNodeIP,
		OrderNumber:     device.OrderNumber,
	}
	if device.OOBAccessible != nil {
		resp.OOBAccessible = *device.OOBAccessible
	}
	if device.USiteID != nil {
		resp.ServerUSite = &ServerUSiteSimplify{ID: *device.USiteID}
	}
	if device.OOBPassword != "" {
		password, err := utils.AESDecrypt(device.OOBPassword, []byte(conf.Crypto.Key))
		if err != nil {
			log.Errorf("解密失败，密码为[%s]", device.OOBPassword)
			return nil, err
		}
		resp.OOBPassword = string(password)
	}
	//idc
	if idc, err := repo.GetIDCByID(device.IDCID); err == nil {
		resp.IDC.Name = idc.Name
	}

	if device.StoreRoomID != 0 {
		//库房
		if room, err := repo.GetStoreRoomByID(device.StoreRoomID); err == nil && room != nil {
			resp.StoreRoom.Name = room.Name
		}
		//货架
		if cabinet, err := repo.GetVirtualCabinetByID(device.VCabinetID); err == nil && cabinet != nil {
			resp.VCabinet.Number = cabinet.Number
		}
	}
	if device.ServerRoomID != 0 {
		//机房管理单元
		if room, err := repo.GetServerRoomByID(device.ServerRoomID); err == nil {
			resp.ServerRoom.Name = room.Name
		}
		//机架
		if cabinet, err := repo.GetServerCabinetByID(device.CabinetID); err == nil {
			resp.ServerCabinet.Number = cabinet.Number
		}
		//机位
		if device.USiteID != nil {
			if u, err := repo.GetServerUSiteByID(*device.USiteID); err == nil {
				resp.ServerUSite.Number = u.Number
				resp.ServerUSite.PhysicalArea = u.PhysicalArea
			}
		}	
		resp.TOR, _ = repo.GetTORBySN(device.SN)
		//装机参数
		devSett, _ := repo.GetDeviceSettingBySN(device.SN)
		if devSett != nil {
			resp.IntranetIP = devSett.IntranetIP
			resp.ExtranetIP = devSett.ExtranetIP
			if devSett.InstallType == model.InstallationPXE {
				sysTpl, _ := repo.GetSystemTemplateBySN(devSett.SN)
				if sysTpl != nil {
					resp.OS = sysTpl.Name
				}
			} else if devSett.InstallType == model.InstallationImage {
				imageTpl, _ := repo.GetImageTemplateBySN(devSett.SN)
				if imageTpl != nil {
					resp.OS = imageTpl.Name
				}
			}
		}
	}
	return
}

func GetDeviceQuerys(log logger.Logger, repo model.Repo, param string) (*model.DeviceQueryParamResp, error) {
	return repo.GetDeviceQuerys(param)
}

// GetDeviceBySN 查询采集到的设备信息
func GetDeviceBySN(log logger.Logger, repo model.Repo, sn string) (*device.Device, error) {
	m, err := repo.GetDeviceBySN(sn)
	if err != nil {
		return nil, err
	}

	var dev device.Device
	dev.SN = m.SN
	dev.Vendor = m.Vendor
	dev.Model = m.DevModel
	dev.Arch = m.Arch
	dev.BootOSIP = m.BootOSIP
	dev.BootOSMac = m.BootOSMac
	dev.NICDevice = m.NICDevice

	if m.CPU != "" {
		if err = json.Unmarshal([]byte(m.CPU), &dev.CPU); err != nil {
			log.Error(err)
			return nil, err
		}
	}

	if m.Memory != "" {
		if err = json.Unmarshal([]byte(m.Memory), &dev.Memory); err != nil {
			log.Error(err)
			return nil, err
		}
	}

	if m.Disk != "" {
		if err = json.Unmarshal([]byte(m.Disk), &dev.Disk); err != nil {
			log.Error(err)
			return nil, err
		}
	}

	if m.DiskSlot != "" {
		if err = json.Unmarshal([]byte(m.DiskSlot), &dev.DiskSlot); err != nil {
			log.Error(err)
			return nil, err
		}
	}

	if m.NIC != "" {
		if err = json.Unmarshal([]byte(m.NIC), &dev.NIC); err != nil {
			log.Error(err)
			return nil, err
		}
	}

	if m.Motherboard != "" {
		if err = json.Unmarshal([]byte(m.Motherboard), &dev.Motherboard); err != nil {
			log.Error(err)
			return nil, err
		}
	}

	if m.RAID != "" {
		if err = json.Unmarshal([]byte(m.RAID), &dev.RAID); err != nil {
			log.Error(err)
			return nil, err
		}
	}

	if m.OOB != "" {
		if err = json.Unmarshal([]byte(m.OOB), &dev.OOB); err != nil {
			log.Error(err)
			return nil, err
		}
	}

	if m.BIOS != "" {
		if err = json.Unmarshal([]byte(m.BIOS), &dev.BIOS); err != nil {
			log.Error(err)
			return nil, err
		}
	}

	if m.Fan != "" {
		if err = json.Unmarshal([]byte(m.Fan), &dev.Fan); err != nil {
			log.Error(err)
			return nil, err
		}
	}

	if m.Power != "" {
		if err = json.Unmarshal([]byte(m.Power), &dev.Power); err != nil {
			log.Error(err)
			return nil, err
		}
	}

	if m.HBA != "" {
		if err = json.Unmarshal([]byte(m.HBA), &dev.HBA); err != nil {
			log.Error(err)
			return nil, err
		}
	}

	if m.PCI != "" {
		if err = json.Unmarshal([]byte(m.PCI), &dev.PCI); err != nil {
			log.Error(err)
			return nil, err
		}
	}

	// TODO 如何呈现交换机(包含在LLDP信息中)信息
	if m.LLDP != "" {
		if err = json.Unmarshal([]byte(m.LLDP), &dev.LLDP); err != nil {
			log.Error(err)
			return nil, err
		}
	}

	if m.Extra != "" {
		if err = json.Unmarshal([]byte(m.Extra), &dev.Extra); err != nil {
			log.Error(err)
			return nil, err
		}
	}
	return &dev, nil
}

// GetCombinedDeviceBySN 返回指定SN的设备详情
func GetCombinedDeviceBySN(repo model.Repo, conf *config.Config, sn string) (dev *CombinedDevice, err error) {
	items, err := repo.GetCombinedDevices(&model.CombinedDeviceCond{SN: sn}, nil, nil)
	if err != nil {
		return nil, err
	}
	if len(items) <= 0 {
		return nil, gorm.ErrRecordNotFound
	}
	var mcd *model.CombinedDevice
	for i := range items {
		if items[i].SN == sn {
			mcd = items[i]
		}
	}

	dev, err = convert2Combind(repo, conf, mcd)
	if err != nil {
		return
	}
	//dev.BootOSIP, _ = GetIPSettingBySN(repo, sn)
	return dev, nil
}

// convert2Combind 将模型层的CombinedDevice转换成服务层的CombinedDevice。
func convert2Combind(repo model.Repo, conf *config.Config, mcd *model.CombinedDevice) (*CombinedDevice, error) {
	scd := CombinedDevice{
		DevicePageResp: DevicePageResp{
			ID:            mcd.ID,
			FixedAssetNum: mcd.FixedAssetNumber,
			SN:            mcd.SN,
			Vendor:        mcd.Vendor,
			Model:         mcd.DevModel,
			PowerStatus:   mcd.PowerStatus,
			Arch:          mcd.Arch,
			Usage:         mcd.Usage,
			Category:      mcd.Category,
			OOBIP:         mcd.OOBIP,
			OOBUser:       mcd.OOBUser,
			IDC:           &IDCSimplify{ID: mcd.IDCID},
			ServerRoom:    &ServerRoomSimplify{ID: mcd.ServerRoomID},
			ServerCabinet: &ServerCabinetSimplify{ID: mcd.CabinetID},
			//ServerUSite:     &ServerUSiteSimplify{ID: *mcd.USiteID},
			HardwareRemark:  mcd.HardwareRemark,
			RAIDRemark:      mcd.RAIDRemark,
			StartedAt:       times.ISOTime(mcd.StartedAt).ToTimeStr(),
			OnShelveAt:      times.ISOTime(mcd.OnShelveAt).ToTimeStr(),
			OperationStatus: mcd.OperationStatus,
			IntranetIP:      mcd.IntranetIP, //业务IP
			ExtranetIP:      mcd.ExtranetIP,
			IntranetIPv6:      mcd.IntranetIPv6, //业务IPv6
			ExtranetIPv6:      mcd.ExtranetIPv6,			
			OriginNode:      mcd.OriginNode,
			OriginNodeIP:    mcd.OriginNodeIP,
			OrderNumber:     mcd.OrderNumber,
			CreatedAt:       times.ISOTime(mcd.CreatedAt).ToDateStr(),
			UpdatedAt:       times.ISOTime(mcd.UpdatedAt).ToDateStr(),
		},
		DeployStatus:    mcd.DeployStatus,
		InstallProgress: mcd.InstallProgress,
		BootOSIP:        mcd.BootOSIP, //bootos临时IP
		BootOSMac:       mcd.BootOSMac,
		Remark:          mcd.Remark,
		OS:              mcd.OS,
		ImageTemplate: ImageTemplateSimplify{
			ID:   mcd.ImageTemplateID,
			Name: mcd.ImageTemplateName,
		},
		HardwareTemplate: HardwareTemplateSimplify{
			ID:   mcd.HardwareTemplateID,
			Name: mcd.HardwareTemplateName,
		},
		Hostname: mcd.Hostname,
		Inspection: InspectionSimplify{
			RunStatus: mcd.InspectionRunStatus,
			Result:    mcd.InspectionResult,
			Remark:    mcd.InspectionRemark,
		},
		DeviceLifecycleDetailPage: DeviceLifecycleDetailPage{
			AssetBelongs	 			:	mcd.AssetBelongs,	 					
			Owner			 			:	mcd.Owner,			 					
			IsRental		 			:	mcd.IsRental,
			MaintenanceServiceProvider	:	mcd.MaintenanceServiceProvider,
			MaintenanceService			:	mcd.MaintenanceService,		
			LogisticsService			:	mcd.LogisticsService,				
			MaintenanceServiceDateBegin :	times.ISOTime(mcd.MaintenanceServiceDateBegin).ToDateStr(),
			MaintenanceServiceDateEnd   :	times.ISOTime(mcd.MaintenanceServiceDateEnd).ToDateStr(),        
			MaintenanceServiceStatus	:	mcd.MaintenanceServiceStatus,
			DeviceRetiredDate       	:	times.ISOTime(mcd.DeviceRetiredDate).ToDateStr(),			
		},
	}
	if mcd.USiteID != nil {
		scd.DevicePageResp.ServerUSite = &ServerUSiteSimplify{ID: *mcd.USiteID}
	}
	if mcd.OOBUser == "" || mcd.OOBPassword == "" {
		history, _ := repo.GetLastOOBHistoryBySN(mcd.SN)
		mcd.OOBUser = history.UsernameNew
		mcd.OOBPassword = history.PasswordNew
	}
	if mcd.OOBPassword != "" {
		password, err := utils.AESDecrypt(mcd.OOBPassword, []byte(conf.Crypto.Key))
		if err != nil {
			return nil, err
		}
		scd.DevicePageResp.OOBPassword = string(password)
	}

	//idc
	if idc, err := repo.GetIDCByID(mcd.IDCID); err == nil {
		scd.DevicePageResp.IDC.Name = idc.Name
	}
	//server_room
	if room, err := repo.GetServerRoomByID(mcd.ServerRoomID); err == nil {
		scd.DevicePageResp.ServerRoom.Name = room.Name
	}
	//server_cabinet
	if cabinet, err := repo.GetServerCabinetByID(mcd.CabinetID); err == nil {
		scd.DevicePageResp.ServerCabinet.Number = cabinet.Number
	}
	//server_usite
	if mcd.USiteID != nil {
		if u, err := repo.GetServerUSiteByID(*mcd.USiteID); err == nil {
			scd.DevicePageResp.ServerUSite.Number = u.Number
		}
	}

	//store_room
	if mcd.StoreRoomID != 0 {
		if stroom, _ := repo.GetStoreRoomByID(mcd.StoreRoomID); stroom != nil {
			scd.DevicePageResp.StoreRoom = &StoreRoomSimplify{ID: mcd.StoreRoomID, Name: stroom.Name}
		}
	}

	if mcd.VCabinetID != 0 {
		if vc, _ := repo.GetVirtualCabinetByID(mcd.VCabinetID); vc != nil {
			scd.DevicePageResp.VCabinet = &VCabinetSimplify{ID: mcd.VCabinetID, Number: vc.Number}
		}
	}

	if mcd.CPU != "" {
		_ = json.Unmarshal([]byte(mcd.CPU), &scd.CPU)
	}
	if mcd.Memory != "" {
		_ = json.Unmarshal([]byte(mcd.Memory), &scd.Memory)
	}
	if mcd.Disk != "" {
		_ = json.Unmarshal([]byte(mcd.Disk), &scd.LogicDisk)
	}
	if mcd.DiskSlot != "" {
		_ = json.Unmarshal([]byte(mcd.DiskSlot), &scd.PhysicalDrive)
	}
	if mcd.NIC != "" {
		_ = json.Unmarshal([]byte(mcd.NIC), &scd.NIC)
	}
	if mcd.Motherboard != "" {
		_ = json.Unmarshal([]byte(mcd.Motherboard), &scd.Motherboard)
	}
	if mcd.RAID != "" {
		_ = json.Unmarshal([]byte(mcd.RAID), &scd.RAID)
	}
	if mcd.OOB != "" {
		_ = json.Unmarshal([]byte(mcd.OOB), &scd.OOB)
	}
	if mcd.BIOS != "" {
		_ = json.Unmarshal([]byte(mcd.BIOS), &scd.BIOS)
	}
	if mcd.Fan != "" {
		_ = json.Unmarshal([]byte(mcd.Fan), &scd.Fan)
	}
	if mcd.Power != "" {
		_ = json.Unmarshal([]byte(mcd.Power), &scd.Power)
	}
	if mcd.HBA != "" {
		_ = json.Unmarshal([]byte(mcd.HBA), &scd.HBA)
	}
	if mcd.PCI != "" {
		_ = json.Unmarshal([]byte(mcd.PCI), &scd.PCI)
	}
	if mcd.LLDP != "" {
		_ = json.Unmarshal([]byte(mcd.LLDP), &scd.LLDP)
	}
	if mcd.DeviceLifecycleDeatail.LifecycleLog != "" {
		_ = json.Unmarshal([]byte(mcd.DeviceLifecycleDeatail.LifecycleLog), &scd.DeviceLifecycleDetailPage.LifecycleLog)
	}
	return &scd, nil
}

//UpdateDeviceBySN 修改指定SN的物理设备
func UpdateDeviceBySN(log logger.Logger, repo model.Repo, conf *config.Config, reqData *UpdateDevicesReq) (*model.Device, error) {
	mod, err := repo.GetDeviceByID(reqData.ID)
	if err != nil {
		return nil, err
	}
	//记录一下原机位
	originUsiteID := mod.USiteID

	mod.FixedAssetNumber = reqData.FixedAssetNum
	mod.SN = reqData.SN
	mod.Vendor = reqData.Vendor
	mod.DevModel = reqData.Model
	mod.Usage = reqData.Usage
	mod.Category = reqData.Category
	mod.IDCID = reqData.IDCID
	mod.ServerRoomID = reqData.ServerRoomID
	mod.CabinetID = reqData.CabinetID
	mod.USiteID = &reqData.USiteID
	mod.StoreRoomID = reqData.StoreRoomID
	mod.VCabinetID = reqData.VCabinetID
	mod.HardwareRemark = reqData.HardwareRemark
	mod.RAIDRemark = reqData.RAIDRemark
	mod.StartedAt = time.Now()
	mod.StartedAt = time.Now()

	if reqData.StartedAt != "" {
		startedAt, _ := time.Parse(times.DateTimeLayout, reqData.StartedAt)
		mod.StartedAt = startedAt
	}

	if reqData.OnShelveAt != "" {
		onShelveAt, _ := time.Parse(times.DateTimeLayout, reqData.OnShelveAt)
		mod.OnShelveAt = onShelveAt
	}
	mod.OperationStatus = reqData.OperationStatus
	if mod.OOBInit == "" {
		mod.OOBInit = model.EmptyJSONObject
	}
	_, err = repo.SaveDevice(mod)
	if err != nil {
		return nil, err
	}

	//如果机位变了，修改机位占用状态,A->B，则A释放，B占用
	if originUsiteID != mod.USiteID {
		if originUsiteID != nil {
			if _, err = repo.BatchUpdateServerUSitesStatus([]uint{*originUsiteID}, model.USiteStatFree); err != nil {
				log.Errorf("update server_usite status failed, usite_id :%s", reqData.USiteID)
			}
		}
		if _, err = repo.BatchUpdateServerUSitesStatus([]uint{reqData.USiteID}, model.USiteStatUsed); err != nil {
			log.Errorf("update server_usite status failed, usite_id :%s", reqData.USiteID)
		}
	}
	return mod, nil
}

//UpdateDeviceOperationStatusBySN 修改指定SN的物理设备的部署状态
func UpdateDeviceOperationStatusBySN(log logger.Logger, repo model.Repo, conf *config.Config, reqData *UpdateDevicesOperationStatusReq) (*model.Device, error) {
	mod, err := repo.GetDeviceByFixAssetNumber(reqData.FixedAssetNum)
	if err != nil || mod == nil {
		mod, err = repo.GetDeviceBySN(reqData.SN)
		if err != nil {
			return nil, err
		}
	}

	mod.UpdatedAt = time.Now()
	mod.Updater = reqData.LoginName
	mod.OperationStatus = reqData.OperationStatus

	_, err = repo.SaveDevice(mod)
	if err != nil {
		return nil, err
	}
	return mod, nil
}

//BatchUpdateDevices 批量修改设备状态，用途等
func BatchUpdateDevices(log logger.Logger, repo model.Repo, conf *config.Config, reqData *BatchUpdateDevicesReq) (affected int, err error) {
	for _, d := range reqData.Devices {
		mod, err := repo.GetDeviceByFixAssetNumber(d.FixedAssetNum)
		if err != nil || mod == nil {
			mod, err = repo.GetDeviceBySN(d.SN)
			if err != nil {
				return affected, err
			}
		}

		mod.UpdatedAt = time.Now()
		mod.Updater = reqData.LoginName
		mod.OperationStatus = d.OperationStatus
		mod.Usage = d.Usage
		mod.HardwareRemark = d.HardwareRemark

		_, err = repo.UpdateDevice(mod)
		if err != nil {
			return affected, err
		}
		affected++
	}
	return
}

//DeleteDevices  批量删除物理设备
func DeleteDevices(log logger.Logger, repo model.Repo, reqData *DeleteDevicesReq) (totalAffected int64, err error) {
	//统计一下关联订单号的数量，用于更新订单状态
	var mOrderAmount = make(map[string]int, 0)
	//以ID为参数
	if len(reqData.IDs) != 0 {
		for _, id := range reqData.IDs {
			log.Debugf("deleting device id: %d ", id)
			dev, err := repo.GetDeviceByID(id)
			if err != nil {
				log.Errorf("delete device id: %d not exist", id)
				continue
				//return totalAffected, err
			}
			//设备的状态必须是待部署OR特殊设备
			if dev.Category != "特殊设备" && dev.OperationStatus != model.DevOperStatPreDeploy {
				log.Errorf("device SN: %s operation status: %s is not pre_deploy, can't be deleted", dev.SN, dev.OperationStatus)
				err = fmt.Errorf("[业务校验]设备SN: %s运营状态: %s不是待部署(pre_deploy)，不允许删除", dev.SN, dev.OperationStatus)
				return totalAffected, err
			}
			//释放机位资源
			if dev.USiteID != nil {
				_, err = BatchUpdateServerUSitesStatus(repo, []uint{*dev.USiteID}, model.USiteStatFree)
				if err != nil {
					log.Errorf("free usite: %d fail", dev.USiteID)
					return totalAffected, err
				}
			}
			//释放IP资源（若存在）
			if err = repo.UnassignIPsBySN(dev.SN); err != nil {
				log.Error("free IP sources fail")
				return totalAffected, err
			}
			//统计删除的设备涉及的订单以及数量
			if dev.OrderNumber != "" {
				mOrderAmount[dev.OrderNumber]++
			}
			//删除装机参数
			if _, err = repo.DeleteDeviceSettingBySN(dev.SN); err != nil {
				log.Errorf("device SN: %s delete device_setting fail", dev.SN)
				return totalAffected, err
			}
			//删除设备生命周期记录
			if _, err = repo.RemoveDeviceLifecycleBySN(dev.SN); err != nil {
				log.Errorf("device SN: %s delete device_lifecycle fail", dev.SN)
				return totalAffected, err
			}
			if centos6.IsPXEUEFI(log, repo, dev.SN) {
				_ = centos6.DropConfigurations(log, repo, dev.SN) // TODO 为支持centos6的UEFI方式安装而临时增加的逻辑，后续会删除。
			}
	
			affected, err := repo.RemoveDeviceByID(id)
			if err == nil {
				totalAffected += affected
			}
		}
	}
	// 以SN为参数
	if len(reqData.SNs) != 0 {
		for _, sn := range reqData.SNs {
			log.Debugf("deleting device sn: %s ", sn)
			dev, err := repo.GetDeviceBySN(sn)
			if err != nil {
				log.Errorf("delete device sn: %s  not exist", sn)
				continue
			}
			//标准设备的状态必须是待部署OR特殊设备
			if dev.Category != "特殊设备" && dev.OperationStatus != model.DevOperStatPreDeploy {
				log.Errorf("device SN: %s operation status: %s is not pre_deploy, can't be deleted", dev.SN, dev.OperationStatus)
				err = fmt.Errorf("[业务校验]设备SN: %s运营状态: %s不是待部署(pre_deploy)，不允许删除", dev.SN, dev.OperationStatus)
				return totalAffected, err
			}
			//特殊设备的状态必须是已上架
			if dev.Category == "特殊设备" && dev.OperationStatus != model.DevOperStatOnShelve {
				log.Errorf("special device SN: %s operation status: %s is not on_shelve, can't be deleted", dev.SN, dev.OperationStatus)
				err = fmt.Errorf("[业务校验]特殊设备SN: %s运营状态: %s不是待部署(pre_deploy)，不允许删除", dev.SN, dev.OperationStatus)
				return totalAffected, err
			}			
			//释放机位资源
			if dev.USiteID != nil {
				_, err = BatchUpdateServerUSitesStatus(repo, []uint{*dev.USiteID}, model.USiteStatFree)
				if err != nil {
					log.Errorf("free usite: %d fail", dev.USiteID)
					return totalAffected, err
				}
			}
			//释放IP资源（若存在）
			if err = repo.UnassignIPsBySN(dev.SN); err != nil {
				log.Error("free IP sources fail")
				return totalAffected, err
			}
			//统计删除的设备涉及的订单以及数量
			if dev.OrderNumber != "" {
				mOrderAmount[dev.OrderNumber]++
			}
			//删除装机参数
			if _, err = repo.DeleteDeviceSettingBySN(dev.SN); err != nil {
				log.Errorf("device SN: %s delete device_setting fail", dev.SN)
				return totalAffected, err
			}
			//删除设备生命周期记录
			if _, err = repo.RemoveDeviceLifecycleBySN(dev.SN); err != nil {
				log.Errorf("device SN: %s delete device_lifecycle fail", dev.SN)
				return totalAffected, err
			}	
			if centos6.IsPXEUEFI(log, repo, dev.SN) {
				_ = centos6.DropConfigurations(log, repo, dev.SN) // TODO 为支持centos6的UEFI方式安装而临时增加的逻辑，后续会删除。
			}
	
			affected, err := repo.RemoveDeviceBySN(sn)
			if err == nil {
				totalAffected += affected
			}
		}
	}
	//更新关联的订单到货数量和订单状态
	for orderNum, deleteCount := range mOrderAmount {
		if err = UpdateOrderByDelete(log, repo, orderNum, deleteCount); err != nil {
			log.Errorf("faile to update order(%s) when deleting pre_deploy devices ", orderNum)
			err = fmt.Errorf("删除待部署设备更新订单%s失败", orderNum)
			return totalAffected, err
		}
	}

	return
}

// GetDevicePowerStatusBySN 根据SN获取设置开电状态
func GetDevicePowerStatusBySN(log logger.Logger, repo model.Repo, conf *config.Config, sn string) (powerStatus string, err error) {
	dev, err := repo.GetDeviceBySN(sn)
	if err != nil {
		return "", err
	}

	domain := conf.Server.OOBDomain
	key := conf.Crypto.Key

	if dev.OOBPassword == "" || dev.OOBUser == "" {
		log.Warnf("设备带外用户或密码为空，尝试找回，[SN:%s]", sn)
		history, err := repo.GetLastOOBHistoryBySN(sn)
		if err != nil {
			log.Errorf("find back oob history by sn:%s fail,%s", sn, err.Error())
			return fmt.Sprintf("find back oob history by sn:%s fail", sn),
				errors.New("带外信息(用户密码)不存在")
		}
		dev.OOBUser = history.UsernameNew
		dev.OOBPassword = history.PasswordNew
	}

	password, err := utils.AESDecrypt(dev.OOBPassword, []byte(key))
	if err != nil {
		return "", fmt.Errorf("解密失败，%s", dev.OOBPassword)
	}

	isOn, err := OOBPowerStatus(log, oob.TransferHostname2IP(log, repo, dev.SN, utils.GetOOBHost(dev.SN, dev.Vendor, domain)),
		dev.OOBUser, string(password), dev.OOBPassword)
	if err != nil {
		return "", err
	}

	powerStatus = model.PowerStatusOff
	if isOn {
		powerStatus = model.PowerStatusOn
	}

	if dev.PowerStatus != powerStatus {
		dev.PowerStatus = powerStatus
		repo.SaveDevice(dev)
	}

	return
}

// CollectedDevice 采集到的设备信息结构体
type CollectedDevice struct {
	device.Device
	OriginNode   string
	OriginNodeIP string
}

// FieldMap 请求字段映射
func (reqData *CollectedDevice) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.SN:          "sn",
		&reqData.Vendor:      "vendor",
		&reqData.Model:       "model",
		&reqData.Arch:        "arch",
		&reqData.BootOSIP:    "bootos_ip",
		&reqData.BootOSMac:   "bootos_mac",
		&reqData.NICDevice:   "nic_device",
		&reqData.CPU:         "cpu",
		&reqData.Memory:      "memory",
		&reqData.Disk:        "disk",
		&reqData.DiskSlot:    "disk_slot",
		&reqData.NIC:         "nic",
		&reqData.Motherboard: "motherboard",
		&reqData.OOB:         "oob",
		&reqData.BIOS:        "bios",
		&reqData.RAID:        "raid",
		&reqData.Power:       "power",
		&reqData.Fan:         "fan",
		&reqData.PCI:         "pci",
		&reqData.HBA:         "hba",
		&reqData.LLDP:        "lldp",
		&reqData.Extra:       "extra",
	}
}

// Validate 结构体数据校验
func (reqData *CollectedDevice) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if reqData.SN == "" {
		errs.Add([]string{"sn"}, binding.RequiredError, "设备序列号'SN'不能为空")
		return errs
	}
	if reqData.Vendor == "" {
		errs.Add([]string{"vendor"}, binding.RequiredError, "厂商名'vendor'不能为空")
		return errs
	}
	return errs
}

func (reqData *CollectedDevice) nicJSON() (json string) {
	if reqData.NIC == nil {
		return json
	}
	return string(reqData.NIC.ToJSON())
}

func (reqData *CollectedDevice) cpuJSON() (json string) {
	if reqData.CPU == nil {
		return json
	}
	return string(reqData.CPU.ToJSON())
}

func (reqData *CollectedDevice) memJSON() (json string) {
	if reqData.Memory == nil {
		return json
	}
	return string(reqData.Memory.ToJSON())
}

func (reqData *CollectedDevice) diskJSON() (json string) {
	if reqData.Disk == nil {
		return json
	}
	return string(reqData.Disk.ToJSON())
}

func (reqData *CollectedDevice) boardJSON() (json string) {
	if reqData.Motherboard == nil {
		return json
	}
	return string(reqData.Motherboard.ToJSON())
}

func (reqData *CollectedDevice) diskSlotJSON() (json string) {
	if reqData.DiskSlot == nil {
		return json
	}
	return string(reqData.DiskSlot.ToJSON())
}

func (reqData *CollectedDevice) oobJSON() (json string) {
	if reqData.OOB == nil {
		return json
	}
	return string(reqData.OOB.ToJSON())
}

func (reqData *CollectedDevice) biosJSON() (json string) {
	if reqData.BIOS == nil {
		return json
	}
	return string(reqData.BIOS.ToJSON())
}

func (reqData *CollectedDevice) raidJSON() (json string) {
	if reqData.RAID == nil || len(reqData.RAID.Items) <= 0 {
		return json
	}
	return string(reqData.RAID.ToJSON())
}

func (reqData *CollectedDevice) fanJSON() (json string) {
	if reqData.Fan == nil {
		return json
	}
	return string(reqData.Fan.ToJSON())
}

func (reqData *CollectedDevice) powerJSON() (json string) {
	if reqData.Power == nil {
		return json
	}
	return string(reqData.Power.ToJSON())
}

func (reqData *CollectedDevice) hbaJSON() (json string) {
	if reqData.HBA == nil {
		return json
	}
	return string(reqData.HBA.ToJSON())
}

func (reqData *CollectedDevice) pciJSON() (json string) {
	if reqData.PCI == nil {
		return json
	}
	return string(reqData.PCI.ToJSON())
}

func (reqData *CollectedDevice) lldpJSON() (json string) {
	if reqData.LLDP == nil {
		return json
	}
	return string(reqData.LLDP.ToJSON())
}

func (reqData *CollectedDevice) extraJSON() (json string) {
	if reqData.Extra == nil {
		return json
	}
	return string(reqData.Extra.ToJSON())
}

// SaveCollectedDevice 保存采集到的设备信息
func SaveCollectedDevice(log logger.Logger, repo model.Repo, reqData *CollectedDevice) (err error) {
	var cpusum, memsum, disksum uint
	if reqData.CPU != nil && reqData.CPU.TotalCores > 0 {
		cpusum = uint(reqData.CPU.TotalCores)
	}
	if reqData.Memory != nil && reqData.Memory.TotalSizeMB > 0 {
		memsum = uint(reqData.Memory.TotalSizeMB)
	}
	if reqData.Disk != nil && reqData.Disk.TotalSizeGB > 0 {
		disksum = uint(reqData.Disk.TotalSizeGB)
	}

	dev := model.CollectedDevice{
		OriginNode:   reqData.OriginNode,
		OriginNodeIP: reqData.OriginNodeIP,
		SN:           reqData.SN,
		Vendor:       reqData.Vendor,
		ModelName:    reqData.Model,
		Arch:         reqData.Arch,
		BootOSIP:     reqData.BootOSIP,
		BootOSMac:    reqData.BootOSMac,
		CPUSum:       cpusum,
		CPU:          reqData.cpuJSON(),
		MemorySum:    memsum,
		Memory:       reqData.memJSON(),
		DiskSum:      disksum,
		Disk:         reqData.diskJSON(),
		DiskSlot:     reqData.diskSlotJSON(),
		NIC:          reqData.nicJSON(),
		NICDevice:    reqData.NICDevice,
		Motherboard:  reqData.boardJSON(),
		RAID:         reqData.raidJSON(),
		OOB:          reqData.oobJSON(),
		BIOS:         reqData.biosJSON(),
		Fan:          reqData.fanJSON(),
		Power:        reqData.powerJSON(),
		HBA:          reqData.hbaJSON(),
		PCI:          reqData.pciJSON(),
		LLDP:         reqData.lldpJSON(),
		Extra:        reqData.extraJSON(),
	}
	return repo.SaveCollectedDeviceBySN(&dev)
}

// CheckUnique 检查字段重复性, 传入列表，及需要排查的字段
//检查的基本思路是将被检查元素作为map的key,遍历源数据，发现重复则value+1，最后检查是否value大于1
func CheckUnique(importDatas [][]string) (valid bool, err error) {
	//校验SN,固资编号，机位编号（机架编号+机位编号）唯一性
	//有几个检查的字段，就定义几个数组元素来做检查
	checkers := make([]map[interface{}]int, 3)
	for i := 0; i < 3; i++ {
		checkers[i] = make(map[interface{}]int)
	}
	for i, d := range importDatas {
		if i == 0 {
			continue
		}
		if d[0] != "" {
			checkers[0][d[0]]++
		}
		checkers[1][d[1]]++
		checkers[2][fmt.Sprintf("%s+%s+%s", d[5], d[6], d[7])]++
	}

	for k, v := range checkers[0] {
		if v > 1 {
			return false, fmt.Errorf("[数据校验]固资编号值:%v不允许重复", k)
		}
	}
	for k, v := range checkers[1] {
		if v > 1 {
			return false, fmt.Errorf("[数据校验]序列号(SN)值:%v不允许重复", k)
		}
	}
	for k, v := range checkers[2] {
		if v > 1 {
			return false, fmt.Errorf("[数据校验]机位编号值(机房名称+机架编号+机位编号):%v不允许重复", k)
		}
	}
	return true, nil
}

// DevicePageByTorReq 物理机分页列表搜索字段
type DevicePageByTorReq struct {
	//用途 'TDSQL','APP','CVM','TGW','NAS','Other'
	Usage string `json:"usage"`
	// Tor分组
	Tor      string `json:"tor"`
	Page     int64
	PageSize int64
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *DevicePageByTorReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.Usage:    "usage",
		&reqData.Tor:      "tor",
		&reqData.Page:     "page",
		&reqData.PageSize: "page_size",
	}
}

// Validate DevicePageByTorReq
func (reqData *DevicePageByTorReq) Validate(request *http.Request, errs binding.Errors) binding.Errors {
	if reqData.Tor == "" {
		errs.Add([]string{"tor"}, binding.RequiredError, "TOR不能为空")
		return errs
	}

	// 默认为CVM
	if reqData.Usage == "" {
		reqData.Usage = "CVM"
	}
	return errs
}

// GetDeviceByTorPage 根据TOR分组查询物理机分页列表
func GetDeviceByTorPage(log logger.Logger, repo model.Repo, conf *config.Config, reqData *DevicePageByTorReq) (pg *page.Page, err error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 100 {
		reqData.PageSize = 20
	}
	if reqData.Page < 0 {
		reqData.Page = 1
	}

	var devs []*model.Device
	if reqData.Tor != "" {
		// 根据tor找到对应的网络设备
		tors := strings.Split(reqData.Tor, ",")
		nd, err := repo.GetNetworkDeviceByTORS(tors...)
		if err != nil {
			return nil, err
		}
		// 根据对于的网络设备找到对应的机位
		if len(nd) > 0 {
			var ndName []string
			for k := range nd {
				ndName = append(ndName, nd[k].Name)
			}
			suid, err := repo.GetServerUsiteByNetworkDeviceName(ndName)
			if err != nil {
				return nil, err
			}
			if len(suid) > 0 {
				devs, err = repo.GetDevicesByUSiteIDS(suid, reqData.Usage)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	pager := page.NewPager(reflect.TypeOf(&DevicePageResp{}), reqData.Page, reqData.PageSize, int64(len(devs)))
	for i := range devs {
		dev, err := conver2DevicePagesResp(log, repo, conf, devs[i])
		if err != nil {
			return nil, err
		}
		pager.AddRecords(dev)
	}
	return pager.BuildPage(), nil
}

//OperationStatusTransfer 运行状态值和数据库存储值的转换
func OperationStatusTransfer(status string, reverse bool) string {
	mStatus := map[string]string{
		"运行中(需告警)":  "run_with_alarm",
		"运行中(无需告警)": "run_without_alarm",
		"重装中":       "reinstalling",
		"搬迁中":       "moving",
		"待退役":       "pre_retire",
		"退役中":       "retiring",
		"已退役":       "retired",
		"待部署":       "pre_deploy",
		"已上架":       "on_shelve",
		"回收中":       "recycling",
		"维护中":       "maintaining",
		"待搬迁":       "pre_move",
		"库房中":       "in_store",
	}
	if !reverse {
		if val, ok := mStatus[status]; ok {
			return val
		}
	} else {
		for key, val := range mStatus {
			if val == status {
				return key
			}
		}
	}
	return ""
}

//全局计数器
var GlobalCount = make(map[string]int, 0)

// 自动生成固资编号
func GenFixedAssetNumber(repo model.Repo) (string, error) {
	//格式：WDEV+YYMM+0001
	//现在是几月：1906
	prefix := "WDEV"
	m := time.Now().Format("0601")
	if count, ok := GlobalCount[m]; ok && count != 0 {
		GlobalCount[m]++
	} else {
		//如果是服务重启，则需要从db里重新获取下最大的固资编号
		//如果是新月份，计数从1开始计
		currentMax, err := repo.GetMaxFixedAssetNumber(m)
		if currentMax == "" || err != nil {
			GlobalCount[m] = 1
		} else {
			currentMaxInt, err := strconv.Atoi(strings.TrimPrefix(currentMax, prefix+m))
			if err != nil {
				return "", err
			}
			GlobalCount[m] = currentMaxInt + 1
		}
	}
	return fmt.Sprintf("%s%s%04d", prefix, m, GlobalCount[m]), nil
}

func validOperationStatus(s string) bool {
	switch s {
	case model.DevOperStatRunWithoutAlarm:
	case model.DevOperStatRunWithAlarm:
	case model.DevOperStateRetired:
	case model.DevOperStatMoving:
	case model.DevOperStatOnShelve:
	case model.DevOperStatPreDeploy:
	case model.DevOperStatPreRetire:
	case model.DevOperStatRetiring:
	case model.DevOperStatRecycling:
	case model.DevOperStatPreMove:
	case model.DevOperStatMaintaining:
	case model.DevOperStatReinstalling:
	default:
		return false
	}
	return true
}


//NewDevice API新增设备所需字段
type NewDevice struct {
	//固资编号
	FixedAssetNum string `json:"fixed_asset_number"`
	//序列号
	SN string `json:"sn"`
	//厂商
	Vendor string `json:"vendor"`
	//型号
	Model string `json:"model"`
	// CPU架构
	//Arch string `json:"arch"`
	//用途 'TDSQL','APP','CVM','TGW','NAS','Other'
	Usage string `json:"usage"`
	//分类
	Category string `json:"category"`
	//机房管理单元名称
	ServerRoomName string `json:"server_room_name"`
	//机架编号
	CabinetNum string `json:"server_cabinet_number"`
	//机位编号
	USiteNum string `json:"server_usite_number"`
	//库房管理单元名称-与机房管理单元-机架编号-机位互斥
	StoreRoomName string `json:"store_room_name"`
	//虚拟货架编号-与机房管理单元-机架编号-机位互斥
	VCabinetNum string `json:"virtual_cabinet_number"`
	//硬件说明
	HardwareRemark string `json:"hardware_remark"`
	//RAID说明
	RAIDRemark string `json:"raid_remark"`
	// OOB初始用户密码,':'分隔,主要针对的是旧机器导入使用，新机器如果带外是出厂默认，可以缺省
	OOBInit string `json:"oob_init"`
	//OriginNodeIP proxy节点IP
	//OriginNodeIP string `json:"origin_node_ip"` //废弃
	//启用时间
	StartedAt string `json:"started_at"`
	//上架时间
	OnShelveAt string `json:"onshelve_at"`
	//关联订单号(非必填)
	OrderNumber string `json:"order_number"`
	// 非标的特殊设备
	IsSpecialDevice	bool			`json:"is_special_device"`
	// 非标的特殊设备:是否分配IPv4 yes|no
	NeedExtranetIPv4 	string 			`json:"need_extranet_ipv4"`
	NeedIntranetIPv4 	string 			`json:"need_intranet_ipv4"`
	// 非标的特殊设备:操作系统
	OS string `json:"os"`
	//设备生命周期相关,详情参考 model.DeviceLifecycle
	AssetBelongs	 				string		`json:"asset_belongs"`
	Owner			 				string		`json:"owner"`
	IsRental		 				string		`json:"is_rental"`
	MaintenanceServiceProvider		string		`json:"maintenance_service_provider"`
	MaintenanceService				string		`json:"maintenance_service"`
	LogisticsService				string		`json:"logistics_service"`
	MaintenanceServiceDateBegin     string 		`json:"maintenance_service_date_begin"`
	MaintenanceServiceDateEnd       string 		`json:"maintenance_service_date_end"`
	MaintenanceServiceStatus		string		`json:"maintenance_service_status"`
	// 数据校验用
	ErrMsgContent string `json:"content"`
	//以下字段是通过名称关联到
	idcID        uint //数据中心ID
	serverRoomID uint //机房管理单元ID
	storeRoomID uint //机房管理单元ID
	cabinetID    uint //机架ID
	vcabinetID  uint //机架ID
	uSiteID      uint //机位ID
}


// NewDevicesReq 请求API新增设备结构体
type NewDevicesReq struct {
	NewDevices		NewDevicesList					
	CurrentUser 	*model.CurrentUser
}

type NewDevicesList []*NewDevice


//checkLength NewDevice数据做字段长度校验
func (newDev *NewDevice) checkLength() {
	leg := len(newDev.SN)
	if leg == 0 || leg > 255 {
		var br string
		if newDev.ErrMsgContent != "" {
			br = "<br />"
		}
		newDev.ErrMsgContent += br + fmt.Sprintf("必填项校验:序列号长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(newDev.Vendor)
	if leg == 0 || leg > 255 {
		var br string
		if newDev.ErrMsgContent != "" {
			br = "<br />"
		}
		newDev.ErrMsgContent += br + fmt.Sprintf("必填项校验:厂商长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(newDev.Model)
	if leg == 0 || leg > 255 {
		var br string
		if newDev.ErrMsgContent != "" {
			br = "<br />"
		}
		newDev.ErrMsgContent += br + fmt.Sprintf("必填项校验:型号长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(newDev.Usage)
	if leg == 0 || leg > 255 {
		var br string
		if newDev.ErrMsgContent != "" {
			br = "<br />"
		}
		newDev.ErrMsgContent += br + fmt.Sprintf("必填项校验:用途长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(newDev.Category)
	if leg == 0 || leg > 255 {
		var br string
		if newDev.ErrMsgContent != "" {
			br = "<br />"
		}
		newDev.ErrMsgContent += br + fmt.Sprintf("必填项校验:设备类型长度为(%d)(不能为空，不能大于255)", leg)
	}
	if len(newDev.ServerRoomName) != 0 && len(newDev.StoreRoomName) != 0 {
		var br string
		if newDev.ErrMsgContent != "" {
			br = "<br />"
		}
		newDev.ErrMsgContent += br + fmt.Sprintf("必填项校验:机房管理单元名称与库房管理单元名称不能同时存在")
	} else if len(newDev.ServerRoomName) == 0 && len(newDev.StoreRoomName) == 0 {
		var br string
		if newDev.ErrMsgContent != "" {
			br = "<br />"
		}
		newDev.ErrMsgContent += br + fmt.Sprintf("必填项校验:机房管理单元名称与库房管理单元名称不能同时为空")
	} else if len(newDev.ServerRoomName) != 0 {
		leg = len(newDev.ServerRoomName)
		if leg > 255 {
			var br string
			if newDev.ErrMsgContent != "" {
				br = "<br />"
			}
			newDev.ErrMsgContent += br + fmt.Sprintf("必填项校验:机房管理单元名称长度为(%d)(不能大于255)", leg)
		}
		leg = len(newDev.CabinetNum)
		if leg == 0 || leg > 255 {
			var br string
			if newDev.ErrMsgContent != "" {
				br = "<br />"
			}
			newDev.ErrMsgContent += br + fmt.Sprintf("必填项校验:机架编号长度为(%d)(不能为空，不能大于255)", leg)
		}
		leg = len(newDev.USiteNum)
		if leg == 0 || leg > 255 {
			var br string
			if newDev.ErrMsgContent != "" {
				br = "<br />"
			}
			newDev.ErrMsgContent += br + fmt.Sprintf("必填项校验:机位编号长度为(%d)(不能为空，不能大于255)", leg)
		}
	} else if len(newDev.StoreRoomName) != 0 {
		leg = len(newDev.StoreRoomName)
		if leg > 255 {
			var br string
			if newDev.ErrMsgContent != "" {
				br = "<br />"
			}
			newDev.ErrMsgContent += br + fmt.Sprintf("必填项校验:库房管理单元名称长度为(%d)(不能大于255)", leg)
		}
		leg = len(newDev.VCabinetNum)
		if leg == 0 || leg > 255 {
			var br string
			if newDev.ErrMsgContent != "" {
				br = "<br />"
			}
			newDev.ErrMsgContent += br + fmt.Sprintf("必填项校验:虚拟货架编号长度为(%d)(不能为空，不能大于255)", leg)
		}
	}
	leg = len(newDev.StartedAt)
	if leg == 0 || leg > 255 {
		var br string
		if newDev.ErrMsgContent != "" {
			br = "<br />"
		}
		newDev.ErrMsgContent += br + fmt.Sprintf("必填项校验:启用时间长度为(%d)(不能为空，不能大于255)", leg)
	}
	if !strings.Contains(newDev.StartedAt, "-") {
		var br string
		if newDev.ErrMsgContent != "" {
			br = "<br />"
		}
		newDev.ErrMsgContent += br + "启用时间格式须为：YYYY-MM-DD"
	}
	leg = len(newDev.OnShelveAt)
	if leg == 0 || leg > 255 {
		var br string
		if newDev.ErrMsgContent != "" {
			br = "<br />"
		}
		newDev.ErrMsgContent += br + fmt.Sprintf("必填项校验:上架时间长度为(%d)(不能为空，不能大于255)", leg)
	}
	if !strings.Contains(newDev.OnShelveAt, "-") {
		var br string
		if newDev.ErrMsgContent != "" {
			br = "<br />"
		}
		newDev.ErrMsgContent += br + "上架时间格式须为：YYYY-MM-DD"
	}
	leg = len(newDev.MaintenanceServiceDateBegin)
	if leg == 0 || leg > 255 {
		var br string
		if newDev.ErrMsgContent != "" {
			br = "<br />"
		}
		newDev.ErrMsgContent += br + fmt.Sprintf("必填项校验:维保起始日期长度为(%d)(不能为空，不能大于255)", leg)
	}
	if !strings.Contains(newDev.MaintenanceServiceDateBegin, "-") {
		var br string
		if newDev.ErrMsgContent != "" {
			br = "<br />"
		}
		newDev.ErrMsgContent += br + "维保起始日期格式须为：YYYY-MM-DD"
	}
	leg = len(newDev.MaintenanceServiceDateEnd)
	if leg == 0 || leg > 255 {
		var br string
		if newDev.ErrMsgContent != "" {
			br = "<br />"
		}
		newDev.ErrMsgContent += br + fmt.Sprintf("必填项校验:维保截止日期长度为(%d)(不能为空，不能大于255)", leg)
	}
	if !strings.Contains(newDev.MaintenanceServiceDateEnd, "-") {
		var br string
		if newDev.ErrMsgContent != "" {
			br = "<br />"
		}
		newDev.ErrMsgContent += br + "维保截止日期格式须为：YYYY-MM-DD"
	}
	//硬件说明
	//RAID说明
	//带外说明
	//暂时不做检验
}


//validate NewDevice数据基本验证
func (newDev *NewDevice) validate(repo model.Repo) error {
	//机房-机架-机位校验
	if newDev.ServerRoomName != "" {
		srs, err := repo.GetServerRoomByName(newDev.ServerRoomName)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == gorm.ErrRecordNotFound || srs == nil {
			var br string
			if newDev.ErrMsgContent != "" {
				br = "<br />"
			}
			newDev.ErrMsgContent += br + fmt.Sprintf("机房名(%s)不存在", newDev.ServerRoomName)
		} else if srs != nil {
			newDev.idcID = srs.IDCID
			newDev.serverRoomID = srs.ID
		}
	    //机架
	    cabinet, err := repo.GetServerCabinetByNumber(newDev.serverRoomID, newDev.CabinetNum)
	    if err != nil && err != gorm.ErrRecordNotFound {
	    	return err
	    }
	    if err == gorm.ErrRecordNotFound || cabinet == nil {
	    	var br string
	    	if newDev.ErrMsgContent != "" {
	    		br = "<br />"
	    	}
	    	newDev.ErrMsgContent += br + fmt.Sprintf("机架编号(%s)不存在", newDev.CabinetNum)
	    } else {
	    	newDev.cabinetID = cabinet.ID
	    }
	    //机位
	    uSite, err := repo.GetServerUSiteByNumber(newDev.cabinetID, newDev.USiteNum)
	    if err != nil && err != gorm.ErrRecordNotFound {
	    	return err
	    }
	    if err == gorm.ErrRecordNotFound || uSite == nil {
	    	var br string
	    	if newDev.ErrMsgContent != "" {
	    		br = "<br />"
	    	}
	    	newDev.ErrMsgContent += br + fmt.Sprintf("机位编号(%s)不存在, 机架编号(%s)", newDev.USiteNum, newDev.CabinetNum)
	    } else {
	    	dev, _ := repo.GetDeviceBySN(newDev.SN)
	    	if !CheckUSiteFree(repo, uSite.ID, dev) {
	    		var br string
	    		if newDev.ErrMsgContent != "" {
	    			br = "<br />"
	    		}
	    		newDev.ErrMsgContent += br + fmt.Sprintf("机位编号(%s), 机架编号(%s)被占用或不可用", newDev.USiteNum, newDev.CabinetNum)
	    	} else {
	    		newDev.uSiteID = uSite.ID
	    	}
	    }		
	} else if newDev.StoreRoomName != "" {
		//库房校验
		srs, err := repo.GetStoreRoomByName(newDev.StoreRoomName)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == gorm.ErrRecordNotFound || srs == nil {
			var br string
			if newDev.ErrMsgContent != "" {
				br = "<br />"
			}
			newDev.ErrMsgContent += br + fmt.Sprintf("库房(%s)不存在", newDev.StoreRoomName)
		} else if srs != nil {
			newDev.idcID = srs.IDCID
			newDev.storeRoomID = srs.ID
		}
	    //虚拟货架
	    cabinet, err := repo.GetVirtualCabinets(&model.VirtualCabinet{
	    	StoreRoomID: newDev.storeRoomID,
	    	Number:      newDev.VCabinetNum,
	    }, nil, nil)
	    if err != nil && err != gorm.ErrRecordNotFound {
	    	return err
	    }
	    if err == gorm.ErrRecordNotFound || len(cabinet) <= 0 {
	    	var br string
	    	if newDev.ErrMsgContent != "" {
	    		br = "<br />"
	    	}
	    	newDev.ErrMsgContent += br + fmt.Sprintf("虚拟货架(%s)不存在", newDev.VCabinetNum)
	    } else if len(cabinet) > 0 {
	    	newDev.vcabinetID = cabinet[0].ID
	    }		
	} else {
		var br string
		if newDev.ErrMsgContent != "" {
			br = "<br />"
		}
		newDev.ErrMsgContent += br + fmt.Sprintf("SN（%s）必须指定机房管理单元或库房管理单元", newDev.SN)
	}

	if newDev.OOBInit != "" {
		if !strings.Contains(newDev.OOBInit, ":") {
			var br string
			if newDev.ErrMsgContent != "" {
				br = "<br />"
			}
			newDev.ErrMsgContent += br + fmt.Sprintf("带外用户密码初始值:%s格式不正确，应以':'分隔", newDev.OOBInit)
		}
	}
	// 订单关联字段校验：IDC\设备类型\设备用途
	if newDev.OrderNumber != "" {
		order, err := repo.GetOrderByNumber(newDev.OrderNumber)
		if err != nil {
			var br string
			if newDev.ErrMsgContent != "" {
				br = "<br />"
			}
			newDev.ErrMsgContent += br + fmt.Sprintf("订单(订单号:%s)不存在", newDev.OrderNumber)
		}
		if order != nil {
			if order.IDCID != newDev.idcID {
				var br string
				if newDev.ErrMsgContent != "" {
					br = "<br />"
				}
				newDev.ErrMsgContent += br + fmt.Sprintf("订单(订单号:%s)与设备（SN:%s）不属于同一个数据中心", newDev.OrderNumber, newDev.SN)
			}
			if order.Category != newDev.Category {
				var br string
				if newDev.ErrMsgContent != "" {
					br = "<br />"
				}
				newDev.ErrMsgContent += br + fmt.Sprintf("订单(订单号:%s)的设备类型（%s） 与 设备（SN:%s）的设备类型（%s）不匹配", newDev.OrderNumber, order.Category, newDev.SN, newDev.Category)
			}
			if order.Usage != newDev.Usage {
				var br string
				if newDev.ErrMsgContent != "" {
					br = "<br />"
				}
				newDev.ErrMsgContent += br + fmt.Sprintf("订单(订单号:%s)的用途（%s） 与 设备（SN:%s）的用途（%s）不匹配", newDev.OrderNumber, order.Usage, newDev.SN, newDev.Usage)
			}
		}
	}
	// 是否租赁
	if newDev.IsRental != model.YES && newDev.IsRental != model.NO {
		var br string
		if newDev.ErrMsgContent != "" {
			br = "<br />"
		}
		newDev.ErrMsgContent += br + fmt.Sprintf("是否租赁 必须为 %s or %s", model.YES, model.NO)
	}
	// 维保状态校验
	if newDev.MaintenanceServiceStatus != model.MaintenanceServiceStatusUnderWarranty && newDev.MaintenanceServiceStatus != model.MaintenanceServiceStatusOutOfWarranty && newDev.MaintenanceServiceStatus != model.MaintenanceServiceStatusInactive {
		var br string
		if newDev.ErrMsgContent != "" {
			br = "<br />"
		}
		newDev.ErrMsgContent += br + fmt.Sprintf("维保状态 必须为 %s or %s or %s", model.MaintenanceServiceStatusUnderWarranty, model.MaintenanceServiceStatusOutOfWarranty, model.MaintenanceServiceStatusInactive)
	}
	return nil
}


// SaveNewDevices 保存（新增）
func SaveNewDevices(log logger.Logger, repo model.Repo, conf *config.Config, reqData *NewDevicesReq) (succeedSNs []string, totalAffected int64, err error) {
	if len(reqData.NewDevices) == 0 {
		log.Error("NewDevices 参数不能为空")
		return succeedSNs, totalAffected, fmt.Errorf("NewDevices 参数不能为空")
	}
	
	// 用于导入后批量更新带外
	var devices []*model.Device
	// 用于关联订单的信息更新
	var mOrderAmount = make(map[string]int, 0)

	for _, newDev := range reqData.NewDevices {

		log.Infof("Begin to save new device SN: %s ", newDev.SN)
		defer log.Infof("End to save new device SN: %s ", newDev.SN)
		//必填项校验
		newDev.checkLength()
		//数据校验并转化数据中心、机房管理单元、机架、机位名称为对应的ID
		err = newDev.validate(repo)
		if err != nil {
			return succeedSNs, totalAffected, err
		}
		if newDev.ErrMsgContent != "" {
			log.Errorf("设备检验参数失败：%s", newDev.ErrMsgContent)
			return succeedSNs, totalAffected, fmt.Errorf("设备检验参数失败：%s", newDev.ErrMsgContent)
		}
		//Device 设备结构体
		mod := &model.Device{
			SN:             strings.TrimSpace(newDev.SN),	//去除首尾空白
			Vendor:         newDev.Vendor,
			DevModel:       newDev.Model,
			Usage:          newDev.Usage,
			Category:       newDev.Category,
			IDCID:          newDev.idcID,
			//ServerRoomID:   newDev.serverRoomID,
			//CabinetID:      newDev.cabinetID,
			//USiteID:        &newDev.uSiteID,
			//StoreRoomID:    newDev.storeRoomID,
			//VCabinetID:     newDev.vcabinetID,
			HardwareRemark: newDev.HardwareRemark,
			RAIDRemark:     newDev.RAIDRemark,
			OOBInit:        "{}",
			PowerStatus:    model.PowerStatusOff,
			OrderNumber:    newDev.OrderNumber,
			Creator:        reqData.CurrentUser.LoginName,
			//JSON type的字段需要默认赋空值
			CPU:         "{}",
			Memory:      "{}",
			Disk:        "{}",
			DiskSlot:    "{}",
			NIC:         "{}",
			Motherboard: "{}",
			RAID:        "{}",
			OOB:         "{}",
			BIOS:        "{}",
			Fan:         "{}",
			Power:       "{}",
			HBA:         "{}",
			PCI:         "{}",
			Switch:      "{}",
			LLDP:        "{}",
			Extra:       "{}",
		}
		//处理所有字段的多余空格字符
		utils.StructTrimSpace(mod)
		//转换日期格式		
		mod.StartedAt, _ = time.Parse(times.DateLayout, newDev.StartedAt)
		if newDev.OnShelveAt != "" {
			mod.OnShelveAt, _ = time.Parse(times.DateLayout, newDev.OnShelveAt)
		} else {
			mod.OnShelveAt = time.Now()
		}
		if newDev.ServerRoomName != "" {
			mod.ServerRoomID = newDev.serverRoomID
			mod.CabinetID = newDev.cabinetID
			mod.USiteID = &newDev.uSiteID
			//特殊设备默认无须部署，直接置为已上架
			if newDev.IsSpecialDevice {
				mod.OperationStatus = model.DevOperStatOnShelve
			} else {
				mod.OperationStatus = model.DevOperStatPreDeploy
			}
		} else if newDev.StoreRoomName != "" {
			mod.StoreRoomID = newDev.storeRoomID
			mod.VCabinetID = newDev.vcabinetID
			mod.OperationStatus = model.DevOperStatInStore
		} else {
			err = fmt.Errorf("SN:%s 未指定机房管理单元或者库房管理单元", newDev.SN)
			return succeedSNs, totalAffected, err
		}
		if newDev.OOBInit != "" {
			words := strings.Split(newDev.OOBInit, ":")
			if len(words) == 2 {
				ou := OOBUser{
					Username: words[0],
					Password: words[1],
				}
				if b, err := json.Marshal(ou); err == nil {
					mod.OOBInit = string(b)
				}
			}
		}
		//新增设备场景未指定固资号时，自动生成固资编号
		if newDev.FixedAssetNum == "" {
			newDev.FixedAssetNum, err = GenFixedAssetNumber(repo)
			mod.FixedAssetNumber = newDev.FixedAssetNum
			if err != nil {
				log.Errorf("generate fixed_asset number for SN:%s fail", newDev.SN)
				return succeedSNs, totalAffected, fmt.Errorf("自动生成固资号失败：%v", err)
			}
			log.Debugf("自动生成固资号: %s", newDev.FixedAssetNum)
		} else {
			log.Debugf("非系统自动生成固资号: %s", newDev.FixedAssetNum)
			mod.FixedAssetNumber = newDev.FixedAssetNum
		}
		
		//默认电源状态 OFF
		if mod.PowerStatus == "" {
			mod.PowerStatus = model.PowerStatusOff
		}
		// 仅记录必要字段到“设备新增”
		optDetail, err := convert2DetailOfOperationTypeAdd(repo, *mod)
		if err != nil {
			log.Errorf("Fail to convert Detail of OperationTypeAdd: %v", err)
		}
		// DeviceLifecycle 变更记录
		deviceLifecycleLog := []model.ChangeLog {
			{
				OperationUser:		reqData.CurrentUser.LoginName,
				OperationType:		model.OperationTypeAdd,
				OperationDetail:	optDetail,
				OperationTime:		times.ISOTime(time.Now()).ToTimeStr(),
			},
		}
		b, _ := json.Marshal(deviceLifecycleLog)
		// SaveDeviceLifecycleReq 结构体
		saveDevLifecycleReq := &SaveDeviceLifecycleReq {
			DeviceLifecycleBase: DeviceLifecycleBase{
				FixedAssetNumber: 				newDev.FixedAssetNum,
				SN:             				newDev.SN,
				AssetBelongs:	 				newDev.AssetBelongs,			
				Owner:				 			newDev.Owner,
				IsRental:		 				newDev.IsRental,		
				MaintenanceServiceProvider: 	newDev.MaintenanceServiceProvider, 		
				MaintenanceService:				newDev.MaintenanceService,				
				LogisticsService:				newDev.LogisticsService,
				MaintenanceServiceStatus:		newDev.MaintenanceServiceStatus,
				LifecycleLog:					string(b),
			},
		}
		saveDevLifecycleReq.MaintenanceServiceDateBegin, _ = time.Parse(times.DateLayout, newDev.MaintenanceServiceDateBegin)
		saveDevLifecycleReq.MaintenanceServiceDateEnd, _ = time.Parse(times.DateLayout, newDev.MaintenanceServiceDateEnd)
		// 通过订单编号获取资产归属、负责人、维保服务等内容
		// 若无订单编号则以参数输入为准
		if newDev.OrderNumber != "" {
			order, err := repo.GetOrderByNumber(newDev.OrderNumber)
			if err != nil {
				log.Errorf("订单(订单号:%s)不存在", newDev.OrderNumber)
				return succeedSNs, totalAffected, err
			}
			if order != nil {
				saveDevLifecycleReq.AssetBelongs = order.AssetBelongs	 			
				saveDevLifecycleReq.Owner = order.Owner
				saveDevLifecycleReq.IsRental = order.IsRental
				saveDevLifecycleReq.MaintenanceServiceProvider = order.MaintenanceServiceProvider
				saveDevLifecycleReq.MaintenanceService = order.MaintenanceService
				saveDevLifecycleReq.LogisticsService = order.LogisticsService
				saveDevLifecycleReq.MaintenanceServiceDateBegin = order.MaintenanceServiceDateBegin
				saveDevLifecycleReq.MaintenanceServiceDateEnd = order.MaintenanceServiceDateEnd
			}
		}
		//插入或者更新
		//查询是否已经存在
		dev, err := repo.GetDeviceBySN(newDev.SN)
		if err != nil && err != gorm.ErrRecordNotFound {
			return succeedSNs, totalAffected, err
		} 
		if dev != nil {
			// 特殊设备仅允许已上架、库房中重复导入
			if newDev.IsSpecialDevice {
				if dev.OperationStatus == model.DevOperStatInStore || dev.OperationStatus == model.DevOperStatOnShelve {
					log.Debugf("SpecialDevice SN %s already exist.Status is %s . Allow to update.", newDev.SN, dev.OperationStatus)
					mod.ID = dev.ID
					mod.CreatedAt = dev.CreatedAt
					mod.UpdatedAt = time.Now()
					mod.OperationStatus = dev.OperationStatus
				} else {
					log.Debugf("SpecialDevice SN %s already exist. Status is %s . Update is forbidden.", newDev.SN, dev.OperationStatus)
					err = fmt.Errorf("特殊设备（SN:%s）导入失败，仅允许运营状态为[库房中|已上架]重复导入更新", newDev.SN)
					return succeedSNs, totalAffected, err
				}
			} else {
				if dev.OperationStatus == model.DevOperStatInStore || dev.OperationStatus == model.DevOperStatPreDeploy {
					log.Debugf("Device SN %s already exist.Status is %s . Allow to update.", newDev.SN, dev.OperationStatus)
					mod.ID = dev.ID
					mod.OperationStatus = dev.OperationStatus
					mod.CreatedAt = dev.CreatedAt
					mod.UpdatedAt = time.Now()
				} else {
					log.Debugf("Device SN %s already exist. Status is %s . Update is forbidden.", newDev.SN, dev.OperationStatus)
					err = fmt.Errorf("设备（SN:%s）导入失败，仅允许运营状态为[库房中|待部署]重复导入更新", newDev.SN)
					return succeedSNs, totalAffected, err
				}
			}
		} else {
			// 首次新增时统计关联订单到货数量
			if newDev.OrderNumber != "" {
				mOrderAmount[newDev.OrderNumber]++ 
			}
		}
		// DeviceLifecycle 查询是否已经存在
		devLifecycle, err := repo.GetDeviceLifecycleBySN(newDev.SN)
		if err != nil && err != gorm.ErrRecordNotFound {
			return succeedSNs, totalAffected, err
		} 
		if devLifecycle != nil {
			log.Debugf("DeviceLifecycle SN %s already exist.Update it.", newDev.SN)
			saveDevLifecycleReq.ID = devLifecycle.ID
		}
		// 首次新增需先保存设备，后续才可支持分配IP
		if _, err = repo.SaveDevice(mod); err != nil {
			log.Debug(err)
			return succeedSNs, totalAffected, err
		}
		// 保存或更新 DeviceLifecycle
		if err = SaveDeviceLifecycle(log, repo, saveDevLifecycleReq); err != nil {
			log.Debug(err)
			return succeedSNs, totalAffected, err
		}
		// 特殊设备默认无需进行OS部署，指定用途、OS、运营状态
		if newDev.IsSpecialDevice {
			mod.Usage = model.DevUsageSpecialDev
			//分配IP
			var inIP, exIP *model.IP
			if newDev.NeedIntranetIPv4 == model.YES {
				if inIP, err = repo.AssignIntranetIP(newDev.SN); err != nil {
					err = fmt.Errorf("SN:%s分配内网失败，err：%v", newDev.SN, err)
					return succeedSNs, totalAffected, err
				}
			} else {
				//需要尝试释放IP
				if _, err = repo.ReleaseIP(newDev.SN, model.Intranet); err != nil {
					err = fmt.Errorf("SN:%s释放内网IP失败，err：%v", newDev.SN, err)
					return succeedSNs, totalAffected, err
				}
			}
		
			if newDev.NeedExtranetIPv4 == model.YES {
				if exIP, err = repo.AssignExtranetIP(newDev.SN); err != nil {
					err = fmt.Errorf("SN:%s分配外网失败，err：%v", newDev.SN, err)
					return succeedSNs, totalAffected, err
				}
			} else {
				if _, err = repo.ReleaseIP(newDev.SN, model.Extranet); err != nil {
					err = fmt.Errorf("SN:%s释放外网IP失败，err：%v", newDev.SN, err)
					return succeedSNs, totalAffected, err
				}
			}
			// 根据OS名称关联系统模板ID
			// sysTpl, _ := repo.GetSystemTemplateByName(newDev.OS)
			// 特殊设备无需经过系统部署，自定义写入操作系统名称
			sysTpl, err := repo.GetSystemTemplateByName(newDev.OS)
			if err != nil {
				log.Errorf("get system template by name %s failed(%v) ,adding one", newDev.OS, err)
				tpl := model.SystemTemplate{
					Family:   	"Custom",
					BootMode: 	"uefi",
					Name:     	newDev.OS,
					PXE:        "#NULL",
					Content:    "#NULL",
					OSLifecycle: model.OSTesting,
					Arch:		 model.OSARCHUNKNOWN,
				}
				_, err := repo.SaveSystemTemplate(&tpl)
				if err != nil {
					log.Errorf("add system template by OS name %s fail,%v", newDev.OS, err)
					err = fmt.Errorf("SN:%s 新增操作系统名称（%s）失败：%v", newDev.SN, newDev.OS, err)
					return succeedSNs, totalAffected, err					
				} else {
					//重新获取新增模板ID
					sysTpl, err = repo.GetSystemTemplateByName(newDev.OS)
					if err != nil {
						log.Errorf("get system template by OS name %s fail,%v", newDev.OS, err)
						err = fmt.Errorf("SN:%s 获取操作系统名称（%s）失败：%v", newDev.SN, newDev.OS, err)
						return succeedSNs, totalAffected, err								
					}
				}
			}
		
			//Add device setting(模拟一条数据)
			ds, err := repo.GetDeviceSettingBySN(newDev.SN)
			if ds == nil || err != nil {
				ds = &model.DeviceSetting{
					SN:              newDev.SN,
					Status:          model.InstallStatusSucc, //这个值固定
					InstallType:     model.InstallationPXE,   //这个值固定
					InstallProgress: 1.0,                     //这个值固定
					NeedIntranetIPv6: model.NO,				//默认仅分配IPV4
					NeedExtranetIPv6: model.NO,				//默认仅分配IPV4
				}
			}
			if sysTpl != nil {
				ds.SystemTemplateID = sysTpl.ID
			}
			ds.NeedExtranetIP = newDev.NeedExtranetIPv4
			if inIP != nil {
				ds.IntranetIP = inIP.IP
			}
			if exIP != nil {
				ds.ExtranetIP = exIP.IP
			}
			if err = repo.SaveDeviceSetting(ds); err != nil {
				log.Debug(err)
				return succeedSNs, totalAffected, err
			}
			//特殊设备场景再次更新相关信息
			dev, err := repo.GetDeviceBySN(newDev.SN)
			if err != nil && err != gorm.ErrRecordNotFound {
				return succeedSNs, totalAffected, err
			} 
			if dev != nil {
				if dev.OperationStatus == model.DevOperStatPreDeploy || dev.OperationStatus == model.DevOperStatOnShelve {
					log.Debugf("SpecialDevice SN %s already exist.Status is %s . Allow to update.", newDev.SN, dev.OperationStatus)
					mod.ID = dev.ID
					mod.CreatedAt = dev.CreatedAt
					mod.UpdatedAt = time.Now()
				}
			}
			if _, err = repo.SaveDevice(mod); err != nil {
				log.Debug(err)
				return succeedSNs, totalAffected, err
			}
		}
		//rollback
		defer func() {
			if err != nil {
				_, _ = repo.BatchUpdateServerUSitesStatus([]uint{*mod.USiteID}, model.USiteStatFree)
				_, _ = repo.DeleteDeviceSettingBySN(mod.SN)
				_, _ = repo.RemoveDeviceByID(mod.ID)
				_, _ = repo.RemoveDeviceLifecycleByID(devLifecycle.ID)
				mod = nil
				saveDevLifecycleReq = nil
			}
		}()
		// 若密码为默认出厂密码，则修改
		if checkOriginPassword(log, conf, mod) {
			devices = append(devices, mod)
		}
		// 修改机位占用状态
		if mod.USiteID != nil {
			if _, err = repo.BatchUpdateServerUSitesStatus([]uint{*mod.USiteID}, model.USiteStatUsed); err != nil {
				log.Errorf("update server_usite status failed, usite_num :%s", newDev.USiteNum)
				return succeedSNs, totalAffected, err
			}
		}
		// 统计成功
		totalAffected++
		succeedSNs = append(succeedSNs, newDev.SN)
	}
	//更新关联的订单到货数量和订单状态
	for orderNum, arrivalCount := range mOrderAmount {
		if err = UpdateOrderByArrival(log, repo, orderNum, arrivalCount); err != nil {
			log.Errorf("UpdateOrderByArrival failed, orderNum :%s", orderNum)
			return succeedSNs, totalAffected, err
		}
	}
	// 新增设备成功后，批量修带外改密码
	go batchUpdateOOBPassword(log, repo, conf, devices)
	return succeedSNs, totalAffected, nil
}


// MoveDevice 设备搬迁的详细参数
type MoveDevice struct {
	SN       		string 	`json:"sn"`
	
	//搬迁到机架，上架
	DstServerRoomName 	string	`json:"dst_server_room_name"`   //目的机房
	DstCabinetNumber    string	`json:"dst_cabinet_number"`     //目的机架
	DstUSiteNumber      string	`json:"dst_usite_number"`       //目的机位
	DstPhysicalArea 	string  `json:"dst_physical_area"`      //目标物理区域

	//搬迁到货架，入库
	DstStoreRoomName 	string	`json:"dst_store_room_name"`        //目的库房
	DstVCabinetNumber  	string	`json:"dst_virtual_cabinet_number"` //目的货架

	//搬迁的类型
	MigType 		string 	`json:"mig_type"` //库房->机架，机架->机架,机架->库房
}


// 设备搬迁 请求参数
type BatchMoveDevicesList []*MoveDevice

type BatchMoveDevicesReq struct {
	Devices 		BatchMoveDevicesList
	CurrentUser 	*model.CurrentUser
}


// 搬迁详细参数校验
func (reqData *BatchMoveDevicesReq) validate(repo model.Repo) error {
	// 目标机位唯一性校验
	var dstUsiteIDsForCheck []uint

	for _, data := range reqData.Devices {
		dev, err := repo.GetDeviceBySN(data.SN)
		if gorm.IsRecordNotFoundError(err) {
			return fmt.Errorf("[搬迁参数校验]设备SN: %s 不存在", data.SN)
		}		
		if err != nil {
			return err
		}
		// 状态必须是[搬迁中]
		if dev.OperationStatus != model.DevOperStatMoving {
			return fmt.Errorf("[搬迁参数校验]设备SN: %s 运营状态: %s 不是搬迁中", data.SN, dev.OperationStatus)
		}
		// 目标位置存在性、合法性校验
		if data.MigType == model.MigTypeUsite2Usite || data.MigType == model.MigTypeStore2Usite {
			if data.DstServerRoomName == "" || data.DstCabinetNumber == "" || data.DstUSiteNumber == "" || data.DstPhysicalArea == "" {
				return fmt.Errorf("[搬迁参数校验]设备SN: %s 目标位置参数存在空值", data.SN)
			}
			cond := model.CombinedServerUSite{
				ServerRoomName:  data.DstServerRoomName,
				PhysicalArea:    data.DstPhysicalArea,
				CabinetNumber:   data.DstCabinetNumber,
				USiteNumber:     data.DstUSiteNumber,
			}
			items, err := repo.GetServerUSiteByCond(&cond, nil, nil)
			if err != nil {
				return err
			}
			if len(items) <= 0 {
				return fmt.Errorf("[搬迁参数校验]设备SN: %s  未获取到目标机位 (%s-%s-%s-%s)", data.SN, data.DstServerRoomName, data.DstCabinetNumber, data.DstUSiteNumber, data.DstPhysicalArea)
			} else if len(items) > 1 {
				return fmt.Errorf("[搬迁参数校验]设备SN: %s  获取到多个目标机位 (%s-%s-%s-%s)", data.SN, data.DstServerRoomName, data.DstCabinetNumber, data.DstUSiteNumber, data.DstPhysicalArea)
			}
			if items[0] != nil {
				if items[0].Status == model.USiteStatDisabled || items[0].Status == model.USiteStatUsed {
					return fmt.Errorf("[搬迁参数校验]设备SN: %s  目标机位(%s-%s-%s-%s) 不可用", data.SN, data.DstServerRoomName, data.DstCabinetNumber, data.DstUSiteNumber, data.DstPhysicalArea)
				}				
				dstUsiteIDsForCheck = append(dstUsiteIDsForCheck, items[0].ID)
			}		
		}
	}

	// 目标机位唯一性校验
	if len(dstUsiteIDsForCheck) >= 2 {
		for i := 0; i < len(dstUsiteIDsForCheck); i++ {
			for j := i + 1; j < len(dstUsiteIDsForCheck); j++ {
				if dstUsiteIDsForCheck[i] == dstUsiteIDsForCheck[j] {
					return fmt.Errorf("存在重复的目标机位，请检查")
				}
			}
		}
	}
	return nil
}


// 设备搬迁（释放IP、释放旧机位、占用目标机位、更新状态）
func BatchMoveDevices(log logger.Logger, repo model.Repo, conf *config.Config, reqData *BatchMoveDevicesReq) (succeedSNs []string, totalAffected int64, err error) {
	// 参数校验
	if err := reqData.validate(repo); err != nil {
		return nil, 0, err
	}
	// 存储目标机位ID
	var dstUsiteIDs []uint
	//存储需回收的机位
	var usiteFreeIDs []uint
	var usiteDisabledIDs []uint	
	//1.释放IP的同时保留IP资源，保留期根据config文件配置,0则不保留
	reserveDay := conf.IP.ReserveDay						
	for _, dev := range reqData.Devices {
		if dev.MigType != model.MigTypeStore2Usite && dev.MigType != model.MigTypeStore2Store {
			if reserveDay <= 0 {
				//直接释放IP资源
				if _, err = repo.ReleaseIP(dev.SN, model.Intranet); err != nil {
					log.Error("device(sn=%s) release ip fail", dev.SN)
				}
				//外网IP、IPv6不一定有
				_, _ = repo.ReleaseIP(dev.SN, model.Extranet)
				_, _ = repo.ReleaseIPv6(dev.SN, model.IPScopeIntranet)
				_, _ = repo.ReleaseIPv6(dev.SN, model.IPScopeExtranet)

				//同时删除装机记录
				if _, err := repo.DeleteDeviceSettingBySN(dev.SN); err != nil {
					log.Error("device(sn=%s) clear device setting fail", dev.SN)
				}
			}
			if reserveDay > 0 {
				now := time.Now()
				releaseDate := now.AddDate(0, 0, reserveDay)
				if _, err = repo.ReserveIP(dev.SN, model.Intranet, releaseDate); err != nil {
					log.Error("device(sn=%s) reserve ip fail", dev.SN)
				}
				//外网IP不一定有
				_, _ = repo.ReserveIP(dev.SN, model.Extranet, releaseDate)
				//同时删除装机记录
				if _, err := repo.DeleteDeviceSettingBySN(dev.SN); err != nil {
					log.Error("device(sn=%s) clear device setting fail", dev.SN)
				}
			}
		}

		oriDev, err := repo.GetDeviceBySN(dev.SN)
		if err != nil {
			log.Error("get origin device(sn=%s) usite fail", dev.SN)
			return succeedSNs, totalAffected, err
		}
		// 变更记录
		var details []string
		details = append(details, fmt.Sprintf("序列号:%v", dev.SN))

		// 	搬迁类型细分							
		switch dev.MigType {
		//机架->机架
		case model.MigTypeUsite2Usite:
			cond := model.CombinedServerUSite{
				ServerRoomName:  dev.DstServerRoomName,
				PhysicalArea:    dev.DstPhysicalArea,
				CabinetNumber:   dev.DstCabinetNumber,
				USiteNumber:     dev.DstUSiteNumber,
			}
			items, err := repo.GetServerUSiteByCond(&cond, nil, nil)
			if err != nil {
				log.Error("get dst server(sn=%s) dst ustie (%s-%s-%s-%s) fail(%s)", dev.SN, dev.DstServerRoomName, dev.DstCabinetNumber, dev.DstUSiteNumber, dev.DstPhysicalArea, err.Error())
				return succeedSNs, totalAffected, err
			}
			if items[0] != nil {
				dstUsiteIDs = append(dstUsiteIDs, items[0].ID)
			}
			if oriDev.USiteID != nil {
				if cabinet, err := repo.GetServerCabinetByID(oriDev.CabinetID); err == nil {
					// 若上层机架为[已锁定]，则对应机位不可释放
					if cabinet.Status == model.CabinetStatLocked {
						usiteDisabledIDs = append(usiteDisabledIDs, *oriDev.USiteID)
					} else {
						usiteFreeIDs = append(usiteFreeIDs, *oriDev.USiteID)
					}
				}				
			}

			// 设备生命周期记录操作内容
			oriServerRoom, _ := repo.GetServerRoomByID(oriDev.ServerRoomID)
			oriCabinet, _ := repo.GetServerCabinetByID(oriDev.CabinetID)
			//oriUsite, _ := repo.GetServerUSiteByID(*oriDev.USiteID)
			details = append(details, fmt.Sprintf("搬迁类型:机架->机架"))
			details = append(details, fmt.Sprintf("机房管理单元: %v -> %v", oriServerRoom.Name, dev.DstServerRoomName))
			details = append(details, fmt.Sprintf("机架: %v -> %v", oriCabinet.Number, dev.DstCabinetNumber))
			if oriDev.USiteID != nil {
				oriUsite, _ := repo.GetServerUSiteByID(*oriDev.USiteID)
				details = append(details, fmt.Sprintf("机位: %v -> %v", oriUsite.Number, dev.DstUSiteNumber))
			} else {
				details = append(details,"机位：获取失败")
			}
			//将物理机状态更新为[待部署]
			if _, err = repo.UpdateDeviceBySN(&model.Device{
				SN:              dev.SN,
				IDCID:           items[0].IDCID,
				ServerRoomID:    items[0].ServerRoomID,
				CabinetID:       items[0].ServerCabinetID,
				USiteID:         &items[0].ID,
				OperationStatus: model.DevOperStatPreDeploy,
			}); err != nil {
				log.Error("update device(sn=%s) status pre_deploy fail", dev.SN)
				return succeedSNs, totalAffected, err
			}
			// 统计成功
			totalAffected++
			succeedSNs = append(succeedSNs, dev.SN)
			// 保存设备生命周期记录操作内容
			devLog := model.ChangeLog {
				OperationUser:		reqData.CurrentUser.Name,
				OperationType:		model.OperationTypeMove,
				OperationDetail:	strings.Join(details, "；"),
				OperationTime:		times.ISOTime(time.Now()).ToTimeStr(),
			}
			adll := &AppendDeviceLifecycleLogReq{
				SN:					dev.SN,
				LifecycleLog: 		devLog,
			}
			if err = AppendDeviceLifecycleLogBySN(log, repo, adll);err != nil {
				log.Error("LifecycleLog: append device lifecycle log (sn:%v) fail(%s)", dev.SN, err.Error())
			}

		//机架->库房
		case model.MigTypeUsite2Store:
			if oriDev.USiteID != nil {
				if cabinet, err := repo.GetServerCabinetByID(oriDev.CabinetID); err == nil {
					// 若上层机架为[已锁定]，则对应机位不可释放
					if cabinet.Status == model.CabinetStatLocked {
						usiteDisabledIDs = append(usiteDisabledIDs, *oriDev.USiteID)
					} else {
						usiteFreeIDs = append(usiteFreeIDs, *oriDev.USiteID)
					}
				}				
				//originUsiteIDs = append(originUsiteIDs, *oriDev.USiteID)
			}			
			cond := model.CombinedStoreRoomVirtualCabinet {
				StoreRoomName:			dev.DstStoreRoomName,
				VirtualCabinetNumber: 	dev.DstVCabinetNumber,
			}
			items, err := repo.GetVirtualCabinetsByCond(&cond, nil, nil)
			if err != nil {
				log.Error("server(sn=%s) get dst store room ustie (%s-%s) fail(%s)", dev.SN, dev.DstStoreRoomName, dev.DstVCabinetNumber, err.Error())
				return succeedSNs, totalAffected, err
			}
			if len(items) <= 0 {
				log.Error("server(sn=%s) no dst store room ustie (%s-%s) found", dev.SN, dev.DstStoreRoomName, dev.DstVCabinetNumber)
				return succeedSNs, totalAffected, err
			} else if len(items) > 1 {
				log.Error("server(sn=%s) more than one dst store room ustie (%s-%s) found", dev.SN, dev.DstStoreRoomName, dev.DstVCabinetNumber)
				return succeedSNs, totalAffected, err
			}
			// 设备生命周期记录操作内容
			oriServerRoom, _ := repo.GetServerRoomByID(oriDev.ServerRoomID)
			oriCabinet, _ := repo.GetServerCabinetByID(oriDev.CabinetID)
			details = append(details, fmt.Sprintf("搬迁类型:机架->库房"))
			details = append(details, fmt.Sprintf("机房管理单元->库房: %v -> %v", oriServerRoom.Name, dev.DstStoreRoomName))
			details = append(details, fmt.Sprintf("机架->货架: %v -> %v", oriCabinet.Number, dev.DstVCabinetNumber))
			//将物理机状态更新为[库存中]
			oriDev.IDCID = 0
			oriDev.ServerRoomID = 0
			oriDev.CabinetID = 0
			oriDev.USiteID = nil //&zero
			oriDev.StoreRoomID = items[0].StoreRoomID
			oriDev.VCabinetID = items[0].ID
			oriDev.OperationStatus = model.DevOperStatInStore
			if _, err = repo.SaveDevice(oriDev); err != nil {
				log.Error("update device(sn=%s) status pre_deploy fail", dev.SN)
				return succeedSNs, totalAffected, err
			}
			// 统计成功
			totalAffected++
			succeedSNs = append(succeedSNs, dev.SN)
			// 保存设备生命周期记录操作内容
			devLog := model.ChangeLog {
				OperationUser:		reqData.CurrentUser.Name,
				OperationType:		model.OperationTypeMove,
				OperationDetail:	strings.Join(details, "；"),
				OperationTime:		times.ISOTime(time.Now()).ToTimeStr(),
			}
			adll := &AppendDeviceLifecycleLogReq{
				SN:					dev.SN,
				LifecycleLog: 		devLog,
			}
			if err = AppendDeviceLifecycleLogBySN(log, repo, adll);err != nil {
				log.Error("LifecycleLog: append device lifecycle log (sn:%v) fail(%s)", dev.SN, err.Error())
			}			
		//库房->机架
		case model.MigTypeStore2Usite:
			cond := model.CombinedServerUSite{
				ServerRoomName:  dev.DstServerRoomName,
				PhysicalArea:    dev.DstPhysicalArea,
				CabinetNumber:   dev.DstCabinetNumber,
				USiteNumber:     dev.DstUSiteNumber,
			}
			items, err := repo.GetServerUSiteByCond(&cond, nil, nil)
			if err != nil {
				log.Error("server(sn=%s) get dst ustie (%s-%s-%s-%s) fail(%s)", dev.SN, dev.DstServerRoomName, dev.DstCabinetNumber, dev.DstUSiteNumber, dev.DstPhysicalArea, err.Error())
				return succeedSNs, totalAffected, err
			}
			if items[0] != nil {
				dstUsiteIDs = append(dstUsiteIDs, items[0].ID)
			}
			// 设备生命周期记录操作内容
			oriStoreRoom, _ := repo.GetStoreRoomByID(oriDev.StoreRoomID)
			oriVCabinet, _ := repo.GetVirtualCabinetByID(oriDev.VCabinetID)
			details = append(details, fmt.Sprintf("搬迁类型:库房->机架"))
			details = append(details, fmt.Sprintf("库房->机房管理单元: %v -> %v", oriStoreRoom.Name, dev.DstServerRoomName))
			details = append(details, fmt.Sprintf("货架->机架: %v -> %v", oriVCabinet.Number, dev.DstCabinetNumber))
			details = append(details, fmt.Sprintf("机位:  -> %v", dev.DstUSiteNumber))

			//将物理机状态更新为[待部署]
			oriDev.IDCID = items[0].IDCID
			oriDev.ServerRoomID = items[0].ServerRoomID
			oriDev.CabinetID = items[0].ServerCabinetID
			oriDev.USiteID = &items[0].ID
			oriDev.StoreRoomID = 0
			oriDev.VCabinetID = 0
			now := time.Now()
			oriDev.StartedAt = now
			oriDev.OnShelveAt = now
			oriDev.OperationStatus = model.DevOperStatPreDeploy
			if _, err = repo.SaveDevice(oriDev); err != nil {
				log.Error("update device(sn=%s) status pre_deploy fail", dev.SN)
				return succeedSNs, totalAffected, err
			}
			// 统计成功
			totalAffected++
			succeedSNs = append(succeedSNs, dev.SN)
			// 保存设备生命周期记录操作内容
			devLog := model.ChangeLog {
				OperationUser:		reqData.CurrentUser.Name,
				OperationType:		model.OperationTypeMove,
				OperationDetail:	strings.Join(details, "；"),
				OperationTime:		times.ISOTime(time.Now()).ToTimeStr(),
			}
			adll := &AppendDeviceLifecycleLogReq{
				SN:					dev.SN,
				LifecycleLog: 		devLog,
			}
			if err = AppendDeviceLifecycleLogBySN(log, repo, adll);err != nil {
				log.Error("LifecycleLog: append device lifecycle log (sn:%v) fail(%s)", dev.SN, err.Error())
			}				
		//库房->库房
		case model.MigTypeStore2Store:
			cond := model.CombinedStoreRoomVirtualCabinet {
				StoreRoomName:			dev.DstStoreRoomName,
				VirtualCabinetNumber: 	dev.DstVCabinetNumber,
			}
			items, err := repo.GetVirtualCabinetsByCond(&cond, nil, nil)
			if err != nil {
				log.Error("server(sn=%s) get dst store room ustie (%s-%s) fail(%s)", dev.SN, dev.DstStoreRoomName, dev.DstVCabinetNumber, err.Error())
				return succeedSNs, totalAffected, err
			}
			if len(items) <= 0 {
				log.Error("server(sn=%s) no dst store room ustie (%s-%s) found", dev.SN, dev.DstStoreRoomName, dev.DstVCabinetNumber)
				return succeedSNs, totalAffected, err
			} else if len(items) > 1 {
				log.Error("server(sn=%s) more than one dst store room ustie (%s-%s) found", dev.SN, dev.DstStoreRoomName, dev.DstVCabinetNumber)
				return succeedSNs, totalAffected, err
			}
			// 设备生命周期记录操作内容
			oriStoreRoom, _ := repo.GetStoreRoomByID(oriDev.StoreRoomID)
			oriVCabinet, _ := repo.GetVirtualCabinetByID(oriDev.VCabinetID)
			details = append(details, fmt.Sprintf("搬迁类型:库房->库房"))
			details = append(details, fmt.Sprintf("库房: %v -> %v", oriStoreRoom.Name, dev.DstStoreRoomName))
			details = append(details, fmt.Sprintf("货架: %v -> %v", oriVCabinet.Number, dev.DstVCabinetNumber))
			// 设备记录更新
			oriDev.IDCID = 0
			oriDev.StoreRoomID = items[0].StoreRoomID
			oriDev.VCabinetID = items[0].ID
			oriDev.OperationStatus = model.DevOperStatInStore
			if _, err = repo.UpdateDeviceBySN(oriDev); err != nil {
				log.Error("update device(sn=%s) status pre_deploy fail", dev.SN)
				return succeedSNs, totalAffected, err
			}
			// 统计成功
			totalAffected++
			succeedSNs = append(succeedSNs, dev.SN)
			// 保存设备生命周期记录操作内容
			devLog := model.ChangeLog {
				OperationUser:		reqData.CurrentUser.Name,
				OperationType:		model.OperationTypeMove,
				OperationDetail:	strings.Join(details, "；"),
				OperationTime:		times.ISOTime(time.Now()).ToTimeStr(),
			}
			adll := &AppendDeviceLifecycleLogReq{
				SN:					dev.SN,
				LifecycleLog: 		devLog,
			}
			if err = AppendDeviceLifecycleLogBySN(log, repo, adll);err != nil {
				log.Error("LifecycleLog: append device lifecycle log (sn:%v) fail(%s)", dev.SN, err.Error())
			}				
		}
	}
	//3.将目标机位标记为[已占用]
	if _, err = repo.BatchUpdateServerUSitesStatus(dstUsiteIDs, model.USiteStatUsed); err != nil {
		log.Error("occupy usites(ids=%v) fail", dstUsiteIDs)
		return succeedSNs, totalAffected, err
	}
	//4.将原机位标记为空闲 or 不可用
	if _, err = repo.BatchUpdateServerUSitesStatus(usiteFreeIDs, model.USiteStatFree); err != nil {
		log.Error("free usites(ids=%v) fail", usiteFreeIDs)
		return succeedSNs, totalAffected, err
	}
	if _, err = repo.BatchUpdateServerUSitesStatus(usiteDisabledIDs, model.USiteStatDisabled); err != nil {
		log.Error("disable usites(ids=%v) fail", usiteDisabledIDs)
		return succeedSNs, totalAffected, err
	}	
	return succeedSNs, totalAffected, nil
}


// 设备退役 请求参数
type BatchRetireDevicesList []string

type BatchRetireDevicesReq struct {
	SNs 			BatchRetireDevicesList
	CurrentUser 	*model.CurrentUser
}

// 设备退役（释放IP、释放旧机位、更新状态）
func BatchRetireDevices(log logger.Logger, repo model.Repo, conf *config.Config, reqData *BatchRetireDevicesReq) (succeedSNs []string, totalAffected int64, err error) {
	//存储需释放的机位
	var usiteFreeIDs []uint
	var usiteDisabledIDs []uint
	//1.释放IP的同时保留IP资源，保留期根据config文件配置,0则不保留
	reserveDay := conf.IP.ReserveDay	
	for _, sn := range reqData.SNs {
		dev, err := repo.GetDeviceBySN(sn)
		if err != nil {
			log.Errorf("get device by sn:%s fail, %s", sn, err.Error())
			return succeedSNs, totalAffected, err
		}
		// 状态必须是[退役中] 
		if dev.OperationStatus != model.DevOperStatRetiring {
			return succeedSNs, totalAffected, fmt.Errorf("[退役参数校验]设备SN: %s 运营状态: %s 不是退役中", sn, dev.OperationStatus)
		}
		if dev != nil && dev.USiteID != nil {
			if cabinet, err := repo.GetServerCabinetByID(dev.CabinetID); err == nil {
				// 若上层机架为[已锁定]，则对应机位不可释放
				if cabinet.Status == model.CabinetStatLocked {
					usiteDisabledIDs = append(usiteDisabledIDs, *dev.USiteID)
				} else {
					usiteFreeIDs = append(usiteFreeIDs, *dev.USiteID)
				}
			}
			//usiteIDs = append(usiteIDs, *dev.USiteID)
		}
		// 变更记录
		optDetail, err := convert2DetailOfOperationTypeRetire(repo, dev)
		if err != nil {
			log.Errorf("Fail to convert Detail of OperationTypeRetire: %v", err)
		}
		devLog := model.ChangeLog {
			OperationUser:		reqData.CurrentUser.Name,
			OperationType:		model.OperationTypeRetire,
			OperationDetail:	optDetail,
			OperationTime:		times.ISOTime(time.Now()).ToTimeStr(),
		}
		adll := &AppendDeviceLifecycleLogReq{
			SN:					dev.SN,
			LifecycleLog: 		devLog,
		}
		if err = AppendDeviceLifecycleLogBySN(log, repo, adll);err != nil {
			log.Error("LifecycleLog: append device lifecycle log (sn:%v) fail(%s)", dev.SN, err.Error())
		}							
		// 释放机位并修改状态
		dev.USiteID = nil
		dev.PowerStatus = model.PowerStatusOff
		dev.OperationStatus = model.DevOperStateRetired
		if _, err := repo.SaveDevice(dev); err != nil {
			log.Error(err)
			return succeedSNs, totalAffected, err
		}
		// 释放IP资源
		if reserveDay <= 0 {
			//直接释放IP资源
			if _, err = repo.ReleaseIP(dev.SN, model.Intranet); err != nil {
				log.Error("device(sn=%s) release ip fail", dev.SN)
			}
			//外网IP、IPv6不一定有
			_, _ = repo.ReleaseIP(dev.SN, model.Extranet)
			_, _ = repo.ReleaseIPv6(dev.SN, model.IPScopeIntranet)
			_, _ = repo.ReleaseIPv6(dev.SN, model.IPScopeExtranet)

			//同时删除装机记录
			if _, err := repo.DeleteDeviceSettingBySN(dev.SN); err != nil {
				log.Error("device(sn=%s) clear device setting fail", dev.SN)
			}
		}
		if reserveDay > 0 {
			now := time.Now()
			releaseDate := now.AddDate(0, 0, reserveDay)
			if _, err = repo.ReserveIP(dev.SN, model.Intranet, releaseDate); err != nil {
				log.Error("device(sn=%s) reserve ip fail", dev.SN)
			}
			//外网IP不一定有
			_, _ = repo.ReserveIP(dev.SN, model.Extranet, releaseDate)
			//同时删除装机记录
			if _, err := repo.DeleteDeviceSettingBySN(dev.SN); err != nil {
				log.Error("device(sn=%s) clear device setting fail", dev.SN)
			}
		}
		// 统计成功
		totalAffected++
		succeedSNs = append(succeedSNs, dev.SN)
		// 记录设备退役日期
		saveDevLifecycleReq := &SaveDeviceLifecycleReq {
			DeviceLifecycleBase: DeviceLifecycleBase{
				DeviceRetiredDate:		time.Now(),
			},
		}
		// DeviceLifecycle 查询是否已经存在
		devLifecycle, err := repo.GetDeviceLifecycleBySN(dev.SN)
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Error(err)
		} 
		if devLifecycle != nil {
			saveDevLifecycleReq.ID = devLifecycle.ID
			// 保存或更新 DeviceLifecycle
			if err = SaveDeviceLifecycle(log, repo, saveDevLifecycleReq); err != nil {
				log.Error(err)
			}
		}
	}
	if _, err = repo.BatchUpdateServerUSitesStatus(usiteFreeIDs, model.USiteStatFree); err != nil {
		log.Error("free usites(ids=%v) fail", usiteFreeIDs)
		return succeedSNs, totalAffected, err
	}
	if _, err = repo.BatchUpdateServerUSitesStatus(usiteDisabledIDs, model.USiteStatDisabled); err != nil {
		log.Error("disable usites(ids=%v) fail", usiteDisabledIDs)
		return succeedSNs, totalAffected, err
	}
	return succeedSNs, totalAffected, nil
}