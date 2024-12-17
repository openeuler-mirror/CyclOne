package route

import (
	"net/http"

	"github.com/go-chi/chi"

	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/server/hwserver/service"
	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
)

// AreYouOK 用于检测服务是否可用
func AreYouOK(w http.ResponseWriter, r *http.Request) {
	render.Text(w, http.StatusOK, []byte("pong"))
}

// Collect 执行设备信息采集。采集期间需阻塞，直到采集完毕并响应采集到的设备信息。
func Collect(w http.ResponseWriter, r *http.Request) {
	conf, _ := middleware.HWConfigFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	dev, err := service.CollectDevice(conf, log)
	if err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK,
		myhttp.NewRespBody(myhttp.Success, "操作成功", myhttp.DumpContent(dev)),
	)
}

// ApplySettings 硬件配置实施。
func ApplySettings(w http.ResponseWriter, r *http.Request) {
	conf, _ := middleware.HWConfigFromContext(r.Context())
	log, _ := middleware.LoggerFromContext(r.Context())

	if err := service.NewSettingWorker(conf, log, chi.URLParam(r, "sn")).Apply(); err != nil {
		HandleErr(r.Context(), w, err)
		return
	}

	render.JSON(w, http.StatusOK, myhttp.SucceedRespBody("操作成功"))
}
