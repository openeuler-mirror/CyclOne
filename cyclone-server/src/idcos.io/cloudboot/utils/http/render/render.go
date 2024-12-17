package render

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"

	xlsx "idcos.io/cloudboot/utils"
	myhttp "idcos.io/cloudboot/utils/http"
)

// JSON 以JSON格式渲染HTTP Response Body
func JSON(w http.ResponseWriter, code int, body *myhttp.RespBody) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(body)
}

// Text 纯文本格式渲染HTTP Response Body
func Text(w http.ResponseWriter, code int, body []byte) (int, error) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	return w.Write(body)
}

// HTML 以text-html格式请求头返回HTTP Response Body
func HTML(w http.ResponseWriter, code int, body []byte) (int, error) {
	w.Header().Set("Content-Type", "application/text-html; charset=utf-8")
	w.WriteHeader(code)
	return w.Write(body)
}

// CSV 以CSV格式渲染HTTP Response Body
func CSV(w http.ResponseWriter, filename string, records [][]string) error {
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename='%s';filename*=utf-8''%s", filename, filename))
	w.Header().Add("Content-Type", "application/octet-stream")
	return csv.NewWriter(w).WriteAll(records)
}

//XLSX 以xlsx格式渲染HTTP Response Body
func XLSX(w http.ResponseWriter, filename string, records [][]string) error {
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s;" /*filename*=utf-8''%s*/, filename+"_export.xlsx" /*, filename*/))
	w.Header().Add("Content-Type", "application/octet-stream")

	file, err := xlsx.WriteToXLSX(filename, records)
	if err != nil {
		return err
	}
	err = file.Write(w)
	return err
}
