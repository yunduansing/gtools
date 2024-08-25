package controllers

import (
	"github.com/yunduansing/gtools/examples/gintest/apiContext"
	"github.com/yunduansing/gtools/examples/gintest/model"
	"time"
)

func GetUser(c *apiContext.ApiContext) model.Response {
	c.Ctx.Log.Info(c.Ctx.Ctx, "get user")
	<-time.After(time.Millisecond * 300)

	return model.Response{
		Code: 200,
		Msg:  "ok",
		Data: struct {
			Id   int64  `json:"id"`
			Name string `json:"name"`
		}{
			Id:   1,
			Name: "张三",
		},
		RequestId: c.Ctx.GetRequestId(),
	}
}
