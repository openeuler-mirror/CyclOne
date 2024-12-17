package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
)

// ComponentLog 组件日志
type ComponentLog struct {
	SN        string
	Component string
	LogData   []byte
}

// SaveComponentLogReq 保存设备组件日志请求结构体
type SaveComponentLogReq struct {
	SN           string
	Component    string
	LogData      []byte
	DataPath     string
	OriginNode   string
	OriginNodeIP string
}

// SaveComponentLog 保存设备组件日志
func SaveComponentLog(repo model.Repo, log logger.Logger, reqData *SaveComponentLogReq) (err error) {
	// 如果为proxy下的机器，则日志存储目录为 rootpath/proxy_ip/log/${SN}/
	if reqData.OriginNodeIP != "" {
		reqData.OriginNodeIP = strings.Replace(reqData.OriginNodeIP, ".", "_", -1)
		reqData.DataPath = filepath.Join(reqData.DataPath, "proxy_"+reqData.OriginNodeIP)
	}

	// 创建目标设备的日志持久化目录
	fp := filepath.Join(reqData.DataPath, "log", reqData.SN)
	err = os.MkdirAll(fp, 0755)
	if err != nil {
		log.Errorf("%s\n", err)
		return err
	}

	// 日志文件命名格式：${组件名}_${时间戳}.log
	filename := fmt.Sprintf("%s_%s.log", reqData.Component, time.Now().Format("200601021504"))
	err = ioutil.WriteFile(fp+"/"+filename, reqData.LogData, 0755)
	if err != nil {
		log.Errorf("%s\n", err)
		return err
	}

	return nil
}
