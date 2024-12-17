package route

import (
	"net/http"

	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/server/cloudbootserver/service"
	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
)

// GetDataDicts 根据type参数查询筛选数据字典信息
func GetDataDicts(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	d, err := service.GetDataDict(log, repo, r.URL.Query().Get("type"))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
			"items": d,
		}),
	)
}

//
////AddDataDicts 增加数据字典
//func AddDataDicts(w http.ResponseWriter, r *http.Request) {
//	repo, _ := middleware.RepoFromContext(r.Context())
//	log, _ := middleware.PanicFileFromContext(r.Context())
//
//	var reqData []*service.DataDict
//	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
//		return
//	}
//	err := service.AddDataDicts(log, repo, reqData)
//	if err != nil {
//		HandleErr(r.Context(), w, err)
//		return
//	}
//
//	render.JSON(w, http.StatusOK,
//		myhttp.SucceedRespBody("操作成功"),
//	)
//}
//
////DelDataDicts
//func DelDataDicts(w http.ResponseWriter, r *http.Request) {
//	repo, _ := middleware.RepoFromContext(r.Context())
//	log, _ := middleware.PanicFileFromContext(r.Context())
//
//	var reqData []*service.DelDataDictReq
//	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
//		return
//	}
//	err := service.DelDataDicts(log, repo, reqData)
//	if err != nil {
//		HandleErr(r.Context(), w, err)
//		return
//	}
//
//	render.JSON(w, http.StatusOK,
//		myhttp.SucceedRespBody("操作成功")),
//	)
//}
//
////UpdateDataDicts
//func UpdateDataDicts(w http.ResponseWriter, r *http.Request) {
//	repo, _ := middleware.RepoFromContext(r.Context())
//	log, _ := middleware.PanicFileFromContext(r.Context())
//
//	var reqData []*service.DataDict
//	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
//		return
//	}
//	err := service.UpdateDataDicts(log, repo, reqData)
//	if err != nil {
//		HandleErr(r.Context(), w, err)
//		return
//	}
//
//	render.JSON(w, http.StatusOK,
//		myhttp.SucceedRespBody("操作成功")),
//	)
//}
