package service

import (
	"errors"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"

	"idcos.io/cloudboot/limiter"

	"github.com/astaxie/beego/httplib"
	"github.com/jinzhu/gorm"
	"github.com/voidint/binding"
	"github.com/voidint/page"
    
	"idcos.io/cloudboot/config"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/server/cloudbootserver/types/setting"
	"idcos.io/cloudboot/utils/centos6"
	"idcos.io/cloudboot/utils/collection"
	mystrings "idcos.io/cloudboot/utils/strings"
	"idcos.io/cloudboot/utils/times"
)

// SaveDeviceSettingItem 保存设备装机参数条目
type SaveDeviceSettingItem struct {
	// 设备序列号
	SN string `json:"sn"`
	// InstallType 安装方式
	// Enum: image,pxe
	InstallType string `json:"install_type"`
	// 操作系统安装模板
	OSTemplateName string `json:"os_template_name"`
	// 硬件配置模板
	HardwareTemplateName string `json:"hardware_template_name"`
	// NeedExtranetIP 是否需要外网IP。可选值：yes-是; no-否;
	// Enum: yes,no
	NeedExtranetIP 			string `json:"need_extranet_ip"`
	// 是否需要内外网IPv6    Enum: yes,no
	NeedExtranetIPv6		string `json:"need_extranet_ipv6"`
	NeedIntranetIPv6		string `json:"need_intranet_ipv6"`
}

// SaveDeviceSettingsReq 保存设备装机参数请求结构体
type SaveDeviceSettingsReq struct {
	Settings    DeviceSettings
	CurrentUser *model.CurrentUser // 操作人
}

// DeviceSettings 设备装机参数集合
type DeviceSettings []*SaveDeviceSettingItem

// Validate 装机参数校验
func (reqData DeviceSettings) Validate(r *http.Request, errs binding.Errors) binding.Errors {
	for i := range reqData {
		if errs = reqData[i].validateOne(r, errs); errs.Len() > 0 {
			return errs
		}
	}
	return errs
}

// validateOne 校验单台设备的装机参数
func (sett *SaveDeviceSettingItem) validateOne(r *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(r.Context())

	if sett.SN == "" {
		errs.Add([]string{"sn"}, binding.RequiredError, "SN不能为空")
		return errs
	}
	devSett, _ := repo.GetDeviceSettingBySN(sett.SN)
	if devSett != nil && devSett.Status == model.InstallStatusIng {
		errs.Add([]string{"sn"}, binding.RequiredError, fmt.Sprintf("[%s]正在部署，不能重复提交", sett.SN))
		return errs
	}
	if sett.InstallType == "" {
		errs.Add([]string{"install_type"}, binding.RequiredError, "安装方式(install_type)不能为空")
		return errs
	}
	if !collection.InSlice(sett.InstallType, []string{model.InstallationPXE, model.InstallationImage}) {
		errs.Add([]string{"install_type"}, binding.BusinessError, fmt.Sprintf("[%s]无效的安装方式(install_type)值: %s", sett.SN, sett.InstallType))
		return errs
	}

	if sett.OSTemplateName == "" {
		errs.Add([]string{"os_template_name"}, binding.RequiredError, fmt.Sprintf("[%s]操作系统模板(os_template_name)不能为空", sett.SN))
		return errs
	}

	switch sett.InstallType {
	case model.InstallationImage:
		if _, err := repo.GetImageTemplateByName(sett.OSTemplateName); err == gorm.ErrRecordNotFound {
			errs.Add([]string{"os_template_name"}, binding.BusinessError, fmt.Sprintf("[%s]无效的操作系统模板(os_template_name): %s", sett.SN, sett.OSTemplateName))
			return errs
		}
	case model.InstallationPXE:
		if _, err := repo.GetSystemTemplateByName(sett.OSTemplateName); err == gorm.ErrRecordNotFound {
			errs.Add([]string{"os_template_name"}, binding.BusinessError, fmt.Sprintf("[%s]无效的操作系统模板(os_template_name): %s", sett.SN, sett.OSTemplateName))
			return errs
		}
	}

	if sett.HardwareTemplateName == "" {
		errs.Add([]string{"hardware_template_name"}, binding.RequiredError, fmt.Sprintf("[%s]硬件配置模板(hardware_template_name)不能为空", sett.SN))
		return errs
	}

	if _, err := repo.GetHardwareTemplateByName(sett.HardwareTemplateName); err == gorm.ErrRecordNotFound {
		errs.Add([]string{"hardware_template_name"}, binding.BusinessError, fmt.Sprintf("[%s]无效的硬件配置模板(hardware_template_name): %s", sett.SN, sett.HardwareTemplateName))
		return errs
	}

	if sett.NeedExtranetIP == "" {
		errs.Add([]string{"need_extranet_ip"}, binding.RequiredError, fmt.Sprintf("[%s]是否配置外网IP(need_extranet_ip)不能为空", sett.SN))
		return errs
	}
	if !collection.InSlice(sett.NeedExtranetIP, []string{model.YES, model.NO}) {
		errs.Add([]string{"need_extranet_ip"}, binding.BusinessError, fmt.Sprintf("[%s]无效的是否配置外网IP(need_extranet_ip)值: %s", sett.SN, sett.NeedExtranetIP))
		return errs
	}
	// 校验当前设备状态，若设备信息存在且状态为已退役、待退役，则不允许安装操作系统。
	dev, err := repo.GetDeviceBySN(sett.SN)
	if err == gorm.ErrRecordNotFound {
		errs.Add([]string{"sn"}, binding.BusinessError, fmt.Sprintf("[%s]请先录入物理机信息", sett.SN))
		return errs
	}
	if err != nil {
		errs.Add([]string{"sn"}, binding.SystemError, "系统内部发生错误")
		return errs
	}
	if dev.OperationStatus == model.DevOperStatPreRetire || dev.OperationStatus == model.DevOperStatRetiring|| dev.OperationStatus == model.DevOperStateRetired {
		errs.Add([]string{"sn"}, binding.BusinessError, fmt.Sprintf("[%s]待退役/退役中/已退役物理机无法上架部署", sett.SN))
		return errs
	}
	// 校验待分配的IP网段是否存在
	_, err = repo.GetIntranetIPNetworksBySN(sett.SN)
	if err == gorm.ErrRecordNotFound {
		errs.Add([]string{"sn"}, binding.BusinessError, fmt.Sprintf("[%s]无法找到合适的内网网段并分配IP", sett.SN))
		return errs
	}
	if err != nil {
		errs.Add([]string{"sn"}, binding.BusinessError, err.Error())
		return errs
	}
	if sett.NeedExtranetIP == model.YES {
		_, err = repo.GetExtranetIPNetworksBySN(sett.SN)
		if err == gorm.ErrRecordNotFound {
			errs.Add([]string{"sn"}, binding.BusinessError, fmt.Sprintf("[%s]无法找到合适的外网网段并分配IP", sett.SN))
			return errs
		}
		if err != nil {
			errs.Add([]string{"sn"}, binding.BusinessError, err.Error())
			return errs
		}
	}
	if sett.NeedIntranetIPv6 == model.YES {
		_, err := repo.GetIPv6NetworkBySN(sett.SN, model.Intranet)
		if err == gorm.ErrRecordNotFound {
			errs.Add([]string{"sn"}, binding.BusinessError, fmt.Sprintf("[%s]无法找到合适的IPv6内网网段并分配IP", sett.SN))
			return errs
		}
		if err != nil {
			errs.Add([]string{"sn"}, binding.BusinessError, err.Error())
			return errs
		}
	}
	if sett.NeedExtranetIPv6 == model.YES {
		_, err := repo.GetIPv6NetworkBySN(sett.SN, model.Extranet)
		if err == gorm.ErrRecordNotFound {
			errs.Add([]string{"sn"}, binding.BusinessError, fmt.Sprintf("[%s]无法找到合适的IPv6外网网段并分配IP", sett.SN))
			return errs
		}
		if err != nil {
			errs.Add([]string{"sn"}, binding.BusinessError, err.Error())
			return errs
		}
	}
	if tor, _ := repo.GetTORBySN(sett.SN); tor == "" {
		errs.Add([]string{"sn"}, binding.BusinessError, fmt.Sprintf("[%s]无法找到TOR", sett.SN))
		return errs
	}

	// TODO 校验网段内是否还有剩余IP可供分配
	return errs
}

// SaveDeviceSettings 批量保存设备装机参数。
// 多台设备的装机参数之间并不包含事务，保存成功的设备的序列号(SN)将被返回。
func SaveDeviceSettings(log logger.Logger, repo model.Repo, conf *config.Config, lim limiter.Limiter, reqData *SaveDeviceSettingsReq) (succeeds []string, err error) {
	for i := range reqData.Settings {
		if err = saveDeviceSetting(log, repo, conf, lim, reqData.Settings[i], reqData.CurrentUser.LoginName); err != nil {
			return succeeds, err
		}
		succeeds = append(succeeds, reqData.Settings[i].SN)
		// 变更记录
		optDetail, err := convert2DetailOfOperationTypeOSInstall(repo, *reqData.Settings[i])
		if err != nil {
			log.Errorf("Fail to convert Detail of OperationTypeOSInstall: %v", err)
		}
		devLog := model.ChangeLog {
			OperationUser:		reqData.CurrentUser.Name,
			OperationType:		model.OperationTypeOSInstall,
			OperationDetail:	optDetail,
			OperationTime:		times.ISOTime(time.Now()).ToTimeStr(),
		}
		adll := &AppendDeviceLifecycleLogReq{
			SN:					reqData.Settings[i].SN,
			LifecycleLog: 		devLog,
		}
		if err = AppendDeviceLifecycleLogBySN(log, repo, adll);err != nil {
			log.Error("LifecycleLog: append device lifecycle log (sn:%v) fail(%s)", reqData.Settings[i].SN, err.Error())
		}

	}
	return succeeds, nil
}

// saveDeviceSetting 保存单台设备的装机参数
func saveDeviceSetting(log logger.Logger, repo model.Repo, conf *config.Config, lim limiter.Limiter, sett *SaveDeviceSettingItem, LoginName string) (err error) {
	log.Infof("Start saving device(%s) settings", sett.SN)
	defer log.Infof("End saving device(%s) settings", sett.SN)

	ds := model.DeviceSetting{
		SN:             	sett.SN,
		InstallType:    	sett.InstallType,
		Status:         	model.InstallStatusPre,
		NeedExtranetIP: 	sett.NeedExtranetIP,
		NeedIntranetIPv6: 	sett.NeedIntranetIPv6,
		NeedExtranetIPv6: 	sett.NeedExtranetIPv6,
	}

	// 获取硬件配置模板ID
	hwTemplate, err := repo.GetHardwareTemplateByName(sett.HardwareTemplateName)
	if err != nil {
		return nil
	}
	ds.HardwareTemplateID = hwTemplate.ID
	// 查询旧装机记录
	old, err := repo.GetDeviceSettingBySN(sett.SN)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if old != nil {
		ds.ID = old.ID
		ds.CreatedAt = old.CreatedAt
		ds.Updater = LoginName
		// 旧安装模板ID 置0 
		ds.ImageTemplateID = 0
		ds.SystemTemplateID = 0
	} else {
		ds.Creator = LoginName
	}

	// 获取操作系统安装模板ID
	switch sett.InstallType {
	case model.InstallationImage:
		imgTemplate, err := repo.GetImageTemplateByName(sett.OSTemplateName)
		if err != nil {
			return nil
		}
		ds.ImageTemplateID = imgTemplate.ID
	case model.InstallationPXE:
		sysTemplate, err := repo.GetSystemTemplateByName(sett.OSTemplateName)
		if err != nil {
			return nil
		}
		ds.SystemTemplateID = sysTemplate.ID
	}

	log.Infof("Start assigning IP to the device(%s)", sett.SN)
	// 记录本次已分配的IP ID以便回滚
	var ipIDsForRollback []uint

	//根据需要分配IP，假如有装机记录，并且有IP，则不用重复分配了
	if old == nil || old.IntranetIP == "" {
		// 获取内网业务IP
		intraIP, err := repo.AssignIntranetIP(sett.SN)
		if err != nil {
			return err
		}
		ds.IntranetIPNetworkID = intraIP.IPNetworkID
		ds.IntranetIP = intraIP.IP
		log.Infof("Device(%s) is assigned to intranet ip: %s", sett.SN, intraIP.IP)
		ipIDsForRollback = append(ipIDsForRollback, intraIP.ID)
	} else {
		ds.IntranetIPNetworkID = old.IntranetIPNetworkID
		ds.IntranetIP = old.IntranetIP
	}

	if sett.NeedExtranetIP == model.YES {
		// 获取外网业务IP
		log.Infof("Device(%s) need extranet ip", sett.SN)
		if old == nil || old.ExtranetIP == "" {
			extraIP, err := repo.AssignExtranetIP(sett.SN)
			if err != nil {
				return err
			}
			ds.ExtranetIPNetworkID = extraIP.IPNetworkID
			ds.ExtranetIP = extraIP.IP
			log.Infof("Device(%s) is assigned to extranet ip: %s", sett.SN, extraIP.IP)
			ipIDsForRollback = append(ipIDsForRollback, extraIP.ID)
		} else {
			ds.ExtranetIPNetworkID = old.ExtranetIPNetworkID
			ds.ExtranetIP = old.ExtranetIP
		}
	} else {
		// 尝试释放之前已经占用的外网业务IP
		_, _ = repo.ReleaseIP(sett.SN, model.IPScopeExtranet)
	}

	// NeedIntranetIPv6
	if sett.NeedIntranetIPv6 == model.YES {
		// 获取内网业务IPv6
		log.Infof("Device(%s) need intranet ipv6", sett.SN)
		if old == nil || old.IntranetIPv6 == "" {
			intraIPv6, err := repo.AssignIPv6(sett.SN, model.IPScopeIntranet)
			if err != nil {
				return err
			}
			ds.IntranetIPv6NetworkID = intraIPv6.IPNetworkID
			ds.IntranetIPv6 = intraIPv6.IP
			log.Infof("Device(%s) is assigned to intranet ipv6: %s", sett.SN, intraIPv6.IP)
			ipIDsForRollback = append(ipIDsForRollback, intraIPv6.ID)
		} else {
			ds.IntranetIPv6NetworkID = old.IntranetIPv6NetworkID
			ds.IntranetIPv6 = old.IntranetIPv6
		}
	} else {
		// 尝试释放之前已经占用的内网业务IPv6
		log.Infof("Device(%s) no need intranet ipv6", sett.SN)
		_, _ = repo.ReleaseIPv6(sett.SN, model.IPScopeIntranet)
	}

	// NeedExtranetIPv6
	if sett.NeedExtranetIPv6 == model.YES {
		// 获取外网业务IPv6
		log.Infof("Device(%s) need extranet ipv6", sett.SN)
		if old == nil || old.ExtranetIPv6 == "" {
			extraIPv6, err := repo.AssignIPv6(sett.SN, model.IPScopeExtranet)
			if err != nil {
				return err
			}
			ds.ExtranetIPv6NetworkID = extraIPv6.IPNetworkID
			ds.ExtranetIPv6 = extraIPv6.IP
			log.Infof("Device(%s) is assigned to extranet ipv6: %s", sett.SN, extraIPv6.IP)
			ipIDsForRollback = append(ipIDsForRollback, extraIPv6.ID)
		} else {
			ds.ExtranetIPv6NetworkID = old.ExtranetIPv6NetworkID
			ds.ExtranetIPv6 = old.ExtranetIPv6
		}
	} else {
		// 尝试释放之前已经占用的外网业务IPv6
		log.Infof("Device(%s) no need extranet ipv6", sett.SN)
		_, _ = repo.ReleaseIPv6(sett.SN, model.IPScopeExtranet)
	}

	// 保存装机参数，失败时回滚已分配的IP
	if err = repo.SaveDeviceSetting(&ds); err != nil {
		for _, id := range ipIDsForRollback {
			_ = repo.UnassignIP(id)
			log.Debugf("Unassign IP (ID: %v) when failed to save device setting of SN: %s", id, sett.SN)
		}
		return err
	}
	if old != nil && old.ID > 0 {
		_, _ = repo.UpdateDeviceLogType(old.ID, model.DeviceLogInstallType, model.DeviceLogHistoryType)
	}

	log.Infof("Device(%s) settings have been saved. Going to install os.", sett.SN)
	// 进入装机部署
	go startInstallation(log, repo, conf, lim, sett.SN)
	return nil
}

// startInstallation 带外远程开机并开始安装流程
func startInstallation(log logger.Logger, repo model.Repo, conf *config.Config, lim limiter.Limiter, sn string) (err error) {
	// 1、获取进入bootos所需令牌
	var bucket limiter.Bucket
	var token limiter.Token

	if conf.DHCPLimiter.Enable {
		bucket, err = lim.Route(sn)
		if err != nil {
			return err
		}
		token, err = bucket.Acquire(sn, time.Second*time.Duration(conf.DHCPLimiter.WaitingTimeout))
		if err != nil {
			return err
		}
	}

	// 2、修改装机开始时间、装机状态
	now := time.Now()
	if _, err = repo.UpdateDeviceSettingBySN(&model.DeviceSetting{
		SN: sn,
		InstallationStartTime: &now,
		Status:                model.InstallStatusIng,
	}); err != nil {
		if conf.DHCPLimiter.Enable {
			_ = bucket.Return(sn, token)
		}
		return err
	}

	log.Infof("Remote boot and start the OS installation process of the device(SN: %s )", sn)
	var sb strings.Builder
	//TODO 优化项：如果设备在bootos,则跳过重启的过程
	//如果设备已上电，PXE重启，如果关机该命令会有问题，此时再尝试开机
	out, err := OperateOOBPower(log, repo, sn, conf.Crypto.Key, conf.Server.OOBDomain, PowerRestart, true)
	if err != nil {
		log.Errorf("PXE-PowerRestart重启失败(%s)，下一步尝试开机",err.Error())
		sb.WriteString(fmt.Sprintf("\nPXE-PowerRestart重启失败(%s)，下一步尝试开机\n", err.Error()))
		out, err = OperateOOBPower(log, repo, sn, conf.Crypto.Key, conf.Server.OOBDomain, PowerOn, true)
		sb.WriteString(strings.TrimSpace(out))
		if err != nil {
			log.Errorf("PXE-PowerOn开机失败(%s)",err.Error())
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

	url := fmt.Sprintf("%s/api/cloudboot/v1/devices/%s/installations/progress", fmt.Sprintf("http://localhost:%d", conf.Server.HTTPPort), sn)
	log.Infof("POST %s", url)
	log.Infof("Request body: %s", reqBody)

	resp, err := httplib.Post(url).
		Header("Content-Type", "application/json").
		Header("Accept", "application/json").
		Body(reqBody).Response()
	if err != nil {
		log.Error(err)
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Infof("Response body: %s", respBody)
	return err
}

//GetNetworkSettingBySNResp 根据SN请求设备业务网络配置响应结构
type GetNetworkSettingBySNResp struct {
	// IP来源
	IPSource        string `json:"ip_source"`
	BondingRequired string `json:"bonding_required"`
	//分配给该设备的IP对象数组
	Items []model.IPAndIPNetworkUnion `json:"items"`
}

// GetNetworkSettingBySN 返回指定SN的设备业务IP配置信息
func GetNetworkSettingBySN(log logger.Logger, repo model.Repo, sn string) (ipnet *GetNetworkSettingBySNResp, err error) {
	//根据SN获取设备信息
	dev, err := repo.GetDeviceBySN(sn)
	if err != nil {
		return nil, err
	}
	// 根据SN获取装机参数
	ds, err := repo.GetDeviceSettingBySN(sn)
	if err != nil {
		log.Errorf("get device setting fail:%v", err)
		return nil, err
	}
	if ds != nil {
		// 可能有多IP,只分别返回第一个作为配置项
		in := strings.Split(ds.IntranetIP, ",")[0]
		ex := strings.Split(ds.ExtranetIP, ",")[0]
		inv6 := strings.Split(ds.IntranetIPv6, ",")[0]
		exv6 := strings.Split(ds.ExtranetIPv6, ",")[0]
		nwsRet := make([]model.IPAndIPNetworkUnion, 0, 4)

		//通过IP反查网段信息
		if in != "" {
			ip, _ := repo.GetIPByIP(in, model.IPScopeIntranet)
			if ip.IPNetworkID > 0 {
				netSegment, err := repo.GetIPNetworkByID(ip.IPNetworkID)
				if err != nil {
					log.Errorf("query intranet net segment by id:%d fail,%v", ip.IPNetworkID, err)
				} else if netSegment != nil {
					scope := model.IPScopeIntranet
					//TGW 网段网关有2个，逗号分隔，与关联的2台交换机顺序一一对应
					//故获取网关时需要根据设备对应的机位-交换机端口进行顺序获取
					if netSegment.Category == model.TGWIntranet || netSegment.Category == model.TGWExtranet {
						switcher, err := repo.GetIntranetSwitchBySN(sn)
						if err != nil {
							log.Errorf("query intranet switch by dev sn:%s fail:%v", sn, err)
						}
						ipGatewayList := strings.Split(netSegment.Gateway, ",")
						var switchs []string
						_ = json.Unmarshal([]byte(netSegment.Switches), &switchs)
						// 仅处理2个网关对应2台交换机设备的场景
						if len(ipGatewayList) == 2 && len(switchs) == 2 {
							for k := range switchs {
								if switchs[k] == switcher.FixedAssetNumber {
									nwsRet = append(nwsRet, model.IPAndIPNetworkUnion{
										IP:      in,
										Netmask: netSegment.Netmask,
										Gateway: ipGatewayList[k],
										Scope:   &scope,
										Version: model.IPv4,
									})
									break
								}
							}
						} else {
							nwsRet = append(nwsRet, model.IPAndIPNetworkUnion{
								IP:      in,
								Netmask: netSegment.Netmask,
								Gateway: netSegment.Gateway,
								Scope:   &scope,
								Version: model.IPv4,
							})
						}
					} else { // 非TGW网段仅有一个网关IP
						nwsRet = append(nwsRet, model.IPAndIPNetworkUnion{
							IP:      in,
							Netmask: netSegment.Netmask,
							Gateway: netSegment.Gateway,
							Scope:   &scope,
							Version: model.IPv4,
						})
					}
				}
			}
		}
		if ex != "" {
			ip, _ := repo.GetIPByIP(ex, model.IPScopeExtranet)
			if ip.IPNetworkID > 0 {
				netSegment, err := repo.GetIPNetworkByID(ip.IPNetworkID)
				if err != nil {
					log.Errorf("query extranet net segment by id:%d fail,%v", ip.IPNetworkID, err)
				} else if netSegment != nil {
					scope := model.IPScopeExtranet
					if netSegment.Category == model.TGWIntranet || netSegment.Category == model.TGWExtranet {
						switcher, err := repo.GetExtranetSwitchBySN(sn)
						if err != nil {
							log.Errorf("query extranet switch by dev sn:%s fail:%v", sn, err)
						}
						ipGatewayList := strings.Split(netSegment.Gateway, ",")
						var switchs []string
						_ = json.Unmarshal([]byte(netSegment.Switches), &switchs)
						if len(ipGatewayList) == 2 && len(switchs) == 2 {
							for k := range switchs {
								if switchs[k] == switcher.FixedAssetNumber {
									nwsRet = append(nwsRet, model.IPAndIPNetworkUnion{
										IP:      ex,
										Netmask: netSegment.Netmask,
										Gateway: ipGatewayList[k],
										Scope:   &scope,
										Version: model.IPv4,
									})
									break
								}
							}
						} else {
							nwsRet = append(nwsRet, model.IPAndIPNetworkUnion{
								IP:      ex,
								Netmask: netSegment.Netmask,
								Gateway: netSegment.Gateway,
								Scope:   &scope,
								Version: model.IPv4,
							})
						}
					} else { // 非TGW网段仅有一个网关IP
						nwsRet = append(nwsRet, model.IPAndIPNetworkUnion{
							IP:      ex,
							Netmask: netSegment.Netmask,
							Gateway: netSegment.Gateway,
							Scope:   &scope,
							Version: model.IPv4,
						})
					}
				}
			}
		}

		//IPv6
		if inv6 != "" {
			ip, _ := repo.GetIPByIP(inv6, model.IPScopeIntranet)
			if ip.IPNetworkID > 0 {
				netSegment, err := repo.GetIPNetworkByID(ip.IPNetworkID)
				if err != nil {
					log.Errorf("query intranet net segment by id:%d fail,%v", ip.IPNetworkID, err)
				} else if netSegment != nil {
					scope := model.IPScopeIntranet
					nwsRet = append(nwsRet, model.IPAndIPNetworkUnion{
						IP:      inv6,
						Netmask: netSegment.Netmask,
						Gateway: netSegment.Gateway,
						Scope:   &scope,
						Version: model.IPv6,
					})
				}
			}
		}
		if exv6 != "" {
			ip, _ := repo.GetIPByIP(exv6, model.IPScopeExtranet)
			if ip.IPNetworkID > 0 {
				netSegment, err := repo.GetIPNetworkByID(ip.IPNetworkID)
				if err != nil {
					log.Errorf("query extranet net segment by id:%d fail,%v", ip.IPNetworkID, err)
				} else if netSegment != nil {
					scope := model.IPScopeExtranet
					nwsRet = append(nwsRet, model.IPAndIPNetworkUnion{
						IP:      exv6,
						Netmask: netSegment.Netmask,
						Gateway: netSegment.Gateway,
						Scope:   &scope,
						Version: model.IPv6,
					})
				}
			}
		}

		ipnet = &GetNetworkSettingBySNResp{
			IPSource: model.IPSourceStatic,
			Items:    nwsRet,
		}
	}

	//根据机位获取网络区域信息
	if dev.USiteID != nil {
		us, err := repo.GetServerUSiteByID(*dev.USiteID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		if err == gorm.ErrRecordNotFound {
			return ipnet, nil
		}
		if strings.Contains(strings.ToLower(us.PhysicalArea), "bonding") {
			ipnet.BondingRequired = "yes"
		}
		if strings.Contains(strings.ToUpper(us.PhysicalArea), "TGW") {
			ipnet.BondingRequired = "no"
		}
	}
	if ipnet.BondingRequired != "yes" {
		ipnet.BondingRequired = "no"
	}

	return ipnet, nil
}

// GetOSUserSettingBySNReq 查询设备OS用户配置参数请求结构体
type GetOSUserSettingBySNReq struct {
	SN string `json:"sn"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *GetOSUserSettingBySNReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.SN: "sn",
	}
}

// Validate 结构体数据校验
func (reqData *GetOSUserSettingBySNReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())

	sett, err := repo.GetDeviceSettingBySN(reqData.SN)
	if gorm.IsRecordNotFoundError(err) {
		errs.Add([]string{"sn"}, binding.BusinessError, "该物理机无历史部署记录，无法提供配置参数。")
		return errs
	}
	if err != nil {
		errs.Add([]string{"sn"}, binding.SystemError, "系统内部错误")
		return errs
	}

	switch sett.InstallType {
	case model.InstallationPXE:
		if tpl, _ := repo.GetSystemTemplateByID(uint(sett.SystemTemplateID)); tpl == nil || tpl.ID <= 0 {
			errs.Add([]string{"sn"}, binding.BusinessError, "该物理机未指定正确的系统安装模板")
			return errs
		}
	case model.InstallationImage:
		if tpl, _ := repo.GetImageTemplateByID(uint(sett.ImageTemplateID)); tpl == nil || tpl.ID <= 0 {
			errs.Add([]string{"sn"}, binding.BusinessError, "该物理机未指定正确的镜像安装模板")
			return errs
		}
	default:
		errs.Add([]string{"sn"}, binding.BusinessError, "该物理机未指定正确的安装方式。")
		return errs
	}
	return errs
}

// GetOSUserSettingBySN 返回指定设备的操作系统用户配置参数
func GetOSUserSettingBySN(log logger.Logger, repo model.Repo, reqData *GetOSUserSettingBySNReq) (*setting.OSUserSettingItem, error) {
	// 从装机关联的系统安装模板/镜像安装模板中获取用户名、密码
	sett, err := repo.GetDeviceSettingBySN(reqData.SN)
	if err != nil {
		return nil, err
	}

	switch sett.InstallType {
	case model.InstallationPXE:
		tpl, err := repo.GetSystemTemplateByID(uint(sett.SystemTemplateID))
		if err != nil {
			return nil, err
		}
		return &setting.OSUserSettingItem{
			Username: tpl.Username,
			Password: tpl.Password,
		}, nil

	case model.InstallationImage:
		tpl, err := repo.GetImageTemplateByID(uint(sett.ImageTemplateID))
		if err != nil {
			return nil, err
		}
		return &setting.OSUserSettingItem{
			Username: tpl.Username,
			Password: tpl.Password,
		}, nil
	}
	log.Warnf("Invalid installation type: %s=%s", reqData.SN, sett.InstallType)

	return new(setting.OSUserSettingItem), nil
}

// ReinstallsReq 批量重装请求参数
type ReinstallsReq struct {
	SNs []string `json:"sns"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *ReinstallsReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.SNs: "sns",
	}
}

// Validate 结构体数据校验
func (reqData *ReinstallsReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())
	for i := range reqData.SNs {
		_, err := repo.GetDeviceBySN(reqData.SNs[i])
		if gorm.IsRecordNotFoundError(err) {
			errs.Add([]string{"sns"}, binding.BusinessError, fmt.Sprintf("物理机(%s)信息不存在", reqData.SNs[i]))
			return errs
		}
		if err != nil {
			errs.Add([]string{"sns"}, binding.SystemError, "系统内部错误")
			return errs
		}

		sett, err := repo.GetDeviceSettingBySN(reqData.SNs[i])
		if gorm.IsRecordNotFoundError(err) {
			errs.Add([]string{"sns"}, binding.BusinessError, fmt.Sprintf("物理机(%s)无历史部署记录，请重新申请上架部署。", reqData.SNs[i]))
			return errs
		}
		if err != nil {
			errs.Add([]string{"sns"}, binding.SystemError, "系统内部错误")
			return errs
		}
		if sett.Status == model.InstallStatusIng {
			errs.Add([]string{"sns"}, binding.BusinessError, fmt.Sprintf("物理机(%s)当前正在上架部署，请稍后再试。", reqData.SNs[i]))
			return errs
		}

		if tor, _ := repo.GetTORBySN(reqData.SNs[i]); tor == "" {
			errs.Add([]string{"sn"}, binding.BusinessError, fmt.Sprintf("物理机(%s)无法找到TOR", reqData.SNs[i]))
			return errs
		}
	}
	return errs
}

// Reinstalls 批量重装指定设备
func Reinstalls(log logger.Logger, repo model.Repo, conf *config.Config, lim limiter.Limiter, reqData *ReinstallsReq) (err error) {
	for i := range reqData.SNs {
		devSett, err := repo.GetDeviceSettingBySN(reqData.SNs[i])
		if err != nil {
			return err
		}
		if _, err = repo.UpdateInstallStatusAndProgressByID(devSett.ID, model.InstallStatusPre, 0.0); err != nil {
			return err
		}

		if centos6.IsPXEUEFI(log, repo, reqData.SNs[i]) {
			_ = centos6.DropConfigurations(log, repo, reqData.SNs[i]) // TODO 为支持centos6的UEFI方式安装而临时增加的逻辑，后续会删除。
		}

		if _, err = repo.UpdateDeviceLogType(devSett.ID, model.DeviceLogInstallType, model.DeviceLogHistoryType); err != nil {
			return err
		}

		go startInstallation(log, repo, conf, lim, devSett.SN)
	}
	return nil
}

// CancelInstallsReq 批量取消安装请求结构体
type CancelInstallsReq struct {
	SNs []string `json:"sns"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *CancelInstallsReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.SNs: "sns",
	}
}

// Validate 结构体数据校验
func (reqData *CancelInstallsReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())

	for i := range reqData.SNs {
		_, err := repo.GetDeviceBySN(reqData.SNs[i])
		if gorm.IsRecordNotFoundError(err) {
			errs.Add([]string{"sns"}, binding.BusinessError, fmt.Sprintf("物理机(%s)信息不存在", reqData.SNs[i]))
			return errs
		}
		if err != nil {
			errs.Add([]string{"sns"}, binding.SystemError, "系统内部错误")
			return errs
		}

		_, err = repo.GetDeviceSettingBySN(reqData.SNs[i])
		if gorm.IsRecordNotFoundError(err) {
			errs.Add([]string{"sns"}, binding.BusinessError, fmt.Sprintf("物理机(%s)并未部署，无法取消。", reqData.SNs[i]))
			return errs
		}
		if err != nil {
			errs.Add([]string{"sns"}, binding.SystemError, "系统内部错误")
			return errs
		}
	}
	return errs
}

// CancelInstalls 批量取消安装指定设备
func CancelInstalls(log logger.Logger, repo model.Repo, reqData *CancelInstallsReq) (err error) {
	for i := range reqData.SNs {
		devSett, err := repo.GetDeviceSettingBySN(reqData.SNs[i])
		if err != nil {
			return err
		}
		if devSett.Status == model.InstallStatusSucc || devSett.Status == model.InstallStatusFail {
			continue
		}

		// 重置设备SN的装机参数，进度设置为0，装机状态设置为失败
		if _, err = repo.UpdateInstallStatusAndProgressByID(devSett.ID, model.InstallStatusFail, 0.0); err != nil {
			return err
		}

		if centos6.IsPXEUEFI(log, repo, reqData.SNs[i]) {
			_ = centos6.DropConfigurations(log, repo, reqData.SNs[i]) // TODO 为支持centos6的UEFI方式安装而临时增加的逻辑，后续会删除。
		}

		//if _, err = repo.UpdateDeviceLogType(devSett.ID, model.DeviceLogInstallType, model.DeviceLogHistoryType); err != nil {
		//	return err
		//}
	}
	return nil
}

// RemoveDeviceSettingsReq 批量删除设备装机参数请求结构体
type RemoveDeviceSettingsReq struct {
	SNs []string `json:"sns"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *RemoveDeviceSettingsReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.SNs: "sns",
	}
}

// Validate 结构体数据校验
func (reqData *RemoveDeviceSettingsReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())

	for i := range reqData.SNs {
		_, err := repo.GetDeviceSettingBySN(reqData.SNs[i])
		if gorm.IsRecordNotFoundError(err) {
			errs.Add([]string{"sns"}, binding.BusinessError, fmt.Sprintf("物理机(%s)上架部署记录不存在，请重试。", reqData.SNs[i]))
			return errs
		}
		if err != nil {
			errs.Add([]string{"sns"}, binding.SystemError, "系统内部错误")
			return errs
		}
	}
	return errs
}

// RemoveDeviceSettings 批量删除设备装机参数
func RemoveDeviceSettings(log logger.Logger, repo model.Repo, reqData *RemoveDeviceSettingsReq) (err error) {
	for i := range reqData.SNs {
		sett, err := repo.GetDeviceSettingBySN(reqData.SNs[i])
		if gorm.IsRecordNotFoundError(err) {
			continue
		}
		if err != nil {
			return err
		}

		//释放IP资源
		if err = repo.UnassignIPsBySN(sett.SN); err != nil {
			log.Errorf("release ip for sn:%s fail", sett.SN)
			return err
		}

		if _, err = repo.DeleteDeviceSettingByID(sett.ID); err != nil {
			return nil
		}

		if centos6.IsPXEUEFI(log, repo, reqData.SNs[i]) {
			_ = centos6.DropConfigurations(log, repo, reqData.SNs[i]) // TODO 为支持centos6的UEFI方式安装而临时增加的逻辑，后续会删除。
		}

		if _, err = repo.UpdateDeviceLogType(sett.ID, model.DeviceLogInstallType, model.DeviceLogHistoryType); err != nil {
			return err
		}
	}
	return nil
}

// GetDeviceSettingPageReq 装机设备查询条件列表
type GetDeviceSettingPageReq struct {
	// 源节点名
	OriginNode string `json:"origin_node"`
	// 所属数据中心ID
	IDCID string `json:"idc_id"`
	// 所属机房ID
	ServerRoomID string `json:"server_room_id"`
	// 所属机架ID
	ServerCabinetID string `json:"server_cabinet_id"`
	// 所属机位ID
	ServerUsiteID string `json:"server_usite_id"`
	// 机房名
	ServerRoomName string `json:"server_room_name"`
	// 机架编号
	ServerCabinetNumber string `json:"server_cabinet_number"`
	// 设备序列号(多个SN用英文逗号分隔)
	Sn string `json:"sn"`
	// 固资号
	FN string `json:"fixed_asset_number"`
	// 设备类型
	Category string `json:"category"`	
	// 内外IP
	IntranetIP string `json:"intranet_ip"`
	// 外网IP
	ExtranetIP string `json:"extranet_ip"`
	// 硬件配置模板ID
	HardwareTemplateID uint `json:"hardware_template_id"`
	// 镜像配置模板ID
	ImageTemplateID uint `json:"image_template_id"`
	// 装机状态。可选值  pre_install-等待安装; installing-正在安装; failure-安装失败; success-安装成功;
	Status string `json:"status"`
	// 分页页号
	Page int64 `json:"page"`
	// 分页大小。默认值:10。阈值 100
	PageSize int64 `json:"page_size"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *GetDeviceSettingPageReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDCID:               "idc_id",
		&reqData.ServerRoomID:        "server_room_id",
		&reqData.ServerCabinetID:     "server_cabinet_id",
		&reqData.ServerUsiteID:       "server_usite_id",
		&reqData.Sn:                  "sn",
		&reqData.FN:                  "fixed_asset_number",
		&reqData.Category:            "category",		
		&reqData.ServerRoomName:      "server_room_name",
		&reqData.ServerCabinetNumber: "server_cabinet_number",
		&reqData.IntranetIP:          "intranet_ip",
		&reqData.ExtranetIP:          "extranet_ip",
		&reqData.HardwareTemplateID:  "hardware_template_id",
		&reqData.ImageTemplateID:     "image_template_id",
		&reqData.Status:              "status",
		&reqData.Page:                "page",
		&reqData.PageSize:            "page_size",
	}
}

// GetDeviceSettingPage 按条件查询机位分页列表
func GetDeviceSettingPage(log logger.Logger, repo model.Repo, reqData *GetDeviceSettingPageReq) (pg *page.Page, err error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 100 {
		reqData.PageSize = 10
	}
	if reqData.Page < 0 {
		reqData.Page = 0
	}

	cond := model.CombineDeviceSetting{
		IDCID:               mystrings.Multi2UintSlice(reqData.IDCID),
		ServerRoomID:        mystrings.Multi2UintSlice(reqData.ServerRoomID),
		ServerCabinetID:     mystrings.Multi2UintSlice(reqData.ServerCabinetID),
		ServerUsiteID:       mystrings.Multi2UintSlice(reqData.ServerUsiteID),
		Sn:                  reqData.Sn,
		FN:                  reqData.FN,
		Category:            reqData.Category,		
		ServerRoomName:      reqData.ServerRoomName,
		ServerCabinetNumber: reqData.ServerCabinetNumber,
		ExtranetIP:          reqData.ExtranetIP,
		IntranetIP:          reqData.IntranetIP,
		HardwareTemplateID:  reqData.HardwareTemplateID,
		ImageTemplateID:     reqData.ImageTemplateID,
		Status:              reqData.Status,
	}

	totalRecords, err := repo.CountDeviceSettingCombines(&cond)
	if err != nil {
		return nil, err
	}

	pager := page.NewPager(reflect.TypeOf(&DeviceSettingResp{}), reqData.Page, reqData.PageSize, totalRecords)
	orderBy := []model.OrderByPair{{
		Name:      "updated_at",
		Direction: model.DESC,
	}}
	items, err := repo.GetDeviceSettingCombinesByCond(&cond, orderBy, pager.BuildLimiter())
	if err != nil {
		return nil, err
	}

	for i := range items {
		resp, err := convert2DeviceSettingResp(repo, log, items[i])
		if err != nil {
			continue
		}
		if resp == nil {
			continue
		}
		pager.AddRecords(resp)
	}
	return pager.BuildPage(), nil
}

func convert2DeviceSettingResp(repo model.Repo, log logger.Logger, item *model.DeviceSetting) (*DeviceSettingResp, error) {
	deviceSetting := DeviceSettingResp{
		ID:              item.ID,
		Sn:              item.SN,
		Status:          item.Status,
		InstallType:     item.InstallType,
		InstallProgress: item.InstallProgress,
		ExtranetIP:      item.ExtranetIP,
		IntranetIP:      item.IntranetIP,
		UpdatedAt:       item.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	device, _ := repo.GetDeviceBySN(item.SN)
	if device == nil {
		log.Errorf("can not get device, sn: %s", item.SN)
		return nil, fmt.Errorf("can not get device, sn: %s", item.SN)
	}

	deviceSetting.Device.ID = device.ID
	deviceSetting.Device.SN = device.SN
	deviceSetting.Device.FixedAssetNumber = device.FixedAssetNumber
	deviceSetting.Device.Category = device.Category
	deviceSetting.Sn = device.SN

	if idc, _ := repo.GetIDCByID(device.IDCID); idc != nil {
		deviceSetting.IDC.ID = idc.ID
		deviceSetting.IDC.Name = idc.Name
	}

	if room, _ := repo.GetServerRoomByID(device.ServerRoomID); room != nil {
		deviceSetting.ServerRoom.ID = room.ID
		deviceSetting.ServerRoom.Name = room.Name
	}

	if cabinet, _ := repo.GetServerCabinetByID(device.CabinetID); cabinet != nil {
		deviceSetting.ServerCabinet.ID = cabinet.ID
		deviceSetting.ServerCabinet.Number = cabinet.Number
	}

	if device.USiteID != nil {
		if usite, _ := repo.GetServerUSiteByID(*device.USiteID); usite != nil {
			deviceSetting.ServerUSite.ID = usite.ID
			deviceSetting.ServerUSite.Number = usite.Number
			deviceSetting.ServerUSite.PhysicalArea = usite.PhysicalArea
		}
	}
	if imageTemplate, _ := repo.GetImageTemplateByID(item.ImageTemplateID); imageTemplate != nil {
		deviceSetting.ImageTemplate.ID = imageTemplate.ID
		deviceSetting.ImageTemplate.Name = imageTemplate.Name
	}
	if systemTemplate, _ := repo.GetSystemTemplateByID(item.SystemTemplateID); systemTemplate != nil {
		deviceSetting.SystemTemplate.ID = systemTemplate.ID
		deviceSetting.SystemTemplate.Name = systemTemplate.Name
	}
	if hardwareTemplate, _ := repo.GetHardwareTemplateByID(item.HardwareTemplateID); hardwareTemplate != nil {
		deviceSetting.HardwareTemplate.ID = hardwareTemplate.ID
		deviceSetting.HardwareTemplate.Name = hardwareTemplate.Name
	}

	deviceSetting.DHCPToken, _ = repo.GetTokenBySN(item.SN)
	deviceSetting.TOR, _ = repo.GetTORBySN(item.SN)

	return &deviceSetting, nil
}

// DeviceInstallStaticCountResp 装机信息统计resp
type DeviceInstallStaticCountResp struct {
	// 物理机总数
	TotalDevices int64 `json:"total_devices"`
	// 等待安装物理机数量
	PreInstallCount int64 `json:"preinstall_count"`
	// 正在安装的物理机数量
	InstallingCount int64 `json:"installing_count"`
	// 安装失败的物理机数量
	FailureCount int64 `json:"failure_count"`
	// 安装成功的物理机数量
	SuccessCount int64 `json:"success_count"`
}

// CountDeviceInstallStatic 装机信息统计
// 统计物理机总数、 等待安装物理机数量、 正在安装的物理机数量、安装失败的物理机数量、安装成功的物理机数量
func CountDeviceInstallStatic(log logger.Logger, repo model.Repo) (*DeviceInstallStaticCountResp, error) {

	resp := DeviceInstallStaticCountResp{}

	// 统计物理机总数
	totalDevices, err := repo.CountDevices(&model.Device{})
	if err != nil {
		return nil, err
	}

	// 等待安装物理机数量
	preInstall, err := repo.CountDeviceSettingByStatus(model.InstallStatusPre)
	if err != nil {
		return nil, err

	}

	// 正在安装的物理机数量
	installing, err := repo.CountDeviceSettingByStatus(model.InstallStatusIng)
	if err != nil {
		return nil, err
	}

	// 安装失败的物理机数量
	failure, err := repo.CountDeviceSettingByStatus(model.InstallStatusFail)
	if err != nil {
		return nil, err
	}

	// 安装成功的物理机数量
	success, err := repo.CountDeviceSettingByStatus(model.InstallStatusSucc)
	if err != nil {
		return nil, err
	}

	resp.TotalDevices = totalDevices
	resp.PreInstallCount = preInstall
	resp.InstallingCount = installing
	resp.FailureCount = failure
	resp.SuccessCount = success

	return &resp, nil
}

// DeviceSettingResp 装机参数返回体
type DeviceSettingResp struct {
	// 数据中心
	IDC struct {
		// 数据中心ID
		ID uint `json:"id"`
		// 数据中心名称
		Name string `json:"name"`
	} `json:"idc"`

	// 机房
	ServerRoom struct {
		// 机房ID
		ID uint `json:"id"`
		// 机房名称
		Name string `json:"name"`
	} `json:"server_room"`

	// 机架
	ServerCabinet struct {
		// 机架ID
		ID uint `json:"id"`
		// 机架编号
		Number string `json:"number"`
	} `json:"server_cabinet"`

	// 机位
	ServerUSite struct {
		// 机位ID
		ID uint `json:"id"`
		// 机位编号
		Number string `json:"number"`
		// 物理区域
		PhysicalArea string `json:"physical_area"`
	} `json:"server_usite"`

	// 物理机
	Device struct {
		// 设备ID
		ID uint `json:"id"`
		// 设备SN
		SN string `json:"sn"`
		// 设备SN
		FixedAssetNumber string `json:"fixed_asset_number"`
		// 设备类型
		Category         string    `gorm:"column:category"`	
	} `json:"device"`

	// 硬件模板
	HardwareTemplate struct {
		// 硬件模板ID
		ID uint `json:"id"`
		// 硬件模板名称
		Name string `json:"name"`
	} `json:"hardware_template"`

	// 镜像模板
	ImageTemplate struct {
		// 镜像模板ID
		ID uint `json:"id"`
		// 镜像模板名称
		Name string `json:"name"`
	} `json:"image_template"`

	// 系统模板
	SystemTemplate struct {
		// 系统模板ID
		ID uint `json:"id"`
		// 系统模板名称
		Name string `json:"name"`
	} `json:"system_template"`

	// TOR tor
	TOR string `json:"tor"`
	// 设备参数主键
	ID uint `json:"id"`
	// 安装进度
	InstallProgress float64 `json:"install_progress"`
	// 安装类型： image/pxe
	InstallType string `json:"install_type"`
	// 设备序列号(多个SN用英文逗号分隔)
	Sn string `json:"sn"`
	// 装机状态。可选值  pre_install-等待安装; installing-正在安装; failure-安装失败; success-安装成功;
	Status string `json:"status"`
	// 外网IP
	ExtranetIP string `json:"extranet_ip"`
	// 内网IP
	IntranetIP string `json:"intranet_ip"`
	// DHCP IP令牌
	DHCPToken string `json:"dhcp_token"`
	// UpdatedAt 修改时间
	UpdatedAt string `json:"updated_at"`
}

// DeviceSetting 设备装机参数结构体
type DeviceSetting struct {
	// 设备序列号
	SN string `json:"sn"`
	// InstallType 安装方式
	// Enum: image,pxe
	InstallType string `json:"install_type"`
	// 操作系统安装模板
	OSTemplate *setting.OSTemplateSetting `json:"os_template"`
	// 硬件配置模板
	HardwareTemplate *setting.HardwareTemplateSetting `json:"hardware_template"`
	// NeedExtranetIP 是否需要外网IP。可选值：yes-是; no-否;
	// Enum: yes,no
	NeedExtranetIP string `json:"need_extranet_ip"`
	// 内网业务IP
	IntranetIP struct {
		IP      string `json:"ip"`
		Netmask string `json:"netmask"`
		Gateway string `json:"gateway"`
	} `json:"intranet_ip"`
	// 外网业务IP
	ExtranetIP struct {
		IP      string `json:"ip"`
		Netmask string `json:"netmask"`
		Gateway string `json:"gateway"`
	} `json:"extranet_ip"`
}

// GetDeviceSettingBySN 查询指定设备的装机参数
func GetDeviceSettingBySN(log logger.Logger, repo model.Repo, sn string) (*DeviceSetting, error) {
	sett, err := repo.GetDeviceSettingBySN(sn)
	if err != nil {
		return nil, err
	}

	// 获取硬件模板
	var hwTemplateID int
	var hwTemplateName string

	if hwTemplate, _ := repo.GetHardwareTemplateByID(sett.HardwareTemplateID); hwTemplate != nil {
		hwTemplateID, hwTemplateName = int(hwTemplate.ID), hwTemplate.Name
	}

	// 获取操作系统安装模板ID
	var osID int
	var osFamily, osName, osBootMode string

	switch sett.InstallType {
	case model.InstallationImage:
		if osTemplate, _ := repo.GetImageTemplateByID(sett.ImageTemplateID); osTemplate != nil {
			osID, osFamily, osName, osBootMode = int(osTemplate.ID), osTemplate.Family, osTemplate.Name, osTemplate.BootMode
		}

	case model.InstallationPXE:
		if osTemplate, _ := repo.GetSystemTemplateByID(sett.SystemTemplateID); osTemplate != nil {
			osID, osFamily, osName, osBootMode = int(osTemplate.ID), osTemplate.Family, osTemplate.Name, osTemplate.BootMode
		}
	}

	var intraNetmask, intraGateway string
	if ipnet, _ := repo.GetIPNetworkByID(sett.IntranetIPNetworkID); ipnet != nil {
		intraNetmask, intraGateway = ipnet.Netmask, ipnet.Gateway
	}

	var extraNetmask, extraGateway string
	if ipnet, _ := repo.GetIPNetworkByID(sett.ExtranetIPNetworkID); ipnet != nil {
		extraNetmask, extraGateway = ipnet.Netmask, ipnet.Gateway
	}

	ds := DeviceSetting{
		SN:          sett.SN,
		InstallType: sett.InstallType,
		OSTemplate: &setting.OSTemplateSetting{
			ID:       osID,
			Family:   osFamily,
			Name:     osName,
			BootMode: osBootMode,
		},
		HardwareTemplate: &setting.HardwareTemplateSetting{
			ID:   hwTemplateID,
			Name: hwTemplateName,
		},
		NeedExtranetIP: sett.NeedExtranetIP,
	}
	ds.IntranetIP.IP = sett.IntranetIP
	ds.IntranetIP.Netmask = intraNetmask
	ds.IntranetIP.Gateway = intraGateway
	ds.ExtranetIP.IP = sett.ExtranetIP
	ds.ExtranetIP.Netmask = extraNetmask
	ds.ExtranetIP.Gateway = extraGateway

	return &ds, nil
}

// ImageTemplate 镜像安装模板
type ImageTemplate struct {
	ID         uint          `json:"id"`
	Family     string        `json:"family"`      // 操作系统族系
	Name       string        `json:"name"`        // 模板名
	BootMode   string        `json:"boot_mode"`   // 启动模式
	URL        string        `json:"url"`         // PXE 引导模板内容
	Username   string        `json:"username"`    // 操作系统用户名
	Password   string        `json:"password"`    // 操作系统用户密码
	PreScript  string        `json:"pre_script"`  // 前置脚本
	PostScript string        `json:"post_script"` // 后置脚本
	CreatedAt  times.ISOTime `json:"created_at"`  // 创建时间
	UpdatedAt  times.ISOTime `json:"updated_at"`  // 更新时间
	Disks      []struct {
		Name       string `json:"name"`
		Partitions []struct {
			Name       string `json:"name"`
			Size       string `json:"size"`
			Fstype     string `json:"fstype"`
			Mountpoint string `json:"mountpoint"`
		} `json:"partitions"`
	} `json:"disks"`
}

// GetImageTemplateBySN 根据SN查询设备的系统模板信息
func GetImageTemplateBySN(log logger.Logger, repo model.Repo, sn string) (*ImageTemplate, error) {
	mod, err := repo.GetImageTemplateBySN(sn)
	if err != nil {
		return nil, err
	}
	mod.Partition = strings.Replace(mod.Partition, "\\", "", -1)
	tpl := ImageTemplate{
		ID:         mod.ID,
		CreatedAt:  times.ISOTime(mod.CreatedAt),
		UpdatedAt:  times.ISOTime(mod.UpdatedAt),
		Name:       mod.Name,
		URL:        mod.ImageURL,
		Family:     mod.Family,
		BootMode:   mod.BootMode,
		Username:   mod.Username,
		Password:   mod.Password,
		PreScript:  base64.StdEncoding.EncodeToString([]byte(mystrings.DOS2UNIX(mod.PreScript))),
		PostScript: base64.StdEncoding.EncodeToString([]byte(mystrings.DOS2UNIX(mod.PostScript))),
	}
	if err := json.Unmarshal([]byte(mod.Partition), &tpl.Disks); err != nil {
		log.Error(err)
		return nil, err
	}
	return &tpl, nil
}

// AutoReinstallsReq 批量重装请求参数
type AutoReinstallsReq struct {
	SNs                []string  `json:"sns"`
	RecycleReinstall   string    `json:"recycle_reinstall"`  // yes or no 影响重装后设备的运营状态
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *AutoReinstallsReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.SNs: "sns",
		&reqData.RecycleReinstall: "recycle_reinstall",
	}
}

// Validate 结构体数据校验
func (reqData *AutoReinstallsReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())
	for i := range reqData.SNs {
		// 设备信息校验
		dev, err := repo.GetDeviceBySN(reqData.SNs[i])
		if gorm.IsRecordNotFoundError(err) {
			errs.Add([]string{"sns"}, binding.BusinessError, fmt.Sprintf("Device of SN:%s not exist", reqData.SNs[i]))
			return errs
		}
		// 回收重装场景校验设备运营状态=待部署
		if reqData.RecycleReinstall == "yes" {
			if dev.OperationStatus != model.DevOperStatPreDeploy {
				errs.Add([]string{"sns"}, binding.BusinessError, fmt.Sprintf("Status of SN:%s is not predeploy", reqData.SNs[i]))
				return errs
			}
		}
		if tor, _ := repo.GetTORBySN(reqData.SNs[i]); tor == "" {
			errs.Add([]string{"sn"}, binding.BusinessError, fmt.Sprintf("SN:%s can not find TOR", reqData.SNs[i]))
			return errs
		}
		if err != nil {
			errs.Add([]string{"sns"}, binding.SystemError, "System error")
			return errs
		}		
		
		// 装机参数记录校验
		sett, err := repo.GetDeviceSettingBySN(reqData.SNs[i])
		if gorm.IsRecordNotFoundError(err) {
			continue
		}
		if sett.Status == model.InstallStatusIng || sett.Status == model.InstallStatusPre {
			errs.Add([]string{"sns"}, binding.BusinessError, fmt.Sprintf("Device SN:%s is installing or pre install", reqData.SNs[i]))
			return errs
		}
		if err != nil {
			errs.Add([]string{"sns"}, binding.SystemError, "Failed to get device setting")
			return errs
		}		
	}
	return errs
}

// AutoReinstalls 调用规则引擎自动生成设备的装机参数，并发起部署
func AutoReinstalls(log logger.Logger, repo model.Repo, conf *config.Config, lim limiter.Limiter, reqData *AutoReinstallsReq) (succeeds []string, err error) {
	log.Infof("Start auto reinstall generating device %v settings", reqData.SNs)
	defer log.Infof("End auto reinstall saving device %v settings", reqData.SNs)

	for i := range reqData.SNs {
		// 获取设备信息
		dev, err := repo.GetDeviceBySN(reqData.SNs[i])
		if err != nil {
			log.Errorf("Fail to get device of sn: %s , error: %s", reqData.SNs[i], err)
			continue
		}
		// 装机参数对象
		ds := model.DeviceSetting{
			SN:             dev.SN,
			InstallType:    model.InstallationImage,
			Status:         model.InstallStatusPre,
			NeedIntranetIPv6: model.NO,				//默认仅分配IPV4
			NeedExtranetIPv6: model.NO,				//默认仅分配IPV4
		}		
		// 查询装机记录
		old, err := repo.GetDeviceSettingBySN(dev.SN)
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Errorf("Fail to get exist devicesetting of device(sn: %s ), error: %s", dev.SN, err)
			continue
		}
		if old != nil {
			ds.ID = old.ID
			ds.CreatedAt = old.CreatedAt
			ds.Updater = "auto_reinstall_api"
		} else {
			ds.Creator = "auto_reinstall_api"
		}
		// 根据 USiteID 查询获取 PhysicalArea
		var ServerUSite *model.ServerUSite
		ServerUSite, err = repo.GetServerUSiteByID(*dev.USiteID)
		if err != nil {
			log.Errorf("Fail to get serverusite of device(sn: %s ), error: %s", dev.SN, err)
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
		items, err := repo.GetDeviceCategorys(&cond, nil, nil)
		if err != nil {
			log.Errorf("Fail to get category of device(sn: %s ), error: %s", dev.SN, err)
			continue
		}
		isFITIEcoProduct := "unknown"
		if len(items) != 0 {
			isFITIEcoProduct = items[0].IsFITIEcoProduct
			log.Debugf("device[SN:%s] category[%s] is FITIEcoProduct", dev.SN, dev.Category)
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
		result, err := RuleInduction(log, repo, &df)
		if err != nil {
			log.Errorf("Fail to get rule action results of device(sn: %s ), error: %s", dev.SN, err)
			continue
		}		
		// 获取硬件配置模板ID
		hwTemplate, err := repo.GetHardwareTemplateByName(result.RaidResult)
		if err != nil {
			log.Errorf("Fail to GetHardwareTemplateByName of device(sn: %s name: %s ), error: %s", dev.SN, result.RaidResult, err)
			continue
		}
		ds.HardwareTemplateID = hwTemplate.ID
		// 获取操作系统模版ID
		// 优先image 次之pxe
		imgTemplate, err := repo.GetImageTemplateByName(result.OSResult)
		if err != nil {
			log.Errorf("Fail to GetImageTemplateByName of device(sn: %s name: %s ), error: %s ,try SystemTemplate", dev.SN, result.OSResult, err)
			//continue
			sysTemplate, err := repo.GetSystemTemplateByName(result.OSResult)
			if err != nil {
				log.Errorf("Fail to GetSystemTemplateByName of device(sn: %s name: %s ), error: %s", dev.SN, result.OSResult, err)
				continue
			}
			ds.SystemTemplateID = sysTemplate.ID
			ds.InstallType = model.InstallationPXE
		}
		ds.ImageTemplateID = imgTemplate.ID
		// 获取网络配置
		if result.NetworkResult == "la_wa" {
			// 获取内网业务IP
			intraIP, err := repo.AssignIntranetIP(dev.SN)
			if err != nil {
				log.Errorf("AutoGenDeviceSetting fail to AssignIntranetIP of device(sn: %s ), error: %s", dev.SN, err)
				continue
			}
			ds.IntranetIPNetworkID = intraIP.IPNetworkID
			ds.IntranetIP = intraIP.IP
			log.Infof("Device(%s) is assigned to intranet ip: %s", dev.SN, intraIP.IP)
			// 获取外网业务IP
			extraIP, err := repo.AssignExtranetIP(dev.SN)
			if err != nil {
				log.Errorf("AutoGenDeviceSetting fail to AssignExtranetIP of device(sn: %s ), error: %s", dev.SN, err)
				continue
			}
			ds.ExtranetIPNetworkID = extraIP.IPNetworkID
			ds.ExtranetIP = extraIP.IP
			ds.NeedExtranetIP = "yes"
			log.Infof("Device(%s) is assigned to extranet ip: %s", dev.SN, extraIP.IP)
		} else if result.NetworkResult == "la" {
			// 获取内网业务IP
			intraIP, err := repo.AssignIntranetIP(dev.SN)
			if err != nil {
				log.Errorf("AutoGenDeviceSetting fail to AssignIntranetIP of device(sn: %s ), error: %s", dev.SN, err)
				continue
			}
			ds.IntranetIPNetworkID = intraIP.IPNetworkID
			ds.IntranetIP = intraIP.IP
			ds.NeedExtranetIP = "no"
			log.Infof("Device(%s) is assigned to intranet ip: %s", dev.SN, intraIP.IP)
		} else {
			// 网络配置规则推理结果不符合la la_wa 时跳过
			log.Errorf("Fail to AssignIP for device(sn: %s  expected network_conf: %s ), error: %s", dev.SN, result.NetworkResult, err)
			continue			
		}

		// 当RAID、OS、NETWORK均得到对应参数时，保存装机参数
		if err = repo.SaveDeviceSetting(&ds); err != nil {
			log.Errorf("Device(%s) settings(auto) have not been saved err: %s", dev.SN, err)
			return succeeds, err
		}

		// 更新保存旧的装机日志
		if old != nil && old.ID > 0 {
			_, _ = repo.UpdateDeviceLogType(old.ID, model.DeviceLogInstallType, model.DeviceLogHistoryType)
		}		
		log.Infof("Device(%s) settings(auto) have been saved", dev.SN)

		// 更新设备状态=重装中
		mod := &model.Device{
			SN:              reqData.SNs[i],
			OperationStatus: model.DevOperStatReinstalling,
		}
		// remark字段标记，初始的状态，物理机重装场景恢复运营状态
		if reqData.RecycleReinstall == "no" {
			mod.Remark = dev.OperationStatus
		}
		_, err = repo.UpdateDeviceBySN(mod)
		if err != nil {
			log.Errorf("Failt to update status of device %s to reinstalling", reqData.SNs[i])
		}
		succeeds = append(succeeds, reqData.SNs[i])
		// 开始部署
		go startInstallation(log, repo, conf, lim, dev.SN)
	}
	return succeeds, nil
}

// SaveDeviceSettingsAndReinstalls 批量保存设备装机参数并进入系统部署，完成部署后恢复初始的运营状态
func SaveDeviceSettingsAndReinstalls(log logger.Logger, repo model.Repo, conf *config.Config, lim limiter.Limiter, reqData *SaveDeviceSettingsReq) (succeeds []string, err error) {
	for i := range reqData.Settings {
		// 获取设备信息
		dev, err := repo.GetDeviceBySN(reqData.Settings[i].SN)
		if err != nil {
			log.Errorf("Fail to get device of sn: %s , error: %s", reqData.Settings[i].SN, err)
			continue
		}
		// 保存装机参数并进入系统部署
		if err = saveDeviceSetting(log, repo, conf, lim, reqData.Settings[i], reqData.CurrentUser.LoginName); err != nil {
			return succeeds, err
		}

		// 更新设备状态=重装中
		mod := &model.Device{
			SN:              reqData.Settings[i].SN,
			OperationStatus: model.DevOperStatReinstalling,
			Remark: 		 dev.OperationStatus,			// remark字段临时记录初始的状态，重装场景恢复运营状态
		}
		_, err = repo.UpdateDeviceBySN(mod)
		if err != nil {
			log.Errorf("Failt to update status of device %s to reinstalling", reqData.Settings[i].SN)
		}
		succeeds = append(succeeds, reqData.Settings[i].SN)
	}
	return succeeds, nil
}


// UpdateDeviceSettingReq
type UpdateDeviceSettingReq struct {
	SN 						string		`json:"sn"`
	OSSystemTemplateName	string		`json:"os"`
}


// UpdateDeviceSetting 更新装机参数。
// 支持裸金属设备、特殊设备等自定义操作系统名称（优先显示system_template_name,后显示image_template_name）
func UpdateDeviceSetting(log logger.Logger, repo model.Repo, reqData *UpdateDeviceSettingReq)  error {
	if reqData.OSSystemTemplateName != "" && reqData.SN != "" {
		log.Debugf("begin to update device setting (system template name:%s for SN:%s)", reqData.OSSystemTemplateName ,reqData.SN)
		// 操作系统名称校验,不存在则新增一条记录
		sysTpl, err := repo.GetSystemTemplateByName(reqData.OSSystemTemplateName)
		if err != nil {
			log.Errorf("get system template by name %s failed(%v) ,adding one", reqData.OSSystemTemplateName, err)
			tpl := model.SystemTemplate{
				Family:   	"Custom",
				BootMode: 	"uefi",
				Name:     	reqData.OSSystemTemplateName,
				PXE:        "#NULL",
				Content:    "#NULL",
				OSLifecycle: model.OSTesting,
				Arch:		 model.OSARCHUNKNOWN,		
			}
			_, err := repo.SaveSystemTemplate(&tpl)
			if err != nil {
				log.Errorf("add system template by name %s failed(%v)", reqData.OSSystemTemplateName, err)
				return err
			} else {
				//重新获取新增模板ID
				sysTpl, err = repo.GetSystemTemplateByName(reqData.OSSystemTemplateName)
				if err != nil {
					log.Errorf("get system template by name %s failed(%v)", reqData.OSSystemTemplateName, err)
				}
			}
		}
		//校验是否有装机记录,仅允许更新已有装机记录的系统名称
		ds, err := repo.GetDeviceSettingBySN(reqData.SN)
		if err != nil {
			log.Errorf("get device(SN:%s) setting err:%v", reqData.SN, err)
			return err
		}
		if ds == nil {
			log.Errorf("device(SN:%s) setting not found ,do nothing", reqData.SN)
			return errors.New("不存在装机记录，无法更新操作系统名称")
		} else {
			// 关联到PXE模板的系统名称
			ds.InstallType = model.InstallationPXE
			ds.SystemTemplateID = sysTpl.ID
			ds.ImageTemplateID = 0
			ds.UpdatedAt = time.Now()
			// 更新0值需使用save
			if err = repo.SaveDeviceSetting(ds); err != nil {
				log.Errorf("failed to device setting (system template name:%s for SN:%s), err:%v", reqData.OSSystemTemplateName, reqData.SN, err)
				return err
			}
			log.Infof("updated device setting (system template name:%s for SN:%s)", reqData.OSSystemTemplateName, reqData.SN)
			return nil
		}
	}
	return errors.New("不允许为空")
}


// SaveDeviceSettingItemWithoutInstall 保存设备装机参数条目
type SaveDeviceSettingItemWithoutInstall struct {
	// 设备序列号
	SN string `json:"sn"`
	// 操作系统安装模板
	OSTemplateName string `json:"os_template_name"`
}

// SaveDeviceSettingsWithoutInstallsReq 保存设备装机参数请求结构体
type SaveDeviceSettingsWithoutInstallsReq struct {
	Settings    DeviceSettingsWithoutInstalls
	CurrentUser *model.CurrentUser // 操作人
}

// DeviceSettings 设备装机参数集合
type DeviceSettingsWithoutInstalls []*SaveDeviceSettingItemWithoutInstall

// SaveDeviceSettingsWithoutInstalls 批量保存设备装机参数，不进入部署，直接运营状态置为‘已上架’（裸金属、特殊设备交付）
func SaveDeviceSettingsWithoutInstalls(log logger.Logger, repo model.Repo, conf *config.Config, lim limiter.Limiter, reqData *SaveDeviceSettingsWithoutInstallsReq) (succeeds []string, err error) {
	for i := range reqData.Settings {
		// 保存装机参数
		if err = saveDeviceSettingWithoutInstalls(log, repo, conf, lim, reqData.Settings[i], reqData.CurrentUser.LoginName); err != nil {
			return succeeds, err
		}

		// 更新设备状态=已上架
		mod := &model.Device{
			SN:              reqData.Settings[i].SN,
			OperationStatus: model.DevOperStatOnShelve,
		}
		_, err = repo.UpdateDeviceBySN(mod)
		if err != nil {
			log.Errorf("Failt to update status of device %s to onshelf", reqData.Settings[i].SN)
		}
		succeeds = append(succeeds, reqData.Settings[i].SN)
	}
	return succeeds, nil
}

// saveDeviceSettingWithoutInstalls 保存单台设备的装机参数(忽略部署)
func saveDeviceSettingWithoutInstalls(log logger.Logger, repo model.Repo, conf *config.Config, lim limiter.Limiter, sett *SaveDeviceSettingItemWithoutInstall, LoginName string) (err error) {
	log.Infof("Start saving device(%s) settings without installation", sett.SN)
	defer log.Infof("End saving device(%s) settings without installation", sett.SN)

	// 获取设备信息
	dev, err := repo.GetDeviceBySN(sett.SN)
	if err != nil {
		log.Errorf("Failed to get device of sn: %s , error: %s", sett.SN, err)
		return err
	}
    // 根据 USiteID 查询获取 PhysicalArea
    var ServerUSite *model.ServerUSite
    ServerUSite, err = repo.GetServerUSiteByID(*dev.USiteID)
    if err != nil {
    	log.Errorf("Fail to get serverusite of device(sn: %s ), error: %s", dev.SN, err)
    	return err
    }	
	// 初始化装机模板数据	
	ds := model.DeviceSetting{
		SN:             	sett.SN,
		InstallType:    	model.InstallationPXE,     //忽略部署，默认PXE
		Status:         	model.InstallStatusSucc,  //忽略部署，默认成功
		NeedExtranetIP:     model.NO,
		NeedIntranetIPv6:   model.NO,				  //默认仅分配IPV4
		NeedExtranetIPv6:   model.NO,				  //默认仅分配IPV4
		HardwareTemplateID: 0,                        //忽略部署，默认0
	}
	// 查询旧装机记录
	old, err := repo.GetDeviceSettingBySN(sett.SN)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if old != nil {
		ds.ID = old.ID
		ds.CreatedAt = old.CreatedAt
		ds.Updater = LoginName
		// 旧安装模板ID 置0 
		ds.ImageTemplateID = 0
		ds.SystemTemplateID = 0
	} else {
		ds.Creator = LoginName
	}
	// 操作系统名称校验,不存在则新增一条记录
	sysTemplate, err := repo.GetSystemTemplateByName(sett.OSTemplateName)
	if err != nil {
		log.Errorf("get system template by name %s failed(%v) ,adding one", sett.OSTemplateName, err)
		tpl := model.SystemTemplate{
			Family:   	"Custom",
			BootMode: 	"uefi",
			Name:     	sett.OSTemplateName,
			PXE:        "#NULL",
			Content:    "#NULL",
			OSLifecycle: model.OSTesting,
			Arch:		 model.OSARCHUNKNOWN,
		}
		_, err := repo.SaveSystemTemplate(&tpl)
		if err != nil {
			log.Errorf("add system template by name %s failed(%v)", sett.OSTemplateName, err)
			return err
		} else {
			//重新获取新增模板ID
			sysTemplate, err = repo.GetSystemTemplateByName(sett.OSTemplateName)
			if err != nil {
				log.Errorf("get system template by name %s failed(%v)", sett.OSTemplateName, err)
				return err
			}
		}
	}
	ds.SystemTemplateID = sysTemplate.ID

	log.Infof("Start assigning IPv4 to the device(%s)", sett.SN)
	// 记录本次已分配的IP ID以便回滚
	var ipIDsForRollback []uint

	//根据需要分配IP，假如有装机记录，并且有IP，则不用重复分配了
	if old == nil || old.IntranetIP == "" {
		// 获取内网业务IP
		intraIP, err := repo.AssignIntranetIP(sett.SN)
		if err != nil {
			return err
		}
		ds.IntranetIPNetworkID = intraIP.IPNetworkID
		ds.IntranetIP = intraIP.IP
		log.Infof("Device(%s) is assigned to intranet ip: %s", sett.SN, intraIP.IP)
		ipIDsForRollback = append(ipIDsForRollback, intraIP.ID)
	} else {
		ds.IntranetIPNetworkID = old.IntranetIPNetworkID
		ds.IntranetIP = old.IntranetIP
	}
	// 规则引擎的事实对象
	df := DeviceFact{
		SN:                   dev.SN,
		Category:             dev.Category,
		PhysicalArea:         ServerUSite.PhysicalArea,
		Arch:                 dev.Arch,
		Vendor:               dev.Vendor,
	}
	// 规则推理得到装机参数-network
	result, err := RuleTypeInduction(log, repo, &df, model.DeviceSettingRuleNetwork)
	if err != nil {
		log.Errorf("Fail to get rule action results of device(sn: %s ), error: %s", dev.SN, err)
		return err
	}
	// 获取网络配置
	if result.NetworkResult == "la_wa" {
		ds.NeedExtranetIP = model.YES
		// 获取外网业务IP
		log.Infof("Device(%s) need extranet ip", sett.SN)
		if old == nil || old.ExtranetIP == "" {
			extraIP, err := repo.AssignExtranetIP(sett.SN)
			if err != nil {
				return err
			}
			ds.ExtranetIPNetworkID = extraIP.IPNetworkID
			ds.ExtranetIP = extraIP.IP
			log.Infof("Device(%s) is assigned to extranet ip: %s", sett.SN, extraIP.IP)
			ipIDsForRollback = append(ipIDsForRollback, extraIP.ID)
		} else {
			ds.ExtranetIPNetworkID = old.ExtranetIPNetworkID
			ds.ExtranetIP = old.ExtranetIP
		}
	} else {
		// 尝试释放之前已经占用的外网业务IP
		_, _ = repo.ReleaseIP(sett.SN, model.IPScopeExtranet)
	}

	// 保存装机参数，失败时回滚已分配的IP
	if err = repo.SaveDeviceSetting(&ds); err != nil {
		for _, id := range ipIDsForRollback {
			_ = repo.UnassignIP(id)
			log.Debugf("Unassign IP (ID: %v) when failed to save device setting of SN: %s", id, sett.SN)
		}
		return err
	}
	if old != nil && old.ID > 0 {
		_, _ = repo.UpdateDeviceLogType(old.ID, model.DeviceLogInstallType, model.DeviceLogHistoryType)
	}

	log.Infof("Device(%s) settings have been saved.", sett.SN)
	return nil
}


// SetInstallsOKReq 批量设置部署状态=success
type SetInstallsOKReq struct {
	SNs []string `json:"sns"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *SetInstallsOKReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.SNs: "sns",
	}
}

// Validate 结构体数据校验
func (reqData *SetInstallsOKReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())
	for i := range reqData.SNs {
		_, err := repo.GetDeviceBySN(reqData.SNs[i])
		if gorm.IsRecordNotFoundError(err) {
			errs.Add([]string{"sns"}, binding.BusinessError, fmt.Sprintf("物理机(%s)信息不存在", reqData.SNs[i]))
			return errs
		}
		if err != nil {
			errs.Add([]string{"sns"}, binding.SystemError, "系统内部错误")
			return errs
		}

		sett, err := repo.GetDeviceSettingBySN(reqData.SNs[i])
		if gorm.IsRecordNotFoundError(err) {
			errs.Add([]string{"sns"}, binding.BusinessError, fmt.Sprintf("物理机(%s)无历史部署记录。", reqData.SNs[i]))
			return errs
		}
		if err != nil {
			errs.Add([]string{"sns"}, binding.SystemError, "系统内部错误")
			return errs
		}
		if sett.Status != model.InstallStatusFail {
			errs.Add([]string{"sns"}, binding.BusinessError, fmt.Sprintf("物理机(%s)部署状态非失败，仅适用于部署失败的物理机。", reqData.SNs[i]))
			return errs
		}

	}
	return errs
}

// SetInstallsOK 批量设置部署状态=success
func SetInstallsOK(log logger.Logger, repo model.Repo, conf *config.Config, lim limiter.Limiter, reqData *SetInstallsOKReq) (err error) {
	for i := range reqData.SNs {
		devSett, err := repo.GetDeviceSettingBySN(reqData.SNs[i])
		if err != nil {
			return err
		}
		if _, err = repo.UpdateInstallStatusAndProgressByID(devSett.ID, model.InstallStatusSucc, 1.0); err != nil {
			return err
		}
		//如果备注字段有运行状态，则说明是需要还原该状态
		d, err := repo.GetDeviceBySN(reqData.SNs[i])
		if d == nil {
			log.Error("device:%s not exist,%v", reqData.SNs[i], err)
		} else if d != nil && d.Remark != "" && validOperationStatus(d.Remark) {
			d.OperationStatus = d.Remark
			d.Remark = ""
		} else {
			d.OperationStatus = model.DevOperStatOnShelve // 系统安装完成后运营状态改为'已上架'
		}
		if _, err := repo.SaveDevice(d); err != nil {
			log.Error("save device:%d fail,%v", d, err)
		}
		//更新装机日志
		_, err = repo.SaveDeviceLog(&model.DeviceLog{
			DeviceSettingID: devSett.ID,
			Title:           "设置部署状态",
			LogType:         model.DeviceLogInstallType,
			Content:         "SUCCESS",
			SN:              reqData.SNs[i],
		})
		if err != nil {
			log.Errorf("save device log err:%v", err)
		}
		// 归还令牌
		if conf.DHCPLimiter.Enable {
			if bucket, _ := lim.Route(reqData.SNs[i]); bucket != nil {
				if token, _ := repo.GetTokenBySN(reqData.SNs[i]); token != "" {
					_ = bucket.Return(reqData.SNs[i], limiter.Token(token))
				}
			}
		}	
	}
	return nil
}