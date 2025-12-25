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
	userService      UserServiceInterface
	jwtWang          *auth.JWTWang
	whitelistManager *auth.WhitelistManager
}

// NewAuthService 创建认证服务
func NewAuthService(userService UserServiceInterface, jwtWang *auth.JWTWang, whitelistManager *auth.WhitelistManager) *AuthService {
	return &AuthService{
		userService:      userService,
		jwtWang:          jwtWang,
		whitelistManager: whitelistManager,
	}
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, req *LoginRequest, source string) (*LoginResponse, error) {
	// 1. 验证用户名密码
	user, err := s.userService.ValidatePassword(ctx, req.LoginID, req.Password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// 2. 检查用户状态
	if !user.IsActive() {
		return nil, ErrAccountSuspended
	}

	// 3. 生成 Access Token（包含 loginID、userName、role 和 department）
	accessToken, jti, err := s.jwtWang.GenerateAccessToken(user.LoginID, user.Username, string(user.Role), user.Department)
	if err != nil {
		return nil, ErrTokenGenerateFailed
	}

	// 4. 添加到白名单（一个 loginId + source 只能有一个有效 JWT）
	if s.whitelistManager != nil {
		if err := s.whitelistManager.AddToWhitelist(user.LoginID, source, jti); err != nil {
			// 白名单添加失败，记录日志但不影响登录流程
			println("Failed to add JWT to whitelist:", err.Error())
		}
	}

	// 5. 生成 Refresh Token（只包含 loginID）
	refreshToken, err := s.jwtWang.GenerateRefreshToken(user.LoginID)
	if err != nil {
		return nil, ErrTokenGenerateFailed
	}

	// 6. 构造扁平化响应（直接返回 username、role 和 department）
	return &LoginResponse{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		Username:              user.Username,
		Role:                  string(user.Role),
		Department:            user.Department,
		RequirePasswordChange: user.HasInitPass,
	}, nil
}

// RefreshToken 刷新 Token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string, source string) (*RefreshTokenResponse, error) {
	// 1. 验证 Refresh Token
	loginID, err := s.jwtWang.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// 2. 获取用户信息（需要 userName、role 和 department 生成新 token）
	user, err := s.userService.GetByLoginID(ctx, loginID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// 3. 检查用户状态
	if !user.IsActive() {
		return nil, ErrAccountSuspended
	}

	// 4. 生成新的 Access Token（包含 department）
	newAccessToken, jti, err := s.jwtWang.GenerateAccessToken(user.LoginID, user.Username, string(user.Role), user.Department)
	if err != nil {
		return nil, ErrTokenGenerateFailed
	}

	// 5. 更新白名单（替换旧的 JWT）
	if s.whitelistManager != nil {
		if err := s.whitelistManager.AddToWhitelist(user.LoginID, source, jti); err != nil {
			// 白名单更新失败，记录日志但不影响刷新流程
			println("Failed to update JWT whitelist:", err.Error())
		}
	}

	return &RefreshTokenResponse{
		AccessToken: newAccessToken,
	}, nil
}

// Logout 用户登出
func (s *AuthService) Logout(loginID string, source string) error {
	if s.whitelistManager != nil {
		return s.whitelistManager.RemoveFromWhitelist(loginID, source)
	}
	return nil
}