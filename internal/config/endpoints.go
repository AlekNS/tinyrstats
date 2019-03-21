package config

import (
	"regexp"

	"github.com/spf13/viper"
)

// EndpointsSettings configuration for all possible access endpoints.
type EndpointsSettings struct {
	// HTTPBind using for binding to the interface and port.
	HTTPBind string
}

// FromViperEndpointsSettings fill-ups configuration structure from viper.
func FromViperEndpointsSettings(v *viper.Viper) *EndpointsSettings {
	const (
		khttpbind = "endpoints.http.bind"
	)

	// Setup defaults
	v.SetDefault(khttpbind, ":8080")

	// Validations
	// simplified check
	var r = regexp.MustCompile(`^[^:]*:[\d]{2,5}$`)
	if r.MatchString(v.GetString(khttpbind)) == false {
		panic(khttpbind + " has invalid value")
	}

	return &EndpointsSettings{
		HTTPBind: v.GetString(khttpbind),
	}
}
