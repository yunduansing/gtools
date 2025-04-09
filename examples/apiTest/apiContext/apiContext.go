package apiContext

import (
	"apiTest/service"
	"github.com/gin-gonic/gin"
	"github.com/yunduansing/gtools/context"
)

type ApiContext struct {
	Ctx        *context.Context
	GinContext *gin.Context
	UserInfo   *service.UserLoginDataCache
}
