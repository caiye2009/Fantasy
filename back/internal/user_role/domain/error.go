package domain

import "errors"

var (
	ErrUserRoleNotFound     = errors.New("user_role not found")
	ErrUserRoleCodeExists   = errors.New("role_code already exists")
	ErrUserRoleInvalidInput = errors.New("invalid input")
)
