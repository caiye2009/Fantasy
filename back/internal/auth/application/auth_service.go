package application

import (
	"context"

	"back/pkg/auth"
	userDomain "back/internal/user/domain"
)

// UserServiceInterface User 服务接口
type UserServiceInterface interface {
	GetByLoginID(ctx context.Context, loginID string) (*userDomain.User, error)
	ValidatePassword(ctx context.Context, loginID, password string) (*userDomain.User, error)
}

// AuthService 认证应用服务
type AuthService struct {
	userService UserServiceInterface
	jwtWang     *auth.JWTWang
}

// NewAuthService 创建认证服务
func NewAuthService(userService UserServiceInterface, jwtWang *auth.JWTWang) *AuthService {
	return &AuthService{
		userService: userService,
		jwtWang:     jwtWang,
	}
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 1. 验证用户名密码
	user, err := s.userService.ValidatePassword(ctx, req.LoginID, req.Password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	
	// 2. 检查用户状态
	if !user.IsActive() {
		return nil, ErrAccountSuspended
	}
	
	// 3. 生成 Access Token（包含 loginID、userName 和 role）
	accessToken, err := s.jwtWang.GenerateAccessToken(user.LoginID, user.Username, string(user.Role))
	if err != nil {
		return nil, ErrTokenGenerateFailed
	}

	// 4. 生成 Refresh Token（只包含 loginID）
	refreshToken, err := s.jwtWang.GenerateRefreshToken(user.LoginID)
	if err != nil {
		return nil, ErrTokenGenerateFailed
	}
	
	// 5. 构造扁平化响应（直接返回 username 和 role）
	return &LoginResponse{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		Username:              user.Username,
		Role:                  string(user.Role),
		RequirePasswordChange: user.HasInitPass,
	}, nil
}

// RefreshToken 刷新 Token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error) {
	// 1. 验证 Refresh Token
	loginID, err := s.jwtWang.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// 2. 获取用户信息（需要 userName 和 role 生成新 token）
	user, err := s.userService.GetByLoginID(ctx, loginID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// 3. 检查用户状态
	if !user.IsActive() {
		return nil, ErrAccountSuspended
	}

	// 4. 生成新的 Access Token
	newAccessToken, err := s.jwtWang.GenerateAccessToken(user.LoginID, user.Username, string(user.Role))
	if err != nil {
		return nil, ErrTokenGenerateFailed
	}

	return &RefreshTokenResponse{
		AccessToken: newAccessToken,
	}, nil
}