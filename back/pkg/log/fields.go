package log

import (
	"time"
	"go.uber.org/zap"
)

// Field 日志字段
type Field struct {
	zapField zap.Field
}

// String 创建字符串字段
func String(key, value string) Field {
	return Field{zapField: zap.String(key, value)}
}

// Int 创建整数字段
func Int(key string, value int) Field {
	return Field{zapField: zap.Int(key, value)}
}

// Int64 创建 int64 字段
func Int64(key string, value int64) Field {
	return Field{zapField: zap.Int64(key, value)}
}

// Uint 创建无符号整数字段
func Uint(key string, value uint) Field {
	return Field{zapField: zap.Uint(key, value)}
}

// Float64 创建浮点数字段
func Float64(key string, value float64) Field {
	return Field{zapField: zap.Float64(key, value)}
}

// Bool 创建布尔字段
func Bool(key string, value bool) Field {
	return Field{zapField: zap.Bool(key, value)}
}

// Duration 创建时间间隔字段
func Duration(key string, value time.Duration) Field {
	return Field{zapField: zap.Duration(key, value)}
}

// Time 创建时间字段
func Time(key string, value time.Time) Field {
	return Field{zapField: zap.Time(key, value)}
}

// Error 创建错误字段
func Error(err error) Field {
	return Field{zapField: zap.Error(err)}
}

// Any 创建任意类型字段
func Any(key string, value interface{}) Field {
	return Field{zapField: zap.Any(key, value)}
}