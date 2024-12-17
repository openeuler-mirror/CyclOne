package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/astaxie/beego/httplib"
	"github.com/urfave/cli"

	"idcos.io/cloudboot/build"
	"idcos.io/cloudboot/config"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/server/cloudbootserver/types/setting"
	"idcos.io/cloudboot/server/cloudbootserver/util"
	"idcos.io/cloudboot/utils"
	mystrings "idcos.io/cloudboot/utils/strings"
	"idcos.io/cloudboot/utils/win"
)

var (
	systemTemplateFile = "X:\\Windows\\System32\\unattended.xml"

	cachePath          = "c:/firstboot"
	nicFile            = filepath.Join(cachePath, "nic.txt")
	deviceSettingFile  = filepath.Join(cachePath, "deviceSetting.json")
	networkSettingFile = filepath.Join(cachePath, "networkSetting.json")
	preInstallScript   = filepath.Join(cachePath, "preinstall.cmd")
	postInstallScript  = filepath.Join(cachePath, "postinstall.cmd")
)

// Options peconfig组件命令行选项值
type Options struct {
	LogLevel string
	LogFile  string
	Domain   string
}

func main() {
	var opts Options
	var log logger.Logger

	app := cli.NewApp()
	app.Version = build.Version()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "log-level",
			Value:       "debug",
			Usage:       "log level, optional values: debug|info|warn|error",
			Destination: &opts.LogLevel,
		}, cli.StringFlag{
			Name:        "log-file",
			Value:       "peconfig.log",
			Usage:       "log file path name",
			Destination: &opts.LogFile,
		},
		cli.StringFlag{
			Name:        "domain",
			Value:       "osinstall",
			Usage:       "api server domain name",
			Destination: &opts.Domain,
		},
	}
	app.Before = func(ctx *cli.Context) error {
		log = logger.NewBeeLogger(&config.Logger{
			Level:          opts.LogLevel,
			LogFile:        opts.LogFile,
			ConsoleEnabled: false,
			RotateEnabled:  false,
		})
		return nil
	}
	app.Action = func(c *cli.Context) {
		if err := NewPEConfiger(log, &opts).Run(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}

	if err := app.Run(os.Args); err != nil {
		log.Error(err)
	}
}

// InstallationError 操作系统安装过程中发生的错误
type InstallationError struct {
	Title string
	Cause string
}

// NewInstallationError 返回操作系统安装错误实例
func NewInstallationError(title string, cause string) error {
	return &InstallationError{
		Title: title,
		Cause: cause,
	}
}

func (e *InstallationError) Error() string {
	return e.Cause
}

// PEConfiger WinPE配置器
type PEConfiger struct {
	log     logger.Logger
	opts    *Options
	sn      string
	srvAddr string
}

// NewPEConfiger 返回WinPE配置器实例
func NewPEConfiger(log logger.Logger, opts *Options) *PEConfiger {
	return &PEConfiger{
		log:  log,
		opts: opts,
	}
}

// handleErr 对安装过程中产生的致命错误进行处理。
// 对于非安装错误(InstallationError)，则不进行任何处理。
func (pec *PEConfiger) handleErr(err error) {
	if ierr, ok := err.(*InstallationError); ok {
		pec.reportProgress(-1, ierr.Title, ierr.Error())
	}
}

const (
	defaultLoopInterval = 60
)

func (pec *PEConfiger) inQueue() (bool, error) {
	stat, err := pec.getInstallationStatus()
	if err != nil {
		return false, err
	}
	//设备状态必须为等待安装或者正在安装，才允许继续执行。
	if stat.Status != model.InstallStatusPre && stat.Status != model.InstallStatusIng {
		pec.log.Errorf("%s is not in the installation queue\n", pec.sn)
		return false, err
	}
	return true, nil
}

// waitToEnterQueue 一直等待，直到进入安装队列。
func (pec *PEConfiger) waitToEnterQueue() {
	t := time.NewTicker(time.Duration(defaultLoopInterval) * time.Second) // 轮询间隔时间
LOOP:
	for {
		select {
		case <-t.C:
			if yes, err := pec.inQueue(); err == nil && yes {
				t.Stop()
				break LOOP
			}
		}
	}
}

// Run 运行WinPE配置器
func (pec *PEConfiger) Run() (err error) {
	defer func(e *error) {
		if err := recover(); err != nil {
			pec.handleErr(*e)
		}
	}(&err)

	time.Sleep(30 * time.Second) // 提供足够的时间进行网络初始化

	fmt.Println("[1/7] Get SN")
	if err = pec.loadSN(); err != nil {
		return NewInstallationError("Failed to get SN", err.Error())
	}
	fmt.Printf("SN ==> %s\n\n", pec.sn)

	fmt.Println("[2/7] Lookup IP")
	ip, _ := pec.lookupIP(pec.opts.Domain)
	if ip != "" {
		pec.srvAddr = ip
	} else {
		pec.srvAddr = pec.opts.Domain
	}
	fmt.Printf("%s ==> %s\n\n", "osinstall", ip)

	if !utils.PingLoop(pec.srvAddr, 30, 2) {
		return NewInstallationError(
			fmt.Sprintf("The cloudboot server(%s) is unreachable", pec.srvAddr),
			"ping timeout",
		)
	}

	fmt.Printf("[3/7] Get installation status\n")
	var stat *InstallationStatus
	stat, err = pec.getInstallationStatus()
	if err != nil {
		pec.log.Errorf("Failed to get installation status, Error: %s\n", err.Error())
	}
	fmt.Printf("status ==> %s %s %f\n\n", stat.Type, stat.Status, stat.Progress)
	//设备状态必须为等待安装或者正在安装，才允许继续执行。
	if stat.Status != model.InstallStatusPre && stat.Status != model.InstallStatusIng {
		pec.log.Errorf("%s is not in the installation queue\n", pec.sn)
	}

	for {
		pec.Post() //上传日志
		//轮询是否进入安装队列
		pec.waitToEnterQueue()
		fmt.Printf("Enter installation queue\n")

		pec.reportProgress(0.6, "Start installation", fmt.Sprintf("Start installation by %s", stat.Type))

		switch stat.Type {
		case model.InstallationPXE:
			fmt.Printf("[4/7] Installation by system template\n")
			fmt.Printf("pull system template\n")
			var data []byte
			if data, err = pec.pullSystemTemplate(); err != nil {
				pec.log.Errorf("Failed to pull system template , Error: %s \n", err.Error())
				continue
			}

			fmt.Printf("persist system template")
			if err = ioutil.WriteFile(systemTemplateFile, data, 0744); err != nil {
				pec.log.Errorf("Failed to persist system template , Error: %s\n", err.Error())
				continue
			}

			if err = pec.installBySystemTemplate(data); err != nil {
				pec.log.Errorf("Installation failed(system template) , Error: %s\n", err.Error())
				continue
			}
		case model.InstallationImage:
			fmt.Printf("[4/7] Installation by image template\n")
			fmt.Printf("pull image template\n")
			var tpl *ImageTemplate
			tpl, err = pec.pullImageTemplate()
			if err != nil {
				pec.log.Errorf("Failed to pull image template , Error: %s\n", err.Error())
				continue
			}

			if err = pec.doPartitions(tpl.Disks); err != nil {
				pec.log.Errorf("Failed to partition or quick format , Error: %s\n", err.Error())
				continue
			}

			if err = pec.applyImage(tpl.URL); err != nil {
				pec.log.Errorf("Failed to apply-image , Error: %s \n", err.Error())
				continue
			}

			if err = pec.bcdboot(); err != nil {
				pec.log.Errorf("Failed to bcdboot , Error: %s \n", err.Error())
				continue
			}

			if err = pec.addDrivers(tpl.Name); err != nil {
				pec.log.Errorf("Failed to add drivers , Error: %s \n", err.Error())
				continue
			}

			_ = os.MkdirAll(cachePath, 0755)

			if tpl.PreScript != "" {
				pec.log.Debug("persist pre-install script ")
				if err = ioutil.WriteFile(preInstallScript, []byte(tpl.PreScript), 0755); err != nil {
					pec.log.Debug("x\n")
				} else {
					pec.log.Debug("√\n")
				}
			}

			if tpl.PostScript != "" {
				pec.log.Debug("persist post-install script ")
				if err = ioutil.WriteFile(postInstallScript, []byte(tpl.PostScript), 0755); err != nil {
					pec.log.Debug("x\n")
				} else {
					pec.log.Debug("√\n")
				}
			}
		}

		fmt.Printf("[5/7] Copy %s directory\n", `Z:\windows\firstboot`)
		if err = pec.copyFirstBoot(`Z:\windows\firstboot`); err != nil {
			pec.log.Errorf("Failed to copy 'firstboot' directory , Error: %s\n", err.Error())
			continue
		}

		fmt.Printf("[6/7] Persistent device settings\n")
		//_ = pec.saveDeviceSettings()
		if errs := pec.saveDeviceSettings(); errs != nil {
			pec.log.Errorf("Failed to save device settings, Error: %s\n", errs)
			break
		}

		pec.reportProgress(0.7, "Diskpart", "diskpart")
		pec.reportProgress(0.75, "Change hostname", "change hostname")
		pec.reportProgress(0.8, "Change network configuration", "change network")
		pec.reportProgress(0.9, "Change registry configuration", "change registry")
		pec.reportProgress(1.0, "The installation is complete", "finish")

		fmt.Printf("[7/7] Reboot after 10s\n")
		time.Sleep(10 * time.Second)
		if err = pec.reboot(); err != nil {
			pec.log.Errorf("Reboot failed , Error: %s\n", err.Error())
			continue
		} else {
			break
		}
	}
	return nil
}

func (pec *PEConfiger) isVM() (bool, error) {
	output, err := win.ExecOutputWithLog(pec.log, `systeminfo`)
	if err != nil {
		return false, err
	}
	return regexp.MustCompile(`(?i)VMware|VirtualBox|KVM|Xen|Parallels`).MatchString(string(output)), nil
}

func (pec *PEConfiger) loadSN() (err error) {
	if vm, _ := pec.isVM(); vm {
		pec.sn, err = win.MacAddress(pec.log)
		return err
	}
	pec.sn, err = win.PhysicalMachineSN(pec.log)
	return err
}

func (pec *PEConfiger) lookupIP(domain string) (ip string, err error) {
	output, err := win.ExecOutputWithLog(pec.log, fmt.Sprintf("ping %s", domain))
	if err != nil {
		return "", err
	}

	arr := regexp.MustCompile(`(.+)(\s)(\d+)\.(\d+)\.(\d+)\.(\d+)([:|\s])(.+)TTL`).FindStringSubmatch(string(output))
	if len(arr) != 9 {
		return "", nil
	}
	return fmt.Sprintf("%s.%s.%s.%s",
		strings.TrimSpace(arr[3]),
		strings.TrimSpace(arr[4]),
		strings.TrimSpace(arr[5]),
		strings.TrimSpace(arr[6]),
	), nil
}

// InstallationStatus 设备操作系统安装状态
type InstallationStatus struct {
	Type     string  `json:"type"`
	Status   string  `json:"status"`
	Progress float64 `json:"progress"`
}

// getInstallationStatus 查询设备装机状态
func (pec *PEConfiger) getInstallationStatus() (status *InstallationStatus, err error) {
	var respData struct {
		Status  string             `json:"status"`
		Message string             `json:"message"`
		Content InstallationStatus `json:"content"`
	}

	url := fmt.Sprintf("http://%s/api/cloudboot/v1/devices/%s/installations/status", pec.srvAddr, pec.sn)
	if err = httplib.Get(url).ToJSON(&respData); err != nil {
		pec.log.Error(err)
		return nil, err
	}

	if respData.Status != "success" {
		return nil, fmt.Errorf("unexpected response value(status=%s, message=%s)", respData.Status, respData.Message)
	}
	return &respData.Content, nil
}

// reportProgress 上报安装进度
func (pec *PEConfiger) reportProgress(progress float64, title, logMsg string) error {
	var reqData struct {
		InstallProgress float64 `json:"progress"`
		InstallLog      string  `json:"log"`
		Title           string  `json:"title"`
	}
	reqData.InstallProgress = progress
	reqData.Title = title
	reqData.InstallLog = base64.StdEncoding.EncodeToString([]byte(logMsg)) // base64编码

	var respData struct {
		Status  string
		Message string
		Content struct {
			Result string
		}
	}

	url := fmt.Sprintf("http://%s/api/cloudboot/v1/devices/%s/installations/progress?lang=en-US", pec.srvAddr, pec.sn)
	reqBody, _ := json.Marshal(&reqData)
	pec.log.Debugf("POST %s\rRequest body ==>\n%s", url, reqBody)

	if err := httplib.Post(url).Body(reqBody).Header("Content-Type", "application/json").ToJSON(&respData); err != nil {
		pec.log.Error(err)
		return err
	}

	respBody, _ := json.Marshal(&respData)
	pec.log.Debugf("Response body ==>\n%s", respBody)

	if strings.ToLower(respData.Status) != "success" {
		return errors.New(respData.Message)
	}
	return nil
}

// pullSystemTemplate 拉取系统模板内容
func (pec *PEConfiger) pullSystemTemplate() (data []byte, err error) {
	url := fmt.Sprintf("http://%s/api/cloudboot/v1/devices/%s/settings/system-template?type=raw", pec.srvAddr, pec.sn)
	pec.log.Debugf("GET %s", url)

	resp, err := http.Get(url)
	if err != nil {
		pec.log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		pec.log.Error(err)
		return nil, err
	}
	return []byte(mystrings.UNIX2DOS(string(data))), nil
}

func (pec *PEConfiger) installBySystemTemplate(systemTemplate []byte) (err error) {
	arr := regexp.MustCompile(`<Path>(.*)\\install.wim</Path>`).FindStringSubmatch(string(systemTemplate))
	if len(arr) != 2 {
		return errors.New("'install.wim' not found in system template")
	}

	cmdArgs := fmt.Sprintf(`%s /unattend:unattended.xml /noreboot`, filepath.Join(arr[1], "setup.exe"))
	// fmt.Printf("start installation ==> %s\n", cmdArgs)
	_, err = win.ExecOutputWithLog(pec.log, cmdArgs)
	return err
}

// getThenSave 请求URL并将response body内容写入目标文件。
func (pec *PEConfiger) getThenSave(url, destFile string) (err error) {
	pec.log.Debugf("GET %s", url)
	return httplib.Get(url).ToFile(destFile)
}

// Partition 分区
type Partition struct {
	Name       string `json:"name"`
	Size       string `json:"size"`
	Fstype     string `json:"fstype"`
	Mountpoint string `json:"mountpoint"`
}

// Disk 逻辑磁盘
type Disk struct {
	Name       string      `json:"name"`
	Partitions []Partition `json:"partitions"`
}

// ImageTemplate 镜像安装模板
type ImageTemplate struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Disks      []Disk `json:"disks"`
	PreScript  string `json:"pre_script"`
	PostScript string `json:"post_script"`
}

// pullImageTemplate 从服务端拉取镜像模板
func (pec *PEConfiger) pullImageTemplate() (*ImageTemplate, error) {
	var respData struct {
		Status  string
		Message string
		Content ImageTemplate
	}
	url := fmt.Sprintf("http://%s/api/cloudboot/v1/devices/%s/settings/image-template", pec.srvAddr, pec.sn)

	if err := httplib.Get(url).ToJSON(&respData); err != nil {
		pec.log.Error(err)
		return nil, err
	}

	if respData.Status != "success" {
		return nil, fmt.Errorf("failed to pull the image template: %s", respData.Message)
	}

	if respData.Content.PreScript != "" {
		if preData, err := base64.StdEncoding.DecodeString(respData.Content.PreScript); err == nil {
			respData.Content.PreScript = mystrings.UNIX2DOS(string(preData))
		}
	}

	if respData.Content.PostScript != "" {
		if postData, err := base64.StdEncoding.DecodeString(respData.Content.PostScript); err == nil {
			respData.Content.PostScript = mystrings.UNIX2DOS(string(postData))
		}
	}
	return &respData.Content, nil
}

func (pec *PEConfiger) doPartitions(disks []Disk) (err error) {
	// TODO 待优化
	if len(disks) > 0 {
		return DiskSlice(disks).ToDiskPartConfigurations().Apply(pec.log)
	}
	p, err := pec.getPartitionSetting()
	if err != nil {
		return pec.formatDiskWithSignalFile()
	}
	return p.Apply(pec.log)
}

func (pec *PEConfiger) applyImage(imageURL string) error {
	cmdArgs := fmt.Sprintf(`Dism /apply-image /imagefile:%s /index:1 /ApplyDir:C:\`, imageURL)
	// fmt.Printf("apply image ==> %s\n", cmdArgs)
	_, err := win.ExecOutputWithLog(pec.log, cmdArgs)
	return err
}

func (pec *PEConfiger) bcdboot() error {
	cmdArgs := `bcdboot C:\windows /s C: /l zh-CN /v`
	// fmt.Printf("bcdboot ==> %s\n", cmdArgs)
	_, err := win.ExecOutputWithLog(pec.log, cmdArgs)
	return err
}

func (pec *PEConfiger) addDrivers(imageTPLName string) error {
	imageTPLName = strings.ToLower(imageTPLName)
	osName := "2012r2"
	if strings.Contains(imageTPLName, "2008") {
		osName = "2008r2"
	} else if strings.Contains(imageTPLName, "win7") {
		osName = "win7"
	} else if strings.Contains(imageTPLName, "win10") {
		osName = "win10"
	}
	cmdArgs := fmt.Sprintf(`Dism /Image:C:\ /Add-Driver /Driver:Z:\windows\drivers\%s\ /Recurse /forceunsigned`, osName)
	// fmt.Printf("add drivers ==> %s\n", cmdArgs)
	_, err := win.ExecOutputWithLog(pec.log, cmdArgs)
	return err
}

func (pec *PEConfiger) copyFirstBoot(srcDir string) error {
	destDir := `C:\firstboot`
	// 可能存在CD-ROM\EFI等其他设备占用盘符
	if utils.CheckFileIsExist("d:/windows") {
		destDir = `D:\firstboot`
	}
	if utils.CheckFileIsExist("e:/windows") {
		destDir = `E:\firstboot`
	}
	if utils.CheckFileIsExist("f:/windows") {
		destDir = `F:\firstboot`
	}

	cmdArgs := fmt.Sprintf(`xcopy /s /e /y /i %s %s`, srcDir, destDir)
	_, err := win.ExecOutputWithLog(pec.log, cmdArgs)
	return err
}

func (pec *PEConfiger) saveDeviceSettings() (errs []error) {
	// 可能存在CD-ROM\EFI等其他设备占用盘符
	if utils.CheckFileIsExist("d:/windows") {
		cachePath          = "d:/firstboot"
		deviceSettingFile  = filepath.Join(cachePath, "deviceSetting.json")
		networkSettingFile = filepath.Join(cachePath, "networkSetting.json")
	}
	if utils.CheckFileIsExist("e:/windows") {
		cachePath          = "e:/firstboot"
		deviceSettingFile  = filepath.Join(cachePath, "deviceSetting.json")
		networkSettingFile = filepath.Join(cachePath, "networkSetting.json")
	}
	if utils.CheckFileIsExist("f:/windows") {
		cachePath          = "f:/firstboot"
		deviceSettingFile  = filepath.Join(cachePath, "deviceSetting.json")
		networkSettingFile = filepath.Join(cachePath, "networkSetting.json")
	}
	pairs := map[string]string{
		fmt.Sprintf("http://%s/api/cloudboot/v1/devices/%s/settings", pec.srvAddr, pec.sn):          deviceSettingFile,
		fmt.Sprintf("http://%s/api/cloudboot/v1/devices/%s/settings/networks", pec.srvAddr, pec.sn): networkSettingFile,
	}
	for url, dest := range pairs {
		fmt.Printf("Save %s ", dest)
		if err := pec.getThenSave(url, dest); err != nil {
			fmt.Println("x")
			errs = append(errs, err)
		} else {
			fmt.Println("√")
		}
	}
	return errs
}

func (pec *PEConfiger) reboot() error {
	_, err := win.ExecOutputWithLog(pec.log, `wpeutil reboot`)
	return err
}

// TODO 待优化
func (pec *PEConfiger) getPartitionSetting() (win.DiskPartConfigurations, error) {
	var part win.DiskPartConfigurations
	var url = fmt.Sprintf("http://%s/api/cloudboot/v1/devices/%s/settings/partitions", pec.srvAddr, pec.sn)
	pec.log.Debug(url)
	// 规避风险: 不可信Url作为输入提供给HTTP请求
	domainCheck := "osinstall.idcos.com"
	if !strings.HasSuffix(url, domainCheck) {
		return part, fmt.Errorf("URL %s 不满足域名 %s\n", url, domainCheck)
	}	
	resp, err := http.Get(url)
	if err != nil {
		return part, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return part, fmt.Errorf("http status code: %d", resp.StatusCode)
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return part, fmt.Errorf("call url: %s failed", url)
	}
	pec.log.Debug(string(body))
	var jsonResp struct {
		Content map[string][]setting.PartitionSettingItem
		Status  string
	}
	errJSON := json.Unmarshal(body, &jsonResp)
	if errJSON != nil {
		return part, errJSON
	}
	if jsonResp.Status == "success" {
		if items, ok := jsonResp.Content["items"]; ok {
			ps := make([]win.DiskPartConfiguration, 0)
			part := make([]win.PartConfiguration, 0)
			markIdx := 0
			for i := range items {
				part = append(part, win.PartConfiguration{
					Size:       items[i].Size,
					FSType:     items[i].FS,
					Mountpoint: items[i].Mountpoint,
				})
				if i < len(items)-1 && (items[i+1].Disk != items[i].Disk) {
					ps = append(ps, win.DiskPartConfiguration{
						Disk:       (items[i]).Disk,
						Partitions: part[markIdx : i+1],
					})
					markIdx = i + 1
				} else if i == len(items)-1 {
					ps = append(ps, win.DiskPartConfiguration{
						Disk:       (items[i]).Disk,
						Partitions: part[markIdx:],
					})
				}
			}
			return ps, nil
		}
	}
	return nil, nil
}

// TODO 待优化
func (pec *PEConfiger) formatDiskWithSignalFile() error {
	var cmd = `fsutil fsinfo drives`
	outputBytes, err := win.ExecOutputWithLog(pec.log, cmd)
	if err != nil {
		return err
	}
	//output:Drives: C:\ D:\ E:\ Z:\
	var output = string(outputBytes)
	//先替换Drives: 再通过 :\ 分割后遍历
	fix := util.SubStrByByte(output, strings.IndexAny(output, ":"))
	output = strings.Replace(output, fix+":", "", -1)
	list := strings.Split(output, ":\\")
	for _, volume := range list {
		volume = strings.TrimSpace(volume)
		if volume == "" {
			continue
		}

		var file = fmt.Sprintf(`%s:/format.txt`, volume)
		//格式化标识文件不存在，则略过
		if !utils.CheckFileIsExist(file) {
			continue
		}

		//格式化
		cmd = fmt.Sprintf(`format %s: /quick`, volume)
		_, err := win.ExecOutputWithLog(pec.log, cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

// loadData 加载日志
func (pec *PEConfiger) loadData() ([]byte, error) {
	data, err := ioutil.ReadFile("peconfig.log")
	if err != nil {
		pec.log.Error(err)
		return nil, err
	}
	return data, nil
}

// Post 向服务端发送日志
func (pec *PEConfiger) Post() (err error) {
	data, err := pec.loadData()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("http://%s/api/cloudboot/v1/devices/%s/components/%s/logs?lang=en-US",
		pec.srvAddr, pec.sn,
		"peconfig")

	resp, err := httplib.Post(url).Header("Accept", "application/json").Body(data).DoRequest()
	if err != nil {
		pec.log.Error(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("post log failed: %s", http.StatusText(resp.StatusCode))
	}
	return nil
}
