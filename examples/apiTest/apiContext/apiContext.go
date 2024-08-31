package apiContext

import (
	"github.com/gin-gonic/gin"
	"github.com/yunduansing/gtools/context"
)

type ApiContext struct {
	Ctx        *context.Context
	GinContext *gin.Context
}
