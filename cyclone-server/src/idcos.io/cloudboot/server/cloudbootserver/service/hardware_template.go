package service

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/voidint/binding"
	"github.com/voidint/page"
	"idcos.io/cloudboot/model"
)

type (
	// HardwareTplResp 硬件配置模板返回结构
	HardwareTplResp struct {
		ID uint `json:"id"`
		//是否是内置模板'Yes|No'
		Builtin string `json:"builtin"`
		//模板名称
		Name string `json:"name"`
		//硬件制造厂商
		Vendor string `json:"vendor"`
		//型号名称
		ModelName string `json:"model_name"`
		//模板内容(数据结构)
		Data Data `json:"data"`
		//创建时间
		CreatedAt string `json:"created_at"`
		//修改时间
		UpdatedAt string `json:"updated_at"`
	}

	// Data 硬件配置下发数据
	Data []struct {
		Category string            `json:"category"`
		Action   string            `json:"action"`
		Metadata map[string]string `json:"metadata"`
	}

	// SaveHardwareTplReq 硬件配置模板保存结构体
	SaveHardwareTplReq struct {
		// 主键
		ID uint `json:"id"`
		//是否是内置模板'Yes|No'
		Builtin string `json:"builtin"`
		//模板名称
		Name string `json:"name"`
		//硬件制造厂商
		Vendor string `json:"vendor"`
		//型号名称
		ModelName string `json:"model_name"`
		//模板内容(数据结构)
		Data Data `json:"data"`
	}
)

// HardwareTplPageReq 获取硬件配置模板分页过滤参数结构
type HardwareTplPageReq struct {
	//是否是内置模板'Yes|No'
	Builtin string `json:"builtin"`
	//模板名称
	Name string `json:"name"`
	//硬件制造厂商
	Vendor string `json:"vendor"`
	//型号名称
	ModelName string `json:"model_name"`
	// 分页页号
	Page int64 `json:"page"`
	// 分页大小
	PageSize int64 `json:"page_size"`
}

// FieldMap 请求字段映射
func (reqData *HardwareTplPageReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.Builtin:   "builtin",
		&reqData.Name:      "name",
		&reqData.Vendor:    "vendor",
		&reqData.ModelName: "model_name",
		&reqData.Page:      "page",
		&reqData.PageSize:  "page_size",
	}
}

// FieldMap 请求字段映射
func (reqData *SaveHardwareTplReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.Builtin:   "builtin",
		&reqData.Name:      "name",
		&reqData.Vendor:    "vendor",
		&reqData.ModelName: "model_name",
		&reqData.Data:      "data",
		&reqData.ID:        "id",
	}
}

// RemoveHardwareTemplate 删除指定ID的硬件模板信息
func RemoveHardwareTemplate(repo model.Repo, id uint) (affected int64, err error) {
	return repo.RemoveHardwareTemplateByID(id)

}

// GetHardwareTplPage 查询硬件配置模板分页列表
func GetHardwareTplPage(repo model.Repo, reqData *HardwareTplPageReq) (pg *page.Page, err error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 100 {
		reqData.PageSize = 20
	}
	if reqData.Page < 0 {
		reqData.Page = 0
	}

	cond := model.HardwareTplCond{
		Builtin:   reqData.Builtin,
		Name:      reqData.Name,
		Vendor:    reqData.Vendor,
		ModelName: reqData.ModelName,
	}

	totalRecords, err := repo.CountHardwareByCond(&cond)
	if err != nil {
		return nil, err
	}

	pager := page.NewPager(reflect.TypeOf(&HardwareTplResp{}), reqData.Page, reqData.PageSize, totalRecords)
	items, err := repo.GetHardwaresByCond(&cond, pager.BuildLimiter())
	if err != nil {
		return nil, err
	}
	for _, mod := range items {
		d := Data{}
		_ = json.Unmarshal([]byte(mod.Data), &d)
		item := &HardwareTplResp{
			ID:        mod.ID,
			Builtin:   mod.Builtin,
			Name:      mod.Name,
			Vendor:    mod.Vendor,
			ModelName: mod.ModelName,
			Data:      d,
			CreatedAt: mod.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: mod.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		pager.AddRecords(item)
	}
	return pager.BuildPage(), nil
}

// SaveHardwareTemplate 保存硬件配置模板信息
func SaveHardwareTemplate(repo model.Repo, reqData *SaveHardwareTplReq) (id uint, err error) {
	dataStr, err := json.Marshal(reqData.Data)
	if err != nil {
		return 0, err
	}

	hardwareTmpl := model.HardwareTemplate{
		Builtin:   model.NO,
		Name:      reqData.Name,
		Vendor:    reqData.Vendor,
		ModelName: reqData.ModelName,
		Data:      string(dataStr),
	}

	hardwareTmpl.ID = reqData.ID

	return repo.SaveHardwareTemplate(&hardwareTmpl)
}
