package centos6

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"idcos.io/cloudboot/utils/sh"
)

// dhcpFilename DHCP配置文件路径
const dhcpFilename = "/etc/dhcp/dhcpd.conf"

const dhcpItemFormat = `
### begin %s ###
host %s {
	hardware ethernet %s;
	filename "BOOTX64.efi";
}
### end %s ###`

// AddOneToDHCP 返回添加了配置项的dhcp配置
func AddOneToDHCP(sn, mac string, src []byte) (dst []byte) {
	buf := bytes.NewBuffer(src)
	buf.WriteByte('\n')
	buf.WriteString(strings.TrimSpace(fmt.Sprintf(dhcpItemFormat, sn, sn, mac, sn)))
	return buf.Bytes()
}

// RmOneFromDHCP 返回移除了目标配置项的dhcp配置
func RmOneFromDHCP(sn string, src []byte) (dst []byte) {
	txt, sBegin, sEnd := string(src), fmt.Sprintf("### begin %s ###", sn), fmt.Sprintf("### end %s ###", sn)
	beginIdx, endIdx := strings.Index(txt, sBegin), strings.Index(txt, sEnd)

	if beginIdx < 0 || endIdx < 0 {
		return src
	}
	return []byte(txt[:beginIdx] + txt[endIdx+len(sEnd):])
}

// OverwriteDHCP 覆写dhcp配置并重启dhcpd进程
func OverwriteDHCP(dhcp []byte) (err error) {
	if err = ioutil.WriteFile(dhcpFilename, dhcp, 0644); err != nil {
		return err
	}
	_, err = RestartDHCP()
	return err
}

// RestartDHCP 使用systemctl重启dhcpd进程
func RestartDHCP() (out []byte, err error) {
	return sh.ExecOutput(nil, "systemctl restart dhcpd")
}
