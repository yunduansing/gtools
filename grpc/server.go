package grpctool

import (
	"fmt"
	"github.com/yunduansing/gtools/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
)

type ServerConfig struct {
	Port   int
	TlsPem string
	TlsKey string
}

type GrpcServerHandler func(server *grpc.Server)

func RunWithTls(c ServerConfig, servers ...GrpcServerHandler) error {
	lis, err := net.Listen("tcp", ":"+fmt.Sprint(c.Port))

	if err != nil {
		panic(err)
	}
	//grpc tls
	var opts []grpc.ServerOption
	cre, err := credentials.NewServerTLSFromFile(c.TlsPem, c.TlsKey)
	if err != nil {
		panic(err)
	}
	opts = append(opts, grpc.Creds(cre))
	//tls end
	s := grpc.NewServer(opts...)
	//register rpc server handler
	//RegisterSchedulerServer(s, &Impl{})
	for _, f := range servers {
		f(s)
	}

	logger.Info("server listening at %v", lis.Addr())
	err = s.Serve(lis)
	return err
}

func Run(c ServerConfig, servers ...GrpcServerHandler) error {
	lis, err := net.Listen("tcp", ":"+fmt.Sprint(c.Port))

	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	for _, f := range servers {
		f(s)
	}
	logger.Info("server listening at %v", lis.Addr())
	err = s.Serve(lis)
	return err
}
