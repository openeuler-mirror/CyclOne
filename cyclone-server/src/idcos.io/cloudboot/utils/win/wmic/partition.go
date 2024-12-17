package wmic

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"idcos.io/cloudboot/logger"
	winutil "idcos.io/cloudboot/utils/win"
)

// ErrKeywordsNotFound 标准输出中用于解析用的关键字不存在
var ErrKeywordsNotFound = errors.New("keywords not found")

// PartitionSizeByIndex 根据索引查询分区的大小（单位字节）
func PartitionSizeByIndex(log logger.Logger, index int) (byteSize int64, err error) {
	cmdAndArgs := fmt.Sprintf("wmic partition where Index=%d get Size /value", index)
	output, err := winutil.ExecOutput(log, cmdAndArgs)
	if err != nil {
		if log != nil {
			log.Errorf("Exec %q err: %s\noutput:\n%s", cmdAndArgs, err, string(output))
		}
		return 0, err
	}
	if log != nil {
		log.Infof("%s ==>\n%s", cmdAndArgs, string(output))
	}

	outstr := strings.TrimSpace(string(output))
	idx := strings.Index(outstr, "Size=")
	if idx < 0 {
		return 0, ErrKeywordsNotFound
	}

	sSize := outstr[idx+len("Size="):]
	if byteSize, err = strconv.ParseInt(sSize, 10, 64); err != nil && log != nil {
		log.Errorf("Parse string value %q to integer value error: %s", sSize, err)
	}
	return byteSize, err
}
