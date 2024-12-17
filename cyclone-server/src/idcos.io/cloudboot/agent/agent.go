package agent

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/httplib"

	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/server/cloudbootserver/service"
	"idcos.io/cloudboot/server/cloudbootserver/types/device"
	"idcos.io/cloudboot/utils/sh"
)

// Agent agent data struct
type Agent struct {
	Cmdline
	log          logger.Logger
	collected    *device.Device         // 采集到的设备信息
	settings     *service.DeviceSetting // 装机参数
	hwSrvBaseURL string
}

// Cmdline /proc/cmdline文件信息映射结构体
type Cmdline struct {
	ServerAddr     string
	HWServerPort   int
	DevelopeMode   string
	LoopInterval   int
	PreInstallURL  string
	PostInstallURL string
}

var (
	// 匹配proc/cmdline文件内信息的正则表达式
	srvAddrReg   = regexp.MustCompile(`SERVER_ADDR=([^ ]+)`)
	hwSrvPortReg = regexp.MustCompile(`HW_SERVER_PORT=([^ ]+)`)
	intervalReg  = regexp.MustCompile(`LOOP_INTERVAL=([^ ]+)`)
	devReg       = regexp.MustCompile(`DEVELOPER=([^ ]+)`)
	preReg       = regexp.MustCompile(`PRE=([^ ]+)`)
	postReg      = regexp.MustCompile(`POST=([^ ]+)`)
)

const (
	defaultLoopInterval = 60
	defaultHWSrvPort    = 8081
)

// LoadCmdline 读取/proc/cmdline文件并解析到结构体
func LoadCmdline(log logger.Logger) (*Cmdline, error) {
	output, err := sh.ExecOutputWithLog(log, "cat /proc/cmdline")
	if err != nil {
		return nil, err
	}
	var cLine Cmdline
	var pair []string

	// 加载SERVER_ADDR: required
	pair = srvAddrReg.FindStringSubmatch(string(output))
	if len(pair) != 2 { // k=v结构
		return nil, errors.New("'SERVER_ADDR' is missing")
	}
	cLine.ServerAddr = strings.TrimSpace(pair[1])

	// 加载HW_SERVER_PORT: optional
	if pair = hwSrvPortReg.FindStringSubmatch(string(output)); len(pair) == 2 {
		if port, err := strconv.Atoi(strings.TrimSpace(pair[1])); err == nil && port > 0 {
			cLine.HWServerPort = port
		}
	} else {
		cLine.HWServerPort = defaultHWSrvPort
	}

	// 加载LOOP_INTERVAL: optional
	if pair = intervalReg.FindStringSubmatch(string(output)); len(pair) == 2 {
		cLine.LoopInterval = parseInterval(strings.TrimSpace(pair[1]), defaultLoopInterval)
	} else {
		log.Warnf("%q is missing, use default value %d", "LOOP_INTERVAL", defaultLoopInterval)
		cLine.LoopInterval = defaultLoopInterval
	}

	// 加载DEVELOPER: optional
	if pair = devReg.FindStringSubmatch(string(output)); len(pair) == 2 {
		cLine.DevelopeMode = strings.TrimSpace(pair[1])
	}

	// 加载PRE: optional
	if pair = preReg.FindStringSubmatch(string(output)); len(pair) == 2 {
		cLine.PreInstallURL = strings.TrimSpace(pair[1])
	}

	// 加载POST: optional
	if pair = postReg.FindStringSubmatch(string(output)); len(pair) == 2 {
		cLine.PostInstallURL = strings.TrimSpace(pair[1])
	}
	return &cLine, nil
}

// New create agent
func New(log logger.Logger) (*Agent, error) {
	cLine, err := LoadCmdline(log)
	if err != nil {
		return nil, err
	}
	log.Infof("%#v", cLine)

	return &Agent{
		log:          log,
		hwSrvBaseURL: fmt.Sprintf("http://localhost:%d/", cLine.HWServerPort),
		Cmdline:      *cLine,
	}, nil
}

var (
	// ErrSNNotFound 设备序列号信息未采集到
	ErrSNNotFound = errors.New("SN not found")
)

// Run 运行Agent
func (agent *Agent) Run() (err error) {
	agent.runPreInstall()

	if err = agent.runHWServer(); err != nil {
		return err
	}

	agent.waitHWServerOK()

	if err := agent.collect(); err != nil {
		return err
	}

	if agent.collected == nil || agent.collected.SN == "" {
		agent.log.Error("Unable to collect SN")
		return ErrSNNotFound
	}

	agent.reportProgress(0.1, "进入BootOS", "已进入BootOS")
	_ = agent.postDevice()
	agent.reportProgress(0.2, "采集并上报设备信息", "完成")

	for {
		// 1、 一直等待，直到设备允许进入安装队列
		agent.waitToEnterQueue()
		agent.log.Debug("Enter installation queue")

		// 2、从服务端加载设备的装机参数
		if err = agent.loadDeviceSettings(); err != nil {
			agent.reportProgress(-1, "无法从cloudboot-server获取设备装机配置参数信息", err.Error())
			continue
		}
		agent.reportProgress(0.25, "拉取设备操作系统安装配置", "完成")

		// 3、调用hw-server组件API进行硬件配置实施
		if err = agent.applyHardwareSettings(); err != nil {
			agent.reportProgress(-1, "硬件配置失败", "硬件配置失败")
			continue
		}

		// _ = agent.enableIPMILan()

		agent.runPostInstall()

		// 4、安装并重启
		if agent.isByImage() && !agent.isByWindowsImage() {
			// Linux镜像安装
			agent.reportProgress(0.55, "Linux镜像安装", "开始镜像安装")
			if _, err = agent.installByLinuxImage(); err != nil {
				agent.reportProgress(-1, "Linux镜像安装失败", err.Error())
				continue
			}

			if err = agent.reboot(); err != nil {
				agent.reportProgress(-1, "重启失败", err.Error())
				continue
			}
			break

		} else {
			if agent.isCentOS6PXE() {
				if err = agent.genPXE4CentOS6(); err != nil {
					agent.reportProgress(-1, "生成CentOS6.x操作系统(UEFI)特制PXE文件", err.Error())
					continue
				}
				agent.reportProgress(0.5, "生成CentOS6.x操作系统(UEFI)特制PXE文件", "完成")
			}

			time.Sleep(3 * time.Second)

			agent.reportProgress(0.55, "BootOS重启", "重启中...")
			if err = agent.rebootFromPXE(); err != nil {
				agent.reportProgress(-1, "重启失败", err.Error())
				continue
			}
			break
		}
	}
	return nil
}

// isCentOS6PXE 当前设备是否是通过PXE方式安装centos6
func (agent *Agent) isCentOS6PXE() bool {
	return agent.settings != nil &&
		agent.settings.InstallType == model.InstallationPXE &&
		agent.settings.OSTemplate != nil &&
		agent.settings.OSTemplate.BootMode == model.BootModeUEFI &&
		strings.Contains(strings.ToLower(agent.settings.OSTemplate.Name), "centos6")
}

// genPXE4CentOS6 为安装centos6.x生成PXE文件
func (agent *Agent) genPXE4CentOS6() (err error) {
	url := fmt.Sprintf("%s/api/cloudboot/v1/devices/%s/centos6/uefi/pxe", agent.ServerAddr, agent.collected.SN)
	reqData := map[string]string{
		"mac": agent.collected.BootOSMac,
	}
	var respData struct {
		Status  string
		Message string
	}
	if err = agent.doPOSTUnmarshal(url, &reqData, &respData); err != nil {
		return err
	}
	if strings.ToLower(respData.Status) != "success" {
		return fmt.Errorf("Status: %s, Message: %s", respData.Status, respData.Message)
	}
	return nil
}

// isByImage 是否是镜像安装方式安装OS
func (agent *Agent) isByImage() (yes bool) {
	return agent.settings != nil && agent.settings.InstallType == model.InstallationImage
}

// isByWindowsImage 是否是windows镜像安装方式安装OS
func (agent *Agent) isByWindowsImage() (yes bool) {
	return agent.isByImage() && agent.settings.OSTemplate != nil && strings.HasPrefix(strings.ToLower(agent.settings.OSTemplate.Family), "win")
}

// BootMode 引导模式
type BootMode string

var (
	// LegacyBIOS 传统的BIOS引导模式
	LegacyBIOS BootMode = "legacy_bios"
	// UEFI UEFI引导模式
	UEFI BootMode = "uefi"
)

// bootMode 返回系统引导模式
func (agent *Agent) bootMode() (mode BootMode) {
	info, err := os.Stat("/sys/firmware/efi")
	if err == nil && info.IsDir() {
		return UEFI
	}
	return LegacyBIOS
}

// waitToEnterQueue 一直等待，直到进入安装队列。
func (agent *Agent) waitToEnterQueue() {
	t := time.NewTicker(time.Duration(agent.LoopInterval) * time.Second) // 轮询间隔时间
LOOP:
	for {
		select {
		case <-t.C:
			if yes, err := agent.inQueue(agent.collected.SN); err == nil && yes {
				t.Stop()
				break LOOP
			}
		}
	}
}

// inQueue 判断设备是否在安装队列
func (agent *Agent) inQueue(sn string) (yes bool, err error) {
	var respData struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Content struct {
			Result bool `json:"result"`
		} `json:"content"`
	}
	url := fmt.Sprintf("%s/api/cloudboot/v1/devices/%s/is-in-install-list", agent.ServerAddr, sn)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		agent.log.Error(err)
		return false, err
	}
	return respData.Content.Result, nil
}

// loadDeviceSettings 从服务端加载当前设备的装机参数信息
func (agent *Agent) loadDeviceSettings() (err error) {
	agent.log.Infof("Start loading device settings")
	if agent.collected == nil || agent.collected.SN == "" {
		return ErrCollectDeviceInfoFirst
	}

	url := fmt.Sprintf("%s/api/cloudboot/v1/devices/%s/settings", agent.ServerAddr, agent.collected.SN)
	var respData struct {
		Status  string
		Message string
		Content service.DeviceSetting
	}
	if err = httplib.Get(url).ToJSON(&respData); err != nil {
		agent.log.Error(err)
		return err
	}
	agent.log.Debugf("%#v", respData)

	if strings.ToLower(respData.Status) != "success" {
		agent.log.Errorf("GET %s, response status=%s, response message=%s", url, respData.Status, respData.Message)
		return fmt.Errorf("response status=%s, response message=%s", respData.Status, respData.Message)
	}

	agent.settings = &respData.Content
	return nil
}

// reportProgress 上报执行结果
func (agent *Agent) reportProgress(installProgress float64, title, installLog string) bool {
	var reqData struct {
		InstallProgress float64 `json:"progress"`
		InstallLog      string  `json:"log"`
		Title           string  `json:"title"`
	}
	reqData.InstallProgress = installProgress
	reqData.Title = title
	reqData.InstallLog = base64.StdEncoding.EncodeToString([]byte(installLog)) // base64编码

	var respData struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	url := fmt.Sprintf("%s/api/cloudboot/v1/devices/%s/installations/progress", agent.ServerAddr, agent.collected.SN)
	if err := agent.doPOSTUnmarshal(url, &reqData, &respData); err != nil {
		return false
	}
	//return strings.ToLower(respData.Status) == "success"
	if strings.ToLower(respData.Status) == "success" {
		return true
	}
	panic("report progress fail, respData.Status: " + respData.Status)
}

const (
	pxeboot4uefi = "ipmitool chassis bootdev pxe options=efiboot; ipmitool chassis power cycle; sleep 5s; ipmitool power reset"
	pxeboot4bios = "ipmitool chassis bootdev pxe; ipmitool chassis power cycle; sleep 5s; ipmitool power reset"
)

// rebootFromPXE 从PXE重启系统
func (agent *Agent) rebootFromPXE() error {
	var cmds string
	switch agent.bootMode() {
	case UEFI:
		cmds = pxeboot4uefi
	case LegacyBIOS:
		cmds = pxeboot4bios
	}

	out0, err0 := sh.ExecOutputWithLog(agent.log, cmds)
	if err0 == nil {
		return nil
	}

	// 进行重启的第二选择尝试
	cmds2th := `fdisk -lu | awk '/^Disk.*bytes/ { gsub(/:/, ""); system("dd if=/dev/zero of="$2" bs=512 count=1") }'; reboot -f`
	out1, err1 := sh.ExecOutputWithLog(agent.log, cmds2th)
	if err1 == nil {
		return nil
	}
	return fmt.Errorf("%s==>\n%s\n%s\n%s==>\n%s\n%s", cmds, err0, string(out0), cmds2th, err1, string(out1))
}

// reboot 重启bootos
func (agent *Agent) reboot() error {
	cmdArgs := "reboot"
	if output, err := sh.ExecOutputWithLog(agent.log, cmdArgs); err != nil {
		return fmt.Errorf("reboot error: \n#%s\n%v\n%s", cmdArgs, err, string(output))
	}
	return nil
}

// installByLinuxImage 调用命令进行Linux镜像安装
func (agent *Agent) installByLinuxImage() (output []byte, err error) {
	return sh.ExecOutputWithLog(agent.log, `/usr/local/bin/linuxinstall`)
}

// doPOSTUnmarshal 将指定数据序列化成JSON并通过HTTP POST发送到远端，并将JSON格式的响应信息反序列化到respData中。
// respData必须为非nil的指针类型，否则将返回对应的错误。
func (agent *Agent) doPOSTUnmarshal(url string, reqData, respData interface{}) error {
	respBody, err := agent.doPOST(url, reqData)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(respBody, respData); err != nil {
		agent.log.Error(err)
		return err
	}
	return nil
}

// doPOST 将指定数据序列化成JSON并通过HTTP POST发送到远端
func (agent *Agent) doPOST(url string, reqData interface{}) ([]byte, error) {
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		agent.log.Error(err)
		return nil, err
	}

	agent.log.Debugf("POST %s, request body: %s", url, reqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		agent.log.Error(err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		agent.log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	agent.log.Debugf("POST %s, response body: %s", url, respBody)

	if resp.StatusCode != http.StatusOK {
		agent.log.Errorf("POST %s, response status code: %d", url, resp.StatusCode)
		return nil, fmt.Errorf("http status code: %d", resp.StatusCode)
	}
	return respBody, nil
}

// wgetO 模拟'wget -X GET -O'，请求一个HTTP URL并将响应信息以文件形式保存。
func (agent *Agent) wgetO(url, dstFilename string, perm os.FileMode) (err error) {
	// 规避风险: 不可信Url作为输入提供给HTTP请求
	domainCheck := "osinstall.idcos.com"
	if !strings.HasSuffix(url, domainCheck) {
		return fmt.Errorf("URL %s 不满足域名 %s\n", url, domainCheck)
	}
	resp, err := http.Get(url)
	if err != nil {
		agent.log.Error(err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		agent.log.Error(err)
		return err
	}
	if err = ioutil.WriteFile(dstFilename, body, perm); err != nil {
		agent.log.Error(err)
		return err
	}
	return nil
}

func parseInterval(interval string, defValue int) int {
	i, err := strconv.Atoi(strings.TrimSpace(interval))
	if err != nil || i <= 0 {
		return defValue
	}
	return i
}

// InstallationStatus 设备操作系统安装状态
type InstallationStatus struct {
	Type     string  `json:"type"`
	Status   string  `json:"status"`
	Progress float64 `json:"progress"`
}

// getInstallationStatus 查询设备装机状态
func (agent *Agent) getInstallationStatus() (status *InstallationStatus, err error) {
	var respData struct {
		Status  string             `json:"status"`
		Message string             `json:"message"`
		Content InstallationStatus `json:"content"`
	}

	url := fmt.Sprintf("%s/api/cloudboot/v1/devices/%s/installations/status", agent.ServerAddr, agent.collected.SN)
	if err = httplib.Get(url).ToJSON(&respData); err != nil {
		agent.log.Error(err)
		return nil, err
	}

	if respData.Status != "success" {
		return nil, fmt.Errorf("unexpected response value(status=%s, message=%s)", respData.Status, respData.Message)
	}
	return &respData.Content, nil
}
