package client

import (
	context2 "context"
	"errors"
	"github.com/yunduansing/gtools/context"
	grpctool "github.com/yunduansing/gtools/grpc"
	"github.com/yunduansing/gtools/logger"
	"google.golang.org/grpc"
	"protocol/middleware"
	userpb "protocol/user"
	"sync"
)

var (
	userClient *grpc.ClientConn
	userOnce   sync.Once
)

func NewUserClient() userpb.UserServiceClient {
	userOnce.Do(
		func() {
			conn, err := grpctool.Init(
				grpctool.ClientConfig{
					Address:   "localhost",
					Port:      8080,
					ServerPem: "",
				}, grpc.WithUnaryInterceptor(middleware.UnaryReqTimeInterceptor),
			)
			if err != nil {
				logger.GetLogger().Panic(context2.Background(), "create user protocol client conn error:", err)
				panic(err)
			}
			userClient = conn
		},
	)
	return userpb.NewUserServiceClient(userClient)
}

func GetUser(ctx context.Context, req *userpb.GetUserReq) (*userpb.User, int64, error) {
	c := NewUserClient()
	resp, err := c.GetUser(ctx.Ctx, req)
	if err != nil {
		logger.GetLogger().Panic(ctx., "get user error:", err)
		return nil, -1, err
	}
	if resp.Code != 0 {
		return nil, resp.Code, errors.New(resp.Msg)
	}
	return resp.GetUserData(), 0, nil
}
