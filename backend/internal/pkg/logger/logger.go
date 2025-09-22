package logger

import (
	"insight/config"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

var Logger *zap.Logger
var once sync.Once

func InitLogger() {
	once.Do(func() { Logger = createZapLog() })
}

func createZapLog() *zap.Logger {
	Config := config.GetConfig()
	if Config.System.Debug == true {
		if Logger, err := zap.NewDevelopment(); err == nil {
			return Logger
		} else {
			panic("Init Logger Failed, " + err.Error())
		}
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}

	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	baseFile, err := os.Getwd()
	if err != nil {
		panic("Failed to get working directory: " + err.Error())
	}
	filename := filepath.Join(baseFile, "/logs", Config.Logger.FileName)
	var writer zapcore.WriteSyncer
	if Config.Logger.DefaultDivision == "size" {
		// 按文件大小切割日志
		writer = zapcore.AddSync(getLumberJackWriter(filename, *Config))
	} else {
		// 按天切割日志
		writer = zapcore.AddSync(getRotateWriter(filename, *Config))
	}
	zapCore := zapcore.NewCore(encoder, writer, zap.InfoLevel)
	//zap.AddStacktrace(zap.WarnLevel)
	return zap.New(zapCore, zap.AddCaller())
}

// getRotateWriter 按日期切割日志
func getRotateWriter(filename string, config config.Config) io.Writer {
	maxAge := time.Duration(config.Logger.DivisionTime.MaxAge)
	rotationTime := time.Duration(config.Logger.DivisionTime.RotationTime)
	hook, err := rotatelogs.New(
		filename+".%Y%m%d",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*maxAge),
		rotatelogs.WithRotationTime(time.Hour*rotationTime), // 默认一天
	)
	if err != nil {
		panic(err)
	}
	return hook
}

// getLumberJackWriter 按文件切割日志
func getLumberJackWriter(filename string, config config.Config) io.Writer {
	// 日志切割配置
	return &lumberjack.Logger{
		Filename:   filename,                              // 日志文件位置
		MaxSize:    config.Logger.DivisionSize.MaxSize,    // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: config.Logger.DivisionSize.MaxBackups, // 保留旧文件的最大个数
		MaxAge:     config.Logger.DivisionSize.MaxAge,     // 保留旧文件的最大天数
		Compress:   config.Logger.DivisionSize.Compress,   // 是否压缩/归档旧文件
	}
}

func createInMemoryLogCore() (*zap.Logger, *observer.ObservedLogs) {
	core, logs := observer.New(zapcore.InfoLevel)
	return zap.New(core), logs
}
