package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"reflect"
)

type RedisCli struct {
	redis.UniversalClient
}

func (r *RedisCli) Get(ctx context.Context, key string) (res string, err error) {
	res, err = r.UniversalClient.Get(ctx, key).Result()
	return
}

type Config struct {
	Addrs            []string
	DB               int
	UserName         string
	Password         string
	MasterName       string //when redis sentinel
	SentinelPassword string //when redis sentinel
}

func New(c Config) *RedisCli {
	var r = new(RedisCli)
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:            c.Addrs,
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
func (r *RedisCli) HSetFromStruct(ctx context.Context, key string, data interface{}) *redis.IntCmd {
	mapData := make(map[string]string)
	d := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	for i := 0; i < d.NumField(); i++ {
		mapData[d.Field(i).Name] = fmt.Sprint(v.Field(i).Interface())
	}
	return r.HSet(ctx, key, mapData)
}

// HSetFromStructByPip  使用pipeline把struct按hash结构存入redis
func (r *RedisCli) HSetFromStructByPip(ctx context.Context, pip *redis.Pipeliner, key string, data interface{}) *redis.IntCmd {
	mapData := make(map[string]string)
	d := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	for i := 0; i < d.NumField(); i++ {
		mapData[d.Field(i).Name] = fmt.Sprint(v.Field(i).Interface())
	}
	return (*pip).HSet(ctx, key, mapData)
}

func (r *RedisCli) DistLock(ctx context.Context, key string) (bool, error) {
	return false, nil
}

func (r *RedisCli) DistUnLock(ctx context.Context, key string) (bool, error) {
	return false, nil
}
