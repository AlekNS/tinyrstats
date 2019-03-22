package services

import (
	"bytes"
	"context"
	"strings"
	"time"

	"github.com/alekns/tinyrstats/internal/config"
	"github.com/alekns/tinyrstats/internal/monitor"
	"github.com/alekns/tinyrstats/pkg/helpers/network"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

const contentLengthReaderBufferSize = 1 << 16

// HTTPHealthService .
type HTTPHealthService struct {
	settings *config.TasksSettings
	logger   log.Logger
}

func getTimeWithMilliseconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// CheckStatus .
func (rh *HTTPHealthService) CheckStatus(ctx context.Context, request *monitor.HealthTask) (*monitor.HealthTaskStatus, error) {
	logger := log.With(rh.logger,
		"method", "CheckStatus",
		"taskId", request.ID,
		"taskMethod", request.Method,
		"taskUrl", request.URL)

	if !(strings.HasPrefix(request.URL, "http:") || strings.HasPrefix(request.URL, "https:")) {
		level.Error(logger).Log("err", monitor.ErrResourceRequestNotHTTP.Error())
		return nil, monitor.ErrResourceRequestNotHTTP
	}

	status := &monitor.HealthTaskStatus{
		LastTime: getTimeWithMilliseconds(),
	}

	level.Debug(logger).Log("msg", "request HTTP resource")

	timeout := request.Timeout
	if timeout < 1 {
		timeout = int64(rh.settings.DefaultTimeout)
	}

	// @TODO: Use DNS resolve name first method
	response, err := network.HTTPRequestAndGetResponse(ctx,
		time.Duration(timeout)*time.Millisecond,
		request.Method,
		request.URL,
		bytes.NewReader([]byte(request.Body)),
		request.Headers)

	status.ResponseTime = getTimeWithMilliseconds() - status.LastTime

	level.Debug(logger).Log("msg", "complete requesting of HTTP resource status")

	if response != nil {
		status.StatusCode = response.StatusCode
		status.Headers = response.Header

		response.Body.Close()
	}

	if err != nil {
		status.Error = &monitor.HealthTaskError{
			Text: err.Error(),
		}
		if strings.Contains(err.Error(), "context deadline exceeded") {
			status.Error.IsTimeout = true
			status.Error.Text = "timeout"
		} else if strings.Contains(err.Error(), ": dial tcp: lookup ") {
			status.Error.IsDNSError = true
			status.Error.Text = "dns error"
		} else if err != network.ErrHTTPClientError && err != network.ErrHTTPServerError {
			return status, err
		}
	}

	return status, nil
}

// NewHTTPHealthService .
func NewHTTPHealthService(settings *config.TasksSettings, logger log.Logger) *HTTPHealthService {
	return &HTTPHealthService{
		settings: settings,
		logger:   log.With(logger, "service", "HTTPHealthService"),
	}
}
