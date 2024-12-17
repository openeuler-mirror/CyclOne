package service

import (
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
	mytimes "idcos.io/cloudboot/utils/times"
)

// SystemTemplateResp 系统安装模板
type SystemTemplateResp struct {
	ID        uint   `json:"id"`
	Family    string `json:"family"`     // 操作系统族系
	Name      string `json:"name"`       // 模板名
	BootMode  string `json:"boot_mode"`  // 启动模式
	PXE       string `json:"pxe"`        // PXE 引导模板内容
	Username  string `json:"username"`   // 操作系统用户名
	Password  string `json:"password"`   // 操作系统用户密码
	Content   string `json:"content"`    // 模板内容
	CreatedAt string `json:"created_at"` // 创建时间
	UpdatedAt string `json:"updated_at"` // 更新时间
	// OS生命周期：Testing|Active(Default)|Active|Containment|EOL
	OSLifecycle     string `json:"os_lifecycle"`
	//  OS架构平台：x86_64|aarch64	
	Arch            string `json:"arch"`	
}

// GetSystemTemplatePageReq 查询系统安装模板分页请求结构体
type GetSystemTemplatePageReq struct {
	Family   string `json:"family"`    // 操作系统族系
	Name     string `json:"name"`      // 模板名
	BootMode string `json:"boot_mode"` // 启动模式
	// OS生命周期：Testing|Active(Default)|Active|Containment|EOL
	OSLifecycle     string `json:"os_lifecycle"`
	//  OS架构平台：x86_64|aarch64	
	Arch            string `json:"arch"`	
	Page     int64  `json:"page"`
	PageSize int64  `json:"page_size"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *GetSystemTemplatePageReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.Family:   "family",
		&reqData.Name:     "name",
		&reqData.BootMode: "boot_mode",
		&reqData.OSLifecycle:   "os_lifecycle",
		&reqData.Arch:          "arch",	
		&reqData.Page:     "page",
		&reqData.PageSize: "page_size",
	}
}

// GetSystemTemplatePage 按条件查询系统安装模板分页列表
func GetSystemTemplatePage(log logger.Logger, repo model.Repo, reqData *GetSystemTemplatePageReq) (pg *page.Page, err error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 100 {
		reqData.PageSize = 10
	}
	if reqData.Page < 0 {
		reqData.Page = 0
	}

	cond := model.SystemTemplate{
		Family:   reqData.Family,
		BootMode: reqData.BootMode,
		Name:     reqData.Name,
	}

	totalRecords, err := repo.CountSystemTemplates(&cond)
	if err != nil {
		return nil, err
	}

	pager := page.NewPager(reflect.TypeOf(&SystemTemplateResp{}), reqData.Page, reqData.PageSize, totalRecords)
	items, err := repo.GetSystemTemplates(&cond, model.OneOrderBy("name", model.ASC), pager.BuildLimiter())
	if err != nil {
		return nil, err
	}

	for i := range items {
		pager.AddRecords(convert2SystemTemplate(items[i]))
	}
	return pager.BuildPage(), nil
}

// GetSystemTemplateByID 返回指定ID的系统模板
func GetSystemTemplateByID(log logger.Logger, repo model.Repo, id uint) (*SystemTemplateResp, error) {
	tpl, err := repo.GetSystemTemplateByID(id)
	if err != nil {
		return nil, err
	}
	return convert2SystemTemplate(tpl), nil
}

// convert2SystemTemplate 将model层的系统安装模板对象转化成service层的系统安装模板对象
func convert2SystemTemplate(src *model.SystemTemplate) *SystemTemplateResp {
	if src == nil {
		return nil
	}
	return &SystemTemplateResp{
		ID:           src.ID,
		Family:       src.Family,
		Name:         src.Name,
		BootMode:     src.BootMode,
		PXE:          src.PXE,
		Username:     src.Username,
		Password:     src.Password,
		Content:      src.Content,
		CreatedAt:    src.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    src.UpdatedAt.Format("2006-01-02 15:04:05"),
		OSLifecycle:  src.OSLifecycle,
		Arch:         src.Arch,
	}
}

// SaveSystemTemplateReq 保存(新增/更新)系统安装模板请求结构体
type SaveSystemTemplateReq struct {
	// 注：更新时该属性大于0
	ID uint `json:"id"`
	// 操作系统族系
	Family string `json:"family"`
	// 模板名
	Name string `json:"name"`
	// 启动模式
	BootMode string `json:"boot_mode"`
	// PXE 引导模板内容
	PXE string `json:"pxe"`
	// 操作系统用户名
	Username string `json:"username"`
	// 操作系统用户密码
	Password string `json:"password"`
	// 模板内容
	Content string `json:"content"`
	// OS生命周期：Testing|Active(Default)|Active|Containment|EOL
	OSLifecycle     string `json:"os_lifecycle"`
	//  OS架构平台：x86_64|aarch64	
	Arch            string `json:"arch"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *SaveSystemTemplateReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.Family:        "family",
		&reqData.Name:          "name",
		&reqData.BootMode:      "boot_mode",
		&reqData.PXE:           "pxe",
		&reqData.Username:      "username",
		&reqData.Password:      "password",
		&reqData.Content:       "content",
		&reqData.OSLifecycle:   "os_lifecycle",
		&reqData.Arch:          "arch",		
	}
}

// Validate 结构体数据校验
func (reqData *SaveSystemTemplateReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())

	if reqData.Family == "" {
		errs.Add([]string{"family"}, binding.RequiredError, "操作系统不能为空")
		return errs
	}

	if reqData.Name == "" {
		errs.Add([]string{"name"}, binding.RequiredError, "名称不能为空")
		return errs
	}

	// 校验系统安装模板唯一性
	items, err := repo.GetSystemTemplates(&model.SystemTemplate{
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

	if reqData.PXE == "" {
		errs.Add([]string{"pxe"}, binding.RequiredError, "PXE模板不能为空")
		return errs
	}
	if reqData.Content == "" {
		errs.Add([]string{"content"}, binding.RequiredError, "系统模板不能为空")
		return errs
	}
	return errs
}

// SaveSystemTemplate 保存系统安装模板
func SaveSystemTemplate(log logger.Logger, repo model.Repo, reqData *SaveSystemTemplateReq) (id uint, err error) {
	tpl := model.SystemTemplate{
		Family:   reqData.Family,
		BootMode: reqData.BootMode,
		Name:     reqData.Name,
		Username: reqData.Username,
		Password: reqData.Password,
		PXE:      mystrings.DOS2UNIX(reqData.PXE),
		Content:  mystrings.DOS2UNIX(reqData.Content),
		OSLifecycle:     reqData.OSLifecycle,
		Arch:            reqData.Arch,
	}
	tpl.ID = reqData.ID
	affected, err := repo.SaveSystemTemplate(&tpl)
	if err != nil {
		return 0, err
	}
	if affected <= 0 && tpl.ID > 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return tpl.ID, nil
}

// RemoveSystemTemplate 删除指定ID的系统安装模板
func RemoveSystemTemplate(log logger.Logger, repo model.Repo, id uint) (err error) {
	affected, err := repo.RemoveSystemTemplate(id)
	if err != nil {
		return err
	}
	if affected <= 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// SystemTemplate 系统安装模板
type SystemTemplate struct {
	ID                 uint            `json:"id"`
	Family             string          `json:"family"`     // 操作系统族系
	Name               string          `json:"name"`       // 模板名
	BootMode           string          `json:"boot_mode"`  // 启动模式
	PXE                string          `json:"pxe"`        // PXE 引导模板内容
	Username           string          `json:"username"`   // 操作系统用户名
	Password           string          `json:"password"`   // 操作系统用户密码
	Content            string          `json:"content"`    // 模板内容
	OSLifecycle        string          `json:"os_lifecycle"` // OS生命周期：Testing|Active(Default)|Active|Containment|EOL
	Arch               string          `json:"arch"`  //  OS架构平台：x86_64|aarch64	
	CreatedAt mytimes.ISOTime          `json:"created_at"` // 创建时间
	UpdatedAt mytimes.ISOTime          `json:"updated_at"` // 更新时间
}

// GetSystemTemplateBySN 根据SN查询设备的系统模板信息
func GetSystemTemplateBySN(log logger.Logger, repo model.Repo, sn string) (*SystemTemplate, error) {
	tpl, err := repo.GetSystemTemplateBySN(sn)
	if err != nil {
		return nil, err
	}
	return &SystemTemplate{
		ID:          tpl.ID,
		Family:      tpl.Family,
		Name:        tpl.Name,
		BootMode:    tpl.BootMode,
		PXE:         tpl.PXE,
		Username:    tpl.Username,
		Password:    tpl.Username,
		Content:     mystrings.DOS2UNIX(tpl.Content),
		OSLifecycle: tpl.OSLifecycle,
		Arch:        tpl.Arch,
		CreatedAt:   mytimes.ISOTime(tpl.CreatedAt),
		UpdatedAt:   mytimes.ISOTime(tpl.UpdatedAt),
	}, nil
}
