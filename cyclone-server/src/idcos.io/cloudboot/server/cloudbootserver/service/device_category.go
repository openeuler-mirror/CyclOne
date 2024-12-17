package service

import (
	"fmt"
	"net/http"

	"github.com/voidint/binding"

	"reflect"

	"github.com/voidint/page"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils/times"
)

//SaveDeviceCategoryReq 保存设备类型请求参数
type SaveDeviceCategoryReq struct {
	DeviceCategoryBase
	ID uint `json:"id"`
	// 用户登录名
	LoginName string `json:"-"`
}

// 基本字段
type DeviceCategoryBase struct {
	//Catetory 设备类型
	Category 						string `json:"category"`
	// 硬件配置
	Hardware 						string `json:"hardware"`
	// 处理器生产商，如：Intel(R) Corporation\HiSilicon
	CentralProcessorManufacturer	string	`json:"central_processor_manufacture"`
	// 处理器架构,如：x86_64\aarch64
	CentralProcessorArch			string	`json:"central_processor_arch"`
	//功率
	Power 							string `json:"power"`
	//设备所占用的U数 机柜参数Unit 1U = 44.45mm
	Unit 							uint	`json:"unit"`
	//是否金融信创生态产品: yes or no
	IsFITIEcoProduct				string	`json:"is_fiti_eco_product"`
	// Remark 备注
	Remark 							string `json:"remark"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *SaveDeviceCategoryReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.ID:       							"id",
		&reqData.Category: 							"category",
		&reqData.Hardware: 							"hardware",
		&reqData.CentralProcessorManufacturer: 		"central_processor_manufacture",
		&reqData.CentralProcessorArch: 				"central_processor_arch",
		&reqData.Power:    							"power",
		&reqData.Unit:    							"unit",
		&reqData.IsFITIEcoProduct:    				"is_fiti_eco_product",
		&reqData.Remark:   							"remark",
	}
}

// Validate 结构体数据校验
func (reqData *SaveDeviceCategoryReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())
	//必须的基本数据不能为空
	if errs = reqData.baseValidate(req, errs); errs != nil {
		return errs
	}

	//更新设备类型信息，校验指定ID的设备类型是否存在
	if reqData.ID > 0 {
		if _, err := repo.GetDeviceCategoryByID(reqData.ID); errs != nil {
			errs.Add([]string{"id"}, binding.RequiredError, fmt.Sprintf("查询设备类型id(%d)出现错误: %s", reqData.ID, err.Error()))
			return errs
		}
	} else {
		if c, _ := repo.GetDeviceCategoryByName(reqData.Category); c != nil {
			errs.Add([]string{"id"}, binding.RequiredError, fmt.Sprintf("设备类型(%s)已存在", reqData.Category))
			return errs
		}
	}
	// 校验
	if reqData.IsFITIEcoProduct != model.NO && reqData.IsFITIEcoProduct != model.YES {
		errs.Add([]string{"is_fiti_eco_product"}, binding.RequiredError, fmt.Sprintf("是否金融信创生态产品(%s)不合法", reqData.IsFITIEcoProduct))
		return errs
	}
	return nil
}

//baseValidate 必要参数不能为空
func (reqData *SaveDeviceCategoryReq) baseValidate(req *http.Request, errs binding.Errors) binding.Errors {
	if reqData.Category == "" {
		errs.Add([]string{"category"}, binding.RequiredError, "设备类型不能为空")
		return errs
	}
	if reqData.Hardware == "" {
		errs.Add([]string{"hardware"}, binding.RequiredError, "硬件配置不能为空")
		return errs
	}
	return errs
}

type DelDeviceCategoryReq struct {
	IDs []uint `json:"ids"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *DelDeviceCategoryReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDs: "ids",
	}
}

// Validate 结构体数据校验
func (reqData *DelDeviceCategoryReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())
	for _, id := range reqData.IDs {
		if _, err := repo.GetDeviceCategoryByID(id); err != nil {
			errs.Add([]string{"id"}, binding.RequiredError, fmt.Sprintf("设备类型id(%d)不存在", id))
			return errs
		}
	}
	return nil
}

//SaveDeviceCategory 保存设备类型
func SaveDeviceCategory(log logger.Logger, repo model.Repo, reqData *SaveDeviceCategoryReq) error {
	sr := model.DeviceCategory{
		Category: 						reqData.Category,
		Hardware: 						reqData.Hardware,
		CentralProcessorManufacturer: 	reqData.CentralProcessorManufacturer,
		CentralProcessorArch: 			reqData.CentralProcessorArch,
		Power:    						reqData.Power,
		Unit:	  						reqData.Unit,
		IsFITIEcoProduct: 				reqData.IsFITIEcoProduct,
		Remark:   						reqData.Remark,
		Creator:  						reqData.LoginName,
	}
	sr.Model.ID = reqData.ID

	_, err := repo.SaveDeviceCategory(&sr)
	if err != nil {
		return err
	}

	reqData.ID = sr.Model.ID
	return err
}

//RemoveDeviceCategorys 删除指定ID的设备类型
func RemoveDeviceCategorys(log logger.Logger, repo model.Repo, reqData *DelDeviceCategoryReq) (affected int64, err error) {
	for _, id := range reqData.IDs {
		_, err := repo.RemoveDeviceCategoryByID(id)
		if err != nil {
			log.Errorf("delete order(id=%d) fail,err:%v", id, err)
			return affected, err
		}

		affected++
	}
	return affected, err
}

//GetDeviceCategoryPageReq 获取设备类型分页请求参数
type GetDeviceCategoryPageReq struct {
	DeviceCategoryBase
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *GetDeviceCategoryPageReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.Category: 							"category",
		&reqData.Hardware: 							"hardware",
		&reqData.CentralProcessorManufacturer: 		"central_processor_manufacture",
		&reqData.CentralProcessorArch: 				"central_processor_arch",
		&reqData.Power:    							"power",
		&reqData.Unit:     							"unit",
		&reqData.IsFITIEcoProduct:    				"is_fiti_eco_product",
		&reqData.Remark:   							"remark",
		&reqData.Page:     							"page",
		&reqData.PageSize: 							"page_size",	
	}
}

//DeviceCategoryResp 设备类型分页查询信息
type DeviceCategoryResp struct {
	Category string `json:"category"`
	// 硬件配置
	Hardware string `json:"hardware"`
	CentralProcessorManufacturer	string	`json:"central_processor_manufacture"`
	// 处理器架构,如：x86_64\aarch64
	CentralProcessorArch			string	`json:"central_processor_arch"`
	//功率
	Power string `json:"power"`
	//设备所占用的U数 机柜参数Unit 1U = 44.45mm
	Unit 		uint	`json:"unit"`	
	//是否金融信创生态产品: yes or no
	IsFITIEcoProduct				string	`json:"is_fiti_eco_product"`	
	// Remark 备注
	Remark string `json:"remark"`
	//设备类型ID。
	ID        uint          `json:"id"`
	Creator   string        `json:"creator"`
	CreatedAt times.ISOTime `json:"created_at"`
	UpdatedAt times.ISOTime `json:"updated_at"`
}

//GetDeviceCategorysPage 获取设备类型分页
func GetDeviceCategorysPage(log logger.Logger, repo model.Repo, reqData *GetDeviceCategoryPageReq) (*page.Page, error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 1000 {
		reqData.PageSize = 20
	}
	if reqData.Page < 0 {
		reqData.Page = 0
	}

	cond := model.DeviceCategory{
		Category: reqData.Category,
		Hardware: reqData.Hardware,
		Power:    reqData.Power,
		Remark:   reqData.Remark,
	}

	totalRecords, err := repo.CountDeviceCategorys(&cond)
	if err != nil {
		return nil, err
	}

	pager := page.NewPager(reflect.TypeOf(&DeviceCategoryResp{}), reqData.Page, reqData.PageSize, totalRecords)
	items, err := repo.GetDeviceCategorys(&cond, model.OneOrderBy("id", model.DESC), pager.BuildLimiter())
	if err != nil {
		return nil, err
	}
	for i := range items {
		item, err := convert2DeviceCategoryResult(log, repo, items[i])
		if err != nil {
			return nil, err
		}
		if item != nil {
			pager.AddRecords(item)
		}
	}
	return pager.BuildPage(), nil
}

func convert2DeviceCategoryResult(log logger.Logger, repo model.Repo, mod *model.DeviceCategory) (*DeviceCategoryResp, error) {
	if mod == nil {
		return nil, nil
	}
	result := DeviceCategoryResp{
		ID:        mod.ID,
		CreatedAt: times.ISOTime(mod.CreatedAt),
		UpdatedAt: times.ISOTime(mod.UpdatedAt),
		Category:  mod.Category,
		Hardware:  mod.Hardware,
		CentralProcessorManufacturer: mod.CentralProcessorManufacturer,
		CentralProcessorArch: mod.CentralProcessorArch,
		Power:     mod.Power,
		Unit: 	   mod.Unit,
		IsFITIEcoProduct: mod.IsFITIEcoProduct,
		Remark:    mod.Remark,
		Creator:   mod.Creator,
	}
	return &result, nil
}

//GetDeviceCategoryByID 获取指定ID的设备类型的详细信息
func GetDeviceCategoryByID(log logger.Logger, repo model.Repo, id uint) (*DeviceCategoryResp, error) {
	items, err := repo.GetDeviceCategoryByID(id)
	if err != nil {
		return nil, err
	}

	item, err := convert2DeviceCategoryResult(log, repo, items)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func GetDeviceCategoryQuerys(log logger.Logger, repo model.Repo, param string) (*model.DeviceQueryParamResp, error) {
	return repo.GetDeviceCategoryQuerys(param)
}
