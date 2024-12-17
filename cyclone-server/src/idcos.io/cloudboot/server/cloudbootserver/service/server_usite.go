package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/voidint/binding"
	"github.com/voidint/page"

	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils"
	"idcos.io/cloudboot/utils/collection"
	strings2 "idcos.io/cloudboot/utils/strings"
	"idcos.io/cloudboot/utils/upload"
	"idcos.io/cloudboot/utils/user"
)

// SaveServerUSiteReq 保存机位(U位)请求结构体
type SaveServerUSiteReq struct {
	// 所属机房ID
	ServerRoomID uint `json:"server_room_id"`
	// 所属机架(柜)ID
	ServerCabinetID uint `json:"server_cabinet_id"`
	// 机位(U位)ID。若id=0，则新增。若id>0，则修改。
	ID uint `json:"id"`
	// 机位(U位)编号
	Number string `json:"number"`
	// 机位(U位)高度
	Height uint `json:"height"`
	// 起始U数
	Beginning uint `json:"beginning"`
	// 备注
	Remark string `json:"remark"`
	// 物理区域
	PhysicalArea string `json:"physical_area"`
	// 管理网交换机
	OobnetSwitches []*model.SwitchInfo `json:"oobnet_switches"`
	// 内网交换机
	IntranetSwitches []*model.SwitchInfo `json:"intranet_switches"`
	// 外网交换机
	ExtranetSwitches []*model.SwitchInfo `json:"extranet_switches"`
	// 内外网端口速率
	LAWAPortRate	string	`json:"la_wa_port_rate"`
	LoginUser        *model.CurrentUser  `json:"-"`
}

// ServerUSiteResp 机位返回信息体
type ServerUSiteResp struct {
	IDC struct {
		// 数据中心ID
		ID uint `json:"id"`
		// 数据中心名称
		Name string `json:"name"`
	} `json:"idc"`

	ServerRoom struct {
		// 机房ID
		ID uint `json:"id"`
		// 机房名称
		Name string `json:"name"`
	} `json:"server_room"`

	ServerCabinet struct {
		// 机架ID
		ID uint `json:"id"`
		// 机架编号
		Number string `json:"number"`
	} `json:"server_cabinet"`

	// 机位(U位)ID。若id=0，则新增。若id>0，则修改。
	ID uint `json:"id"`
	// 机位(U位)编号
	Number string `json:"number"`
	// 机位(U位)高度
	Height uint `json:"height"`
	// 起始U数
	Beginning uint `json:"beginning"`
	// 物理区域
	PhysicalArea string `json:"physical_area"`
	// 管理网交换机
	OobnetSwitches []*model.SwitchInfo `json:"oobnet_switches"`
	// 内网交换机
	IntranetSwitches []*model.SwitchInfo `json:"intranet_switches"`
	// 外网交换机
	ExtranetSwitches []*model.SwitchInfo `json:"extranet_switches"`
	// 内外网端口速率
	LAWAPortRate	string	`json:"la_wa_port_rate"`
	// 备注
	Remark string `json:"remark"`
	// 创建时间
	CreatedAt string `json:"created_at"`
	// 修改时间
	UpdatedAt string `json:"updated_at"`
	// 机位状态
	// USiteStatFree = "free" 机位(U位)状态-空闲
	// USiteStatPreOccupied = "pre_occupied" 机位(U位)状态-预占用
	// USiteStatUsed = "used" 机位(U位)状态-已使用
	// USiteStatDisabled = "disabled" 机位(U位)状态-不可用
	Status *string `json:"status"`
	// 状态
	Creator string `json:"creator"`
}

// FieldMap 请求字段映射
func (reqData *SaveServerUSiteReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.ServerRoomID:     "server_room_id",
		&reqData.ServerCabinetID:  "server_cabinet_id",
		&reqData.ID:               "id",
		&reqData.Number:           "number",
		&reqData.Height:           "height",
		&reqData.Beginning:        "beginning",
		&reqData.PhysicalArea:     "logical_area",
		&reqData.Remark:           "remark",
		&reqData.OobnetSwitches:   "oobnet_switches",
		&reqData.IntranetSwitches: "intranet_switches",
		&reqData.ExtranetSwitches: "extranet_switches",
		&reqData.LAWAPortRate:     "la_wa_port_rate",
	}
}

// Validate 结构体数据校验
func (reqData *SaveServerUSiteReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())

	if reqData.ServerRoomID == 0 {
		errs.Add([]string{"server_room_id"}, binding.RequiredError, "机房管理单元ID不能为空")
		return errs
	}
	if reqData.ServerCabinetID == 0 {
		errs.Add([]string{"network_area_id"}, binding.RequiredError, "机架ID不能为空")
		return errs
	}
	if reqData.Number == "" {
		errs.Add([]string{"number"}, binding.RequiredError, "机位编号不能为空")
		return errs
	}
	if reqData.Height <= 0 {
		errs.Add([]string{"height"}, binding.RequiredError, "机位高度须大于0")
		return errs
	}
	if reqData.Beginning <= 0 {
		errs.Add([]string{"network_rate"}, binding.RequiredError, "起始U数须大于0")
		return errs
	}

	// 校验指定ID的机架(柜)是否存在
	cabinet, err := repo.GetServerCabinetByID(reqData.ServerCabinetID)
	if err == gorm.ErrRecordNotFound {
		errs.Add([]string{"idc_id", "server_room_id", "server_cabinet_id"}, binding.BusinessError, "该机架不存在")
		return errs
	}
	// 校验机位高度不能超过机架高度
	// 机架高度为N，事际上是有N+1个机位
	if cabinet.Height+1 < reqData.Height {
		errs.Add([]string{"height"}, binding.BusinessError, fmt.Sprintf("机位高度 %d, 不能超过所在机架高度 %d", reqData.Height, cabinet.Height))
		return errs
	}
	// 校验机位开始U位不能超过机架高度
	if cabinet.Height < reqData.Beginning {
		errs.Add([]string{"beginning"}, binding.BusinessError, fmt.Sprintf("机位开始U位 %d, 不能超过所在机架高度 %d", reqData.Beginning, cabinet.Height))
		return errs
	}
	// 校验机位开始U位+机位调试不能超过机架高度
	if cabinet.Height < reqData.Beginning+(reqData.Height-1) {
		errs.Add([]string{"beginning"}, binding.BusinessError, fmt.Sprintf("机位开始U位 %d, U位高度 %d, 不能超过所在机架高度 %d", reqData.Beginning, reqData.Height, cabinet.Height))
		return errs
	}

	if err != nil {
		errs.Add([]string{"idc_id", "server_room_id", "server_cabinet_id"}, binding.SystemError, "系统内部错误")
		return errs
	}
	if cabinet.ServerRoomID != reqData.ServerRoomID || cabinet.ID != reqData.ServerCabinetID {
		errs.Add([]string{"server_cabinet_id", "server_room_id"}, binding.BusinessError, "机房管理单元、机架不匹配")
		return errs
	}

	cond := model.CombinedServerUSite{
		ServerCabinetID: []uint{reqData.ServerCabinetID},
		ServerRoomID:    []uint{reqData.ServerRoomID},
	}
	usites, err := repo.GetServerUSiteByCond(&cond, nil, nil)
	for _, usite := range usites {
		if usite.Number == reqData.Number && reqData.ID != usite.ID {
			errs.Add([]string{"number"}, binding.BusinessError, fmt.Sprintf("机位(%s)在所属机架内已经存在", reqData.Number))
			return errs
		}
	}

	// 内外网端口速率校验：GE\10GE\25GE\40GE
	if _, ok := model.PortRateMap[reqData.LAWAPortRate]; !ok {
		errs.Add([]string{"la_wa_port_rate"}, binding.BusinessError, fmt.Sprintf("机位(%s)关联内外网端口速率（%s）不正确，应为：GE|10GE|25GE|40GE", reqData.Number, reqData.LAWAPortRate))
		return errs
	}

	// 如果新增，做唯一性校验
	if reqData.ID == 0 {
		cond := model.CombinedServerUSite{
			CabinetNumber:   reqData.Number,
			ServerRoomID:    []uint{reqData.ServerRoomID},
			ServerCabinetID: []uint{reqData.ServerCabinetID},
		}
		usites, err := repo.GetServerUSiteByCond(&cond, nil, nil)
		if err != nil {
			errs.Add([]string{"repeat"}, binding.SystemError, "系统内部错误")
			return errs
		}
		for _, usite := range usites {
			if usite.Number == reqData.Number {
				errs.Add([]string{"number"}, binding.BusinessError, fmt.Sprintf("机位Number %s在所属机架内已经存在", reqData.Number))
				return errs
			}
		}
	}

	return errs
}

// 批量修改机位状态入参结构体（机房管理单元+机架+机位）
type ServerUSiteForUpdate struct {
	ServerRoomName         string    `json:"server_room_name"`
	ServerCabinetNumber    string    `json:"server_cabinet_number"`
	ServerUsiteNumber      string    `json:"server_usite_number"`
}

// UpdateServerUSiteStatusReq 批量修改机位状态入参结构体
type UpdateServerUSiteStatusReq struct {
	// 状态
	Status    string             		`json:"status"`
	// id列表
	IDs       []uint             		`json:"ids"`
	// 机位结构体参数
	USites    []*ServerUSiteForUpdate	`json:"usites"`
	LoginUser *model.CurrentUser 		`json:"-"`
}

// FieldMap 请求字段映射
func (reqData *UpdateServerUSiteStatusReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDs:    "ids",
		&reqData.Status: "status",
		&reqData.USites: "usites",
	}
}

// Validate 结构体数据校验
func (reqData *UpdateServerUSiteStatusReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if len(reqData.IDs) == 0 && len(reqData.USites) == 0 {
		errs.Add([]string{"目标机位"}, binding.RequiredError, "不能为空")
		return errs
	}
	if reqData.Status == "" {
		errs.Add([]string{"status"}, binding.RequiredError, "状态值不能为空")
		return errs
	}

	if !collection.InSlice(reqData.Status, []string{model.USiteStatFree, model.USiteStatPreOccupied, model.USiteStatUsed, model.USiteStatDisabled}) {
		errs.Add([]string{"status"}, binding.RequiredError, "无效的状态值")
		return errs
	}
	return errs
}

// GetServerUSitePageReq 机位分页查询
type GetServerUSitePageReq struct {
	// 所属数据中心ID
	IDCID string `json:"idc_id"`
	// 所属机房ID
	ServerRoomID   string `json:"server_room_id"`
	ServerRoomName string
	// 所属机架ID
	ServerCabinetID string `json:"server_cabinet_id"`
	// 网络区域ID
	NetAreaID string `json:"network_area_id"`
	// 物理区域
	PhysicalArea string `json:"physical_area"`
	// 机架编号
	CabinetNumber string `json:"cabinet_number"`
	// 机位编号
	USiteNumber string `json:"usite_number"`
	// 机位调试
	Height string `json:"height"`
	// 机位状态
	// USiteStatFree = "free" 机位(U位)状态-空闲
	// USiteStatPreOccupied = "pre_occupied" 机位(U位)状态-预占用
	// USiteStatUsed = "used" 机位(U位)状态-已使用
	// USiteStatDisabled = "disabled" 机位(U位)状态-不可用
	Status string `json:"status"`
	// 内外网端口速率
	LAWAPortRate	string	`json:"la_wa_port_rate"`	
	// 分页页号
	Page int64 `json:"page"`
	// 分页大小
	PageSize int64 `json:"page_size"`
	// 从UAM上获取用户信息的钩子
	GetNameFromUAM user.GetNameFromUAM `json:"-"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *GetServerUSitePageReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDCID:           "idc_id",
		&reqData.ServerRoomID:    "server_room_id",
		&reqData.ServerRoomName:  "server_room_name",
		&reqData.ServerCabinetID: "server_cabinet_id",
		&reqData.CabinetNumber:   "cabinet_number",
		&reqData.USiteNumber:     "usite_number",
		&reqData.NetAreaID:       "network_area_id",
		&reqData.PhysicalArea:    "physical_area",
		&reqData.Height:          "height",
		&reqData.Status:          "status",
		&reqData.LAWAPortRate:    "la_wa_port_rate",
		&reqData.Page:            "page",
		&reqData.PageSize:        "page_size",
	}
}

// GetServerUSitePage 按条件查询机位分页列表
func GetServerUSitePage(log logger.Logger, repo model.Repo, reqData *GetServerUSitePageReq) (pg *page.Page, err error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 100 {
		reqData.PageSize = 10
	}
	if reqData.Page < 0 {
		reqData.Page = 0
	}

	cond := model.CombinedServerUSite{
		IDCID:           strings2.Multi2UintSlice(reqData.IDCID),
		ServerRoomID:    strings2.Multi2UintSlice(reqData.ServerRoomID),
		ServerCabinetID: strings2.Multi2UintSlice(reqData.ServerCabinetID),
		ServerRoomName:  reqData.ServerRoomName,
		NetAreaID:       strings2.Multi2UintSlice(reqData.NetAreaID),
		PhysicalArea:    reqData.PhysicalArea,
		CabinetNumber:   reqData.CabinetNumber,
		USiteNumber:     reqData.USiteNumber,
		Height:          reqData.Height,
		Status:          reqData.Status,
		LAWAPortRate:	 reqData.LAWAPortRate,
	}

	totalRecords, err := repo.CountServerUSite(&cond)
	if err != nil {
		return nil, err
	}

	pager := page.NewPager(reflect.TypeOf(&ServerUSiteResp{}), reqData.Page, reqData.PageSize, totalRecords)
	items, err := repo.GetServerUSiteByCond(&cond, model.OneOrderBy("id", model.DESC), pager.BuildLimiter())
	if err != nil {
		return nil, err
	}

	for i := range items {
		usite := convert2ServerUsite(repo, items[i])

		if usite != nil {
			_, name, _ := reqData.GetNameFromUAM(usite.Creator)
			usite.CreatedAt = name
		}

		pager.AddRecords(convert2ServerUsite(repo, items[i]))
	}
	return pager.BuildPage(), nil
}

func convert2ServerUsite(repo model.Repo, item *model.ServerUSite) *ServerUSiteResp {
	uSite := ServerUSiteResp{
		ID:               item.ID,
		Number:           item.Number,
		Height:           item.Height,
		Beginning:        item.Beginning,
		PhysicalArea:     item.PhysicalArea,
		OobnetSwitches:   toSwitchInfoArray(item.OobnetSwitches),
		IntranetSwitches: toSwitchInfoArray(item.IntranetSwitches),
		ExtranetSwitches: toSwitchInfoArray(item.ExtranetSwitches),
		Remark:           item.Remark,
		CreatedAt:        item.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        item.UpdatedAt.Format("2006-01-02 15:04:05"),
		Status:           &item.Status,
		LAWAPortRate:	  item.LAWAPortRate,
		Creator:          item.Creator,
	}

	if idc, _ := repo.GetIDCByID(item.IDCID); idc != nil {
		uSite.IDC.ID = idc.ID
		uSite.IDC.Name = idc.Name
	}

	if room, _ := repo.GetServerRoomByID(item.ServerRoomID); room != nil {
		uSite.ServerRoom.ID = room.ID
		uSite.ServerRoom.Name = room.Name
	}

	if cabinet, _ := repo.GetServerCabinetByID(item.ServerCabinetID); cabinet != nil {
		uSite.ServerCabinet.ID = cabinet.ID
		uSite.ServerCabinet.Number = cabinet.Number
	}

	return &uSite
}

// SaveServerUSite 保存机位(U位)
func SaveServerUSite(log logger.Logger, repo model.Repo, reqData *SaveServerUSiteReq) (err error) {
	var originUSite *model.ServerUSite
	if reqData.ID > 0 {
		originUSite, err = repo.GetServerUSiteByID(reqData.ID)
		if err != nil {
			return err
		}
	}

	// 通过机架去查找idc_id
	cabinet, _ := repo.GetServerCabinetByID(reqData.ServerCabinetID)

	uSite := model.ServerUSite{
		IDCID:            cabinet.IDCID,
		ServerRoomID:     reqData.ServerRoomID,
		ServerCabinetID:  reqData.ServerCabinetID,
		Number:           reqData.Number,
		Beginning:        reqData.Beginning,
		OobnetSwitches:   toJSONString(reqData.OobnetSwitches),
		IntranetSwitches: toJSONString(reqData.IntranetSwitches),
		ExtranetSwitches: toJSONString(reqData.ExtranetSwitches),
		Height:           reqData.Height,
		PhysicalArea:     reqData.PhysicalArea,
		Status:           model.USiteStatFree,
		LAWAPortRate:     reqData.LAWAPortRate,
		Remark:           reqData.Remark,
		Creator:          reqData.LoginUser.LoginName,
	}
	uSite.ID = reqData.ID

	if originUSite != nil {
		uSite.Status = originUSite.Status
	}

	_, err = repo.SaveServerUSite(&uSite)
	if reqData.ID == 0 {
		reqData.ID = uSite.ID
	}
	return err
}

// GetServerUSiteByID 根据ID查询机位信息
func GetServerUSiteByID(repo model.Repo, id uint) (resp *ServerUSiteResp, err error) {
	uSite, err := repo.GetServerUSiteByID(id)
	if err != nil {
		return nil, err
	}

	idc, err := repo.GetIDCByID(uSite.IDCID)
	if err != nil {
		return nil, err
	}

	room, err := repo.GetServerRoomByID(uSite.ServerRoomID)
	if err != nil {
		return nil, err
	}

	cabinet, err := repo.GetServerCabinetByID(uSite.ServerCabinetID)
	if err != nil {
		return nil, err
	}

	resp = &ServerUSiteResp{
		ID:               uSite.ID,
		Number:           uSite.Number,
		Height:           uSite.Height,
		Beginning:        uSite.Beginning,
		Remark:           uSite.Remark,
		UpdatedAt:        uSite.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedAt:        uSite.CreatedAt.Format("2006-01-02 15:04:05"),
		Creator:          uSite.Creator,
		Status:           &uSite.Status,
		LAWAPortRate:     uSite.LAWAPortRate,
		PhysicalArea:     uSite.PhysicalArea,
		OobnetSwitches:   toSwitchInfoArray(uSite.OobnetSwitches),
		IntranetSwitches: toSwitchInfoArray(uSite.IntranetSwitches),
		ExtranetSwitches: toSwitchInfoArray(uSite.ExtranetSwitches),
	}
	resp.IDC.ID = idc.ID
	resp.IDC.Name = idc.Name
	resp.ServerRoom.ID = room.ID
	resp.ServerRoom.Name = room.Name
	resp.ServerCabinet.ID = cabinet.ID
	resp.ServerCabinet.Number = cabinet.Number

	return resp, nil

}

type UsiteTreeResp struct {
	Roots []*RootCabinet `json:"roots"`
}

type RootCabinet struct {
	CabinetID            uint         `json:"cabinet_id"`
	CabinetNumber        string       `json:"cabinet_number"`
	AvailableUsitesCount uint         `json:"available_usites_count"`
	Leaves               []*LeafUsite `json:"leaves"`
}
type LeafUsite struct {
	UsiteID     uint   `json:"usite_id"`
	UsiteNumber string `json:"usite_number"`
	Status      string `json:"status"`
}

type UsiteTreeReq struct {
	IDCID        uint   `json:"idc_id"`
	ServerRoomID uint   `json:"server_room_id"`
	PhysicalArea string `json:"physical_area"`
	UsiteStatus  string `json:"usite_status"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *UsiteTreeReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDCID:        "idc_id",
		&reqData.ServerRoomID: "server_room_id",
		&reqData.PhysicalArea: "physical_area",
		&reqData.UsiteStatus:  "usite_status",
	}
}

// Validate 结构体数据校验
func (reqData *UsiteTreeReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if reqData.IDCID == 0 {
		errs.Add([]string{"idc_id"}, binding.RequiredError, "需要指定数据中心ID")
		return errs
	}
	if reqData.ServerRoomID == 0 {
		errs.Add([]string{"status"}, binding.RequiredError, "需要指定机房管理单元ID")
		return errs
	}
	if reqData.PhysicalArea == "" {
		errs.Add([]string{"physical_area"}, binding.RequiredError, "需要指定物理区域")
		return errs
	}
	return errs
}

// GetUsiteTree 机位树
func GetUsiteTree(log logger.Logger, repo model.Repo, reqData *UsiteTreeReq) (tree *UsiteTreeResp, err error) {
	//查询指定数据中心，机房下的机架信息
	cabinets, err := repo.GetCabinetOrderByFreeUsites(&model.ServerCabinet{
		IDCID:        reqData.IDCID,
		ServerRoomID: reqData.ServerRoomID,
	}, reqData.PhysicalArea)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	tree = &UsiteTreeResp{}
	tree.Roots = make([]*RootCabinet, 0, len(cabinets))

	for _, c := range cabinets {
		//获取每个机架下的机位信息
		cond := &model.CombinedServerUSite{
			ServerCabinetID: []uint{c.ID},
			PhysicalArea:    reqData.PhysicalArea,
		}
		if reqData.UsiteStatus != "" {
			cond.Status = reqData.UsiteStatus
		}
		usites, err := repo.GetServerUSiteByCond(cond, model.OneOrderBy("number", model.ASC), nil)
		if err != nil {
			log.Error(err)
			continue
		}
		leaf := make([]*LeafUsite, 0, len(usites))
		for _, u := range usites {
			leaf = append(leaf, &LeafUsite{
				UsiteID:     u.ID,
				UsiteNumber: u.Number,
				Status:      u.Status,
			})
		}

		tree.Roots = append(tree.Roots, &RootCabinet{
			CabinetID:            c.ID,
			CabinetNumber:        c.Number,
			AvailableUsitesCount: c.AvailableUsiteCount,
			Leaves:               leaf,
		})
	}
	return
}

// BatchUpdateServerUSitesStatus 批量更新机位状态信息
func BatchUpdateServerUSitesStatus(repo model.Repo, ids []uint, status string) (affected int64, err error) {
	return repo.BatchUpdateServerUSitesStatus(ids, status)
}

// BatchUpdateServerUSitesStatusByCond 根据其他定制性条件进行批量更新机位状态信息
func BatchUpdateServerUSitesStatusByCond(log logger.Logger, repo model.Repo, susfu []*ServerUSiteForUpdate, status string) (affected int64, err error) {
	// 根据 机房管理单元-机架-机位 更新状态
	if len(susfu) !=0 {
		// 拼接 ServerUSiteForUpdate 成 机房管理单元-机架-机位
		var usiteslice []string
		for _, u := range susfu {
			// 根据条件查询获取机位
			condcheck := model.CombinedServerUSite{
				ServerRoomNameCabinetNumUSiteNum: fmt.Sprintf("%s%s%s", u.ServerRoomName, u.ServerCabinetNumber, u.ServerUsiteNumber),
			}
			usite, err := repo.GetServerUSiteByCond(&condcheck, nil, nil)
			if err != nil {
				log.Errorf(err.Error())
				return 0, err
			}
			if len(usite) == 0 {
				log.Errorf("机位不存在：%s", condcheck.ServerRoomNameCabinetNumUSiteNum)
				return 0, fmt.Errorf("机位不存在：%s", condcheck.ServerRoomNameCabinetNumUSiteNum)
			}
			if len(usite) > 1 {
				log.Errorf("机位不唯一：%s", condcheck.ServerRoomNameCabinetNumUSiteNum)
				return 0, fmt.Errorf("机位不唯一：%s", condcheck.ServerRoomNameCabinetNumUSiteNum)
			}
			usiteslice = append(usiteslice, fmt.Sprintf("%s%s%s", u.ServerRoomName, u.ServerCabinetNumber, u.ServerUsiteNumber))
		}
		// 根据条件查询获取机位
		cond := model.CombinedServerUSite{
			ServerRoomNameCabinetNumUSiteNumSlice: usiteslice,
		}
		usites, err := repo.GetServerUSiteByCond(&cond, nil, nil)
		if err != nil {
			log.Errorf(err.Error())
			return 0, err
		}
		// 查询的机位数量必须与参数一致
		if len(usites) == len(susfu) {
			var ids []uint
			switch status {
			// 空闲 -> 预占用
			case model.USiteStatPreOccupied:
				for k := range usites {
					if usites[k].Status != model.USiteStatFree {
						serverRoom, _ := repo.GetServerRoomByID(usites[k].ServerRoomID)
						cabinet, _ := repo.GetServerCabinetByID(usites[k].ServerCabinetID)
						log.Errorf("存在非空闲的机位（%v-%v-%v），不允许置为预占用", serverRoom.Name, cabinet.Number, usites[k].Number)
						return 0, fmt.Errorf("存在非空闲的机位（%v-%v-%v），不允许置为预占用", serverRoom.Name, cabinet.Number, usites[k].Number)
					}
					ids = append(ids, usites[k].ID)
				}
			// 预占用 -> 空闲
			case model.USiteStatFree:
				for k := range usites {
					if usites[k].Status != model.USiteStatPreOccupied {
						serverRoom, _ := repo.GetServerRoomByID(usites[k].ServerRoomID)
						cabinet, _ := repo.GetServerCabinetByID(usites[k].ServerCabinetID)
						log.Errorf("存在非预占用的机位（%v-%v-%v），不允许置为空闲", serverRoom.Name, cabinet.Number, usites[k].Number)
						return 0, fmt.Errorf("存在非预占用的机位（%v-%v-%v），不允许置为空闲", serverRoom.Name, cabinet.Number, usites[k].Number)
					}
					ids = append(ids, usites[k].ID)
				}
			}
			// 校验通过，通过ID批量更新机位状态
			return repo.BatchUpdateServerUSitesStatus(ids, status)
		} else {
			log.Errorf("查询获取的机位数量(%v)与请求传入的数量(%v)不一致",len(usites), len(susfu))
			return 0, fmt.Errorf("查询获取的机位数量(%v)与请求传入的数量(%v)不一致",len(usites), len(susfu))
		}
	}
	log.Errorf("请求传入的机位数量为0，未更新任何机位状态")
	return 0, fmt.Errorf("请求传入的机位数量为0，未更新任何机位状态")
}

// DeleteServerUSitePort 删除机位端口号
func DeleteServerUSitePort(repo model.Repo, id uint) (affected int64, err error) {

	//USite check
	usite, err := repo.GetServerUSiteByID(id)
	if err != nil {
		return 0, err
	}

	if usite == nil {
		return 0, fmt.Errorf("机位信息不存在")
	}

	// check device
	cond := &model.CombinedDeviceCond{
		USiteID: []uint{usite.ID},
	}
	devices, err := repo.GetCombinedDevices(cond, nil, nil)
	if err != nil {
		return 0, err
	}

	if len(devices) > 0 {
		return 0, fmt.Errorf("指定编号(%s)的机位下存在设备信息，无法删除", usite.Number)
	}

	return repo.DeleteServerUSitePort(id)
}

// RemoveServerUSiteByID 删除机位
func RemoveServerUSiteByID(repo model.Repo, id uint) (affected int64, err error) {
	// 校验机位是否被使用
	usite, err := repo.GetServerUSiteByID(id)
	if err != nil {
		return 0, err
	}

	if model.USiteStatUsed == usite.Status {
		return 0, fmt.Errorf("机位(%s)已经被使用，无法删除", usite.Number)
	}

	items, err := repo.GetDevicesByUSiteID(id)
	if err != nil {
		return 0, err
	}

	if len(items) > 0 {
		return 0, fmt.Errorf("机位(%s)已经分配设备信息，无法删除", usite.Number)
	}

	return repo.RemoveServerUSiteByID(id)
}

//ServerUSiteForImport 为了导入
type ServerUSiteForImport struct {
	ServerCabinetNumber string `json:"server_cabinet_number"`
	ServerRoomName      string `json:"server_room_name"`
	Number              string `json:"number"`
	Height              int    `json:"height"`
	PhysicalArea        string `json:"physical_area"`
	Beginning           int    `json:"beginning"`
	Remark              string `json:"remark"`
	Status              string `json:"status"`
	IDCID               uint   `json:"idc_id"`
	ServerRoomID        uint   `json:"server_room_id"`
	ServerCabinetID     uint   `json:"server_cabinet_id"`
	Creator             string `json:"creator"`
	Content             string `json:"content"`
	ID uint `json:"-"`
}

//ServerUSitePortForImport 为了导入
type ServerUSitePortForImport struct {
	ID                     uint   `json:"id"`
	ServerCabinetNumber    string `json:"server_cabinet_number"`
	ServerRoomName         string `json:"server_room_name"`
	ServerRoomID           uint   `json:"server_room_id"`
	ServerCabinetID        uint   `json:"server_cabinet_id"`
	IDCID                  uint   `json:"idc_id"`
	Number                 string `json:"number"`
	OobnetSwitchNamePort   string `json:"oobnet_switch_name_port"`
	IntranetNamePort       string `json:"intranet_switch_name_port"`
	ExtranetSwitchNamePort string `json:"extranet_switch_name_port"`
	LAWAPortRate   		   string `json:"la_wa_port_rate"`		// 内外网端口速率：GE\10GE\25GE\40GE
	Content                string `json:"content"`
}

//checkLength 对导入文件中的数据做基本验证
func (susfi *ServerUSiteForImport) checkLength() {
	leg := len(susfi.ServerRoomName)
	if leg == 0 || leg > 255 {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("必填项校验:机房管理单元名称length(%d)", leg)
	}
	leg = len(susfi.ServerCabinetNumber)
	if leg == 0 || leg > 255 {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("必填项校验:机架编号length(%d)", leg)
	}
	leg = len(susfi.Number)
	if leg == 0 || leg > 255 {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("必填项校验:机位编号length(%d)", leg)
	}
	if susfi.Height <= 0 {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("必填项校验:机位高度(%d)", susfi.Height)
	}
	if susfi.Beginning <= 0 {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("必填项校验:起始U数(%d)", susfi.Beginning)
	}
	leg = len(susfi.PhysicalArea)
	if leg == 0 || leg > 255 {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("必填项校验: 物理区域length(%d)", leg)
	}
}

//validate 对导入文件中的数据做基本验证
func (susfi *ServerUSiteForImport) validate(repo model.Repo) error {
	//机房校验
	srs, err := repo.GetServerRoomByName(susfi.ServerRoomName)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound || srs == nil {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("机房名(%s)不存在", susfi.ServerRoomName)
		return nil
	}

	susfi.IDCID = srs.IDCID
	susfi.ServerRoomID = srs.ID

	//机架(柜)校验
	cabinets, err := repo.GetServerCabinetID(&model.ServerCabinet{Number: susfi.ServerCabinetNumber,
		ServerRoomID: susfi.ServerRoomID,
		IDCID:        susfi.IDCID,
	})
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound || len(cabinets) != 1 {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("机架(柜)(%s)不存在", susfi.ServerCabinetNumber)
		return nil
	}
	susfi.ServerCabinetID = cabinets[0]

	cabinet, _ := repo.GetServerCabinetByID(susfi.ServerCabinetID)

	// 校验机位高度不能超过机架高度
	if cabinet.Height < uint(susfi.Height) {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("机位高度 %d, 不能超过所在机架高度 %d", susfi.Height, cabinet.Height)
		return nil
	}
	// 校验机位开始U位不能超过机架高度
	if cabinet.Height < uint(susfi.Beginning) {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("机位开始U位 %d, 不能超过所在机架高度 %d", susfi.Beginning, cabinet.Height)
		return nil
	}
	// 校验机位开始U位+机位调试不能超过机架高度
	if cabinet.Height < uint(susfi.Beginning+susfi.Height) {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("机位开始U位 %d, U位调试 %d, 不能超过所在机架高度 %d", susfi.Beginning, susfi.Height, cabinet.Height)
		return nil
	}

	if netArea, err := repo.GetNetworkAreaByID(cabinet.NetworkAreaID); err != nil || netArea == nil {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("机架：%s网络区域ID：%d不存在", cabinet.Number, cabinet.NetworkAreaID)
		return nil
	} else if netArea != nil {
		var physicalArea = make([]model.ParamListINT, 0)
		err = json.Unmarshal([]byte(netArea.PhysicalArea), &physicalArea)
		if err != nil {
			return err
		}
		physicalAreaValid := false
		for _, phyArea := range physicalArea {
			if phyArea.Name == susfi.PhysicalArea {
				physicalAreaValid = true
				break
			}
		}
		if !physicalAreaValid {
			var br string
			if susfi.Content != "" {
				br = "<br />"
			}
			susfi.Content += br + fmt.Sprintf("物理区域:%s非法,不属于网络区域:%s", susfi.PhysicalArea, netArea.Name)
			return nil
		}
	}

	// 如果新增，做唯一性校验
	cond := model.CombinedServerUSite{
		IDCID:           []uint{susfi.IDCID},
		USiteNumber:     susfi.Number,
		ServerRoomID:    []uint{susfi.ServerRoomID},
		ServerCabinetID: []uint{susfi.ServerCabinetID},
	}
	usites, err := repo.GetServerUSiteByCond(&cond, nil, nil)
	if err != nil {
		return err
	}
	if len(usites) != 0 {
		for k := range usites {
			if usites[k].Number == susfi.Number {
				susfi.ID = usites[k].ID
				break
			}
		}
	} else {
		susfi.ID = 0
	}

	//逻辑区域校验： 待定
	//TODO
	return nil
}

//checkLength 对导入文件中的数据做基本验证
func (susfi *ServerUSitePortForImport) checkLength() {
	leg := len(susfi.ServerRoomName)
	if leg == 0 || leg > 255 {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("必填项校验:机房管理单元名称length(%d)", leg)
	}
	leg = len(susfi.ServerCabinetNumber)
	if leg == 0 || leg > 255 {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("必填项校验:机架编号length(%d)", leg)
	}
	leg = len(susfi.Number)
	if leg == 0 || leg > 255 {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("必填项校验:机位编号length(%d)", leg)
	}
	leg = len(susfi.OobnetSwitchNamePort)
	if leg == 0 || leg > 255 {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("必填项校验:带外交换机设备端口(%s)", susfi.OobnetSwitchNamePort)
	}
	leg = len(susfi.IntranetNamePort)
	if leg == 0 || leg > 255 {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("必填项校验:内网交换机设备端口(%s)", susfi.IntranetNamePort)
	}
	leg = len(susfi.LAWAPortRate)
	if leg == 0 || leg > 255 {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("必填项校验:内外网端口速率(%s)", susfi.IntranetNamePort)
	}
}

//validate 对导入文件中的数据做基本验证
func (susfi *ServerUSitePortForImport) validate(repo model.Repo) error {
	// 添加机位端口名称格式校验
	if susfi.ExtranetSwitchNamePort != "" && !strings.Contains(susfi.ExtranetSwitchNamePort, "_") {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("外网交换机名称端口号(%s)格式不正确，多个交换机用英文分号(;)分隔，设备名称与端口用下划线(_)分隔", susfi.ExtranetSwitchNamePort)
	}

	if !strings.Contains(susfi.IntranetNamePort, "_") {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("内网交换机名称端口号(%s)格式不正确，多个交换机用英文分号(;)分隔，设备名称与端口用下划线(_)分隔", susfi.IntranetNamePort)
	}

	if !strings.Contains(susfi.OobnetSwitchNamePort, "_") {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("管理网交换机名称端口号(%s)格式不正确，多个交换机用英文分号(;)分隔，设备名称与端口用下划线(_)分隔", susfi.OobnetSwitchNamePort)
	}
	// 内外网端口速率校验：GE\10GE\25GE\40GE
	if _, ok := model.PortRateMap[susfi.LAWAPortRate]; !ok {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("内外网端口速率（%s）不正确，应为：GE|10GE|25GE|40GE", susfi.LAWAPortRate)
	}
	//机房校验
	srs, err := repo.GetServerRoomByName(susfi.ServerRoomName)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound || srs == nil {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("机房名(%s)不存在", susfi.ServerRoomName)
		return nil
	}

	susfi.IDCID = srs.IDCID
	susfi.ServerRoomID = srs.ID

	//机架(柜)校验
	scid, err := repo.GetServerCabinetID(&model.ServerCabinet{Number: susfi.ServerCabinetNumber,
		ServerRoomID: susfi.ServerRoomID,
		IDCID:        susfi.IDCID,
	})
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound || len(scid) != 1 {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("机房%s 中不存在 机架(柜)(%s)", srs.Name, susfi.ServerCabinetNumber)
		return nil
	}
	susfi.ServerCabinetID = scid[0]

	// 机位校验
	usite, err := repo.GetServerUSiteByNumber(susfi.ServerCabinetID, susfi.Number)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("机位(U位)(%s)不存在", susfi.Number)
		return nil
	}

	if usite.ServerCabinetID != susfi.ServerCabinetID {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("[%s]机位的机架与[%s]机架不一致", usite.Number, susfi.ServerCabinetNumber)
		return nil
	}
	if usite.ServerRoomID != susfi.ServerRoomID {
		var br string
		if susfi.Content != "" {
			br = "<br />"
		}
		susfi.Content += br + fmt.Sprintf("[%s]机位的机房与[%s]机房不一致", usite.Number, susfi.ServerRoomName)
		return nil
	}
	return nil
}

//ImportServerUSitePreview 导入预览
func ImportServerUSitePreview(log logger.Logger, repo model.Repo, reqData *upload.ImportReq) (map[string]interface{}, error) {
	ra, err := utils.ParseDataFromXLSX(upload.UploadDir + reqData.FileName)
	if err != nil {
		return nil, err
	}
	length := len(ra)

	var success []*ServerUSiteForImport
	var failure []*ServerUSiteForImport
	for i := 1; i < length; i++ {
		row := &ServerUSiteForImport{}
		if len(ra[i]) < 8 {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "导入文件列长度不对（应为8列）"
			failure = append(failure, row)
			continue
		}

		row.ServerRoomName = strings.TrimSpace(ra[i][0])
		row.ServerCabinetNumber = strings.TrimSpace(ra[i][1])
		row.Number = strings.TrimSpace(ra[i][2])
		row.Height, _ = strconv.Atoi(strings.TrimSpace(ra[i][3]))
		row.PhysicalArea = strings.TrimSpace(ra[i][4])
		row.Beginning, _ = strconv.Atoi(strings.TrimSpace(ra[i][5]))
		row.Status = strings.TrimSpace(ra[i][6])
		row.Remark = strings.TrimSpace(ra[i][7])

		//必填项校验
		row.checkLength()
		//机房和机架(柜)校验
		err := row.validate(repo)
		if err != nil {
			return nil, err
		}

		if row.Content != "" {
			failure = append(failure, row)
		} else {
			success = append(success, row)
		}
	}

	var data []*ServerUSiteForImport
	if len(failure) > 0 {
		data = failure
	} else {
		data = success
	}
	var result []*ServerUSiteForImport
	for i := 0; i < len(data); i++ {
		if uint(i) >= reqData.Offset && uint(i) < (reqData.Offset+reqData.Limit) {
			result = append(result, data[i])
		}
	}
	if len(failure) > 0 {
		return map[string]interface{}{"status": "failure",
			"message":      "导入机位(U位)错误",
			"record_count": len(data),
			"content":      result,
		}, nil
	}
	return map[string]interface{}{"status": "success",
		"message":       "操作成功",
		"import_status": true,
		"record_count":  len(data),
		"content":       result,
	}, nil
}

//ImportServerUSitePortsPreview 导入预览
func ImportServerUSitePortsPreview(log logger.Logger, repo model.Repo, reqData *upload.ImportReq) (map[string]interface{}, error) {
	ra, err := utils.ParseDataFromXLSX(upload.UploadDir + reqData.FileName)
	if err != nil {
		return nil, err
	}
	length := len(ra)

	var success []*ServerUSitePortForImport
	var failure []*ServerUSitePortForImport
	for i := 1; i < length; i++ {
		row := &ServerUSitePortForImport{}
		if len(ra[i]) < 7 {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "导入文件列长度不对（应为7列）"
			failure = append(failure, row)
			continue
		}

		row.ServerRoomName = strings.TrimSpace(ra[i][0])
		row.ServerCabinetNumber = strings.TrimSpace(ra[i][1])
		row.Number = strings.TrimSpace(ra[i][2])
		row.OobnetSwitchNamePort = strings.TrimSpace(ra[i][3])
		row.IntranetNamePort = strings.TrimSpace(ra[i][4])
		row.ExtranetSwitchNamePort = strings.TrimSpace(ra[i][5])
		row.LAWAPortRate = strings.TrimSpace(ra[i][6])

		//中文分号替换
		row.OobnetSwitchNamePort = strings.Replace(row.OobnetSwitchNamePort, "；", ";", -1)
		row.IntranetNamePort = strings.Replace(row.IntranetNamePort, "；", ";", -1)
		row.ExtranetSwitchNamePort = strings.Replace(row.ExtranetSwitchNamePort, "；", ";", -1)

		//必填项校验
		row.checkLength()
		//机房和机架(柜)校验
		err := row.validate(repo)
		if err != nil {
			return nil, err
		}

		if row.Content != "" {
			failure = append(failure, row)
		} else {
			success = append(success, row)
		}
	}

	var data []*ServerUSitePortForImport
	if len(failure) > 0 {
		data = failure
	} else {
		data = success
	}
	var result []*ServerUSitePortForImport
	for i := 0; i < len(data); i++ {
		if uint(i) >= reqData.Offset && uint(i) < (reqData.Offset+reqData.Limit) {
			result = append(result, data[i])
		}
	}
	if len(failure) > 0 {
		return map[string]interface{}{"status": "failure",
			"message":      "导入机位(U位)端口错误",
			"record_count": len(data),
			"content":      result,
		}, nil
	}
	return map[string]interface{}{"status": "success",
		"message":       "操作成功",
		"import_status": true,
		"record_count":  len(data),
		"content":       result,
	}, nil
}

//ImportServerUSite 将导入机位(U位)放到数据库
func ImportServerUSite(log logger.Logger, repo model.Repo, reqData *upload.ImportReq) error {
	ra, err := utils.ParseDataFromXLSX(upload.UploadDir + reqData.FileName)
	if err != nil {
		return err
	}
	length := len(ra)
	for i := 1; i < length; i++ {
		row := &ServerUSiteForImport{}
		if len(ra[i]) < 8 {
			continue
		}

		row.ServerRoomName = strings.TrimSpace(ra[i][0])
		row.ServerCabinetNumber = strings.TrimSpace(ra[i][1])
		row.Number = strings.TrimSpace(ra[i][2])
		row.Height, _ = strconv.Atoi(strings.TrimSpace(ra[i][3]))
		row.PhysicalArea = strings.TrimSpace(ra[i][4])
		row.Beginning, _ = strconv.Atoi(strings.TrimSpace(ra[i][5]))
		row.Status = usiteStatusTransfer(ra[i][6], false)
		row.Remark = strings.TrimSpace(ra[i][7])

		//必填项校验
		row.checkLength()
		//机房和机架(柜)校验
		err := row.validate(repo)
		if err != nil {
			return err
		}

		sus := &model.ServerUSite{
			IDCID:           row.IDCID,
			ServerRoomID:    row.ServerRoomID,
			ServerCabinetID: row.ServerCabinetID,
			Height:          uint(row.Height),
			Beginning:       uint(row.Beginning),
			PhysicalArea:    row.PhysicalArea,
			Number:          row.Number,
			Remark:          row.Remark,
			LAWAPortRate:    model.PortRateDefault,
			Status:          model.USiteStatFree,
			Creator:         reqData.UserName, //row.Creator,
		}

		sus.ID = row.ID

		if row.Status != "" {
			sus.Status = row.Status
		}

		//插入或者更新
		if _, err = repo.SaveServerUSite(sus); err != nil {
			return err
		}
	}
	defer os.Remove(upload.UploadDir + reqData.FileName)
	return nil
}

//ImportServerUSitePort 将导入机位(U位)放到数据库
func ImportServerUSitePort(log logger.Logger, repo model.Repo, reqData *upload.ImportReq) error {
	ra, err := utils.ParseDataFromXLSX(upload.UploadDir + reqData.FileName)
	if err != nil {
		return err
	}
	length := len(ra)
	for i := 1; i < length; i++ {
		row := &ServerUSitePortForImport{}
		if len(ra[i]) < 7 {
			continue
		}

		row.ServerRoomName = strings.TrimSpace(ra[i][0])
		row.ServerCabinetNumber = strings.TrimSpace(ra[i][1])
		row.Number = strings.TrimSpace(ra[i][2])
		row.OobnetSwitchNamePort = strings.TrimSpace(ra[i][3])
		row.IntranetNamePort = strings.TrimSpace(ra[i][4])
		row.ExtranetSwitchNamePort = strings.TrimSpace(ra[i][5])
		row.LAWAPortRate = strings.TrimSpace(ra[i][6])

		//必填项校验
		row.checkLength()
		//机房和机架(柜)校验
		err := row.validate(repo)
		if err != nil {
			return err
		}
		
		//中文分号替换
		row.OobnetSwitchNamePort = strings.Replace(row.OobnetSwitchNamePort, "；", ";", -1)
		row.IntranetNamePort = strings.Replace(row.IntranetNamePort, "；", ";", -1)
		row.ExtranetSwitchNamePort = strings.Replace(row.ExtranetSwitchNamePort, "；", ";", -1)		

		// 处理机位端口号
		OobnetSwitches := processSwitchNamePort(row.OobnetSwitchNamePort)
		IntranetSwitches := processSwitchNamePort(row.IntranetNamePort)
		ExtranetSwitches := processSwitchNamePort(row.ExtranetSwitchNamePort)

		sus := &model.ServerUSite{
			IDCID:            row.IDCID,
			ServerRoomID:     row.ServerRoomID,
			ServerCabinetID:  row.ServerCabinetID,
			OobnetSwitches:   toJSONString(OobnetSwitches),
			IntranetSwitches: toJSONString(IntranetSwitches),
			ExtranetSwitches: toJSONString(ExtranetSwitches),
			LAWAPortRate:	  row.LAWAPortRate,
		}

		//查询是否已经存在
		usite, err := repo.GetServerUSiteByNumber(row.ServerCabinetID, row.Number)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		sus.ID = usite.ID
		sus.Number = usite.Number
		sus.PhysicalArea = usite.PhysicalArea
		sus.Height = usite.Height
		sus.Status = usite.Status
		sus.Remark = usite.Remark
		sus.Beginning = usite.Beginning

		//插入或者更新
		if _, err = repo.SaveServerUSite(sus); err != nil {
			return err
		}
	}
	defer os.Remove(upload.UploadDir + reqData.FileName)
	return nil
}

// processSwitchNamePort 处理机位端口号
func processSwitchNamePort(namePort string) []*model.SwitchInfo {
	var OobnetSwitches []*model.SwitchInfo

	if strings.Contains(namePort, "_") && strings.Contains(namePort, ";") {
		for _, oobnetSwitch := range strings.Split(namePort, ";") {
			switchInfo := &model.SwitchInfo{
				Name: oobnetSwitch[:strings.LastIndex(oobnetSwitch, "_")],
				Port: oobnetSwitch[strings.LastIndex(oobnetSwitch, "_")+1:],
			}
			OobnetSwitches = append(OobnetSwitches, switchInfo)
		}
	}

	if strings.Contains(namePort, "_") && !strings.Contains(namePort, ";") {
		switchInfo := &model.SwitchInfo{
			Name: namePort[:strings.LastIndex(namePort, "_")],
			Port: namePort[strings.LastIndex(namePort, "_")+1:],
		}
		OobnetSwitches = append(OobnetSwitches, switchInfo)
	}

	return OobnetSwitches
}

//CheckUSiteFree 检查机位是否空闲
func CheckUSiteFree(repo model.Repo, uSiteID uint, dev *model.Device) bool {
	//如果设备的机位和没有变化，则不做检验，返回true
	if dev != nil && dev.USiteID != nil && *dev.USiteID == uSiteID {
		return true
	}
	uSite, err := repo.GetServerUSiteByID(uSiteID)
	if err != nil {
		return false
	}
	return uSite.Status == model.USiteStatFree || uSite.Status == model.USiteStatPreOccupied
}

//toJSONString 转换json为字符串
func toJSONString(arg interface{}) string {
	bytes, _ := json.Marshal(arg)
	return string(bytes)
}

//toSwitchInfoArray 转换jsonArray为交换机信息数组
func toSwitchInfoArray(str string) []*model.SwitchInfo {
	var items []*model.SwitchInfo
	json.Unmarshal([]byte(str), &items)
	return items
}

//usiteStatusTransfer 运行状态值和数据库存储值的转换
func usiteStatusTransfer(status string, reverse bool) string {
	mStatus := map[string]string{
		"空闲":  "free",
		"预占用": "pre_occupied",
		"不可用": "disabled",
		"已使用": "used",
	}
	if !reverse {
		if val, ok := mStatus[status]; ok {
			return val
		}
	} else {
		for key, val := range mStatus {
			if val == status {
				return key
			}
		}
	}
	return ""
}

// PhysicalAreaConnd 物理区域搜索条件
type PhysicalAreaConnd struct {
	NetworkAreaID uint `json:"network_area_id"`
	//Name          string `json:"name"`
}

// FieldMap 请求字段映射
func (reqData *PhysicalAreaConnd) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.NetworkAreaID: "network_area_id",
		//&reqData.Name:          "name",
	}
}

// Validate 结构体数据校验
func (reqData *PhysicalAreaConnd) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	return nil
}

//GetPhysicalAreas 物理区域列表
func GetPhysicalAreas(log logger.Logger, repo model.Repo, conn PhysicalAreaConnd) (*model.DeviceQueryParamResp, error) {

	var physicalArea []model.ParamListINT

	if conn.NetworkAreaID > 0 {
		netArea, err := repo.GetNetworkAreaByID(conn.NetworkAreaID)
		if err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(netArea.PhysicalArea), &physicalArea)
	} else {
		netAreas, err := repo.GetNetworkAreasByCond(nil)
		if err != nil {
			return nil, err
		}

		for i := range netAreas {
			var pArea []model.ParamListINT
			_ = json.Unmarshal([]byte(netAreas[i].PhysicalArea), &pArea)
			physicalArea = append(physicalArea, pArea...)
		}
	}
	resp := &model.DeviceQueryParamResp{
		ParamName: "physical_area",
		//List:      physicalArea,
	}
	for _, l := range physicalArea {
		resp.List = append(resp.List, model.ParamList{
			ID:   strconv.Itoa(int(l.ID)),
			Name: l.Name,
		})
	}
	return resp, nil

}


// BatchUpdateServerUSitesRemarkReq
type BatchUpdateServerUSitesRemarkReq struct {
	// 备注
	Remark      string             		`json:"remark"`
	// id列表
	IDs       	[]uint             		`json:"ids"`
}

// BatchUpdateServerUSitesRemark 批量更新机位备注信息
func BatchUpdateServerUSitesRemark(repo model.Repo, ids []uint, remark string) (affected int64, err error) {
	return repo.BatchUpdateServerUSitesRemark(ids, remark)
}