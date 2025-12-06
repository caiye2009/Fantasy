package application

// JWTWangAdapter JWT 工具适配器
type JWTWangAdapter struct {
	jwtWang JWTWangInterface
}

// JWTWangInterface pkg/auth.JWTWang 的接口抽象
type JWTWangInterface interface {
	GenerateAccessToken(loginID, role string) (string, error)
	GenerateRefreshToken(loginID string) (string, error)
	ValidateRefreshToken(token string) (string, error)
}

// NewJWTWangAdapter 创建适配器
func NewJWTWangAdapter(jwtWang JWTWangInterface) JWTGenerator {
	return &JWTWangAdapter{jwtWang: jwtWang}
}

// GenerateAccessToken 生成 Access Token
func (a *JWTWangAdapter) GenerateAccessToken(loginID, role string) (string, error) {
	return a.jwtWang.GenerateAccessToken(loginID, role)
}

// GenerateRefreshToken 生成 Refresh Token
func (a *JWTWangAdapter) GenerateRefreshToken(loginID string) (string, error) {
	return a.jwtWang.GenerateRefreshToken(loginID)
}

// ValidateRefreshToken 验证 Refresh Token
func (a *JWTWangAdapter) ValidateRefreshToken(token string) (string, error) {
	return a.jwtWang.ValidateRefreshToken(token)
}