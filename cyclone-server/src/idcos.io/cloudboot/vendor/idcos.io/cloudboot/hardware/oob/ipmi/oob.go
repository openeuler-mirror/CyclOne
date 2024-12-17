package ipmi

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"idcos.io/cloudboot/hardware"
	"idcos.io/cloudboot/hardware/oob"
	"idcos.io/cloudboot/logger"
)

const (
	// name 处理器名称
	name = "IPMI"
	// tool 硬件配置工具
	tool = "ipmitool"
)

func init() {
	oob.Register(name, new(worker))
}

type worker struct {
	hardware.Base
	mux sync.Mutex
}

// SetDebug 设置是否开启debug。若开启debug，会将关键日志信息写入console。
func (worker *worker) SetDebug(debug bool) {
	worker.Base.SetDebug(debug)
}

// SetLog 更换日志实现。默认情况下内部无日志实现。
func (worker *worker) SetLog(log logger.Logger) {
	worker.mux.Lock()
	defer worker.mux.Unlock()

	worker.Base.SetLog(log)
}

// SetDHCP 设置IP来源是DHCP
func (worker *worker) SetDHCP() error {
	channel, err := worker.Channel()
	if err != nil {
		return err
	}
	_, err = worker.Base.ExecByShell(tool, "lan", "set", strconv.Itoa(channel), "ipsrc", "dhcp")
	return err
}

// SetStaticIP 设置IP来源是静态IP
func (worker *worker) SetStaticIP(ip, netmask, gateway string) error {
	channel, err := worker.Channel()
	if err != nil {
		return err
	}
	ch := strconv.Itoa(channel)
	if _, err := worker.Base.ExecByShell(tool, "lan", "set", ch, "ipsrc", "static"); err != nil {
		return err
	}
	if _, err := worker.Base.ExecByShell(tool, "lan", "set", ch, "ipaddr", ip); err != nil {
		return err
	}
	if _, err := worker.Base.ExecByShell(tool, "lan", "set", ch, "netmask", netmask); err != nil {
		return err
	}
	if _, err := worker.Base.ExecByShell(tool, "lan", "set", ch, "defgw", "ipaddr", gateway); err != nil {
		return err
	}
	return nil
}

// GenerateUser 生成用户带外帐号
func (worker *worker) GenerateUser(sett *oob.UserSettingItem) error {
	var channel, userID int
	user, err := worker.findUserByName(sett.Username)
	if err != nil {
		return err
	}
	if user == nil { // 目标用户不存在
		channel, err = worker.Channel()
		if err != nil {
			return err
		}

		userID, err = worker.newUserID()
		if err != nil {
			return err
		}

		// ipmitool user set name $userid "$_user"
		if _, err = worker.Base.ExecByShell(tool, "user", "set", "name", strconv.Itoa(userID), fmt.Sprintf("%q", sett.Username)); err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)

	} else { // 目标用户已存在
		channel = user.Channel
		userID = user.ID
	}

	// ipmitool user set password $userid "$_pw"
	if _, err = worker.Base.ExecByShell(tool, "user", "set", "password", strconv.Itoa(userID), fmt.Sprintf("%q", sett.Password)); err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)

	// ipmitool user enable $userid
	if _, err = worker.Base.ExecByShell(tool, "user", "enable", strconv.Itoa(userID)); err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)

	// ipmitool user priv $userid $privilege_level $channel
	var output []byte
	if output, err = worker.Base.ExecByShell(tool, "user", "priv", strconv.Itoa(userID), strconv.Itoa(sett.PrivilegeLevel), strconv.Itoa(channel)); err != nil {
		err = errors.New(string(output))
		return err
	}
	time.Sleep(500 * time.Millisecond)

	// ipmitool channel setaccess $channel $userid callin=on ipmi=on link=on privilege=$privilege_level
	if _, err = worker.Base.ExecByShell(tool, "channel", "setaccess", strconv.Itoa(channel), strconv.Itoa(userID), "callin=on", "ipmi=on", "link=on", fmt.Sprintf("privilege=%d", sett.PrivilegeLevel)); err != nil {
		// TODO 暂时性处理办法：去掉'link=on'后重试
		if _, err = worker.Base.ExecByShell(tool, "channel", "setaccess", strconv.Itoa(channel), strconv.Itoa(userID), "callin=on", "ipmi=on", fmt.Sprintf("privilege=%d", sett.PrivilegeLevel)); err != nil {
			return err
		}
		return err
	}
	return nil
}

// Users 返回OOB用户列表
func (worker *worker) Users() ([]oob.User, error) {
	channel, err := worker.Channel()
	if err != nil {
		return nil, err
	}

	output, err := worker.Base.ExecByShell(tool, "user", "list", strconv.Itoa(channel))
	if err != nil {
		return nil, err
	}

	users, err := ParseUsers(output)
	if err != nil {
		return nil, err
	}

	for i := range users {
		users[i].Channel = channel
		users[i].Access, _ = worker.userAccess(channel, users[i].ID)
		if users[i].Access != nil {
			users[i].Name = users[i].Access.UserName // 解析ipmitool user list $channel的输出可能无法得到正确的用户名
		}
	}
	return users, nil
}

// BMCColdReset (冷)重启BMC
func (worker *worker) BMCColdReset() error {
	_, _ = worker.Base.ExecByShell(tool, "mc", "reset", "cold")
	return nil // 假设每次执行'ipmitool mc reset cold'都能达到预期效果，丢弃error。
}

// channel 返回channel。
func (worker *worker) Channel() (int, error) {
	// 试错法查找channel
	for i := 0; i <= 10; i++ {
		out, _ := worker.Base.ExecByShell(tool, "lan", "print", strconv.Itoa(i))
		if len(out) > 0 && !strings.Contains(string(out), "Invalid channel") { // 通过exit code判断并不能保证一定准确
			return i, nil
		}
	}
	return 0, oob.ErrChannelNotFound
}

// BMC 返回OOB的BMC信息
func (worker *worker) BMC() (*oob.BMC, error) {
	output, err := worker.Base.ExecByShell(tool, "mc", "info")
	if err != nil {
		return nil, err
	}

	var bmc oob.BMC
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "Firmware Revision") {
			bmc.FirmwareReversion = worker.extractValue(line, ":")
		} else if strings.HasPrefix(line, "IPMI Version") {
			bmc.IPMIVersion = worker.extractValue(line, ":")
		} else if strings.HasPrefix(line, "Vendor ID") {
			bmc.ManufacturerID = worker.extractValue(line, ":")
		} else if strings.HasPrefix(line, "Vendor Name") {
			bmc.ManufacturerName = worker.extractValue(line, ":")
		}
	}
	return &bmc, nil
}

// Network 返回OOB网络信息
func (worker *worker) Network() (*oob.Network, error) {
	channel, err := worker.Channel()
	if err != nil {
		return nil, err
	}

	output, _ := worker.Base.ExecByShell(tool, "lan", "print", strconv.Itoa(channel)) // 舍弃error。因为部分机型下执行该命令，即使命令输出正确，exit code也会是一个非0值。

	var network oob.Network
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "IP Address Source") {
			network.IPSrc = worker.extractValue(line, ":")
		} else if strings.HasPrefix(line, "IP Address") {
			network.IP = worker.extractValue(line, ":")
		} else if strings.HasPrefix(line, "Subnet Mask") {
			network.Netmask = worker.extractValue(line, ":")
		} else if strings.HasPrefix(line, "MAC Address") {
			network.Mac = worker.extractValue(line, ":")
		} else if strings.HasPrefix(line, "Default Gateway IP") {
			network.Gateway = worker.extractValue(line, ":")
		}
	}
	return &network, nil
}

// PostCheck OOB配置实施后置检查
func (worker *worker) PostCheck(sett *oob.Setting) (items []hardware.CheckingItem) {
	if sett == nil {
		return nil
	}
	if sett.Network != nil {
		items = append(items, worker.checkNetwork(sett.Network)...)
	}
	if sett.User != nil {
		items = append(items, worker.checkUser(sett.User)...)
	}
	return items
}
