package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yunduansing/gtools/examples/gintest/config"
	"github.com/yunduansing/gtools/examples/gintest/routers"
	"github.com/yunduansing/gtools/opentelemetry/metrics"
	"github.com/yunduansing/gtools/opentelemetry/tracing"
)

func main() {
	config.InitConfig()
	metrics.InitMetricsProvider()
	tracing.InitOtelTracer(config.Config.ServiceName)
	r := gin.Default()
	routers.Register(r)
	r.Run() // listen and serve on 0.0.0.0:8080
}
