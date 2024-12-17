package route

import (
	"net/http"

	"github.com/voidint/binding"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/server/cloudbootserver/service"
	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
)

// GetTemplatesByCond 根据查询条件查询镜像和系统模板信息
func GetTemplatesByCond(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.TemplateReq
	if binding.Form(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	resp, err := service.GetTemplatesByCond(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	osResp := service.TemplateRecordsResp{
		Records: resp,
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(&osResp)),
	)
}

