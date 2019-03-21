package runner

//go:generate mockgen -source=interfaces.go -package=runner -destination=interfaces_mocks.go

import (
	"context"
)

type (
	// Consumer receive values from processor.
	Consumer interface {
		Accept(ctx context.Context, results ...interface{}) error
	}

	// ErrorHandler react on the error incidents.
	ErrorHandler interface {
		OnError(err error) error
	}

	// Processor for processing task.
	Processor interface {
		Enqueue(ctx context.Context, task interface{}) error
	}
)
