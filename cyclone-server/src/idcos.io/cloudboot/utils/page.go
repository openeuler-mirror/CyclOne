package utils

import (
	"fmt"
	"reflect"
)

// Page 业务逻辑层分页对象
type Page struct {
	Offset       int64         `json:"offset"`
	Limit        int64         `json:"limit"`
	TotalRecords int64         `json:"recordCount"`
	Records      []interface{} `json:"list"`
}

// Pager 业务逻辑层分页接口
type Pager interface {
	// AddRecords 向当前分页中增加行记录
	AddRecords(records ...interface{}) error
	// BuildPage 构建业务逻辑层分页对象
	BuildPage() *Page
}

// pagerImpl 业务逻辑层默认分页接口实现
type pagerImpl struct {
	offset       int64
	limit        int64
	totalRecords int64
	records      []interface{}
	elemType     reflect.Type
}

// NewPager 新建业务逻辑层分页对象。
// elemType 分页中每条记录所代表的对象的反射类型
func NewPager(elemType reflect.Type, offset, limit, totalRecords int64) Pager {
	if offset <= 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 10
	}

	return &pagerImpl{
		offset:       offset,
		limit:        limit,
		totalRecords: totalRecords,
		records:      make([]interface{}, 0, limit),
		elemType:     elemType,
	}
}

// isAcceptableElem 判断对象反射类型是否与分页记录对象反射类型一致。
func (p *pagerImpl) isAcceptableElem(k interface{}) bool {
	return reflect.TypeOf(k) == p.elemType
}

// AddRows 添加多条记录
func (p *pagerImpl) AddRecords(records ...interface{}) error {
	for _, record := range records {
		if !p.isAcceptableElem(record) {
			return fmt.Errorf("invalid element: %#v", record)
		}
	}
	p.records = append(p.records, records...)
	return nil
}

// BuildPage 构建业务逻辑层分页对象
func (p *pagerImpl) BuildPage() (page *Page) {
	return &Page{
		Offset:       p.offset,
		Limit:        p.limit,
		TotalRecords: p.totalRecords,
		Records:      p.records,
	}
}

// EmptyPage 空分页对象
func EmptyPage(offset, limit int64) *Page {
	return &Page{
		Offset:       offset,
		Limit:        limit,
		TotalRecords: 0,
		Records:      make([]interface{}, 0),
	}
}
