package route

import (
	"context"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/voidint/binding"

	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
)

// HandleValidateErrs 数据校验错误处理。
// 将参数绑定校验框架(github.com/voidint/binding)产生的校验错误，以一定格式写入http response body。
func HandleValidateErrs(errs binding.Errors, w http.ResponseWriter) bool {
	if errs.Len() <= 0 {
		return false
	}

	if errs.Has(binding.SystemError) {
		render.JSON(w, http.StatusInternalServerError, myhttp.ErrRespBody(errs.Error()))
	} else {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(errs.Error()))
	}
	return true
}

// HandleErr service层的业务逻辑错误处理
func HandleErr(ctx context.Context, w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	switch err {
	case gorm.ErrRecordNotFound:
		render.JSON(w, http.StatusNotFound, myhttp.ErrRespBody("资源不存在"))

	default:
		render.JSON(w, http.StatusInternalServerError, myhttp.ErrRespBody(err.Error()))
	}
}
