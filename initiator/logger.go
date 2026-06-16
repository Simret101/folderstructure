package initiator

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"folderstructure/platform/logger"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func getLogFilePath(level string) string {
	now := time.Now()
	baseDir := viper.GetString("app.log_path")
	if baseDir == "" {
		baseDir = "logs"
	}

	dir := filepath.Join(
		baseDir,
		now.Format("2006"),
		now.Format("01"),
		now.Format("02"),
		level,
	)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Printf("Failed to create log directory: %v", err)
		return filepath.Join(baseDir, level+".log")
	}

	return filepath.Join(dir, level+".log")
}

func NewFileLoggerWithLevel() *zap.Logger {
	lvl := zapcore.Level(viper.GetInt("logger.level"))
	if lvl < zapcore.DebugLevel || lvl > zapcore.FatalLevel {
		lvl = zapcore.InfoLevel
	}

	infoPath := getLogFilePath("info")
	errorPath := getLogFilePath("error")

	fileEncoderConfig := zap.NewProductionEncoderConfig()
	fileEncoderConfig.TimeKey = "time"
	fileEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	infoWriter := &lumberjack.Logger{
		Filename:   infoPath,
		MaxSize:    50,
		MaxBackups: 10,
		MaxAge:     7,
		Compress:   true,
	}

	errorWriter := &lumberjack.Logger{
		Filename:   errorPath,
		MaxSize:    50,
		MaxBackups: 20,
		MaxAge:     14,
		Compress:   true,
	}

	infoCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(infoWriter),
		zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= lvl && l < zapcore.ErrorLevel
		}),
	)

	errorCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(errorWriter),
		zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= zapcore.ErrorLevel
		}),
	)

	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout),
		zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= lvl
		}),
	)

	core := zapcore.NewTee(infoCore, errorCore, consoleCore)

	zapLogger := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	if viper.GetBool("logger.sampling") {
		zapLogger = zapLogger.WithOptions(
			zap.WrapCore(func(c zapcore.Core) zapcore.Core {
				return zapcore.NewSamplerWithOptions(c, time.Second, 100, 10)
			}),
		)
	}

	return zapLogger
}

func InitLogger() logger.Logger {
	zapLogger := NewFileLoggerWithLevel()
	return logger.New(zapLogger)
}
