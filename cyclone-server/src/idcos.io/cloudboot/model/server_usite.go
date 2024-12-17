package model

import (
	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
)

const (
	// USiteStatFree 机位(U位)状态-空闲
	USiteStatFree = "free"
	// USiteStatPreOccupied 机位(U位)状态-预占用
	USiteStatPreOccupied = "pre_occupied"
	// USiteStatUsed 机位(U位)状态-已使用
	USiteStatUsed = "used"
	// USiteStatDisabled 机位(U位)状态-不可用
	USiteStatDisabled = "disabled"
	// 端口速率枚举值
	PortRateDefault = "25GE"
	PortRateGE = "GE"
	PortRate10GE = "10GE"
	PortRate25GE = "25GE"
	PortRate40GE = "40GE"
)

var PortRateMap = map[string]string{
	"GE":   PortRateGE,
	"10GE": PortRate10GE,
	"25GE": PortRate25GE,
	"40GE": PortRate40GE,
}

// CombinedServerUSite 机位分页查询
type CombinedServerUSite struct {
	IDCID           []uint `json:"idc_id"`            	// 所属数据中心ID
	ServerRoomID    []uint `json:"server_room_id"`    	// 所属机房ID
	ServerCabinetID []uint `json:"server_cabinet_id"` 	// 所属机架ID
	NetAreaID       []uint `json:"network_area_id"`   	// 网络区域ID
	ServerRoomName string
	ServerRoomNameCabinetNumUSiteNumSlice  []string    		// 机房管理单元-机架-机位
	ServerRoomNameCabinetNumUSiteNum  string    		    // 机房管理单元-机架-机位
	PhysicalArea   string `json:"physical_area"`  		// 物理区域
	CabinetNumber  string `json:"cabinet_number"` 		// 机架编号
	USiteNumber    string `json:"usite_number"`   		// 机位编号
	Height         string `json:"height"`         		// 机位调试
	Status         string `json:"status"`         		// 机位状态
	LAWAPortRate   string `json:"la_wa_port_rate"`		// 内外网端口速率：GE\10GE\25GE\40GE
}

// SwitchInfo 交换机信息
type SwitchInfo struct {
	Name string `json:"name"`
	Port string `json:"port"`
}

// ServerUSite 机位(U位)
type ServerUSite struct {
	gorm.Model
	IDCID            uint   `gorm:"column:idc_id"`
	ServerRoomID     uint   `gorm:"column:server_room_id"`
	ServerCabinetID  uint   `gorm:"column:server_cabinet_id"`
	Number           string `gorm:"column:number"`
	Beginning        uint   `gorm:"column:beginning"`
	Height           uint   `gorm:"column:height"`
	PhysicalArea     string `gorm:"column:physical_area"`
	OobnetSwitches   string `gorm:"oobnet_switches"`
	IntranetSwitches string `gorm:"intranet_switches"`
	ExtranetSwitches string `gorm:"extranet_switches"`
	LAWAPortRate	 string `gorm:"column:la_wa_port_rate"`
	Status           string `gorm:"column:status"`
	Remark           string `gorm:"column:remark"`
	Creator          string `gorm:"column:creator"`
}

// BeforeSave 机位信息保存时防止空字符串信息
func (usite *ServerUSite) BeforeSave() (err error) {
	replaceIfBlank(&usite.OobnetSwitches, EmptyJSONArray)
	replaceIfBlank(&usite.IntranetSwitches, EmptyJSONArray)
	replaceIfBlank(&usite.ExtranetSwitches, EmptyJSONArray)
	return
}

// TableName 指定数据库表名
func (ServerUSite) TableName() string {
	return "server_usite"
}

// IServerUSite 机位(U位)持久化接口
type IServerUSite interface {
	//GetServerUSiteCountByServerCabinetID 统计机架的机位数
	GetServerUSiteCountByServerCabinetID(id uint) (count int64, err error)
	// SaveServerUSite 保存机位(U位)
	SaveServerUSite(*ServerUSite) (affected int64, err error)
	// GetServerUSiteByID 返回指定ID的机位(U位)
	GetServerUSiteByID(id uint) (*ServerUSite, error)
	// GetServerUSiteByNumber 返回指定编号的机位(U位)
	GetServerUSiteByNumber(cabinetID uint, number string) (*ServerUSite, error)
	// BatchUpdateServerUSitesStatus 批量更新机位状态信息
	BatchUpdateServerUSitesStatus(id []uint, status string) (affected int64, err error)
	// BatchUpdateServerUSitesRemark 批量更新机位备注信息
	BatchUpdateServerUSitesRemark(id []uint, remark string) (affected int64, err error)	
	// DeleteServerUSitePort 删除机位端口号
	DeleteServerUSitePort(id uint) (affected int64, err error)
	// RemoveServerUSiteByID 删除机位
	RemoveServerUSiteByID(id uint) (affected int64, err error)
	// CountServerUSite 统计机位数量
	CountServerUSite(cond *CombinedServerUSite) (count int64, err error)
	//GetServerUSiteID 返回满足过滤条件的机位(U位) 不支持模糊查找
	GetServerUSiteID(cond *ServerUSite) (id []uint, err error)
	// CountServerUSite 查询机位信息
	GetServerUSiteByCond(cond *CombinedServerUSite, orderby OrderBy, limiter *page.Limiter) (items []*ServerUSite, err error)
	// 根据网络设备名，获取对应的机位ID
	GetServerUsiteByNetworkDeviceName(name []string) (nd []uint, err error)
	GetPhysicalAreas() (*DeviceQueryParamResp, error)
}
