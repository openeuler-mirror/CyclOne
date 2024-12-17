package route

import (
	"net/http"

	"github.com/voidint/binding"

	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/server/cloudbootserver/service"
	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
)

// GetUserPage 按条件查询当前用户所在租户下的用户分页列表
func GetUserPage(w http.ResponseWriter, r *http.Request) {
	conf, _ := middleware.ConfigFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	reqData := service.GetUserPageReq{
		User: user,
	}
	if binding.Form(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	pg, err := service.GetUserPage(log, conf, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(pg)),
	)
}

// GetUserByToken 返回指定token的用户信息
func GetUserByToken(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.LoginUserFromContext(r.Context())

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(user)),
	)
}

// ChangeUserPassword 修改当前用户
func ChangeUserPassword(w http.ResponseWriter, r *http.Request) {
	conf, _ := middleware.ConfigFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	reqData := service.ChangeUserPasswordReq{
		User: user,
	}
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if err := service.ChangeUserPassword(log, conf, &reqData); err != nil {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}
	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}
