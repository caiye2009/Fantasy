package domain

import "errors"

var (
	ErrFabricNotFound  = errors.New("fabric not found")
	ErrFabricInvalid   = errors.New("fabric invalid")
)