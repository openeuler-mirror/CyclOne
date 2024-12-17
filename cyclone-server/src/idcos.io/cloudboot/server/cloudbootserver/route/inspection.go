package route

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/voidint/binding"

	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/server/cloudbootserver/service"
	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
)

// GetInspectionPage 查询硬件巡检分页
func GetInspectionPage(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.GetInspectionsPageReq
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	pg, err := service.GetInspectionsPage(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(pg)),
	)
}

// GetInspectionBySN 查询指定设备某一次硬件巡检明细
func GetInspectionBySN(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())

	reqData := service.GetInspectionBySNReq{
		SN: strings.TrimSpace(chi.URLParam(r, "sn")),
	}
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	resp, err := service.GetInspectionBySN(repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success,
		"操作成功", resp))
}

// GetInspectionStartTimesBySN 查询指定设备的历史硬件巡检开始时间列表
func GetInspectionStartTimesBySN(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())

	reqData := service.GetInspectionStartTimesBySNReq{
		SN: strings.TrimSpace(chi.URLParam(r, "sn")),
	}
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	starts, err := service.GetInspectionStartTimesBySN(repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success,
		"操作成功", map[string]interface{}{
			"start_time": starts,
		}))
}

// GetInspectionStatistics 查询硬件巡检结果状态(正常、警告、异常)统计
func GetInspectionStatistics(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.GetGetInspectionStatisticsReq
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	stats, err := service.GetGetInspectionStatistics(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
			"statistics": stats,
		}),
	)
}

// GetInspectionRecordsPage 查询物理机巡检记录分页列表
func GetInspectionRecordsPage(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.InspectionRecordsPageReq
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	pg, err := service.GetInspectionRecordsPage(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(pg)),
	)
}
