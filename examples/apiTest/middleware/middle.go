package middleware

import (
	"bytes"
	context2 "context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/yunduansing/gtools/context"
	"github.com/yunduansing/gtools/examples/apiTest/apiContext"
	"github.com/yunduansing/gtools/examples/apiTest/config"
	"github.com/yunduansing/gtools/examples/apiTest/model"
	"github.com/yunduansing/gtools/logger"
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
		logger.GetLogger().Error(c.Request.Context(), "request limit err:", err)
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

		xRequestId := c.GetHeader("X-Request-Id")
		if len(xRequestId) > 0 {
			requestId = xRequestId
		} else {
			requestId = context.GenRequestIdByUUID()
		}
		var headers = make(map[string]string)
		for k, header := range c.Request.Header {
			headers[k] = header[0]
		}

		c.Set("requestId", requestId)
		ctx := context2.WithValue(c.Request.Context(), "requestId", requestId)
		myCtx := context.NewContext(ctx, context.WithRequestId(requestId),
			context.WithRequestTime(start.Format(time.DateTime)))
		cc := apiContext.ApiContext{
			Ctx:        &myCtx,
			GinContext: c,
		}

		reqData, _ := io.ReadAll(c.Request.Body)
		var req any
		_ = json.Unmarshal(reqData, &req)
		myCtx.Log.WithField("URL", c.Request.URL).
			WithField("Method", c.Request.Method).
			WithField("ClientIP", c.ClientIP()).
			WithField("Headers", headers).
			WithField("Req", req).
			Info(ctx, "api请求")
		c.Request.Body = io.NopCloser(bytes.NewBuffer(reqData))

		//c.Next()

		resp := handler(&cc)
		resp.RequestId = requestId
		myCtx.Log.WithField("URL", c.Request.URL).
			WithField("Cost", time.Since(start)).
			WithField("Method", c.Request.Method).
			WithField("Resp", resp).
			WithField("ClientIP", c.ClientIP()).
			WithField("ResponseStatus", c.Writer.Status()).
			Info(ctx, "api响应")
		c.JSON(200, resp)
	}

}
