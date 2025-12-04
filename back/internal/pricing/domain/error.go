package domain

import "errors"

var (
	ErrTargetNotFound      = errors.New("target not found")
	ErrTargetTypeRequired  = errors.New("target_type is required")
	ErrInvalidTargetType   = errors.New("invalid target_type")
	ErrTargetIDRequired    = errors.New("target_id is required")
	ErrSupplierIDRequired  = errors.New("supplier_id is required")
	ErrInvalidPrice        = errors.New("price must be greater than 0")
	ErrPriceNotFound       = errors.New("price not found")
)