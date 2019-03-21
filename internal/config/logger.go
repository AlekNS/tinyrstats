package config

import (
	"github.com/spf13/viper"
)

// LoggerSettings .
type LoggerSettings struct {
	// ConsoleLevel defines logger verbosity.
	ConsoleLevel string
}

// FromViperLoggerSettings fill-up configuration structure from viper.
func FromViperLoggerSettings(v *viper.Viper) *LoggerSettings {
	const (
		kconsolelevel = "logging.console.level"
	)

	// Setup defaults
	v.SetDefault(kconsolelevel, "info")

	// Validations
	level := v.GetString(kconsolelevel)
	switch level {
	case "debug", "info", "warn", "error":
	default:
		panic(kconsolelevel + " has invalid value")
	}

	return &LoggerSettings{
		ConsoleLevel: level,
	}
}
