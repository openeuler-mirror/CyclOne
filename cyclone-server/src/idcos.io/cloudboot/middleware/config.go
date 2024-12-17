package middleware

import (
	"context"
	"net/http"

	"idcos.io/cloudboot/config"
)

// ctxConfigKey 注入的*config.Config对应的查询Key
var ctxConfigKey uint8

// ConfigFromContext 从ctx中获取model.Repo
func ConfigFromContext(ctx context.Context) (*config.Config, bool) {
	conf, ok := ctx.Value(&ctxConfigKey).(*config.Config)
	return conf, ok
}

// InjectConfig 注入*config.Config
func InjectConfig(conf *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(
				context.WithValue(r.Context(), &ctxConfigKey, conf),
			))
		})
	}
}
