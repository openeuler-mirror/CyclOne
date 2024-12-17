package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/voidint/binding"
	"github.com/voidint/page"

	"os"

	"idcos.io/cloudboot/config"
	"idcos.io/cloudboot/limiter"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils"
	"idcos.io/cloudboot/utils/times"
	"idcos.io/cloudboot/utils/upload"
	myuser "idcos.io/cloudboot/utils/user"
)

// Approval 审批
type Approval struct {
	ID        uint                   `json:"id"`         // 审批单ID
	Title     string                 `json:"title"`      // 审批单标题
	Type      string                 `json:"type"`       // 审批类型
	Metadata  map[string]interface{} `json:"metadata"`   // 审批单元数据
	Initiator string                 `json:"initiator"`  // 审批发起人ID
	Approvers []string               `json:"approvers"`  // 审批人ID构成的JSON数组字符串
	Cc        []string               `json:"cc"`         // 抄送人ID构成的JSON数组字符串
	StartTime *time.Time             `json:"start_time"` // 审批开始时间
	EndTime   *time.Time             `json:"end_time"`   // 审批结束时间
	Status    string                 `json:"status"`     // 审批单状态
	Steps     []*ApprovalStep        `json:"steps"`      // 审批步骤
}

// ApprovalStep 审批步骤
type ApprovalStep struct {
	ID         uint              `json:"id"`          // 审批步骤ID
	ApprovalID uint              `json:"approval_id"` // 所属审批单ID
	Approver   string            `json:"approver"`    // 审批步骤审批人ID
	Title      string            `json:"title"`       // 审批步骤标题
	Action     string            `json:"action"`      // 审批动作
	Remark     string            `json:"remark"`      // 审批批注
	StartTime  *time.Time        `json:"start_time"`  // 审批步骤开始时间
	EndTime    *time.Time        `json:"end_time"`    // 审批步骤结束时间
	Hooks      []*model.StepHook `json:"hooks"`       // 当前审批步骤同意后执行的钩子对象数组字符串
}

// approvalLimit 审批提交数据的限制，暂定不能多于10个
const approvalLimit = 10
const approvalLimit50 = 50
const defaultMailFrom = "wbk_wecloud@webank.com"

// SubmitIDsApprovalReq 根据ID修改状态的请求体
type SubmitIDsApprovalReq struct {
	SubmitApprovalCommon
	IDs []uint `json:"ids"` //ids
}

// ToJSON 将请求结构体转换成JSON
func (reqData *SubmitIDsApprovalReq) ToJSON() []byte {
	b, _ := json.Marshal(reqData.IDs)
	return b
}

// ApproversJSON 返回审批人ID构成的JSON数组字符串
func (reqData *SubmitIDsApprovalReq) ApproversJSON() []byte {
	b, _ := json.Marshal(reqData.Approvers)
	return b
}

// HooksJSON 返回审批步骤钩子构成的字符串
func (reqData *SubmitIDsApprovalReq) HooksJSON(hook ...string) []byte {
	hooks := []*model.StepHook{
		{
			ID:              hook[0],
			Description:     "",
			ContinueOnError: true,
		},
	}
	b, _ := json.Marshal(hooks)
	return b
}

// FieldMap 字段映射
func (reqData *SubmitIDsApprovalReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDs:       "ids",
		&reqData.Approvers: "approvers",
	}
}

// SubmitCabinetOfflineApprovalReq 提交机架下线审批请求参数结构体
type SubmitCabinetOfflineApprovalReq struct {
	SubmitApprovalCommon
	CabinetIDs []uint `json:"ids"` // 机架ID列表
}

type SubmitApprovalCommon struct {
	Approvers   []string           `json:"approvers"` // 审批人ID列表
	CurrentUser *model.CurrentUser `json:"-"`
	FrontData   string             `json:"front_data"` //前端传入的元数据
	Remark      string
	// 从UAM上获取用户信息的钩子
	GetEmailFromUAM myuser.GetEmailFromUAM `json:"-"`	
}

// ToJSON 将请求结构体转换成JSON
func (reqData *SubmitCabinetOfflineApprovalReq) ToJSON() []byte {
	b, _ := json.Marshal(reqData.CabinetIDs)
	return b
}

// ApproversJSON 返回审批人ID构成的JSON数组字符串
func (reqData *SubmitCabinetOfflineApprovalReq) ApproversJSON() []byte {
	b, _ := json.Marshal(reqData.Approvers)
	return b
}

// HooksJSON 返回审批步骤钩子构成的字符串
func (reqData *SubmitCabinetOfflineApprovalReq) HooksJSON() []byte {
	hooks := []*model.StepHook{
		{
			ID:              model.HookCabinetOffline, // 机架下线钩子
			Description:     "",
			ContinueOnError: true,
		},
	}
	b, _ := json.Marshal(hooks)
	return b
}

// FieldMap 字段映射
func (reqData *SubmitCabinetOfflineApprovalReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.CabinetIDs: "ids",
		&reqData.Approvers:  "approvers",
	}
}

// Validate 结构体数据校验
func (reqData *SubmitApprovalCommon) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if len(reqData.Approvers) > 1 && reqData.Approvers[0] == reqData.Approvers[1] {
		errs.Add([]string{"approvers"}, binding.BusinessError, "审批人和实施人不能相同")
		return errs
	}
	if len(reqData.Approvers) > 0 && reqData.CurrentUser.ID == reqData.Approvers[0] {
		errs.Add([]string{"approvers"}, binding.BusinessError, "审批人不能是当前登录用户")
		return errs
	}
	return errs
}

// SubmitIDCAbolishApprovalReq 提交数据中心裁撤请求体
type SubmitIDCAbolishApprovalReq struct {
	SubmitIDsApprovalReq
}

// Validate 结构体数据校验
func (reqData *SubmitIDCAbolishApprovalReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if errsRet := reqData.SubmitApprovalCommon.Validate(req, errs); len(errsRet) != 0 {
		return errsRet
	}
	repo, _ := middleware.RepoFromContext(req.Context())

	if len(reqData.IDs) > approvalLimit {
		errs.Add([]string{"ids"}, binding.BusinessError,
			fmt.Sprintf("提交数据量(%d)须不大于(%d)条", len(reqData.IDs), approvalLimit))
		return errs
	}
	// 校验数据中心下的机房是否'已裁撤'状态
	for _, id := range reqData.IDs {
		idc, err := repo.GetIDCByID(id)
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			errs.Add([]string{"ids"}, binding.SystemError, "系统内部发生错误")
			return errs
		}
		if idc == nil {
			errs.Add([]string{"ids"}, binding.BusinessError, fmt.Sprintf("数据中心(ID:%d)不存在", id))
			return errs
		}
		srvRooms, err := repo.GetServerRooms(&model.ServerRoomCond{IDCID: []uint{id}}, nil, nil)
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			errs.Add([]string{"ids"}, binding.SystemError, "系统内部发生错误")
			return errs
		}
		for _, srvRoom := range srvRooms {
			if srvRoom.Status != model.RoomStatAbolished {
				errs.Add([]string{"ids"}, binding.SystemError, fmt.Sprintf("数据中心(%s)申请裁撤失败，其下有'未裁撤'机房(%s)", idc.Name, srvRoom.Name))
				return errs
			}
		}
	}
	// TODO 校验审批人是否存在
	return errs
}

// SubmitIDCAbolishApproval 提交数据中心裁撤审批
func SubmitIDCAbolishApproval(log logger.Logger, repo model.Repo, conf *config.Config, reqData *SubmitIDCAbolishApprovalReq) (approvalID uint, err error) {
	now := time.Now()

	approval := model.Approval{
		Title:      fmt.Sprintf("%s的数据中心裁撤审批", reqData.CurrentUser.Name),
		Type:       model.ApprovalTypeIDCAbolish,
		Metadata:   string(reqData.ToJSON()),
		FrontData:  reqData.FrontData,
		Initiator:  reqData.CurrentUser.ID,
		Approvers:  string(reqData.ApproversJSON()),
		Remark:     reqData.Remark,
		StartTime:  &now,
		IsRejected: model.NO,
		Status:     model.ApprovalStatusApproval,
	}

	steps := make([]*model.ApprovalStep, 0, len(reqData.Approvers))
	for i, uid := range reqData.Approvers {
		step := &model.ApprovalStep{
			Approver: uid,
			Title:    approval.Title,
		}
		if i == 0 {
			step.StartTime = &now // 为审批单的第一个步骤设置开始时间
		}
		if i == len(reqData.Approvers)-1 { // 在审批最后一步添加钩子
			step.Hooks = string(reqData.HooksJSON(model.HookIDCAbolish))
		}
		steps = append(steps, step)
	}

	if err = repo.SubmitApproval(&approval, steps...); err != nil {
		return 0, err
	}
	// 获取审批人Email并调用邮件发送api
	var aprs []string
	err = json.Unmarshal([]byte(approval.Approvers), &aprs)
	if err != nil {
		log.Debug(err)
	}
	for k, _ := range aprs {
		loginID, _, email, _ := reqData.SubmitApprovalCommon.GetEmailFromUAM(aprs[k])
		if email != "" {
			sendmailreq := middleware.SendMailReq {
				From:    defaultMailFrom,
				To:      email,
				Title:	 approval.Title,
				Content: approval.Title,
				BodyFormat: "1",
				Priority: "0",
			}
			go middleware.SendMail(log, repo, conf, &sendmailreq)
		} else {
			log.Errorf("数据中心裁撤审批邮件通知失败: 用户 %s 的 Email 为空", loginID)
		}
	}
	// 邮件发送不影响审批单提交
	return approval.ID, nil
}

// SubmitServerRoomAbolishApprovalReq 提交机房管理单元裁撤请求体
type SubmitServerRoomAbolishApprovalReq struct {
	SubmitIDsApprovalReq
}

// Validate 结构体数据校验
func (reqData *SubmitServerRoomAbolishApprovalReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if errsRet := reqData.SubmitApprovalCommon.Validate(req, errs); len(errsRet) != 0 {
		return errsRet
	}
	repo, _ := middleware.RepoFromContext(req.Context())

	if len(reqData.IDs) > approvalLimit {
		errs.Add([]string{"ids"}, binding.BusinessError,
			fmt.Sprintf("提交数据量(%d)须不大于(%d)条", len(reqData.IDs), approvalLimit))
		return errs
	}
	// 校验机房下的机架是否'已下线'状态
	for _, id := range reqData.IDs {
		cabinets, err := repo.GetServerCabinets(&model.ServerCabinetCond{ServerRoomID: []uint{id}}, nil, nil)
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			errs.Add([]string{"ids"}, binding.SystemError, "系统内部发生错误")
			return errs
		}
		for _, c := range cabinets {
			if c.Status != model.CabinetStatOffline {
				errs.Add([]string{"ids"}, binding.SystemError, fmt.Sprintf("机房(ID:%d)有'未下线'机架(ID:%s)", id, c.Number))
				return errs
			}
		}
	}
	// TODO 校验审批人是否存在
	return errs
}

// SubmitServerRoomAbolishApproval 提交机房管理单元裁撤审批
func SubmitServerRoomAbolishApproval(log logger.Logger, repo model.Repo, conf *config.Config, reqData *SubmitServerRoomAbolishApprovalReq) (approvalID uint, err error) {
	now := time.Now()

	approval := model.Approval{
		Title:      fmt.Sprintf("%s的机房裁撤审批", reqData.CurrentUser.Name),
		Type:       model.ApprovalTypeServerRoomAbolish,
		Metadata:   string(reqData.ToJSON()),
		FrontData:  reqData.FrontData,
		Initiator:  reqData.CurrentUser.ID,
		Approvers:  string(reqData.ApproversJSON()),
		Remark:     reqData.Remark,
		StartTime:  &now,
		IsRejected: model.NO,
		Status:     model.ApprovalStatusApproval,
	}

	steps := make([]*model.ApprovalStep, 0, len(reqData.Approvers))
	for i, uid := range reqData.Approvers {
		step := &model.ApprovalStep{
			Approver: uid,
			Title:    approval.Title,
		}
		if i == 0 {
			step.StartTime = &now // 为审批单的第一个步骤设置开始时间
		}
		if i == len(reqData.Approvers)-1 { // 在审批最后一步添加钩子
			step.Hooks = string(reqData.HooksJSON(model.HookServerRoomAbolish))
		}
		steps = append(steps, step)
	}

	if err = repo.SubmitApproval(&approval, steps...); err != nil {
		return 0, err
	}
	// 获取审批人Email并调用邮件发送api
	var aprs []string
	err = json.Unmarshal([]byte(approval.Approvers), &aprs)
	if err != nil {
		log.Debug(err)
	}
	for k, _ := range aprs {
		loginID, _, email, _ := reqData.SubmitApprovalCommon.GetEmailFromUAM(aprs[k])
		if email != "" {
			sendmailreq := middleware.SendMailReq {
				From:    defaultMailFrom,
				To:      email,
				Title:	 approval.Title,
				Content: approval.Title,
				BodyFormat: "1",
				Priority: "0",
			}
			go middleware.SendMail(log, repo, conf, &sendmailreq)
		} else {
			log.Errorf("审批：%s  邮件通知失败: 用户 %s 的 Email 为空", approval.Title, loginID)
		}
	}
	// 邮件发送不影响审批单提交
	return approval.ID, nil
}

// SubmitNetAreaOfflineApprovalReq 提交网络区域下线请求体
type SubmitNetAreaOfflineApprovalReq struct {
	SubmitIDsApprovalReq
}

// Validate 结构体数据校验
func (reqData *SubmitNetAreaOfflineApprovalReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if errsRet := reqData.SubmitApprovalCommon.Validate(req, errs); len(errsRet) != 0 {
		return errsRet
	}
	//repo, _ := middleware.RepoFromContext(req.Context())

	if len(reqData.IDs) > approvalLimit {
		errs.Add([]string{"ids"}, binding.BusinessError,
			fmt.Sprintf("提交数据量(%d)须不大于(%d)条", len(reqData.IDs), approvalLimit))
		return errs
	}
	// 校验
	// TODO 待完善
	// TODO 校验审批人是否存在
	return errs
}

// SubmitNetAreaOfflineApproval 提交网络区域下线审批
func SubmitNetAreaOfflineApproval(log logger.Logger, repo model.Repo, conf *config.Config, reqData *SubmitNetAreaOfflineApprovalReq) (approvalID uint, err error) {
	now := time.Now()

	approval := model.Approval{
		Title:      fmt.Sprintf("%s的网络区域下线审批", reqData.CurrentUser.Name),
		Type:       model.ApprovalTypeNetAreaOffline,
		Metadata:   string(reqData.ToJSON()),
		FrontData:  reqData.FrontData,
		Initiator:  reqData.CurrentUser.ID,
		Approvers:  string(reqData.ApproversJSON()),
		Remark:     reqData.Remark,
		StartTime:  &now,
		IsRejected: model.NO,
		Status:     model.ApprovalStatusApproval,
	}

	steps := make([]*model.ApprovalStep, 0, len(reqData.Approvers))
	for i, uid := range reqData.Approvers {
		step := &model.ApprovalStep{
			Approver: uid,
			Title:    approval.Title,
		}
		if i == 0 {
			step.StartTime = &now // 为审批单的第一个步骤设置开始时间
		}
		if i == len(reqData.Approvers)-1 { // 在审批最后一步添加钩子
			step.Hooks = string(reqData.HooksJSON(model.HookNetAreaOffline))
		}
		steps = append(steps, step)
	}

	if err = repo.SubmitApproval(&approval, steps...); err != nil {
		return 0, err
	}
	// 获取审批人Email并调用邮件发送api
	var aprs []string
	err = json.Unmarshal([]byte(approval.Approvers), &aprs)
	if err != nil {
		log.Debug(err)
	}
	for k, _ := range aprs {
		loginID, _, email, _ := reqData.SubmitApprovalCommon.GetEmailFromUAM(aprs[k])
		if email != "" {
			sendmailreq := middleware.SendMailReq {
				From:    defaultMailFrom,
				To:      email,
				Title:	 approval.Title,
				Content: approval.Title,
				BodyFormat: "1",
				Priority: "0",
			}
			go middleware.SendMail(log, repo, conf, &sendmailreq)
		} else {
			log.Errorf("审批：%s  邮件通知失败: 用户 %s 的 Email 为空", approval.Title, loginID)
		}
	}
	// 邮件发送不影响审批单提交	
	return approval.ID, nil
}

// SubmitIPUnassignApprovalReq 提交IP回收请求体
type SubmitIPUnassignApprovalReq struct {
	SubmitIDsApprovalReq
}

// Validate 结构体数据校验
func (reqData *SubmitIPUnassignApprovalReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if errsRet := reqData.SubmitApprovalCommon.Validate(req, errs); len(errsRet) != 0 {
		return errsRet
	}
	repo, _ := middleware.RepoFromContext(req.Context())

	if len(reqData.IDs) > approvalLimit {
		errs.Add([]string{"ids"}, binding.BusinessError,
			fmt.Sprintf("提交数据量(%d)须不大于(%d)条", len(reqData.IDs), approvalLimit))
		return errs
	}
	// 校验
	for _, id := range reqData.IDs {
		ip, _ := repo.GetIPByID(id)
		if ip.IsUsed == model.IPNotUsed {
			errs.Add([]string{"ids"}, binding.BusinessError,
				fmt.Sprintf("IP(%s)未分配，无需释放", ip.IP))
			return errs
		}
	}
	// TODO 校验审批人是否存在
	return errs
}

// SubmitIPUnassignApproval 提交IP回收审批
func SubmitIPUnassignApproval(log logger.Logger, repo model.Repo, conf *config.Config, reqData *SubmitIPUnassignApprovalReq) (approvalID uint, err error) {
	now := time.Now()

	approval := model.Approval{
		Title:      fmt.Sprintf("%s的IP回收审批", reqData.CurrentUser.Name),
		Type:       model.ApprovalTypeIPUnassign,
		Metadata:   string(reqData.ToJSON()),
		FrontData:  reqData.FrontData,
		Initiator:  reqData.CurrentUser.ID,
		Approvers:  string(reqData.ApproversJSON()),
		Remark:     reqData.Remark,
		StartTime:  &now,
		IsRejected: model.NO,
		Status:     model.ApprovalStatusApproval,
	}

	steps := make([]*model.ApprovalStep, 0, len(reqData.Approvers))
	for i, uid := range reqData.Approvers {
		step := &model.ApprovalStep{
			Approver: uid,
			Title:    approval.Title,
		}
		if i == 0 {
			step.StartTime = &now // 为审批单的第一个步骤设置开始时间
		}
		if i == len(reqData.Approvers)-1 { // 在审批最后一步添加钩子
			step.Hooks = string(reqData.HooksJSON(model.HookIPUnassign))
		}
		steps = append(steps, step)
	}

	if err = repo.SubmitApproval(&approval, steps...); err != nil {
		return 0, err
	}
	// 获取审批人Email并调用邮件发送api
	var aprs []string
	err = json.Unmarshal([]byte(approval.Approvers), &aprs)
	if err != nil {
		log.Debug(err)
	}
	for k, _ := range aprs {
		loginID, _, email, _ := reqData.SubmitApprovalCommon.GetEmailFromUAM(aprs[k])
		if email != "" {
			sendmailreq := middleware.SendMailReq {
				From:    defaultMailFrom,
				To:      email,
				Title:	 approval.Title,
				Content: approval.Title,
				BodyFormat: "1",
				Priority: "0",
			}
			go middleware.SendMail(log, repo, conf, &sendmailreq)
		} else {
			log.Errorf("审批：%s  邮件通知失败: 用户 %s 的 Email 为空", approval.Title, loginID)
		}
	}
	// 邮件发送不影响审批单提交
	return approval.ID, nil
}

// SubmitDevicePowerOffApprovalReq 提交物理机关机请求体
type SubmitDevicePowerOffApprovalReq struct {
	SubmitDeviceRetirementApprovalReq
}

// Validate 结构体数据校验
func (reqData *SubmitDevicePowerOffApprovalReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if errsRet := reqData.SubmitApprovalCommon.Validate(req, errs); len(errsRet) != 0 {
		return errsRet
	}
	//repo, _ := middleware.RepoFromContext(req.Context())

	if len(reqData.SNs) > approvalLimit {
		errs.Add([]string{"ids"}, binding.BusinessError,
			fmt.Sprintf("提交数据量(%d)须不大于(%d)条", len(reqData.SNs), approvalLimit))
		return errs
	}
	// 校验
	//for _, id := range reqData.IDs {
	//	ip, _ := repo.GetIPByID(id)
	//	if ip.SN == "" {
	//		errs.Add([]string{"ids"}, binding.BusinessError,
	//			fmt.Sprintf("IP(%s)未分配，无需释放", ip.IP))
	//		return errs
	//	}
	//}
	// TODO 校验审批人是否存在
	return errs
}

// SubmitDevicePowerOffApproval 提交物理机关机审批
func SubmitDevicePowerOffApproval(log logger.Logger, repo model.Repo, conf *config.Config, reqData *SubmitDevicePowerOffApprovalReq) (approvalID uint, err error) {
	now := time.Now()

	approval := model.Approval{
		Title:      fmt.Sprintf("%s的物理机关机审批", reqData.CurrentUser.Name),
		Type:       model.ApprovalTypeDevPowerOff,
		Metadata:   string(reqData.ToJSON()),
		FrontData:  reqData.FrontData,
		Initiator:  reqData.CurrentUser.ID,
		Approvers:  string(reqData.ApproversJSON()),
		Remark:     reqData.Remark,
		StartTime:  &now,
		IsRejected: model.NO,
		Status:     model.ApprovalStatusApproval,
	}

	steps := make([]*model.ApprovalStep, 0, len(reqData.Approvers))
	for i, uid := range reqData.Approvers {
		step := &model.ApprovalStep{
			Approver: uid,
			Title:    approval.Title,
		}
		if i == 0 {
			step.StartTime = &now // 为审批单的第一个步骤设置开始时间
		}
		if i == len(reqData.Approvers)-1 { // 在审批最后一步添加钩子
			step.Hooks = string(reqData.HooksJSON(model.HookDevPowerOff))
		}
		steps = append(steps, step)
	}

	if err = repo.SubmitApproval(&approval, steps...); err != nil {
		return 0, err
	}
	// 获取审批人Email并调用邮件发送api
	var aprs []string
	err = json.Unmarshal([]byte(approval.Approvers), &aprs)
	if err != nil {
		log.Debug(err)
	}
	for k, _ := range aprs {
		loginID, _, email, _ := reqData.SubmitApprovalCommon.GetEmailFromUAM(aprs[k])
		if email != "" {
			sendmailreq := middleware.SendMailReq {
				From:    defaultMailFrom,
				To:      email,
				Title:	 approval.Title,
				Content: approval.Title,
				BodyFormat: "1",
				Priority: "0",
			}
			go middleware.SendMail(log, repo, conf, &sendmailreq)
		} else {
			log.Errorf("审批：%s  邮件通知失败: 用户 %s 的 Email 为空", approval.Title, loginID)
		}
	}
	// 邮件发送不影响审批单提交	
	return approval.ID, nil
}

// SubmitDeviceRestartApproval 提交物理机关机审批
func SubmitDeviceRestartApproval(log logger.Logger, repo model.Repo, conf *config.Config, reqData *SubmitDevicePowerOffApprovalReq) (approvalID uint, err error) {
	now := time.Now()

	approval := model.Approval{
		Title:      fmt.Sprintf("%s的物理机重启审批", reqData.CurrentUser.Name),
		Type:       model.ApprovalTypeDevRestart,
		Metadata:   string(reqData.ToJSON()),
		FrontData:  reqData.FrontData,
		Initiator:  reqData.CurrentUser.ID,
		Approvers:  string(reqData.ApproversJSON()),
		Remark:     reqData.Remark,
		StartTime:  &now,
		IsRejected: model.NO,
		Status:     model.ApprovalStatusApproval,
	}

	steps := make([]*model.ApprovalStep, 0, len(reqData.Approvers))
	for i, uid := range reqData.Approvers {
		step := &model.ApprovalStep{
			Approver: uid,
			Title:    approval.Title,
		}
		if i == 0 {
			step.StartTime = &now // 为审批单的第一个步骤设置开始时间
		}
		if i == len(reqData.Approvers)-1 { // 在审批最后一步添加钩子
			step.Hooks = string(reqData.HooksJSON(model.HookDevRestart))
		}
		steps = append(steps, step)
	}

	if err = repo.SubmitApproval(&approval, steps...); err != nil {
		return 0, err
	}
	return approval.ID, nil
}

// Validate 结构体数据校验
func (reqData *SubmitCabinetOfflineApprovalReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if errsRet := reqData.SubmitApprovalCommon.Validate(req, errs); len(errsRet) != 0 {
		return errsRet
	}
	repo, _ := middleware.RepoFromContext(req.Context())

	if len(reqData.CabinetIDs) > approvalLimit {
		errs.Add([]string{"ids"}, binding.BusinessError,
			fmt.Sprintf("提交数据量(%d)须不大于(%d)条", len(reqData.CabinetIDs), approvalLimit))
		return errs
	}
	// 校验机架是否存在且是否处于'关电'状态
	for _, id := range reqData.CabinetIDs {
		cabinet, err := repo.GetServerCabinetByID(id)
		if gorm.IsRecordNotFoundError(err) {
			errs.Add([]string{"ids"}, binding.BusinessError, fmt.Sprintf("机架(id=%d)不存在", id))
			return errs
		}
		if err != nil {
			errs.Add([]string{"ids"}, binding.SystemError, "系统内部发生错误")
			return errs
		}
		if cabinet.Status != model.CabinetStatEnabled {
			errs.Add([]string{"ids"}, binding.BusinessError, fmt.Sprintf("机架(%s)未启用", cabinet.Number))
			return errs
		}
		if cabinet.IsPowered == model.YES {
			errs.Add([]string{"ids"}, binding.BusinessError, fmt.Sprintf("请先对机架(%s)进行关电", cabinet.Number))
			return errs
		}
	}
	// TODO 校验审批人是否存在
	return errs
}

// SubmitCabinetOfflineApproval 提交机架下线审批
func SubmitCabinetOfflineApproval(log logger.Logger, repo model.Repo, conf *config.Config, reqData *SubmitCabinetOfflineApprovalReq) (approvalID uint, err error) {
	now := time.Now()

	approval := model.Approval{
		Title:      fmt.Sprintf("%s的机架下线审批", reqData.CurrentUser.Name),
		Type:       model.ApprovalTypeCabinetOffline,
		Metadata:   string(reqData.ToJSON()),
		FrontData:  reqData.FrontData,
		Initiator:  reqData.CurrentUser.ID,
		Approvers:  string(reqData.ApproversJSON()),
		Remark:     reqData.Remark,
		StartTime:  &now,
		IsRejected: model.NO,
		Status:     model.ApprovalStatusApproval,
	}

	steps := make([]*model.ApprovalStep, 0, len(reqData.Approvers))
	for i, uid := range reqData.Approvers {
		step := &model.ApprovalStep{
			Approver: uid,
			Title:    approval.Title,
		}
		if i == 0 {
			step.StartTime = &now // 为审批单的第一个步骤设置开始时间
		}
		if i == len(reqData.Approvers)-1 { // 在审批最后一步添加钩子
			step.Hooks = string(reqData.HooksJSON())
		}
		steps = append(steps, step)
	}

	if err = repo.SubmitApproval(&approval, steps...); err != nil {
		return 0, err
	}
	// 获取审批人Email并调用邮件发送api
	var aprs []string
	err = json.Unmarshal([]byte(approval.Approvers), &aprs)
	if err != nil {
		log.Debug(err)
	}
	for k, _ := range aprs {
		loginID, _, email, _ := reqData.SubmitApprovalCommon.GetEmailFromUAM(aprs[k])
		if email != "" {
			sendmailreq := middleware.SendMailReq {
				From:    defaultMailFrom,
				To:      email,
				Title:	 approval.Title,
				Content: approval.Title,
				BodyFormat: "1",
				Priority: "0",
			}
			go middleware.SendMail(log, repo, conf, &sendmailreq)
		} else {
			log.Errorf("审批：%s  邮件通知失败: 用户 %s 的 Email 为空", approval.Title, loginID)
		}
	}
	// 邮件发送不影响审批单提交	
	return approval.ID, nil
}

// SubmitCabinetPowerOffApprovalReq 提交机架关电审批请求参数结构体
type SubmitCabinetPowerOffApprovalReq struct {
	SubmitApprovalCommon
	CabinetIDs []uint `json:"ids"` // 机架ID列表
	//Approvers   []string           `json:"approvers"` // 审批人ID列表
	//CurrentUser *model.CurrentUser `json:"-"`
	//FrontData   string             `json:"front_data"` //前端传入的元数据
	//Remark      string
}

// ToJSON 将请求结构体转换成JSON
func (reqData *SubmitCabinetPowerOffApprovalReq) ToJSON() []byte {
	b, _ := json.Marshal(reqData.CabinetIDs)
	return b
}

// ApproversJSON 返回审批人ID构成的JSON数组字符串
func (reqData *SubmitCabinetPowerOffApprovalReq) ApproversJSON() []byte {
	b, _ := json.Marshal(reqData.Approvers)
	return b
}

// HooksJSON 返回审批步骤钩子构成的字符串
func (reqData *SubmitCabinetPowerOffApprovalReq) HooksJSON() []byte {
	hooks := []*model.StepHook{
		{
			ID:              model.HookCabinetPowerOff, // 机架关电钩子
			Description:     "",
			ContinueOnError: true,
		},
	}
	b, _ := json.Marshal(hooks)
	return b
}

// FieldMap 字段映射
func (reqData *SubmitCabinetPowerOffApprovalReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.CabinetIDs: "ids",
		&reqData.Approvers:  "approvers",
	}
}

// Validate 结构体数据校验
func (reqData *SubmitCabinetPowerOffApprovalReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if errsRet := reqData.SubmitApprovalCommon.Validate(req, errs); len(errsRet) != 0 {
		return errsRet
	}
	repo, _ := middleware.RepoFromContext(req.Context())
	if len(reqData.CabinetIDs) > approvalLimit {
		errs.Add([]string{"ids"}, binding.BusinessError,
			fmt.Sprintf("提交数据量(%d)须不大于(%d)条", len(reqData.CabinetIDs), approvalLimit))
		return errs
	}
	// 校验机架是否存在且是否处于'开电'状态
	for _, id := range reqData.CabinetIDs {
		cabinet, err := repo.GetServerCabinetByID(id)
		if gorm.IsRecordNotFoundError(err) {
			errs.Add([]string{"ids"}, binding.BusinessError, fmt.Sprintf("机架(id=%d)不存在", id))
			return errs
		}
		if err != nil {
			errs.Add([]string{"ids"}, binding.SystemError, "系统内部发生错误")
			return errs
		}
		if cabinet.IsPowered != model.YES {
			errs.Add([]string{"ids"}, binding.BusinessError, fmt.Sprintf("机架(%s)未开电", cabinet.Number))
			return errs
		}
		//检查同组tor下的2个机架是否关联有物理机设备，如果有则不允许关电
		// 检查自己的所有机位是否空闲
		selfUs, err := repo.GetServerUSiteByCond(&model.CombinedServerUSite{ServerCabinetID: []uint{id}}, nil, nil)
		if err != nil {
			errs.Add([]string{"ids"}, binding.BusinessError, fmt.Sprintf("检查机架(%s)所有机位占用情况失败", cabinet.Number))
			return errs
		}
		for _, u := range selfUs {
			if u.Status != model.USiteStatFree && u.Status != model.USiteStatDisabled {
				errs.Add([]string{"ids"}, binding.BusinessError, fmt.Sprintf("机架(%s)机位(%s)被占用", cabinet.Number, u.Number))
				return errs
			}
		}
		//检查同TOR下的另一个机架的所有机位是否空闲
		peerNetworkDevice, err := repo.GetPeerNetworkDeviceByCabinetID(id)
		if err != nil {
			//return ?
		} else if peerNetworkDevice != nil {
			peerUs, err := repo.GetServerUSiteByCond(&model.CombinedServerUSite{ServerCabinetID: []uint{peerNetworkDevice.ServerCabinetID}}, nil, nil)
			if err != nil {
				errs.Add([]string{"ids"}, binding.BusinessError, fmt.Sprintf("检查同组TOR的对端机架(%d)所有机位占用情况失败", peerNetworkDevice.ServerCabinetID))
				return errs
			}
			for _, u := range peerUs {
				if u.Status != model.USiteStatFree && u.Status != model.USiteStatDisabled {
					cabinet, _ := repo.GetServerCabinetByID(peerNetworkDevice.ServerCabinetID)
					if cabinet != nil {
						errs.Add([]string{"ids"}, binding.BusinessError, fmt.Sprintf("同组TOR的对端机架(%s)机位(%s)被占用", cabinet.Number, u.Number))
					} else {
						errs.Add([]string{"ids"}, binding.BusinessError, fmt.Sprintf("同组TOR的对端机架(ID:%d)机位(%s)被占用", peerNetworkDevice.ServerCabinetID, u.Number))
					}
					return errs
				}
			}
		}
	}
	// TODO 校验审批人是否存在
	return errs
}

// SubmitCabinetPowerOffApproval 提交机架关电审批
func SubmitCabinetPowerOffApproval(log logger.Logger, repo model.Repo, conf *config.Config, reqData *SubmitCabinetPowerOffApprovalReq) (approvalID uint, err error) {
	now := time.Now()

	approval := model.Approval{
		Title:      fmt.Sprintf("%s的机架关电审批", reqData.CurrentUser.Name),
		Type:       model.ApprovalTypeCabinetPowerOff,
		Metadata:   string(reqData.ToJSON()),
		FrontData:  reqData.FrontData,
		Initiator:  reqData.CurrentUser.ID,
		Approvers:  string(reqData.ApproversJSON()),
		Remark:     reqData.Remark,
		StartTime:  &now,
		IsRejected: model.NO,
		Status:     model.ApprovalStatusApproval,
	}

	steps := make([]*model.ApprovalStep, 0, len(reqData.Approvers))
	for i, uid := range reqData.Approvers {
		step := &model.ApprovalStep{
			Approver: uid,
			Title:    approval.Title, //"审批",
		}
		if i == 0 {
			step.StartTime = &now // 为审批单的第一个步骤设置开始时间
		}
		if i == len(reqData.Approvers)-1 { // 在审批最后一步添加钩子
			step.Hooks = string(reqData.HooksJSON())
		}
		steps = append(steps, step)
	}

	if err = repo.SubmitApproval(&approval, steps...); err != nil {
		return 0, err
	}
	// 获取审批人Email并调用邮件发送api
	var aprs []string
	err = json.Unmarshal([]byte(approval.Approvers), &aprs)
	if err != nil {
		log.Debug(err)
	}
	for k, _ := range aprs {
		loginID, _, email, _ := reqData.SubmitApprovalCommon.GetEmailFromUAM(aprs[k])
		if email != "" {
			sendmailreq := middleware.SendMailReq {
				From:    defaultMailFrom,
				To:      email,
				Title:	 approval.Title,
				Content: approval.Title,
				BodyFormat: "1",
				Priority: "0",
			}
			go middleware.SendMail(log, repo, conf, &sendmailreq)
		} else {
			log.Errorf("审批：%s  邮件通知失败: 用户 %s 的 Email 为空", approval.Title, loginID)
		}
	}
	// 邮件发送不影响审批单提交	
	return approval.ID, nil
}

// SubmitDeviceMigrationApprovalReq 提物理机搬迁审批请求参数结构体
type SubmitDeviceMigrationApprovalReq struct {
	SubmitApprovalCommon
	Data []*SubmitDeviceMigrationApprovalReqData `json:"data"` //搬迁的参数
	//Approvers   []string                               `json:"approvers"` // 审批人ID列表
	//CurrentUser *model.CurrentUser                     `json:"-"`
	//FrontData   string                                 `json:"front_data"` //前端传入的元数据
	//Remark      string
}

// SubmitDeviceMigrationApprovalReqData //搬迁的参数
type SubmitDeviceMigrationApprovalReqData struct {
	SN       string `json:"sn"`
	DstIDCID uint   `json:"dst_idc_id"` //目的数据中心

	//搬迁到机架，上架
	DstServerRoomID uint `json:"dst_server_room_id"` //目的机房
	DstCabinetID    uint `json:"dst_cabinet_id"`     //目的机架
	DstUSiteID      uint `json:"dst_usite_id"`       //目的机位

	//搬迁到货架，入库
	DstStoreRoomID uint `json:"dst_store_room_id"`      //目的库房
	DstVCabinetID  uint `json:"dst_virtual_cabinet_id"` //目的货架

	//搬迁的类型
	MigType string `json:"mig_type"` //库房->机架，机架->机架,机架->库房
}

// ToJSON 将请求结构体转换成JSON
func (reqData *SubmitDeviceMigrationApprovalReq) ToJSON() []byte {
	b, _ := json.Marshal(reqData.Data)
	return b
}

// ApproversJSON 返回审批人ID构成的JSON数组字符串
func (reqData *SubmitDeviceMigrationApprovalReq) ApproversJSON() []byte {
	b, _ := json.Marshal(reqData.Approvers)
	return b
}

// Step1HooksJSON 返回审批步骤钩子构成的字符串
func (reqData *SubmitDeviceMigrationApprovalReq) Step1HooksJSON() []byte {
	hooks := []*model.StepHook{
		{
			ID:              model.HookDeviceMigrationPowerOff, // 物理机搬迁关电钩子
			Description:     "",
			ContinueOnError: true,
		},
		{
			ID:              model.HookDeviceMigrationReserveIP, // 物理机搬迁保留IP钩子
			Description:     "",
			ContinueOnError: true,
		},
	}
	b, _ := json.Marshal(hooks)
	return b
}

// Step2HooksJSON 返回审批步骤钩子构成的字符串
func (reqData *SubmitDeviceMigrationApprovalReq) Step2HooksJSON() []byte {
	hooks := []*model.StepHook{
		{
			ID:              model.HookDeviceMigration, // 物理机搬迁钩子
			Description:     "",
			ContinueOnError: true,
		},
	}
	b, _ := json.Marshal(hooks)
	return b
}

// FieldMap 字段映射
func (reqData *SubmitDeviceMigrationApprovalReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.Data:      "data",
		&reqData.Approvers: "approvers",
	}
}

// Validate 结构体数据校验
func (reqData *SubmitDeviceMigrationApprovalReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if errsRet := reqData.SubmitApprovalCommon.Validate(req, errs); len(errsRet) != 0 {
		return errsRet
	}
	repo, _ := middleware.RepoFromContext(req.Context())
	if len(reqData.Data) > approvalLimit50 {
		errs.Add([]string{"data"}, binding.BusinessError,
			fmt.Sprintf("提交数据量(%d)须不大于(%d)条", len(reqData.Data), approvalLimit50))
		return errs
	}
	// 目标机位唯一性校验
	var dstUsiteIDsForCheck []uint
	for _, data := range reqData.Data {
		dstUsiteIDsForCheck = append(dstUsiteIDsForCheck, data.DstUSiteID)
		//目标机房是否存在
		//检查物理机的当前运营状态
		dev, err := repo.GetDeviceBySN(data.SN)
		if gorm.IsRecordNotFoundError(err) {
			errs.Add([]string{"sns"}, binding.BusinessError, fmt.Sprintf("设备(sn=%s)不存在", data.SN))
			return errs
		}
		if err != nil {
			errs.Add([]string{"sns"}, binding.SystemError, "系统内部发生错误")
			return errs
		}
		if dev.OperationStatus != model.DevOperStatPreDeploy &&
			dev.OperationStatus != model.DevOperStatOnShelve &&
			dev.OperationStatus != model.DevOperStatPreMove &&
			dev.OperationStatus != model.DevOperStatInStore {
			errs.Add([]string{"sns"}, binding.BusinessError,
				fmt.Sprintf("设备(sn=%s)运营状态必须是已上架/待部署/待搬迁/库房中(当前状态：%s)", data.SN,
					OperationStatusTransfer(dev.OperationStatus, true)))
			return errs
		}
		if data.DstServerRoomID != 0 { //搬迁到机架
			if dev.OperationStatus == model.DevOperStatInStore || dev.StoreRoomID != 0 {
				data.MigType = model.MigTypeStore2Usite //从库房到机架，此时不用埋关电和释放IP的钩子函数
			} else {
				data.MigType = model.MigTypeUsite2Usite
			}
			//检查目标机架是否存在
			c, err := repo.GetServerCabinetByID(data.DstCabinetID)
			if gorm.IsRecordNotFoundError(err) {
				errs.Add([]string{"ids"}, binding.BusinessError, fmt.Sprintf("机架(id=%d)不存在", data.DstUSiteID))
				return errs
			}
			if err != nil {
				errs.Add([]string{"ids"}, binding.SystemError, "系统内部发生错误")
				return errs
			}
			//检查目标机位是否存在且状态不为已使用或者不可用
			usite, err := repo.GetServerUSiteByID(data.DstUSiteID)
			if gorm.IsRecordNotFoundError(err) {
				errs.Add([]string{"ids"}, binding.BusinessError, fmt.Sprintf("机位(id=%d)不存在", data.DstUSiteID))
				return errs
			}
			if err != nil {
				errs.Add([]string{"ids"}, binding.SystemError, "系统内部发生错误")
				return errs
			}
			if usite.Status == model.USiteStatDisabled || usite.Status == model.USiteStatUsed {
				errs.Add([]string{"ids"}, binding.SystemError, fmt.Sprintf("机位%s(机架:%s)已被使用或不可用",
					usite.Number, c.Number))
				return errs
			}
		} else if data.DstStoreRoomID != 0 {
			if dev.OperationStatus == model.DevOperStatInStore || dev.StoreRoomID != 0 {
				data.MigType = model.MigTypeStore2Store
			} else {
				data.MigType = model.MigTypeUsite2Store
			}

			//检查库房单元是否存在
			_, err = repo.GetStoreRoomByID(data.DstStoreRoomID)
			if gorm.IsRecordNotFoundError(err) {
				errs.Add([]string{"ids"}, binding.BusinessError, fmt.Sprintf("库房单元(id=%d)不存在", data.DstStoreRoomID))
				return errs
			}
			if err != nil {
				errs.Add([]string{"ids"}, binding.SystemError, "系统内部发生错误")
				return errs
			}
			//检查虚拟货架是否存在
			_, err = repo.GetVirtualCabinetByID(data.DstVCabinetID)
			if gorm.IsRecordNotFoundError(err) {
				errs.Add([]string{"ids"}, binding.BusinessError, fmt.Sprintf("虚拟货架(id=%d)不存在", data.DstVCabinetID))
				return errs
			}
			if err != nil {
				errs.Add([]string{"ids"}, binding.SystemError, "系统内部发生错误")
				return errs
			}
		}
	}
	// 目标机位唯一性校验
	if len(dstUsiteIDsForCheck) >= 2 {
		for i := 0; i < len(dstUsiteIDsForCheck); i++ {
			for j := i + 1; j < len(dstUsiteIDsForCheck); j++ {
				if dstUsiteIDsForCheck[i] == dstUsiteIDsForCheck[j] {
					errs.Add([]string{"ids"}, binding.BusinessError, "存在重复的目标机位，请检查")
					return errs
				}
			}
		}
	}
	// TODO 校验审批人是否存在
	return errs
}

// SubmitDeviceMigrationApproval 提交物理机搬迁审批
func SubmitDeviceMigrationApproval(log logger.Logger, repo model.Repo, conf *config.Config, reqData *SubmitDeviceMigrationApprovalReq) (approvalID uint, err error) {
	now := time.Now()

	approval := model.Approval{
		Title:      fmt.Sprintf("%s的物理机搬迁审批", reqData.CurrentUser.Name),
		Type:       model.ApprovalTypeDeviceMigration,
		Metadata:   string(reqData.ToJSON()),
		FrontData:  reqData.FrontData,
		Initiator:  reqData.CurrentUser.ID,
		Approvers:  string(reqData.ApproversJSON()),
		Remark:     reqData.Remark,
		StartTime:  &now,
		IsRejected: model.NO,
		Status:     model.ApprovalStatusApproval,
	}

	steps := make([]*model.ApprovalStep, 0, len(reqData.Approvers))
	for i, uid := range reqData.Approvers {
		step := &model.ApprovalStep{
			Approver: uid,
			Title:    approval.Title,
		}
		if i == 0 {
			step.StartTime = &now // 为审批单的第一个步骤设置开始时间
			step.Hooks = string(reqData.Step1HooksJSON())
		}
		if i == len(reqData.Approvers)-1 { // 在审批最后一步添加钩子
			step.Hooks = string(reqData.Step2HooksJSON())
		}
		steps = append(steps, step)
	}

	if err = repo.SubmitApproval(&approval, steps...); err != nil {
		return 0, err
	}

	//预占用机位
	var usiteIDs []uint
	for _, dev := range reqData.Data {
		usiteIDs = append(usiteIDs, dev.DstUSiteID)
	}
	_, _ = repo.BatchUpdateServerUSitesStatus(usiteIDs, model.USiteStatPreOccupied)

	// 获取审批人Email并调用邮件发送api
	var aprs []string
	err = json.Unmarshal([]byte(approval.Approvers), &aprs)
	if err != nil {
		log.Debug(err)
	}
	for k, _ := range aprs {
		loginID, _, email, _ := reqData.SubmitApprovalCommon.GetEmailFromUAM(aprs[k])
		if email != "" {
			sendmailreq := middleware.SendMailReq {
				From:    defaultMailFrom,
				To:      email,
				Title:	 approval.Title,
				Content: approval.Title,
				BodyFormat: "1",
				Priority: "0",
			}
			go middleware.SendMail(log, repo, conf, &sendmailreq)
		} else {
			log.Errorf("审批：%s  邮件通知失败: 用户 %s 的 Email 为空", approval.Title, loginID)
		}
	}
	// 邮件发送不影响审批单提交	
	return approval.ID, nil
}

// ImportMigrationApprovalReq 导入待搬迁物理机的字段结构
type ImportMigrationApprovalReq struct {
	// 序列号
	SN string `json:"sn"`
	// 目标IDC
	DstIDCName string `json:"dst_idc_name"`
	// 目标机房
	DstServerRoomName string `json:"dst_server_room_name"`
	// 目标机架
	DstCabinetNum string `json:"dst_cabinet_number"`
	// 目标机位
	DstUSiteNum string `json:"dst_usite_number"`
	// 目标库房
	DstStoreRoom string `json:"dst_store_room_name"`
	// 目标虚拟货架
	DstVCabinetNum string `json:"dst_virtual_cabinet_number"`
	// 搬迁类型，
	MigType string `json:"mig_type"`

	ErrMsgContent string `json:"content"`

	//私有部分
	DstIDCID        uint                   `json:"dst_idc_id"`
	DstServerRoomID uint                   `json:"dst_server_room_id"`
	DstCabinetID    uint                   `json:"dst_cabinet_id"`
	DstUsiteID      uint                   `json:"dst_usite_id"`
	DstStoreRoomID  uint                   `json:"dst_store_room_id"`
	DstVCabinetID   uint                   `json:"dst_virtual_cabinet_id"`
	Fn              string                 `json:"fixed_asset_number"` //固资
	MModel          string                 `json:"model"`              //型号
	Category        string                 `json:"category"`           //类型
	OperationStatus string                 `json:"operation_status"`
	IDC             *IDCSimplify           `json:"idc"`
	ServerRoom      *ServerRoomSimplify    `json:"server_room"`
	ServerCabinet   *ServerCabinetSimplify `json:"server_cabinet"`
	ServerUSite     *ServerUSiteSimplify   `json:"server_usite"`
}

//checkLength 对导入文件中的数据做字段长度校验
func (impDevReq *ImportMigrationApprovalReq) checkLength() {
	leg := len(impDevReq.SN)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:序列号长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.DstIDCName)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:数据中心名称长度为(%d)(不能为空，不能大于255)", leg)
	}
	legServerRoom := len(impDevReq.DstServerRoomName)
	lenStoreRoom := len(impDevReq.DstStoreRoom)
	if (legServerRoom == 0 && lenStoreRoom == 0) || legServerRoom > 255 || lenStoreRoom > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + "必填项校验:机房管理单元、库房管理单元不能同时为空，不能大于255)"
	}
	if lenStoreRoom != 0 && legServerRoom != 0 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + "机房和库房只需二选一"
	}

	leg = len(impDevReq.DstCabinetNum)
	lenVCabinetNum := len(impDevReq.DstVCabinetNum)
	if (leg == 0 && lenVCabinetNum == 0) || leg > 255 || lenVCabinetNum > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + "必填项校验:机架编号、虚拟货架不能同时为空，且不能大于255"
	}
	if leg != 0 && lenVCabinetNum != 0 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + "机架和虚拟货架只需二选一"
	}

	leg = len(impDevReq.DstUSiteNum)
	if (legServerRoom != 0 && leg == 0) || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:机位编号长度为(%d)(不能为空，不能大于255)", leg)
	}
}

//validate 对导入文件中的数据做基本验证
func (impDevReq *ImportMigrationApprovalReq) validate(repo model.Repo) error {
	d, err := repo.GetDeviceBySN(impDevReq.SN)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	} else if err == gorm.ErrRecordNotFound || d == nil {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("设备(SN:%s)不存在", impDevReq.SN)
		return nil
	}
	if d != nil {
		impDevReq.Fn = d.FixedAssetNumber
		impDevReq.MModel = d.DevModel
		impDevReq.Category = d.Category
		impDevReq.OperationStatus = d.OperationStatus
		impDevReq.IDC = &IDCSimplify{ID: d.IDCID}
		impDevReq.ServerRoom = &ServerRoomSimplify{ID: d.ServerRoomID}
		impDevReq.ServerCabinet = &ServerCabinetSimplify{ID: d.CabinetID}
		impDevReq.ServerUSite = &ServerUSiteSimplify{}
		//conver2DevicePagesResp
		//idc
		if idc, err := repo.GetIDCByID(d.IDCID); err == nil {
			impDevReq.IDC.Name = idc.Name
		}

		if d.OperationStatus != model.DevOperStatInStore || d.StoreRoomID != 0 { //常规已上架的流程
			//server_room
			if room, err := repo.GetServerRoomByID(d.ServerRoomID); err == nil {
				impDevReq.ServerRoom.Name = room.Name
			}
			//server_cabinet
			if cabinet, err := repo.GetServerCabinetByID(d.CabinetID); err == nil {
				impDevReq.ServerCabinet.Number = cabinet.Number
			}
			//server_usite
			if d.USiteID != nil {
				if u, err := repo.GetServerUSiteByID(*d.USiteID); err == nil {
					impDevReq.ServerUSite.ID = u.ID
					impDevReq.ServerUSite.Number = u.Number
					impDevReq.ServerUSite.PhysicalArea = u.PhysicalArea
				}
			}
			//impDevReq.TOR, _ = repo.GetTORBySN(d.SN)

		} else { //库房中的
			if room, err := repo.GetStoreRoomByID(d.StoreRoomID); err == nil && room != nil {
				impDevReq.ServerRoom.ID = room.ID
				impDevReq.ServerRoom.Name = room.Name
			}
			//virtual_cabinet
			if cabinet, err := repo.GetVirtualCabinetByID(d.VCabinetID); err == nil && cabinet != nil {
				impDevReq.ServerCabinet.ID = cabinet.ID
				impDevReq.ServerCabinet.Number = cabinet.Number
			}

		}

	}

	//数据中心
	idc, err := repo.GetIDCByName(impDevReq.DstIDCName)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	} else if err == gorm.ErrRecordNotFound || idc == nil {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("数据中心(%s)不存在", impDevReq.DstIDCName)
	}
	if idc != nil {
		impDevReq.DstIDCID = idc.ID
	}

	//导入到机架
	if impDevReq.DstServerRoomName != "" {
		if d.OperationStatus == model.DevOperStatInStore {
			impDevReq.MigType = model.MigTypeStore2Usite //从库房到机架，此时不用埋关电和释放IP的钩子函数
		} else {
			impDevReq.MigType = model.MigTypeUsite2Usite
		}
		//机房校验
		srs, err := repo.GetServerRoomByName(impDevReq.DstServerRoomName)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == gorm.ErrRecordNotFound || srs == nil {
			var br string
			if impDevReq.ErrMsgContent != "" {
				br = "<br />"
			}
			impDevReq.ErrMsgContent += br + fmt.Sprintf("机房名(%s)不存在", impDevReq.DstServerRoomName)
		} else if srs.IDCID != idc.ID {
			var br string
			if impDevReq.ErrMsgContent != "" {
				br = "<br />"
			}
			impDevReq.ErrMsgContent += br + fmt.Sprintf("机房名(%s)不属于数据中心(%s)", impDevReq.DstServerRoomName, impDevReq.DstIDCName)
		} else {
			impDevReq.DstIDCID = srs.IDCID
			impDevReq.DstServerRoomID = srs.ID
		}

		//机架
		cabinet, err := repo.GetServerCabinetByNumber(impDevReq.DstServerRoomID, impDevReq.DstCabinetNum)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == gorm.ErrRecordNotFound || cabinet == nil {
			var br string
			if impDevReq.ErrMsgContent != "" {
				br = "<br />"
			}
			impDevReq.ErrMsgContent += br + fmt.Sprintf("机架编号(%s)不存在", impDevReq.DstCabinetNum)
		} else {
			impDevReq.DstCabinetID = cabinet.ID
		}
		//机位
		uSite, err := repo.GetServerUSiteByNumber(impDevReq.DstCabinetID, impDevReq.DstUSiteNum)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == gorm.ErrRecordNotFound || uSite == nil {
			var br string
			if impDevReq.ErrMsgContent != "" {
				br = "<br />"
			}
			impDevReq.ErrMsgContent += br + fmt.Sprintf("机位编号(%s)不存在, 机架编号(%s)", impDevReq.DstUSiteNum, impDevReq.DstCabinetNum)
		} else {
			if uSite.Status == model.USiteStatDisabled || uSite.Status == model.USiteStatUsed {
				var br string
				if impDevReq.ErrMsgContent != "" {
					br = "<br />"
				}
				impDevReq.ErrMsgContent += br + fmt.Sprintf("机位编号(%s)已使用或不可用", uSite.Number)
			}
			impDevReq.DstUsiteID = uSite.ID
		}

	} else {
		//导入到库房
		if d.OperationStatus == model.DevOperStatInStore {
			impDevReq.MigType = model.MigTypeStore2Store
		} else {
			impDevReq.MigType = model.MigTypeUsite2Store
		}
		//检查库房单元是否存在
		storeRoom, err := repo.GetStoreRoomByName(impDevReq.DstStoreRoom)
		if gorm.IsRecordNotFoundError(err) {
			var br string
			if impDevReq.ErrMsgContent != "" {
				br = "<br />"
			}
			impDevReq.ErrMsgContent += br + fmt.Sprintf("库房单元(%s)不存在", impDevReq.DstStoreRoom)
		}
		if err != nil {
			var br string
			if impDevReq.ErrMsgContent != "" {
				br = "<br />"
			}
			impDevReq.ErrMsgContent += br + "系统内部发生错误"
		}
		if storeRoom != nil {
			if storeRoom.IDCID != idc.ID {
				var br string
				if impDevReq.ErrMsgContent != "" {
					br = "<br />"
				}
				impDevReq.ErrMsgContent += br + fmt.Sprintf("库房(%s)不属于idc(%s)", storeRoom.Name, idc.Name)
			}
			impDevReq.DstStoreRoomID = storeRoom.ID
			//impDevReq.DstIDCID = storeRoom.IDCID
		}
		//检查虚拟货架是否存在
		vc, _ := repo.GetVirtualCabinets(&model.VirtualCabinet{Number: impDevReq.DstVCabinetNum}, nil, nil)
		if len(vc) == 0 {
			var br string
			if impDevReq.ErrMsgContent != "" {
				br = "<br />"
			}
			impDevReq.ErrMsgContent += br + fmt.Sprintf("虚拟货架(%s)不存在", impDevReq.DstVCabinetNum)
		} else {
			impDevReq.DstVCabinetID = vc[0].ID
		}
	}
	return nil
}

//ImportMigrationApprovalPriview 预览导入的待搬迁物理机
func ImportMigrationApprovalPriview(log logger.Logger, repo model.Repo, reqData *ImportPreviewReq) (map[string]interface{}, error) {
	fileName := upload.UploadDir + reqData.FileName
	ra, err := utils.ParseDataFromXLSX(fileName)
	if err != nil {
		return nil, err
	}
	length := len(ra)
	if length > approvalLimit50+1 {
		return nil, fmt.Errorf("导入数量不允许超过%d条", approvalLimit50)
	}

	var success []*ImportMigrationApprovalReq
	var failure []*ImportMigrationApprovalReq

	for i := 1; i < length; i++ {
		row := &ImportMigrationApprovalReq{}
		if len(ra[i]) < 7 {
			var br string
			if row.ErrMsgContent != "" {
				br = "<br />"
			}
			row.ErrMsgContent += br + "导入文件列长度不对（应为7列）"
			failure = append(failure, row)
			continue
		}
		row.SN = ra[i][0]
		row.DstIDCName = ra[i][1]
		row.DstServerRoomName = ra[i][2]
		row.DstCabinetNum = ra[i][3]
		row.DstUSiteNum = ra[i][4]
		row.DstStoreRoom = ra[i][5]
		row.DstVCabinetNum = ra[i][6]

		utils.StructTrimSpace(row)

		//字段存在性校验
		row.checkLength()

		//数据有效性校验
		err := row.validate(repo)
		if err != nil {
			return nil, err
		}

		if row.ErrMsgContent != "" {
			failure = append(failure, row)
		} else {
			success = append(success, row)
		}
	}

	var data []*ImportMigrationApprovalReq
	if len(failure) > 0 {
		data = failure
	} else {
		data = success
	}
	var result []*ImportMigrationApprovalReq
	for i := 0; i < len(data); i++ {
		if uint(i) >= reqData.Offset && uint(i) < (reqData.Offset+reqData.Limit) {
			result = append(result, data[i])
		}
	}
	if len(failure) > 0 {
		_ = os.Remove(fileName)
		return map[string]interface{}{"status": "failure",
			"message":       "导入服务器错误",
			"total_records": len(data),
			"content":       result,
		}, nil
	}
	return map[string]interface{}{"status": "success",
		"message":       "操作成功",
		"import_status": true,
		"total_records": len(data),
		"content":       result,
	}, nil
}

// ImportMigrationApproval 将导入的待搬迁机器生成审批单
func ImportMigrationApproval(log logger.Logger, repo model.Repo, conf *config.Config, user *model.CurrentUser, reqData *ImportApprovalReq) (approvalID uint, err error) {
	fileName := upload.UploadDir + reqData.FileName
	ra, err := utils.ParseDataFromXLSX(fileName)
	if err != nil {
		return
	}
	length := len(ra)

	//把临时文件删了
	err = os.Remove(fileName)
	if err != nil {
		log.Warnf("remove tmp file: %s fail", fileName)
		return
	}

	//front_data 这个数据是为了和前端提交的数据保持一致
	type FrontData struct {
		FN        string `json:"fixed_asset_number"`
		SN        string `json:"sn"`
		Model     string `json:"model"`
		IDC       string `json:"idc"`
		SR        string `json:"server_room_name"`
		Cabinet   string `json:"server_cabinet_number"`
		Usite     string `json:"server_usite_number"`
		StoreRoom string `json:"store_room_name"`
		VCabinet  string `json:"virtual_cabinet_number"`
	}
	fd := make([]*FrontData, 0, len(ra))

	var data []*SubmitDeviceMigrationApprovalReqData
	for i := 1; i < length; i++ {
		row := &ImportMigrationApprovalReq{}
		row.SN = ra[i][0]
		row.DstIDCName = ra[i][1]
		row.DstServerRoomName = ra[i][2]
		row.DstCabinetNum = ra[i][3]
		row.DstUSiteNum = ra[i][4]
		row.DstStoreRoom = ra[i][5]
		row.DstVCabinetNum = ra[i][6]

		utils.StructTrimSpace(row)

		_ = row.validate(repo)

		data = append(data, &SubmitDeviceMigrationApprovalReqData{
			SN:       row.SN,
			DstIDCID: row.DstIDCID,
			//导入到机房
			DstServerRoomID: row.DstServerRoomID,
			DstCabinetID:    row.DstCabinetID,
			DstUSiteID:      row.DstUsiteID,
			//导入到库房
			DstStoreRoomID: row.DstStoreRoomID,
			DstVCabinetID:  row.DstVCabinetID,
			MigType:        row.MigType,
		})

		fd = append(fd, &FrontData{
			FN:        row.Fn,
			SN:        row.SN,
			Model:     row.MModel,
			IDC:       row.DstIDCName,
			SR:        row.DstServerRoomName,
			Cabinet:   row.DstCabinetNum,
			Usite:     row.DstUSiteNum,
			StoreRoom: row.DstStoreRoom,
			VCabinet:  row.DstVCabinetNum,
		})
	}

	fdByte, err := json.Marshal(fd)
	if err != nil {
		return 0, err
	}
	newReqData := SubmitDeviceMigrationApprovalReq{
		SubmitApprovalCommon: SubmitApprovalCommon{
			CurrentUser: user,
			Approvers:   reqData.Approvers,
			FrontData:   string(fdByte),
			Remark:      "导入",
		},
		Data: data,
	}

	return SubmitDeviceMigrationApproval(log, repo, conf, &newReqData)
}

// SubmitDeviceRetirementApprovalReq 提物理机退役审批请求参数结构体
type SubmitDeviceRetirementApprovalReq struct {
	SubmitApprovalCommon
	SNs []string `json:"sns"`
	//Approvers   []string           `json:"approvers"` // 审批人ID列表
	//CurrentUser *model.CurrentUser `json:"-"`
	//FrontData   string             `json:"front_data"` //前端传入的元数据
	//Remark      string
}

// ToJSON 将请求结构体转换成JSON
func (reqData *SubmitDeviceRetirementApprovalReq) ToJSON() []byte {
	b, _ := json.Marshal(reqData.SNs)
	return b
}

// ApproversJSON 返回审批人ID构成的JSON数组字符串
func (reqData *SubmitDeviceRetirementApprovalReq) ApproversJSON() []byte {
	b, _ := json.Marshal(reqData.Approvers)
	return b
}

// ToJSON 将请求结构体转换成JSON
func (reqData *SubmitDeviceRetirementApprovalReq) HooksJSON(hook ...string) []byte {
	hooks := []*model.StepHook{
		{
			ID:              hook[0], // 物理机退役关电钩子
			Description:     "",
			ContinueOnError: true,
		},
	}
	b, _ := json.Marshal(hooks)
	return b
}

// Step1HooksJSON 返回审批步骤钩子构成的字符串
func (reqData *SubmitDeviceRetirementApprovalReq) Step1HooksJSON() []byte {
	hooks := []*model.StepHook{
		{
			ID:              model.HookDeviceRetirementPowerOff, // 物理机退役关电钩子
			Description:     "",
			ContinueOnError: true,
		},
		{
			ID:              model.HookDeviceRetirementReserveIP, // 物理机退役保留IP钩子
			Description:     "",
			ContinueOnError: true,
		},
	}
	b, _ := json.Marshal(hooks)
	return b
}

// Step2HooksJSON 返回审批步骤钩子构成的字符串
func (reqData *SubmitDeviceRetirementApprovalReq) Step2HooksJSON() []byte {
	hooks := []*model.StepHook{
		{
			ID:              model.HookDeviceRetirement, // 物理机退役修改状态钩子
			Description:     "",
			ContinueOnError: true,
		},
	}
	b, _ := json.Marshal(hooks)
	return b
}

// FieldMap 字段映射
func (reqData *SubmitDeviceRetirementApprovalReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.SNs:       "sns",
		&reqData.Approvers: "approvers",
	}
}

// Validate 结构体数据校验
func (reqData *SubmitDeviceRetirementApprovalReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if errsRet := reqData.SubmitApprovalCommon.Validate(req, errs); len(errsRet) != 0 {
		return errsRet
	}
	repo, _ := middleware.RepoFromContext(req.Context())
	if len(reqData.SNs) > approvalLimit50 {
		errs.Add([]string{"sns"}, binding.BusinessError,
			fmt.Sprintf("提交数据量(%d)须不大于(%d)条", len(reqData.SNs), approvalLimit50))
		return errs
	}
	for _, sn := range reqData.SNs {
		//检查物理机的当前运营状态
		dev, err := repo.GetDeviceBySN(sn)
		if gorm.IsRecordNotFoundError(err) {
			errs.Add([]string{"sns"}, binding.BusinessError, fmt.Sprintf("设备(sn=%s)不存在", sn))
			return errs
		}
		if err != nil {
			errs.Add([]string{"sns"}, binding.SystemError, "系统内部发生错误")
			return errs
		}
		if dev.OperationStatus != model.DevOperStatPreRetire {
			errs.Add([]string{"sns"}, binding.BusinessError,
				fmt.Sprintf("设备(sn=%s)运营状态必须是待退役(当前状态：%s)", sn, OperationStatusTransfer(dev.OperationStatus, true)))
			return errs
		}
	}
	// TODO 校验审批人是否存在
	return errs
}

// SubmitDeviceRetirementApproval 提交物理机退役审批
func SubmitDeviceRetirementApproval(log logger.Logger, repo model.Repo, conf *config.Config, reqData *SubmitDeviceRetirementApprovalReq) (approvalID uint, err error) {
	now := time.Now()

	approval := model.Approval{
		Title:      fmt.Sprintf("%s的物理机退役审批", reqData.CurrentUser.Name),
		Type:       model.ApprovalTypeDeviceRetirement,
		Metadata:   string(reqData.ToJSON()),
		FrontData:  reqData.FrontData,
		Initiator:  reqData.CurrentUser.ID,
		Approvers:  string(reqData.ApproversJSON()),
		Remark:     reqData.Remark,
		StartTime:  &now,
		IsRejected: model.NO,
		Status:     model.ApprovalStatusApproval,
	}

	steps := make([]*model.ApprovalStep, 0, len(reqData.Approvers))
	for i, uid := range reqData.Approvers {
		step := &model.ApprovalStep{
			Approver: uid,
			Title:    approval.Title,
		}
		if i == 0 {
			step.StartTime = &now // 为审批单的第一个步骤设置开始时间
			step.Hooks = string(reqData.Step1HooksJSON())
		}
		if i == len(reqData.Approvers)-1 { // 在审批最后一步添加钩子
			step.Hooks = string(reqData.Step2HooksJSON())
		}
		steps = append(steps, step)
	}

	if err = repo.SubmitApproval(&approval, steps...); err != nil {
		return 0, err
	}
	// 获取审批人Email并调用邮件发送api
	var aprs []string
	err = json.Unmarshal([]byte(approval.Approvers), &aprs)
	if err != nil {
		log.Debug(err)
	}
	for k, _ := range aprs {
		loginID, _, email, _ := reqData.SubmitApprovalCommon.GetEmailFromUAM(aprs[k])
		if email != "" {
			sendmailreq := middleware.SendMailReq {
				From:    defaultMailFrom,
				To:      email,
				Title:	 approval.Title,
				Content: approval.Title,
				BodyFormat: "1",
				Priority: "0",
			}
			go middleware.SendMail(log, repo, conf, &sendmailreq)
		} else {
			log.Errorf("审批：%s  邮件通知失败: 用户 %s 的 Email 为空", approval.Title, loginID)
		}
	}
	// 邮件发送不影响审批单提交
	return approval.ID, nil
}

// SubmitDeviceReInstallationApprovalReq 提物理机重装审批请求参数结构体
type SubmitDeviceReInstallationApprovalReq struct {
	SubmitApprovalCommon
	Settings OSReinstallSetting `json:"settings"`
}

// ToJSON 将请求结构体转换成JSON
func (reqData *SubmitDeviceReInstallationApprovalReq) ToJSON() []byte {
	b, _ := json.Marshal(reqData.Settings)
	return b
}

// ApproversJSON 返回审批人ID构成的JSON数组字符串
func (reqData *SubmitDeviceReInstallationApprovalReq) ApproversJSON() []byte {
	b, _ := json.Marshal(reqData.Approvers)
	return b
}

// HooksJSON 返回审批步骤钩子构成的字符串
func (reqData *SubmitDeviceReInstallationApprovalReq) HooksJSON() []byte {
	hooks := []*model.StepHook{
		{
			ID:              model.HookDeviceOSReinstallation, // 物理机重装钩子
			Description:     "",
			ContinueOnError: true,
		},
	}
	b, _ := json.Marshal(hooks)
	return b
}

// FieldMap 字段映射
func (reqData *SubmitDeviceReInstallationApprovalReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.Settings:  "settings",
		&reqData.Approvers: "approvers",
	}
}

// Validate 结构体数据校验
func (reqData *SubmitDeviceReInstallationApprovalReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if errsRet := reqData.SubmitApprovalCommon.Validate(req, errs); len(errsRet) != 0 {
		return errsRet
	}
	repo, _ := middleware.RepoFromContext(req.Context())
	if len(reqData.Settings) > approvalLimit50 {
		errs.Add([]string{"settings"}, binding.BusinessError,
			fmt.Sprintf("提交数据量(%d)须不大于(%d)条", len(reqData.Settings), approvalLimit50))
		return errs
	}
	for _, sett := range reqData.Settings {
		//检查物理机的当前运营状态
		dev, err := repo.GetDeviceBySN(sett.SN)
		if gorm.IsRecordNotFoundError(err) {
			errs.Add([]string{"sns"}, binding.BusinessError, fmt.Sprintf("设备(sn=%s)不存在", sett.SN))
			return errs
		}
		if err != nil {
			errs.Add([]string{"sns"}, binding.SystemError, "系统内部发生错误")
			return errs
		}
		if dev.OperationStatus != model.DevOperStatRunWithAlarm &&
			dev.OperationStatus != model.DevOperStatRunWithoutAlarm &&
			dev.OperationStatus != model.DevOperStatOnShelve {
			errs.Add([]string{"sns"}, binding.BusinessError,
				fmt.Sprintf("设备(sn=%s)运营状态必须是运营中或已上架(当前状态：%s)", sett.SN,
					OperationStatusTransfer(dev.OperationStatus, true)))
			return errs
		}
		sett.OriginStatus = dev.OperationStatus
	}
	// 做参数校验，不然出了错误会把人整死的
	errs = reqData.Settings.Validate(req, errs)
	// TODO 校验审批人是否存在
	return errs
}

// SubmitDeviceReInstallationApproval 提交物理机重装审批
func SubmitDeviceReInstallationApproval(log logger.Logger, repo model.Repo, conf *config.Config, reqData *SubmitDeviceReInstallationApprovalReq) (approvalID uint, err error) {
	now := time.Now()

	approval := model.Approval{
		Title:      fmt.Sprintf("%s的物理机重装审批", reqData.CurrentUser.Name),
		Type:       model.ApprovalTypeDeviceOSReinstallation,
		Metadata:   string(reqData.ToJSON()),
		FrontData:  reqData.FrontData,
		Initiator:  reqData.CurrentUser.ID,
		Approvers:  string(reqData.ApproversJSON()),
		Remark:     reqData.Remark,
		StartTime:  &now,
		IsRejected: model.NO,
		Status:     model.ApprovalStatusApproval,
	}

	steps := make([]*model.ApprovalStep, 0, len(reqData.Approvers))
	for i, uid := range reqData.Approvers {
		step := &model.ApprovalStep{
			Approver: uid,
			Title:    approval.Title,
		}
		if i == 0 {
			step.StartTime = &now // 为审批单的第一个步骤设置开始时间
		}
		if i == len(reqData.Approvers)-1 { // 在审批最后一步添加钩子
			step.Hooks = string(reqData.HooksJSON())
		}
		steps = append(steps, step)
	}

	if err = repo.SubmitApproval(&approval, steps...); err != nil {
		return 0, err
	}
	// 获取审批人Email并调用邮件发送api
	var aprs []string
	err = json.Unmarshal([]byte(approval.Approvers), &aprs)
	if err != nil {
		log.Debug(err)
	}
	for k, _ := range aprs {
		loginID, _, email, _ := reqData.SubmitApprovalCommon.GetEmailFromUAM(aprs[k])
		if email != "" {
			sendmailreq := middleware.SendMailReq {
				From:    defaultMailFrom,
				To:      email,
				Title:	 approval.Title,
				Content: approval.Title,
				BodyFormat: "1",
				Priority: "0",
			}
			go middleware.SendMail(log, repo, conf, &sendmailreq)
		} else {
			log.Errorf("审批：%s  邮件通知失败: 用户 %s 的 Email 为空", approval.Title, loginID)
		}
	}
	// 邮件发送不影响审批单提交
	return approval.ID, nil
}

// SubmitDeviceRecycleApprovalReq 回收审批的请求结构
type SubmitDeviceRecycleApprovalReq struct {
	SubmitDeviceRetirementApprovalReq
	ApprovalType string             `json:"approval_type"` //device_recycle_pre_move|device_recycle_pre_retire|device_recycle_reinstall
	Settings     OSReinstallSetting `json:"settings"`      //当ApprovalType=reinstall时，需要指定此参数
}

type OSReinstallSetting []*OSReinstallSettingItem

type OSReinstallSettingItem struct {
	SaveDeviceSettingItem
	OriginStatus string `json:"origin_status"` //记录初始的设备运行状态，以便还原
}

// Validate 装机参数校验
func (reqData OSReinstallSetting) Validate(r *http.Request, errs binding.Errors) binding.Errors {
	for i := range reqData {
		if errs = reqData[i].validateOne(r, errs); errs.Len() > 0 {
			return errs
		}
	}
	return errs
}

// ToJSON 将请求结构体转换成JSON
func (reqData *SubmitDeviceRecycleApprovalReq) ToJSON() (b []byte) {
	switch reqData.ApprovalType {
	case model.ApprovalTypeDeviceRecycleReinstall:
		b, _ = json.Marshal(reqData.Settings)
	case model.ApprovalTypeDeviceRecyclePreMove:
		fallthrough
	case model.ApprovalTypeDeviceRecyclePreRetire:
		b, _ = json.Marshal(reqData.SNs)
	}

	return b
}

// ApproversJSON 返回审批人ID构成的JSON数组字符串
func (reqData *SubmitDeviceRecycleApprovalReq) ApproversJSON() []byte {
	b, _ := json.Marshal(reqData.Approvers)
	return b
}

// HooksJSON 返回审批步骤钩子构成的字符串
func (reqData *SubmitDeviceRecycleApprovalReq) HooksJSON() (b []byte) {
	switch reqData.ApprovalType {
	case model.ApprovalTypeDeviceRecycleReinstall:
		hooks := []*model.StepHook{
			{
				ID:              model.HookDeviceRecycleReinstall, // 物理机回收重装钩子
				Description:     "",
				ContinueOnError: true,
			},
		}
		b, _ = json.Marshal(hooks)
	case model.ApprovalTypeDeviceRecyclePreMove:
		hooks := []*model.StepHook{
			{
				ID:              model.HookDeviceRecyclePreMove, // 物理机回收待搬迁钩子
				Description:     "",
				ContinueOnError: true,
			},
		}
		b, _ = json.Marshal(hooks)
	case model.ApprovalTypeDeviceRecyclePreRetire:
		hooks := []*model.StepHook{
			{
				ID:              model.HookDeviceRecyclePreRetire, // 物理机回收待退役钩子
				Description:     "",
				ContinueOnError: true,
			},
		}
		b, _ = json.Marshal(hooks)
	}

	return b
}

// FieldMap 字段映射
func (reqData *SubmitDeviceRecycleApprovalReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.Settings:     "settings",
		&reqData.Approvers:    "approvers",
		&reqData.ApprovalType: "approval_type",
		&reqData.SNs:          "sns",
	}
}

// Validate 结构体数据校验
func (reqData *SubmitDeviceRecycleApprovalReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if errsRet := reqData.SubmitApprovalCommon.Validate(req, errs); len(errsRet) != 0 {
		return errsRet
	}
	if reqData.ApprovalType != model.ApprovalTypeDeviceRecycleReinstall &&
		reqData.ApprovalType != model.ApprovalTypeDeviceRecyclePreRetire &&
		reqData.ApprovalType != model.ApprovalTypeDeviceRecyclePreMove {
		errs.Add([]string{"approval_type"}, binding.BusinessError,
			fmt.Sprintf("审批类型值(%s)不合法", reqData.ApprovalType))
		return errs
	}
	repo, _ := middleware.RepoFromContext(req.Context())
	if reqData.ApprovalType == model.ApprovalTypeDeviceRecycleReinstall &&
		len(reqData.Settings) > approvalLimit50 {
		errs.Add([]string{"settings"}, binding.BusinessError,
			fmt.Sprintf("提交数据量(%d)须不大于(%d)条", len(reqData.Settings), approvalLimit50))
		return errs
	} else if reqData.ApprovalType != model.ApprovalTypeDeviceRecycleReinstall &&
		len(reqData.SNs) > approvalLimit50 {
		errs.Add([]string{"settings"}, binding.BusinessError,
			fmt.Sprintf("提交数据量(%d)须不大于(%d)条", len(reqData.Settings), approvalLimit50))
		return errs
	}
	for _, sett := range reqData.Settings {
		//检查物理机的当前运营状态
		dev, err := repo.GetDeviceBySN(sett.SN)
		if gorm.IsRecordNotFoundError(err) {
			errs.Add([]string{"sns"}, binding.BusinessError, fmt.Sprintf("设备(sn=%s)不存在", sett.SN))
			return errs
		}
		if err != nil {
			errs.Add([]string{"sns"}, binding.SystemError, "系统内部发生错误")
			return errs
		}
		if dev.OperationStatus != model.DevOperStatRunWithAlarm &&
			dev.OperationStatus != model.DevOperStatRunWithoutAlarm &&
			dev.OperationStatus != model.DevOperStatOnShelve &&
			dev.OperationStatus != model.DevOperStatRecycling {
			errs.Add([]string{"sns"}, binding.BusinessError,
				fmt.Sprintf("设备(sn=%s)运营状态必须是运营中或已上架或回收中(当前状态：%s)", sett.SN,
					OperationStatusTransfer(dev.OperationStatus, true)))
			return errs
		}
		sett.OriginStatus = dev.OperationStatus
	}
	// 做参数校验，不然出了错误会把人整死的
	if reqData.ApprovalType == model.ApprovalTypeDeviceRecycleReinstall {
		errs = reqData.Settings.Validate(req, errs)
	}
	// TODO 校验审批人是否存在
	return errs
}

// SubmitDeviceRecycleApproval 提交物理机回收审批
func SubmitDeviceRecycleApproval(log logger.Logger, repo model.Repo, conf *config.Config, reqData *SubmitDeviceRecycleApprovalReq) (approvalID uint, err error) {
	now := time.Now()

	approval := model.Approval{
		//Title:      fmt.Sprintf("%s的物理机回收审批", reqData.CurrentUser.Name),
		Type:       reqData.ApprovalType,
		Metadata:   string(reqData.ToJSON()),
		FrontData:  reqData.FrontData,
		Initiator:  reqData.CurrentUser.ID,
		Approvers:  string(reqData.ApproversJSON()),
		Remark:     reqData.Remark,
		StartTime:  &now,
		IsRejected: model.NO,
		Status:     model.ApprovalStatusApproval,
	}

	switch reqData.ApprovalType {
	case model.ApprovalTypeDeviceRecycleReinstall:
		approval.Title = fmt.Sprintf("%s的物理机回收重装审批", reqData.CurrentUser.Name)
	case model.ApprovalTypeDeviceRecyclePreMove:
		approval.Title = fmt.Sprintf("%s的物理机回收待搬迁审批", reqData.CurrentUser.Name)
	case model.ApprovalTypeDeviceRecyclePreRetire:
		approval.Title = fmt.Sprintf("%s的物理机回收待退役审批", reqData.CurrentUser.Name)
	}
	steps := make([]*model.ApprovalStep, 0, len(reqData.Approvers))
	for i, uid := range reqData.Approvers {
		step := &model.ApprovalStep{
			Approver: uid,
			Title:    approval.Title,
		}
		if i == 0 {
			step.StartTime = &now // 为审批单的第一个步骤设置开始时间
		}
		if i == len(reqData.Approvers)-1 { // 在审批最后一步添加钩子
			step.Hooks = string(reqData.HooksJSON())
		}
		steps = append(steps, step)
	}

	if err = repo.SubmitApproval(&approval, steps...); err != nil {
		return 0, err
	}
	// 获取审批人Email并调用邮件发送api
	var aprs []string
	err = json.Unmarshal([]byte(approval.Approvers), &aprs)
	if err != nil {
		log.Debug(err)
	}
	for k, _ := range aprs {
		loginID, _, email, _ := reqData.SubmitApprovalCommon.GetEmailFromUAM(aprs[k])
		if email != "" {
			sendmailreq := middleware.SendMailReq {
				From:    defaultMailFrom,
				To:      email,
				Title:	 approval.Title,
				Content: approval.Title,
				BodyFormat: "1",
				Priority: "0",
			}
			go middleware.SendMail(log, repo, conf, &sendmailreq)
		} else {
			log.Errorf("审批：%s  邮件通知失败: 用户 %s 的 Email 为空", approval.Title, loginID)
		}
	}
	// 邮件发送不影响审批单提交	
	return approval.ID, nil
}

// GetMyApprovalPageReq 获取我发起的审批单请求参数
type GetMyApprovalPageReq struct {
	// 申请单类型
	Type string `json:"type"`
	// 申请单状态
	Status string `json:"status"`
	// 分页页号
	Page int64 `json:"page"`
	// 分页大小。默认值:10。阈值: 100。
	PageSize int64 `json:"page_size"`

	// 当前用户ID
	CurrentUserID string `json:"-"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *GetMyApprovalPageReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.Type:     "type",
		&reqData.Status:   "status",
		&reqData.Page:     "page",
		&reqData.PageSize: "page_size",
	}
}

func checkApprovalType(at string) bool {
	if at != "" {
		// 这里的校验要与model类型保持一致
		switch at {
		case model.ApprovalTypeIDCAbolish:
		case model.ApprovalTypeServerRoomAbolish:
		case model.ApprovalTypeNetAreaOffline:
		case model.ApprovalTypeIPUnassign:
		case model.ApprovalTypeDevPowerOff:
		case model.ApprovalTypeDevRestart:
		case model.ApprovalTypeCabinetPowerOff:
		case model.ApprovalTypeCabinetOffline:
		case model.ApprovalTypeDeviceOSReinstallation:
		case model.ApprovalTypeDeviceMigration:
		case model.ApprovalTypeDeviceRetirement:
		case model.ApprovalTypeDeviceRecycleReinstall:
		case model.ApprovalTypeDeviceRecyclePreMove:
		case model.ApprovalTypeDeviceRecyclePreRetire:
		default:
			return false
		}
	}
	return true
}

func checkApprovalStatus(as string) bool {
	if as != "" {
		// 这里的校验要与model状态保持一致
		switch as {
		case model.ApprovalStatusApproval:
		case model.ApprovalStatusCompleted:
		case model.ApprovalStatusRevoked:
		default:
			return false
		}
	}
	return true

}

// Validate 结构体数据校验
func (reqData *GetMyApprovalPageReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	user, _ := middleware.LoginUserFromContext(req.Context())
	if user != nil {
		reqData.CurrentUserID = user.ID
	}

	if !checkApprovalType(reqData.Type) {
		errs.Add([]string{"type"}, binding.RequiredError, "审批单类型值错误")
		return errs
	}

	if !checkApprovalStatus(reqData.Status) {
		errs.Add([]string{"status"}, binding.RequiredError, "审批单状态值错误")
		return errs
	}

	return errs
}

// GetMyApprovalPage 获取我发起的申请单分页
func GetMyApprovalPage(log logger.Logger, repo model.Repo, reqData *GetMyApprovalPageReq) (*page.Page, error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 100 {
		reqData.PageSize = 20
	}
	if reqData.Page < 0 {
		reqData.Page = 0
	}

	cond := model.Approval{
		Type:   reqData.Type,
		Status: reqData.Status,
	}

	totalRecords, err := repo.CountInitiatedApprovals(reqData.CurrentUserID, &cond)
	if err != nil {
		return nil, err
	}

	pager := page.NewPager(reflect.TypeOf(&ApprovalPage{}), reqData.Page, reqData.PageSize, totalRecords)
	items, err := repo.GetInitiatedApprovals(reqData.CurrentUserID, &cond, model.OneOrderBy("id", model.DESC), pager.BuildLimiter())
	if err != nil {
		return nil, err
	}
	for i := range items {
		item, err := convert2MyApprovalPageResult(log, repo, items[i])
		if err != nil {
			return nil, err
		}
		if item != nil {
			pager.AddRecords(item)
		}
	}
	return pager.BuildPage(), nil
}

// ApprovalPage 我发起的审批分页查询信息
type ApprovalPage struct {
	// 审批单ID
	ID uint `json:"id"`
	// 审批类型
	Type string `json:"type"`
	// 申请单标题
	Title string `json:"title"`
	// 审批发起时间
	StartTime times.ISOTime `json:"start_time"`
	// 审批结束时间
	EndTime times.ISOTime `json:"end_time"`
	// 审批是否被拒绝
	IsRejected string `json:"is_rejected"`
	// 审批状态
	Status string `json:"status"`
}

func convert2MyApprovalPageResult(log logger.Logger, repo model.Repo, appr *model.Approval) (*ApprovalPage, error) {
	if appr == nil {
		return nil, nil
	}
	ap := &ApprovalPage{
		ID:         appr.ID,
		Title:      appr.Title,
		Type:       appr.Type,
		Status:     appr.Status,
		IsRejected: appr.IsRejected,
	}
	if appr.StartTime != nil {
		ap.StartTime = times.ISOTime(*appr.StartTime)
	}
	if appr.EndTime != nil {
		ap.EndTime = times.ISOTime(*appr.EndTime)
	}
	return ap, nil
}

// GetApproveByMePageReq 获取待我审批请求参数
type GetApproveByMePageReq struct {
	// 申请单类型
	Type string `json:"type"`
	// 分页页号
	Page int64 `json:"page"`
	// 分页大小。默认值:10。阈值: 100。
	PageSize int64 `json:"page_size"`

	// 当前用户ID
	CurrentUserID string `json:"-"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *GetApproveByMePageReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.Type:     "type",
		&reqData.Page:     "page",
		&reqData.PageSize: "page_size",
	}
}

// Validate 结构体数据校验
func (reqData *GetApproveByMePageReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	user, _ := middleware.LoginUserFromContext(req.Context())
	if user != nil {
		reqData.CurrentUserID = user.ID
	}

	if !checkApprovalType(reqData.Type) {
		errs.Add([]string{"type"}, binding.RequiredError, "审批单类型值错误")
		return errs
	}
	return errs
}

// GetApproveByMePage 获取待我审批的申请单分页
func GetApproveByMePage(log logger.Logger, repo model.Repo, reqData *GetApproveByMePageReq) (*page.Page, error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 100 {
		reqData.PageSize = 20
	}
	if reqData.Page < 0 {
		reqData.Page = 0
	}

	cond := model.Approval{
		Type:   reqData.Type,
		Status: model.ApprovalStatusApproval,
	}

	totalRecords, err := repo.CountPendingApprovals(reqData.CurrentUserID, &cond)
	if err != nil {
		return nil, err
	}

	pager := page.NewPager(reflect.TypeOf(&ApprovalPage{}), reqData.Page, reqData.PageSize, totalRecords)
	items, err := repo.GetPendingApprovals(reqData.CurrentUserID, &cond, model.OneOrderBy("id", model.DESC), pager.BuildLimiter())
	if err != nil {
		return nil, err
	}

	for i := range items {
		ap := ApprovalPage{
			ID:         items[i].ID,
			Title:      items[i].Title,
			Type:       items[i].Type,
			IsRejected: items[i].IsRejected,
			Status:     items[i].Status,
		}
		if items[i].StartTime != nil {
			ap.StartTime = times.ISOTime(*items[i].StartTime)
		}
		if items[i].EndTime != nil {
			ap.EndTime = times.ISOTime(*items[i].EndTime)
		}

		pager.AddRecords(&ap)
	}
	return pager.BuildPage(), nil
}

// GetApprovedByMeReq 获取我审批完成的请求参数
type GetApprovedByMeReq struct {
	// 申请单类型
	Type string `json:"type"`
	// 申请单状态
	Status string `json:"status"`
	// 分页页号
	Page int64 `json:"page"`
	// 分页大小。默认值:10。阈值: 100。
	PageSize int64 `json:"page_size"`

	// 当前用户ID
	CurrentUserID string `json:"-"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *GetApprovedByMeReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.Type:     "type",
		&reqData.Status:   "status",
		&reqData.Page:     "page",
		&reqData.PageSize: "page_size",
	}
}

// Validate 结构体数据校验
func (reqData *GetApprovedByMeReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	user, _ := middleware.LoginUserFromContext(req.Context())
	if user != nil {
		reqData.CurrentUserID = user.ID
	}

	if !checkApprovalType(reqData.Type) {
		errs.Add([]string{"type"}, binding.RequiredError, "审批单类型值错误")
		return errs
	}

	if !checkApprovalStatus(reqData.Status) {
		errs.Add([]string{"status"}, binding.RequiredError, "审批单状态值错误")
		return errs
	}

	return errs
}

// GetApprovedByMePage 获取被我审批的申请单分页
func GetApprovedByMePage(log logger.Logger, repo model.Repo, reqData *GetApprovedByMeReq) (*page.Page, error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 100 {
		reqData.PageSize = 20
	}
	if reqData.Page < 0 {
		reqData.Page = 0
	}

	cond := model.Approval{
		Type:   reqData.Type,
		Status: reqData.Status,
	}

	totalRecords, err := repo.CountApprovedApprovals(reqData.CurrentUserID, &cond)
	if err != nil {
		return nil, err
	}

	pager := page.NewPager(reflect.TypeOf(&ApprovalPage{}), reqData.Page, reqData.PageSize, totalRecords)
	items, err := repo.GetApprovedApprovals(reqData.CurrentUserID, &cond, model.OneOrderBy("id", model.DESC), pager.BuildLimiter())
	if err != nil {
		return nil, err
	}

	for i := range items {
		ap := ApprovalPage{
			ID:         items[i].ID,
			Title:      items[i].Title,
			Type:       items[i].Type,
			IsRejected: items[i].IsRejected,
			Status:     items[i].Status,
		}
		if items[i].StartTime != nil {
			ap.StartTime = times.ISOTime(*items[i].StartTime)
		}
		if items[i].EndTime != nil {
			ap.EndTime = times.ISOTime(*items[i].EndTime)
		}
		pager.AddRecords(&ap)
	}
	return pager.BuildPage(), nil
}

// RevokeApprovalReq 取消申请单请求
type RevokeApprovalReq struct {
	// 审批单ID
	ID int `json:"id"`

	// 当前用户ID
	CurrentUserID string `json:"-"`
}

// Validate 结构体数据校验
func (reqData *RevokeApprovalReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())
	user, _ := middleware.LoginUserFromContext(req.Context())

	if user == nil {
		errs.Add([]string{"user"}, binding.RequiredError, "不存在登录用户")
		return errs
	}

	if reqData.ID <= 0 {
		errs.Add([]string{"approval"}, binding.RequiredError, "指定的申请单不存在")
		return errs
	}

	// 查询当前用户的申请单是否存在
	item, err := repo.GetApprovalByID(uint(reqData.ID))
	if gorm.IsRecordNotFoundError(err) {
		errs.Add([]string{"approval"}, binding.RequiredError, "当前用户下指定的申请单不存在")
		return errs
	}
	if err != nil {
		errs.Add([]string{"approval"}, binding.SystemError, "系统内部错误")
		return errs
	}
	if item.Initiator != user.ID {
		errs.Add([]string{"approval"}, binding.RequiredError, "当前用户下指定的申请单不存在")
		return errs
	}
	if item.Status != model.ApprovalStatusApproval {
		errs.Add([]string{"approval"}, binding.RequiredError, "申请单已经被拒绝或者已经完成,不能撤销")
		return errs
	}

	return errs
}

// RevokeApproval 取消我发起的申请单
func RevokeApproval(log logger.Logger, repo model.Repo, reqData *RevokeApprovalReq) error {
	// 查询当前用户的申请单是否存在
	item, err := repo.GetApprovalByID(uint(reqData.ID))
	if gorm.IsRecordNotFoundError(err) {
		return nil
	}
	if err != nil {
		log.Errorf("RevokeApproval: %s\n", err.Error())
		return err
	}
	// 如果为设备搬迁，将预占用机位释放
	if item.Type == model.ApprovalTypeDeviceMigration {
		var data []*SubmitDeviceMigrationApprovalReqData
		if err = json.Unmarshal([]byte(item.Metadata), &data); err != nil {
			log.Errorf("unmarshal metadata err:%s", err.Error())
			return err
		}
		var dstUsiteIDs []uint
		for _, dev := range data {
			//检查目标机位是否存在且状态为预占用
			usite, err := repo.GetServerUSiteByID(dev.DstUSiteID)
			if gorm.IsRecordNotFoundError(err) {
				log.Errorf("serverusite not found err:%s", err.Error())
				return err
			}
			if err != nil {
				log.Errorf("check serverusite err:%s", err.Error())
				return err
			}
			if usite.Status == model.USiteStatPreOccupied {
				dstUsiteIDs = append(dstUsiteIDs, dev.DstUSiteID)
			}			
		}
		
		// 仅将状态为预占用的目标机位释放为空闲
		if _, err = repo.BatchUpdateServerUSitesStatus(dstUsiteIDs, model.USiteStatFree); err != nil {
			log.Error("free usites(ids=%v) fail", dstUsiteIDs)
			return err
		}
	}
	return repo.RevokeApproval(uint(reqData.ID))
}

// GetApprovalByIDReq 查询申请单请求
type GetApprovalByIDReq struct {
	// 审批单ID
	ID int `json:"id"`
	// 当前用户ID
	CurrentUserID string `json:"-"`

	// 从UAM上获取用户信息的钩子
	GetNameFromUAM myuser.GetNameFromUAM `json:"-"`
}

// Validate 结构体数据校验
func (reqData *GetApprovalByIDReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())
	user, _ := middleware.LoginUserFromContext(req.Context())

	if user == nil {
		errs.Add([]string{"user"}, binding.RequiredError, "不存在登录用户")
		return errs
	}

	if reqData.ID <= 0 {
		errs.Add([]string{"approval"}, binding.RequiredError, "指定的申请单不存在")
		return errs
	}

	// 查询当前用户的申请单是否存在
	item, err := repo.GetApprovalByID(uint(reqData.ID))
	if gorm.IsRecordNotFoundError(err) {
		errs.Add([]string{"approval"}, binding.RequiredError, "当前用户下指定的申请单不存在")
		return errs
	}
	if err != nil {
		errs.Add([]string{"approval"}, binding.SystemError, "系统内部错误")
		return errs
	}
	if item.Initiator != user.ID && !strings.Contains(item.Approvers, user.ID) && !strings.Contains(item.Cc, user.ID) {
		errs.Add([]string{"approval"}, binding.RequiredError, "当前用户下没有权限查看指定的申请单详情")
		return errs
	}

	return errs

}

// UserInfo 用户信息
type UserInfo struct {
	ID        string `json:"id"`
	LoginName string `json:"login_name"`
	Name      string `json:"name"`
}

// ApprovalStepGetByIDRes 审批步骤
type ApprovalStepGetByIDRes struct {
	ID           uint          `json:"id"`            // 审批步骤ID
	Approver     UserInfo      `json:"approver"`      // 审批步骤审批人ID
	NextApprover UserInfo      `json:"next_approver"` // 下一个审批人
	Title        string        `json:"title"`         // 审批步骤标题
	Action       string        `json:"action"`        // 审批动作
	Remark       string        `json:"remark"`        // 审批批注
	StartTime    times.ISOTime `json:"start_time"`    // 审批步骤开始时间
	EndTime      times.ISOTime `json:"end_time"`      // 审批步骤结束时间
	Hooks        string        `json:"hooks"`         // 当前审批步骤同意后执行的钩子对象数组字符串
}

// ApprovalGetByIDRes 审批单详情
type ApprovalGetByIDRes struct {
	ID         uint                     `json:"id"`          // 审批单ID
	Title      string                   `json:"title"`       // 审批单标题
	FrontData  string                   `json:"front_data"`  // 审批单元数据快照
	Remark     string                   `json:"remark"`      //备注信息
	Type       string                   `json:"type"`        // 审批类型
	Metadata   string                   `json:"metadata"`    // 审批单元数据
	Initiator  UserInfo                 `json:"initiator"`   // 审批发起人
	Approvers  []UserInfo               `json:"approvers"`   // 审批人ID构成的JSON数组字符串
	Cc         []UserInfo               `json:"cc"`          // 抄送人ID构成的JSON数组字符串
	StartTime  times.ISOTime            `json:"start_time"`  // 审批开始时间
	EndTime    times.ISOTime            `json:"end_time"`    // 审批结束时间
	IsRejected *string                  `json:"is_rejected"` // 审批单是否被拒绝
	Status     string                   `json:"status"`      // 审批单状态
	Steps      []ApprovalStepGetByIDRes `json:"steps"`       // 审批流程
}

// GetApprovalByID 获取指定ID的申请单详情
func GetApprovalByID(log logger.Logger, repo model.Repo, reqData *GetApprovalByIDReq) (*ApprovalGetByIDRes, error) {

	iap, err := repo.GetApprovalByID(uint(reqData.ID))
	if gorm.IsRecordNotFoundError(err) && iap == nil {
		return nil, nil
	}
	if err != nil || iap == nil {
		return nil, err
	}
	agbir := &ApprovalGetByIDRes{
		ID:         iap.ID,
		Title:      iap.Title,
		Type:       iap.Type,
		Metadata:   iap.Metadata,
		IsRejected: &iap.IsRejected,
		Status:     iap.Status,
		FrontData:  iap.FrontData,
		Remark:     iap.Remark,
	}
	if iap.StartTime != nil {
		agbir.StartTime = times.ISOTime(*iap.StartTime)
	}
	if iap.EndTime != nil {
		agbir.EndTime = times.ISOTime(*iap.EndTime)
	}

	// 申请单发起人用户信息
	loginName, name, _ := reqData.GetNameFromUAM(iap.Initiator)
	agbir.Initiator = UserInfo{
		ID:        iap.Initiator,
		LoginName: loginName,
		Name:      name,
	}

	// 申请单审批人
	var aprs []string
	err = json.Unmarshal([]byte(iap.Approvers), &aprs)
	if err != nil {
		log.Debug(err)
		return agbir, nil
	}
	for k := range aprs {
		loginName, name, _ := reqData.GetNameFromUAM(aprs[k])
		agbir.Approvers = append(agbir.Approvers, UserInfo{
			ID:        aprs[k],
			LoginName: loginName,
			Name:      name,
		})
	}

	// 申请单抄送人
	var ccrs []string
	err = json.Unmarshal([]byte(iap.Cc), &ccrs)
	if err != nil {
		log.Debug(err)
		return agbir, nil
	}
	for k := range ccrs {
		loginName, name, _ := reqData.GetNameFromUAM(ccrs[k])
		agbir.Cc = append(agbir.Cc, UserInfo{
			ID:        ccrs[k],
			LoginName: loginName,
			Name:      name,
		})
	}

	ists, err := repo.GetApprovalStepByApprovalID(iap.ID)
	if err != nil {
		return nil, err
	}
	for k := range ists {
		// 当前审批人信息
		loginName, name, _ := reqData.GetNameFromUAM(ists[k].Approver)
		// 下一个审批人信息
		next := false
		nexter := ""
		for kk := range aprs {
			if next {
				nexter = aprs[kk]
			}
			if aprs[kk] == ists[k].Approver {
				next = true
			}
		}
		nloginName, nname, _ := reqData.GetNameFromUAM(nexter)
		//
		asgbir := ApprovalStepGetByIDRes{
			ID: ists[k].ID,
			Approver: UserInfo{
				ID:        ists[k].Approver,
				LoginName: loginName,
				Name:      name,
			},
			Remark: ists[k].Remark,
			Title:  ists[k].Title,
			NextApprover: UserInfo{
				ID:        nexter,
				LoginName: nloginName,
				Name:      nname,
			},
			Hooks: ists[k].Hooks,
		}
		if ists[k].StartTime != nil {
			asgbir.StartTime = times.ISOTime(*ists[k].StartTime)
		}
		if ists[k].EndTime != nil {
			asgbir.EndTime = times.ISOTime(*ists[k].EndTime)
		}
		if ists[k].Action != nil {
			asgbir.Action = *ists[k].Action
		}

		agbir.Steps = append(agbir.Steps, asgbir)
	}
	return agbir, nil
}

//ApproveReq 审批的请求参数
type ApproveReq struct {
	Action      string             `json:"action"` //reject|agree
	Remark      string             `json:"remark"`
	CurrentUser *model.CurrentUser `json:"-"`
	ApprovalID  uint               `json:"approval_id"`
	StepID      uint               `json:"approval_step_id"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *ApproveReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.Action:     "action",
		&reqData.Remark:     "remark",
		&reqData.ApprovalID: "approval_id",
		&reqData.StepID:     "approval_step_id",
	}
}

// Validate 结构体数据校验
func (reqData *ApproveReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())

	approval, err := repo.GetApprovalByID(reqData.ApprovalID)
	if err == gorm.ErrRecordNotFound {
		errs.Add([]string{"approval"}, binding.RequiredError, "申请单不存在")
		return errs
	}
	if approval.Status != model.ApprovalStatusApproval {
		errs.Add([]string{"approval"}, binding.BusinessError, "申请单已完成/撤销")
		return errs
	}
	if reqData.Action == "" {
		errs.Add([]string{"approval"}, binding.RequiredError, "审批动作值为空")
		return errs
	}
	if reqData.Action != model.ApprovalActionAgree && reqData.Action != model.ApprovalActionReject {
		errs.Add([]string{"approval"}, binding.BusinessError, "审批动作值非法")
		return errs
	}
	step, err := repo.GetApprovalStepByID(reqData.StepID)
	if err == gorm.ErrRecordNotFound {
		errs.Add([]string{"approval_step"}, binding.RequiredError, "审批步骤不存在")
		return errs
	}
	if step.Approver != reqData.CurrentUser.ID {
		errs.Add([]string{"approval_step"}, binding.BusinessError,
			fmt.Sprintf("当前用户(%s)不是该步骤审批人(%s)", reqData.CurrentUser.LoginName, step.Approver))
		return errs
	}
	return nil
}

// Approve 审批
func Approve(log logger.Logger, repo model.Repo, conf *config.Config, lim limiter.Limiter, reqData *ApproveReq) error {
	now := time.Now()
	mod := model.ApprovalStep{
		ID:         reqData.StepID,
		Action:     &reqData.Action,
		ApprovalID: reqData.ApprovalID,
		Remark:     reqData.Remark,
		EndTime:    &now,
	}
	_, err := repo.UpdateApprovalStep(&mod)
	if err != nil {
		log.Error("update approval step fail", err.Error())
		return err
	}
	//若当前审批步骤结果为agree
	//将下一步的start_time设置为当前时间。若当前已经是最后一个审批步骤，则结束整个审批单。
	//调用hooks
	if reqData.Action == model.ApprovalActionAgree {
		//获取下一个步骤
		steps, err := repo.GetApprovalStepByApprovalID(reqData.ApprovalID)
		if err != nil {
			log.Errorf("get steps by approval_id:%d fail:%s", reqData.ApprovalID, err.Error())
			return err
		}
		for i, step := range steps {
			if step.ID == reqData.StepID {
				//如果已经是最后一个
				if i == len(steps)-1 {
					approval := model.Approval{
						ID:         reqData.ApprovalID,
						EndTime:    &now,
						IsRejected: model.NO,
						Status:     model.ApprovalStatusCompleted,
					}
					_, err := repo.UpdateApproval(&approval)
					if err != nil {
						log.Errorf("update by approval_id:%d fail:%s", reqData.ApprovalID, err.Error())
						return err
					}
				} else {
					//下一个步骤
					mod := model.ApprovalStep{
						ID:        steps[i+1].ID,
						StartTime: &now,
					}
					_, err := repo.UpdateApprovalStep(&mod)
					if err != nil {
						log.Errorf("update next approval step id:%d fail,%s", steps[i+1].ID, err.Error())
						return err
					}
				}
				//hooks
				if step.Hooks != "" {
					//获取审批单的传入参数
					approval, err := repo.GetApprovalByID(step.ApprovalID)
					if err != nil {
						log.Error("get approval metadata fail", err.Error())
						return err
					}
					//执行调用钩子方法
					hooks := make([]*model.StepHook, 0)
					err = json.Unmarshal([]byte(step.Hooks), &hooks)
					if err != nil {
						log.Errorf("unmarshal hooks err:%s", err.Error())
						return err
					}
					for _, hook := range hooks {
						isRecycleOSReinstall := false
						switch hook.ID {
						case model.HookDevPowerOff:
							var sns []string
							if err = json.Unmarshal([]byte(approval.Metadata), &sns); err != nil {
								log.Errorf("unmarshal metadata err:%s", err.Error())
								return err
							}
							if _, err = BatchOperateOOBPower(log, repo, PowerOff, conf, false, sns); err != nil {
								hook.Results = append(hook.Results, &model.StepHookResult{
									TargetID:   fmt.Sprintf("%v", sns),
									TargetType: "power off device",
									ExecResult: "failure",
								})
								log.Error("power off device (sns:%v) fail", sns)
							}
							// 变更记录
							for i := range sns {
								devLog := model.ChangeLog {
										OperationUser:		reqData.CurrentUser.Name,
										OperationType:		model.OperationTypePowerControl,
										OperationDetail:	fmt.Sprintf("物理机关电（SN: %v）", sns[i]),
										OperationTime:		times.ISOTime(time.Now()).ToTimeStr(),
								}
								adll := &AppendDeviceLifecycleLogReq{
									SN:					sns[i],
									LifecycleLog: 		devLog,
								}
								if err = AppendDeviceLifecycleLogBySN(log, repo, adll);err != nil {
									log.Error("LifecycleLog: append device lifecycle log (sn:%v) fail(%s)", sns[i], err.Error())
								}
							}
						case model.HookDevRestart:
							var sns []string
							if err = json.Unmarshal([]byte(approval.Metadata), &sns); err != nil {
								log.Errorf("unmarshal metadata err:%s", err.Error())
								return err
							}
							if _, err = BatchOperateOOBPower(log, repo, PowerRestart, conf, false, sns); err != nil {
								hook.Results = append(hook.Results, &model.StepHookResult{
									TargetID:   fmt.Sprintf("%v", sns),
									TargetType: "power restart device",
									ExecResult: "failure",
								})
								log.Error("power restart device (sns:%v) fail", sns)
								//return ?
							}
							// 变更记录
							for i := range sns {
								devLog := model.ChangeLog {
									OperationUser:		reqData.CurrentUser.Name,
									OperationType:		model.OperationTypePowerControl,
									OperationDetail:	fmt.Sprintf("物理机重启（SN: %v）", sns[i]),
									OperationTime:		times.ISOTime(time.Now()).ToTimeStr(),
								}
								adll := &AppendDeviceLifecycleLogReq{
									SN:					sns[i],
									LifecycleLog: 		devLog,
								}
								if err = AppendDeviceLifecycleLogBySN(log, repo, adll);err != nil {
									log.Error("LifecycleLog: append device lifecycle log (sn:%v) fail(%s)", sns[i], err.Error())
								}
							}			
						case model.HookIPUnassign:
							var ids []uint
							if err = json.Unmarshal([]byte(approval.Metadata), &ids); err != nil {
								log.Errorf("unmarshal metadata err:%s", err.Error())
								return err
							}
							for _, id := range ids {
								ip, _ := repo.GetIPByID(id)
								if ip.SN != "" {
									ipn, err := repo.GetIPNetworkByID(ip.IPNetworkID)
									if err != nil && err != gorm.ErrRecordNotFound {
										log.Errorf("get ipnetwork of ip:%s fail,%v", ip.IP, err)
										return err
									}
									ds, err := repo.GetDeviceSettingBySN(ip.SN)
									if err != nil && err != gorm.ErrRecordNotFound {
										log.Errorf("get device setting of sn:%s fail,%v", ip.SN, err)
									} else if ds != nil {
										ds.SN = ip.SN
										if ipn.Version == model.IPv4 {
											if ip.Scope != nil && *ip.Scope == model.IPScopeIntranet {
												ds.IntranetIP = removeIP(ds.IntranetIP, ip.IP)
											} 
											if ip.Scope != nil && *ip.Scope == model.IPScopeExtranet {
												ds.ExtranetIP = removeIP(ds.ExtranetIP, ip.IP)
											}
										} else if ipn.Version == model.IPv6 {
											if ip.Scope != nil && *ip.Scope == model.IPScopeIntranet {
												ds.IntranetIPv6 = removeIP(ds.IntranetIPv6, ip.IP)
											} 
											if ip.Scope != nil && *ip.Scope == model.IPScopeExtranet {
												ds.ExtranetIPv6 = removeIP(ds.ExtranetIPv6, ip.IP)
											}
										}
										err = repo.SaveDeviceSetting(ds)
										if err != nil {
											log.Error(err)
											return err
										}
										// 变更记录
										devLog := model.ChangeLog {
											OperationUser:		reqData.CurrentUser.Name,
											OperationType:		model.OperationTypeUpdate,
											OperationDetail:	fmt.Sprintf("物理机IP回收（SN: %v, IP: %v）", ip.SN, ip.IP),
											OperationTime:		times.ISOTime(time.Now()).ToTimeStr(),
										}
										adll := &AppendDeviceLifecycleLogReq{
											SN:					ip.SN,
											LifecycleLog: 		devLog,
										}
										if err = AppendDeviceLifecycleLogBySN(log, repo, adll);err != nil {
											log.Error("LifecycleLog: append device lifecycle log (sn:%v) fail(%s)", ip.SN, err.Error())
										}
									}
								}
								if err = repo.UnassignIP(id); err != nil {
									hook.Results = append(hook.Results, &model.StepHookResult{
										TargetID:   fmt.Sprintf("%v", ids),
										TargetType: "unssign ip",
										ExecResult: "failure",
									})
									log.Error("unssign ip (id:%v) fail", ids)
									//return ?
								}
							}
						case model.HookNetAreaOffline:
							var ids []uint
							if err = json.Unmarshal([]byte(approval.Metadata), &ids); err != nil {
								log.Errorf("unmarshal metadata err:%s", err.Error())
								return err
							}
							if _, err = repo.UpdateNetworkAreaStatus(model.NetworkAreaStatOffline, ids...); err != nil {
								hook.Results = append(hook.Results, &model.StepHookResult{
									TargetID:   fmt.Sprintf("%v", ids),
									TargetType: "network-area",
									ExecResult: "failure",
								})
								log.Error("offline network-area (id:%v) fail", ids)
								//return ?
							}
						case model.HookServerRoomAbolish:
							var ids []uint
							if err = json.Unmarshal([]byte(approval.Metadata), &ids); err != nil {
								log.Errorf("unmarshal metadata err:%s", err.Error())
								return err
							}
							if _, err = repo.UpdateServerRoomStatus(model.RoomStatAbolished, ids...); err != nil {
								hook.Results = append(hook.Results, &model.StepHookResult{
									TargetID:   fmt.Sprintf("%v", ids),
									TargetType: "server-room",
									ExecResult: "failure",
								})
								log.Error("abolish server-room (id:%v) fail", ids)
								//return ?
							}
						case model.HookIDCAbolish:
							var ids []uint
							if err = json.Unmarshal([]byte(approval.Metadata), &ids); err != nil {
								log.Errorf("unmarshal metadata err:%s", err.Error())
								return err
							}
							if _, err = repo.UpdateIDCStatus(model.IDCStatAbolished, ids...); err != nil {
								hook.Results = append(hook.Results, &model.StepHookResult{
									TargetID:   fmt.Sprintf("%v", ids),
									TargetType: "idc",
									ExecResult: "failure",
								})
								log.Error("abolish idc (id:%v) fail", ids)
								//return ?
							}
						case model.HookCabinetPowerOff:
							var cabinetIDs []uint
							if err = json.Unmarshal([]byte(approval.Metadata), &cabinetIDs); err != nil {
								log.Errorf("unmarshal metadata err:%s", err.Error())
								return err
							}
							for _, id := range cabinetIDs {
								if _, err = repo.PowerOffServerCabinetByID(id); err != nil {
									hook.Results = append(hook.Results, &model.StepHookResult{
										TargetID:   fmt.Sprintf("%d", id),
										TargetType: "cabinet",
										ExecResult: "failure",
									})
									log.Error("power off cabinet(id=%d) fail", id)
									//return ?
								}

							}
						case model.HookCabinetOffline:
							var cabinetIDs []uint
							if err = json.Unmarshal([]byte(approval.Metadata), &cabinetIDs); err != nil {
								log.Errorf("unmarshal metadata err:%s", err.Error())
								return err
							}
							if _, err = repo.UpdateServerCabinetStatus(cabinetIDs, model.CabinetStatOffline); err != nil {
								hook.Results = append(hook.Results, &model.StepHookResult{
									TargetID:   fmt.Sprintf("%v", cabinetIDs),
									TargetType: "cabinet",
									ExecResult: "failure",
								})
								log.Error("offline cabinet(ids=%v) fail", cabinetIDs)
								//return ?
							}
						case model.HookDeviceMigrationPowerOff:
							var data []*SubmitDeviceMigrationApprovalReqData
							if err = json.Unmarshal([]byte(approval.Metadata), &data); err != nil {
								log.Errorf("unmarshal metadata err:%s", err.Error())
								return err
							}
							for _, dev := range data {
								if dev.MigType != model.MigTypeStore2Usite && dev.MigType != model.MigTypeStore2Store {
									if _, err = OperateOOBPower(log, repo, dev.SN, conf.Crypto.Key, conf.Server.OOBDomain,
										PowerOff, false); err != nil {
										hook.Results = append(hook.Results, &model.StepHookResult{
											TargetID:   fmt.Sprintf("%s", dev.SN),
											TargetType: "device",
											ExecResult: "failure",
										})
										log.Error("power off device(sn=%s) fail", dev.SN)
									}
								}
								//将运营状态更新为搬迁中
								if _, err = repo.UpdateDeviceBySN(&model.Device{
									SN:              dev.SN,
									OperationStatus: model.DevOperStatMoving,
								}); err != nil {
									hook.Results = append(hook.Results, &model.StepHookResult{
										TargetID:   fmt.Sprintf("%s", dev.SN),
										TargetType: "device",
										ExecResult: "failure",
									})
									log.Error("update device(sn=%s) status pre_deploy fail", dev.SN)
								}
							}
						case model.HookDeviceMigrationReleaseIP:
							var data []*SubmitDeviceMigrationApprovalReqData
							if err = json.Unmarshal([]byte(approval.Metadata), &data); err != nil {
								log.Errorf("unmarshal metadata err:%s", err.Error())
								return err
							}
							for _, dev := range data {
								if dev.MigType != model.MigTypeStore2Usite && dev.MigType != model.MigTypeStore2Store {
									//释放IP资源
									if _, err = repo.ReleaseIP(dev.SN, model.Intranet); err != nil {
										hook.Results = append(hook.Results, &model.StepHookResult{
											TargetID:   fmt.Sprintf("%s", dev.SN),
											TargetType: "device",
											ExecResult: "failure",
										})
										log.Error("device(sn=%s) release ip fail", dev.SN)
									}
									//外网IP不一定有
									_, _ = repo.ReleaseIP(dev.SN, model.Extranet)

									//同时删除装机记录
									if _, err := repo.DeleteDeviceSettingBySN(dev.SN); err != nil {
										log.Error("device(sn=%s) clear device setting fail", dev.SN)
									}
								}
							}
						case model.HookDeviceMigrationReserveIP:
							var data []*SubmitDeviceMigrationApprovalReqData
							if err = json.Unmarshal([]byte(approval.Metadata), &data); err != nil {
								log.Errorf("unmarshal metadata err:%s", err.Error())
								return err
							}
							//保留IP资源，保留期根据config文件配置,0则不保留
							reserveDay := conf.IP.ReserveDay							
							for _, dev := range data {
								if dev.MigType != model.MigTypeStore2Usite && dev.MigType != model.MigTypeStore2Store {
									if reserveDay <= 0 {
										//释放IP资源
									    if _, err = repo.ReleaseIP(dev.SN, model.Intranet); err != nil {
									    	hook.Results = append(hook.Results, &model.StepHookResult{
									    		TargetID:   fmt.Sprintf("%s", dev.SN),
									    		TargetType: "device",
									    		ExecResult: "failure",
									    	})
									    	log.Error("device(sn=%s) release ip fail", dev.SN)
									    }
									    //外网IP不一定有
									    _, _ = repo.ReleaseIP(dev.SN, model.Extranet)
    
									    //同时删除装机记录
									    if _, err := repo.DeleteDeviceSettingBySN(dev.SN); err != nil {
									    	log.Error("device(sn=%s) clear device setting fail", dev.SN)
									    }
									}
									if reserveDay > 0 {
										releaseDate := now.AddDate(0, 0, reserveDay)
										if _, err = repo.ReserveIP(dev.SN, model.Intranet, releaseDate); err != nil {
											hook.Results = append(hook.Results, &model.StepHookResult{
												TargetID:   fmt.Sprintf("%s", dev.SN),
												TargetType: "device",
												ExecResult: "failure",
											})
											log.Error("device(sn=%s) reserve ip fail", dev.SN)
										}
										//外网IP不一定有
										_, _ = repo.ReserveIP(dev.SN, model.Extranet, releaseDate)
	
										//同时删除装机记录
										if _, err := repo.DeleteDeviceSettingBySN(dev.SN); err != nil {
											log.Error("device(sn=%s) clear device setting fail", dev.SN)
										}
									}
								}
							}						
						case model.HookDeviceMigration:
							var data []*SubmitDeviceMigrationApprovalReqData
							if err = json.Unmarshal([]byte(approval.Metadata), &data); err != nil {
								log.Errorf("unmarshal metadata err:%s", err.Error())
								return err
							}
							var dstUsiteIDs []uint
							//存储需释放的机位
							var usiteFreeIDs []uint
							var usiteDisabledIDs []uint

							for _, dev := range data {
								//获取原机位并标记为空闲
								oriDev, err := repo.GetDeviceBySN(dev.SN)
								if err != nil {
									hook.Results = append(hook.Results, &model.StepHookResult{
										TargetID:   fmt.Sprintf("%s", dev.SN),
										TargetType: "device",
										ExecResult: "failure",
									})
									log.Error("get origin device(sn=%s) usite fail", dev.SN)
								}
								// 变更记录
								optDetail, err := convert2DetailOfOperationTypeMove(repo, dev)
								if err != nil {
									log.Errorf("Fail to convert Detail of OperationTypeMove: %v", err)
								}
								devLog := model.ChangeLog {
									OperationUser:		reqData.CurrentUser.Name,
									OperationType:		model.OperationTypeMove,
									OperationDetail:	optDetail,
									OperationTime:		times.ISOTime(time.Now()).ToTimeStr(),
								}
								adll := &AppendDeviceLifecycleLogReq{
									SN:					dev.SN,
									LifecycleLog: 		devLog,
								}
								if err = AppendDeviceLifecycleLogBySN(log, repo, adll);err != nil {
									log.Error("LifecycleLog: append device lifecycle log (sn:%v) fail(%s)", dev.SN, err.Error())
								}
								// 	搬迁类型细分							
								switch dev.MigType {
								case model.MigTypeUsite2Usite:
									dstUsiteIDs = append(dstUsiteIDs, dev.DstUSiteID)
									//获取原机位并标记为空闲 或 不可用
									if oriDev.USiteID != nil {
										if cabinet, err := repo.GetServerCabinetByID(oriDev.CabinetID); err == nil {
											// 若上层机架为[已锁定]，则对应机位不可释放
											if cabinet.Status == model.CabinetStatLocked {
												usiteDisabledIDs = append(usiteDisabledIDs, *oriDev.USiteID)
											} else {
												usiteFreeIDs = append(usiteFreeIDs, *oriDev.USiteID)
											}
										}				
									}									

									//将物理机状态更新为[待部署]
									if _, err = repo.UpdateDeviceBySN(&model.Device{
										SN:              dev.SN,
										IDCID:           dev.DstIDCID,
										ServerRoomID:    dev.DstServerRoomID,
										CabinetID:       dev.DstCabinetID,
										USiteID:         &dev.DstUSiteID,
										OperationStatus: model.DevOperStatPreDeploy,
									}); err != nil {
										hook.Results = append(hook.Results, &model.StepHookResult{
											TargetID:   fmt.Sprintf("%s", dev.SN),
											TargetType: "device",
											ExecResult: "failure",
										})
										log.Error("update device(sn=%s) status pre_deploy fail", dev.SN)
									}
								case model.MigTypeUsite2Store:
									//获取原机位并标记为空闲 或 不可用
									if oriDev.USiteID != nil {
										if cabinet, err := repo.GetServerCabinetByID(oriDev.CabinetID); err == nil {
											// 若上层机架为[已锁定]，则对应机位不可释放
											if cabinet.Status == model.CabinetStatLocked {
												usiteDisabledIDs = append(usiteDisabledIDs, *oriDev.USiteID)
											} else {
												usiteFreeIDs = append(usiteFreeIDs, *oriDev.USiteID)
											}
										}				
									}

									//将物理机状态更新为[待部署]
									oriDev.IDCID = dev.DstIDCID
									oriDev.ServerRoomID = 0
									oriDev.CabinetID = 0
									oriDev.USiteID = nil //&zero
									oriDev.StoreRoomID = dev.DstStoreRoomID
									oriDev.VCabinetID = dev.DstVCabinetID
									oriDev.OperationStatus = model.DevOperStatInStore
									if _, err = repo.SaveDevice(oriDev); err != nil {
										hook.Results = append(hook.Results, &model.StepHookResult{
											TargetID:   fmt.Sprintf("%s", dev.SN),
											TargetType: "device",
											ExecResult: "failure",
										})
										log.Error("update device(sn=%s) status pre_deploy fail", dev.SN)
									}
								case model.MigTypeStore2Usite:
									//占用目标机位
									dstUsiteIDs = append(dstUsiteIDs, dev.DstUSiteID)

									//将物理机状态更新为[待部署]
									oriDev.IDCID = dev.DstIDCID
									oriDev.ServerRoomID = dev.DstServerRoomID
									oriDev.CabinetID = dev.DstCabinetID
									oriDev.USiteID = &dev.DstUSiteID
									oriDev.StoreRoomID = 0
									oriDev.VCabinetID = 0
									now := time.Now()
									oriDev.StartedAt = now
									oriDev.OnShelveAt = now
									oriDev.OperationStatus = model.DevOperStatPreDeploy
									if _, err = repo.SaveDevice(oriDev); err != nil {
										hook.Results = append(hook.Results, &model.StepHookResult{
											TargetID:   fmt.Sprintf("%s", dev.SN),
											TargetType: "device",
											ExecResult: "failure",
										})
										log.Error("update device(sn=%s) status pre_deploy fail", dev.SN)
									}
								case model.MigTypeStore2Store:
									//只修改这仨
									oriDev.IDCID = dev.DstIDCID
									oriDev.StoreRoomID = dev.DstStoreRoomID
									oriDev.VCabinetID = dev.DstVCabinetID
									if _, err = repo.UpdateDeviceBySN(oriDev); err != nil {
										hook.Results = append(hook.Results, &model.StepHookResult{
											TargetID:   fmt.Sprintf("%s", dev.SN),
											TargetType: "device",
											ExecResult: "failure",
										})
										log.Error("update device(sn=%s) status pre_deploy fail", dev.SN)
									}
								}
							}
							//将目标机位标记为[已占用]
							if _, err = repo.BatchUpdateServerUSitesStatus(dstUsiteIDs, model.USiteStatUsed); err != nil {
								hook.Results = append(hook.Results, &model.StepHookResult{
									TargetID:   fmt.Sprintf("%v", dstUsiteIDs),
									TargetType: "device",
									ExecResult: "failure",
								})
								log.Error("occupy usites(ids=%v) fail", dstUsiteIDs)
							}
							//将原机位标记为空闲
							if _, err = repo.BatchUpdateServerUSitesStatus(usiteFreeIDs, model.USiteStatFree); err != nil {
								hook.Results = append(hook.Results, &model.StepHookResult{
									TargetID:   fmt.Sprintf("%v", usiteFreeIDs),
									TargetType: "device",
									ExecResult: "failure",
								})
								log.Error("free usites(ids=%v) fail", usiteFreeIDs)
							}
							//将原机位标记为不可用
							if _, err = repo.BatchUpdateServerUSitesStatus(usiteDisabledIDs, model.USiteStatDisabled); err != nil {
								hook.Results = append(hook.Results, &model.StepHookResult{
									TargetID:   fmt.Sprintf("%v", usiteDisabledIDs),
									TargetType: "device",
									ExecResult: "failure",
								})
								log.Error("disable usites(ids=%v) fail", usiteDisabledIDs)
							}							

						case model.HookDeviceRetirementPowerOff:
							var SNs []string
							if err = json.Unmarshal([]byte(approval.Metadata), &SNs); err != nil {
								log.Errorf("unmarshal metadata err:%s", err.Error())
								return err
							}
							for _, sn := range SNs {
								if _, err = OperateOOBPower(log, repo, sn, conf.Crypto.Key, conf.Server.OOBDomain,
									PowerOff, false); err != nil {
									hook.Results = append(hook.Results, &model.StepHookResult{
										TargetID:   fmt.Sprintf("%s", sn),
										TargetType: "device",
										ExecResult: "failure",
									})
									log.Error("power off device(sn=%s) fail", sn)
								}
							}
						case model.HookDeviceRetirementReleaseIP:
							var SNs []string
							if err = json.Unmarshal([]byte(approval.Metadata), &SNs); err != nil {
								log.Errorf("unmarshal metadata err:%s", err.Error())
								return err
							}
							for _, sn := range SNs {
								//释放IP资源
								if _, err = repo.ReleaseIP(sn, model.Intranet); err != nil {
									hook.Results = append(hook.Results, &model.StepHookResult{
										TargetID:   fmt.Sprintf("%s", sn),
										TargetType: "device",
										ExecResult: "failure",
									})
									log.Error("device(sn=%s) release ip fail", sn)
								}
								//外网IP不一定有
								_, _ = repo.ReleaseIP(sn, model.Extranet)
								_, _ = repo.ReleaseIPv6(sn, model.IPScopeIntranet)
								_, _ = repo.ReleaseIPv6(sn, model.IPScopeExtranet)
								//同时删除装机记录
								if _, err := repo.DeleteDeviceSettingBySN(sn); err != nil {
									log.Error("device(sn=%s) clear device setting fail", sn)
								}								
							}
						case model.HookDeviceRetirementReserveIP:
							var SNs []string
							if err = json.Unmarshal([]byte(approval.Metadata), &SNs); err != nil {
								log.Errorf("unmarshal metadata err:%s", err.Error())
								return err
							}
							//保留IP资源，保留期根据config文件配置,0则不保留
							reserveDay := conf.IP.ReserveDay							
							for _, sn := range SNs {
								if reserveDay <= 0{
									if _, err = repo.ReleaseIP(sn, model.Intranet); err != nil {
										hook.Results = append(hook.Results, &model.StepHookResult{
											TargetID:   fmt.Sprintf("%s", sn),
											TargetType: "device",
											ExecResult: "failure",
										})
										log.Error("device(sn=%s) release ip fail", sn)
									}
									//外网IP不一定有
									_, _ = repo.ReleaseIP(sn, model.Extranet)
									_, _ = repo.ReleaseIPv6(sn, model.IPScopeIntranet)
									_, _ = repo.ReleaseIPv6(sn, model.IPScopeExtranet)
									//同时删除装机记录
									if _, err := repo.DeleteDeviceSettingBySN(sn); err != nil {
										log.Error("device(sn=%s) clear device setting fail", sn)
									}																	
								}
								if reserveDay > 0 {
									releaseDate := now.AddDate(0, 0, reserveDay)
									if _, err = repo.ReserveIP(sn, model.Intranet, releaseDate); err != nil {
										hook.Results = append(hook.Results, &model.StepHookResult{
											TargetID:   fmt.Sprintf("%s", sn),
											TargetType: "device",
											ExecResult: "failure",
										})
										log.Error("device(sn=%s) release ip fail", sn)
									}
									//外网IP不一定有
									_, _ = repo.ReserveIP(sn, model.Extranet, releaseDate)
									//同时删除装机记录
									if _, err := repo.DeleteDeviceSettingBySN(sn); err != nil {
										log.Error("device(sn=%s) clear device setting fail", sn)
									}		
								}				
							}							
						case model.HookDeviceRetirement:
							//修改运营状态，同时释放机位资源
							var SNs []string
							if err = json.Unmarshal([]byte(approval.Metadata), &SNs); err != nil {
								log.Errorf("unmarshal metadata err:%s", err.Error())
								return err
							}
							//var usiteIDs []uint
							//存储需释放的机位
							var usiteFreeIDs []uint
							var usiteDisabledIDs []uint
							for _, sn := range SNs {
								dev, err := repo.GetDeviceBySN(sn)
								if err != nil {
									log.Errorf("get device by sn:%s fail, %s", sn, err.Error())
									continue
								}
								if dev != nil && dev.USiteID != nil {
									//usiteIDs = append(usiteIDs, *dev.USiteID)
									if cabinet, err := repo.GetServerCabinetByID(dev.CabinetID); err == nil {
										// 若上层机架为[已锁定]，则对应机位不可释放
										if cabinet.Status == model.CabinetStatLocked {
											usiteDisabledIDs = append(usiteDisabledIDs, *dev.USiteID)
										} else {
											usiteFreeIDs = append(usiteFreeIDs, *dev.USiteID)
										}
									}
								}
								// 变更记录
								optDetail, err := convert2DetailOfOperationTypeRetire(repo, dev)
								if err != nil {
									log.Errorf("Fail to convert Detail of OperationTypeRetire: %v", err)
								}
								devLog := model.ChangeLog {
									OperationUser:		reqData.CurrentUser.Name,
									OperationType:		model.OperationTypeRetire,
									OperationDetail:	optDetail,
									OperationTime:		times.ISOTime(time.Now()).ToTimeStr(),
								}
								adll := &AppendDeviceLifecycleLogReq{
									SN:					dev.SN,
									LifecycleLog: 		devLog,
								}
								if err = AppendDeviceLifecycleLogBySN(log, repo, adll);err != nil {
									log.Error("LifecycleLog: append device lifecycle log (sn:%v) fail(%s)", dev.SN, err.Error())
								}								
								// 释放机位并修改状态
								dev.USiteID = nil
								dev.PowerStatus = model.PowerStatusOff
								dev.OperationStatus = model.DevOperStateRetired
								if _, err := repo.SaveDevice(dev); err != nil {
									log.Error(err)
									//usiteIDs = append(usiteIDs, *dev.USiteID)
								}
								// 记录设备退役日期
								saveDevLifecycleReq := &SaveDeviceLifecycleReq {
									DeviceLifecycleBase: DeviceLifecycleBase{
										DeviceRetiredDate:		time.Now(),
									},
								}
								// DeviceLifecycle 查询是否已经存在
								devLifecycle, err := repo.GetDeviceLifecycleBySN(dev.SN)
								if err != nil && err != gorm.ErrRecordNotFound {
									log.Error(err)
								} 
								if devLifecycle != nil {
									saveDevLifecycleReq.ID = devLifecycle.ID
									// 保存或更新 DeviceLifecycle
									if err = SaveDeviceLifecycle(log, repo, saveDevLifecycleReq); err != nil {
										log.Error(err)
									}
								}
							}
							if _, err = repo.BatchUpdateServerUSitesStatus(usiteFreeIDs, model.USiteStatFree); err != nil {
								hook.Results = append(hook.Results, &model.StepHookResult{
									TargetID:   fmt.Sprintf("%v", usiteFreeIDs),
									TargetType: "device",
									ExecResult: "failure",
								})
								log.Error("free usites(ids=%v) fail", usiteFreeIDs)
							}
							if _, err = repo.BatchUpdateServerUSitesStatus(usiteDisabledIDs, model.USiteStatDisabled); err != nil {
								hook.Results = append(hook.Results, &model.StepHookResult{
									TargetID:   fmt.Sprintf("%v", usiteDisabledIDs),
									TargetType: "device",
									ExecResult: "failure",
								})
								log.Error("free usites(ids=%v) fail", usiteDisabledIDs)
							}							
						case model.HookDeviceRecyclePreRetire:
							var SNs []string
							if err = json.Unmarshal([]byte(approval.Metadata), &SNs); err != nil {
								log.Errorf("unmarshal metadata err:%s", err.Error())
								return err
							}
							for _, sn := range SNs {
								// 变更记录
								optDetail, err := convert2DetailOfOperationTypeUpdate(repo, model.Device{
									SN: sn, OperationStatus: model.DevOperStatPreRetire,
								})
								if err != nil {
									log.Errorf("Fail to convert Detail of OperationTypeUpdate: %v", err)
								}
								devLog := model.ChangeLog {
									OperationUser:		reqData.CurrentUser.Name,
									OperationType:		model.OperationTypeUpdate,
									OperationDetail:	optDetail,
									OperationTime:		times.ISOTime(time.Now()).ToTimeStr(),
								}
								adll := &AppendDeviceLifecycleLogReq{
									SN:					sn,
									LifecycleLog: 		devLog,
								}
								if err = AppendDeviceLifecycleLogBySN(log, repo, adll);err != nil {
									log.Error("LifecycleLog: append device lifecycle log (sn:%v) fail(%s)", sn, err.Error())
								}								
								if _, err := repo.UpdateDeviceBySN(&model.Device{
									SN: sn, OperationStatus: model.DevOperStatPreRetire,
								}); err != nil {
									log.Errorf("recycle device sn:%s fail, %s", sn, err.Error())
								}
							}
						case model.HookDeviceRecyclePreMove:
							var SNs []string
							if err = json.Unmarshal([]byte(approval.Metadata), &SNs); err != nil {
								log.Errorf("unmarshal metadata err:%s", err.Error())
								return err
							}
							for _, sn := range SNs {
								// 变更记录
								optDetail, err := convert2DetailOfOperationTypeUpdate(repo, model.Device{
									SN: sn, OperationStatus: model.DevOperStatPreMove,
								})
								if err != nil {
									log.Errorf("Fail to convert Detail of OperationTypeUpdate: %v", err)
								}
								devLog := model.ChangeLog {
									OperationUser:		reqData.CurrentUser.Name,
									OperationType:		model.OperationTypeUpdate,
									OperationDetail:	optDetail,
									OperationTime:		times.ISOTime(time.Now()).ToTimeStr(),
								}
								adll := &AppendDeviceLifecycleLogReq{
									SN:					sn,
									LifecycleLog: 		devLog,
								}
								if err = AppendDeviceLifecycleLogBySN(log, repo, adll);err != nil {
									log.Error("LifecycleLog: append device lifecycle log (sn:%v) fail(%s)", sn, err.Error())
								}								
								if _, err := repo.UpdateDeviceBySN(&model.Device{
									SN: sn, OperationStatus: model.DevOperStatPreMove,
								}); err != nil {
									log.Errorf("recycle device sn:%s fail, %s", sn, err.Error())
								}
							}
						case model.HookDeviceRecycleReinstall:
							isRecycleOSReinstall = true
							fallthrough
						case model.HookDeviceOSReinstallation:
							var settingsExt OSReinstallSetting = make([]*OSReinstallSettingItem, 0)
							if err = json.Unmarshal([]byte(approval.Metadata), &settingsExt); err != nil {
								log.Errorf("unmarshal metadata err:%s", err.Error())
								return err
							}
							var settings DeviceSettings = make([]*SaveDeviceSettingItem, 0, len(settingsExt))
							for _, s := range settingsExt {
								settings = append(settings, &s.SaveDeviceSettingItem)
							}
							saveReq := SaveDeviceSettingsReq{
								Settings:    settings,
								CurrentUser: reqData.CurrentUser,
							}
							succeed, err := SaveDeviceSettings(log, repo, conf, lim, &saveReq)
							if err != nil {
								hook.Results = append(hook.Results, &model.StepHookResult{
									TargetID:   fmt.Sprintf("%v", settings),
									TargetType: "device",
									ExecResult: "failure",
								})
								log.Error("save device settings(params=%v) fail", settings)
							}
							if len(succeed) != len(settings) {
								hook.Results = append(hook.Results, &model.StepHookResult{
									TargetID:   fmt.Sprintf("%v", settings),
									TargetType: "device",
									ExecResult: "failure",
								})
								log.Error("save device settings(params=%v) partly fail", settings)
							}

							// 更新设备状态

							for k := range settingsExt {
								mod := &model.Device{
									SN:              settingsExt[k].SN,
									OperationStatus: model.DevOperStatReinstalling,
								}
								// 这里实现不太好，但目前没想到比较好的方式，
								// 重装需要将运行状态改回初始的运行状态，而回收重装成功后统一成已上架
								// 所以这里用remark字段标记，初始的状态
								if isRecycleOSReinstall == false {
									mod.Remark = settingsExt[k].OriginStatus
								}
								_, err = repo.UpdateDeviceBySN(mod)
								if err != nil {
									hook.Results = append(hook.Results, &model.StepHookResult{
										TargetID:   fmt.Sprintf("%v", settingsExt[k]),
										TargetType: "device",
										ExecResult: "failure",
									})
									log.Error("update device status (params=%v)fail", settings[k])
								}
							}
						}
						if hook.ContinueOnError == false && len(hook.Results) != 0 {
							break
						}
					}
				}
				break
			}
		}

	} else if reqData.Action == model.ApprovalActionReject { //若当前审批步骤结果为reject，则结束整个审批单。
		//获取审批单的传入参数
		oriApproval, err := repo.GetApprovalByID(reqData.ApprovalID)
		if err != nil {
			log.Error("get approval metadata fail", err.Error())
			return err
		}
		//如果是搬迁的审批，需要释放预占用的机位
		if oriApproval.Type == model.ApprovalTypeDeviceMigration {
			var data []*SubmitDeviceMigrationApprovalReqData
			if err = json.Unmarshal([]byte(oriApproval.Metadata), &data); err != nil {
				log.Errorf("unmarshal metadata err:%s", err.Error())
				return err
			}
			var dstUsiteIDs []uint
			for _, dev := range data {
				dstUsiteIDs = append(dstUsiteIDs, dev.DstUSiteID)
				//将运营状态更新为待搬迁
				if _, err = repo.UpdateDeviceBySN(&model.Device{
					SN:              dev.SN,
					OperationStatus: model.DevOperStatPreMove,
				}); err != nil {
					log.Error("update device(sn=%s) status pre_deploy failed when resuming DeviceMigrationApprovalReq", dev.SN)
				}
			}
			//将目标机位(预占用)恢复为空闲
			if _, err = repo.BatchUpdateServerUSitesStatus(dstUsiteIDs, model.USiteStatFree); err != nil {
				log.Error("free usites(ids=%v) failed when resuming DeviceMigrationApprovalReq", dstUsiteIDs)
			}
		}
		approval := model.Approval{
			ID:         reqData.ApprovalID,
			EndTime:    &now,
			IsRejected: model.YES,
			Status:     model.ApprovalStatusFailure,
		}
		_, err = repo.UpdateApproval(&approval)
		if err != nil {
			log.Errorf("update by approval_id:%d fail:%s", reqData.ApprovalID, err.Error())
			return err
		}
	}
	return nil
}
