package services

import (
	"context"
	"testing"
	"time"

	"github.com/alekns/tinyrstats/internal/config"
	"github.com/alekns/tinyrstats/internal/monitor"
	"github.com/alekns/tinyrstats/pkg/helpers/runner"
	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTaskExecutorSpec(t *testing.T) {
	Convey("Given created task executor", t, func(c C) {
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

		taskExec := NewConcurretTaskExecutor(&config.Settings{
			Tasks: &config.TasksSettings{
				MaxConcurrency: 2,
				MaxPending:     4,
				TaskQueueSize:  2,
			},
		}, log.NewNopLogger(), mockConsumer)
		taskExec.Start(ctx)
		defer taskExec.Stop()

		c.Convey("When task is accepted", func(c C) {

			err := taskExec.Accept(ctx, task)
			So(err, ShouldBeNil)

			c.Convey("Then receiver consumer should receive same task after half second", func(c C) {
				time.Sleep(time.Millisecond * 250)
			})
		})
	})
}
