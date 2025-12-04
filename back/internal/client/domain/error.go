package domain

import "errors"

var (
	ErrClientNotFound    = errors.New("client not found")
	ErrClientNameEmpty   = errors.New("client name cannot be empty")
	ErrClientNameInvalid = errors.New("client name must be between 2 and 100 characters")
	ErrClientNameExists  = errors.New("client name already exists")
	ErrContactTooLong    = errors.New("contact cannot exceed 50 characters")
	ErrPhoneTooLong      = errors.New("phone cannot exceed 20 characters")
	ErrPhoneInvalid      = errors.New("invalid phone format")
	ErrEmailTooLong      = errors.New("email cannot exceed 100 characters")
	ErrEmailInvalid      = errors.New("invalid email format")
	ErrAddressTooLong    = errors.New("address cannot exceed 200 characters")
)