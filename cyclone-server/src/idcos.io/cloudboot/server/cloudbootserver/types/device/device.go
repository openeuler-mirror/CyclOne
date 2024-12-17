package device

import (
	"fmt"
	"strings"

	"idcos.io/cloudboot/hardware/collector"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils/collection"
)

// Device 设备信息采集上报请求数据
type Device struct {
	SN          string                 `json:"sn"`          // 设备序列号
	Vendor      string                 `json:"vendor"`      // 设备厂商名
	Model       string                 `json:"model"`       // 设备型号
	Arch        string                 `json:"arch"`        // CPU硬件架构
	BootOSIP    string                 `json:"bootos_ip"`   // bootos ip地址
	BootOSMac   string                 `json:"bootos_mac"`  // bootos ip对应的mac地址
	NICDevice   string                 `json:"nic_device"`  // 根据NIC计算获取
	CPU         *collector.CPU         `json:"cpu"`         // CPU
	Memory      *collector.Memory      `json:"memory"`      // 内存
	Disk        *collector.Disk        `json:"disk"`        // 逻辑磁盘
	DiskSlot    *collector.DiskSlot    `json:"disk_slot"`   // 磁盘槽位(物理驱动器)
	NIC         *collector.NIC         `json:"nic"`         // 网卡
	Motherboard *collector.Motherboard `json:"motherboard"` // 主板
	OOB         *collector.OOB         `json:"oob"`         // 带外
	BIOS        *collector.BIOS        `json:"bios"`        // BIOS
	RAID        *collector.RAID        `json:"raid"`        // RAID
	Power       *collector.Power       `json:"power"`       // 电源
	Fan         *collector.Fan         `json:"fan"`         // 风扇
	PCI         *collector.PCI         `json:"pci"`         // PCI插槽
	HBA         *collector.HBA         `json:"hba"`         // HBA
	LLDP        *collector.LLDP        `json:"lldp"`        // LLDP
	Extra       *collector.Extra       `json:"extra"`       // 自定义扩展字段
}

// Setup 设置
func (reqData *Device) Setup() {
	if reqData.OOB != nil && reqData.OOB.Network != nil {
		if reqData.OOB.Network.IPSrc == "Static Address" {
			reqData.OOB.Network.IPSrc = model.IPSourceStatic
		} else {
			reqData.OOB.Network.IPSrc = model.IPSourceDHCP
		}
	}

	// 计算IP、MAC、NICDevice
	var uniqNICs []string
	if reqData.NIC != nil {
		for i := range reqData.NIC.Items {
			nic := fmt.Sprintf("%s %s", reqData.NIC.Items[i].Company, reqData.NIC.Items[i].ModelName)
			if !collection.InSlice(nic, uniqNICs) {
				uniqNICs = append(uniqNICs, nic)
			}

			if reqData.NIC.Items[i].IP != "" {
				reqData.BootOSIP = reqData.NIC.Items[i].IP
				reqData.BootOSMac = reqData.NIC.Items[i].Mac
			}
		}
	}
	reqData.NICDevice = strings.Join(uniqNICs, "\n")
}
