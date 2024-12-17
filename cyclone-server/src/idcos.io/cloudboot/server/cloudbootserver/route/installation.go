package route

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/voidint/binding"

	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/server/cloudbootserver/service"
	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
)

// ReportInstallProgress 上报安装进度
func ReportInstallProgress(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	logger, _ := middleware.LoggerFromContext(r.Context())
	lim, _ := middleware.DHCPLimiterFromContext(r.Context())

	reqData := service.InstallProgressReq{
		SN: chi.URLParam(r, "sn"),
	}
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	if err := service.ReportInstallProgress(logger, repo, conf, lim, &reqData); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// GetInstallationStatus 查询设备安装状态
func GetInstallationStatus(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	stat, err := service.GetInstallationStatus(log, repo, chi.URLParam(r, "sn"))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(stat)),
	)
}

// IsInInstallList 查询指定设备是否在装机队列
func IsInInstallList(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	logger, _ := middleware.LoggerFromContext(r.Context())

	inList, _ := service.GetIsInInstallListBySN(logger, repo, chi.URLParam(r, "sn"))
	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"result": inList,
	}))
}
