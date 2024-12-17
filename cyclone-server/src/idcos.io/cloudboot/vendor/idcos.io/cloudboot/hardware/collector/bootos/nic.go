package bootos

import (
	"bufio"
	"bytes"
	"net"
	"regexp"
	"strings"

	"idcos.io/cloudboot/hardware/collector"
)

// NIC 采集并返回当前设备的网卡信息
func (c *bootosC) NIC() (*collector.NIC, error) {
	items, err := c.nicItemsBase()
	if err != nil {
		return nil, err
	}
	pairs, err := c.busDesignationPairs()
	if err != nil {
		return nil, err
	}
	c.populateNICItems(&items, pairs)
	return &collector.NIC{
		Items: items,
	}, nil
}

// nicItemsBase 返回包含基本信息(名称、mac、ip)的网卡列表
func (c *bootosC) nicItemsBase() (items []collector.NICItem, err error) {
	nics, err := net.Interfaces()
	if err != nil {
		if log := c.GetLog(); log != nil {
			log.Error(err)
		}
		return nil, err
	}
	for i := range nics {
		if strings.HasPrefix(strings.ToLower(nics[i].Name), "lo") ||
			strings.HasPrefix(strings.ToLower(nics[i].Name), "vnet") {
			continue
		}
		ip, _ := c.findIPByNICName(nics[i].Name)

		items = append(items, collector.NICItem{
			Name: nics[i].Name,
			Mac:  nics[i].HardwareAddr.String(),
			IP:   ip,
		})
	}
	return items, nil
}

var ipReg = regexp.MustCompile(`^(\d+.){3}\d+$`)

func (c *bootosC) findIPByNICName(name string) (ip string, err error) {
	// ip -4 -o address show
	out, err := c.ExecByShell("ip", "-4", "-o", "address", "show")
	if err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.Contains(line, name) {
			continue
		}
		fields := strings.Fields(line)
		for i := range fields {
			idx := strings.Index(fields[i], "/")
			if idx < 0 {
				continue
			}
			ip = fields[i][:idx]
			if ipReg.MatchString(ip) {
				return ip, nil
			}
		}
	}
	return "", nil
}

func (c *bootosC) busDesignationPairs() (pairs map[string]string, err error) {
	output, err := c.Base.ExecByShell("dmidecode", "-t", "slot")
	if err != nil {
		return nil, err
	}
	pairs = make(map[string]string)
	blocks := bytes.Split(output, []byte{'\n', '\n'})
	for i := range blocks {
		bus, designation := c.parseBusDesignation(blocks[i])
		if bus == "" || designation == "" {
			continue
		}
		pairs[bus] = designation
	}
	return pairs, nil
}

func (c *bootosC) parseBusDesignation(out []byte) (bus, designation string) {
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "Current Usage:") && c.extractValue(line, colonSeparator) != "In Use" {
			return "", ""
		} else if strings.HasPrefix(line, "Bus Address:") {
			busRaw := c.extractValue(line, colonSeparator)
			if strings.Contains(busRaw, ".") {
				bus = strings.Split(busRaw, ".")[0]
			}
		} else if strings.HasPrefix(line, "Designation:") {
			designation = c.extractValue(line, colonSeparator)
		}
	}
	return bus, designation
}

// populateNICItems 将列表中的网卡信息填充完整
func (c *bootosC) populateNICItems(items *[]collector.NICItem, pairs map[string]string) {
	for i := range *items {
		_ = c.populateNICItem(&((*items)[i]), pairs)
	}
}

func (c *bootosC) populateNICItem(item *collector.NICItem, pairs map[string]string) (err error) {
	defer c.populateNICStatusAndSpeed(item)

	output, err := c.Base.ExecByShell("ethtool", "-i", item.Name)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "firmware-version:") {
			item.FirmwareVersion = c.extractValue(line, colonSeparator)
		}
		if strings.HasPrefix(line, "bus-info:") {
			item.Businfo = c.extractValue(line, colonSeparator)
		}
	}

	businfo := item.Businfo // 0000:02:00.0
	if !strings.Contains(businfo, ".") || !strings.Contains(businfo, ":") {
		return nil
	}
	subBusDot := strings.Split(businfo, ".")[0] //0000:02:00
	if designation, ok := pairs[subBusDot]; ok {
		item.Side = "outside"
		item.Designation = designation
	} else {
		item.Side = "inside"
	}
	subBusColon := strings.Join(strings.Split(businfo, ":")[1:], ":") //02:00.0
	out, err := c.Base.ExecByShell("lspci", "|", "grep", subBusColon)
	if err != nil {
		return err
	}
	// 01:00.0 Ethernet controller: Broadcom Corporation NetXtreme BCM5720 Gigabit Ethernet PCIe
	if fields := strings.Fields(strings.TrimSpace(string(out))); len(fields) > 5 { // TODO 使用更加严谨的方式采集厂商和型号信息
		item.Company = strings.Join(fields[3:5], " ")
		item.ModelName = strings.Join(fields[5:], " ")
	}
	return nil
}

// populateNICStatusAndSpeed 填充网卡列表中的网卡速率信息
func (c *bootosC) populateNICStatusAndSpeed(item *collector.NICItem) (err error) {
	output, err := c.Base.ExecByShell("ethtool", item.Name)
	if err != nil {
		return err
	}
	if strings.Contains(string(output), "10000base") {
		item.Speed = "10Gbit/s"
	} else if strings.Contains(string(output), "1000base") {
		item.Speed = "1Gbit/s"
	} else if strings.Contains(string(output), "100base") {
		item.Speed = "100Mbit/s"
	} else if strings.Contains(string(output), "10base") {
		item.Speed = "10Mbit/s"
	}

	//获取网卡的连接状态信息
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "Link detected:") {
			item.Link = c.extractValue(line, colonSeparator)
		}
	}
	return nil
}
