package httputils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/yunduansing/gtools/logger"
	"net/http"
	"strings"
	"time"
)

var fastClient = &fasthttp.Client{
	MaxConnWaitTimeout: 1 * time.Minute,
}

func doRequest(ctx context.Context, url string, method, contentType string, reqData interface{}, headers map[string]string) (res []byte, statusCode int, err error) {
	req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()
	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}
	req.Header.Set("Content-Type", contentType)
	req.SetRequestURI(url)

	if reqData != nil {
		data, _ := json.Marshal(reqData)
		req.SetBody(data)
	}

	req.Header.SetMethod(method)
	start := time.Now()
	err = fastClient.Do(req, resp)
	end := time.Since(start)
	if resp != nil {
		logger.Errorf("url=%s http.request.method=%s http.response.statusCode=%d http.request.time=%s", url, method, resp.StatusCode(), end.String())
	} else {
		logger.Infof("url=%s http.request.method=%s http.request.time=%s", url, method, end.String())
	}

	if err != nil {
		return nil, 500, err
	}
	if resp == nil {
		logger.Infof("url=%s http.request.method=%s err=resp对象为空", url, method)
		return nil, 500, errors.New("http请求错误")
	}
	respBody := resp.Body()

	resp.ReleaseBody(len(respBody))

	return respBody, resp.StatusCode(), err
}

func PostJson(ctx context.Context, url string, reqData interface{}, headers map[string]string) (res []byte, statusCode int, err error) {
	return doRequest(ctx, url, http.MethodPost, "application/json", reqData, headers)
}

func Get(ctx context.Context, url string, reqData map[string]interface{}, headers map[string]string) (res []byte, statusCode int, err error) {
	if reqData != nil && len(reqData) > 0 {
		if !strings.Contains(url, "?") {
			url += "?"
		}
		for k, v := range reqData {
			url += k + "=" + fmt.Sprint(v) + "&"
		}
		url = strings.TrimRight(url, "&")
	}
	return doRequest(ctx, url, http.MethodGet, "application/json", nil, headers)
}
