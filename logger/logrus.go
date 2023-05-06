package logger

import "github.com/sirupsen/logrus"

type logrusLog struct {
	*logrus.Logger
}

func (l logrusLog) info(v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l logrusLog) infof(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l logrusLog) error(v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l logrusLog) errorf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l logrusLog) panic(v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l logrusLog) panicf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l logrusLog) warn(v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l logrusLog) warnf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l logrusLog) debug(v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l logrusLog) debugf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func newLogrusLog(c Config) *logrusLog {
	log := logrus.New()
	log.SetLevel(getLogrusLevel(c.Level))
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
