package service

import (
	"net/http"
	"reflect"

	"github.com/voidint/binding"
	"github.com/voidint/page"

	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
)

//GetOperateLogPageReq 获取操作记录分页请求参数
type GetOperateLogPageReq struct {
	// 请求方式
	HTTPMethod string `json:"http_method"`
	//  路由
	URL string `json:"url"`
	// 源数据
	Source string `json:"source"`
	// 目标数据
	Destination string `json:"destination"`
	//操作类型
	CategoryCode string `json:"column:category_code"`
	//操作类型名称
	CategoryName string `json:"column:category_name"`
	// 分页页号
	Page int64 `json:"page"`
	// 分页大小。默认值:10。阈值: 100。
	PageSize int64 `json:"page_size"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *GetOperateLogPageReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.HTTPMethod:   "http_method",
		&reqData.URL:          "url",
		&reqData.Source:       "source",
		&reqData.Destination:  "destination",
		&reqData.Page:         "page",
		&reqData.PageSize:     "page_size",
		&reqData.CategoryCode: `category_code`,
		&reqData.CategoryName: `category_name`,
	}
}

//GetOperateLogWithPage 获取操作记录分页
func GetOperateLogWithPage(log logger.Logger, repo model.Repo, reqData *GetOperateLogPageReq) (*page.Page, error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 100 {
		reqData.PageSize = 20
	}
	if reqData.Page < 0 {
		reqData.Page = 0
	}

	cond := model.OperateLog{
		HTTPMethod:   reqData.HTTPMethod,
		URL:          reqData.URL,
		Source:       reqData.Source,
		CategoryCode: reqData.CategoryCode,
		CategoryName: reqData.CategoryName,
		Destination:  reqData.Destination,
	}

	totalRecords, err := repo.CountOperateLog(&cond)
	if err != nil {
		return nil, err
	}

	pager := page.NewPager(reflect.TypeOf(&model.OperateLog{}), reqData.Page, reqData.PageSize, totalRecords)
	items, err := repo.GetOperateLogByCond(&cond, model.OneOrderBy("id", model.DESC), pager.BuildLimiter())
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		pager.AddRecords(item)
	}

	return pager.BuildPage(), nil
}
