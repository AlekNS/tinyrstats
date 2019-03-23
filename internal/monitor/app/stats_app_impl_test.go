package app

import (
	"context"
	"testing"

	"github.com/alekns/tinyrstats/internal/monitor"
	"github.com/go-kit/kit/log"

	"github.com/alekns/tinyrstats/internal/config"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStatsAppImplSpec(t *testing.T) {
	Convey("Given created stats application", t, func(c C) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockStatsSvc := monitor.NewMockStatsService(mockCtrl)
		mockStatsSvc.EXPECT().
			GetAllHosts().
			Return(monitor.StatsHostsInfo{"host1": 10}).
			Times(1)
		mockStatsSvc.EXPECT().
			GetMinMax().
			Return(int32(3), int32(5)).
			Times(1)

		statsApp := newStatsApp(&config.Settings{},
			log.NewNopLogger(), mockStatsSvc, NewSyncEventsImpl())

		c.Convey("When query was called success", func(c C) {
			result, err := statsApp.QueryBy(context.TODO(), &monitor.QueryCallStatistic{})

			c.Convey("Then no errors was happend and valid results are received", func(c C) {
				So(err, ShouldBeNil)
				So(result, ShouldNotBeNil)

				So(result.TotalCount, ShouldEqual, 10)
				So(len(result.Resources), ShouldEqual, 1)
				So(result.MinResponseCount, ShouldEqual, int32(3))
				So(result.MaxResponseCount, ShouldEqual, int32(5))
			})
		})
	})
}
