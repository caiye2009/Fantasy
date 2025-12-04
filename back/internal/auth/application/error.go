package application

import "errors"

var (
	ErrInvalidCredentials = errors.New("工号或密码错误")
	ErrAccountSuspended   = errors.New("账号已被停用")
	ErrTokenGenerateFailed = errors.New("生成token失败")
)