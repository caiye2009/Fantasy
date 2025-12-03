package user

import "back/pkg/fields"

type User struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	LoginID      string `gorm:"size:50;not null;uniqueIndex" json:"login_id"`
	Username     string `gorm:"size:100;not null" json:"username"`
	PasswordHash string `gorm:"size:255;not null" json:"-"`
	Email        string `gorm:"size:100" json:"email"`
	Role         string `gorm:"size:50;not null;index" json:"role"`
	Status       string `gorm:"size:20;default:'active';index" json:"status"`
	HasInitPass  bool   `gorm:"default:true" json:"has_init_pass"`
	fields.DTFields
}

func (User) TableName() string {
	return "users"
}

// CreateUserRequest 创建用户请求 (不接受密码)
type CreateUserRequest struct {
	LoginID  string `json:"login_id" binding:"required,min=4,max=20"`
	Username string `json:"username" binding:"required,min=2,max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
	Role     string `json:"role" binding:"required,oneof=admin hr sales follower assistant user"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,min=2,max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
	Role     string `json:"role" binding:"omitempty,oneof=admin hr sales follower assistant user"`
	Status   string `json:"status" binding:"omitempty,oneof=active inactive suspended"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required_without=IsInitPass"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

type UserResponse struct {
	ID          uint   `json:"id"`
	LoginID     string `json:"login_id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	Status      string `json:"status"`
	HasInitPass bool   `json:"has_init_pass"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:          u.ID,
		LoginID:     u.LoginID,
		Username:    u.Username,
		Email:       u.Email,
		Role:        u.Role,
		Status:      u.Status,
		HasInitPass: u.HasInitPass,
		CreatedAt:   u.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   u.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}