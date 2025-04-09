package model

type Response struct {
	Code      int64  `json:"code"`
	Msg       string `json:"msg"`
	Data      any    `json:"data"`
	RequestId string `json:"requestId"`
}
