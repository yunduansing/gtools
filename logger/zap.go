package logger

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/yunduansing/gtools/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

var Logger *zap.Logger

type zapLog struct {
	*zap.Logger
}

func (log *zapLog) close() {
	//TODO implement me
	panic("implement me")
}

func getZapLogLevel(level string) zapcore.Level {
	switch level {
	case "info":
		return zapcore.InfoLevel
	case "debug":
		return zapcore.DebugLevel
	case "error":
		return zapcore.ErrorLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	case "warn":
		return zapcore.WarnLevel
	}
	return zapcore.InfoLevel
}

func getMsg(v ...interface{}) string {
	list := v[0].([]interface{})
	var msg strings.Builder
	for _, item := range list {
		contentItem := getLogContent(item)
		msg.WriteString(contentItem + " ")
	}
	return msg.String()
}

func (log *zapLog) info(v ...interface{}) {
	msg := getMsg(v)
	log.Info(msg)
}

func (log *zapLog) infof(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (log *zapLog) error(v ...interface{}) {
	msg := getMsg(v)
	log.Error(msg)
}

func (log *zapLog) errorf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (log *zapLog) panic(v ...interface{}) {
	msg := getMsg(v)
	log.Panic(msg)
}

func (log *zapLog) panicf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (log *zapLog) warn(v ...interface{}) {
	msg := getMsg(v)
	log.Warn(msg)
}

func (log *zapLog) warnf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (log *zapLog) debug(v ...interface{}) {
	msg := getMsg(v)
	log.Debug(msg)
}

func (log *zapLog) debugf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func newZapLog(c Config) *zapLog {
	log := getLogWriter(c)

	return &zapLog{log}
}

func getLogContent(content interface{}) string {
	switch v := content.(type) {
	case string:
		return v
	case int64, int, int32, float64, float32, bool:
		return fmt.Sprint(v)
	case zap.Field:
		if v.Integer > 0 {
			return fmt.Sprintf("%s:%d", v.Key, v.Integer)
		} else if v.Interface != nil {
			return fmt.Sprintf("%s:%s", v.Key, utils.ToJson(v.Interface))
		}
		return fmt.Sprintf("%s:%s", v.Key, v.String)
	case KeyPair:
		switch v.Val.(type) {
		case int64, int, int32, float64, float32, bool:
			return fmt.Sprintf("%s:", v.Key) + fmt.Sprint(v.Val)
		}
		return fmt.Sprintf("%s:%s", v.Key, utils.ToJson(v.Val))
	}
	return ""
}

func InitLogger() {
	Logger = getLogWriter(Config{
		Level: "info",
	})
}

//func Info(summary string, fields ...zap.Field) {
//	Logger.Info(summary, fields...)
//}
//
//func Error(summary string, fields ...zap.Field) {
//	Logger.Error(summary, fields...)
//}
//
//func Debug(summary string, fields ...zap.Field) {
//	Logger.Debug(summary, fields...)
//}
//
//func Fatal(summary string, fields ...zap.Field) {
//	Logger.Fatal(summary, fields...)
//}
//
//func Warn(summary string, fields ...zap.Field) {
//	Logger.Warn(summary, fields...)
//}
//
//func Panic(summary string, fields ...zap.Field) {
//	Logger.Panic(summary, fields...)
//}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(c Config) *zap.Logger {
	var coreArr []zapcore.Core

	//获取编码器
	encoderConfig := zap.NewProductionEncoderConfig()            //NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        //指定时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder //按级别显示不同颜色，不需要的话取值zapcore.CapitalLevelEncoder就可以了
	//encoderConfig.EncodeCaller = zapcore.FullCallerEncoder      	//显示完整文件路径
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	//日志级别
	errorPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //error级别
		return lev >= zap.ErrorLevel
	})
	infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //info和debug级别,debug级别是最低的
		return lev < zap.ErrorLevel && lev >= getZapLogLevel(c.Level)
	})

	//info文件writeSyncer
	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/info.log", c.FilePath), //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    c.MaxSize,                              //文件大小限制,单位MB
		MaxBackups: c.BackupNum,                            //最大保留日志文件数量
		MaxAge:     c.MaxAge,                               //日志文件保留天数
		Compress:   c.Compress,                             //是否压缩处理
	})
	infoFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), infoPriority) //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志

	//error文件writeSyncer
	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/error.log", c.FilePath), //日志文件存放目录
		MaxSize:    c.MaxSize,                               //文件大小限制,单位MB
		MaxBackups: c.BackupNum,                             //最大保留日志文件数量
		MaxAge:     c.MaxAge,                                //日志文件保留天数
		Compress:   c.Compress,                              //是否压缩处理
	})
	errorFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer, zapcore.AddSync(os.Stdout)), errorPriority) //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志

	coreArr = append(coreArr, infoFileCore)
	coreArr = append(coreArr, errorFileCore)

	var options []zap.Option
	if len(c.ServiceName) > 0 {
		options = append(options, zap.Fields(zap.String("service", c.ServiceName)))
	}
	options = append(options, []zap.Option{zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)}...)

	return zap.New(zapcore.NewTee(coreArr...), options...) //zap.AddCaller()为显示文件名和行号，可省略
}
