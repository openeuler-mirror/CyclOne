package route

import (
	"net/http"

	"github.com/go-chi/chi"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/server/cloudbootserver/service"
	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
)

// GetHardwareSettingBySN 查询设备的硬件配置参数(RAID、BIOS、OOB)
func GetHardwareSettingBySN(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	items, err := service.GetHardwareSettingsBySN(log, repo, chi.URLParam(r, "sn"))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(items)),
	)
}

// GetHardwareInfoBySN 查询设备的硬件配置参数(序列号、厂商、型号、设备类型、硬件备注)
func GetHardwareInfoBySN(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	items, err := service.GetHardwareInfoBySN(log, repo, chi.URLParam(r, "sn"))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(items)),
	)
}
