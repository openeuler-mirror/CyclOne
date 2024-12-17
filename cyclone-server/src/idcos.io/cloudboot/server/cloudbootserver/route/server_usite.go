package route

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/voidint/binding"

	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/server/cloudbootserver/service"
	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
	"idcos.io/cloudboot/utils/upload"
	"idcos.io/cloudboot/utils/user"
)

// SaveServerUSite 保存(新增/修改)机位(U位)
func SaveServerUSite(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	loginUser, _ := middleware.LoginUserFromContext(r.Context())

	var reqData service.SaveServerUSiteReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	reqData.LoginUser = loginUser

	if err := service.SaveServerUSite(log, repo, &reqData); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"id": reqData.ID,
	}))
}

// UpdateServerUSite 保存(新增/修改)机位(U位)
func UpdateServerUSite(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	loginUser, _ := middleware.LoginUserFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error(err)
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}

	reqData := service.SaveServerUSiteReq{
		ID:        uint(id),
		LoginUser: loginUser,
	}

	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if err := service.SaveServerUSite(log, repo, &reqData); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"id": reqData.ID,
	}))
}

// GetServerUSitePage 查询机位(U位)分页列表
func GetServerUSitePage(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	loginUser, _ := middleware.LoginUserFromContext(r.Context())

	var reqData service.GetServerUSitePageReq
	if binding.Form(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	reqData.GetNameFromUAM = user.GetUserByLoginName(log, loginUser.Token, conf.UAM.RootEndpoint)
	pg, err := service.GetServerUSitePage(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(pg)),
	)
}

// GetServerUSiteByID 查询指定ID的机位(U位)信息详情
func GetServerUSiteByID(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	uSite, err := service.GetServerUSiteByID(repo, uint(id))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(uSite)),
	)
}

// GetUsiteTree 查询U位目录树
func GetUsiteTree(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	//conf, _ := middleware.ConfigFromContext(r.Context())
	//loginUser, _ := middleware.LoginUserFromContext(r.Context())

	var reqData service.UsiteTreeReq
	if binding.Form(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	//reqData.GetNameFromUAM = user.GetUserByLoginName(log, loginUser.Token, conf.UAM.RootEndpoint)
	tree, err := service.GetUsiteTree(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(tree)),
	)
}

// BatchUpdateServerUSitesStatus 批量修改机位状态入参结构体
func BatchUpdateServerUSitesStatus(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	loginUser, _ := middleware.LoginUserFromContext(r.Context())

	var reqData service.UpdateServerUSiteStatusReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	reqData.LoginUser = loginUser

	affected, err := service.BatchUpdateServerUSitesStatus(repo, reqData.IDs, reqData.Status)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"affected": affected,
	}))
}

// BatchUpdateServerUSitesStatusByCond 根据其他定制性条件进行批量更新机位状态信息
func BatchUpdateServerUSitesStatusByCond(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	loginUser, _ := middleware.LoginUserFromContext(r.Context())

	var reqData service.UpdateServerUSiteStatusReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	reqData.LoginUser = loginUser

	affected, err := service.BatchUpdateServerUSitesStatusByCond(log, repo, reqData.USites, reqData.Status)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"affected": affected,
	}))
}

// DeleteServerUSitePort 删除机位端口号
func DeleteServerUSitePort(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error(err)
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}

	affected, err := service.DeleteServerUSitePort(repo, uint(id))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	if affected == 0 {
		render.JSON(w, http.StatusNotFound, myhttp.ErrRespBody(errors.New("未查寻到有有效的记录信息").Error()))
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"id": id,
	}))
}

// DeleteServerUSite 删除机位
func DeleteServerUSite(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error(err)
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}

	affected, err := service.RemoveServerUSiteByID(repo, uint(id))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	if affected == 0 {
		render.JSON(w, http.StatusNotFound, myhttp.ErrRespBody(errors.New("未查寻到有有效的记录信息").Error()))
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"id": id,
	}))
}

//UploadServerUSite 加载导入机位(U位)文件
func UploadServerUSite(w http.ResponseWriter, r *http.Request) {
	//预防乱码
	w.Header().Add("Content-type", "text/html; charset=utf-8")

	//解析并生成临时文件，为后续的工作做准备
	filename, err := upload.GenerateTempFile(r, "server-usite")
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"result": filename,
	}))
}

//UploadServerUSitePort 加载导入机位(U位)端口号文件
func UploadServerUSitePort(w http.ResponseWriter, r *http.Request) {
	//预防乱码
	w.Header().Add("Content-type", "text/html; charset=utf-8")

	//解析并生成临时文件，为后续的工作做准备
	filename, err := upload.GenerateTempFile(r, "server-usite-port")
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"result": filename,
	}))
}

//ImportServerUSitePriview 导入机架(柜)文件预览
func ImportServerUSitePriview(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData upload.ImportReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	result, err := service.ImportServerUSitePreview(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", result))
}

//ImportServerUSitePortsPriview 导入机架(柜)文件预览
func ImportServerUSitePortsPriview(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData upload.ImportReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	result, err := service.ImportServerUSitePortsPreview(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", result))
}

//ImportServerUSite  导入机架(柜)文件
func ImportServerUSite(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	var reqData upload.ImportReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	reqData.UserName = user.LoginName //user.Name
	err := service.ImportServerUSite(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

//ImportServerUSitePort  导入机架(柜)文件
func ImportServerUSitePort(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData upload.ImportReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	err := service.ImportServerUSitePort(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// GetPhysicalAreas 物理区域列表
// @Discard
func GetPhysicalAreas(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.PhysicalAreaConnd
	if binding.Form(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	p, err := service.GetPhysicalAreas(log, repo, reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(p)),
	)
}

// BatchUpdateServerUSitesRemark 批量更新机位备注信息
func BatchUpdateServerUSitesRemark(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())

	var reqData service.BatchUpdateServerUSitesRemarkReq
	if err := myhttp.DecodeJSON(r, &reqData); err != nil {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}	

	affected, err := service.BatchUpdateServerUSitesRemark(repo, reqData.IDs, reqData.Remark)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"affected": affected,
	}))
}