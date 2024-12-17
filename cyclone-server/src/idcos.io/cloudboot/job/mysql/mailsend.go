package mysql

import (
	"strings"
	"github.com/jinzhu/gorm"
	"idcos.io/cloudboot/utils/times"
	"errors"
	"encoding/base64"
	"bytes"
	"fmt"
	"time"
	"encoding/json"

	"idcos.io/cloudboot/config"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils"
)

const (
	// 邮件任务Title
	MailTitleDevicesExpired = "过保设备信息"
	MailTitleDevicesPreDeploy = "待部署设备信息"
)

// MailSendJob
type MailSendJob struct {
	id   string // 任务ID(全局唯一)
	log  logger.Logger
	repo model.Repo
	conf *config.Config
}

// NewMailSendJob 实例化任务管理器
func NewMailSendJob(log logger.Logger, repo model.Repo, conf *config.Config, jobid string) *MailSendJob {
	return &MailSendJob{
		log:  log,
		repo: repo,
		conf: conf,
		id:   jobid,
	}
}

// 内置任务，自动发送邮件
// 根据job的target字段（json）进行区分发送内容
func (j *MailSendJob) Run() {
	j.log.Debugf("Start Job: Mail send")
	defer j.log.Debugf("Mail send is completed")

	defer func() {
		if err := recover(); err != nil {
			j.log.Errorf("Mail send job panic: \n%s", err)
		}
	}()

	// 1、根据任务ID获取任务明细
	mjob, err := j.repo.GetJobByID(j.id)
	if err != nil {
		j.log.Error(err)
		return
	}
	j.log.Debugf("[%s]MailSendJob: %+v", j.id, mjob)

	// 2、提取出邮件任务的请求参数: map [From To CC Title] []string
	var target map[string][]string
	if err = json.Unmarshal([]byte(mjob.Target), &target); err != nil {
		j.log.Error(err)
		return
	}
	// 校验参数
	if len(target["From"]) != 1 {
		j.log.Error("发件人邮箱未定义或存在多个")
		return
	}

	if len(target["Title"]) != 1 {
		j.log.Error("邮件标题未定义或存在多个")
		return
	} else {
		// 3、根据邮件标题调用不同的邮件内容
		switch target["Title"][0] {
		case MailTitleDevicesExpired:
			if err = j.MailSendDevicesExpired(target); err != nil {
				j.log.Error(err)
				return
			}
		case MailTitleDevicesPreDeploy:
			if err = j.MailSendDevicesPreDeploy(target); err != nil {
				j.log.Error(err)
				return
			}
		default:
			j.log.Error("邮件任务未定义")
			return
		}
	}
}

//DeviceSendmailCondition 物理机邮件发送列表搜索字段
type DeviceSendmailCondition struct {
	//运营状态:运营中(需告警),运营中(无需告警),重装中,搬迁中,待退役,退役中,已退役,待部署,已上架,回收中
	OperationStatus string `json:"operation_status"`
	//启用日期
	StartedAt  string `json:"started_at"`
	// 预部署状态物理机（没有任何安装记录的物理机）。
	//PreDeployed bool `json:"pre_deployed"`
	//Page           int64 `json:"-"`
	//PageSize       int64 `json:"-"`
}

// 用于邮件发送的设备信息结构体，注意IP为敏感字段，不可发送
type DeviceForSendmail struct {
	//固资编号
	FixedAssetNum string `json:"fixed_asset_number"`
	//序列号
	SN string `json:"sn"`
	//厂商
	Vendor string `json:"vendor"`
	//型号
	Model string `json:"model"`
	// 硬件架构
	Arch string `json:"arch"`
	//用途 
	Usage string `json:"usage"`
	//设备类型
	Category string `json:"category"`
	// TOR
	//TOR string `json:"tor"`
	//硬件说明
	HardwareRemark string `json:"hardware_remark"`
	//RAID说明
	//RAIDRemark string `json:"raid_remark"`
	//启用时间
	StartedAt string `json:"started_at"`
	//运营状态:运营中(需告警),运营中(无需告警),重装中,搬迁中,待退役,退役中,已退役,待部署,已上架,回收中
	OperationStatus string `json:"operation_status"`
	//数据中心
	IDC *IDCSimplify `json:"idc"`
	//机房管理单元
	ServerRoom *ServerRoomSimplify `json:"server_room"`
	//机架
	ServerCabinet *ServerCabinetSimplify `json:"server_cabinet"`
	//机位
	ServerUSite *ServerUSiteSimplify `json:"server_usite"`
	// from DeviceLifecycle [负责人 维保截止日期 维保状态]
	Owner			 				string		`json:"owner"`
	MaintenanceServiceDateEnd       string  	`json:"maintenance_service_date_end"`
	MaintenanceServiceStatus		string		`json:"maintenance_service_status"`	
}

// IDCSimplify 数据中心信息
type IDCSimplify struct {
	//数据中心ID
	ID uint `json:"id"`
	//数据中心名称
	Name string `json:"name"`
}

// ServerRoomSimplify 机房管理单元
type ServerRoomSimplify struct {
	//机房管理单元ID
	ID uint `json:"id"`
	//机房管理单元名称
	Name string `json:"name"`
}

// ServerCabinetSimplify 机架
type ServerCabinetSimplify struct {
	//机架ID
	ID uint `json:"id"`
	//机架编号
	Number string `json:"number"`
}

// ServerUSiteSimplify 机位
type ServerUSiteSimplify struct {
	//机位ID
	ID uint `json:"id"`
	//机位编号
	Number string `json:"number"`
	// 物理区域
	PhysicalArea string `json:"physical_area"`
}

// AttachmentDevices 邮件附件设备信息集合
type AttachmentDevices []*DeviceForSendmail

// ToTableRecords 生成用于表格显示的二维字符串切片
func (items AttachmentDevices) ToTableRecords() (records [][]string) {
	records = make([][]string, 0, len(items))

	for i := range items {
		idcName := ""
		if items[i].IDC != nil {
			idcName = items[i].IDC.Name
		}
		serverRoomName := ""
		if items[i].ServerRoom != nil {
			serverRoomName = items[i].ServerRoom.Name
		}
		serverCabinetNumber := ""
		if items[i].ServerCabinet != nil {
			serverCabinetNumber = items[i].ServerCabinet.Number
		}
		serverUsiteNumber := ""
		physicalArea := ""
		if items[i].ServerUSite != nil {
			serverUsiteNumber = items[i].ServerUSite.Number
			physicalArea = items[i].ServerUSite.PhysicalArea
		}		

		records = append(records, []string{
			items[i].FixedAssetNum,
			items[i].SN,
			items[i].Vendor,
			items[i].Model,
			items[i].Arch,
			items[i].Usage,
			items[i].Category,
			items[i].HardwareRemark,
			items[i].StartedAt,
			items[i].OperationStatus,
			idcName,
			serverRoomName,
			physicalArea,
			serverCabinetNumber,
			serverUsiteNumber,
			items[i].Owner,
			items[i].MaintenanceServiceDateEnd,
			items[i].MaintenanceServiceStatus,
		})
	}
	return records
}

// 通过 DeviceSendmailCondition 条件查询获取用于邮件发送的设备信息
func (j *MailSendJob) GetDeviceForSendmail(log logger.Logger, repo model.Repo, cond *DeviceSendmailCondition) (rst []*DeviceForSendmail, err error) {
	rst = make([]*DeviceForSendmail, 0)
	if cond.StartedAt != "" {
		items,err := repo.GetDeviceByStartedAt(cond.StartedAt)
		if err != nil {
			log.Errorf("[FAILED]获取启用日期为 %s 前的设备信息, err: %s", cond.StartedAt, err.Error())
			return nil,err
		}
		for i := range items {
			if items[i].OperationStatus != "retired" {
				item, err := convert2SendmailResult(log, repo, items[i])
				if err != nil {
					return nil, err
				}
				if item != nil {
					rst = append(rst, item)
				}
			}
		}
		return rst, nil
	} else if cond.OperationStatus != "" {
		items,err := repo.GetDevices(&model.Device{OperationStatus: cond.OperationStatus}, nil, nil)
		if err != nil {
			log.Errorf("[FAILED]获取运营状态为 %s 前的设备信息, err: %s", cond.OperationStatus, err.Error())
			return nil,err
		}
		for i := range items {
			item, err := convert2SendmailResult(log, repo, items[i])
			if err != nil {
				return nil, err
			}
			if item != nil {
				rst = append(rst, item)
			}
		}
		return rst, nil
	} else {
		return nil,err
	}
}

// 将model/device.GetDevices 的返回结果进行转换
func convert2SendmailResult(log logger.Logger, repo model.Repo, d *model.Device) (*DeviceForSendmail, error) {
	if d == nil {
		return nil, nil
	}
	//映射运营状态:运营中(需告警),运营中(无需告警),重装中,搬迁中,待退役,退役中,已退役,待部署,已上架,回收中'
	OperationStatusMap :=map[string]string {
		"run_with_alarm":"运营中(需告警)",
		"run_without_alarm":"运营中(无需告警)",
		"reinstalling":"重装中",
		"moving":"搬迁中",
		"pre_retire":"待退役",
		"retiring":"退役中",		
		"retired":"已退役",
		"pre_deploy":"待部署",
		"on_shelve":"已上架",
		"recycling":"回收中",
	}
	// 维保状态枚举值:在保-under_warranty;过保-out_of_warranty;未激活-inactive'
	MaintenanceServiceStatusMap :=map[string]string {
		"under_warranty":"在保",
		"out_of_warranty":"过保",
		"inactive":"未激活",
	}	

	result := DeviceForSendmail{
		FixedAssetNum:       d.FixedAssetNumber,
		SN:                  d.SN,
		Vendor:              d.Vendor,
		Model:               d.DevModel,
		Arch:                d.Arch,
		Usage:               d.Usage,
		Category:            d.Category,
		HardwareRemark:      d.HardwareRemark,
		StartedAt:           times.ISOTime(d.StartedAt).ToDateStr(),
		OperationStatus:     OperationStatusMap[d.OperationStatus],
		IDC:                 &IDCSimplify{ID: d.IDCID},
		ServerRoom:          &ServerRoomSimplify{ID: d.ServerRoomID},
		ServerCabinet: 		 &ServerCabinetSimplify{ID: d.CabinetID},
	}

	idc, err := repo.GetIDCByID(d.IDCID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if idc != nil {
		result.IDC.Name = idc.Name
	}

	sroom, err := repo.GetServerRoomByID(d.ServerRoomID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if sroom != nil {
		result.ServerRoom.Name = sroom.Name
	}

	if cabinet, err := repo.GetServerCabinetByID(d.CabinetID); err == nil {
		result.ServerCabinet.Number = cabinet.Number
	}

	if d.USiteID != nil {
		result.ServerUSite = &ServerUSiteSimplify{ID: *d.USiteID}
		if u, err := repo.GetServerUSiteByID(*d.USiteID); err == nil {
			result.ServerUSite.Number = u.Number
			result.ServerUSite.PhysicalArea = u.PhysicalArea
		}
	}

	// DeviceLifecycle 查询是否已经存在
	devLifecycle, err := repo.GetDeviceLifecycleBySN(d.SN)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if devLifecycle != nil {
		result.Owner = devLifecycle.Owner
		result.MaintenanceServiceDateEnd = times.ISOTime(devLifecycle.MaintenanceServiceDateEnd).ToDateStr()
		result.MaintenanceServiceStatus = MaintenanceServiceStatusMap[devLifecycle.MaintenanceServiceStatus]
	} else {
		result.Owner = "Undefined"
		result.MaintenanceServiceDateEnd = "Undefined"
		result.MaintenanceServiceStatus = "Undefined"
	}
	return &result, nil
}

// 邮件Title:"过保设备信息"  内容：发送启用日期 >= 4.5 year 设备信息
func (j *MailSendJob) MailSendDevicesExpired(smq map[string][]string) (err error) {

	sendmailreq := middleware.SendMailReq {
		From:    strings.Join(smq["From"],";"),
		To:      strings.Join(smq["To"],";"),
		CC:		 strings.Join(smq["CC"],";"),
		Title:	 strings.Join(smq["Title"],"-"),
		BodyFormat: "1",
		Priority: "1",
	}

	//根据 now 获取 4.5 年前的日期并作为‘启用日期’条件查询设备
	currentTime := time.Now()
	pastTime := currentTime.AddDate(-4, -6, 0).Format("2006-01-02")
	var cond DeviceSendmailCondition
	cond.StartedAt = pastTime

	// 根据启用日期条件获取 4.5 前的设备信息
	deviceData,err := j.GetDeviceForSendmail(j.log, j.repo, &cond)
	if err != nil {
		j.log.Errorf("获取启用日期大于4.5年的设备信息失败, err: %s", err.Error())
		return err
	}

	if len(deviceData) == 0 {
		j.log.Error("启用日期大于4.5年的设备数量为： 0")
		return errors.New("启用日期大于4.5年的设备数量为： 0")
	}
	j.log.Debugf("启用日期大于4.5年的设备数量为： %v", len(deviceData))
	sendmailreq.Content = fmt.Sprintf("[From Dolphin] 共计 %v 台设备启用日期大于4.5年，详细信息见附件。", len(deviceData))
	
	// 设备信息写入buffer
	file, err := utils.WriteToXLSX(utils.FileDeviceSendMail, AttachmentDevices(deviceData).ToTableRecords())
	if err != nil {
		j.log.Errorf("设备信息生成xls失败, err: %s", err.Error())
		return err
	}
	buf := bytes.Buffer{}
	err = file.Write(&buf)
	if err != nil {
		return err
	}

	//读取buffer内容并编码 作为附件发送
	bufferRead := buf.Bytes()
	// 附件的参数格式： [{"name": "device_expire.xlsx", "Data": base64data}]
	base64data := base64.StdEncoding.EncodeToString(bufferRead)
	sendmailreq.Attachments = []middleware.MailAttachment {{Name:"device_expire.xlsx", Data:base64data}}


	//调用邮件发送api请求
	 err = middleware.SendMail(j.log, j.repo, j.conf, &sendmailreq)
	if err != nil {
		j.log.Errorf("SendMail fail, err: %s", err.Error())
		return err
	}
	return nil
}

// 触发邮件发送，发送待部署状态设备信息通知一线进行处理
func (j *MailSendJob) MailSendDevicesPreDeploy(smq map[string][]string) (err error) {

	sendmailreq := middleware.SendMailReq {
		From:    strings.Join(smq["From"],";"),
		To:      strings.Join(smq["To"],";"),
		CC:		 strings.Join(smq["CC"],";"),
		Title:	 strings.Join(smq["Title"],"-"),
		BodyFormat: "1",
		Priority: "1",
	}

	// 获取运营状态为待部署设备信息
	var cond DeviceSendmailCondition
	cond.OperationStatus = "pre_deploy"
	deviceData,err := j.GetDeviceForSendmail(j.log, j.repo, &cond)
	if err != nil {
		j.log.Errorf("get pre_deploy devices fail,%v", err)
		return err
	}

	if len(deviceData) == 0 {
		j.log.Error("状态为待部署的设备数量为： 0")
		return errors.New("状态为待部署的设备数量为： 0")
	}
	j.log.Debugf("状态为待部署的设备数量为： %v", len(deviceData))
	sendmailreq.Content = fmt.Sprintf("From dolphin 共计 %v 台设备，详细信息见附件。", len(deviceData))
	
	// 设备信息写入buffer
	file, err := utils.WriteToXLSX(utils.FileDeviceSendMail, AttachmentDevices(deviceData).ToTableRecords())
	if err != nil {
		j.log.Errorf("[FAILED]设备信息生成xlsx, err: %s", err.Error())
		return err
	}
	buf := bytes.Buffer{}
	err = file.Write(&buf)
	if err != nil {
		return err
	}

	//读取buffer内容并编码 作为附件发送
	bufferRead := buf.Bytes()
	// 附件的参数格式： [{"name": "device_pre_deploy.xlsx", "Data": base64data}]
	base64data := base64.StdEncoding.EncodeToString(bufferRead)
	sendmailreq.Attachments = []middleware.MailAttachment {{Name:"device_pre_deploy.xlsx", Data:base64data}}

	//调用邮件发送api请求
	err = middleware.SendMail(j.log, j.repo, j.conf, &sendmailreq)
	if err != nil {
		j.log.Errorf("SendMail fail, err: %s", err.Error())
		return err
	}
	return nil
}