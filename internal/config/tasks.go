package config

import "github.com/spf13/viper"

// TasksSettings configure tasks.
type TasksSettings struct {
	// DefaultTimeout is used when timeout for task was not specified.
	DefaultTimeout int

	// MaxPending is task count that could be in pending.
	MaxPending int
	// QueueSize is buffer channel size.
	QueueSize int
	// MaxConcurrency is count of concurrent workers.
	MaxConcurrency int

	// RepositoryBucketsCount need to increase concurrency.
	RepositoryBucketsCount int
}

// FromViperTasksSettings fill-up configuration structure from viper.
func FromViperTasksSettings(v *viper.Viper) *TasksSettings {
	const (
		kdefaulttimeout = "tasks.defaults.timeout"

		kmaxpending     = "tasks.maxpending"
		kqueuesize      = "tasks.queuesize"
		kmaxconcurrency = "tasks.maxconcurrency"

		krepbuckets = "tasks.repository.bucketscount"
	)

	// Setup defaults
	v.SetDefault(kdefaulttimeout, 5000)

	v.SetDefault(kmaxpending, 512)
	v.SetDefault(kqueuesize, 256)
	v.SetDefault(kmaxconcurrency, 128)

	v.SetDefault(krepbuckets, 32)

	// Validations
	if v.GetInt(kdefaulttimeout) < 10 {
		panic(kdefaulttimeout + " has too low value")
	}
	if v.GetInt(kmaxpending) < 1 {
		panic(kmaxpending + " has too low value")
	}
	if v.GetInt(kqueuesize) < 1 {
		panic(kqueuesize + " has too low value")
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
		QueueSize:      v.GetInt(kqueuesize),
		MaxConcurrency: v.GetInt(kmaxconcurrency),

		RepositoryBucketsCount: v.GetInt(krepbuckets),
	}
}
