package webhttp

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/alekns/tinyrstats/internal/monitor"
	"github.com/pkg/errors"

	"github.com/alekns/tinyrstats/internal/presentation/external/endpoints"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// SetupHTTPServerHandlers pre-fill router by endpoints.
func SetupHTTPServerHandlers(logger log.Logger,
	route *mux.Route,
	endpointsSet *endpoints.Set) {

	router := route.Subrouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeErrorResponse),
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

func encodeErrorResponse(ctx context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("abnormal error content, encodeErrorResponse with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	statusCode := http.StatusBadRequest
	switch err {
	case monitor.ErrTaskNotFound:
		statusCode = http.StatusNotFound
	default:
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"statusCode": statusCode,
		"error":      err.Error(),
	})
}

func decodeQueryTaskRequest(ctx context.Context, req *http.Request) (request interface{}, err error) {
	values := req.URL.Query()

	const kresponsetime = "responsetime"

	// Validate parameters
	host := values.Get("resource")
	responseTime := values.Get(kresponsetime)

	if len(host) == 0 && len(responseTime) == 0 {
		return nil, ErrNoParameters
	}

	if len(host) != 0 && len(responseTime) != 0 {
		return nil, ErrTooManyParameters
	}

	responseTimeType := monitor.QueryResponseMinTime
	if len(responseTime) != 0 {
		switch responseTime {
		case "max":
			responseTimeType = monitor.QueryResponseMaxTime
		case "min":
		default:
			return nil, errors.Wrap(ErrInvalidParameterValue, kresponsetime+" support only min or max")
		}
	}

	return &monitor.QueryTask{
		ByHost:         host,
		ByResponseTime: responseTimeType,
	}, nil
}

func decodeStatisticsQueryRequest(ctx context.Context, req *http.Request) (request interface{}, err error) {
	return &monitor.QueryCallStatistic{}, nil
}
