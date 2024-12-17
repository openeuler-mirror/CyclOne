package service

import (
	"idcos.io/cloudboot/hardware/collector"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/server/cloudbootserver/types/device"
	"idcos.io/cloudboot/server/hwserver/config"
)

// CollectDevice 采集设备信息并丢弃采集过程中的错误。
func CollectDevice(conf *config.Configuration, log logger.Logger) (*device.Device, error) {
	c := collector.SelectCollector(collector.DefaultCollector)
	if c == nil {
		log.Warnf("Unregistered collector: %s", collector.DefaultCollector)
		return nil, collector.ErrUnregisteredCollector
	}
	c.SetLog(log)

	var dev device.Device
	_, dev.SN, dev.Vendor, dev.Model, dev.Arch, _ = c.BASE()
	dev.CPU, _ = c.CPU()
	dev.Memory, _ = c.Memory()
	dev.Disk, _ = c.Disk()
	dev.DiskSlot, _ = c.DiskSlot()
	dev.NIC, _ = c.NIC()
	dev.Motherboard, _ = c.Motherboard()
	dev.OOB, _ = c.OOB()
	dev.BIOS, _ = c.BIOS()
	dev.RAID, _ = c.RAID()
	dev.Fan, _ = c.Fan()
	dev.Power, _ = c.Power()
	dev.PCI, _ = c.PCI()
	dev.HBA, _ = c.HBA()
	dev.LLDP, _ = c.LLDP()
	dev.Setup()
	return &dev, nil
}
