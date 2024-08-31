package context

import (
	c "context"
	"github.com/yunduansing/gtools/logger"
	"github.com/yunduansing/gtools/utils"
)

type Context struct {
	Ctx         c.Context
	requestId   string
	Log         *logger.Logger
	requestTime string
}

type Option func(Context)

func NewContext(ctx c.Context, opts ...Option) Context {
	myCtx := Context{}
	myCtx.Ctx = ctx
	for _, opt := range opts {
		opt(myCtx)
	}

	myCtx.Log = logger.GetLogger()
	return myCtx
}

func (ctx *Context) GetRequestId() string {
	return ctx.requestId
}

func (ctx *Context) GetRequestTime() string {
	return ctx.requestTime
}

func WithRequestId(requestId string) Option {
	return func(ctx Context) {
		ctx.Ctx = c.WithValue(ctx.Ctx, "requestId", requestId)
		ctx.requestId = requestId
	}
}

func WithRequestTime(requestTime string) Option {
	return func(ctx Context) {
		ctx.Ctx = c.WithValue(ctx.Ctx, "requestTime", requestTime)
		ctx.requestTime = requestTime
	}
}

func GenRequestIdByUUID() string {
	return utils.UUID()
}
