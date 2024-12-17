package service

import (
	"net/http"
	"reflect"
	"time"

	"github.com/voidint/binding"
	"github.com/voidint/page"

	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
)

//GetAPILogPageReq 获取API记录分页请求参数
type GetAPILogPageReq struct {
	//操作时间起
	CreatedAtStart time.Time `json:"created_at_start"`
	//操作时间止
	CreatedAtEnd time.Time `json:"created_at_end"`
	//捷足先登人
	Operator string `json:"operator"`
	//接口描述
	Description string `json:"description"`
	//API信息
	API string `json:"api"`
	//请求方法
	Method string `json:"method"`
	//API状态
	Status string `json:"status"`
	//耗时起
	Cost1 float64 `json:"cost1"`
	//耗时止
	Cost2 float64 `json:"cost2"`
	// 分页页号
	Page int64 `json:"page"`
	// 分页大小。默认值:10。阈值: 100。
	PageSize int64 `json:"page_size"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *GetAPILogPageReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.CreatedAtStart: "created_at_start",
		&reqData.CreatedAtEnd:   "created_at_end",
		&reqData.Operator:       "operator",
		&reqData.Description:    "description",
		&reqData.API:            "api",
		&reqData.Method:         "method",
		&reqData.Status:         "status",
		&reqData.Cost1:          "cost1",
		&reqData.Cost2:          "cost2",
		&reqData.Page:           "page",
		&reqData.PageSize:       "page_size",
	}
}

//GetAPILogWithPage 获取API记录分页
func GetAPILogWithPage(log logger.Logger, repo model.Repo, reqData *GetAPILogPageReq) (*page.Page, error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 100 {
		reqData.PageSize = 20
	}
	if reqData.Page < 0 {
		reqData.Page = 0
	}

	cond := model.APILogCond{
		CreatedAtStart: reqData.CreatedAtStart,
		CreatedAtEnd:   reqData.CreatedAtEnd,
		Operator:       reqData.Operator,
		Description:    reqData.Description,
		API:            reqData.API,
		Method:         reqData.Method,
		Status:         reqData.Status,
		Cost1:          reqData.Cost1,
		Cost2:          reqData.Cost2,
	}

	totalRecords, err := repo.CountAPILog(&cond)
	if err != nil {
		return nil, err
	}

	pager := page.NewPager(reflect.TypeOf(&model.APILog{}), reqData.Page, reqData.PageSize, totalRecords)
	items, err := repo.GetAPILogByCond(&cond, model.OneOrderBy("id", model.DESC), pager.BuildLimiter())
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		pager.AddRecords(item)
	}

	return pager.BuildPage(), nil
}
