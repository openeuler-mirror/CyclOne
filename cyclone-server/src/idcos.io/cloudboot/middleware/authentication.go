package middleware

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"idcos.io/cloudboot/model"
	myhttp "idcos.io/cloudboot/utils/http"
	"idcos.io/cloudboot/utils/http/render"
	uam "idcos.io/uam-go"
)

const (
	// GET GET方法
	GET = http.MethodGet
	// POST POST方法
	POST = http.MethodPost
	// PUT PUT方法
	PUT = http.MethodPut
	// DELETE DELETE方法
	DELETE = http.MethodDelete
)

// APIMeta API元数据信息
type APIMeta struct {
	Method string
	URIReg *regexp.Regexp
}

// Match 返回此request是否与当前的API元数据匹配
func (api APIMeta) Match(r *http.Request) bool {
	return r.Method == api.Method && api.URIReg.MatchString(ReqURI(r))
}

// UnAuthAPIs 无需登录认证的API元信息集合
type UnAuthAPIs []APIMeta

// ShouldNotBeChecked 返回是否需要对此request请求进行登录认证的布尔值
func (items UnAuthAPIs) ShouldNotBeChecked(r *http.Request) bool {
	for i := range items {
		if items[i].Match(r) {
			return true
		}
	}
	return false
}

// ReqURI 获取请求的URL“主体“
func ReqURI(r *http.Request) string {
	reqURI := r.URL.RequestURI()
	if idx := strings.Index(reqURI, "?"); idx > 0 {
		reqURI = reqURI[:idx]
	}
	return strings.TrimSuffix(reqURI, "/")
}

var (
	// unAuthAPIs 无需登录认证的API集合
	unAuthAPIs = UnAuthAPIs([]APIMeta{
		{POST, regexp.MustCompile("^/api/cloudboot/v1/devices/collections$")},
		{GET, regexp.MustCompile("^/api/cloudboot/v1/devices/[0-9A-Za-z_]+/collections$")},

		{GET, regexp.MustCompile("^/api/cloudboot/v1/devices/[0-9A-Za-z_]+/settings/networks$")},
		{GET, regexp.MustCompile("^/api/cloudboot/v1/devices/[0-9A-Za-z_]+/settings/os-users$")},
		{GET, regexp.MustCompile("^/api/cloudboot/v1/devices/[0-9A-Za-z_]+/settings/hardwares$")},
		{GET, regexp.MustCompile("^/api/cloudboot/v1/devices/[0-9A-Za-z_]+/settings/hardwareinfo$")},
		{GET, regexp.MustCompile("^/api/cloudboot/v1/devices/[0-9A-Za-z_]+/settings/system-template$")},
		{GET, regexp.MustCompile("^/api/cloudboot/v1/devices/[0-9A-Za-z_]+/settings/image-template$")},
		{GET, regexp.MustCompile("^/api/cloudboot/v1/devices/[0-9A-Za-z_]+/settings$")},
		{GET, regexp.MustCompile("^/api/cloudboot/v1/devices/[0-9A-Za-z_]+/is-in-install-list$")},
		{POST, regexp.MustCompile("^/api/cloudboot/v1/devices/[0-9A-Za-z_]+/installations/progress$")},
		{GET, regexp.MustCompile("^/api/cloudboot/v1/devices/[0-9A-Za-z_]+/installations/status$")},

		{GET, regexp.MustCompile("^/api/cloudboot/v1/devices/[0-9A-Za-z_:%]+/pxe$")},
		{POST, regexp.MustCompile("^/api/cloudboot/v1/devices/[0-9A-Za-z_]+/pxe$")},
		{POST, regexp.MustCompile("^/api/cloudboot/v1/devices/[0-9A-Za-z_]+/centos6/uefi/pxe$")},

		{GET, regexp.MustCompile("^/api/cloudboot/v1/permissions/codes$")},
		{GET, regexp.MustCompile("^/api/cloudboot/v1/system/login/settings$")},
		{GET, regexp.MustCompile("^/api/cloudboot/v1/orders/export$")},
		{GET, regexp.MustCompile("^/api/cloudboot/v1/samba/settings$")},
	})
)

// ctxLoginUserKey 登录用户指针在上下文中的查询Key
var ctxLoginUserKey uint8

// Authenticator 用户登录认证中间件
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if unAuthAPIs.ShouldNotBeChecked(r) {
			next.ServeHTTP(w, r)
			return
		}
		user, err := whoami(r)
		if err != nil {
			render.JSON(w, http.StatusInternalServerError, myhttp.ErrRespBody(err.Error()))
			return
		}

		if user == nil {
			conf, _ := ConfigFromContext(r.Context())
			render.JSON(w, http.StatusUnauthorized,
				myhttp.NewRespBody(myhttp.Failure, http.StatusText(http.StatusUnauthorized), map[string]interface{}{
					"channel": "uam_sso",
					"url":     fmt.Sprintf("%s/login", strings.TrimSuffix(conf.UAM.RootEndpoint, "/")),
				}),
			)
			return
		}

		r.Header.Add("Authorization", user.Token)

		r = r.WithContext(context.WithValue(r.Context(), &ctxLoginUserKey, user))
		next.ServeHTTP(w, r)
	})
}

// LoginUserFromContext 从ctx中获取登录用户信息
func LoginUserFromContext(ctx context.Context) (*model.CurrentUser, bool) {
	user, ok := ctx.Value(&ctxLoginUserKey).(*model.CurrentUser)
	return user, ok
}

// whoami 返回当前用户信息
func whoami(r *http.Request) (user *model.CurrentUser, err error) {
	log, _ := LoggerFromContext(r.Context())
	conf, _ := ConfigFromContext(r.Context())

	var token string
	// 尝试从HTTP Header中获取token
	if token = r.Header.Get("Authorization"); token == "" {
		if ck, _ := r.Cookie("access-token"); ck != nil {
			token = ck.Value
		}
	}

	uamUser, err := uam.NewClient(conf.UAM.RootEndpoint, token, uam.LogOption(log)).AuthUser(model.AppID4UAM)
	if err == uam.ErrUnauthenticatedUser {
		return nil, nil // 未认证，返回空user指针。
	}
	if err != nil {
		return nil, err // 调用uam发生未知错误
	}

	if uamUser.Status == uam.UserStatusDisabled {
		return nil, nil // 帐号被禁用
	}
	return &model.CurrentUser{
		Token:           token,
		ID:              uamUser.ID,
		LoginName:       uamUser.LoginID,
		Name:            uamUser.Name,
		Email:           uamUser.Email,
		Status:          uamUser.Status,
		UAMPortalURL:    conf.UAM.PortalURL,
		PermissionCodes: filterBootPermission(uamUser),
		Tenant: &model.Tenant{
			ID:   uamUser.Tenant.ID,
			Name: uamUser.Tenant.Name,
		},
		Department: &model.Department{
			ID:   uamUser.Department.ID,
			Name: uamUser.Department.Name,
		},
	}, nil
}

//filterBootPermission 从用户权限当中过滤出来BOOT的权限
func filterBootPermission(uamUser *uam.User) []string {
	permissions := uamUser.PermissionCodes(model.MenuPermissionType)
	permissions = append(permissions, uamUser.PermissionCodes(model.ButtonPermissionType)...)
	permissions = append(permissions, uamUser.PermissionCodes(model.DataPermissionType)...)
	permissions = append(permissions, uamUser.PermissionCodes(model.APIPermissionType)...)
	return permissions
}
