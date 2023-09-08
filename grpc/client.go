package grpctool

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type ClientConfig struct {
	Address   string
	Port      int
	ServerPem string
}

func Init(c ClientConfig) (*grpc.ClientConn, error) {
	target := c.Address
	if c.Port > 0 {
		target += ":" + fmt.Sprint(c.Port)
	}
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return conn, err
}

func InitWithTls(c ClientConfig) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	//no tls grpc.WithTransportCredentials(insecure.NewCredentials())
	cre, err := credentials.NewClientTLSFromFile(c.ServerPem, "go-grpc-example")
	if err != nil {
		log.Fatalf("Failed to create TLS credentials: %v", err)
	}

	opts = append(opts, grpc.WithTransportCredentials(cre))
	target := c.Address
	if c.Port > 0 {
		target += ":" + fmt.Sprint(c.Port)
	}
	conn, err := grpc.Dial(target, opts...)
	return conn, err
}
