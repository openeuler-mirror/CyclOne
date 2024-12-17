package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"math/rand"
	"time"

	"sync"

	"github.com/jinzhu/gorm"
	"github.com/voidint/binding"
	"idcos.io/cloudboot/config"
	"idcos.io/cloudboot/hardware/oob/ipmi"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils"
	"idcos.io/cloudboot/utils/oob"
	"idcos.io/cloudboot/utils/sh"
)

// OOBPowerBatchOperateReq 带外管理批量请求结构体
type OOBPowerBatchOperateReq struct {
	Sns []string `json:"sns"`
}

// FieldMap 请求字段映射
func (reqData *OOBPowerBatchOperateReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.Sns: "sns",
	}
}

const (
	// PowerRestart 重启
	PowerRestart = "restart"
	// PowerOn 开机
	PowerOn = "on"
	// PowerOff 关机
	PowerOff = "off"
	// OperatePowerOn 服务器开机
	OperatePowerOn = "ipmitool -I lanplus -H $HOST -U $USER -P $PASSWD power on"
	// OperatePowerOff 服务器关机
	OperatePowerOff = "ipmitool -I lanplus -H $HOST -U $USER -P $PASSWD power off"
	// OperatePowerRestart 服务器重启
	OperatePowerRestart = "ipmitool -I lanplus -H $HOST -U $USER -P $PASSWD power reset"
	// OperatePowerStatus 服务器状态
	OperatePowerStatus = "ipmitool -I lanplus -H $HOST -U $USER -P $PASSWD power status"
	// OperatePowerPXE PXE启动 不能指定特定端口,端口指定可借助racadm等工具实现，网卡的pxe功能需要在bios中开启
	OperatePowerPXE = "ipmitool -I lanplus -H $HOST -U $USER -P $PASSWD chassis bootdev pxe  options=efiboot"
	//CmdSleep2s 命令之间睡2秒
	CmdSleep2s = "sleep 2s"

	// Domain 带外注册在DNS上的域名, 与hostname(SN)一起充当HOST的角色
	//Domain = "oob.webank.com"

	// OOBVendorConfig 带外管理厂商 用户名密码配置
	OOBVendorConfig = `
{
  "Dell": [
    {
      "password": "admin",
      "username": "albert"
    },
    {
      "password": "calvin",
      "username": "root"
    }
  ],
  "戴尔": [
    {
      "password": "admin",
      "username": "albert"
    },
    {
      "password": "calvin",
      "username": "root"
    }
  ],
  "HP": [
    {
      "password": "12345678",
      "username": "Administrator"
    }
  ],
  "惠普": [
    {
      "password": "12345678",
      "username": "Administrator"
    }
  ],  
  "Huawei": [
    {
      "password": "Huawei12#$",
      "username": "root"
	},
    {
		"password": "Admin@9000",
		"username": "Administrator"
	}
  ],
  "华为": [
    {
      "password": "Huawei12#$",
      "username": "root"
	},
    {
		"password": "Admin@9000",
		"username": "Administrator"
	}
  ],  
  "default": [
    {
      "password": "admin",
      "username": "albert"
    },
    {
      "password": "root",
      "username": "root"
	},
    {
		"password": "admin",
		"username": "root"
	},
    {
		"password": "jumpm3",
		"username": "albert"
	}	
  ],
  "ibm": [
    {
      "password": "PASSW0RD",
      "username": "USERID"
    }
  ]
}`
)

// OOBUser 带外管理用户结构体
type OOBUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// OOBInfo 带外信息
type OOBInfo struct {
	Username string `json:"username"`
	IP       string `json:"ip"`
	Password string `json:"password"`
}

// OOBOperate 带外管理操作结构体
type OOBOperate struct {
	HOST   string `json:"host"`
	User   string `json:"user"`
	Passwd string `json:"passwd"`
}

// UpdateOOBPasswordReq 修改带外密码请求体
type UpdateOOBPasswordReq struct {
	//设备序列号
	SN string `json:"sn"`
	//带外用户名
	Username string `json:"oob_user_name"`
	//带外原密码，置空时，新密码必须是正确的当前密码，否则修改会失败
	PasswordOld string `json:"oob_password_old"`
	//带外新密码，置空时且旧密码有效，则修改为一个随机新密码
	PasswordNew string `json:"oob_password_new"`
	//备注信息
	Remark string
}

// OOBHistoryReq 带外用户密码修改历史请求参数结构
type OOBHistoryReq struct {
	//设备序列号
	SN string
	//带外原用户名
	UsernameOld string
	//带外新用户名
	UsernameNew string
	//带外原密码
	PasswordOld string
	//带外新密码
	PasswordNew string
	//记录创建者ID
	Creator string
	//备注
	Remark string
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *UpdateOOBPasswordReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.Username:    "oob_user_name",
		&reqData.PasswordOld: "oob_password_old",
		&reqData.PasswordNew: "oob_password_new",
	}
}

// GetOOBInfoBySn 根据sn获取设备带外信息
func GetOOBInfoBySn(log logger.Logger, repo model.Repo, conf *config.Config, sn string) (user *OOBInfo, err error) {
	user = &OOBInfo{}
	device, err := repo.GetDeviceBySN(sn)
	if err != nil || device == nil {
		return user, err
	}
	user.IP = oob.TransferHostname2IP(log, repo, sn, utils.GetOOBHost(sn, device.Vendor, conf.Server.OOBDomain))

	// 查看用户及密码，从device_oob_history这张表找，可以找到已经从平台删除的设备的信息
	u, err := FindOOBByHistory(log, repo, conf, sn)
	if err != nil && err != gorm.ErrRecordNotFound {
		return user, err
	}
	if u == nil {
		//return nil, fmt.Errorf("带外用户名或密码为空，请设置")
		log.Warnf("device SN: %s oob user or password empty", sn)
	} else {
		user.Username = u.Username
		user.Password = u.Password

		//查看device表里面的用户名和密码是否存在，若不存在，则同步一下
		device, _ := repo.GetDeviceBySN(sn)
		if device.OOBIP == "" || device.OOBIP != user.IP {
			device.OOBIP = user.IP
		}
		if device.OOBUser == "" || device.OOBPassword == "" {
			device.OOBUser = user.Username
			encryptedPassword, err := utils.AESEncrypt(user.Password, []byte(conf.Crypto.Key))
			if err != nil {
				log.Errorf("SN: %s sync oob info fail, encrypt new password fail", sn)
				return nil, errors.New("新密码加密失败")
			}
			device.OOBPassword = encryptedPassword
		}
		repo.UpdateDevice(device)
	}
	return user, nil //忽略查询错误
}

// GetOOBUserBySN 根据sn设备带外账户信息
func GetOOBUserBySN(log logger.Logger, repo model.Repo, sn, key string) (user *OOBUser, err error) {
	m, err := repo.GetDeviceBySN(sn)
	if err != nil {
		return nil, err
	}
	user = &OOBUser{}
	user.Username = m.OOBUser
	pwd, err := utils.AESDecrypt(m.OOBPassword, []byte(key))
	if err != nil {
		return nil, err
	}
	user.Password = string(pwd)
	return user, err
}

// BatchOperateOOBPower 带外管理批量操作
func BatchOperateOOBPower(log logger.Logger, repo model.Repo, operate string, conf *config.Config, isPxe bool, sns []string) (output string, err error) {
	var errMsg []string
	var outMsg []string

	for _, sn := range sns {
		out, err := OperateOOBPower(log, repo, sn, conf.Crypto.Key, conf.Server.OOBDomain, operate, isPxe)

		if err != nil {
			log.Debugf("%s device operate failure, err: %s", sn, err.Error())
			errMsg = append(errMsg, err.Error())
			continue
		}

		log.Debugf("%s device operate success, output: %s", sn, out)
		outMsg = append(outMsg, fmt.Sprintf("设备(%s)操作成功", sn))
	}
	if len(errMsg) > 0 {
		return "", fmt.Errorf(strings.Join(errMsg, "\n"))
	}

	return strings.Join(outMsg, "\n"), nil
}

// OperateOOBPower 带外管理的操作
func OperateOOBPower(log logger.Logger, repo model.Repo, sn, key, oobDomain, operate string, isPxe bool) (output string, err error) {
	m, err := repo.GetDeviceBySN(sn)
	if err != nil {
		return "", err
	}

	if m.OOBUser == "" || m.OOBPassword == "" {
		log.Warnf("设备带外用户或密码为空，尝试找回，[SN:%s]", m.SN)
		history, err := repo.GetLastOOBHistoryBySN(sn)
		if err != nil {
			log.Errorf("find back oob history by sn:%s fail,%s", sn, err.Error())
			return fmt.Sprintf("find back oob history by sn:%s fail", sn), fmt.Errorf("设备用户名密码为空，无法操作")
		} else {
			m.OOBUser = history.UsernameNew
			m.OOBPassword = history.PasswordNew
		}
	}

	oobHost := utils.GetOOBHost(m.SN, m.Vendor, oobDomain)
	oobIP := oob.TransferHostname2IP(log, repo, m.SN, oobHost)
	if oobIP == "" {
		return "操作失败", errors.New("未获取到带外IP")
	}
	oobUser := m.OOBUser
	oobPassword, err := utils.AESDecrypt(m.OOBPassword, []byte(key))
	if err != nil {
		log.Debugf("descrypt password failure, err: %s", err.Error())
		return "", err
	}
	// check is power
	isPowerOn, err := OOBPowerStatus(log, oobIP, oobUser, string(oobPassword), m.OOBPassword)
	if err != nil {
		return "", err
	}

	var cmd string
	switch operate {
	case PowerOff:
		cmd += OperatePowerOff
	case PowerOn:
		if isPxe {
			cmd = OperatePowerPXE + " && " + CmdSleep2s + " && "
		}
		cmd += OperatePowerOn
	case PowerRestart:
		// 设备在关闭状态下无法重启
		if !isPowerOn {
			return "", fmt.Errorf("设备为关机状态，无法重启")
		}

		if isPxe {
			cmd = OperatePowerPXE + " && " + CmdSleep2s + " && "
		}
		cmd += OperatePowerRestart
	}

	cmd = replaceCmd(cmd, oobIP, oobUser, string(oobPassword))

	DesensitizePasswordLog(log, fmt.Sprintf("start to exec cmd: [%s]", cmd), string(oobPassword), m.OOBPassword)

	out, err := sh.ExecOutputWithLog(log, cmd)
	if err != nil {
		log.Debugf("exec [%s] done , err: [%s], stdout: [%s]", sh.CmdDesensitization(cmd), err.Error(), sh.CmdDesensitization(string(out)))
		return "", ProcessStdoutMsg(out, oobIP, oobUser)
	}

	DesensitizePasswordLog(log, fmt.Sprintf("exec [%s] done，output: [%s]", cmd, string(out)), string(oobPassword), m.OOBPassword)

	// update device power_status
	status := model.PowerStatusOn

	if operate == PowerOff {
		status = model.PowerStatusOff
	}
	m.PowerStatus = status
	repo.UpdateDevice(m)

	return string(out), nil
}

// ProcessStdoutMsg 处理错误信息
func ProcessStdoutMsg(output []byte, host, username string) error {
	str := string(output)

	if strings.Contains(str, "Could not open socket!") {
		return fmt.Errorf("无法连接目标主机(%s)", host)
	}

	if strings.Contains(str, "command not found") {
		return fmt.Errorf("命令不存在")
	}

	if strings.Contains(str, "Unable to establish IPMI v2 / RMCP+ session") {
		return fmt.Errorf("用户名(%s)或密码错误", username)
	}
	return fmt.Errorf("其他错误:%s", string(output))
}

// DesensitizePasswordLog 脱敏密码输出
func DesensitizePasswordLog(log logger.Logger, str, oldPass, newPass string) {
	log.Debugf("%s", strings.Replace(str, oldPass, newPass, -1))
}

// OOBPowerStatus 检查OOBPower状态
func OOBPowerStatus(log logger.Logger, oobHost, oobUser, oobPassword, oldPassword string) (bool, error) {
	if oobHost == "" {
		return false, errors.New("未获取到带外IP")
	}
	cmd := replaceCmd(OperatePowerStatus, oobHost, oobUser, oobPassword)

	output, err := sh.ExecOutputWithLog(log, cmd)
	if err != nil {
		return false, ProcessStdoutMsg(output, oobHost, oobUser)
	}

	return strings.Contains(string(output), "Chassis Power is on"), nil
}

// replaceCmd 将HOST、USER、PASSWD进行替换
func replaceCmd(cmd, oobHost, oobUser, oobPassword string) string {
	cmd = strings.Replace(cmd, "$HOST", oobHost, -1)
	cmd = strings.Replace(cmd, "$USER", oobUser, -1)
	cmd = strings.Replace(cmd, "$PASSWD", oobPassword, -1)
	return cmd
}

// UpdateOOBPasswordBySN 更新带外密码
// 旧密码为空，新密码为用户指定的正确密码，直接修改库
// 旧密码有值，且正确，新密码修改为用户指定值，或随机值
// 旧密码有值，且不正确，修改失败
func UpdateOOBPasswordBySN(log logger.Logger, repo model.Repo, reqData *UpdateOOBPasswordReq,
	conf *config.Config) (user *OOBUser, err error) {

	_, err = utils.AESEncrypt(reqData.PasswordOld, []byte(conf.Crypto.Key))
	if err != nil {
		return nil, err
	}

	user = &OOBUser{}
	//直接用新密码测试连通性，测试通过直接改表
	if reqData.PasswordOld == "" {
		return nil, errors.New("[数据校验]未指定旧密码，更改失败")
	}

	dev, err := repo.GetDeviceBySN(reqData.SN)
	if err != nil {
		log.Errorf("SN:%s update new oob password fail，Err:%s", dev.SN, err.Error())
		return
	}

	//获取用户名对应的userID
	cmd := fmt.Sprintf("ipmitool -I lanplus -H %s -U %s -P %s user list",
		oob.TransferHostname2IP(log, repo, reqData.SN, utils.GetOOBHost(reqData.SN, dev.Vendor, conf.Server.OOBDomain)),
		reqData.Username, reqData.PasswordOld)
	out, err := sh.ExecOutputWithLog(log, cmd)
	if err != nil || strings.Contains(string(out), "Error: Unable to establish LAN session") {
		log.Errorf("SN: %s change oob password fail, old oob password is wrong")
		return nil, errors.New("[数据校验]旧密码不正确，更改失败")
	}
	users, err := ipmi.ParseUsers(out)
	if err != nil {
		log.Errorf("SN: %s change oob password fail, parse oob user fail", reqData.SN)
		return nil, errors.New("解析带外用户列表失败")
	}
	userID := 0
	for _, u := range users {
		if reqData.Username == u.Name {
			userID = u.ID
			break
		}
	}
	//如果新密码为空，则自动生成一个随机密码
	if reqData.PasswordNew == "" {
		reqData.PasswordNew = GenPassword()
	}
	cmd = fmt.Sprintf("ipmitool -I lanplus -H %s -U %s -P %s user set password %d %s",
		oob.TransferHostname2IP(log, repo, reqData.SN, utils.GetOOBHost(reqData.SN, dev.Vendor, conf.Server.OOBDomain)),
		reqData.Username,
		reqData.PasswordOld, userID, reqData.PasswordNew)
	out, err = sh.ExecOutputWithLog(log, cmd)
	if err != nil || strings.Contains(string(out), "Error: Unable to establish LAN session") {
		log.Errorf("SN: %s change oob password fail, %s", reqData.SN, err.Error())
		return nil, err
	}

	//更改成功,修改数据库
	user.Username = reqData.Username
	user.Password = reqData.PasswordNew
	encryptedPassword, err := utils.AESEncrypt(user.Password, []byte(conf.Crypto.Key))
	if err != nil {
		log.Errorf("SN: %s change oob password fail, encrypt new password fail", reqData.SN)
		return nil, errors.New("新密码加密失败")
	}

	dev.OOBUser = user.Username
	dev.OOBPassword = encryptedPassword

	yes := model.YES
	dev.OOBAccessible = &yes
	_, err = repo.UpdateDevice(dev)
	if err != nil {
		log.Errorf("SN:%s update new oob password fail，Err:%s", dev.SN, err.Error())
		return
	}

	//将密码修改记录写到device_oob_history这张表-->改成Mysql触发器
	//AddOOBHistory(log, repo, conf, &OOBHistoryReq{
	//	SN:          reqData.SN,
	//	UsernameOld: dev.OOBUser,
	//	UsernameNew: reqData.Username,
	//	PasswordOld: reqData.PasswordOld,
	//	PasswordNew: reqData.PasswordNew,
	//	//Creator:     "",
	//	Remark: reqData.Remark,
	//})

	return
}

//GenPassword 生成8位密码
func GenPassword() string {
	//生成规则 8位数字字母随机组合
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	ascii := "~@#%^" //!$这几个符号在shell中有特殊含义，别引入了
	length := len(str)
	result := make([]byte, 0, 8)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 8; i++ {
		if i == 4 {
			result = append(result, ascii[r.Intn(5)])
		} else {
			result = append(result, str[r.Intn(length)]) //52英文字母，10数字
		}
	}
	return string(result)
}

//GetDefaultOOBUserPassword 根据机器厂商查询默认出厂带外用户及密码
//IBM=USERID:PASSW0RD
//HP=Administrator:12345678
//DELL=albert:admin\root:calvin
//HUAWEI=root:Huawei12#$
//OTHER=albert:admin\root:root
//返回规则改成：如果连通测试不通过，则返回nil
func GetDefaultOOBUserPassword(log logger.Logger, repo model.Repo, config *config.Config, sn, vendor string) (user *OOBUser, err error) {
	vendor = strings.ToLower(vendor)

	var mapper map[string][]*OOBUser
	err = json.Unmarshal([]byte(OOBVendorConfig), &mapper)

	oobIP := oob.TransferHostname2IP(log, repo, sn, utils.GetOOBHost(sn, vendor, config.Server.OOBDomain))
	for k, oobusers := range mapper {
		if strings.ToLower(k) != strings.ToLower(vendor) {
			continue
		}

		for _, oobUser := range oobusers {
			if oobPingTest(log, oobIP, oobUser.Username, oobUser.Password, config.Crypto.Key) {
				user = oobUser
				break
			}
		}
		return
	}

	for _, oobUser := range mapper["default"] {
		if oobPingTest(log, oobIP, oobUser.Username, oobUser.Password, config.Crypto.Key) {
			user = oobUser
			break
		}
	}
	return
}

// oobPingTest 测试带外是否可以正常使用
func oobPingTest(log logger.Logger, host, username, password, oldPassword string) bool {
	cmd := fmt.Sprintf("ipmitool -I lanplus -H %s -U %s -P %s user list", host, username, password)
	out, err := sh.ExecDesensitizeOutputWithLog(log, cmd, password, oldPassword)
	if err != nil || strings.Contains(string(out), "Error: Unable to establish LAN session") {
		return false
	}
	return true
}

// AddOOBHistory 增加带外修改历史记录
//func AddOOBHistory(log logger.Logger, repo model.Repo, conf *config.Config, oobHistory *OOBHistoryReq) (affected int64, err error) {
//	mod := model.OOBHistory{
//		SN:          oobHistory.SN,
//		UsernameOld: oobHistory.UsernameOld,
//		UsernameNew: oobHistory.UsernameOld,
//		//PasswordOld:oobHistory.PasswordOld,
//		//PasswordNew:oobHistory.PasswordNew,
//		Remark:  oobHistory.Remark,
//		Creator: oobHistory.Creator,
//	}
//	if mod.PasswordOld, err = utils.AESEncrypt(oobHistory.PasswordOld, []byte(conf.Crypto.Key)); err != nil {
//		log.Errorf("encrypt old oob password fail, %s", err.Error())
//		return 0, err
//	}
//	if mod.PasswordNew, err = utils.AESEncrypt(oobHistory.PasswordNew, []byte(conf.Crypto.Key)); err != nil {
//		log.Errorf("encrypt new oob password fail, %s", err.Error())
//		return 0, err
//	}
//
//	return repo.AddOOBHistory(&mod)
//}

// FindOOBByHistory 通过历史记录找回设备当前带外用户密码
func FindOOBByHistory(log logger.Logger, repo model.Repo, conf *config.Config, sn string) (oob *OOBUser, err error) {
	oob = new(OOBUser)
	history, err := repo.GetLastOOBHistoryBySN(sn)
	if err != nil {
		log.Errorf("find back oob history by sn:%s fail,%s", sn, err.Error())
		return nil, err
	}
	oob.Username = history.UsernameNew
	password, err := utils.AESDecrypt(history.PasswordNew, []byte(conf.Crypto.Key))
	if err != nil {
		log.Errorf("decrypt oob new password:%s fail,%s", history.PasswordNew, err.Error())
		return nil, err
	}
	oob.Password = string(password)
	return
}

type ReAccessOOBReq struct {
	SNs []string `json:"sns"`
}

// FieldMap 请求字段映射
func (reqData *ReAccessOOBReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.SNs: "sns",
	}
}

func ReAccessOOB(log logger.Logger, repo model.Repo, conf *config.Config, reqData *ReAccessOOBReq) (err error) {
	//找出所有的未纳管的机器
	isYes := model.NO
	devs := make([]*model.Device, 0)
	if len(reqData.SNs) != 0 {
		for _, sn := range reqData.SNs {
			d, err := repo.GetDeviceBySN(sn)
			if err != nil || d == nil {
				log.Errorf("sn:%s not exist", err)
				continue
			} else {
				devs = append(devs, d)
			}
		}
	} else {
		devs, err = repo.GetDevices(&model.Device{OOBAccessible: &isYes}, nil, nil)
		if err != nil {
			log.Errorf("get unaccessible devices fail,%v", err)
			return err
		}
	}
	var wg sync.WaitGroup

	sem := make(chan struct{}, 50) // 最多允许50个并发同时执行

	for i := range devs {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()

			sem <- struct{}{}        // 获取信号
			defer func() { <-sem }() // 释放信号

			dev := devs[j]
			oobIP := oob.TransferHostname2IP(log, repo, dev.SN, utils.GetOOBHost(dev.SN, dev.Vendor, conf.Server.OOBDomain))
			//如果IP都没解析到，就不要再往下试用户密码了
			if oobIP == "" {
				log.Warnf("SN: %s get oob ip from dns fail", dev.SN)
				return
			}
			if dev.OOBPassword != "" {
				//设备表中的密码需要解密后使用
				oldOOBPasswordByte, err := utils.AESDecrypt(dev.OOBPassword, []byte(conf.Crypto.Key))
				if err != nil {
					log.Warnf("SN: %s decrypt old password：%s fail", dev.OOBPassword)
				} else if oobPingTest(log, oobIP, dev.OOBUser, string(oldOOBPasswordByte), dev.OOBPassword) {
					isYesLocal := model.YES
					dev.OOBAccessible = &isYesLocal
					if _, err = repo.UpdateDevice(dev); err != nil {
						log.Errorf("update oob accessible status fail,SN:%s,err:%v", dev.SN, err)
					}
					return //success 1
				}
			}

			// 2 获取一个出厂默认的用户密码
			defaultUser, err := GetDefaultOOBUserPassword(log, repo, conf, dev.SN, dev.Vendor)
			if defaultUser != nil {
				UpdateOOBPasswordBySN(log, repo, &UpdateOOBPasswordReq{
					SN:          dev.SN,
					Username:    defaultUser.Username,
					PasswordOld: defaultUser.Password,
					PasswordNew: GenPassword(),
				}, conf)
				return //success 2
			}
			log.Infof("SN：%s get default oob user fail,err:%v", dev.SN, err)

			// 3 尝试oob导入账户
			ouInit := OOBUser{}
			_ = json.Unmarshal([]byte(dev.OOBInit), &ouInit)
			if ouInit.Password != "" {
				encryptedPassword, err := utils.AESEncrypt(ouInit.Password, []byte(conf.Crypto.Key))
				if err != nil {
					log.Errorf("SN: %s encrypt imported oob password fail", dev.SN)
					return
				}
				if oobPingTest(log, oobIP, ouInit.Username, ouInit.Password, encryptedPassword) {
					isYesLocal := model.YES
					dev.OOBAccessible = &isYesLocal
					dev.OOBUser = ouInit.Username
					dev.OOBPassword = encryptedPassword
					if _, err = repo.UpdateDevice(dev); err != nil {
						log.Errorf("update oob accessible status fail,SN:%s,err:%v", dev.SN, err)
					}
					return //success 3
				}
			}
		}(i)
	}
	wg.Wait()

	return nil
}


// OOBInspectionOperateReq 带外巡检请求结构体
type OOBInspectionOperateReq struct {
	SN 				string 		`json:"sn"` // 序列号
	IP 				string 		`json:"ip"` // 内网IP
	DataType	 	string 		`json:"data_type"`  //校验枚举值 [sensor sel all]
}

type OOBInspectionOperateResp struct {
	DataType	 	string 					`json:"data_type"`  //校验枚举值 [sensor sel all]
	SensorData		[]*model.SensorData		`json:"sensor_data"` 
	SelData			[]*model.SelData		`json:"sel_data"` 
}
// FieldMap 请求字段映射
func (reqData *OOBInspectionOperateReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.SN: "sn",
		&reqData.IP: "ip",
		&reqData.DataType: "data_type",
	}
}

// OOBInspectionOperate 带外巡检操作[ipmimonitoring|ipmi-sel]
func OOBInspectionOperate(log logger.Logger, repo model.Repo, conf *config.Config, reqData *OOBInspectionOperateReq) (*OOBInspectionOperateResp, error) {
	// 优先根据SN巡检，其次根据内网IP
	var sn string
	if reqData.SN != "" {
		sn = reqData.SN
	} else {
		if reqData.IP != "" {
			cond := model.IPPageCond{
				IP:               reqData.IP,
			}
			items, err := repo.GetIPs(&cond, model.OneOrderBy("id", model.DESC), nil)
			if err != nil {
				return nil, err
			}
			for i := range items {
				sn = items[i].SN
			}
		}
	}
	
	m, err := repo.GetDeviceBySN(sn)
	if err != nil {
		return nil, err
	}

	if m.OOBUser == "" || m.OOBPassword == "" {
		log.Warnf("设备带外用户或密码为空，尝试找回，[SN:%s]", m.SN)
		history, err := repo.GetLastOOBHistoryBySN(sn)
		if err != nil {
			log.Errorf("find back oob history by sn:%s fail,%s", sn, err.Error())
			return nil, fmt.Errorf("设备用户名密码为空，无法操作")
		} else {
			m.OOBUser = history.UsernameNew
			m.OOBPassword = history.PasswordNew
		}
	}

	oobHost := utils.GetOOBHost(m.SN, m.Vendor, conf.Server.OOBDomain)
	oobIP := oob.TransferHostname2IP(log, repo, m.SN, oobHost)
	if oobIP == "" {
		return nil, errors.New("未获取到带外IP")
	}
	oobUser := m.OOBUser
	oobPassword, err := utils.AESDecrypt(m.OOBPassword, []byte(conf.Crypto.Key))
	if err != nil {
		log.Debugf("descrypt password failure, err: %s", err.Error())
		return nil, err
	}
	// 仅对开电状态的机器发起巡检
	isPowerOn, err := OOBPowerStatus(log, oobIP, oobUser, string(oobPassword), m.OOBPassword)
	if err != nil {
		return nil, err
	}
	if !isPowerOn {
		return nil, fmt.Errorf("设备为关电状态，无法巡检")
	}
	
	switch reqData.DataType {
	case "sensor":
		var result = &OOBInspectionOperateResp{
			DataType:	"sensor",
			SensorData: nil,
			SelData:	nil,
		}
		// 调用ipmimonitoring工具采集传感器数据
		ipmiMonitor, err := oob.NewIPMImonitoring(log, oobIP, oobUser, string(oobPassword))
		if err != nil {
			return nil, err
		}
		sensors, err := ipmiMonitor.CollectSensorData()
		if err != nil {
			log.Errorf("IPMI-Sensor data collection failed: %s", err.Error())
			return nil, err
		}
		if sensors != nil {
			result.SensorData = sensors
			return result, nil
		}
		return nil, fmt.Errorf("巡检数据（%s）为空", reqData.DataType)
	case "sel":
		var result = &OOBInspectionOperateResp{
			DataType:	"sel",
			SensorData: nil,
			SelData:	nil,
		}
		// 调用ipmi-sel工具采集系统事件日志
		ipmiSel, err := oob.NewIPMIsel(log, oobIP, oobUser, string(oobPassword))
		if err != nil {
			return nil, err
		}
		sel, err := ipmiSel.CollectSelData()
		if err != nil {
			log.Errorf("IPMI-Sel data collection failed: %s", err.Error())
			return nil, err
		}
		if sel != nil {
			result.SelData = sel
			return result, nil
		}
		return nil, fmt.Errorf("巡检数据（%s）为空", reqData.DataType)
	case "all":
		var result = &OOBInspectionOperateResp{
			DataType:	"all",
			SensorData: nil,
			SelData:	nil,
		}
		// 调用ipmimonitoring工具采集传感器数据
		ipmiMonitor, err := oob.NewIPMImonitoring(log, oobIP, oobUser, string(oobPassword))
		if err != nil {
			return nil, err
		}
		sensors, err := ipmiMonitor.CollectSensorData()
		if err != nil {
			log.Errorf("IPMI-Sensor data collection failed: %s", err.Error())
			return nil, err
		}
		if sensors != nil {
			result.SensorData = sensors
		} else {
			log.Warnf("IPMI-Sensor data is empty")
		}
		// 调用ipmi-sel工具采集系统事件日志
		ipmiSel, err := oob.NewIPMIsel(log, oobIP, oobUser, string(oobPassword))
		if err != nil {
			return result, err
		}
		sel, err := ipmiSel.CollectSelData()
		if err != nil {
			log.Errorf("IPMI-Sel data collection failed: %s", err.Error())
			return result, err
		}
		if sel != nil {
			result.SelData = sel
			return result, nil
		} else {
			log.Warnf("IPMI-Sel data is empty")
		}
	}
	return nil, fmt.Errorf("巡检数据（%s）为空", reqData.DataType)
}