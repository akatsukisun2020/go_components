package logger

import (
	"context"

	"github.com/google/uuid"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

type Field struct {
	Key   string
	Value string
}

// UdfLoggerContext 自定义日志key // value:map[string]string
type UdfLoggerContext struct {
	Fields []Field
}

type UdfLoggerContextKey struct{} // context的value需要的是 "comparable"

var udfLogHander *logger

// UdfLoggerConfig 自定义日志配置
type UdfLoggerConfig struct {
	FileName     string // 日志文件
	MaxFileSize  int64  // 最大文件长度，单位MB
	MaxBackups   int32  // 最大保留时长,单位天
	LevelEnabler int32  // 日志级别,控制输出的日志的最低级别, -1是debug日志
}

type logger struct {
	sugarLogger *zap.SugaredLogger
}

// InitLogger 日志初始化
func InitLogger(conf *UdfLoggerConfig) {
	if udfLogHander != nil {
		return
	}
	if conf == nil {
		conf = &UdfLoggerConfig{
			FileName:     "log/grpc.log",
			MaxFileSize:  1024,
			MaxBackups:   7,
			LevelEnabler: -1,
		}
	}

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   conf.FileName,
		MaxSize:    int(conf.MaxFileSize),
		MaxBackups: int(conf.MaxBackups),
		LocalTime:  true,
	})

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		w, // 输出到文件，文件名参数化
		zapcore.Level(conf.LevelEnabler),
	)

	// AddCallerSkip 表示打印的caller的栈的层次,1刚好能打印业务代码，因为底层包装了很多层
	zaplogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	udfLogHander = &logger{
		sugarLogger: zaplogger.Sugar(),
	}

	return
}

// main 退出时执行,日志刷盘
func DeferLogger() {
	if udfLogHander == nil {
		return
	}

	defer udfLogHander.sugarLogger.Sync()
}

func DebugContextf(ctx context.Context, template string, args ...interface{}) {
	values, ok := ctx.Value(UdfLoggerContextKey{}).(UdfLoggerContext)
	if ok || len(values.Fields) == 0 {

		var fields []interface{}
		for _, v := range values.Fields {
			fields = append(fields, v.Key, v.Value)
		}
		udfLogHander.sugarLogger.With(fields...).Debugf(template, args...)
	} else {
		udfLogHander.sugarLogger.Debugf(template, args...)
	}
}

func InfoContextf(ctx context.Context, template string, args ...interface{}) {
	values, ok := ctx.Value(UdfLoggerContextKey{}).(UdfLoggerContext)
	if ok || len(values.Fields) == 0 {
		var fields []interface{}
		for _, v := range values.Fields {
			fields = append(fields, v.Key, v.Value)
		}
		udfLogHander.sugarLogger.With(fields...).Infof(template, args...)
	} else {
		udfLogHander.sugarLogger.Infof(template, args...)
	}
}

func ErrorContextf(ctx context.Context, template string, args ...interface{}) {
	values, ok := ctx.Value(UdfLoggerContextKey{}).(UdfLoggerContext)
	if ok || len(values.Fields) == 0 {
		var fields []interface{}
		for _, v := range values.Fields {
			fields = append(fields, v.Key, v.Value)
		}
		udfLogHander.sugarLogger.With(fields...).Errorf(template, args...)
	} else {
		udfLogHander.sugarLogger.Errorf(template, args...)
	}
}

func FatalContextf(ctx context.Context, template string, args ...interface{}) {
	values, ok := ctx.Value(UdfLoggerContextKey{}).(UdfLoggerContext)
	if ok || len(values.Fields) == 0 {
		var fields []interface{}
		for _, v := range values.Fields {
			fields = append(fields, v.Key, v.Value)
		}
		udfLogHander.sugarLogger.With(fields...).Fatalf(template, args...)
	} else {
		udfLogHander.sugarLogger.Fatalf(template, args...)
	}
}

// UnaryLoggerInterceptor 日志拦截器,功能是，单机日志单个请求的链路能够串起来
func UnaryLoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	ctx = context.WithValue(ctx, UdfLoggerContextKey{}, UdfLoggerContext{Fields: []Field{{
		"trace_id", uuid.New().String(),
	}}})
	ctx = context.WithValue(ctx, "lala", "haha")

	m, err := handler(ctx, req)
	if err != nil {
		udfLogHander.sugarLogger.Errorf("RPC failed with error %v", err)
	}

	return m, err
}
