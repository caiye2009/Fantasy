package domain

import "errors"

var (
	ErrMaterialNotFound = errors.New("material not found")
	ErrMaterialInvalid  = errors.New("material invalid")
)
