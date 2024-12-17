package bytes

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestByte2GBRounding(t *testing.T) {
	Convey("将字节转化为整数吉字节", t, func() {
		So(Byte2GBRounding(GB), ShouldEqual, 1)
		So(Byte2GBRounding(7*GB), ShouldEqual, 7)
		So(Byte2GBRounding(1024*MB), ShouldEqual, 1)
		So(Byte2GBRounding(1024*1024*KB), ShouldEqual, 1)
		So(Byte2GBRounding(1024*1024*1024*B), ShouldEqual, 1)
		So(Byte2GBRounding(1023*MB), ShouldEqual, 0)
		So(Byte2GBRounding(1025*MB), ShouldEqual, 1)
		So(Byte2GBRounding(2047*MB), ShouldEqual, 1)
		So(Byte2GBRounding(2048*MB), ShouldEqual, 2)
	})
}

func TestByte2GB(t *testing.T) {
	Convey("将字节转化为吉字节", t, func() {
		So(Byte2GB(GB), ShouldEqual, 1)
		So(Byte2GB(7*GB), ShouldEqual, 7)
		So(Byte2GB(1024*MB), ShouldEqual, 1)
		So(Byte2GB(1024*1024*KB), ShouldEqual, 1)
		So(Byte2GB(1024*1024*1024*B), ShouldEqual, 1)

		So(Byte2GB(1023*MB), ShouldEqual, float64(1023)/float64(1024))
		So(Byte2GB(1025*MB), ShouldEqual, float64(1025)/float64(1024))
		So(Byte2GB(2047*MB), ShouldEqual, float64(2047)/float64(1024))
		So(Byte2GB(2048*MB), ShouldEqual, float64(2048)/float64(1024))
	})
}

func TestParse2Byte(t *testing.T) {
	Convey("容量值解析", t, func() {
		So(TB, ShouldEqual, 1099511627776)
		So(GB, ShouldEqual, 1073741824)
		So(MB, ShouldEqual, 1048576)
		So(KB, ShouldEqual, 1024)
		So(B, ShouldEqual, 1)

		var size Byte
		var err error

		size, err = Parse2Byte("7.5", "MB")
		So(err, ShouldEqual, ErrMalformedSizeStringValue)

		size, err = Parse2Byte("7", "AB")
		So(err, ShouldEqual, ErrMalformedUnitStringValue)

		size, err = Parse2Byte("7", "B")
		So(err, ShouldBeNil)
		So(size, ShouldEqual, Byte(7)*B)

		size, err = Parse2Byte("7", "kb")
		So(err, ShouldBeNil)
		So(size, ShouldEqual, Byte(7)*KB)

		size, err = Parse2Byte("7", "mb")
		So(err, ShouldBeNil)
		So(size, ShouldEqual, Byte(7)*MB)

		size, err = Parse2Byte("7", "GB")
		So(err, ShouldBeNil)
		So(size, ShouldEqual, Byte(7)*GB)

		size, err = Parse2Byte("7", "TB")
		So(err, ShouldBeNil)
		So(size, ShouldEqual, Byte(7)*TB)
	})
}
