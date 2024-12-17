package user

import (
	"errors"

	"idcos.io/cloudboot/logger"
	uam "idcos.io/uam-go"
)

const (
	// appID4UAM UAM系统【权限资源管理】-【系统名称】
	appID4UAM = "cloudboot"
	// permissionCategory4UAM UAM系统【权限资源管理】-【编码】
	permissionCategory4UAM = "cloudboot"
	// tenantID4UAM UAM系统默认租户ID
	tenantID4UAM = "default"
)

var (
	// ErrNotExistGetUserFromUAM 错误
	ErrNotExistGetUserFromUAM = errors.New("UAM上不存在的用户")

	// ErrParamGetUserFromUAM 参数错误
	ErrParamGetUserFromUAM = errors.New("参数错误")
)

// GetNameFromUAM 从UAM上获取合适的UAM信息
type GetNameFromUAM func(string) (string, string, error)

// GetUUIDFromUAM 从UAM上获取合适的UAM信息
type GetUUIDFromUAM func(string) (string, error)

// GetEmailFromUAM 从UAM上获取合适的UAM信息
type GetEmailFromUAM func(string) (string, string, string, error)

// GetUsersByUUID 根据UUID从UAM服务器上找到匹配的用户信息
func GetUsersByUUID(log logger.Logger, accessToken, uamEndPoint string) GetNameFromUAM {
	return func(uuid string) (loginName, name string, err error) {
		if accessToken == "" || uamEndPoint == "" || uuid == "" {
			return "", "", ErrParamGetUserFromUAM
		}

		user, err := uam.NewClient(uamEndPoint, accessToken, uam.LogOption(log)).GetUserByID(uuid)
		if err != nil {
			return "", "", err
		}

		if user != nil {
			return user.LoginID, user.Name, nil
		}

		return "", "", ErrNotExistGetUserFromUAM
	}
}

// GetEmailByUUID 根据UUID从UAM服务器上找到匹配的用户信息
func GetEmailByUUID(log logger.Logger, accessToken, uamEndPoint string) GetEmailFromUAM {
	return func(uuid string) (loginName, name, email string, err error) {
		if accessToken == "" || uamEndPoint == "" || uuid == "" {
			return "", "", "", ErrParamGetUserFromUAM
		}

		user, err := uam.NewClient(uamEndPoint, accessToken, uam.LogOption(log)).GetUserByID(uuid)
		if err != nil {
			return "", "", "", err
		}

		if user != nil {
			return user.LoginID, user.Name, user.Email, nil
		}
		return "", "", "", ErrNotExistGetUserFromUAM
	}
}

// GetUUIDByLoginName 根据LoginName从UAM服务器上找到匹配的用户UUID
func GetUUIDByLoginName(log logger.Logger, accessToken, uamEndPoint string) GetUUIDFromUAM {
	return func(loginName string) (uuid string, err error) {
		if accessToken == "" || uamEndPoint == "" || loginName == "" {
			return "", ErrParamGetUserFromUAM
		}
		opts := []uam.SetOptionFunc{
			uam.LogOption(log),
			uam.TenantOption(tenantID4UAM),
		}
		user, err := uam.NewClient(uamEndPoint, accessToken, opts...).GetUserByLoginID(loginName)
		if err != nil {
			return "", err
		}

		if user != nil {
			return user.ID, nil
		}
		return "", ErrNotExistGetUserFromUAM
	}
}

// GetUserByLoginName 根据LoginName从UAM服务器上找到匹配的用户
func GetUserByLoginName(log logger.Logger, accessToken, uamEndPoint string) GetNameFromUAM {
	return func(loginName string) (login, name string, err error) {
		if accessToken == "" || uamEndPoint == "" || loginName == "" {
			return "", "", ErrParamGetUserFromUAM
		}

		user, err := uam.NewClient(uamEndPoint, accessToken, uam.LogOption(log)).GetUserByID(loginName)
		if err != nil {
			return "", "", err
		}

		if user != nil {
			return user.LoginID, user.Name, nil
		}

		return "", "", ErrNotExistGetUserFromUAM
	}
}
