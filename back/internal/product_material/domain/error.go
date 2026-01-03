package domain

import "errors"

var (
	ErrProductMaterialNotFound = errors.New("product_material not found")
	ErrProductMaterialInvalid  = errors.New("product_material invalid")
)
