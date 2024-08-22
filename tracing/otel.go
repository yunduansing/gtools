package tracing

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func InitOtelTracer(serviceName string) {
	// Create a tracer. Usually, tracer is a global variable.
	tracer = otel.Tracer(serviceName)
}

func TraceFunc(ctx context.Context, spanName string, do func(span trace.Span)) {

	// Create a root span (a trace) to measure some operation.
	_, span := tracer.Start(ctx, spanName)
	// End the span when the operation we are measuring is done.
	defer span.End()

	do(span)
}
