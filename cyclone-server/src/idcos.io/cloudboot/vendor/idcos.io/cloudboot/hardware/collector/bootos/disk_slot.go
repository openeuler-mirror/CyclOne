package bootos

import (
	"idcos.io/cloudboot/hardware/collector"
	"idcos.io/cloudboot/hardware/raid"
)

// DiskSlot 采集并返回当前设备的磁盘槽位(物理驱动器)信息
func (c *bootosC) DiskSlot() (*collector.DiskSlot, error) {
	name, err := raid.Whoami()
	if err != nil {
		if log := c.Base.GetLog(); log != nil {
			log.Error(err)
		}
		return nil, err
	}
	worker := raid.SelectWorker(name)
	if worker == nil {
		if log := c.Base.GetLog(); log != nil {
			log.Errorf("Can't find the RAID worker by name: %s", name)
		}
		return nil, nil
	}

	ctrls, err := worker.Controllers()
	if err != nil {
		return nil, err
	}

	var items []raid.PhysicalDrive
	for i := range ctrls {
		drives, err := worker.PhysicalDrives(ctrls[i].ID)
		if err != nil {
			return nil, err
		}
		items = append(items, drives...)
	}
	return &collector.DiskSlot{
		Items: items,
	}, nil
}
