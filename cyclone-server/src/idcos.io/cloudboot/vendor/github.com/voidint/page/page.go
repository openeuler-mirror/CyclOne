package page

import (
	"fmt"
	"reflect"
)

// Limiter database limit clause in MySQL.
//
// eg. SELECT * FROM `user` LIMIT 0,2
//
// eg. SELECT * FROM `user` LIMIT 2 OFFSET 0
type Limiter struct {
	Offset int64
	Limit  int64
}

// Page business layer pagination struct.
type Page struct {
	Page         int64         `json:"page"` // page number
	PageSize     int64         `json:"page_size"`
	TotalPages   int64         `json:"total_pages"`
	TotalRecords int64         `json:"total_records"`
	Records      []interface{} `json:"records"`
}

// Pager business layer pagination interface.
type Pager interface {
	// BuildLimiter get database limit clause for data access layer.
	BuildLimiter() *Limiter
	// AddRecords append record to collection.
	AddRecords(records ...interface{}) error
	// BuildPage get page object for serialization.
	BuildPage() *Page
}

// NewPager return business layer pagination instance.
func NewPager(elemType reflect.Type, page, pageSize, totalRecords int64) Pager {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if totalRecords < 0 {
		totalRecords = 0
	}
	capacity := pageSize
	if pageSize > totalRecords {
		capacity = totalRecords
	}

	totalPages, _ := calcTotalPages(pageSize, totalRecords)
	return &pagerImpl{
		page:         page,
		pageSize:     pageSize,
		totalPages:   totalPages,
		totalRecords: totalRecords,
		records:      make([]interface{}, 0, capacity),
		elemType:     elemType,
	}
}

type pagerImpl struct {
	page         int64
	pageSize     int64
	totalPages   int64
	totalRecords int64
	records      []interface{}
	elemType     reflect.Type
}

func (p *pagerImpl) BuildLimiter() *Limiter {
	return &Limiter{
		Limit:  p.pageSize,
		Offset: (p.page - 1) * p.pageSize,
	}
}

func (p *pagerImpl) isAcceptableElem(k interface{}) bool {
	return reflect.TypeOf(k) == p.elemType
}

func (p *pagerImpl) AddRecords(records ...interface{}) error {
	for _, record := range records {
		if !p.isAcceptableElem(record) {
			return fmt.Errorf("invalid element: %#v", record)
		}
	}
	p.records = append(p.records, records...)
	return nil
}

func (p *pagerImpl) BuildPage() *Page {
	return &Page{
		Page:         p.page,
		PageSize:     p.pageSize,
		TotalPages:   p.totalPages,
		TotalRecords: p.totalRecords,
		Records:      p.records,
	}
}

// EmptyPage return empty record object.
func EmptyPage(page, pageSize int64) *Page {
	return &Page{
		Page:         page,
		PageSize:     pageSize,
		TotalPages:   0,
		TotalRecords: 0,
		Records:      make([]interface{}, 0),
	}
}
