package wmic

import (
	"errors"
	"io/ioutil"
	"testing"

	"bou.ke/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/utils"
	winutil "idcos.io/cloudboot/utils/win"
)

func TestPartitionSizeByIndex(t *testing.T) {
	Convey("根据索引查询分区的大小（单位字节）", t, func() {
		Convey("命令行查询获得正整数值", func() {
			monkey.Patch(winutil.ExecOutput, func(log logger.Logger, cmdAndArgs string) (output []byte, err error) {
				outputUTF8, err := ioutil.ReadFile("./testdata/partition.output")
				if err != nil {
					return nil, err
				}
				return []byte(utils.UTF82GBK(string(outputUTF8))), nil
			})
			defer monkey.Unpatch(winutil.ExecOutput)

			size, err := PartitionSizeByIndex(log4test, 0)
			So(err, ShouldBeNil)
			So(size, ShouldEqual, 367001600)
		})

		Convey("命令行查询获得非整数值", func() {
			Convey("命令执行的标准输出中不包含解析关键字", func() {
				monkey.Patch(winutil.ExecOutput, func(log logger.Logger, cmdAndArgs string) (output []byte, err error) {
					return []byte(utils.UTF82GBK(string(`helloworld`))), nil
				})
				defer monkey.Unpatch(winutil.ExecOutput)

				size, err := PartitionSizeByIndex(log4test, 0)
				So(err, ShouldEqual, ErrKeywordsNotFound)
				So(size, ShouldEqual, 0)
			})

			Convey("命令执行的标准输出中包含解析关键字，但对应的值不是整数", func() {
				monkey.Patch(winutil.ExecOutput, func(log logger.Logger, cmdAndArgs string) (output []byte, err error) {
					return []byte(utils.UTF82GBK(string("C:\firstboot>wmic partition where Index=0 get Size /value\r\nSize=xxxx"))), nil
				})
				defer monkey.Unpatch(winutil.ExecOutput)

				size, err := PartitionSizeByIndex(log4test, 0)
				So(err, ShouldNotBeNil)
				So(size, ShouldEqual, 0)
			})
		})

		Convey("命令行查询进程意外退出", func() {
			var ErrExec = errors.New("exec error")
			monkey.Patch(winutil.ExecOutput, func(log logger.Logger, cmdAndArgs string) (output []byte, err error) {
				return nil, ErrExec
			})
			defer monkey.Unpatch(winutil.ExecOutput)

			size, err := PartitionSizeByIndex(log4test, 0)
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, ErrExec)
			So(size, ShouldEqual, 0)
		})
	})
}
