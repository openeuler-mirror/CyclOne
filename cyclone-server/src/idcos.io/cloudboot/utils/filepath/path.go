package filepath

import (
	"os/user"
	"strings"
)

const (
	// HomeDirFlag 当前用户家目录标识符
	HomeDirFlag = "~"
)

// Rel2Abs 将~转化为用户家目录
func Rel2Abs(raw string) (string, error) {
	raw = strings.TrimSpace(raw)

	if !strings.HasPrefix(raw, HomeDirFlag) {
		return raw, nil
	}
	user, err := user.Current()
	if err != nil {
		return raw, err
	}
	return strings.Replace(raw, HomeDirFlag, user.HomeDir, 1), nil
}
