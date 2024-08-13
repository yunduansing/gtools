package httptool

import (
	"context"
	"errors"
	"github.com/sony/gobreaker"
	"github.com/yunduansing/gtools/breaker"
	"github.com/yunduansing/gtools/logger"
	"net/http"
	"time"
)

const DefaultTimeout = time.Second * 30

type HttpTool struct {
	Retry            int           //重试次数，0表示不重试
	Interval         time.Duration //重试间隔，默认为0即不等待立即重试
	*breaker.Breaker               //熔断器，若为空则表示不熔断，默认不启用，需要单独调用set方法启用
	Timeout          time.Duration //超时设置，默认30s
	req              IRequest
}

func NewHttpTool(retry int, retryWait, timeout time.Duration, req IRequest) *HttpTool {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &HttpTool{Retry: retry, Interval: retryWait, Timeout: timeout, req: req}
}

func (h *HttpTool) SetBreaker(b *breaker.Breaker) {
	h.Breaker = b
}

func (h *HttpTool) Get(ctx context.Context, url string, reqData map[string]interface{}, headers map[string]string) (res []byte, statusCode int, err error) {
	for i := 0; i < h.Retry; i++ {
		if h.Breaker != nil {
			_, err = h.Breaker.Execute(func() (interface{}, error) {
				res, statusCode, err = h.req.Get(ctx, url, reqData, headers)
				return nil, err
			})
			if errors.Is(err, gobreaker.ErrOpenState) {
				return nil, http.StatusInternalServerError, err
			} else if errors.Is(err, gobreaker.ErrTooManyRequests) {
				return nil, http.StatusInternalServerError, err
			}
		} else {
			res, statusCode, err = h.req.Get(ctx, url, reqData, headers)
		}

		if err == nil {
			return
		}
		if h.Retry > 0 && h.Interval > 0 && i != h.Retry-1 {
			<-time.After(h.Interval)
			logger.Infof(ctx, "url=%s http.request.method=%s retry=%d retryWait=%d", url, "GET", i+1, h.Interval)
		}
	}
	return

}

func (h *HttpTool) PostJson(ctx context.Context, url string, reqData interface{}, headers map[string]string) (res []byte, statusCode int, err error) {
	for i := 0; i < h.Retry; i++ {
		if h.Breaker != nil {
			_, err = h.Breaker.Execute(func() (interface{}, error) {
				res, statusCode, err = h.req.PostJson(ctx, url, reqData, headers)
				return nil, err
			})
			if errors.Is(err, gobreaker.ErrOpenState) {
				return nil, http.StatusInternalServerError, err
			} else if errors.Is(err, gobreaker.ErrTooManyRequests) {
				return nil, http.StatusInternalServerError, err
			}
		} else {
			res, statusCode, err = h.req.PostJson(ctx, url, reqData, headers)
		}

		if err == nil {
			return
		}
		if h.Retry > 0 && h.Interval > 0 && i != h.Retry-1 {

			<-time.After(h.Interval)
			logger.Infof(ctx, "url=%s http.request.method=%s retry=%d retryWait=%d", url, "POST", i+1, h.Interval)
		}
	}
	return
}
