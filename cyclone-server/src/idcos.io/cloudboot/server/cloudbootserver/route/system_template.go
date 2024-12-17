package route

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/voidint/binding"

	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/server/cloudbootserver/service"
	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
)

// GetSystemTemplateByID 查询指定ID的系统安装模板
func GetSystemTemplateByID(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error(err)
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}
	tpl, err := service.GetSystemTemplateByID(log, repo, uint(id))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(tpl)),
	)
}

// GetSystemTemplatePage 按条件查询系统安装模板分页列表
func GetSystemTemplatePage(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.GetSystemTemplatePageReq
	if binding.Form(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	pg, err := service.GetSystemTemplatePage(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(pg)),
	)
}

// AddSystemTemplate 新增系统安装模板
func AddSystemTemplate(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.SaveSystemTemplateReq
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	resp, err := service.SaveSystemTemplate(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
			"id": resp,
		}),
	)
}

// SaveSystemTemplate  新增系统模板
func SaveSystemTemplate(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.SaveSystemTemplateReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	resp, err := service.SaveSystemTemplate(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
			"id": resp,
		}),
	)
}

// UpdateSystemTemplateByID 更新系统安装模板
func UpdateSystemTemplateByID(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error(err)
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}

	reqData := service.SaveSystemTemplateReq{
		ID: uint(id),
	}
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	resp, err := service.SaveSystemTemplate(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
			"id": resp,
		}),
	)
}

// RemoveSystemTemplateByID 删除系统安装模板
func RemoveSystemTemplateByID(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error(err)
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}

	if err := service.RemoveSystemTemplate(log, repo, uint(id)); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// GetSystemTemplateBySN 根据SN查询系统安装模板
func GetSystemTemplateBySN(w http.ResponseWriter, r *http.Request) {
	log, _ := middleware.LoggerFromContext(r.Context())
	repo, _ := middleware.RepoFromContext(r.Context())

	typ := strings.TrimSpace(r.FormValue("type"))
	if typ == "" {
		typ = "raw"
	}

	tpl, err := service.GetSystemTemplateBySN(log, repo, chi.URLParam(r, "sn"))
	if err != nil {
		if typ == "raw" {
			render.Text(w, http.StatusNotFound, []byte(http.StatusText(http.StatusNotFound)))
		} else {
			HandleErr(r.Context(), w, err)
		}
		return
	}

	if typ == "raw" {
		render.Text(w, http.StatusOK, []byte(tpl.Content))
		return
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(tpl)),
	)
}
