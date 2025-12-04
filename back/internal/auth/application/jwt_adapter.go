package application

// JWTWangAdapter JWT 工具适配器
type JWTWangAdapter struct {
	jwtWang JWTWangInterface
}

// JWTWangInterface pkg/auth.JWTWang 的接口抽象
type JWTWangInterface interface {
	GenerateToken(loginID, role string) (string, error)
}

// NewJWTWangAdapter 创建适配器
func NewJWTWangAdapter(jwtWang JWTWangInterface) JWTGenerator {
	return &JWTWangAdapter{jwtWang: jwtWang}
}

// GenerateToken 生成 Token
func (a *JWTWangAdapter) GenerateToken(loginID, role string) (string, error) {
	return a.jwtWang.GenerateToken(loginID, role)
}