package monitor

//go:generate mockgen -source=interfaces.go -package=monitor -destination=interfaces_mocks.go

import (
	"context"

	"github.com/alekns/tinyrstats/pkg/helpers/subscribs"
)

type (
	// TaskRepository .
	TaskRepository interface {
		// GetByResponseTimeMinOrMax .
		GetByResponseTimeMinOrMax(context.Context, bool) (*Task, error)
		// GetByID .
		GetByID(context.Context, TaskID) (*Task, error)

		// Save .
		Save(context.Context, *Task) error

		// Delete .
		Delete(context.Context, TaskID) error
		// DeleteAll .
		DeleteAll(context.Context)
	}

	// HealthService .
	HealthService interface {
		CheckStatus(context.Context, *HealthTask) (*HealthTaskStatus, error)
	}

	// ScheduleTaskService .
	ScheduleTaskService interface {
		// Schedule .
		Schedule(context.Context, TaskID, *ScheduleHealthTask) error
		// Cancel .
		Cancel(context.Context, TaskID) error
		// CancelAll .
		CancelAll(context.Context) error
	}

	// Events all processable events.
	Events interface {
		// TaskQueriedByURL .
		TaskQueriedByURL() subscribs.EventHandler
		// TaskQueriedByMinResponse .
		TaskQueriedByMinResponse() subscribs.EventHandler
		// TaskQueriedByMaxResponse .
		TaskQueriedByMaxResponse() subscribs.EventHandler
	}

	// StatsHostsInfo .
	StatsHostsInfo = map[string]int

	// StatsService gathering statistics of the requests
	StatsService interface {
		GetAllHosts() StatsHostsInfo
		GetMinMax() (int32, int32)

		AddHost(string, int)
		AddMinMax(bool, int32)

		DeleteHost(string)
	}

	// TaskApp .
	TaskApp interface {
		// CreateAndRun .
		CreateAndRun(context.Context, *CreateTaskCommand) (*CreateTaskResult, error)
		// QueryBy .
		QueryBy(context.Context, *QueryTask) (*QueryTaskResult, error)
	}

	// StatsApp .
	StatsApp interface {
		// QueryBy .
		QueryBy(context.Context, *QueryCallStatistic) (*QueryCallStatisticResult, error)
	}
)
