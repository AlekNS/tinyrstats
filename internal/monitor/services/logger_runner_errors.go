package services

import (
	"github.com/alekns/tinyrstats/pkg/helpers/runner"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type (
	loggerRunnerErrors struct {
		logger log.Logger
	}
)

// OnProduceError .
func (le *loggerRunnerErrors) OnError(err error) error {
	level.Error(le.logger).Log("err", err.Error())
	return err
}

// NewLoggerRunnerException .
func NewLoggerRunnerException(logger log.Logger,
	serviceName string) runner.ErrorHandler {

	return &loggerRunnerErrors{
		logger: log.With(logger, "service", serviceName),
	}
}
