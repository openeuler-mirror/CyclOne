package route

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/voidint/binding"

	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/server/cloudbootserver/service"
	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
)

// GetJobByID 查询指定ID的任务
func GetJobByID(w http.ResponseWriter, r *http.Request) {
	conf, _ := middleware.ConfigFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	jobmgr, _ := middleware.JobManagerFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	job, err := service.GetJobByID(log, conf, jobmgr, &service.GetJobByIDReq{
		ID:          chi.URLParam(r, "job_id"),
		CurrentUser: user,
	})
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(job)),
	)
}

// GetJobPage 查询满足过滤条件的任务分页列表
func GetJobPage(w http.ResponseWriter, r *http.Request) {
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	jobmgr, _ := middleware.JobManagerFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	reqData := service.GetJobPageReq{
		CurrentUser: user,
	}
	if binding.Form(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	pg, err := service.GetJobPage(log, conf, jobmgr, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(pg)),
	)
}

// PauseJob 暂停运行中的目标定时任务
func PauseJob(w http.ResponseWriter, r *http.Request) {
	jobmgr, _ := middleware.JobManagerFromContext(r.Context())

	if err := service.PauseJob(jobmgr, chi.URLParam(r, "job_id")); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.SucceedRespBody("操作成功"),
	)
}

// UnpauseJob 继续已暂停的目标定时任务
func UnpauseJob(w http.ResponseWriter, r *http.Request) {
	jobmgr, _ := middleware.JobManagerFromContext(r.Context())

	if err := service.UnpauseJob(jobmgr, chi.URLParam(r, "job_id")); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.SucceedRespBody("操作成功"),
	)
}

// RemoveJob 删除非内置任务
func RemoveJob(w http.ResponseWriter, r *http.Request) {
	jobmgr, _ := middleware.JobManagerFromContext(r.Context())

	if err := service.RemoveJob(jobmgr, chi.URLParam(r, "job_id")); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.SucceedRespBody("操作成功"),
	)
}

// AddInspectionJob 新增硬件巡检任务
func AddInspectionJob(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	jobmgr, _ := middleware.JobManagerFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	reqData := service.AddInspectionJobReq{
		OriginNode: myhttp.ExtractOriginNodeWithDefault(r, "master"),
		Creator:    user.ID,
	}
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	newid, err := service.AddInspectionJob(log, repo, conf, jobmgr, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功",
		map[string]interface{}{
			"job_id": newid,
		},
	))
}
