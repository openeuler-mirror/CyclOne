package wmic

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"idcos.io/cloudboot/logger"
	winutil "idcos.io/cloudboot/utils/win"
)

var naturalNumReg = regexp.MustCompile("^[0-9]+$") // 自然数正则表达式

var (
	// ErrDiskNotFound 指定的磁盘不存在错误
	ErrDiskNotFound = errors.New("disk not found")
)

// DiskSize 获取指定磁盘(如)的容量大小(byte)。
// 若指定的磁盘不存在，则返回ErrDiskNotFound错误。
func DiskSize(log logger.Logger, disk string) (size int64, err error) {
	output, err := logicalDiskGet(log, disk, "Size")
	if err != nil {
		return 0, err
	}
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
			return 0, err
		}
		line = strings.TrimSpace(line)
		if strings.Contains(line, "No Instance(s) Available") {
			return 0, ErrDiskNotFound
		}
		if naturalNumReg.MatchString(line) {
			return strconv.ParseInt(line, 10, 64)
		}
	}
	return size, nil
}

// 示例：wmic logicaldisk where name="C:" get Size,Name
func logicalDiskGet(log logger.Logger, disk string, properties ...string) (output []byte, err error) {
	cmdAndArgs := fmt.Sprintf("wmic logicaldisk where name=%q get %s", disk, strings.Join(properties, ","))
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
