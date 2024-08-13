package redistool

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
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
