package context

import (
	c "context"
	"github.com/yunduansing/gtools/logger"
	"github.com/yunduansing/gtools/utils"
	"google.golang.org/grpc/metadata"
)

type Context struct {
	Ctx         c.Context
	requestId   string
	Log         *logger.Logger
	requestTime string
}

type Option func(*Context)

func NewContext(ctx c.Context, opts ...Option) Context {
	myCtx := Context{}
	myCtx.Ctx = ctx
	for _, opt := range opts {
		opt(&myCtx)
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
	return func(ctx *Context) {
		ctx.Ctx = c.WithValue(ctx.Ctx, "requestId", requestId)
		ctx.requestId = requestId
	}
}

func WithRequestTime(requestTime string) Option {
	return func(ctx *Context) {
		ctx.Ctx = c.WithValue(ctx.Ctx, "requestTime", requestTime)
		ctx.requestTime = requestTime
	}
}

func GenRequestIdByUUID() string {
	return utils.UUID()
}

// cc 是最终返回的 context
func (ctx *Context) GetContext() c.Context {
	// 拿原来的 context
	baseCtx := ctx.Ctx

	// 从已有 context 提取 metadata（保留 traceparent、grpc-trace-bin 等）
	md, ok := metadata.FromOutgoingContext(baseCtx)
	if !ok {
		md = metadata.MD{}
	} else {
		md = md.Copy()
	}

	// 追加你自己的业务字段
	md.Set("requestid", ctx.GetRequestId())
	md.Set("requesttime", ctx.GetRequestTime())

	// 合并后的 context 返回
	return metadata.NewOutgoingContext(baseCtx, md)
}
