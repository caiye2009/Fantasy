package auth

import (
	"context"
)

// ContextKey 自定义 context key 类型
type ContextKey string

const (
	// ContextKeyLoginID context 中存储 loginId 的 key
	ContextKeyLoginID ContextKey = "loginId"
	// ContextKeyRole context 中存储 role 的 key
	ContextKeyRole ContextKey = "role"
)

// GetLoginID 从 context 中获取 loginId
func GetLoginID(ctx context.Context) (string, bool) {
	loginId, ok := ctx.Value(ContextKeyLoginID).(string)
	return loginId, ok
}

// GetRole 从 context 中获取 role
func GetRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(ContextKeyRole).(string)
	return role, ok
}

// MustGetLoginID 从 context 中获取 loginId（如果不存在则返回空字符串）
func MustGetLoginID(ctx context.Context) string {
	loginId, _ := GetLoginID(ctx)
	return loginId
}

// MustGetRole 从 context 中获取 role（如果不存在则返回空字符串）
func MustGetRole(ctx context.Context) string {
	role, _ := GetRole(ctx)
	return role
}

// SetLoginID 设置 loginId 到 context
func SetLoginID(ctx context.Context, loginId string) context.Context {
	return context.WithValue(ctx, ContextKeyLoginID, loginId)
}

// SetRole 设置 role 到 context
func SetRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, ContextKeyRole, role)
}
