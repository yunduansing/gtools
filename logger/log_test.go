package logger

import (
	"context"
	"github.com/yunduansing/gtools/utils"
	"testing"
)

type contextTest struct {
	context.Context
	Log ILog
}

func newContext(ctx context.Context) *contextTest {
	logger.WithField("requestId", ctx.Value("requestId").(string))
	return &contextTest{
		Context: ctx,
		Log:     logger,
	}
}

func TestInitLog(t *testing.T) {
	initLogrusLog()
	ctx := context.WithValue(context.Background(), "requestId", utils.UUID())
	c := newContext(ctx)
	c.Log.Info("下单请求，成功。")
}

func initLogrusLog() {
	InitLog(Config{
		LogType:     "logrus",
		ServiceName: "gtools-logrus-log-test",
		Level:       "Debug",
		FilePath:    "./log/",
	})

}

func initZapLog() {
	InitLog(Config{
		LogType:     "zap",
		ServiceName: "gtools-zap-log-test",
		Level:       "Debug",
		FilePath:    "./log",
	})
}
