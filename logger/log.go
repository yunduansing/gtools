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

func InitLog(c Config) {
	if len(c.FilePath) == 0 {
		c.FilePath = "./logs"
	}
	switch c.LogType {
	case LogTypeLogrus:
		logger = newLogrusLog(c)
	//case "zap":
	//logger = newZapLog(c)
	default:
		logger = newZapLog(c)
	}
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	l.ctx = ctx
	return l
}

func (l *Logger) WithField(key string, val any) *Logger {
	l.kv = append(l.kv, KeyPair{
		Key: key,
		Val: val,
	})
	return l
}

func GetLogger() *Logger {
	if logger == nil {
		once.Do(func() {
			logger = newZapLog(Config{
				Level:       "Info",
				FilePath:    "./logs",
				LogType:     "zap",
				ServiceName: "",
			})
		})
	}

	return p.Get().(*Logger)
}

func (l *Logger) Info(ctx context.Context, v ...interface{}) {
	if len(l.kv) > 0 {
		ctx = context.WithValue(ctx, "kv", l.kv)
	}
	defer func() {
		l.kv = nil
		p.Put(l)
	}()
	l.logger.Info(ctx, v...)
}
func (l *Logger) Infof(ctx context.Context, format string, v ...interface{}) {
	if len(l.kv) > 0 {
		ctx = context.WithValue(ctx, "kv", l.kv)
	}
	defer func() {
		l.kv = nil
		p.Put(l)
	}()
	l.logger.Infof(ctx, format, v...)
}
func (l *Logger) Error(ctx context.Context, v ...interface{}) {
	if len(l.kv) > 0 {
		ctx = context.WithValue(ctx, "kv", l.kv)
	}
	defer func() {
		l.kv = nil
		p.Put(l)
	}()
	l.logger.Error(ctx, v...)
}
func (l *Logger) Errorf(ctx context.Context, format string, v ...interface{}) {
	if len(l.kv) > 0 {
		ctx = context.WithValue(ctx, "kv", l.kv)
	}
	defer func() {
		l.kv = nil
		p.Put(l)
	}()
	l.logger.Errorf(ctx, format, v)
}
func (l *Logger) Panic(ctx context.Context, v ...interface{}) {
	if len(l.kv) > 0 {
		ctx = context.WithValue(ctx, "kv", l.kv)
	}
	defer p.Put(l)
	l.logger.Panic(ctx, v...)
}
func (l *Logger) Panicf(ctx context.Context, format string, v ...interface{}) {
	if len(l.kv) > 0 {
		ctx = context.WithValue(ctx, "kv", l.kv)
	}
	defer func() {
		l.kv = nil
		p.Put(l)
	}()
	l.logger.Panicf(ctx, format, v)
}
func (l *Logger) Warn(ctx context.Context, v ...interface{}) {
	if len(l.kv) > 0 {
		ctx = context.WithValue(ctx, "kv", l.kv)
	}
	defer func() {
		l.kv = nil
		p.Put(l)
	}()
	l.logger.Warn(ctx, v...)
}
func (l *Logger) Warnf(ctx context.Context, format string, v ...interface{}) {
	if len(l.kv) > 0 {
		ctx = context.WithValue(ctx, "kv", l.kv)
	}
	defer func() {
		l.kv = nil
		p.Put(l)
	}()
	l.logger.Warnf(ctx, format, v)
}
func (l *Logger) Debug(ctx context.Context, v ...interface{}) {
	if len(l.kv) > 0 {
		ctx = context.WithValue(ctx, "kv", l.kv)
	}
	defer func() {
		l.kv = nil
		p.Put(l)
	}()
	l.logger.Debug(ctx, v...)
}
func (l *Logger) Debugf(ctx context.Context, format string, v ...interface{}) {
	if len(l.kv) > 0 {
		ctx = context.WithValue(ctx, "kv", l.kv)
	}
	defer func() {
		l.kv = nil
		p.Put(l)
	}()
	l.logger.Debugf(ctx, format, v)
}

func (l *Logger) Close(f func()) {
	f()
	//GetLogger().close()
}
