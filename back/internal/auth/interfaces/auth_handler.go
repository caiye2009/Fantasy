package interfaces

import (
	"errors"
	"net/http"
	
	"github.com/gin-gonic/gin"
	
	"back/internal/auth/application"
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
// @Success      200 {object} application.LoginResponse "登录成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      401 {object} map[string]string "认证失败"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req application.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}
	
	resp, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		// 根据错误类型返回不同的状态码
		if errors.Is(err, application.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, application.ErrAccountSuspended) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败"})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}

// Logout 用户登出（可选）
// @Summary      用户登出
// @Description  用户登出（前端清除 Token）
// @Tags         认证
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string "登出成功"
// @Security     Bearer
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// JWT 是无状态的，登出通常在前端完成
	// 如果需要服务端维护黑名单，可以在这里实现
	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

// RefreshToken 刷新 Token（可选）
// @Summary      刷新访问令牌
// @Description  使用旧 Token 刷新获取新 Token
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        request body map[string]string true "旧 Token"
// @Success      200 {object} map[string]string "刷新成功"
// @Failure      401 {object} map[string]string "Token 无效"
// @Security     Bearer
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}
	
	newToken, err := h.service.RefreshToken(c.Request.Context(), req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token 无效"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"access_token": newToken,
	})
}

// RegisterAuthHandlers 注册路由
func RegisterAuthHandlers(rg *gin.RouterGroup, service *application.AuthService) {
	handler := NewAuthHandler(service)
	
	// 公开路由（不需要认证）
	rg.POST("/auth/login", handler.Login)
	rg.POST("/auth/logout", handler.Logout)
	rg.POST("/auth/refresh", handler.RefreshToken)
}