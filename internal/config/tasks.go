package config

import "github.com/spf13/viper"

// TasksSettings .
type TasksSettings struct {
	DefaultTimeout int

	MaxPending     int
	TaskQueueSize  int
	MaxConcurrency int

	RepositoryBucketsCount int
}

// FromViperTasksSettings fill-up configuration structure from viper.
func FromViperTasksSettings(v *viper.Viper) *TasksSettings {
	const (
		kdefaulttimeout = "tasks.defaults.timeout"

		kmaxpending     = "tasks.maxpending"
		ktaskqueue      = "tasks.taskqueuesize"
		kmaxconcurrency = "tasks.maxconcurrency"

		krepbuckets = "tasks.repository.bucketscount"
	)

	// Setup defaults
	v.SetDefault(kdefaulttimeout, 5000)

	v.SetDefault(kmaxpending, 512)
	v.SetDefault(ktaskqueue, 256)
	v.SetDefault(kmaxconcurrency, 128)

	v.SetDefault(krepbuckets, 32)

	// Validations
	if v.GetInt(kdefaulttimeout) < 10 {
		panic(kdefaulttimeout + " has too low value")
	}
	if v.GetInt(kmaxpending) < 1 {
		panic(kmaxpending + " has too low value")
	}
	if v.GetInt(ktaskqueue) < 1 {
		panic(ktaskqueue + " has too low value")
	}
	if v.GetInt(kmaxconcurrency) < 1 {
		panic(kmaxconcurrency + " has too low value")
	}
	if v.GetInt(krepbuckets) < 1 {
		panic(krepbuckets + " has too low value")
	}

	return &TasksSettings{
		DefaultTimeout: v.GetInt(kdefaulttimeout),

		MaxPending:     v.GetInt(kmaxpending),
		TaskQueueSize:  v.GetInt(ktaskqueue),
		MaxConcurrency: v.GetInt(kmaxconcurrency),

		RepositoryBucketsCount: v.GetInt(krepbuckets),
	}
}
