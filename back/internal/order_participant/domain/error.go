package domain

import "errors"

var (
	ErrOrderParticipantNotFound = errors.New("order_participant not found")
	ErrOrderParticipantInvalid  = errors.New("order_participant invalid")
)
