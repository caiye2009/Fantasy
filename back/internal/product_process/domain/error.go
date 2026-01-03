package domain

import "errors"

var (
	ErrProductProcessNotFound = errors.New("product_process not found")
	ErrProductProcessInvalid  = errors.New("product_process invalid")
)
