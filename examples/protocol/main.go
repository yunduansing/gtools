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
	tracing.InitTracer("localhost:4318", tracing.ExporterTempo, "grpc-test", "dev", "id")
	tracing.InitOtelTracer("grpc-test")
	handler := otelgrpc.NewServerHandler(otelgrpc.WithTracerProvider(otel.GetTracerProvider()))
	err := grpctool.Run(
		grpctool.ServerConfig{Port: 8080}, func(s *grpc.Server) {
			userpb.RegisterUserServiceServer(s, &server.Server{})
		}, grpc.UnaryInterceptor(middleware.UnaryRespTimeServerInterceptor), grpc.StatsHandler(handler),
	)

	if err != nil {
		logger.GetLogger().Panic(context.Background(), "run grpc server failed：", err)
		return
	}
}
