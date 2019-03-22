package webhttp

import (
	"context"
	"errors"
	"net/http"

	"github.com/alekns/tinyrstats/internal/monitor"

	"github.com/alekns/tinyrstats/internal/presentation/external/endpoints"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

type httpServerTransport struct {
	checkStatus http.Handler
}

// SetupHTTPServerHandlers .
func SetupHTTPServerHandlers(logger log.Logger,
	route *mux.Route,
	endpointsSet *endpoints.Set) {

	router := route.Subrouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
	}

	router.Methods("GET").Path("/tasks/actions/query").Handler(httptransport.NewServer(
		endpointsSet.TaskQuery,
		decodeQueryTaskRequest,
		httptransport.EncodeJSONResponse,
		options...,
	))

	router.Methods("GET").Path("/statistics/tasks/queries").Handler(httptransport.NewServer(
		endpointsSet.StatisticsQuery,
		decodeStatisticsQueryRequest,
		httptransport.EncodeJSONResponse,
		options...,
	))
}

func decodeQueryTaskRequest(_ context.Context, req *http.Request) (request interface{}, err error) {
	values := req.URL.Query()
	host := values.Get("resource")
	responseTime := values.Get("responsetime")

	if len(host) == 0 && len(responseTime) == 0 {
		return nil, errors.New("no params")
	}

	if len(host) != 0 && len(responseTime) != 0 {
		return nil, errors.New("only one request supported")
	}

	responseTimeType := monitor.QueryResponseMinTime
	if len(responseTime) != 0 {
		switch values.Get("responsetime") {
		case "max":
			responseTimeType = monitor.QueryResponseMaxTime
		case "min":
		default:
			return nil, errors.New("only min max supported")
		}
	}

	return &monitor.QueryTask{
		ByHost:         host,
		ByResponseTime: responseTimeType,
	}, nil
}

func decodeStatisticsQueryRequest(_ context.Context, req *http.Request) (request interface{}, err error) {
	return &monitor.QueryCallStatistic{}, nil
}
