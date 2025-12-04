package domain

import "errors"

var (
	ErrSupplierNotFound    = errors.New("supplier not found")
	ErrSupplierNameEmpty   = errors.New("supplier name cannot be empty")
	ErrSupplierNameInvalid = errors.New("supplier name must be between 2 and 100 characters")
	ErrSupplierNameExists  = errors.New("supplier name already exists")
	ErrContactTooLong      = errors.New("contact cannot exceed 50 characters")
	ErrPhoneTooLong        = errors.New("phone cannot exceed 20 characters")
	ErrPhoneInvalid        = errors.New("invalid phone format")
	ErrEmailTooLong        = errors.New("email cannot exceed 100 characters")
	ErrEmailInvalid        = errors.New("invalid email format")
	ErrAddressTooLong      = errors.New("address cannot exceed 200 characters")
)