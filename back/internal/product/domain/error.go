package domain

import "errors"

var (
	ErrProductNotFound = errors.New("product not found")
	ErrProductInvalid  = errors.New("product invalid")
)
