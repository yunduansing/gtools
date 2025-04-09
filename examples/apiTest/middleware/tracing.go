package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/yunduansing/gtools/logger"
	"github.com/yunduansing/gtools/opentelemetry/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func ApiRequestTracing(c *gin.Context) {

	tracing.TraceFunc(
		c.Request.Context(), c.Request.RequestURI, func(ctx context.Context, span trace.Span) {
			xRequestId := c.GetHeader("X-Request-Id")
			if len(xRequestId) > 0 {
				span.SetAttributes(attribute.String("X-Request-Id", xRequestId))
				logger.GetLogger().WithField("TraceId", span.SpanContext().TraceID().String()).Info(
					c.Request.Context(), "请求方携带了request id，因此api请求日志记录下trace id",
				)
			} else {
				xRequestId = span.SpanContext().TraceID().String()
				c.Request.Header.Set("X-Request-Id", xRequestId)
			}
			ctx = context.WithValue(ctx, "requestId", xRequestId)
			c.Request = c.Request.WithContext(ctx)

			c.Next()
		},
	)
}
