package worker

//go:generate mockgen -source=interfaces.go -package=worker -destination=interfaces_mocks.go

import "context"

type (
	// TaskRepository .
	TaskRepository interface {
		GetByResponseTimeMinOrMax(ctx context.Context, bool isNeedMax) (*Task, error)
		GetByID(ctx context.Context, taskID TaskID) (*Task, error)

		Save(ctx context.Context, task *Task) error

		Delete(ctx context.Context, taskID TaskID) error
		DeleteAll(ctx context.Context)
	}

	// HealthService .
	HealthService interface {
		CheckStatus(context.Context, *HealthTask) (*HealthTaskStatus, error)
	}

	// ScheduleTaskService .
	ScheduleTaskService interface {
		Schedule(context.Context, *ScheduleHealthTask) error
		Cancel(context.Context, TaskID) error
		CancelAll(context.Context) error
	}
)
