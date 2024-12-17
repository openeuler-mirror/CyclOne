package wmic

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"idcos.io/cloudboot/logger"
	winutil "idcos.io/cloudboot/utils/win"
)

// DNSSetting 查询DNS设置
func DNSSetting(log logger.Logger) (setting string, err error) {
	cmdAndArgs := "wmic nicconfig where IPEnabled=True get DNSServerSearchOrder /value"
	output, err := winutil.ExecOutput(log, cmdAndArgs)
	if err != nil {
		if log != nil {
			log.Errorf("Exec %q err: %s\noutput:\n%s", cmdAndArgs, err, string(output))
		}
		return "", err
	}
	if log != nil {
		log.Infof("%s ==>\n%s", cmdAndArgs, string(output))
	}
	// 示例：DNSServerSearchOrder={"192.168.0.1","10.0.0.1"}
	rd := bufio.NewReader(bytes.NewBuffer(output))
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
		if strings.Contains(line, "DNSServerSearchOrder=") {
			pair := strings.Split(line, "=")
			if len(pair) < 2 {
				return "", nil
			}
			return strings.TrimSpace(pair[1]), nil
		}
	}
	return "", nil
}
