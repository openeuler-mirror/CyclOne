package wmic

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"bou.ke/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"idcos.io/cloudboot/config"
	"idcos.io/cloudboot/logger"
)

var log4test logger.Logger

func init() {
	log4test = logger.NewBeeLogger(&config.Logger{
		Level:   "debug",
		LogFile: filepath.Join(os.TempDir(), "gotest.log"),
	})
}

func TestGetNIC(t *testing.T) {
	Convey("执行命令行命令查询本地网卡信息", t, func() {
		Convey("本地有多块网卡", func() {
			monkey.Patch(nicGet, func(log logger.Logger, where string, properties ...string) (output []byte, err error) {
				outputUTF8, err := ioutil.ReadFile("./testdata/nic_get_more.output")
				if err != nil {
					return nil, err
				}
				return []byte(string(outputUTF8)), nil
			})
			defer monkey.Unpatch(nicGet)

			nics, err := GetNIC(log4test)
			So(err, ShouldBeNil)
			So(len(nics), ShouldEqual, 2)
			So(nics[0].MacAddr, ShouldEqual, "08:00:27:FB:FC:47")
			So(nics[0].NetConnID, ShouldEqual, "以太网")
			So(nics[1].MacAddr, ShouldEqual, "08:00:28:FB:FC:48")
			So(nics[1].NetConnID, ShouldEqual, "test")
		})

		Convey("本地仅有一块网卡", func() {
			monkey.Patch(nicGet, func(log logger.Logger, where string, properties ...string) (output []byte, err error) {
				outputUTF8, err := ioutil.ReadFile("./testdata/nic_get_one.output")
				if err != nil {
					return nil, err
				}
				return []byte(string(outputUTF8)), nil
			})
			defer monkey.Unpatch(nicGet)

			nics, err := GetNIC(log4test)
			So(err, ShouldBeNil)
			So(len(nics), ShouldEqual, 1)
			So(nics[0].MacAddr, ShouldEqual, "08:00:27:FB:FC:47")
			So(nics[0].NetConnID, ShouldEqual, "以太网")
		})
	})
}
