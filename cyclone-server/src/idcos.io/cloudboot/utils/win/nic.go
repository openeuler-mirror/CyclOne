package win

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/utils"
)

// RenameNICName 修改网卡名
func RenameNICName(log logger.Logger, oldName string, newName string) error {
	if log != nil {
		log.Infof("Reset network connection name from %q to %q", oldName, newName)
	}
	cmdAndArgs := fmt.Sprintf(`netsh interface set interface name="%s" newname="%s"`, oldName, newName)
	output, err := ExecOutput(log, utils.UTF82GBK(cmdAndArgs))
	if err != nil {
		if log != nil {
			log.Errorf("Exec %q err: %s\noutput:\n%s", cmdAndArgs, err, string(output))
		}
		return err
	}
	outputUTF8 := utils.GBK2UTF8(string(output))
	if log != nil {
		log.Infof("%s ==>\n%s", cmdAndArgs, outputUTF8)
	}
	return nil
}

// DisableNIC 禁用本地连接（网卡）
func DisableNIC(log logger.Logger, name string) (err error) {
	return enableOrDisableNIC(log, name, false)
}

// DisableNICs 批量禁用本地连接（网卡）
func DisableNICs(log logger.Logger, names ...string) (err error) {
	log.Infof("Start disabling NICs: %v", names)
	for _, name := range names {
		if err = DisableNIC(log, name); err != nil {
			return err
		}
	}
	return nil
}

// EnableNIC 启用本地连接（网卡）
func EnableNIC(log logger.Logger, name string) (err error) {
	return enableOrDisableNIC(log, name, true)
}

// EnableNICs 批量启用本地连接（网卡）
func EnableNICs(log logger.Logger, names ...string) (err error) {
	for _, name := range names {
		if err = EnableNIC(log, name); err != nil {
			return err
		}
	}
	return nil
}

func enableOrDisableNIC(log logger.Logger, name string, enable bool) (err error) {
	stat := "disabled"
	if enable {
		stat = "enabled"
	}
	cmdAndArgs := fmt.Sprintf(`netsh interface set interface name="%s" admin=%s`, name, stat)
	output, err := ExecOutput(log, utils.UTF82GBK(cmdAndArgs))
	if err != nil {
		if log != nil {
			log.Errorf("Exec %q err: %s\noutput:\n%s", cmdAndArgs, err, string(output))
		}
		return err
	}
	outputUTF8 := utils.GBK2UTF8(string(output))
	if log != nil {
		log.Infof("%s ==>\n%s", cmdAndArgs, outputUTF8)
	}
	return nil
}

// GetTeam0Name 通过命令行查询team0网卡名称（对应的网络连接名称）
func GetTeam0Name(log logger.Logger) (name string, err error) {
	outUTF8, err := ExecOutputWithLog(log, `wmic nic where (Name like '%%Team0%%') get NetConnectionID /value`)
	if err != nil {
		return "", err
	}
	rd := bufio.NewReader(bytes.NewBuffer(outUTF8))
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			if log != nil {
				log.Error(err)
			}
			return "", err
		}
		line = strings.TrimSpace(line)
		if strings.Contains(line, "NetConnectionID=") {
			return strings.TrimPrefix(line, "NetConnectionID="), nil
		}
	}
	return "", nil
}
