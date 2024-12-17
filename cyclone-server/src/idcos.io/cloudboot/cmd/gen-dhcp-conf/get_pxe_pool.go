package main

import (
	"strconv"
	//"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//"net/url"
	"os"
	"strings"
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
	//vlan
	Vlan string `json:"vlan"`
	// ipv4 ipv6
	Version string `json:"version"`
}

// Content 请求网段信息返回内容结构
type Content struct {
	Page         int             `json:"page"`
	PageSize     int             `json:"page_size"`
	Records      []IPNetworkPage `json:"records"`
	TotalPages   int             `json:"total_pages"`
	TotalRecords int             `json:"total_records"`
}

// IPNetworkResp 请求网段信息返回信息结构
type IPNetworkResp struct {
	Status  string  `json:"status"`
	Message string  `json:"message"`
	Content Content `json:"content"`
}

// GetIPNetworkPageReq 查询网段分页请求结构体
type GetIPNetworkPageReq struct {
	//机房
	ServerRoomName string `json:"server_room_name"`
	//CIDR网段名称
	CIDR string `json:"cidr"`
	Switches string `json:"switches"`
	//网段类别。可选值: ilo-服务器ILO; tgw_intranet-服务器TGW内网; tgw_extranet-服务器TGW外网; intranet-服务器普通内网; extranet-服务器普通外网; v_intranet-服务器虚拟化内网;
	Category string `json:"category"`
	//分页页号
	Page int64 `json:"page"`
	//分页大小
	PageSize int64 `json:"page_size"`
}

var (
	token = ""
)

func getEnv() string {
	dhcpToken := strings.TrimSpace(os.Getenv("DHCP_TOKEN"))
	if dhcpToken != "" {
		token = dhcpToken
	}
	return token
}

// GetIPNetworks 返回满足过滤条件的网段
func GetIPNetworks(srname, cidr, switches, category, port, ip string) (items []IPNetworkPage, err error) {
	//if !checkCategory(category) {
	//	return nil, nil
	//}
	var data []byte
	var ipnr IPNetworkResp
	url := fmt.Sprintf("http://%s:%s/api/cloudboot/v1/ip-networks", ip, port)
	if data, err = DoGet(url, &GetIPNetworkPageReq{
		Category: category,
		ServerRoomName: srname,
		CIDR: cidr,
		Switches: switches,
		Page: 1,
		PageSize: 1000,
	}); err == nil {
		err = json.Unmarshal(data, &ipnr)
		if err != nil {
			log.Printf("JSON解析请求数据出错: %s\n", err.Error())
		}
	} else {
		log.Printf("获取网段信息出错: %s\n", err.Error())
	}

	return ipnr.Content.Records, nil
}

// DoGet  发送get请求
func DoGet(myurl string, reqData *GetIPNetworkPageReq) ([]byte, error) {
	log.Printf("GET %s, reqData: %v", myurl, reqData)
	req, err := http.NewRequest("GET", myurl, nil)
	if err != nil {
		return nil, err
	}
	// 拼接参数
	q := req.URL.Query()
	q.Add("page",strconv.FormatInt(reqData.Page,10))
	q.Add("page_size",strconv.FormatInt(reqData.PageSize,10))
	if reqData.ServerRoomName != "" {
		q.Add("server_room_name", reqData.ServerRoomName)
	}
	if reqData.CIDR != "" {
		q.Add("cidr", reqData.CIDR)
	}
	if reqData.Switches != "" {
		q.Add("switches", reqData.Switches)
	}	
	q.Add("category", reqData.Category)
	req.URL.RawQuery = q.Encode()
	log.Printf("GET %s, raw quey: %s", myurl, req.URL.RawQuery)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", getEnv())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	log.Printf("GET %s, response body: %s", myurl, respBody)

	if resp.StatusCode != http.StatusOK {
		log.Printf("GET %s, response status code: %d", myurl, resp.StatusCode)
		return nil, fmt.Errorf("http status code: %d", resp.StatusCode)
	}
	return respBody, nil
}
