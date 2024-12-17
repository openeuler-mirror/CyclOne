package route

import (
	"net/http"

	"idcos.io/cloudboot/middleware"
	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
)

// GetSambaConf 获取Server Samba配置
func GetSambaConf(w http.ResponseWriter, r *http.Request) {
	//locale, _ := middleware.LocaleFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())

	_ = render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(&conf.Samba)),
	)
}
