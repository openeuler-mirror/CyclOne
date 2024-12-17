package middleware

import (
	"context"
	"net/http"

	"idcos.io/cloudboot/server/hwserver/config"
)

// ctxHWConfigKey 注入的*config.Configuration对应的查询Key
var ctxHWConfigKey uint8

// HWConfigFromContext 从ctx中获取*config.Configuration
func HWConfigFromContext(ctx context.Context) (*config.Configuration, bool) {
	conf, ok := ctx.Value(&ctxHWConfigKey).(*config.Configuration)
	return conf, ok
}

// InjectHWConfig 注入*config.Configuration
func InjectHWConfig(conf *config.Configuration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(
				context.WithValue(r.Context(), &ctxHWConfigKey, conf),
			))
		})
	}
}
