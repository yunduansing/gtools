package logger

type ILog interface {
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	close()
	WithField(field, value string)
}
