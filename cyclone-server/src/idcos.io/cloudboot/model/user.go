package model

import "idcos.io/cloudboot/utils/collection"

const (
	// AppID4UAM UAM系统【权限资源管理】-【系统名称】
	AppID4UAM = "cloudboot"
	// PermissionCategory4UAM UAM系统【权限资源管理】-【编码】
	PermissionCategory4UAM = "cloudboot_webank"
	// TenantID4UAM UAM系统默认租户ID
	TenantID4UAM = "default"

	//MenuPermissionType 菜单权限类型
	MenuPermissionType = "cloudboot_menu_permission"
	//ButtonPermissionType 按钮权限类型
	ButtonPermissionType = "cloudboot_button_permission"
	//DataPermissionType 数据权限类型
	DataPermissionType = "cloudboot_data_permission"
	//APIPermissionType API接口权限类型
	APIPermissionType = "cloudboot_api_permission"
)

// Tenant 租户
type Tenant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Department 部门
type Department struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CurrentUser 当前用户信息 // TODO 嵌入User结构体
type CurrentUser struct {
	Token           string      `json:"token"`
	ID              string      `json:"id"`
	LoginName       string      `json:"login_name"`
	Name            string      `json:"name"`
	Email           string      `json:"email"`
	Status          string      `json:"status"`
	UAMPortalURL    string      `json:"uam_portal_url"`
	PermissionCodes []string    `json:"permissions"`
	Tenant          *Tenant     `json:"tenant"`
	Department      *Department `json:"department"`
}

// Allow 若当前用户自身所拥有的权限码与目标权限码中的任意一个重合，则返回true。反之，返回false。
func (u *CurrentUser) Allow(codes ...string) bool {
	if len(u.PermissionCodes) <= 0 {
		return false
	}
	return collection.SSliceContainsAny(u.PermissionCodes, codes)
}

// User 用户基本信息
type User struct {
	ID         string      `json:"id"`
	LoginName  string      `json:"login_name"`
	Name       string      `json:"name"`
	Email      string      `json:"email"`
	Status     string      `json:"status"`
	Tenant     *Tenant     `json:"tenant"`
	Department *Department `json:"department"`
}
