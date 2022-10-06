package logger

import "github.com/sirupsen/logrus"

type logrusLog struct {
	*logrus.Logger
}

func newLogrusLog(c Config) *logrusLog {
	log := logrus.New()
	log.SetLevel(getLevel(c.Level))
	return &logrusLog{log}
}

func getLevel(level string) logrus.Level {
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

func (log logrusLog) Info(v ...interface{}) {
	log.Info(v)
}
func (log logrusLog) Infof(format string, v ...interface{}) {
	log.Infof(format, v)
}
func (log logrusLog) Error(v ...interface{}) {
	log.Error(v)
}
func (log logrusLog) Errorf(format string, v ...interface{}) {
	log.Errorf(format, v)
}
func (log logrusLog) Panic(v ...interface{}) {
	log.Panic(v)
}
func (log logrusLog) Panicf(format string, v ...interface{}) {
	log.Panicf(format, v)
}
func (log logrusLog) Warn(v ...interface{}) {

}
func (log logrusLog) Warnf(format string, v ...interface{}) {

}
func (log logrusLog) Debug(v ...interface{}) {

}
func (log logrusLog) Debugf(format string, v ...interface{}) {

}
