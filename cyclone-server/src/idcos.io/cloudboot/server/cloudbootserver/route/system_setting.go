package route

import (
	"net/http"

	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/server/cloudbootserver/service"
	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
)

// GetSystemLoginSetting 查询系统登录配置
func GetSystemLoginSetting(w http.ResponseWriter, r *http.Request) {
	conf, _ := middleware.ConfigFromContext(r.Context())

	sett, err := service.GetSystemLoginSetting(conf)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(sett)),
	)
}
