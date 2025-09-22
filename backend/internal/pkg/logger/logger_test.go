package logger

import (
	"insight/config"
	"insight/config/autoload"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// 创建一个简单的开发日志记录器用于测试
func TestDevelopmentLogger(t *testing.T) {
	// 创建内存日志核心
	core, logs := observer.New(zapcore.InfoLevel)
	testLogger := zap.New(core)

	// 保存原始日志记录器并在测试后恢复
	origLogger := Logger
	defer func() { Logger = origLogger }()

	// 设置测试日志记录器
	Logger = testLogger

	// 测试记录消息
	Logger.Info("Test info message")
	Logger.Error("Test error message")

	// 验证日志记录
	assert.Equal(t, 2, logs.Len(), "Should have recorded 2 log messages")
	assert.Equal(t, "Test info message", logs.All()[0].Message)
	assert.Equal(t, "Test error message", logs.All()[1].Message)
}

// 测试单例模式
func TestLoggerSingleton(t *testing.T) {
	// 模拟 once.Do 行为
	once.Do(func() {
		// 使用简单的日志记录器而非调用 createZapLog
		logger, _ := zap.NewDevelopment()
		Logger = logger
	})

	// 保存引用
	firstLogger := Logger

	// 再次"初始化"，不应该有效果
	once.Do(func() {
		panic("This should not be executed")
	})

	// 验证仍然是相同的实例
	assert.Equal(t, firstLogger, Logger, "Logger should remain the same instance")
}

// 测试内存日志输出
func TestInMemoryLogger(t *testing.T) {
	// 创建内存日志记录器
	logger, logs := createInMemoryLogCore()

	// 写入各种级别的日志
	logger.Debug("Debug message") // 这个不应该被记录，因为级别是 Info
	logger.Info("Info message")
	logger.Warn("Warning message")
	logger.Error("Error message")

	// 验证日志数量
	assert.Equal(t, 3, logs.Len(), "Should have 3 log entries (Info, Warn, Error)")

	// 验证日志内容和级别
	allLogs := logs.All()
	assert.Equal(t, "Info message", allLogs[0].Message, "First log should be Info")
	assert.Equal(t, zapcore.InfoLevel, allLogs[0].Level, "First log level should be Info")

	assert.Equal(t, "Warning message", allLogs[1].Message, "Second log should be Warning")
	assert.Equal(t, zapcore.WarnLevel, allLogs[1].Level, "Second log level should be Warn")

	assert.Equal(t, "Error message", allLogs[2].Message, "Third log should be Error")
	assert.Equal(t, zapcore.ErrorLevel, allLogs[2].Level, "Third log level should be Error")
}

// 测试结构化日志字段
func TestStructuredLogging(t *testing.T) {
	// 创建内存日志记录器
	logger, logs := createInMemoryLogCore()

	// 写入带有结构化字段的日志
	logger.Info("User logged in",
		zap.String("username", "testuser"),
		zap.String("ip", "192.168.1.1"),
		zap.Int64("timestamp", time.Now().Unix()),
	)

	// 验证日志包含结构化字段
	allLogs := logs.All()
	assert.Equal(t, 1, logs.Len(), "Should have 1 log entry")

	fields := allLogs[0].Context
	assert.Equal(t, 3, len(fields), "Log should have 3 fields")

	// 验证字段名称和类型
	assert.Equal(t, "username", fields[0].Key, "First field should be username")
	assert.Equal(t, "testuser", fields[0].String, "Username should be testuser")

	assert.Equal(t, "ip", fields[1].Key, "Second field should be ip")
	assert.Equal(t, "192.168.1.1", fields[1].String, "IP should be 192.168.1.1")

	assert.Equal(t, "timestamp", fields[2].Key, "Third field should be timestamp")
}

// 测试轮转写入器创建
func TestRotateWriters(t *testing.T) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "logger_test")
	assert.NoError(t, err, "Should create temp directory")
	defer os.RemoveAll(tmpDir)

	// 创建测试配置
	testConfig := config.Config{
		System: autoload.SystemConfig{
			Debug: false,
		},
		Logger: autoload.LoggerConfig{
			FileName: "test.log",
			DivisionSize: autoload.DivisionSize{
				MaxSize:    1,
				MaxBackups: 3,
				MaxAge:     7,
				Compress:   false,
			},
			DivisionTime: autoload.DivisionTime{
				MaxAge:       7,
				RotationTime: 24,
			},
		},
	}

	// 测试路径
	testPath := filepath.Join(tmpDir, "test.log")

	// 测试大小轮转
	sizeWriter := getLumberJackWriter(testPath, testConfig)
	assert.NotNil(t, sizeWriter, "Size rotation writer should be created")

	// 测试时间轮转
	timeWriter := getRotateWriter(testPath, testConfig)
	assert.NotNil(t, timeWriter, "Time rotation writer should be created")
}

// 测试使用项目中的getLumberJackWriter写入日志
func TestGetLumberJackWriterWithFile(t *testing.T) {
	// 创建临时目录用于测试
	tmpDir, err := os.MkdirTemp("", "logger_project_test")
	assert.NoError(t, err, "Should create temp directory")
	defer os.RemoveAll(tmpDir)

	// 设置日志文件路径
	logFilePath := filepath.Join(tmpDir, "project.log")

	// 创建测试配置
	testConfig := config.Config{
		System: autoload.SystemConfig{
			Debug: false,
		},
		Logger: autoload.LoggerConfig{
			FileName: "test.log",
			DivisionSize: autoload.DivisionSize{
				MaxSize:    1,
				MaxBackups: 3,
				MaxAge:     7,
				Compress:   false,
			},
		},
	}

	// 获取项目中的日志写入器
	writer := getLumberJackWriter(logFilePath, testConfig)

	// 创建编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// 创建Core
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(writer),
		zap.InfoLevel,
	)

	// 创建Logger
	projectLogger := zap.New(core)

	// 写入测试日志
	testMessage := "This is a test using project's writer"
	projectLogger.Info(testMessage, zap.String("source", "project_test"))
	projectLogger.Sync() // 确保写入磁盘

	// 读取日志文件内容
	content, err := os.ReadFile(logFilePath)
	assert.NoError(t, err, "Should read log file")

	// 验证日志内容
	assert.Contains(t, string(content), testMessage, "Log file should contain test message")
	assert.Contains(t, string(content), "source=project_test", "Log file should contain field")

	// 显示文件内容
	t.Logf("Successfully wrote and verified project log file at: %s", logFilePath)
	t.Logf("Log content: %s", string(content))
}
