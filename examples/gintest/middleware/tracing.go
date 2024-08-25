package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/yunduansing/gtools/opentelemetry/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func ApiRequestTracing(c *gin.Context) {
	tracing.TraceFunc(c.Request.Context(), c.Request.RequestURI, func(span trace.Span) {
		xRequestId := c.GetHeader("X-Request-Id")
		if len(xRequestId) > 0 {
			span.SetAttributes(attribute.String("X-Request-Id", xRequestId))
		} else {
			c.Request.Header.Set("X-Request-Id", span.SpanContext().TraceID().String())
		}
		c.Next()
	})
}
