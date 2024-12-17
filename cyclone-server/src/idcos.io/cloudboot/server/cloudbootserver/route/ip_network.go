package route

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/voidint/binding"

	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/server/cloudbootserver/service"
	"idcos.io/cloudboot/utils"
	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
	"idcos.io/cloudboot/utils/upload"
)

// GetIPNetworkByID 查询指定ID网段
func GetIPNetworkByID(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error(err)
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}
	one, err := service.GetIPNetworkByID(log, repo, uint(id))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(one)),
	)
}

// GetIPNetworkPage 按条件查询网段分页列表
func GetIPNetworkPage(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.GetIPNetworkPageReq
	if binding.Form(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	pg, err := service.GetIPNetworkPage(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(pg)),
	)
}

// SaveIPNetwork 新增网段
func SaveIPNetwork(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	var reqData service.SaveIPNetworkReq
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if user != nil {
		reqData.LoginName = user.LoginName
	}

	if err := service.SaveIPNetwork(log, repo, &reqData); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"id": reqData.ID,
	}))
}

// UpdateIPNetwork 更新网段
func UpdateIPNetwork(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error(err)
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}

	reqData := service.SaveIPNetworkReq{
		ID: id,
	}

	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if err := service.SaveIPNetwork(log, repo, &reqData); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"id": reqData.ID,
	}))
}

// RemoveIPNetworkByID 删除网段
func RemoveIPNetworkByID(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error(err)
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}

	//删除前校验，存在子数据不允许删除
	msg := service.RemoveIPNetworkValidte(log, repo, uint(id))
	if msg != "" {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(msg))
		return
	}

	if err := service.RemoveIPNetwork(log, repo, uint(id)); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}


// RemoveIPNetworks 批量删除网段
func RemoveIPNetworks(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.DelIPNetworkReq
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	affected, err := service.RemoveIPNetworks(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"affected": affected,
	}))
}


// GetIPSPage 按条件查询网段分页列表
func GetIPSPage(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.GetIPSPageReq
	if binding.Form(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	pg, err := service.GetIPSPage(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(pg)),
	)
}

//ExportIP 导出IP信息
func ExportIP(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())

	reqData := service.GetIPSPageReq{} //ExportDevicesReq{}
	if binding.Bind(r, &reqData).Handle(w) {
		return
	}

	items, err := service.GetExportIP(log, repo, conf, &reqData)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = render.XLSX(w, utils.FileIP, service.ExportIPInfo(items).ToTableRecords())
	if err != nil {
		log.Error(err)
	}
}

//AssignIP 手动分配IP
func AssignIP(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.AssignIPReq
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	err := service.AssignIP(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}


// AssignIPv4 web界面按钮分配ipv4
func AssignIPv4(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.AssignIPv4Req
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	err := service.AssignIPv4(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}


// AssignIPv6 web界面按钮分配ipv6
func AssignIPv6(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.AssignIPv6Req
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	err := service.AssignIPv6(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}


//UnassignIP 手动解除与物理机绑定
func UnassignIP(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.UnassignIPReq
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	err := service.UnassignIP(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// DisableIP 更新IP状态，主要是禁用一个IP
func DisableIP(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.DisableIPReq
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	err := service.DisableIP(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// UploadIPNetworks 导入文件上传
func UploadIPNetworks(w http.ResponseWriter, r *http.Request) {
	//预防乱码
	w.Header().Add("Content-type", "text/html; charset=utf-8")

	//解析并生成临时文件，为后续的工作做准备
	filename, err := upload.GenerateTempFile(r, "ip-networks")
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"result": filename,
	}))
}

// ImportIPNetworksPreview 导入设备文件预览
func ImportIPNetworksPreview(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.ImportPreviewReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	result, err := service.ImportIPNetworksPreview(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", result))
}

//ImportIPNetworks 导入物理机
func ImportIPNetworks(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	var reqData service.ImportPreviewReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if user != nil {
		reqData.LoginName = user.LoginName
	}

	err := service.ImportIPNetworks(log, repo, conf, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}
