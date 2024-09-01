package rateLimiter

import (
	"github.com/go-redis/redis_rate/v10"
)

func NewLimiter(rdb *redis.RedisCli) *redis_rate.Limiter {
	limiter := redis_rate.NewLimiter(rdb)
	return limiter
}
