package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"reflect"

	"idcos.io/cloudboot/config"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/utils/ping"
)

// var Logger *logs.BeeLogger
var Logger logger.Logger
var rootPath = "c:/firstboot"
var logPath = path.Join(rootPath, "log")
var logFile = path.Join(logPath, "setup.log")

var rootPathPeAgent = "X:\\Windows\\System32"
var logFilePeAgent = path.Join(rootPathPeAgent, "pe-agent.log")

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// @Deprecated
func InitFileLog() {
	os.MkdirAll(logPath, 0666)
	Logger = logger.NewBeeLogger(&config.Logger{
		Level:          "debug",
		LogFile:        logFile,
		ConsoleEnabled: true,
	})
}

// @Deprecated
func InitPeAgentLog() {
	Logger = logger.NewBeeLogger(&config.Logger{
		Level:          "debug",
		LogFile:        logFilePeAgent,
		ConsoleEnabled: true,
	})
}

// @Deprecated
func InitConsoleLog() {
	Logger = logger.NewBeeLogger(&config.Logger{
		Level:          "debug",
		LogFile:        logFilePeAgent,
		ConsoleEnabled: true,
	})
}

// ExecCmd 执行 command
func ExecCmd(scriptFile, cmdString string) ([]byte, error) {

	// 生成临时文件
	if CheckFileIsExist(scriptFile) {
		os.Remove(scriptFile)
	}
	file, err := os.Create(scriptFile)
	if err != nil {
		return nil, err
	}
	//defer os.Remove(file.Name())
	defer file.Close()

	if _, err = file.WriteString(cmdString); err != nil {
		return nil, err
	}
	file.Close()

	return ExecScript(scriptFile)
}

// ExecScript exec script
func ExecScript(scriptPath string) ([]byte, error) {

	var cmd = exec.Command("cmd", "/c", scriptPath)
	return cmd.Output()
}

// ExecCmd 执行 command
func ExecCmdOutputWithLogfile(scriptFile, cmdString string) ([]byte, error) {
	// 生成临时文件
	if CheckFileIsExist(scriptFile) {
		os.Remove(scriptFile)
	}
	file, err := os.Create(scriptFile)
	if err != nil {
		return nil, err
	}
	//defer os.Remove(file.Name())
	defer file.Close()

	if _, err = file.WriteString(cmdString); err != nil {
		return nil, err
	}
	file.Close()

	_ = os.MkdirAll(filepath.Join("c:/", "firstboot", "log"), 0777)
	log := filepath.Join(filepath.Join("c:/", "firstboot", "log", "stdout.log"))
	errExec := ExecScriptOutputWithLogfile(scriptFile, log)
	if errExec != nil {
		return nil, errExec
	}

	fi, err := os.Open(log)
	if err != nil {
		return nil, err
	}
	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	return fd, nil
}

func ExecScriptOutputWithLogfile(scriptPath string, log string) error {
	if CheckFileIsExist(log) {
		os.Remove(log)
	}
	cmd := exec.Command("cmd", "/c", scriptPath)
	stdout, err := os.OpenFile(log, os.O_CREATE|os.O_WRONLY, 0600)
	cmd.Stdout = stdout
	if err != nil {
		return err
	}
	defer stdout.Close()
	// 执行命令
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// ExecCmd 执行 command
//func ExecCmdStdoutConsole(scriptFile, cmdString string) ([]byte, error) {
//	// 生成临时文件
//	if CheckFileIsExist(scriptFile) {
//		os.Remove(scriptFile)
//	}
//	file, err := os.Create(scriptFile)
//	if err != nil {
//		return nil, err
//	}
//	//defer os.Remove(file.Name())
//	defer file.Close()
//
//	if _, err = file.WriteString(cmdString); err != nil {
//		return nil, err
//	}
//	file.Close()
//
//	output, errExec := ExecScriptStdoutConsole(scriptFile)
//	if errExec != nil {
//		return output, errExec
//	}
//	return output, nil
//}

func ExecScriptStdoutConsole(scriptPath string) ([]byte, error) {
	var output string
	cmd := exec.Command("cmd", "/c", scriptPath)
	w := bytes.NewBuffer(nil)
	cmd.Stdout = os.Stdout
	cmd.Stderr = w
	cmd.Start() // Start开始执行c包含的命令，但并不会等待该命令完成即返回。
	// Wait方法会返回命令的返回状态码并在命令返回后释放相关的资源。
	if err := cmd.Wait(); err != nil {
		str := fmt.Sprintf("stderr:\n%s\n\nerror:\n%s", string(w.Bytes()), err.Error())
		// fmt.Println(str)
		return []byte(output), errors.New(str)
	}
	return []byte(output), nil
}

// CallRestAPI 调用restful api
func CallRestAPI(url string, jsonReq interface{}) ([]byte, error) {
	var req *http.Request
	var resp *http.Response
	var err error
	var reqBody []byte

	if reqBody, err = json.Marshal(jsonReq); err != nil {
		return nil, err
	}

	fmt.Printf("Request BODY: %s \n", string(reqBody))
	if req, err = http.NewRequest("POST", url, bytes.NewBuffer(reqBody)); err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	if resp, err = http.DefaultClient.Do(req); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code: %d", resp.StatusCode)
	}

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("call url: %s failed", url)
	}

	return body, nil
}

// ReportProgress 上报执行结果
func ReportProgress(installProgress float64, sn, title, installLog string, host string) bool {
	var url = fmt.Sprintf("http://%s/api/cloudboot/v1/devices/%s/installations/progress", host, sn)
	Logger.Debug("ReportProgress url:%s\n", url)
	var jsonReq struct {
		//SN              string  `json:"sn"`
		InstallProgress float64 `json:"progress"`
		InstallLog      string  `json:"log"`
		Title           string  `json:"title"`
	}
	//jsonReq.SN = sn
	jsonReq.InstallProgress = installProgress
	jsonReq.Title = title
	jsonReq.InstallLog = base64.StdEncoding.EncodeToString([]byte(installLog)) // base64编码
	Logger.Debug("SN: %s\n", sn)
	Logger.Debug("InstallProgress: %f\n", jsonReq.InstallProgress)
	Logger.Debug("InstallLog: %s\n", jsonReq.InstallLog)
	Logger.Debug("Title: %s\n", jsonReq.Title)

	var jsonResp struct {
		Status  string
		Message string
		Content struct {
			Result string
		}
	}

	var ret, err = CallRestAPI(url, jsonReq)
	Logger.Debug("ReportProgress api result:%s\n", string(ret))
	if err != nil {
		Logger.Error(err.Error())
		return false
	}

	if err := json.Unmarshal(ret, &jsonResp); err != nil {
		Logger.Error(err.Error())
		return false
	}

	if jsonResp.Status != "success" {
		return false
	}
	return true
}

// PingLoop return when success
func PingLoop(host string, pkgCnt int, timeout int) bool {
	for i := 0; i < pkgCnt; i++ {
		if ping.Ping(host, timeout) {
			return true
		}
	}
	return false
}

func ListDir(dirPth string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if !fi.IsDir() {
			continue
		}
		files = append(files, dirPth+PthSep+fi.Name())
	}
	return files, nil
}

func ListFiles(dirPth string, suffix string, onlyReturnFileName bool) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			if onlyReturnFileName == true {
				files = append(files, fi.Name())
			} else {
				files = append(files, dirPth+PthSep+fi.Name())
			}
		}
	}
	return files, nil
}

// StructTrimSpace 遍历结构体的字段并找出string类型字段进行trimspace操作
// 入参得是一个结构体指针
func StructTrimSpace(i interface{}) {
	ref := reflect.ValueOf(i).Elem()
	for i := 0; i < ref.NumField(); i++ {
		field := ref.Field(i)
		switch field.Kind() {
		case reflect.String:
			newVal := strings.TrimSpace(field.Interface().(string))
			canSetVal := reflect.ValueOf(newVal)
			field.Set(canSetVal)
		}
	}
}

// GetOOBHost 获取OOB主机信息
// 理论上一个SN对应一个主机名。
// 但在实际环境中，某些存量设备
// hp的机器，hostname可能是ILO+$SN(如：ILO123456)
// dell的机器，hostname可能是rac-$SN(如rac-123456)
// 故返回多个，上层用重试法确定实际值
func GetOOBHost(sn, vendor, domain string) []string {
	snPrefixILO := "ILO"
	snPrefixDRAC := "rac-"
	sns := make([]string, 1, 2)
	if domain != "" {
		domain = "." + domain
	}
	vendorLower := strings.ToLower(vendor)
	sns[0] = sn + domain
	if strings.Contains(vendorLower, "dell") || strings.Contains(vendorLower, "戴尔") {
		sns = append(sns, snPrefixDRAC+sn+domain)
	} else if strings.Contains(vendorLower, "hp") || strings.Contains(vendorLower, "惠普") {
		sns = append(sns, snPrefixILO+sn+domain)
	}
	return sns
}

//ToJsonString 转换json为字符串
func ToJsonString(arg interface{}) string {
	bytes, _ := json.Marshal(arg)
	return string(bytes)
}
