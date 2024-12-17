package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/astaxie/beego/httplib"
	"github.com/axgle/mahonia"
	"github.com/urfave/cli"

	"idcos.io/cloudboot/build"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/server/cloudbootserver/service"
	"idcos.io/cloudboot/utils"
)

func main() {
	app := cli.NewApp()
	app.Version = build.Version()
	app.Action = func(c *cli.Context) {
		run(c)
	}
	app.Run(os.Args)
}

var (
	rootPath          = "c:/firstboot"
	serverHost        = "osinstall" //cloudboot server host
	scriptFile        = filepath.Join(rootPath, "temp-script.cmd")
	preInstallScript  = filepath.Join(rootPath, "preInstall.cmd")
	postInstallScript = filepath.Join(rootPath, "postInstall.cmd")

	devSettFile        = filepath.Join(rootPath, "deviceSetting.json")
	networkSettingFile = filepath.Join(rootPath, "networkSetting.json")
	nicFile            = filepath.Join(rootPath, "nic.txt")
)

func run(c *cli.Context) error {
	utils.InitFileLog()
	enc := mahonia.NewEncoder("gbk")
	if utils.CheckFileIsExist(preInstallScript) {
		content, err := getFileContent(preInstallScript)
		if err != nil {
			utils.Logger.Error(err.Error())
		}
		utils.Logger.Debug("run pre install script:" + preInstallScript + "\n" + content)
		b, err := utils.ExecScript(preInstallScript)
		if err != nil {
			utils.Logger.Error("preinstall error: %s", err.Error())
		}
		var output = string(b)
		output = enc.ConvertString(output)
		utils.Logger.Debug(output)
	}
	//init cloudboot server host
	serverIP := getDomainLookupIP(serverHost)
	if serverIP != "" {
		serverHost = serverIP
	}

	var sn = getSN()
	isVM := isVirtualMachine()
	if isVM {
		sn = getMacAddress()
	}

	if sn == "" {
		utils.Logger.Error("get sn failed!")
	}

	//deviceSettings, err := loadDeviceSetting()
	//if err != nil {
	//	utils.Logger.Error(err.Error())
	//}

	// TODO 修改用户及密码

	// 执行post脚本，拿到改网络之前，为的是在post里执行检查脚本的结果，可以上报到server from wangsu@20161201
	if utils.CheckFileIsExist(postInstallScript) {
		content, err := getFileContent(postInstallScript)
		if err != nil {
			utils.Logger.Error(err.Error())
		}
		utils.Logger.Debug("run post install script:" + postInstallScript + "\n" + content)
		b, err := utils.ExecScript(postInstallScript)
		if err != nil {
			utils.Logger.Error("postInstall error: %s", err.Error())
		}
		var output = string(b)
		//output = enc.ConvertString(output)
		utils.Logger.Debug(output)
		utils.Logger.Debug("execute post script finish")
	}

	//修改网络
	//_ = deviceSettings
	//if err = changeNetwork(); err != nil {
	//	utils.Logger.Error(err.Error())
	//}

	//TODO 修改路由
	//if err = changeReg(); err != nil {
	//	utils.Logger.Error(err.Error())
	//}

	// 这两行代码要和"idcos.io/cloudboot/utils" 下的utils.go中的代码保持一致，否则会找不到日志文件
	var logPath = path.Join(rootPath, "log")
	var logFile = path.Join(logPath, "setup.log")
	// 重启之前上传日志
	PostLog(serverHost, sn, logFile)

	utils.Logger.Debug("start reboot")
	reboot()

	return nil
}
func changeNetwork() error {
	ns, err := loadNetworkSetting()
	if err != nil {
		utils.Logger.Errorf("load networkSetting.json fail,%v", err)
		return err
	} else if ns != nil {
		//如果不做bonding，eth0-外网IP，eth1内网IP
		//如果有多IP，只取列表的第一个下发配置
		intranetIndex, extranetIndex := -1, -1 //标记需要配置的内外网IP在数组中的位置
		for k, item := range ns.Items {
			if item.Scope != nil && *item.Scope == model.Intranet && intranetIndex != -1 {
				intranetIndex = k
			} else if item.Scope != nil && *item.Scope == model.Extranet && extranetIndex != -1 {
				extranetIndex = k
			}
		}
		if ns.BondingRequired == model.NO {
			if intranetIndex != -1 {
				err = changeIP("本地连接2",
					ns.Items[intranetIndex].IP,
					ns.Items[intranetIndex].Netmask,
					ns.Items[intranetIndex].Gateway)
				if err != nil {
					utils.Logger.Errorf("intranet ip:%s config fail,%v", ns.Items[intranetIndex], err)
					//return
				}
			}
			if extranetIndex != -1 {
				err = changeIP("本地连接1",
					ns.Items[extranetIndex].IP,
					ns.Items[extranetIndex].Netmask,
					ns.Items[extranetIndex].Gateway)
				if err != nil {
					utils.Logger.Errorf("extranet ip:%s config fail,%v", ns.Items[extranetIndex], err)
					//return
				}
			}
		} else {
			//如果需要bonding... TODO
			if err := teaming(ns); err != nil {
				utils.Logger.Errorf("teaming fail:%v", err)
				return err
			}
		}
	}

	return nil
}

type NicAdapter struct {
	ID   string
	Name string
}

func teaming(netSett *service.GetNetworkSettingBySNResp) error {
	//采集本地网卡信息
	raw := getRawLocalhostNicInfo()
	rows := strings.Split(raw, "\r\n")

	type LNic struct {
		Hostname     string
		Mac          string
		HardwareName string
		Name         string
	}
	var lnics []LNic
	for k, v := range rows {
		if k == 0 {
			continue
		}
		arr := strings.Split(v, ",")
		utils.Logger.Info(v)
		utils.Logger.Info(fmt.Sprintf("index:%d,array length:%d", k, len(arr)))
		if len(arr) != 4 {
			utils.Logger.Error("nic info wrong:")
			continue
		}
		var lnic LNic
		lnic.Hostname = strings.TrimSpace(arr[0])
		lnic.Mac = strings.TrimSpace(arr[1])
		lnic.Mac = strings.ToLower(lnic.Mac)
		lnic.HardwareName = strings.TrimSpace(arr[2])
		lnic.Name = strings.TrimSpace(arr[3])
		lnics = append(lnics, lnic)
	}

	//等待10S，等待驱动加载完成
	utils.Logger.Debug("sleep 10s:")
	time.Sleep(10 * time.Second)

	//采集网卡适配器信息

	var adapters []NicAdapter
	raw = getRawNicAdapterInfo()
	var r = `(\d)\)(\s)(.*)`
	reg := regexp.MustCompile(r)
	matchs := reg.FindAllStringSubmatch(raw, -1)
	for _, match := range matchs {
		var adapter NicAdapter
		adapter.ID = strings.TrimSpace(match[1])
		adapter.ID = strings.Trim(adapter.ID, "\r\n")
		adapter.Name = strings.TrimSpace(match[3])
		adapter.Name = strings.Trim(adapter.Name, "\r\n")
		adapters = append(adapters, adapter)
	}

	if err := teamingForEach(netSett, adapters); err != nil {
		return err
	}

	return nil
}

//单组teaming配置
func teamingForEach(netSett *service.GetNetworkSettingBySNResp, adapters []NicAdapter) error {
	for index := 0; index < 2; index++ {
		//创建 Teaming // 做team, eth2,3做bond0配外网，eth0,1做bond1,配内网
		cmd := fmt.Sprintf(`PROSetCL.exe Team_Create %s,%s Team%d SFT`, adapters[2-2*index].ID, adapters[3-2*index].ID, index)
		errCmd := runCmd(cmd)
		if errCmd != nil {
			utils.Logger.Errorf("team create error :%s", errCmd.Error())
		}

		//设置teaming参数
		cmd = fmt.Sprintf(`PROSetCL.exe Team_SetSetting %d FailbackEnabled 禁用`, index+1)
		errCmd = runCmd(cmd)
		if errCmd != nil {
			utils.Logger.Errorf("team create error :%s", errCmd.Error())
		}
		cmd = fmt.Sprintf(`PROSetCL.exe Team_SetSetting %d ConnMonEnabled 启用`, index+1)
		errCmd = runCmd(cmd)
		if errCmd != nil {
			utils.Logger.Errorf("team create error :%s", errCmd.Error())
		}

		//设置网关
		cmd = fmt.Sprintf(`PROSetCL.exe Team_SetSetting %d ConnMonClients %s`, index+1, netSett.Items[index].Gateway)
		errCmd = runCmd(cmd)
		if errCmd != nil {
			utils.Logger.Errorf("team create error :%s", errCmd.Error())
		}

		cmd = fmt.Sprintf(`PROSetCL.exe Team_SetAdapterPriority %s 主适配器`, adapters[2-2*index].ID)
		errCmd = runCmd(cmd)
		if errCmd != nil {
			utils.Logger.Errorf("team create error :%s", errCmd.Error())
		}

		cmd = fmt.Sprintf(`PROSetCL.exe Team_SetAdapterPriority %s 次适配器`, adapters[3-2*index].ID)
		errCmd = runCmd(cmd)
		if errCmd != nil {
			utils.Logger.Errorf("team create error :%s", errCmd.Error())
		}
	}
	return nil
}

func getFileContent(path string) (string, error) {
	fi, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		return "", err
	}
	return string(fd), nil
}

// TODO
func configNic() (string, error) {
	//做teaming
	// 1.查询Server接口 getnicinfobysn，获取逻辑名  2.查询本地网卡信息  3.重置网卡名
	//通过API获取网卡逻辑命名信息
	// var url = fmt.Sprintf("http://%s/api/osinstall/v1/device/getNicInfoBySn?sn=%s", host, sn)
	// utils.Logger.Debug(url)
	// resp, err := http.Get(url)
	// if err != nil {
	// 	return "", err
	// }
	// defer resp.Body.Close()
	// if resp.StatusCode != http.StatusOK {
	// 	return "", fmt.Errorf("http status code: %d", resp.StatusCode)
	// }
	// var body []byte
	// body, err = ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return "", fmt.Errorf("call url: %s failed", url)
	// }

	cache, err := getFileContent(nicFile)
	if err != nil {
		return "", err
	}
	body := []byte(cache)

	type Nic struct {
		Name string
		Mac  string
		Type string
	}
	var jsonResp struct {
		Content []Nic
	}
	if err := json.Unmarshal(body, &jsonResp); err != nil {
		return "", err
	}
	nics := jsonResp.Content

	//采集本地网卡信息
	raw := getRawLocalhostNicInfo()
	rows := strings.Split(raw, "\r\n")

	type LNic struct {
		Hostname     string
		Mac          string
		HardwareName string
		Name         string
	}
	var lnics []LNic
	for k, v := range rows {
		if k == 0 {
			continue
		}
		arr := strings.Split(v, ",")
		utils.Logger.Error(v)
		utils.Logger.Error(fmt.Sprintf("index:%d,array length:%d", k, len(arr)))
		if len(arr) != 4 {
			utils.Logger.Error("nic info wrong:")
			continue
		}
		var lnic LNic
		lnic.Hostname = strings.TrimSpace(arr[0])
		lnic.Mac = strings.TrimSpace(arr[1])
		lnic.Mac = strings.ToLower(lnic.Mac)
		lnic.HardwareName = strings.TrimSpace(arr[2])
		lnic.Name = strings.TrimSpace(arr[3])
		lnics = append(lnics, lnic)

		for _, nic := range nics {
			if lnic.Mac == nic.Mac {
				err := changeNicName(lnic.Name, nic.Name)
				if err != nil {
					utils.Logger.Error("change nic name error:", err.Error())
				}
			}
		}
	}

	//等待10S，等待驱动加载完成
	utils.Logger.Debug("sleep 10s:")
	time.Sleep(10 * time.Second)

	//采集网卡适配器信息
	type NicAdapter struct {
		ID   string
		Name string
	}
	var adapters []NicAdapter
	raw = getRawNicAdapterInfo()
	var r = `(\d)\)(\s)(.*)`
	reg := regexp.MustCompile(r)
	matchs := reg.FindAllStringSubmatch(raw, -1)
	for _, match := range matchs {
		var adapter NicAdapter
		adapter.ID = strings.TrimSpace(match[1])
		adapter.ID = strings.Trim(adapter.ID, "\r\n")
		adapter.Name = strings.TrimSpace(match[3])
		adapter.Name = strings.Trim(adapter.Name, "\r\n")
		adapters = append(adapters, adapter)
	}

	//创建 Teaming
	var masterAdapterID string
	var slaveAdapterID string
	var masterMac string
	for _, nic := range nics {
		if nic.Type == "master" {
			for _, lnic := range lnics {
				if nic.Mac == lnic.Mac {
					for _, adapter := range adapters {
						if adapter.Name == lnic.HardwareName {
							masterAdapterID = adapter.ID
						}
					}
				}
			}
			masterMac = nic.Mac
		}
		if nic.Type == "slave" {
			for _, lnic := range lnics {
				if nic.Mac == lnic.Mac {
					for _, adapter := range adapters {
						if adapter.Name == lnic.HardwareName {
							slaveAdapterID = adapter.ID
						}
					}
				}
			}
		}
	}
	//创建 Teaming
	cmd := fmt.Sprintf(`PROSetCL.exe Team_Create %s,%s Team0 SFT`, masterAdapterID, slaveAdapterID)
	errCmd := runCmd(cmd)
	if errCmd != nil {
		utils.Logger.Error("team create error :", errCmd.Error())
	}

	//设置teaming参数
	cmd = `PROSetCL.exe Team_SetSetting 1 FailbackEnabled 禁用`
	errCmd = runCmd(cmd)
	if errCmd != nil {
		utils.Logger.Error("team create error :", errCmd.Error())
	}
	cmd = `PROSetCL.exe Team_SetSetting 1 ConnMonEnabled 已启用`
	errCmd = runCmd(cmd)
	if errCmd != nil {
		utils.Logger.Error("team create error :", errCmd.Error())
	}

	//获取网关
	// raw, err = getRawDeviceNetworkInfo(host, sn)
	// if err != nil {
	// 	utils.Logger.Error("get gateway from api error :", err.Error())
	// }
	// type JsonRespGatewayContent struct {
	// 	Gateway string
	// }
	// var jsonRespGateway struct {
	// 	Value JsonRespGatewayContent
	// }
	// err = json.Unmarshal([]byte(raw), &jsonRespGateway)
	// if err != nil {
	// 	utils.Logger.Error("format json error :", err.Error())
	// }
	restInfo, err := loadDeviceSetting()
	if err != nil {
		utils.Logger.Error(err.Error())
	}
	//设置网关
	cmd = fmt.Sprintf(`PROSetCL.exe Team_SetSetting 1 ConnMonClients %s`, restInfo.IntranetIP.Gateway)
	errCmd = runCmd(cmd)
	if errCmd != nil {
		utils.Logger.Error("team create error :", errCmd.Error())
	}

	cmd = fmt.Sprintf(`PROSetCL.exe Team_SetAdapterPriority %s 主适配器`, masterAdapterID)
	errCmd = runCmd(cmd)
	if errCmd != nil {
		utils.Logger.Error("team create error :", errCmd.Error())
	}

	cmd = fmt.Sprintf(`PROSetCL.exe Team_SetAdapterPriority %s 次适配器`, slaveAdapterID)
	errCmd = runCmd(cmd)
	if errCmd != nil {
		utils.Logger.Error("team create error :", errCmd.Error())
	}
	return masterMac, nil
}

func runCmd(cmd string) error {
	enc := mahonia.NewEncoder("gbk")
	cmd = enc.ConvertString(cmd)
	utils.Logger.Debug(cmd)
	outputBytes, err := utils.ExecCmdOutputWithLogfile(scriptFile, cmd)
	if err != nil {
		return err
	}
	//output := enc.ConvertString(string(outputBytes))
	output := string(outputBytes)
	utils.Logger.Debug(string(output))
	return nil
}

func getRawNicAdapterInfo() string {
	cmd := `PROSetCL.exe Adapter_Enumerate`
	var output string
	utils.Logger.Debug(cmd)
	outputBytes, err := utils.ExecCmdOutputWithLogfile(scriptFile, cmd)
	if err != nil {
		utils.Logger.Error(err.Error())
	} else {
		// enc := mahonia.NewDecoder("gbk")
		// output = enc.ConvertString(string(outputBytes))
		output = string(outputBytes)
		utils.Logger.Debug(output)
	}
	return output
}

func getRawLocalhostNicInfo() string {
	cmd := `wmic nic where "PhysicalAdapter=TRUE" get MACAddress,Name,NetConnectionID /format:csv`
	var output string
	utils.Logger.Debug(cmd)
	outputBytes, err := utils.ExecCmd(scriptFile, cmd)
	if err != nil {
		utils.Logger.Error(err.Error())
	} else {
		// enc := mahonia.NewDecoder("gbk")
		// output = enc.ConvertString(string(outputBytes))
		output = string(outputBytes)
		utils.Logger.Debug(output)
	}
	return output
}

// 不推荐使用
func changeNicName(oldName string, newName string) error {
	enc := mahonia.NewEncoder("gbk")
	newName = enc.ConvertString(newName)
	// oldName = enc.ConvertString(oldName)
	cmd := fmt.Sprintf(`netsh interface set interface name="%s" newname="%s"`, oldName, newName)
	//cmd = enc.ConvertString(cmd)
	utils.Logger.Debug(cmd)
	output, err := utils.ExecCmd(scriptFile, cmd)
	if err != nil {
		return err
	}
	utils.Logger.Debug(string(output))
	return nil
}

// 查看本机 SN
func getSN() string {
	var cmd = `wmic bios get SerialNumber /VALUE`
	var r = `SerialNumber=(\S+)`
	var output string
	utils.Logger.Debug(cmd)
	if outputBytes, err := utils.ExecCmd(scriptFile, cmd); err != nil {
		utils.Logger.Error(err.Error())
	} else {
		output = string(outputBytes)
		utils.Logger.Debug(output)
	}

	reg := regexp.MustCompile(r)
	var regResult = reg.FindStringSubmatch(output)
	if regResult == nil || len(regResult) != 2 {
		return ""
	}

	// fmt.Println(strings.Trim(regResult[1], "\r\n"))
	var result string
	result = strings.Trim(regResult[1], "\r\n")
	result = strings.TrimSpace(result)
	return result
}

//是否是虚拟机
func isVirtualMachine() bool {
	var cmd = `systeminfo`
	var output string
	utils.Logger.Debug(cmd)
	if outputBytes, err := utils.ExecCmd(scriptFile, cmd); err != nil {
		utils.Logger.Error(err.Error())
	} else {
		output = string(outputBytes)
		utils.Logger.Debug(output)
	}

	isValidate, err := regexp.MatchString(`(?i)VMware|VirtualBox|KVM|Xen|Parallels`, output)
	if err != nil {
		utils.Logger.Error(err.Error())
		return false
	}
	return isValidate
}

// 获取Mac地址
func getMacAddress() string {
	var cmd = `wmic nic where "NetConnectionStatus=2" get MACAddress /VALUE`
	var r = `(?i)MACAddress=(\S+)`
	var output string
	utils.Logger.Debug(cmd)
	if outputBytes, err := utils.ExecCmd(scriptFile, cmd); err != nil {
		utils.Logger.Error(err.Error())
	} else {
		output = string(outputBytes)
		utils.Logger.Debug(output)
	}

	reg := regexp.MustCompile(r)
	var regResult = reg.FindStringSubmatch(output)
	if regResult == nil || len(regResult) != 2 {
		return ""
	}

	var result string
	result = strings.Trim(regResult[1], "\r\n")
	result = strings.TrimSpace(result)
	return result
}

//get domain's lookup ip
func getDomainLookupIP(domain string) string {
	var cmd = `ping ` + domain
	var r = `(.+)(\s)(\d+)\.(\d+)\.(\d+)\.(\d+)([:|\s])(.+)TTL`
	var output string
	utils.Logger.Debug(cmd)
	if outputBytes, err := utils.ExecCmd(scriptFile, cmd); err != nil {
		utils.Logger.Error(err.Error())
	} else {
		output = string(outputBytes)
		utils.Logger.Debug(output)
	}

	reg := regexp.MustCompile(r)
	var regResult = reg.FindStringSubmatch(output)
	if regResult == nil || len(regResult) != 9 {
		return ""
	}

	var result = fmt.Sprintf("%s.%s.%s.%s", strings.TrimSpace(regResult[3]),
		strings.TrimSpace(regResult[4]),
		strings.TrimSpace(regResult[5]),
		strings.TrimSpace(regResult[6]))
	return result
}

// 网卡名称
func getNicInterfaceIndex(mac string) string {
	enc := mahonia.NewDecoder("gbk")
	// var cmd = fmt.Sprintf(`wmic nic where (MACAddress="%s" and Name="%s") get InterfaceIndex /value`, mac, name)
	//var cmd = fmt.Sprintf(`wmic nic where (MACAddress="%s" and Name like '%%%%%%%%Team0%%%%%%%%') get InterfaceIndex /value`, mac)
	var cmd = fmt.Sprintf(`wmic nic where (MACAddress="%s" AND netConnectionStatus=2) get InterfaceIndex /value`, mac)
	var r = `InterfaceIndex=(.*)`
	var output string
	utils.Logger.Debug(cmd)
	if outputBytes, err := utils.ExecCmd(scriptFile, cmd); err != nil {
		utils.Logger.Error(err.Error())
	} else {
		output = enc.ConvertString(string(outputBytes))
		//output = string(outputBytes)
		utils.Logger.Debug(output)
	}

	reg := regexp.MustCompile(r)
	var regResult = reg.FindStringSubmatch(output)
	if regResult == nil || len(regResult) != 2 {
		return ""
	}
	utils.Logger.Info("Nic Interface Index:" + regResult[1])
	// fmt.Println(strings.Trim(regResult[1], "\r\n"))
	return regResult[1]
}

// loadDeviceSetting 从文件导入全量硬件配置
func loadDeviceSetting() (*service.DeviceSetting, error) {
	cache, err := getFileContent(devSettFile)
	if err != nil {
		return nil, err
	}
	body := []byte(cache)
	var jsonResp struct {
		Status  string                `json:"status"`
		Message string                `json:"message"`
		Content service.DeviceSetting `json:"content"`
	}

	if err := json.Unmarshal(body, &jsonResp); err != nil {
		utils.Logger.Errorf("load deviceSetting.json failed, %s", err.Error())
		return nil, err
	}
	return &jsonResp.Content, nil
}

// loadNetworkSetting 从文件导入全量硬件配置
func loadNetworkSetting() (*service.GetNetworkSettingBySNResp, error) {
	cache, err := getFileContent(devSettFile)
	if err != nil {
		return nil, err
	}
	body := []byte(cache)
	var jsonResp struct {
		Status  string                            `json:"status"`
		Message string                            `json:"message"`
		Content service.GetNetworkSettingBySNResp `json:"content"`
	}

	if err := json.Unmarshal(body, &jsonResp); err != nil {
		utils.Logger.Errorf("load deviceSetting.json failed, %s", err.Error())
		return nil, err
	}
	return &jsonResp.Content, nil
}

// 修改 IP
func changeIP(nic, ip, netmask, gateway string) error {
	var cmd = fmt.Sprintf(`netsh interface ipv4 set address name="%s" source=static addr=%s mask=%s gateway=%s`, nic, ip, netmask, gateway)
	enc := mahonia.NewEncoder("gbk")
	cmd = enc.ConvertString(cmd)
	utils.Logger.Debug(cmd)

	output, err := utils.ExecCmd(scriptFile, cmd)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil
	}
	utils.Logger.Debug(string(output))
	return nil
}

// 修改DNS
func changeDNS(nic, dns string) error {
	var cmd = fmt.Sprintf(`netsh interface ipv4 set dnsservers name="%s" static %s primary`, nic, dns)
	enc := mahonia.NewEncoder("gbk")
	cmd = enc.ConvertString(cmd)
	utils.Logger.Debug(cmd)

	output, err := utils.ExecCmd(scriptFile, cmd)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil
	}
	utils.Logger.Debug(string(output))

	return nil
}

// 修改注册表
func changeReg() error {
	var cmd1 = `reg add "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" /v AutoAdminLogon /t reg_sz /d 0 /f`
	utils.Logger.Debug(cmd1)
	output, err := utils.ExecCmd(scriptFile, cmd1)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil
	}
	utils.Logger.Debug(string(output))

	var cmd2 = `reg add "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" /v Defaultpassword /t reg_sz /d "" /f`
	utils.Logger.Debug(cmd2)
	output, err = utils.ExecCmd(scriptFile, cmd2)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil
	}
	utils.Logger.Debug(string(output))
	return nil
}

// 重启
func reboot() error {
	var cmd = fmt.Sprintf(`shutdown -f -r -t 10`)
	utils.Logger.Debug(cmd)
	output, err := utils.ExecCmd(scriptFile, cmd)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil
	}
	utils.Logger.Debug(string(output))
	return nil
}

// loadData 加载日志
func loadData(logpath string) ([]byte, error) {
	data, err := ioutil.ReadFile(logpath)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	return data, nil
}

// PostLog 向服务端发送日志
func PostLog(serverHost, sn, logpath string) (err error) {
	data, err := loadData(logpath)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("http://%s/api/cloudboot/v1/devices/%s/components/%s/logs?lang=en-US",
		serverHost, sn,
		"winconfig")

	resp, err := httplib.Post(url).Header("Accept", "application/json").Body(data).DoRequest()
	if err != nil {
		utils.Logger.Error(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("post log failed: %s", http.StatusText(resp.StatusCode))
	}
	return nil
}
