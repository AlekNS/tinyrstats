package logger

import (
	"github.com/go-kit/kit/log"
	kitlevel "github.com/go-kit/kit/log/level"
)

// SetLevelLogger setup gokit logger with specified level and default fields.
func SetLevelLogger(logger log.Logger, level string) log.Logger {
	if level == "debug" {
		logger = kitlevel.NewFilter(logger, kitlevel.AllowDebug())
	} else if level == "info" {
		logger = kitlevel.NewFilter(logger, kitlevel.AllowInfo())
	} else if level == "warn" {
		logger = kitlevel.NewFilter(logger, kitlevel.AllowWarn())
	} else {
		logger = kitlevel.NewFilter(logger, kitlevel.AllowError())
	}

	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	return logger
}
