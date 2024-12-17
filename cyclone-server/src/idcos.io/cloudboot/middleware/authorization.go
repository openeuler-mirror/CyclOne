package middleware

import (
	"fmt"
	"idcos.io/cloudboot/model"
	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
	"net/http"
	"regexp"
)

// AuthorizationAPI 鉴权API
type AuthorizationAPI struct {
	API   *APIMeta
	Codes []string
}

// AuthorizationAPIs 待鉴权的API元信息集合
type AuthorizationAPIs []*AuthorizationAPI

// ShouldBeChecked 返回是否需要对此request请求进行鉴权的布尔值与权限码集合
func (items AuthorizationAPIs) ShouldBeChecked(r *http.Request) (yes bool, codes []string) {
	if user, _ := LoginUserFromContext(r.Context()); user == nil {
		return false, nil // 未经过登录认证的API无须鉴权
	}
	for i := range items {
		if items[i].API.Match(r) {
			return true, items[i].Codes
		}
	}
	return false, nil
}

// 待鉴权API集合
var authorizationAPIs AuthorizationAPIs

// InitAuthorizationAPIs 初始化鉴权API
func InitAuthorizationAPIs(repo model.Repo) error {
	items, err := repo.GetSystemSetting4AuthorizationAPIs()
	if err != nil {
		return fmt.Errorf("unable to load system setting 'authorization': %s", err.Error())
	}

	authAPIs := make([]*AuthorizationAPI, 0, len(items))
	for i := range items {
		reg, err := regexp.CompilePOSIX(items[i].API.URIReg)
		if err != nil {
			return fmt.Errorf("invalid regex(%s) expression in 'authorization' system setting", items[i].API.URIReg)
		}
		authAPIs = append(authAPIs, &AuthorizationAPI{
			API: &APIMeta{
				Method: items[i].API.Method,
				URIReg: reg,
			},
			Codes: items[i].Codes,
		})
	}
	authorizationAPIs = AuthorizationAPIs(authAPIs)
	return nil
}

// Authorization 用户鉴权中间件
func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shouldBeCheck, codes := authorizationAPIs.ShouldBeChecked(r)
		if !shouldBeCheck {
			next.ServeHTTP(w, r)
			return
		}

		user, _ := LoginUserFromContext(r.Context())

		// 当前用户只要拥有目标权限码中的任意一个，该用户即为合法用户。
		if !user.Allow(codes...) {
			render.JSON(w, http.StatusForbidden, myhttp.NewRespBody(myhttp.Failure, "用户权限不足", map[string]interface{}{
				"allowed_permission_codes": codes,
				"user_permission_codes":    user.PermissionCodes,
			}))
			return
		}
		next.ServeHTTP(w, r)
	})
}
