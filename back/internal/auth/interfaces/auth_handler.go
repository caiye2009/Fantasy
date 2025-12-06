package interfaces

import (
	"errors"	
	"github.com/gin-gonic/gin"
	
	"back/internal/auth/application"
	"back/pkg/fields"
)

// AuthHandler 认证 Handler
type AuthHandler struct {
	service *application.AuthService
}

// NewAuthHandler 创建 Handler
func NewAuthHandler(service *application.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Login 用户登录
// @Summary      用户登录
// @Description  通过登录ID和密码进行用户认证
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        request body application.LoginRequest true "登录请求参数"
// @Success      200 {object} fields.Response{data=application.LoginResponse} "登录成功"
// @Failure      200 {object} fields.Response "请求参数错误"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req application.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fields.Error(c, "请求参数错误")
		return
	}
	
	resp, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, application.ErrInvalidCredentials) {
			fields.Error(c, err.Error())
			return
		}
		if errors.Is(err, application.ErrAccountSuspended) {
			fields.Error(c, err.Error())
			return
		}
		fields.Error(c, "登录失败")
		return
	}
	
	fields.Success(c, resp)
}

// Logout 用户登出
// @Summary      用户登出
// @Description  用户登出（前端清除 Token）
// @Tags         认证
// @Accept       json
// @Produce      json
// @Success      200 {object} fields.Response "登出成功"
// @Security     Bearer
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	fields.Success(c, gin.H{"message": "登出成功"})
}

// RefreshToken 刷新 Token
// @Summary      刷新访问令牌
// @Description  使用 Refresh Token 获取新的 Access Token
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        request body application.RefreshTokenRequest true "Refresh Token"
// @Success      200 {object} fields.Response{data=application.RefreshTokenResponse} "刷新成功"
// @Failure      200 {object} fields.Response "Token 无效"
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req application.RefreshTokenRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		fields.Error(c, "请求参数错误")
		return
	}
	
	resp, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, application.ErrInvalidToken) {
			fields.Error(c, "Token 无效或已过期")
			return
		}
		fields.Error(c, "刷新 Token 失败")
		return
	}
	
	fields.Success(c, resp)
}

// RegisterAuthHandlers 注册路由
func RegisterAuthHandlers(rg *gin.RouterGroup, service *application.AuthService) {
	handler := NewAuthHandler(service)
	
	rg.POST("/auth/login", handler.Login)
	rg.POST("/auth/logout", handler.Logout)
	rg.POST("/auth/refresh", handler.RefreshToken)
}