package win

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"reflect"
	"testing"
	"time"

	mybytes "idcos.io/cloudboot/utils/bytes"

	"bou.ke/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"idcos.io/cloudboot/utils"
)

func TestDiskNo(t *testing.T) {
	Convey("转换并获取磁盘编号", t, func() {
		var no int
		var err error

		no, err = DiskNo("")
		So(err, ShouldEqual, ErrMalformedDisk)
		So(no, ShouldEqual, 0)

		no, err = DiskNo("-1")
		So(err, ShouldEqual, ErrMalformedDisk)
		So(no, ShouldEqual, 0)

		no, err = DiskNo("/dev/sda99")
		So(err, ShouldEqual, ErrMalformedDisk)
		So(no, ShouldEqual, 0)

		no, err = DiskNo("0")
		So(err, ShouldBeNil)
		So(no, ShouldEqual, 0)

		no, err = DiskNo("7")
		So(err, ShouldBeNil)
		So(no, ShouldEqual, 7)

		no, err = DiskNo("/dev/sdb")
		So(err, ShouldBeNil)
		So(no, ShouldEqual, 1)

		no, err = DiskNo("/dev/hde")
		So(err, ShouldBeNil)
		So(no, ShouldEqual, 4)
	})
}

func TestCMD(t *testing.T) {
	Convey("返回windows分区的执行命令", t, func() {
		conf := &DiskPartConfiguration{
			Disk: 0,
			Partitions: []PartConfiguration{
				{
					Size:       "1024",
					FSType:     "ntfs",
					Mountpoint: "C",
				},
			},
		}
		So(string(conf.CMD()), ShouldEqual, fmt.Sprintf(
			"select disk %d\r\nclean\r\ncreate partition primary size=%s\r\nformat quick fs=%s\r\nassign letter=%s\r\nactive\r\n",
			conf.Disk,
			conf.Partitions[0].Size,
			conf.Partitions[0].FSType,
			conf.Partitions[0].Mountpoint,
		))

		conf = &DiskPartConfiguration{
			Disk: 0,
			Partitions: []PartConfiguration{
				{
					Size:       "free",
					FSType:     "ntfs",
					Mountpoint: "D",
				},
			},
		}
		So(string(conf.CMD()), ShouldEqual, fmt.Sprintf(
			"select disk %d\r\nclean\r\ncreate partition primary\r\nformat quick fs=%s\r\nassign letter=%s\r\n",
			conf.Disk,
			conf.Partitions[0].FSType,
			conf.Partitions[0].Mountpoint,
		))

		conf = &DiskPartConfiguration{
			Disk: 1,
			GPT:  true,
			Partitions: []PartConfiguration{
				{
					Size:       "free",
					FSType:     "ntfs",
					Mountpoint: "D",
				},
			},
		}
		So(string(conf.CMD()), ShouldEqual, fmt.Sprintf(
			"select disk %d\r\nclean\r\nconvert gpt\r\ncreate partition primary\r\nformat quick fs=%s\r\nassign letter=%s\r\n",
			conf.Disk,
			conf.Partitions[0].FSType,
			conf.Partitions[0].Mountpoint,
		))

		conf = &DiskPartConfiguration{
			Disk: 1,
			GPT:  true,
			Partitions: []PartConfiguration{
				{
					Size:       "free",
					FSType:     "ntfs",
					Mountpoint: "C",
				},
			},
		}
		So(string(conf.CMD()), ShouldEqual, fmt.Sprintf(
			"select disk %d\r\nclean\r\nconvert gpt\r\ncreate partition primary\r\nformat quick fs=%s\r\nassign letter=%s\r\n",
			conf.Disk,
			conf.Partitions[0].FSType,
			conf.Partitions[0].Mountpoint,
		))
	})
}

func Test_listVolumes(t *testing.T) {
	Convey("返回计算机的磁盘卷列表（不包含已移除的卷和系统保留卷）", t, func() {
		var cmd *exec.Cmd
		monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Output", func(_ *exec.Cmd) ([]byte, error) {
			return ioutil.ReadFile("./testdata/diskpart_list_vol.txt")
		})
		defer monkey.UnpatchAll()

		vols, err := new(DiskPartConfiguration).listVolumes(log4test)
		So(err, ShouldBeNil)
		So(len(vols), ShouldEqual, 1)
		So(vols[0].ID, ShouldEqual, "1")
		So(vols[0].Letter, ShouldEqual, "C")
	})
}

func Test_rmVolumeLetter(t *testing.T) {
	Convey("移除指定卷的盘符", t, func() {
		var actualCmdArgs string
		var cmd *exec.Cmd
		monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Output", func(cmd *exec.Cmd) ([]byte, error) {
			b, _ := ioutil.ReadFile(cmd.Args[len(cmd.Args)-1])
			actualCmdArgs = string(b)
			return nil, nil
		})
		defer monkey.UnpatchAll()

		vol := &volume{
			ID:     "1",
			Letter: "D",
		}
		dc := new(DiskPartConfiguration)
		So(dc.rmVolumeLetter(log4test, vol), ShouldBeNil)
		So(actualCmdArgs, ShouldEqual, fmt.Sprintf("select vol %s\r\nremove letter=%s", vol.ID, vol.Letter))
	})
}

func Test_exec(t *testing.T) {
	Convey("执行diskpart相关命令", t, func() {
		cmdArgs := "hello world"
		timestamp := fmt.Sprintf("%d", time.Now().Unix())
		filename := fmt.Sprintf("diskpart_%s.txt", timestamp)

		monkey.Patch(utils.UUID, func() string {
			return timestamp
		})
		defer monkey.UnpatchAll()

		var cmd *exec.Cmd
		monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Output", func(_ *exec.Cmd) ([]byte, error) {
			b, err := ioutil.ReadFile(filename)
			return []byte(fmt.Sprintf("Your command is %q", b)), err
		})

		out, err := new(DiskPartConfiguration).exec(log4test, cmdArgs)
		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, fmt.Sprintf("Your command is %q", cmdArgs))
	})
}

func Test_diskSize(t *testing.T) {
	Convey("获取当前磁盘容量值", t, func() {
		Convey("预期当前磁盘(0)容量为25GB", func() {
			var cmd *exec.Cmd
			monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Output", func(_ *exec.Cmd) ([]byte, error) {
				return ioutil.ReadFile("./testdata/diskpart_list_disk_gb.txt")
			})
			defer monkey.UnpatchAll()

			size, err := new(DiskPartConfiguration).diskSize(log4test)
			So(err, ShouldBeNil)
			So(size, ShouldEqual, 25*mybytes.GB)
		})

		Convey("预期当前磁盘(1)容量为2TB", func() {
			var cmd *exec.Cmd
			monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Output", func(_ *exec.Cmd) ([]byte, error) {
				return ioutil.ReadFile("./testdata/diskpart_list_disk_tb.txt")
			})
			defer monkey.UnpatchAll()

			df := &DiskPartConfiguration{
				Disk: 1,
			}
			size, err := df.diskSize(log4test)
			So(err, ShouldBeNil)
			So(size, ShouldEqual, 2*mybytes.TB)
		})

		Convey("预期执行diskpart脚本发生错误", func() {
			var ErrExec = errors.New("exec error")
			var cmd *exec.Cmd
			monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Output", func(_ *exec.Cmd) ([]byte, error) {
				return nil, ErrExec
			})
			defer monkey.UnpatchAll()

			size, err := new(DiskPartConfiguration).diskSize(log4test)
			So(err, ShouldEqual, ErrExec)
			So(size, ShouldEqual, 0)
		})

		Convey("预期解析diskpart脚本输出发生错误", func() {
			var cmd *exec.Cmd
			monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Output", func(_ *exec.Cmd) ([]byte, error) {
				return ioutil.ReadFile("./testdata/diskpart_list_disk_not_selected.txt")
			})
			defer monkey.UnpatchAll()

			size, err := new(DiskPartConfiguration).diskSize(log4test)
			So(err, ShouldEqual, ErrDiskNotSelected)
			So(size, ShouldEqual, 0)
		})
	})
}

func Test_setup(t *testing.T) {
	Convey("分区结构体初始化", t, func() {
		Convey("预期当前磁盘(0)使用MBR分区", func() {
			var cmd *exec.Cmd
			monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Output", func(_ *exec.Cmd) ([]byte, error) {
				return ioutil.ReadFile("./testdata/diskpart_list_disk_gb.txt")
			})
			defer monkey.UnpatchAll()

			df := DiskPartConfiguration{
				Disk: 0,
			}
			So(df.setup(log4test), ShouldBeNil)
			So(df.GPT, ShouldBeFalse)
		})

		Convey("预期当前磁盘(1)使用GPT分区", func() {
			var cmd *exec.Cmd
			monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Output", func(_ *exec.Cmd) ([]byte, error) {
				return ioutil.ReadFile("./testdata/diskpart_list_disk_tb.txt")
			})
			defer monkey.UnpatchAll()

			df := DiskPartConfiguration{
				Disk: 1,
			}
			So(df.setup(log4test), ShouldBeNil)
			So(df.GPT, ShouldBeTrue)
		})
	})

}
