package domain

import "errors"

var (
	ErrOrderProductNotFound = errors.New("order_product not found")
	ErrOrderProductInvalid  = errors.New("order_product invalid")
)
