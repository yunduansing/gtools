package logger

type ILog interface {
	info(v ...interface{})
	infof(format string, v ...interface{})
	error(v ...interface{})
	errorf(format string, v ...interface{})
	panic(v ...interface{})
	panicf(format string, v ...interface{})
	warn(v ...interface{})
	warnf(format string, v ...interface{})
	debug(v ...interface{})
	debugf(format string, v ...interface{})
	close()
}
