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

// ReturnLimiterTokenBySN 归还目标设备的限流器令牌
func ReturnLimiterTokenBySN(w http.ResponseWriter, r *http.Request) {
	log, _ := middleware.LoggerFromContext(r.Context())
	repo, _ := middleware.RepoFromContext(r.Context())
	lim, _ := middleware.DHCPLimiterFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())

	//if !conf.DHCPLimiter.Enable {
	//	render.JSON(w, http.StatusAccepted, myhttp.ErrRespBody("限流器已关闭，请修改限流器配置后再重试。")) // TODO 寻找一个更合适的状态码
	//	return
	//}

	reqData := service.ReturnLimiterTokenBySNReq{
		SN: chi.URLParam(r, "sn"),
	}

	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	if err := service.ReturnLimiterTokenBySN(log, repo, lim, conf, &reqData); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// ReturnLimiterTokens 归还指定设备列表的限流器令牌
func ReturnLimiterTokens(w http.ResponseWriter, r *http.Request) {
	log, _ := middleware.LoggerFromContext(r.Context())
	lim, _ := middleware.DHCPLimiterFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	repo, _ := middleware.RepoFromContext(r.Context())

	//if !conf.DHCPLimiter.Enable {
	//	render.JSON(w, http.StatusAccepted, myhttp.ErrRespBody("限流器已关闭，请修改限流器配置后再重试。")) // TODO 寻找一个更合适的状态码
	//	return
	//}

	reqData := service.ReturnLimiterTokensReq{}

	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	if succeed, err := service.ReturnLimiterTokens(log, repo, lim, conf, &reqData); err != nil || succeed == 0 {
		render.JSON(w, http.StatusOK, myhttp.ErrRespBody("操作失败"))
		return
	} else if succeed != len(reqData.Tokens) {
		render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("部分成功"))
		return
	}
	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}
