package bootos

import (
	"bufio"
	"strings"

	"idcos.io/cloudboot/hardware/collector"
)

// Backplane 采集并返回背板信息。
func (c *bootosC) Backplane() (*collector.Backplane, error) {
	output, err := c.ExecByShell("racadm", "swinventory") // TODO 仅适用于Dell，暂时不知道通用的采集方法。
	if err != nil {
		return nil, err
	}

	var blocks []string
	all := strings.Split(string(output), "\n\n")
	for i := range all {
		if strings.Contains(all[i], "ComponentType = FIRMWAR") &&
			strings.Contains(all[i], "ElementName = Backplane") &&
			strings.Contains(all[i], "Current Version =") {
			blocks = append(blocks, all[i])
		}
	}

	if len(blocks) <= 0 {
		return nil, nil
	}
	var bp collector.Backplane
	for i := range blocks {
		var item collector.BackplaneItem
		scanner := bufio.NewScanner(strings.NewReader(blocks[i]))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "ElementName =") {
				item.Name = c.extractValue(line, eqSeprator)
			}
			if strings.HasPrefix(line, "Current Version =") {
				item.FirmwareVersion = c.extractValue(line, eqSeprator)
			}
		}
		if err := scanner.Err(); err != nil {
			if log := c.GetLog(); log != nil {
				log.Error(err)
			}
			return nil, err
		}
		bp.Items = append(bp.Items, item)
	}
	return &bp, nil
}
