package application

import (
	"context"
	
	userDomain "back/internal/user/domain"
)

// UserServiceInterface User 服务接口
type UserServiceInterface interface {
	GetByLoginID(ctx context.Context, loginID string) (*userDomain.User, error)
	ValidatePassword(ctx context.Context, loginID, password string) (*userDomain.User, error)
}

// JWTGenerator JWT 生成器接口
type JWTGenerator interface {
	GenerateToken(loginID, role string) (string, error)
}

// AuthService 认证应用服务
type AuthService struct {
	userService UserServiceInterface
	jwtGen      JWTGenerator
}

// NewAuthService 创建认证服务
func NewAuthService(userService UserServiceInterface, jwtGen JWTGenerator) *AuthService {
	return &AuthService{
		userService: userService,
		jwtGen:      jwtGen,
	}
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 1. 验证用户名密码（通过 UserService）
	user, err := s.userService.ValidatePassword(ctx, req.LoginID, req.Password)
	if err != nil {
		// 不暴露具体错误原因（安全考虑）
		return nil, ErrInvalidCredentials
	}
	
	// 2. 检查用户状态
	if !user.IsActive() {
		return nil, ErrAccountSuspended
	}
	
	// 3. 生成 JWT Token
	token, err := s.jwtGen.GenerateToken(user.LoginID, string(user.Role))
	if err != nil {
		return nil, ErrTokenGenerateFailed
	}
	
	// 4. 构造响应
	return &LoginResponse{
		AccessToken: token,
		User: &UserInfo{
			ID:       user.ID,
			LoginID:  user.LoginID,
			Username: user.Username,
			Email:    user.Email,
			Role:     string(user.Role),
		},
		RequirePasswordChange: user.HasInitPass,
	}, nil
}

// ValidateToken 验证 Token（可选，用于中间件）
func (s *AuthService) ValidateToken(token string) (string, error) {
	// 如果需要在这里验证 Token，可以注入 JWT 验证器
	// 通常在中间件层完成
	return "", nil
}

// RefreshToken 刷新 Token（可选）
func (s *AuthService) RefreshToken(ctx context.Context, oldToken string) (string, error) {
	// 实现 Token 刷新逻辑
	return "", nil
}