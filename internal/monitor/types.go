package worker

import (
	"net/http"
)

type (
	// TaskID .
	TaskID string

	Task struct {
		*HealthTask

		ID     TaskID            `json:"id"`
		Status *HealthTaskStatus `json:"status,omitempty"`
	}

	// HealthTask .
	HealthTask struct {
		Timeout int64       `json:"timeout"`

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
		Task     *HealthTask `json:"task,omitempty"`
	}
)
