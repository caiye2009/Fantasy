package interfaces

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"back/pkg/endpoint"
	"back/internal/user/application"
	"back/internal/user/domain"
)

// UserHandler 用户 Handler
type UserHandler struct {
	service *application.UserService
}

// NewUserHandler 创建 Handler
func NewUserHandler(service *application.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Create 创建用户
// @Summary      创建用户
// @Description  创建新用户并返回登录ID和默认密码
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        request body application.CreateUserRequest true "用户信息"
// @Success      200 {object} application.CreateUserResponse "创建成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /user [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req application.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	resp, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, domain.ErrLoginIDDuplicate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "工号已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}

// Get 获取用户详情
// @Summary      获取用户详情
// @Description  根据用户ID获取用户详细信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id path int true "用户ID"
// @Success      200 {object} application.UserResponse "获取成功"
// @Failure      404 {object} map[string]string "用户不存在"
// @Security     Bearer
// @Router       /user/{id} [get]
func (h *UserHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	resp, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}

// List 获取用户列表
// @Summary      获取用户列表
// @Description  获取所有用户列表
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        limit query int false "每页数量" default(10)
// @Param        offset query int false "偏移量" default(0)
// @Success      200 {object} application.UserListResponse "获取成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /user [get]
func (h *UserHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	resp, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}

// Update 更新用户信息
// @Summary      更新用户信息
// @Description  根据用户ID更新用户信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id path int true "用户ID"
// @Param        request body application.UpdateUserRequest true "更新的用户信息"
// @Success      200 {object} map[string]string "更新成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /user/{id} [post]
func (h *UserHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	var req application.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := h.service.Update(c.Request.Context(), uint(id), &req); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// Delete 删除用户
// @Summary      删除用户
// @Description  根据用户ID删除用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id path int true "用户ID"
// @Success      200 {object} map[string]string "删除成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /user/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			return
		}
		if errors.Is(err, domain.ErrCannotDeleteAdmin) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "不能删除管理员用户"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ChangePassword 修改密码
// @Summary      修改密码
// @Description  当前登录用户修改自己的密码
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        request body application.ChangePasswordRequest true "修改密码参数"
// @Success      200 {object} map[string]string "密码修改成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      401 {object} map[string]string "未登录"
// @Failure      404 {object} map[string]string "用户不存在"
// @Security     Bearer
// @Router       /user/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	// 1. 从 Context 获取登录用户 ID
	loginID, exists := c.Get("login_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	
	// 2. 查询用户
	user, err := h.service.GetByLoginID(c.Request.Context(), loginID.(string))
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// 3. 解析请求
	var req application.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 4. 修改密码
	if err := h.service.ChangePassword(c.Request.Context(), user.ID, &req); err != nil {
		if errors.Is(err, domain.ErrCurrentPasswordIncorrect) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "当前密码错误"})
			return
		}
		if errors.Is(err, domain.ErrPasswordTooShort) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "密码长度至少6位"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

// GetRoutes 获取路由定义
func (h *UserHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/user", Handler: h.Create, Domain: "user", Action: "create"},
		{Method: "GET", Path: "/user/:id", Handler: h.Get, Domain: "", Action: ""},
		{Method: "GET", Path: "/user", Handler: h.List, Domain: "", Action: ""},
		{Method: "PUT", Path: "/user/:id", Handler: h.Update, Domain: "user", Action: "update"},
		{Method: "DELETE", Path: "/user/:id", Handler: h.Delete, Domain: "user", Action: "delete"},
		{Method: "PUT", Path: "/user/:id/password", Handler: h.ChangePassword, Domain: "user", Action: "changePassword"},
	}
}