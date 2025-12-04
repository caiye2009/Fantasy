package application

// LoginRequest 登录请求
type LoginRequest struct {
	LoginID  string `json:"login_id" binding:"required"`
	Password string `json:"password" binding:"required,min=3"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken           string    `json:"access_token"`
	User                  *UserInfo `json:"user"`
	RequirePasswordChange bool      `json:"require_password_change"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID       uint   `json:"id"`
	LoginID  string `json:"login_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}