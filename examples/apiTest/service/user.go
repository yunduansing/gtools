package service

import (
	"apiTest/model"
	"github.com/yunduansing/gtools/context"
	"protocol/client"
	userpb "protocol/user"
)

type UserService struct {
	ctx *context.Context
}

func NewUserService(ctx *context.Context) *UserService {
	return &UserService{ctx: ctx}
}

func (u *UserService) GetUser(id int64) (res *model.User, code int64) {
	u.ctx.Log.Info(u.ctx.Ctx, "get user with id ", id)
	//res, err := model.FindUserById(u.ctx, id)
	//if err != nil {
	//	u.ctx.Log.Error(u.ctx.Ctx, "get user with id err：", err)
	//}
	getUser, code, err := client.GetUser(u.ctx, &userpb.GetUserReq{})
	if err != nil {
		return nil, -1
	}
	return &model.User{UserId: getUser.UserId, UserName: getUser.UserName, Phone: getUser.Phone}, 0
}
