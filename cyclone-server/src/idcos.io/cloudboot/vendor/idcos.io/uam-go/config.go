package uam

// Config 配置项
type Config struct {
	// RootEndpoint UAM系统的根URL路径(不以'/'结尾)
	RootEndpoint string
	// TenantID 租户ID
	TenantID string
	// Token UAM系统的access token
	Token string
	// Log 日志实现
	Log Logger
}

const (
	// DefaultTenantID 默认租户名
	DefaultTenantID = "default"
)

// LoadDefault 加载默认的配置项
func (c *Config) LoadDefault() {
	if c.TenantID == "" {
		c.TenantID = DefaultTenantID
	}

	if c.Log == nil {
		c.Log = defaultLog
	}
}

// SetOptionFunc 函数式选项设置器
type SetOptionFunc func(c *Config)

// TenantOption 返回设置租户ID选项的函数
func TenantOption(tenantID string) SetOptionFunc {
	return func(c *Config) {
		c.TenantID = tenantID
	}
}

// LogOption 返回设置日志选项的函数
func LogOption(log Logger) SetOptionFunc {
	return func(c *Config) {
		c.Log = log
	}
}
