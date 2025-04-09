package routers

import (
	"apiTest/config"
	"apiTest/controllers"
	"apiTest/middleware"
	"github.com/gin-gonic/gin"
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

	c.GET("/api/v1/user", middleware.WrapRequestMiddle(controllers.GetUser))
	c.POST("/api/v1/user", middleware.WrapRequestMiddle(controllers.AddOrUpdateUser))
	c.POST("/api/v1/user/login", middleware.WrapRequestMiddle(controllers.UserLogin))
}
