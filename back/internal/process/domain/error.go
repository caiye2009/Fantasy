package domain

import "errors"

var (
	ErrProcessNotFound = errors.New("process not found")
	ErrProcessInvalid  = errors.New("process invalid")
)
