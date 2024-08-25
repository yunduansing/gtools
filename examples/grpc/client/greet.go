package client

import (
	"github.com/yunduansing/gtools/examples/grpc/middleware"
	grpctool "github.com/yunduansing/gtools/grpc"
	"google.golang.org/grpc"
)

func NewGreetClient() {
	grpctool.Init(grpctool.ClientConfig{
		Address:   "",
		Port:      0,
		ServerPem: "",
	}, grpc.WithUnaryInterceptor(middleware.UnaryReqTimeInterceptor))
}
