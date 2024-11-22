package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

func init() {
	// 创建一个自定义的 zap 配置
	config := zap.NewProductionConfig()

	// 修改配置项，比如设置日志级别、输出路径等
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	config.OutputPaths = []string{"stdout"}

	// 创建一个自定义的 encoder 配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 创建一个自定义的 zap logger
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		zap.InfoLevel,
	)
	//logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	logger = zap.New(core)
	sugar = logger.Sugar()
}

// Info  记录格式化的信息级别的日志
func Info(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}

// Error 记录格式化的错误级别的日志
func Error(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}

// Panic 记录格式化的 panic 级别的日志
func Panic(template string, args ...interface{}) {
	sugar.Panicf(template, args...)
}

// Sync 刷新日志缓冲区
func Sync() {
	_ = sugar.Sync()
}
