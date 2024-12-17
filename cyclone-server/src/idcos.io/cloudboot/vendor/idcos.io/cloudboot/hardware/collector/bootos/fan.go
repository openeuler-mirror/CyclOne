package bootos

import (
	"bufio"
	"bytes"
	"strings"

	"idcos.io/cloudboot/hardware/collector"
)

// Fan 采集并返回当前设备的所有风扇信息。
func (c *bootosC) Fan() (*collector.Fan, error) {
	output, err := c.Base.ExecByShell("ipmitool", "sdr", "list", "|", "grep", "-i", `"Fan"`, "|", "grep", "-i", `"RPM"`)
	if err != nil {
		return nil, err
	}

	var items []collector.FanItem
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Split(line, "|")
		if len(fields) < 2 {
			continue
		}
		items = append(items, collector.FanItem{
			ID:    strings.TrimSpace(fields[0]),
			Speed: strings.TrimSpace(fields[1]),
		})
	}
	return &collector.Fan{
		Items: items,
	}, nil
}
