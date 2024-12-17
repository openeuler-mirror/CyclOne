package service

import (
	"fmt"
	"net/http"
	"reflect"
	"net/url"
	"os"
	"strings"
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/voidint/binding"
	"github.com/voidint/page"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils"
	strings2 "idcos.io/cloudboot/utils/strings"
	"idcos.io/cloudboot/utils/upload"
)

// NetworkDevicePageReq 网络设备分页请求体
type NetworkDevicePageReq struct {
	// 所属数据中心ID
	IDCID string `json:"idc_id"`
	// 所属机房
	ServerRoomID   string `json:"server_room_id"`
	ServerRoomName string `json:"server_room_name"`
	// 所属机架
	ServerCabinetID     string `json:"server_cabinet_id"`
	ServerCabinetNumber string `json:"server_cabinet_number"`
	// 资产编号
	FixedAssetNumber string `json:"fixed_asset_number"`
	// SN号
	SN string `json:"sn"`
	// 名称
	Name string `json:"name"`
	// 产品型号
	ModelNumber string `json:"model"`
	// 厂商
	Vendor string `json:"vendor"`
	// 操作系统
	OS string `json:"os"`
	// 类型 switch-交换机;
	Type string `json:"type"`
	// TOR
	TOR string `json:"tor"`
	// 用途
	Usage string `json:"usage"`
	// 状态
	Status string `json:"status"`	
	// 页号
	Page int64 `json:"page"`
	// 页大小
	PageSize int64 `json:"page_size"`
}

// SaveNetworkDeviceReq 网络设备保存接口结构体
type SaveNetworkDeviceReq struct {
	// 所属数据中心ID
	IDCID uint `json:"idc_id"`
	// 所属机房
	ServerRoomID uint `json:"server_room_id"`
	// 所属机架
	ServerCabinetID uint `json:"server_cabinet_id"`
	// 资产编号
	FixedAssetNumber string `json:"fixed_asset_number"`
	// SN号
	SN string `json:"sn"`
	// 名称
	Name string `json:"name"`
	// 产品型号
	ModelNumber string `json:"model"`
	// 厂商
	Vendor string `json:"vendor"`
	// 操作系统
	OS string `json:"os"`
	// 类型 switch-交换机;
	Type string `json:"type"`
	// TOR
	TOR string `json:"tor"`
	// 用途
	Usage string `json:"usage"`
	// 状态
	Status string `json:"status"`
	// 用户登录名
	LoginName string `json:"-"`
}


// FieldMap 请求参数与结构体字段建立映射
func (reqData *NetworkDevicePageReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDCID:               "idc_id",
		&reqData.ServerRoomID:        "server_room_id",
		&reqData.ServerRoomName:      "server_room_name",
		&reqData.ServerCabinetID:     "server_cabinet_id",
		&reqData.ServerCabinetNumber: "server_cabinet_number",
		&reqData.FixedAssetNumber:    "fixed_asset_number",
		&reqData.SN:                  "sn",
		&reqData.Name:                "name",
		&reqData.ModelNumber:         "model",
		&reqData.Vendor:              "vendor",
		&reqData.TOR:                 "tor",
		&reqData.OS:                  "os",
		&reqData.Type:                "type",
		&reqData.Usage:               "usage",
		&reqData.Status:              "status",
		&reqData.Page:                "page",
		&reqData.PageSize:            "page_size",
	}
}


// FieldMap 请求参数与结构体字段建立映射
func (reqData *SaveNetworkDeviceReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDCID:            "idc_id",
		&reqData.ServerRoomID:     "server_room_id",
		&reqData.ServerCabinetID:  "server_cabinet_id",
		&reqData.FixedAssetNumber: "fixed_asset_number",
		&reqData.SN:               "sn",
		&reqData.Name:             "name",
		&reqData.ModelNumber:      "model",
		&reqData.Vendor:           "vendor",
		&reqData.OS:               "os",
		&reqData.Type:             "type",
		&reqData.Usage:            "usage",
		&reqData.Status:           "status",
	}
}

//
// Validate 结构体数据校验
func (reqData *SaveNetworkDeviceReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())

	//校验名称是否已经存在
	cond := &model.NetworkDeviceCond{
		Name: reqData.Name,
	}
	networks, _ := repo.GetNetworkDevicesByCond(cond, nil, nil)
	if len(networks) > 0 {
		errs.Add([]string{"name"}, binding.BusinessError, fmt.Sprintf("名称(%s)信息已经存在", reqData.SN))
		return errs
	}

	idc, err := repo.GetIDCByID(reqData.IDCID)
	if err == gorm.ErrRecordNotFound || idc == nil {
		errs.Add([]string{"idc_id"}, binding.BusinessError, "数据中心不存在")
		return errs
	}
	room, err := repo.GetServerRoomByID(reqData.ServerRoomID)
	if err == gorm.ErrRecordNotFound || room == nil {
		errs.Add([]string{"server_room_id"}, binding.BusinessError, "机房信息不存在")
		return errs
	}
	cabinet, err := repo.GetServerCabinetByID(reqData.ServerCabinetID)
	if err == gorm.ErrRecordNotFound || cabinet == nil {
		errs.Add([]string{"server_cabinet"}, binding.BusinessError, "机架信息不存在")
		return errs
	}

	items, err := repo.GetNetworkDeviceBySN(reqData.SN)
	if len(items) > 0 {
		errs.Add([]string{"sn"}, binding.BusinessError, fmt.Sprintf("SN(%s)信息已经存在", reqData.SN))
		return errs
	}

	items, err = repo.GetNetworkDeviceByFixedAssetNumber(reqData.FixedAssetNumber)
	if len(items) > 0 {
		errs.Add([]string{"fixed_asset_number"}, binding.BusinessError, fmt.Sprintf("资产编号(%s)信息已经存在", reqData.FixedAssetNumber))
		return errs
	}

	return errs
}

// NetworkDeviceResp 网络设备返回体
type NetworkDeviceResp struct {
	// 数据中心
	IDC struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"idc"`
	// 机房信息
	ServerRoom struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"server_room"`
	// 机架信息
	ServerCabinet struct {
		ID     uint   `json:"id"`
		Number string `json:"number"`
	} `json:"server_cabinet"`
	// 主键
	ID uint `json:"id"`
	// 资产编号
	FixedAssetNumber string `json:"fixed_asset_number"`
	// SN号
	SN string `json:"sn"`
	// 名称
	Name string `json:"name"`
	// 产品型号
	ModelNumber string `json:"model"`
	// 厂商
	Vendor string `json:"vendor"`
	// 操作系统
	OS string `json:"os"`
	// 类型
	Type string `json:"type"`
	// TOR
	TOR string `json:"tor"`
	// 用途
	Usage string `json:"usage"`
	// 状态
	Status string `json:"status"`	
	// 创建时间
	CreatedAt string `json:"created_at"`
	// 修改时间
	UpdatedAt string `json:"updated_at"`
}

// GetNetworkDevicesPage 返回满足过滤条件的网络设备(不支持模糊查找)
func GetNetworkDevicesPage(repo model.Repo, reqData *NetworkDevicePageReq) (pg *page.Page, err error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 10000 { //下拉需要
		reqData.PageSize = 10
	}
	if reqData.Page < 0 {
		reqData.Page = 0
	}

	cond := model.NetworkDeviceCond{
		IDCID:               strings2.Multi2UintSlice(reqData.IDCID),
		ServerRoomID:        strings2.Multi2UintSlice(reqData.ServerRoomID),
		ServerRoomName:      reqData.ServerRoomName,
		ServerCabinetID:     strings2.Multi2UintSlice(reqData.ServerCabinetID),
		ServerCabinetNumber: reqData.ServerCabinetNumber,
		FixedAssetNumber:    reqData.FixedAssetNumber,
		SN:                  reqData.SN,
		Name:                reqData.Name,
		ModelNumber:         reqData.ModelNumber,
		Vendor:              reqData.Vendor,
		OS:                  reqData.OS,
		Type:                reqData.Type,
		Usage: 				 reqData.Usage,
		Status:              reqData.Status,
	}

	//TOR包含+号，特殊处理下
	cond.TOR, _ = url.QueryUnescape(reqData.TOR)

	totalRecords, err := repo.CountNetworkDevices(&cond)
	if err != nil {
		return nil, err
	}

	pager := page.NewPager(reflect.TypeOf(&NetworkDeviceResp{}), reqData.Page, reqData.PageSize, totalRecords)
	items, err := repo.GetNetworkDevicesByCond(&cond, model.TwoOrderBy("idc_id", model.ASC, "server_cabinet_id", model.ASC),
		pager.BuildLimiter())
	if err != nil {
		return nil, err
	}

	for i := range items {
		pager.AddRecords(convert2NetworkDeviceResp(repo, items[i]))
	}

	return pager.BuildPage(), nil
}

// RemoveNetworkDeviceByID 删除指定ID的网络设备
func RemoveNetworkDeviceByID(repo model.Repo, id uint) (err error) {

	device, err := repo.GetNetworkDeviceByID(id)

	if err != nil {
		return err
	}

	// 校验网络设备是否被引用
	ipNets, err := repo.GetIPNetworksBySwitchNumber(device.FixedAssetNumber)
	if err != nil {
		return err
	}

	if len(ipNets) > 0 {
		return fmt.Errorf("网络设备(%s, %s)已经分配了网段，无法删除", device.Name, device.FixedAssetNumber)
	}

	return repo.RemoveNetworkDeviceByID(id)

}


// 批量删除网络设备请求结构体
type DelNetworkDeviceReq struct {
	IDs []uint `json:"ids"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *DelNetworkDeviceReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDs: "ids",
	}
}

// Validate 结构体数据校验
func (reqData *DelNetworkDeviceReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())
	for _, id := range reqData.IDs {

		device, err := repo.GetNetworkDeviceByID(id)
		if err != nil {
			errs.Add([]string{"id"}, binding.RequiredError, fmt.Sprintf("网络设备ID(%d)不存在", id))
			return errs
		}
	    // 校验网络设备是否被引用
	    ipNets, err := repo.GetIPNetworksBySwitchNumber(device.FixedAssetNumber)
	    if err != nil {
			errs.Add([]string{"id"}, binding.RequiredError, fmt.Sprintf("网络设备(%s)是否关联网段失败: %s", device.FixedAssetNumber, err.Error()))
	    	return errs
	    }
	    if len(ipNets) > 0 {
			errs.Add([]string{"id"}, binding.RequiredError, fmt.Sprintf("网络设备(%s, %s)存在关联网段,无法删除", device.Name, device.FixedAssetNumber))
	    	return errs
	    }	
	}
	return nil
}

//RemoveNetworkDevices 删除指定ID的网络设备
func RemoveNetworkDevices(log logger.Logger, repo model.Repo, reqData *DelNetworkDeviceReq) (affected int64, err error) {
	for _, id := range reqData.IDs {
		err = repo.RemoveNetworkDeviceByID(id)
		if err != nil {
			log.Errorf("delete network device(id=%d) fail,err:%v", id, err)
			return affected, err
		}

		affected++
	}
	return affected, err
}


// GetNetworkDeviceByID 查询指定ID的网络设备
func GetNetworkDeviceByID(repo model.Repo, id uint) (network *NetworkDeviceResp, err error) {
	networkDevice, err := repo.GetNetworkDeviceByID(id)
	if err != nil {
		return nil, err
	}
	return convert2NetworkDeviceResp(repo, networkDevice), nil
}

// SaveNetworkDevice 保存网络区域
func SaveNetworkDevice(log logger.Logger, repo model.Repo, reqData *SaveNetworkDeviceReq) (networkDevice *model.NetworkDevice, err error) {
	na := model.NetworkDevice{
		IDCID:            reqData.IDCID,
		ServerRoomID:     reqData.ServerRoomID,
		ServerCabinetID:  reqData.ServerCabinetID,
		Name:             reqData.Name,
		FixedAssetNumber: reqData.FixedAssetNumber,
		SN:               reqData.SN,
		ModelNumber:      reqData.ModelNumber,
		Vendor:           reqData.Vendor,
		OS:               reqData.OS,
		Type:             reqData.Type,
		TOR:              reqData.TOR,
		Usage:            reqData.Usage,
		Status:           reqData.Status,
	}
	return repo.SaveNetworkDevice(&na)
}


// convert2NetworkDeviceResp 转换设备信息
func convert2NetworkDeviceResp(repo model.Repo, item *model.NetworkDevice) *NetworkDeviceResp {
	device := NetworkDeviceResp{
		ID:               item.ID,
		CreatedAt:        item.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        item.UpdatedAt.Format("2006-01-02 15:04:05"),
		FixedAssetNumber: item.FixedAssetNumber,
		SN:               item.SN,
		Name:             item.Name,
		ModelNumber:      item.ModelNumber,
		Vendor:           item.Vendor,
		Type:             item.Type,
		TOR:              item.TOR,
		OS:               item.OS,
		Usage:            item.Usage,
		Status:           item.Status,
	}

	if idc, _ := repo.GetIDCByID(item.IDCID); idc != nil {
		device.IDC.ID = idc.ID
		device.IDC.Name = idc.Name
	}

	if room, _ := repo.GetServerRoomByID(item.ServerRoomID); room != nil {
		device.ServerRoom.ID = room.ID
		device.ServerRoom.Name = room.Name
	}
	if cabinet, _ := repo.GetServerCabinetByID(item.ServerCabinetID); cabinet != nil {
		device.ServerCabinet.ID = cabinet.ID
		device.ServerCabinet.Number = cabinet.Number
	}

	return &device
}

//TOR组相关的逻辑，由于TOR是交换机的一个属性，暂时把放到这里
//GetTORBySN 根据设备SN查询所属的TOR组
func GetTORBySN(log logger.Logger, repo model.Repo, SN string) (tor string, err error) {
	netDev, err := repo.GetIntranetSwitchBySN(SN)
	if err != nil {
		log.Errorf("get tor by sn:%s fail,%s", SN, err.Error())
		return "", err
	}
	return netDev.TOR, nil
}

//GetAllTORs 查询所有的TOR组名称
func GetAllTORs(log logger.Logger, repo model.Repo) []string {
	tors, err := repo.GetTORs()
	if err != nil {
		log.Errorf("get all tors fail,%s", err.Error())
		return nil
	}
	return tors
}

type NetworkDeviceForImport struct {
	ID             uint
	IDCName        string `json:"idc_name"`
	IDCID          uint
	ServerRoomName string `json:"server_room_name"`
	ServerRoomID   uint
	// 所属机架
	ServerCabinetNumber string `json:"server_cabinet_number"`
	ServerCabinetID     uint
	// 资产编号
	FixedAssetNumber string `json:"fixed_asset_number"`
	// SN号
	SN string `json:"sn"`
	// 名称
	Name string `json:"name"`
	// 产品型号
	ModelNumber string `json:"model"`
	// 厂商
	Vendor string `json:"vendor"`
	// 操作系统
	OS string `json:"os"`
	// 类型 switch-交换机;
	Type string `json:"type"`
	// TOR
	TOR string `json:"tor"`
	// 用途
	Usage string `json:"usage"`

	Content string `json:"content"`
}

//checkLength 对导入文件中的数据做基本验证
func (srfi *NetworkDeviceForImport) checkLength() {
	leg := len(srfi.IDCName)
	if leg == 0 || leg > 255 {
		var br string
		if srfi.Content != "" {
			br = "<br />"
		}
		srfi.Content += br + fmt.Sprintf("必填项校验:数据中心长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(srfi.ServerRoomName)
	if leg == 0 || leg > 255 {
		var br string
		if srfi.Content != "" {
			br = "<br />"
		}
		srfi.Content += br + fmt.Sprintf("必填项校验:机房管理单元长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(srfi.ServerCabinetNumber)
	if leg == 0 || leg > 255 {
		var br string
		if srfi.Content != "" {
			br = "<br />"
		}
		srfi.Content += br + fmt.Sprintf("必填项校验:所属机架编号长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(srfi.FixedAssetNumber)
	if leg == 0 || leg > 255 {
		var br string
		if srfi.Content != "" {
			br = "<br />"
		}
		srfi.Content += br + fmt.Sprintf("必填项校验:固资编号长度为(%d)(不能为空)", leg)
	}
	leg = len(srfi.SN)
	if leg == 0 || leg > 255 {
		var br string
		if srfi.Content != "" {
			br = "<br />"
		}
		srfi.Content += br + fmt.Sprintf("必填项校验:序列号长度为(%d)(不能为空)", leg)
	}
	leg = len(srfi.Name)
	if leg == 0 || leg > 255 {
		var br string
		if srfi.Content != "" {
			br = "<br />"
		}
		srfi.Content += br + fmt.Sprintf("必填项校验:设备名称长度为(%d)(不能为空)", leg)
	}
	leg = len(srfi.ModelNumber)
	if leg == 0 || leg > 255 {
		var br string
		if srfi.Content != "" {
			br = "<br />"
		}
		srfi.Content += br + fmt.Sprintf("必填项校验:型号长度为(%d)(不能为空)", leg)
	}
	leg = len(srfi.Vendor)
	if leg == 0 || leg > 255 {
		var br string
		if srfi.Content != "" {
			br = "<br />"
		}
		srfi.Content += br + fmt.Sprintf("必填项校验:厂商长度为(%d)(不能为空)", leg)
	}
	leg = len(srfi.OS)
	if leg == 0 || leg > 255 {
		var br string
		if srfi.Content != "" {
			br = "<br />"
		}
		srfi.Content += br + fmt.Sprintf("必填项校验:操作系统长度为(%d)(不能为空)", leg)
	}
	leg = len(srfi.Type)
	if leg == 0 || leg > 255 {
		var br string
		if srfi.Content != "" {
			br = "<br />"
		}
		srfi.Content += br + fmt.Sprintf("必填项校验:类型长度为(%d)(不能为空)", leg)
	}
	leg = len(srfi.TOR)
	if leg == 0 || leg > 255 {
		var br string
		if srfi.Content != "" {
			br = "<br />"
		}
		srfi.Content += br + fmt.Sprintf("必填项校验:TOR长度为(%d)(不能为空)", leg)
	}
	leg = len(srfi.Usage)
	if leg == 0 || leg > 255 {
		var br string
		if srfi.Content != "" {
			br = "<br />"
		}
		srfi.Content += br + fmt.Sprintf("必填项校验:用途长度为(%d)(不能为空)", leg)
	}
}

//validate 对导入文件中的数据做基本验证
func (srfi *NetworkDeviceForImport) validate(repo model.Repo) (int, error) {
	//数据中心校验
	idc, err := repo.GetIDCByName(srfi.IDCName)
	if err != nil && err != gorm.ErrRecordNotFound {
		return upload.Return, err
	}
	if err == gorm.ErrRecordNotFound || idc == nil {
		var br string
		if srfi.Content != "" {
			br = "<br />"
		}
		srfi.Content += br + fmt.Sprintf("数据中心(%s)不存在", srfi.IDCName)
		return upload.Return, nil
	}

	//机房校验
	srs, err := repo.GetServerRoomByName(srfi.ServerRoomName)
	if err != nil && err != gorm.ErrRecordNotFound {
		return upload.Return, err
	}
	if err == gorm.ErrRecordNotFound || srs == nil {
		var br string
		if srfi.Content != "" {
			br = "<br />"
		}
		srfi.Content += br + fmt.Sprintf("机房管理单元(%s)不存在", srfi.ServerRoomName)
		return upload.Return, nil
	} else {
		srfi.ServerRoomID = srs.ID
	}
	srfi.IDCID = idc.ID
	cabinet, err := repo.GetServerCabinetByNumber(srfi.ServerRoomID, srfi.ServerCabinetNumber)
	if err == gorm.ErrRecordNotFound || cabinet == nil {
		var br string
		if srfi.Content != "" {
			br = "<br />"
		}
		srfi.Content += br + fmt.Sprintf("机架编号(%s)对应的机架不存在", srfi.ServerCabinetNumber)
		return upload.Return, nil
	} else {
		srfi.ServerCabinetID = cabinet.ID
	}

	items, err := repo.GetNetworkDeviceByFixedAssetNumber(srfi.FixedAssetNumber)
	if len(items) > 0 {
		srfi.ID = items[0].ID
	}

	//类型目前只有交换机
	if srfi.Type != "交换机" {
		return upload.Return, errors.New("类型须为：交换机")
	}
	items, err = repo.GetNetworkDevicesByCond(&model.NetworkDeviceCond{Name: srfi.Name}, nil, nil)
	if len(items) > 0 {
		if items[0].FixedAssetNumber != srfi.FixedAssetNumber {
			return upload.Return, fmt.Errorf("与固资号(%s)的名称重复", items[0].FixedAssetNumber)
		}
	}

	return upload.DO, nil
}

//ImportNetworkDevicePreview 导入预览
func ImportNetworkDevicePreview(log logger.Logger, repo model.Repo, reqData *upload.ImportReq) (map[string]interface{}, error) {
	ra, err := utils.ParseDataFromXLSX(upload.UploadDir + reqData.FileName)
	if err != nil {
		return nil, err
	}
	length := len(ra)

	var success []*NetworkDeviceForImport
	var failure []*NetworkDeviceForImport
	for i := 1; i < length; i++ {
		row := &NetworkDeviceForImport{}
		if len(ra[i]) < 12 {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "导入文件列长度不对（应为12列）"
			failure = append(failure, row)
			continue
		}

		row.IDCName = strings.TrimSpace(ra[i][0])
		row.ServerRoomName = strings.TrimSpace(ra[i][1])
		row.ServerCabinetNumber = strings.TrimSpace(ra[i][2])
		row.FixedAssetNumber = strings.TrimSpace(ra[i][3])
		row.SN = strings.TrimSpace(ra[i][4])
		row.Name = strings.TrimSpace(ra[i][5])
		row.ModelNumber = strings.TrimSpace(ra[i][6])
		row.Vendor = strings.TrimSpace(ra[i][7])
		row.OS = strings.TrimSpace(ra[i][8])
		row.Type = strings.TrimSpace(ra[i][9])
		row.TOR = strings.TrimSpace(ra[i][10])
		row.Usage = strings.TrimSpace(ra[i][11])

		//必填项校验
		row.checkLength()
		//机房和数据中心校验
		_, err := row.validate(repo)
		if err != nil {
			return nil, err
		}

		if row.Content != "" {
			failure = append(failure, row)
		} else {
			success = append(success, row)
		}
	}

	var data []*NetworkDeviceForImport
	if len(failure) > 0 {
		data = failure
	} else {
		data = success
	}
	var result []*NetworkDeviceForImport
	for i := 0; i < len(data); i++ {
		if uint(i) >= reqData.Offset && uint(i) < (reqData.Offset+reqData.Limit) {
			result = append(result, data[i])
		}
	}
	if len(failure) > 0 {
		_ = os.Remove(upload.UploadDir + reqData.FileName)
		return map[string]interface{}{"status": "failure",
			"message":       "导入网络设备错误",
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

//ImportNetworkDevices 将导入机房放到数据库
func ImportNetworkDevices(log logger.Logger, repo model.Repo, reqData *upload.ImportReq) error {
	ra, err := utils.ParseDataFromXLSX(upload.UploadDir + reqData.FileName)
	if err != nil {
		return err
	}
	length := len(ra)
	for i := 1; i < length; i++ {
		row := &NetworkDeviceForImport{}
		if len(ra[i]) < 12 {
			continue
		}

		row.IDCName = strings.TrimSpace(ra[i][0])
		row.ServerRoomName = strings.TrimSpace(ra[i][1])
		row.ServerCabinetNumber = strings.TrimSpace(ra[i][2])
		row.FixedAssetNumber = strings.TrimSpace(ra[i][3])
		row.SN = strings.TrimSpace(ra[i][4])
		row.Name = strings.TrimSpace(ra[i][5])
		row.ModelNumber = strings.TrimSpace(ra[i][6])
		row.Vendor = strings.TrimSpace(ra[i][7])
		row.OS = strings.TrimSpace(ra[i][8])
		row.Type = strings.TrimSpace(ra[i][9])
		row.TOR = strings.TrimSpace(ra[i][10])
		row.Usage = strings.TrimSpace(ra[i][11])

		//必填项校验
		row.checkLength()
		//机房和数据中心校验
		isSave, err := row.validate(repo)
		if err != nil {
			return err
		}
		if isSave == upload.Continue {
			continue
		}

		sr := &model.NetworkDevice{
			IDCID:            row.IDCID,
			ServerRoomID:     row.ServerRoomID,
			ServerCabinetID:  row.ServerCabinetID,
			FixedAssetNumber: row.FixedAssetNumber,
			SN:               row.SN,
			Name:             row.Name,
			ModelNumber:      row.ModelNumber,
			Vendor:           row.Vendor,
			OS:               row.OS,
			Type:             "switch", //row.Type,
			TOR:              row.TOR,
			Usage:            row.Usage,
		}
		sr.Model.ID = row.ID

		_, err = repo.SaveNetworkDevice(sr)
		if err != nil {
			return err
		}
	}
	defer os.Remove(upload.UploadDir + reqData.FileName)
	return nil
}
