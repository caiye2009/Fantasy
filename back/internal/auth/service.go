package auth

import (
	"errors"
	"back/internal/user"
	"back/pkg/auth"
)

type AuthService struct {
	userService *user.UserService
	jwtWang     *auth.JWTWang
}

func NewAuthService(userService *user.UserService, jwtWang *auth.JWTWang) *AuthService {
	return &AuthService{
		userService: userService,
		jwtWang:     jwtWang,
	}
}

func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 根据 login_id 查询用户
	u, err := s.userService.GetByLoginID(req.LoginID)
	if err != nil {
		return nil, errors.New("工号或密码错误")
	}

	// 验证密码
	if err := s.userService.ValidatePassword(u, req.Password); err != nil {
		return nil, errors.New("工号或密码错误")
	}

	// 检查用户状态
	if u.Status != "active" {
		return nil, errors.New("账号已被停用")
	}

	// 生成 JWT (使用 login_id)
	token, err := s.jwtWang.GenerateToken(u.LoginID, u.Role)
	if err != nil {
		return nil, errors.New("生成token失败")
	}

	return &LoginResponse{
		AccessToken: token,
		User: &UserInfo{
			ID:       u.ID,
			LoginID:  u.LoginID,
			Username: u.Username,
			Email:    u.Email,
			Role:     u.Role,
		},
		RequirePasswordChange: u.HasInitPass,
	}, nil
}