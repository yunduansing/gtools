package jaeger

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"time"
)

// tracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func tracerProvider(url, serviceName, environment, id string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp,tracesdk.WithBatchTimeout(time.Second)),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName+"-"+environment+":"+id),
			attribute.String("environment", environment),
			attribute.String("ID", id),
		)),
	)
	return tp, nil
}

func NewJaeger(ctx context.Context, url, serviceName, environment, id string) {
	tp, err := tracerProvider(url, serviceName, environment, id)
	if err != nil {
		return
	}

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)

	//ctx, cancel := context.WithCancel(ctx)
	//sig := make(chan os.Signal, 1)
	//select {
	//case <-ctx.Done():
	//case <-sig:
	//
	//}
	//if err := tp.Shutdown(ctx); err != nil {
	//	log.Fatal(err)
	//}
	//cancel()
}

func StartFromContext(ctx context.Context, tracer, spanName string) (context.Context, trace.Span) {
	t := otel.Tracer(tracer)
	return t.Start(ctx, spanName)
}