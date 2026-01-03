package domain

import "errors"

var (
	ErrProductionOrderNotFound = errors.New("production_order not found")
	ErrProductionOrderInvalid  = errors.New("production_order invalid")
)
