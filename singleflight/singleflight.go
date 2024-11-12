package singleflight

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	context2 "github.com/yunduansing/gtools/context"
	redistool "github.com/yunduansing/gtools/redis"
	"github.com/yunduansing/gtools/utils"
	"go.opentelemetry.io/otel/trace"
	"math/rand"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

// SingleFlight 用于实现单次请求处理的逻辑
// 防 缓存击穿、缓存穿透和缓存雪崩
func SingleFlight(ctx context2.Context, store *redistool.Client, key string, res any, fn func() (any, error), lockTimeout time.Duration, cacheTTL time.Duration) (err error) {
	// 首先检查缓存是否存在
	store.WrapDoWithTracing(ctx.Ctx, key, func(ctx1 context.Context, span trace.Span) error {
		var get string
		get, err = store.Get(ctx1, key)
		if err != nil && !errors.Is(err, redis.Nil) {
			ctx.Log.Error(ctx1, "get from redis:", err)
			return err
		}
		if len(get) > 0 && err == nil {
			if get == "null" { //处理value实际不存在的情况
				return nil
			}
			err = json.Unmarshal(utils.StringToByte(get), &res)
			if err != nil {
				ctx.Log.Error(ctx1, "Unmarshal data from redis:", err)
				return err
			}
			return nil
		}

		// 如果缓存不存在，尝试获取锁
		lock := redistool.NewRedisLockWithContext(ctx, store.UniversalClient, key+":lock")
		if lockTimeout > 0 {
			lock.SetExpire(int(lockTimeout.Seconds()))
		}

		var locked bool
		locked, err = lock.Acquire(0, 0)
		if err != nil {
			return err
		}

		if !locked {
			return nil
		}

		defer lock.Release()
		// 获取锁成功，再次检查缓存是否存在
		get, err = store.Get(ctx1, key)
		if err != nil && !errors.Is(err, redis.Nil) {
			ctx.Log.Error(ctx1, "get from redis:", err)
			return err
		}
		if len(get) > 0 && err == nil {
			err = json.Unmarshal(utils.StringToByte(get), &res)
			if err != nil {
				ctx.Log.Error(ctx1, "Unmarshal data from redis:", err)
				return err
			}
			return nil
		}

		// 缓存仍然不存在，执行函数并缓存结果
		res, err = fn()
		if err != nil {
			return err
		}
		if res == nil { //不存在的数据，缓存3-5s
			ttl := random.Intn(6)
			if ttl < 3 {
				ttl = ttl + 3
			}
			cacheTTL = time.Duration(ttl) * time.Second
		}
		if err = store.Set(ctx.Ctx, key, utils.ToJsonString(res), cacheTTL).Err(); err != nil {
			return err
		}

		return nil
	})

	return err
}
