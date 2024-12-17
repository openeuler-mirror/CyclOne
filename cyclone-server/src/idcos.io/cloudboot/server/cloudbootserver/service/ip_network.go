package service

import (
	"time"
	"errors"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/voidint/binding"
	"github.com/voidint/page"

	"idcos.io/cloudboot/config"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils"
	"idcos.io/cloudboot/utils/network"
	"idcos.io/cloudboot/utils/ping"
	strings2 "idcos.io/cloudboot/utils/strings"
	"idcos.io/cloudboot/utils/times"
	"idcos.io/cloudboot/utils/upload"
	nw "idcos.io/cloudboot/utils/network"
)

var (
	ipReg = regexp.MustCompile("^((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)$")
	ipv6Reg = regexp.MustCompile("^\\s*((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)(\\.(25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)(\\.(25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)(\\.(25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)(\\.(25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)(\\.(25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)(\\.(25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)(\\.(25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)){3}))|:)))(%.+)?\\s*$")
)

//IDCOrServerRoom 网段分页查询专用
type IDCOrServerRoom struct {
	//索引
	ID int `json:"id"`
	//名称
	Name string `json:"name"`
}

//SwitchDetail 交换机详情
type SwitchDetail struct {
	// 资产编号
	FixedAssetNumber string `json:"fixed_asset_number"`
	// 名称
	Name string `json:"name"`
	// Tor分组
	Tor string `json:"tor"`
}

// IPNetworkPage 网段
type IPNetworkPage struct {
	//数据中心
	IDC IDCOrServerRoom `json:"idc"`
	//机房
	ServerRoom IDCOrServerRoom `json:"server_room"`
	//网段ID
	ID int `json:"id"`
	//CIDR网段
	CIDR string `json:"cidr"`
	//网段类别(ilo-服务器ILO; tgw_intranet-服务器TGW内网; tgw_extranet-服务器TGW外网; intranet-服务器普通内网; extranet-服务器普通外网; v_intranet-服务器虚拟化内网;)
	Category string `json:"category"`
	//掩码
	Netmask string `json:"netmask"`
	//网关
	Gateway string `json:"gateway"`
	//业务IP资源池
	IPPool string `json:"ip_pool"`
	//带外IP资源池
	PXEPool string `json:"pxe_pool"`
	//交换机设备
	Switchs []SwitchDetail `json:"switchs"`
	//网络区域
	NetworkArea   string `json:"network_area"`
	NetworkAreaID uint   `json:"-"`
	//vlan
	Vlan string `json:"vlan"`
	//IP版本
	Version string `json:"version"`
	// 创建时间
	CreatedAt times.ISOTime `json:"created_at"`
	// 更新时间
	UpdatedAt times.ISOTime `json:"updated_at"`
}

// GetIPNetworkPageReq 查询网段分页请求结构体
type GetIPNetworkPageReq struct {
	//机房ID
	ServerRoomID   string `json:"server_room_id"`
	ServerRoomName string `json:"server_room_name"`
	//CIDR网段，支持模糊查询
	CIDR string `json:"cidr"`
	//网段类别。可选值: ilo-服务器ILO; tgw_intranet-服务器TGW内网; tgw_extranet-服务器TGW外网; intranet-服务器普通内网; extranet-服务器普通外网; v_intranet-服务器虚拟化内网;
	Category        string `json:"category"`
	Switches        string `json:"switches"`
	NetworkAreaID   string `json:"network_area_id"`
	NetworkAreaName string `json:"network_area_name"`
	//分页页号
	Page int64 `json:"page"`
	//分页大小
	PageSize int64 `json:"page_size"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *GetIPNetworkPageReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.ServerRoomID:    "server_room_id",
		&reqData.ServerRoomName:  "server_room_name",
		&reqData.Category:        "category",
		&reqData.CIDR:            "cidr",
		&reqData.Switches:        "switches",
		&reqData.NetworkAreaID:   "network_area_id",
		&reqData.NetworkAreaName: "network_area_name",
		&reqData.Page:            "page",
		&reqData.PageSize:        "page_size",
	}
}

//Validate 请求数据校验
func (reqData *GetIPNetworkPageReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())

	if reqData.ServerRoomID != "" {
		srs := strings2.Multi2UintSlice(reqData.ServerRoomID)
		for _, sr := range srs {
			if _, err := repo.GetServerRoomByID(sr); err == gorm.ErrRecordNotFound {
				errs.Add([]string{"server_room_id"}, binding.BusinessError, http.StatusText(http.StatusNotFound))
				return errs
			}
		}
	}

	// 校验CIDR
	if reqData.CIDR != "" {
		cidrs := strings2.MultiLines2Slice(reqData.CIDR)
		for _, cidr := range cidrs {
			cidrArr := strings.Split(cidr, "/")
			if len(cidrArr) != 2 || (!ipReg.MatchString(cidrArr[0]) && !ipv6Reg.MatchString(cidrArr[0])) {
				errs.Add([]string{"cidr"}, binding.BusinessError, "无效CIDR(正确格式为:IP/掩码长度)")
				return errs
			}
		}
	}

	// 校验Category
	if reqData.Category != "" {
		category := strings.Split(strings.TrimSpace(reqData.Category), ",")
		for k := range category {
			categoryIsGood := true
			switch category[k] {
			case model.ILO:
			case model.TGWIntranet:
			case model.TGWExtranet:
			case model.Intranet:
			case model.Extranet:
			case model.VIntranet:
			case model.VExtranet:
			default:
				categoryIsGood = false
			}
			if !categoryIsGood {
				errs.Add([]string{"category"}, binding.RequiredError, fmt.Sprintf("必填函数校验: category只能为 %s",
					"可选值: ilo-服务器ILO; tgw_intranet-服务器TGW内网; tgw_extranet-服务器TGW外网; intranet-服务器普通内网; extranet-服务器普通外网; v_intranet-服务器虚拟化内网"))
				return errs
			}
		}
	}
	return errs
}

//按固定格式提取固资号
func ExtractFN(name string) string {
	reg := regexp.MustCompile(`\([\w|\p{Han}]+\)?`)
	rawFN := reg.FindAllString(name, -1)
	fns := make([]string, 0, len(rawFN))
	for _, fn := range rawFN {
		fns = append(fns, strings.Trim(fn, "()"))
	}
	return strings.Join(fns, commaSep)
}

// GetIPNetworkPage 按条件查询业务网段分页列表
func GetIPNetworkPage(log logger.Logger, repo model.Repo, reqData *GetIPNetworkPageReq) (pg *page.Page, err error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 1000 {
		log.Warn("page size > 1000 or <= 0 , re-write to 10")
		reqData.PageSize = 10
	}
	if reqData.Page < 0 {
		reqData.Page = 0
	}

	//网络区域检索时，分开处理
	isFilterNetworkArea := false
	if reqData.NetworkAreaName != "" {
		isFilterNetworkArea = true
	}

	cond := model.IPNetworkCond{
		CIDR:           reqData.CIDR,
		Category:       reqData.Category,
		ServerRoomID:   strings2.Multi2UintSlice(reqData.ServerRoomID),
		Switches:       reqData.Switches,
		ServerRoomName: reqData.ServerRoomName,
	}

	totalRecords, err := repo.CountIPNetworks(&cond)
	if err != nil {
		return nil, err
	}

	pager := page.NewPager(reflect.TypeOf(&IPNetworkPage{}), reqData.Page, reqData.PageSize, totalRecords)
	var items []*model.IPNetwork
	if isFilterNetworkArea {
		items, err = repo.GetIPNetworks(&cond, model.OneOrderBy("id", model.DESC), nil)
	} else {
		items, err = repo.GetIPNetworks(&cond, model.OneOrderBy("id", model.DESC), pager.BuildLimiter())
	}
	if err != nil {
		return nil, err
	}

	for i := range items {
		ipnp, err := convert2IPNetwork(items[i], log, repo)
		if err != nil {
			return nil, err
		}
		if ipnp != nil {
			//网络区域的过滤放到这里
			if isFilterNetworkArea {
				netAreas := strings2.MultiLines2Slice(reqData.NetworkAreaName)
				for _, netArea := range netAreas {
					if strings.Contains(ipnp.NetworkArea, netArea) {
						pager.AddRecords(ipnp)
						break
					}
				}
			} else {
				pager.AddRecords(ipnp)
			}
		}
	}

	//这里可能有bug，当查询网络区域时，计数和分页可能不对
	if isFilterNetworkArea {
		p := pager.BuildPage()
		p.TotalRecords = int64(len(p.Records))
		return p, nil
	}
	return pager.BuildPage(), nil
}

// convert2IPNetwork 将model层的业务网段对象转化成service层的业务网段对象
func convert2IPNetwork(src *model.IPNetwork, log logger.Logger, repo model.Repo) (*IPNetworkPage, error) {
	if src == nil {
		return nil, nil
	}

	ipnp := &IPNetworkPage{
		ID:        int(src.ID),
		CIDR:      src.CIDR,
		Netmask:   src.Netmask,
		Category:  src.Category,
		Gateway:   src.Gateway,
		IPPool:    src.IPPool,
		PXEPool:   src.PXEPool,
		Vlan:      src.Vlan,
		Version:   src.Version,
		CreatedAt: times.ISOTime(src.CreatedAt),
		UpdatedAt: times.ISOTime(src.UpdatedAt),
	}

	sr, err := repo.GetServerRoomByID(src.ServerRoomID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if sr != nil {
		ipnp.ServerRoom = IDCOrServerRoom{
			ID:   int(sr.ID),
			Name: sr.Name,
		}

		idc, err := repo.GetIDCByID(sr.IDCID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}

		if idc != nil {
			ipnp.IDC = IDCOrServerRoom{
				ID:   int(idc.ID),
				Name: idc.Name,
			}
		}
	}

	//查询交换机设备详情
	var switchs []string
	_ = json.Unmarshal([]byte(src.Switches), &switchs)
	for k := range switchs {
		nd, err := repo.GetNetworkDeviceByFixedAssetNumber(switchs[k])

		if err != nil {
			return nil, err
		}
		if len(nd) == 1 {
			ipnp.Switchs = append(ipnp.Switchs, SwitchDetail{
				FixedAssetNumber: nd[0].FixedAssetNumber,
				Name:             nd[0].Name,
				Tor:              nd[0].TOR,
			})
			//网络区域的获取方法是,交换机>机架>网络区域>
			if ipnp.NetworkArea == "" {
				if cabinet, err := repo.GetServerCabinetByID(nd[0].ServerCabinetID); err != nil {
					log.Errorf("get cabinet by id:%d fail", nd[0].ServerCabinetID)
					continue
				} else if cabinet != nil {
					if na, err := repo.GetNetworkAreaByID(cabinet.NetworkAreaID); err != nil {
						log.Errorf("get net_area by id:%d fail", cabinet.NetworkAreaID)
						continue
					} else if na != nil {
						ipnp.NetworkArea = na.Name
						ipnp.NetworkAreaID = na.ID
					}
				}
			}
		}
		if len(nd) > 1 {
			return nil, fmt.Errorf("网络设备数据查询不唯一,固资编号为(%s)", switchs[k])
		}

	}

	return ipnp, nil
}

// GetIPNetworkByID 返回指定ID的网段
func GetIPNetworkByID(log logger.Logger, repo model.Repo, id uint) (ipnet *IPNetworkPage, err error) {
	one, err := repo.GetIPNetworkByID(id)
	if err != nil {
		return nil, err
	}

	ipnet, err = convert2IPNetwork(one, log, repo)
	if err != nil {
		return nil, err
	}
	return ipnet, nil
}

// SaveIPNetworkReq 保存(新增/更新)业务网段请求结构体
type SaveIPNetworkReq struct {
	// (required): 所属机房管理单元ID
	ServerRoomID int `json:"server_room_id"`
	// 网段ID。若id=0，则新增。若id>0，则修改
	ID int `json:"id"`
	//(required): CIDR网段
	CIDR string `json:"cidr"`
	//(required): 网段类别。可选值: ilo-服务器ILO; tgw_intranet-服务器TGW内网; tgw_extranet-服务器TGW外网; intranet-服务器普通内网; extranet-服务器普通外网; v_intranet-服务器虚拟化内网;
	Category string `json:"category"`
	//(required): 掩码
	Netmask string `json:"netmask"`
	// (required): 网关
	Gateway string `json:"gateway"`
	//(required): 业务IP池
	IPPool string `json:"ip_pool"`
	//(required): PXE IP资源池
	PXEPool string `json:"pxe_pool"`
	//(required): 网段作用范围内的交换机固定资产编号字符串数组
	Switchs []string `json:"switchs"`
	//(required): VLAN
	Vlan  string `json:"vlan"`
	// Version : ipv4 ipv6
	Version string `json:"version"`
	IDCID uint
	// 用户登录名
	LoginName string `json:"-"`
}


// FieldMap 请求参数与结构体字段建立映射
func (reqData *SaveIPNetworkReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.ServerRoomID: "server_room_id",
		&reqData.ID:           "id",
		&reqData.Gateway:      "gateway",
		&reqData.Vlan:         "vlan",
		&reqData.CIDR:         "cidr",
		&reqData.Category:     "category",
		&reqData.Netmask:      "netmask",
		&reqData.Switchs:      "switchs",
		&reqData.IPPool:       "ip_pool",
		&reqData.PXEPool:      "pxe_pool",
		&reqData.Version:      "version",
	}
}

// Validate 结构体数据校验
func (reqData *SaveIPNetworkReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())

	if reqData.ServerRoomID > 0 {
		sr, err := repo.GetServerRoomByID(uint(reqData.ServerRoomID))
		if err == gorm.ErrRecordNotFound {
			errs.Add([]string{"server_room_id"}, binding.BusinessError, http.StatusText(http.StatusNotFound))
			return errs
		}
		reqData.IDCID = sr.IDCID
	}

	// 校验CIDR
	if reqData.CIDR == "" {
		errs.Add([]string{"cidr"}, binding.RequiredError, fmt.Sprintln("必填函数校验: CIDR为空"))
		return errs
	}

	// 校验Category
	categoryIsGood := true
	switch reqData.Category {
	case "ilo":
	case "tgw_intranet":
	case "tgw_extranet":
	case "intranet":
	case "extranet":
	case "v_intranet":
	case "v_extranet":
	default:
		categoryIsGood = false
	}
	if !categoryIsGood {
		errs.Add([]string{"category"}, binding.RequiredError, fmt.Sprintf("必填函数校验: category只能为 %s",
			"可选值: ilo-服务器ILO; tgw_intranet-服务器TGW内网; tgw_extranet-服务器TGW外网; intranet-服务器普通内网; extranet-服务器普通外网; v_intranet-服务器虚拟化内网; v_extranet-服务器虚拟化外网"))
		return errs
	}

	// 获取已有全量网段信息并校验
	items, err := repo.GetIPNetworks(&model.IPNetworkCond{
		CIDR: reqData.CIDR,
	}, nil, nil)
	if err != nil {
		errs.Add([]string{"cidr"}, binding.SystemError, "内部访问数据库错误")
		return errs
	}
	for _, item := range items {
		if (reqData.ID == 0 && (item.CIDR == reqData.CIDR)) || // 新增时，CIDR不能重复
			((reqData.ID > 0 && item.CIDR == reqData.CIDR) && (uint(reqData.ID) != item.ID)) { // 更新时，CIDR不能重复（除了自身外）
			errs.Add([]string{"cidr"}, binding.BusinessError, fmt.Sprintf("网段CIDR:%s已经存在", reqData.CIDR))
			return errs
		}
	}

	// 数据校验 - IPv4
	if reqData.Version == model.IPv4 {
		cidrArr := strings.Split(reqData.CIDR, "/")
		if len(cidrArr) != 2 || !ipReg.MatchString(cidrArr[0]) {
			errs.Add([]string{"cidr"}, binding.BusinessError, "IPv4-CIDR校验：无效的网段")
			return errs
		}
		// IPv4  掩码长度转换成详细数值
		maskLen, _ := strconv.Atoi(cidrArr[1])
		netMask := network.GetCidrIPMask(maskLen)
		if reqData.Netmask == "" {
			errs.Add([]string{"mask"}, binding.RequiredError, fmt.Sprintln("IPv4必填项校验: 掩码不能为空"))
			return errs
		}
		if !ipReg.MatchString(reqData.Netmask) {
			errs.Add([]string{"mask"}, binding.BusinessError, fmt.Sprintf("IPv4掩码(%s)不正确", reqData.Netmask))
			return errs
		}
	
		if reqData.Netmask != netMask {
			errs.Add([]string{"mask"}, binding.BusinessError, fmt.Sprintf("IPv4掩码(%s)与CIDR指定(%s)不匹配", reqData.Netmask, netMask))
			return errs
		}
		// 校验网关
		if reqData.Gateway == "" {
			errs.Add([]string{"gateway"}, binding.RequiredError, fmt.Sprintln("IPv4必填项校验: 网关不能为空"))
			return errs
		}
		// 校验非TGW网段Gateway为网段第一个IP
		if reqData.Category !="tgw_intranet" && reqData.Category !="tgw_extranet" {
			min, _ := nw.GetCidrIPRange(reqData.CIDR)
			if reqData.Gateway != min {
				errs.Add([]string{"gateway"}, binding.RequiredError, fmt.Sprintf("IPv4必填项校验: 非TGW网段网关%s校验失败，应为网段第一个IP: %s", reqData.Gateway, min))
				return errs
			}
		}
		// 资源池校验
		if reqData.IPPool == "" {
			errs.Add([]string{"ip_pool"}, binding.RequiredError, fmt.Sprintln("IPv4必填项校验: ip资源池不能为空"))
			return errs
		}
		ipPool := strings.Split(reqData.IPPool, ",")
		if len(ipPool) < 2 || !ipReg.MatchString(ipPool[0]) || !ipReg.MatchString(ipPool[1]) {
			errs.Add([]string{"ip_pool"}, binding.BusinessError, fmt.Sprintf("IPv4资源池(%s)不正确", reqData.IPPool))
			return errs
		}
		if !network.CIDRContainsIP(reqData.CIDR, ipPool[0]) || !network.CIDRContainsIP(reqData.CIDR, ipPool[1]) {
			errs.Add([]string{"ip_pool"}, binding.BusinessError, fmt.Sprintf("IPv4资源池(%s)不在指定的CIDR内", reqData.IPPool))
			return errs
		}
	}

	// 数据校验 - IPv6
	if reqData.Version == model.IPv6 {
		cidrArr := strings.Split(reqData.CIDR, "/")
		if len(cidrArr) != 2 || !ipv6Reg.MatchString(cidrArr[0]) {
			errs.Add([]string{"cidr"}, binding.BusinessError, "IPv6-CIDR校验：无效的网段")
			return errs
		}
		// 校验网关
		if reqData.Gateway == "" {
			errs.Add([]string{"gateway"}, binding.RequiredError, fmt.Sprintln("IPv6必填项校验: 网关不能为空"))
			return errs
		}		
	}

	// 重复校验
	if reqData.ID == 0 {
		totalRecords, err := repo.CountIPNetworks(&model.IPNetworkCond{
			CIDR:         reqData.CIDR,
			Category:     reqData.Category,
			ServerRoomID: []uint{uint(reqData.ServerRoomID)},
		})
		if err != nil {
			errs.Add([]string{"count_network"}, binding.SystemError, "内部访问数据库错误")
			return errs
		}
		if totalRecords > 0 {
			errs.Add([]string{"network_repeat"}, binding.BusinessError, fmt.Sprintf("重复添加网段(%s)", reqData.CIDR))
			return errs
		}
	}
	// 交换机校验 TODO
	return errs
}

// SaveIPNetwork 保存(新增/更新)网段及其网段内的IP
func SaveIPNetwork(log logger.Logger, repo model.Repo, reqData *SaveIPNetworkReq) (err error) {
	swb, _ := json.Marshal(reqData.Switchs)
	sw := string(swb)
	ipnet := &model.IPNetwork{
		IDCID:        reqData.IDCID,
		ServerRoomID: uint(reqData.ServerRoomID),
		Category:     reqData.Category,
		Gateway:      reqData.Gateway,
		CIDR:         reqData.CIDR,
		Netmask:      reqData.Netmask,
		Vlan:         reqData.Vlan,
		Version:      reqData.Version,
		Switches:     sw,
		IPPool:       reqData.IPPool,
		PXEPool:      reqData.PXEPool,
		Creator:      reqData.LoginName,
	}
	ipnet.ID = uint(reqData.ID)

	_, err = repo.SaveIPNetwork(ipnet)
	reqData.ID = int(ipnet.ID)
	return err
}


//RemoveIPNetworkValidte 删除操作校验
func RemoveIPNetworkValidte(log logger.Logger, repo model.Repo, id uint) string {
	//统计该网段是否存在已分配的IP
	count, _ := repo.CountIPs(&model.IPPageCond{
		IPNetworkID: []uint{id},
		IsUsed:      model.IPUsed,
	})
	if count > 0 {
		return fmt.Sprintf("网段下面的IP已经分配(%d)个,不允许删除", count)
	}
	return ""
}

// RemoveIPNetwork 删除指定ID的网段
func RemoveIPNetwork(log logger.Logger, repo model.Repo, id uint) (err error) {
	affected, err := repo.RemoveIPNetworkByID(id)
	if err != nil {
		return err
	}
	if affected <= 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// 批量删除网段请求结构体
type DelIPNetworkReq struct {
	IDs []uint `json:"ids"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *DelIPNetworkReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDs: "ids",
	}
}

// Validate 结构体数据校验
func (reqData *DelIPNetworkReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())
	for _, id := range reqData.IDs {
		if _, err := repo.GetIPNetworkByID(id); err != nil {
			errs.Add([]string{"id"}, binding.RequiredError, fmt.Sprintf("网段ID(%d)不存在", id))
			return errs
		}
	    //统计该网段是否存在已分配的IP
	    count, _ := repo.CountIPs(&model.IPPageCond{
	    	IPNetworkID: []uint{id},
	    	IsUsed:      model.IPUsed,
	    })
	    if count > 0 {
			errs.Add([]string{"id"}, binding.RequiredError, fmt.Sprintf("网段ID(%d)存在(%d)个已分配的IP,不允许删除", id, count))
	    	return errs
	    }		
	}
	return nil
}

//RemoveIPNetworks 删除指定ID的网段
func RemoveIPNetworks(log logger.Logger, repo model.Repo, reqData *DelIPNetworkReq) (affected int64, err error) {
	for _, id := range reqData.IDs {
		_, err := repo.RemoveIPNetworkByID(id)
		if err != nil {
			log.Errorf("delete ip_network(id=%d) fail,err:%v", id, err)
			return affected, err
		}

		affected++
	}
	return affected, err
}

//IPNetworkForIPSPage 为了查询IP分页
type IPNetworkForIPSPage struct {
	//IP网段索引
	ID int `json:"id"`
	//CIDR网段
	CIDR string `json:"cidr"`
	//掩码
	Netmask string `json:"netmask"`
	//网关
	Gateway string `json:"gateway"`
	//网段类别
	Category string `json:"category"`
	//IP版本
	Version string `json:"version"`	
}

// IPSPage 网段
type IPSPage struct {
	//网段信息
	IPNetwork IPNetworkForIPSPage `json:"ip_network"`
	//索引
	ID int `json:"id"`
	//ip地址，支持模糊查询
	IP string `json:"ip"`
	//网段类别。可选值: pxe-PXE用IP; business-业务用IP;
	Category string `json:"category"`
	//是否已经被使用。可选值: yes-是; no-否;
	IsUsed string `json:"is_used"`
	//占用IP的设备序列号。(支持模糊查询)
	SN string `json:"sn"`
	//固资编号
	FixedAssetNumber string `json:"fixed_asset_number"`
	//内外网
	Scope string `json:"scope"`
	//Remark string `json:"remark"`
	// 创建时间
	CreatedAt times.ISOTime `json:"created_at"`
	// 更新时间
	UpdatedAt times.ISOTime `json:"updated_at"`
}

// GetIPSPageReq 查询网段分页请求结构体
type GetIPSPageReq struct {
	//网段ID
	IPNetworkID string `json:"ip_network_id"`
	//CIDR网段
	CIDR string `json:"cidr"`
	//ip地址，支持模糊查询
	IP string `json:"ip"`
	//网段类别。
	Category string `json:"ipnetwork_category"`
	//是否已经被使用。可选值: yes-是; no-否;
	IsUsed string `json:"is_used"`
	//占用IP的设备序列号。(支持模糊查询)
	SN string `json:"sn"`
	//固资编号
	FixedAssetNumber string `json:"fixed_asset_number"`
	//内外网
	Scope string `json:"scope"`
	//分页页号
	Page int64 `json:"page"`
	//分页大小
	PageSize int64 `json:"page_size"`
	// 用于勾选并导出
	ID               string
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *GetIPSPageReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IPNetworkID:      "ip_network_id",
		&reqData.CIDR:             "cidr",
		&reqData.IP:               "ip",
		&reqData.Category:         "ipnetwork_category",
		&reqData.IsUsed:           "is_used",
		&reqData.SN:               "sn",
		&reqData.FixedAssetNumber: "fixed_asset_number",
		&reqData.Scope:            "scope",
		&reqData.Page:             "page",
		&reqData.PageSize:         "page_size",
		&reqData.ID:               "id",
	}
}

//Validate 请求数据校验
func (reqData *GetIPSPageReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	//repo, _ := middleware.RepoFromContext(req.Context())
	//if reqData.IPNetworkID !="" {
	//	if _, err := repo.GetIPNetworkByID(uint(reqData.IPNetworkID)); err == gorm.ErrRecordNotFound {
	//		errs.Add([]string{"ip_network_id"}, binding.BusinessError, http.StatusText(http.StatusNotFound))
	//		return errs
	//	}
	//}

	if reqData.IP != "" {
		ip := reqData.IP

		// 格式校验
		ips := strings.Split(ip, ".")
		iplen := len(ips)
		for i := 0; i < iplen; i++ {
			if !strings.Contains(ips[i], "*") {
				seg, _ := strconv.Atoi(ips[i])
				if seg > 255 {
					errs.Add([]string{"ip"}, binding.RequiredError, fmt.Sprintf("IP(%s)格式不对", reqData.IP))
					return errs
				}
			}
		}

		// 长度不够补齐
		if iplen < 4 {
			for i := iplen; i < 4; i++ {
				ip = ip + ".*"
			}
		}

		if strings.Contains(ip, "*") {
			ip = strings.Replace(ip, "*", "%", -1)
		}

		reqData.IP = ip
	}

	// 校验Category
	//if reqData.Category != "" {
	//	categoryIsGood := true
	//	switch reqData.Category {
	//	case model.PXEIP:
	//	case model.BusinessIP:
	//	default:
	//		categoryIsGood = false
	//	}
	//	if !categoryIsGood {
	//		errs.Add([]string{"category"}, binding.RequiredError, fmt.Sprintf("category只能为 %s",
	//			"pxe-PXE用IP; business-业务用IP"))
	//		return errs
	//	}
	//
	//}

	// 校验IsUsed
	if reqData.IsUsed != "" {
		vals := strings2.MultiLines2Slice(reqData.IsUsed)
		for _, val := range vals {
			usedIsGood := true
			switch val {
			case model.IPNotUsed:
			case model.IPUsed:
			case model.IPDisabled:
			default:
				usedIsGood = false
			}
			if !usedIsGood {
				errs.Add([]string{"is_used"}, binding.RequiredError, fmt.Sprintf("is_used只能为 %s",
					"yes-是; no-否; disabled-不可用"))
				return errs
			}
		}
	}

	//校验scope
	if reqData.Scope != "" {
		switch reqData.Scope {
		case model.IPScopeIntranet:
		case model.IPScopeExtranet:
		default:
			errs.Add([]string{"scope"}, binding.RequiredError, fmt.Sprintf("scope %s",
				"intranet-内网/extranet-外网"))
			return errs
		}
	}
	//SN校验 TODO
	return errs
}

// GetIPSPage 按条件查询业务网段分页列表
func GetIPSPage(log logger.Logger, repo model.Repo, reqData *GetIPSPageReq) (pg *page.Page, err error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 100 {
		reqData.PageSize = 10
	}
	if reqData.Page < 0 {
		reqData.Page = 0
	}

	cond := model.IPPageCond{
		IPNetworkID:      strings2.Multi2UintSlice(reqData.IPNetworkID),
		CIDR:             reqData.CIDR,
		IP:               reqData.IP,
		Category:         reqData.Category,
		IsUsed:           reqData.IsUsed,
		SN:               reqData.SN,
		FixedAssetNumber: reqData.FixedAssetNumber,
		Scope:            &reqData.Scope,
	}

	totalRecords, err := repo.CountIPs(&cond)
	if err != nil {
		return nil, err
	}

	pager := page.NewPager(reflect.TypeOf(&IPSPage{}), reqData.Page, reqData.PageSize, totalRecords)
	items, err := repo.GetIPs(&cond, model.OneOrderBy("id", model.DESC), pager.BuildLimiter())
	if err != nil {
		return nil, err
	}

	for i := range items {
		ipnp, err := convert2IPS(items[i], repo)
		if err != nil {
			return nil, err
		}
		if ipnp != nil {
			pager.AddRecords(ipnp)
		}
	}
	return pager.BuildPage(), nil
}

//GetExportIP 获取导出IP列表
func GetExportIP(log logger.Logger, repo model.Repo, conf *config.Config, reqData *GetIPSPageReq) (ips []*IPSPage, err error) {
	cond := model.IPPageCond{
		ID:               strings2.Multi2UintSlice(reqData.ID),
		IPNetworkID:      strings2.Multi2UintSlice(reqData.IPNetworkID),
		CIDR:             reqData.CIDR,
		IP:               reqData.IP,
		Category:         reqData.Category,
		IsUsed:           reqData.IsUsed,
		SN:               reqData.SN,
		FixedAssetNumber: reqData.FixedAssetNumber,
		Scope:            &reqData.Scope,
	}

	items, err := repo.GetIPs(&cond, model.OneOrderBy("id", model.DESC), nil)
	if err != nil {
		return nil, err
	}
	ips = make([]*IPSPage, 0, len(items))
	for i := range items {
		item, err := convert2IPExport(items[i], repo)
		if err != nil {
			log.Error(err)
		}
		ips = append(ips, item)
	}
	return
}

// convert2IPS 将model层的业务网段对象转化成service层的业务网段对象
func convert2IPS(src *model.IPCombined, repo model.Repo) (*IPSPage, error) {
	if src == nil {
		return nil, nil
	}

	ipsp := &IPSPage{
		ID:               int(src.ID),
		IP:               src.IP.IP,
		Category:         src.Category,
		IsUsed:           src.IsUsed,
		SN:               src.SN,
		FixedAssetNumber: src.FixedAssetNumber,
		CreatedAt:        times.ISOTime(src.CreatedAt),
		UpdatedAt:        times.ISOTime(src.UpdatedAt),
	}
	if src.Scope != nil {
		ipsp.Scope = *src.Scope
	}
	if src.Remark != nil && src.IsUsed == model.IPDisabled {
		ipsp.FixedAssetNumber = *src.Remark
	}
	ipn, err := repo.GetIPNetworkByID(src.IPNetworkID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if ipn != nil {
		ipsp.IPNetwork = IPNetworkForIPSPage{
			ID:       int(ipn.ID),
			CIDR:     ipn.CIDR,
			Netmask:  ipn.Netmask,
			Gateway:  ipn.Gateway,
			Category: ipn.Category,
			Version:  ipn.Version,
		}
	}

	return ipsp, nil
}

// convert2IPExport 将model层的业务网段对象转化成excel导出详细内容
func convert2IPExport(src *model.IPCombined, repo model.Repo) (*IPSPage, error) {
	if src == nil {
		return nil, nil
	}
    
	IPNetworkCategoryMap := map[string]string {
		"ilo":            "服务器ILO",
		"tgw_intranet":   "服务器TGW内网",
		"tgw_extranet":   "服务器TGW外网",
		"intranet":       "服务器普通内网",
		"extranet":       "服务器普通外网",
		"v_intranet":     "服务器虚拟化内网",
		"v_extranet":     "服务器虚拟化外网",
	}
	IPCategoryMap := map[string]string {
		"pxe":        "PXE IP",
		"business":   "业务 IP",
	}
	YesNoMap := map[string]string {
		"yes":      "是",
		"no":       "否",
		"disabled": "不可用",
	}
	IPScopeMap := map[string]string {
		"intranet":    "内网(LA)",
		"extranet":    "外网(WA)",
	}

	ipsp := &IPSPage{
		ID:               int(src.ID),
		IP:               src.IP.IP,
		Category:         IPCategoryMap[src.Category],
		IsUsed:           YesNoMap[src.IsUsed],
		SN:               src.SN,
		FixedAssetNumber: src.FixedAssetNumber,
		CreatedAt:        times.ISOTime(src.CreatedAt),
		UpdatedAt:        times.ISOTime(src.UpdatedAt),
	}
	if src.Scope != nil {
		ipsp.Scope = IPScopeMap[*src.Scope]
	}
	if src.Remark != nil && src.IsUsed == model.IPDisabled {
		ipsp.FixedAssetNumber = *src.Remark
	}
	ipn, err := repo.GetIPNetworkByID(src.IPNetworkID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if ipn != nil {
		ipsp.IPNetwork = IPNetworkForIPSPage{
			ID:       int(ipn.ID),
			CIDR:     ipn.CIDR,
			Netmask:  ipn.Netmask,
			Gateway:  ipn.Gateway,
			Category: IPNetworkCategoryMap[ipn.Category],
			Version:  ipn.Version,
		}
	}

	return ipsp, nil
}

// AssignIPReq 手动分配IP请求结构
type AssignIPReq struct {
	//(required)IP的ID
	ID int `json:"id"`
	//(required): 设备序列号
	SN string `json:"sn"`
	//(required): IP种类("extranet"/"intranet")
	Scope string `json:"scope"`

	NetworkID uint   `json:"-"`
	IP        string `json:"-"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *AssignIPReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.SN:    "sn",
		&reqData.ID:    "id",
		&reqData.Scope: "scope",
	}
}

// Validate 结构体数据校验
func (reqData *AssignIPReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())
	log, _ := middleware.LoggerFromContext(req.Context())
	log.Debugf("id: %v, sn: %v, scope: %v\n", reqData.ID, reqData.SN, reqData.Scope)
	// 基本参数必要性检查
	if reqData.ID <= 0 {
		errs.Add([]string{"IP"}, binding.RequiredError, "必须指定IP的ID")
		return errs
	}
	if reqData.SN == "" {
		errs.Add([]string{"SN"}, binding.RequiredError, "必须指定物理机序列号")
		return errs
	}
	if reqData.Scope == "" {
		errs.Add([]string{"scope"}, binding.RequiredError, fmt.Sprintf("scope只能为: %s",
			"intranet-内网/extranet-外网"))
		return errs
	}

	//校验scope
	if reqData.Scope != "" {
		switch reqData.Scope {
		case model.IPScopeIntranet:
		case model.IPScopeExtranet:
		default:
			errs.Add([]string{"scope"}, binding.BusinessError, fmt.Sprintf("scope只能为: %s",
				"intranet-内网/extranet-外网"))
			return errs
		}
	}

	// 判断IP有效性
	ip, err := repo.GetIPByID(uint(reqData.ID))
	if gorm.IsRecordNotFoundError(err) || ip == nil {
		errs.Add([]string{"ip"}, binding.BusinessError, "指定的ip不存在,请刷新页面后重新尝试")
		return errs
	}
	if err != nil {
		errs.Add([]string{"ip"}, binding.SystemError, "系统内部错误")
		return errs
	}
	if ip.IsUsed == model.IPUsed {
		errs.Add([]string{"ip"}, binding.BusinessError, fmt.Sprintf("ip(%s)已经被设备(%s)占用", reqData.IP, ip.SN))
		return errs
	}
	if ip.IsUsed == model.IPDisabled {
		errs.Add([]string{"ip"}, binding.BusinessError, fmt.Sprintf("ip(%s)不可用(disabled)", reqData.IP))
		return errs
	}
	reqData.NetworkID = ip.IPNetworkID
	reqData.IP = ip.IP

	ipn, err := repo.GetIPNetworkByID(reqData.NetworkID)
	if gorm.IsRecordNotFoundError(err) || ip == nil {
		errs.Add([]string{"ipnetwork"}, binding.BusinessError, "指定的ip所在的网段存不存在,请刷新页面后重新尝试")
		return errs
	}
	if err != nil {
		errs.Add([]string{"ipnetwork"}, binding.SystemError, "系统内部错误")
		return errs
	}

	if strings.Contains(ipn.Category, model.IPScopeIntranet) {
		if reqData.Scope != model.IPScopeIntranet {
			errs.Add([]string{"SN"}, binding.BusinessError, "内网IP无法作为外网IP分配")
			return errs
		}
	}
	if strings.Contains(ipn.Category, model.IPScopeExtranet) {
		if reqData.Scope != model.IPScopeExtranet {
			errs.Add([]string{"SN"}, binding.BusinessError, "外网IP无法作为内网IP分配")
			return errs
		}
	}

	// 判断物理机是否存在
	dev, err := repo.GetDeviceBySN(reqData.SN)
	if gorm.IsRecordNotFoundError(err) {
		errs.Add([]string{"SN"}, binding.BusinessError, fmt.Sprintf("物理机(%s)不存在，请录入", reqData.SN))
		return errs
	}
	if dev.OperationStatus == model.DevOperStatInStore {
		errs.Add([]string{"SN"}, binding.BusinessError, fmt.Sprintf("%s物理机(%s)不允许直接分配IP",
			OperationStatusTransfer(dev.OperationStatus, true), reqData.SN))
		return errs
	}

	if err != nil {
		errs.Add([]string{"sns"}, binding.SystemError, "系统内部错误")
		return errs
	}

	// 校验网段是否一致
	if reqData.Scope == model.IPScopeIntranet {
		intranet, err := repo.GetIntranetIPNetworksBySN(reqData.SN)
		if gorm.IsRecordNotFoundError(err) || intranet == nil {
			errs.Add([]string{"ip_network"}, binding.BusinessError, fmt.Sprintf("物理机(%s)关联的网段存在问题(请检查对应的网络设备与网段信息是否匹配)", reqData.SN))
			return errs
		}
		if err != nil {
			errs.Add([]string{"ip"}, binding.SystemError, "系统内部错误")
			return errs
		}
		isError := true
		for k := range intranet {
			if intranet[k].Category == model.VIntranet || intranet[k].Category == model.VExtranet || intranet[k].Category == model.ILO {
				continue
			}
			if ip.IPNetworkID == intranet[k].ID {
				isError = false
			}
		}
		if isError {
			errs.Add([]string{"ip_network"}, binding.BusinessError, fmt.Sprintf("物理机(%s)所在的内网网段与指定IP(%s)所在内网网段不一致", reqData.SN, ip.IP))
			return errs
		}
	}

	if reqData.Scope == model.IPScopeExtranet {
		extranet, err := repo.GetExtranetIPNetworksBySN(reqData.SN)
		if gorm.IsRecordNotFoundError(err) || extranet == nil {
			errs.Add([]string{"ip_network"}, binding.BusinessError, fmt.Sprintf("物理机(%s)关联的网段存在问题(请检查对应的网络设备与网段信息是否匹配)", reqData.SN))
			return errs
		}
		if err != nil {
			errs.Add([]string{"ip"}, binding.SystemError, "系统内部错误")
			return errs
		}
		isError := true
		for k := range extranet {
			if extranet[k].Category == model.VIntranet || extranet[k].Category == model.VExtranet || extranet[k].Category == model.ILO {
				continue
			}
			if ip.IPNetworkID == extranet[k].ID {
				isError = false
			}
		}
		if isError {
			errs.Add([]string{"ip"}, binding.BusinessError, fmt.Sprintf("物理机(%s)所在的外网网段与指定IP(%s)所在外网网段不一致", reqData.SN, ip.IP))
			return errs
		}
	}

	return errs
}

// AssignIP 手动分配IP
func AssignIP(log logger.Logger, repo model.Repo, reqData *AssignIPReq) (err error) {
	//先ping检查，如果IP可以ping通，则失败
	if err = ping.PingTest(reqData.IP); err == nil {
		return fmt.Errorf("IP:%s连通性(ping)测试失败", reqData.IP)
	}

	//var ds model.DeviceSetting
	ds, err := repo.GetDeviceSettingBySN(reqData.SN)
	if ds == nil {
		ds = &model.DeviceSetting{
			SN:             reqData.SN,
			NeedExtranetIP: model.NO,
			NeedIntranetIPv6: model.NO,				//默认仅分配IPV4
			NeedExtranetIPv6: model.NO,				//默认仅分配IPV4
			InstallType:    model.InstallationPXE, 	 //这个值固定
			Status: 		model.InstallStatusSucc, //这个值固定
		}
	}
	//ds.SN = reqData.SN
	if reqData.Scope == model.IPScopeIntranet {
		//关联的网段唯一
		ds.IntranetIPNetworkID = reqData.NetworkID
		ds.IntranetIP = addIP(ds.IntranetIP, reqData.IP)
	} else if reqData.Scope == model.IPScopeExtranet {
		ds.ExtranetIPNetworkID = reqData.NetworkID
		ds.ExtranetIP = addIP(ds.ExtranetIP, reqData.IP)
		ds.NeedExtranetIP = model.YES
	}
	err = repo.SaveDeviceSetting(ds)
	if err != nil {
		log.Error(err)
		return err
	}

	//顺序放到最后，前面都成功了才占用
	err = repo.AssignIP(reqData.SN, reqData.Scope, uint(reqData.ID))
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

//在IP字段中追加一个IP，以逗号分隔
func addIP(oriIP, newIP string) string {
	ips := strings.Split(oriIP, commaSep)
	if len(ips) == 1 && ips[0] == "" {
		return newIP
	}
	ips = append(ips, newIP)
	return strings.Join(ips, commaSep)
}

//在IP字段中删除一个IP，比如某个IP被释放了
func removeIP(oriIP, delIP string) string {
	ips := strings.Split(oriIP, ",")
	newList := make([]string, 0, len(ips))
	for _, ip := range ips {
		if ip != delIP {
			newList = append(newList, ip)
		}
	}
	if len(newList) == 0 {
		return ""
	}
	return strings.Join(newList, ",")
}

type DisableIPReq struct {
	// IP的ID列表。
	IDs    []uint `json:"ids"`
	Remark string `json:"remark"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *DisableIPReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDs: "ids",
	}
}

// Validate 结构体数据校验
func (reqData *DisableIPReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	// 基本参数必要性检查
	if len(reqData.IDs) <= 0 {
		errs.Add([]string{"IP"}, binding.RequiredError, "必须指定IP的有效ID")
		return errs
	}

	repo, _ := middleware.RepoFromContext(req.Context())

	// 判断IP有效性
	for _, id := range reqData.IDs {
		ip, err := repo.GetIPByID(uint(id))
		if gorm.IsRecordNotFoundError(err) || ip == nil {
			errs.Add([]string{"ip"}, binding.BusinessError, fmt.Sprintf("ip(id:%d)不存在", id))
			return errs
		}
		if ip.IsUsed == model.IPUsed {
			errs.Add([]string{"ip"}, binding.BusinessError, fmt.Sprintf("ip(%s)已分配，请先释放", ip.IP))
			return errs
		}
	}
	return nil
}

// DisableIP 手动禁用IP
func DisableIP(log logger.Logger, repo model.Repo, reqData *DisableIPReq) (err error) {
	for _, id := range reqData.IDs {
		ip, err := repo.GetIPByID(uint(id))
		if gorm.IsRecordNotFoundError(err) || ip == nil {
			return err
		}
		ip.IsUsed = model.IPDisabled
		ip.Remark = &reqData.Remark
		if affected, err := repo.SaveIP(ip); err != nil || affected == 0 {
			log.Error(err)
			return err
		}
	}
	return nil
}

// UnassignIPReq 手动取消分配IP请求结构
type UnassignIPReq struct {
	// IP的ID。
	ID int `json:"id"`

	NetworkID uint   `json:"-"`
	IP        string `json:"-"`
	Scope     string `json:"-"`
	SN        string `json:"-"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *UnassignIPReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.ID: "id",
	}
}

// Validate 结构体数据校验
func (reqData *UnassignIPReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	// 基本参数必要性检查
	if reqData.ID <= 0 {
		errs.Add([]string{"IP"}, binding.RequiredError, "必须指定IP的有效ID")
		return errs
	}

	repo, _ := middleware.RepoFromContext(req.Context())

	// 判断IP有效性
	ip, err := repo.GetIPByID(uint(reqData.ID))
	if gorm.IsRecordNotFoundError(err) || ip == nil {
		errs.Add([]string{"ip"}, binding.BusinessError, fmt.Sprintf("ip(%d)不存在", reqData.ID))
		return errs
	}
	if err != nil {
		errs.Add([]string{"ip"}, binding.SystemError, "系统内部错误")
		return errs
	}
	if ip.IsUsed != model.IPUsed {
		errs.Add([]string{"ip"}, binding.BusinessError, fmt.Sprintf("ip(%d)未被占用，无需释放", reqData.ID))
		return errs
	}
	reqData.NetworkID = ip.IPNetworkID
	reqData.IP = ip.IP

	if ip.Scope != nil {
		reqData.Scope = *ip.Scope
	}
	reqData.SN = ip.SN

	dev, _ := repo.GetDeviceSettingBySN(reqData.SN)
	if dev == nil {
		reqData.SN = ""
	}
	return errs
}

// UnassignIP 手动取消分配IP
func UnassignIP(log logger.Logger, repo model.Repo, reqData *UnassignIPReq) (err error) {
	err = repo.UnassignIP(uint(reqData.ID))
	if err != nil {
		log.Error(err)
		return err
	}
	// 如果物理机没有配置装机参数，手动取消分配，不需要修改
	if reqData.SN != "" {
		ipn, err := repo.GetIPNetworkByID(uint(reqData.NetworkID))
		log.Debugf("UnassignIP of IPNetwork:%+v", ipn)
		ds, err := repo.GetDeviceSettingBySN(reqData.SN)
		if err != nil {
			log.Error(err)
			return err
		}
		ds.SN = reqData.SN
		if reqData.Scope == model.IPScopeIntranet {
			if ipn.Version == model.IPv4 {
				log.Debugf("Removing %s from %s", reqData.IP, ds.IntranetIP)
				ds.IntranetIP = removeIP(ds.IntranetIP, reqData.IP)
			} else {
				log.Debugf("Removing %s from %s", reqData.IP, ds.IntranetIPv6)
				ds.IntranetIPv6 = removeIP(ds.IntranetIPv6, reqData.IP)
			}
		}
		if reqData.Scope == model.IPScopeExtranet {
			if ipn.Version == model.IPv4 {
				log.Debugf("Removing %s from %s", reqData.IP, ds.ExtranetIP)
				ds.ExtranetIP = removeIP(ds.ExtranetIP, reqData.IP)
			} else {
				log.Debugf("Removing %s from %s", reqData.IP, ds.ExtranetIPv6)
				ds.ExtranetIPv6 = removeIP(ds.ExtranetIPv6, reqData.IP)
			}
		}
		_, err = repo.UpdateDeviceSettingBySN(ds)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

//ImportIPNetworksReq 设备导入Excel表对应字段
type ImportIPNetworksReq struct {
	IDCID          uint
	ServerRoomID   uint   `json:"server_room_id"`
	ServerRoomName string `json:"server_room_name"`
	//网段ID
	ID uint `json:"id"`
	//CIDR网段
	CIDR string `json:"cidr"`
	//网段类别(ilo-服务器ILO; tgw_intranet-服务器TGW内网; tgw_extranet-服务器TGW外网; intranet-服务器普通内网; extranet-服务器普通外网; v_intranet-服务器虚拟化内网;)
	Category string `json:"category"`
	//掩码
	Netmask string `json:"netmask"`
	//网关
	Gateway string `json:"gateway"`
	//业务IP资源池
	IPPool string `json:"ip_pool"`
	//带外IP资源池
	PXEPool string `json:"pxe_pool"`
	//交换机设备
	Switchs string `json:"switchs"`
	//vlan
	Vlan string `json:"vlan"`
	//IP版本
	Version string `json:"version"`
	ErrMsgContent string
}

//checkLength 对导入文件中的数据做字段长度校验
func (impDevReq *ImportIPNetworksReq) checkLength() {
	leg := len(impDevReq.ServerRoomName)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:机房管理单元长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.CIDR)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:网段(CIDR)长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.Category)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:类型长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.Netmask)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:掩码长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.Gateway)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:默认网关长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.IPPool)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:IP池长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.Switchs)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:覆盖交换机长度为(%d)(不能为空，不能大于255)", leg)
	}
	leg = len(impDevReq.Version)
	if leg == 0 || leg > 255 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("必填项校验:IP版本长度为(%d)(不能为空，不能大于255)", leg)
	}
}

//validate 对导入文件中的数据做基本验证
func (impDevReq *ImportIPNetworksReq) validate(repo model.Repo) error {
	//机房校验
	srs, err := repo.GetServerRoomByName(impDevReq.ServerRoomName)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound || srs == nil {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("机房名(%s)不存在", impDevReq.ServerRoomName)
	} else {
		impDevReq.ServerRoomID = srs.ID
		impDevReq.IDCID = srs.IDCID
	}
	// CIDR网段校验
	if strings.Index(impDevReq.CIDR, "/") == -1 {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("CIDR:%s格式错误", impDevReq.CIDR)
	} else {
		if n, err := repo.GetIPNetworks(&model.IPNetworkCond{CIDR: impDevReq.CIDR}, nil, nil); err != nil && err != gorm.ErrRecordNotFound {
			var br string
			if impDevReq.ErrMsgContent != "" {
				br = "<br />"
			}
			impDevReq.ErrMsgContent += br + fmt.Sprintf("查询CIDR:%s错误", impDevReq.CIDR)
		} else if len(n) > 0 {
			var br string
			if impDevReq.ErrMsgContent != "" {
				br = "<br />"
			}
			impDevReq.ErrMsgContent += br + fmt.Sprintf("网段：%s已存在", impDevReq.CIDR)
		}
	}
	// 网段类别校验
	if IPTypeTransate(impDevReq.Category, false) == "" {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + fmt.Sprintf("网段类型：%s不合法", impDevReq.Category)
	}
	if err := checkSwitches(strings2.MultiLines2Slice(impDevReq.Switchs), repo); err != nil {
		var br string
		if impDevReq.ErrMsgContent != "" {
			br = "<br />"
		}
		impDevReq.ErrMsgContent += br + err.Error()
	}
	//TODO
	return nil
}

//ImportIPNetworksPreview 导入预览
func ImportIPNetworksPreview(log logger.Logger, repo model.Repo, reqData *ImportPreviewReq) (map[string]interface{}, error) {
	ra, err := utils.ParseDataFromXLSX(upload.UploadDir + reqData.FileName)
	if err != nil {
		return nil, err
	}
	length := len(ra)

	var success []*ImportIPNetworksReq
	var failure []*ImportIPNetworksReq
	var errContent []string

	for i := 1; i < length; i++ {
		row := &ImportIPNetworksReq{}
		if len(ra[i]) < 10 {
			var br string
			if row.ErrMsgContent != "" {
				br = "<br />"
			}
			row.ErrMsgContent += br + "导入文件列长度应为10列"
			failure = append(failure, row)
			continue
		}
		row.ServerRoomName = ra[i][0]
		row.CIDR = ra[i][1]
		row.Category = ra[i][2]
		row.Netmask = ra[i][3]
		row.Gateway = ra[i][4]
		row.IPPool = ra[i][5]
		row.PXEPool = ra[i][6]
		row.Switchs = ra[i][7]
		row.Vlan = ra[i][8]
		row.Version = ra[i][9]

		utils.StructTrimSpace(row)

		//字段存在性校验
		row.checkLength()

		//数据有效性校验
		err := row.validate(repo)
		if err != nil {
			return nil, err
		}

		row.Category = IPTypeTransate(row.Category, false)

		if row.ErrMsgContent != "" {
			failure = append(failure, row)
			errContent = append(errContent, row.ErrMsgContent)
		} else {
			success = append(success, row)
		}
	}

	var data []*ImportIPNetworksReq
	if len(failure) > 0 {
		data = failure
	} else {
		data = success
	}
	var result []*ImportIPNetworksReq
	for i := 0; i < len(data); i++ {
		if uint(i) >= reqData.Offset && uint(i) < (reqData.Offset+reqData.Limit) {
			result = append(result, data[i])
		}
	}
	if len(errContent) > 0 {
		_ = os.Remove(upload.UploadDir + reqData.FileName)
		return map[string]interface{}{"status": "failure",
			"message":       strings.Join(errContent, "\n"),
			"total_records": len(errContent),
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

// ImportIPNetworks 将设备放到数据库
func ImportIPNetworks(log logger.Logger, repo model.Repo, conf *config.Config, reqData *ImportPreviewReq) error {
	fileName := upload.UploadDir + reqData.FileName
	ra, err := utils.ParseDataFromXLSX(fileName)
	if err != nil {
		return err
	}
	//把临时文件删了
	err = os.Remove(fileName)
	if err != nil {
		log.Warnf("ImportIPNetworks - remove tmp file: %s fail", fileName)
		return err
	}
	length := len(ra)

	for i := 1; i < length; i++ {
		row := &ImportIPNetworksReq{
			ServerRoomName: ra[i][0],
			CIDR:           ra[i][1],
			Category:       ra[i][2],
			Netmask:        ra[i][3],
			Gateway:        ra[i][4],
			IPPool:         ra[i][5],
			PXEPool:        ra[i][6],
			Switchs:        ra[i][7],
			Vlan:           ra[i][8],
			Version:		ra[i][9],
		}
		if len(ra[i]) < 9 {
			continue
		}

		//处理所有字段的多余空格字符
		utils.StructTrimSpace(row)

		//必填项校验
		row.checkLength()

		//机房和网络区域校验
		err := row.validate(repo)
		if err != nil {
			return err
		}
		sws := strings2.MultiLines2Slice(row.Switchs)
		swByte, err := json.Marshal(sws)
		if err != nil {
			log.Error("json marshal switches:%v err", row.Switchs)
		}
		ipnet := &model.IPNetwork{
			IDCID:        row.IDCID,
			ServerRoomID: row.ServerRoomID,
			Category:     IPTypeTransate(row.Category, false),
			CIDR:         row.CIDR,
			Netmask:      row.Netmask,
			Gateway:      row.Gateway,
			IPPool:       row.IPPool,
			PXEPool:      row.PXEPool,
			Switches:     string(swByte),
			Vlan:         row.Vlan,
			Version:	  row.Version,
		}

		//插入或者更新
		if _, err = repo.SaveIPNetwork(ipnet); err != nil {
			return err
		}

	}

	return nil
}

func checkSwitches(sws []string, repo model.Repo) error {
	cond := &model.NetworkDeviceCond{
		FixedAssetNumber: strings.Join(sws, ","),
	}

	res, _ := repo.GetNetworkDevicesByCond(cond, model.OrderBy{}, nil)
	if len(res) == len(sws) {
		return nil
	}
	var errSn []string

	for i := range sws {
		errSn = append(errSn, sws[i])
		for j := range res {
			if res[j].SN != sws[i] {
				continue
			}
			errSn = errSn[:len(errSn)-1]
			break
		}
	}
	if len(errSn) > 0 {
		return fmt.Errorf("交换机【%s】不存在，请检查数据或格式是否正确", strings.Join(errSn, ","))
	}
	return nil
}

// IPTypeTransate 网段类型的中英文翻译(Ch-> En by default)
func IPTypeTransate(status string, reverse bool) string {
	mStatus := map[string]string{
		"服务器ILO":   "ilo",
		"服务器TGW内网": "tgw_intranet",
		"服务器TGW外网": "tgw_extranet",
		"服务器普通内网":  "intranet",
		"服务器普通外网":  "extranet",
		"服务器虚拟化内网": "v_intranet",
		"服务器虚拟化外网": "v_extranet",
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


// AssignIPv6Req 手动分配IPv6请求结构
type AssignIPv6Req struct {
	//(required): 设备序列号
	SN string `json:"sn"`
	//(required): IP种类("extranet"/"intranet")
	Scope string `json:"scope"`
}


// FieldMap 请求参数与结构体字段建立映射
func (reqData *AssignIPv6Req) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.SN:    		"sn",
		&reqData.Scope: 		"scope",
	}
}

// Validate 结构体数据校验
func (reqData *AssignIPv6Req) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())
	log, _ := middleware.LoggerFromContext(req.Context())
	log.Debugf("sn: %v, scope: %v\n", reqData.SN, reqData.Scope)
	// 基本参数必要性检查
	if reqData.SN == "" {
		errs.Add([]string{"SN"}, binding.RequiredError, "必须指定物理机序列号")
		return errs
	}
	if reqData.Scope == "" {
		errs.Add([]string{"scope"}, binding.RequiredError, fmt.Sprintf("scope只能为: %s",
			"intranet-内网/extranet-外网"))
		return errs
	}

	//校验scope
	if reqData.Scope != "" {
		switch reqData.Scope {
		case model.IPScopeIntranet:
		case model.IPScopeExtranet:
		default:
			errs.Add([]string{"scope"}, binding.BusinessError, fmt.Sprintf("scope只能为: %s",
				"intranet-内网/extranet-外网"))
			return errs
		}
	}

	// 判断物理机是否存在
	dev, err := repo.GetDeviceBySN(reqData.SN)
	if gorm.IsRecordNotFoundError(err) {
		errs.Add([]string{"SN"}, binding.BusinessError, fmt.Sprintf("物理机(%s)不存在，请录入", reqData.SN))
		return errs
	}
	if dev.OperationStatus == model.DevOperStatInStore {
		errs.Add([]string{"SN"}, binding.BusinessError, fmt.Sprintf("%s物理机(%s)不允许直接分配IP",
			OperationStatusTransfer(dev.OperationStatus, true), reqData.SN))
		return errs
	}

	if err != nil {
		errs.Add([]string{"sns"}, binding.SystemError, "系统内部错误")
		return errs
	}

	// 校验网段是否一致
	if reqData.Scope == model.IPScopeIntranet {
		intranet, err := repo.GetIntranetIPNetworksBySN(reqData.SN)
		if gorm.IsRecordNotFoundError(err) || intranet == nil {
			errs.Add([]string{"ip_network"}, binding.BusinessError, fmt.Sprintf("物理机(%s)关联的网段存在问题(请检查对应的网络设备与网段信息是否匹配)", reqData.SN))
			return errs
		}
		if err != nil {
			errs.Add([]string{"ip"}, binding.SystemError, "系统内部错误")
			return errs
		}
		isError := true
		for k := range intranet {
			if intranet[k].Category == model.VIntranet || intranet[k].Category == model.VExtranet || intranet[k].Category == model.ILO {
				continue
			}
			if intranet[k].Version == model.IPv6 {
				isError = false
			}
		}
		if isError {
			errs.Add([]string{"ip_network"}, binding.BusinessError, fmt.Sprintf("物理机(%s)关联的内网IPv6网段不存在", reqData.SN))
			return errs
		}
	}

	if reqData.Scope == model.IPScopeExtranet {
		extranet, err := repo.GetExtranetIPNetworksBySN(reqData.SN)
		if gorm.IsRecordNotFoundError(err) || extranet == nil {
			errs.Add([]string{"ip_network"}, binding.BusinessError, fmt.Sprintf("物理机(%s)关联的网段存在问题(请检查对应的网络设备与网段信息是否匹配)", reqData.SN))
			return errs
		}
		if err != nil {
			errs.Add([]string{"ip"}, binding.SystemError, "系统内部错误")
			return errs
		}
		isError := true
		for k := range extranet {
			if extranet[k].Category == model.VIntranet || extranet[k].Category == model.VExtranet || extranet[k].Category == model.ILO {
				continue
			}
			if extranet[k].Version == model.IPv6 {
				isError = false
			}
		}
		if isError {
			errs.Add([]string{"ip"}, binding.BusinessError, fmt.Sprintf("物理机(%s)关联的外网IPv6网段不存在", reqData.SN))
			return errs
		}
	}

	return errs
}


// AssignIPv6 手动分配IPv6（根据选定的设备以及网段，自动计算并分配一个IPv6）
func AssignIPv6(log logger.Logger, repo model.Repo, reqData *AssignIPv6Req) (err error) {
	// 计算并分配一个ipv6
	var ipv6Assign *model.IP
	if ipnetwork, err := repo.GetIPv6NetworkBySN(reqData.SN, reqData.Scope); err != nil {
		log.Error(err)
		return err
	} else {
		log.Debugf("GetIPv6NetworkBySN ipv6:%+v", ipnetwork)
		var scope string
		if strings.Contains(ipnetwork.Category, "intranet") {
			scope = model.IPScopeIntranet
		} else if strings.Contains(ipnetwork.Category, "extranet") {
			scope = model.IPScopeExtranet
		}

		// 根据网段ID获取可分配的空闲IPv6
		ipv6Assign, err = repo.GetAvailableIPByIPNetworkID(ipnetwork.ID)
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			log.Error(err)
			return err
		}
		log.Debugf("GetAvailableIPByIPNetworkID ipv6:%+v", ipv6Assign)
		// 无空闲IPv6时，获取该网段最后一个IP并计算下一个IPv6地址
		if ipv6Assign == nil {
			log.Debugf("No available IPv6 of CIDR: %s .  Will generate one.", ipnetwork.CIDR)
			if ipv6Last, err := repo.GetLastIPv6ByIPNetworkID(ipnetwork.ID);err != nil && !gorm.IsRecordNotFoundError(err){
				log.Error(err)
				return err
			} else if ipv6Last != nil { // 最后一个IP不为空则获取下一个
				log.Debugf("Last IPv6 of CIDR: %s is %+v", ipnetwork.CIDR, ipv6Last)
				if ipv6Next, err := network.GetNextIPv6OfCIDR(ipv6Last.IP, ipnetwork.CIDR); err != nil {
					log.Error(err)
					return err
				} else if ipv6Next != "" {
					log.Debugf("Next IPv6 of CIDR: %s is %s", ipnetwork.CIDR, ipv6Next)
					ipv6Assign = &model.IP {
						IP: 			ipv6Next,
						IPNetworkID:	ipnetwork.ID,
						Scope:			&scope,
						Category:		model.BusinessIP,
						IsUsed:			model.IPNotUsed,
						ReleaseDate: 	time.Now(),
					}
					err = repo.CreateIP(ipv6Assign)
					if err != nil {
						log.Error(err)
						return err
					}
				}
			} else {  // 获取不到最后一个IP则当新网段分配处理
				log.Debugf("No IPv6 exist of CIDR: %s .  Will generate one.", ipnetwork.CIDR)
				if ipv6First, err := network.GetFirstIPv6OfCIDR(ipnetwork.CIDR); err != nil {
					log.Error(err)
					return err
				} else if ipv6First != "" {
					log.Debugf("New IPv6 of CIDR: %s is %s", ipnetwork.CIDR, ipv6First)
					ipv6Assign = &model.IP {
						IPNetworkID:	ipnetwork.ID,
						IP: 			ipv6First,
						Scope:			&scope,
						Category:		model.BusinessIP,
						IsUsed:			model.IPNotUsed,
						ReleaseDate: 	time.Now(),
					}
					err = repo.CreateIP(ipv6Assign)
					if err != nil {
						log.Error(err)
						return err
					}

				}
 			}
		}
	}

	if ipv6Assign != nil {
		log.Debugf("Get available ipv6:%+v", ipv6Assign)
	} else {
		log.Error("Failed to get available ipv6.exit.")
		return errors.New("Failed to get available ipv6.")
	}

	//var ds model.DeviceSetting
	ds, err := repo.GetDeviceSettingBySN(reqData.SN)
	if ds == nil {
		ds = &model.DeviceSetting{
			SN:             reqData.SN,
			NeedExtranetIP: model.NO,
			InstallType:    model.InstallationPXE,
			Status: 		model.InstallStatusSucc,
		}
	}
	//ds.SN = reqData.SN
	if reqData.Scope == model.IPScopeIntranet {
		ds.IntranetIPv6NetworkID = ipv6Assign.IPNetworkID
		ds.IntranetIPv6 = addIP(ds.IntranetIPv6, ipv6Assign.IP)
		ds.NeedIntranetIPv6 = model.YES
	} else if reqData.Scope == model.IPScopeExtranet {
		ds.ExtranetIPv6NetworkID = ipv6Assign.IPNetworkID
		ds.ExtranetIPv6 = addIP(ds.ExtranetIP, ipv6Assign.IP)
		ds.NeedExtranetIPv6 = model.YES
	}
	err = repo.SaveDeviceSetting(ds)
	if err != nil {
		log.Error(err)
		return err
	}

	//顺序放到最后，前面都成功了才占用
	err = repo.AssignIPByIP(reqData.SN, reqData.Scope, ipv6Assign.IP)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}


// AssignIPv4Req 手动为设备分配IPv4请求结构
type AssignIPv4Req struct {
	//(required): 设备序列号
	SN string `json:"sn"`
	//(required): IP种类("extranet"/"intranet")
	Scope string `json:"scope"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *AssignIPv4Req) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.SN:    "sn",
		&reqData.Scope: "scope",
	}
}

// Validate 结构体数据校验
func (reqData *AssignIPv4Req) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())
	log, _ := middleware.LoggerFromContext(req.Context())
	log.Debugf("sn: %v, scope: %v\n", reqData.SN, reqData.Scope)
	// 基本参数必要性检查
	if reqData.SN == "" {
		errs.Add([]string{"SN"}, binding.RequiredError, "必须指定物理机序列号")
		return errs
	}
	if reqData.Scope == "" {
		errs.Add([]string{"scope"}, binding.RequiredError, fmt.Sprintf("scope只能为: %s",
			"intranet-内网/extranet-外网"))
		return errs
	}
	//校验scope
	if reqData.Scope != "" {
		switch reqData.Scope {
		case model.IPScopeIntranet:
		case model.IPScopeExtranet:
		default:
			errs.Add([]string{"scope"}, binding.BusinessError, fmt.Sprintf("scope只能为: %s",
				"intranet-内网/extranet-外网"))
			return errs
		}
	}

	// 判断物理机是否存在
	dev, err := repo.GetDeviceBySN(reqData.SN)
	if gorm.IsRecordNotFoundError(err) {
		errs.Add([]string{"SN"}, binding.BusinessError, fmt.Sprintf("物理机(%s)不存在，请录入", reqData.SN))
		return errs
	}
	if dev.OperationStatus == model.DevOperStatInStore {
		errs.Add([]string{"SN"}, binding.BusinessError, fmt.Sprintf("%s物理机(%s)不允许直接分配IP",
			OperationStatusTransfer(dev.OperationStatus, true), reqData.SN))
		return errs
	}

	if err != nil {
		errs.Add([]string{"sns"}, binding.SystemError, "系统内部错误")
		return errs
	}

	// 校验网段是否一致
	if reqData.Scope == model.IPScopeIntranet {
		intranet, err := repo.GetIntranetIPNetworksBySN(reqData.SN)
		if gorm.IsRecordNotFoundError(err) || intranet == nil {
			errs.Add([]string{"ip_network"}, binding.BusinessError, fmt.Sprintf("物理机(%s)关联的网段存在问题(请检查对应的网络设备与网段信息是否匹配)", reqData.SN))
			return errs
		}
		if err != nil {
			errs.Add([]string{"ip"}, binding.SystemError, "系统内部错误")
			return errs
		}
		isError := true
		for k := range intranet {
			if intranet[k].Category == model.VIntranet || intranet[k].Category == model.VExtranet || intranet[k].Category == model.ILO {
				continue
			}
			if intranet[k].Version == model.IPv4 {
				isError = false
			}
		}
		if isError {
			errs.Add([]string{"ip_network"}, binding.BusinessError, fmt.Sprintf("物理机(%s)关联的内网IPv4网段不存在", reqData.SN))
			return errs
		}
	}

	if reqData.Scope == model.IPScopeExtranet {
		extranet, err := repo.GetExtranetIPNetworksBySN(reqData.SN)
		if gorm.IsRecordNotFoundError(err) || extranet == nil {
			errs.Add([]string{"ip_network"}, binding.BusinessError, fmt.Sprintf("物理机(%s)关联的网段存在问题(请检查对应的网络设备与网段信息是否匹配)", reqData.SN))
			return errs
		}
		if err != nil {
			errs.Add([]string{"ip"}, binding.SystemError, "系统内部错误")
			return errs
		}
		isError := true
		for k := range extranet {
			if extranet[k].Category == model.VIntranet || extranet[k].Category == model.VExtranet || extranet[k].Category == model.ILO {
				continue
			}
			if extranet[k].Version == model.IPv4 {
				isError = false
			}
		}
		if isError {
			errs.Add([]string{"ip"}, binding.BusinessError, fmt.Sprintf("物理机(%s)关联的外网IPv4网段不存在", reqData.SN))
			return errs
		}
	}

	return errs
}

// AssignIPv4 手动分配IPv4
func AssignIPv4(log logger.Logger, repo model.Repo, reqData *AssignIPv4Req) (err error) {
	var ipv4Assign *model.IP
	//var ds model.DeviceSetting
	ds, err := repo.GetDeviceSettingBySN(reqData.SN)
	if ds == nil {
		ds = &model.DeviceSetting{
			SN:             reqData.SN,
			NeedExtranetIP: model.NO,
			NeedIntranetIPv6: model.NO,
			NeedExtranetIPv6: model.NO,
			InstallType:    model.InstallationPXE, 	 //这个值固定
			Status: 		model.InstallStatusSucc, //这个值固定
		}
	}
	// 获取SN对应的网段
	if reqData.Scope == model.IPScopeIntranet {
		intranet, err := repo.GetIntranetIPNetworksBySN(reqData.SN)
		if gorm.IsRecordNotFoundError(err) || intranet == nil {
			log.Error(err)
			return err
		}
		if err != nil {
			log.Error(err)
			return err
		}
		for k := range intranet {
			if intranet[k].Category == model.Intranet && intranet[k].Version == model.IPv4 {
				ipv4Assign, err = repo.GetAvailableIPByIPNetworkID(intranet[k].ID)
				if err != nil && !gorm.IsRecordNotFoundError(err) {
					log.Error(err)
					return err
				}
				log.Debugf("GetAvailableIPByIPNetworkID Intranet ipv4:%+v", ipv4Assign)
				break
			}
		}
	}
	if reqData.Scope == model.IPScopeExtranet {
		extranet, err := repo.GetExtranetIPNetworksBySN(reqData.SN)
		if gorm.IsRecordNotFoundError(err) || extranet == nil {
			log.Error(err)
			return err
		}
		if err != nil {
			log.Error(err)
			return err
		}
		for k := range extranet {
			if extranet[k].Category == model.Extranet && extranet[k].Version == model.IPv4 {
				ipv4Assign, err = repo.GetAvailableIPByIPNetworkID(extranet[k].ID)
				if err != nil && !gorm.IsRecordNotFoundError(err) {
					log.Error(err)
					return err
				}
				log.Debugf("GetAvailableIPByIPNetworkID Extranet ipv4:%+v", ipv4Assign)
				break
			}
		}
	}	

	//先ping检查，如果IP可以ping通，则失败
	if err = ping.PingTest(ipv4Assign.IP); err == nil {
		return fmt.Errorf("IP:%s连通性(ping)测试失败", ipv4Assign.IP)
	}

	//ds.SN = reqData.SN
	if reqData.Scope == model.IPScopeIntranet {
		//关联的网段唯一
		ds.IntranetIPNetworkID = ipv4Assign.IPNetworkID
		ds.IntranetIP = addIP(ds.IntranetIP, ipv4Assign.IP)
	} else if reqData.Scope == model.IPScopeExtranet {
		ds.ExtranetIPNetworkID = ipv4Assign.IPNetworkID
		ds.ExtranetIP = addIP(ds.ExtranetIP, ipv4Assign.IP)
		ds.NeedExtranetIP = model.YES
	}
	err = repo.SaveDeviceSetting(ds)
	if err != nil {
		log.Error(err)
		return err
	}

	//顺序放到最后，前面都成功了才占用
	err = repo.AssignIP(reqData.SN, reqData.Scope, uint(ipv4Assign.ID))
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}