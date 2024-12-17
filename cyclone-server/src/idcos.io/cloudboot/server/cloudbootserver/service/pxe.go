package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"idcos.io/cloudboot/utils/pxe"

	"github.com/jinzhu/gorm"
	"github.com/voidint/binding"

	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils/centos6"
)

// GenPXEReq 生成pxe文件请求参数
type GenPXEReq struct {
	SN  string `json:"-"`
	MAC string `json:"mac"`
}

// FieldMap 请求字段映射
func (reqData *GenPXEReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.MAC: "mac",
	}
}

// Validate 装机参数校验
func (reqData *GenPXEReq) Validate(r *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(r.Context())

	if reqData.SN == "" {
		errs.Add([]string{"sn"}, binding.RequiredError, "SN不能为空")
		return errs
	}
	_, err := repo.GetDeviceBySN(reqData.SN)
	if gorm.IsRecordNotFoundError(err) {
		errs.Add([]string{"sn"}, binding.BusinessError, http.StatusText(http.StatusNotFound))
		return errs
	}
	if err != nil {
		errs.Add([]string{"sn"}, binding.SystemError, "系统内部错误")
		return errs
	}
	// TODO 校验目标设备是否包含指定mac地址的网卡
	sett, err := repo.GetDeviceSettingBySN(reqData.SN)
	if gorm.IsRecordNotFoundError(err) {
		errs.Add([]string{"sn"}, binding.BusinessError, "该物理机还未部署上架")
		return errs
	}
	if err != nil {
		errs.Add([]string{"sn"}, binding.SystemError, "系统内部错误")
		return errs
	}
	if sett.InstallType == model.InstallationImage {
		errs.Add([]string{"sn"}, binding.BusinessError, "该物理机使用镜像方式安装，无需生成PXE文件。")
		return errs
	}
	// TODO 校验当前设备是否处于操作系统安装过程中
	_, err = repo.GetSystemTemplateBySN(reqData.SN)
	if gorm.IsRecordNotFoundError(err) {
		errs.Add([]string{"sn"}, binding.BusinessError, fmt.Sprintf("未找到id为%d的系统模板", sett.SystemTemplateID))
		return errs
	}
	return errs
}

const (
	// pxeDir 生成PXE文件的目录
	pxeDir = "/var/lib/tftpboot"
	// dhcpFilename DHCP配置文件路径
	dhcpFilename = "/etc/dhcp/dhcpd.conf"
)

// GenPXE 为指定设备在本地目录生成PXE文件
func GenPXE(log logger.Logger, repo model.Repo, reqData *GenPXEReq) (filename string, err error) {
	sysTemplate, err := repo.GetSystemTemplateBySN(reqData.SN)
	if err != nil {
		return "", err
	}
	if err = os.MkdirAll(pxeDir, 0755); err != nil {
		log.Error(err)
		return "", err
	}

	filename = filepath.Join(pxeDir, fmt.Sprintf("grub.cfg-01-%s", strings.Replace(reqData.MAC, ":", "-", -1)))
	if err = ioutil.WriteFile(filename, []byte(sysTemplate.PXE), 0644); err != nil {
		log.Error(err)
		return "", err
	}
	return filename, nil
}

// GenPXE4CentOS6UEFIReq 请求结构体
type GenPXE4CentOS6UEFIReq struct {
	SN  string `json:"-"`
	MAC string `json:"mac"`
}

// FieldMap 请求字段映射
func (reqData *GenPXE4CentOS6UEFIReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.MAC: "mac",
	}
}

// GenPXE4CentOS6UEFI 为CentOS6.x（UEFI）生成PXE文件、修改dhcp配置文件、重启dhcp服务。
func GenPXE4CentOS6UEFI(log logger.Logger, repo model.Repo, reqData *GenPXE4CentOS6UEFIReq) (filename string, err error) {
	// 1、获取设备对应的PXE内容
	sysTemplate, err := repo.GetSystemTemplateBySN(reqData.SN)
	if err != nil {
		return "", err
	}

	// 2、在指定目录生成PXE文件
	filename, err = pxe.GenFile(reqData.MAC, []byte(strings.Replace(sysTemplate.PXE, "{sn}", reqData.SN, -1)))
	if err != nil {
		log.Errorf("Generate pxe file error: %s", err.Error())
		return "", err
	}

	// 3、读取原DHCP配置文件并以此生成新的配置文件
	dhcpSRC, err := ioutil.ReadFile(dhcpFilename)
	if err != nil {
		log.Error(err)
		return "", err
	}
	log.Infof("DHCP configuration(before): \n%s", dhcpSRC)

	dhcpDST := centos6.AddOneToDHCP(reqData.SN, reqData.MAC, dhcpSRC)

	log.Infof("DHCP configuration(after): \n%s", dhcpDST)

	// 4、将dhcp配置持久化并重启服务
	if err = centos6.OverwriteDHCP(dhcpDST); err != nil {
		log.Errorf("Restart dhcpd error: %s", err.Error())
		log.Warn("Try to rollback the dhcp configuration file")
		_ = centos6.OverwriteDHCP(dhcpSRC) // 尝试回滚配置
		return "", err
	}
	return filename, nil
}

// GetPXEReq 获取PXE请求结构体
type GetPXEReq struct {
	SN         string `json:"-"`
	Arch       string `json:"arch"`
	Cpuvendor  string `json:"cpuvendor"`
}

// FieldMap 请求字段映射
func (reqData *GetPXEReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.Arch: "arch",
		&reqData.Cpuvendor: "cpuvendor",
	}
}

// Validate 装机参数校验
func (reqData *GetPXEReq) Validate(r *http.Request, errs binding.Errors) binding.Errors {
	// repo, _ := middleware.RepoFromContext(r.Context())
	if reqData.SN == "" {
		errs.Add([]string{"sn"}, binding.RequiredError, "SN不能为空")
		return errs
	}

	//if reqData.Arch == "" {
	//	reqData.Arch = X8664
	//}
//
	//if reqData.Arch != ARM64 && reqData.Arch != X8664 {
	//	errs.Add([]string{"arch"}, binding.BusinessError, fmt.Sprintf("不支持该硬件架构: %s", reqData.Arch))
	//	return errs
	//}

	return errs
}

const (
	// ARM64 CPU硬件架构-arm64 通过iPXE.efi启动固件获取${buildarch}
	ARM64 = "arm64"
	// X8664 CPU硬件架构-x86_64 通过iPXE.efi启动固件获取${buildarch}
	X8664 = "x86_64"
	// CPU厂商 通过iPXE.efi启动固件获取${cpuvendor}(当前该特性仅支持x86_64架构)
	CPU_Vendor_Intel = "intel"
	CPU_Vendor_Hygon = "hygon"
)

// 支持根据CPU厂商进行组合bootos名称： bootos_x86_64]_[intel|hygon|]
const (
	bootOSARM64PXE    = "bootos_arm64" //ARM64架构，但CPU厂商不在兼容列表时返回该默认模板
	bootOSX8664PXE    = "bootos_x86_64" //X8664架构，但CPU厂商不在兼容列表时返回该默认模板
	bootOSDefault     = "bootos_default" //CPU架构不在兼容列表时，返回该模板
	winPE2012X8664PXE = "winpe2012_x86_64" // win server默认模板
	localPXE          = "local"
)

// GetPXE 根据SN查询设备的PXE引导文件
// 按照如下规则返回不同的模板内容 issues/133
// 未采集设备信息，有或无装机参数，则返回名为bootos的系统安装模板。
// 已采集设备信息，无装机参数，则返回名为bootos的系统安装模板。
// 已采集设备信息，有装机参数，装机状态为pre_install/failure，则返回名为bootos的系统安装模板。
// 已采集设备信息，有装机参数，装机状态为pre_hwcheck返回名为bootos的系统安装模板。
// 已采集设备信息，有装机参数，装机状态为success，则返回名为local的系统安装模板。
// 已采集设备信息，有装机参数，装机状态为success，部署状态为offline, 则返回名为bootos的系统安装模板。
// 已采集设备信息，有装机参数，装机状态为installing，且进度小于0.5，则返回bootos的系统安装模板。
// 已采集设备信息，有装机参数，装机状态为installing，则返回装机参数中指定的系统安装模板。
func GetPXE(log logger.Logger, repo model.Repo, reqData *GetPXEReq) (pxe string, err error) {
	defer func() {
		//PXE模板中的{sn}占位符替换成具体的SN
		pxe = strings.Replace(pxe, "{sn}", reqData.SN, -1)
	}()
	log.Debugf("Device(SN: %s, ARCH: %s, CPUVENDOR(not supported @arm64): %s) req bootos", reqData.SN, reqData.Arch, reqData.Cpuvendor)
	// 判断是否采集硬件信息
	dev, err := repo.GetDeviceBySNOrMAC(reqData.SN)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return "", err
	}
	if dev == nil {
		log.Debugf("Device(%s) has not been collected, boot from bootos", reqData.SN)
		return getBootOSPXE(log, repo, reqData.Arch, reqData.Cpuvendor)
	}
	reqData.SN = dev.SN

	// 判断是否提交装机参数
	sett, err := repo.GetDeviceSettingBySN(reqData.SN)
	if gorm.IsRecordNotFoundError(err) {
		log.Debugf("Device(%s) has not been submit settings, boot from bootos", reqData.SN)
		return getBootOSPXE(log, repo, reqData.Arch, reqData.Cpuvendor)
	}
	if err != nil {
		return "", err
	}

	log.Infof("Device(%s) settings:\ninstall_type ==> %s\nsystem_template_id ==> %d\nimage_template_id ==> %d\nstatus ==> %s\nprogress ==> %f",
		sett.SN,
		sett.InstallType,
		sett.SystemTemplateID,
		sett.ImageTemplateID,
		sett.Status,
		sett.InstallProgress,
	)

	switch sett.Status {
	case model.InstallStatusPre, model.InstallStatusFail:
		// log.Debugf("Device(%s) has not been submit settings, boot from bootos", reqData.SN)
		return getBootOSPXE(log, repo, reqData.Arch, reqData.Cpuvendor)

	case model.InstallStatusSucc:
		if dev.OperationStatus == model.DevOperStatReinstalling {
			log.Debugf("Device(%s) is being reinstall, boot from bootos", reqData.SN)
			return getBootOSPXE(log, repo, reqData.Arch, reqData.Cpuvendor)
		}
		log.Debugf("Device(%s) has been install success, boot from local", reqData.SN)
		return getPXEByName(log, repo, localPXE)

	case model.InstallStatusIng:
		if sett.InstallProgress < 0.5 {
			log.Debugf("Device(%s) boot from bootos for subsequent hardware configuration", reqData.SN)
			return getBootOSPXE(log, repo, reqData.Arch, reqData.Cpuvendor)
		} else if sett.InstallProgress > 0.8 {
			return getPXEByName(log, repo, localPXE)
		}

		log.Debugf("Device(%s) is installing, boot from tempalte", reqData.SN)
		if sett.InstallType == model.InstallationImage {
			imgTemplate, _ := repo.GetImageTemplateBySN(reqData.SN)

			if imgTemplate != nil && strings.HasPrefix(strings.ToLower(imgTemplate.Name), "win") {
				log.Debugf("Device(%s) is installing by windows image,boot from winpe2012", reqData.SN)
				return getPXEByName(log, repo, winPE2012X8664PXE)
			}
		}

		sysTemplate, err := repo.GetSystemTemplateBySN(reqData.SN)
		if err != nil {
			return "", err
		}
		return sysTemplate.PXE, nil
	}

	log.Debugf("BootOS by default when get pxe template failed")
	return getBootOSPXE(log, repo, reqData.Arch, reqData.Cpuvendor)
}

// getBootOSPXE 获取进入BootOS的PXE模板
func getBootOSPXE(log logger.Logger, repo model.Repo, arch string, cpuvendor string) (pxe string, err error) {
	switch arch {
	case ARM64:
		log.Debugf("get bootos by name: %s", bootOSARM64PXE)
		return getPXEByName(log, repo, bootOSARM64PXE)
	case X8664:
		if strings.Contains(strings.ToLower(cpuvendor), CPU_Vendor_Hygon) {
			bootOS := "bootos_" + X8664 + "_" + CPU_Vendor_Hygon
			log.Debugf("get bootos by name: %s", bootOS)
			return getPXEByName(log, repo, bootOS)
		} else if strings.Contains(strings.ToLower(cpuvendor), CPU_Vendor_Intel) {
			bootOS := "bootos_" + X8664 + "_" + CPU_Vendor_Intel
			log.Debugf("get bootos by name: %s", bootOS)
			return getPXEByName(log, repo, bootOS)
		} else {
			log.Debugf("get bootos by name: %s", bootOSX8664PXE)
			return getPXEByName(log, repo, bootOSX8664PXE)
		}
	default:
		log.Debugf("get bootos by name: %s", bootOSDefault)
		return getPXEByName(log, repo, bootOSDefault)
	}
}

func getPXEByName(log logger.Logger, repo model.Repo, name string) (pxe string, err error) {
	tpl, err := repo.GetSystemTemplateByName(name)
	if err != nil {
		log.Errorf("Get system template(%s) failed: %s", name, err.Error())
		return "", err
	}
	return tpl.PXE, nil
}
