package network

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

var (
	//errScope IP列表传递参数有问题
	errScope = errors.New("IP池区间不正确，必须以逗号分隔")
	errCIDRNotContains = errors.New("IPv6不属于该CIDR")
)

//GetIPListByMinAndMaxIP 获取IP列表
func GetIPListByMinAndMaxIP(scope string) ([]string, error) {
	minAndMax := strings.Split(scope, ",")
	if len(minAndMax) < 2 {
		return nil, errScope
	}
	var result []string
	list1 := strings.Split(minAndMax[0], ".")
	list2 := strings.Split(minAndMax[1], ".")
	min3, err := strconv.Atoi(list1[2])
	if err != nil {
		return result, err
	}

	min4, err := strconv.Atoi(list1[3])
	if err != nil {
		return result, err
	}
	if min4 == 0 {
		min4 = 1
	}

	max3, err := strconv.Atoi(list2[2])
	if err != nil {
		return result, err
	}

	max4, err := strconv.Atoi(list2[3])
	if err != nil {
		return result, err
	}

	// 如果第3级一样, 第四级决定IP列表
	if min3 == max3 {
		for j := min4; j <= max4; j++ {
			ip := list1[0] + "." + list1[1] + "." + fmt.Sprintf("%d", min3) + "." + fmt.Sprintf("%d", j)
			result = append(result, ip)
		}
		return result, nil
	}

	// 如果第3级区间至少为1
	for i := min3; i <= max3; i++ {
		// 最后一个区间
		if max3-i == 0 {
			for j := 1; j <= max4; j++ {
				ip := list1[0] + "." + list1[1] + "." + fmt.Sprintf("%d", i) + "." + fmt.Sprintf("%d", j)
				result = append(result, ip)
			}
		} else if min3-i == 0 {
			for j := min4; j <= 255; j++ {
				ip := list1[0] + "." + list1[1] + "." + fmt.Sprintf("%d", i) + "." + fmt.Sprintf("%d", j)
				result = append(result, ip)
			}
		} else {
			for j := 1; j <= 255; j++ {
				ip := list1[0] + "." + list1[1] + "." + fmt.Sprintf("%d", i) + "." + fmt.Sprintf("%d", j)
				result = append(result, ip)
			}
		}
	}
	return result, nil
}

// GetCidrIPFirstAndLast 获取CIDR的最大最小IP
func GetCidrIPFirstAndLast(cidr string) (string, string) {
	ip := strings.Split(cidr, "/")[0]
	ipSegs := strings.Split(ip, ".")
	maskLen, _ := strconv.Atoi(strings.Split(cidr, "/")[1])
	seg3MinIP, seg3MaxIP := GetIPSeg3Range(ipSegs, maskLen)
	seg4MinIP, seg4MaxIP := GetIPSeg4Range(ipSegs, maskLen)
	ipPrefix := ipSegs[0] + "." + ipSegs[1] + "."
	seg4MinIP = seg4MinIP - 1
	seg4MaxIP = seg4MaxIP + 1
	return ipPrefix + strconv.Itoa(seg3MinIP) + "." + strconv.Itoa(seg4MinIP),
		ipPrefix + strconv.Itoa(seg3MaxIP) + "." + strconv.Itoa(seg4MaxIP)
}

// GetCidrIPRange 获取CIDR的最大最小IP
func GetCidrIPRange(cidr string) (string, string) {
	ip := strings.Split(cidr, "/")[0]
	ipSegs := strings.Split(ip, ".")
	maskLen, _ := strconv.Atoi(strings.Split(cidr, "/")[1])
	seg3MinIP, seg3MaxIP := GetIPSeg3Range(ipSegs, maskLen)
	seg4MinIP, seg4MaxIP := GetIPSeg4Range(ipSegs, maskLen)
	ipPrefix := ipSegs[0] + "." + ipSegs[1] + "."

	return ipPrefix + strconv.Itoa(seg3MinIP) + "." + strconv.Itoa(seg4MinIP),
		ipPrefix + strconv.Itoa(seg3MaxIP) + "." + strconv.Itoa(seg4MaxIP)
}

// GetCidrFirst2IPList 获取CIDR的前2个IP List
func GetCidrFirst2IPList(cidr string) (result []string) {
	ip := strings.Split(cidr, "/")[0]
	ipSegs := strings.Split(ip, ".")
	maskLen, _ := strconv.Atoi(strings.Split(cidr, "/")[1])
	seg3MinIP, _ := GetIPSeg3Range(ipSegs, maskLen)
	seg4MinIP, _ := GetIPSeg4Range(ipSegs, maskLen)
	ipPrefix := ipSegs[0] + "." + ipSegs[1] + "."
	result = append(result, (ipPrefix + strconv.Itoa(seg3MinIP) + "." + strconv.Itoa(seg4MinIP)))
	result = append(result, (ipPrefix + strconv.Itoa(seg3MinIP) + "." + strconv.Itoa(seg4MinIP+1)))
	return
}

// GetCidrLast10IPRange 获取CIDR的最后10个IP range
func GetCidrLast10IPRange(cidr string) (string) {
	ip := strings.Split(cidr, "/")[0]
	ipSegs := strings.Split(ip, ".")
	maskLen, _ := strconv.Atoi(strings.Split(cidr, "/")[1])
	_, seg3MaxIP := GetIPSeg3Range(ipSegs, maskLen)
	seg4MinIP, seg4MaxIP := GetIPSeg4Range(ipSegs, maskLen)
	ipPrefix := ipSegs[0] + "." + ipSegs[1] + "."
	if (seg4MaxIP - seg4MinIP) < 10 {
		return ""
	}

	return ipPrefix + strconv.Itoa(seg3MaxIP) + "." + strconv.Itoa(seg4MaxIP-10) + "," +
		ipPrefix + strconv.Itoa(seg3MaxIP) + "." + strconv.Itoa(seg4MaxIP)
}

//GetCidrIPRouteAndSubNet 获取CIDR的路由IP
func GetCidrIPRouteAndSubNet(cidr string) (string, string, string, error) {
	ip, subnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", "", "", err
	}
	ipSegs := strings.Split(ip.String(), ".")
	maskLen, _ := strconv.Atoi(strings.Split(cidr, "/")[1])
	seg3MinIP, _ := GetIPSeg3Range(ipSegs, maskLen)
	seg4MinIP, _ := GetIPSeg4Range(ipSegs, maskLen)
	ipPrefix := ipSegs[0] + "." + ipSegs[1] + "."

	return net.IP(subnet.Mask).String(), ipPrefix + strconv.Itoa(seg3MinIP) + "." + strconv.Itoa(seg4MinIP), subnet.IP.String(), nil
}

//GetCidrHostNum 计算得到CIDR地址范围内可拥有的主机数量
func GetCidrHostNum(maskLen int) uint {
	cidrIPNum := uint(0)
	var i = uint(32 - maskLen - 1)
	for ; i >= 1; i-- {
		cidrIPNum += 1 << i
	}
	return cidrIPNum
}

//GetCidrIPMask 获取Cidr的掩码
func GetCidrIPMask(maskLen int) string {
	// ^uint32(0)二进制为32个比特1，通过向左位移，得到CIDR掩码的二进制
	cidrMask := ^uint32(0) << uint(32-maskLen)
	//fmt.Println(fmt.Sprintf("%b \n", cidrMask))
	//计算CIDR掩码的四个片段，将想要得到的片段移动到内存最低8位后，将其强转为8位整型，从而得到
	cidrMaskSeg1 := uint8(cidrMask >> 24)
	cidrMaskSeg2 := uint8(cidrMask >> 16)
	cidrMaskSeg3 := uint8(cidrMask >> 8)
	cidrMaskSeg4 := uint8(cidrMask & uint32(255))

	return fmt.Sprint(cidrMaskSeg1) + "." + fmt.Sprint(cidrMaskSeg2) + "." + fmt.Sprint(cidrMaskSeg3) + "." + fmt.Sprint(cidrMaskSeg4)
}

//GetIPSeg3Range 得到第三段IP的区间（第一片段.第二片段.第三片段.第四片段）
func GetIPSeg3Range(ipSegs []string, maskLen int) (int, int) {
	if maskLen > 24 {
		segIP, _ := strconv.Atoi(ipSegs[2])
		return segIP, segIP
	}
	ipSeg, _ := strconv.Atoi(ipSegs[2])
	return GetIPSegRange(uint8(ipSeg), uint8(24-maskLen))
}

//GetIPSeg4Range 得到第四段IP的区间（第一片段.第二片段.第三片段.第四片段）
func GetIPSeg4Range(ipSegs []string, maskLen int) (int, int) {
	ipSeg, _ := strconv.Atoi(ipSegs[3])
	segMinIP, segMaxIP := GetIPSegRange(uint8(ipSeg), uint8(32-maskLen))
	return segMinIP + 1, segMaxIP - 1
}

//GetIPSegRange 根据用户输入的基础IP地址和CIDR掩码计算一个IP片段的区间
func GetIPSegRange(userSegIP, offset uint8) (int, int) {
	var ipSegMax uint8 = 255
	netSegIP := ipSegMax << offset
	segMinIP := netSegIP & userSegIP
	segMaxIP := userSegIP&(255<<offset) | ^(255 << offset)
	return int(segMinIP), int(segMaxIP)
}

//CIDRContainsIP 查看IP是否包含在CIDR内
func CIDRContainsIP(cidr, source string) bool {
	_, subnet, _ := net.ParseCIDR(cidr)
	ip := net.ParseIP(source)
	return subnet.Contains(ip)
}

// 根据起始ipv6地址以及对应的CIDR，计算返回下一个地址
func GetNextIPv6OfCIDR(ipv6Begin, cidr string) (string, error) {
	_, ipnetwork, _ := net.ParseCIDR(cidr)
	ip := net.ParseIP(ipv6Begin)
	// ip[j] type uint8
	for j := len(ip)-1; j>=0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
	if ipnetwork.Contains(ip) {
		return ip.String(), nil
	} else {
		return "",errCIDRNotContains
	}
}

// 根据CIDR，计算返回第一个地址（使用于新网段场景，默认返回网段第4个地址）
func GetFirstIPv6OfCIDR(cidr string) (string, error) {
	ip, ipnetwork, _ := net.ParseCIDR(cidr)
	// 循环4次得到第一个可用地址
	for i := 0; i < 4; i++ {
		for j := len(ip)-1; j>=0; j-- {
			ip[j]++
			if ip[j] > 0 {
				break
			}
		}
	} 
	if ipnetwork.Contains(ip) {
		return ip.String(), nil
	} else {
		return "",errCIDRNotContains
	}
}