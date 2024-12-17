package uam

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	resty "gopkg.in/resty.v1"
)

const (
	// UserStatusEnabled 用户帐号状态-可用（激活）
	UserStatusEnabled = "ENABLED"
	// UserStatusDisabled 用户帐号状态-禁用
	UserStatusDisabled = "DISABLED"
)

// Tenant 租户
type Tenant struct {
	ID   string `json:"tenantId"`
	Name string `json:"tenantName"`
}

// Department 部门
type Department struct {
	ID   string `json:"deptId"`
	Name string `json:"deptName"`
}

// User 用户
type User struct {
	Tenant
	Department
	Token       string              `json:"token"`
	ID          string              `json:"id"`
	LoginID     string              `json:"loginId"`
	Name        string              `json:"name"`
	Email       string              `json:"email"`
	IsActive    string              `json:"isActive"`
	Status      string              `json:"status"`
	Permissions map[string][]string `json:"permissions"`
	Roles       []string            `json:"roleIds"`
	Groups      []string            `json:"userGroups"`
}

// PermissionCodes 返回指定'权限资源类型'的权限码
func (u *User) PermissionCodes(category string) []string {
	if len(u.Permissions) <= 0 {
		return []string{}
	}
	codes, ok := u.Permissions[category]
	if !ok {
		return []string{}
	}
	return codes
}

var (
	// ErrUnauthenticatedUser 未认证用户错误
	ErrUnauthenticatedUser = errors.New("[UAM]unauthenticated user")

	// ErrCallUAM 调用UAM发生非预期错误
	ErrCallUAM = errors.New("[UAM]call UAM error")
)

// AuthUser 返回指定系统下指定Token的用户信息。若Token非法，则返回ErrUnauthorizedUser错误。
func (cli *Client) AuthUser(appID string) (*User, error) {
	url := fmt.Sprintf("%s/rbac/api/authInfo", cli.conf.RootEndpoint)
	cli.log.Debugf("Fetch user ==> GET %s?appId=%s&token=***", url, appID)

	var respData struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Errors  string `json:"errors"`
		User    User   `json:"content"`
	}

	resp, err := resty.New().NewRequest().
		SetHeader(cli.tokenHeader()).
		SetHeader(cli.acceptJSONHeader()).
		SetQueryParam(appIDParam, appID).
		SetQueryParam(tokenParam, cli.conf.Token).
		SetResult(&respData).
		Get(url)

	if err != nil {
		cli.log.Error(err)
		return nil, wrapErr(err)
	}

	if resp.StatusCode() == http.StatusBadRequest || resp.StatusCode() == http.StatusUnauthorized {
		cli.log.Warnf("Fetch user ==> GET %s?appId=%s&token=***\nHTTP Status Code: %d\nResponse body: %s", url, appID, resp.StatusCode(), resp.Body())
		return nil, ErrUnauthenticatedUser
	}
	if resp.StatusCode() != http.StatusOK {
		cli.log.Warnf("Fetch user ==> GET %s?appId=%s&token=***\nHTTP Status Code: %d\nResponse body: %s", url, appID, resp.StatusCode(), resp.Body())
		return nil, ErrCallUAM
	}

	if respData.Status != success {
		cli.log.Warnf("Fetch user ==> GET %s?appId=%s&token=***\nHTTP Status Code: %d\nResponse body: %s", url, appID, resp.StatusCode(), resp.Body())
		return nil, ErrUnauthenticatedUser // uam封装不友好，若token格式非法，如token=123，会响应200状态码。
	}
	return &respData.User, nil
}

// ChangeUserPassword 修改用户密码。
// userID-用户ID，必选。
// oldPwd-旧密码，必选。
// newPwd-新密码，必选。
func (cli *Client) ChangeUserPassword(userID, oldPwd, newPwd string) error {
	url := fmt.Sprintf("%s/rbac/api/modifyPW", cli.conf.RootEndpoint)
	cli.log.Debugf("Change user password ==> POST %s", url)

	var respData struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	resp, err := resty.New().NewRequest().
		SetHeader(cli.tokenHeader()).
		SetHeader(cli.contentTypeJSONHeader()).
		SetHeader(cli.acceptJSONHeader()).
		SetBody(map[string]string{
			"userId":      userID,
			"oldPassword": oldPwd,
			"newPassword": newPwd,
		}).
		SetResult(&respData).
		Post(url)

	if err != nil {
		cli.log.Error(err)
		return wrapErr(err)
	}

	if resp.StatusCode() == http.StatusUnauthorized {
		cli.log.Warnf("Change user password ==> POST %s\nHTTP Status Code: %d\nResponse body: %s", url, resp.StatusCode(), resp.Body())
		return ErrUnauthenticatedUser
	}

	if respData.Status != success {
		cli.log.Warnf("Change user password ==> POST %s\nHTTP Status Code: %d\nResponse body: %s", url, resp.StatusCode(), resp.Body())
		return wrapErrString(respData.Message)
	}
	return nil
}

// GetUsers 查询当前租户下满足过滤条件的用户列表。
// deptID-部门ID，可选。
// name-用户名，可选。
// limit-分页限制，可选。默认pageNo为1，pageSize为10。
func (cli *Client) GetUsers(deptID, name string, limit *PageLimiter) (pageNo, pageSize, totalPages, totalRecords int64, users []*User, err error) {
	if limit == nil {
		limit = genDefaultPageLimiter()
	}
	url := fmt.Sprintf("%s/rbac/api/pageList?tenantId=%s&deptId=%s&name=%s&pageNo=%d&pageSize=%d",
		cli.conf.RootEndpoint, cli.conf.TenantID, deptID, url.QueryEscape(name), limit.PageNo, limit.PageSize)

	cli.log.Debugf("Fetch users ==> GET %s", url)

	var respData struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Content struct {
			PageNo       int64   `json:"pageNo"`
			PageSize     int64   `json:"pageSize"`
			TotalPages   int64   `json:"totalPage"`
			TotalRecords int64   `json:"totalCount"`
			Records      []*User `json:"list"`
		} `json:"content"`
	}

	resp, err := resty.New().NewRequest().
		SetHeader(cli.tokenHeader()).
		SetHeader(cli.acceptJSONHeader()).
		SetResult(&respData).
		Get(url)

	if err != nil {
		cli.log.Error(err)
		return pageNo, pageSize, totalPages, totalRecords, nil, wrapErr(err)
	}

	if resp.StatusCode() == http.StatusUnauthorized {
		cli.log.Warnf("Fetch users ==> GET %s\nHTTP Status Code: %d\nResponse body: %s", url, resp.StatusCode(), resp.Body())
		return pageNo, pageSize, totalPages, totalRecords, nil, ErrUnauthenticatedUser
	}

	if respData.Status != success {
		cli.log.Warnf("Fetch users ==> GET %s\nHTTP Status Code: %d\nResponse body: %s", url, resp.StatusCode(), resp.Body())
		return pageNo, pageSize, totalPages, totalRecords, nil, wrapErrString(respData.Message)
	}

	return respData.Content.PageNo,
		respData.Content.PageSize,
		respData.Content.TotalPages,
		respData.Content.TotalRecords,
		respData.Content.Records, nil
}

// ErrUserNotFound 用户不存在
var ErrUserNotFound = errors.New("[UAM]user not found")

// GetUserByID 返回指定ID的用户。
// id-用户ID，必选。
func (cli *Client) GetUserByID(id string) (use *User, err error) {
	url := fmt.Sprintf("%s/rbac/api/account/id?id=%s", cli.conf.RootEndpoint, id)

	cli.log.Debugf("Fetch user ==> GET %s", url)

	var respData struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Content User   `json:"content"`
	}

	resp, err := resty.New().NewRequest().
		SetHeader(cli.tokenHeader()).
		SetHeader(cli.acceptJSONHeader()).
		SetResult(&respData).
		Get(url)

	if err != nil {
		cli.log.Error(err)
		return nil, wrapErr(err)
	}

	if resp.StatusCode() == http.StatusUnauthorized {
		cli.log.Warnf("Fetch user ==> GET %s\nHTTP Status Code: %d\nResponse body: %s", url, resp.StatusCode(), resp.Body())
		return nil, ErrUnauthenticatedUser
	}

	if respData.Status != success {
		cli.log.Warnf("Fetch user ==> GET %s\nHTTP Status Code: %d\nResponse body: %s", url, resp.StatusCode(), resp.Body())
		return nil, ErrUserNotFound // UAM不能明确区分'用户不存在错误'与其他错误(如数据库错误)
	}

	return &respData.Content, nil
}

// GetUserByLoginID 返回当前租户下指定登录名的用户。
// loginid-用户登录名，必选。
func (cli *Client) GetUserByLoginID(loginID string) (use *User, err error) {
	url := fmt.Sprintf("%s/rbac/api/account/getByAccountNo?tenantId=%s&accountNo=%s", cli.conf.RootEndpoint, cli.conf.TenantID, loginID)

	cli.log.Debugf("Fetch user ==> GET %s", url)

	var respData struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Content User   `json:"content"`
	}

	resp, err := resty.New().NewRequest().
		SetHeader(cli.tokenHeader()).
		SetHeader(cli.acceptJSONHeader()).
		SetResult(&respData).
		Get(url)

	if err != nil {
		cli.log.Error(err)
		return nil, wrapErr(err)
	}

	if resp.StatusCode() == http.StatusUnauthorized {
		cli.log.Warnf("Fetch user ==> GET %s\nHTTP Status Code: %d\nResponse body: %s", url, resp.StatusCode(), resp.Body())
		return nil, ErrUnauthenticatedUser
	}

	if respData.Status != success {
		cli.log.Warnf("Fetch user ==> GET %s\nHTTP Status Code: %d\nResponse body: %s", url, resp.StatusCode(), resp.Body())
		return nil, ErrUserNotFound // UAM不能明确区分'用户不存在错误'与其他错误(如数据库错误)
	}

	return &respData.Content, nil
}
