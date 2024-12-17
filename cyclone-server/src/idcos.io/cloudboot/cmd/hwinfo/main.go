package main

import (
	"bytes"
	"strings"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"fmt"
	"os"

	"github.com/urfave/cli"
	"idcos.io/cloudboot/hardware/collector"
	"idcos.io/cloudboot/server/cloudbootserver/types/device"
	"idcos.io/cloudboot/hardware/raid"

	_ "idcos.io/cloudboot/hardware/collector/bootos"
	_ "idcos.io/cloudboot/hardware/oob/ipmi"
	_ "idcos.io/cloudboot/hardware/raid/adaptecsmartraid"
	_ "idcos.io/cloudboot/hardware/raid/avagomegaraid"
	_ "idcos.io/cloudboot/hardware/raid/hpsmartarray"
	_ "idcos.io/cloudboot/hardware/raid/lsimegaraid"
	_ "idcos.io/cloudboot/hardware/raid/lsisas2"
	_ "idcos.io/cloudboot/hardware/raid/lsisas3"
)


var (
	name         = "hwinfo"
	usage        = "for testing hwinfo collect and post it."
	postServer   = ""
	version      = "1.1.0" // 2023-02-09 新增上报采集数据功能(可选)
)


func main() {
	app := cli.NewApp()
	app.Name = name
	app.Usage = usage
	app.Version = version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "p, postserver",
			Value: postServer,
			Usage: "指定上报采集数据到服务端的IP (默认：仅测试采集不上报数据)",
		},
	}
	app.Action = func(ctx *cli.Context) error {
		postServer = ctx.String("p")
		return Main(postServer)
	}

	app.Run(os.Args)
}


func Main(pServer string) error {
	fmt.Println("===========本程序仅用于设备信息采集模块的自测及上报更新===========")

	c := collector.SelectCollector(collector.DefaultCollector)
	if c == nil {
		fmt.Fprintf(os.Stderr, "Unregistered collector: %s\n", collector.DefaultCollector)
		return collector.ErrUnregisteredCollector
	}

	var dev device.Device

	fmt.Printf("\n===========BASE===========\n")
	isVM, sn, manufacturer, modelName, arch, err := c.BASE()
	//_, dev.SN, dev.Vendor, dev.Model, dev.Arch, _ = c.BASE()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Collect Base error: %s\n", err.Error())
	}
	fmt.Println("isVM: ", isVM)
	fmt.Println("SN: ", sn)
	fmt.Println("Manufacturer: ", manufacturer)
	fmt.Println("Model: ", modelName)
	fmt.Println("Arch: ", arch)
	dev.SN = sn
	dev.Vendor = manufacturer
	dev.Model = modelName
	dev.Arch = arch

	fmt.Printf("\n===========CPU===========\n")
	dev.CPU, err = c.CPU()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Collect CPU error: %s\n", err.Error())
	}
	_ = json.NewEncoder(os.Stdout).Encode(dev.CPU)

	fmt.Printf("\n===========Memory===========\n")
	dev.Memory, err = c.Memory()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Collect memory error: %s\n", err.Error())
	}
	_ = json.NewEncoder(os.Stdout).Encode(dev.Memory)

	fmt.Printf("\n===========Motherboard===========\n")
	dev.Motherboard, err = c.Motherboard()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Collect motherboard error: %s\n", err.Error())
	}
	_ = json.NewEncoder(os.Stdout).Encode(dev.Motherboard)

	fmt.Printf("\n===========Disk===========\n")
	dev.Disk, err = c.Disk()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Collect logical disk error: %s\n", err.Error())
	}
	_ = json.NewEncoder(os.Stdout).Encode(dev.Disk)

	fmt.Printf("\n===========Disk Slot===========\n")
	name, err := raid.Whoami()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Get crontroller by lspci error: %s\n", err.Error())
	}
	fmt.Println("Controller(lscpi): ", name)
	worker := raid.SelectWorker(name)
	if worker == nil {
		fmt.Fprintf(os.Stderr, "Get worker of crontroller error: %s\n", err.Error())
	} else {
		fmt.Println("Controller using cmdline: ", worker.GetCMDLine())
	}
	dev.DiskSlot, err = c.DiskSlot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Collect disk slot error: %s\n", err.Error())
	}
	_ = json.NewEncoder(os.Stdout).Encode(dev.DiskSlot)

	fmt.Printf("\n===========NIC===========\n")
	dev.NIC, err = c.NIC()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Collect nic error: %s\n", err.Error())
	}
	_ = json.NewEncoder(os.Stdout).Encode(dev.NIC)

	fmt.Printf("\n===========OOB===========\n")
	dev.OOB, err = c.OOB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Collect oob error: %s\n", err.Error())
	}
	_ = json.NewEncoder(os.Stdout).Encode(dev.OOB)

	fmt.Printf("\n===========BIOS===========\n")
	dev.BIOS, err = c.BIOS()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Collect bios error: %s\n", err.Error())
	}
	_ = json.NewEncoder(os.Stdout).Encode(dev.BIOS)

	fmt.Printf("\n===========LLDP===========\n")
	fmt.Println("LLDP using cmdline: lldpctl -f json")
	dev.LLDP, err = c.LLDP()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Collect lldp error: %s\n", err.Error())
	}
	_ = json.NewEncoder(os.Stdout).Encode(dev.LLDP)

	fmt.Printf("\n===========PCI===========\n")
	dev.PCI, err = c.PCI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Collect pci error: %s\n", err.Error())
	}
	_ = json.NewEncoder(os.Stdout).Encode(dev.PCI)

	fmt.Printf("\n===========Fan===========\n")
	dev.Fan, err = c.Fan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Collect fan error: %s\n", err.Error())
	}
	_ = json.NewEncoder(os.Stdout).Encode(dev.Fan)

	fmt.Printf("\n===========RAID===========\n")
	dev.RAID, err = c.RAID()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Collect raid error: %s\n", err.Error())
	}
	_ = json.NewEncoder(os.Stdout).Encode(dev.RAID)

	fmt.Printf("\n===========Power===========\n")
	dev.Power, err = c.Power()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Collect power error: %s\n", err.Error())
	}
	_ = json.NewEncoder(os.Stdout).Encode(dev.Power)

	fmt.Printf("\n===========iDRAC===========\n")
	idrac, err := c.IDRAC()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Collect idrac error: %s\n", err.Error())
	}
	_ = json.NewEncoder(os.Stdout).Encode(idrac)

	fmt.Printf("\n===========iLO===========\n")
	ilo, err := c.ILO()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Collect ilo error: %s\n", err.Error())
	}
	_ = json.NewEncoder(os.Stdout).Encode(ilo)

	fmt.Printf("\n===========HBA===========\n")
	hba, err := c.HBA()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Collect hba error: %s\n", err.Error())
	}
	if hba != nil {
		_ = json.NewEncoder(os.Stdout).Encode(hba)
		dev.HBA = hba
	} else {
		fmt.Println("HBA not found.")
	}

	fmt.Printf("\n===========PostCollectedData===========\n")
	if pServer != "" {
		dev.Setup()
		url := fmt.Sprintf("http://%s/api/cloudboot/v1/devices/collections?lang=en-US", pServer)
		fmt.Println("PostURL: ", url)
		var respData struct {
			Status  string
			Message string
		}
		if err = doPOSTUnmarshal(url, &dev, &respData); err != nil {
			return err
		}
		if strings.ToLower(respData.Status) != "success" {
			fmt.Printf("status: %s, message: %s", respData.Status, respData.Message)
			return fmt.Errorf("status: %s, message: %s", respData.Status, respData.Message)
		}
		fmt.Printf("status: %s, message: %s", respData.Status, respData.Message)
		return nil
	} else {
		fmt.Println("PostServer not specify.Do not post data.")
	}
	return nil
}

// doPOSTUnmarshal 将指定数据序列化成JSON并通过HTTP POST发送到远端，并将JSON格式的响应信息反序列化到respData中。
// respData必须为非nil的指针类型，否则将返回对应的错误。
func doPOSTUnmarshal(url string, reqData, respData interface{}) error {
	respBody, err := doPOST(url, reqData)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(respBody, respData); err != nil {
		fmt.Printf("json unmarshal error: %s", err.Error())
		return err
	}
	return nil
}

// doPOST 将指定数据序列化成JSON并通过HTTP POST发送到远端
func doPOST(url string, reqData interface{}) ([]byte, error) {
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		fmt.Printf("json marshal error: %s", err.Error())
		return nil, err
	}

	fmt.Printf("POST %s, request body: %s", url, reqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("json marshal error: %s", err.Error())
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("json marshal error: %s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("POST %s, response body: %s", url, respBody)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("POST %s, response status code: %d", url, resp.StatusCode)
		return nil, fmt.Errorf("http status code: %d", resp.StatusCode)
	}
	return respBody, nil
}