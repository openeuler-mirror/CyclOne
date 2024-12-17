package bootos

import (
	"idcos.io/cloudboot/hardware/collector"
	"idcos.io/cloudboot/hardware/raid"
)

// RAID 采集并返回当前设备的RAID信息。
func (c *bootosC) RAID() (*collector.RAID, error) {
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
	return &collector.RAID{
		Items: ctrls,
	}, nil
}
