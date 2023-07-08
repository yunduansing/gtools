package logger

import "sync"

type Config struct {
	Level       string `json:",default=info,options=debug|info|warn|error|panic|fatal"` //日志级别，默认为info
	FilePath    string `json:",default=/log,optional"`                                  //日志文件路径
	LogType     string `json:",default=zap,options=logrus|zap,optional"`                //日志类型，目前支持zap和logrus
	ServiceName string `json:",optional"`                                               //所属服务
	MaxSize     int    `json:",default=10,optional"`                                    //日志文件最大数量
	MaxAge      int    `json:",default=30,optional"`                                    //最大保留天数
	BackupNum   int    `json:",default=100,optional"`                                   //最大保留日志文件数量
	Compress    bool   `json:",default=false,optional"`                                 //是否压缩
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

type KeyPair struct {
	Key string      `json:"key"`
	Val interface{} `json:"val"`
}

var (
	logger ILog
	once   sync.Once
)

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

func getLogger() ILog {
	if logger == nil {
		once.Do(func() {
			logger = newZapLog(Config{
				Level:       "info",
				FilePath:    "log",
				LogType:     "zap",
				ServiceName: "",
			})
		})
	}
	return logger
}

func Info(v ...interface{}) {
	getLogger().info(v...)
}
func Infof(format string, v ...interface{}) {
	getLogger().infof(format, v)
}
func Error(v ...interface{}) {
	getLogger().error(v...)
}
func Errorf(format string, v ...interface{}) {
	getLogger().errorf(format, v)
}
func Panic(v ...interface{}) {
	getLogger().panic(v...)
}
func Panicf(format string, v ...interface{}) {
	getLogger().panicf(format, v)
}
func Warn(v ...interface{}) {
	getLogger().warn(v...)
}
func Warnf(format string, v ...interface{}) {
	getLogger().warnf(format, v)
}
func Debug(v ...interface{}) {
	getLogger().debug(v...)
}
func Debugf(format string, v ...interface{}) {
	getLogger().debugf(format, v)
}

func Close(f func()) {
	f()
	getLogger().close()
}
