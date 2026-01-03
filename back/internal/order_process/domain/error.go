package domain

import "errors"

var (
	ErrOrderProcessNotFound = errors.New("order_process not found")
	ErrOrderProcessInvalid  = errors.New("order_process invalid")
)
