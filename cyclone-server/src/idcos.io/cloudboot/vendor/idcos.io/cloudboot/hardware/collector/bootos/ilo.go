package bootos

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"

	"idcos.io/cloudboot/hardware/collector"
)

const (
	// iLOTool iLO配置管理工具
	iLOTool = "hponcfg"

	// getFirmwareVersion4iLO 请求查询iLO固件版本的XML
	getFirmwareVersion4iLO = `
	<?xml version="1.0"?>
	<RIBCL VERSION="2.0">
        <LOGIN USER_LOGIN="admin" PASSWORD="Cyclone@1234">
            <RIB_INFO MODE="read">
            	<GET_FW_VERSION/>
			</RIB_INFO>
        </LOGIN>
    </RIBCL>
	`
)

var (
	iLOReqXML  = filepath.Join(os.TempDir(), "ilo_req.xml")
	iLORespXML = filepath.Join(os.TempDir(), "ilo_resp.xml")
)

type fwVersionResp struct {
	XMLName xml.Name `xml:"GET_FW_VERSION"`
	Version string   `xml:"FIRMWARE_VERSION,attr"`
	Date    string   `xml:"FIRMWARE_DATE,attr"`
}

// ILO 采集并返回HP iLO信息。
func (c *bootosC) ILO() (ilo *collector.ILO, err error) {
	if err = ioutil.WriteFile(iLOReqXML, []byte(getFirmwareVersion4iLO), 0644); err != nil {
		if log := c.Base.GetLog(); log != nil {
			log.Error(err)
		}
		return nil, err
	}

	if _, err = c.Base.ExecByShell(iLOTool, "-f", iLOReqXML, "-l", iLORespXML); err != nil {
		return nil, err
	}

	bResp, err := ioutil.ReadFile(iLORespXML)
	if err != nil {
		if log := c.Base.GetLog(); log != nil {
			log.Error(err)
		}
		return nil, err
	}
	var resp fwVersionResp
	if err = xml.Unmarshal(bResp, &resp); err != nil {
		if log := c.Base.GetLog(); log != nil {
			log.Error(err)
		}
		return nil, err
	}
	return &collector.ILO{
		FirmwareDate:    resp.Date,
		FirmwareVersion: resp.Version,
	}, nil
}
