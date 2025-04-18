package controllers

import (
	"apiTest/apiContext"
	"apiTest/model"
	"apiTest/service"
	"errors"
	model2 "github.com/yunduansing/gtools/examples/coupon/model"
	"time"
)

func GetUser(c *apiContext.ApiContext) model.Response {
	c.Ctx.Log.WithField("controller", "User").Info(c.Ctx.Ctx, "get user")
	<-time.After(time.Millisecond * 300)
	userCtx := service.NewUserService(c.Ctx)
	res, code := userCtx.GetUser(1)
	return model.Response{
		Code:      code,
		Msg:       "ok",
		Data:      res,
		RequestId: c.Ctx.GetRequestId(),
	}
}

func AddOrUpdateUser(c *apiContext.ApiContext) model.Response {
	var resp = model.Response{}
	var req model2.User
	err := c.GinContext.Bind(&req)
	if err != nil {
		c.Ctx.Log.Error(c.Ctx.Ctx, "parse req param error：", err)
		resp.Code = -1
		resp.Msg = err.Error()
		return resp
	}
	c.Ctx.Log.Error(c.Ctx.Ctx, errors.New("1111"))
	c.Ctx.Log.WithField("controller", "User").Info(c.Ctx.Ctx, "AddOrUpdateUser")
	<-time.After(time.Millisecond * 300)
	//userCtx := service.NewUserService(c.Ctx)
	return model.Response{
		Code:      200,
		Msg:       "ok",
		RequestId: c.Ctx.GetRequestId(),
	}
}
