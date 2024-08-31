package controllers

import (
	"github.com/yunduansing/gtools/examples/apiTest/apiContext"
	"github.com/yunduansing/gtools/examples/apiTest/model"
	"github.com/yunduansing/gtools/examples/apiTest/service"
)

func UserLogin(c *apiContext.ApiContext) (resp model.Response) {

	var req service.UserLoginReq
	err := c.GinContext.Bind(&req)
	if err != nil {
		c.Ctx.Log.Error(c.Ctx.Ctx, "parse req param error：", err)
		resp.Code = -1
		resp.Msg = err.Error()
		return resp
	}
	if len(req.Phone) == 0 {
		resp.Code = -1
		resp.Msg = "手机号不能为空"
		return resp
	}
	loginData, err := service.NewUserLoginService(c.Ctx).UserLogin(&service.UserLoginReq{
		Phone: req.Phone,
		Code:  "",
	})
	if err != nil {
		model.GetErrorResponse(-1, err, &resp)
		return resp
	}
	resp.Data = loginData
	return resp
}
