package endpoints

import (
	"github.com/alekns/tinyrstats/internal/config"
	"github.com/alekns/tinyrstats/internal/monitor/app"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

// Set .
type Set struct {
	TaskQuery       endpoint.Endpoint
	StatisticsQuery endpoint.Endpoint
}

// NewSet .
func NewSet(settings *config.EndpointsSettings,
	logger log.Logger,
	registry app.Registry) *Set {

	// healthCheckerSvc = WrapLoggingHealthCheckerService(healthCheckerSvc, logger)
	// healthCheckerSvc = WrapMetricsHealthCheckerService(healthCheckerSvc)

	// healthCheckStatusEndpoint := makeResourceCheckStatusEndpoint(settings,
	// 	tracer,
	// 	nc,
	// 	healthCheckerSvc)
	// healthCheckStatusEndpoint = WrapTracerToRequestForHealthCheckerService()(healthCheckStatusEndpoint)
	// healthCheckStatusEndpoint = opentracing.TraceServer(tracer, "RequestHealthCheck")(healthCheckStatusEndpoint)

	taskQueryEndpoint := makeTaskQueryEndpoint(registry)
	statisticsQueryEndpoint := makeStatisticsQueryEndpoint(registry)

	return &Set{
		TaskQuery:       taskQueryEndpoint,
		StatisticsQuery: statisticsQueryEndpoint,
	}
}
