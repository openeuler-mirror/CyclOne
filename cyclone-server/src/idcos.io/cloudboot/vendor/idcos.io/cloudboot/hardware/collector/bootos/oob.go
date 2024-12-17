package bootos

import (
	"strconv"

	"idcos.io/cloudboot/hardware/collector"
	"idcos.io/cloudboot/hardware/oob"
)

// OOB 采集并返回当前设备的OOB信息。
func (c *bootosC) OOB() (*collector.OOB, error) {
	worker := oob.SelectWorker(oob.DefaultWorker)
	if worker == nil {
		if log := c.Base.GetLog(); log != nil {
			log.Errorf("Can't find the OOB worker by name: %s", oob.DefaultWorker)
		}
		return nil, nil
	}
	var oob collector.OOB
	if bmc, _ := worker.BMC(); bmc != nil {
		oob.Firmware = bmc.FirmwareReversion
	}

	if network, _ := worker.Network(); network != nil {
		oob.Network = &collector.OOBNetwork{
			IPSrc:   network.IPSrc,
			IP:      network.IP,
			Netmask: network.Netmask,
			Gateway: network.Gateway,
		}
	}

	if users, _ := worker.Users(); len(users) > 0 {
		oobusers := make([]collector.OOBUser, 0, len(users))
		for i := range users {
			user := collector.OOBUser{
				ID:   users[i].ID,
				Name: users[i].Name,
			}
			if users[i].Access != nil {
				user.PrivilegeLevel = strconv.Itoa(users[i].Access.PrivilegeLevel)
			}
			oobusers = append(oobusers, user)
		}
		oob.User = oobusers
	}
	return &oob, nil
}
