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
)

// SaveStoreRoom 保存(新增/修改)库房
func SaveStoreRoom(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	var reqData service.SaveStoreRoomReq
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	if user != nil {
		reqData.LoginName = user.LoginName
	}

	err := service.SaveStoreRoom(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"id": reqData.ID,
	}))
}

// RemoveStoreRoom 移除库房
func RemoveStoreRoom(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error(err)
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}

	//删除前校验，存在子数据不允许删除
	msg := service.RemoveStoreRoomValidate(log, repo, uint(id))
	if msg != "" {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(msg))
		return
	}

	err = service.RemoveStoreRoomByID(log, repo, uint(id))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// GetStoreRooms 查询机房分页列表
func GetStoreRooms(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.GetStoreRoomPageReq
	if binding.Form(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	pg, err := service.GetStoreRoomWithPage(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(pg)),
	)
}

// GetStoreRoom 查询指定ID的机房信息详情
func GetStoreRoom(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error(err)
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}

	sr, err := service.GetStoreRoomByID(log, repo, uint(id))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(sr)),
	)
}

//UploadStoreRoom 加载导入网络区域文件
func UploadStoreRoom(w http.ResponseWriter, r *http.Request) {
	//预防乱码
	w.Header().Add("Content-type", "text/html; charset=utf-8")

	//解析并生成临时文件，为后续的工作做准备
	filename, err := upload.GenerateTempFile(r, "store_room")
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"result": filename,
	}))

	return
}

//ImportStoreRoomPriview 导入库房文件预览
func ImportStoreRoomPriview(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData upload.ImportReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}

	result, err := service.ImportStoreRoomPriview(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", result))
}

//ImportStoreRoom  导入库房文件
func ImportStoreRoom(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	var reqData upload.ImportReq
	if binding.Json(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	reqData.UserName = user.LoginName

	err := service.ImportStoreRoom(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// 虚拟货架
// SaveVirtualCabinet 保存(新增/修改)虚拟货架
func SaveVirtualCabinet(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())
	user, _ := middleware.LoginUserFromContext(r.Context())

	var reqData service.SaveVirtualCabinetReq
	if binding.Bind(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	if user != nil {
		reqData.LoginName = user.LoginName
	}

	err := service.SaveVirtualCabinet(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.NewRespBody(myhttp.Success, "操作成功", map[string]interface{}{
		"id": reqData.ID,
	}))
}

// RemoveVirtualCabinet 移除库房
func RemoveVirtualCabinet(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error(err)
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}

	//删除前校验，存在子数据不允许删除
	valid, err := service.RemoveVirtualCabinetValidate(log, repo, uint(id))
	if !valid {
		render.JSON(w, http.StatusBadRequest, myhttp.ErrRespBody(err.Error()))
		return
	}

	err = service.RemoveVirtualCabinetByID(log, repo, uint(id))
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// GetVirtualCabinets 查询机房分页列表
func GetVirtualCabinets(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var reqData service.GetVirtualCabinetPageReq
	if binding.Form(r, &reqData).CustomHandle(HandleValidateErrs, w) {
		return
	}
	pg, err := service.GetVirtualCabinetWithPage(log, repo, &reqData)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(pg)),
	)
}
