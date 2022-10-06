package logger

import (
	"go.uber.org/zap"
	"testing"
)

type Req struct {
	Id int `json:"id"`
}

func TestInfo(t *testing.T) {
	InitLogger()
	req := Req{1}
	Info("xxxx", zap.String("s", "s"), zap.Any("req", req))
}
