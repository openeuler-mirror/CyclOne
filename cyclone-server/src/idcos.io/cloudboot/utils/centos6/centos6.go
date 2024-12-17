package centos6

import (
	"io/ioutil"
	"strings"

	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils/pxe"
)

// IsPXEUEFI 判断设备是否是PXE方式安装centos6(uefi)
func IsPXEUEFI(log logger.Logger, repo model.Repo, sn string) bool {
	sett, _ := repo.GetDeviceSettingBySN(sn)
	if sett == nil || sett.InstallType != model.InstallationPXE {
		return false
	}

	osTemplate, _ := repo.GetSystemTemplateByID(sett.SystemTemplateID)

	return osTemplate != nil &&
		osTemplate.BootMode == model.BootModeUEFI &&
		strings.Contains(strings.ToLower(osTemplate.Name), "centos6")
}

// DropConfigurations 移除和centos6安装相关的特殊配置
func DropConfigurations(log logger.Logger, repo model.Repo, sn string) error {
	log.Infof("============== Remove configurations for %s(CentOS 6 + PXE + UEFI) ==============", sn)

	// 移除dhcp配置文件中相关配置项并重启dhcpd
	dhcpSRC, err := ioutil.ReadFile(dhcpFilename)
	if err == nil && len(dhcpSRC) > 0 {
		log.Info("After the OS installation is complete, try to modify the dhcp configuration file.")
		if err = OverwriteDHCP(RmOneFromDHCP(sn, dhcpSRC)); err != nil {
			log.Errorf("Restart dhcpd error: %s", err.Error())
			log.Warn("Try to rollback the dhcp configuration file")
			_ = OverwriteDHCP(dhcpSRC) // 尝试回滚配置
		}
	}
	// 删除PXE文件
	if mac, _ := getBootOSMacBySN(log, repo, sn); mac != "" {
		log.Info("Try to remove the PXE file")
		if err = pxe.RemoveFile(mac); err != nil {
			log.Error(err)
		}
	}
	return err
}

func getBootOSMacBySN(log logger.Logger, repo model.Repo, sn string) (mac string, err error) {
	dev, err := repo.GetDeviceBySN(sn)
	if err != nil {
		return "", err
	}
	return dev.BootOSMac, nil
}
