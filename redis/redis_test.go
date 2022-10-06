package redis

import (
	"context"
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	cli := New(Config{
		Addrs:    []string{"localhost:6379"},
		Password: "pass",
	})
	r, err := cli.Get(context.Background(), "KEY")
	log.Print("result:", r, " err:", err)
}
