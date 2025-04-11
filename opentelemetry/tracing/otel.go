package tracing

import (
	"context"
	"github.com/yunduansing/gtools/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const (
	ExporterJaeger = "jaeger"
	ExporterTempo  = "tempo"
)

var tracer trace.Tracer

func InitOtelTracer(serviceName string) {
	// Create a tracer. Usually, tracer is a global variable.
	tracer = otel.Tracer(serviceName)
}

func InitTracer(endpoint, exporter, serviceName, env, id string) {
	switch exporter {
	case ExporterJaeger:
		_, err := tracerProviderJaeger(endpoint, serviceName, env, id)
		if err != nil {
			logger.GetLogger().Panicf(context.Background(), "tracer provider jaeger error:%v", err)
			return
		}
	case ExporterTempo:
		_, err := tranceProviderTempo(context.Background(), endpoint, serviceName)
		if err != nil {
			logger.GetLogger().Panicf(context.Background(), "tracer provider tempo error:%v", err)
			return
		}
	}
}

func TraceFunc(ctx context.Context, spanName string, do func(ctx context.Context, span trace.Span)) {

	// Create a root span (a trace) to measure some operation.
	ctx, span := tracer.Start(ctx, spanName)
	// End the span when the operation we are measuring is done.
	defer span.End()

	do(ctx, span)
}
