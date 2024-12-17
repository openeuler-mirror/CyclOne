package http

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
)

const (
	// Success 成功状态
	Success = "success"
	// Failure 非成功状态
	Failure = "failure"
)

// RespBody 响应结构体
type RespBody struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Content map[string]interface{} `json:"content,omitempty"`
}

// ErrRespBody 返回包含失败信息的响应体指针
func ErrRespBody(msg string) *RespBody {
	return &RespBody{
		Status:  Failure,
		Message: msg,
	}
}

// SucceedRespBody 返回包含成功信息的响应体指针
func SucceedRespBody(msg string) *RespBody {
	return &RespBody{
		Status:  Success,
		Message: msg,
	}
}

// NewRespBody 返回响应体指针
func NewRespBody(status, msg string, content map[string]interface{}) *RespBody {
	return &RespBody{
		Status:  status,
		Message: msg,
		Content: content,
	}
}

// DecodeJSON 将Request Body中的JSON内容反序列化到reqData对象中。
// reqData应为结构体的指针。
func DecodeJSON(r *http.Request, reqData interface{}) error {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return nil
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(reqData)
}

const jsonOmit = "-"

// DumpContent 返回由content中的属性和值组成的map。content只能是struct或者struct指针，否则将返回nil。
func DumpContent(content interface{}) map[string]interface{} {
	if content == nil {
		return nil
	}

	contentType := reflect.TypeOf(content)
	contentVal := reflect.ValueOf(content)

	switch contentType.Kind() {
	case reflect.Ptr:
		val := reflect.Indirect(contentVal) //取出指针所指向的实例的value
		tp := val.Type()                    // 通过value来获取真实的实例的Type
		if tp.Kind() != reflect.Struct {
			return nil
		}
		contentType = tp
		contentVal = val
	case reflect.Struct:
	// nothing to do
	default:
		return nil
	}

	mp := make(map[string]interface{}, contentType.NumField())
	for i, nf := 0, contentType.NumField(); i < nf; i++ {
		fName := contentType.Field(i).Name
		fVal := contentVal.FieldByName(fName)
		alias := contentType.Field(i).Tag.Get("json")
		if alias == jsonOmit {
			continue
		}
		if len(alias) > 0 {
			mp[alias] = fVal.Interface()
		} else {
			mp[fName] = fVal.Interface()
		}
	}
	return mp
}
