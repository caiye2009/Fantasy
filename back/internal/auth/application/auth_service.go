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
	GenerateAccessToken(loginID, role string) (string, error)
	GenerateRefreshToken(loginID string) (string, error)
	ValidateRefreshToken(token string) (loginID string, err error)
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
	// 1. 验证用户名密码
	user, err := s.userService.ValidatePassword(ctx, req.LoginID, req.Password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	
	// 2. 检查用户状态
	if !user.IsActive() {
		return nil, ErrAccountSuspended
	}
	
	// 3. 生成 Access Token（包含 loginID 和 role）
	accessToken, err := s.jwtGen.GenerateAccessToken(user.LoginID, string(user.Role))
	if err != nil {
		return nil, ErrTokenGenerateFailed
	}
	
	// 4. 生成 Refresh Token（只包含 loginID）
	refreshToken, err := s.jwtGen.GenerateRefreshToken(user.LoginID)
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
	loginID, err := s.jwtGen.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}
	
	// 2. 获取用户信息（需要 role 生成新 token）
	user, err := s.userService.GetByLoginID(ctx, loginID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	
	// 3. 检查用户状态
	if !user.IsActive() {
		return nil, ErrAccountSuspended
	}
	
	// 4. 生成新的 Access Token
	newAccessToken, err := s.jwtGen.GenerateAccessToken(user.LoginID, string(user.Role))
	if err != nil {
		return nil, ErrTokenGenerateFailed
	}
	
	return &RefreshTokenResponse{
		AccessToken: newAccessToken,
	}, nil
}