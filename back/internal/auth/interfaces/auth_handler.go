package interfaces

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"back/pkg/endpoint"
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
// @Failure      200 {object} map[string]string "请求参数错误"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req application.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 从 header 获取 source（h5 或 mobile）
	source := c.GetHeader("Source")
	if source == "" {
		source = "h5" // 默认值
	}

	resp, err := h.service.Login(c.Request.Context(), &req, source)
	if err != nil {
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

// Logout 用户登出
// @Summary      用户登出
// @Description  用户登出（清除 JWT 白名单）
// @Tags         认证
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string "登出成功"
// @Security     Bearer
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// 从 context 获取 loginId（auth 中间件已设置）
	loginID, exists := c.Get("loginId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// 从 header 获取 source
	source := c.GetHeader("Source")
	if source == "" {
		source = "h5" // 默认值
	}

	// 调用服务清除白名单
	if err := h.service.Logout(loginID.(string), source); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "登出失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

// RefreshToken 刷新 Token
// @Summary      刷新访问令牌
// @Description  使用 Refresh Token 获取新的 Access Token
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        request body application.RefreshTokenRequest true "Refresh Token"
// @Success      200 {object} application.RefreshTokenResponse "刷新成功"
// @Failure      200 {object} map[string]string "Token 无效"
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req application.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 从 header 获取 source（h5 或 mobile）
	source := c.GetHeader("Source")
	if source == "" {
		source = "h5" // 默认值
	}

	resp, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken, source)
	if err != nil {
		if errors.Is(err, application.ErrInvalidToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token 无效或已过期"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "刷新 Token 失败"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetPublicRoutes 获取公开路由定义（不需要认证）
func (h *AuthHandler) GetPublicRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/auth/login", Handler: h.Login, Domain: "auth", Action: "login"},
		{Method: "POST", Path: "/auth/refresh", Handler: h.RefreshToken, Domain: "auth", Action: "refresh"},
	}
}

// GetProtectedRoutes 获取受保护路由定义（需要认证）
func (h *AuthHandler) GetProtectedRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/auth/logout", Handler: h.Logout, Domain: "auth", Action: "logout"},
	}
}