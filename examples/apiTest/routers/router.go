package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/yunduansing/gtools/examples/apiTest/config"
	"github.com/yunduansing/gtools/examples/apiTest/controllers"
	"github.com/yunduansing/gtools/examples/apiTest/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func Register(c *gin.Engine) {
	c.Use(otelgin.Middleware("api test"))
	if config.Config.IsMetricsOpen {
		c.Use(middleware.ApiRequestCounterMetrics(), middleware.ApiRequestDurationMetrics())
	}
	if config.Config.IsTracingOpen {
		c.Use(middleware.ApiRequestTracing)
	}
	if config.Config.IsRequestLimiterOpen {
		c.Use(middleware.RequestLimiter)
	}
	c.GET(
		"/ping", func(c *gin.Context) {
			c.JSON(
				200, gin.H{
					"message": "pong",
				},
			)
		},
	)

	c.GET("/api/v1/user/", middleware.WrapRequestMiddle(controllers.GetUser))
	c.POST("/api/v1/user/", middleware.WrapRequestMiddle(controllers.AddOrUpdateUser))
	c.POST("/api/v1/user/login", middleware.WrapRequestMiddle(controllers.UserLogin))
}
