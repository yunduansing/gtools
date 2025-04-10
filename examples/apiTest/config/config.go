package config

import (
	"github.com/go-redis/redis_rate/v10"
	"github.com/yunduansing/gtools/logger"
	redistool "github.com/yunduansing/gtools/redis"
)

const (
	EnvDev  = "dev"
	EnvProd = "prod"
	EnvTest = "test"
	EnvUat  = "uat"
)

var (
	Config         ServiceConfig
	Redis          *redistool.Client
	Host           string
	Port           int
	Limiter        *redis_rate.Limiter
	LimitPerSecond int
	IsLimiterOpen  bool
	Uptrace        UptraceConfig
)

func InitConfig() {
	Config = ServiceConfig{
		ServiceName:          "apiTest",
		IsMetricsOpen:        true,
		IsTracingOpen:        true,
		IsRequestLimiterOpen: true,
		Env:                  EnvDev,
	}
	Uptrace = UptraceConfig{
		Version: "v1.0.0",
		Dsn:     "http://project2_secret_token@192.168.2.46:14317/1",
	}
	logger.InitLog(logger.Config{ServiceName: Config.ServiceName, FilePath: "./logs"})
	initRedis()
	Host = "0.0.0.0"
	Port = 8888
	LimitPerSecond = 10000
	initLimiter()
}

func initRedis() {
	Redis = redistool.New(redistool.Config{
		Addr:     []string{"192.168.2.44:16379"},
		Password: "",
		DB:       5,
	})
}

func initLimiter() {
	IsLimiterOpen = true
	Limiter = redis_rate.NewLimiter(Redis)
}
