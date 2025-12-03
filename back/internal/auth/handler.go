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