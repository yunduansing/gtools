package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/yunduansing/gtools/examples/apiTest/config"
	"github.com/yunduansing/gtools/examples/apiTest/middleware"
	"github.com/yunduansing/gtools/examples/apiTest/routers"
	"github.com/yunduansing/gtools/examples/apiTest/service"
	"github.com/yunduansing/gtools/logger"
	"github.com/yunduansing/gtools/opentelemetry/metrics"
	"github.com/yunduansing/gtools/opentelemetry/tracing"
)

func main() {
	config.InitConfig()
	middleware.Init()
	metrics.InitMetricsProvider("192.168.2.46:4318", config.Config.ServiceName, config.Uptrace.Version, "http://project2_secret_token@192.168.2.46:14317/1")
	tracing.InitOtelTracer(config.Config.ServiceName)
	middleware.InitUptrace(config.Config.ServiceName, config.Uptrace.Version, config.Uptrace.Dsn, config.Config.Env)
	service.Init()
	r := gin.Default()
	routers.Register(r)
	logger.GetLogger().WithField("config", config.Config).Info(context.Background(), "run api server")
	err := r.Run()
	if err != nil {
		logger.GetLogger().Error(context.Background(), "run api server error:", err)
		return
	} // listen and serve on 0.0.0.0:8080
}
