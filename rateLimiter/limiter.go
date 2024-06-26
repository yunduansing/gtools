package rateLimiter

import (
	"github.com/go-redis/redis_rate/v10"
	"github.com/yunduansing/gtools/redistool"
)

func NewLimiter(rdb *redistool.RedisCli) *redis_rate.Limiter {
	limiter := redis_rate.NewLimiter(rdb)
	return limiter
}
