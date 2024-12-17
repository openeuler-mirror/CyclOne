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

// SaveDeviceSettings 批量保存设备装机参数
func SaveDeviceSettings(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())
	lim, _ := middleware.DHCPLimiterFromContext(r.Context())

	var reqData service.DeviceSettings
	if err := myhttp.DecodeJSON(r, &reqData); err != nil {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}

	if errs := reqData.Validate(r, []binding.Error{}); errs.Len() > 0 {
		HandleValidateErrs(errs, w)
		return
	}

	succeeds, err := service.SaveDeviceSettings(log, repo, conf, lim, &service.SaveDeviceSettingsReq{
		Settings:    reqData,
		CurrentUser: user,
	})
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
			"succeeds": succeeds,
		}),
	)
}

// SaveDeviceSettings 批量保存设备装机参数并重新安装系统
func SaveDeviceSettingsAndReinstalls(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())
	lim, _ := middleware.DHCPLimiterFromContext(r.Context())

	var reqData service.DeviceSettings
	if err := myhttp.DecodeJSON(r, &reqData); err != nil {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}

	if errs := reqData.Validate(r, []binding.Error{}); errs.Len() > 0 {
		HandleValidateErrs(errs, w)
		return
	}

	succeeds, err := service.SaveDeviceSettingsAndReinstalls(log, repo, conf, lim, &service.SaveDeviceSettingsReq{
		Settings:    reqData,
		CurrentUser: user,
	})
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
			"succeeds": succeeds,
		}),
	)
}

// SaveDeviceSettingsWithoutInstalls 批量保存设备装机参数(忽略部署）
func SaveDeviceSettingsWithoutInstalls(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())
	lim, _ := middleware.DHCPLimiterFromContext(r.Context())

	var reqData service.DeviceSettingsWithoutInstalls
	if err := myhttp.DecodeJSON(r, &reqData); err != nil {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}

	succeeds, err := service.SaveDeviceSettingsWithoutInstalls(log, repo, conf, lim, &service.SaveDeviceSettingsWithoutInstallsReq{
		Settings:    reqData,
		CurrentUser: user,
	})
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
			"succeeds": succeeds,
		}),
	)
}

// Reinstalls 批量重装设备
func Reinstalls(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	lim, _ := middleware.DHCPLimiterFromContext(r.Context())

	var reqData service.ReinstallsReq
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if err := service.Reinstalls(log, repo, conf, lim, &reqData); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// Reinstalls 批量重装设备（根据规则引擎自动生成装机参数）
func AutoReinstalls(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	lim, _ := middleware.DHCPLimiterFromContext(r.Context())

	var reqData service.AutoReinstallsReq
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	succeeds, err := service.AutoReinstalls(log, repo, conf, lim, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
			"succeeds": succeeds,
		}),
	)	
}

// CancelInstalls 批量取消安装设备
func CancelInstalls(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.CancelInstallsReq
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if err := service.CancelInstalls(log, repo, &reqData); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// RemoveDeviceSettings 批量删除设备装机参数
func RemoveDeviceSettings(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.RemoveDeviceSettingsReq
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if err := service.RemoveDeviceSettings(log, repo, &reqData); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// GetDeviceSettingPage 查询设备装机参数分页列表
func GetDeviceSettingPage(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.GetDeviceSettingPageReq
	if binding.Form(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	pg, err := service.GetDeviceSettingPage(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(pg)),
	)
}

// GetDeviceSettingBySN 查询指定sn的装机参数信息
func GetDeviceSettingBySN(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	sett, err := service.GetDeviceSettingBySN(log, repo, chi.URLParam(r, "sn"))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(sett)),
	)
}

// GetNetworkSettingBySN 返回指定设备的装机参数(业务网络配置)
func GetNetworkSettingBySN(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	setting, err := service.GetNetworkSettingBySN(log, repo, chi.URLParam(r, "sn"))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(setting)),
	)
}

// GetOSUserSettingsBySN 返回指定设备的操作系统用户配置参数
func GetOSUserSettingsBySN(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	reqData := service.GetOSUserSettingBySNReq{
		SN: chi.URLParam(r, "sn"),
	}

	if errs := reqData.Validate(r, []binding.Error{}); errs.Len() > 0 {
		HandleValidateErrs(errs, w)
		return
	}

	sett, err := service.GetOSUserSettingBySN(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(sett)),
	)
}

//CountDeviceInstallStatic 装机信息统计
func CountDeviceInstallStatic(w http.ResponseWriter, r *http.Request) {
	log, _ := middleware.LoggerFromContext(r.Context())
	repo, _ := middleware.RepoFromContext(r.Context())

	count, err := service.CountDeviceInstallStatic(log, repo)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(count)),
	)
}


//UpdateDeviceSetting
func UpdateDeviceSetting(w http.ResponseWriter, r *http.Request) {
	log, _ := middleware.LoggerFromContext(r.Context())
	repo, _ := middleware.RepoFromContext(r.Context())
	
	var reqData service.UpdateDeviceSettingReq
	
	if err := myhttp.DecodeJSON(r, &reqData); err != nil {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}

	err := service.UpdateDeviceSetting(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}


// SetInstallsOK 批量设置部署状态=success
func SetInstallsOK(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	lim, _ := middleware.DHCPLimiterFromContext(r.Context())

	var reqData service.SetInstallsOKReq
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if err := service.SetInstallsOK(log, repo, conf, lim, &reqData); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}