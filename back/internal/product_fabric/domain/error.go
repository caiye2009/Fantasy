package domain

import "errors"

var (
	ErrProductFabricNotFound = errors.New("product_fabric not found")
	ErrProductFabricInvalid  = errors.New("product_fabric invalid")
)