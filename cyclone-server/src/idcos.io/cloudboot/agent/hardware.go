package agent

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/httplib"
)

const hwServerPath = "/usr/bin/hw-server"

// arch 返回当前
func (agent *Agent) arch() string {
	out, _ := exec.Command("uname", "-m").Output()
	return strings.TrimSpace(string(out))
}

// runHWServer 从远端下载hw-server并启动该服务
func (agent *Agent) runHWServer() (err error) {
	url := fmt.Sprintf("%s/hw-tools/%s/hw-server", agent.Cmdline.ServerAddr, agent.arch())
	if err = agent.dumpHWServer(url, hwServerPath); err != nil {
		return err
	}
	_ = os.Chmod(hwServerPath, 0755)

	cmdArgs := []string{
		hwServerPath,
		"--base-URL", agent.Cmdline.ServerAddr,
		"--port", strconv.Itoa(agent.Cmdline.HWServerPort),
	}
	agent.log.Debugf("Run hw-server: %s", strings.Join(cmdArgs, " "))
	return exec.Command(cmdArgs[0], cmdArgs[1:]...).Start()
}

// dumpHWServer 从远端下载hw-server组件并将其转储到指定位置
func (agent *Agent) dumpHWServer(url string, dstFilename string) (err error) {
	agent.log.Debugf("Start dumping hw-server from %s to %s", url, dstFilename)
	if err = httplib.Get(url).ToFile(dstFilename); err != nil {
		agent.log.Errorf("Dumping hw-server error: %s", err)
		return err
	}
	return nil
}

// waitHWServerOK 阻塞直到hw-server组件处于ready状态。
func (agent *Agent) waitHWServerOK() {
	t := time.NewTicker(time.Second)
LOOP:
	for {
		select {
		case <-t.C:
			if agent.hwServerOK() {
				t.Stop()
				break LOOP
			}
		}
	}
	agent.log.Infof("hw-server is ready")
}

// hwServerOK 请求hw-serverAPI以确定该服务是否处于ready状态。
func (agent *Agent) hwServerOK() bool {
	url := fmt.Sprintf("%s/api/cloudboot/hw/v1/ping", strings.TrimSuffix(agent.hwSrvBaseURL, "/"))
	agent.log.Debugf("GET %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return resp.StatusCode == http.StatusOK && string(body) == "pong"
}

// applyHardwareSettings 调用hw-server组件API进行硬件配置实施
func (agent *Agent) applyHardwareSettings() (err error) {
	if agent.collected == nil || agent.collected.SN == "" {
		return ErrCollectDeviceInfoFirst
	}

	agent.log.Infof("Start applying hardware settings")

	var respData struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	url := fmt.Sprintf("%s/api/cloudboot/hw/v1/devices/%s/settings/applyings?lang=en-US", strings.TrimSuffix(agent.hwSrvBaseURL, "/"), agent.collected.SN)
	agent.log.Debugf("POST %s", url)

	if err = httplib.Post(url).Header("Accept", "application/json").SetTimeout(time.Minute, time.Minute*60).ToJSON(&respData); err != nil { // TODO 依旧存在超时风险
		agent.log.Errorf("POST %s error: %s", url, err)
		return err
	}

	if respData.Status != "success" {
		agent.log.Errorf("POST %s, response error: %s", url, respData.Message)
		return errors.New(respData.Message)
	}
	return nil
}
