package middleware

import (
	"context"
	"net/http"

	"idcos.io/cloudboot/model"
)

// ctxRepoKey 注入的model.Repo对应的查询Key
var ctxRepoKey uint8

// RepoFromContext 从ctx中获取model.Repo
func RepoFromContext(ctx context.Context) (model.Repo, bool) {
	repo, ok := ctx.Value(&ctxRepoKey).(model.Repo)
	return repo, ok
}

// InjectRepo 注入model.Repo
func InjectRepo(repo model.Repo) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), &ctxRepoKey, repo))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
