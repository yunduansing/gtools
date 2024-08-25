package config

import (
	"github.com/go-redis/redis_rate/v10"
	"github.com/yunduansing/gtools/logger"
	"github.com/yunduansing/gtools/redistool"
)

var (
	Config         ServiceConfig
	Redis          *redistool.Client
	Port           int
	Limiter        *redis_rate.Limiter
	LimitPerSecond int
	IsLimiterOpen  bool
)

func InitConfig() {
	logger.InitLog(logger.Config{})
	initRedis()
	Port = 8080
	LimitPerSecond = 1
	initLimiter()
}

func initRedis() {
	Redis = redistool.New(redistool.Config{
		Addr:     []string{""},
		Password: "",
		DB:       0,
	})
}

func initLimiter() {
	IsLimiterOpen = true
	Limiter = redis_rate.NewLimiter(Redis)
}
