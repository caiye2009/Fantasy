package domain

import "errors"

var (
	ErrSupplierNotFound = errors.New("supplier not found")
	ErrSupplierInvalid  = errors.New("supplier invalid")
)
