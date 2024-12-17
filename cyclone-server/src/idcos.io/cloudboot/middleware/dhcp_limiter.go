package middleware

import (
	"context"
	"net/http"

	"idcos.io/cloudboot/limiter"
)

// ctxDHCPLimiterKey 注入的limiter.Limiter对应的查询Key
var ctxDHCPLimiterKey uint8

// DHCPLimiterFromContext 从ctx中获取DHCPLimiter
func DHCPLimiterFromContext(ctx context.Context) (limiter.Limiter, bool) {
	lim, ok := ctx.Value(&ctxDHCPLimiterKey).(limiter.Limiter)
	return lim, ok
}

// InjectDHCPLimiter 注入limiter.Limiter
func InjectDHCPLimiter(lim limiter.Limiter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), &ctxDHCPLimiterKey, lim))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
