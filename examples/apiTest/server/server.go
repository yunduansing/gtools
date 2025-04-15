package main

import (
	"apiTest/config"
	"apiTest/middleware"
	"apiTest/routers"
	"apiTest/service"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yunduansing/gtools/logger"
	"github.com/yunduansing/gtools/opentelemetry/metrics"
	"github.com/yunduansing/gtools/opentelemetry/tracing"
	"os"
)

func main() {
	config.InitConfig()
	middleware.InitMetrics()
	metrics.InitMetricsProvider(
		"192.168.2.46:4317", config.Config.ServiceName, config.Uptrace.Version,
		"uptrace-dsn=http://project2_secret_token@192.168.2.46:14317/1",
	)
	os.Setenv("uptrace-dsn", "http://project2_secret_token@192.168.2.46:14317/1")
	tracing.InitTracer("localhost:4318", tracing.ExporterTempo, config.Config.ServiceName, "dev", "id")
	tracing.InitOtelTracer(config.Config.ServiceName)
	middleware.InitUptrace(config.Config.ServiceName, config.Uptrace.Version, config.Uptrace.Dsn, config.Config.Env)
	service.Init()
	r := gin.Default()
	routers.Register(r)
	logger.GetLogger().WithField("config", config.Config).Info(context.Background(), "run api server")
	err := r.Run(fmt.Sprintf("%s:%d", config.Host, config.Port))
	if err != nil {
		logger.GetLogger().Error(context.Background(), "run api server error:", err)
		return
	} // listen and serve on 0.0.0.0:8080
}
