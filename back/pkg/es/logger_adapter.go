package es

import (
	applog "back/pkg/log"
)

// LoggerAdapter 日志适配器
type LoggerAdapter struct {
	logger *applog.Logger
}

// NewLoggerAdapter 创建日志适配器
func NewLoggerAdapter(logger *applog.Logger) Logger {
	return &LoggerAdapter{logger: logger}
}

// Info 信息日志
func (a *LoggerAdapter) Info(msg string, fields ...interface{}) {
	// 将 fields 转换为 log.Field
	logFields := make([]applog.Field, 0, len(fields)/2)
	for i := 0; i < len(fields)-1; i += 2 {
		if key, ok := fields[i].(string); ok {
			logFields = append(logFields, applog.Any(key, fields[i+1]))
		}
	}
	a.logger.Info(msg, logFields...)
}

// Error 错误日志
func (a *LoggerAdapter) Error(msg string, fields ...interface{}) {
	// 将 fields 转换为 log.Field
	logFields := make([]applog.Field, 0, len(fields)/2)
	for i := 0; i < len(fields)-1; i += 2 {
		if key, ok := fields[i].(string); ok {
			logFields = append(logFields, applog.Any(key, fields[i+1]))
		}
	}
	a.logger.Error(msg, logFields...)
}