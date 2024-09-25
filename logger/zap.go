package logger

import (
	"context"
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/yunduansing/gtools/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"time"
)

type zapLog struct {
	*zap.Logger
	kv []KeyPair
}

func (log *zapLog) WithField(field, value string) {
	log.kv = append(log.kv, KeyPair{
		Key: field,
		Val: value,
	})
}

func (log *zapLog) close() {

}

func getZapLogLevel(level string) zapcore.Level {
	switch level {
	case "Info":
		return zapcore.InfoLevel
	case "Debug":
		return zapcore.DebugLevel
	case "Error":
		return zapcore.ErrorLevel
	case "Panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	case "Warn":
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

func getKVFromCtx(ctx context.Context) []KeyPair {
	var kv []KeyPair
	if ctx == nil {
		return kv
	}
	if ctx.Value("kv") != nil {
		kv = ctx.Value("kv").([]KeyPair)
	}
	if ctx.Value("requestId") != nil {
		kv = append(kv, KeyPair{"RequestId", ctx.Value("requestId")})
	}
	return kv
}

func (log *zapLog) Info(ctx context.Context, v ...interface{}) {
	msg := getMsg(v)

	var fields []zap.Field
	var kv = getKVFromCtx(ctx)

	for _, with := range kv {
		fields = append(fields, zap.Any(with.Key, with.Val))
	}

	log.Logger.Info(msg, fields...)
}

func (log *zapLog) Infof(ctx context.Context, format string, v ...interface{}) {
	var fields []zap.Field
	var kv = getKVFromCtx(ctx)

	for _, with := range kv {
		fields = append(fields, zap.Any(with.Key, with.Val))
	}

	log.Logger.Info(fmt.Sprintf(format, v...), fields...)
}

func (log *zapLog) Error(ctx context.Context, v ...interface{}) {
	msg := getMsg(v)

	var fields []zap.Field
	var kv = getKVFromCtx(ctx)

	for _, with := range kv {
		fields = append(fields, zap.Any(with.Key, with.Val))
	}

	log.Logger.Error(msg, fields...)
}

func (log *zapLog) Errorf(ctx context.Context, format string, v ...interface{}) {
	var fields []zap.Field
	var kv = getKVFromCtx(ctx)

	for _, with := range kv {
		fields = append(fields, zap.Any(with.Key, with.Val))
	}

	log.Logger.Error(fmt.Sprintf(format, v...), fields...)
}

func (log *zapLog) Panic(ctx context.Context, v ...interface{}) {
	msg := getMsg(v)

	var fields []zap.Field
	var kv = getKVFromCtx(ctx)

	for _, with := range kv {
		fields = append(fields, zap.Any(with.Key, with.Val))
	}

	log.Logger.Panic(msg, fields...)
}

func (log *zapLog) Panicf(ctx context.Context, format string, v ...interface{}) {
	var fields []zap.Field
	var kv = getKVFromCtx(ctx)

	for _, with := range kv {
		fields = append(fields, zap.Any(with.Key, with.Val))
	}

	log.Logger.Panic(fmt.Sprintf(format, v...), fields...)
}

func (log *zapLog) Warn(ctx context.Context, v ...interface{}) {
	msg := getMsg(v)

	var fields []zap.Field
	var kv = getKVFromCtx(ctx)

	for _, with := range kv {
		fields = append(fields, zap.Any(with.Key, with.Val))
	}

	log.Logger.Warn(msg, fields...)
}

func (log *zapLog) Warnf(ctx context.Context, format string, v ...interface{}) {
	var fields []zap.Field
	var kv = getKVFromCtx(ctx)

	for _, with := range kv {
		fields = append(fields, zap.Any(with.Key, with.Val))
	}

	log.Logger.Warn(fmt.Sprintf(format, v...), fields...)
}

func (log *zapLog) Debug(ctx context.Context, v ...interface{}) {
	msg := getMsg(v)

	var fields []zap.Field
	var kv = getKVFromCtx(ctx)

	for _, with := range kv {
		fields = append(fields, zap.Any(with.Key, with.Val))
	}

	log.Logger.Debug(msg, fields...)
}

func (log *zapLog) Debugf(ctx context.Context, format string, v ...interface{}) {
	var fields []zap.Field
	var kv = getKVFromCtx(ctx)

	for _, with := range kv {
		fields = append(fields, zap.Any(with.Key, with.Val))
	}

	log.Logger.Debug(fmt.Sprintf(format, v...), fields...)
}

func newZapLog(c Config) *zapLog {
	if len(c.FilePath) == 0 {
		c.FilePath = "./logs"
	}
	log := getLogWriter(c)

	return &zapLog{Logger: log}
}

func getLogContent(content interface{}) string {
	switch v := content.(type) {
	case error:
		return v.Error()
	case string:
		return v
	case int64, int, int32, float64, float32, bool:
		return fmt.Sprint(v)
	case zap.Field:
		if v.Integer > 0 {
			return fmt.Sprintf("%s:%d", v.Key, v.Integer)
		} else if v.Interface != nil {
			return fmt.Sprintf("%s:%s", v.Key, utils.ToJsonString(v.Interface))
		}
		return fmt.Sprintf("%s:%s", v.Key, v.String)
	case KeyPair:
		switch v.Val.(type) {
		case int64, int, int32, float64, float32, bool:
			return fmt.Sprintf("%s:", v.Key) + fmt.Sprint(v.Val)
		}
		return fmt.Sprintf("%s:%s", v.Key, utils.ToJsonString(v.Val))
	}
	return utils.ToJsonString(content)
}

func getLogWriter(c Config) *zap.Logger {
	var coreArr []zapcore.Core

	//获取编码器
	encoderConfig := zap.NewProductionEncoderConfig()                               //NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(utils.ChineseTimeLayout) //指定时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder                    //按级别显示不同颜色，不需要的话取值zapcore.CapitalLevelEncoder就可以了
	//encoderConfig.EncodeCaller = zapcore.FullCallerEncoder      	//显示完整文件路径
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime + ".000")
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	encoderFileConfig := zap.NewProductionEncoderConfig()                               //NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	encoderFileConfig.EncodeTime = zapcore.TimeEncoderOfLayout(utils.ChineseTimeLayout) //指定时间格式
	//encoderFileConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder                    //按级别显示不同颜色，不需要的话取值zapcore.CapitalLevelEncoder就可以了
	//encoderConfig.EncodeCaller = zapcore.FullCallerEncoder      	//显示完整文件路径
	encoderFileConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime + ".000")
	encoderFile := zapcore.NewConsoleEncoder(encoderFileConfig)

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
	infoFileCore := zapcore.NewCore(encoderFile, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer), infoPriority) //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志

	//error文件writeSyncer
	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/error.log", c.FilePath), //日志文件存放目录
		MaxSize:    c.MaxSize,                               //文件大小限制,单位MB
		MaxBackups: c.BackupNum,                             //最大保留日志文件数量
		MaxAge:     c.MaxAge,                                //日志文件保留天数
		Compress:   c.Compress,                              //是否压缩处理
	})
	errorFileCore := zapcore.NewCore(encoderFile, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer), errorPriority) //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志

	coreArr = append(coreArr, infoFileCore)
	coreArr = append(coreArr, errorFileCore)

	consoleCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.DebugLevel)
	coreArr = append(coreArr, consoleCore)

	var options []zap.Option
	if len(c.ServiceName) > 0 {
		options = append(options, zap.Fields(zap.String("service", c.ServiceName)))
	}
	options = append(options, []zap.Option{zap.AddCaller(), zap.AddCallerSkip(2),
		zap.AddStacktrace(zapcore.ErrorLevel)}...)

	return zap.New(zapcore.NewTee(coreArr...), options...) //zap.AddCaller()为显示文件名和行号，可省略
}
