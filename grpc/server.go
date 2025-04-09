package grpctool

import (
	"context"
	"fmt"
	"github.com/yunduansing/gtools/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"net"
	"time"
)

type ServerConfig struct {
	Port          int
	TlsPem        string
	TlsKey        string
	IsTracingOpen bool
}

type GrpcServerHandler func(server *grpc.Server)

func RunWithTls(c ServerConfig, serverRegister GrpcServerHandler, opts ...grpc.ServerOption) error {
	lis, err := net.Listen("tcp", ":"+fmt.Sprint(c.Port))

	if err != nil {
		panic(err)
	}
	//protocol tls

	cre, err := credentials.NewServerTLSFromFile(c.TlsPem, c.TlsKey)
	if err != nil {
		panic(err)
	}
	opts = append(opts, grpc.Creds(cre))
	opts = append(opts, grpc.KeepaliveEnforcementPolicy(kaep))
	opts = append(opts, grpc.KeepaliveParams(kasp))
	//tls end
	s := grpc.NewServer(opts...)
	//register rpc server handler
	//RegisterSchedulerServer(s, &Impl{})
	serverRegister(s)

	logger.GetLogger().Infof(context.TODO(), "rpc server listening at %+v", lis.Addr())
	err = s.Serve(lis)
	return err
}

var kaep = keepalive.EnforcementPolicy{
	MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
	PermitWithoutStream: true,            // Allow pings even when there are no active streams
}

var kasp = keepalive.ServerParameters{
	MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
	MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
	MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
	Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
	Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
}

func Run(c ServerConfig, serverRegister GrpcServerHandler, opts ...grpc.ServerOption) error {
	lis, err := net.Listen("tcp", ":"+fmt.Sprint(c.Port))

	if err != nil {
		panic(err)
	}
	opts = append(opts, grpc.KeepaliveEnforcementPolicy(kaep))
	opts = append(opts, grpc.KeepaliveParams(kasp))
	s := grpc.NewServer(opts...)
	serverRegister(s)
	logger.GetLogger().Infof(context.TODO(), "rpc server listening at %+v", lis.Addr())
	err = s.Serve(lis)
	return err
}
