package domain

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidSize     = errors.New("invalid page size")
	ErrSizeTooLarge    = errors.New("page size exceeds maximum limit of 100")
	ErrInvalidFrom     = errors.New("invalid from offset")
	ErrIndexNotFound   = errors.New("index not found")
	ErrSearchFailed    = errors.New("search operation failed")
	ErrInvalidQuery    = errors.New("invalid search query")
)

// ErrInvalidFilterField 无效的过滤字段错误
func ErrInvalidFilterField(field, index string) error {
	return fmt.Errorf("field '%s' is not filterable in index '%s'", field, index)
}

// ErrInvalidSortField 无效的排序字段错误
func ErrInvalidSortField(field, index string) error {
	return fmt.Errorf("field '%s' is not sortable in index '%s'", field, index)
}