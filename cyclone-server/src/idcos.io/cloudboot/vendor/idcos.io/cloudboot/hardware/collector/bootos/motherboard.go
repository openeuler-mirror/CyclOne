package bootos

import (
	"strings"

	"idcos.io/cloudboot/hardware/collector"
)

// Motherboard 采集并返回当前设备的主板信息
func (c *bootosC) Motherboard() (*collector.Motherboard, error) {
	var err error
	var board collector.Motherboard

	board.SerialNumber, err = c.boardSN()
	if err != nil {
		return nil, err
	}
	board.Manufacturer, err = c.boardManufacturer()
	if err != nil {
		return nil, err
	}
	board.ProductName, err = c.boardProductName()
	if err != nil {
		return nil, err
	}
	return &board, nil
}

// boardSN 返回主板序列号
func (c *bootosC) boardSN() (sn string, err error) {
	output, err := c.Base.ExecByShell("dmidecode", "-s", "baseboard-serial-number")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// boardManufacturer 返回主板厂商名
func (c *bootosC) boardManufacturer() (name string, err error) {
	output, err := c.Base.ExecByShell("dmidecode", "-s", "baseboard-manufacturer")
	name = strings.TrimSpace(string(output))
	if name != "" {
		return name, err
	}
	// dmidecode | grep -A16 "System Information$" | grep "Vendor:" | sed -n '1p'
	output, err = c.Base.ExecByShell("dmidecode", "|", "grep", "-A16", `"System Information$"`, "|", "grep", `"Vendor:"`, "|", "sed", "-n", `'1p'`)
	return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(string(output)), "Vendor:")), err
}

// boardProductName 返回主板产品名
func (c *bootosC) boardProductName() (name string, err error) {
	output, err := c.Base.ExecByShell("dmidecode", "-s", "baseboard-product-name")
	name = strings.TrimSpace(string(output))
	if name != "" {
		return name, err
	}
	// dmidecode | grep -A16 "System Information$" | grep "Product Name:" | sed -n '1p'
	output, err = c.Base.ExecByShell("dmidecode", "|", "grep", "-A16", `"System Information$"`, "|", "grep", `"Product Name:"`, "|", "sed", "-n", `'1p'`)
	return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(string(output)), "Product Name:")), err
}
