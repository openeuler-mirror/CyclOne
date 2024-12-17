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
)

// GetOOBUserBySN 根据sn查询带外账户信息
func GetOOBUserBySN(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())

	user, err := service.GetOOBInfoBySn(log, repo, conf, chi.URLParam(r, "sn"))
	if err != nil {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(user)),
	)
}

// OOBPowerOn 带外管理批量开机
func OOBPowerOn(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())

	var reqData service.OOBPowerBatchOperateReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	output, err := service.BatchOperateOOBPower(log, repo, service.PowerOn, conf, false, reqData.Sns)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
			"output": output,
		}),
	)
}

// OOBPowerOff 带外管理批量关机
func OOBPowerOff(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())

	var reqData service.OOBPowerBatchOperateReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	output, err := service.BatchOperateOOBPower(log, repo, service.PowerOff, conf, false, reqData.Sns)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
			"output": output,
		}),
	)
}

// OOBPowerRestart 带外管理批量重启
func OOBPowerRestart(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())

	var reqData service.OOBPowerBatchOperateReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	output, err := service.BatchOperateOOBPower(log, repo, service.PowerRestart, conf, false, reqData.Sns)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
			"output": output,
		}),
	)
}

// OOBPowerPxeRestart 带外管理批量PXE重启
func OOBPowerPxeRestart(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())

	var reqData service.OOBPowerBatchOperateReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	output, err := service.BatchOperateOOBPower(log, repo, service.PowerRestart, conf, true, reqData.Sns)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
			"output": output,
		}),
	)
}

// DevicePowerStatus 查看设备电源状态
func DevicePowerStatus(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())

	sn := chi.URLParam(r, "sn")

	output, err := service.GetDevicePowerStatusBySN(log, repo, conf, sn)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
			"power_status": output,
		}),
	)
}

//UpdateOOBPasswordBySN 根据sn更改带外密码
func UpdateOOBPasswordBySN(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())

	req := new(service.UpdateOOBPasswordReq)
	req.SN = chi.URLParam(r, "sn")

	if binding.Bind(r, req).CustomHandle(HandleValidateErrs, w) {
		return
	}
	//req.Remark = model.OOBHistoryRemarkManu
	_, err := service.UpdateOOBPasswordBySN(log, repo, req, conf)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}
	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

//ReAccessOOB 批量重新纳管带外，检查带外是否通
func ReAccessOOB(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	req := service.ReAccessOOBReq{}
	if binding.Bind(r, &req).CustomHandle(HandleValidateErrs, w) {
		return
	}
	err := service.ReAccessOOB(log, repo, conf, &req)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}
	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

//ExportOOB 导出带外信息
func ExportOOB(w http.ResponseWriter, r *http.Request) {
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

	err = render.XLSX(w, utils.FileOOB, service.ExportOOBInfo(items).ToTableRecords())
	if err != nil {
		log.Error(err)
	}
}

// OOBInspectionOperate 带外巡检操作
func OOBInspectionOperate(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())

	var reqData service.OOBInspectionOperateReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	output, err := service.OOBInspectionOperate(log, repo, conf, &reqData)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(output)),
	)
}