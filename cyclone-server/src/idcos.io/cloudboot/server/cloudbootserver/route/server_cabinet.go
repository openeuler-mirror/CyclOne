package route

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/voidint/binding"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/server/cloudbootserver/service"
	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
	"idcos.io/cloudboot/utils/upload"
)

// UpdateServerCabinetStatus 保存(新增/修改)机架(柜)状态（是否启用）
func UpdateServerCabinetStatus(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.UpdateServerCabinetStatusReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if err := service.UpdateServerCabinetStatus(log, repo, &reqData); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// AcceptServerCabinet 机架验收
func AcceptServerCabinet(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.AcceptServerCabinetStatusReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if err := service.AcceptServerCabinet(log, repo, &reqData); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// EnableServerCabinet 启用
func EnableServerCabinet(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.EnableServerCabinetStatusReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if err := service.EnableServerCabinet(log, repo, &reqData); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// ReconstructServerCabinet 重建（将已下线的重新建设，状态改为建设中）
func ReconstructServerCabinet(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.ReconstructServerCabinetStatusReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if err := service.ReconstructServerCabinet(log, repo, &reqData); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// SaveServerCabinet 保存(新增)机架(柜)
func SaveServerCabinet(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	var reqData service.SaveServerCabinetReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if user != nil {
		reqData.LoginName = user.LoginName
	}

	if err := service.SaveServerCabinet(log, repo, &reqData); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"id": reqData.ID,
	}))
}

// UpdateServerCabinet 保存(修改)机架(柜)
func UpdateServerCabinet(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error(err)
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}

	reqData := service.SaveServerCabinetReq{
		ID: uint(id),
	}

	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if err = service.SaveServerCabinet(log, repo, &reqData); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"id": reqData.ID,
	}))
}

// RemoveServerCabinetByID 移除机架(柜)
func RemoveServerCabinetByID(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	id, err := strconv.Atoi((chi.URLParam(r, "id")))
	if err != nil {
		log.Error(err)
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}

	//删除前校验，存在子数据不允许删除
	msg := service.RemoveServerCabinetValidte(log, repo, uint(id))
	if msg != "" {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(msg))
		return
	}

	err = service.RemoveServerCabinetByID(log, repo, uint(id))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// GetServerCabinetPage 查询机架(柜)分页列表
func GetServerCabinetPage(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.GetServerCabinetPageReq
	if binding.Form(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	pg, err := service.GetServerCabinetPage(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(pg)),
	)
}

// GetServerCabinetByID 查询指定ID的机架(柜)信息详情
func GetServerCabinetByID(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	id, err := strconv.Atoi((chi.URLParam(r, "id")))
	if err != nil {
		log.Error(err)
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}

	sbc, err := service.GetServerByCabinetID(log, repo, uint(id))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(sbc)),
	)
}

//UploadServerCabinet 加载导入机架(柜)文件
func UploadServerCabinet(w http.ResponseWriter, r *http.Request) {
	//预防乱码
	w.Header().Add("Content-type", "text/html; charset=utf-8")

	//解析并生成临时文件，为后续的工作做准备
	filename, err := upload.GenerateTempFile(r, "server-cabinet")
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"result": filename,
	}))

}

//ImportServerCabinetPriview 导入机架(柜)文件预览
func ImportServerCabinetPriview(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData upload.ImportReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	result, err := service.ImportServerCabinetPriview(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", result))
}

//ImportServerCabinet  导入机架(柜)文件
func ImportServerCabinet(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData upload.ImportReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	err := service.ImportServerCabinet(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// PowerOnServerCabinetByID 机架(柜)开电API
func PowerOnServerCabinetByID(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())

	var reqData service.CabinetPowerBatchOperateReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	cabinets, err := service.PowerOnServerCabinetByID(repo, reqData.IDS)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	if cabinets == 0 {
		HandleErr(r.Context(), w, fmt.Errorf("未查询到有效的机架信息，%d", reqData.IDS))
		return
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(cabinets)),
	)
}

// PowerOffServerCabinetByID 机架(柜)关电API
func PowerOffServerCabinetByID(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	one, err := service.PowerOffServerCabinetByID(repo, uint(id))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	if one == 0 {
		HandleErr(r.Context(), w, fmt.Errorf("未查询到有效的机架信息，%d", id))
		return
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(one)),
	)
}

// BatchUpdateServerCabinetsType 批量更新机架(柜)类型
func BatchUpdateServerCabinetType(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())

	var reqData service.UpdateServerCabinetTypeReq
	if err := myhttp.DecodeJSON(r, &reqData); err != nil {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}	

	affected, err := service.BatchUpdateServerCabinetType(repo, reqData.IDs, reqData.Type)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"affected": affected,
	}))
}

// BatchUpdateServerCabinetRemark 批量更新机架(柜)备注信息
func BatchUpdateServerCabinetRemark(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())

	var reqData service.UpdateServerCabinetRemarkReq
	if err := myhttp.DecodeJSON(r, &reqData); err != nil {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}	

	affected, err := service.BatchUpdateServerCabinetRemark(repo, reqData.IDs, reqData.Remark)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"affected": affected,
	}))
}