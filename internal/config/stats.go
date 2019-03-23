package config

import "github.com/spf13/viper"

// StatsSettings configure stats application.
type StatsSettings struct {
	// BucketsCount need to increase access concurrency.
	BucketsCount int
}

// FromViperStatsSettings fill-up configuration structure from viper.
func FromViperStatsSettings(v *viper.Viper) *StatsSettings {
	const (
		kbuckets = "stats.bucketscount"
	)

	// Setup defaults
	v.SetDefault(kbuckets, 32)

	// Validations
	if v.GetInt(kbuckets) < 1 {
		panic(kbuckets + " has too low value")
	}

	return &StatsSettings{
		BucketsCount: v.GetInt(kbuckets),
	}
}
