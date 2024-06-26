package sharedcall

import (
	"context"
	"fmt"
	"github.com/yunduansing/gtools/redistool"
	"sync"
)

type SharedCall struct {
	sync.Mutex
	c map[string]*Call
	r *redistool.RedisCli
}

type Call struct {
	exec func()
}

func NewSharedCall(redisCli *redistool.RedisCli) *SharedCall {
	return &SharedCall{
		c: make(map[string]*Call),
		r: redisCli,
	}
}

// Do shared call,
func (s *SharedCall) Do(ctx context.Context, key string, exec func() (res interface{}, err error)) (res interface{}, err error) {
	ok, err := s.r.DistLock(ctx, key)
	if err != nil {
		return nil, err
	}
	if !ok {

	}
	if c, ok := s.c[key]; ok {
		fmt.Println(c)
	}
	return execCall(exec)
}

func execCall(exec func() (res interface{}, err error)) (res interface{}, err error) {
	return exec()
}
