package model

import (
	"time"

	"github.com/voidint/page"
)

const (
	// ApprovalTypeIDCAbolish 审批单类型-数据中心裁撤
	ApprovalTypeIDCAbolish = "idc_abolish"
	// ApprovalTypeServerRoomAbolish 审批单类型-数据中心裁撤
	ApprovalTypeServerRoomAbolish = "server_room_abolish"
	// ApprovalTypeNetAreaOffline 审批单类型-网络区域下线
	ApprovalTypeNetAreaOffline = "network_area_offline"
	// ApprovalTypeIPUnassign 审批单类型-IP回收
	ApprovalTypeIPUnassign = "ip_unassign"
	// ApprovalTypeDevPowerOff 审批单类型-物理机关机
	ApprovalTypeDevPowerOff = "device_power_off"
	// ApprovalTypeDevRestart 审批单类型-物理机重启
	ApprovalTypeDevRestart = "device_restart"
	// ApprovalTypeCabinetPowerOff 审批单类型-机架关电
	ApprovalTypeCabinetPowerOff = "cabinet_power_off"
	// ApprovalTypeCabinetOffline 审批单类型-机架下架
	ApprovalTypeCabinetOffline = "cabinet_offline"
	// ApprovalTypeDeviceOSReinstallation 审批单类型-物理机OS重装
	ApprovalTypeDeviceOSReinstallation = "device_os_reinstallation"
	// ApprovalTypeDeviceMigration 审批单类型-物理机搬迁
	ApprovalTypeDeviceMigration = "device_migration"
	// ApprovalTypeDeviceRetirement 审批单类型-物理机退役(报废)
	ApprovalTypeDeviceRetirement = "device_retirement"
	// ApprovalTypeDeviceRecycleReinstall 审批单类型-物理机回收(重装)
	ApprovalTypeDeviceRecycleReinstall = "device_recycle_reinstall"
	// ApprovalTypeDeviceRecyclePreMove 审批单类型-物理机回收(待搬迁)
	ApprovalTypeDeviceRecyclePreMove = "device_recycle_pre_move"
	// ApprovalTypeDeviceRecyclePreRetire 审批单类型-物理机回收(待退役)
	ApprovalTypeDeviceRecyclePreRetire = "device_recycle_pre_retire"
)

const (
	// ApprovalStatusApproval 审批单状态-审批中
	ApprovalStatusApproval = "approval"
	// ApprovalStatusCompleted 审批单状态-已完成
	ApprovalStatusCompleted = "completed"
	// ApprovalStatusRevoked 审批单状态-已撤销
	ApprovalStatusRevoked = "revoked"
	// ApprovalStatusFailure 审批单状态-失败
	ApprovalStatusFailure = "failure"
)

const (
	MigTypeStore2Usite = "store_to_usite" //库房->机架
	MigTypeUsite2Usite = "usite_to_usite" //机架->机架
	MigTypeUsite2Store = "usite_to_store" //机架->库房
	MigTypeStore2Store = "store_to_store" //库房->库房
)

// Approval 审批
type Approval struct {
	ID         uint       `gorm:"column:id"`          // 审批单ID
	Title      string     `gorm:"column:title"`       // 审批单标题
	Type       string     `gorm:"column:type"`        // 审批类型
	Metadata   string     `gorm:"column:metadata"`    // 审批单元数据
	FrontData  string     `gorm:"column:front_data"`  // 审批单元数据快照
	Initiator  string     `gorm:"column:initiator"`   // 审批发起人ID
	Approvers  string     `gorm:"column:approvers"`   // 审批人ID构成的JSON数组字符串
	Cc         string     `gorm:"column:cc"`          // 抄送人ID构成的JSON数组字符串
	Remark     string     `gorm:"column:remark"`      //备注信息
	StartTime  *time.Time `gorm:"column:start_time"`  // 审批开始时间
	EndTime    *time.Time `gorm:"column:end_time"`    // 审批结束时间
	IsRejected string     `gorm:"column:is_rejected"` // 审批单是否被拒绝
	Status     string     `gorm:"column:status"`      // 审批单状态
}

// TableName 指定数据库表名
func (Approval) TableName() string {
	return "approval"
}

const (
	// ApprovalActionAgree 审批动作-同意
	ApprovalActionAgree = "agree"
	// ApprovalActionReject 审批动作-拒绝
	ApprovalActionReject = "reject"
)

// ApprovalStep 审批步骤
type ApprovalStep struct {
	ID         uint       `gorm:"column:id"`          // 审批步骤ID
	ApprovalID uint       `gorm:"column:approval_id"` // 所属审批单ID
	Approver   string     `gorm:"column:approver"`    // 审批步骤审批人ID
	Title      string     `gorm:"column:title"`       // 审批步骤标题
	Action     *string    `gorm:"column:action"`      // 审批动作
	Remark     string     `gorm:"column:remark"`      // 审批批注
	StartTime  *time.Time `gorm:"column:start_time"`  // 审批步骤开始时间
	EndTime    *time.Time `gorm:"column:end_time"`    // 审批步骤结束时间
	Hooks      string     `gorm:"column:hooks"`       // 当前审批步骤同意后执行的钩子对象数组字符串
}

// TableName 指定数据库表名
func (ApprovalStep) TableName() string {
	return "approval_step"
}

// IApproval 审批功能数据库操作接口
type IApproval interface {
	// SubmitApproval 提交审批单及其审批步骤
	SubmitApproval(approval *Approval, steps ...*ApprovalStep) (err error)
	// RevokeApproval 撤销目标审批单
	RevokeApproval(approvalID uint) (err error)
	// Approve 审批
	Approve(approvalID, stepID uint, action, remark string) error
	// UpdateApproval 修改审批单
	UpdateApproval(mod *Approval) (affected int64, err error)
	// GetApprovalByID 查询指定审批单
	GetApprovalByID(approvalID uint) (approval *Approval, err error)
	// GetApprovalStepByApprovalID 查询指定审批单的审批步骤明细
	GetApprovalStepByApprovalID(approvalID uint) (steps []*ApprovalStep, err error)
	// GetInitiatedApprovals 查询'我发起的'审批分页列表
	GetInitiatedApprovals(currentUserID string, cond *Approval, orderby OrderBy, limiter *page.Limiter) (items []*Approval, err error)
	// GetPendingApprovals 查询'待我审批'的审批分页列表
	GetPendingApprovals(currentUserID string, cond *Approval, orderby OrderBy, limiter *page.Limiter) (items []*Approval, err error)
	// GetApprovedApprovals 查询'我已审批的'的审批分页列表
	GetApprovedApprovals(currentUserID string, cond *Approval, orderby OrderBy, limiter *page.Limiter) (items []*Approval, err error)
	// CountApprovals 统计审批单个数
	CountInitiatedApprovals(currentUserID string, cond *Approval) (int64, error)
	// CountApprovals 统计审批单个数
	CountPendingApprovals(currentUserID string, cond *Approval) (int64, error)
	// CountApprovals 统计审批单个数
	CountApprovedApprovals(currentUserID string, cond *Approval) (int64, error)
	// GetApprovalStepByCond 根据条件查询指定审批单的审批步骤明细
	GetApprovalStepByCond(cond *ApprovalStep, isApproved bool) (steps []*ApprovalStep, err error)
	//GetApprovalStepByID 根据审批步骤ID查询审批步骤
	GetApprovalStepByID(stepID uint) (step *ApprovalStep, err error)
	// UpdateApprovalStep 修改审批步骤
	UpdateApprovalStep(mod *ApprovalStep) (affected int64, err error)
}

const (
	// HookIDCAbolish 数据中心裁撤钩子
	HookIDCAbolish = "idc_abolish_hook"
	// HookServerRoomAbolish 数据中心裁撤钩子
	HookServerRoomAbolish = "server_room_abolish_hook"
	// HookNetAreaOffline 网络区域下线
	HookNetAreaOffline = "net_area_offline_hook"
	// HookIPUnassign IP回收
	HookIPUnassign = "ip_unassign_hook"
	// HookDevPowerOff 物理机关机钩子
	HookDevPowerOff = "device_power_off_hook"
	// HookDevRestart 物理机重启钩子
	HookDevRestart = "device_restart_hook"
	// HookCabinetPowerOff 机架关电钩子
	HookCabinetPowerOff = "cabinet_power_off_hook"
	// HookCabinetOffline 机架下线钩子
	HookCabinetOffline = "cabinet_offline_hook"
	// HookDeviceMigration 物理机搬迁钩子
	HookDeviceMigration = "device_migration_hook"
	// HookDeviceMigrationPowerOff 物理机搬迁关电钩子
	HookDeviceMigrationPowerOff = "device_migration_poweroff_hook"
	// HookDeviceMigrationReleaseIP 物理机搬迁释放IP钩子
	HookDeviceMigrationReleaseIP = "device_migration_release_ip_hook"
	// HookDeviceMigrationReserveIP 物理机搬迁保留IP钩子
	HookDeviceMigrationReserveIP = "device_migration_reserve_ip_hook"	
	// HookDeviceRetirementPowerOff 物理机退役远程关机钩子
	HookDeviceRetirementPowerOff = "device_retirement_poweroff_hook"
	// HookDeviceRetirementReleaseIP 物理机退役释放IP钩子
	HookDeviceRetirementReleaseIP = "device_retirement_release_ip_hook"
	// HookDeviceRetirementReserveIP 物理机退役释放IP钩子
	HookDeviceRetirementReserveIP = "device_retirement_reserve_ip_hook"	
	// HookDeviceRetirement 物理机退役钩子(修改状态)
	HookDeviceRetirement = "device_retirement_hook"
	// HookDeviceOSReinstallation 物理机OS重装钩子
	HookDeviceOSReinstallation = "device_os_reinstallation_hook"
	// HookDeviceRecycleReinstall 物理机回收钩子
	HookDeviceRecycleReinstall = "device_recycle_reinstall_hook"
	// HookDeviceRecyclePreMove 物理机回收待搬迁钩子
	HookDeviceRecyclePreMove = "device_recycle_pre_move_hook"
	// HookDeviceRecyclePreRetire 物理机回收待退役钩子
	HookDeviceRecyclePreRetire = "device_recycle_pre_retire_hook"
)

// StepHook 审批步骤钩子
type StepHook struct {
	ID              string            `json:"id"`                // 字符串标识(唯一)
	Description     string            `json:"description"`       // 描述
	ContinueOnError bool              `json:"continue_on_error"` // 执行发生错误后是否继续执行后续钩子
	Results         []*StepHookResult `json:"results"`           // 执行明细
}

// StepHookResult 钩子执行明细
type StepHookResult struct {
	TargetID   string `json:"target_id"`   // 目标对象标识字符串
	TargetType string `json:"target_type"` // 目标对象类型。可选值: cabinet-机架; physical_machine-物理机;
	ExecResult string `json:"exec_result"` // 执行结果。可选值: success-成功; failure-失败;
	ExecInfo   string `json:"exec_info"`   // 针对目标对象的执行结果信息
}

// BeforeSave 保存审批信息前的钩子方法。
// 防止将空字符串写入类型为JSON的数据库字段中引发报错。
func (app *Approval) BeforeSave() (err error) {
	replaceIfBlank(&app.Metadata, EmptyJSONObject)
	replaceIfBlank(&app.FrontData, EmptyJSONObject)
	replaceIfBlank(&app.Approvers, EmptyJSONArray)
	replaceIfBlank(&app.Cc, EmptyJSONArray)
	return
}

// BeforeSave 保存审批步骤信息前的钩子方法。
// 防止将空字符串写入类型为JSON的数据库字段中引发报错。
func (step *ApprovalStep) BeforeSave() (err error) {
	replaceIfBlank(&step.Hooks, EmptyJSONArray)
	return
}
