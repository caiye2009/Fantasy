package domain

import "errors"

var (
	// 订单基础错误
	ErrOrderNotFound           = errors.New("order not found")
	ErrOrderNoEmpty            = errors.New("order number cannot be empty")
	ErrOrderNoTooLong          = errors.New("order number cannot exceed 50 characters")
	ErrOrderNoDuplicate        = errors.New("order number already exists")
	ErrClientIDRequired        = errors.New("client_id is required")
	ErrProductIDRequired       = errors.New("product_id is required")
	ErrInvalidQuantity         = errors.New("quantity must be greater than 0")
	ErrInvalidUnitPrice        = errors.New("unit_price cannot be negative")
	ErrCreatedByRequired       = errors.New("created_by is required")
	ErrInvalidOrderStatus      = errors.New("invalid order status for this operation")
	ErrCannotConfirm           = errors.New("can only confirm pending orders")
	ErrCannotStartProduction   = errors.New("can only start production for confirmed orders")
	ErrCannotComplete          = errors.New("can only complete orders in production")
	ErrCannotCancelCompleted   = errors.New("cannot cancel completed orders")
	ErrCannotUpdateCompleted   = errors.New("cannot update completed or cancelled orders")
	ErrCannotDeleteCompleted   = errors.New("cannot delete completed orders")

	// 参与者相关错误
	ErrParticipantNotFound     = errors.New("participant not found")
	ErrParticipantAlreadyExists = errors.New("participant already exists")
	ErrParticipantRequired     = errors.New("participant is required")
	ErrInvalidParticipantRole  = errors.New("invalid participant role")

	// 进度相关错误
	ErrProgressNotFound        = errors.New("progress not found")
	ErrProgressNotExists       = errors.New("progress does not exist")
	ErrProgressAlreadyExists   = errors.New("progress already exists")
	ErrInvalidProgressType     = errors.New("invalid progress type")

	// 事件相关错误
	ErrEventNotFound           = errors.New("event not found")
	ErrInvalidEventType        = errors.New("invalid event type")

	// 权限相关错误
	ErrPermissionDenied        = errors.New("permission denied")
	ErrNotParticipant          = errors.New("user is not a participant of this order")

	// 分配相关错误
	ErrDepartmentRequired      = errors.New("department is required")
	ErrDepartmentAlreadyAssigned = errors.New("department already assigned")
	ErrPersonnelAlreadyAssigned = errors.New("personnel already assigned")
	ErrInvalidFabricTarget     = errors.New("fabric target quantity must be greater than 0")
)