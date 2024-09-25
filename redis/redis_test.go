package redistool

import (
	"context"
	"github.com/yunduansing/gtools/logger"
	"go.opentelemetry.io/otel/trace"
	"testing"
)

func TestNew(t *testing.T) {
	cli := New(Config{
		Addr:     []string{"localhost:6379"},
		Password: "pass",
	})
	r, err := cli.Get(context.Background(), "KEY")
	if err != nil {
		logger.GetLogger().Error(context.Background(), "redis.Get error", err)
		return
	}
	t.Log("result:", r)
}

func TestClient_WrapDoWithTracing(t *testing.T) {
	cli := New(Config{
		Addr:     []string{"localhost:6379"},
		Password: "pass",
	})
	cli.WrapDoWithTracing(context.Background(), "redis.Get", func(ctx context.Context, span trace.Span) error {
		r, err := cli.Get(ctx, "KEY")
		if err != nil {
			logger.GetLogger().Error(ctx, "redis.Get error", err)
			return err
		}
		t.Log("result:", r)
		return nil
	})
}
