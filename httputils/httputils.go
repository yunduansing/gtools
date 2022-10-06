package httputils

import (
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

const timeout = 30 * time.Second

// HttpProxy http代理
func HttpProxy(serviceUrl string, w http.ResponseWriter, r *http.Request) {
	targetUrl, _ := url.Parse(serviceUrl)
	proxy := httputil.NewSingleHostReverseProxy(targetUrl)
	proxy.ServeHTTP(w, r)
}

// HttpGet 发送http get请求
func HttpGet(path string, headers map[string]string) (data []byte, statusCode int, err error) {
	client := &http.Client{Timeout: timeout}
	var req = new(http.Request)
	req.Method = http.MethodGet
	targetUrl, _ := url.Parse(path)
	req.URL = targetUrl
	if req.Header == nil {
		req.Header = make(map[string][]string)
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 500, err
	}
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	return data, resp.StatusCode, err
}

// HttpPost 发送http post请求
func HttpPost(url string, headers map[string]string, body []byte) (data []byte, statusCode int, err error) {
	return HttpRequest(http.MethodPost, url, headers, body)
}

// HttpPut 发送http put请求
func HttpPut(url string, headers map[string]string, body []byte) (data []byte, statusCode int, err error) {
	return HttpRequest(http.MethodPut, url, headers, body)
}

// HttpDelete 发送http delete请求
func HttpDelete(url string, headers map[string]string, body []byte) (data []byte, statusCode int, err error) {
	return HttpRequest(http.MethodDelete, url, headers, body)
}

// HttpRequest 发送http请求
func HttpRequest(method, path string, headers map[string]string, body []byte) (data []byte, statusCode int, err error) {
	client := &http.Client{Timeout: 30 * time.Second}
	var req = new(http.Request)
	req.Method = method
	targetUrl, _ := url.Parse(path)
	req.URL = targetUrl
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	if body != nil {
		req.Body.Read(body)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 500, err
	}
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	return data, resp.StatusCode, err
}

// TransmitRequest 转发http请求
func TransmitRequest(path string, origin *http.Request) (data []byte, statusCode int, contentType string, err error) {
	client := &http.Client{Timeout: 30 * time.Second}
	var req = new(http.Request)
	req.Method = origin.Method
	targetUrl, _ := url.Parse(path)
	req.URL = targetUrl
	req.Header = origin.Header

	req.Body = origin.Body
	resp, err := client.Do(req)
	if err != nil {
		return nil, resp.StatusCode, resp.Header.Get("Content-Type"), err
	}
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	return data, resp.StatusCode, resp.Header.Get("Content-Type"), err
}

func SendRequest(req *http.Request) (data []byte, statusCode int, err error) {
	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	return data, resp.StatusCode, err
}
