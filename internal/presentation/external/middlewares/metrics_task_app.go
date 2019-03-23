package middlewares

import (
	"context"
	"fmt"
	"time"

	"github.com/alekns/tinyrstats/internal"
	"github.com/alekns/tinyrstats/internal/monitor"
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

type metricsTaskAppDecorator struct {
	queryByCount      metrics.Counter
	createAndRunCount metrics.Counter

	origInstance monitor.TaskApp
}

func (mw *metricsTaskAppDecorator) CreateAndRun(ctx context.Context, cmd *monitor.CreateTaskCommand) (result *monitor.CreateTaskResult, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "CreateAndRun", "error", fmt.Sprint(err != nil)}
		mw.createAndRunCount.With(lvs...).Add(1)
	}(time.Now())

	return mw.origInstance.CreateAndRun(ctx, cmd)
}

func (mw *metricsTaskAppDecorator) QueryBy(ctx context.Context, query *monitor.QueryTask) (result *monitor.QueryTaskResult, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "QueryBy", "error", fmt.Sprint(err != nil)}
		mw.queryByCount.With(lvs...).Add(1)
	}(time.Now())

	return mw.origInstance.QueryBy(ctx, query)
}

// WrapMetricsTaskApp is wrapper for TaskApp metrics.
func WrapMetricsTaskApp(instance monitor.TaskApp) monitor.TaskApp {
	return &metricsTaskAppDecorator{
		createAndRunCount: kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: internal.ServiceName,
			Subsystem: "task_app",
			Name:      "createandrun_count",
			Help:      "Number of TaskApp.CreateAndRun requests.",
		}, []string{"method", "error"}),
		queryByCount: kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: internal.ServiceName,
			Subsystem: "task_app",
			Name:      "queryby_count",
			Help:      "Number of TaskApp.QueryBy requests.",
		}, []string{"method", "error"}),
		origInstance: instance,
	}
}
