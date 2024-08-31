package client

import (
	context2 "context"
	"github.com/yunduansing/gtools/examples/protocol/middleware"
	grpctool "github.com/yunduansing/gtools/grpc"
	"github.com/yunduansing/gtools/logger"
	"google.golang.org/grpc"
	"sync"
)

var (
	userClient *grpc.ClientConn
	userOnce   sync.Once
)

func NewUserClient() *grpc.ClientConn {
	userOnce.Do(func() {
		conn, err := grpctool.Init(grpctool.ClientConfig{
			Address:   "",
			Port:      0,
			ServerPem: "",
		}, grpc.WithUnaryInterceptor(middleware.UnaryReqTimeInterceptor))
		if err != nil {
			logger.Panic(context2.Background(), "create user protocol client conn error:", err)
			panic(err)
		}
		userClient = conn
	})
	return userClient
}
