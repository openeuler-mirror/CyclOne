package oob

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils/ping"
	"idcos.io/cloudboot/utils/sh"
)

var (
	// ErrMissingOOBInfo OOB信息缺失
	ErrMissingOOBInfo = errors.New("missing oob information")
	// ErrOOBIPUnreachable 带外IP不可达
	ErrOOBIPUnreachable = errors.New("oob ip is unreachable")
	// ErrMissingIPMImonitoring 缺失ipmimonitoring采集工具
	ErrMissingIPMImonitoring = errors.New("missing 'ipmimonitoring' tool")
	// ErrUsernamePassword 带外用户名、密码不匹配
	ErrUsernamePassword = errors.New("username and password do not match")
)

// IPMImonitoring IPMI监控工具
type IPMImonitoring struct {
	log      logger.Logger
	ip       string
	username string
	password string
}

// NewIPMImonitoring 实例化IPMI监控工具
func NewIPMImonitoring(log logger.Logger, ip, username, password string) (*IPMImonitoring, error) {
	if ip == "" || username == "" || password == "" {
		return nil, ErrMissingOOBInfo
	}
	return &IPMImonitoring{
		log:      log,
		ip:       ip,
		username: username,
		password: password,
	}, nil
}

// lookPath 在PATH环境变量下查找ipmimonitoring
func (im *IPMImonitoring) lookPath() (err error) {
	if _, err = exec.LookPath("ipmimonitoring"); err != nil {
		im.log.Warn("Tool 'ipmimonitoring' not installed correctly")
		return ErrMissingIPMImonitoring
	}
	return nil
}

// ping 检测网络是否可达
func (im *IPMImonitoring) ping() (err error) {
	if !ping.Ping(im.ip, 5) {
		im.log.Warnf("The OOB IP(%s) is unreachable.", im.ip)
		return ErrOOBIPUnreachable
	}
	return nil
}

// CollectSensorData 返回传感器数据
func (im *IPMImonitoring) CollectSensorData() (items []*model.SensorData, err error) {
	if err = im.lookPath(); err != nil {
		return nil, err
	}
	// -W authcap 各厂商实现IPMI协议存在不同，部分型号需要指定 workaround
	cmd := fmt.Sprintf("ipmimonitoring --interpret-oem-data --ignore-not-available-sensors --ignore-unrecognized-events --output-sensor-state --entity-sensor-names --comma-separated-output --no-header-output -D LAN_2_0 -W authcap -h %s -u %s -p %s",
		im.ip,
		im.username,
		im.password,
	)
	out, err := sh.ExecOutputWithLog(im.log, cmd)
	if err != nil {
		if pingerr := im.ping(); pingerr != nil {
			return nil, ErrOOBIPUnreachable
		}
		return nil, ErrUsernamePassword
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Split(line, ",")
		if len(fields) != 7 {
			continue
		}
		items = append(items, &model.SensorData{
			ID:      fields[0],
			Name:    fields[1],
			Type:    fields[2],
			State:   fields[3],
			Reading: fields[4],
			Units:   fields[5],
			Event:   strings.Replace(fields[6], "'", "", -1),
		})
	}
	return items, scanner.Err()
}
