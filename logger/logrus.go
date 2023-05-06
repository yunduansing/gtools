package logger

import "github.com/sirupsen/logrus"

type logrusLog struct {
	*logrus.Logger
}

func (l *logrusLog) info(v ...interface{}) {
	l.Info(v...)
}

func (l *logrusLog) infof(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *logrusLog) error(v ...interface{}) {
	l.Error(v...)
}

func (l *logrusLog) errorf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *logrusLog) panic(v ...interface{}) {
	l.Panic(v...)
}

func (l *logrusLog) panicf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *logrusLog) warn(v ...interface{}) {
	l.Warn(v...)
}

func (l *logrusLog) warnf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *logrusLog) debug(v ...interface{}) {
	l.Debug(v...)
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
	log.SetFormatter(&logrus.TextFormatter{})
	log.SetReportCaller(true)
	log.AddHook(&DefaultFieldHook{c})
	return &logrusLog{log}
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
