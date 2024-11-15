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
// 防 缓存击穿、缓存穿透
func SingleFlight(ctx context2.Context, store *redistool.Client, key string, res any, fn func() (any, error), lockTimeout time.Duration, cacheTTL time.Duration) (err error) {
	store.WrapDoWithTracing(ctx.Ctx, key, func(ctx1 context.Context, span trace.Span) error {
		var exists, locked bool
		// 检查缓存
		if exists, err = checkCache(ctx, store, key, &res); exists || err != nil {
			return err
		}

		// 尝试获取锁
		lock := redistool.NewRedisLockWithContext(ctx, store.UniversalClient, key+":lock")
		if lockTimeout > 0 {
			lock.SetExpire(int(lockTimeout.Seconds()))
		}

		locked, err = lock.Acquire(100, 0)
		if err != nil {
			ctx.Log.Errorf(ctx1, "get lock %s:lock fail:%s", key, err)
			return err
		}
		if !locked {
			return nil
		}
		defer lock.Release()

		// 再次检查缓存
		if exists, err = checkCache(ctx, store, key, &res); exists || err != nil {
			return err
		}

		// 执行函数获取结果并缓存
		res, err = fn()
		if err != nil {
			return err
		}
		if res == nil {
			cacheTTL = time.Duration(random.Intn(3)+3) * time.Second
		}
		return store.Set(ctx.Ctx, key, utils.ToJsonString(res), cacheTTL).Err()
	})
	return err
}

// checkCache 封装缓存检查逻辑，避免重复代码
func checkCache(ctx context2.Context, store *redistool.Client, key string, res any) (exists bool, err error) {
	get, err := store.Get(ctx.Ctx, key)
	if err != nil && !errors.Is(err, redis.Nil) {
		ctx.Log.Errorf(ctx.Ctx, "checkCache %s err:%s", key, err)
		return false, err
	}
	if get == "null" {
		return true, nil
	}
	if errors.Is(err, redis.Nil) || len(get) == 0 {
		return false, nil
	}
	err = json.Unmarshal(utils.StringToByte(get), &res)
	return err == nil, err
}
