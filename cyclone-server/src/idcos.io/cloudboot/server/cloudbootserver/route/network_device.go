package route

import (
	"net/http"
	"strconv"

	myhttp "idcos.io/cloudboot/utils/http"

	"github.com/go-chi/chi"
	"github.com/voidint/binding"
	"idcos.io/cloudboot/middleware"

	"idcos.io/cloudboot/server/cloudbootserver/service"

	"idcos.io/cloudboot/utils/http/render"
	"idcos.io/cloudboot/utils/upload"
)

// DeleteNetworkDeviceByID 移除网络设备
func DeleteNetworkDeviceByID(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error(err)
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}
	if err := service.RemoveNetworkDeviceByID(repo, uint(id)); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// RemoveNetworkDevices 批量删除网络设备
func RemoveNetworkDevices(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.DelNetworkDeviceReq
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	affected, err := service.RemoveNetworkDevices(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"affected": affected,
	}))
}

// GetNetworkDeviceByID 查询指定ID的网络设备
func GetNetworkDeviceByID(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error(err)
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}
	resp, err := service.GetNetworkDeviceByID(repo, uint(id))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(resp)))
}

// GetNetworkDevicePage 查询网络设备分页列表
func GetNetworkDevicePage(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())

	var reqData service.NetworkDevicePageReq
	if binding.Form(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	pg, err := service.GetNetworkDevicesPage(repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(pg)),
	)
}

// SaveNetworkDevice 查询网络设备分页列表
func SaveNetworkDevice(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	var reqData service.SaveNetworkDeviceReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if user != nil {
		reqData.LoginName = user.LoginName
	}

	device, err := service.SaveNetworkDevice(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
			"fixed_asset_number": device.FixedAssetNumber,
			"id":                 device.ID,
		}),
	)
}


// UploadNetworkDevices 网络设备导入文件上传
func UploadNetworkDevices(w http.ResponseWriter, r *http.Request) {
	//预防乱码
	w.Header().Add("Content-type", "text/html; charset=utf-8")

	//解析并生成临时文件，为后续的工作做准备
	filename, err := upload.GenerateTempFile(r, "network-device")
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"result": filename,
	}))
}

// ImportNetworkDevicesPreview 导入设备文件预览
func ImportNetworkDevicesPreview(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData upload.ImportReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	result, err := service.ImportNetworkDevicePreview(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", result))
}

//ImportNetworkDevices 导入
func ImportNetworkDevices(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	var reqData upload.ImportReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if user != nil {
		reqData.UserName = user.LoginName
	}

	err := service.ImportNetworkDevices(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}
