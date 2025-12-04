package domain

import "errors"

var (
	ErrInvalidSize     = errors.New("invalid page size")
	ErrSizeTooLarge    = errors.New("page size exceeds maximum limit of 100")
	ErrInvalidFrom     = errors.New("invalid from offset")
	ErrIndexNotFound   = errors.New("index not found")
	ErrSearchFailed    = errors.New("search operation failed")
	ErrInvalidQuery    = errors.New("invalid search query")
)