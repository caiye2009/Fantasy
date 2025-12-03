package auth

import (
	"back/pkg/fields"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *AuthService
}

func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login godoc
// @Summary      用户登录
// @Description  通过登录ID和密码进行用户认证
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "登录请求参数"
// @Success      200 {object} fields.Response{data=LoginResponse} "登录成功"
// @Failure      400 {object} fields.Response "请求参数错误"
// @Failure      401 {object} fields.Response "认证失败"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fields.Error(c, err.Error())
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		fields.Error(c, err.Error())
		return
	}

	fields.Success(c, resp)
}