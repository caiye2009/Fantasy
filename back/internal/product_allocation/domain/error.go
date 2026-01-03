package domain

import "errors"

var (
	ErrProductAllocationNotFound = errors.New("product_allocation not found")
	ErrProductAllocationInvalid  = errors.New("product_allocation invalid")
)
