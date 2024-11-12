package singleflight

import (
	"context"
	context2 "github.com/yunduansing/gtools/context"
	"github.com/yunduansing/gtools/opentelemetry/tracing"
	redistool "github.com/yunduansing/gtools/redis"
	"testing"
	"time"
)

type User struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func TestSingleFlight(t *testing.T) {
	tracing.InitOtelTracer("TestSingleFlight")
	rdb := redistool.New(redistool.Config{
		Addr:             []string{"192.168.6.41:6371", "192.168.6.42:6373", "192.168.6.43:6375"},
		DB:               0,
		UserName:         "",
		Password:         "123456",
		MasterName:       "",
		SentinelPassword: "",
	})

	key := "test:user:1"

	var user User

	err := SingleFlight(context2.NewContext(context.Background()), rdb, key, &user, func() (r any, err error) {
		<-time.After(time.Second)
		return nil, nil
	}, time.Second, time.Minute)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(user)
}
