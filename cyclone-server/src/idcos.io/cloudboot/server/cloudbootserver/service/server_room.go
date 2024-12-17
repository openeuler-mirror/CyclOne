package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/voidint/binding"
	"github.com/voidint/page"

	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils"
	strings2 "idcos.io/cloudboot/utils/strings"
	"idcos.io/cloudboot/utils/times"
	"idcos.io/cloudboot/utils/upload"
)

//SaveServerRoomReq 保存机房请求参数
type SaveServerRoomReq struct {
	//(required) 所属数据中心ID
	IDCID uint `json:"idc_id"`
	//机房ID。若id=0，则新增。若id>0，则修改。
	ID uint `json:"id"`
	//required) 机房名称
	Name string `json:"name"`
	//(required)一级机房ID
	FirstServerRoom uint `json:"first_server_room"`
	//(required) 所属城市
	City string `json:"city"`
	//(required) 地址
	Address string `json:"address"`
	//(required) 机房管理单元负责人
	ServerRoomManager string `json:"server_room_manager"`
	//(required) 供应商负责人
	VendorManager string `json:"vendor_manager"`
	//(required) 网络资产负责人
	NetworkAssetManager string `json:"network_asset_manager"`
	//(required) 7*24小时保障电话
	SupportPhoneNumber string `json:"support_phone_number"`
	// 用户登录名
	LoginName string `json:"-"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *SaveServerRoomReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDCID:               "idc_id",
		&reqData.ID:                  "id",
		&reqData.Name:                "name",
		&reqData.FirstServerRoom:     "first_server_room",
		&reqData.City:                "city",
		&reqData.Address:             "address",
		&reqData.ServerRoomManager:   "server_room_manager",
		&reqData.VendorManager:       "vendor_manager",
		&reqData.NetworkAssetManager: "network_asset_manager",
		&reqData.SupportPhoneNumber:  "support_phone_number",
	}
}

var (
	//中国区手机/固话号码
	chinaTelephoneNumberReq = regexp.MustCompile("^(((\\+\\d{2}-)?0\\d{2,3}-\\d{7,8})|((\\+\\d{2}-)?(\\d{2,3}-)?([1][3,4,5,7,8][0-9]\\d{8})))$")
)

// Validate 结构体数据校验
func (reqData *SaveServerRoomReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	//必须的基本数据不能为空
	if errs = reqData.baseValidate(req, errs); errs != nil {
		return errs
	}

	//更新机房信息，校验指定ID的机房是否存在
	if reqData.ID > 0 {
		if errs = reqData.serverRoomValidate(req, errs); errs != nil {
			return errs
		}
	}

	//校验IDC数据
	if errs = reqData.checkIDCValidate(req, errs); errs != nil {
		return errs
	}

	//校验指定7*24小时电话格式是否正确
	phlist := strings.Split(reqData.SupportPhoneNumber, ",")
	for k := range phlist {
		if !chinaTelephoneNumberReq.MatchString(phlist[k]) {
			errs.Add([]string{"phone"}, binding.RequiredError, fmt.Sprintf("7*24小时电话(%s)格式不正确", phlist[k]))
			return errs
		}
	}
	// TODO 校验指定各块的负责人是否存在
	return errs
}

func (reqData *SaveServerRoomReq) checkIDCValidate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())

	idc, err := repo.GetIDCByID(reqData.IDCID)
	if err != nil {
		errs.Add([]string{"idc"}, binding.RequiredError, fmt.Sprintf("获取数据中心id(%d)出现错误: %s", reqData.IDCID, err.Error()))
		return errs
	}

	var fsr []IDCFirstServerRoom
	if idc.FirstServerRoom != "" {
		err := json.Unmarshal([]byte(idc.FirstServerRoom), &fsr)
		if err != nil {
			errs.Add([]string{"idc"}, binding.RequiredError, fmt.Sprintf("获取数据中心id(%d)(一级机房)出现错误: %s", reqData.IDCID, err.Error()))
			return errs
		}
	}

	isExist := false
	for k := range fsr {
		if reqData.FirstServerRoom == fsr[k].ID {
			isExist = true
			break
		}
	}
	if !isExist {
		errs.Add([]string{"idc"}, binding.RequiredError, fmt.Sprintf("获取数据中心(一级机房)不存在ID为 %d 的", reqData.FirstServerRoom))
		return errs
	}
	return errs

}

func (reqData *SaveServerRoomReq) serverRoomValidate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())

	var cond model.ServerRoomCond

	cond.ID = []uint{reqData.ID}

	items, err := repo.GetServerRooms(&cond, nil, nil)
	if err != nil {
		errs.Add([]string{"server_room"}, binding.RequiredError, fmt.Sprintf("更新机房ID %d不存在", reqData.ID))
		return errs
	}
	if len(items) < 1 {
		errs.Add([]string{"server_room"}, binding.RequiredError, fmt.Sprintf("更新机房ID %d不存在", reqData.ID))
		return errs
	}

	//校验更新机房名是否已经存在
	if reqData.Name != "" {
		item, _ := repo.GetServerRoomByName(reqData.Name)
		if item != nil && item.ID != reqData.ID {
			errs.Add([]string{"server_room"}, binding.RequiredError, fmt.Sprintf("更新机房 指定的机房名%s已经存在", reqData.Name))
			return errs
		}
	}
	return errs
}

//baseValidate 必要参数不能为空
func (reqData *SaveServerRoomReq) baseValidate(req *http.Request, errs binding.Errors) binding.Errors {
	if reqData.IDCID == 0 {
		errs.Add([]string{"idc_id"}, binding.RequiredError, "数据中心ID不能为空")
		return errs
	}

	if reqData.Name == "" {
		errs.Add([]string{"name"}, binding.RequiredError, "机房名称不能为空")
		return errs
	}

	if reqData.City == "" {
		errs.Add([]string{"city"}, binding.RequiredError, "机房所属城市不能为空")
		return errs
	}

	if reqData.Address == "" {
		errs.Add([]string{"address"}, binding.RequiredError, "机房地址不能为空")
		return errs
	}

	if reqData.ServerRoomManager == "" {
		errs.Add([]string{"server_room_manager"}, binding.RequiredError, "机房管理单元负责人不能为空")
		return errs
	}

	if reqData.VendorManager == "" {
		errs.Add([]string{"vendor_manager"}, binding.RequiredError, "供应商负责人不能为空")
		return errs
	}

	if reqData.NetworkAssetManager == "" {
		errs.Add([]string{"network_asset_manager"}, binding.RequiredError, "网络资产负责人不能为空")
		return errs
	}

	if reqData.SupportPhoneNumber == "" {
		errs.Add([]string{"support_phone_number"}, binding.RequiredError, "7*24小时保障电话不能为空")
		return errs
	}
	return errs
}

//SaveServerRoom 保存机房
func SaveServerRoom(log logger.Logger, repo model.Repo, reqData *SaveServerRoomReq) error {
	sr := model.ServerRoom{
		IDCID:               reqData.IDCID,
		Name:                reqData.Name,
		FirstServerRoom:     reqData.FirstServerRoom,
		City:                reqData.City,
		Address:             reqData.Address,
		ServerRoomManager:   reqData.ServerRoomManager,
		VendorManager:       reqData.VendorManager,
		NetworkAssetManager: reqData.NetworkAssetManager,
		SupportPhoneNumber:  reqData.SupportPhoneNumber,
		Creator:             reqData.LoginName,
	}
	if reqData.ID == 0 {
		sr.Status = model.RoomStatUnderConstruction
		_, err := repo.SaveServerRoom(&sr)
		if err != nil {
			return err
		}
	}
	if reqData.ID > 0 {
		sr.ID = uint(reqData.ID)
		_, err := repo.UpdateServerRoom([]*model.ServerRoom{&sr})
		if err != nil {
			return err
		}
	}
	reqData.ID = sr.ID
	return nil
}

//UpdateServerRoomStateReq 更新机房状态请求参数
type UpdateServerRoomStateReq struct {
	//(required): 目标机房状态。可选值: under_construction-建设中; accepted-已验收; production-已投产; abolished-已裁撤;
	Status string `json:"status"`
	//(required): 机房ID列表
	IDs []int `json:"ids"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *UpdateServerRoomStateReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.Status: "status",
		&reqData.IDs:    "ids",
	}
}

// Validate 结构体数据校验
func (reqData *UpdateServerRoomStateReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if reqData.Status == "" {
		errs.Add([]string{"status"}, binding.RequiredError, "目标机房状态不能为空")
		return errs
	}

	//这里的校验要与model状态保持一致
	if reqData.Status != model.RoomStatUnderConstruction &&
		reqData.Status != model.RoomStatAccepted &&
		reqData.Status != model.RoomStatProduction {
		//reqData.Status != model.RoomStatAbolished
		errs.Add([]string{"status"}, binding.RequiredError, "状态值错误(只能为:under_construction/accepted/production中的一种)")
		return errs
	}
	return errs

}

//UpdateServerRoomStatus 批量修改机房状态
func UpdateServerRoomStatus(log logger.Logger, repo model.Repo, reqData *UpdateServerRoomStateReq) error {
	var srs []*model.ServerRoom
	for k := range reqData.IDs {
		//实现方式不好，如果机房比较多会造成大量变量，对GC造成压力(这里sr会发生逃逸，放到堆上)
		sr := &model.ServerRoom{
			Status: reqData.Status,
		}
		sr.ID = uint(reqData.IDs[k])
		srs = append(srs, sr)
	}
	if _, err := repo.UpdateServerRoom(srs); err != nil {
		return err
	}
	return nil
}

//RemoveServerRoomValidte 删除操作校验
func RemoveServerRoomValidte(log logger.Logger, repo model.Repo, id uint) string {
	//统计机房中的机架(柜数)
	count, _ := repo.GetServerCabinetCountByServerRoomID(id)
	if count > 0 {
		return fmt.Sprintf("机房下面存在(%d)个机架,不允许删除", count)
	}
	count, _ = repo.CountNetworkAreas(&model.NetworkAreaCond{
		ServerRoomID: []uint{id},
	})
	if count > 0 {
		return fmt.Sprintf("机房下面存在(%d)个网络区域,不允许删除", count)
	}
	count, _ = repo.CountIPNetworks(&model.IPNetworkCond{ServerRoomID: []uint{id}})
	if count > 0 {
		return fmt.Sprintf("机房下面存在(%d)网段,不允许删除", count)
	}
	return ""
}

//RemoveServerRoomByID 删除指定ID的机房
func RemoveServerRoomByID(log logger.Logger, repo model.Repo, id uint) error {
	_, err := repo.RemoveServerRoomByID(id)
	return err
}

//GetServerRoomPageReq 获取机房分页请求参数
type GetServerRoomPageReq struct {
	// 所属数据中心ID
	IDCID string `json:"idc_id"`
	//  机房名称(支持模糊匹配)
	Name string `json:"name"`
	// 一级机房ID
	FirstServerRoom string `json:"first_server_room"`
	// 城市(支持模糊匹配)
	City string `json:"city"`
	// 地址(支持模糊匹配)
	Address string `json:"address"`
	// 机房管理单元负责人(支持模糊匹配)
	ServerRoomManager string `json:"server_room_manager"`
	// 供应商负责人(支持模糊匹配)
	VendorManager string `json:"vendor_manager"`
	// 网络资产负责人(支持模糊匹配)
	NetworkAssetManager string `json:"network_asset_manager"`
	// 7*24小时保障电话(支持模糊匹配)
	SupportPhoneNumber string `json:"support_phone_number"`
	//状态。可选值: under_construction-建设中; accepted-已验收; production-已投产; abolished-已裁撤;
	Status string `json:"status"`
	// 分页页号
	Page int64 `json:"page"`
	// 分页大小。默认值:10。阈值: 100。
	PageSize int64 `json:"page_size"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *GetServerRoomPageReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDCID:               "idc_id",
		&reqData.Status:              "status",
		&reqData.Name:                "name",
		&reqData.FirstServerRoom:     "first_server_room",
		&reqData.City:                "city",
		&reqData.Address:             "address",
		&reqData.ServerRoomManager:   "server_room_manager",
		&reqData.VendorManager:       "vendor_manager",
		&reqData.NetworkAssetManager: "network_asset_manager",
		&reqData.SupportPhoneNumber:  "support_phone_number",
		&reqData.Page:                "page",
		&reqData.PageSize:            "page_size",
	}
}

//IDCFirstServerRoom 一级机房
type IDCFirstServerRoom struct {
	//一级机房ID
	ID uint `json:"id"`
	//一级机房名
	Name string `json:"name"`
}

//IDCForServerRoomPage 数据中心
type IDCForServerRoomPage struct {
	//数据中心ID
	ID uint `json:"id"`
	//数据中心名
	Name string `json:"name"`
	//一级机房
	FirstServerRoom []IDCFirstServerRoom `json:"first_server_room"`
	//厂商
	Vendor string `json:"vendor"`
	//状态
	Status string `json:"status"`
}

//ServerRoom 机房分页查询信息
type ServerRoom struct {
	//机房ID
	ID uint `json:"id"`
	//创建时间
	CreatedAt times.ISOTime `json:"created_at"`
	//更新时间
	UpdatedAt times.ISOTime `json:"updated_at"`
	//机房名
	Name string `json:"name"`
	//一级机房ID
	FirstServerRoom IDCFirstServerRoom `json:"first_server_room"`
	//机架数
	CabinetCount int64 `json:"cabinet_count"`
	//所在城市
	City string `json:"city"`
	//地址
	Address string `json:"address"`
	//机房管理人
	ServerRoomManager string `json:"server_room_manager"`
	//供应商负责人
	VendorManager string `json:"vendor_manager"`
	//网络资产管理人
	NetworkAssetManager string `json:"network_asset_manager"`
	//支撑电话号码
	SupportPhoneNumber string `json:"support_phone_number"`
	//机房状态
	Status string `json:"status"`
	//创建人
	Creator string `json:"creator"`
	//数据中心
	IDC IDCForServerRoomPage `json:"idc"`
}

//GetServerRoomWithPage 获取机房分页
func GetServerRoomWithPage(log logger.Logger, repo model.Repo, reqData *GetServerRoomPageReq) (*page.Page, error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 100 {
		reqData.PageSize = 20
	}
	if reqData.Page < 0 {
		reqData.Page = 0
	}

	cond := model.ServerRoomCond{
		IDCID:               strings2.Multi2UintSlice(reqData.IDCID),
		Name:                reqData.Name,
		FirstServerRoom:     strings2.Multi2UintSlice(reqData.FirstServerRoom),
		City:                reqData.City,
		Address:             reqData.Address,
		ServerRoomManager:   reqData.ServerRoomManager,
		VendorManager:       reqData.VendorManager,
		NetworkAssetManager: reqData.NetworkAssetManager,
		SupportPhoneNumber:  reqData.SupportPhoneNumber,
		Status:              reqData.Status,
	}

	totalRecords, err := repo.CountServerRooms(&cond)
	if err != nil {
		return nil, err
	}

	pager := page.NewPager(reflect.TypeOf(&ServerRoom{}), reqData.Page, reqData.PageSize, totalRecords)
	items, err := repo.GetServerRooms(&cond, model.OneOrderBy("id", model.DESC), pager.BuildLimiter())
	if err != nil {
		return nil, err
	}
	for i := range items {
		item, err := convert2ServerRoomResult(log, repo, items[i])
		if err != nil {
			return nil, err
		}
		if item != nil {
			pager.AddRecords(item)
		}
	}
	return pager.BuildPage(), nil
}

func convert2ServerRoomResult(log logger.Logger, repo model.Repo, srcp *model.ServerRoom) (*ServerRoom, error) {
	if srcp == nil {
		return nil, nil
	}

	idc, err := repo.GetIDCByID(srcp.IDCID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	result := ServerRoom{
		ID:                  srcp.ID,
		CreatedAt:           times.ISOTime(srcp.CreatedAt),
		UpdatedAt:           times.ISOTime(srcp.UpdatedAt),
		Name:                srcp.Name,
		City:                srcp.City,
		Address:             srcp.Address,
		ServerRoomManager:   srcp.ServerRoomManager,
		VendorManager:       srcp.VendorManager,
		NetworkAssetManager: srcp.NetworkAssetManager,
		SupportPhoneNumber:  srcp.SupportPhoneNumber,
		Status:              srcp.Status,
		Creator:             srcp.Creator,
	}

	if idc != nil {
		var fsr []IDCFirstServerRoom
		if idc.FirstServerRoom != "" {
			err := json.Unmarshal([]byte(idc.FirstServerRoom), &fsr)
			if err != nil {
				return nil, err
			}
		}
		result.IDC = IDCForServerRoomPage{
			ID:              srcp.IDCID,
			Name:            idc.Name,
			FirstServerRoom: fsr,
			Vendor:          idc.Vendor,
			Status:          idc.Status,
		}

		for k := range fsr {
			if fsr[k].ID == srcp.FirstServerRoom {
				result.FirstServerRoom = fsr[k]
			}
		}
	}

	//统计机房中的机架(柜数)
	count, err := repo.GetServerCabinetCountByServerRoomID(result.ID)
	if err != nil {
		return nil, err
	}
	result.CabinetCount = count

	return &result, nil
}

//GetServerRoomByID 获取指定ID的机房的详细信息
func GetServerRoomByID(log logger.Logger, repo model.Repo, id uint) (*ServerRoom, error) {
	items, err := repo.GetServerRoomByID(id)
	if err != nil {
		return nil, err
	}

	item, err := convert2ServerRoomResult(log, repo, items)
	if err != nil {
		return nil, err
	}
	return item, nil
}

//ServerRoomForImport 为了导入
type ServerRoomForImport struct {
	//文件导入数据
	ID uint `json:"-"`
	//机房名
	Name string `json:"name"`
	//数据中心名
	IDCName string `json:"idc_name"`
	//一级机房名
	FirstServerRoomName string `json:"first_server_room_name"`
	//所在城市
	City string `json:"city"`
	//地址
	Address string `json:"address"`
	//机房负责人
	ServerRoomManager string `json:"server_room_manager"`
	//供应商负责人
	VendorManager string `json:"vendor_manager"`
	//网络资产负责人
	NetworkAssetManager string `json:"network_asset_manager"`
	//支撑电话号码
	SupportPhoneNumber string `json:"support_phone_number"`

	//创建人
	Creator string `json:"creator"`
	//机房ID
	IDCID uint `json:"idc_id"`
	//一级机房
	FirstServerRoom uint `json:"first_server_room"`
	//附加内容
	Content string `json:"content"`
}

//checkLength 对导入文件中的数据做基本验证
func (srfi *ServerRoomForImport) checkLength() {
	leg := len(srfi.Name)
	if leg == 0 || leg > 255 {
		var br string
		if srfi.Content != "" {
			br = "<br />"
		}
		srfi.Content += br + fmt.Sprintf("必填项校验:机房管理单元长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(srfi.IDCName)
	if leg == 0 || leg > 255 {
		var br string
		if srfi.Content != "" {
			br = "<br />"
		}
		srfi.Content += br + fmt.Sprintf("必填项校验:数据中心长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(srfi.FirstServerRoomName)
	if leg == 0 || leg > 255 {
		var br string
		if srfi.Content != "" {
			br = "<br />"
		}
		srfi.Content += br + fmt.Sprintf("必填项校验:所属一级机房长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(srfi.City)
	if leg == 0 || leg > 255 {
		var br string
		if srfi.Content != "" {
			br = "<br />"
		}
		srfi.Content += br + fmt.Sprintf("必填项校验:机房管理单元所属城市长度为(%d)(不能为空)", leg)
	}
	leg = len(srfi.Address)
	if leg == 0 || leg > 255 {
		var br string
		if srfi.Content != "" {
			br = "<br />"
		}
		srfi.Content += br + fmt.Sprintf("必填项校验:机房管理单元地址长度为(%d)(不能为空)", leg)
	}
	leg = len(srfi.ServerRoomManager)
	if leg == 0 || leg > 255 {
		var br string
		if srfi.Content != "" {
			br = "<br />"
		}
		srfi.Content += br + fmt.Sprintf("必填项校验:机房管理单元负责人长度为(%d)(不能为空)", leg)
	}
	leg = len(srfi.VendorManager)
	if leg == 0 || leg > 255 {
		var br string
		if srfi.Content != "" {
			br = "<br />"
		}
		srfi.Content += br + fmt.Sprintf("必填项校验:供应商负责人长度为(%d)(不能为空)", leg)
	}
	leg = len(srfi.NetworkAssetManager)
	if leg == 0 || leg > 255 {
		var br string
		if srfi.Content != "" {
			br = "<br />"
		}
		srfi.Content += br + fmt.Sprintf("必填项校验:网络资产负责人长度为(%d)(不能为空)", leg)
	}
	leg = len(srfi.SupportPhoneNumber)
	if leg == 0 || leg > 255 {
		var br string
		if srfi.Content != "" {
			br = "<br />"
		}
		srfi.Content += br + fmt.Sprintf("必填项校验:7*24小时保障电话长度为(%d)(不能为空)", leg)
	}
}

//validate 对导入文件中的数据做基本验证
func (srfi *ServerRoomForImport) validate(repo model.Repo) (int, error) {
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
		return upload.Continue, nil
	}

	//一级机房校验
	fsrExist := false
	var fsr []IDCFirstServerRoom
	if idc.FirstServerRoom != "" {
		err := json.Unmarshal([]byte(idc.FirstServerRoom), &fsr)
		if err != nil {
			return upload.Return, err
		}
	}
	for k := range fsr {
		if fsr[k].Name == srfi.FirstServerRoomName {
			fsrExist = true
			srfi.FirstServerRoom = fsr[k].ID
		}
	}
	if !fsrExist {
		var br string
		if srfi.Content != "" {
			br = "<br />"
		}
		srfi.Content += br + fmt.Sprintf("在数据中心(%s)不存在一级机房(%s)", srfi.IDCName, srfi.FirstServerRoomName)
		return upload.Continue, nil
	}

	//机房校验
	srs, err := repo.GetServerRoomByName(srfi.Name)
	if err != nil && err != gorm.ErrRecordNotFound {
		return upload.Return, err
	}
	if err == gorm.ErrRecordNotFound || srs == nil {
		srfi.ID = 0
	} else {
		srfi.IDCID = srs.IDCID
		srfi.ID = srs.ID

		//如果通过导入修改数据中心，这块可以不做
		//
		//  if srfi.IDCID != idc.ID{
		// 	var br string
		// 	if srfi.Content != "" {
		// 		br = "<br />"
		// 	}
		// 	srfi.Content += br + fmt.Sprintf("已经存在的机房(%s)与录入的数据中心(%s)不一致", srfi.FirstServerRoomName, srfi.IDCName)
		// 	return
		// }
	}
	//使用表中数据查出来的数据中心ID，如果通过Excel导入修改数据中心要这样做
	srfi.IDCID = idc.ID

	//校验指定7*24小时电话格式是否正确
	phlist := strings.Split(srfi.SupportPhoneNumber, ",")
	for k := range phlist {
		if !chinaTelephoneNumberReq.MatchString(phlist[k]) {
			var br string
			if srfi.Content != "" {
				br = "<br />"
			}
			srfi.Content += br + fmt.Sprintf("7*24小时电话(%s)格式不正确", phlist[k])
		}
	}

	return upload.DO, nil
}

//ImportServerRoomPriview 导入预览
func ImportServerRoomPriview(log logger.Logger, repo model.Repo, reqData *upload.ImportReq) (map[string]interface{}, error) {
	ra, err := utils.ParseDataFromXLSX(upload.UploadDir + reqData.FileName)
	if err != nil {
		return nil, err
	}
	length := len(ra)

	var success []*ServerRoomForImport
	var failure []*ServerRoomForImport
	for i := 1; i < length; i++ {
		row := &ServerRoomForImport{}
		if len(ra[i]) < 9 {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "导入文件列长度不对（应为9列）"
			failure = append(failure, row)
			continue
		}

		row.Name = strings.TrimSpace(ra[i][0])
		row.IDCName = strings.TrimSpace(ra[i][1])
		row.FirstServerRoomName = strings.TrimSpace(ra[i][2])
		row.City = strings.TrimSpace(ra[i][3])
		row.Address = strings.TrimSpace(ra[i][4])
		row.ServerRoomManager = strings.TrimSpace(ra[i][5])
		row.VendorManager = strings.TrimSpace(ra[i][6])
		row.NetworkAssetManager = strings.TrimSpace(ra[i][7])
		row.SupportPhoneNumber = strings.TrimSpace(ra[i][8])

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

	var data []*ServerRoomForImport
	if len(failure) > 0 {
		data = failure
	} else {
		data = success
	}
	var result []*ServerRoomForImport
	for i := 0; i < len(data); i++ {
		if uint(i) >= reqData.Offset && uint(i) < (reqData.Offset+reqData.Limit) {
			result = append(result, data[i])
		}
	}
	if len(failure) > 0 {
		return map[string]interface{}{"status": "failure",
			"message":       "导入机房错误",
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

//ImportServerRoom 将导入机房放到数据库
func ImportServerRoom(log logger.Logger, repo model.Repo, reqData *upload.ImportReq) error {
	ra, err := utils.ParseDataFromXLSX(upload.UploadDir + reqData.FileName)
	if err != nil {
		return err
	}
	length := len(ra)
	for i := 1; i < length; i++ {
		row := &ServerRoomForImport{}
		if len(ra[i]) < 9 {
			continue
		}

		row.Name = strings.TrimSpace(ra[i][0])
		row.IDCName = strings.TrimSpace(ra[i][1])
		row.FirstServerRoomName = strings.TrimSpace(ra[i][2])
		row.City = strings.TrimSpace(ra[i][3])
		row.Address = strings.TrimSpace(ra[i][4])
		row.ServerRoomManager = strings.TrimSpace(ra[i][5])
		row.VendorManager = strings.TrimSpace(ra[i][6])
		row.NetworkAssetManager = strings.TrimSpace(ra[i][7])
		row.SupportPhoneNumber = strings.TrimSpace(ra[i][8])

		//必填项校验
		row.checkLength()
		//机房和数据中心校验
		isSave, err := row.validate(repo)
		if err != nil {
			return err
		}

		//不能获取IDC
		if isSave == upload.Continue {
			continue
		}

		sr := &model.ServerRoom{
			IDCID:               row.IDCID,
			Name:                row.Name,
			FirstServerRoom:     row.FirstServerRoom,
			City:                row.City,
			Address:             row.Address,
			ServerRoomManager:   row.ServerRoomManager,
			VendorManager:       row.VendorManager,
			NetworkAssetManager: row.NetworkAssetManager,
			SupportPhoneNumber:  row.SupportPhoneNumber,
			Creator:             row.Creator,
		}
		if row.ID == 0 {
			sr.Status = model.RoomStatUnderConstruction
			_, err := repo.SaveServerRoom(sr)
			if err != nil {
				return err
			}
		}
		if row.ID > 0 {
			sr.ID = uint(row.ID)
			_, err := repo.UpdateServerRoom([]*model.ServerRoom{sr})
			if err != nil {
				return err
			}
		}
	}
	defer os.Remove(upload.UploadDir + reqData.FileName)
	return nil
}
