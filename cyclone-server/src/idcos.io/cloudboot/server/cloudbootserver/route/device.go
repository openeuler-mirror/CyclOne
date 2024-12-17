package route

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/voidint/binding"

	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/server/cloudbootserver/service"
	"idcos.io/cloudboot/utils"
	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
	"idcos.io/cloudboot/utils/upload"
)

// SaveCollectedDevice 保存采集到的物理机信息
func SaveCollectedDevice(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	reqData := new(service.CollectedDevice)
	//reqData.SN = chi.URLParam(r, "sn")
	reqData.OriginNode = myhttp.ExtractOriginNodeWithDefault(r, "master")
	reqData.OriginNodeIP = myhttp.ExtractOriginNodeIP(r)
	if binding.Json(r, reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	log.Debugf("%s ==> %s: %s, %s:%s",
		reqData.SN,
		myhttp.XForwardedOrigin,
		r.Header.Get(myhttp.XForwardedOrigin),
		myhttp.XForwardedFor,
		r.Header.Get(myhttp.XForwardedFor),
	)

	if err := service.SaveCollectedDevice(log, repo, reqData); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// UploadDevices 文件导入文件上传
func UploadDevices(w http.ResponseWriter, r *http.Request) {
	//预防乱码
	w.Header().Add("Content-type", "text/html; charset=utf-8")

	//解析并生成临时文件，为后续的工作做准备
	filename, err := upload.GenerateTempFile(r, "device")
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"result": filename,
	}))
}

// ImportDevicesPreview 导入设备文件预览
func ImportDevicesPreview(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.ImportPreviewReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	result, err := service.ImportDevicesPreview(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", result))
}

//ImportDevices 导入物理机
func ImportDevices(w http.ResponseWriter, r *http.Request) {
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

	err := service.ImportDevices(log, repo, conf, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// UploadDevices2Store 文件导入文件上传
func UploadDevices2Store(w http.ResponseWriter, r *http.Request) {
	//预防乱码
	w.Header().Add("Content-type", "text/html; charset=utf-8")

	//解析并生成临时文件，为后续的工作做准备
	filename, err := upload.GenerateTempFile(r, "store-device")
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"result": filename,
	}))
}

// ImportDevices2StorePreview 导入设备文件预览
func ImportDevices2StorePreview(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.ImportPreviewReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	result, err := service.ImportDevices2StorePreview(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", result))
}

//ImportDevices2Store 导入物理机到库房的货架上
func ImportDevices2Store(w http.ResponseWriter, r *http.Request) {
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

	err := service.ImportDevices2Store(log, repo, conf, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// UploadStockDevices 存量物理机导入文件上传
func UploadStockDevices(w http.ResponseWriter, r *http.Request) {
	//预防乱码
	w.Header().Add("Content-type", "text/html; charset=utf-8")

	//解析并生成临时文件，为后续的工作做准备
	filename, err := upload.GenerateTempFile(r, "stock-device")
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"result": filename,
	}))
}

// ImportStockDevicesPreview 导入设备文件预览
func ImportStockDevicesPreview(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.ImportPreviewReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	result, err := service.ImportStockDevicesPreview(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", result))
}

//ImportStockDevices 导入物理机
func ImportStockDevices(w http.ResponseWriter, r *http.Request) {
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

	err := service.ImportStockDevices(log, repo, conf, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// GetDevicePage 查询物理机分页列表
func GetDevicePage(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())

	req := new(service.DevicePageReq)
	if binding.Form(r, req).CustomHandle(HandleValidateErrs, w) {
		return
	}

	pg, err := service.GetDevicePage(log, repo, conf, req)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(pg)),
	)
}

// GetDeviceBySN 根据sn查询采集到的设备信息
func GetDeviceBySN(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	dev, err := service.GetDeviceBySN(log, repo, chi.URLParam(r, "sn"))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(dev)),
	)
}

// GetDeviceQuerys 设备的查询（过滤）参数列表
func GetDeviceQuerys(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	//|idc|server_room|server_cabinet|physical_area|op_status|usage|category|vendor
	p, err := service.GetDeviceQuerys(log, repo, chi.URLParam(r, "param_name"))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(p)),
	)
}

// GetCombinedDeviceBySN 根据sn查询设备信息及其若干装机参数
func GetCombinedDeviceBySN(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())

	dev, err := service.GetCombinedDeviceBySN(repo, conf, chi.URLParam(r, "sn"))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(dev)),
	)
}

// UpdateDevice 更改设备信息
func UpdateDevice(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	req := new(service.UpdateDevicesReq)
	if binding.Bind(r, req).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if user != nil {
		req.LoginName = user.LoginName
	}

	dev, err := service.UpdateDeviceBySN(log, repo, conf, req)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{"id": dev.ID}),
	)
}

// UpdateDeviceOperationStatus 更改设备信息操作状态
func UpdateDeviceOperationStatus(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	req := new(service.UpdateDevicesOperationStatusReq)
	if binding.Bind(r, req).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if user != nil {
		req.LoginName = user.LoginName
	}

	dev, err := service.UpdateDeviceOperationStatusBySN(log, repo, conf, req)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{"id": dev.ID}),
	)
}

// UpdateDevices 批量更改设备信息
func UpdateDevices(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	req := new(service.BatchUpdateDevicesReq)
	if binding.Bind(r, req).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if user != nil {
		req.LoginName = user.LoginName
	}

	affected, err := service.BatchUpdateDevices(log, repo, conf, req)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{"affected": affected}),
	)
}

// DeleteDevices 批量删除设备
func DeleteDevices(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	req := new(service.DeleteDevicesReq)
	if binding.Bind(r, req).CustomHandle(HandleValidateErrs, w) {
		return
	}
	totalAffected, err := service.DeleteDevices(log, repo, req)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{"ids": req.IDs, "sns": req.SNs, "total_affected": totalAffected}),
	)
}

// GetDevicePageByTor 根据指定的TOR查询物理机分页列表
func GetDevicePageByTor(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())

	req := new(service.DevicePageByTorReq)
	if binding.Form(r, req).CustomHandle(HandleValidateErrs, w) {
		return
	}

	pg, err := service.GetDeviceByTorPage(log, repo, conf, req)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(pg)),
	)
}

// ExportCombinedDevices 导出设备信息详情
func ExportCombinedDevices(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())

	reqData := service.DevicePageReq{} //ExportDevicesReq{}
	if binding.Bind(r, &reqData).Handle(w) {
		return
	}

	items, err := service.GetExportDevices(log, repo, conf, &reqData)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, item := range items {
		item.OperationStatus = service.OperationStatusTransfer(item.OperationStatus, true)
		item.PowerStatus = service.ConvertPowerStatus(item.PowerStatus)
	}

	err = render.XLSX(w, utils.FileDevice, service.ExportedDevices(items).ToTableRecords())
	if err != nil {
		log.Error(err)
	}
}

/////////////////////////////特殊设备/////////////////////////////////
// UploadSpecialDevices 文件导入文件上传
func UploadSpecialDevices(w http.ResponseWriter, r *http.Request) {
	//预防乱码
	w.Header().Add("Content-type", "text/html; charset=utf-8")

	//解析并生成临时文件，为后续的工作做准备
	filename, err := upload.GenerateTempFile(r, "special-device")
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"result": filename,
	}))
}

// ImportSpecialDevicesPreview 导入特殊设备文件预览
func ImportSpecialDevicesPreview(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.ImportPreviewReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	result, err := service.ImportSpecialDevicesPreview(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", result))
}

//ImportSpecialDevices 导入特殊设备
func ImportSpecialDevices(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	var reqData service.ImportPreviewReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if user != nil {
		reqData.LoginName = user.LoginName
	}

	err := service.ImportSpecialDevices(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

//SaveSpecialDevices 新增特殊设备
func SaveSpecialDevice(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	var reqData service.SpecialDeviceReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if user != nil {
		reqData.LoginName = user.LoginName
	}

	mod, err := service.SaveSpecialDevices(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{"id": mod.ID}))
}


//SaveNewDevices 新增设备
func SaveNewDevices(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())

	var reqData service.NewDevicesList
	if err := myhttp.DecodeJSON(r, &reqData); err != nil {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}
	succeedSNs, totalAffected, err := service.SaveNewDevices(log, repo, conf, &service.NewDevicesReq{
		NewDevices:    reqData,
		CurrentUser:   user,
	})
	if err != nil {
		//HandleErr(r.Context(), w, err)
		//render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		render.JSON(w, http.StatusOK,
			myhttp.NewRespBody(myhttp.Failure, "操作失败", map[string]interface{}{"succeed_sns": succeedSNs, "total_affected": totalAffected, "detail": err.Error()}),
		)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{"succeed_sns": succeedSNs, "total_affected": totalAffected, "detail":"success"}),
	)
}


//BatchMoveDevices 设备搬迁
func BatchMoveDevices(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())

	var reqData service.BatchMoveDevicesList
	if err := myhttp.DecodeJSON(r, &reqData); err != nil {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}
	succeedSNs, totalAffected, err := service.BatchMoveDevices(log, repo, conf, &service.BatchMoveDevicesReq{
		Devices:    reqData,
		CurrentUser:   user,
	})
	if err != nil {
		render.JSON(w, http.StatusOK,
			myhttp.NewRespBody(myhttp.Failure, "操作失败", map[string]interface{}{"succeed_sns": succeedSNs, "total_affected": totalAffected, "detail": err.Error()}),
		)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{"succeed_sns": succeedSNs, "total_affected": totalAffected, "detail":"success"}),
	)
}


//BatchMoveDevices 设备退役
func BatchRetireDevices(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())

	var reqData service.BatchRetireDevicesList
	if err := myhttp.DecodeJSON(r, &reqData); err != nil {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}
	succeedSNs, totalAffected, err := service.BatchRetireDevices(log, repo, conf, &service.BatchRetireDevicesReq{
		SNs:    	reqData,
		CurrentUser:   user,
	})
	if err != nil {
		render.JSON(w, http.StatusOK,
			myhttp.NewRespBody(myhttp.Failure, "操作失败", map[string]interface{}{"succeed_sns": succeedSNs, "total_affected": totalAffected, "detail": err.Error()}),
		)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{"succeed_sns": succeedSNs, "total_affected": totalAffected, "detail":"success"}),
	)
}