package application

import "errors"

var (
	ErrInvalidCredentials  = errors.New("用户名或密码错误")
	ErrAccountSuspended    = errors.New("账户已被停用")
	ErrTokenGenerateFailed = errors.New("生成 Token 失败")
	ErrInvalidToken        = errors.New("Token 无效或已过期")
	ErrUserNotFound        = errors.New("用户不存在")
)