package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/utils/win"
)

// ErrRoutesNotFound route信息不存在
var ErrRoutesNotFound = errors.New("routes not found")

func writeCMDs(w io.Writer, cmds []string) (n int64, err error) {
	var buf bytes.Buffer
	for i := range cmds {
		buf.WriteString(cmds[i])
		buf.WriteByte('\r')
		buf.WriteByte('\n')
	}
	return buf.WriteTo(w)
}

// ApplyRoutesCMDs 返回配置路由的命令切片
func ApplyRoutesCMDs(log logger.Logger, routesFile, masterMacAddr, slaveMacAddr string) (cmds []string, err error) {
	routes, err := loadRoutesFromFile(log, routesFile)
	if err != nil {
		return nil, err
	}

	if len(routes) == 0 {
		log.Errorf("Routes not found in file of %s", routesFile)
		return nil, ErrRoutesNotFound
	}

	var masterIFID, slaveIFID string
	masterIFID, err = getInterfaceID(log, masterMacAddr)
	if err != nil {
		return nil, err
	}

	if slaveMacAddr != "" {
		slaveIFID, err = getInterfaceID(log, slaveMacAddr)
		if err != nil {
			return nil, err
		}
	}
	log.Infof("InterfaceID4master: %s, InterfaceID4slave: %s", masterIFID, slaveIFID)
	switch len(routes) {
	case 1: //仅配置主IP路由
		routes[0].InterfaceID = masterIFID
		return []string{
			fmt.Sprintf("route -p ADD %s MASK %s %s  IF %s", routes[0].Dest, routes[0].Mask, routes[0].Gateway, routes[0].InterfaceID),
		}, nil
	default: // 配置主/从IP路由
		routes[0].InterfaceID = masterIFID
		routes[1].InterfaceID = slaveIFID
		return []string{
			fmt.Sprintf("route -p ADD %s MASK %s %s  IF %s", routes[0].Dest, routes[0].Mask, routes[0].Gateway, routes[0].InterfaceID),
			fmt.Sprintf("route -p ADD %s MASK %s %s  IF %s", routes[1].Dest, routes[1].Mask, routes[1].Gateway, routes[1].InterfaceID),
		}, nil
	}
}

// ApplyRoutes 将路由表应用到本机
func ApplyRoutes(log logger.Logger, routesFile, masterMacAddr, slaveMacAddr string) (err error) {
	log.Infof("Start applying routes(macAddr4master: %q, macAddr4slave: %q)", masterMacAddr, slaveMacAddr)
	routes, err := loadRoutesFromFile(log, routesFile)
	if err != nil {
		return err
	}

	if len(routes) == 0 {
		log.Errorf("Routes not found in file of %s", routesFile)
		return ErrRoutesNotFound
	}

	var masterIFID, slaveIFID string
	masterIFID, err = getInterfaceID(log, masterMacAddr)
	if err != nil {
		return err
	}

	if slaveMacAddr != "" {
		slaveIFID, err = getInterfaceID(log, slaveMacAddr)
		if err != nil {
			return err
		}
	}
	log.Infof("InterfaceID4master: %s, InterfaceID4slave: %s", masterIFID, slaveIFID)
	switch len(routes) {
	case 1: //仅配置主IP路由
		routes[0].InterfaceID = masterIFID
		return addRoute(log, &routes[0])
	default: // 配置主/从IP路由
		routes[0].InterfaceID = masterIFID
		routes[1].InterfaceID = slaveIFID
		if err = addRoute(log, &routes[0]); err != nil {
			return err
		}
		return addRoute(log, &routes[1])
	}
}

// addRoute 添加一条永久路由
func addRoute(log logger.Logger, route *Route) (err error) {
	if route == nil || route.InterfaceID == "" {
		log.Infof("Discard add route operation: %#v", route)
		return nil
	}
	cmdAndArgs := fmt.Sprintf("route -p ADD %s MASK %s %s  IF %s", route.Dest, route.Mask, route.Gateway, route.InterfaceID)
	_, err = win.ExecOutputWithLog(log, cmdAndArgs)
	return err
}

// Route 路由信息
type Route struct {
	Dest        string `json:"dest"`    // 目标网段
	Mask        string `json:"mask"`    // 子网掩码
	Gateway     string `json:"gateway"` // 网关地址
	InterfaceID string `json:"-"`       // 网络接口ID
}

func loadRoutesFromFile(log logger.Logger, srcFile string) (routes []Route, err error) {
	b, err := ioutil.ReadFile(srcFile)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var wapper struct {
		Content []Route
	}

	if err = json.Unmarshal(b, &wapper); err != nil {
		log.Error(err)
		return nil, err
	}
	return wapper.Content, nil
}

func getInterfaceID(log logger.Logger, macAddr string) (id string, err error) {
	outUTF8, err := win.ExecOutputWithLog(log, `route print interface`)
	if err != nil {
		return "", err
	}
	return parseInterfaceID(log, outUTF8, macAddr)
}

func parseInterfaceID(log logger.Logger, outUTF8 []byte, macAddr string) (id string, err error) {
	addr := strings.ToLower(strings.Replace(macAddr, ":", " ", -1))
	rd := bufio.NewReader(bytes.NewBuffer(outUTF8))
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error(err)
			return "", err
		}
		line = strings.ToLower(strings.TrimSpace(line))
		if strings.Contains(line, addr) {
			if idx := strings.Index(line, "..."); idx > 0 {
				return strings.TrimSpace(line[:idx]), nil
			}
		}
	}
	return "", nil
}
