package raid

const (
	// Sep 驱动器分隔符
	Sep = ","
	// 空格分隔符（AdaptecSamrtRAID)
	SpaceSep = " "
)

// Setting RAID配置参数
type Setting struct {
	Controllers []ControllerSetting `json:"controllers"`
}

// ControllerSetting 单块RAID卡控制器配置
type ControllerSetting struct {
	Index     uint           `json:"index"`               // RAID控制器索引号。注意，此'索引号'代表的是RAID卡的顺序号(多块RAID卡场景下)，自'0'始。0表示首块RAID卡，以此类推。
	Clear     string         `json:"clear,omitempty"`     // 是否擦除当前RAID卡的配置。可选值: ON|OFF
	DiskInit  string         `json:"disk_init,omitempty"` // 是否初始化逻辑磁盘。可选值: ON|OFF
	Hotspares string         `json:"hotspares,omitempty"` // 全局热备盘物理驱动器。多个物理驱动用英文逗号分隔。
	JBODs     string         `json:"jbods,omitempty"`     // 待设置为直通盘的物理驱动器。多个物理驱动用英文逗号分隔。
	Arrays    []ArraySetting `json:"arrays"`              // 独立冗余磁盘阵列
}

// ArraySetting 单组RAID配置参数
type ArraySetting struct {
	Level     string `json:"level,omitempty"`     // RAID级别
	Drives    string `json:"drives,omitempty"`    // 物理驱动器标识。多个物理驱动用英文逗号分隔。
	Hotspares string `json:"hotspares,omitempty"` // 局部热备盘物理驱动器（保留属性，暂未使用）。多个物理驱动用英文逗号分隔。
}
