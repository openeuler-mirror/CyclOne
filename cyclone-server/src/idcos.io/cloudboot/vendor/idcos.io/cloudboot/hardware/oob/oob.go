package oob

import (
	"errors"

	"idcos.io/cloudboot/hardware"
	"idcos.io/cloudboot/logger"
)

const (
	// DefaultWorker 默认处理器名称
	DefaultWorker = "IPMI"
)

const (
	// DHCP IP来源-DHCP
	DHCP = "dhcp"
	// Static IP来源-静态IP
	Static = "static"
)

var (
	// ErrChannelNotFound impi channel未发现
	ErrChannelNotFound = errors.New("channel not found")
	// ErrUnknownHardware 未知的OOB硬件类型
	ErrUnknownHardware = errors.New("unknown OOB hardware")
)

// Worker OOB处理器
type Worker interface {
	NetworkWorker
	UserWorker
	BMCWorker

	// SetDebug 设置是否开启debug。若开启debug，会将关键日志信息写入console。
	SetDebug(debug bool)
	// SetLog 更换日志实现。默认情况下内部无日志实现。
	SetLog(log logger.Logger)
	// Channel 返回impi channel
	Channel() (int, error)
	// PostCheck OOB配置实施后置检查
	PostCheck(sett *Setting) []hardware.CheckingItem
}

// BMC OOB的BMC信息
type BMC struct {
	FirmwareReversion string
	IPMIVersion       string
	ManufacturerID    string
	ManufacturerName  string
}

// Access OOB channel access
type Access struct {
	Channel        int
	MaxUserIDs     int
	EnabledUserIDs int
	Accesses       []UserAccess
}

// UserAccess OOB用户Access
type UserAccess struct {
	UserID             int
	UserName           string
	FixedName          string
	AccessAvailable    string
	LinkAuthentication string
	IPMIMessaging      string
	PrivilegeLevel     int
}

// Network OOB网络信息
type Network struct {
	IPSrc   string // 可选值: static|dhcp
	Mac     string // IP对应Mac地址
	IP      string // BootOSIP
	Netmask string // 子网掩码
	Gateway string // 默认网关
}

// NetworkWorker OOB网络模块处理器
type NetworkWorker interface {
	// SetDHCP 设置IP来源是DHCP
	SetDHCP() error
	// SetStaticIP 设置IP来源是静态IP
	SetStaticIP(ip, netmask, gateway string) error
	// Network 返回OOB网络信息
	Network() (*Network, error)
}

// User OOB用户信息
type User struct {
	Channel int
	ID      int
	Name    string
	Access  *UserAccess
}

// UserWorker OOB用户模块处理器
type UserWorker interface {
	// GenerateUser 生成用户帐号。
	// 若用户（以用户名为准）未存在，则新增帐号。
	// 若用户已经存在则修改用户密码、权限级别等属性。
	GenerateUser(sett *UserSettingItem) error
	// Users 返回OOB用户列表
	Users() ([]User, error)
}

// BMCWorker BMC模块处理器
type BMCWorker interface {
	// BMC 返回OOB的BMC信息
	BMC() (*BMC, error)
	// BMCColdReset (冷)重启BMC
	BMCColdReset() error
}

// Whoami 返回当前的BIOS固件对应的处理器名
func Whoami() (worker string, err error) {
	return DefaultWorker, nil // 暂时只有一个实现
}
