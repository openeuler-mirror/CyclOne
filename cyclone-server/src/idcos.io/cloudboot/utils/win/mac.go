package win

import (
	"regexp"
	"strings"

	"idcos.io/cloudboot/logger"
)

// MacAddress 通过本地命令查询当前设备mac地址
func MacAddress(log logger.Logger) (addr string, err error) {
	cmdArgs := `wmic nic where "NetConnectionStatus=2" get MACAddress /VALUE`
	outputUTF8, err := ExecOutputWithLog(log, cmdArgs)
	if err != nil {
		return "", err
	}

	results := regexp.MustCompile(`(?i)MACAddress=(\S+)`).FindStringSubmatch(string(outputUTF8))
	if len(results) != 2 {
		return "", nil
	}
	return strings.TrimSpace(results[1]), nil
}
