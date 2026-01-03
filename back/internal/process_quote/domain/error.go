package domain

import "errors"

var (
	ErrProcessQuoteNotFound = errors.New("process_quote not found")
	ErrProcessQuoteInvalid  = errors.New("process_quote invalid")
)
