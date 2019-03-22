package services

import (
	"context"
	"testing"

	"github.com/alekns/tinyrstats/internal/monitor"
	"github.com/alekns/tinyrstats/pkg/helpers/runner"
	"github.com/golang/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHealthServiceConsumerSpec(t *testing.T) {
	Convey("Given created health service consumer", t, func(c C) {
		mockCtrl := gomock.NewController(t)

		task := &monitor.HealthTask{
			Method: "GET",
		}
		result := &monitor.HealthTaskStatus{}
		mockHealthSvc := monitor.NewMockHealthService(mockCtrl)
		mockHealthSvc.EXPECT().
			CheckStatus(gomock.Any(), gomock.Eq(task)).
			Return(result, nil).
			Times(1)
		mockConsumer := runner.NewMockConsumer(mockCtrl)
		mockConsumer.EXPECT().
			Accept(gomock.Any(), gomock.Eq(task), gomock.Eq(result)).
			Return(nil).
			Times(1)

		defer mockCtrl.Finish()

		healthConsumer := NewHealthServiceConsumer(mockHealthSvc, mockConsumer)

		c.Convey("Consumer should receive success and valid values", func(c C) {
			err := healthConsumer.Accept(context.Background(), task)

			c.So(err, ShouldBeNil)
		})
	})
}
