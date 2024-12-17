package win

import (
	"errors"
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
func TestPhysicalMachineSN(t *testing.T) {
	Convey("获取当前设备SN", t, func() {
		Convey("命令执行发生错误", func() {
			var ErrExec = errors.New("exec error")
			monkey.Patch(ExecOutputWithLog, func(_ logger.Logger, _ string) ([]byte, error) {
				return nil, ErrExec
			})
			defer monkey.UnpatchAll()

			sn, err := PhysicalMachineSN(log4test)
			So(err, ShouldEqual, ErrExec)
			So(sn, ShouldBeEmpty)
		})

		Convey("命令执行输出的格式与预期一致", func() {
			monkey.Patch(ExecOutputWithLog, func(_ logger.Logger, _ string) ([]byte, error) {
				return ioutil.ReadFile("./testdata/sn.txt")
			})
			defer monkey.UnpatchAll()

			sn, err := PhysicalMachineSN(log4test)
			So(err, ShouldBeNil)
			So(sn, ShouldEqual, "0")
		})
	})
}
