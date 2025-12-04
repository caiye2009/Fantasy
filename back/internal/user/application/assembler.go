package application

import "back/internal/user/domain"

// ToUser DTO → Domain Model
func ToUser(req *CreateUserRequest) *domain.User {
	return &domain.User{
		LoginID:  req.LoginID,
		Username: req.Username,
		Email:    req.Email,
		Role:     domain.UserRole(req.Role),
		Status:   domain.UserStatusActive,
	}
}

// ToUserResponse Domain Model → DTO
func ToUserResponse(u *domain.User) *UserResponse {
	return &UserResponse{
		ID:          u.ID,
		LoginID:     u.LoginID,
		Username:    u.Username,
		Email:       u.Email,
		Role:        string(u.Role),
		Status:      string(u.Status),
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