package application

// LoginRequest 登录请求
type LoginRequest struct {
	LoginID  string `json:"loginId" binding:"required"`
	Password string `json:"password" binding:"required,min=3"`
}

// LoginResponse 登录响应（扁平化）
type LoginResponse struct {
	AccessToken           string `json:"accessToken"`           // 改为小驼峰
	RefreshToken          string `json:"refreshToken"`          // 新增
	Username              string `json:"username"`              // 扁平化，直接返回用户名
	Role                  string `json:"role"`                  // 扁平化，直接返回角色
	RequirePasswordChange bool   `json:"requirePasswordChange"` // 改为小驼峰
}

// RefreshTokenRequest 刷新 Token 请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// RefreshTokenResponse 刷新 Token 响应
type RefreshTokenResponse struct {
	AccessToken string `json:"accessToken"`
}