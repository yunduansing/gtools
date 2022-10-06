package logger

type Config struct {
	Level    string `json:",default=info,options=debug|info|warn|error|panic|fatal"`
	FileName string `json:",optional"`
	LogType  string `json:",default=logrus,options=logrus|zap,optional"`
}

type LogLevel int

const (
	LevelFatal LogLevel = 1
	LevelPanic LogLevel = 2
	LevelError LogLevel = 3
	LevelWarn  LogLevel = 4
	LevelInfo  LogLevel = 5
	LevelDebug LogLevel = 6
)

var logger ILog

func InitLog(c Config) {
	switch c.LogType {
	case "logrus":
		logger = newLogrusLog(c)
	case "zap":
		logger = newZapLog(c)
	default:
		logger = newLogrusLog(c)
	}
}

func Info(v ...interface{}) {
	logger.Info(v)
}
func Infof(format string, v ...interface{}) {
	logger.Infof(format, v)
}
func Error(v ...interface{}) {
	logger.Error(v)
}
func Errorf(format string, v ...interface{}) {
	logger.Errorf(format, v)
}
func Panic(v ...interface{}) {
	logger.Panic(v)
}
func Panicf(format string, v ...interface{}) {
	logger.Panicf(format, v)
}
func Warn(v ...interface{}) {
	logger.Warn(v)
}
func Warnf(format string, v ...interface{}) {
	logger.Warnf(format, v)
}
func Debug(v ...interface{}) {
	logger.Debug(v)
}
func Debugf(format string, v ...interface{}) {
	logger.Debugf(format, v)
}
