package services

import (
	"context"
	"testing"

	"github.com/alekns/tinyrstats/internal/monitor"
	"github.com/golang/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"
)

func TestResultsTaskRepositoryConsumerSpec(t *testing.T) {
	Convey("Given created results task repository consumer", t, func(c C) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockTaskRep := monitor.NewMockTaskRepository(mockCtrl)
		mockTaskRep.EXPECT().
			Save(gomock.Any(), gomock.Eq(&monitor.Task{
				HealthTask: monitor.HealthTask{
					ID: "1",
				},
				Status: &monitor.HealthTaskStatus{
					StatusCode: 200,
				},
			})).
			Return(nil).
			Times(1)

		consumer := NewResultsTaskRepositoryConsumer(mockTaskRep)

		c.Convey("Consumer should call task repository save", func(c C) {
			So(consumer.Accept(ctx,
				&monitor.HealthTask{ID: "1"}, &monitor.HealthTaskStatus{StatusCode: 200}), ShouldBeNil)
		})
	})
}
