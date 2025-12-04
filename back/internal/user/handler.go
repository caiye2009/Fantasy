package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *UserService
}

func NewUserHandler(userService *UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Create godoc
// @Summary      创建用户
// @Description  创建新用户并返回登录ID和默认密码
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        request body CreateUserRequest true "用户信息"
// @Success      200 {object} map[string]interface{} "创建成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /user [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "用户创建成功",
		"login_id": user.LoginID,
		"password": "123",
		"user":     user.ToResponse(),
	})
}

// Get godoc
// @Summary      获取用户详情
// @Description  根据用户ID获取用户详细信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id path int true "用户ID"
// @Success      200 {object} UserResponse "获取成功"
// @Failure      404 {object} map[string]string "用户不存在"
// @Security     Bearer
// @Router       /user/{id} [get]
func (h *UserHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	user, err := h.userService.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// List godoc
// @Summary      获取用户列表
// @Description  获取所有用户列表
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        limit query int false "每页数量" default(10)
// @Param        offset query int false "偏移量" default(0)
// @Success      200 {object} map[string]interface{} "获取成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /user [get]
func (h *UserHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	users, err := h.userService.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]*UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"total": len(responses),
		"users": responses,
	})
}

// Update godoc
// @Summary      更新用户信息
// @Description  根据用户ID更新用户信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id path int true "用户ID"
// @Param        request body UpdateUserRequest true "更新的用户信息"
// @Success      200 {object} map[string]string "更新成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /user/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.Update(c.Request.Context(), uint(id), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// Delete godoc
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

	if err := h.userService.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ChangePassword godoc
// @Summary      修改密码
// @Description  当前登录用户修改自己的密码
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        request body ChangePasswordRequest true "修改密码参数"
// @Success      200 {object} map[string]string "密码修改成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      401 {object} map[string]string "未登录"
// @Failure      404 {object} map[string]string "用户不存在"
// @Security     Bearer
// @Router       /user/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	loginID, exists := c.Get("login_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	user, err := h.userService.GetByLoginID(c.Request.Context(), loginID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.ChangePassword(c.Request.Context(), user.ID, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}