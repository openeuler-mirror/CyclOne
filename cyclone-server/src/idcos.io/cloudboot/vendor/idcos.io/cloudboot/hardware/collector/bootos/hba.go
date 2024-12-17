package bootos

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"idcos.io/cloudboot/hardware/collector"
)

// HBA 采集并返回当前设备的HBA卡信息。若当前设备不包含HBA卡，则返回值都为nil。
func (c *bootosC) HBA() (*collector.HBA, error) {
	if !c.hasHBA() {
		return nil, nil
	}

	hosts, err := c.hbaHosts()
	if err != nil {
		return nil, err
	}

	hba := collector.HBA{
		Items: make([]collector.HBAItem, 0, len(hosts)),
	}
	for i := range hosts {
		var item collector.HBAItem
		item.Host = hosts[i]
		item.FirmwareVersion, _ = c.hbaFwVersion(hosts[i])
		item.WWPNs, _ = c.hbaWWPNs(hosts[i])
		item.WWNNs, _ = c.hbaWWNNs(hosts[i])
		hba.Items = append(hba.Items, item)
	}
	return &hba, nil
}

// hbaHosts 返回/sys/class/fc_host/目录下第一层的所有子目录名称
func (c *bootosC) hbaHosts() (hosts []string, err error) {
	flist, err := ioutil.ReadDir("/sys/class/fc_host/")
	if err != nil {
		if log := c.Base.GetLog(); log != nil {
			log.Error(err)
		}
		return nil, err
	}
	for i := range flist {
		hosts = append(hosts, flist[i].Name())
	}
	return hosts, nil
}

// hasHBA 检查当前设备是否含有HBA卡
func (c *bootosC) hasHBA() bool {
	output, _ := c.Base.ExecByShell("lspci", "|", "grep", `"Fibre Channel"`)
	return strings.TrimSpace(string(output)) != ""
}

// hbaFwVersion 采集指定HBA卡的固件版本号
func (c *bootosC) hbaFwVersion(host string) (version string, err error) {
	// /sys/class/scsi_host/$host/fw_version
	b, err := ioutil.ReadFile(filepath.Join("/sys", "class", "scsi_host", host, "fw_version"))
	if err != nil {
		if log := c.Base.GetLog(); log != nil {
			log.Error(err)
		}
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

func (c *bootosC) hbaWWPNs(host string) (wwpns []string, err error) {
	b, err := ioutil.ReadFile(filepath.Join("/sys", "class", "fc_host", host, "port_name"))
	if err != nil {
		if log := c.Base.GetLog(); log != nil {
			log.Error(err)
		}
		return nil, err
	}
	return []string{
		strings.TrimSpace(string(b)), // TODO 理论上一块HBA卡有多个wwpn，暂时不清楚此文件内容格式。
	}, nil
}

func (c *bootosC) hbaWWNNs(host string) (wwnns []string, err error) {
	b, err := ioutil.ReadFile(filepath.Join("/sys", "class", "scsi_host", host, "device", "fc_host", host, "node_name"))
	if err != nil {
		if log := c.Base.GetLog(); log != nil {
			log.Error(err)
		}
		return nil, err
	}
	return []string{
		strings.TrimSpace(string(b)), // TODO 理论上一块HBA卡有一个或者多个wwnn，暂时不清楚此文件内容格式。
	}, nil
}
