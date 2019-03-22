package external

import (
	"context"
	"net/http"

	"github.com/alekns/tinyrstats/internal/presentation/external/webhttp"

	"github.com/alekns/tinyrstats/internal/config"
	"github.com/alekns/tinyrstats/internal/monitor/app"
	"github.com/alekns/tinyrstats/internal/presentation/external/endpoints"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	grouprun "github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

// BootstrapAndServe .
func BootstrapAndServe(rootContext context.Context,
	settings *config.Settings,
	cmd *cobra.Command,
	rootLogger log.Logger,
	registry app.Registry) {

	logger := log.With(rootLogger, "method", "BootstrapAndServe")
	level.Info(logger).Log("msg", "bootstrap and serve")

	// Prepare Metrics
	var metricsHandler = http.NewServeMux()
	metricsServer := &http.Server{
		Addr:    settings.Endpoints.PromMetricsHTTPBind,
		Handler: metricsHandler,
	}
	metricsHandler.Handle("/metrics", promhttp.Handler())

	// Init endpoints
	endpointSet := endpoints.NewSet(settings.Endpoints, logger, registry)

	// HTTP transport
	httpHandlers := mux.NewRouter()
	webhttp.SetupHTTPServerHandlers(logger, httpHandlers.PathPrefix("/api/v1"), endpointSet)
	httpServer := &http.Server{
		Addr:    settings.Endpoints.HTTPBind,
		Handler: httpHandlers,
		// @TODO: Timeouts
	}

	// Run infrastructure
	var g grouprun.Group
	g.Add(func() error {
		level.Info(logger).Log("transport", "http", "bind", settings.Endpoints.HTTPBind)
		return httpServer.ListenAndServe()
	}, func(error) {
		httpServer.Close()
	})

	g.Add(func() error {
		level.Info(logger).Log("metrics", "prometheus", "bind", settings.Endpoints.PromMetricsHTTPBind)

		return metricsServer.ListenAndServe()
	}, func(error) {
		metricsServer.Shutdown(rootContext)
	})

	g.Add(func() error {
		select {
		case <-rootContext.Done():
			level.Info(logger).Log("msg", "catch context cancellation")
			return rootContext.Err()
		}
		return nil
	}, func(error) {})

	// Run and Serve
	errText := ""
	if err := g.Run(); err != nil {
		errText = err.Error()
	}

	level.Info(logger).Log("msg", "stop", "err", errText)
}
