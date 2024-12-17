package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/astaxie/beego/httplib"

	"idcos.io/cloudboot/build"
	"idcos.io/cloudboot/config"
)

// logger 全局终端日志
var logger = log.New(os.Stdout, "", log.Lshortfile)

var (
	domain             string
	rootPath           string
	copySrcFile        string
	defaultCopySrcFile = filepath.Join(fmt.Sprintf("Z:%c", os.PathSeparator), "windows", "firstboot", "peconfig.exe")
	defaultRootPath    = filepath.Join(fmt.Sprintf("X:%c", os.PathSeparator), "Windows", "System32")
)

func main() {
	flag.StringVar(&domain, "domain", "osinstall", "目标域名")
	flag.StringVar(&copySrcFile, "copy-file-from", defaultCopySrcFile, "待拷贝的源文件")
	flag.StringVar(&rootPath, "root-path", defaultRootPath, "目标程序放置的根目录")
	flag.Parse()

	fmt.Println("peagent", build.Version())
	fmt.Println()

	time.Sleep(30 * time.Second) // 提供足够的时间进行网络初始化

	ip, err := LookupIP(domain)
	if err != nil {
		logger.Println(err)
		return
	}
	fmt.Println("Lookup domain found IP: ",ip)

	retries := 3 // Samba挂载失败时重试次数
	for i := 0; i < retries; i++ {
		//default samba
		sambaSetting := &config.Samba{
			Server:   ip,
			User:     "",
			Password: "",
		}
		fmt.Println("Set samba IP to ",sambaSetting.Server)
		if s, err := getSambaSettings(ip); err == nil && s != nil {
			sambaSetting.User = s.User
			sambaSetting.Password = s.Password
		}
		if err = MountSamba(sambaSetting); err != nil {
			time.Sleep(time.Duration((i+1)*5) * time.Second)
			continue
		}
		break
	}
	if err != nil {
		logger.Println(err)
		return
	}

	if err := Copy(copySrcFile, rootPath); err != nil {
		logger.Println(err)
		return
	}

	if err := RunEXE(filepath.Join(rootPath, "peconfig.exe")); err != nil {
		logger.Println(err)
		return
	}
}

var pingReg = regexp.MustCompile(`(.+)(\s)(\d+)\.(\d+)\.(\d+)\.(\d+)([:|\s])(.+)TTL`)

// ErrIPNotFound 无法查找到IP
var ErrIPNotFound = errors.New("ip not found")

// LookupIP 根据域名查找IP
func LookupIP(domain string) (ip string, err error) {
	output, err := ExecOutput(fmt.Sprintf("ping %s", domain))
	if err != nil {
		return "", err
	}
	regResult := pingReg.FindStringSubmatch(string(output))
	if regResult == nil || len(regResult) != 9 {
		return "", ErrIPNotFound
	}
	return fmt.Sprintf("%s.%s.%s.%s",
		strings.TrimSpace(regResult[3]),
		strings.TrimSpace(regResult[4]),
		strings.TrimSpace(regResult[5]),
		strings.TrimSpace(regResult[6]),
	), nil
}

// MountSamba 挂载Samba目录
func MountSamba(samba *config.Samba) (err error) {
	if SambaMounted() {
		return nil
	}

	var cmdAndArgs string
	if samba.User == "" {
		cmdAndArgs = fmt.Sprintf(`net use Z: \\%s\samba`, samba.Server)
	} else {
		cmdAndArgs = fmt.Sprintf(`net use Z: \\%s\samba /user:%s %s`, samba.Server, samba.User, samba.Password)
	}
	_, err = ExecOutputDesensitization(cmdAndArgs, samba.Password)
	return err
}

// SambaMounted 判断Samba是否已经挂载
func SambaMounted() (mounted bool) {
	_, err := ExecOutput(`net use Z:`)
	return err == nil
}

// Copy 拷贝文件
func Copy(src string, dest string) (err error) {
	logger.Printf("Start copying files from %s to %s", src, dest)
	_, err = ExecOutput(fmt.Sprintf("xcopy /s /e /y /i %s %s", src, dest))
	return err
}

// RunEXE 运行可执行文件
func RunEXE(filepath string, args ...string) error {
	cmd := exec.Command(filepath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ExecOutput windows系统下，执行命令字符串cmdAndArgs，并将命令执行的标准输出和标准错误输出都通过字节切片output返回。
func ExecOutput(cmdAndArgs string) (output []byte, err error) {
	scriptFile, err := GenTempScript([]byte(cmdAndArgs))
	if err != nil {
		return nil, err
	}
	defer os.Remove(scriptFile)

	output, err = exec.Command("cmd", "/c", scriptFile).Output()
	logger.Printf("%s ==>\n%s\n", cmdAndArgs, output)

	return output, err
}

// ExecOutput windows系统下，执行命令字符串cmdAndArgs，并将命令执行的标准输出和标准错误输出都通过字节切片output返回。
func ExecOutputDesensitization(cmdAndArgs string, args ...string) (output []byte, err error) {
	scriptFile, err := GenTempScript([]byte(cmdAndArgs))
	if err != nil {
		return nil, err
	}
	defer os.Remove(scriptFile)

	output, err = exec.Command("cmd", "/c", scriptFile).Output()

	opLog := string(output)
	for _, arg := range args {
		cmdAndArgs = strings.Replace(cmdAndArgs, arg, "******", -1)
		opLog = strings.Replace(opLog, arg, "******", -1)
	}
	logger.Printf("%s ==>\n%s\n", cmdAndArgs, opLog)

	return output, err
}

// GenTempScript 在系统临时目录生成bat脚本文件
func GenTempScript(content []byte) (scriptFile string, err error) {
	scriptFile = filepath.Join(os.TempDir(), fmt.Sprintf("%d.bat", time.Now().Unix()))
	if err = ioutil.WriteFile(scriptFile, content, 0744); err != nil {
		return "", err
	}
	return scriptFile, nil
}

// getSambaSettings 查询samba服务的连接配置
func getSambaSettings(server string) (*config.Samba, error) {
	var respData struct {
		Status  string       `json:"status"`
		Message string       `json:"message"`
		Content config.Samba `json:"content"`
	}

	url := fmt.Sprintf("http://%s/api/cloudboot/v1/samba/settings", server)
	if err := httplib.Get(url).ToJSON(&respData); err != nil {
		return nil, err
	}

	if respData.Status != "success" {
		return nil, fmt.Errorf("unexpected response value(status=%s, message=%s)", respData.Status, respData.Message)
	}
	return &respData.Content, nil
}
