package domain

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserInvalid  = errors.New("user invalid")
)
