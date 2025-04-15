package context

import (
	c "context"
	"github.com/yunduansing/gtools/logger"
	"github.com/yunduansing/gtools/utils"
	"go.opentelemetry.io/otel/propagation"
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

	// 确保 OpenTelemetry 的传播信息也被注入
	baseCtx = metadata.NewOutgoingContext(baseCtx, md)

	// 显式调用 OpenTelemetry 的 propagator 注入
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	carrier := propagation.HeaderCarrier{}
	propagator.Inject(baseCtx, carrier)

	// 将 OpenTelemetry 的 headers 合并到 metadata
	for k, v := range carrier {
		md[k] = v
	}

	return metadata.NewOutgoingContext(baseCtx, md)
}
