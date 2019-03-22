package services

import (
	"context"
	"testing"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/alekns/tinyrstats/pkg/helpers/runner"

	"github.com/alekns/tinyrstats/internal/config"
	"github.com/alekns/tinyrstats/internal/monitor"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestScheduleTaskSpec(t *testing.T) {
	Convey("Given created interval scheduler", t, func(c C) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		task := &monitor.HealthTask{
			ID: "1",
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockConsumer := runner.NewMockConsumer(mockCtrl)
		mockConsumer.EXPECT().
			Accept(gomock.Any(), gomock.Eq(task)).
			Return(nil).
			Times(1)

		scheduler := NewIntervalScheduler(log.NewNopLogger(),
			&config.SchedulerSettings{
				MaxConcurrency:  4,
				DefaultInterval: -10, // for immediately
			},
			mockConsumer)
		defer scheduler.Stop(ctx)

		So(scheduler.Schedule(ctx, "1", &monitor.ScheduleHealthTask{
			Interval: 0, // immediately
			Task:     task,
		}), ShouldBeNil)

		c.So(scheduler.Start(ctx), ShouldBeNil)

		c.Convey("Then after half second consumer should receive valid values", func(c C) {
			time.Sleep(250 * time.Millisecond)
		})
	})
}
