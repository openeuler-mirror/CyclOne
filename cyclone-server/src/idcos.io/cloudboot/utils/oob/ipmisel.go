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
	//ErrMissingOOBInfo = errors.New("missing oob information")
	// ErrOOBIPUnreachable 带外IP不可达
	//ErrOOBIPUnreachable = errors.New("oob ip is unreachable")
	// ErrMissingIPMIsel 缺失IPMIsel采集工具
	ErrMissingIPMIsel = errors.New("missing 'ipmi-sel' tool")
	// ErrUsernamePassword 带外用户名、密码不匹配
	//ErrUsernamePassword = errors.New("username and password do not match")
)

// IPMIsel IPMI监控工具
type IPMIsel struct {
	log      logger.Logger
	ip       string
	username string
	password string
}

// NewIPMIsel 实例化IPMI监控工具
func NewIPMIsel(log logger.Logger, ip, username, password string) (*IPMIsel, error) {
	if ip == "" || username == "" || password == "" {
		return nil, ErrMissingOOBInfo
	}
	return &IPMIsel{
		log:      log,
		ip:       ip,
		username: username,
		password: password,
	}, nil
}

// lookPath 在PATH环境变量下查找IPMIsel
func (im *IPMIsel) lookPath() (err error) {
	if _, err = exec.LookPath("ipmi-sel"); err != nil {
		im.log.Warn("Tool 'ipmi-sel' not installed correctly")
		return ErrMissingIPMIsel
	}
	return nil
}

// ping 检测网络是否可达
func (im *IPMIsel) ping() (err error) {
	if !ping.Ping(im.ip, 5) {
		im.log.Warnf("The OOB IP(%s) is unreachable.", im.ip)
		return ErrOOBIPUnreachable
	}
	return nil
}

// CollectSelData 返回系统事件日志
func (im *IPMIsel) CollectSelData() (items []*model.SelData, err error) {
	if err = im.lookPath(); err != nil {
		return nil, err
	}
	// -W authcap 各厂商实现IPMI协议存在不同，部分型号需要指定 workaround
	cmd := fmt.Sprintf("ipmi-sel --tail=50 --output-event-state --non-abbreviated-units --interpret-oem-data --system-event-only --comma-separated-output --no-header-output -D LAN_2_0 -W authcap -h %s -u %s -p %s",
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
		items = append(items, &model.SelData{
			ID:      fields[0],
			Date:    fields[1],
			Time:    fields[2],
			Name:    fields[3],
			Type:    fields[4],
			State:   fields[5],
			Event:   fields[6],
		})
	}
	return items, scanner.Err()
}
