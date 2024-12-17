package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
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
	"idcos.io/cloudboot/utils/times"
	"idcos.io/cloudboot/utils/upload"
)

// SaveNetworkAreaReq 保存网络区域请求结构体
type SaveNetworkAreaReq struct {
	IDCID uint `json:"-"`
	// 所属机房ID
	// Required: true
	ServerRoomID uint `json:"server_room_id"`
	// 网络区域ID。若id=0，则新增。若id>0，则修改。
	ID uint `json:"id"`
	// 网络区域名称
	// Required: true
	Name string `json:"name"`
	// 关联物理区域
	// Required: true
	PhysicalArea []string `json:"physical_area"`
	// 状态。状态。nonproduction-未投产; production-已投产; offline-已下线(回收)
	// Required: true
	// Enum: nonproduction,production,offline
	Status string `json:"status"`
	// 用户登录名
	LoginName string `json:"-"`
}

// FieldMap 请求字段映射
func (reqData *SaveNetworkAreaReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.ServerRoomID: "server_room_id",
		&reqData.ID:           "id",
		&reqData.Name:         "name",
		&reqData.PhysicalArea: "physical_area",
		&reqData.Status:       "status",
	}
}

// Validate 结构体数据校验
func (reqData *SaveNetworkAreaReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())

	if reqData.ServerRoomID == 0 {
		errs.Add([]string{"server_room_id"}, binding.RequiredError, "机房管理单元ID不能为空")
		return errs
	}
	// 校验指定数据中心下的指定机房是否真实存在
	room, err := repo.GetServerRoomByID(reqData.ServerRoomID)
	if err == gorm.ErrRecordNotFound {
		errs.Add([]string{"idc_id", "server_room_id"}, binding.BusinessError, "该机房管理单元不存在")
		return errs
	}
	if err != nil {
		errs.Add([]string{"idc_id", "server_room_id", "name"}, binding.SystemError, "系统内部错误")
		return errs
	}
	reqData.IDCID = room.IDCID

	if reqData.Name == "" {
		errs.Add([]string{"name"}, binding.RequiredError, "网络区域名称不能为空")
		return errs
	}

	cond := model.NetworkAreaCond{
		IDCID:        []uint{reqData.IDCID},
		ServerRoomID: []uint{reqData.ServerRoomID},
		Name:         reqData.Name, // 模糊匹配
	}

	items, err := repo.GetNetworkAreas(&cond, nil, nil)
	if err != nil {
		errs.Add([]string{"idc_id", "server_room_id", "name"}, binding.SystemError, "系统内部错误")
		return errs
	}
	for _, item := range items {
		if (reqData.ID == 0 && item.Name == reqData.Name) || // 新增时，网络区域名称不能重复。
			(reqData.ID > 0 && item.Name == reqData.Name && reqData.ID != item.ID) { // 更新时，网络区域名称不能重复（除了自身外）。
			errs.Add([]string{"name"}, binding.BusinessError, "同一个机房管理单元内网络区域名称不能重复")
			return errs
		}
	}

	if len(reqData.PhysicalArea) == 0 {
		errs.Add([]string{"physical_area"}, binding.RequiredError, "关联物理区域不能为空")
		return errs
	}
	if reqData.Status == "" {
		errs.Add([]string{"status"}, binding.RequiredError, "状态不能为空")
		return errs
	}

	if !collection.InSlice(reqData.Status, []string{model.NetworkAreaStatNonProduction, model.NetworkAreaStatProduction, model.NetworkAreaStatOffline}) {
		errs.Add([]string{"status"}, binding.RequiredError, "无效的状态值")
		return errs
	}
	return errs
}

// PhysicalArea 物理区域
type PhysicalArea struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// convert2PhysicalAreas 将物理区域字符串切片转成物理区域对象切片
func convert2PhysicalAreas(items []string) []*PhysicalArea {
	if len(items) == 0 {
		return []*PhysicalArea{}
	}
	areas := make([]*PhysicalArea, 0, len(items))
	for i := range items {
		areas = append(areas, &PhysicalArea{
			ID:   uint(i + 1),
			Name: items[i],
		})
	}
	return areas
}

// PhysicalAreas 物理区域集合
type PhysicalAreas []*PhysicalArea

// ToJSON 转化成JSON字节切片
func (items PhysicalAreas) ToJSON() []byte {
	b, _ := json.Marshal(items)
	return b
}

// SaveNetworkArea 保存网络区域
func SaveNetworkArea(log logger.Logger, repo model.Repo, reqData *SaveNetworkAreaReq) (err error) {
	if reqData.ID > 0 {
		_, err = repo.GetNetworkAreaByID(reqData.ID)
		if err != nil {
			return err
		}
	}

	na := model.NetworkArea{
		IDCID:        reqData.IDCID,
		ServerRoomID: reqData.ServerRoomID,
		Name:         reqData.Name,
		Status:       reqData.Status,
		Creator:      reqData.LoginName,
	}
	na.ID = reqData.ID
	na.PhysicalArea = string(PhysicalAreas(convert2PhysicalAreas(reqData.PhysicalArea)).ToJSON())
	_, err = repo.SaveNetworkArea(&na)
	if reqData.ID == 0 {
		reqData.ID = na.ID
	}
	return err
}

// NetworkArea 网络区域
type NetworkArea struct {
	// 数据中心
	IDC struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"idc"`
	// 机房管理单元
	ServerRoom struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"server_room"`
	// 主键ID
	ID uint `json:"id"`
	// 网络区域名
	Name string `json:"name"`
	// 物理区域
	PhysicalArea []*PhysicalArea `json:"physical_area"`
	// Status 状态
	Status string `json:"status"`
	// 创建时间
	CreatedAt string `json:"created_at"`
	// 修改时间
	UpdatedAt string `json:"updated_at"`
}

// GetNetworkAreaPageReq 查询网络区域分页请求结构体
type GetNetworkAreaPageReq struct {
	// 数据中心ID
	IDCID string `json:"idc_id"`
	// 机房管理单元ID
	ServerRoomID string `json:"server_room_id"`
	// 机房名
	ServerRoomName string `json:"server_room_name"`
	// 网络区域名称
	Name string `json:"name"`
	// 关联物理区域
	PhysicalArea string `json:"physical_area"`
	// 状态。可选值: nonproduction-未投产; production-已投产; offline-已下线(回收);
	// Enum: nonproduction,production,offline
	Status string `json:"status"`
	// 分页页号
	Page int64 `json:"page"`
	// 分页大小
	PageSize int64 `json:"page_size"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *GetNetworkAreaPageReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDCID:          "idc_id",
		&reqData.ServerRoomID:   "server_room_id",
		&reqData.ServerRoomName: "server_room_name",
		&reqData.Name:           "name",
		&reqData.PhysicalArea:   "physical_area",
		&reqData.Status:         "status",
		&reqData.Page:           "page",
		&reqData.PageSize:       "page_size",
	}
}

// GetNetworkAreaPage 按条件查询网络区域分页列表
func GetNetworkAreaPage(log logger.Logger, repo model.Repo, reqData *GetNetworkAreaPageReq) (pg *page.Page, err error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 100 {
		reqData.PageSize = 10
	}
	if reqData.Page < 0 {
		reqData.Page = 0
	}

	cond := model.NetworkAreaCond{
		IDCID:          strings2.Multi2UintSlice(reqData.IDCID),
		ServerRoomID:   strings2.Multi2UintSlice(reqData.ServerRoomID),
		Name:           reqData.Name,
		PhysicalArea:   reqData.PhysicalArea,
		Status:         reqData.Status,
		ServerRoomName: reqData.ServerRoomName,
	}

	totalRecords, err := repo.CountNetworkAreas(&cond)
	if err != nil {
		return nil, err
	}

	pager := page.NewPager(reflect.TypeOf(&NetworkArea{}), reqData.Page, reqData.PageSize, totalRecords)
	items, err := repo.GetNetworkAreas(&cond, model.OneOrderBy("updated_at", model.DESC), pager.BuildLimiter())
	if err != nil {
		return nil, err
	}

	for i := range items {
		pager.AddRecords(convert2NetworkArea(repo, items[i]))
	}
	return pager.BuildPage(), nil
}

func convert2NetworkArea(repo model.Repo, src *model.NetworkArea) *NetworkArea {
	dst := NetworkArea{
		ID:        src.ID,
		Name:      src.Name,
		Status:    src.Status,
		CreatedAt: src.CreatedAt.Format(times.DateLayout),
		UpdatedAt: src.UpdatedAt.Format(times.DateLayout),
	}

	if idc, _ := repo.GetIDCByID(src.IDCID); idc != nil {
		dst.IDC.ID, dst.IDC.Name = idc.ID, idc.Name
	}

	if room, _ := repo.GetServerRoomByID(src.ServerRoomID); room != nil {
		dst.ServerRoom.ID, dst.ServerRoom.Name = room.ID, room.Name
	}

	if src.PhysicalArea != "" {
		_ = json.Unmarshal([]byte(src.PhysicalArea), &dst.PhysicalArea)
	}
	return &dst
}

// GetNetworkAreaByID 返回指定ID的网络区域
func GetNetworkAreaByID(repo model.Repo, id uint) (*NetworkArea, error) {
	one, err := repo.GetNetworkAreaByID(id)
	if err != nil {
		return nil, err
	}
	return convert2NetworkArea(repo, one), nil
}

// UpdateNetworkAreasStatusReq 批量修改网络区域状态API
type UpdateNetworkAreasStatusReq struct {
	// 状态。状态。nonproduction-未投产; production-已投产; offline-已下线(回收)
	// Required: true
	// Enum: nonproduction,production,offline
	Status string `json:"status"`
	// 网络区域ID列表
	// Required: true
	IDs []uint `json:"ids"`
}

// FieldMap 请求字段映射
func (reqData *UpdateNetworkAreasStatusReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDs:    "ids",
		&reqData.Status: "status",
	}
}

// Validate 结构体数据校验
func (reqData *UpdateNetworkAreasStatusReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if len(reqData.IDs) == 0 {
		errs.Add([]string{"ids"}, binding.RequiredError, "网络区域ID不能为空")
		return errs
	}
	if reqData.Status == "" {
		errs.Add([]string{"status"}, binding.RequiredError, "状态值不能为空")
		return errs
	}

	if !collection.InSlice(reqData.Status, []string{model.NetworkAreaStatNonProduction, model.NetworkAreaStatProduction /*, model.NetworkAreaStatOffline*/}) {
		errs.Add([]string{"status"}, binding.RequiredError, "无效的状态值")
		return errs
	}
	return errs
}

// UpdateNetworkAreasStatus 批量修改网络区域状态
func UpdateNetworkAreasStatus(repo model.Repo, reqData *UpdateNetworkAreasStatusReq) (affected int64, err error) {
	return repo.UpdateNetworkAreaStatus(reqData.Status, reqData.IDs...)
}

// RemoveNetworkAreaReq 移除网络区域请求结构体
type RemoveNetworkAreaReq struct {
	ID uint `json:"-"`
}

// Validate 结构体数据校验
func (reqData *RemoveNetworkAreaReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())
	// 删除前须确保其下不含任何机架
	count, err := repo.CountServerCabinets(&model.ServerCabinetCond{
		NetworkAreaID: []uint{reqData.ID},
	})
	if err != nil {
		errs.Add([]string{"id"}, binding.SystemError, "系统内部错误")
		return errs
	}
	if count > 0 {
		errs.Add([]string{"id"}, binding.RequiredError, "请先删除该网络区域下的所有机架")
		return errs
	}
	return errs
}

// RemoveNetworkArea 删除指定ID的业务网段
func RemoveNetworkArea(repo model.Repo, reqData *RemoveNetworkAreaReq) (err error) {
	affected, err := repo.RemoveNetworkAreaByID(reqData.ID)
	if err != nil {
		return err
	}
	if affected <= 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

//NetworkAreaForImport 为了导入
type NetworkAreaForImport struct {
	//文件导入数据
	// 网络区域名
	Name string `json:"name"`
	// 机房管理单元
	ServerRoomName string `json:"server_room_name"`
	// Status 状态
	Status string `json:"status"`
	// 物理区域
	PhysicalArea string `json:"physical_area"`

	//创建人
	Creator string `json:"creator"`
	// 主键ID
	ID uint `json:"id"`
	//机房ID
	IDCID uint `json:"idc_id"`
	//一级机房
	ServerRoomID uint `json:"server_room_id"`
	//附加内容
	Content string `json:"content"`
}

//checkLength 对导入文件中的数据做基本验证
func (nafi *NetworkAreaForImport) checkLength() {
	leg := len(nafi.Name)
	if leg == 0 || leg > 255 {
		var br string
		if nafi.Content != "" {
			br = "<br />"
		}
		nafi.Content += br + fmt.Sprintf("必填项校验:区域名称长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(nafi.ServerRoomName)
	if leg == 0 || leg > 255 {
		var br string
		if nafi.Content != "" {
			br = "<br />"
		}
		nafi.Content += br + fmt.Sprintf("必填项校验:机房管理单元长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(nafi.Status)
	if leg == 0 || leg > 255 {
		var br string
		if nafi.Content != "" {
			br = "<br />"
		}
		nafi.Content += br + fmt.Sprintf("必填项校验:状态长度为(%d)(不能为空)", leg)
	}
	leg = len(nafi.PhysicalArea)
	if leg == 0 {
		var br string
		if nafi.Content != "" {
			br = "<br />"
		}
		nafi.Content += br + fmt.Sprintf("必填项校验:关联物理区域长度为(%d)(不能为空)", leg)
	}
}

var (
	networkAreaTypeMap = map[string]string{
		"未投产": model.NetworkAreaStatNonProduction,
		"已投产": model.NetworkAreaStatProduction,
		"回收":  model.NetworkAreaStatOffline,
		"":    "",
	}
)

//validate 对导入文件中的数据做基本验证
func (nafi *NetworkAreaForImport) validate(log logger.Logger, repo model.Repo) (int, error) {
	//机房校验
	sr, err := repo.GetServerRoomByName(nafi.ServerRoomName)
	if err != nil && err != gorm.ErrRecordNotFound {
		return upload.Return, err
	}
	if err == gorm.ErrRecordNotFound || sr == nil {
		var br string
		if nafi.Content != "" {
			br = "<br />"
		}
		nafi.Content += br + fmt.Sprintf("机房管理单元(%s)不存在", nafi.ServerRoomName)
		return upload.Continue, nil
	}
	nafi.IDCID = sr.IDCID
	nafi.ServerRoomID = sr.ID

	//网络区域校验
	nabc, err := repo.GetNetworkAreasByCond(&model.NetworkArea{Name: nafi.Name,
		ServerRoomID: nafi.ServerRoomID,
		IDCID:        nafi.IDCID,
	})
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Debug(err.Error())
		return upload.Return, err
	}
	if len(nabc) > 1 {
		var br string
		if nafi.Content != "" {
			br = "<br />"
		}
		nafi.Content += br + fmt.Sprintf("网络区域表存在冲突，网络区域(%s)在同一数据中心，同一机房下不能存在多个", nafi.Name)
		return upload.Continue, nil
	}
	if len(nabc) == 1 {
		nafi.ID = nabc[0].ID
	}

	if networkAreaTypeMap[nafi.Status] == "" {
		var br string
		if nafi.Content != "" {
			br = "<br />"
		}
		nafi.Content += br + fmt.Sprintf("网络区域类型只能为(%s)", "未投产/已投产/回收")
		return upload.Continue, nil
	}

	nafi.Status = networkAreaTypeMap[nafi.Status]
	//关联物理区域不做校验

	return upload.DO, nil
}

//ImportNetworkAreaPriview 导入预览
func ImportNetworkAreaPriview(log logger.Logger, repo model.Repo, reqData *upload.ImportReq) (map[string]interface{}, error) {
	ra, err := utils.ParseDataFromXLSX(upload.UploadDir + reqData.FileName)
	if err != nil {
		return nil, err
	}
	length := len(ra)

	var success []*NetworkAreaForImport
	var failure []*NetworkAreaForImport
	for i := 1; i < length; i++ {
		row := &NetworkAreaForImport{}
		if len(ra[i]) < 4 {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "导入文件列长度不对（应为4列）"
			failure = append(failure, row)
			continue
		}

		row.Name = strings.TrimSpace(ra[i][0])
		row.ServerRoomName = strings.TrimSpace(ra[i][1])
		row.Status = strings.TrimSpace(ra[i][2])
		row.PhysicalArea = strings.TrimSpace(ra[i][3])

		//必填项校验
		row.checkLength()
		//机房和网络区域校验
		_, err := row.validate(log, repo)
		if err != nil {
			return nil, err
		}

		if row.Content != "" {
			failure = append(failure, row)
		} else {
			success = append(success, row)
		}
	}

	var data []*NetworkAreaForImport
	if len(failure) > 0 {
		data = failure
	} else {
		data = success
	}
	var result []*NetworkAreaForImport
	for i := 0; i < len(data); i++ {
		if uint(i) >= reqData.Offset && uint(i) < (reqData.Offset+reqData.Limit) {
			result = append(result, data[i])
		}
	}
	if len(failure) > 0 {
		return map[string]interface{}{"status": "failure",
			"message":       "导入网络区域错误",
			"import_status": false,
			"record_count":  len(data),
			"content":       result,
		}, nil
	}
	return map[string]interface{}{"status": "success",
		"message":       "操作成功",
		"import_status": true,
		"record_count":  len(data),
		"content":       result,
	}, nil
}

//ImportNetworkArea 将导入网络区域放到数据库
func ImportNetworkArea(log logger.Logger, repo model.Repo, reqData *upload.ImportReq) error {
	ra, err := utils.ParseDataFromXLSX(upload.UploadDir + reqData.FileName)
	if err != nil {
		return err
	}
	length := len(ra)
	for i := 1; i < length; i++ {
		row := &NetworkAreaForImport{}
		if len(ra[i]) < 4 {
			continue
		}

		row.Name = strings.TrimSpace(ra[i][0])
		row.ServerRoomName = strings.TrimSpace(ra[i][1])
		row.Status = strings.TrimSpace(ra[i][2])
		row.PhysicalArea = strings.TrimSpace(ra[i][3])

		//必填项校验
		row.checkLength()
		//机房和网络区域校验
		isSave, err := row.validate(log, repo)
		if err != nil {
			return err
		}

		//不能获取机房信息
		if isSave == upload.Continue {
			continue
		}

		na := model.NetworkArea{
			IDCID:        row.IDCID,
			ServerRoomID: row.ServerRoomID,
			Name:         row.Name,
			Status:       row.Status,
			Creator:      "", // TODO 待补充
		}
		na.ID = row.ID
		pa := strings.Split(row.PhysicalArea, ",")
		na.PhysicalArea = string(PhysicalAreas(convert2PhysicalAreas(pa)).ToJSON())
		log.Errorf("%v", na)
		if _, err = repo.SaveNetworkArea(&na); err != nil {
			return err
		}
	}
	defer os.Remove(upload.UploadDir + reqData.FileName)
	return nil
}
