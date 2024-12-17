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
	"idcos.io/cloudboot/utils/upload"
	myuser "idcos.io/cloudboot/utils/user"
)

// SubmitIDCAbolishApproval 提交数据中心裁撤审批
func SubmitIDCAbolishApproval(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	var token string
	if user != nil {
		token = user.Token
	}

	reqData := service.SubmitIDCAbolishApprovalReq{
		SubmitIDsApprovalReq: service.SubmitIDsApprovalReq{
			SubmitApprovalCommon: service.SubmitApprovalCommon{
				CurrentUser: user,
				GetEmailFromUAM: myuser.GetEmailByUUID(log, token, conf.UAM.RootEndpoint),
			},
		},
	}

	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	id, err := service.SubmitIDCAbolishApproval(log, repo, conf, &reqData)
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

// SubmitServerRoomAbolishApproval 提交机房管理单元裁撤审批
func SubmitServerRoomAbolishApproval(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	var token string
	if user != nil {
		token = user.Token
	}

	reqData := service.SubmitServerRoomAbolishApprovalReq{
		SubmitIDsApprovalReq: service.SubmitIDsApprovalReq{
			SubmitApprovalCommon: service.SubmitApprovalCommon{
				CurrentUser: user,
				GetEmailFromUAM: myuser.GetEmailByUUID(log, token, conf.UAM.RootEndpoint),
			},
		},
	}
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	id, err := service.SubmitServerRoomAbolishApproval(log, repo, conf, &reqData)
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

//SubmitNetAreaOfflineApproval 提交网络区域下线审批
func SubmitNetAreaOfflineApproval(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())
	var token string
	if user != nil {
		token = user.Token
	}
	reqData := service.SubmitNetAreaOfflineApprovalReq{
		SubmitIDsApprovalReq: service.SubmitIDsApprovalReq{
			SubmitApprovalCommon: service.SubmitApprovalCommon{
				CurrentUser: user,
				GetEmailFromUAM: myuser.GetEmailByUUID(log, token, conf.UAM.RootEndpoint),
			},
		},
	}
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	id, err := service.SubmitNetAreaOfflineApproval(log, repo, conf, &reqData)
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

// SubmitIPUnassignApproval IP回收（取消分配）
func SubmitIPUnassignApproval(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	var token string
	if user != nil {
		token = user.Token
	}

	reqData := service.SubmitIPUnassignApprovalReq{
		SubmitIDsApprovalReq: service.SubmitIDsApprovalReq{
			SubmitApprovalCommon: service.SubmitApprovalCommon{
				CurrentUser: user,
				GetEmailFromUAM: myuser.GetEmailByUUID(log, token, conf.UAM.RootEndpoint),
			},
		},
	}
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	id, err := service.SubmitIPUnassignApproval(log, repo, conf, &reqData)
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

// SubmitDevicePowerOffApproval 物理机关机
func SubmitDevicePowerOffApproval(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())
	
	var token string
	if user != nil {
		token = user.Token
	}

	reqData := service.SubmitDevicePowerOffApprovalReq{
		SubmitDeviceRetirementApprovalReq: service.SubmitDeviceRetirementApprovalReq{
			SubmitApprovalCommon: service.SubmitApprovalCommon{
				CurrentUser: user,
				GetEmailFromUAM: myuser.GetEmailByUUID(log, token, conf.UAM.RootEndpoint),
			},
		},
	}
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	id, err := service.SubmitDevicePowerOffApproval(log, repo, conf, &reqData)
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

// SubmitDeviceRestartApproval 物理机重启
func SubmitDeviceRestartApproval(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	var token string
	if user != nil {
		token = user.Token
	}

	reqData := service.SubmitDevicePowerOffApprovalReq{
		SubmitDeviceRetirementApprovalReq: service.SubmitDeviceRetirementApprovalReq{
			SubmitApprovalCommon: service.SubmitApprovalCommon{
				CurrentUser: user,
				GetEmailFromUAM: myuser.GetEmailByUUID(log, token, conf.UAM.RootEndpoint),
			},
		},
	}
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	id, err := service.SubmitDeviceRestartApproval(log, repo, conf, &reqData)
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

// SubmitCabinetOfflineApproval 提交机架下线审批
func SubmitCabinetOfflineApproval(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	var token string
	if user != nil {
		token = user.Token
	}
	
	reqData := service.SubmitCabinetOfflineApprovalReq{
		SubmitApprovalCommon: service.SubmitApprovalCommon{
			CurrentUser: user,
			GetEmailFromUAM: myuser.GetEmailByUUID(log, token, conf.UAM.RootEndpoint),
		},
	}
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	id, err := service.SubmitCabinetOfflineApproval(log, repo, conf, &reqData)
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

// SubmitCabinetPowerOffApproval 提交机架关电审批
func SubmitCabinetPowerOffApproval(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	var token string
	if user != nil {
		token = user.Token
	}
	
	reqData := service.SubmitCabinetPowerOffApprovalReq{
		SubmitApprovalCommon: service.SubmitApprovalCommon{
			CurrentUser: user,
			GetEmailFromUAM: myuser.GetEmailByUUID(log, token, conf.UAM.RootEndpoint),
		},
	}
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	id, err := service.SubmitCabinetPowerOffApproval(log, repo, conf, &reqData)
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

// SubmitDeviceMigrationApproval 提交物理机搬迁审批
func SubmitDeviceMigrationApproval(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	var token string
	if user != nil {
		token = user.Token
	}

	reqData := service.SubmitDeviceMigrationApprovalReq{
		SubmitApprovalCommon: service.SubmitApprovalCommon{
			CurrentUser: user,
			GetEmailFromUAM: myuser.GetEmailByUUID(log, token, conf.UAM.RootEndpoint),
		},
	}
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	id, err := service.SubmitDeviceMigrationApproval(log, repo, conf, &reqData)
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

// UploadMigrationApproval 文件导入文件上传
func UploadMigrationApproval(w http.ResponseWriter, r *http.Request) {
	//预防乱码
	w.Header().Add("Content-type", "text/html; charset=utf-8")

	//解析并生成临时文件，为后续的工作做准备
	filename, err := upload.GenerateTempFile(r, "migration_approval_")
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"result": filename,
	}))
}

// ImportMigrationApprovalPriview 导入设备文件预览
func ImportMigrationApprovalPriview(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.ImportPreviewReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	result, err := service.ImportMigrationApprovalPriview(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", result))
}

//ImportMigrationApproval 导入待搬迁物理机
func ImportMigrationApproval(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	var reqData service.ImportApprovalReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if user != nil {
		reqData.LoginName = user.LoginName
	}

	id, err := service.ImportMigrationApproval(log, repo, conf, user, &reqData)
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

// SubmitDeviceRetirementApproval 提交物理机退役审批
func SubmitDeviceRetirementApproval(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	var token string
	if user != nil {
		token = user.Token
	}

	reqData := service.SubmitDeviceRetirementApprovalReq{
		SubmitApprovalCommon: service.SubmitApprovalCommon{
			CurrentUser: user,
			GetEmailFromUAM: myuser.GetEmailByUUID(log, token, conf.UAM.RootEndpoint),
		},
	}
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	id, err := service.SubmitDeviceRetirementApproval(log, repo, conf, &reqData)
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

// SubmitDeviceOSReInstallationApproval 提交物理机OS重装
func SubmitDeviceOSReInstallationApproval(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	var token string
	if user != nil {
		token = user.Token
	}

	reqData := service.SubmitDeviceReInstallationApprovalReq{
		SubmitApprovalCommon: service.SubmitApprovalCommon{
			CurrentUser: user,
			GetEmailFromUAM: myuser.GetEmailByUUID(log, token, conf.UAM.RootEndpoint),
		},
	}
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	id, err := service.SubmitDeviceReInstallationApproval(log, repo, conf, &reqData)
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

// SubmitDeviceRecycleApproval 提交物理机回收审批
func SubmitDeviceRecycleApproval(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	var token string
	if user != nil {
		token = user.Token
	}

	reqData := service.SubmitDeviceRecycleApprovalReq{
		SubmitDeviceRetirementApprovalReq: service.SubmitDeviceRetirementApprovalReq{
			SubmitApprovalCommon: service.SubmitApprovalCommon{
				CurrentUser: user,
				GetEmailFromUAM: myuser.GetEmailByUUID(log, token, conf.UAM.RootEndpoint),
			},
		},
	}
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	id, err := service.SubmitDeviceRecycleApproval(log, repo, conf, &reqData)
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

// Approve 审批指定的问题单
func Approve(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	lim, _ := middleware.DHCPLimiterFromContext(r.Context())

	approveID, _ := strconv.Atoi(chi.URLParam(r, "approval_id"))
	stepID, _ := strconv.Atoi(chi.URLParam(r, "approval_step_id"))
	reqData := service.ApproveReq{
		CurrentUser: user,
		ApprovalID:  uint(approveID),
		StepID:      uint(stepID),
	}
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	err := service.Approve(log, repo, conf, lim, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// GetMyApprovalPage  我发起的审批
func GetMyApprovalPage(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.GetMyApprovalPageReq
	if binding.Form(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	pg, err := service.GetMyApprovalPage(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(pg)),
	)
}

// GetApproveByMePage 获取待我审批的
func GetApproveByMePage(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.GetApproveByMePageReq
	if binding.Form(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	pg, err := service.GetApproveByMePage(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(pg)),
	)
}

// GetApprovedByMePage 获取我审批完成的
func GetApprovedByMePage(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.GetApprovedByMeReq
	if binding.Form(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	pg, err := service.GetApprovedByMePage(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(pg)),
	)
}

// RevokeApproval 取消我发起的审批单
func RevokeApproval(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	id, err := strconv.Atoi((chi.URLParam(r, "approval_id")))
	if err != nil {
		log.Error(err)
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}

	var be binding.Errors
	reqData := &service.RevokeApprovalReq{
		ID: id,
	}
	if HandleValidateErrs(reqData.Validate(r, be), w) {
		return
	}

	err = service.RevokeApproval(log, repo, reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// GetApprovalByID 查询指定ID的申请单信息详情
func GetApprovalByID(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())

	var token string
	if user != nil {
		token = user.Token
	}

	id, err := strconv.Atoi(chi.URLParam(r, "approval_id"))
	if err != nil {
		log.Error(err)
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}

	var be binding.Errors
	reqData := &service.GetApprovalByIDReq{
		ID:             id,
		GetNameFromUAM: myuser.GetUsersByUUID(log, token, conf.UAM.RootEndpoint),
	}
	if HandleValidateErrs(reqData.Validate(r, be), w) {
		return
	}

	sr, err := service.GetApprovalByID(log, repo, reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(sr)),
	)
}
