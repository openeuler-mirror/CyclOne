package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/voidint/binding"

	"reflect"

	"os"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils"
	strings2 "idcos.io/cloudboot/utils/strings"
	"idcos.io/cloudboot/utils/times"
	"idcos.io/cloudboot/utils/upload"
)

//SaveStoreRoomReq 保存库房请求参数
type SaveStoreRoomReq struct {
	//(required) 所属数据中心ID
	IDCID uint `json:"idc_id"`
	//库房ID。若id=0，则新增。若id>0，则修改。
	ID uint `json:"id"`
	//required) 库房名称
	Name string `json:"name"`
	//(required)一级机房
	FirstServerRoom string `json:"first_server_room"`
	//(required) 所属城市
	City string `json:"city"`
	//(required) 地址
	Address string `json:"address"`
	//(required) 库房管理单元负责人
	StoreRoomManager string `json:"store_room_manager"`
	//(required) 供应商负责人
	VendorManager string `json:"vendor_manager"`
	// 用户登录名
	LoginName string `json:"-"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *SaveStoreRoomReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDCID:            "idc_id",
		&reqData.ID:               "id",
		&reqData.Name:             "name",
		&reqData.FirstServerRoom:  "first_server_room",
		&reqData.City:             "city",
		&reqData.Address:          "address",
		&reqData.StoreRoomManager: "store_room_manager",
		&reqData.VendorManager:    "vendor_manager",
	}
}

// Validate 结构体数据校验
func (reqData *SaveStoreRoomReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	//必须的基本数据不能为空
	if errs = reqData.baseValidate(req, errs); errs != nil {
		return errs
	}

	//更新库房信息，校验指定ID的库房是否存在
	if reqData.ID > 0 {
		if errs = reqData.storeRoomValidate(req, errs); errs != nil {
			return errs
		}
	}

	//校验IDC数据
	if errs = reqData.checkIDCValidate(req, errs); errs != nil {
		return errs
	}

	return errs
}

func (reqData *SaveStoreRoomReq) checkIDCValidate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())

	idc, err := repo.GetIDCByID(reqData.IDCID)
	if err != nil {
		errs.Add([]string{"idc"}, binding.RequiredError, fmt.Sprintf("获取数据中心id(%d)出现错误: %s", reqData.IDCID, err.Error()))
		return errs
	}
	reqData.IDCID = idc.ID

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
		if reqData.FirstServerRoom == fsr[k].Name {
			isExist = true
			break
		}
	}
	if !isExist {
		errs.Add([]string{"idc"}, binding.RequiredError, fmt.Sprintf("获取数据中心(一级机房：%s)不存在", reqData.FirstServerRoom))
		return errs
	}
	return errs

}

func (reqData *SaveStoreRoomReq) storeRoomValidate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())

	var cond model.StoreRoomCond

	cond.ID = []uint{reqData.ID}

	items, err := repo.GetStoreRooms(&cond, nil, nil)
	if err != nil {
		errs.Add([]string{"store_room"}, binding.RequiredError, fmt.Sprintf("更新库房ID %d不存在", reqData.ID))
		return errs
	}
	if len(items) < 1 {
		errs.Add([]string{"store_room"}, binding.RequiredError, fmt.Sprintf("更新库房ID %d不存在", reqData.ID))
		return errs
	}

	//校验更新库房名是否已经存在
	if reqData.Name != "" {
		item, _ := repo.GetStoreRooms(&model.StoreRoomCond{Name: reqData.Name}, nil, nil)
		if item != nil && len(item) > 0 && item[0].ID != reqData.ID {
			errs.Add([]string{"store_room"}, binding.RequiredError, fmt.Sprintf("更新库房 指定的库房名%s已经存在", reqData.Name))
			return errs
		}
	}
	return errs
}

//baseValidate 必要参数不能为空
func (reqData *SaveStoreRoomReq) baseValidate(req *http.Request, errs binding.Errors) binding.Errors {
	if reqData.IDCID == 0 {
		errs.Add([]string{"idc_id"}, binding.RequiredError, "数据中心ID不能为空")
		return errs
	}

	if reqData.Name == "" {
		errs.Add([]string{"name"}, binding.RequiredError, "库房名称不能为空")
		return errs
	}

	if reqData.City == "" {
		errs.Add([]string{"city"}, binding.RequiredError, "库房所属城市不能为空")
		return errs
	}

	if reqData.Address == "" {
		errs.Add([]string{"address"}, binding.RequiredError, "库房地址不能为空")
		return errs
	}

	if reqData.StoreRoomManager == "" {
		errs.Add([]string{"store_room_manager"}, binding.RequiredError, "库房管理单元负责人不能为空")
		return errs
	}

	if reqData.VendorManager == "" {
		errs.Add([]string{"vendor_manager"}, binding.RequiredError, "供应商负责人不能为空")
		return errs
	}
	return errs
}

//SaveStoreRoom 保存库房
func SaveStoreRoom(log logger.Logger, repo model.Repo, reqData *SaveStoreRoomReq) error {
	sr := model.StoreRoom{
		IDCID:            reqData.IDCID,
		Name:             reqData.Name,
		FirstServerRoom:  reqData.FirstServerRoom,
		City:             reqData.City,
		Address:          reqData.Address,
		StoreRoomManager: reqData.StoreRoomManager,
		VendorManager:    reqData.VendorManager,
		Creator:          reqData.LoginName,
	}
	if reqData.ID == 0 {
		sr.Status = model.RoomStatAccepted
		_, err := repo.SaveStoreRoom(&sr)
		if err != nil {
			return err
		}
	}
	if reqData.ID > 0 {
		sr.ID = uint(reqData.ID)
		_, err := repo.UpdateStoreRoom([]*model.StoreRoom{&sr})
		if err != nil {
			return err
		}
	}
	reqData.ID = sr.ID
	return nil
}

type StoreRoomImportReq struct {
	SaveStoreRoomReq
	IDCName string `json:"idc_name"`
	//FirstServerRoom string `json:"first_server_room_name"`
	Content string `json:"content"`
}

//checkLength 对导入文件中的数据做基本验证
func (nafi *StoreRoomImportReq) checkLength() {
	leg := len(nafi.IDCName)
	if leg == 0 || leg > 255 {
		var br string
		if nafi.Content != "" {
			br = "<br />"
		}
		nafi.Content += br + fmt.Sprintf("必填项校验:数据中心名称长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(nafi.FirstServerRoom)
	if leg == 0 || leg > 255 {
		var br string
		if nafi.Content != "" {
			br = "<br />"
		}
		nafi.Content += br + fmt.Sprintf("必填项校验:一级名称长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(nafi.Name)
	if leg == 0 || leg > 255 {
		var br string
		if nafi.Content != "" {
			br = "<br />"
		}
		nafi.Content += br + fmt.Sprintf("必填项校验:库房名称长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(nafi.City)
	if leg == 0 || leg > 255 {
		var br string
		if nafi.Content != "" {
			br = "<br />"
		}
		nafi.Content += br + fmt.Sprintf("必填项校验:城市长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(nafi.Address)
	if leg == 0 || leg > 255 {
		var br string
		if nafi.Content != "" {
			br = "<br />"
		}
		nafi.Content += br + fmt.Sprintf("必填项校验:地址长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(nafi.StoreRoomManager)
	if leg == 0 || leg > 255 {
		var br string
		if nafi.Content != "" {
			br = "<br />"
		}
		nafi.Content += br + fmt.Sprintf("必填项校验:负责人长度为(%d)(不能为空，不能大于255)", leg)
	}
}

//validate 对导入文件中的数据做基本验证
func (nafi *StoreRoomImportReq) validate(log logger.Logger, repo model.Repo) (int, error) {
	//机房校验
	idc, err := repo.GetIDCByName(nafi.IDCName)
	if err != nil && err != gorm.ErrRecordNotFound {
		return upload.Return, err
	}
	if err == gorm.ErrRecordNotFound || idc == nil {
		var br string
		if nafi.Content != "" {
			br = "<br />"
		}
		nafi.Content += br + fmt.Sprintf("数据中心(%s)不存在", nafi.FirstServerRoom)
		return upload.Continue, nil
	}
	nafi.IDCID = idc.ID

	sr, err := repo.GetStoreRoomByName(nafi.Name)
	if sr != nil {
		nafi.ID = sr.ID
	}

	return upload.DO, nil
}

//ImportStoreRoomPriview 导入预览
func ImportStoreRoomPriview(log logger.Logger, repo model.Repo, reqData *upload.ImportReq) (map[string]interface{}, error) {
	ra, err := utils.ParseDataFromXLSX(upload.UploadDir + reqData.FileName)
	if err != nil {
		return nil, err
	}
	length := len(ra)

	var success []*StoreRoomImportReq
	var failure []*StoreRoomImportReq
	for i := 1; i < length; i++ {
		row := &StoreRoomImportReq{}
		if len(ra[i]) < 7 {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "导入文件列长度不对（应为7列）"
			failure = append(failure, row)
			continue
		}
		row.IDCName = strings.TrimSpace(ra[i][0])
		row.FirstServerRoom = strings.TrimSpace(ra[i][1])
		row.Name = strings.TrimSpace(ra[i][2])
		row.City = strings.TrimSpace(ra[i][3])
		row.Address = strings.TrimSpace(ra[i][4])
		row.StoreRoomManager = strings.TrimSpace(ra[i][5])
		row.VendorManager = strings.TrimSpace(ra[i][6])

		//必填项校验
		row.checkLength()
		//机房和网络区域校验
		_, err := row.validate(log, repo)
		if err != nil {
			return nil, err
		}

		if row.Content != "" {
			failure = append(failure, row)
		} else {
			success = append(success, row)
		}
	}

	var data []*StoreRoomImportReq
	if len(failure) > 0 {
		data = failure
	} else {
		data = success
	}
	var result []*StoreRoomImportReq
	for i := 0; i < len(data); i++ {
		if uint(i) >= reqData.Offset && uint(i) < (reqData.Offset+reqData.Limit) {
			result = append(result, data[i])
		}
	}
	if len(failure) > 0 {
		return map[string]interface{}{"status": "failure",
			"message":       "导入库房错误",
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

//ImportStoreRoom 将导入网络区域放到数据库
func ImportStoreRoom(log logger.Logger, repo model.Repo, reqData *upload.ImportReq) error {
	ra, err := utils.ParseDataFromXLSX(upload.UploadDir + reqData.FileName)
	if err != nil {
		return err
	}
	length := len(ra)
	for i := 1; i < length; i++ {
		row := &StoreRoomImportReq{}
		if len(ra[i]) < 7 {
			continue
		}

		row.IDCName = strings.TrimSpace(ra[i][0])
		row.FirstServerRoom = strings.TrimSpace(ra[i][1])
		row.Name = strings.TrimSpace(ra[i][2])
		row.City = strings.TrimSpace(ra[i][3])
		row.Address = strings.TrimSpace(ra[i][4])
		row.StoreRoomManager = strings.TrimSpace(ra[i][5])
		row.VendorManager = strings.TrimSpace(ra[i][6])

		//必填项校验
		row.checkLength()
		//机房和网络区域校验
		isSave, err := row.validate(log, repo)
		if err != nil {
			return err
		}

		//不能获取机房信息
		if isSave == upload.Continue {
			continue
		}
		row.SaveStoreRoomReq.LoginName = reqData.UserName

		if err = SaveStoreRoom(log, repo, &row.SaveStoreRoomReq); err != nil {
			return err
		}
	}
	defer os.Remove(upload.UploadDir + reqData.FileName)
	return nil
}

//RemoveStoreRoomValidate 删除操作校验
func RemoveStoreRoomValidate(log logger.Logger, repo model.Repo, storeRoomID uint) string {
	//统计库房中的机架(柜数)
	count, _ := repo.CountVirtualCabinets(&model.VirtualCabinet{StoreRoomID: storeRoomID})
	if count > 0 {
		return fmt.Sprintf("库房下面存在(%d)个虚拟货架,须先删除", count)
	}
	return ""
}

//RemoveStoreRoomByID 删除指定ID的库房
func RemoveStoreRoomByID(log logger.Logger, repo model.Repo, id uint) error {
	_, err := repo.RemoveStoreRoomByID(id)
	return err
}

//GetStoreRoomPageReq 获取库房分页请求参数
type GetStoreRoomPageReq struct {
	// 所属数据中心ID
	IDCID string `json:"idc_id"`
	//  库房名称(支持模糊匹配)
	Name string `json:"name"`
	// 一级机房ID
	FirstServerRoom string `json:"first_server_room"`
	// 城市(支持模糊匹配)
	City string `json:"city"`
	// 地址(支持模糊匹配)
	Address string `json:"address"`
	// 库房管理单元负责人(支持模糊匹配)
	StoreRoomManager string `json:"store_room_manager"`
	// 供应商负责人(支持模糊匹配)
	VendorManager string `json:"vendor_manager"`
	// 分页页号
	Page int64 `json:"page"`
	// 分页大小。默认值:10。阈值: 100。
	PageSize int64 `json:"page_size"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *GetStoreRoomPageReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDCID: "idc_id",
		//&reqData.Status:              "status",
		&reqData.Name:             "name",
		&reqData.FirstServerRoom:  "first_server_room",
		&reqData.City:             "city",
		&reqData.Address:          "address",
		&reqData.StoreRoomManager: "store_room_manager",
		&reqData.VendorManager:    "vendor_manager",
		&reqData.Page:             "page",
		&reqData.PageSize:         "page_size",
	}
}

//StoreRoom 库房分页查询信息
type StoreRoom struct {
	//库房ID
	ID uint `json:"id"`
	//创建时间
	CreatedAt times.ISOTime `json:"created_at"`
	//更新时间
	UpdatedAt times.ISOTime `json:"updated_at"`
	//库房名
	Name string `json:"name"`
	//一级机房ID
	FirstServerRoom IDCFirstServerRoom `json:"first_server_room"`
	//机架数
	CabinetCount int64 `json:"cabinet_count"`
	//所在城市
	City string `json:"city"`
	//地址
	Address string `json:"address"`
	//库房管理人
	StoreRoomManager string `json:"store_room_manager"`
	//供应商负责人
	VendorManager string `json:"vendor_manager"`
	//库房状态
	Status string `json:"status"`
	//创建人
	Creator string `json:"creator"`
	//数据中心
	IDC IDCForServerRoomPage `json:"idc"`
}

//GetStoreRoomWithPage 获取库房分页
func GetStoreRoomWithPage(log logger.Logger, repo model.Repo, reqData *GetStoreRoomPageReq) (*page.Page, error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 100 {
		reqData.PageSize = 20
	}
	if reqData.Page < 0 {
		reqData.Page = 0
	}

	cond := model.StoreRoomCond{
		IDCID:            strings2.Multi2UintSlice(reqData.IDCID),
		Name:             reqData.Name,
		FirstServerRoom:  reqData.FirstServerRoom,
		City:             reqData.City,
		Address:          reqData.Address,
		StoreRoomManager: reqData.StoreRoomManager,
		VendorManager:    reqData.VendorManager,
		//Status:           reqData.Status,
	}

	totalRecords, err := repo.CountStoreRooms(&cond)
	if err != nil {
		return nil, err
	}

	pager := page.NewPager(reflect.TypeOf(&StoreRoom{}), reqData.Page, reqData.PageSize, totalRecords)
	items, err := repo.GetStoreRooms(&cond, model.OneOrderBy("id", model.DESC), pager.BuildLimiter())
	if err != nil {
		return nil, err
	}
	for i := range items {
		item, err := convert2StoreRoomResult(log, repo, items[i])
		if err != nil {
			return nil, err
		}
		if item != nil {
			pager.AddRecords(item)
		}
	}
	return pager.BuildPage(), nil
}

func convert2StoreRoomResult(log logger.Logger, repo model.Repo, srcp *model.StoreRoom) (*StoreRoom, error) {
	if srcp == nil {
		return nil, nil
	}

	idc, err := repo.GetIDCByID(srcp.IDCID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	result := StoreRoom{
		ID:               srcp.ID,
		CreatedAt:        times.ISOTime(srcp.CreatedAt),
		UpdatedAt:        times.ISOTime(srcp.UpdatedAt),
		Name:             srcp.Name,
		City:             srcp.City,
		Address:          srcp.Address,
		StoreRoomManager: srcp.StoreRoomManager,
		VendorManager:    srcp.VendorManager,
		Status:           srcp.Status,
		Creator:          srcp.Creator,
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
			if fsr[k].Name == srcp.FirstServerRoom {
				result.FirstServerRoom = fsr[k]
			}
		}
	}

	//统计库房中的虚拟货架数量
	count, err := repo.CountVirtualCabinets(&model.VirtualCabinet{StoreRoomID: result.ID})
	if err != nil {
		return nil, err
	}
	result.CabinetCount = count

	return &result, nil
}

//GetStoreRoomByID 获取指定ID的库房的详细信息
func GetStoreRoomByID(log logger.Logger, repo model.Repo, id uint) (*StoreRoom, error) {
	items, err := repo.GetStoreRoomByID(id)
	if err != nil {
		return nil, err
	}

	item, err := convert2StoreRoomResult(log, repo, items)
	if err != nil {
		return nil, err
	}
	return item, nil
}

////////////////////////////////////////////////// 虚拟货架  /////////////////////////////////////

//SaveVirtualCabinetReq 保存虚拟货架请求参数
type SaveVirtualCabinetReq struct {
	// 库房ID
	StoreRoomID uint `json:"store_room_id"`
	//库房ID。若id=0，则新增。若id>0，则修改。
	ID uint `json:"id"`
	//(required) 编号
	Number string `json:"number"`
	//备注
	Remark string
	// 用户登录名
	LoginName string `json:"-"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *SaveVirtualCabinetReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.StoreRoomID: "store_room_id",
		//&reqData.ID:          "id",
		&reqData.Number: "name",
		&reqData.Remark: "remark",
	}
}

//SaveVirtualCabinet 新增虚拟货架
func SaveVirtualCabinet(log logger.Logger, repo model.Repo, reqData *SaveVirtualCabinetReq) error {
	sr := model.VirtualCabinet{
		StoreRoomID: reqData.StoreRoomID,
		Number:      reqData.Number,
		Remark:      reqData.Remark,
		Status:      model.CabinetStatEnabled,
		Creator:     reqData.LoginName,
	}
	_, err := repo.SaveVirtualCabinet(&sr)
	if err != nil {
		return err
	}
	reqData.ID = sr.ID
	return nil
}

// RemoveVirtualCabinetValidate 删除之前的数据校验
func RemoveVirtualCabinetValidate(log logger.Logger, repo model.Repo, vCabinetID uint) (bool, error) {
	devs, err := repo.GetDevices(&model.Device{VCabinetID: vCabinetID}, nil, nil)
	if err != nil {
		log.Errorf("count devices of virtual cabinet:%d fail", vCabinetID)
		return false, err
	}
	if len(devs) != 0 {
		return false, fmt.Errorf("虚拟机位(id:%d)上关联有物理机(SN:%s),须先搬离再删除", vCabinetID, devs[0].SN)
	}
	return true, nil
}

//RemoveVirtualCabinetByID 删除指定ID的虚拟货架
func RemoveVirtualCabinetByID(log logger.Logger, repo model.Repo, id uint) error {
	_, err := repo.RemoveVirtualCabinetByID(id)
	return err
}

type GetVirtualCabinetPageReq struct {
	Page int64 `json:"page"`
	// 分页大小。默认值:10。阈值 100
	PageSize int64 `json:"page_size"`
	// StoreRoomID
	StoreRoomID uint `json:"store_room_id"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *GetVirtualCabinetPageReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.StoreRoomID: "store_room_id",
		&reqData.Page:        "page",
		&reqData.PageSize:    "page_size",
	}
}

//VirtualCabinet 虚拟货架信息
type VirtualCabinet struct {
	//ID
	ID uint `json:"id"`
	//创建时间
	CreatedAt times.ISOTime `json:"created_at"`
	//更新时间
	UpdatedAt   times.ISOTime `json:"updated_at"`
	StoreRoomID uint          `json:"store_room_id"`
	Number      string        `json:"number"`
	Remark      string        `json:"remark"`
	Status      string        `json:"status"`
	Creator     string        `json:"creator"`
}

//GetVirtualCabinetWithPage 获取虚拟货架分页
func GetVirtualCabinetWithPage(log logger.Logger, repo model.Repo, reqData *GetVirtualCabinetPageReq) (*page.Page, error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 100 {
		reqData.PageSize = 20
	}
	if reqData.Page < 0 {
		reqData.Page = 0
	}

	cond := model.VirtualCabinet{
		StoreRoomID: reqData.StoreRoomID,
	}

	totalRecords, err := repo.CountVirtualCabinets(&cond)
	if err != nil {
		return nil, err
	}

	pager := page.NewPager(reflect.TypeOf(&VirtualCabinet{}), reqData.Page, reqData.PageSize, totalRecords)
	items, err := repo.GetVirtualCabinets(&cond, model.OneOrderBy("id", model.DESC), pager.BuildLimiter())
	if err != nil {
		return nil, err
	}
	for _, i := range items {
		result := VirtualCabinet{
			ID:          i.ID,
			CreatedAt:   times.ISOTime(i.CreatedAt),
			UpdatedAt:   times.ISOTime(i.UpdatedAt),
			StoreRoomID: i.StoreRoomID,
			Number:      i.Number,
			Status:      i.Status,
			Remark:      i.Remark,
			Creator:     i.Creator,
		}
		pager.AddRecords(&result)
	}
	return pager.BuildPage(), nil
}
