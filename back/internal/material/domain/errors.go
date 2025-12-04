package domain

import "errors"

var (
	ErrMaterialNotFound      = errors.New("material not found")
	ErrMaterialNameEmpty     = errors.New("material name cannot be empty")
	ErrMaterialNameInvalid   = errors.New("material name must be between 2 and 100 characters")
	ErrMaterialSpecTooLong   = errors.New("material spec cannot exceed 200 characters")
	ErrMaterialUnitTooLong   = errors.New("material unit cannot exceed 20 characters")
	ErrMaterialCannotDelete  = errors.New("material has related records, cannot delete")
)