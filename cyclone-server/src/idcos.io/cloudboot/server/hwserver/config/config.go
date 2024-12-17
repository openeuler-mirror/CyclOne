package config

// Configuration hw-server配置结构体
type Configuration struct {
	SN           string // 当前设备sn
	Manufacturer string // 当前设备厂商名

	Logger struct {
		Level string
		Dir   string
	}

	HTTPServer struct {
		Enabled bool
		Port    int
		BaseURL string
	}
}
