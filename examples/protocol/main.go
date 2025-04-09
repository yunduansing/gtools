package main

import (
	"context"
	grpctool "github.com/yunduansing/gtools/grpc"
	"github.com/yunduansing/gtools/logger"
	"github.com/yunduansing/gtools/opentelemetry/tracing"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"protocol/middleware"
	"protocol/server"
	userpb "protocol/user"
)

func main() {
	tracing.InitOtelTracer("grpc test")
	handler := otelgrpc.NewServerHandler(otelgrpc.WithTracerProvider(otel.GetTracerProvider()))
	err := grpctool.Run(
		grpctool.ServerConfig{Port: 8080}, func(s *grpc.Server) {
			userpb.RegisterUserServiceServer(s, &server.Server{})
		}, grpc.StatsHandler(handler), grpc.UnaryInterceptor(middleware.UnaryRespTimeServerInterceptor),
	)

	if err != nil {
		logger.GetLogger().Panic(context.Background(), "run grpc server failed：", err)
		return
	}
}
