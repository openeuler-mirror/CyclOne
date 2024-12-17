package middleware

import (
	"context"
	"net/http"

	"idcos.io/cloudboot/job"
)

// ctxJobManagerKey 注入的job.Manager对应的查询Key
var ctxJobManagerKey uint8

// JobManagerFromContext 从ctx中获取JobManager
func JobManagerFromContext(ctx context.Context) (job.Manager, bool) {
	jmgr, ok := ctx.Value(&ctxJobManagerKey).(job.Manager)
	return jmgr, ok
}

// InjectJobManager 注入job.Manager
func InjectJobManager(jmgr job.Manager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), &ctxJobManagerKey, jmgr))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
