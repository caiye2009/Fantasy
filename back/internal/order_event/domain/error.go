package domain

import "errors"

var (
	ErrOrderEventNotFound = errors.New("order_event not found")
	ErrOrderEventInvalid  = errors.New("order_event invalid")
)
