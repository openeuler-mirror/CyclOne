package agent

import (
	"errors"
	"fmt"
	"strings"

	"github.com/astaxie/beego/httplib"
	"idcos.io/cloudboot/server/cloudbootserver/types/device"
)

var (
	// ErrCollectDeviceInfoFirst 请先采集设备信息
	ErrCollectDeviceInfoFirst = errors.New("please collect the device information first")
)

// postDevice 上报设备信息
func (agent *Agent) postDevice() (err error) {
	url := fmt.Sprintf("%s/api/cloudboot/v1/devices/collections?lang=en-US", agent.ServerAddr)
	var respData struct {
		Status  string
		Message string
	}
	if err = agent.doPOSTUnmarshal(url, agent.collected, &respData); err != nil {
		return err
	}
	if strings.ToLower(respData.Status) != "success" {
		return fmt.Errorf("status: %s, message: %s", respData.Status, respData.Message)
	}
	return nil
}

// collect 调用hw-server的API完成设备信息采集
func (agent *Agent) collect() (err error) {
	url := fmt.Sprintf("%s/api/cloudboot/hw/v1/devices/collections?lang=en-US", strings.TrimSuffix(agent.hwSrvBaseURL, "/"))

	var respData struct {
		Status  string        `json:"status"`
		Message string        `json:"message"`
		Content device.Device `json:"content"`
	}
	if err = httplib.Get(url).Header("Accept", "application/json").ToJSON(&respData); err != nil {
		agent.log.Errorf("GET %s ==>\n%s", url, err.Error())
		return err
	}
	if respData.Status != "success" {
		return errors.New(respData.Message)
	}
	agent.collected = &respData.Content
	agent.collected.Setup()
	return nil
}
