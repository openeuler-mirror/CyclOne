package service

import "strconv"

// ExportedDevices 导出设备信息集合
type ExportedDevices []*DevicePageResp

// ToTableRecords 生成用于表格显示的二维字符串切片
func (items ExportedDevices) ToTableRecords() (records [][]string) {
	records = make([][]string, 0, len(items))

	for i := range items {
		idcName := ""
		if items[i].IDC != nil {
			idcName = items[i].IDC.Name
		}
		serverRoomName := ""
		if items[i].ServerRoom != nil {
			serverRoomName = items[i].ServerRoom.Name
		}
		cabinetNum := ""
		if items[i].ServerCabinet != nil {
			cabinetNum = items[i].ServerCabinet.Number
		}
		usiteNum := ""
		if items[i].ServerUSite != nil {
			usiteNum = items[i].ServerUSite.Number
		}
		records = append(records, []string{
			strconv.Itoa(int(items[i].ID)),
			items[i].FixedAssetNum,
			items[i].SN,
			items[i].Vendor,
			items[i].Model,
			items[i].Arch,
			items[i].Usage,
			items[i].Category,
			idcName,
			serverRoomName,
			cabinetNum,
			usiteNum,
			items[i].OOBIP,
			items[i].OOBUser,
			items[i].OOBPassword,
			items[i].PowerStatus,
			items[i].HardwareRemark,
			items[i].RAIDRemark,
			items[i].StartedAt,
			items[i].OnShelveAt,
			items[i].OperationStatus,
			items[i].IntranetIP,
			items[i].ExtranetIP,
			items[i].OS,
		})
	}
	return records
}

// ExportOOBInfo 导出带外信息
type ExportOOBInfo []*DevicePageResp

// ToTableRecords 生成用于表格显示的二维字符串切片
func (items ExportOOBInfo) ToTableRecords() (records [][]string) {
	records = make([][]string, 0, len(items))

	for i := range items {
		records = append(records, []string{
			items[i].FixedAssetNum,
			items[i].SN,
			items[i].IntranetIP,
			items[i].OOBIP,
			items[i].OOBUser,
			items[i].OOBPassword,
		})
	}
	return records
}

// ExportIPInfo 导出IP信息
type ExportIPInfo []*IPSPage

// ToTableRecords 生成用于表格显示的二维字符串切片
func (items ExportIPInfo) ToTableRecords() (records [][]string) {
	records = make([][]string, 0, len(items))

	for i := range items {
		records = append(records, []string{
			items[i].IP,
			items[i].IPNetwork.CIDR,
			items[i].IPNetwork.Gateway,
			items[i].IPNetwork.Netmask,
			items[i].IPNetwork.Category,
			items[i].IPNetwork.Version,
			items[i].IsUsed,
			items[i].Category,
			items[i].Scope,
			items[i].SN,
			items[i].FixedAssetNumber,
		})
	}
	return records
}