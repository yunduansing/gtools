package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/yunduansing/gtools/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"time"
)

var (
	apiLatencyFloat64Histogram metric.Float64Histogram
	apiRequestCounter          metric.Int64Counter
)

func InitMetrics() {
	meter := otel.GetMeterProvider().Meter("api-metrics")
	var err error
	apiLatencyFloat64Histogram, err = meter.Float64Histogram(
		"api_request_latency",
		metric.WithDescription("API Request Latency"),
		//metric.WithUnit(metric.Seo)
	)
	if err != nil {
		logger.GetLogger().Errorf(context.Background(), "Failed to create histogram:%s", err)

	}

	apiRequestCounter, err = meter.Int64Counter(
		"api_request_counter",
		metric.WithDescription("API Request Counter"),
	)
	if err != nil {
		logger.GetLogger().Errorf(context.Background(), "Failed to create counter:%s", err)

	}
}

func getApiRequestLabels(c *gin.Context, t time.Time) attribute.Set {
	return attribute.NewSet(
		attribute.String("api.method", c.Request.Method),
		attribute.String("api.path", c.Request.URL.Path),
		attribute.String("api.time", t.String()),
		attribute.String("api.requestId", c.Request.Header.Get("X-Request-Id")),
	)
}

func ApiRequestDurationMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		labels := getApiRequestLabels(c, start)

		apiLatencyFloat64Histogram.Record(
			c.Request.Context(),
			time.Since(start).Seconds(),
			metric.WithAttributeSet(labels),
		)
	}
}

func ApiRequestCounterMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		labels := getApiRequestLabels(c, time.Now())
		apiRequestCounter.Add(c.Request.Context(), 1, metric.WithAttributeSet(labels))
		c.Next()
	}
}
