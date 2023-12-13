package logger

import (
	"testing"
)

func TestError(t *testing.T) {
	InitLog(Config{
		LogType:     "logrus",
		ServiceName: "gtools-logrus-log-test",
		Level:       "Debug",
		FilePath:    "./log/",
	})
	req := Req{1, "234234234uljfljrlerj"}
	//Info("gtools测试info级别日志", zap.String("method", "logger.TestInfo"), zap.Any("req", req))
	//Debug("gtools测试debug级别日志", req)
	Error("gtools测试error级别日志", req, testReqHttp())
}
