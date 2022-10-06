package httputils

import (
	"github.com/valyala/fasthttp"
	"io"
	"net/http"
	"time"
)

var fastClient *fasthttp.Client

func FastSend(target string, headers map[string]string, origin *http.Request) (data []byte, statusCode int, contentType string, err error) {
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

	if origin != nil {
		for k, v := range origin.Header {
			for i, _ := range v {
				req.Header.Add(k, v[i])
			}
		}
	}
	req.SetRequestURI(target)
	body, err := io.ReadAll(origin.Body)
	if err != nil {
		return nil, 0, "Application/json", err
	}
	req.SetBody(body)
	req.Header.SetMethod(origin.Method)
	c := &fasthttp.Client{
		MaxConnsPerHost: 10000,
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    30 * time.Second,
	}
	err = c.Do(req, resp)
	if err != nil {

		return nil, 500, "", err
	}
	return resp.Body(), resp.StatusCode(), "Application/json", err
}
