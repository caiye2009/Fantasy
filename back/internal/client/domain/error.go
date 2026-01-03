package domain

import "errors"

var (
	ErrClientNotFound = errors.New("client not found")
	ErrClientInvalid  = errors.New("client invalid")
)
