package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/voidint/binding"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
)

// TemplateReq 模板参数
type TemplateReq struct {
	// 安装方式。可选值: image-镜像安装模板; pxe-PXE安装模板;
	InstallType string `json:"install_type"`
	// 操作系统族系
	Family string `json:"family"`
	// 启动模式。可选值: legacy_bios-传统BIOS模式; uefi-UEFI模式;
	BootMode string `json:"boot_mode"`
	// 模板名(支持模糊查询)
	Name string `json:"name"`
	OSLifecycle        string          `json:"os_lifecycle"` // OS生命周期：Testing|Active(Default)|Active|Containment|EOL
	Arch               string          `json:"arch"`  //  OS架构平台：x86_64|aarch64
}

// SaveTemplateReq 新增、修改模板接口
type SaveTemplateReq struct {
	ID uint `json:"id"`
	// 安装方式。可选值: image-镜像安装模板; pxe-PXE安装模板;
	InstallType string `json:"install_type"`
	// 操作系统族系
	Family string `json:"family"`
	// 启动模式。可选值: legacy_bios-传统BIOS模式; uefi-UEFI模式;
	BootMode string `json:"boot_mode"`
	// 模板名(支持模糊查询)
	Name string `json:"name"`
	OSLifecycle        string          `json:"os_lifecycle"` // OS生命周期：Testing|Active(Default)|Active|Containment|EOL
	Arch               string          `json:"arch"`  //  OS架构平台：x86_64|aarch64	
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *SaveTemplateReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.InstallType:   "install_type",
		&reqData.Family:        "family",
		&reqData.Name:          "name",
		&reqData.BootMode:      "boot_mode",
		&reqData.ID:            "id",
		&reqData.OSLifecycle:   "os_lifecycle",
		&reqData.Arch:          "arch",
	}
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *TemplateReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.InstallType:   "install_type",
		&reqData.Family:        "family",
		&reqData.Name:          "name",
		&reqData.BootMode:      "boot_mode",
		&reqData.OSLifecycle:   "os_lifecycle",
		&reqData.Arch:          "arch",
	}
}

// Validate 结构体数据校验
func (reqData *TemplateReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {

	if reqData.InstallType != "" {
		if !(reqData.InstallType == model.InstallationPXE || reqData.InstallType == model.InstallationImage) {
			errs.Add([]string{"os_template"}, binding.RequiredError, fmt.Sprintf("安装类型参数不正确 %s, 可选值: image-镜像安装模板; pxe-PXE安装模板", reqData.InstallType))
		}
	}

	if reqData.BootMode != "" {
		if !(reqData.BootMode == model.BootModeUEFI || reqData.BootMode == model.BootModeBIOS) {
			errs.Add([]string{"os_template"}, binding.RequiredError, fmt.Sprintf("启动模式参数不正确 %s, 可选值: legacy_bios-传统BIOS模式; uefi-UEFI模式", reqData.InstallType))

		}
	}
	return errs
}

// OsTemplateResp 模板参数返回值
type OsTemplateResp struct {
	// 安装方式。可选值: image-镜像安装模板; pxe-PXE安装模板;
	InstallType string `json:"install_type"`
	// 操作系统族系
	Family string `json:"family"`
	// 启动模式。可选值: legacy_bios-传统BIOS模式; uefi-UEFI模式;
	BootMode string `json:"boot_mode"`
	// 模板名(支持模糊查询)
	Name string `json:"name"`
	// 主键
	ID uint `json:"id"`
	OSLifecycle        string          `json:"os_lifecycle"` // OS生命周期：Testing|Active(Default)|Active|Containment|EOL
	Arch               string          `json:"arch"`  //  OS架构平台：x86_64|aarch64
}

// TemplateRecordsResp 操作系统模板返回体
type TemplateRecordsResp struct {
	Records []*OsTemplateResp `json:"records"`
}

// GetTemplatesByCond 根据查询条件查询模板信息
func GetTemplatesByCond(log logger.Logger, repo model.Repo, cond *TemplateReq) (items []*OsTemplateResp, err error) {

	commonCond := model.CommonTemplateCond{
		Family:        cond.Family,
		BootMode:      cond.BootMode,
		Name:          cond.Name,
		OSLifecycle:   cond.OSLifecycle,
		Arch:          cond.Arch,
	}

	if cond.InstallType == model.InstallationImage {
		imageTmps, err := repo.GetImageTemplatesByCond(&commonCond)
		if err != nil {
			return nil, err
		}

		for _, imageTmp := range imageTmps {
			items = append(items, imageTemplateConvert(imageTmp))
		}

		return items, nil
	}

	if cond.InstallType == model.InstallationPXE {
		systemTmps, err := repo.GetSystemTemplatesByCond(&commonCond)
		if err != nil {
			return nil, err
		}
		for _, systemTmp := range systemTmps {
			if systemTmp.Name == "bootos" || systemTmp.Name == "local" {
				continue
			}
			items = append(items, systemTemplateConvert(systemTmp))
		}
		return items, nil
	}

	imageTemplates, err := repo.GetImageTemplatesByCond(&commonCond)
	if err != nil {
		return nil, err
	}

	for _, imageTmp := range imageTemplates {
		items = append(items, imageTemplateConvert(imageTmp))
	}

	systemTemplates, err := repo.GetSystemTemplatesByCond(&commonCond)
	if err != nil {
		return nil, err
	}
	for _, systemTmp := range systemTemplates {
		if strings.Contains(systemTmp.Name, "bootos") || systemTmp.Name == "local" || systemTmp.Name == "winpe2012_x86_64" {
			continue
		}
		items = append(items, systemTemplateConvert(systemTmp))
	}

	for i, resp := range items {
		resp.ID = uint(i + 1)
	}

	return items, nil
}

// imageTemplateConvert 镜像模板转换
func imageTemplateConvert(template *model.ImageTemplate) *OsTemplateResp {
	return &OsTemplateResp{
		InstallType:     model.InstallationImage,
		Family:          template.Family,
		BootMode:        template.BootMode,
		Name:            template.Name,
		ID:              template.ID,
		OSLifecycle:     template.OSLifecycle,
		Arch:            template.Arch,
	}
}

// systemTemplateConvert 系统模板转换
func systemTemplateConvert(template *model.SystemTemplate) *OsTemplateResp {
	return &OsTemplateResp{
		InstallType:     model.InstallationPXE,
		Family:          template.Family,
		BootMode:        template.BootMode,
		Name:            template.Name,
		ID:              template.ID,
		OSLifecycle:     template.OSLifecycle,
		Arch:            template.Arch,
	}
}
