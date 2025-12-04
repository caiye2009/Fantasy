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
	ErrCannotDeleteAdmin         = errors.New("cannot delete admin user")
)