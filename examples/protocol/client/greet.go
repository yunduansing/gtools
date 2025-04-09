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
	greetClient *grpc.ClientConn
	greetOnce   sync.Once
)

func NewGreetClient() *grpc.ClientConn {
	userOnce.Do(
		func() {
			conn, err := grpctool.Init(
				grpctool.ClientConfig{
					Address:   "",
					Port:      0,
					ServerPem: "",
				}, grpc.WithUnaryInterceptor(middleware.UnaryReqTimeInterceptor),
			)
			if err != nil {
				logger.GetLogger().Panic(context2.Background(), "create user protocol client conn error:", err)
				panic(err)
			}
			greetClient = conn
		},
	)
	return greetClient
}
