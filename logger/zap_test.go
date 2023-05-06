package logger

import (
	"go.uber.org/zap"
	"testing"
)

type Req struct {
	Id    int    `json:"id"`
	ReqId string `json:"req_id"`
}

func TestInfo(t *testing.T) {
	InitLog(Config{
		LogType:     "logrus",
		ServiceName: "gtools-test",
	})
	req := Req{1, "234234234uljfljrlerj"}
	Info("gtools测试info级别日志", zap.String("method", "logger.TestInfo"), zap.Any("req", req))
	Error("gtools测试error级别日志", zap.String("method", "logger.TestInfo"), zap.Any("req", req))
}
