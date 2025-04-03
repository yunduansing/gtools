package logger

import (
	"context"
	"sync"
)

var (
	logger ILog
	once   sync.Once
	p      sync.Pool
)

type Logger struct {
	logger ILog
	ctx    context.Context
	kv     []KeyPair
}

const (
	LogTypeZap    = "zap"
	LogTypeLogrus = "logrus"
)

func init() {
	p.New = func() any {
		return &Logger{logger: logger}
	}
}

// 初始化日志
func InitLog(c Config) {
	if c.FilePath == "" {
		c.FilePath = "./logs"
	}
	if c.LogType == LogTypeLogrus {
		logger = newLogrusLog(c)
	} else {
		logger = newZapLog(c)
	}
}

// 绑定 Context
func (l *Logger) WithContext(ctx context.Context) *Logger {
	l.ctx = ctx
	return l
}

// 追加 Key-Value 字段
func (l *Logger) WithField(key string, val any) *Logger {
	l.kv = append(l.kv, KeyPair{Key: key, Val: val})
	return l
}

// 获取 Logger 实例
func GetLogger() *Logger {
	if logger == nil {
		once.Do(
			func() {
				logger = newZapLog(
					Config{
						Level:       "Info",
						FilePath:    "./logs",
						LogType:     "zap",
						ServiceName: "",
					},
				)
			},
		)
	}

	logInstance := p.Get().(*Logger)
	logInstance.logger = logger
	return logInstance
}

// 统一处理日志方法，减少重复代码
func (l *Logger) logWithContext(
	ctx context.Context, logFunc func(ctx context.Context, v ...interface{}), v ...interface{},
) {
	if len(l.kv) > 0 {
		ctx = context.WithValue(ctx, "kv", l.kv)
	}
	defer func() {
		l.kv = nil
		p.Put(l)
	}()
	logFunc(ctx, v...)
}

// 统一处理日志格式化方法
func (l *Logger) logfWithContext(
	ctx context.Context, logFunc func(ctx context.Context, format string, v ...interface{}), format string,
	v ...interface{},
) {
	if len(l.kv) > 0 {
		ctx = context.WithValue(ctx, "kv", l.kv)
	}
	defer func() {
		l.kv = nil
		p.Put(l)
	}()
	logFunc(ctx, format, v...)
}

// 具体日志方法调用
func (l *Logger) Info(ctx context.Context, v ...interface{}) {
	l.logWithContext(ctx, l.logger.Info, v...)
}
func (l *Logger) Infof(ctx context.Context, format string, v ...interface{}) {
	l.logfWithContext(ctx, l.logger.Infof, format, v...)
}
func (l *Logger) Error(ctx context.Context, v ...interface{}) {
	l.logWithContext(ctx, l.logger.Error, v...)
}
func (l *Logger) Errorf(ctx context.Context, format string, v ...interface{}) {
	l.logfWithContext(ctx, l.logger.Errorf, format, v...)
}
func (l *Logger) Panic(ctx context.Context, v ...interface{}) {
	l.logWithContext(ctx, l.logger.Panic, v...)
}
func (l *Logger) Panicf(ctx context.Context, format string, v ...interface{}) {
	l.logfWithContext(ctx, l.logger.Panicf, format, v...)
}
func (l *Logger) Warn(ctx context.Context, v ...interface{}) {
	l.logWithContext(ctx, l.logger.Warn, v...)
}
func (l *Logger) Warnf(ctx context.Context, format string, v ...interface{}) {
	l.logfWithContext(ctx, l.logger.Warnf, format, v...)
}
func (l *Logger) Debug(ctx context.Context, v ...interface{}) {
	l.logWithContext(ctx, l.logger.Debug, v...)
}
func (l *Logger) Debugf(ctx context.Context, format string, v ...interface{}) {
	l.logfWithContext(ctx, l.logger.Debugf, format, v...)
}

// 关闭 Logger
func (l *Logger) Close(f func()) {
	f()
}
