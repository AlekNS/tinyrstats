package tracer

import (
	"io"

	"github.com/go-kit/kit/log"

	opentracing "github.com/opentracing/opentracing-go"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	logkit "github.com/uber/jaeger-lib/client/log/go-kit"
)

// NewJaeger initializes jarger tracer.
func NewJaeger(logger log.Logger, serviceName string) (opentracing.Tracer, io.Closer, error) {
	cfg, err := jaegerconfig.FromEnv()

	if err != nil {
		return nil, nil, err
	}

	tracer, closer, err := cfg.New(
		serviceName,
		jaegerconfig.Logger(logkit.NewLogger(logger)),
	)
	if err != nil {
		return nil, nil, err
	}

	return tracer, closer, nil
}
