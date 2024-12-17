package bootos

import (
	"bufio"
	"bytes"
	"strings"

	"idcos.io/cloudboot/hardware/collector"
	mybytes "idcos.io/cloudboot/utils/bytes"
)

// Memory 采集并返回当前设备内存信息
func (c *bootosC) Memory() (*collector.Memory, error) {
	var err error
	var mem collector.Memory
	// 采集所有内存条信息
	mem.Items, err = c.memoryDevices()
	if err != nil {
		return nil, err
	}
	mem.TotalSize, mem.TotalSizeMB, err = c.calcMemTotalSize(mem.Items)
	if err != nil {
		return nil, err
	}
	return &mem, nil
}

// memoryDevices 采集已经插在插槽上的内存条信息列表
func (c *bootosC) memoryDevices() (items []collector.MemoryDevice, err error) {
	output, err := c.Base.ExecByShell("dmidecode", "-t", "memory")
	if err != nil {
		return nil, err
	}
	array := bytes.Split(output, []byte("Memory Device"))
	for i := range array {
		item := c.parseMemoryDevice(array[i])
		if item == nil {
			continue
		}
		items = append(items, *item)
	}
	return items, nil
}

func (c *bootosC) parseMemoryDevice(out []byte) *collector.MemoryDevice {
	var memDev collector.MemoryDevice
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.Contains(line, ":") {
			continue
		}

		if strings.HasPrefix(line, "Size:") {
			if strings.Contains(line, "No Module Installed") {
				return nil // 该插槽上未插内存条
			}
			memDev.Size = c.extractValue(line, colonSeparator)

		} else if strings.HasPrefix(line, "Speed:") {
			memDev.Speed = c.extractValue(line, colonSeparator)

		} else if strings.HasPrefix(line, "Locator:") {
			memDev.Locator = c.extractValue(line, colonSeparator)

		} else if strings.HasPrefix(line, "Type:") {
			memDev.Type = c.extractValue(line, colonSeparator)
		}
	}
	if memDev.Locator == "" && memDev.Size == "" && memDev.Speed == "" && memDev.Type == "" {
		return nil
	}
	return &memDev
}

// calcDiskTotalSize 计算逻辑磁盘总容量。
// bsize 单位为Byte的容量值。
// mbsize 单位为MB的容量值。
func (c *bootosC) calcMemTotalSize(items []collector.MemoryDevice) (bsize int64, mbsize int, err error) {
	for i := range items {
		fields := strings.Fields(items[i].Size)
		if len(fields) != 2 {
			continue
		}
		fields[0] = strings.TrimSpace(fields[0])
		fields[1] = strings.TrimSpace(fields[1])
		size, err := mybytes.Parse2Byte(fields[0], fields[1])
		if err != nil {
			if log := c.Base.GetLog(); log != nil {
				log.Errorf("Parse2Byte(%s, %s) error: %s", fields[0], fields[1], err.Error())
			}
			return 0, 0, err
		}
		bsize += int64(size)
	}
	return bsize, mybytes.Byte2MBRounding(mybytes.Byte(bsize)), nil
}
