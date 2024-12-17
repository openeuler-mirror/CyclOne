package service

import (
	"idcos.io/cloudboot/utils/times"
	"strings"
	"fmt"
	"reflect"
	"net/http"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"time"

	"idcos.io/cloudboot/config"
	"github.com/voidint/binding"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
)


//SaveDeviceCategoryReq 保存设备生命周期详情请求参数
type SaveDeviceLifecycleReq struct {
	DeviceLifecycleBase
	ID uint `json:"id"`
	// 用户登录名
	LoginName string `json:"-"`
}

// 基本字段
type DeviceLifecycleBase struct {
	FixedAssetNumber 				string    	`json:"fixed_asset_number"`
	SN               				string    	`json:"sn"`
	AssetBelongs	 				string		`json:"asset_belongs"`
	Owner			 				string		`json:"owner"`
	IsRental		 				string		`json:"is_rental"`
	MaintenanceServiceProvider		string		`json:"maintenance_service_provider"`
	MaintenanceService				string		`json:"maintenance_service"`
	LogisticsService				string		`json:"logistics_service"`
	MaintenanceServiceDateBegin     time.Time 	`json:"maintenance_service_date_begin"`
	MaintenanceServiceDateEnd       time.Time 	`json:"maintenance_service_date_end"`
	MaintenanceServiceStatus		string		`json:"maintenance_service_status"`
	DeviceRetiredDate       		time.Time 	`json:"device_retired_date"`
	LifecycleLog					string		`json:"lifecycle_log"`
}

//SaveDeviceLifecycle 保存/更新 设备生命周期记录
func SaveDeviceLifecycle(log logger.Logger, repo model.Repo, reqData *SaveDeviceLifecycleReq) error {
	// DeviceLifecycle 结构体
	modDevLifecycle := &model.DeviceLifecycle{
		FixedAssetNumber: 				reqData.FixedAssetNumber,
		SN:             				reqData.SN,
		AssetBelongs:	 				reqData.AssetBelongs,			
		Owner:				 			reqData.Owner,			
		IsRental:		 				reqData.IsRental,
		MaintenanceServiceProvider: 	reqData.MaintenanceServiceProvider, 		
		MaintenanceService:				reqData.MaintenanceService,				
		LogisticsService:				reqData.LogisticsService,
		MaintenanceServiceDateBegin:	reqData.MaintenanceServiceDateBegin, 		
		MaintenanceServiceDateEnd:		reqData.MaintenanceServiceDateEnd, 		
		MaintenanceServiceStatus:		reqData.MaintenanceServiceStatus,
		DeviceRetiredDate:              reqData.DeviceRetiredDate, // sql_mode=STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE 时因parseTime=True强制转换零值导致无法写入DB
		LifecycleLog:					reqData.LifecycleLog,
	}

	modDevLifecycle.Model.ID = reqData.ID

	_, err := repo.SaveDeviceLifecycle(modDevLifecycle)
	if err != nil {
		log.Errorf("SaveDeviceLifecycle failed: %s", err.Error())
		return err
	}

	reqData.ID = modDevLifecycle.Model.ID
	return err
}

// 设备详情页展示结构体
type DeviceLifecycleDetailPage struct {
	AssetBelongs	 				string					`json:"asset_belongs"`
	Owner			 				string					`json:"owner"`
	IsRental		 				string					`json:"is_rental"`
	MaintenanceServiceProvider		string					`json:"maintenance_service_provider"`
	MaintenanceService				string					`json:"maintenance_service"`
	LogisticsService				string					`json:"logistics_service"`
	MaintenanceServiceDateBegin     string 					`json:"maintenance_service_date_begin"`
	MaintenanceServiceDateEnd       string 					`json:"maintenance_service_date_end"`
	MaintenanceServiceStatus		string					`json:"maintenance_service_status"`
	DeviceRetiredDate       		string 					`json:"device_retired_date"`
	LifecycleLog					[]model.ChangeLog		`json:"lifecycle_log"`
}


// 更新追加设备生命周期变更记录 请求结构体
type AppendDeviceLifecycleLogReq struct {
	SN					string				`json:"sn"`
	LifecycleLog		model.ChangeLog		`json:"lifecycle_log"`
}

//AppendDeviceLifecycleLog 更新追加设备生命周期变更记录
func AppendDeviceLifecycleLogBySN(log logger.Logger, repo model.Repo, reqData *AppendDeviceLifecycleLogReq) error {
	// DeviceLifecycle 查询是否已经存在
	devLifecycle, err := repo.GetDeviceLifecycleBySN(reqData.SN)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error(err.Error())
		return err
	}
	if devLifecycle != nil {
		log.Debugf("Begin to AppendDeviceLifecycleLogBySN:%s", reqData.SN)
		defer log.Debugf("End to AppendDeviceLifecycleLogBySN:%s", reqData.SN)
		// 获取当前的生命周期日志记录
		var devLL []model.ChangeLog
		if devLifecycle.LifecycleLog != "" {
			if err = json.Unmarshal([]byte(devLifecycle.LifecycleLog), &devLL);err != nil {
				log.Error(err.Error())
				return err
			}
		}
		// 追加
		devLL = append(devLL, reqData.LifecycleLog)
		b, err := json.Marshal(devLL)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		// DeviceLifecycle 结构体
		modDevLifecycle := &model.DeviceLifecycle{
			SN:             				reqData.SN,
			LifecycleLog:					string(b),
		}
		if err = repo.UpdateDeviceLifecycleBySN(modDevLifecycle);err != nil {
			log.Errorf("UpdateDeviceLifecycleBySN failed: %s", err.Error())
			return err
		}
	}
	return nil
}


// 转换结构体属性值以便存于 ChangeLog.OperationDetail
func convert2DetailOfOperationTypeAdd(repo model.Repo, mdev model.Device) (string, error) {
	
	getType := reflect.TypeOf(mdev)
	getValue := reflect.ValueOf(mdev)
	var details []string

	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i)
		value := getValue.Field(i).Interface()

		switch field.Name {
		case "FixedAssetNumber":
			details = append(details, fmt.Sprintf("固资编号:%v", value))
		case "SN":
			details = append(details, fmt.Sprintf("序列号:%v", value))
		case "Vendor":
			details = append(details, fmt.Sprintf("厂商:%v", value))
		case "Category":
			details = append(details, fmt.Sprintf("设备类型:%v", value))
		case "HardwareRemark":
			details = append(details, fmt.Sprintf("硬件备注:%v", value))
		case "ServerRoomID":
			if unitValue, ok := value.(uint); ok {
				serverRoom, _ := repo.GetServerRoomByID(unitValue)
				if serverRoom != nil {
					details = append(details, fmt.Sprintf("机房管理单元:%v", serverRoom.Name))
				}
			} else {
				details = append(details, fmt.Sprintf("机房管理单元:%v", value))
			} 
		case "CabinetID":
			if unitValue, ok := value.(uint); ok {
				cabinet, _ := repo.GetServerCabinetByID(unitValue)
				if cabinet != nil {
					details = append(details, fmt.Sprintf("机架编号:%v", cabinet.Number))
				}
			} else {
				details = append(details, fmt.Sprintf("机架编号:%v", value))
			}
		case "USiteID":
			if unitValue, ok := value.(*uint); ok {
				if  unitValue != nil { // 避免空指针
					usite, _ := repo.GetServerUSiteByID(*unitValue)
					if usite != nil {
						details = append(details, fmt.Sprintf("机位编号:%v", usite.Number))
					}
				}
			} else {
				details = append(details, fmt.Sprintf("机位编号:%v", value))
			}
	    case "StoreRoomID":
	    	if unitValue, ok := value.(uint); ok {
	    		storeRoom, _ := repo.GetStoreRoomByID(unitValue)
	    		if storeRoom != nil {
	    			details = append(details, fmt.Sprintf("库房管理单元:%v", storeRoom.Name))
	    		}
	    	} else {
	    		details = append(details, fmt.Sprintf("库房管理单元:%v", value))
	    	} 
	    case "VCabinetID":
	    	if unitValue, ok := value.(uint); ok {
	    		vcabinet, _ := repo.GetVirtualCabinetByID(unitValue)
	    		if vcabinet != nil {
	    			details = append(details, fmt.Sprintf("货架编号:%v", vcabinet.Number))
	    		}
	    	} else {
	    		details = append(details, fmt.Sprintf("货架编号:%v", value))
			}
		}
	}

	if len(details) > 0 {
		return strings.Join(details, "；"), nil
	} else {
		return "", fmt.Errorf("新增设备记录（SN：%v）字段详情转化失败", mdev.SN)
	}
}


// 转换结构体属性值以便存于 ChangeLog.OperationDetail
func convert2DetailOfOperationTypeMove(repo model.Repo, dev *SubmitDeviceMigrationApprovalReqData) (string, error) {
	var details []string

	oriDev, err := repo.GetDeviceBySN(dev.SN)
	if err != nil {
		return "", fmt.Errorf("搬迁设备记录（SN：%v）获取旧详情失败", dev.SN)
	}

	oriServerRoom, _ := repo.GetServerRoomByID(oriDev.ServerRoomID)
	oriCabinet, _ := repo.GetServerCabinetByID(oriDev.CabinetID)
	//oriUsite, _ := repo.GetServerUSiteByID(*oriDev.USiteID)
	oriStoreRoom, _ := repo.GetStoreRoomByID(oriDev.StoreRoomID)
	oriVCabinet, _ := repo.GetVirtualCabinetByID(oriDev.VCabinetID)
	dstServerRoom, _ := repo.GetServerRoomByID(dev.DstServerRoomID)
	dstCabinet, _ := repo.GetServerCabinetByID(dev.DstCabinetID)
	//dstUsite, _ := repo.GetServerUSiteByID(dev.DstUSiteID)
	dstStoreRoom, _ := repo.GetStoreRoomByID(dev.DstStoreRoomID)
	dstVCabinet, _ := repo.GetVirtualCabinetByID(dev.DstVCabinetID)

	details = append(details, fmt.Sprintf("序列号:%v", dev.SN))

	switch dev.MigType {
	case model.MigTypeUsite2Usite:
		details = append(details, fmt.Sprintf("搬迁类型:机架->机架"))
		// 机房管理单元
		if oriServerRoom != nil && dstServerRoom != nil {
			details = append(details, fmt.Sprintf("机房管理单元: %v -> %v", oriServerRoom.Name, dstServerRoom.Name))
		}
		// 机架
		if oriCabinet != nil && dstCabinet != nil {
			details = append(details, fmt.Sprintf("机架: %v -> %v", oriCabinet.Number, dstCabinet.Number))
		}
		// 机位
		if oriDev.USiteID != nil && dev.DstUSiteID != 0 {
			oriUsite, _ := repo.GetServerUSiteByID(*oriDev.USiteID)
			dstUsite, _ := repo.GetServerUSiteByID(dev.DstUSiteID)
			if oriUsite != nil && dstUsite != nil {
				details = append(details, fmt.Sprintf("机位: %v -> %v", oriUsite.Number, dstUsite.Number))
			}
		} else {
			details = append(details, "机位: 获取失败")
		}

	case model.MigTypeUsite2Store:
		details = append(details, fmt.Sprintf("搬迁类型:机架->库房"))
		// 库房
		if oriServerRoom != nil && dstStoreRoom != nil {
			details = append(details, fmt.Sprintf("机房管理单元->库房: %v -> %v", oriServerRoom.Name, dstStoreRoom.Name))
		}
		// 货架
		if oriCabinet != nil && dstVCabinet != nil {
			details = append(details, fmt.Sprintf("机架->货架: %v -> %v", oriCabinet.Number, dstVCabinet.Number))
		}

	case model.MigTypeStore2Usite:
		details = append(details, fmt.Sprintf("搬迁类型:库房->机架"))
		if oriStoreRoom != nil && dstServerRoom != nil {
			details = append(details, fmt.Sprintf("库房->机房管理单元: %v -> %v", oriStoreRoom.Name, dstServerRoom.Name))
		}
		if oriVCabinet != nil && dstCabinet != nil {
			details = append(details, fmt.Sprintf("货架->机架: %v -> %v", oriVCabinet.Number, dstCabinet.Number))
		}
		if dev.DstUSiteID != 0 {
			dstUsite, _ := repo.GetServerUSiteByID(dev.DstUSiteID)
			if dstUsite != nil {
				details = append(details, fmt.Sprintf("机位: %v", dstUsite.Number))
			}
		} else {
			details = append(details, "机位: 获取失败")
		}		

	case model.MigTypeStore2Store:
		details = append(details, fmt.Sprintf("搬迁类型:库房->库房"))
		if oriStoreRoom != nil && dstStoreRoom != nil {
			details = append(details, fmt.Sprintf("库房: %v -> %v", oriStoreRoom.Name, dstStoreRoom.Name))
		}
		if oriVCabinet != nil && dstVCabinet != nil {
			details = append(details, fmt.Sprintf("货架: %v -> %v", oriVCabinet.Number, dstVCabinet.Number))
		}
	}

	if len(details) > 0 {
		return strings.Join(details, "；"), nil
	} else {
		return "", fmt.Errorf("搬迁设备记录（SN：%v）字段详情转化失败", dev.SN)
	}
}


// 转换结构体属性值以便存于 ChangeLog.OperationDetail
func convert2DetailOfOperationTypeRetire(repo model.Repo, mdev *model.Device) (string, error) {
	var details []string

	details = append(details, fmt.Sprintf("序列号:%v", mdev.SN))
	
	if oriServerRoom, _ := repo.GetServerRoomByID(mdev.ServerRoomID); oriServerRoom != nil {
		details = append(details, fmt.Sprintf("机房管理单元: %v", oriServerRoom.Name))
	}

	if oriCabinet, _ := repo.GetServerCabinetByID(mdev.CabinetID); oriCabinet != nil {
		details = append(details, fmt.Sprintf("机架: %v", oriCabinet.Number))
	}
	
	if mdev.USiteID != nil {
		if oriUsite, _ := repo.GetServerUSiteByID(*mdev.USiteID); oriUsite != nil {
			details = append(details, fmt.Sprintf("释放机位: %v", oriUsite.Number))
		}
	}

	if oriStoreRoom, _ := repo.GetStoreRoomByID(mdev.StoreRoomID); oriStoreRoom != nil {
		details = append(details, fmt.Sprintf("库房: %v", oriStoreRoom.Name))
	}
	
	if oriVCabinet, _ := repo.GetVirtualCabinetByID(mdev.VCabinetID);oriVCabinet != nil {
		details = append(details, fmt.Sprintf("货架: %v", oriVCabinet.Number))
	}

	if len(details) > 0 {
		return strings.Join(details, "；"), nil
	} else {
		return "", fmt.Errorf("退役设备记录（SN：%v）字段详情转化失败", mdev.SN)
	}
}


func convert2DetailOfOperationTypeUpdate(repo model.Repo, mdev model.Device) (string, error) {
	var details []string
	getType := reflect.TypeOf(mdev)
	getValue := reflect.ValueOf(mdev)
	// 仅处理根据SN更新字段的场景
	oriDev, err := repo.GetDeviceBySN(mdev.SN)
	if err != nil {
		return "", fmt.Errorf("更新设备记录（SN：%v）获取旧详情失败", mdev.SN)
	}

	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i)
		value := getValue.Field(i).Interface()

		switch field.Name {
		case "FixedAssetNumber":
			if strValue, ok := value.(string); ok {
				if strValue != "" {
					details = append(details, fmt.Sprintf("固资编号:%v -> %v", oriDev.FixedAssetNumber, value))
				}
			} else {
				details = append(details, fmt.Sprintf("固资编号:%v", value))
			}
		case "OperationStatus":
			if strValue, ok := value.(string); ok {
				if strValue != "" {
					details = append(details, fmt.Sprintf("运营状态:%v -> %v", OperationStatusTransfer(oriDev.OperationStatus, true), OperationStatusTransfer(strValue, true)))
				}
			} else {
				details = append(details, fmt.Sprintf("运营状态:%v", value))
			}
		case "Usage":
			if strValue, ok := value.(string); ok {
				if strValue != "" {
					details = append(details, fmt.Sprintf("用途:%v -> %v", oriDev.Usage, value))
				}
			} else {
				details = append(details, fmt.Sprintf("用途:%v", value))
			}
		}
	}

	if len(details) > 0 {
		return strings.Join(details, "；"), nil
	} else {
		return "", fmt.Errorf("退役设备记录（SN：%v）字段详情转化失败", mdev.SN)
	}

}


func convert2DetailOfOperationTypeOSInstall(repo model.Repo, setting SaveDeviceSettingItem) (string, error) {
	var details []string

	getType := reflect.TypeOf(setting)
	getValue := reflect.ValueOf(setting)

	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i)
		value := getValue.Field(i).Interface()

		switch field.Name {
		case "InstallType":
			details = append(details, fmt.Sprintf("安装方式:%v ", value))
		case "OSTemplateName":
			details = append(details, fmt.Sprintf("操作系统安装模板:%v ", value))
		case "HardwareTemplateName":
			details = append(details, fmt.Sprintf("硬件配置模板:%v ", value))
		case "NeedExtranetIP":
			details = append(details, fmt.Sprintf("是否需要外网IPv4:%v ", value))
		case "NeedExtranetIPv6":
			details = append(details, fmt.Sprintf("是否需要内网IPv6:%v ", value))
		case "NeedIntranetIPv6":
			details = append(details, fmt.Sprintf("是否需要外网IPv6:%v ", value))
		}
	}

	if len(details) > 0 {
		return strings.Join(details, "；"), nil
	} else {
		return "", fmt.Errorf("部署设备记录（SN：%v）字段详情转化失败", setting.SN)
	}

}


// GetDeviceLifecycleBySN
func GetDeviceLifecycleBySN(log logger.Logger, repo model.Repo, sn string) (*DeviceLifecycleDetailPage, error) {
	devLifecycle, err := repo.GetDeviceLifecycleBySN(sn)
	if err != nil {
		log.Errorf("fail to get device lifecycle by sn(%s), err: %v", sn, err)
		return nil, err
	}
	devL := &DeviceLifecycleDetailPage {
		AssetBelongs	 			:	devLifecycle.AssetBelongs,	 					
		Owner			 			:	devLifecycle.Owner,			 					
		IsRental		 			:	devLifecycle.IsRental,
		MaintenanceServiceProvider	:	devLifecycle.MaintenanceServiceProvider,
		MaintenanceService			:	devLifecycle.MaintenanceService,		
		LogisticsService			:	devLifecycle.LogisticsService,				
		MaintenanceServiceDateBegin :	times.ISOTime(devLifecycle.MaintenanceServiceDateBegin).ToDateStr(),
		MaintenanceServiceDateEnd   :	times.ISOTime(devLifecycle.MaintenanceServiceDateEnd).ToDateStr(),        
		MaintenanceServiceStatus	:	devLifecycle.MaintenanceServiceStatus,
		DeviceRetiredDate       	:	times.ISOTime(devLifecycle.DeviceRetiredDate).ToDateStr(),			
	}
	return devL, nil
}

//UpdateDeviceLifecycleReq
type UpdateDeviceLifecycleReq struct {
	// 用户登录名
	LoginName string `json:"-"`
	// 不支持更新的字段
	SN               				string    	`json:"sn"`
	ID								uint 		`json:"id"`
	FixedAssetNumber 				string    	`json:"fixed_asset_number"`
	DeviceRetiredDate       		time.Time 	`json:"device_retired_date"`
	LifecycleLog					string		`json:"lifecycle_log"`	
	//更新字段
	AssetBelongs	 				string		`json:"asset_belongs"`
	Owner			 				string		`json:"owner"`
	IsRental		 				string		`json:"is_rental"`
	MaintenanceServiceProvider		string		`json:"maintenance_service_provider"`
	MaintenanceService				string		`json:"maintenance_service"`
	LogisticsService				string		`json:"logistics_service"`
	MaintenanceServiceDateBegin     string 		`json:"maintenance_service_date_begin"`
	MaintenanceServiceDateEnd       string 		`json:"maintenance_service_date_end"`
	MaintenanceServiceStatus		string		`json:"maintenance_service_status"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *UpdateDeviceLifecycleReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.AssetBelongs: 							"asset_belongs",
		&reqData.Owner: 								"owner",
		&reqData.IsRental:    							"is_rental",
		&reqData.MaintenanceServiceProvider: 			"maintenance_service_provider",
		&reqData.MaintenanceService: 					"maintenance_service",
		&reqData.LogisticsService:    					"logistics_service",
		&reqData.MaintenanceServiceDateBegin: 			"maintenance_service_date_begin",
		&reqData.MaintenanceServiceDateEnd: 			"maintenance_service_date_end",
	}
}

// Validate 结构体数据校验
func (reqData *UpdateDeviceLifecycleReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())
	//SN不能为空
	if reqData.SN != "" {
		// DeviceLifecycle 查询是否已经存在
		devLifecycle, err := repo.GetDeviceLifecycleBySN(reqData.SN)
		if err != nil && err != gorm.ErrRecordNotFound {
			errs.Add([]string{"sn"}, binding.RequiredError, fmt.Sprintf("查询设备SN(%s)生命周期信息失败", reqData.SN))
			return errs
		}
		if devLifecycle != nil {
			// 通过SN转换为实际更新的model.ID
			reqData.ID = devLifecycle.ID
		} else {
			errs.Add([]string{"sn"}, binding.RequiredError, fmt.Sprintf("查询设备SN(%s)生命周期信息不存在", reqData.SN))
			return errs
		}

	} else {
		errs.Add([]string{"sn"}, binding.RequiredError, fmt.Sprintf("查询设备SN(%s)信息不允许为空", reqData.SN))
		return errs
	}
	// 校验是否租赁字段
	if reqData.IsRental != "" {
		if reqData.IsRental != model.YESRental &&  reqData.IsRental != model.NORental {
			errs.Add([]string{"is_rental"}, binding.RequiredError, fmt.Sprintf("是否租赁(%s)必须为 yes 或 no", reqData.IsRental))
			return errs
		}
	}
	if reqData.MaintenanceServiceDateBegin == "" {
		errs.Add([]string{"maintenance_service_date_begin"}, binding.RequiredError, fmt.Sprintf("维保起始日期(%s)不允许为空", reqData.MaintenanceServiceDateBegin))
		return errs
	}
	if !strings.Contains(reqData.MaintenanceServiceDateBegin, "-") {
		errs.Add([]string{"maintenance_service_date_begin"}, binding.RequiredError, fmt.Sprintf("维保起始日期(%s)格式须为：YYYY-MM-DD", reqData.MaintenanceServiceDateBegin))
		return errs
	}
	if reqData.MaintenanceServiceDateEnd == "" {
		errs.Add([]string{"maintenance_service_date_end"}, binding.RequiredError, fmt.Sprintf("维保截止日期(%s)不允许为空", reqData.MaintenanceServiceDateEnd))
		return errs
	}
	if !strings.Contains(reqData.MaintenanceServiceDateEnd, "-") {
		errs.Add([]string{"maintenance_service_date_end"}, binding.RequiredError, fmt.Sprintf("维保截止日期(%s)格式须为：YYYY-MM-DD", reqData.MaintenanceServiceDateEnd))
		return errs
	}	
	return nil
}

//UpdateDeviceLifecycleBySN 更新 设备生命周期记录(for 单台设备延保)
func UpdateDeviceLifecycleBySN(log logger.Logger, repo model.Repo, reqData *UpdateDeviceLifecycleReq) error {
	// DeviceLifecycle 结构体
	modDevLifecycle := &model.DeviceLifecycle{
		SN:             				reqData.SN,
		AssetBelongs:	 				reqData.AssetBelongs,			
		Owner:				 			reqData.Owner,			
		IsRental:		 				reqData.IsRental,
		MaintenanceServiceProvider: 	reqData.MaintenanceServiceProvider,
		MaintenanceService:				reqData.MaintenanceService,	
		LogisticsService:				reqData.LogisticsService,
		//MaintenanceServiceStatus:		reqData.MaintenanceServiceStatus,
		//DeviceRetiredDate:			reqData.DeviceRetiredDate,
		//LifecycleLog:					reqData.LifecycleLog,
	}

	modDevLifecycle.Model.ID = reqData.ID
	modDevLifecycle.MaintenanceServiceDateBegin, _ = time.Parse(times.DateLayout, reqData.MaintenanceServiceDateBegin)
	modDevLifecycle.MaintenanceServiceDateEnd, _ = time.Parse(times.DateLayout, reqData.MaintenanceServiceDateEnd)

	now := time.Now()
	currentDate := now.Format("2006-01-02")
	if now.Before(modDevLifecycle.MaintenanceServiceDateBegin) {
		log.Debugf("current_date=%s < maintenance_service_begin_date=%s. Update Device(SN:%s) set status=%s -> %s", 
			currentDate, modDevLifecycle.MaintenanceServiceDateBegin, reqData.SN, modDevLifecycle.MaintenanceServiceStatus, model.MaintenanceServiceStatusInactive)
		modDevLifecycle.MaintenanceServiceStatus = model.MaintenanceServiceStatusInactive
	} else if now.After(modDevLifecycle.MaintenanceServiceDateEnd) {
		log.Debugf("current_date=%s > maintenance_service_end_date=%s. Update Device(%s) set status=%s -> %s", 
			currentDate, modDevLifecycle.MaintenanceServiceDateEnd, reqData.SN, modDevLifecycle.MaintenanceServiceStatus, model.MaintenanceServiceStatusOutOfWarranty)
		modDevLifecycle.MaintenanceServiceStatus = model.MaintenanceServiceStatusOutOfWarranty
	} else {
		log.Debugf("current_date=%s between  maintenance_service_begin_date=%s and maintenance_service_end_date=%s. Update Device(%s) set status=%s -> %s",
			currentDate, modDevLifecycle.MaintenanceServiceDateBegin, modDevLifecycle.MaintenanceServiceDateEnd, reqData.SN, modDevLifecycle.MaintenanceServiceStatus, model.MaintenanceServiceStatusUnderWarranty)
		modDevLifecycle.MaintenanceServiceStatus = model.MaintenanceServiceStatusUnderWarranty
	}
	
	_, err := repo.SaveDeviceLifecycle(modDevLifecycle)
	if err != nil {
		log.Errorf("SaveDeviceLifecycle failed: %s", err.Error())
		return err
	}
	return err
}

//BatchUpdateDeviceLifecycles 批量修改设备生命周期记录结构
type BatchUpdateDeviceLifecycles struct {
	DeviceLifecycles []UpdateDeviceLifecycles `json:"device_lifecycles"`
	// 用户登录名
	LoginName string `json:"-"`
}

// UpdateDeviceLifecycles 批量修改设备生命周期记录结构
type UpdateDeviceLifecycles struct {
	SN                              string      `json:"sn"`
	AssetBelongs	 				string		`json:"asset_belongs"`
	Owner			 				string		`json:"owner"`
	IsRental		 				string		`json:"is_rental"`
	MaintenanceServiceProvider		string		`json:"maintenance_service_provider"`
	MaintenanceService				string		`json:"maintenance_service"`
	LogisticsService				string		`json:"logistics_service"`
	MaintenanceServiceDateBegin     string 		`json:"maintenance_service_date_begin"`
	MaintenanceServiceDateEnd       string 		`json:"maintenance_service_date_end"`
	MaintenanceServiceStatus		string		`json:"maintenance_service_status"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *BatchUpdateDeviceLifecycles) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.DeviceLifecycles: "device_lifecycles",
	}
}

// Validate 对修改的数据做基本校验
func (reqData *BatchUpdateDeviceLifecycles) Validate(request *http.Request, errs binding.Errors) binding.Errors {
	for _, d := range reqData.DeviceLifecycles {
		if d.SN == "" {
			errs.Add([]string{"设备SN号"}, binding.RequiredError, fmt.Sprintf("设备SN不能为空"))
			return errs
		}
		// 校验是否租赁字段
		if d.IsRental != "" {
			if d.IsRental != model.YESRental &&  d.IsRental != model.NORental {
				errs.Add([]string{"is_rental"}, binding.RequiredError, fmt.Sprintf("是否租赁(%s)必须为 yes 或 no", d.IsRental))
				return errs
			}
		}
		if d.MaintenanceServiceDateBegin != "" && !strings.Contains(d.MaintenanceServiceDateBegin, "-") {
			errs.Add([]string{"maintenance_service_date_begin"}, binding.RequiredError, fmt.Sprintf("维保起始日期(%s)格式须为：YYYY-MM-DD", d.MaintenanceServiceDateBegin))
			return errs
		}
		if d.MaintenanceServiceDateEnd != "" && !strings.Contains(d.MaintenanceServiceDateEnd, "-") {
			errs.Add([]string{"maintenance_service_date_end"}, binding.RequiredError, fmt.Sprintf("维保截止日期(%s)格式须为：YYYY-MM-DD", d.MaintenanceServiceDateEnd))
			return errs
		}
	}

	return nil
}


//BatchUpdateDeviceLifecycleBySN 批量修改设备生命周期记录
func BatchUpdateDeviceLifecycleBySN(log logger.Logger, repo model.Repo, conf *config.Config, reqData *BatchUpdateDeviceLifecycles) (succeedSNs []string, totalAffected int64, err error) {
	for _, d := range reqData.DeviceLifecycles {
		// DeviceLifecycle 查询是否已经存在
		modDevLifecycle, err := repo.GetDeviceLifecycleBySN(d.SN)
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Errorf("GetDeviceLifecycleBySN failed: %s", err.Error())
			return succeedSNs, totalAffected, err
		}
		if modDevLifecycle != nil {
			var details []string
			if d.AssetBelongs != "" {
				details = append(details, fmt.Sprintf("资产归属:%v -> %v", modDevLifecycle.AssetBelongs, d.AssetBelongs))
				modDevLifecycle.AssetBelongs = d.AssetBelongs
			}
			if d.Owner != "" {
				details = append(details, fmt.Sprintf("负责人:%v -> %v", modDevLifecycle.Owner, d.Owner))
				modDevLifecycle.Owner = d.Owner
			}
			if d.IsRental != "" {
				details = append(details, fmt.Sprintf("是否租赁:%v -> %v", modDevLifecycle.IsRental, d.IsRental))
				modDevLifecycle.IsRental = d.IsRental
			}
			if d.MaintenanceServiceProvider != "" {
				details = append(details, fmt.Sprintf("维保服务供应商:%v -> %v", modDevLifecycle.MaintenanceServiceProvider, d.MaintenanceServiceProvider))
				modDevLifecycle.MaintenanceServiceProvider = d.MaintenanceServiceProvider
			}
			if d.MaintenanceService != "" {
				details = append(details, fmt.Sprintf("维保服务内容:%v -> %v", modDevLifecycle.MaintenanceService, d.MaintenanceService))
				modDevLifecycle.MaintenanceService = d.MaintenanceService
			}
			if d.LogisticsService != "" {
				details = append(details, fmt.Sprintf("物流服务内容:%v -> %v", modDevLifecycle.LogisticsService, d.LogisticsService))
				modDevLifecycle.LogisticsService = d.LogisticsService
			}
			if d.MaintenanceServiceDateBegin != "" {
				details = append(details, fmt.Sprintf("维保服务起始日期:%v -> %v", modDevLifecycle.MaintenanceServiceDateBegin, d.MaintenanceServiceDateBegin))
				modDevLifecycle.MaintenanceServiceDateBegin, _ = time.Parse(times.DateLayout, d.MaintenanceServiceDateBegin)
			}
			if d.MaintenanceServiceDateEnd != "" {
				details = append(details, fmt.Sprintf("维保服务截止日期:%v -> %v", modDevLifecycle.MaintenanceServiceDateEnd, d.MaintenanceServiceDateEnd))
				modDevLifecycle.MaintenanceServiceDateEnd, _ = time.Parse(times.DateLayout, d.MaintenanceServiceDateEnd)
			}
		
			now := time.Now()
			currentDate := now.Format("2006-01-02")
			if now.Before(modDevLifecycle.MaintenanceServiceDateBegin) {
				log.Debugf("current_date=%s < maintenance_service_begin_date=%s. Update Device(SN:%s) set status=%s -> %s", 
					currentDate, modDevLifecycle.MaintenanceServiceDateBegin, d.SN, modDevLifecycle.MaintenanceServiceStatus, model.MaintenanceServiceStatusInactive)
				details = append(details, fmt.Sprintf("维保服务状态:%v -> %v", modDevLifecycle.MaintenanceServiceStatus, model.MaintenanceServiceStatusInactive))
				modDevLifecycle.MaintenanceServiceStatus = model.MaintenanceServiceStatusInactive
			} else if now.After(modDevLifecycle.MaintenanceServiceDateEnd) {
				log.Debugf("current_date=%s > maintenance_service_end_date=%s. Update Device(%s) set status=%s -> %s", 
					currentDate, modDevLifecycle.MaintenanceServiceDateEnd, d.SN, modDevLifecycle.MaintenanceServiceStatus, model.MaintenanceServiceStatusOutOfWarranty)
				details = append(details, fmt.Sprintf("维保服务状态:%v -> %v", modDevLifecycle.MaintenanceServiceStatus, model.MaintenanceServiceStatusOutOfWarranty))
				modDevLifecycle.MaintenanceServiceStatus = model.MaintenanceServiceStatusOutOfWarranty
			} else {
				log.Debugf("current_date=%s between  maintenance_service_begin_date=%s and maintenance_service_end_date=%s. Update Device(%s) set status=%s -> %s",
					currentDate, modDevLifecycle.MaintenanceServiceDateBegin, modDevLifecycle.MaintenanceServiceDateEnd, d.SN, modDevLifecycle.MaintenanceServiceStatus, model.MaintenanceServiceStatusUnderWarranty)
				details = append(details, fmt.Sprintf("维保服务状态:%v -> %v", modDevLifecycle.MaintenanceServiceStatus, model.MaintenanceServiceStatusUnderWarranty))
				modDevLifecycle.MaintenanceServiceStatus = model.MaintenanceServiceStatusUnderWarranty
			}
			if len(details) > 0 {
				// 获取当前的生命周期日志记录
				var devLL []model.ChangeLog
				if modDevLifecycle.LifecycleLog != "" {
					if err = json.Unmarshal([]byte(modDevLifecycle.LifecycleLog), &devLL);err != nil {
						log.Error(err.Error())
						return succeedSNs, totalAffected, err
					}
				}			
				devLog := model.ChangeLog {
					OperationUser:		reqData.LoginName,
					OperationType:		model.OperationTypeUpdate,
					OperationDetail:	strings.Join(details, "；"),
					OperationTime:		times.ISOTime(time.Now()).ToTimeStr(),
				}
				// 追加Log
				devLL = append(devLL, devLog)
				b, err := json.Marshal(devLL)
				if err != nil {
					log.Error(err.Error())
					return succeedSNs, totalAffected, err
				}
				modDevLifecycle.LifecycleLog = string(b)
			}
			_, err := repo.SaveDeviceLifecycle(modDevLifecycle)
			if err != nil {
				log.Errorf("SaveDeviceLifecycle failed: %s", err.Error())
				return  succeedSNs, totalAffected, err
			}
			// 统计成功
			totalAffected++
			succeedSNs = append(succeedSNs, d.SN)

		} else {
			log.Errorf("查询设备SN(%s)生命周期信息不存在", d.SN)
			return succeedSNs, totalAffected, fmt.Errorf("查询设备SN(%s)生命周期信息不存在", d.SN)
		}
	}
	return
}