package times

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

const layout = "2006-01-02"

func TestLatestDays(t *testing.T) {
	Convey("返回最近N天的日期列表", t, func() {
		Convey("日期递增", func() {
			var items []time.Time
			bt, _ := time.Parse(layout, "2018-10-17")
			items = LatestDays(bt, 2, true)
			So(len(items), ShouldEqual, 2)
			So(items[0].Format(layout), ShouldEqual, "2018-10-16")
			So(items[1].Format(layout), ShouldEqual, "2018-10-17")

			items = LatestDays(bt, 1, true)
			So(len(items), ShouldEqual, 1)
			So(items[0].Format(layout), ShouldEqual, "2018-10-17")
		})

		Convey("日期递减", func() {
			var items []time.Time
			bt, _ := time.Parse(layout, "2018-10-17")
			items = LatestDays(bt, 2, false)
			So(len(items), ShouldEqual, 2)
			So(items[0].Format(layout), ShouldEqual, "2018-10-17")
			So(items[1].Format(layout), ShouldEqual, "2018-10-16")

			items = LatestDays(bt, 1, false)
			So(len(items), ShouldEqual, 1)
			So(items[0].Format(layout), ShouldEqual, "2018-10-17")
		})
	})
}
