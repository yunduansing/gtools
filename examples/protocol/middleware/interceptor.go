package middleware

import (
	"context"
	"github.com/yunduansing/gtools/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

func UnaryReqTimeInterceptor(
	ctx context.Context, method string, req any, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {

	start := time.Now()

	err := invoker(ctx, method, req, reply, cc, opts...)
	logger.GetLogger().WithField("method", method).
		WithField("duration", time.Since(start).String()).
		WithField("req", req).
		WithField("reply", reply).
		Infof(ctx, "middleware log rpc call")
	return err
}

func UnaryRespTimeServerInterceptor(
	ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.GetLogger().Errorf(ctx, "middleware log rpc call failed with empty metadata")
	} else {
		var requestId = md["requestid"]
		if len(requestId) > 0 {
			ctx = context.WithValue(ctx, "requestId", requestId[0])
		}
	}
	reply, err := handler(ctx, req)
	if err != nil {
		logger.GetLogger().WithField("method", info.FullMethod).
			WithField("duration", time.Since(start).String()).
			WithField("req", req).
			Errorf(ctx, "middleware log rpc call failed with handler error: %v", err)
	} else {
		logger.GetLogger().WithField("method", info.FullMethod).
			WithField("duration", time.Since(start).String()).
			WithField("req", req).
			WithField("reply", reply).
			Infof(ctx, "middleware log rpc call success response")
	}
	return reply, err
}
