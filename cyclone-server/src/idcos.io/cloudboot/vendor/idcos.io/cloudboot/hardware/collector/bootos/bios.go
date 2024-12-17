package bootos

import (
	"bufio"
	"bytes"
	"strings"

	"idcos.io/cloudboot/hardware/collector"
)

// BIOS 采集并返回当前设备的BIOS信息。
func (c *bootosC) BIOS() (*collector.BIOS, error) {
	output, err := c.Base.ExecByShell("dmidecode", "-t", "bios")
	if err != nil {
		return nil, err
	}

	var bios collector.BIOS
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "Vendor:") {
			bios.Vendor = c.extractValue(line, colonSeparator)
		} else if strings.HasPrefix(line, "Version:") {
			bios.Version = c.extractValue(line, colonSeparator)
		}
	}
	return &bios, nil
}
