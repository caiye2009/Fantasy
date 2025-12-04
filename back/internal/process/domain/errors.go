package domain

import "errors"

var (
	ErrProcessNotFound      = errors.New("process not found")
	ErrProcessNameEmpty     = errors.New("process name cannot be empty")
	ErrProcessNameInvalid   = errors.New("process name must be between 2 and 100 characters")
	ErrProcessCannotDelete  = errors.New("process has related records, cannot delete")
)