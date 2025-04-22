package tools

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

// 初始化jaeger客户端
func InitJaegerClient(serviceName string, host string) error {
	cfg := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: host,
		},
	}
	tracer, _, err := cfg.NewTracer()
	if err != nil {
		return err
	}
	opentracing.SetGlobalTracer(tracer)
	return nil
}
