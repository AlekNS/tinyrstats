package endpoints

import (
	"context"

	"github.com/alekns/tinyrstats/internal/monitor"
	"github.com/alekns/tinyrstats/internal/monitor/app"
	"github.com/go-kit/kit/endpoint"
)

func makeStatisticsQueryEndpoint(registry app.Registry) endpoint.Endpoint {
	return func(ctx context.Context, reqRaw interface{}) (interface{}, error) {
		req := reqRaw.(*monitor.QueryCallStatistic)
		statsApp := registry.StatsApp()
		return wrapResponseToData(statsApp.QueryBy(ctx, req))
	}
}
