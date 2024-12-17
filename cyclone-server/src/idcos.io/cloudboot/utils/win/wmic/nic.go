package wmic

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"idcos.io/cloudboot/logger"
	winutil "idcos.io/cloudboot/utils/win"
)

var (
// ErrNICNotFound nic不存在
// ErrNICNotFound = errors.New("nic not found")
)

// NIC 网卡信息
type NIC struct {
	MacAddr   string `json:"MacAddress"`      // MAC地址
	NetConnID string `json:"NetConnectionID"` // 网卡名
}

// FindNICByMacAddr 根据mac地址查找网卡信息结构体指针
func FindNICByMacAddr(nics []NIC, addr string) (nic *NIC) {
	addr = strings.ToLower(addr)
	for i := range nics {
		if strings.ToLower(nics[i].MacAddr) == addr {
			return &nics[i]
		}
	}
	return nil
}

// RmNIC 从集合中移除指定名称的网卡
func RmNIC(nics []NIC, id string) (news []NIC) {
	if id == "" {
		return nics
	}
	id = strings.ToLower(id)
	for i := range nics {
		if strings.ToLower(nics[i].NetConnID) != id {
			news = append(news, nics[i])
		}
	}
	return news
}

// GetNIC 通过执行本地命令查询网卡信息
func GetNIC(log logger.Logger) (list []NIC, err error) {
	output, err := nicGet(log, "PhysicalAdapter=TRUE", "MACAddress", "NetConnectionID")
	if err != nil {
		return nil, err
	}
	return parseNICGetOutput(log, output)
}

// Get10GNICs 查询万兆网卡列表
func Get10GNICs(log logger.Logger) (list []NIC, err error) {
	output, err := nicGet(log, `PhysicalAdapter=TRUE and Name like '%%10G%%'`, "MACAddress", "NetConnectionID")
	if err != nil {
		return nil, err
	}
	return parseNICGetOutput(log, output)
}

// Get1GNICs 查询千兆网卡列表
func Get1GNICs(log logger.Logger) (list []NIC, err error) {
	output, err := nicGet(log, `PhysicalAdapter=TRUE and not Name like '%%10G%%'`, "MACAddress", "Name", "NetConnectionID") // 假定排查掉所有万兆后，剩下的都是千兆网卡。
	if err != nil {
		return nil, err
	}
	return parseNICGetOutput(log, output)
}

// GetNICByMacAddr 根据mac地址查询网卡信息
func GetNICByMacAddr(log logger.Logger, macAddr string) (nic *NIC, err error) {
	output, err := nicGet(log, fmt.Sprintf(`PhysicalAdapter=TRUE and MACAddress=%q`, macAddr), "MACAddress", "NetConnectionID")
	if err != nil {
		return nil, err
	}
	list, err := parseNICGetOutput(log, output)
	if err != nil {
		return nil, err
	}
	if len(list) <= 0 {
		return nil, nil
	}
	log.Debugf("%s --> %#v", macAddr, list[0])
	return &list[0], nil
}

// parseNICGetOutput 解析形如wmic nic where "xxxx" get MACAddress,NetConnectionID /format:csv命令输出信息
func parseNICGetOutput(log logger.Logger, outUTF8 []byte) (list []NIC, err error) {
	rd := bufio.NewReader(bytes.NewBuffer(outUTF8))
	var started bool //标识正文开始
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			if log != nil {
				log.Error(err)
			}
			return nil, err
		}
		line = strings.TrimSpace(line)
		if strings.Contains(line, "没有可用实例") {
			// return nil, ErrNICNotFound
			return nil, nil
		}

		if line == "Node,MACAddress,NetConnectionID" {
			started = true
			continue
		}

		if !started {
			continue
		}

		properties := strings.Split(line, ",")
		if len(properties) != 3 {
			continue
		}

		list = append(list, NIC{
			MacAddr:   strings.TrimSpace(properties[1]),
			NetConnID: strings.TrimSpace(properties[2]),
		})
	}
	return list, nil
}

func nicGet(log logger.Logger, where string, properties ...string) (output []byte, err error) {
	// 示例：wmic nic where "PhysicalAdapter=TRUE" get MACAddress,Name,NetConnectionID /format:csv
	cmdAndArgs := fmt.Sprintf(`wmic nic where "%s" get %s /format:csv`, where, strings.Join(properties, ","))
	return winutil.ExecOutputWithLog(log, cmdAndArgs)
}
