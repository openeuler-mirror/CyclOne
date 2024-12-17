package win

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/utils"
	mybytes "idcos.io/cloudboot/utils/bytes"
)

// PartConfiguration 分区配置信息
type PartConfiguration struct {
	Size       string // 字节数或者free
	FSType     string
	Mountpoint string
}

// DiskPartConfiguration 磁盘分区配置信息
type DiskPartConfiguration struct {
	Disk       int // 磁盘编号，如0、1、2...
	GPT        bool
	Partitions []PartConfiguration
}

var (
	numReg   = regexp.MustCompile("^\\d+$")
	alphaReg = regexp.MustCompile("^[A-Za-z]$")
)

func (conf *DiskPartConfiguration) setup(log logger.Logger) (err error) {
	size, err := conf.diskSize(log)
	if err != nil {
		return err
	}
	conf.GPT = size.GE(2 * mybytes.TB) // 约定: 磁盘容量大于等于2T使用GPT方式分区
	return nil
}

var (
	// ErrDiskNotSelected 磁盘未选中
	ErrDiskNotSelected = errors.New("disk is not selected")
)

// diskSize 返回当前磁盘的容量
func (conf *DiskPartConfiguration) diskSize(log logger.Logger) (size mybytes.Byte, err error) {
	output, err := conf.exec(log, fmt.Sprintf("select disk %d\r\nlist disk", conf.Disk))
	if err != nil {
		return mybytes.Byte(0), err
	}
	var started bool
	rd := bufio.NewReader(bytes.NewBuffer(output))
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return size, err
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "-----") {
			started = true
			continue
		}
		if !started || !strings.HasPrefix(line, "*") { // 选中的目标磁盘记录会以"*"开头
			continue
		}

		line = strings.TrimSpace(strings.TrimPrefix(line, "*"))
		fields := strings.Fields(line)
		if len(fields) < 5 {
			break
		}
		return mybytes.Parse2Byte(strings.TrimSpace(fields[3]), strings.TrimSpace(fields[4]))
	}
	return mybytes.Byte(0), ErrDiskNotSelected
}

type volume struct {
	ID     string
	Letter string
}

// listVolumes 列出当前磁盘的所有卷(不包含已移除的卷和系统保留卷)
func (conf *DiskPartConfiguration) listVolumes(log logger.Logger) (vols []*volume, err error) {
	output, err := conf.exec(log, "list vol")
	if err != nil {
		return nil, err
	}
	var started bool
	rd := bufio.NewReader(bytes.NewBuffer(output))
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "-----") {
			started = true
			continue
		}
		if !started {
			continue
		}
		line = strings.TrimSpace(strings.TrimPrefix(line, "*"))
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		id, letter := strings.TrimSpace(fields[1]), strings.TrimSpace(fields[2])
		if numReg.MatchString(id) && alphaReg.MatchString(letter) {
			vols = append(vols, &volume{
				ID:     id,
				Letter: letter,
			})
		}
	}
	return vols, nil
}

// rmVolumeLetter 移除指定卷的盘符
func (conf *DiskPartConfiguration) rmVolumeLetter(log logger.Logger, vol *volume) (err error) {
	// 选中目标卷，移除目标卷盘符
	_, err = conf.exec(log, fmt.Sprintf("select vol %s\r\nremove letter=%s", vol.ID, vol.Letter))
	return err
}

// exec 将待执行命令写入文件，然后通过'diskpart /s $filepath'方式执行命令。
func (conf *DiskPartConfiguration) exec(log logger.Logger, cmdAndArgs string) (output []byte, err error) {
	filename := fmt.Sprintf("diskpart_%s.txt", utils.UUID())
	if err = ioutil.WriteFile(filename, []byte(cmdAndArgs), 0777); err != nil {
		return nil, err
	}
	defer os.Remove(filename)

	output, err = exec.Command("diskpart", "/s", filename).Output()
	log.Infof("diskpart /s %s\n%s==>\n%s", filename, cmdAndArgs, output)
	return output, err
}

// CMD 返回创建分区的命令内容
func (conf *DiskPartConfiguration) CMD() []byte {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("select disk %d\r\n", conf.Disk))
	buf.WriteString("clean\r\n")
	if conf.GPT {
		buf.WriteString("convert gpt\r\n")
	}
	for i := range conf.Partitions {
		if strings.ToLower(conf.Partitions[i].Size) == "free" {
			buf.WriteString("create partition primary\r\n")
		} else {
			buf.WriteString(fmt.Sprintf("create partition primary size=%s\r\n", conf.Partitions[i].Size))
		}
		buf.WriteString(fmt.Sprintf("format quick fs=%s\r\n", conf.Partitions[i].FSType))
		buf.WriteString(fmt.Sprintf("assign letter=%s\r\n", conf.Partitions[i].Mountpoint))
		if strings.ToUpper(conf.Partitions[i].Mountpoint) == "C" && !conf.GPT {
			buf.WriteString("active\r\n")
		}
	}
	return buf.Bytes()
}

// Apply 通过本地调用diskpart命令执行磁盘分区
func (conf *DiskPartConfiguration) Apply(log logger.Logger) (output []byte, err error) {
	if err = conf.setup(log); err != nil {
		return nil, err
	}

	vols, _ := conf.listVolumes(log)
	for i := range vols {
		for x := range conf.Partitions {
			if strings.ToUpper(vols[i].Letter) == strings.ToUpper(conf.Partitions[x].Mountpoint) {
				_ = conf.rmVolumeLetter(log, vols[i]) // 待分配的盘符和现有盘符冲突，移除现有盘符。
			}
		}
	}

	// 执行实际分区
	return conf.exec(log, string(conf.CMD()))
}

// DiskPartConfigurations 磁盘分区集合（注意集合中元素顺序）
type DiskPartConfigurations []DiskPartConfiguration

// Apply 通过本地调用diskpart命令执行磁盘分区
func (confs DiskPartConfigurations) Apply(log logger.Logger) (err error) {
	for i := range confs {
		if _, err = confs[i].Apply(log); err != nil {
			return err
		}
		if i < len(confs)-1 {
			// 微软建议：两次执行diskpart脚本的间隔期 >= 15s，特别是对于写操作。
			// https://technet.microsoft.com/zh-cn/library/cc766465(v=ws.10).aspx
			time.Sleep(time.Duration(15) * time.Second)
		}
	}
	return nil
}

var (
	// ErrMalformedDisk 磁盘命名格式错误
	ErrMalformedDisk = errors.New("malformed disk")
)

// DiskNo 返回磁盘序号。自0始。
func DiskNo(dev string) (no int, err error) {
	dev = strings.ToLower(strings.TrimSpace(dev))
	if dev == "" {
		return 0, ErrMalformedDisk
	}
	if numReg.MatchString(dev) {
		return strconv.Atoi(dev)
	}

	lastC := dev[len(dev)-1]
	if lastC > 'z' || lastC < 'a' {
		return 0, ErrMalformedDisk
	}

	return int(lastC - 'a'), nil
}
