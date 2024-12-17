package service

import (
	"strconv"
	"idcos.io/cloudboot/utils/times"
	"encoding/json"
	"net/http"

	"fmt"

	"time"

	"os"

	"github.com/jinzhu/gorm"
	"github.com/voidint/binding"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils"
	"idcos.io/cloudboot/utils/upload"
)

//SpecialDevice 设备导入Excel表对应字段
//SN、型号、厂商、硬件备注、机房管理单元、机架、机位、是否分配内网IP、是否分配外网IP。
type SpecialDevice struct {
	//序列号
	SN string `json:"sn"`
	//厂商
	Vendor string `json:"vendor"`
	//型号
	DevModel string `json:"model"`
	// 设备类型 IN [堡垒机|加密机|前置机|SSL服务器|签名服务器|X86服务器|GPU服务器|授时服务器|其他]
	Category string `json:"category"`
	// 设备负责人
	Owner	string `json:"owner"`
	// 维保服务起始日期
	MaintenanceServiceDateBegin    string `json:"maintenance_service_date_begin"`
	// 保修期（月数）
	MaintenanceMonths    int `json:"maintenance_months"`
	// 操作系统名称
	OSReleaseName	string `json:"os_release_name"`
	//机房管理单元名称
	ServerRoomName string `json:"server_room_name"`
	//机架编号
	CabinetNum string `json:"server_cabinet_number"`
	//机位编号
	USiteNum string `json:"server_usite_number"`
	//硬件说明
	HardwareRemark string `json:"hardware_remark"`
	//是否分配外网IP yes|no
	NeedExtranetIP string `json:"need_extranet_ip"`
	//是否分配内网IP yes|no
	NeedIntranetIP string `json:"need_intranet_ip"`
	//以上是对应Excel导入字段，以下字段是通过名称关联到
	// 数据校验用
	ErrMsgContent string `json:"content"`
	//Usage string `json:"usage"` //固定为特殊设备
	ID            uint
	FixedAssetNum string //由规则生成
	IDCID         uint   //数据中心ID
	ServerRoomID  uint   `json:"server_room_id"`    //机房管理单元ID
	CabinetID     uint   `json:"server_cabinet_id"` //机架ID
	USiteID       uint   `json:"server_usite_id"`   //机位ID
	//关联订单号(非必填)
	OrderNumber string		`json:"order_number"`	
}

type SpecialDeviceReq struct {
	SpecialDevice
	LoginName string `json:"-"`
}


// WHEN Usage=特殊设备 -> 设备类型 IN [堡垒机|加密机|前置机|SSL服务器|签名服务器|X86服务器|GPU服务器|授时服务器|其他]
var SpecialDevCategoryList = []string{"堡垒机","加密机","前置机","SSL服务器","签名服务器","X86服务器","GPU服务器","授时服务器","其他"}

// FieldMap 请求字段映射
func (reqData *SpecialDeviceReq) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&reqData.SN:            					"sn",
		&reqData.Vendor:        					"vendor",
		&reqData.DevModel:      					"model",
		&reqData.Category:      					"category",
		&reqData.Owner:       						"owner",
		&reqData.MaintenanceServiceDateBegin:       "maintenance_service_date_begin",
		&reqData.MaintenanceMonths:       			"maintenance_months",
		&reqData.OSReleaseName:       				"os_release_name",
		&reqData.ServerRoomName: 					"server_room_name",
		&reqData.CabinetNum:     					"server_cabinet_number",
		&reqData.USiteNum:       					"server_usite_number",
		&reqData.HardwareRemark: 					"hardware_remark",
		&reqData.NeedIntranetIP: 					"need_intranet_ip",
		&reqData.NeedExtranetIP: 					"need_extranet_ip",
		&reqData.ServerRoomID:   					"server_room_id",
		&reqData.CabinetID:      					"server_cabinet_id",
		&reqData.USiteID:        					"server_usite_id",
		&reqData.OrderNumber:    					"order_number",
	}
}

// Validate 校验入参
func (reqData *SpecialDeviceReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())
	//机房校验
	//字段空值校验
	if reqData.SN == "" {
		errs.Add([]string{"SN"}, binding.RequiredError, "SN不能为空")
		return errs
	}
	if reqData.DevModel == "" {
		errs.Add([]string{"DevModel"}, binding.RequiredError, "型号不能为空")
		return errs
	}
	// WHEN Usage=特殊设备 -> 设备类型 IN [堡垒机|加密机|前置机|SSL服务器|签名服务器|X86服务器|GPU服务器|授时服务器|其他]	
	checkMsg := "特殊设备的设备类型 NOT IN [堡垒机|加密机|前置机|SSL服务器|签名服务器|X86服务器|GPU服务器|授时服务器|其他]"
	for _, devCategory := range SpecialDevCategoryList {
		if reqData.Category == devCategory {
			checkMsg = ""
		}
	}
	if checkMsg != "" {
		errs.Add([]string{"DevCategory"}, binding.RequiredError, checkMsg)
		return errs
	}
	if reqData.Vendor == "" {
		errs.Add([]string{"Vendor"}, binding.RequiredError, "厂商不能为空")
		return errs
	}
	if reqData.OSReleaseName == "" {
		errs.Add([]string{"OSReleaseName"}, binding.RequiredError, "操作系统名称不能为空")
		return errs
	}
	if reqData.MaintenanceMonths <= 0 {
		errs.Add([]string{"MaintenanceMonths"}, binding.RequiredError, "保修期（月数）不能小于0个")
		return errs
	}		
	if reqData.ServerRoomID == 0 {
		errs.Add([]string{"ServerRoomID"}, binding.RequiredError, "机房不能为空")
		return errs
	}
	if reqData.CabinetID == 0 {
		errs.Add([]string{"CabinetID"}, binding.RequiredError, "机架不能为空")
		return errs
	}
	if reqData.USiteID == 0 {
		errs.Add([]string{"USiteID"}, binding.RequiredError, "机位不能为空")
		return errs
	}
	d, err := repo.GetDeviceBySN(reqData.SN)
	if d != nil {
		errs.Add([]string{"SN"}, binding.RequiredError, fmt.Sprintf("SN:%s已存在，不能重复", reqData.SN))
		return errs
		//reqData.id = d.ID
	}
	if reqData.NeedExtranetIP != model.YES && reqData.NeedExtranetIP != model.NO {
		errs.Add([]string{"NeedExtranetIP"}, binding.RequiredError, "是否分配外网IP值不合法（yes|no）")
		return errs
	}
	if reqData.NeedIntranetIP != model.YES && reqData.NeedIntranetIP != model.NO {
		errs.Add([]string{"NeedIntranetIP"}, binding.RequiredError, "是否分配内网IP值不合法（yes|no）")
		return errs
	}
	//机房IDC
	srs, err := repo.GetServerRoomByID(reqData.ServerRoomID)
	if err == gorm.ErrRecordNotFound || srs == nil {
		errs.Add([]string{"ServerRoomID"}, binding.RequiredError, fmt.Sprintf("机房(ID:%d)不存在", reqData.ServerRoomID))
		return errs
	} else {
		reqData.IDCID = srs.IDCID
	}
	//机架
	cabinet, err := repo.GetServerCabinetByID(reqData.CabinetID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errs
	}
	if err == gorm.ErrRecordNotFound || cabinet == nil {
		errs.Add([]string{"cabinet"}, binding.RequiredError, fmt.Sprintf("机架(ID:%d)不存在", reqData.CabinetID))
		return errs
	} else {
		//reqData.CabinetID = cabinet.ID
	}
	//机位
	uSite, err := repo.GetServerUSiteByID(reqData.USiteID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errs
	}
	if err == gorm.ErrRecordNotFound || uSite == nil {
		errs.Add([]string{"USiteNum"}, binding.RequiredError, fmt.Sprintf("机位(ID:%d)不存在", reqData.USiteID))
		return errs
	} else {
		// 设备已存在，说明是先进入bootos或者数据已经导入过，重新导入刷新，此时的机位如果没有变化
		// 则不判断机位占用状态，因为就是被自己占着
		dev, _ := repo.GetDeviceBySN(reqData.SN)
		if !CheckUSiteFree(repo, uSite.ID, dev) {
			errs.Add([]string{"USiteNum"}, binding.RequiredError,
				fmt.Sprintf("机位编号(%s), 机架编号(%s)被占用或不可用", uSite.Number, cabinet.Number))
			return errs
		} else {
			reqData.USiteID = uSite.ID
		}
	}
	// 订单编号
	if reqData.OrderNumber != "" {
		order, err := repo.GetOrderByNumber(reqData.OrderNumber)
		if err != nil && err != gorm.ErrRecordNotFound {
			return errs
		}
		if err == gorm.ErrRecordNotFound || order == nil {
			errs.Add([]string{"OrderNumber"}, binding.RequiredError, fmt.Sprintf("订单(订单号:%s)不存在", reqData.OrderNumber))
			return errs
		}
	}
	return nil
}

// SaveSpecialDevices 保存（新增/修改）
func SaveSpecialDevices(log logger.Logger, repo model.Repo, dev *SpecialDeviceReq) (mod *model.Device, err error) {
	_ = log

	//统计一下关联订单号的数量，用于更新订单状态
	var mOrderAmount = make(map[string]int, 0)

	mod = &model.Device{
		Model:          gorm.Model{ID: uint(dev.ID)},
		SN:             dev.SN,
		Vendor:         dev.Vendor,
		DevModel:       dev.DevModel,
		Category:       dev.Category,
		IDCID:          dev.IDCID,
		ServerRoomID:   dev.ServerRoomID,
		CabinetID:      dev.CabinetID,
		USiteID:        &dev.USiteID,
		HardwareRemark: dev.HardwareRemark,
		Creator:        dev.LoginName,
		OrderNumber:	dev.OrderNumber,
	}
	fn, err := GenFixedAssetNumber(repo)
	if err != nil {
		return nil, err
	}
	mod.FixedAssetNumber = fn
	// WHEN Usage=特殊设备 -> 设备类型 IN [堡垒机|加密机|前置机|SSL服务器|签名服务器|X86服务器|GPU服务器|授时服务器|其他]
	mod.Usage = model.DevUsageSpecialDev
	// mod.Category = model.DevUsageSpecialDev
	mod.OperationStatus = model.DevOperStatOnShelve
	mod.PowerStatus = model.PowerStatusOn
	now := time.Now()
	mod.StartedAt = now
	mod.OnShelveAt = now
	mod.CreatedAt = now

	if mod.USiteID != nil {
		if _, err = repo.BatchUpdateServerUSitesStatus([]uint{*mod.USiteID}, model.USiteStatUsed); err != nil {
			err = fmt.Errorf("占用机位:%s失败,%v", dev.USiteNum, err)
			return
		}
	}

	_, err = repo.SaveDevice(mod)
	if err != nil {
		return
	}

	//分配IP
	var inIP, exIP *model.IP
	if dev.NeedIntranetIP == model.YES {
		if inIP, err = repo.AssignIntranetIP(dev.SN); err != nil {
			err = fmt.Errorf("SN:%s分配内网失败，err：%v", dev.SN, err)
			return
		}
	} else {
		//需要尝试释放IP
		if _, err = repo.ReleaseIP(dev.SN, model.Intranet); err != nil {
			err = fmt.Errorf("SN:%s释放内网IP失败，err：%v", dev.SN, err)
			return
		}
	}

	if dev.NeedExtranetIP == model.YES {
		if exIP, err = repo.AssignExtranetIP(dev.SN); err != nil {
			err = fmt.Errorf("SN:%s分配外网失败，err：%v", dev.SN, err)
			return
		}
	} else {
		if _, err = repo.ReleaseIP(dev.SN, model.Extranet); err != nil {
			err = fmt.Errorf("SN:%s释放外网IP失败，err：%v", dev.SN, err)
			return
		}
	}
	// 特殊设备无需经过系统部署，自定义写入操作系统名称
	sysTpl, err := repo.GetSystemTemplateByName(dev.OSReleaseName)
	if err != nil {
		log.Errorf("get system template by name %s failed(%v) ,adding one", dev.OSReleaseName, err)
		tpl := model.SystemTemplate{
			Family:   	"Custom",
			BootMode: 	"uefi",
			Name:     	dev.OSReleaseName,
			PXE:        "#NULL",
			Content:    "#NULL",
			OSLifecycle: model.OSTesting,
			Arch:		 model.OSARCHUNKNOWN,

		}
		_, err := repo.SaveSystemTemplate(&tpl)
		if err != nil {
			log.Errorf("add system template by name %s fail,%v", dev.OSReleaseName, err)
		} else {
			//重新获取新增模板ID
			sysTpl, err = repo.GetSystemTemplateByName(dev.OSReleaseName)
		}
	}

	//Add device setting(模拟一条数据)
	ds, err := repo.GetDeviceSettingBySN(dev.SN)
	if ds == nil || err != nil {
		ds = &model.DeviceSetting{
			SN:              dev.SN,
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
	ds.NeedExtranetIP = dev.NeedExtranetIP
	if inIP != nil {
		ds.IntranetIP = inIP.IP
	}
	if exIP != nil {
		ds.ExtranetIP = exIP.IP
	}
	if err = repo.SaveDeviceSetting(ds); err != nil {
		return
	}
	// 仅记录必要字段到“设备新增”
	optDetail, err := convert2DetailOfOperationTypeAdd(repo, *mod)
	if err != nil {
		log.Errorf("Fail to convert Detail of OperationTypeAdd: %v", err)
	}
	// DeviceLifecycle 变更记录
	deviceLifecycleLog := []model.ChangeLog {
		{
			OperationUser:		dev.LoginName,
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
			Owner:							dev.Owner,
			IsRental:						"no",
			MaintenanceServiceProvider:		dev.Vendor,
			MaintenanceService:				"Undefined",
			LogisticsService:				"Undefined",
			MaintenanceServiceStatus:		model.MaintenanceServiceStatusInactive, //新增场景默认-未激活
			LifecycleLog:					string(b),
		},
	}

	if dev.MaintenanceServiceDateBegin != "" {
		t, err := time.Parse(times.DateLayout, dev.MaintenanceServiceDateBegin)
		if err != nil {
			log.Errorf("parse maintenance time %s err:%v , using current time for maintenance-date", dev.MaintenanceServiceDateBegin, err)
			saveDevLifecycleReq.MaintenanceServiceDateBegin = now
			saveDevLifecycleReq.MaintenanceServiceDateEnd = now.AddDate(0, dev.MaintenanceMonths, 0)
		} else {
			saveDevLifecycleReq.MaintenanceServiceDateBegin = t
			saveDevLifecycleReq.MaintenanceServiceDateEnd = t.AddDate(0, dev.MaintenanceMonths, 0)
		}
	} else {
		saveDevLifecycleReq.MaintenanceServiceDateBegin = now
		saveDevLifecycleReq.MaintenanceServiceDateEnd = now.AddDate(0, dev.MaintenanceMonths, 0)
	}

	// 通过订单编号获取资产归属、负责人、维保服务等内容
	// 若无订单编号则以参数输入为准
	if dev.OrderNumber != "" {
		order, err := repo.GetOrderByNumber(dev.OrderNumber)
		if err != nil {
			log.Errorf("订单(订单号:%s)不存在", dev.OrderNumber)
			return nil, err
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
			mOrderAmount[dev.OrderNumber]++
		}
	}
	// DeviceLifecycle 查询是否已经存在
	devLifecycle, err := repo.GetDeviceLifecycleBySN(dev.SN)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} 
	if devLifecycle != nil {
		log.Debugf("DeviceLifecycle SN %s already exist.Update it.", dev.SN)
		saveDevLifecycleReq.ID = devLifecycle.ID
	}
	//更新关联的订单到货数量和订单状态
	for orderNum, arrivalCount := range mOrderAmount {
		if err = UpdateOrderByArrival(log, repo, orderNum, arrivalCount); err != nil {
			return nil, err
		}
	}
	// 保存或更新 DeviceLifecycle
	if err = SaveDeviceLifecycle(log, repo, saveDevLifecycleReq); err != nil {
		log.Debug(err)
		return nil, err
	}	
	//如果没有成功，要回滚!!
	defer func() {
		if err != nil {
			//rollback
			_, _ = repo.BatchUpdateServerUSitesStatus([]uint{*mod.USiteID}, model.USiteStatFree)
			_, _ = repo.DeleteDeviceSettingBySN(mod.SN)
			_, _ = repo.RemoveDeviceByID(mod.ID)
			_, _ = repo.RemoveDeviceLifecycleByID(devLifecycle.ID)
			mod = nil
			saveDevLifecycleReq = nil
		}
	}()

	return mod, nil
}

func ImportSpecialDevicesPreview(log logger.Logger, repo model.Repo, reqData *ImportPreviewReq) (map[string]interface{}, error) {
	ra, err := utils.ParseDataFromXLSX(upload.UploadDir + reqData.FileName)
	if err != nil {
		return nil, err
	}
	length := len(ra)

	var success []*SpecialDeviceReq
	var failure []*SpecialDeviceReq

	//if valid, err := CheckUnique(ra); !valid {
	//	return nil, err
	//}
	for i := 1; i < length; i++ {
		row := &SpecialDeviceReq{}
		if len(ra[i]) < 15 {
			var br string
			if row.ErrMsgContent != "" {
				br = "<br />"
			}
			row.ErrMsgContent += br + "导入文件列长度不对（应为15列）"
			failure = append(failure, row)
			continue
		}
		row.SN = ra[i][0]
		row.Vendor = ra[i][1]
		row.DevModel = ra[i][2]
		row.ServerRoomName = ra[i][3]
		row.CabinetNum = ra[i][4]
		row.USiteNum = ra[i][5]
		row.HardwareRemark = ra[i][6]
		row.NeedIntranetIP = ra[i][7]
		row.NeedExtranetIP = ra[i][8]
		row.OrderNumber = ra[i][9]
		row.Category = ra[i][10]
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
		row.OSReleaseName = ra[i][14]


		utils.StructTrimSpace(&row.SpecialDevice)

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
	}

	var data []*SpecialDeviceReq
	if len(failure) > 0 {
		data = failure
	} else {
		data = success
	}
	var result []*SpecialDeviceReq
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

//checkLength 对导入文件中的数据做字段长度校验
func (impDevReq *SpecialDeviceReq) checkLength() {
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
	leg = len(impDevReq.DevModel)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:型号长度为(%d)(不能为空，不能大于255)", leg)
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
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:机架长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.USiteNum)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:机位长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.NeedIntranetIP)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:是否分配内网IP长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.NeedExtranetIP)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:是否分配外网IP长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.Category)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:设备类型长度为(%d)(不能为空，不能大于255)", leg)
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
	leg = len(impDevReq.OSReleaseName)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:操作系统名称长度为(%d)(不能为空，不能大于255)", leg)
	}	
}

//validate 对导入文件中的数据做基本验证
func (impDevReq *SpecialDeviceReq) validate(repo model.Repo) error {
	d, err := repo.GetDeviceBySN(impDevReq.SN)
	if d != nil {
		impDevReq.ID = d.ID
	}
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
		impDevReq.IDCID = srs.IDCID
		impDevReq.ServerRoomID = srs.ID
	}

	//机架
	cabinet, err := repo.GetServerCabinetByNumber(impDevReq.ServerRoomID, impDevReq.CabinetNum)
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
		impDevReq.CabinetID = cabinet.ID
	}
	//机位
	uSite, err := repo.GetServerUSiteByNumber(impDevReq.CabinetID, impDevReq.USiteNum)
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
			impDevReq.USiteID = uSite.ID
		}
	}
	if impDevReq.NeedExtranetIP != model.YesCh && impDevReq.NeedExtranetIP != model.NoCh {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("是否分配外网IP值%s不合法（yes|no）", impDevReq.NeedExtranetIP)
	}
	if impDevReq.NeedIntranetIP != model.YesCh && impDevReq.NeedIntranetIP != model.NoCh {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("是否分配内网IP值%s不合法（yes|no）", impDevReq.NeedIntranetIP)
	}
	// 订单编号校验
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
			if order.IDCID != impDevReq.IDCID {
				var br string
				if impDevReq.ErrMsgContent != "" {
					br = "<br />"
				}
				impDevReq.ErrMsgContent += br + fmt.Sprintf("订单(订单号:%s)与设备（SN:%s）不属于同一个数据中心", impDevReq.OrderNumber, impDevReq.SN)
			}
		}
	}

	return nil
}

func ImportSpecialDevices(log logger.Logger, repo model.Repo, reqData *ImportPreviewReq) error {
	fileName := upload.UploadDir + reqData.FileName
	ra, err := utils.ParseDataFromXLSX(fileName)
	if err != nil {
		return err
	}
	length := len(ra)

	_ = os.Remove(fileName)

	for i := 1; i < length; i++ {
		row := &SpecialDeviceReq{}
		if len(ra[i]) < 15 {
			var br string
			if row.ErrMsgContent != "" {
				br = "<br />"
			}
			row.ErrMsgContent += br + "导入文件列长度不对（应为15列）"
			continue
		}
		row.SN = ra[i][0]
		row.Vendor = ra[i][1]
		row.DevModel = ra[i][2]
		row.ServerRoomName = ra[i][3]
		row.CabinetNum = ra[i][4]
		row.USiteNum = ra[i][5]
		row.HardwareRemark = ra[i][6]
		row.OrderNumber = ra[i][9]
		row.Category = ra[i][10]
		row.Owner = ra[i][11]
		row.OSReleaseName = ra[i][14]
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
		row.OSReleaseName = ra[i][14]

		if ra[i][7] == model.YesCh {
			row.NeedIntranetIP = model.YES
		} else {
			row.NeedIntranetIP = model.NO
		}
		if ra[i][8] == model.YesCh {
			row.NeedExtranetIP = model.YES
		} else {
			row.NeedExtranetIP = model.NO
		}
		row.LoginName = reqData.LoginName
		utils.StructTrimSpace(&row.SpecialDevice)
		row.validate(repo)

		if _, err := SaveSpecialDevices(log, repo, row); err != nil {
			log.Errorf("import special device(SN:%s) fail,%v", row.SN, err)
			return err
		}

	}
	return nil
}

