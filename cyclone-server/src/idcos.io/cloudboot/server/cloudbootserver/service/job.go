package service

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/voidint/binding"
	"github.com/voidint/page"

	"idcos.io/cloudboot/config"
	"idcos.io/cloudboot/job"
	"idcos.io/cloudboot/job/mysql"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils"
	"idcos.io/cloudboot/utils/collection"
	"idcos.io/cloudboot/utils/times"
)

// Job 任务
type Job struct {
	ID         string              `json:"id"`          // 任务ID（保证全局唯一）
	Builtin    string              `json:"builtin"`     // 是否是内建任务
	Title      string              `json:"title"`       // 任务标题
	Category   string              `json:"category"`    // 任务类型。可选值:inspection、installation_timeout、release_ip
	Rate       string              `json:"rate"`        // 任务执行频率。可选值:immediately-立刻执行; fixed_rate-固定频率(周期性)执行;
    //   - Full crontab specs, e.g. "* * * * * ?"
    //   - Descriptors, e.g. "@midnight", "@every 1h30m"
	// Split on whitespace.  We require 5 or 6 fields.
	// (second) (minute) (hour) (day of month) (month) (day of week, optional)
	Cron       string              `json:"cron"`        // cron表达式。若为一次性任务，则该值为空。
	CronRender string              `json:"cron_render"` // cron表达式UI渲染信息
	Target     map[string][]string `json:"target"`      // 任务作用目标。map中的value暂定为[]string，可能发生变化。
	Status     string              `json:"status"`      // 任务状态。
	Creator    struct {
		ID        string `json:"id"`
		LoginName string `json:"login_name"`
		Name      string `json:"name"`
	} `json:"creator"` // 任务创建者
	CreatedAt times.ISOTime `json:"created_at"` // 创建时间
	UpdatedAt times.ISOTime `json:"updated_at"` // 更新时间
}

// GetJobPageReq 查询满足条件的任务分页列表请求结构体
type GetJobPageReq struct {
	// 任务创建人ID
	Creator string `json:"creator"`
	// 标题(支持模糊查询)
	Title string `json:"title"`
	// 是否是内建任务。可选值: yes-是; no-否;
	// Enum: yes,no
	Builtin string `json:"builtin"`
	// 类别。可选值: inspection-硬件巡检任务; installation_timeout-装机超时检查任务;release_ip-释放IP任务;
	// Enum: inspection,installation_timeout,release_ip
	Category string `json:"category"`
	// 执行频率。可选值: immediately-立刻; fixed_rate-定时;
	// Enum: immediately,fixed_rate
	Rate string `json:"rate"`
	// 状态。可选值: running-运行中; paused-已暂停; stoped-已停止; deleted-已删除; 支持多个值（用英文逗号分隔）组合过滤。
	// Enum: running,paused,stoped,deleted
	Status string `json:"status"`
	// 分页页号
	Page int64 `json:"page"`
	// 分页大小
	PageSize int64 `json:"page_size"`
	// 当前用户
	CurrentUser *model.CurrentUser `json:"-"`
}

// FieldMap 请求字段映射
func (reqData *GetJobPageReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.Creator:  "creator",
		&reqData.Title:    "title",
		&reqData.Builtin:  "builtin",
		&reqData.Category: "category",
		&reqData.Rate:     "rate",
		&reqData.Status:   "status",
		&reqData.Page:     "page",
		&reqData.PageSize: "page_size",
	}
}

// Validate 结构体数据校验
func (reqData *GetJobPageReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if reqData.Builtin != "" &&
		reqData.Builtin != job.Builtin &&
		reqData.Builtin != job.NoBuiltin {
		errs.Add([]string{"builtin"}, binding.BusinessError, fmt.Sprintf("无效的字段值builtin=%s", reqData.Builtin))
		return errs
	}

	if reqData.Category != "" &&
		reqData.Category != job.CategoryInspection &&
		reqData.Category != job.CategoryInstallationTime {
		errs.Add([]string{"category"}, binding.BusinessError, fmt.Sprintf("无效的字段值category=%s", reqData.Category))
		return errs
	}

	if reqData.Rate != "" &&
		reqData.Rate != job.RateFixedRate &&
		reqData.Rate != job.RateImmediately {
		errs.Add([]string{"rate"}, binding.BusinessError, fmt.Sprintf("无效的字段值rate=%s", reqData.Rate))
		return errs
	}

	if reqData.Status != "" &&
		reqData.Status != job.Running &&
		reqData.Status != job.Paused &&
		reqData.Status != job.Stoped &&
		reqData.Status != job.Deleted {
		errs.Add([]string{"status"}, binding.BusinessError, fmt.Sprintf("无效的字段值status=%s", reqData.Status))
		return errs
	}

	return errs
}

// GetJobPage 返回满足过滤条件的任务分页列表
func GetJobPage(log logger.Logger, conf *config.Config, jobmgr job.Manager, reqData *GetJobPageReq) (pg *page.Page, err error) {
	if reqData.PageSize <= 0 {
		reqData.PageSize = 10
	}
	if reqData.PageSize > 100 {
		reqData.PageSize = 100
	}
	if reqData.Page < 0 {
		reqData.Page = 0
	}

	if reqData.Status == "" {
		reqData.Status = fmt.Sprintf("%s,%s,%s", job.Running, job.Stoped, job.Paused) // 默认不加载'已删除'状态的任务。
	}

	cond := job.Job{
		Builtin:  reqData.Builtin,
		Title:    reqData.Title,
		Category: reqData.Category,
		Rate:     reqData.Rate,
		Status:   reqData.Status,
		Creator:  reqData.Creator,
	}

	totalRecords, items, err := jobmgr.GetJobs(&cond, reqData.Page, reqData.PageSize)
	if err != nil {
		return nil, err
	}

	pager := page.NewPager(reflect.TypeOf(&Job{}), reqData.Page, reqData.PageSize, totalRecords)
	for i := range items {
		pager.AddRecords(convert2Job(log, conf, reqData.CurrentUser.Token, items[i]))
	}
	return pager.BuildPage(), nil
}

// GetJobByIDReq 查询目标任务请求结构体
type GetJobByIDReq struct {
	ID string `json:"id"`
	// 当前用户
	CurrentUser *model.CurrentUser `json:"-"`
}

// GetJobByID 查询目标任务明细
func GetJobByID(log logger.Logger, conf *config.Config, jobmgr job.Manager, reqData *GetJobByIDReq) (*Job, error) {
	cjob, err := jobmgr.GetJobByID(reqData.ID)
	if err != nil {
		return nil, err
	}
	return convert2Job(log, conf, reqData.CurrentUser.Token, cjob), nil
}

// convert2Job 将模型层的CombinedDevice转换成服务层的CombinedDevice。
func convert2Job(log logger.Logger, conf *config.Config, token string, item *job.Job) *Job {
	j := Job{
		ID:         item.ID,
		Builtin:    item.Builtin,
		Title:      item.Title,
		Category:   item.Category,
		Rate:       item.Rate,
		Cron:       item.CronExpression,
		CronRender: item.CronRender,
		Target:     item.Target,
		Status:     item.Status,
		CreatedAt:  times.ISOTime(item.CreatedAt),
		UpdatedAt:  times.ISOTime(item.UpdatedAt),
	}
	if u, _ := GetUserByID(log, conf, token, item.Creator); u != nil {
		j.Creator.ID = u.ID
		j.Creator.LoginName = u.LoginName
		j.Creator.Name = u.Name
	}
	return &j

}

// PauseJob 暂停运行中的目标定时任务
func PauseJob(jobmgr job.Manager, jobid string) (err error) {
	return jobmgr.Pause(jobid)
}

// UnpauseJob 继续已暂停的目标定时任务
func UnpauseJob(jobmgr job.Manager, jobid string) (err error) {
	return jobmgr.Unpause(jobid)
}

// RemoveJob 删除非内置任务
func RemoveJob(jobmgr job.Manager, jobid string) (err error) {
	return jobmgr.Remove(jobid)
}

// AddInspectionJobReq 新增硬件巡检任务请求结构体
type AddInspectionJobReq struct {
	// 源节点（分布式部署条件下）
	OriginNode string `json:"-"`
	// 任务创建者/操作人
	Creator string `json:"-"`
	// 任务标题
	Title string `json:"title"`
	// 目标设备
	SN []string `json:"sn"`
	// 执行频率
	Rate string `json:"rate"`
	// cron表达式
	Cron string `json:"cron"`
	// cron表达式UI渲染信息
	CronRender string `json:"cron_render"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *AddInspectionJobReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.SN:         "sn",
		&reqData.Title:      "title",
		&reqData.Rate:       "rate",
		&reqData.Cron:       "cron",
		&reqData.CronRender: "cron_render",
	}
}

// Validate 结构体数据校验
func (reqData *AddInspectionJobReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())

	// 校验SN是否在库
	var cond model.Device
	for _, sn := range reqData.SN {
		cond.SN = sn
		if count, _ := repo.CountDevices(&cond); count <= 0 {
			errs.Add([]string{"sn"}, binding.BusinessError, fmt.Sprintf("设备(%s)不存在", sn))
			return errs
		}
	}

	if reqData.Rate == "" {
		errs.Add([]string{"rate"}, binding.RequiredError, fmt.Sprintf("基本参数要求(巡检频率不能为空)"))
		return errs
	}

	if !collection.InSlice(reqData.Rate, []string{job.RateImmediately, job.RateFixedRate}) {
		errs.Add([]string{"rate"}, binding.BusinessError, fmt.Sprintf("基本参数要求(巡检频率(%s)只能为:(%s/%s))", reqData.Rate, job.RateImmediately, job.RateFixedRate))
		return errs
	}

	// 固定周期任务必须填写cron表达式
	if reqData.Rate == job.RateFixedRate && reqData.Cron == "" {
		errs.Add([]string{"cron"}, binding.RequiredError, fmt.Sprintf("基本参数要求(固定周期情况下Cron不能为空)"))
		return errs
	}

	// TODO 校验cron表达式有效性
	return errs
}

var (
	// ErrInspectionDevicesNotSpecified 未指定硬件巡检设备
	ErrInspectionDevicesNotSpecified = errors.New("inspection devices not specified")
)

// AddInspectionJob 新增硬件巡检任务并返回任务ID
func AddInspectionJob(log logger.Logger, repo model.Repo, conf *config.Config, jobmgr job.Manager, reqData *AddInspectionJobReq) (jobid string, err error) {
	// 若未指定SN列表，则获取历史已巡检的SN列表。
	//if len(reqData.SN) <= 0 {
	//	if reqData.SN, err = repo.GetInspectedSN(); err != nil {
	//		return "", err
	//	}
	//}
	//if len(reqData.SN) <= 0 {
	//	return "", ErrInspectionDevicesNotSpecified
	//}
	// 若未指定SN列表，则为空（实际巡检任务会全量巡检）
	var target map[string][]string
	if len(reqData.SN) > 0 {
		targetPairs, err := getOriginNodeDevicesPairs(repo, reqData.SN...)
		if err != nil {
			return "", err
		}
		target = targetPairs
	}

	if reqData.Rate == job.RateImmediately {
		reqData.Cron = ""
	}

	// 2、生成全局唯一的任务ID
	jobid = utils.UUID()

	// 3、将任务提交至任务管理器
	if err = jobmgr.Submit(&job.Job{
		ID:             jobid,
		Creator:        reqData.Creator,
		Builtin:        job.NoBuiltin,
		Title:          fmt.Sprintf("[硬件巡检]%s", reqData.Title),
		Category:       job.CategoryInspection,
		Rate:           reqData.Rate,
		CronExpression: reqData.Cron,
		CronRender:     reqData.CronRender,
		Target:         target,
		Status:         job.Running,
		CronJob:        mysql.NewInspectionJob(log, repo, conf, jobid),
	}); err != nil {
		return "", err
	}
	return jobid, nil
}

// getOriginNodeDevicesPairs 根据config文件 NODE IP 构成的键值对
// 实现各机房管理单元NODE节点对应自身机房的SN，多个NODE节点时轮询分配
func getOriginNodeDevicesPairs(repo model.Repo, sns ...string) (pairs map[string][]string, err error) {
	pairs = make(map[string][]string)
	var nodeIPs []string
	var nodeIPLength int
	for i, sn := range sns {
		dev, err := repo.GetDeviceBySN(sn)
		if err != nil {
			return nil, err
		}
		// 忽略待退役、已退役设备
		if dev.OperationStatus == model.DevOperStatPreRetire || dev.OperationStatus == model.DevOperStatRetiring || dev.OperationStatus == model.DevOperStateRetired {
			continue
		}

		if dev.ServerRoomID != 0 {
			if nodeIPLength = len(middleware.MapDistributeNode.MDistribute[dev.ServerRoomID]); nodeIPLength != 0 {
				nodeIPs = middleware.MapDistributeNode.MDistribute[dev.ServerRoomID]
				_, ok := pairs[nodeIPs[i%nodeIPLength]]
				if !ok {
					pairs[nodeIPs[i%nodeIPLength]] = []string{}
				}
				pairs[nodeIPs[i%nodeIPLength]] = append(pairs[nodeIPs[i%nodeIPLength]], dev.SN)
			} else {
				_, ok := pairs["master"]
				if !ok {
					pairs["master"] = []string{}
				}
				pairs["master"] = append(pairs["master"], dev.SN)
			}
		}
	}
	return pairs, nil
}
