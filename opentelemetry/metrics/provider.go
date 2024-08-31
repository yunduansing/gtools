package metrics

import (
	"github.com/yunduansing/gtools/logger"
	"golang.org/x/net/context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

func InitMetricsProvider(exporterEndpoint, serviceName, version string) {
	// 初始化一个meterProvider，并设置为全局
	// 实际应用中，应该使用MeterProvider的实现来连接到监控系统
	// 这里使用的是No-op Meters提供的MeterProvider，仅用于示例
	ctx := context.Background()
	exporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint(exporterEndpoint),
		otlpmetrichttp.WithInsecure(),
		otlpmetrichttp.WithCompression(otlpmetrichttp.GzipCompression),
	)
	if err != nil {
		panic(err)
	}

	reader := sdkmetric.NewPeriodicReader(
		exporter,
		sdkmetric.WithInterval(15*time.Second),
	)

	resource, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("service.version", version),
		))
	if err != nil {
		logger.GetLogger().Panic(ctx, "create new Meter Provider error：", err)
		panic(err)
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetric.WithResource(resource),
	)
	otel.SetMeterProvider(provider)

}
