package domain

import "errors"

var (
	ErrMaterialShrinkageNotFound = errors.New("material_shrinkage material not found")
	ErrSupplierShrinkageNotFound  = errors.New("material_shrinkage supplier not found")
)
