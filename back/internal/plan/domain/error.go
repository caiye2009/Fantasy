package domain

import "errors"

var (
	ErrPlanNotFound              = errors.New("plan not found")
	ErrPlanNoEmpty               = errors.New("plan number cannot be empty")
	ErrPlanNoTooLong             = errors.New("plan number cannot exceed 50 characters")
	ErrPlanNoDuplicate           = errors.New("plan number already exists")
	ErrOrderIDRequired           = errors.New("order_id is required")
	ErrProductIDRequired         = errors.New("product_id is required")
	ErrInvalidQuantity           = errors.New("quantity must be greater than 0")
	ErrCreatedByRequired         = errors.New("created_by is required")
	ErrCannotStartPlan           = errors.New("can only start planned status")
	ErrCannotCompletePlan        = errors.New("can only complete in_progress status")
	ErrCannotCancelCompletedPlan = errors.New("cannot cancel completed plan")
	ErrCannotUpdateCompletedPlan = errors.New("cannot update completed plan")
	ErrCannotDeleteCompletedPlan = errors.New("cannot delete completed plan")
)