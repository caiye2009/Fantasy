package application

import "back/internal/user/domain"

// ToUser DTO → Domain Model (不包含 LoginID，由 Service 自动生成)
func ToUser(req *CreateUserRequest, loginID string) *domain.User {
	return &domain.User{
		LoginID:    loginID,
		Username:   req.Username,
		Department: req.Department,
		Email:      req.Email,
		Role:       req.Role,
		Status:     domain.UserStatusActive,
	}
}

// ToUserResponse Domain Model → DTO
func ToUserResponse(u *domain.User) *UserResponse {
	return &UserResponse{
		ID:          u.ID,
		LoginID:     u.LoginID,
		Username:    u.Username,
		Department:  u.Department,
		Email:       u.Email,
		Role:        u.Role,
		Status:      u.Status,
		HasInitPass: u.HasInitPass,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

// ToUserListResponse Domain Models → List DTO
func ToUserListResponse(users []*domain.User, total int64) *UserListResponse {
	responses := make([]*UserResponse, len(users))
	for i, u := range users {
		responses[i] = ToUserResponse(u)
	}
	
	return &UserListResponse{
		Total: total,
		Users: responses,
	}
}