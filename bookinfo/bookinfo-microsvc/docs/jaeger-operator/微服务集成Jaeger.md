# 微服务集成Jaeger

安装依赖库
```
go get github.com/uber/jaeger-client-go
go get github.com/opentracing/opentracing-go
```

初始化Jaeger Tracer
```
package main

import (
    "github.com/uber/jaeger-client-go"
    "github.com/uber/jaeger-client-go/config"
    "github.com/opentracing/opentracing-go"
    "log"
)

func InitJaeger(service string) (opentracing.Tracer, error) {
    cfg := &config.Configuration{
        ServiceName: service,
        Sampler: &config.SamplerConfig{
            Type:  "const",
            Param: 1,
        },
        Reporter: &config.ReporterConfig{
            LogSpans:           true,
            LocalAgentHostPort: "localhost:6831", // Jaeger agent 地址
        },
    }

    tracer, _, err := cfg.NewTracer()
    if err != nil {
        return nil, err
    }
    opentracing.SetGlobalTracer(tracer)
    return tracer, nil
}

func main() {
    tracer, err := InitJaeger("ServiceA")
    if err != nil {
        log.Fatal(err)
    }
    defer tracer.Close()

    // 服务逻辑...
}
```

调用时，将trace信息写入到Header
```go
func callServiceB(ctx context.Context) {
    span, _ := opentracing.StartSpanFromContext(ctx, "callServiceB")
    defer span.Finish()

    req, _ := http.NewRequest("GET", "http://service-b/api", nil)
    injectErr := opentracing.GlobalTracer().Inject(
        span.Context(),
        opentracing.HTTPHeaders,
        opentracing.HTTPHeadersCarrier(req.Header),
    )
    if injectErr != nil {
        log.Println("Could not inject span context into header", injectErr)
    }

    // 继续发起请求
    // ...
}
```
接收时，提取trace信息
```go

func handler(w http.ResponseWriter, r *http.Request) {
    wireContext, _ := opentracing.GlobalTracer().Extract(
        opentracing.HTTPHeaders,
        opentracing.HTTPHeadersCarrier(r.Header),
    )
    span := opentracing.StartSpan("handlerB", opentracing.ChildOf(wireContext))
    defer span.Finish()

    // 继续业务逻辑或调用服务 
    // ...
}
```
