package bootos

import (
	"bytes"
	"strings"

	"idcos.io/cloudboot/hardware/collector"
)

// BASE 采集并返回当前设备的基本信息
func (c *bootosC) BASE() (isVM bool, sn, vendor, model, arch string, err error) {
	// 不再支持VM了
	//isVM, err = c.isVM()
	//if err != nil {
	//	return isVM, sn, vendor, model, arch, err
	//}

	sn, err = c.sn(isVM)
	if err != nil {
		return isVM, sn, vendor, model, arch, err
	}

	vendor, _ = c.manufacturer()
	model, _ = c.productName()
	arch, _ = c.arch()

	return isVM, sn, vendor, model, arch, nil
}

func (c *bootosC) sn(isVM bool) (string, error) {
	if isVM {
		output, err := c.Base.ExecByShell("sed", "q", "/sys/class/net/*/address")
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(output)), nil
	}

	output, err := c.Base.ExecByShell("dmidecode", "-s", "system-serial-number", "2>/dev/null", "|", "awk", `'/^[^#]/ { print $1 }'`)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func (c *bootosC) isVM() (bool, error) {
	output, err := c.Base.ExecByShell("dmidecode")
	if err != nil {
		return false, err
	}
	return bytes.Contains(output, []byte("VMware")) ||
		bytes.Contains(output, []byte("VirtualBox")) ||
		bytes.Contains(output, []byte("KVM")) ||
		bytes.Contains(output, []byte("Xen")) ||
		bytes.Contains(output, []byte("Parallels")), nil
}

// productName 返回设备的产品名称
func (c *bootosC) productName() (string, error) {
	// dmidecode -s system-product-name | awk '/^[^#]/'
	output, err := c.Base.ExecByShell("dmidecode", "-s", "system-product-name", "|", "awk", `'/^[^#]/'`)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// arch 采集CPU硬件架构
func (c *bootosC) arch() (string, error) {
	output, err := c.Base.ExecByShell("uname -m")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// manufacturer 返回设备厂商名称
func (c *bootosC) manufacturer() (string, error) {
	// dmidecode -s system-manufacturer | awk '/^[^#]/ { print $1 }'
	output, err := c.Base.ExecByShell("dmidecode", "-s", "system-manufacturer")
	if err != nil {
		return "", err
	}

	name := strings.TrimSpace(string(output))
	nameLower := strings.ToLower(name)

	if nameLower == "" {
		return "", nil
	} else if strings.Contains(nameLower, "dell") {
		return collector.Dell, nil

	} else if strings.Contains(nameLower, "super") && strings.Contains(nameLower, "micro") {
		return collector.Supermicro, nil

	} else if strings.Contains(nameLower, "great") && strings.Contains(nameLower, "wall") {
		return collector.Greatwall, nil

	} else if strings.Contains(nameLower, "h3c") {
		return collector.H3C, nil

	} else if strings.Contains(nameLower, "lenovo") {
		return collector.Lenovo, nil

	} else if strings.Contains(nameLower, "huawei") {
		return collector.Huawei, nil

	} else if strings.Contains(nameLower, "xfusion") {
		return collector.XFusion, nil

	} else if strings.Contains(nameLower, "inspur") {
		return collector.Inspur, nil

	} else if strings.Contains(nameLower, "sugon") {
		return collector.Sugon, nil

	} else if strings.Contains(nameLower, "unis") {
		return collector.UNIS, nil

	} else if strings.Contains(nameLower, "suma") {
		return collector.Suma, nil
	
	} else if strings.Contains(nameLower, "hp") || strings.Contains(nameLower, "hpe") {
		return collector.HP, nil
	} else {
		return strings.Fields(name)[0], nil
	}
}
