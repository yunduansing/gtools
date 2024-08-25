package middleware

import (
	"bytes"
	context2 "context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/yunduansing/gtools/context"
	"github.com/yunduansing/gtools/examples/gintest/apiContext"
	"github.com/yunduansing/gtools/examples/gintest/config"
	"github.com/yunduansing/gtools/examples/gintest/model"
	"github.com/yunduansing/gtools/logger"
	"github.com/yunduansing/gtools/utils"
	"io"
	"net/http"
	"time"
)

type GinAction func(c *apiContext.ApiContext) model.Response

func RequestLimiter(c *gin.Context) {
	if !config.IsLimiterOpen {
		c.Next()
		return
	}
	res, err := config.Limiter.Allow(c.Request.Context(), c.Request.RequestURI, redis_rate.PerSecond(config.LimitPerSecond))
	if err != nil {
		logger.Error(c.Request.Context(), "request limit err:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		c.Abort()
		return
	}
	if res.Allowed == 0 {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": "too many requests"})
		c.Abort()
		return
	}
	c.Next()
}

func WrapRequestMiddle(handler GinAction) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		var requestId string
		ctx := context2.WithValue(c.Request.Context(), "requestId", requestId)
		xRequestId := c.GetHeader("X-Request-Id")
		if len(xRequestId) > 0 {
			requestId = xRequestId
		} else {
			requestId = context.GenRequestIdByUUID()
		}

		c.Set("requestId", requestId)
		myCtx := context.NewContext(ctx, context.WithRequestId(requestId),
			context.WithRequestTime(start.Format(time.DateTime)))
		cc := apiContext.ApiContext{
			Ctx:        &myCtx,
			GinContext: c,
		}

		reqData, _ := io.ReadAll(c.Request.Body)
		myCtx.Log.Infof(ctx, "url=%s  req=%s", c.Request.URL.RequestURI(), utils.ByteToString(reqData))
		c.Request.Body = io.NopCloser(bytes.NewBuffer(reqData))

		c.Next()

		resp := handler(&cc)
		logger.Infof(ctx, "Cost=%s URL=%s Method=%s ClientIP=%s ResponseStatus=%d resp=%s", time.Since(start),
			c.Request.URL, c.Request.Method, c.ClientIP(), c.Writer.Status(), utils.ToJsonString(resp))
	}

}
