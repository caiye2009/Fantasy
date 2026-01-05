package domain

import "errors"

var (
	ErrFabricProcessNotFound = errors.New("fabric_process not found")
	ErrFabricProcessInvalid  = errors.New("fabric_process invalid")
)