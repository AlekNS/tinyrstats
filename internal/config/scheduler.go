package config

import "github.com/spf13/viper"

// SchedulerSettings .
type SchedulerSettings struct {
	// DefaultInterval used when task interval == 0
	DefaultInterval int

	// MaxConcurrency defines parallel scheduler timers
	// with individual heap.
	MaxConcurrency int
}

// FromViperSchedulerSettings fill-up configuration structure from viper.
func FromViperSchedulerSettings(v *viper.Viper) *SchedulerSettings {
	const (
		kdefaultinterval = "scheduler.defaults.interval"
		kmaxconcurrency  = "scheduler.maxconcurrency"
	)

	// Setup defaults
	v.SetDefault(kdefaultinterval, 60)
	v.SetDefault(kmaxconcurrency, 16)

	// Validations
	if v.GetInt(kdefaultinterval) < 1 {
		panic(kdefaultinterval + " has too low value")
	}

	if v.GetInt(kmaxconcurrency) < 1 {
		panic(kmaxconcurrency + " has too low value")
	}

	return &SchedulerSettings{
		DefaultInterval: v.GetInt(kdefaultinterval),
		MaxConcurrency:  v.GetInt(kmaxconcurrency),
	}
}
