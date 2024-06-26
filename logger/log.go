package logger

import (
	"context"
	"sync"
)

var (
	logger ILog
	once   sync.Once
)

const (
	LogTypeZap    = "zap"
	LogTypeLogrus = "logrus"
)

func InitLog(c Config) {
	switch c.LogType {
	case LogTypeLogrus:
		logger = newLogrusLog(c)
	//case "zap":
	//logger = newZapLog(c)
	default:
		logger = newZapLog(c)
	}
}

func GetLogger() ILog {
	if logger == nil {
		once.Do(func() {
			logger = newZapLog(Config{
				Level:       "Info",
				FilePath:    "log",
				LogType:     "zap",
				ServiceName: "",
			})
		})
	}
	return logger
}

func Info(ctx context.Context, v ...interface{}) {
	GetLogger().Info(ctx, v...)
}
func Infof(ctx context.Context, format string, v ...interface{}) {
	GetLogger().Infof(ctx, format, v...)
}
func Error(ctx context.Context, v ...interface{}) {
	GetLogger().Error(ctx, v...)
}
func Errorf(ctx context.Context, format string, v ...interface{}) {
	GetLogger().Errorf(ctx, format, v)
}
func Panic(ctx context.Context, v ...interface{}) {
	GetLogger().Panic(ctx, v...)
}
func Panicf(ctx context.Context, format string, v ...interface{}) {
	GetLogger().Panicf(ctx, format, v)
}
func Warn(ctx context.Context, v ...interface{}) {
	GetLogger().Warn(ctx, v...)
}
func Warnf(ctx context.Context, format string, v ...interface{}) {
	GetLogger().Warnf(ctx, format, v)
}
func Debug(ctx context.Context, v ...interface{}) {
	GetLogger().Debug(ctx, v...)
}
func Debugf(ctx context.Context, format string, v ...interface{}) {
	GetLogger().Debugf(ctx, format, v)
}

func Close(f func()) {
	f()
	GetLogger().close()
}
