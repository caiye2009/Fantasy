package auth

type LoginRequest struct {
    LoginID  string `json:"loginId" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
    AccessToken           string    `json:"accessToken"`
    User                  *UserInfo `json:"user"`
    RequirePasswordChange bool      `json:"requirePasswordChange"`
}

type UserInfo struct {
    ID       uint   `json:"id"`
    LoginID  string `json:"loginId"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Role     string `json:"role"`
}
