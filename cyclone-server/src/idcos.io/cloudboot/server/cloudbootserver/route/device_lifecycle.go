package route

import (
	"net/http"
	"github.com/voidint/binding"
	"github.com/go-chi/chi"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/server/cloudbootserver/service"
	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
)


// GetDeviceLifecycleBySN
func GetDeviceLifecycleBySN(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	devLifecycle, err := service.GetDeviceLifecycleBySN(log, repo, chi.URLParam(r, "sn"))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(devLifecycle)),
	)
}

// UpdateDeviceLifecycleBySN
func UpdateDeviceLifecycleBySN(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	reqData := service.UpdateDeviceLifecycleReq{
		SN:        chi.URLParam(r, "sn"),
	}
	if user != nil {
		reqData.LoginName = user.LoginName
	}

	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	err := service.UpdateDeviceLifecycleBySN(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"sn": reqData.SN,
	}))
}

// BatchUpdateDeviceLifecycleBySN 批量更改设备生命周期信息
func BatchUpdateDeviceLifecycleBySN(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	req := new(service.BatchUpdateDeviceLifecycles)
	if binding.Bind(r, req).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if user != nil {
		req.LoginName = user.LoginName
	}

	succeedSNs, totalAffected, err := service.BatchUpdateDeviceLifecycleBySN(log, repo, conf, req)
	if err != nil {
		render.JSON(w, http.StatusOK,
			myhttp.NewRespBody(myhttp.Failure, "操作失败", map[string]interface{}{"succeed_sns": succeedSNs, "total_affected": totalAffected, "detail": err.Error()}),
		)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{"succeed_sns": succeedSNs, "total_affected": totalAffected, "detail":"success"}),
	)
}