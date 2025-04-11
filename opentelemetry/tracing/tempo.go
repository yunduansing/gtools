package tracing

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
)

type TempoConfig struct {
	Endpoint    string `json:"endpoint"`
	ServiceName string `json:"serviceName"`
}

func tracerProviderTempo(ctx context.Context, endpoint, serviceName string) (*sdktrace.TracerProvider, error) {
	// 创建 OTLP Trace Exporter（通过 HTTP 协议发送）
	exporter, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithInsecure(), // Tempo 默认无 TLS
		otlptracehttp.WithEndpoint(endpoint),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}
