package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/urfave/cli"

	"idcos.io/cloudboot/build"
	nw "idcos.io/cloudboot/utils/network"
)

var (
	name     = "gen-dhcp-conf"
	usage    = "for server pxe,bmc,ilo. you should export DHCP_TOKEN as os env before running."
	toFile   = ""
	port     = "8083"
	dhcpFile = "/etc/dhcp/dhcpd.conf"
	category = "ilo,intranet,tgw_intranet"
	server_room_name = ""
	cidr 	 = ""
	switches = ""
	ip       = "localhost"
)

func main() {
	app := cli.NewApp()
	app.Name = name
	app.Usage = usage
	app.Version = build.Version()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "t, to",
			Value: toFile,
			Usage: "生成的文件放到哪里 (default: 输出到控制台)",
		},
		cli.StringFlag{
			Name:  "c, category",
			Value: category,
			Usage: "指定具体网段类别",
		},
		cli.StringFlag{
			Name:  "srn, server-room-name",
			Value: server_room_name,
			Usage: "指定机房管理单元涉及的网段",
		},
		cli.StringFlag{
			Name:  "ci, cidr",
			Value: server_room_name,
			Usage: "指定网段名称（如：192.168.0.1/26）",
		},
		cli.StringFlag{
			Name:  "sw, switches",
			Value: server_room_name,
			Usage: "指定交换机固资编号",
		},				
		cli.StringFlag{
			Name:  "p, port",
			Value: port,
			Usage: "cloudboot-server 端口号",
		},
		cli.StringFlag{
			Name:  "i, ip",
			Value: ip,
			Usage: "cloudboot-server IP地址",
		},
		cli.StringFlag{
			Name:  "d, dhcp",
			Value: dhcpFile,
			Usage: "dhcp配置文件路径",
		},
	}
	app.Action = func(ctx *cli.Context) error {
		toFile = ctx.String("t")
		category = ctx.String("c")
		server_room_name = ctx.String("srn")
		cidr = ctx.String("ci")
		switches = ctx.String("sw")
		port = ctx.String("p")
		ip = ctx.String("i")
		dhcpFile = ctx.String("d")

		if !checkAvailable(category, dhcpFile) {
			return errors.New("")
		}

		conf := GenConf(ip, port, toFile, category, server_room_name, cidr, switches)
		return Main(conf)
	}

	app.Run(os.Args)
}

var (
	// DhcpConf dhcp ip配置
	DhcpConfILO = `subnet %s netmask %s {
    option routers %s;
    default-lease-time 7200;
    max-lease-time 7200;
    pool {
        failover peer "failover-dhcp";
        range %s %s;
    }
}`
	DhcpConfPXE = `subnet %s netmask %s {
    option routers %s;
    default-lease-time 600;
    max-lease-time 600;
    pool {
        failover peer "failover-dhcp";
        range %s %s;
    }
}`
)

//PxeRange PxeRange
type PxeRange struct {
	Start string
	End   string
}

// ipToUint ip4
func ipToUint(ip string) (int, error) {
	list := strings.Split(ip, ".")
	min1, err := strconv.Atoi(list[0])
	if err != nil {
		return 1, err
	}
	min2, err := strconv.Atoi(list[1])
	if err != nil {
		return 1, err
	}
	min3, err := strconv.Atoi(list[2])
	if err != nil {
		return 1, err
	}
	min4, err := strconv.Atoi(list[3])
	if err != nil {
		return 1, err
	}

	return min1*256*256*256 + min2*256*256 + min3*256 + min4, nil
}

func uintToIP(ipu int) string {
	ipu1 := ipu / (256 * 256 * 256)
	ipum := ipu % (256 * 256 * 256)
	ipu2 := ipum / (256 * 256)
	ipum = ipum % (256 * 256)
	ipu3 := ipum / 256
	ipu4 := ipum % 256

	return fmt.Sprintf("%d.%d.%d.%d", ipu1, ipu2, ipu3, ipu4)
}

func checkIntersection(pxeF, pxeS string) (string, bool) {
	var start, end int
	intersection := false
	pxeFS := strings.Split(pxeF, ",")
	if len(pxeFS) < 2 {
		return "", intersection
	}
	pxeSS := strings.Split(pxeS, ",")
	if len(pxeSS) < 2 {
		return "", intersection
	}
	pxeFSS, _ := ipToUint(pxeFS[0])
	pxeFSE, _ := ipToUint(pxeFS[1])
	pxeSSS, _ := ipToUint(pxeSS[0])
	pxeSSE, _ := ipToUint(pxeSS[1])

	start = pxeSSS
	if pxeFSS <= pxeSSS {
		intersection = true
		start = pxeFSS
	}

	end = pxeSSE
	if pxeSSE <= pxeFSE {
		intersection = true
		end = pxeFSE
	}
	return uintToIP(start) + "," + uintToIP(end), intersection
}

func checkRepeat(items []IPNetworkPage) map[string]string {
	// 筛查pxe池和CIDR完全一样的
	cidrPxeRange := make(map[string]bool, 1)
	for k := range items {
		if !cidrPxeRange[items[k].CIDR+":"+items[k].PXEPool] {
			cidrPxeRange[items[k].CIDR+":"+items[k].PXEPool] = true
		}
	}

	// 对于有交集的pxe池合并
	var netmask, min, key string
	cidrMapPxeRange := make(map[string]string, 1)
	for k := range cidrPxeRange {
		cprSplit := strings.Split(k, ":")
		if len(cprSplit) < 2 {
			continue
		}
		netmask, _, min, _ = nw.GetCidrIPRouteAndSubNet(cprSplit[0])
		key = netmask + ":" + min
		if cidrMapPxeRange[key] != "" {
			result, intersection := checkIntersection(cidrMapPxeRange[key], cprSplit[1])
			if intersection {
				cidrMapPxeRange[key] = result
			} else {
				cidrMapPxeRange[key] = cprSplit[1]
			}
		} else {
			cidrMapPxeRange[key] = cprSplit[1]
		}
	}
	return cidrMapPxeRange
}


// Main 入口函数
func Main(conf *Config) error {
	var dhcpConf []string

	// 获取网段分页数据
	items, err := GetIPNetworks(conf.DHCP.ServerRoomName, conf.DHCP.CIDR, conf.DHCP.Switches, conf.DHCP.Category, conf.Repo.Port,conf.Repo.IP)
	if err != nil || items == nil {
		log.Println(err)
		return err
	}

	for k, _ := range items {
		//仅处理ipv4
		if items[k].Version != "ipv4" {
			continue
		}
		//仅处理 服务器普通内网、服务器ILO
		switch items[k].Category {
		case "intranet":
			if items[k].PXEPool != "" {
				pxe := strings.Split(items[k].PXEPool, ",")
				if len(pxe) < 2 {
					continue
				}
				netmask, _, min, _ := nw.GetCidrIPRouteAndSubNet(items[k].CIDR)
				dhcpConf = append(dhcpConf, fmt.Sprintf(DhcpConfPXE, min, netmask, items[k].Gateway, pxe[0], pxe[1]))
			}

		case "ilo":
			if items[k].IPPool != "" {
				ilo := strings.Split(items[k].IPPool, ",")
				if len(ilo) < 2 {
					continue
				}
				netmask, _, min, _ := nw.GetCidrIPRouteAndSubNet(items[k].CIDR)
				dhcpConf = append(dhcpConf, fmt.Sprintf(DhcpConfILO, min, netmask, items[k].Gateway, ilo[0], ilo[1]))
			}
		case "tgw_intranet":
			if items[k].PXEPool != "" {
				pxe := strings.Split(items[k].PXEPool, ",")
				if len(pxe) < 2 {
					continue
				}
				netmask, _, min, _ := nw.GetCidrIPRouteAndSubNet(items[k].CIDR)
				dhcpConf = append(dhcpConf, fmt.Sprintf(DhcpConfPXE, min, netmask, items[k].Gateway, pxe[0], pxe[1]))
			}
		}
	}
	//cpr := checkRepeat(items)
	// 生成dhcp子网配置
	//for k, v := range cpr {
	//	cprSplit := strings.Split(k, ":")
	//	if len(cprSplit) < 1 {
	//		continue
	//	}
	//	pxe := strings.Split(v, ",")
	//	if len(pxe) < 2 {
	//		continue
	//	}
	//	if  == "ilo" {
	//		dhcpConf = append(dhcpConf, fmt.Sprintf(DhcpConfILO, cprSplit[1], cprSplit[0], pxe[0], pxe[0], pxe[1]))
	//	}
	//}
	// 如果数据库中没有pxe配置，直接跳过
	if len(dhcpConf) < 1 {
		return nil
	}
	subnet := strings.Join(dhcpConf, "\n")

	// 读取dhcp配置文件头部内容并拼接
	//if !Exist(conf.DHCP.Path) {
	//	fmt.Printf("文件%s不存在\n", conf.DHCP.Path)
	//	return nil
	//}
	//content, err := ReadFile(conf.DHCP.Path, conf.DHCP.Copy)
	//if err != nil {
	//	log.Println(err)
	//	return err
	//}
	//index := strings.Index(content, "subnet")
	//content = content[:index]
	//content = content + subnet

	// 写入dhcp配置文件或者打印到控制台
	if conf.DHCP.To == "" || !Exist(conf.DHCP.To) {
		fmt.Printf("%s\n", subnet)
		return nil
	}
	err = WriteFile(conf.DHCP.To, subnet)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//ReadFile 读取文件
func ReadFile(path, copy string) (string, error) {
	fi, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		return "", err
	}
	content := string(fd)
	err = GenerTempFile(copy, []byte(content))
	if err != nil {
		return "", err
	}
	return content, nil
}

//WriteFile 写文件
func WriteFile(path, content string) error {
	fi, err := os.Open(name)
	if err != nil {
		return err
	}
	defer fi.Close()
	err = ioutil.WriteFile(path, []byte(content), 766)
	return err
}

//GenerTempFile 生成临时文件
func GenerTempFile(copy string, content []byte) error {
	exten := fmt.Sprintf("dhcp-%s-back", time.Now().Format("2006-01-02-15:04:05"))
	tmpfile, err := ioutil.TempFile(copy, exten)
	if err != nil {
		return err
	}

	if _, err := tmpfile.Write(content); err != nil {
		return err
	}
	log.Printf("备份配置文件路径:%s\n", tmpfile.Name())
	err = tmpfile.Close()
	return err
}

//Exist 判断文件是否存在
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func checkCategory(category []string) bool {
	if category == nil {
		return true
	}
	var result string
	for k := range category {
		switch category[k] {
		case "ilo":
		case "tgw_intranet":
		case "tgw_extranet":
		case "intranet":
		case "extranet":
		case "v_intranet":
		case "v_extranet":
		default:
			result = result + category[k] + ","
		}
	}
	result = strings.TrimRight(result, ",")
	if result != "" {
		fmt.Printf("category(%s)必须为('ilo','tgw_intranet','tgw_extranet','intranet','extranet','v_intranet','v_extranet')\n", result)
		return false
	}
	return true
}

func checkAvailable(category, dhcpfile string) bool {
	existError := true
	if !Exist(dhcpfile) {
		fmt.Printf("DHCP配置文件(%s)不存在\n", dhcpfile)
		existError = false
	}

	if !checkCategory(strings.Split(category, ",")) {
		existError = false
	}
	return existError
}

// Config config 数据结构体
type Config struct {
	Repo Repo
	DHCP DHCP
}

//Repo 数据库配置
type Repo struct {
	Token string
	Port  string
	IP    string
}

//DHCP DHCP路径配置
type DHCP struct {
	Path     string `json:"path"`
	To       string `json:"to"`
	Copy     string `json:"copy"`
	Category string `json:"category"`
	CIDR string `json:"cidr"`
	Switches string `json:"switches"`	
	ServerRoomName string `json:"server_room_name"`
}

// JSONLoader 从json 文件中加配置数据
type JSONLoader struct {
	path string
}

// New 新建JSONLoader
func New(path string) *JSONLoader {
	return &JSONLoader{path}
}

// GenConf 构造配置
func GenConf(ip, port, to, category, server_room_name, cidr, switches string) *Config {

	conf := Config{
		Repo: Repo{
			Token: token,
			Port:  port,
			IP:    ip,
		},
		DHCP: DHCP{
			To:       to,
			Category: category,
			ServerRoomName: server_room_name,
			CIDR: cidr,
			Switches: switches,
			Path:     dhcpFile,
		},
	}

	return &conf
}
