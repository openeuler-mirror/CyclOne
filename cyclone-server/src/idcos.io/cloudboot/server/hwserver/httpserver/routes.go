package httpserver

import (
	"github.com/go-chi/chi"

	"idcos.io/cloudboot/server/hwserver/httpserver/route"
)

func registerHandlers(mux *chi.Mux) {
	mux.Get("/api/cloudboot/hw/v1/ping", route.AreYouOK)
	mux.Get("/api/cloudboot/hw/v1/devices/collections", route.Collect)                    // 触发设备信息采集并返回设备信息
	mux.Post("/api/cloudboot/hw/v1/devices/{sn}/settings/applyings", route.ApplySettings) // 触发硬件配置实施
}
