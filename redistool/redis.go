package redistool

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/yunduansing/gtools/opentelemetry/tracing"
	"go.opentelemetry.io/otel/trace"
	"reflect"
)

type Client struct {
	redis.UniversalClient
}

func (r *Client) Get(ctx context.Context, key string) (res string, err error) {
	res, err = r.UniversalClient.Get(ctx, key).Result()
	return
}

type Config struct {
	Addr             []string
	DB               int
	UserName         string
	Password         string
	MasterName       string //when redis sentinel
	SentinelPassword string //when redis sentinel
}

func New(c Config) *Client {
	var r = new(Client)
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:            c.Addr,
		DB:               c.DB,
		Username:         c.UserName,
		Password:         c.Password,
		SentinelPassword: c.SentinelPassword,
		MasterName:       c.MasterName,
	})
	r.UniversalClient = rdb
	return r
}

// HSetFromStruct 把struct按hash结构存入redis
func (r *Client) HSetFromStruct(ctx context.Context, key string, data interface{}) *redis.IntCmd {
	mapData := make(map[string]string)
	d := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	for i := 0; i < d.NumField(); i++ {
		mapData[d.Field(i).Name] = fmt.Sprint(v.Field(i).Interface())
	}
	return r.HSet(ctx, key, mapData)
}

// HSetFromStructByPip  使用pipeline把struct按hash结构存入redis
func (r *Client) HSetFromStructByPip(ctx context.Context, pip *redis.Pipeliner, key string, data interface{}) *redis.IntCmd {
	mapData := make(map[string]string)
	d := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	for i := 0; i < d.NumField(); i++ {
		mapData[d.Field(i).Name] = fmt.Sprint(v.Field(i).Interface())
	}
	return (*pip).HSet(ctx, key, mapData)
}

// WrapDoWithTracing  使用链路追踪
//
// @spanName: 链路名称
//
// @fn: func(ctx context.Context, span trace.Span)
//
// example:
//
//	cli.WrapDoWithTracing(context.Background(), "redis.Get", func(ctx context.Context, span trace.Span) {
//			r, err := cli.Get(ctx, "KEY")
//			if err != nil {
//				logger.GetLogger().Error(ctx, "redis.Get error", err)
//				span.RecordError(err)
//				return
//			}
//			t.Log("result:", r)
//		})
func (r *Client) WrapDoWithTracing(ctx context.Context, spanName string, fn func(ctx context.Context, span trace.Span) error) {
	tracing.TraceFunc(ctx, spanName, func(ctx context.Context, span trace.Span) {
		err := fn(ctx, span)
		if err != nil {
			span.RecordError(err)
		}

	})
}
