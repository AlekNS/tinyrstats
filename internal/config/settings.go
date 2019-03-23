package config

import (
	"github.com/spf13/viper"
)

// Settings gathering all application settings.
type Settings struct {
	Logger    *LoggerSettings
	Tasks     *TasksSettings
	Scheduler *SchedulerSettings
	Stats     *StatsSettings
	Endpoints *EndpointsSettings
}

// GetSettings reads all from config file and env.
func GetSettings(viper *viper.Viper) *Settings {
	return &Settings{
		Logger:    FromViperLoggerSettings(viper),
		Tasks:     FromViperTasksSettings(viper),
		Scheduler: FromViperSchedulerSettings(viper),
		Stats:     FromViperStatsSettings(viper),
		Endpoints: FromViperEndpointsSettings(viper),
	}
}
