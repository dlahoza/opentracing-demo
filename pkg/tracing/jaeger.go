package tracing

import (
	"log"
	"os"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/rpcmetrics"
	jexpvar "github.com/uber/jaeger-lib/metrics/expvar"
)

// Init creates a new instance of Jaeger tracer.
func Init(serviceName string) opentracing.Tracer {
	logger := log.New(os.Stderr, "[jaeger]", log.LstdFlags)
	metricsFactory := jexpvar.NewFactory(10) // 10 buckets for histograms
	cfg, err := config.FromEnv()
	if err != nil {
		logger.Fatal("cannot parse Jaeger env vars", err)
	}
	cfg.ServiceName = serviceName
	cfg.Sampler.Type = "const"
	cfg.Sampler.Param = 1

	// TODO(ys) a quick hack to ensure random generators get different seeds, which are based on current time.
	time.Sleep(100 * time.Millisecond)
	jaegerLogger := jaegerLoggerAdapter{logger}

	metricsFactory = metricsFactory.Namespace(serviceName, nil)
	tracer, _, err := cfg.NewTracer(
		config.Logger(jaegerLogger),
		config.Metrics(metricsFactory),
		config.Observer(rpcmetrics.NewObserver(metricsFactory, rpcmetrics.DefaultNameNormalizer)),
	)
	if err != nil {
		logger.Fatal("cannot initialize Jaeger Tracer", err)
	}
	return tracer
}
