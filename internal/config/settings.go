package config

import (
	"github.com/spf13/viper"
)

// Settings .
type Settings struct {
	Logger    *LoggerSettings
	Tasks     *TasksSettings
	Scheduler *SchedulerSettings
	Endpoints *EndpointsSettings
}

// GetSettings .
func GetSettings(viper *viper.Viper) *Settings {
	return &Settings{
		Logger:    FromViperLoggerSettings(viper),
		Tasks:     FromViperTasksSettings(viper),
		Scheduler: FromViperSchedulerSettings(viper),
		Endpoints: FromViperEndpointsSettings(viper),
	}
}
