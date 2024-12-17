package route

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/server/cloudbootserver/service"
	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
)

// GetDeviceLogByDeviceSettingID 根据SN号返回安装的系统日志信息
func GetDeviceLogByDeviceSettingID(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "device_setting_id"))
	if err != nil {
		render.JSON(w, http.StatusOK, myhttp.ErrRespBody(err.Error()))
		return
	}

	logs, _ := service.GetDeviceLogByDeviceSettingID(log, repo, uint(id))

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(logs)),
	)
}
