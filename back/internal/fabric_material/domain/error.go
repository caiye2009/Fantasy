package domain

import "errors"

var (
	ErrFabricMaterialNotFound = errors.New("fabric_material not found")
	ErrFabricMaterialInvalid  = errors.New("fabric_material invalid")
)