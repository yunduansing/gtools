package logger

type Config struct {
	Level       string `json:",default=info,options=debug|info|warn|error|panic|fatal"` //日志级别，默认为info
	FileName    string `json:",optional"`                                               //文件名
	LogType     string `json:",default=logrus,options=logrus|zap,optional"`             //日志类型，目前支持zap和logrus
	ServiceName string `json:",optional"`                                               //所属服务
	MaxSize     string `json:",default=10MB,optional"`                                  //日志文件最大数量
	MaxAge      int    `json:",default=30,optional"`                                    //最大保留天数
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
		logger = newZapLog(c)
	}
}

func Info(v ...interface{}) {
	logger.info(v)
}
func Infof(format string, v ...interface{}) {
	logger.infof(format, v)
}
func Error(v ...interface{}) {
	logger.error(v)
}
func Errorf(format string, v ...interface{}) {
	logger.errorf(format, v)
}
func Panic(v ...interface{}) {
	logger.panic(v)
}
func Panicf(format string, v ...interface{}) {
	logger.panicf(format, v)
}
func Warn(v ...interface{}) {
	logger.warn(v)
}
func Warnf(format string, v ...interface{}) {
	logger.warnf(format, v)
}
func Debug(v ...interface{}) {
	logger.debug(v)
}
func Debugf(format string, v ...interface{}) {
	logger.debugf(format, v)
}
