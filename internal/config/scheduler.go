package config

import "github.com/spf13/viper"

// SchedulerSettings .
type SchedulerSettings struct {
	// DefaultInterval used when task interval == 0
	DefaultInterval int

	// BucketsCount defines parallel scheduler timers
	// with individual heap.
	BucketsCount int
}

// FromViperSchedulerSettings fill-up configuration structure from viper.
func FromViperSchedulerSettings(v *viper.Viper) *SchedulerSettings {
	const (
		kdefaultinterval = "scheduler.defaults.interval"
		kbucketscount    = "scheduler.bucketscount"
	)

	// Setup defaults
	v.SetDefault(kdefaultinterval, 60000)
	v.SetDefault(kbucketscount, 16)

	// Validations
	if v.GetInt(kdefaultinterval) < 1000 {
		panic(kdefaultinterval + " has too low value")
	}

	if v.GetInt(kbucketscount) < 1 {
		panic(kbucketscount + " has too low value")
	}

	return &SchedulerSettings{
		DefaultInterval: v.GetInt(kdefaultinterval),
		BucketsCount:    v.GetInt(kbucketscount),
	}
}
