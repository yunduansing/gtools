package context

import (
	"context"
	"github.com/yunduansing/gtools/logger"
	"testing"
	"time"
)

func TestNewContext(t *testing.T) {
	ctx := NewContext(context.TODO(), WithRequestId(GenRequestIdByUUID()))
	ctx.Log.Info(ctx.Ctx, "rpc call ", "duration ", time.Second)

	logger.Info(context.TODO(), "aaaaaaaaaa")
}
