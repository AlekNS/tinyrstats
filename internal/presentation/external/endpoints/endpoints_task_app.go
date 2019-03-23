package endpoints

import (
	"context"

	"github.com/alekns/tinyrstats/internal/presentation/external/middlewares"

	"github.com/alekns/tinyrstats/internal/monitor"

	"github.com/alekns/tinyrstats/internal/monitor/app"
	"github.com/go-kit/kit/endpoint"
)

func makeTaskQueryEndpoint(registry app.Registry) endpoint.Endpoint {
	taskApp := middlewares.WrapMetricsTaskApp(registry.TaskApp())

	return func(ctx context.Context, reqRaw interface{}) (interface{}, error) {
		req := reqRaw.(*monitor.QueryTask)

		return wrapResponseToData(taskApp.QueryBy(ctx, req))
	}
}
