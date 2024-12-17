package raid

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"

	"idcos.io/cloudboot/hardware"
	"idcos.io/cloudboot/logger"
	"strconv"
)

const (
	// AvagoMegaRAID RAID实现-Avago MegaRAID
	AvagoMegaRAID = "AVAGO_MEGA_RAID"
	// LSIMegaRAID RAID实现-LSI MegaRAID
	LSIMegaRAID = "LSI_MEGA_RAID"
	// HPSmartArray RAID实现-HP SmartArray
	HPSmartArray = "HP_SMART_ARRAY"
	// HPESmartArray RAID实现-HPE SmartArray
	HPESmartArray = "HPE_SMART_ARRAY"
	// LSISAS2 RAID实现-LSI SAS2
	LSISAS2 = "LSI_SAS2"
	// LSISAS3 RAID实现-LSI SAS3
	LSISAS3 = "LSI_SAS3"
	// AdaptecSmartRAID RAID实现-Adaptec Smart Storage
	AdaptecSmartRAID = "ADAPTEC_SMART_RAID"

)

// Level RAID级别
type Level string

var (
	// RAID0 RAID级别-RAID0
	RAID0 Level = "raid0"
	// RAID1 RAID级别-RAID1
	RAID1 Level = "raid1"
	// RAID5 RAID级别-RAID5
	RAID5 Level = "raid5"
	// RAID6 RAID级别-RAID6
	RAID6 Level = "raid6"
	// RAID10 RAID级别-RAID10
	RAID10 Level = "raid10"
	// RAID1E RAID级别-RAID1E
	RAID1E Level = "raid1e"
	// RAID50 RAID级别-RAID50
	RAID50 Level = "raid50"
	// RAID60 RAID级别-RAID60
	RAID60 Level = "raid60"
)

// MustSelectLevel 返回指定的RAID级别
func MustSelectLevel(level string) Level {
	l, err := SelectLevel(level)
	if err != nil {
		panic(err.Error())
	}
	return l
}

// SelectLevel 返回指定的RAID级别
func SelectLevel(level string) (Level, error) {
	switch strings.ToLower(level) {
	case string(RAID0):
		return RAID0, nil
	case string(RAID1):
		return RAID1, nil
	case string(RAID10):
		return RAID10, nil
	case string(RAID1E):
		return RAID1E, nil
	case string(RAID5):
		return RAID5, nil
	case string(RAID50):
		return RAID50, nil
	case string(RAID6):
		return RAID6, nil
	case string(RAID60):
		return RAID60, nil
	}
	return "", ErrRAIDLevel
}

// Whoami 返回当前的RAID硬件对应的处理器名
func Whoami() (worker string, err error) {
	output, err := exec.Command("lspci").Output()
	if err != nil {
		return AvagoMegaRAID,nil
	}

	if bytes.Contains(output, []byte("MegaRAID")) {
		arch, _ := exec.Command("uname", "-m").Output()
		if len(arch) > 0 && strings.Contains(string(arch), "aarch64") {
			return AvagoMegaRAID, nil
		}
		return LSIMegaRAID, nil

	} else if bytes.Contains(output, []byte("Smart Array")) {
		return HPSmartArray, nil

	} else if bytes.Contains(output, []byte("Adaptec Smart")) {
		return AdaptecSmartRAID, nil

	} else if bytes.Contains(output, []byte("SAS2004")) ||
		bytes.Contains(output, []byte("SAS2008")) ||
		bytes.Contains(output, []byte("SAS2108")) ||
		bytes.Contains(output, []byte("SAS2208")) ||
		bytes.Contains(output, []byte("SAS2304")) ||
		bytes.Contains(output, []byte("SAS2308")) {
		return LSISAS2, nil

	} else if bytes.Contains(output, []byte("SAS3004")) || bytes.Contains(output, []byte("SAS3008")) {
		return LSISAS3, nil

	}
	//给个默认值
	return AvagoMegaRAID,nil
}

var (
	// ErrRAIDLevel 无效的RAID级别
	ErrRAIDLevel = errors.New("error RAID level")
	// ErrUnknownHardware 未知的RAID硬件
	ErrUnknownHardware = errors.New("unknown RAID hardware")
	// ErrControllerNotFound RAID controller未发现
	ErrControllerNotFound = errors.New("controller not found")
	// ErrPhysicalDriveNotFound 指定的物理驱动器未发现
	ErrPhysicalDriveNotFound = errors.New("specified physical drive not found")
	// ErrUnsupportedLevel 不支持的RAID级别
	ErrUnsupportedLevel = errors.New("unsupported RAID level")
	// ErrInitLogicalDrivesNotSupport 暂不支持初始化逻辑磁盘
	ErrInitLogicalDrivesNotSupport = errors.New("initializing logical drives is not supported")
	// ErrJBODModeNotSupport 不支持JBOD模式
	ErrJBODModeNotSupport = errors.New("JBOD mode is not supported")
	// ErrUnsupportedRAIDManagingMode 暂不支持的RAID管理模式
	ErrUnsupportedRAIDManagingMode = errors.New("unsupported RAID managing mode")
	// ErrInvalidDiskIdentity 无效的硬盘标识
	ErrInvalidDiskIdentity = errors.New("invalid disk identity")
)

// PhysicalDrive 物理驱动器(物理磁盘)
type PhysicalDrive struct {
	ID           string `json:"id"`
	Name         string `json:"name"` // 唯一标识
	RawSize      string `json:"raw_size"`
	MediaType    string `json:"media_type"`
	ControllerID string `json:"controller_id"` // 所属的RAID卡控制器编号
}

// Controller RAID控制器，指代一块RAID卡硬件。
type Controller struct {
	ID              string `json:"id"`               // 唯一标识/编号
	ModelName       string `json:"model_name"`       // 型号
	FirmwareVersion string `json:"firmware_version"` // 固件版本
	Mode            string `json:"mode"`             // 当前所处模式RAID/JBOD(HBA)
}

//实现按RawSize从小到大排序
type PDSlice []struct{
	Index int
	PhysicalDrive
}

func (pds PDSlice) Len() int {
	return len(pds)
}
func (pds PDSlice) Less(i, j int) bool {
	fi, fj := 0.0, 0.0
	if strings.Contains(pds[i].RawSize, "TB") {
		fi, _ = strconv.ParseFloat(strings.TrimSpace(strings.Replace(pds[i].RawSize, "TB", "", -1)), 64)
		fi *= 1000
	} else {
		fi, _ = strconv.ParseFloat(strings.TrimSpace(strings.Replace(pds[i].RawSize, "GB", "", -1)), 64)
	}
	if strings.Contains(pds[j].RawSize, "TB") {
		fj, _ = strconv.ParseFloat(strings.TrimSpace(strings.Replace(pds[j].RawSize, "TB", "", -1)), 64)
		fj *= 1000
	} else {
		fj, _ = strconv.ParseFloat(strings.TrimSpace(strings.Replace(pds[j].RawSize, "GB", "", -1)), 64)
	}
	return fi < fj
}
func (pds PDSlice) Swap(i, j int) {
	pds[i], pds[j] = pds[j], pds[i]
}

// Worker RAID处理器接口
type Worker interface {
	// SetDebug 设置是否开启debug。若开启debug，会将关键日志信息写入console。
	SetDebug(debug bool)
	// SetLog 更换日志实现。默认情况下内部无日志实现。
	SetLog(log logger.Logger)
	// Convert2ControllerID RAID卡控制器索引号转换成RAID卡控制器ID。
	// ctrlIndex RAID卡控制器索引号。0表示首块RAID卡，1表示第二块RAID卡，以此类推。
	// 如'HP SmartArray' RAID卡，若索引号为'0'，那么可能转换获得的RAID卡控制器ID为'1'，原因是这款RAID卡ID是从'1'开始的。
	// 如'LSI SAS3' RAID卡，若索引号为'0'，那么可能转换获得的RAID卡控制器ID依然为'0'，原因是这款RAID卡ID就是从'0'开始的。
	Convert2ControllerID(ctrlIndex uint) (ctrlID string, err error)
	// Controllers 返回设备的RAID卡控制器列表。
	Controllers() (ctrls []Controller, err error)
	// PhysicalDrives 返回指定RAID卡控制器的物理驱动器(物理磁盘)列表。
	// ctrlID 表示RAID卡控制器ID，通过Convert2ControllerID方法可获取到第N块RAID卡的控制器ID。
	PhysicalDrives(ctrlID string) (pds []PhysicalDrive, err error)
	// Clear 擦除指定RAID卡控制器的配置。
	// ctrlID 表示RAID卡控制器ID，通过Convert2ControllerID方法可获取到第N块RAID卡的控制器ID。
	Clear(ctrlID string) error
	// InitLogicalDrives 初始化指定RAID卡控制器下的逻辑驱动器(逻辑磁盘)。
	InitLogicalDrives(ctrlID string) error
	// TranslateLegacyDisks 将目标RAID卡控制器下传统"1|2,3|4-6|7-|all"形式的硬盘标识符转换成实际的物理驱动器。
	TranslateLegacyDisks(ctrlID string, legacyDisks string) (drives []string, err error)
	// CreateArray 在指定RAID卡控制器下的若干物理驱动器上创建指定级别的阵列(逻辑磁盘)。
	// ctrlID 表示RAID卡控制器ID，通过Convert2ControllerID方法可获取到第N块RAID卡的控制器ID。
	// level RAID级别，不同的RAID卡所支持的级别有所差异。
	// drives 物理驱动器列表。PhysicalDrives方法可获得的物理驱动器信息列表，其中'Name'属性即为某个物理驱动器标识。
	CreateArray(ctrlID string, level Level, drives []string) error
	// SetGlobalHotspares 将指定RAID卡下的物理驱动器设置为全局热备盘。
	// ctrlID 表示RAID卡控制器ID，通过Convert2ControllerID方法可获取到第N块RAID卡的控制器ID。
	SetGlobalHotspares(ctrlID string, drives []string) error
	// SetJBODs 将指定RAID卡控制器设置为直通模式，或者将RAID卡控制器下部分物理驱动器设置为直通模式。
	// ctrlID 表示RAID卡控制器ID，通过Convert2ControllerID方法可获取到第N块RAID卡的控制器ID。
	// drives 物理驱动器列表。若物理驱动器列表为空，则意味着将指定的RAID卡下所有的物理驱动器都设置为直通模式。
	SetJBODs(ctrlID string, drives []string) error
	// PostCheck RAID配置实施后置检查
	PostCheck(sett *Setting) []hardware.CheckingItem
	// GetCMDLine 返回 worker 所使用的cmdline
	GetCMDLine() string
}
