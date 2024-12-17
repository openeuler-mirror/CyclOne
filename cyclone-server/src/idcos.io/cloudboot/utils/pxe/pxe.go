package pxe

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// dir 生成PXE文件的目录
const dir = "/var/lib/tftpboot"

// Dir 返回PXE文件根目录
func Dir() string {
	return dir
}

// RemoveFile 删除设备PXE文件
func RemoveFile(mac string) error {
	return os.Remove(filepath.Join(Dir(), fmt.Sprintf("01-%s", strings.Replace(mac, ":", "-", -1))))
}

// GenFile 生成PXE文件
func GenFile(mac string, content []byte) (filename string, err error) {
	if err = os.MkdirAll(Dir(), 0755); err != nil {
		return "", err
	}
	filename = filepath.Join(Dir(), fmt.Sprintf("01-%s", strings.Replace(mac, ":", "-", -1)))
	return filename, ioutil.WriteFile(filename, content, 0644)
}
