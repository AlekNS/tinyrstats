package services

import (
	"testing"

	"github.com/alekns/tinyrstats/internal/config"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStatsServiceSpec(t *testing.T) {
	Convey("Given created stats service", t, func(c C) {

		svc := NewStatsServiceImpl(&config.StatsSettings{
			BucketsCount: 8,
		})

		c.Convey("When add min", func(c C) {

			svc.AddMinMax(false, 1)

			c.Convey("Then getted values should be valid", func(c C) {
				min, max := svc.GetMinMax()
				So(min, ShouldEqual, 1)
				So(max, ShouldEqual, 0)
			})
		})

		c.Convey("When add 2 hosts", func(c C) {

			svc.AddHost("host1", 1)
			svc.AddHost("host1", 1)
			svc.AddHost("host2", 1)

			c.Convey("Then getAll should return two hosts", func(c C) {
				items := svc.GetAllHosts()

				So(len(items), ShouldEqual, 2)

				So(items, ShouldContainKey, "host1")
				So(items, ShouldContainKey, "host2")

				So(items["host1"], ShouldEqual, 2)
				So(items["host2"], ShouldEqual, 1)
			})
		})
	})
}
