package win

import (
	"regexp"
	"strings"

	"idcos.io/cloudboot/logger"
)

// PhysicalMachineSN 查询物理机SN
func PhysicalMachineSN(log logger.Logger) (sn string, err error) {
	cmdArgs := `wmic bios get SerialNumber /VALUE`
	outputUTF8, err := ExecOutputWithLog(log, cmdArgs)
	if err != nil {
		return "", err
	}
	results := regexp.MustCompile(`SerialNumber=(\S+)`).FindStringSubmatch(string(outputUTF8))
	if len(results) != 2 {
		return "", nil
	}
	return strings.TrimSpace(results[1]), nil
}
