package app

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	opentracing "github.com/opentracing/opentracing-go"

	"github.com/alekns/tinyrstats/internal/config"
	"github.com/alekns/tinyrstats/internal/monitor"
)

type statsAppImpl struct {
	settings *config.Settings
	logger   log.Logger

	statsService monitor.StatsService
	events       monitor.Events
}

// QueryBy .
func (sa *statsAppImpl) QueryBy(ctx context.Context, query *monitor.QueryCallStatistic) (*monitor.QueryCallStatisticResult, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "QueryBy")
	defer span.Finish()

	logger := log.With(sa.logger, "method", "request statistic")

	level.Debug(logger).Log("msg", "QueryBy")

	minCount, maxCount := sa.statsService.GetMinMax()
	allHosts := sa.statsService.GetAllHosts()
	totalHosts := 0 // int(minCount) + int(maxCount)
	for _, value := range allHosts {
		totalHosts += value
	}

	span.SetTag("total_resources", totalHosts)

	return &monitor.QueryCallStatisticResult{
		TotalCount:       totalHosts,
		MinResponseCount: int(minCount),
		MaxResponseCount: int(maxCount),
		Resources:        allHosts,
	}, nil
}

func newStatsApp(settings *config.Settings,
	logger log.Logger,
	statsService monitor.StatsService,
	events monitor.Events) *statsAppImpl {

	hostHandler := func(args ...interface{}) {
		statsService.AddHost(args[0].(string), 1)
	}
	events.TaskQueriedByResource().On(&hostHandler)

	minHandler := func(args ...interface{}) {
		statsService.AddMinMax(false, 1)
	}
	events.TaskQueriedByMinResponse().On(&minHandler)

	maxHandler := func(args ...interface{}) {
		statsService.AddMinMax(true, 1)
	}
	events.TaskQueriedByMaxResponse().On(&maxHandler)

	return &statsAppImpl{
		logger:       log.With(logger, "service", "TaskApp"),
		settings:     settings,
		statsService: statsService,
		events:       events,
	}
}
