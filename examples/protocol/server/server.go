package server

import (
	"context"
	userpb "protocol/user"
)

type Server struct {
	userpb.UnimplementedUserServiceServer
}

func (s *Server) mustEmbedUnimplementedUserServiceServer() {
	//TODO implement me
	panic("implement me")
}

func (s *Server) GetUser(ctx context.Context, req *userpb.GetUserReq) (*userpb.UserResponse, error) {
	resp := &userpb.UserResponse{}
	resp.Data = &userpb.UserResponse_UserData{
		UserData: &userpb.User{
			UserId:   1,
			UserName: "lucky",
			Phone:    "18888888888",
		},
	}
	return resp, nil
}
