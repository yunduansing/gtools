package apiContext

import (
	"github.com/gin-gonic/gin"
	"github.com/yunduansing/gtools/context"
	"github.com/yunduansing/gtools/examples/apiTest/service"
)

type ApiContext struct {
	Ctx        *context.Context
	GinContext *gin.Context
	UserInfo   *service.UserLoginDataCache
}
