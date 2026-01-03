package domain

import "errors"

var (
	ErrOrderMaterialNotFound = errors.New("order_material not found")
	ErrOrderMaterialInvalid  = errors.New("order_material invalid")
)
