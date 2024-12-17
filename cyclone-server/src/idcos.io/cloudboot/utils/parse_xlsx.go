package utils

import (
	"strconv"

	"github.com/tealeg/xlsx"
)

//ParseDataFromXLSX 从XLSX格式的配置文件中解析数据
//选择有两种1. ToSlice 解析数据到一个[][][]string中，
//2. ToSliceUnmerged 解析数据到一个[][][]string中，与ToSlice不同的是:
// Example: table where A1:A2 merged.
// | 01.01.2011 | Bread | 20 |
// |            | Fish  | 70 |
// This sheet will be converted to the slice:
// [[01.01.2011 Bread 20]
// 	[01.01.2011 Fish  70] ]
// 为什么是一个三维数组，主要是因为xlsx有sheet，第一维表示一张表
// 现在的要求是一个文件只能有一张表。
func ParseDataFromXLSX(filePath string) ([][]string, error) {
	f, err := xlsx.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	output := [][][]string{}

	for k, sheet := range f.Sheets {
		if k > 0 {
			break
		}
		s := [][]string{}
		var length int = 0
		//统计第一行的列数
		for _, cell := range sheet.Rows[0].Cells {
			str, _ := cell.FormattedValue()
			if err != nil {
				// Recover from strconv.NumError if the value is an empty string,
				// and insert an empty string in the output.
				if numErr, ok := err.(*strconv.NumError); ok && numErr.Num == "" {
					str = ""
				} else {
					return nil, err
				}
			}
			if str != "" {
				length += 1
			}
		}
		//根据第一行的列数进行数据填充，如果不够用空补，如果超过不管
		for _, row := range sheet.Rows {
			if row == nil {
				continue
			}
			var nodata bool = false
			r := []string{}
			for _, cell := range row.Cells {
				str, err := cell.FormattedValue()
				if err != nil {
					// Recover from strconv.NumError if the value is an empty string,
					// and insert an empty string in the output.
					if numErr, ok := err.(*strconv.NumError); ok && numErr.Num == "" {
						str = ""
					} else {
						return nil, err
					}
				}
				r = append(r, str)
			}

			for i := len(row.Cells); i < length; i++ {
				r = append(r, "")
			}
			for i := 0; i < length; i++ {
				if r[i] != "" {
					nodata = true
				}
			}
			if nodata {
				s = append(s, r)
			}
		}
		output = append(output, s)
	}
	return output[0], nil
}

const FileDevice = "device"
const FileOrder = "order"
const FileOOB = "oob_info"
const FileDeviceSendMail = "device_sendmail"
const FileIP = "ip"

var (
	mTitle = map[string][]string{
		FileDevice: []string{
			`ID`,
			`固资编号`,
			`序列号`,
			`厂商`,
			`型号`,
			`CPU架构`,
			`用途`,
			`类型`,
			`数据中心`,
			`机房管理单元`,
			`机架编号`,
			`机位编号`,
			`带外IP`,
			`带外用户`,
			`带外密码`,
			`电源状态`,
			`硬件备注`,
			`RAID备注`,
			`启用时间`,
			`上架时间`,
			`运营状态`,
			`内网IP`,
			`外网IP`,
			`操作系统`,
		},
		FileOrder: []string{
			`订单号`,
			`逻辑区(物理区域)`,
			`用途`,
			`设备类型`,
			`数量`,
			`数据中心`,
			`机房管理单元`,
			`预占机位`,
			`预计到货时间`,
		},
		FileOOB: []string{
			`固资编号`,
			`序列号`,
			`内网IP`,
			`带外IP`,
			`带外用户`,
			`带外密码`,
		},
		FileDeviceSendMail: []string{
			`固资编号`,
			`序列号`,
			`厂商`,
			`型号`,
			`CPU架构`,
			`用途`,
			`设备类型`,
			`硬件备注`,
			`启用时间`,
			`运营状态`,
			`数据中心`,
			`机房管理单元`,
			`物理区域`,
			`机架`,
			`机位`,
			`负责人`,
			`维保截止日期`,
			`维保状态`,
		},
		FileIP: []string{
			`IP地址`,
			`网段名称`,
			`网段网关`,
			`网段掩码`,
			`网段类别`,
			`IP版本`,
			`是否使用`,
			`IP作用范围`,
			`IP用途`,
			`关联SN`,
			`关联固资编号`,
		},
	}
)

//WriteToXLSX 将数据写入到指定文件的xlsx文件中, 现在的要求是一个文件只能有一张表。
func WriteToXLSX(which string, data [][]string) (*xlsx.File, error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		return nil, err
	}
	//Add title
	row = sheet.AddRow()
	for _, vc := range mTitle[which] {
		cell = row.AddCell()
		cell.Value = vc
	}

	//Add data
	for _, vr := range data {
		row = sheet.AddRow()
		for _, vc := range vr {
			cell = row.AddCell()
			cell.Value = vc
		}
	}

	return file, nil
}
