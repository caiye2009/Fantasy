package domain

import "errors"

var (
	ErrUserNotFound              = errors.New("user not found")
	ErrLoginIDEmpty              = errors.New("login_id cannot be empty")
	ErrLoginIDInvalid            = errors.New("login_id must be between 4 and 20 characters")
	ErrLoginIDDuplicate          = errors.New("login_id already exists")
	ErrUsernameEmpty             = errors.New("username cannot be empty")
	ErrUsernameInvalid           = errors.New("username must be between 2 and 50 characters")
	ErrEmailInvalid              = errors.New("invalid email format")
	ErrPasswordRequired          = errors.New("password is required")
	ErrPasswordTooShort          = errors.New("password must be at least 6 characters")
	ErrPasswordIncorrect         = errors.New("incorrect password")
	ErrCurrentPasswordRequired   = errors.New("current password is required")
	ErrCurrentPasswordIncorrect  = errors.New("current password is incorrect")
	ErrUserAlreadyActive         = errors.New("user is already active")
	ErrUserAlreadySuspended      = errors.New("user is already suspended")
	ErrUserStatusInvalid         = errors.New("invalid user status")
	ErrCannotDeleteAdmin         = errors.New("cannot delete admin user")
	ErrInvalidRole               = errors.New("invalid user role")

	// Department errors
	ErrDepartmentNotFound      = errors.New("department not found")
	ErrDepartmentNameEmpty     = errors.New("department name cannot be empty")
	ErrDepartmentNameInvalid   = errors.New("department name must be between 2 and 100 characters")
	ErrDepartmentCodeInvalid   = errors.New("department code must be between 2 and 50 characters")
	ErrDepartmentCodeDuplicate = errors.New("department code already exists")

	// Role errors
	ErrRoleNotFound      = errors.New("role not found")
	ErrRoleNameEmpty     = errors.New("role name cannot be empty")
	ErrRoleNameInvalid   = errors.New("role name must be between 2 and 100 characters")
	ErrRoleCodeEmpty     = errors.New("role code cannot be empty")
	ErrRoleCodeInvalid   = errors.New("role code must be between 2 and 50 characters")
	ErrRoleCodeDuplicate = errors.New("role code already exists")
)