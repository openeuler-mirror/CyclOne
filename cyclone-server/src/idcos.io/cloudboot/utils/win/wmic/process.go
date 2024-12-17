package wmic

import (
	"fmt"
	"strings"

	"idcos.io/cloudboot/logger"
	winutil "idcos.io/cloudboot/utils/win"
)

// IsProcessRunning 判断指定名称的进程是否存在
func IsProcessRunning(log logger.Logger, processName string) (running bool, err error) {
	output, err := processGet(log, "Caption")
	if err != nil {
		return false, err
	}
	return strings.Contains(string(output), processName), nil // TODO 需要一行一行遍历判断是否存在该进程
}

// 示例 wmic process get Caption,Name
func processGet(log logger.Logger, properties ...string) (output []byte, err error) {
	cmdAndArgs := fmt.Sprintf("wmic process get %s", strings.Join(properties, ","))
	output, err = winutil.ExecOutput(log, cmdAndArgs)
	if err != nil {
		if log != nil {
			log.Errorf("Exec %q err: %s\noutput:\n%s", cmdAndArgs, err, string(output))
		}
		return nil, err
	}
	if log != nil {
		log.Infof("%s ==>\n%s", cmdAndArgs, string(output))
	}
	return output, nil
}
