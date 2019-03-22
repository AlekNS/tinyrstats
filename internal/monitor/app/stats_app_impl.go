package app

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/alekns/tinyrstats/internal/config"
	"github.com/alekns/tinyrstats/internal/monitor"
)

type statsAppImpl struct {
	settings *config.Settings
	logger   log.Logger

	events monitor.Events
}

// Create .
func (ta *statsAppImpl) QueryBy(ctx context.Context, query *monitor.QueryStatistic) (*monitor.QueryStatisticResult, error) {
	logger := log.With(ta.logger, "method", "QueryBy")

	level.Debug(logger).Log("msg", "QueryBy")

	return &monitor.QueryStatisticResult{}, nil
}

func newStatsApp(settings *config.Settings,
	logger log.Logger,
	events monitor.Events) *statsAppImpl {

	return &statsAppImpl{
		logger:   log.With(logger, "service", "TaskApp"),
		settings: settings,
		events:   events,
	}
}
