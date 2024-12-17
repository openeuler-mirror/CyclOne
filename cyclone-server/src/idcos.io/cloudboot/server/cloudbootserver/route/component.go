package route

import (
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"

	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/server/cloudbootserver/service"
	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
)

// SavePEConfigLog 保存peconfig组件日志
func SavePEConfigLog(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	//
	conf, _ := middleware.ConfigFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var err error
	var reqData service.SaveComponentLogReq
	reqData.SN = chi.URLParam(r, "sn")
	reqData.Component = model.ComponentPEConfig
	reqData.LogData, err = ioutil.ReadAll(r.Body)
	reqData.DataPath = conf.Server.StorageRootDir
	reqData.OriginNode = myhttp.ExtractOriginNodeWithDefault(r, "master")
	reqData.OriginNodeIP = myhttp.ExtractOriginNodeIP(r)

	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	if err = service.SaveComponentLog(repo, log, &reqData); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// SaveHWServerLog 保存hw-server组件日志
func SaveHWServerLog(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	//
	conf, _ := middleware.ConfigFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var err error
	var reqData service.SaveComponentLogReq
	reqData.SN = chi.URLParam(r, "sn")
	reqData.Component = model.ComponentHWServer
	reqData.LogData, err = ioutil.ReadAll(r.Body)
	reqData.DataPath = conf.Server.StorageRootDir
	reqData.OriginNode = myhttp.ExtractOriginNodeWithDefault(r, "master")
	reqData.OriginNodeIP = myhttp.ExtractOriginNodeIP(r)

	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	if err = service.SaveComponentLog(repo, log, &reqData); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// SaveCloudbootAgentLog 保存cloudboot-agent组件日志
func SaveCloudbootAgentLog(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	//
	conf, _ := middleware.ConfigFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var err error
	var reqData service.SaveComponentLogReq
	reqData.SN = chi.URLParam(r, "sn")
	reqData.Component = model.ComponentAgent
	reqData.LogData, err = ioutil.ReadAll(r.Body)
	reqData.DataPath = conf.Server.StorageRootDir
	reqData.OriginNode = myhttp.ExtractOriginNodeWithDefault(r, "master")
	reqData.OriginNodeIP = myhttp.ExtractOriginNodeIP(r)

	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	if err = service.SaveComponentLog(repo, log, &reqData); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// SaveWinConfigLog 保存winconfig组件日志
func SaveWinConfigLog(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var err error
	var reqData service.SaveComponentLogReq
	reqData.SN = chi.URLParam(r, "sn")
	reqData.Component = model.ComponentWINConfig
	reqData.LogData, err = ioutil.ReadAll(r.Body)
	reqData.DataPath = conf.Server.StorageRootDir
	reqData.OriginNode = myhttp.ExtractOriginNodeWithDefault(r, "master")
	reqData.OriginNodeIP = myhttp.ExtractOriginNodeIP(r)

	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	if err = service.SaveComponentLog(repo, log, &reqData); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// SaveOSConfigLog 保存系统配置日志
func SaveOSConfigLog(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var err error
	var reqData service.SaveComponentLogReq
	reqData.SN = chi.URLParam(r, "sn")
	reqData.Component = model.OSConfigLog
	reqData.LogData, err = ioutil.ReadAll(r.Body)
	reqData.DataPath = conf.Server.StorageRootDir
	reqData.OriginNode = myhttp.ExtractOriginNodeWithDefault(r, "master")
	reqData.OriginNodeIP = myhttp.ExtractOriginNodeIP(r)

	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	if err = service.SaveComponentLog(repo, log, &reqData); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}

// SaveImageCloneLog 保存镜像制作日志
func SaveImageCloneLog(w http.ResponseWriter, r *http.Request) {
	repo, _ := middleware.RepoFromContext(r.Context())
	conf, _ := middleware.ConfigFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	var err error
	var reqData service.SaveComponentLogReq
	reqData.SN = chi.URLParam(r, "sn")
	reqData.Component = model.ImageCloneLog
	reqData.LogData, err = ioutil.ReadAll(r.Body)
	reqData.DataPath = conf.Server.StorageRootDir
	reqData.OriginNode = myhttp.ExtractOriginNodeWithDefault(r, "master")
	reqData.OriginNodeIP = myhttp.ExtractOriginNodeIP(r)

	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	if err = service.SaveComponentLog(repo, log, &reqData); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}
	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}
