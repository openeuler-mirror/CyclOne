package service

import (
	"net/http"

	"fmt"

	"encoding/json"

	"reflect"

	"errors"

	"github.com/jinzhu/gorm"
	"github.com/voidint/binding"
	"github.com/voidint/page"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils/times"
)

type (
	// IDCReq 保存/修改数据中心Body参数结构
	IDCReq struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
		//用途
		Usage string `json:"usage"`
		//一级机房,可以有多个
		FirstServerRoom FirstServerRooms `json:"first_server_room"`
		//供应商
		Vendor string `json:"vendor"`
		//建设中under_construction，已验收accepted，已投产production，已裁撤abolished
		Status string `json:"status"`
		// 用户登录名
		LoginName string `json:"-"`
	}

	// FirstServerRooms 一级机房数组
	FirstServerRooms []string

	// FirstServerRoom 一级机房
	FirstServerRoom struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	// IDCUpdateReq 批量更新状态请求体
	IDCUpdateReq struct {
		//建设中under_construction，已验收accepted，已投产production，已裁撤abolished
		Status string `json:"status"`
		//指定修改的ID列表
		IDs []uint `json:"ids"`
	}

	// IDCResp 数据中心返回结构
	IDCResp struct {
		ID              uint              `json:"id"`
		Name            string            `json:"name"`
		Usage           string            `json:"usage"`
		FirstServerRoom []FirstServerRoom `json:"first_server_room"`
		Vendor          string            `json:"vendor"`
		Status          string            `json:"status"`
		CreatedAt       times.ISOTime     `json:"created_at"`
		UpdatedAt       times.ISOTime     `json:"updated_at"`
	}

	// IDCPageReq 分页查询条件
	IDCPageReq struct {
		Name string `json:"name"`
		//用途：production-生产; disaster_recovery-容灾; pre_production-准生产; testing-测试;',
		Usage string `json:"usage"`
		//一级机房
		FirstServerRoom string `json:"first_server_room"`
		//供应商
		Vendor string `json:"vendor"`
		//建设中under_construction，已验收accepted，已投产production，已裁撤abolished
		Status   string `json:"status"`
		Page     int64  `json:"page"`
		PageSize int64  `json:"page_size"`
	}
)

// FieldMap 请求字段映射
func (reqData *IDCReq) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&reqData.ID:              "id",
		&reqData.Name:            "name",
		&reqData.Usage:           "usage",
		&reqData.FirstServerRoom: "first_server_room",
		&reqData.Vendor:          "vendor",
	}
}

// FieldMap 请求字段映射
func (reqData *IDCUpdateReq) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&reqData.Status: "status",
		&reqData.IDs:    "ids",
	}
}

// FieldMap 请求字段映射
func (reqData *IDCPageReq) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&reqData.Name:            "name",
		&reqData.Usage:           "usage",
		&reqData.FirstServerRoom: "first_server_room",
		&reqData.Vendor:          "vendor",
		&reqData.Status:          "status",
		&reqData.Page:            "page",
		&reqData.PageSize:        "page_size",
	}
}

//ToJSON FirstServerRooms转Json
func (f *FirstServerRooms) ToJSON() []byte {
	rooms := make([]*FirstServerRoom, 0)
	for i, room := range *f {
		rooms = append(rooms, &FirstServerRoom{
			ID:   i,
			Name: room,
		})
	}
	b, _ := json.Marshal(rooms)
	return b
}

// Validate 校验入参
func (reqData *IDCReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())
	if reqData.Name == "" {
		errs.Add([]string{"name"}, binding.RequiredError, "名称不能为空")
		return errs
	}

	if idc, _ := repo.GetIDCByName(reqData.Name); idc != nil && idc.ID != reqData.ID {
		errs.Add([]string{"name"}, binding.RequiredError, "名称已存在,不能重复")
		return errs
	}
	if reqData.Usage == "" {
		errs.Add([]string{"usage"}, binding.RequiredError, "用途不能为空")
		return errs
	}

	usageEnums := []string{model.IDCUsageProduction, model.IDCUsageDisasterRecovery,
		model.IDCUsagePreProduction, model.IDCUsageTesting}
	if !ValidateEnum(reqData.Usage, usageEnums) {
		errs.Add([]string{"usage"}, binding.RequiredError, fmt.Sprintf("用途的取值必须为%v", usageEnums))
		return errs
	}
	if reqData.FirstServerRoom == nil {
		errs.Add([]string{"first_server_room"}, binding.RequiredError, "一级机房不能为空")
		return errs
	}
	if reqData.Vendor == "" {
		errs.Add([]string{"vendor"}, binding.RequiredError, "供应商不能为空")
		return errs
	}
	return errs
}

// Validate 校验入参
func (reqData *IDCUpdateReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())
	statusEnums := []string{model.IDCStatUnderConstruction,
		model.IDCStatAccepted, model.IDCStatProduction /*, model.IDCStatAbolished*/}
	if !ValidateEnum(reqData.Status, statusEnums) {
		errs.Add([]string{"status"}, binding.RequiredError, fmt.Sprintf("数据中心状态取值必须是%v", statusEnums))
		return errs
	}

	if len(reqData.IDs) == 0 {
		errs.Add([]string{"ids"}, binding.RequiredError, "数据中心ID列表为空")
		return errs
	}

	for _, idcID := range reqData.IDs {
		_, err := repo.GetIDCByID(idcID)
		if err == gorm.ErrRecordNotFound {
			errs.Add([]string{"ids"}, binding.RequiredError, fmt.Sprintf("数据中心ID：%d 未找到", idcID))
			return errs
		}
	}
	return errs
}

// ValidateEnum 校验枚举字段的值
func ValidateEnum(val string, valList []string) bool {
	for _, v := range valList {
		if val == v {
			return true
		}
	}
	return false
}

// SaveIDC 保存（新增/修改）
func SaveIDC(log logger.Logger, repo model.Repo, req *IDCReq) (*model.IDC, error) {
	idc := &model.IDC{
		Model:           gorm.Model{ID: uint(req.ID)},
		Name:            req.Name,
		Usage:           req.Usage,
		FirstServerRoom: string(req.FirstServerRoom.ToJSON()),
		Vendor:          req.Vendor,
		Creator:         req.LoginName,
	}
	if req.ID == 0 {
		//add a new record
		idc.Status = model.IDCStatUnderConstruction //set default
		_, mod, err := repo.AddIDC(idc)
		if err != nil {
			return nil, err
		}
		return mod, nil
	}

	//update a new record
	_, err := repo.UpdateIDC(idc)
	if err != nil {
		return nil, err
	}
	return idc, nil
}

// RemoveIDCByID 删除指定ID的数据中心
func RemoveIDCByID(log logger.Logger, repo model.Repo, idcID uint) error {
	//有被关联机房的不能被删除
	serverRoomCountTotal, err := repo.CountServerRooms(&model.ServerRoomCond{
		IDCID: []uint{idcID},
	})
	if serverRoomCountTotal > 0 {
		return errors.New("[业务校验]有机房关联该数据中心，不允许删除")
	}

	_, err = repo.GetIDCByID(idcID)
	if err != nil {
		return err
	}
	_, err = repo.RemoveIDCByID(idcID)
	return err
}

// UpdateIDCStatus 批量更新数据中心状态
func UpdateIDCStatus(log logger.Logger, repo model.Repo, req *IDCUpdateReq) error {
	_, err := repo.UpdateIDCStatus(req.Status, req.IDs...)
	return err
}

// GetIDCByID 查询指定ID的数据中心信息详情
func GetIDCByID(log logger.Logger, repo model.Repo, idcID uint) (*IDCResp, error) {
	mod, err := repo.GetIDCByID(idcID)
	if err != nil {
		return nil, err
	}
	resp := convert2IDCResp(mod)
	return resp, nil
}

// GetIDCPage 查询数据中心分页列表
func GetIDCPage(log logger.Logger, repo model.Repo, reqData *IDCPageReq) (pg *page.Page, err error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 100 {
		reqData.PageSize = 20
	}
	if reqData.Page < 0 {
		reqData.Page = 0
	}

	cond := model.IDC{
		Name:            reqData.Name,
		Usage:           reqData.Usage,
		FirstServerRoom: reqData.FirstServerRoom,
		Vendor:          reqData.Vendor,
		Status:          reqData.Status,
	}

	totalRecords, err := repo.CountIDCs(&cond)
	if err != nil {
		return nil, err
	}

	pager := page.NewPager(reflect.TypeOf(&IDCResp{}), reqData.Page, reqData.PageSize, totalRecords)
	items, err := repo.GetIDCs(&cond, model.OneOrderBy("id", model.DESC), pager.BuildLimiter())
	if err != nil {
		return nil, err
	}
	for i := range items {
		item := convert2IDCResp(items[i])
		pager.AddRecords(item)
	}
	return pager.BuildPage(), nil

}

// convert2IDCResp 按照返回的结构定义转换
func convert2IDCResp(mod *model.IDC) *IDCResp {
	resp := &IDCResp{
		ID:        mod.ID,
		Name:      mod.Name,
		Usage:     mod.Usage,
		Vendor:    mod.Vendor,
		Status:    mod.Status,
		CreatedAt: times.ISOTime(mod.CreatedAt),
		UpdatedAt: times.ISOTime(mod.UpdatedAt),
	}
	_ = json.Unmarshal([]byte(mod.FirstServerRoom), &resp.FirstServerRoom)
	return resp
}
