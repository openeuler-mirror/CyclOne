package bootos

import (
	"bufio"
	"bytes"
	"strings"

	"idcos.io/cloudboot/hardware/collector"
)

// PCI 采集并返回当前设备的所有PCI插槽信息。
func (c *bootosC) PCI() (*collector.PCI, error) {
	output, err := c.Base.ExecByShell("dmidecode", "-t", "slot")
	if err != nil {
		return nil, err
	}

	blocks := bytes.Split(output, []byte{'\n', '\n'})
	items := make([]collector.SlotItem, 0, len(blocks)-1)
	for i := range blocks {
		slot, err := c.parseSlotItem(blocks[i])
		if err != nil {
			return nil, err
		}
		if slot == nil {
			continue
		}
		items = append(items, *slot)
	}
	return &collector.PCI{
		TotalSlots: len(items),
		Items:      items,
	}, nil
}

func (c *bootosC) parseSlotItem(block []byte) (*collector.SlotItem, error) {
	if !bytes.Contains(block, []byte("System Slot Information")) {
		return nil, nil
	}
	var slot collector.SlotItem
	scanner := bufio.NewScanner(bytes.NewReader(block))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "Type:") {
			slot.Type = c.extractValue(line, colonSeparator)

		} else if strings.HasPrefix(line, "Designation:") {
			slot.Designation = c.extractValue(line, colonSeparator)

		} else if strings.HasPrefix(line, "Current Usage:") {
			slot.CurrentUsage = c.extractValue(line, colonSeparator)
		}
	}
	return &slot, nil
}
