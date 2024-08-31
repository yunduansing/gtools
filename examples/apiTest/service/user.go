package service

import (
	"github.com/yunduansing/gtools/context"
	"github.com/yunduansing/gtools/examples/apiTest/model"
)

type UserService struct {
	ctx *context.Context
}

func NewUserService(ctx *context.Context) *UserService {
	return &UserService{ctx: ctx}
}

func (u *UserService) GetUser(id int64) (res *model.User) {
	u.ctx.Log.Info(u.ctx.Ctx, "get user with id ", id)
	res, err := model.FindUserById(u.ctx, id)
	if err != nil {
		u.ctx.Log.Error(u.ctx.Ctx, "get user with id err：", err)
	}
	return
}
