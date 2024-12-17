package model


const (
	// OS生命周期
	OSTesting = "testing"  // 测试中，未提供技术支持
	OSActiveDefault = "active_default" // 默认提供该OS
	OSActive = "active" // 可按需提供该OS
	OSContainment = "containment" // 存量提供技术支持，不再支持新增
	OSEOL = "end_of_life" // 生命周期终止

	// OS架构平台
	OSARCHAARCH64 = "aarch64"
	OSARCHX8664 = "x86_64"
	OSARCHUNKNOWN = "unknown"
)

//OSLifecycleMap = map[string]string {
//	"testing":"Testing",
//	"active_default":"Active(Default)",
//	"active":"Active",
//	"containment":"Containment",
//	"end_of_life":"EOL",
//}

// IOSTemplate 系统安装模板持久化接口
type IOSTemplate interface {
	// 查询指定设备关联的系统安装模板
	GetSystemTemplateBySN(sn string) (*SystemTemplate, error)
	// GetSystemTemplatesByCond 根据条件查询系统安装模板
	GetSystemTemplatesByCond(cond *CommonTemplateCond) ([]*SystemTemplate, error)
	// GetImageTemplatesByCond 根据条件查询镜像安装模板
	GetImageTemplatesByCond(cond *CommonTemplateCond) (templates []*ImageTemplate, err error)
}

// CommonTemplateCond 系统安装模板查询条件
type CommonTemplateCond struct {
	Family          string `json:"family"`
	BootMode        string `json:"boot_mode"`
	Name            string `json:"name"`
	OSLifecycle     string `json:"os_lifecycle"` // OS生命周期：Testing|Active(Default)|Active|Containment|EOL
	Arch            string `json:"arch"`  //  OS架构平台：x86_64|aarch64	
}
