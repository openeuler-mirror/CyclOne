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

// SaveIDC 新增数据中心
func SaveIDC(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	req := new(service.IDCReq)
	if binding.Json(r, req).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if user != nil {
		req.LoginName = user.LoginName
	}

	mod, err := service.SaveIDC(log, repo, req)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{"id": mod.ID}),
	)
}

// UpdateIDC 修改数据中心
func UpdateIDC(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error(err)
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}

	req := service.IDCReq{
		ID: uint(id),
	}
	if binding.Json(r, &req).CustomHandle(HandleValidateErrs, w) {
		return
	}

	mod, err := service.SaveIDC(log, repo, &req)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{"id": mod.ID}),
	)
}

// UpdateIDCStatus 批量修改数据中心状态
func UpdateIDCStatus(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	req := new(service.IDCUpdateReq)
	if binding.Json(r, req).CustomHandle(HandleValidateErrs, w) {
		return
	}

	err := service.UpdateIDCStatus(log, repo, req)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{"ids": req.IDs}),
	)
}

// RemoveIDCByID 移除数据中心
func RemoveIDCByID(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	idcID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	err = service.RemoveIDCByID(log, repo, uint(idcID))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{"id": idcID}),
	)
}

// GetIDCPage 查询数据中心分页列表
func GetIDCPage(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	req := new(service.IDCPageReq)
	if binding.Form(r, req).CustomHandle(HandleValidateErrs, w) {
		return
	}

	pg, err := service.GetIDCPage(log, repo, req)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(pg)),
	)
}

// GetIDCByID 查询指定ID的数据中心信息详情
func GetIDCByID(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	idcID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	idc, err := service.GetIDCByID(log, repo, uint(idcID))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(idc)),
	)
}
