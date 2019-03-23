package monitor

import (
	"net/http"
)

type (
	// TaskID .
	TaskID string

	// Task .
	Task struct {
		HealthTask

		Status *HealthTaskStatus `json:"status,omitempty"`
	}

	// HealthTask .
	HealthTask struct {
		ID TaskID `json:"id"`

		Timeout int64 `json:"timeout"`

		// HTTP part
		// @TODO: use sub-structure for possible customization.

		URL     string      `json:"url"`
		Method  string      `json:"method"`
		Body    string      `json:"body,omitempty"`
		Headers http.Header `json:"headers,omitempty"`
	}

	// HealthTaskStatus .
	HealthTaskStatus struct {
		Error *HealthTaskError `json:"error,omitempty"`

		LastTime     int64 `json:"lastTime"`
		ResponseTime int64 `json:"responseTime,omitempty"`

		// HTTP part
		// @TODO: use sub-structure for possible customization.
		StatusCode int         `json:"statusCode,omitempty"`
		Headers    http.Header `json:"responseHeaders,omitempty"`
	}

	// HealthTaskError .
	HealthTaskError struct {
		Text string `json:"text"`

		IsTimeout  bool `json:"isTimeout,omitempty"`
		IsDNSError bool `json:"isDnsError,omitempty"`
	}

	// ScheduleHealthTask .
	ScheduleHealthTask struct {
		Interval int         `json:"interval,omitempty"`
		Task     *HealthTask `json:"task"`
	}

	// Applications DTO

	// CreateTaskCommand .
	CreateTaskCommand struct {
		ScheduleHealthTask
	}

	// CreateTaskResult .
	CreateTaskResult struct {
		ID string `json:"id"`
	}

	// QueryResponseTimeType .
	QueryResponseTimeType int

	// QueryTask .
	QueryTask struct {
		ByHost string `json:"host"`

		ByResponseTime QueryResponseTimeType `json:"responseTime"`
	}

	// QueryTaskResult .
	QueryTaskResult struct {
		Task
	}

	// QueryCallStatistic is dummy, need for extension in future.
	QueryCallStatistic struct{}

	// QueryCallStatisticResult .
	QueryCallStatisticResult struct {
		TotalCount       int            `json:"totalCount"`
		Resources        StatsHostsInfo `json:"resourcesCounts"`
		MinResponseCount int            `json:"minResponseCount"`
		MaxResponseCount int            `json:"maxResponseCount"`
	}
)
