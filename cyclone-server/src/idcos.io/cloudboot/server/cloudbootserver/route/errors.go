package route

import (
	"context"
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/voidint/binding"

	"idcos.io/cloudboot/limiter"
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
	} else if errs.Has(binding.BusinessError) && strings.Contains(errs.Error(), http.StatusText(http.StatusNotFound)) {
		render.JSON(w, http.StatusNotFound, myhttp.ErrRespBody("资源不存在"))
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

	if strings.Contains(err.Error(), "校验") { // TODO 不好的实现方式，待修改。
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}
	switch err {
	case gorm.ErrRecordNotFound:
		render.JSON(w, http.StatusNotFound, myhttp.ErrRespBody("资源不存在"))
	case limiter.ErrBucketNotFound:
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody("未找到设备对应的令牌桶(TOR)"))
	case limiter.ErrInvalidOrReturnedToken:
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody("设备未曾申请令牌或已归还令牌"))
	default:
		render.JSON(w, http.StatusInternalServerError, myhttp.ErrRespBody(err.Error()))
	}
}
