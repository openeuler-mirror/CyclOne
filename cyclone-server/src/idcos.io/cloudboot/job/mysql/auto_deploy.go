package mysql

import (
	"idcos.io/cloudboot/utils/times"
	"reflect"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/astaxie/beego/httplib"
	"github.com/jinzhu/gorm"

	"idcos.io/cloudboot/limiter"
	"idcos.io/cloudboot/config"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils"
	"idcos.io/cloudboot/utils/oob"
	"idcos.io/cloudboot/utils/sh"	
)

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
)

// AutoDeployJob 自动部署任务
type AutoDeployJob struct {
	id   string // 任务ID(全局唯一)
	log  logger.Logger
	repo model.Repo
	conf *config.Config
}

// NewAutoDeployJob 实例化任务管理器
func NewAutoDeployJob(log logger.Logger, repo model.Repo, conf *config.Config, jobid string) *AutoDeployJob {
	return &AutoDeployJob{
		log:  log,
		repo: repo,
		conf: conf,
		id:   jobid,
	}
}

// Run 任务
// 内置任务，根据装机参数规则引擎生成装机参数，自动对待部署设备发起部署（仅首次，存在失败部署记录则不处理）
func (j *AutoDeployJob) Run() {
	j.log.Debugf("Start auto deploy job")
	defer j.log.Debugf("Auto deploy job is completed")

	defer func() {
		if err := recover(); err != nil {
			j.log.Errorf("Auto deploy job panic: \n%s", err)
		}
	}()

	status := model.DevOperStatPreDeploy
	items, err := j.repo.GetDeviceByOperationStatus(status)
	if err != nil {
		j.log.Errorf("Fail to get predepoy devices: \n%s", err)
	}

	if len(items) <= 0 {
		j.log.Infof("The predepoy devices list is empty")
		return
	}
	err = j.AutoGenDeviceSetting(items)
	if err != nil {
		j.log.Errorf("Auto deploy failed: %s .", err)
	}
}

// AutoGenDeviceSetting 调用规则引擎自动生成设备的装机参数，并发起部署
func (j *AutoDeployJob) AutoGenDeviceSetting(devices []*model.Device) (err error) {
	j.log.Infof("Start auto depoy job generating device settings")
	defer j.log.Infof("End auto depoy job saving device settings")

	for _, dev := range devices {
		// 查询是否有失败的装机记录
		old, err := j.repo.GetDeviceSettingBySN(dev.SN)
		if err != nil && err != gorm.ErrRecordNotFound {
			j.log.Errorf("Fail to get exist devicesetting of device(sn: %s ), error: %s", dev.SN, err)
			continue
		}
		// 仅处理不存在装机记录的首次部署场景
		if old != nil {
			j.log.Errorf("Exist devicesetting of device(sn: %s ), skipping now", dev.SN)
			continue
		} 
		// 根据 USiteID 查询获取 PhysicalArea
		var ServerUSite *model.ServerUSite
		ServerUSite, err = j.repo.GetServerUSiteByID(*dev.USiteID)
		if err != nil {
			j.log.Errorf("Fail to get serverusite of device(sn: %s ), error: %s", dev.SN, err)
			continue
		}

		// 根据 Category 判断 Arch
		if strings.HasPrefix(dev.Category, "GA-") {
			dev.Arch = "arm"
		} else {
			dev.Arch = "x86"
		}

		// 转换 Vendor 厂商属性值(规则限定英文小写)
		lower_vendor := strings.ToLower(dev.Vendor)
		switch lower_vendor {
		case "华为":
			lower_vendor = "huawei"
		case "浪潮":
			lower_vendor = "inspur"
		case "联想":
			lower_vendor = "lenovo"
		case "戴尔":
			lower_vendor = "dell"
		case "惠普":
			lower_vendor = "hp"
		case "华三":
			lower_vendor = "h3c"
		case "超聚变":
			lower_vendor = "xfusion"
		case "紫光恒越":
			lower_vendor = "unis"
		case "中科可控":
			lower_vendor = "suma"
		case "中科曙光":
			lower_vendor = "sugon"		
		}
		// 获取设备类型是否为金融信创生态产品
		cond := model.DeviceCategory{
			Category: dev.Category,
		}
		items, err := j.repo.GetDeviceCategorys(&cond, nil, nil)
		if err != nil {
			j.log.Errorf("Fail to get category of device(sn: %s ), error: %s", dev.SN, err)
			continue
		}
		isFITIEcoProduct := "unknown"
		if len(items) != 0 {
			isFITIEcoProduct = items[0].IsFITIEcoProduct
			j.log.Debugf("device[SN:%s] category[%s] is FITIEcoProduct", dev.SN, dev.Category)
		}
		// 规则引擎的事实对象
		df := DeviceFact{
			SN:                   dev.SN,
			Category:             dev.Category,
			PhysicalArea:         ServerUSite.PhysicalArea,
			Arch:                 dev.Arch,
			Vendor:               lower_vendor,
			IsFITIEcoProduct:     isFITIEcoProduct,
		}
		// 规则推理得到装机参数
		result, err := j.RuleInduction(&df)
		if err != nil {
			j.log.Errorf("Fail to get rule action results of device(sn: %s ), error: %s", dev.SN, err)
			continue
		}
		// 装机参数对象
		ds := model.DeviceSetting{
			SN:             dev.SN,
			InstallType:    model.InstallationImage,
			Status:         model.InstallStatusPre,
			Creator:        "auto_deploy",
			NeedIntranetIPv6: model.NO,				//默认仅分配IPV4
			NeedExtranetIPv6: model.NO,				//默认仅分配IPV4
		}
		// 获取硬件配置模板ID
		hwTemplate, err := j.repo.GetHardwareTemplateByName(result.RaidResult)
		if err != nil {
			j.log.Errorf("Fail to GetHardwareTemplateByName of device(sn: %s name: %s ), error: %s", dev.SN, result.RaidResult, err)
			continue
		}
		ds.HardwareTemplateID = hwTemplate.ID
		// 获取操作系统模版ID
		// 优先image 次之pxe
		imgTemplate, err := j.repo.GetImageTemplateByName(result.OSResult)
		if err != nil {
			j.log.Errorf("Fail to GetImageTemplateByName of device(sn: %s name: %s ), error: %s ,try SystemTemplate", dev.SN, result.OSResult, err)
			//continue
			sysTemplate, err := j.repo.GetSystemTemplateByName(result.OSResult)
			if err != nil {
				j.log.Errorf("Fail to GetSystemTemplateByName of device(sn: %s name: %s ), error: %s", dev.SN, result.OSResult, err)
				continue
			}
			ds.SystemTemplateID = sysTemplate.ID
			ds.InstallType = model.InstallationPXE
		}
		ds.ImageTemplateID = imgTemplate.ID

		// 获取网络配置
		if result.NetworkResult == "la_wa" {
			// 获取内网业务IP
			intraIP, err := j.repo.AssignIntranetIP(dev.SN)
			if err != nil {
				j.log.Errorf("AutoGenDeviceSetting fail to AssignIntranetIP of device(sn: %s ), error: %s", dev.SN, err)
				continue
			}
			ds.IntranetIPNetworkID = intraIP.IPNetworkID
			ds.IntranetIP = intraIP.IP
			j.log.Infof("Device(%s) is assigned to intranet ip: %s", dev.SN, intraIP.IP)
			// 获取外网业务IP
			extraIP, err := j.repo.AssignExtranetIP(dev.SN)
			if err != nil {
				j.log.Errorf("AutoGenDeviceSetting fail to AssignExtranetIP of device(sn: %s ), error: %s", dev.SN, err)
				continue
			}
			ds.ExtranetIPNetworkID = extraIP.IPNetworkID
			ds.ExtranetIP = extraIP.IP
			ds.NeedExtranetIP = "yes"
			j.log.Infof("Device(%s) is assigned to extranet ip: %s", dev.SN, extraIP.IP)
		} else if result.NetworkResult == "la" {
			// 获取内网业务IP
			intraIP, err := j.repo.AssignIntranetIP(dev.SN)
			if err != nil {
				j.log.Errorf("AutoGenDeviceSetting fail to AssignIntranetIP of device(sn: %s ), error: %s", dev.SN, err)
				continue
			}
			ds.IntranetIPNetworkID = intraIP.IPNetworkID
			ds.IntranetIP = intraIP.IP
			ds.NeedExtranetIP = "no"
			j.log.Infof("Device(%s) is assigned to intranet ip: %s", dev.SN, intraIP.IP)
		} else {
			// 网络配置规则推理结果不符合la la_wa 时跳过
			j.log.Errorf("Fail to AssignIP for device(sn: %s  expected network_conf: %s ), error: %s", dev.SN, result.NetworkResult, err)
			continue			
		}
		
		// 当RAID、OS、NETWORK均得到对应参数时，保存装机参数
		if err = j.repo.SaveDeviceSetting(&ds); err != nil {
			j.log.Errorf("Device(%s) settings(auto) have not been saved err: %s", dev.SN, err)
			return nil
		}
		// 变更记录
		optDetail, err := j.convert2DetailOfOperationTypeOSInstall(ds)
		if err != nil {
			j.log.Errorf("Fail to convert Detail of OperationTypeOSInstall: %v", err)
		}
		devLog := model.ChangeLog {
			OperationUser:		"auto_deploy",
			OperationType:		model.OperationTypeOSInstall,
			OperationDetail:	optDetail,
			OperationTime:		times.ISOTime(time.Now()).ToTimeStr(),
		}
		adll := &AppendDeviceLifecycleLogReq{
			SN:					dev.SN,
			LifecycleLog: 		devLog,
		}
		if err = j.AppendDeviceLifecycleLogBySN(adll);err != nil {
			j.log.Error("LifecycleLog: append device lifecycle log (sn:%v) fail(%s)", dev.SN, err.Error())
		}		

		j.log.Infof("Device(%s) settings(auto) have been saved", dev.SN)
		go j.startInstallation(limiter.GlobalLimiter, dev.SN)
	}
	return nil
}


// 事实对象的属性（与规则库所需的属性匹配）
type DeviceFact struct {
	SN                 string      // 设备SN
	Category           string      // 设备类型
	PhysicalArea       string      // 物理区域， 如： ServerFarm bonding区1
	Arch               string      // cpu架构：x86\arm
	Vendor             string      // 设备厂商
	IsFITIEcoProduct   string      // 是否金融信创生态产品
}

// 推理规则前件是否包含逻辑运算符号AND OR
type ConditionHasLogicalOperator struct {
	HasAND    bool
	HasOR     bool
	HasBoth   bool
}

// 规则推理返回结构体
type ActionResult struct {
	RuleType        string    // os raid network
	OSResult        string
	RaidResult      string
	NetworkResult   string
}

// 元规则推理，仅针对单个 RuleP 结构体进行推理，返回 bool
func (j *AutoDeployJob) MetaRuleInduction(fact *DeviceFact, rulep *model.RuleP) (bool) {
	if rulep.Attribute != "" && rulep.Attribute != model.AttributeLogicalOperator {
		var factvalue string
		// 获取fact对应的属性值
		switch rulep.Attribute {
		case model.AttributeCategory:
			factvalue = fact.Category
		case model.AttributePhysicalArea:
			factvalue = fact.PhysicalArea
		case model.AttributeArch:
			factvalue = fact.Arch
		case model.AttributeVendor:
			factvalue = fact.Vendor
		case model.AttributeIsFITIEcoProduct:
			factvalue = fact.IsFITIEcoProduct			
		default:
			j.log.Errorf("Attribute of rule is not defined.")
			return false
		}
		
		// 逻辑处理 equal contains in
		switch rulep.Operator {
		case model.OperatorEqual:
			if factvalue == rulep.Value[0] {
				j.log.Infof("Device %s  Meta_Rule_Hit: %s equal %s ", fact.SN, factvalue, rulep.Value)
				return true
			}
		case model.OperatorContains:
			if strings.Contains(factvalue, rulep.Value[0]) {
				j.log.Infof("Device %s  Meta_Rule_Hit: %s contains %s ", fact.SN, factvalue, rulep.Value)
				return true
			}
		case model.OperatorIN:
			for _, each := range rulep.Value {
				if factvalue == each {
					j.log.Infof("Device %s  Meta_Rule_Hit: %s in %s ", fact.SN, factvalue, rulep.Value)
					return true
				}
			}
		default:
			j.log.Errorf("Operator of rule is not defined.")
			return false
		}
	}
    return false
}

// 前件解析器
func (j *AutoDeployJob) ConditionParser(fact *DeviceFact, rule *model.DeviceSettingRule) (bool, error) {
	if rule.Condition != "" {
		var ruleconditions []model.RuleP
		chlo := ConditionHasLogicalOperator{
			HasAND:  false,
			HasOR:   false,
			HasBoth: false,
		}

		// 解析规则前件
		j.log.Debugf(rule.Condition)
		err := json.Unmarshal([]byte(rule.Condition), &ruleconditions)
		if err != nil {
			j.log.Errorf("Rule condition parser failed (id: %d)", rule.ID)
			j.log.Error(err)
			return false,err
		}
        // fact 与 fact 之间必须通过逻辑运算符连接，故index=奇数可增加校验是否为AND OR
		j.log.Infof("Condition of rule %d : %s", rule.ID, ruleconditions)
		//j.log.Debugf("Condition of rule %d value: %s  type: %T", rule.ID, ruleconditions, ruleconditions)

		// 推理是否包含逻辑运算符 AND OR
		if len(ruleconditions) > 1 {
		    for _, rc := range ruleconditions {
		    	if rc.Attribute == model.AttributeLogicalOperator && rc.Operator == model.OperatorEqual {
		    		switch rc.Value[0] {
		    		case model.OperatorOR:
		    			chlo.HasOR = true
		    		case model.OperatorAND:
		    			chlo.HasAND = true
		    		}
		    	}
		    }
		    if chlo.HasOR && chlo.HasAND {
		    	chlo.HasBoth = true
		    }
	    } else if len(ruleconditions) == 1 {
			for _, rc := range ruleconditions {
				return j.MetaRuleInduction(fact, &rc), nil
			}			
		}
		
		// 前件包含AND + OR 优先处理 OR
		if chlo.HasBoth {
			// 针对包含AND + OR 临时存储结果 [true,and,false,or,true]
        	var logicresult []string			
			for _, rc := range ruleconditions {
				if rc.Attribute == model.AttributeLogicalOperator {
		    		switch rc.Value[0] {
		    		case model.OperatorOR:
		    			logicresult = append(logicresult, "or")
		    		case model.OperatorAND:
						logicresult = append(logicresult, "and")
		    		}
		    	} else {
					if j.MetaRuleInduction(fact, &rc) {
						logicresult = append(logicresult, "true")
					} else {
						logicresult = append(logicresult, "false")
					}
				}
			}
			
			// 逻辑结果数组按and拆分 [[true],[false,or,true]]
			var logicresult_split_by_and [][]string
			i := 0
			for _, v := range logicresult {
				if v == "and" {
					i++
				} else {
					logicresult_split_by_and[i] = append(logicresult_split_by_and[i], v)
				}
			}

			// 处理返回最终bool
			contain_true := false
			for _, splitted_list := range logicresult_split_by_and {
				if len(splitted_list) == 1 {
					if splitted_list[0] == "false"{
						return false,nil
					} else {
						continue
					}
				} else {
					for _, va := range splitted_list {
						if va == "true" {
							contain_true = true
						} else {
							continue
						}
					}
					if contain_true == false {
						return false,nil
					}
				}
			}
			return true, nil
		} else if chlo.HasAND {
			for _, rc := range ruleconditions {
				if rc.Attribute == model.AttributeLogicalOperator {
		    		switch rc.Value[0] {
		    		case model.OperatorOR:
						j.log.Errorf("Not expected OR of rule %d ", rule.ID)
						return false, fmt.Errorf("Not expected OR of rule %d ", rule.ID)
		    		case model.OperatorAND:
						continue
		    		}
		    	} else {
					if j.MetaRuleInduction(fact, &rc) {
						continue
					} else {
						return false, nil
					}
				}
			}
			return true, nil
		} else if chlo.HasOR {
			for _, rc := range ruleconditions {
				if rc.Attribute == model.AttributeLogicalOperator {
		    		switch rc.Value[0] {
					case model.OperatorOR:
						continue
		    		case model.OperatorAND:
						j.log.Errorf("Not expected AND of rule %d ", rule.ID)
						return false, fmt.Errorf("Not expected AND of rule %d ", rule.ID)
		    		}
		    	} else {
					if j.MetaRuleInduction(fact, &rc) {
						return true, nil
					} else {
						continue
					}
				}
			}
			return false, nil
		}
	} 
	j.log.Errorf("Empty of rule %d ", rule.ID)
	return false, fmt.Errorf("Empty of rule %d ", rule.ID)
}


// 规则推理机
func (j *AutoDeployJob) RuleInduction(fact *DeviceFact) (*ActionResult, error) {
	// 根据规则分类提取规则进行处理
	var result = &ActionResult{}
	var ruletypes = []string {model.DeviceSettingRuleOS, model.DeviceSettingRuleRaid, model.DeviceSettingRuleNetwork}
	for _, ruletype := range ruletypes {
		switch ruletype {
		case model.DeviceSettingRuleOS:
			ruleItems, err := j.repo.GetDeviceSettingRulesByType(ruletype)
			if err != nil {
				j.log.Errorf("Fact(device sn: %s ) failed to get os rules.", fact.SN)
				return nil, err
			}
			for _, rule := range ruleItems {
				b, err := j.ConditionParser(fact, rule)
				if err != nil {
					return nil, err
				}
				if b {
					j.log.Infof("Fact(device sn: %s ) successed to get os result.", fact.SN)
					result.OSResult = rule.Action
				}
			}
	
		case model.DeviceSettingRuleRaid:
			ruleItems, err := j.repo.GetDeviceSettingRulesByType(ruletype)
			if err != nil {
				j.log.Errorf("Fact(device sn: %s ) failed to get raid rules.", fact.SN)
				return nil, err
			}
			for _, rule := range ruleItems {
				b, err := j.ConditionParser(fact, rule)
				if err != nil {
					return nil, err
				}
				if b {
					j.log.Infof("Fact(device sn: %s ) successed to get raid result.", fact.SN)
					result.RaidResult = rule.Action
				}
			}
		
		case model.DeviceSettingRuleNetwork:
			ruleItems, err := j.repo.GetDeviceSettingRulesByType(ruletype)
			if err != nil {
				j.log.Errorf("Fact(device sn: %s ) failed to get network rules.", fact.SN)
				return nil, err
			}
			for _, rule := range ruleItems {
				b, err := j.ConditionParser(fact, rule)
				if err != nil {
					return nil, err
				}
				if b {
					j.log.Infof("Fact(device sn: %s ) successed to get network result.", fact.SN)
					result.NetworkResult = rule.Action
				}
			}
		}
	}
	// 仅返回完整的参数推理结果
	if result.OSResult == "" || result.RaidResult == "" || result.NetworkResult == "" {
		j.log.Errorf("Fact(device sn: %s ) do not hit all types of rules.", fact.SN)
		return nil, fmt.Errorf("Fact(device sn: %s ) do not hit all types of rules.", fact.SN)
	}
	return result, nil
}


// startInstallation 带外远程开机并开始安装流程
func (j *AutoDeployJob) startInstallation(lim limiter.Limiter, sn string) (err error) {
	// 1、获取进入bootos所需令牌
	var bucket limiter.Bucket
	var token limiter.Token

	if j.conf.DHCPLimiter.Enable {
		bucket, err = lim.Route(sn)
		if err != nil {
			return err
		}
		token, err = bucket.Acquire(sn, time.Second*time.Duration(j.conf.DHCPLimiter.WaitingTimeout))
		if err != nil {
			return err
		}
	}

	// 2、修改装机开始时间、装机状态
	now := time.Now()
	if _, err = j.repo.UpdateDeviceSettingBySN(&model.DeviceSetting{
		SN: sn,
		InstallationStartTime: &now,
		Status:                model.InstallStatusIng,
	}); err != nil {
		if j.conf.DHCPLimiter.Enable {
			_ = bucket.Return(sn, token)
		}
		return err
	}

	j.log.Infof("Remote boot and start the OS installation process of the device(%s)", sn)
	var sb strings.Builder
	//TODO 优化项：如果设备在bootos,则跳过重启的过程
	//如果设备已上电，PXE重启，如果关机该命令会有问题，此时再尝试开机
	out, err := j.OperateOOBPower(sn, j.conf.Crypto.Key, j.conf.Server.OOBDomain, PowerRestart, true)
	if err != nil {
		j.log.Errorf("PXE-PowerRestart重启失败(%s)，下一步尝试开机",err.Error())
		sb.WriteString(fmt.Sprintf("\nPXE-PowerRestart重启失败(%s)，下一步尝试开机\n", err.Error()))
		out, err = j.OperateOOBPower(sn, j.conf.Crypto.Key, j.conf.Server.OOBDomain, PowerOn, true)
		sb.WriteString(strings.TrimSpace(out))
		if err != nil {
			j.log.Errorf("PXE-PowerOn开机失败(%s)",err.Error())
			sb.WriteString(fmt.Sprintf("PXE-PowerOn开机失败(%s)", err.Error()))
		}
	}
	// 调用自身API进行安装进度上报(进度值为'5%')
	var reqData struct {
		Progress float64 `json:"progress"`
		Log      string  `json:"log"`
		Title    string  `json:"title"`
	}
	if err != nil {
		reqData.Progress = -1.0
	} else {
		reqData.Progress = 0.05
	}
	reqData.Title = "物理机远程开机"
	reqData.Log = base64.StdEncoding.EncodeToString([]byte(sb.String()))

	reqBody, _ := json.Marshal(&reqData)

	url := fmt.Sprintf("%s/api/cloudboot/v1/devices/%s/installations/progress", fmt.Sprintf("http://localhost:%d", j.conf.Server.HTTPPort), sn)
	j.log.Infof("POST %s", url)
	j.log.Infof("Request body: %s", reqBody)

	resp, err := httplib.Post(url).
		Header("Content-Type", "application/json").
		Header("Accept", "application/json").
		Body(reqBody).Response()
	if err != nil {
		j.log.Error(err)
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		j.log.Error(err)
		return err
	}
	j.log.Infof("Response body: %s", respBody)
	return err
}

// OperateOOBPower 带外管理的操作
func (j *AutoDeployJob) OperateOOBPower(sn, key, oobDomain, operate string, isPxe bool) (output string, err error) {
	m, err := j.repo.GetDeviceBySN(sn)
	if err != nil {
		return "", err
	}

	if m.OOBUser == "" || m.OOBPassword == "" {
		j.log.Warnf("设备带外用户或密码为空，尝试找回，[SN:%s]", m.SN)
		history, err := j.repo.GetLastOOBHistoryBySN(sn)
		if err != nil && err == gorm.ErrRecordNotFound {
			j.log.Errorf("find back oob history by sn:%s fail,%s", sn, err.Error())
			return fmt.Sprintf("find back oob history by sn:%s fail", sn), fmt.Errorf("设备用户名密码为空，无法操作")
		} else {
			return fmt.Sprintf("find back oob history by sn:%s fail", sn), err
		}
		m.OOBUser = history.UsernameNew
		m.OOBPassword = history.PasswordNew
	}

	oobHost := utils.GetOOBHost(m.SN, m.Vendor, oobDomain)
	oobIP := oob.TransferHostname2IP(j.log, j.repo, m.SN, oobHost)
	if oobIP == "" {
		return "操作失败", errors.New("未获取到带外IP")
	}
	oobUser := m.OOBUser
	oobPassword, err := utils.AESDecrypt(m.OOBPassword, []byte(key))
	if err != nil {
		j.log.Debugf("descrypt password failure, err: %s", err.Error())
		return "", err
	}

	//TODO ping命令存在问题
	//// check ping
	//if !ping.Ping(oobHost, 5) {
	//	return "", fmt.Errorf("网络无法连通设备 [SN:%s]", m.SN)
	//}

	// check is power
	isPowerOn, err := j.OOBPowerStatus(oobIP, oobUser, string(oobPassword), m.OOBPassword)
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

	cmd = j.replaceCmd(cmd, oobIP, oobUser, string(oobPassword))

	j.DesensitizePasswordLog(fmt.Sprintf("start to exec cmd: [%s]", cmd), string(oobPassword), m.OOBPassword)

	out, err := sh.ExecOutputWithLog(j.log, cmd)
	if err != nil {
		j.log.Debugf("exec [%s] done , err: [%s], stdout: [%s]", sh.CmdDesensitization(cmd), err.Error(), sh.CmdDesensitization(string(out)))
		return "", j.ProcessStdoutMsg(out, oobIP, oobUser)
	}

	j.DesensitizePasswordLog(fmt.Sprintf("exec [%s] done，output: [%s]", cmd, string(out)), string(oobPassword), m.OOBPassword)

	// update device power_status
	status := model.PowerStatusOn

	if operate == PowerOff {
		status = model.PowerStatusOff
	}
	m.PowerStatus = status
	j.repo.UpdateDevice(m)

	return string(out), nil
}

// OOBPowerStatus 检查OOBPower状态
func (j *AutoDeployJob) OOBPowerStatus(oobHost, oobUser, oobPassword, oldPassword string) (bool, error) {
	if oobHost == "" {
		return false, errors.New("未获取到带外IP")
	}
	cmd := j.replaceCmd(OperatePowerStatus, oobHost, oobUser, oobPassword)

	output, err := sh.ExecOutputWithLog(j.log, cmd)
	if err != nil {
		return false, j.ProcessStdoutMsg(output, oobHost, oobUser)
	}

	return strings.Contains(string(output), "Chassis Power is on"), nil
}

// replaceCmd 将HOST、USER、PASSWD进行替换
func (j *AutoDeployJob) replaceCmd(cmd, oobHost, oobUser, oobPassword string) string {
	cmd = strings.Replace(cmd, "$HOST", oobHost, -1)
	cmd = strings.Replace(cmd, "$USER", oobUser, -1)
	cmd = strings.Replace(cmd, "$PASSWD", oobPassword, -1)
	return cmd
}

// ProcessStdoutMsg 处理错误信息
func (j *AutoDeployJob) ProcessStdoutMsg(output []byte, host, username string) error {
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
func (j *AutoDeployJob) DesensitizePasswordLog(str, oldPass, newPass string) {
	j.log.Debugf("%s", strings.Replace(str, oldPass, newPass, -1))
}


// 更新追加设备生命周期变更记录 请求结构体
type AppendDeviceLifecycleLogReq struct {
	SN					string				`json:"sn"`
	LifecycleLog		model.ChangeLog		`json:"lifecycle_log"`
}


//AppendDeviceLifecycleLog 更新追加设备生命周期变更记录
func (j *AutoDeployJob) AppendDeviceLifecycleLogBySN(reqData *AppendDeviceLifecycleLogReq) error {
	// DeviceLifecycle 查询是否已经存在
	devLifecycle, err := j.repo.GetDeviceLifecycleBySN(reqData.SN)
	if err != nil && err != gorm.ErrRecordNotFound {
		j.log.Error(err.Error())
		return err
	}
	if devLifecycle != nil {
		j.log.Debugf("Begin to AppendDeviceLifecycleLogBySN:%s", reqData.SN)
		defer j.log.Debugf("End to AppendDeviceLifecycleLogBySN:%s", reqData.SN)
		// 获取当前的生命周期日志记录
		var devLL []model.ChangeLog
		if devLifecycle.LifecycleLog != "" {
			if err = json.Unmarshal([]byte(devLifecycle.LifecycleLog), &devLL);err != nil {
				j.log.Error(err.Error())
				return err
			}
		}
		// 追加
		devLL = append(devLL, reqData.LifecycleLog)
		b, err := json.Marshal(devLL)
		if err != nil {
			j.log.Error(err.Error())
			return err
		}
		// DeviceLifecycle 结构体
		modDevLifecycle := &model.DeviceLifecycle{
			SN:             				reqData.SN,
			LifecycleLog:					string(b),
		}
		if err = j.repo.UpdateDeviceLifecycleBySN(modDevLifecycle);err != nil {
			j.log.Errorf("UpdateDeviceLifecycleBySN failed: %s", err.Error())
			return err
		}
	}
	return nil
}


func (j *AutoDeployJob) convert2DetailOfOperationTypeOSInstall(setting model.DeviceSetting) (string, error) {
	var details []string

	getType := reflect.TypeOf(setting)
	getValue := reflect.ValueOf(setting)

	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i)
		value := getValue.Field(i).Interface()

		switch field.Name {
		case "InstallType":
			details = append(details, fmt.Sprintf("安装方式:%v ", value))
		case "OSTemplateName":
			details = append(details, fmt.Sprintf("操作系统安装模板:%v ", value))
		case "HardwareTemplateName":
			details = append(details, fmt.Sprintf("硬件配置模板:%v ", value))
		case "NeedExtranetIP":
			details = append(details, fmt.Sprintf("是否需要外网IPv4:%v ", value))
		case "NeedExtranetIPv6":
			details = append(details, fmt.Sprintf("是否需要内网IPv6:%v ", value))
		case "NeedIntranetIPv6":
			details = append(details, fmt.Sprintf("是否需要外网IPv6:%v ", value))
		}
	}

	if len(details) > 0 {
		return strings.Join(details, "；"), nil
	} else {
		return "", fmt.Errorf("部署设备记录（SN：%v）字段详情转化失败", setting.SN)
	}

}