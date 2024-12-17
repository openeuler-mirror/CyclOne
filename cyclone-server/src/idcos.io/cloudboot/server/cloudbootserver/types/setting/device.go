package setting

// NetworkSetting 装机参数-业务网络参数
type NetworkSetting struct {
	Hostname string               `json:"hostname"`  // 主机名
	IPSource string               `json:"ip_source"` // 业务网络类型
	SkipPing string               `json:"skip_ping"` // 是否跳过对业务网络IP的ping检查(仅业务网络类型为static时有效)
	Items    []NetworkSettingItem `json:"items"`
}

// NetworkSettingItem 装机参数-业务网络参数条目
type NetworkSettingItem struct {
	NetworkID int    `json:"network_id"` // 所属网段的ID
	IP        string `json:"ip"`
	Mac       string `json:"mac"`
	Netmask   string `json:"netmask"`
	Gateway   string `json:"gateway"`
	VLAN      string `json:"vlan"`
	Trunk     string `json:"trunk"`
	Bonding   string `json:"bonding"`
	DNS       string `json:"dns"`
}

// PartitionSettingItem 装机参数-分区参数
type PartitionSettingItem struct {
	Disk       int    `json:"disk"`
	Size       string `json:"size"`
	FS         string `json:"fs"`
	Mountpoint string `json:"mountpoint"`
}
// NICBondingSettingItem 网卡绑定信息结构体
type NICBondingSettingItem struct {
	Name string    `json:"bonding_name"` // bonding名称
	Mode string    `json:"bonding_mode"` // bonding模式
	NICs []BondNIC `json:"nics"`         // 需要绑定的网卡，暂时支持2个网卡的绑定。
}

// BondNIC bonding用的网卡信息结构体
type BondNIC struct {
	MacAddr string `json:"mac_addr"` // mac地址
	Type    string `json:"type"`     // 该网卡在绑定中所处的角色/地位。可选值：master|slave
}

// OSTemplateSetting 操作系统模板配置参数
type OSTemplateSetting struct {
	ID       int    `json:"id"`        // 模板ID
	Family   string `json:"family"`    // 操作系统族系
	Name     string `json:"name"`      // 模板名称
	BootMode string `json:"boot_mode"` // 启动模式。可选值: legacy_bios-传统BIOS模式; uefi-UEFI模式;
}

// OSUserSettingItem 装机参数-操作系统管理员用户
type OSUserSettingItem struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
