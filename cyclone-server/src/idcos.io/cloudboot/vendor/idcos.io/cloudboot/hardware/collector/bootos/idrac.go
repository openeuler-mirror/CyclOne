package bootos

import (
	"bufio"
	"bytes"
	"strings"

	"idcos.io/cloudboot/hardware/collector"
)

// IDRAC 采集并返回Dell iDRAC信息。
func (c *bootosC) IDRAC() (*collector.IDRAC, error) {
	// iDRAC并非只适用于Dell，一些贴牌服务器同样有可能使用iDRAC。
	output, err := c.Base.ExecByShell("racadm", "getsysinfo", "-d")
	if err != nil {
		return nil, err
	}
	var idrac collector.IDRAC
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "RAC Date/Time") {
			idrac.RACDateTime = c.extractValue(line, eqSeprator)

		} else if strings.HasPrefix(line, "Firmware Version") {
			idrac.FirmwareVersion = c.extractValue(line, eqSeprator)

		} else if strings.HasPrefix(line, "Firmware Build") {
			idrac.FirmwareBuild = c.extractValue(line, eqSeprator)

		} else if strings.HasPrefix(line, "Last Firmware Update") {
			idrac.LastFirmwareUpdate = c.extractValue(line, eqSeprator)
		}
	}
	return &idrac, scanner.Err()
}
