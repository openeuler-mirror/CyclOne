package config

// Loader 定义统一的配置加载接口
type Loader interface {
	Load() (*Config, error)
	Save(*Config) error
}

// Config config 数据结构体
type Config struct {
	Logger           Logger
	Repo             Repo
	Server           Server
	ExternalService	 ExternalService
	UAM              UAM
	Samba            Samba
	DHCPLimiter      DHCPLimiter
	Crypto           Crypto
	ReverseProxy     ReverseProxy
	DistributedNodes []DistributedNode
	IP               IP
}

// Logger 日志配置结构
type Logger struct {
	Color          bool
	Level          string // 日志级别，可选值：debug、info、warn、error
	LogFile        string // 文件日志路径。若是空字符串则不输出到文件。
	PanicLogFile   string // http请求Panic记录日志
	FilePerm       string // 文件权限，默认0644
	ConsoleEnabled bool   // 是否将日志输出到控制台
	RotateEnabled  bool   // 是否打开日志轮替
}

const (
	// ConsoleLog 控制台日志
	ConsoleLog = "stdout"
)

// Repo 数据库配置结构
type Repo struct {
	Debug          bool   // 是否开启debug模式
	LogDestination string // 可选值：stdout/$(filepath)，分别表示将SQL执行日志写入终端或者写入文件。默认为stdout，若指定了具体的日志文件路径，则将SQL日志写入该文件。
	Connection     string `ini:"connection"`
}

// Server server服务配置
type Server struct {
	HTTPPort           int
	OOBDomain          string
	StorageRootDir     string
}

// ExternalService 调用外部API服务配置
type ExternalService struct {
	ESBBaseURL				string
	ESBAppID				string
	ESBAppToken				string
}

// UAM uam系统配置
type UAM struct {
	RootEndpoint string // UAM后端服务根URL，必选。如'http://127.0.0.1:8092'
	PortalURL    string // UAM前端系统URL，必选。如'http://127.0.0.1:92'
}

// DHCPLimiter DHCP IP请求限流器配置
type DHCPLimiter struct {
	Enable         bool // 是否打开dhcp ip请求限流功能
	Limit          int  // 令牌桶中令牌上限
	WaitingTimeout int  // 等待获取令牌的超时时间(单位秒)
}

// ReverseProxy HTTP反向代理服务设置
type ReverseProxy struct {
	Enable   bool   // 标识当前server是否为proxy server
	HTTPPort int    // 反向代理服务的端口号
	URL      string // 若当前server为proxy server，则需要将所有HTTP请求转发至master server处理。
	Origin   string // 标识当前代理服务器所属地域信息
	IP       string // 当前部署服务器的IP地址（可选）
}

//DistributedNode 各分布式节点的部署IP
type DistributedNode struct {
	NodeName string
	NodeIP   []string
}

// Crypto 加密key
// key可参考以下方式生成
// key := make([]byte, 32)
// if _, err := rand.Read(key); err != nil {
//     panic(err)
// }
// keyStr := base64.StdEncoding.EncodeToString(key)
type Crypto struct {
	Key string `json:"key"`
}

// Samba Samba配置信息，地址、用户名、密码
type Samba struct {
	// 地址
	Server string `json:"server"`
	// 用户名
	User string `json:"user"`
	// 密码
	Password string `json:"password"`
}

// IP保留
type IP struct {
	// 天数配置, 0代表不保留
	ReserveDay int `json:"reserveday"`
}