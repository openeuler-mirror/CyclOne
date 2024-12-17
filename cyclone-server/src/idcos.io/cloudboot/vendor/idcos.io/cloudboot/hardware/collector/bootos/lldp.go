package bootos

import (
	"encoding/json"

	"idcos.io/cloudboot/hardware/collector"
)

// LLDP 采集并返回当前设备LLDP信息。
func (c *bootosC) LLDP() (*collector.LLDP, error) {
	output, err := c.Base.ExecByShell("lldpctl", "-f", "json")
	if err != nil {
		return nil, err
	}
	var lldp collector.LLDP
	return &lldp, json.Unmarshal(output, &lldp)
}
