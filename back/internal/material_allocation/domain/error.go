package domain

import "errors"

var (
	ErrMaterialAllocationNotFound = errors.New("material_allocation not found")
	ErrMaterialAllocationInvalid  = errors.New("material_allocation invalid")
)
