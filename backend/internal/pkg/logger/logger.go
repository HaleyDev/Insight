package logger

import (
	"insight/config"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var once sync.Once

func InitLogger() {
	once.Do(func() { Logger = createZapLog() })
}

func createZapLog() *zap.Logger {

	if config.GetConfig().System.Debug == true {
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
	filename := filepath.Join(baseFile, "/logs", time.Now().Format("2006-01-02")+".log")
	var writer zapcore.WriteSyncer

}
