package service

import (
	"net/http"
	"reflect"

	"github.com/voidint/binding"
	"github.com/voidint/page"

	"idcos.io/cloudboot/config"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
	uam "idcos.io/uam-go"
)

// GetUserPageReq 查询用户分页请求结构体
type GetUserPageReq struct {
	// DepartmentID 部门ID
	DepartmentID string `json:"department_id"`
	// 用户名
	Name string `json:"name"`
	// 分页页号
	Page int64 `json:"page"`
	// 分页大小
	PageSize int64              `json:"page_size"`
	User     *model.CurrentUser `json:"-"`
}

// FieldMap 请求字段映射
func (reqData *GetUserPageReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.DepartmentID: "department_id",
		&reqData.Name:         "name",
		&reqData.Page:         "page",
		&reqData.PageSize:     "page_size",
	}
}

// Validate 结构体数据校验
func (reqData *GetUserPageReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if reqData.User == nil { // 确保先经过登录认证中间件过滤
		panic("unreachable")
	}
	return errs
}

// GetUserPage 按条件查询当前用户所在租户下的用户分页列表
func GetUserPage(log logger.Logger, conf *config.Config, reqData *GetUserPageReq) (pg *page.Page, err error) {
	if reqData.PageSize <= 0 {
		reqData.PageSize = 10
	}
	if reqData.Page <= 0 {
		reqData.Page = 1
	}

	opts := []uam.SetOptionFunc{
		uam.LogOption(log),
		uam.TenantOption(model.TenantID4UAM),
	}
	pageNo, pageSize, _, totalRecords, users, err := uam.NewClient(conf.UAM.RootEndpoint, reqData.User.Token, opts...).GetUsers(reqData.DepartmentID, reqData.Name, &uam.PageLimiter{
		PageNo:   reqData.Page,
		PageSize: reqData.PageSize,
	})
	if err != nil {
		return nil, err
	}
	pager := page.NewPager(reflect.TypeOf(&model.User{}), pageNo, pageSize, totalRecords)
	for i := range users {
		pager.AddRecords(&model.User{
			ID:        users[i].ID,
			LoginName: users[i].LoginID,
			Name:      users[i].Name,
			Email:     users[i].Email,
			Status:    users[i].Status,
			Tenant: &model.Tenant{
				ID:   users[i].Tenant.ID,
				Name: users[i].Tenant.Name,
			},
			Department: &model.Department{
				ID:   users[i].Department.ID,
				Name: users[i].Department.Name,
			},
		})
	}
	return pager.BuildPage(), nil
}

// ChangeUserPasswordReq 修改用户密码请求结构体
type ChangeUserPasswordReq struct {
	OldPassword string             `json:"old_password"`
	NewPassword string             `json:"new_password"`
	User        *model.CurrentUser `json:"-"`
}

// FieldMap 请求字段映射
func (reqData *ChangeUserPasswordReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.OldPassword: "old_password",
		&reqData.NewPassword: "new_password",
	}
}

// Validate 结构体数据校验
func (reqData *ChangeUserPasswordReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if reqData.User == nil { // 确保先经过登录认证中间件过滤
		panic("unreachable")
	}
	if reqData.OldPassword == "" {
		errs.Add([]string{"old_password"}, binding.RequiredError, "旧密码不能为空")
		return errs
	}
	if reqData.NewPassword == "" {
		errs.Add([]string{"new_password"}, binding.RequiredError, "新密码不能为空")
		return errs
	}
	return errs
}

// ChangeUserPassword 修改当前用户密码
func ChangeUserPassword(log logger.Logger, conf *config.Config, reqData *ChangeUserPasswordReq) (err error) {
	return uam.NewClient(conf.UAM.RootEndpoint, reqData.User.Token, uam.LogOption(log)).ChangeUserPassword(reqData.User.ID, reqData.OldPassword, reqData.NewPassword)
}

// GetUserByID 使用当前用户token去查询目标ID的用户信息
func GetUserByID(log logger.Logger, conf *config.Config, token, id string) (*model.User, error) {
	if id == "" {
		return nil, nil
	}
	opts := []uam.SetOptionFunc{
		uam.LogOption(log),
		//uam.TenantOption(conf.UAM.Customer),
	}

	u, err := uam.NewClient(conf.UAM.RootEndpoint, token, opts...).GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:        u.ID,
		LoginName: u.LoginID,
		Name:      u.Name,
		Email:     u.Email,
		Status:    u.Status,
		Tenant: &model.Tenant{
			ID:   u.Tenant.ID,
			Name: u.Tenant.Name,
		},
		Department: &model.Department{
			ID:   u.Department.ID,
			Name: u.Department.Name,
		},
	}, nil
}
