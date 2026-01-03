package domain

import "errors"

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrOrderInvalid  = errors.New("order invalid")
)
