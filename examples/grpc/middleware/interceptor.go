package middleware

import (
	"context"
	"github.com/yunduansing/gtools/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

func UnaryReqTimeInterceptor(ctx context.Context, method string, req any, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	logger.Infof(ctx, "rpc call request duration=%s method=%s req=%+v resp=%+v", time.Since(start), method, req, reply)
	return err
}

func UnaryRespTimeInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Errorf(ctx, "rpc call failed with empty metadata")
	} else {
		var requestId = md["requestid"]
		if len(requestId) > 0 {
			ctx = context.WithValue(ctx, "requestId", requestId[0])
		}
	}
	reply, err := handler(ctx, req)
	if err != nil {
		logger.Errorf(ctx, "rpc call failed with exec handler error: %v", err)
	} else {
		logger.Infof(ctx, "rpc call response duration=%s method=%s req=%+v resp=%+v", time.Since(start), info.FullMethod, req, reply)
	}
	return reply, err
}
