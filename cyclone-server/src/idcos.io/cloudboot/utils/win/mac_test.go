package win

import (
	"errors"
	"io/ioutil"
	"testing"

	"idcos.io/cloudboot/logger"

	"bou.ke/monkey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMacAddress(t *testing.T) {
	Convey("获取当前设备mac地址", t, func() {
		Convey("命令执行发生错误", func() {
			var ErrExec = errors.New("exec error")
			monkey.Patch(ExecOutputWithLog, func(_ logger.Logger, _ string) ([]byte, error) {
				return nil, ErrExec
			})
			defer monkey.UnpatchAll()

			sn, err := MacAddress(log4test)
			So(err, ShouldEqual, ErrExec)
			So(sn, ShouldBeEmpty)
		})

		Convey("命令执行输出的格式与预期一致", func() {
			monkey.Patch(ExecOutputWithLog, func(_ logger.Logger, _ string) ([]byte, error) {
				return ioutil.ReadFile("./testdata/mac_addr.txt")
			})
			defer monkey.UnpatchAll()

			sn, err := MacAddress(log4test)
			So(err, ShouldBeNil)
			So(sn, ShouldEqual, "08:00:27:FB:FC:47")
		})
	})
}
