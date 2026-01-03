package domain

import "errors"

var (
	ErrMaterialQuoteNotFound = errors.New("material_quote not found")
	ErrMaterialQuoteInvalid  = errors.New("material_quote invalid")
)
