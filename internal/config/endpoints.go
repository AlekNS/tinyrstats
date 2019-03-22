package config

import (
	"regexp"

	"github.com/spf13/viper"
)

// EndpointsSettings configuration for all possible access endpoints.
type EndpointsSettings struct {
	// HTTPBind using for binding to the interface and port.
	HTTPBind string
	// PromMetricsHTTPBind using for prometheus metrics.
	PromMetricsHTTPBind string
}

// FromViperEndpointsSettings fill-ups configuration structure from viper.
func FromViperEndpointsSettings(v *viper.Viper) *EndpointsSettings {
	const (
		khttpbind        = "endpoints.http.bind"
		kprommetricsbind = "endpoints.metrics.prometheusbind"
	)

	// Setup defaults
	v.SetDefault(khttpbind, ":8080")
	v.SetDefault(kprommetricsbind, ":8081")

	// Validations
	// simplified check
	var r = regexp.MustCompile(`^[^:]*:[\d]{2,5}$`)

	if r.MatchString(v.GetString(khttpbind)) == false {
		panic(khttpbind + " has invalid value")
	}
	if r.MatchString(v.GetString(kprommetricsbind)) == false {
		panic(kprommetricsbind + " has invalid value")
	}

	return &EndpointsSettings{
		HTTPBind:            v.GetString(khttpbind),
		PromMetricsHTTPBind: v.GetString(kprommetricsbind),
	}
}
