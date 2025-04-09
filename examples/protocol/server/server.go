package server

import (
	"context"
	"github.com/yunduansing/gtools/logger"
	"github.com/yunduansing/gtools/opentelemetry/tracing"
	"go.opentelemetry.io/otel/trace"
	userpb "protocol/user"
)

type Server struct {
	userpb.UnimplementedUserServiceServer
}

func (s *Server) mustEmbedUnimplementedUserServiceServer() {
	//TODO implement me
	panic("implement me")
}

func (s *Server) GetUser(ctx context.Context, req *userpb.GetUserReq) (*userpb.UserResponse, error) {
	tracing.TraceFunc(
		ctx, "UserServer", func(ctx context.Context, span trace.Span) {
			logger.GetLogger().WithField("traceId", span.SpanContext().TraceID().String()).Info(ctx, "call getuser")
		},
	)
	resp := &userpb.UserResponse{}
	resp.Data = &userpb.UserResponse_UserData{
		UserData: &userpb.User{
			UserId:   1,
			UserName: "lucky",
			Phone:    "18888888888",
		},
	}
	return resp, nil
}
