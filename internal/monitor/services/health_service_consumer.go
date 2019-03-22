package services

import (
	"context"

	"github.com/alekns/tinyrstats/internal/monitor"

	"github.com/alekns/tinyrstats/pkg/helpers/runner"
)

type (
	healthServiceConsumer struct {
		healthService           monitor.HealthService
		requestResponseConsumer runner.Consumer
	}
)

// Accept .
func (ct *healthServiceConsumer) Accept(ctx context.Context, results ...interface{}) error {
	result, err := ct.healthService.CheckStatus(ctx, results[0].(*monitor.HealthTask))
	if err != nil {
		return err
	}

	return ct.requestResponseConsumer.Accept(ctx, results[0], result)
}

// NewHealthServiceConsumer .
func NewHealthServiceConsumer(healthService monitor.HealthService,
	requestResponseConsumer runner.Consumer) runner.Consumer {
	return &healthServiceConsumer{
		healthService:           healthService,
		requestResponseConsumer: requestResponseConsumer,
	}
}
