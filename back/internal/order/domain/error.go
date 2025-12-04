package domain

import "errors"

var (
	ErrOrderNotFound           = errors.New("order not found")
	ErrOrderNoEmpty            = errors.New("order number cannot be empty")
	ErrOrderNoTooLong          = errors.New("order number cannot exceed 50 characters")
	ErrOrderNoDuplicate        = errors.New("order number already exists")
	ErrClientIDRequired        = errors.New("client_id is required")
	ErrProductIDRequired       = errors.New("product_id is required")
	ErrInvalidQuantity         = errors.New("quantity must be greater than 0")
	ErrInvalidUnitPrice        = errors.New("unit_price cannot be negative")
	ErrCreatedByRequired       = errors.New("created_by is required")
	ErrCannotConfirm           = errors.New("can only confirm pending orders")
	ErrCannotStartProduction   = errors.New("can only start production for confirmed orders")
	ErrCannotComplete          = errors.New("can only complete orders in production")
	ErrCannotCancelCompleted   = errors.New("cannot cancel completed orders")
	ErrCannotUpdateCompleted   = errors.New("cannot update completed or cancelled orders")
	ErrCannotDeleteCompleted   = errors.New("cannot delete completed orders")
)