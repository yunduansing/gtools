package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"
)

func InitMetricsProvider() {
	// 初始化一个meterProvider，并设置为全局
	// 实际应用中，应该使用MeterProvider的实现来连接到监控系统
	// 这里使用的是No-op Meters提供的MeterProvider，仅用于示例
	metricProvider := metric.NewMeterProvider()
	otel.SetMeterProvider(metricProvider)

	// 创建一个meter
	//meter := otel.GetMeterProvider().Meter("api-metrics")

}
