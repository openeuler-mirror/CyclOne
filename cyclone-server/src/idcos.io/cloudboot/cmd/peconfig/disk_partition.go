package main

import "idcos.io/cloudboot/utils/win"

// DiskSlice 磁盘分区集合（针对镜像安装方式）
type DiskSlice []Disk

// ToDiskPartConfigurations 将磁盘分区集合信息转化成磁盘分区配置信息
func (disks DiskSlice) ToDiskPartConfigurations() win.DiskPartConfigurations {
	items := make([]win.DiskPartConfiguration, 0, len(disks))
	for i := range disks {
		diskNo, _ := win.DiskNo(disks[i].Name)
		items = append(items, win.DiskPartConfiguration{
			Disk:       diskNo,
			Partitions: PartitionSlice(disks[i].Partitions).ToPartConfigurations(),
		})
	}
	return items
}

// PartitionSlice 分区信息（针对镜像安装方式）
type PartitionSlice []Partition

// ToPartConfigurations 将分区信息转化为分区配置信息
func (partitions PartitionSlice) ToPartConfigurations() []win.PartConfiguration {
	items := make([]win.PartConfiguration, 0, len(partitions))
	for i := range partitions {
		items = append(items, win.PartConfiguration{
			Size:       partitions[i].Size,
			FSType:     partitions[i].Fstype,
			Mountpoint: partitions[i].Mountpoint,
		})
	}
	return items
}
