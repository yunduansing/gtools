package singleflight

import (
	"fmt"
	context2 "github.com/yunduansing/gtools/context"
	"sync"
	"time"
)

type SingleFlight struct {
	sync.Mutex
	c      map[string]*Call
	locker *redis.Locker
	ctx    context2.Context
}

type Call struct {
	exec func()
}

func NewSingleFlight(ctx context2.Context, cli *redis.Client, key string) *SingleFlight {
	locker := redis.NewRedisLockWithContext(ctx, cli.UniversalClient, key)
	return &SingleFlight{
		c:      make(map[string]*Call),
		locker: locker,
	}
}

// Do shared call,
func (s *SingleFlight) Do(key string, exec func() (res interface{}, err error)) (res interface{}, err error) {
	for {
		acquireOk, err := s.locker.Acquire(0, time.Second)
		if err != nil {
			s.ctx.Log.Error(s.ctx.Ctx, "acquire redis lock err:", err)
			return nil, err
		}
		if !acquireOk {
			continue
		}
		if c, ok := s.c[key]; ok {
			fmt.Println(c)
		}
		return execCall(exec)
	}
}

func execCall(exec func() (res interface{}, err error)) (res interface{}, err error) {
	return exec()
}
