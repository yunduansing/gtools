package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/yunduansing/gtools/examples/gintest/config"
	"github.com/yunduansing/gtools/examples/gintest/controllers"
	"github.com/yunduansing/gtools/examples/gintest/middleware"
)

func Register(c *gin.Engine) {
	if !config.Config.IsMetricsOpen {
		c.Use(middleware.ApiRequestCounterMetrics(), middleware.ApiRequestDurationMetrics())
	}
	if !config.Config.IsTracingOpen {
		c.Use(middleware.ApiRequestTracing)
	}
	c.Use(middleware.RequestLimiter)
	c.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	c.GET("/user/", middleware.WrapRequestMiddle(controllers.GetUser))
}
