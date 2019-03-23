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

func TestTaskAppImplSpec(t *testing.T) {
	Convey("Given created task application", t, func(c C) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockScheduler := monitor.NewMockScheduleTaskService(mockCtrl)
		mockTaskRep := monitor.NewMockTaskRepository(mockCtrl)

		resultTask := &monitor.Task{}
		resultTask.ID = monitor.TaskID("host1")

		mockTaskRep.EXPECT().
			GetByID(gomock.Any(), resultTask.ID).
			Return(resultTask, nil).
			Times(1)

		taskApp := newTaskApp(&config.Settings{},
			log.NewNopLogger(), NewSyncEventsImpl(), mockScheduler, mockTaskRep)

		c.Convey("When query by resource was called success", func(c C) {
			result, err := taskApp.QueryBy(context.TODO(), &monitor.QueryTask{
				ByHost: "host1",
			})

			c.Convey("Then no errors was happend and valid results are received", func(c C) {
				So(err, ShouldBeNil)
				So(result, ShouldNotBeNil)
				So(result.ID, ShouldEqual, resultTask.ID)
			})
		})
	})
}
