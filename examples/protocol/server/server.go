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
	//TODO implement me
	panic("implement me")
}
