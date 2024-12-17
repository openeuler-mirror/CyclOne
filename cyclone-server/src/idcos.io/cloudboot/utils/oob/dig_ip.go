package oob

import (
	"fmt"
	"regexp"
	"strings"

	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils/sh"
)

//DigDomain 获取IP的dig命令 dig +short dns @dns_server_ip
var DigDomain = "dig +short %s"

func DigCmd(sn, hostname string, log logger.Logger, repo model.Repo) []string {
	digCmd := make([]string, 0)
	if device, err := repo.GetDeviceBySN(sn); err == nil && device.ServerRoomID != 0 {
		if len(middleware.MapDistributeNode.MDistribute[device.ServerRoomID]) != 0 {
			for _, nodeIP := range middleware.MapDistributeNode.MDistribute[device.ServerRoomID] {
				digCmd = append(digCmd, fmt.Sprintf(DigDomain, hostname)+" @"+nodeIP)
			}
			return digCmd
		} else {
			log.Errorf("机房ID:%d没有关联到NodeIP", device.ServerRoomID)
		}
	}
	digCmd = append(digCmd, fmt.Sprintf(DigDomain, hostname)+" @localhost")
	return digCmd
}

//将带外注册域名转成IP
func TransferHostname2IP(log logger.Logger, repo model.Repo, sn string, hostnames []string) string {
	// 查看ip
	// 用生成的不同的hostname重试
	digCmd := ""
	for _, hostname := range hostnames {
		cmds := DigCmd(sn, hostname, log, repo)
		for _, digCmd = range cmds {
			output, err := sh.ExecOutputWithLog(log, digCmd)
			if err != nil {
				log.Warnf("oob ip query cmd:%s fail", digCmd)
				continue
			}
			ipr, _ := regexp.Compile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)
			ipOutput := strings.TrimSpace(string(output))
			if ipr.MatchString(ipOutput) {
				// 更新OOBIP到数据库
				if device, err := repo.GetDeviceBySN(sn); err == nil {
					if ipOutput != device.OOBIP {
						_, _ = repo.UpdateDeviceBySN(&model.Device{
							SN:    sn,
							OOBIP: ipOutput,
						})
					}
				}
				return ipOutput
			}
		}
	}
	log.Errorf("oob ip query cmd: %s fail after 2 attempts", digCmd)
	return ""
}
