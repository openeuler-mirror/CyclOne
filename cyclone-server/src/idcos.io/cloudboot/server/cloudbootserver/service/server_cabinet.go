package service

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/voidint/binding"
	"github.com/voidint/page"

	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils"
	"idcos.io/cloudboot/utils/collection"
	strings2 "idcos.io/cloudboot/utils/strings"
	"idcos.io/cloudboot/utils/times"
	"idcos.io/cloudboot/utils/upload"
)

// SaveServerCabinetReq 保存机架(柜)请求结构体
type SaveServerCabinetReq struct {
	// 所属数据中心ID
	IDCID uint `json:"idc_id"`
	// 所属机房ID
	ServerRoomID uint `json:"server_room_id"`
	// 网络区域ID
	NetworkAreaID uint `json:"network_area_id"`
	// 机架(柜)ID。若id=0，则新增。若id>0，则修改。
	ID uint `json:"id"`
	// 机架(柜)编号
	Number string `json:"number"`
	// 机架(柜)高度
	Height uint `json:"height"`
	// 类型("server","network_device", "reserved")
	Type string `json:"type"`
	// 网络速率
	NetworkRate string `json:"network_rate"`
	// 电流
	Current string `json:"current"`
	// 可用功率
	AvailablePower string `json:"available_power"`
	// 峰值功率
	MaxPower string `json:"max_power"`
	// 备注
	Remark string `json:"remark"`
	// 用户登录名
	LoginName string `json:"-"`
}

// CabinetPowerBatchOperateReq 机架批量操作
type CabinetPowerBatchOperateReq struct {
	IDS []uint `json:"ids"`
}

// FieldMap 请求字段映射
func (reqData *CabinetPowerBatchOperateReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDS: "ids",
	}
}

// FieldMap 请求字段映射
func (reqData *SaveServerCabinetReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDCID:          "idc_id",
		&reqData.ServerRoomID:   "server_room_id",
		&reqData.NetworkAreaID:  "network_area_id",
		&reqData.ID:             "id",
		&reqData.Number:         "number",
		&reqData.Height:         "height",
		&reqData.Type:           "type",
		&reqData.NetworkRate:    "network_rate",
		&reqData.Current:        "current",
		&reqData.AvailablePower: "available_power",
		&reqData.MaxPower:       "max_power",
		&reqData.Remark:         "remark",
	}
}

// Validate 结构体数据校验
func (reqData *SaveServerCabinetReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())

	// if reqData.ServerRoomID == 0 {
	// 	errs.Add([]string{"server_room_id"}, binding.RequiredError, "机房管理单元ID不能为空")
	// 	return errs
	// }
	if reqData.NetworkAreaID == 0 {
		errs.Add([]string{"network_area_id"}, binding.RequiredError, "网络区域ID不能为空")
		return errs
	}
	if reqData.Number == "" {
		errs.Add([]string{"number"}, binding.RequiredError, "机架编号不能为空")
		return errs
	}
	if reqData.Height == 0 {
		errs.Add([]string{"height"}, binding.RequiredError, "机架高度须大于0")
		return errs
	}
	if reqData.Type == "" {
		errs.Add([]string{"type"}, binding.RequiredError, "机架类型不能为空")
		return errs
	}
	if !collection.InSlice(reqData.Type, []string{model.CabinetTypeServer, model.CabinetTypeKvmServer, model.CabinetTypeNetworkDevice, model.CabinetTypeReserved}) {
		errs.Add([]string{"type"}, binding.RequiredError, "无效的机架类型")
		return errs
	}
	if reqData.NetworkRate == "" {
		errs.Add([]string{"network_rate"}, binding.RequiredError, "网络速率不能为空")
		return errs
	}
	if reqData.Current == "" {
		errs.Add([]string{"current"}, binding.RequiredError, "电流不能为空")
		return errs
	}
	if reqData.AvailablePower == "" {
		errs.Add([]string{"available_power"}, binding.RequiredError, "可用功率不能为空")
		return errs
	}
	if reqData.MaxPower == "" {
		errs.Add([]string{"max_power"}, binding.RequiredError, "峰值功率不能为空")
		return errs
	}
	// 校验指定ID的网络区域是否存在
	area, err := repo.GetNetworkAreaByID(reqData.NetworkAreaID)
	if err == gorm.ErrRecordNotFound {
		errs.Add([]string{"idc_id", "server_room_id", "network_area_id"}, binding.BusinessError, "该网络区域不存在")
		return errs
	}
	if err != nil {
		errs.Add([]string{"idc_id", "server_room_id", "network_area_id"}, binding.SystemError, fmt.Sprintf("系统内部错误:%s", err.Error()))
		return errs
	}
	// if area.IDCID != reqData.IDCID || area.ServerRoomID != reqData.ServerRoomID {
	// 	errs.Add([]string{"idc_id", "server_room_id"}, binding.BusinessError, "数据中心、机房管理单元、网络区域三者不匹配")
	// 	return errs
	// }

	// sr, err := repo.GetServerRoomByID(reqData.ServerRoomID)
	// if err != nil {
	// 	errs.Add([]string{"server_room_id"}, binding.SystemError, "系统内部错误")
	// 	return errs
	// }
	reqData.IDCID = area.IDCID
	reqData.ServerRoomID = area.ServerRoomID
	// 校验机架编号在机房内的唯一性
	cond := model.ServerCabinetCond{
		IDCID:        []uint{reqData.IDCID},
		ServerRoomID: []uint{reqData.ServerRoomID},
		Number:       reqData.Number, // 模糊匹配
	}

	items, err := repo.GetServerCabinets(&cond, nil, nil)
	if err != nil {
		errs.Add([]string{"idc_id", "server_room_id", "number"}, binding.SystemError, "系统内部错误")
		return errs
	}
	for _, item := range items {
		if (reqData.ID == 0 && item.Number == reqData.Number) || // 新增时，机架编号不能重复。
			(reqData.ID > 0 && item.Number == reqData.Number && reqData.ID != item.ID) { // 更新时，机架编号不能重复（除了自身外）。
			errs.Add([]string{"number"}, binding.BusinessError, "同一个机房管理单元内机架编号不能重复")
			return errs
		}
	}
	return errs
}

// SaveServerCabinet 保存机架(柜)
func SaveServerCabinet(log logger.Logger, repo model.Repo, reqData *SaveServerCabinetReq) (err error) {
	if reqData.ID > 0 {
		_, err = repo.GetServerCabinetByID(reqData.ID)
		if err != nil {
			return err
		}
	}

	sc := &model.ServerCabinet{
		IDCID:          reqData.IDCID,
		ServerRoomID:   reqData.ServerRoomID,
		NetworkAreaID:  reqData.NetworkAreaID,
		Number:         reqData.Number,
		Height:         reqData.Height,
		Type:           reqData.Type,
		NetworkRate:    reqData.NetworkRate,
		Current:        reqData.Current,
		AvailablePower: reqData.AvailablePower,
		MaxPower:       reqData.MaxPower,
		Remark:         reqData.Remark,
		Creator:        reqData.LoginName,
	}
	sc.ID = reqData.ID

	if sc.ID == 0 {
		sc.IsEnabled = model.NO
		sc.IsPowered = model.CabinetPowerOff
		sc.Status = model.CabinetStatUnderConstruction
		_, err := repo.SaveServerCabinet(sc)
		if err != nil {
			return err
		}
	}
	if sc.ID > 0 {
		_, err := repo.SaveServerCabinet(sc)
		if err != nil {
			return err
		}
	}
	if reqData.ID == 0 {
		reqData.ID = sc.ID
	}
	return err
}

// ServerCabinet 机架(柜)
type ServerCabinet struct {
	//数据中心
	IDC struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"idc"`
	//机房
	ServerRoom struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"server_room"`
	//网络区域
	NetworkArea struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"network_area"`
	//机架(柜)ID
	ID uint `json:"id"`
	//机架(柜)编号
	Number string `json:"number"`
	//机架(柜)高度
	Height uint `json:"height"`
	//机架(柜)类型(("server","network_device", "reserved"))
	Type string `json:"type"`
	//网络速率
	NetworkRate string `json:"network_rate"`
	//电流
	Current string `json:"current"`
	//可用功率
	AvailablePower string `json:"available_power"`
	//最大功率
	MaxPower string `json:"max_power"`
	//是否启用
	IsEnabled string `json:"is_enabled"`
	//启用时间
	EnableTime string `json:"enable_time"`
	//是否上电
	IsPowered string `json:"is_powered"`
	//上电时间
	PowerOnTime string `json:"power_on_time"`
	//关电时间
	PowerOffTime string `json:"power_off_time"`
	//状态
	Status string `json:"status"`
	//创建时间
	CreatedAt string `json:"created_at"`
	//更新时间
	UpdatedAt string `json:"updated_at"`
	//机位总数
	USiteCount int64 `json:"usite_count"`
	//备注
	Remark string `json:"remark"`
}

// GetServerCabinetPageReq 查询机架(柜)分页请求结构体
type GetServerCabinetPageReq struct {
	// 所属数据中心ID,多个以逗号分隔
	IDCID string `json:"idc_id"`
	// 所属机房ID
	ServerRoomID string `json:"server_room_id"`
	// 机架ID
	ServerCabinetID string `json:"server_cabinet_id"`	
	// 网络区域ID
	NetworkAreaID string `json:"network_area_id"`
	// 机架(柜)编号
	Number string `json:"number"`
	// 类型
	Type string `json:"type"`
	// 状态。可选值: under_construction-建设中; not_enabled-未启用; enabled-已启用; offline-已下线;
	Status string `json:"status"`
	// 是否启用(yes/no)
	IsEnabled string `json:"is_enabled"`
	// 是否开电(yes/no)
	IsPowered string `json:"is_powered"`
	// 分页页号
	Page int64 `json:"page"`
	// 分页大小
	PageSize        int64 `json:"page_size"`
	ServerRoomName  string
	NetworkAreaName string
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *GetServerCabinetPageReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDCID:           "idc_id",
		&reqData.ServerRoomID:    "server_room_id",
		&reqData.ServerCabinetID:    "server_cabinet_id",
		&reqData.NetworkAreaID:   "network_area_id",
		&reqData.Number:          "number",
		&reqData.Type:            "type",
		&reqData.Status:          "status",
		&reqData.IsEnabled:       "is_enabled",
		&reqData.IsPowered:       "is_powered",
		&reqData.Page:            "page",
		&reqData.PageSize:        "page_size",
		&reqData.ServerRoomName:  "server_room_name",
		&reqData.NetworkAreaName: "network_area_name",
	}
}

// GetServerCabinetPage 按条件查询机架(柜)分页列表
func GetServerCabinetPage(log logger.Logger, repo model.Repo, reqData *GetServerCabinetPageReq) (pg *page.Page, err error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 1000 {
		reqData.PageSize = 10
	}
	if reqData.Page < 0 {
		reqData.Page = 0
	}

	cond := model.ServerCabinetCond{
		IDCID:           strings2.Multi2UintSlice(reqData.IDCID),
		ServerRoomID:    strings2.Multi2UintSlice(reqData.ServerRoomID),
		ServerCabinetID:    strings2.Multi2UintSlice(reqData.ServerCabinetID),
		NetworkAreaID:   strings2.Multi2UintSlice(reqData.NetworkAreaID),
		Number:          reqData.Number,
		Type:            reqData.Type,
		IsEnabled:       reqData.IsEnabled,
		IsPowered:       reqData.IsPowered,
		Status:          reqData.Status,
		ServerRoomName:  reqData.ServerRoomName,
		NetworkAreaName: reqData.NetworkAreaName,
	}
	////物理区域需要支持多个多选
	//if reqData.NetworkAreaID != "" {
	//	NetworkAreaIDs := mystrings.MultiLines2Slice(reqData.NetworkAreaID)
	//	for _, id := range NetworkAreaIDs {
	//		if idInt, err := strconv.Atoi(id); err != nil {
	//			return nil, fmt.Errorf("network id invalid: %v", err)
	//		} else {
	//			cond.NetworkAreaID = append(cond.NetworkAreaID, uint(idInt))
	//		}
	//	}
	//}
	totalRecords, err := repo.CountServerCabinets(&cond)
	if err != nil {
		return nil, err
	}

	pager := page.NewPager(reflect.TypeOf(&ServerCabinet{}), reqData.Page, reqData.PageSize, totalRecords)
	items, err := repo.GetServerCabinets(&cond, model.OneOrderBy("id", model.DESC), pager.BuildLimiter())
	if err != nil {
		return nil, err
	}

	for i := range items {
		pager.AddRecords(convert2ServerCabinet(repo, items[i]))
	}
	return pager.BuildPage(), nil
}

func convert2ServerCabinet(repo model.Repo, src *model.ServerCabinet) *ServerCabinet {
	dst := ServerCabinet{
		ID:             src.ID,
		Number:         src.Number,
		Height:         src.Height,
		Type:           src.Type,
		NetworkRate:    src.NetworkRate,
		Current:        src.Current,
		AvailablePower: src.AvailablePower,
		MaxPower:       src.MaxPower,
		IsEnabled:      src.IsEnabled,
		IsPowered:      src.IsPowered,
		Status:         src.Status,
		Remark:         src.Remark,
		CreatedAt:      src.CreatedAt.Format(times.DateTimeLayout),
		UpdatedAt:      src.UpdatedAt.Format(times.DateTimeLayout),
	}
	if src.EnableTime != nil {
		dst.EnableTime = src.EnableTime.Format(times.DateLayout)
	}
	if src.PowerOnTime != nil {
		dst.PowerOnTime = src.PowerOnTime.Format(times.DateLayout)
	}
	if src.PowerOffTime != nil {
		dst.PowerOffTime = src.PowerOffTime.Format(times.DateLayout)
	}

	if idc, _ := repo.GetIDCByID(src.IDCID); idc != nil {
		dst.IDC.ID, dst.IDC.Name = idc.ID, idc.Name
	}

	if room, _ := repo.GetServerRoomByID(src.ServerRoomID); room != nil {
		dst.ServerRoom.ID, dst.ServerRoom.Name = room.ID, room.Name
	}

	if area, _ := repo.GetNetworkAreaByID(src.NetworkAreaID); area != nil {
		dst.NetworkArea.ID, dst.NetworkArea.Name = area.ID, area.Name
	}

	//统计机房中的机架(柜数)
	if count, _ := repo.GetServerUSiteCountByServerCabinetID(dst.ID); count > 0 {
		dst.USiteCount = count
	}
	return &dst
}

// PowerOnServerCabinetByID 根据ID将机架(柜)上电
func PowerOnServerCabinetByID(repo model.Repo, ids []uint) (int64, error) {
	one, err := repo.PowerOnServerCabinetByID(ids)
	if err != nil {
		return 0, err
	}
	return one, nil
}

// PowerOffServerCabinetByID 根据ID将机架(柜)下电
func PowerOffServerCabinetByID(repo model.Repo, id uint) (int64, error) {
	one, err := repo.PowerOffServerCabinetByID(id)
	if err != nil {
		return 0, err
	}
	return one, nil
}

//ServerCabinetForImport 为了导入
type ServerCabinetForImport struct {
	//机房名
	ServerRoom string `json:"server_room"`
	//机架编号
	Number string `json:"number"`
	//机架高度
	Height int `json:"height"`
	//机架类型(通用服务器/网络设备/预留)
	Type string `json:"type"`
	//网络速率
	NetworkRate string `json:"network_rate"`
	//网络区域
	NetworkArea string `json:"network_area"`
	//电流
	Current string `json:"current"`
	//可用功率
	AvailablePower string `json:"available_power"`
	//最大功率
	MaxPower string `json:"max_power"`
	//备注
	Remark        string `json:"remark"`
	Content       string `json:"content"`
	Creator       string `json:"creator"`
	IDCID         uint   `json:"idc_id"`
	ServerRoomID  uint   `json:"server_room_id"`
	NetworkAreaID uint   `json:"network_area_id"`
	ID            uint   `json:"id"`
}

var (
	cabinetTypeMap = map[string]string{
		"通用服务器": model.CabinetTypeServer,
		"虚拟化服务器": model.CabinetTypeKvmServer,
		"网络设备":  model.CabinetTypeNetworkDevice,
		"预留":    model.CabinetTypeReserved,
		"":      "",
	}
)

//validate 对导入文件中的数据做基本验证
func (scfi *ServerCabinetForImport) validate(log logger.Logger, repo model.Repo) (int, error) {
	//机架高度
	if scfi.Height < 0 {
		var br string
		if scfi.Content != "" {
			br = "<br />"
		}
		scfi.Content += br + fmt.Sprintf("机架位高度(%d)不能小于0", scfi.Height)
	}

	//机架类型校验
	if cabinetTypeMap[scfi.Type] == "" {
		var br string
		if scfi.Content != "" {
			br = "<br />"
		}
		scfi.Content += br + fmt.Sprintf("机架类型只能为(%s)", "通用服务器/虚拟化服务器/网络设备/预留")
	}

	scfi.Type = cabinetTypeMap[scfi.Type]

	//机房校验
	srs, err := repo.GetServerRoomByName(scfi.ServerRoom)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Debug(err.Error())
		return upload.Return, err
	}
	if err == gorm.ErrRecordNotFound || srs == nil {
		var br string
		if scfi.Content != "" {
			br = "<br />"
		}
		scfi.Content += br + fmt.Sprintf("机房名(%s)不存在", scfi.ServerRoom)
		log.Debug(err.Error())
		return upload.Continue, nil
	}
	scfi.IDCID = srs.IDCID
	scfi.ServerRoomID = srs.ID

	//查询是否已经存在
	id, err := repo.GetServerCabinetID(&model.ServerCabinet{Number: scfi.Number,
		ServerRoomID: scfi.ServerRoomID,
		IDCID:        scfi.IDCID,
	})
	if err != nil {
		log.Debug(err.Error())
		return upload.Return, err
	}
	if len(id) > 1 {
		var br string
		if scfi.Content != "" {
			br = "<br />"
		}
		scfi.Content += br + fmt.Sprintf("机架(柜)(%s)存在冲突", scfi.Number)
		return upload.Continue, nil
	} else if len(id) == 1 {
		scfi.ID = id[0]
	} else {
		scfi.ID = 0
	}

	//网络区域校验 fixbug: 这样做有风险，等网络区域接口
	nabc, err := repo.GetNetworkAreasByCond(&model.NetworkArea{Name: scfi.NetworkArea,
		ServerRoomID: scfi.ServerRoomID,
		IDCID:        scfi.IDCID,
	})
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Debug(err.Error())
		return upload.Return, err
	}
	if len(nabc) < 1 {
		var br string
		if scfi.Content != "" {
			br = "<br />"
		}
		scfi.Content += br + fmt.Sprintf("网络区域(%s)不存在", scfi.NetworkArea)
		return upload.Continue, nil
	}
	if len(nabc) == 1 {
		scfi.NetworkAreaID = nabc[0].ID
	}

	return upload.DO, nil
}

//checkLength 对导入文件中的数据做基本验证
func (scfi *ServerCabinetForImport) checkLength() {
	leg := len(scfi.ServerRoom)
	if leg == 0 || leg > 255 {
		var br string
		if scfi.Content != "" {
			br = "<br />"
		}
		scfi.Content += br + fmt.Sprintf("必填项校验:机房名长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(scfi.Number)
	if leg == 0 || leg > 255 {
		var br string
		if scfi.Content != "" {
			br = "<br />"
		}
		scfi.Content += br + fmt.Sprintf("必填项校验:机架编号长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(scfi.Type)
	if leg == 0 || leg > 255 {
		var br string
		if scfi.Content != "" {
			br = "<br />"
		}
		scfi.Content += br + fmt.Sprintf("必填项校验:机架类型长度为(%d)(不能为空)", leg)
	}
	leg = len(scfi.NetworkRate)
	if leg == 0 || leg > 255 {
		var br string
		if scfi.Content != "" {
			br = "<br />"
		}
		scfi.Content += br + fmt.Sprintf("必填项校验:网络速率长度为(%d)(不能为空)", leg)
	}
	leg = len(scfi.Current)
	if leg == 0 || leg > 255 {
		var br string
		if scfi.Content != "" {
			br = "<br />"
		}
		scfi.Content += br + fmt.Sprintf("必填项校验:电流长度为(%d)(不能为空)", leg)
	}
	leg = len(scfi.AvailablePower)
	if leg == 0 || leg > 255 {
		var br string
		if scfi.Content != "" {
			br = "<br />"
		}
		scfi.Content += br + fmt.Sprintf("必填项校验:可用功率长度为(%d)(不能为空)", leg)
	}
	leg = len(scfi.MaxPower)
	if leg == 0 || leg > 255 {
		var br string
		if scfi.Content != "" {
			br = "<br />"
		}
		scfi.Content += br + fmt.Sprintf("必填项校验:峰值功率长度为(%d)(不能为空)", leg)
	}
}

//ImportServerCabinetPriview 导入预览
func ImportServerCabinetPriview(log logger.Logger, repo model.Repo, reqData *upload.ImportReq) (map[string]interface{}, error) {
	ra, err := utils.ParseDataFromXLSX(upload.UploadDir + reqData.FileName)
	if err != nil {
		return nil, err
	}
	length := len(ra)

	var success []*ServerCabinetForImport
	var failure []*ServerCabinetForImport
	for i := 1; i < length; i++ {
		row := &ServerCabinetForImport{}
		if len(ra[i]) < 10 {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "导入文件列长度不对（应为10列）"
			failure = append(failure, row)
			continue
		}

		row.ServerRoom = strings.TrimSpace(ra[i][0])
		row.Number = strings.TrimSpace(ra[i][1])
		row.Height, _ = strconv.Atoi(strings.TrimSpace(ra[i][2]))
		row.Type = strings.TrimSpace(ra[i][3])
		row.NetworkArea = strings.TrimSpace(ra[i][4])
		row.NetworkRate = strings.TrimSpace(ra[i][5])
		row.Current = strings.TrimSpace(ra[i][6])
		row.AvailablePower = strings.TrimSpace(ra[i][7])
		row.MaxPower = strings.TrimSpace(ra[i][8])
		row.Remark = strings.TrimSpace(ra[i][9])

		//必填项校验
		row.checkLength()
		//机房和网络区域校验
		_, err := row.validate(log, repo)
		if err != nil {
			log.Debug(err.Error())
			return nil, err
		}

		if row.Content != "" {
			failure = append(failure, row)
		} else {
			success = append(success, row)
		}
	}

	var data []*ServerCabinetForImport
	if len(failure) > 0 {
		data = failure
	} else {
		data = success
	}
	var result []*ServerCabinetForImport
	for i := 0; i < len(data); i++ {
		if uint(i) >= reqData.Offset && uint(i) < (reqData.Offset+reqData.Limit) {
			result = append(result, data[i])
		}
	}
	if len(failure) > 0 {
		return map[string]interface{}{"status": "failure",
			"message":       "导入机架错误",
			"import_status": false,
			"record_count":  len(data),
			"content":       result,
		}, nil
	}
	return map[string]interface{}{"status": "success",
		"message":       "操作成功",
		"import_status": true,
		"record_count":  len(data),
		"content":       result,
	}, nil
}

//ImportServerCabinet 将导入机架(柜)放到数据库
func ImportServerCabinet(log logger.Logger, repo model.Repo, reqData *upload.ImportReq) error {
	ra, err := utils.ParseDataFromXLSX(upload.UploadDir + reqData.FileName)
	if err != nil {
		log.Error(err)
		return err
	}
	length := len(ra)
	for i := 1; i < length; i++ {
		row := &ServerCabinetForImport{}
		if len(ra[i]) < 10 {
			continue
		}

		row.ServerRoom = strings.TrimSpace(ra[i][0])
		row.Number = strings.TrimSpace(ra[i][1])
		row.Height, _ = strconv.Atoi(strings.TrimSpace(ra[i][2]))
		row.Type = strings.TrimSpace(ra[i][3])
		row.NetworkArea = strings.TrimSpace(ra[i][4])
		row.NetworkRate = strings.TrimSpace(ra[i][5])
		row.Current = strings.TrimSpace(ra[i][6])
		row.AvailablePower = strings.TrimSpace(ra[i][7])
		row.MaxPower = strings.TrimSpace(ra[i][8])
		row.Remark = strings.TrimSpace(ra[i][9])

		//必填项校验
		row.checkLength()
		//机房和网络区域校验
		save, err := row.validate(log, repo)
		if err != nil {
			log.Error(err)
			return err
		}
		//不能获取机房，OperateLog，网络区域，就不能做保存操作
		if save == upload.Continue {
			continue
		}

		sc := &model.ServerCabinet{
			IDCID:          row.IDCID,
			ServerRoomID:   row.ServerRoomID,
			NetworkAreaID:  row.NetworkAreaID,
			Number:         row.Number,
			Height:         uint(row.Height),
			Type:           row.Type,
			NetworkRate:    row.NetworkRate,
			Current:        row.Current,
			AvailablePower: row.AvailablePower,
			MaxPower:       row.MaxPower,
			Remark:         row.Remark,
			Creator:        row.Creator,
		}
		sc.ID = row.ID
		if row.ID == 0 {
			sc.IsEnabled = model.NO
			sc.IsPowered = model.CabinetPowerOff
			sc.Status = model.CabinetStatUnderConstruction
			_, err := repo.SaveServerCabinet(sc)
			if err != nil {
				return err
			}
		}
		if row.ID > 0 {
			sc.ID = uint(row.ID)
			_, err := repo.SaveServerCabinet(sc)
			if err != nil {
				return err
			}
		}
	}
	defer os.Remove(upload.UploadDir + reqData.FileName)
	return nil
}

//GetServerByCabinetID 查询指定ID的机架(柜)信息详情
func GetServerByCabinetID(log logger.Logger, repo model.Repo, id uint) (*ServerCabinet, error) {
	sc, err := repo.GetServerCabinetByID(id)
	if err != nil {
		return nil, err
	}
	sca := convert2ServerCabinet(repo, sc)
	return sca, nil
}

//RemoveServerCabinetValidte 删除操作校验
func RemoveServerCabinetValidte(log logger.Logger, repo model.Repo, id uint) string {
	sc, _ := repo.GetServerCabinetByID(id)
	if sc == nil {
		return fmt.Sprintf("不存在机架(%d)", id)
	}
	//统计机架下面的中的机位
	count, _ := repo.CountServerUSite(&model.CombinedServerUSite{
		CabinetNumber: sc.Number,
	})
	if count > 0 {
		return fmt.Sprintf("机架下面存在机位(%d),不允许删除", count)
	}
	count, _ = repo.CountNetworkDevices(&model.NetworkDeviceCond{
		ServerCabinetID: []uint{id},
	})
	if count > 0 {
		return fmt.Sprintf("机架下面存在网络设备(%d),不允许删除", count)
	}
	return ""
}

//RemoveServerCabinetByID 删除指定ID的机架(柜)
func RemoveServerCabinetByID(log logger.Logger, repo model.Repo, id uint) error {
	//删除机架(柜)
	_, err := repo.RemoveServerCabinetByID(id)
	return err
}

//UpdateServerCabinetStatusReq 批量更新机架(柜)状态请求
type UpdateServerCabinetStatusReq struct {
	IDS    []uint `json:"ids"`
	Status string `json:"status"`
}

// FieldMap 请求字段映射
func (reqData *UpdateServerCabinetStatusReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDS:    "ids",
		&reqData.Status: "status",
	}
}

//Validate 参数校验
func (reqData *UpdateServerCabinetStatusReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	//这里的校验要与model状态保持一致
	if reqData.Status == "" {
		errs.Add([]string{"status"}, binding.RequiredError, "状态不能为空")
		return errs
	}

	//机架(柜)状态校验
	statusIsGood := true
	switch reqData.Status {
	case model.CabinetStatEnabled:
	case model.CabinetStatLocked:
	default:
		statusIsGood = false
	}
	if !statusIsGood {
		errs.Add([]string{"status"}, binding.RequiredError, "状态更新必须为(已启用 或 已锁定)")
		return errs
	}
	return errs
}

//AcceptServerCabinetStatusReq 批量验收
type AcceptServerCabinetStatusReq struct {
	IDS []uint `json:"ids"`
	//Status string `json:"status"`
}

// FieldMap 请求字段映射
func (reqData *AcceptServerCabinetStatusReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDS: "ids",
		//&reqData.Status: "status",
	}
}

//Validate 参数校验
func (reqData *AcceptServerCabinetStatusReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	//检查当前状态是否：建设中
	repo, _ := middleware.RepoFromContext(req.Context())
	for _, id := range reqData.IDS {
		if cabinet, err := repo.GetServerCabinetByID(id); err != nil {
			if cabinet.Status != model.CabinetStatUnderConstruction {
				errs.Add([]string{"status"}, binding.BusinessError,
					fmt.Sprintf("机架(编号:%s)当前状态：%s，必须是建设中(%s)", cabinet.Number,
						cabinet.Status, model.CabinetStatUnderConstruction))
				return errs
			}
		}
	}
	//检查tor组的两个机架是否关联网络设备
	// TODO
	return errs
}

//EnableServerCabinetStatusReq 批量启用
type EnableServerCabinetStatusReq struct {
	IDS []uint `json:"ids"`
}

// FieldMap 请求字段映射
func (reqData *EnableServerCabinetStatusReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDS: "ids",
	}
}

//Validate 参数校验
func (reqData *EnableServerCabinetStatusReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	//检查当前状态是否：未启用
	repo, _ := middleware.RepoFromContext(req.Context())
	for _, id := range reqData.IDS {
		if cabinet, err := repo.GetServerCabinetByID(id); err != nil {
			if cabinet.Status != model.CabinetStatNotEnabled {
				errs.Add([]string{"status"}, binding.BusinessError,
					fmt.Sprintf("机架(编号:%s)当前状态：%s，必须是未启用(%s)",
						cabinet.Number, cabinet.Status, model.CabinetStatNotEnabled))
				return errs
			}
		}
	}
	return errs
}

//OfflineServerCabinetStatusReq 批量下线
type OfflineServerCabinetStatusReq struct {
	IDS []uint `json:"ids"`
}

// FieldMap 请求字段映射
func (reqData *OfflineServerCabinetStatusReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDS: "ids",
	}
}

//Validate 参数校验
func (reqData *OfflineServerCabinetStatusReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	//检查当前状态是否：已启用
	repo, _ := middleware.RepoFromContext(req.Context())
	for _, id := range reqData.IDS {
		if cabinet, err := repo.GetServerCabinetByID(id); err != nil {
			if cabinet.Status != model.CabinetStatNotEnabled {
				errs.Add([]string{"status"}, binding.BusinessError,
					fmt.Sprintf("机架(编号:%s)当前状态：%s，必须是已启用(%s)",
						cabinet.Number, cabinet.Status, model.CabinetStatEnabled))
				return errs
			}
		}
	}
	return errs
}

//ReconstructServerCabinetStatusReq 批量重建
type ReconstructServerCabinetStatusReq struct {
	IDS []uint `json:"ids"`
}

// FieldMap 请求字段映射
func (reqData *ReconstructServerCabinetStatusReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDS: "ids",
	}
}

//Validate 参数校验
func (reqData *ReconstructServerCabinetStatusReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	//检查当前状态是否：已下线
	repo, _ := middleware.RepoFromContext(req.Context())
	for _, id := range reqData.IDS {
		if cabinet, err := repo.GetServerCabinetByID(id); err != nil {
			if cabinet.Status != model.CabinetStatOffline {
				errs.Add([]string{"status"}, binding.BusinessError,
					fmt.Sprintf("机架(编号:%s)当前状态：%s，必须是已下线(%s)",
						cabinet.Number, cabinet.Status, model.CabinetStatOffline))
				return errs
			}
		}
	}
	return errs
}

//AcceptServerCabinet 验收
func AcceptServerCabinet(log logger.Logger, repo model.Repo, req *AcceptServerCabinetStatusReq) error {
	updateReq := new(UpdateServerCabinetStatusReq)
	updateReq.IDS = req.IDS
	updateReq.Status = model.CabinetStatNotEnabled
	return UpdateServerCabinetStatus(log, repo, updateReq)
}

//EnableServerCabinet 启用
func EnableServerCabinet(log logger.Logger, repo model.Repo, req *EnableServerCabinetStatusReq) error {
	updateReq := new(UpdateServerCabinetStatusReq)
	updateReq.IDS = req.IDS
	updateReq.Status = model.CabinetStatEnabled
	return UpdateServerCabinetStatus(log, repo, updateReq)
}

//ReconstructServerCabinet 重建
func ReconstructServerCabinet(log logger.Logger, repo model.Repo, req *ReconstructServerCabinetStatusReq) error {
	updateReq := new(UpdateServerCabinetStatusReq)
	updateReq.IDS = req.IDS
	updateReq.Status = model.CabinetStatUnderConstruction
	return UpdateServerCabinetStatus(log, repo, updateReq)
}

//UpdateServerCabinetStatus 批量更新机架(柜)状态
func UpdateServerCabinetStatus(log logger.Logger, repo model.Repo, req *UpdateServerCabinetStatusReq) error {
	if _, err := repo.UpdateServerCabinetStatus(req.IDS, req.Status); err != nil {
		return err
	}
	return nil
}

// UpdateServerCabinetTypeReq
type UpdateServerCabinetTypeReq struct {
	// 类型
	Type      string             		`json:"type"`
	// id列表
	IDs       []uint             		`json:"ids"`
}

// BatchUpdateServerCabinetType 批量更新机架(柜)类型
func BatchUpdateServerCabinetType(repo model.Repo, ids []uint, typ string) (affected int64, err error) {
	return repo.UpdateServerCabinetType(ids, typ)
}

// UpdateServerCabinetRemarkReq
type UpdateServerCabinetRemarkReq struct {
	// 备注
	Remark      string             		`json:"remark"`
	// id列表
	IDs       	[]uint             		`json:"ids"`
}

// BatchUpdateServerCabinetRemark 批量更新机架(柜)备注信息
func BatchUpdateServerCabinetRemark(repo model.Repo, ids []uint, remark string) (affected int64, err error) {
	return repo.UpdateServerCabinetRemark(ids, remark)
}