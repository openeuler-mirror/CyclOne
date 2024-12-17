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

// GetHardwareTpls 查询硬件配置模板分页列表
func GetHardwareTpls(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())

	var reqData service.HardwareTplPageReq
	if binding.Form(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	pg, err := service.GetHardwareTplPage(repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(pg)),
	)
}

// GetHardwareTemplateByID 查询指定ID的硬件模板
func GetHardwareTemplateByID(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	template, err := service.GetHardwareSettingsByID(repo, uint(id))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(template)),
	)
}

// RemoveHardTemplateByID 删除指定ID的硬件模板
func RemoveHardTemplateByID(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	affected, err := service.RemoveHardwareTemplate(repo, uint(id))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
			"affected": affected,
		}),
	)
}

// SaveHardwareTemplate 新增、修改硬件配置模板
func SaveHardwareTemplate(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())

	var reqData service.SaveHardwareTplReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	id, err := service.SaveHardwareTemplate(repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
			"id": id,
		}),
	)
}

// UpdateHardwareTemplate 新增、修改硬件配置模板
func UpdateHardwareTemplate(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	var reqData service.SaveHardwareTplReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	reqData.ID = uint(id)

	resp, err := service.SaveHardwareTemplate(repo, &reqData)
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
