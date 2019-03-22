package monitor

//go:generate mockgen -source=interfaces.go -package=monitor -destination=interfaces_mocks.go

import (
	"context"

	"github.com/alekns/tinyrstats/pkg/helpers/subscribs"
)

type (
	// TaskRepository .
	TaskRepository interface {
		GetByResponseTimeMinOrMax(context.Context, bool) (*Task, error)
		GetByID(context.Context, TaskID) (*Task, error)

		Save(context.Context, *Task) error

		Delete(context.Context, TaskID) error
		DeleteAll(context.Context)
	}

	// HealthService .
	HealthService interface {
		CheckStatus(context.Context, *HealthTask) (*HealthTaskStatus, error)
	}

	// ScheduleTaskService .
	ScheduleTaskService interface {
		Schedule(context.Context, TaskID, *ScheduleHealthTask) error
		Cancel(context.Context, TaskID) error
		CancelAll(context.Context) error
	}

	// Events all processable events.
	Events interface {
		TaskQueriedByURL() subscribs.EventHandler
		TaskQueriedByMinResponse() subscribs.EventHandler
		TaskQueriedByMaxResponse() subscribs.EventHandler
	}

	// TaskApp .
	TaskApp interface {
		Create(context.Context, *CreateTaskCommand) (*CreateTaskResult, error)
		QueryBy(context.Context, *QueryTask) (*QueryTaskResult, error)
	}

	// StatsApp .
	StatsApp interface {
		QueryBy(context.Context, *QueryStatistic) (*QueryStatisticResult, error)
	}
)
