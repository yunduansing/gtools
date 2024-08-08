package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	randomLen       = 16
	tolerance       = 500 // milliseconds
	millisPerSecond = 1000
	lockCommand     = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
    return "OK"
else
    return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
end`
	delCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`
)

// A RedisLock is a redis lock.
type RedisLock struct {
	store     redis.Cmdable
	seconds   uint32
	key       string
	id        string
	ctx       context2.Context
	IsDebug   bool
	stopRenew chan struct{}
}

func init() {
	rand.NewSource(time.Now().UnixNano())
}

// NewRedisLock returns a RedisLock.
func NewRedisLock(store redis.Cmdable, key string) *RedisLock {
	return &RedisLock{
		store: store,
		key:   key,
		id:    Randn(randomLen),
		ctx:   context2.NewContext(context.Background()),
	}
}

// NewRedisLockWithContext returns a RedisLock.
func NewRedisLockWithContext(ctx context2.Context, store redis.Cmdable, key string) *RedisLock {
	return &RedisLock{
		store: store,
		key:   key,
		id:    Randn(randomLen),
		ctx:   ctx,
	}
}

// Acquire acquires the lock.
//
// @repeat  重试次数，默认不重试
//
// @wait 重试(repeat>1)时，重试的间隔时间  默认10ms
func (rl *RedisLock) Acquire(repeat int, wait time.Duration) (bool, error) {
	if wait == 0 {
		wait = 10 * time.Millisecond // 默认等待时间
	}

	// 使用for循环进行重试逻辑
	for attempt := 0; attempt <= repeat; attempt++ {
		if attempt > 0 {
			// 使用time.After等待，而不是time.Sleep
			<-time.After(wait)
		}

		seconds := atomic.LoadUint32(&rl.seconds)
		start := time.Now()
		resp, err := rl.store.Eval(context.Background(), lockCommand, []string{rl.key}, []string{
			rl.id, strconv.Itoa(int(seconds)*millisPerSecond + tolerance),
		}).Result()
		if rl.IsDebug {
			rl.ctx.Log.Infof("acquiring lock for key=%s  acquireCount=%d  time cost=%s", rl.key, attempt+1, time.Since(start))
		}

		if errors.Is(err, redis.Nil) {
			// 锁未被获取，继续尝试
			continue
		} else if err != nil {
			rl.ctx.Log.Errorf("Error on acquiring lock for key=%s, err=%s", rl.key, err)
			//return false, err
			continue
		} else if resp == nil {
			rl.ctx.Log.Errorf("Error on acquiring lock for key=%s, err=%s", rl.key, err)
			continue
		}

		reply, ok := resp.(string)
		if ok && reply == "OK" {
			if repeat > 0 {
				rl.ctx.Log.Infof("Success on acquiring lock for key=%s,acquireCount=%d", rl.key, attempt+1)
			}
			return true, nil // 成功获取锁
		}

		rl.ctx.Log.Errorf("Unknown reply when acquiring lock for key=%s: resp=%+v", rl.key, resp)
	}

	return false, nil // 超过重试次数，未能获取锁

}

// AcquireBackoff acquires the lock using 指数退避策略.
//
// @repeat  重试次数，默认不重试
//
// @baseDelay  重试(repeat>1)时，基础延迟时间
//
// @maxDelay 重试(repeat>0)时，最大延迟时间，默认10ms
func (rl *RedisLock) AcquireBackoff(repeat int, baseDelay, maxDelay time.Duration) (bool, error) {
	if baseDelay == 0 {
		baseDelay = time.Millisecond // 默认基础延迟等待时间
	}
	baseDelayTmp := baseDelay
	if maxDelay == 0 {
		maxDelay = 10 * time.Millisecond // 默认最大延迟等待时间
	}

	for attempt := 0; attempt <= repeat; attempt++ {
		if attempt > 0 {
			// 计算退避时间
			delay := time.Duration(rand.Intn(int(baseDelay))) + baseDelay
			if delay > maxDelay {
				delay = maxDelay
			}
			<-time.After(delay)
			if baseDelay >= maxDelay {
				baseDelay = baseDelayTmp
			} else {
				baseDelay *= 2 // 指数增加
			}
		}

		seconds := atomic.LoadUint32(&rl.seconds)
		start := time.Now()
		resp, err := rl.store.Eval(context.Background(), lockCommand, []string{rl.key}, []string{
			rl.id, strconv.Itoa(int(seconds)*millisPerSecond + tolerance),
		}).Result()
		if rl.IsDebug {
			rl.ctx.Log.Infof("acquiring lock for key=%s  acquireCount=%d  time cost=%s", rl.key, attempt+1, time.Since(start))
		}

		if errors.Is(err, redis.Nil) {
			// 锁未被获取，继续尝试
			continue
		} else if err != nil {
			rl.ctx.Log.Errorf("Error on acquiring lock for key=%s, err=%s", rl.key, err)
			continue
		} else if resp == nil {
			rl.ctx.Log.Errorf("Error on acquiring lock for key=%s, err=%s", rl.key, err)
			continue
		}

		reply, ok := resp.(string)
		if ok && reply == "OK" {
			if repeat > 0 {
				rl.ctx.Log.Infof("Success on acquiring lock for key=%s,acquireCount=%d", rl.key, attempt+1)
			}

			return true, nil // 成功获取锁
		}

		rl.ctx.Log.Errorf("Unknown reply when acquiring lock for key=%s: resp=%+v", rl.key, resp)
	}

	return false, nil // 超过重试次数，未能获取锁
}

// Release releases the lock.
func (rl *RedisLock) Release() (bool, error) {
	if rl.stopRenew != nil {
		close(rl.stopRenew)
	}
	resp, err := rl.store.Eval(context.Background(), delCommand, []string{rl.key}, []string{rl.id}).Result()
	if err != nil {
		return false, err
	}

	reply, ok := resp.(int64)
	if !ok {
		return false, nil
	}

	return reply == 1, nil
}

// SetExpire sets the expiration.
func (rl *RedisLock) SetExpire(seconds int) {
	atomic.StoreUint32(&rl.seconds, uint32(seconds))
}
