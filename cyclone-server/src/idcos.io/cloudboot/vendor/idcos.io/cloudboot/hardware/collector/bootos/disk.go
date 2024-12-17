package bootos

import (
	"regexp"
	"strconv"
	"strings"

	"idcos.io/cloudboot/hardware/collector"
	"idcos.io/cloudboot/utils/bytes"
)

// Disk 采集并返回当前设备的逻辑磁盘信息
func (c *bootosC) Disk() (*collector.Disk, error) {
	var disk collector.Disk
	output, err := c.Base.ExecByShell("fdisk", "-l")
	if err != nil {
		return nil, err
	}
	disk.Items = c.parseDiskItems(output)
	disk.TotalSize, disk.TotalSizeGB = c.calcDiskTotalSize(disk.Items)
	return &disk, nil
}

var diskReg = regexp.MustCompile(`Disk\s+(.*):(.*)(\d+)\s+bytes`)

func (c *bootosC) parseDiskItems(output []byte) (items []collector.DiskItem) {
	lines := diskReg.FindAllString(string(output), -1)
	for i := range lines {
		pair := strings.SplitN(lines[i], colonSeparator, -1)
		items = append(items, collector.DiskItem{
			Name: strings.TrimSpace(strings.TrimPrefix(pair[0], "Disk")), // /dev/sda
			Size: strings.TrimSpace(pair[1]),                             // 599.6 GB, 599550590976 bytes
		})
	}
	return items
}

// calcDiskTotalSize 计算逻辑磁盘总容量。
// bsize 单位为Byte的容量值。
// gbsize 单位为GB的容量值。
func (c *bootosC) calcDiskTotalSize(items []collector.DiskItem) (bsize int64, gbsize int) {
	for i := range items {
		begin := strings.Index(items[i].Size, ",")
		end := strings.Index(items[i].Size, "bytes")
		if begin < 0 || begin >= end {
			continue
		}
		ssize := strings.TrimSpace(items[i].Size[begin+1 : end])
		size, _ := strconv.ParseInt(ssize, 10, 64)
		bsize += size
	}
	return bsize, bytes.Byte2GBRounding(bytes.Byte(bsize))
}
