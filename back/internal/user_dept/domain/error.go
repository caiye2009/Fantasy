package domain

import "errors"

var (
	ErrUserDeptNotFound     = errors.New("user_dept not found")
	ErrUserDeptCodeExists   = errors.New("dept_code already exists")
	ErrUserDeptInvalidInput = errors.New("invalid input")
)
