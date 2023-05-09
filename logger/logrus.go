package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
	"time"
)

type logrusLog struct {
	*logrus.Logger
	fErr  *os.File
	fInfo *os.File
}

func (l *logrusLog) close() {
	if l.fErr != nil {
		l.fErr.Close()
	}
	if l.fInfo != nil {
		l.fInfo.Close()
	}
}

func getLogrusMsg(v ...interface{}) string {
	//list := v.([]interface{})
	var msg strings.Builder
	for _, item := range v {
		contentItem := getLogContent(item)
		msg.WriteString(contentItem + " ")
	}
	return msg.String()
}

func (l *logrusLog) info(v ...interface{}) {
	logMsg := getLogrusMsg(v...)
	l.Info(logMsg)
}

func (l *logrusLog) infof(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *logrusLog) error(v ...interface{}) {
	logMsg := getLogrusMsg(v...)
	l.Error(logMsg)
}

func (l *logrusLog) errorf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *logrusLog) panic(v ...interface{}) {
	logMsg := getLogrusMsg(v...)
	l.Panic(logMsg)
}

func (l *logrusLog) panicf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *logrusLog) warn(v ...interface{}) {
	logMsg := getLogrusMsg(v...)
	l.Warn(logMsg)
}

func (l *logrusLog) warnf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *logrusLog) debug(v ...interface{}) {
	logMsg := getLogrusMsg(v...)
	l.Debug(logMsg)
}

func (l *logrusLog) debugf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

type DefaultFieldHook struct {
	c Config
}

func (hook *DefaultFieldHook) Fire(entry *logrus.Entry) error {
	if len(hook.c.ServiceName) > 0 {
		entry.Data["service"] = hook.c.ServiceName
	}

	return nil
}

func (hook *DefaultFieldHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func newLogrusLog(c Config) *logrusLog {
	log := logrus.New()
	log.SetLevel(getLogrusLevel(c.Level))
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetReportCaller(true)
	log.AddHook(&DefaultFieldHook{c})
	//f, err := os.OpenFile(fmt.Sprintf("%s/info.log", c.FilePath), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	//if err != nil {
	//	fmt.Println("Failed to create logfile" + "info.log")
	//	panic(err)
	//}
	//fErr, err := os.OpenFile(fmt.Sprintf("%s/error.log", c.FilePath), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	//if err != nil {
	//	fmt.Println("Failed to create logfile" + "error.log")
	//	panic(err)
	//}
	//log.Out = io.MultiWriter(f, fErr, os.Stdout)
	//log.Formatter = &easy.Formatter{
	//	TimestampFormat: "2006-01-02 15:04:05",
	//	LogFormat:       "[%lvl%]: %time% - %msg%\n",
	//}
	logPath := c.FilePath
	if len(logPath) == 0 {
		logPath = "/log"
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: logWriter(logPath, "info", c.MaxAge, c.BackupNum), // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  logWriter(logPath, "info", c.MaxAge, c.BackupNum),
		logrus.WarnLevel:  logWriter(logPath, "info", c.MaxAge, c.BackupNum),
		logrus.ErrorLevel: logWriter(logPath, "error", c.MaxAge, c.BackupNum),
		logrus.FatalLevel: logWriter(logPath, "error", c.MaxAge, c.BackupNum),
		logrus.PanicLevel: logWriter(logPath, "error", c.MaxAge, c.BackupNum),
	}, &logrus.JSONFormatter{})
	log.AddHook(lfHook)
	return &logrusLog{Logger: log}
}

func logWriter(logPath string, level string, maxAge, backupNum int) *rotatelogs.RotateLogs {
	logFullPath := path.Join(logPath, level)
	logwriter, err := rotatelogs.New(
		logFullPath+".%Y%m%d",
		rotatelogs.WithLinkName(logFullPath),                      // 生成软链，指向最新日志文件
		rotatelogs.WithRotationCount(uint(backupNum)),             // 文件最大保存份数
		rotatelogs.WithRotationTime(24*time.Hour),                 // 日志切割时间间隔
		rotatelogs.WithMaxAge(time.Duration(maxAge)*time.Hour*24), //保留天数
	)
	if err != nil {
		panic(err)
	}
	return logwriter
}

func getLogrusLevel(level string) logrus.Level {
	switch level {
	case "info":
		return logrus.InfoLevel
	case "debug":
		return logrus.DebugLevel
	case "error":
		return logrus.ErrorLevel
	case "panic":
		return logrus.PanicLevel
	case "fatal":
		return logrus.FatalLevel
	case "warn":
		return logrus.WarnLevel
	}
	return logrus.InfoLevel
}
