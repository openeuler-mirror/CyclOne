package service

import (
	"fmt"
	"net/http"

	"github.com/voidint/binding"

	"reflect"

	"encoding/json"
	stringsStd "strings"
	"time"

	"strconv"

	"errors"

	"github.com/jinzhu/gorm"
	"github.com/voidint/page"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils/strings"
	"idcos.io/cloudboot/utils/times"
)

//SaveOrderReq 保存订单请求参数
type SaveOrderReq struct {
	//(required) 所属数据中心ID
	IDCID uint `json:"idc_id"`
	//ServerRoomID 机房管理单元ID
	ServerRoomID uint `json:"server_room_id"`
	//PhysicalArea 物理区域
	PhysicalArea string `json:"physical_area"`
	//订单ID。若id=0，则新增。若id>0，则修改。
	ID uint `json:"id"`
	//Usage 用途
	Usage string `json:"usage"`
	//Catetory 设备类型
	Category string `json:"category"`
	//Amount 数量
	Amount int `json:"amount"`
	//ExpectedArrivalDate 预计到货日期, 格式2006-01-02, 默认订单日期后45天
	ExpectedArrivalDate string `json:"expected_arrival_date"`
	// 预占机架
	//PreOccupiedCabinets string `json:"pre_occupied_cabinets"`
	// 预占机位, 因为机架机位都支持多选，为了方便，机位同时将机架信息带出，JSON形如：{"cabinet_id":1,"cabinet_number":"机架编号1", "usite_id":100,"usite_number":"机位编号100"}
	PreOccupiedUsites string `json:"pre_occupied_usites"`
	// 参考 DeviceLifecycle
	AssetBelongs	 				string		`json:"asset_belongs"`
	Owner			 				string		`json:"owner"`
	IsRental		 				string		`json:"is_rental"`
	MaintenanceServiceProvider		string		`json:"maintenance_service_provider"`
	MaintenanceService				string		`json:"maintenance_service"`
	LogisticsService				string		`json:"logistics_service"`
	MaintenanceServiceDateBegin     string 		`json:"maintenance_service_date_begin"`
	MaintenanceServiceDateEnd       string 		`json:"maintenance_service_date_end"`
	// Remark 备注
	Remark string `json:"remark"`
	// 用户登录名
	LoginName string `json:"-"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *SaveOrderReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDCID:               			"idc_id",
		&reqData.ServerRoomID:        			"server_room_id",
		&reqData.PhysicalArea:        			"physical_area",
		&reqData.ID:                  			"id",
		&reqData.Usage:               			"usage",
		&reqData.Category:            			"category",
		&reqData.Amount:              			"amount",
		&reqData.ExpectedArrivalDate: 			"expected_arrival_date",
		&reqData.PreOccupiedUsites: 			"pre_occupied_usites",
		&reqData.AssetBelongs:					"asset_belongs",
		&reqData.Owner:							"owner",
		&reqData.IsRental:						"is_rental",
		&reqData.MaintenanceServiceProvider:	"maintenance_service_provider",
		&reqData.MaintenanceService:			"maintenance_service",
		&reqData.LogisticsService:				"logistics_service",
		&reqData.MaintenanceServiceDateBegin:	"maintenance_service_date_begin",
		&reqData.MaintenanceServiceDateEnd:		"maintenance_service_date_end",		
		&reqData.Remark:            			"remark",
	}
}

// Validate 结构体数据校验
func (reqData *SaveOrderReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())
	//必须的基本数据不能为空
	if errs = reqData.baseValidate(req, errs); errs != nil {
		return errs
	}

	//更新订单信息，校验指定ID的订单是否存在
	if reqData.ID > 0 {
		if _, err := repo.GetOrderByID(reqData.ID); errs != nil {
			errs.Add([]string{"id"}, binding.RequiredError, fmt.Sprintf("查询订单id(%d)出现错误: %s", reqData.ID, err.Error()))
			return errs
		}
	}

	//校验IDC数据
	if errs = reqData.checkIDCValidate(req, errs); errs != nil {
		return errs
	}

	return errs
}

func (reqData *SaveOrderReq) checkIDCValidate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())

	_, err := repo.GetIDCByID(reqData.IDCID)
	if err != nil {
		errs.Add([]string{"idc"}, binding.RequiredError, fmt.Sprintf("获取数据中心id(%d)出现错误: %s", reqData.IDCID, err.Error()))
		return errs
	}
	sr, err := repo.GetServerRoomByID(reqData.ServerRoomID)
	if err != nil || sr == nil {
		errs.Add([]string{"server_room"}, binding.RequiredError, fmt.Sprintf("获取机房管理单元id(%d)出现错误: %s", reqData.ServerRoomID, err.Error()))
		return errs
	}
	//检查机房下的可用机位是否不小于订单总量
	if freeCount, _ := repo.CountServerUSite(&model.CombinedServerUSite{
		ServerRoomID: []uint{reqData.ServerRoomID},
		Status:       model.USiteStatFree}); freeCount < int64(reqData.Amount) {
		errs.Add([]string{"amount"}, binding.RequiredError, fmt.Sprintf("机房(%s)的可用机位总数(%d)小于订单总量(%d)", sr.Name, freeCount, reqData.Amount))
		return errs
	}
	us := make([]PreOccupiedUsiteData, 0)
	if err := json.Unmarshal([]byte(reqData.PreOccupiedUsites), &us); err != nil {
		errs.Add([]string{"amount"}, binding.BusinessError, fmt.Sprintf("校验机位数量失败:%v", err))
		return errs
	} else if int(len(us)) != reqData.Amount {
		errs.Add([]string{"amount"}, binding.BusinessError, fmt.Sprintf("勾选机位数(%d)不等于订单总数(%d)", len(us), reqData.Amount))
		return errs
	}
	return errs

}

//baseValidate 必要参数不能为空
func (reqData *SaveOrderReq) baseValidate(req *http.Request, errs binding.Errors) binding.Errors {
	if reqData.IDCID == 0 {
		errs.Add([]string{"idc_id"}, binding.RequiredError, "数据中心ID不能为空")
		return errs
	}

	if reqData.ServerRoomID == 0 {
		errs.Add([]string{"name"}, binding.RequiredError, "机房管理单元不能为空")
		return errs
	}

	if reqData.PhysicalArea == "" {
		errs.Add([]string{"physical_area"}, binding.RequiredError, "物理区域不能为空")
		return errs
	}

	if reqData.Usage == "" {
		errs.Add([]string{"usage"}, binding.RequiredError, "用途不能为空")
		return errs
	}

	if reqData.Category == "" {
		errs.Add([]string{"category"}, binding.RequiredError, "设备类型不能为空")
		return errs
	}
	if reqData.Amount == 0 {
		errs.Add([]string{"amount"}, binding.RequiredError, "订单数量不能为空")
		return errs
	}
	if reqData.PreOccupiedUsites == "" {
		errs.Add([]string{"pre_occupied_usites"}, binding.RequiredError, "预占机位不能为空")
		return errs
	}

	return errs
}

type DelOrderReq struct {
	IDs []uint `json:"ids"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *DelOrderReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDs: "ids",
	}
}

// Validate 结构体数据校验
func (reqData *DelOrderReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())
	for _, id := range reqData.IDs {
		if _, err := repo.GetOrderByID(id); err != nil {
			errs.Add([]string{"id"}, binding.RequiredError, fmt.Sprintf("订单id(%d)不存在", id))
			return errs
		}
	}
	return nil
}

type UpdateOrderStatusReq struct {
	ID     uint   `json:"id"`
	Status string `json:"status"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *UpdateOrderStatusReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.ID:     "id",
		&reqData.Status: "status", //"purchasing|partly_arrived|all_arrived|canceled|finished"
	}
}

// Validate 结构体数据校验
func (reqData *UpdateOrderStatusReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())
	if _, err := repo.GetOrderByID(reqData.ID); err != nil {
		errs.Add([]string{"id"}, binding.RequiredError, fmt.Sprintf("订单id(%d)不存在", reqData.ID))
		return errs
	}
	switch reqData.Status {
	case model.OrderStatusPurchasing:
	case model.OrderStatusPartlyArrived:
	case model.OrderStatusAllArrived:
	case model.OrderStatusCanceled:
	case model.OrderStatusfinished:
	default:
		errs.Add([]string{"status"}, binding.RequiredError, fmt.Sprintf("状态值(%s)非法，合法值[%s]",
			"purchasing|partly_arrived|all_arrived|canceled|finished", reqData.Status))
		return errs
	}
	return nil
}

//SaveOrder 保存订单
func SaveOrder(log logger.Logger, repo model.Repo, reqData *SaveOrderReq) error {
	sr := model.Order{
		IDCID:        					reqData.IDCID,
		ServerRoomID: 					reqData.ServerRoomID,
		PhysicalArea: 					reqData.PhysicalArea,
		Usage:    						reqData.Usage,
		Category: 						reqData.Category,
		Amount:   						reqData.Amount,
		AssetBelongs:					reqData.AssetBelongs,
		Owner:							reqData.Owner,
		IsRental:						reqData.IsRental,
		MaintenanceServiceProvider:		reqData.MaintenanceServiceProvider,
		MaintenanceService:				reqData.MaintenanceService,
		LogisticsService:				reqData.LogisticsService,
		PreOccupiedUsites: 				reqData.PreOccupiedUsites,
		Remark:            				reqData.Remark,
		Status:            				model.OrderStatusPurchasing,
		Creator:           				reqData.LoginName,
	}
	if reqData.ID != 0 {
		origin, _ := repo.GetOrderByID(reqData.ID)
		sr.Model = origin.Model
		if origin != nil {
			sr.LeftAmount = origin.LeftAmount + reqData.Amount - origin.Amount
		}
	} else {
		sr.LeftAmount = reqData.Amount //初始状态，未到货数等于订单总数
	}

	if reqData.ExpectedArrivalDate != "" {
		t, err := time.Parse(times.DateLayout, reqData.ExpectedArrivalDate)
		if err != nil {
			log.Errorf("parse time %s err:%v", reqData.ExpectedArrivalDate, err)
			sr.ExpectedArrivalDate = time.Now().AddDate(0, 0, 45)
		} else {
			sr.ExpectedArrivalDate = t
		}
	} else {
		sr.ExpectedArrivalDate = time.Now().AddDate(0, 0, 45)
	}
	if reqData.MaintenanceServiceDateBegin != "" {
		t, err := time.Parse(times.DateLayout, reqData.MaintenanceServiceDateBegin)
		if err != nil {
			log.Errorf("parse time %s err:%v", reqData.MaintenanceServiceDateBegin, err)
			sr.MaintenanceServiceDateBegin = time.Now().AddDate(0, 0, 45)
		} else {
			sr.MaintenanceServiceDateBegin = t
		}
	} else {
		sr.MaintenanceServiceDateBegin = time.Now().AddDate(0, 0, 45)
	}
	if reqData.MaintenanceServiceDateEnd != "" {
		t, err := time.Parse(times.DateLayout, reqData.MaintenanceServiceDateEnd)
		if err != nil {
			log.Errorf("parse time %s err:%v", reqData.MaintenanceServiceDateEnd, err)
			sr.MaintenanceServiceDateEnd = time.Now().AddDate(5, 0, 0)
		} else {
			sr.MaintenanceServiceDateEnd = t
		}
	} else {
		sr.MaintenanceServiceDateEnd = time.Now().AddDate(5, 0, 0)
	}

	sr.Number = GenOrderNumber(repo, reqData.IDCID)

	_, err := repo.SaveOrder(&sr)
	if err != nil {
		return err
	}

	//将机架机位预占用
	us := make([]PreOccupiedUsiteData, 0)
	if err := json.Unmarshal([]byte(reqData.PreOccupiedUsites), &us); err != nil {
		log.Errorf("unmarshal usites(%s) fail，%v", reqData.PreOccupiedUsites, err)
		return err
	}
	ids := make([]uint, len(us))
	for _, u := range us {
		ids = append(ids, u.Value)
	}
	_, err = BatchUpdateServerUSitesStatus(repo, ids, model.USiteStatPreOccupied)

	reqData.ID = sr.Model.ID
	return err
}

//RemoveOrders 删除指定ID的订单
func RemoveOrders(log logger.Logger, repo model.Repo, reqData *DelOrderReq) (affected int64, err error) {
	for _, id := range reqData.IDs {
		//释放预占用机位
		if _, err = ReleasePreOccupiedUsites(log, repo, id); err != nil {
			log.Errorf("release order(id=%d) fail, err:%v", id, err)
			return affected, err
		}
		_, err := repo.RemoveOrderByID(id)
		if err != nil {
			log.Errorf("delete order(id=%d) fail,err:%v", id, err)
			return affected, err
		}

		affected++
	}
	return affected, err
}

// 根据到货数，更新订单的状态
func UpdateOrderByArrival(log logger.Logger, repo model.Repo, orderNum string, arrivalCount int) (err error) {
	if orderNum == "" {
		return errors.New("订单号为空,更新失败")
	}
	order, err := repo.GetOrderByNumber(orderNum)
	if err != nil {
		log.Errorf("get order by number:%s fail,%v", orderNum, err)
		return fmt.Errorf("订单:%s不存在", orderNum)
	}
	switch order.Status {
	case model.OrderStatusCanceled:
		return fmt.Errorf("订单:%s已取消", order.Number)
	case model.OrderStatusfinished:
		return fmt.Errorf("订单:%s已确认完成", order.Number)
	}

	if order != nil {
		delta := order.LeftAmount - arrivalCount
		if delta < 0 {
			log.Errorf("order(number:%s) amount:%d should not be less than arrival count:%d", orderNum, order.Amount, arrivalCount)
			return fmt.Errorf("订单(%s)新增到货数:%d应不大于剩余到货数量数量:%d", orderNum, arrivalCount, order.LeftAmount)
		} else if delta == 0 {
			order.Status = model.OrderStatusAllArrived
		} else {
			order.Status = model.OrderStatusPartlyArrived
		}
		order.LeftAmount = delta
	}
	if _, err := repo.SaveOrder(order); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// 待部署列表删除设备时，根据订单号与数量更新订单的状态
func UpdateOrderByDelete(log logger.Logger, repo model.Repo, orderNum string, deleteCount int) (err error) {
	if orderNum == "" {
		return errors.New("订单号为空,更新失败")
	}
	order, err := repo.GetOrderByNumber(orderNum)
	if err != nil {
		log.Errorf("get order by number:%s fail,%v", orderNum, err)
		return fmt.Errorf("订单:%s不存在", orderNum)
	}
	switch order.Status {
	case model.OrderStatusCanceled:
		return fmt.Errorf("订单:%s已取消", order.Number)
	case model.OrderStatusfinished:
		return fmt.Errorf("订单:%s已确认完成", order.Number)
	}

	if order != nil {
		delta := order.LeftAmount + deleteCount
		if delta > order.Amount {
			log.Errorf("order(number:%s) amount:%d should not be less than (left count + delete count) %d", orderNum, order.Amount, delta)
			return fmt.Errorf("订单(%s)剩余到货数量:%d 与删除设备数之和(%d)不应大于订单设备数(%d)", orderNum, order.LeftAmount, delta, order.Amount)
		} else if delta == order.Amount {
			order.Status = model.OrderStatusPurchasing
		} else {
			order.Status = model.OrderStatusPartlyArrived
		}
		order.LeftAmount = delta
	}
	if _, err := repo.SaveOrder(order); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

//GetOrderPageReq 获取订单分页请求参数
type GetOrderPageReq struct {
	ID           string `json:"id"`
	Number       string `json:"number"`
	PhysicalArea string
	Status       string `json:"status"`
	Usage        string `json:"usage"`
	Page         int64  `json:"page"`
	PageSize     int64  `json:"page_size"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *GetOrderPageReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.ID:           "id",
		&reqData.Number:       "number",
		&reqData.PhysicalArea: "physical_area",
		&reqData.Status:       "status",
		&reqData.Usage:        "usage",
		&reqData.Page:         "page",
		&reqData.PageSize:     "page_size",
	}
}

//OrderResp 订单分页查询信息
type OrderResp struct {
	//数据中心
	IDC *IDCSimplify `json:"idc"`
	//机房管理单元
	ServerRoom *ServerRoomSimplify `json:"server_room"`
	//PhysicalArea 物理区域
	PhysicalArea string `json:"physical_area"`
	//订单ID。
	ID uint `json:"id"`
	//订单编号
	Number string `json:"number"`
	//Usage 用途
	Usage string `json:"usage"`
	//Catetory 设备类型
	Category string `json:"category"`
	//Amount 数量
	Amount int `json:"amount"`
	//LeftAmount未到货数
	LeftAmount int `json:"left_amount"`
	//ExpectedArrivalDate 预计到货日期
	ExpectedArrivalDate string `json:"expected_arrival_date"`
	// 预占机位
	PreOccupiedUsites string `json:"pre_occupied_usites"`
	// 参考 DeviceLifecycle
	AssetBelongs	 				string		`json:"asset_belongs"`
	Owner			 				string		`json:"owner"`
	IsRental		 				string		`json:"is_rental"`
	MaintenanceServiceProvider		string		`json:"maintenance_service_provider"`
	MaintenanceService				string		`json:"maintenance_service"`
	LogisticsService				string		`json:"logistics_service"`
	MaintenanceServiceDateBegin     string 		`json:"maintenance_service_date_begin"`
	MaintenanceServiceDateEnd       string 		`json:"maintenance_service_date_end"`	
	// Remark 备注
	Remark string `json:"remark"`
	// Status 状态
	Status    string        `json:"status"`
	Creator   string        `json:"creator"`
	CreatedAt times.ISOTime `json:"created_at"`
	UpdatedAt times.ISOTime `json:"updated_at"`
}

//GetOrdersPage 获取订单分页
func GetOrdersPage(log logger.Logger, repo model.Repo, reqData *GetOrderPageReq) (*page.Page, error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 100 {
		reqData.PageSize = 20
	}
	if reqData.Page < 0 {
		reqData.Page = 0
	}

	cond := model.OrderCond{
		ID: strings.Multi2UintSlice(reqData.ID),
	}
	cond.Order.Number = reqData.Number
	cond.Order.PhysicalArea = reqData.PhysicalArea
	cond.Order.Usage = reqData.Usage
	cond.Order.Status = reqData.Status

	totalRecords, err := repo.CountOrders(&cond)
	if err != nil {
		return nil, err
	}

	pager := page.NewPager(reflect.TypeOf(&OrderResp{}), reqData.Page, reqData.PageSize, totalRecords)
	items, err := repo.GetOrders(&cond, model.OneOrderBy("id", model.DESC), pager.BuildLimiter())
	if err != nil {
		return nil, err
	}
	for i := range items {
		item, err := convert2OrderResult(log, repo, items[i])
		if err != nil {
			return nil, err
		}
		if item != nil {
			pager.AddRecords(item)
		}
	}
	return pager.BuildPage(), nil
}

//GetExportOrders
func GetExportOrders(log logger.Logger, repo model.Repo, reqData *GetOrderPageReq) ([]*OrderResp, error) {
	rst := make([]*OrderResp, 0)
	cond := model.OrderCond{
		ID: strings.Multi2UintSlice(reqData.ID),
	}
	cond.Order.Number = reqData.Number
	cond.Order.Usage = reqData.Usage
	cond.Order.Status = reqData.Status

	items, err := repo.GetOrders(&cond, model.OneOrderBy("id", model.DESC), nil)
	if err != nil {
		return nil, err
	}
	for i := range items {
		item, err := convert2OrderResult(log, repo, items[i])
		if err != nil {
			return nil, err
		}
		if item != nil {
			rst = append(rst, item)
		}
	}
	return rst, nil
}

func convert2OrderResult(log logger.Logger, repo model.Repo, o *model.Order) (*OrderResp, error) {
	if o == nil {
		return nil, nil
	}
	result := OrderResp{
		IDC:                 			&IDCSimplify{ID: o.IDCID},
		ServerRoom:          			&ServerRoomSimplify{ID: o.ServerRoomID},
		ID:                  			o.Model.ID,
		CreatedAt:           			times.ISOTime(o.CreatedAt),
		UpdatedAt:           			times.ISOTime(o.UpdatedAt),
		PhysicalArea:        			o.PhysicalArea,
		Number:              			o.Number,
		Usage:               			o.Usage,
		Category:            			o.Category,
		Amount:              			o.Amount,
		LeftAmount:          			o.LeftAmount,
		ExpectedArrivalDate: 			o.ExpectedArrivalDate.Format(times.DateLayout),
		AssetBelongs:		 			o.AssetBelongs,	 			
		Owner:				 			o.Owner,
		IsRental:			 			o.IsRental,
		MaintenanceServiceProvider:		o.MaintenanceServiceProvider,
		MaintenanceService:				o.MaintenanceService,
		LogisticsService:				o.LogisticsService,
		MaintenanceServiceDateBegin:	o.MaintenanceServiceDateBegin.Format(times.DateLayout),
		MaintenanceServiceDateEnd:		o.MaintenanceServiceDateEnd.Format(times.DateLayout),
		Remark:  						o.Remark,
		Status:  						o.Status,
		Creator: 						o.Creator,
	}
	result.PreOccupiedUsites, _ = formatUsite(repo, o.PreOccupiedUsites, log)

	idc, err := repo.GetIDCByID(o.IDCID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if idc != nil {
		result.IDC.Name = idc.Name
	}

	sroom, err := repo.GetServerRoomByID(o.ServerRoomID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if sroom != nil {
		result.ServerRoom.Name = sroom.Name
	}

	return &result, nil
}

//GetOrderByID 获取指定ID的订单的详细信息
func GetOrderByID(log logger.Logger, repo model.Repo, id uint) (*OrderResp, error) {
	items, err := repo.GetOrderByID(id)
	if err != nil {
		return nil, err
	}

	item, err := convert2OrderResult(log, repo, items)
	if err != nil {
		return nil, err
	}
	return item, nil
}

//UpdateOrderStatus 获取指定ID的订单的详细信息
func UpdateOrderStatus(log logger.Logger, repo model.Repo, reqData *UpdateOrderStatusReq) error {
	mod := model.Order{
		Status: reqData.Status,
	}
	mod.Model.ID = reqData.ID

	_, err := repo.UpdateOrder(&mod)
	if err != nil {
		return err
	}

	switch reqData.Status {
	//case model.OrderStatusPurchasing: 这个是订单生成的默认状态：采购中
	case model.OrderStatusPartlyArrived:
		//当按订单导入时，但到货数量小于订单数量，此时为部分到货状态。
	case model.OrderStatusAllArrived:
		//当按订单导入时，但到货数量等于订单数量，此时为全部到货状态。
	case model.OrderStatusCanceled:
		//订单被取消
		fallthrough
	case model.OrderStatusfinished:
		//订单已确认
		_, err = ReleasePreOccupiedUsites(log, repo, reqData.ID)
	}
	return err
}

var GlobalDailyCount = make(map[string]int, 1) //每天滚动计数，只缓存一天的数据

//生成订单号的，简单实现版本
func GenOrderNumber(repo model.Repo, idcID uint) string {
	//格式：IDC+$idcID+YYMMDD+001, 如 IDC120190617001
	_ = repo
	//今天是几号：20190617
	date := time.Now().Format("20060102")
	yesterday := time.Now().AddDate(0, 0, -1).Format("20060102")
	if _, ok := GlobalDailyCount[date]; ok {
		GlobalDailyCount[date]++
	} else {
		//TODO：如果是服务重启，从db中恢复最大的订单号
		dbMax, _ := repo.GetMaxOrderNumber(date)
		//如果是新的一天，计数从1开始计
		GlobalDailyCount[date] = dbMax + 1
		//把昨天的计数清零，以免内存逐渐增大
		delete(GlobalDailyCount, yesterday)

	}
	return fmt.Sprintf("IDC%d%s%03d", idcID, date, GlobalDailyCount[date])
}

//OrderStatusTransfer 订单状态值和数据库存储值的转换
func OrderStatusTransfer(status string, reverse bool) string {
	mStatus := map[string]string{
		"采购中":  "purchasing",
		"部分到货": "partly_arrived",
		"全部到货": "all_arrived",
		"已取消":  "canceled",
		"已完成":  "finished",
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

type PreOccupiedUsiteData struct {
	Label string //对应于机位编号
	Value uint   //对应于机位ID
}

// 释放预占用的机位
// 由于预占用的机位跟实际导入的机位有可能不一致，所以
// 当订单取消时，可以直接释放预占用的机位
// 确认订单时，释放当前仍为"预占用"的机位信息。如果机位已被导入设备，此时应该时"已使用"状态, 那是不能释放的。
func ReleasePreOccupiedUsites(log logger.Logger, repo model.Repo, oid uint) (affected int64, err error) {
	order, err := repo.GetOrderByID(oid)
	if err != nil {
		log.Errorf("get order by id(%d) fail, %v", oid, err)
		return 0, err
	}
	us := make([]PreOccupiedUsiteData, 0)
	if err := json.Unmarshal([]byte(order.PreOccupiedUsites), &us); err != nil {
		log.Errorf("unmarshal usites(%s) fail，%v", order.PreOccupiedUsites, err)
		return 0, err
	}
	ids := make([]uint, len(us))
	for _, u := range us {
		cur, err := repo.GetServerUSiteByID(u.Value)
		if err != nil {
			log.Errorf("get usite by id(%d) fail, %v", u.Value, err)
			return 0, err
		}
		if cur.Status == model.USiteStatPreOccupied {
			ids = append(ids, cur.ID)
		}
	}
	return BatchUpdateServerUSitesStatus(repo, ids, model.USiteStatFree)
}

// ExportedOrders 导出订单详情集合
type ExportedOrders []*OrderResp

// ToTableRecords 生成用于表格显示的二维字符串切片
func (items ExportedOrders) ToTableRecords() (records [][]string) {
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
		records = append(records, []string{
			items[i].Number,
			items[i].PhysicalArea,
			items[i].Usage,
			items[i].Category,
			strconv.Itoa(int(items[i].Amount)),
			idcName,
			serverRoomName,
			items[i].PreOccupiedUsites,
			items[i].ExpectedArrivalDate,
		})
	}
	return records
}

// 把预占机位格式化输出
func formatUsite(repo model.Repo, jsonU string, log logger.Logger) (string, error) {
	us := make([]PreOccupiedUsiteData, 0)
	if err := json.Unmarshal([]byte(jsonU), &us); err != nil {
		log.Errorf("unmarshal usites(%s) fail，%v", jsonU, err)
		return "", err
	}
	// 构造一棵机架机位树
	// cabient-number1
	//		usiteA
	//		usiteB
	// cabient-number2
	//		usiteA
	//		usiteB
	tree := make(map[string][]string, 0)
	//ulist := make([]string, 0)
	for _, u := range us {
		cur, err := repo.GetServerUSiteByID(u.Value)
		if err != nil {
			log.Errorf("get usite by id(%d) fail, %v", u.Value, err)
			return "", err
		}
		cabinet, err := repo.GetServerCabinetByID(cur.ServerCabinetID)
		if err != nil {
			log.Errorf("get server cabinet by id(%d) fail, %v", cur.ServerCabinetID, err)
			return "", err
		}
		tree[cabinet.Number] = append(tree[cabinet.Number], cur.Number)
	}

	lines := make([]string, 0)
	for key, vals := range tree {
		lines = append(lines, fmt.Sprintf("%s:", key))
		groupSize := 5 //每行显示几个
		for i := 0; i < len(vals); i++ {
			j := i / groupSize
			if i%groupSize == 0 {
				if groupSize*(i+1) > len(vals) {
					lines = append(lines, fmt.Sprintf("\t%s", stringsStd.Join(vals[j*groupSize:], ";")))
				} else /*if len(vals) > groupSize*i */ {
					lines = append(lines, fmt.Sprintf("\t%s", stringsStd.Join(vals[j*groupSize:(j+1)*groupSize], ";")))
				}
			}
		}
	}
	return stringsStd.Join(lines, "\n"), nil
}
