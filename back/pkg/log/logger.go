package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *Logger

type Logger struct {
	zap *zap.Logger
}

// Init 初始化全局 logger
func Init(level string, format string) error {
	var zapConfig zap.Config

	if format == "json" {
		zapConfig = zap.NewProductionConfig()
	} else {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// 设置日志级别
	switch level {
	case "debug":
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	zapLogger, err := zapConfig.Build()
	if err != nil {
		return err
	}

	globalLogger = &Logger{zap: zapLogger}
	return nil
}

// NewLogger 创建新的 logger 实例
func NewLogger() *Logger {
	if globalLogger == nil {
		// 如果没有初始化,使用默认配置
		zapLogger, _ := zap.NewDevelopment()
		return &Logger{zap: zapLogger}
	}
	return globalLogger
}

// With 创建带上下文的新 logger
func (l *Logger) With(fields ...Field) *Logger {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = f.zapField
	}
	return &Logger{
		zap: l.zap.With(zapFields...),
	}
}

// Info 记录 info 级别日志
func (l *Logger) Info(msg string, fields ...Field) {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = f.zapField
	}
	l.zap.Info(msg, zapFields...)
}

// Warn 记录 warn 级别日志
func (l *Logger) Warn(msg string, fields ...Field) {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = f.zapField
	}
	l.zap.Warn(msg, zapFields...)
}

// Error 记录 error 级别日志
func (l *Logger) Error(msg string, fields ...Field) {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = f.zapField
	}
	l.zap.Error(msg, zapFields...)
}

// Debug 记录 debug 级别日志
func (l *Logger) Debug(msg string, fields ...Field) {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = f.zapField
	}
	l.zap.Debug(msg, zapFields...)
}

// Sync 同步日志 (程序退出前调用)
func (l *Logger) Sync() error {
	return l.zap.Sync()
}