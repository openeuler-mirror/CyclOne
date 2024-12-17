package route

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/voidint/binding"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/server/cloudbootserver/service"
	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
)

// SaveDeviceCategory 保存(新增/修改)设备类型
func SaveDeviceCategory(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	var reqData service.SaveDeviceCategoryReq
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if user != nil {
		reqData.LoginName = user.LoginName
	}

	err := service.SaveDeviceCategory(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"id": reqData.ID,
	}))
}

// RemoveDeviceCategorys 移除设备类型
func RemoveDeviceCategorys(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.DelDeviceCategoryReq
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	affected, err := service.RemoveDeviceCategorys(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"affected": affected,
	}))
}

// GetDeviceCategoryPage 查询设备类型分页列表
func GetDeviceCategoryPage(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.GetDeviceCategoryPageReq
	if binding.Form(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	pg, err := service.GetDeviceCategorysPage(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(pg)),
	)
}

// GetDeviceCategoryByID 查询指定ID的设备类型信息详情
func GetDeviceCategoryByID(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error(err)
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}

	sr, err := service.GetDeviceCategoryByID(log, repo, uint(id))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(sr)),
	)
}

// GetDeviceCategoryQuerys 设备类型的查询（过滤）参数列表
func GetDeviceCategoryQuerys(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	//|category|...
	p, err := service.GetDeviceCategoryQuerys(log, repo, chi.URLParam(r, "param_name"))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(p)),
	)
}
