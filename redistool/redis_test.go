package redistool

import (
	"context"
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	cli := New(Config{
		Addr:     []string{"localhost:6379"},
		Password: "pass",
	})
	r, err := cli.Get(context.Background(), "KEY")
	log.Print("result:", r, " err:", err)
}
