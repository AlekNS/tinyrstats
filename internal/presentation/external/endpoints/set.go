package endpoints

import (
	"github.com/alekns/tinyrstats/internal/config"
	"github.com/alekns/tinyrstats/internal/monitor/app"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kittracing "github.com/go-kit/kit/tracing/opentracing"
	opentracing "github.com/opentracing/opentracing-go"
)

// Set .
type Set struct {
	TaskQuery endpoint.Endpoint

	StatisticsQuery endpoint.Endpoint
}

// NewSet .
func NewSet(settings *config.EndpointsSettings,
	logger log.Logger,
	registry app.Registry) *Set {

	tracer := opentracing.GlobalTracer()

	taskQueryEndpoint := kittracing.TraceServer(tracer, "TaskApp")(
		makeTaskQueryEndpoint(registry))
	statisticsQueryEndpoint := kittracing.TraceServer(tracer, "StatsApp")(
		makeStatisticsQueryEndpoint(registry))

	return &Set{
		TaskQuery:       taskQueryEndpoint,
		StatisticsQuery: statisticsQueryEndpoint,
	}
}
