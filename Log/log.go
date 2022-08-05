package Log

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// global

var Logger *zap.SugaredLogger //项目公共logger

// InitLogger 日志初始化
func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	Logger = logger.Sugar()
}

// getEncoder 设置日志编码方式
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// getLogWriter 设置log地址
func getLogWriter() zapcore.WriteSyncer {
	// 日志文件每1MB会切割并且在当前目录下最多保存5个备份
	logger := lumberjack.Logger{
		Filename:   "./Log/LogFile/Tsubaki.log", //Filename: 日志文件的位置
		MaxSize:    1,                           //MaxSize：在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: 5,                           //MaxBackups：保留旧文件的最大个数
		MaxAge:     30,                          //MaxAges：保留旧文件的最大天数
		Compress:   false,                       //Compress：是否压缩/归档旧文件
	}
	return zapcore.AddSync(&logger)
}
