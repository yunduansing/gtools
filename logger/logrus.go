package logger

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/yunduansing/gtools/utils"
	"os"
	"path"
	"runtime/debug"
	"strings"
	"time"
)

type logrusLog struct {
	*logrus.Logger
	fErr   *os.File
	fInfo  *os.File
	fields []KeyPair
}

func (l *logrusLog) close() {
	if l.fErr != nil {
		l.fErr.Close()
	}
	if l.fInfo != nil {
		l.fInfo.Close()
	}
}

func getLogrusMsg(l *logrusLog, v ...interface{}) string {
	//list := v.([]interface{})
	var msg strings.Builder

	var errList []error
	for _, item := range v {

		switch e := item.(type) {
		case error:
			errList = append(errList, e)
		default:
			contentItem := getLogContent(item)
			msg.WriteString(contentItem + " ")
		}
	}

	//for _, e := range errList {
	//	msg.WriteString(e.Error() + "\n")
	//	msg.WriteString("\n")
	//	stack := Debug.Stack()
	//	fmt.Fprintln(l.Writer(), string(stack))
	//}
	return msg.String()
}

func getLogrusErrorEntry(l *logrusLog, content interface{}) *logrus.Entry {
	switch e := content.(type) {
	case error:
		return l.WithError(e)
	}
	return nil
}

func getLogrusEntry(l *logrusLog, v ...interface{}) *logrus.Entry {
	for _, item := range v {
		res := getLogrusErrorEntry(l, item)
		if res != nil {
			return res
		}
	}
	return l.WithTime(time.Now())
}

// 获取错误堆栈信息
func getErrorStack(l *logrusLog, err error) string {
	var stack strings.Builder
	for {
		if err != nil {
			stack.WriteString(err.Error())
			stack.WriteString("\n")
		}
		debugStack := debug.Stack()
		if debugStack != nil {
			fmt.Fprint(l.Out, stack)
			//stack.WriteString(string(debugStack))
		}
		if err == nil {
			break
		}
		cause, ok := err.(interface{ Causes() []error })
		if !ok {
			break
		}
		if len(cause.Causes()) == 0 {
			break
		}
		err = cause.Causes()[0]
	}
	return stack.String()
}

func (log *logrusLog) WithField(field, value string) {
	log.fields = append(log.fields, KeyPair{field, value})
}

func (l *logrusLog) Info(v ...interface{}) {
	logMsg := getLogrusMsg(l, v...)
	var fields = make(map[string]interface{})
	for _, with := range l.fields {
		fields[with.Key] = with.Val
	}
	l.fields = nil
	l.Logger.WithFields(fields).Info(logMsg)
}

func (l *logrusLog) Infof(format string, v ...interface{}) {
	var fields = make(map[string]interface{})
	for _, with := range l.fields {
		fields[with.Key] = with.Val
	}
	l.fields = nil
	l.Logger.WithFields(fields).Infof(format, v...)
}

func (l *logrusLog) Error(v ...interface{}) {
	logMsg := getLogrusMsg(l, v...)
	var fields = make(map[string]interface{})
	for _, with := range l.fields {
		fields[with.Key] = with.Val
	}
	l.fields = nil
	l.Logger.WithFields(fields).Error(logMsg)
}

func (l *logrusLog) Errorf(format string, v ...interface{}) {
	var fields = make(map[string]interface{})
	for _, with := range l.fields {
		fields[with.Key] = with.Val
	}
	l.fields = nil
	l.Logger.WithFields(fields).Errorf(format, v...)
}

func (l *logrusLog) Panic(v ...interface{}) {
	logMsg := getLogrusMsg(l, v...)
	var fields = make(map[string]interface{})
	for _, with := range l.fields {
		fields[with.Key] = with.Val
	}
	l.fields = nil
	l.Logger.WithFields(fields).Panic(logMsg)
}

func (l *logrusLog) Panicf(format string, v ...interface{}) {
	var fields = make(map[string]interface{})
	for _, with := range l.fields {
		fields[with.Key] = with.Val
	}
	l.fields = nil
	l.Logger.WithFields(fields).Panicf(format, v...)
}

func (l *logrusLog) Warn(v ...interface{}) {
	logMsg := getLogrusMsg(l, v...)
	var fields = make(map[string]interface{})
	for _, with := range l.fields {
		fields[with.Key] = with.Val
	}
	l.fields = nil
	l.Logger.WithFields(fields).Warn(logMsg)
}

func (l *logrusLog) Warnf(format string, v ...interface{}) {
	var fields = make(map[string]interface{})
	for _, with := range l.fields {
		fields[with.Key] = with.Val
	}
	l.fields = nil
	l.Logger.WithFields(fields).Warnf(format, v...)
}

func (l *logrusLog) Debug(v ...interface{}) {
	logMsg := getLogrusMsg(l, v...)
	var fields = make(map[string]interface{})
	for _, with := range l.fields {
		fields[with.Key] = with.Val
	}
	l.fields = nil
	l.Logger.WithFields(fields).Debug(logMsg)
}

func (l *logrusLog) Debugf(format string, v ...interface{}) {
	var fields = make(map[string]interface{})
	for _, with := range l.fields {
		fields[with.Key] = with.Val
	}
	l.fields = nil
	l.Logger.WithFields(fields).Debugf(format, v...)
}

type DefaultFieldHook struct {
	c Config
}

func (hook *DefaultFieldHook) Fire(entry *logrus.Entry) error {
	if len(hook.c.ServiceName) > 0 {
		entry.Data["service"] = hook.c.ServiceName
		entry.Data["logId"] = utils.UUID()
	}

	if e, found := entry.Data[logrus.ErrorKey]; found {
		if _, ok := e.(error); ok {
			stack := debug.Stack()
			fmt.Fprintln(entry.Logger.Out, string(stack))
		}

	}

	return nil
}

func (hook *DefaultFieldHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func newLogrusLog(c Config) *logrusLog {
	log := logrus.New()
	log.SetLevel(getLogrusLevel(c.Level))
	log.SetFormatter(&logrus.TextFormatter{})
	//log.SetReportCaller(true)
	log.AddHook(&DefaultFieldHook{c})
	logPath := c.FilePath
	if len(logPath) == 0 {
		logPath = "./log"
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: logWriter(logPath, "Info", c.MaxAge, c.BackupNum), // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  logWriter(logPath, "Info", c.MaxAge, c.BackupNum),
		logrus.WarnLevel:  logWriter(logPath, "Info", c.MaxAge, c.BackupNum),
		logrus.ErrorLevel: logWriter(logPath, "Error", c.MaxAge, c.BackupNum),
		logrus.FatalLevel: logWriter(logPath, "Error", c.MaxAge, c.BackupNum),
		logrus.PanicLevel: logWriter(logPath, "Error", c.MaxAge, c.BackupNum),
	}, &logrus.TextFormatter{})

	log.AddHook(lfHook)
	return &logrusLog{Logger: log}
}

func logWriter(logPath string, level string, maxAge, backupNum int) *rotatelogs.RotateLogs {
	logFullPath := path.Join(logPath, level)
	rotateLogs, err := rotatelogs.New(
		logFullPath+".log",
		rotatelogs.WithLinkName(logFullPath+".log"),               // 生成软链，指向最新日志文件
		rotatelogs.WithRotationCount(uint(backupNum)),             // 文件最大保存份数
		rotatelogs.WithRotationTime(24*time.Hour),                 // 日志切割时间间隔
		rotatelogs.WithMaxAge(time.Duration(maxAge)*time.Hour*24), //保留天数
	)
	if err != nil {
		panic(err)
	}
	return rotateLogs
}

func getLogrusLevel(level string) logrus.Level {
	switch level {
	case "Info":
		return logrus.InfoLevel
	case "Debug":
		return logrus.DebugLevel
	case "Error":
		return logrus.ErrorLevel
	case "Panic":
		return logrus.PanicLevel
	case "fatal":
		return logrus.FatalLevel
	case "Warn":
		return logrus.WarnLevel
	}
	return logrus.InfoLevel
}
