package grpctool

import (
	"fmt"
	"github.com/yunduansing/gtools/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientConfig struct {
	Address   string
	Port      int
	ServerPem string
}

func Init(c ClientConfig, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	target := c.Address
	if c.Port > 0 {
		target += ":" + fmt.Sprint(c.Port)
	}
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(target, opts...)
	return conn, err
}

func WithClientNoCredentials() grpc.DialOption {
	return grpc.WithTransportCredentials(insecure.NewCredentials())
}

func WithClientTlsCredentials(cre credentials.TransportCredentials) grpc.DialOption {
	return grpc.WithTransportCredentials(cre)
}

func InitWithTls(c ClientConfig, opts ...grpc.DialOption) (*grpc.ClientConn, error) {

	//no tls grpc.WithTransportCredentials(insecure.NewCredentials())
	cre, err := credentials.NewClientTLSFromFile(c.ServerPem, "")
	if err != nil {
		logger.Panicf("Failed to create TLS credentials: %v", err)
	}

	opts = append(opts, WithClientTlsCredentials(cre))
	target := c.Address
	if c.Port > 0 {
		target += ":" + fmt.Sprint(c.Port)
	}
	conn, err := grpc.Dial(target, opts...)
	return conn, err
}
