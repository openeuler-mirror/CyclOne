package service

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/voidint/binding"
	"github.com/voidint/page"

	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils/collection"
	mystrings "idcos.io/cloudboot/utils/strings"
)

// ImageTemplateResp 镜像安装模板
type ImageTemplateResp struct {
	ID         uint   `json:"id"`
	Family     string `json:"family"`      // 操作系统族系
	Name       string `json:"name"`        // 模板名
	BootMode   string `json:"boot_mode"`   // 启动模式
	URL        string `json:"url"`         // PXE 引导模板内容
	Username   string `json:"username"`    // 操作系统用户名
	Password   string `json:"password"`    // 操作系统用户密码
	PreScript  string `json:"pre_script"`  // 前置脚本
	PostScript string `json:"post_script"` // 后置脚本
	CreatedAt  string `json:"created_at"`  // 创建时间
	UpdatedAt  string `json:"updated_at"`  // 更新时间
	// OS生命周期：Testing|Active(Default)|Active|Containment|EOL
	OSLifecycle     string `json:"os_lifecycle"`
	//  OS架构平台：x86_64|aarch64	
	Arch            string `json:"arch"`	
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

// GetImageTemplatePageReq 查询镜像安装模板分页请求结构体
type GetImageTemplatePageReq struct {
	Family   string `json:"family"`    // 操作系统族系
	BootMode string `json:"boot_mode"` // 启动模式
	// OS生命周期：Testing|Active(Default)|Active|Containment|EOL
	OSLifecycle     string `json:"os_lifecycle"`
	//  OS架构平台：x86_64|aarch64	
	Arch            string `json:"arch"`	
	Name     string `json:"name"`      // 模板名
	Page     int64  `json:"page"`
	PageSize int64  `json:"page_size"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *GetImageTemplatePageReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.Family:   "family",
		&reqData.BootMode: "boot_mode",
		&reqData.OSLifecycle:   "os_lifecycle",
		&reqData.Arch:          "arch",	
		&reqData.Name:     "name",
		&reqData.Page:     "page",
		&reqData.PageSize: "page_size",
	}
}

// GetImageTemplatePage 按条件查询镜像安装模板分页列表
func GetImageTemplatePage(log logger.Logger, repo model.Repo, reqData *GetImageTemplatePageReq) (pg *page.Page, err error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 100 {
		reqData.PageSize = 10
	}
	if reqData.Page < 0 {
		reqData.Page = 0
	}

	cond := model.ImageTemplate{
		Family:   reqData.Family,
		BootMode: reqData.BootMode,
		Name:     reqData.Name,
	}

	totalRecords, err := repo.CountImageTemplates(&cond)
	if err != nil {
		return nil, err
	}

	pager := page.NewPager(reflect.TypeOf(&ImageTemplateResp{}), reqData.Page, reqData.PageSize, totalRecords)
	items, err := repo.GetImageTemplates(&cond, model.OneOrderBy("name", model.ASC), pager.BuildLimiter())
	if err != nil {
		return nil, err
	}

	for i := range items {
		tpl, err := convert2ImageTemplate(items[i])
		if err != nil {
			log.Error(err)
			return nil, err
		}
		pager.AddRecords(tpl)
	}
	return pager.BuildPage(), nil
}

// GetImageTemplateByID 返回指定ID的镜像安装模板
func GetImageTemplateByID(log logger.Logger, repo model.Repo, id uint) (*ImageTemplateResp, error) {
	tpl, err := repo.GetImageTemplateByID(id)
	if err != nil {
		return nil, err
	}
	return convert2ImageTemplate(tpl)
}

// convert2ImageTemplate 将model层的镜像安装模板对象转化成service层的镜像安装模板对象
func convert2ImageTemplate(src *model.ImageTemplate) (*ImageTemplateResp, error) {
	if src == nil {
		return nil, nil
	}
	tpl := ImageTemplateResp{
		ID:           src.ID,
		Family:       src.Family,
		Name:         src.Name,
		BootMode:     src.BootMode,
		URL:          src.ImageURL,
		Username:     src.Username,
		Password:     src.Password,
		PreScript:    mystrings.DOS2UNIX(src.PreScript),
		PostScript:   mystrings.DOS2UNIX(src.PostScript),
		CreatedAt:    src.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    src.UpdatedAt.Format("2006-01-02 15:04:05"),
		OSLifecycle:  src.OSLifecycle,
		Arch:         src.Arch,
	}
	if src.Partition != "" {
		if err := json.Unmarshal([]byte(src.Partition), &tpl.Disks); err != nil {
			return nil, err
		}
	}
	return &tpl, nil
}

// SaveImageTemplateReq 保存镜像安装模板请求结构体
type SaveImageTemplateReq struct {
	// 注：更新时该属性大于0
	ID uint `json:"id"`
	// 操作系统族系
	Family string `json:"family"`
	// 模板名
	Name string `json:"name"`
	// 启动模式
	BootMode string `json:"boot_mode"`
	// PXE 引导模板内容
	URL string `json:"url"`
	// 操作系统用户名
	Username string `json:"username"`
	// 操作系统用户密码
	Password string `json:"password"`
	// 前置脚本
	PreScript string `json:"pre_script"`
	// 后置脚本
	PostScript string `json:"post_script"`
	// OS生命周期：Testing|Active(Default)|Active|Containment|EOL
	OSLifecycle     string `json:"os_lifecycle"`
	//  OS架构平台：x86_64|aarch64	
	Arch            string `json:"arch"`	
	// 硬盘信息
	Disks []struct {
		Name       string `json:"name"`
		Partitions []struct {
			Name       string `json:"name"`
			Size       string `json:"size"`
			Fstype     string `json:"fstype"`
			Mountpoint string `json:"mountpoint"`
		} `json:"partitions"`
	} `json:"disks"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *SaveImageTemplateReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.Family:        "family",
		&reqData.Name:          "name",
		&reqData.BootMode:      "boot_mode",
		&reqData.Username:      "username",
		&reqData.Password:      "password",
		&reqData.PreScript:     "pre_script",
		&reqData.PostScript:    "post_script",
		&reqData.Disks:         "disks",
		&reqData.OSLifecycle:   "os_lifecycle",
		&reqData.Arch:          "arch",				
	}
}

// Validate 结构体数据校验
func (reqData *SaveImageTemplateReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())

	if reqData.Family == "" {
		errs.Add([]string{"family"}, binding.RequiredError, "操作系统不能为空")
		return errs
	}

	if reqData.Name == "" {
		errs.Add([]string{"name"}, binding.RequiredError, "名称不能为空")
		return errs
	}
	
	// 校验模板唯一性
	items, err := repo.GetImageTemplates(&model.ImageTemplate{
		Name: reqData.Name,
	}, nil, nil)
	if err != nil {
		errs.Add([]string{"name"}, binding.SystemError, "系统内部错误")
		return errs
	}

	for _, item := range items {
		if (reqData.ID == 0 && strings.ToLower(item.Name) == strings.ToLower(reqData.Name)) || // 新增模板时，模板名不能重名。
			(reqData.ID > 0 && strings.ToLower(item.Name) == strings.ToLower(reqData.Name) && reqData.ID != item.ID) { // 更新模板时，模板名不能重名（除了自身外）。
			errs.Add([]string{"name"}, binding.BusinessError, "名称不能重复")
			return errs
		}
	}

	if reqData.BootMode == "" {
		errs.Add([]string{"boot_mode"}, binding.RequiredError, "启动模式不能为空")
		return errs
	}

	if !collection.InSlice(reqData.BootMode, []string{model.BootModeBIOS, model.BootModeUEFI}) {
		errs.Add([]string{"boot_mode"}, binding.BusinessError, "无效的启动模式")
		return errs
	}
	if reqData.OSLifecycle == "" {
		errs.Add([]string{"os_lifecycle"}, binding.RequiredError, "OS生命周期不能为空")
		return errs
	}

	if !collection.InSlice(reqData.OSLifecycle, []string{model.OSTesting, model.OSActiveDefault, model.OSActive, model.OSContainment, model.OSEOL}) {
		errs.Add([]string{"os_lifecycle"}, binding.BusinessError, "无效的OS生命周期")
		return errs
	}

	if reqData.Arch == "" {
		errs.Add([]string{"arch"}, binding.RequiredError, "OS架构不能为空")
		return errs
	}

	if !collection.InSlice(reqData.Arch, []string{model.OSARCHAARCH64, model.OSARCHX8664}) {
		errs.Add([]string{"arch"}, binding.BusinessError, "无效的OS架构")
		return errs
	}
	if reqData.URL == "" {
		errs.Add([]string{"url"}, binding.RequiredError, "镜像下载地址不能为空")
		return errs
	}
	if len(reqData.Disks) == 0 {
		errs.Add([]string{"disks"}, binding.RequiredError, "分区不能为空")
		return errs
	}
	return errs
}

// SaveImageTemplate 保存镜像安装模板
func SaveImageTemplate(log logger.Logger, repo model.Repo, reqData *SaveImageTemplateReq) (id uint, err error) {
	disks, err := json.Marshal(reqData.Disks)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	tpl := model.ImageTemplate{
		Family:     reqData.Family,
		BootMode:   reqData.BootMode,
		Name:       reqData.Name,
		ImageURL:   reqData.URL,
		Username:   reqData.Username,
		Password:   reqData.Password,
		Partition:  string(disks),
		PreScript:  mystrings.DOS2UNIX(reqData.PreScript),
		PostScript: mystrings.DOS2UNIX(reqData.PostScript),
		OSLifecycle:     reqData.OSLifecycle,
		Arch:            reqData.Arch,
	}
	tpl.ID = reqData.ID
	affected, err := repo.SaveImageTemplate(&tpl)
	if err != nil {
		return 0, err
	}
	if affected <= 0 && tpl.ID > 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return tpl.ID, nil
}

// RemoveImageTemplate 删除指定ID的镜像安装模板
func RemoveImageTemplate(log logger.Logger, repo model.Repo, id uint) (err error) {
	affected, err := repo.RemoveImageTemplate(id)
	if err != nil {
		return err
	}
	if affected <= 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
