package route

import (
	"net/http"

	"github.com/voidint/binding"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/server/cloudbootserver/service"
	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
)

// SaveDeviceSettingRule 保存(新增/修改)规则记录
func SaveDeviceSettingRule(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	var reqData service.SaveDeviceSettingRuleReq
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if user != nil {
		reqData.LoginName = user.LoginName
	}

	err := service.SaveDeviceSettingRule(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"id": reqData.ID,
	}))
}

// RemoveDeviceSettingRules 移除设备类型
func RemoveDeviceSettingRules(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.DelDeviceSettingRuleReq
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	affected, err := service.RemoveDeviceSettingRules(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"affected": affected,
	}))
}

//// GetDeviceSettingRulePage 获取规则记录分页
func GetDeviceSettingRulePage(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.GetDeviceSettingRulePageReq
	if binding.Form(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	pg, err := service.GetDeviceSettingRulePage(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(pg)),
	)
}

// GetDeviceSettingRuleByID 查询指定ID的设备类型信息详情
//func GetDeviceSettingRuleByID(w http.ResponseWriter, r *http.Request) {
//	repo, _ := middleware.RepoFromContext(r.Context())
//	log, _ := middleware.LoggerFromContext(r.Context())
//
//	id, err := strconv.Atoi(chi.URLParam(r, "id"))
//	if err != nil {
//		log.Error(err)
//		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
//		return
//	}
//
//	sr, err := service.GetDeviceSettingRuleByID(log, repo, uint(id))
//	if err != nil {
//		HandleErr(r.Context(), w, err)
//		return
//	}
//	render.JSON(w, http.StatusOK,
//		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(sr)),
//	)
//}

// GetDeviceSettingRuleQuerys 设备类型的查询（过滤）参数列表
//func GetDeviceSettingRuleQuerys(w http.ResponseWriter, r *http.Request) {
//	repo, _ := middleware.RepoFromContext(r.Context())
//	log, _ := middleware.LoggerFromContext(r.Context())
//
//	//|category|...
//	p, err := service.GetDeviceSettingRuleQuerys(log, repo, chi.URLParam(r, "param_name"))
//	if err != nil {
//		HandleErr(r.Context(), w, err)
//		return
//	}
//	render.JSON(w, http.StatusOK,
//		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(p)),
//	)
//}
