package logger

import "sync"

var (
	logger ILog
	once   sync.Once
)

func InitLog(c Config) {
	switch c.LogType {
	case "logrus":
		logger = newLogrusLog(c)
	//case "zap":
	//logger = newZapLog(c)
	default:
		logger = newZapLog(c)
	}
}

func getLogger() ILog {
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

func Info(v ...interface{}) {
	getLogger().Info(v...)
}
func Infof(format string, v ...interface{}) {
	getLogger().Infof(format, v)
}
func Error(v ...interface{}) {
	getLogger().Error(v...)
}
func Errorf(format string, v ...interface{}) {
	getLogger().Errorf(format, v)
}
func Panic(v ...interface{}) {
	getLogger().Panic(v...)
}
func Panicf(format string, v ...interface{}) {
	getLogger().Panicf(format, v)
}
func Warn(v ...interface{}) {
	getLogger().Warn(v...)
}
func Warnf(format string, v ...interface{}) {
	getLogger().Warnf(format, v)
}
func Debug(v ...interface{}) {
	getLogger().Debug(v...)
}
func Debugf(format string, v ...interface{}) {
	getLogger().Debugf(format, v)
}

func Close(f func()) {
	f()
	getLogger().close()
}
