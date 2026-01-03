package domain

import "errors"

var (
	ErrPurchaseOrderNotFound = errors.New("purchase_order not found")
	ErrPurchaseOrderInvalid  = errors.New("purchase_order invalid")
)
