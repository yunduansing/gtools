package httptool

import (
	"context"
	"github.com/yunduansing/gtools/breaker"
	context2 "github.com/yunduansing/gtools/context"
	"github.com/yunduansing/gtools/logger"
	"testing"
	"time"
)

func TestHttpTool_Get(t *testing.T) {
	logger.InitLog(logger.Config{LogType: logger.LogTypeZap})
	c := context2.NewContext(context.TODO(), context2.WithRequestId(context2.GenRequestIdByUUID()))
	cli := NewHttpTool(0, 0, 0, &FastHttpClient{})
	dataBytes, statusCode, err := cli.Get(c.Ctx, "https://google.com", nil, nil)
	t.Log(string(dataBytes), statusCode, err)

	var bc = breaker.Config{
		Sensitivity: 0.9,
		MaxRequest:  0,
	}
	cli = NewHttpTool(1, time.Millisecond*10, time.Second, &FastHttpClient{})
	cli.SetBreaker(breaker.NewBreaker(&bc))
	for i := 0; i < 50; i++ {
		dataBytes, statusCode, err = cli.Get(c.Ctx, "https://google.com", nil, nil)
		t.Log(string(dataBytes), statusCode, err)
	}
}

func TestHttpTool_PostJson(t *testing.T) {
	logger.InitLog(logger.Config{LogType: logger.LogTypeZap})
	c := context2.NewContext(context.TODO(), context2.WithRequestId(context2.GenRequestIdByUUID()))
	cli := NewHttpTool(0, 0, 0, &FastHttpClient{})
	dataBytes, statusCode, err := cli.PostJson(c.Ctx, "https://google.com", nil, nil)
	t.Log(string(dataBytes), statusCode, err)

	var bc = breaker.Config{
		Sensitivity: 0.9,
		MaxRequest:  0,
	}
	cli = NewHttpTool(1, time.Millisecond*10, time.Second, &FastHttpClient{})
	cli.SetBreaker(breaker.NewBreaker(&bc))
	for i := 0; i < 50; i++ {

		dataBytes, statusCode, err = cli.PostJson(c.Ctx, "https://google.com", nil, nil)
		t.Log(string(dataBytes), statusCode, err)
	}
}
