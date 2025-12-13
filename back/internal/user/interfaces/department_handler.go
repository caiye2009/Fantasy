package interfaces

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"back/internal/user/application"
	"back/internal/user/domain"
)

// DepartmentHandler 部门 Handler
type DepartmentHandler struct {
	service *application.DepartmentService
}

// NewDepartmentHandler 创建 Handler
func NewDepartmentHandler(service *application.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: service}
}

// Create 创建部门
// @Summary      创建部门
// @Description  创建新部门
// @Tags         部门管理
// @Accept       json
// @Produce      json
// @Param        request body application.CreateDepartmentRequest true "部门信息"
// @Success      200 {object} application.DepartmentResponse "创建成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /department [post]
func (h *DepartmentHandler) Create(c *gin.Context) {
	var req application.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Create(&req)
	if err != nil {
		if errors.Is(err, domain.ErrDepartmentCodeDuplicate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "部门编码已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Update 更新部门
// @Summary      更新部门
// @Description  更新部门信息
// @Tags         部门管理
// @Accept       json
// @Produce      json
// @Param        id path int true "部门ID"
// @Param        request body application.UpdateDepartmentRequest true "部门信息"
// @Success      200 {object} application.DepartmentResponse "更新成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      404 {object} map[string]string "部门不存在"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /department/{id} [put]
func (h *DepartmentHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req application.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Update(uint(id), &req)
	if err != nil {
		if errors.Is(err, domain.ErrDepartmentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "部门不存在"})
			return
		}
		if errors.Is(err, domain.ErrDepartmentCodeDuplicate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "部门编码已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Get 获取部门详情
// @Summary      获取部门详情
// @Description  根据部门ID获取部门详细信息
// @Tags         部门管理
// @Accept       json
// @Produce      json
// @Param        id path int true "部门ID"
// @Success      200 {object} application.DepartmentResponse "获取成功"
// @Failure      404 {object} map[string]string "部门不存在"
// @Security     Bearer
// @Router       /department/{id} [get]
func (h *DepartmentHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	resp, err := h.service.Get(uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrDepartmentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "部门不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// List 获取部门列表
// @Summary      获取部门列表
// @Description  获取部门列表（支持状态筛选和分页）
// @Tags         部门管理
// @Accept       json
// @Produce      json
// @Param        status query string false "部门状态（active/inactive）"
// @Param        page query int false "页码" default(1)
// @Param        page_size query int false "每页数量" default(10)
// @Success      200 {object} application.DepartmentListResponse "获取成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /department [get]
func (h *DepartmentHandler) List(c *gin.Context) {
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	resp, err := h.service.List(statusPtr, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Deactivate 停用部门
// @Summary      停用部门
// @Description  停用部门（软删除）
// @Tags         部门管理
// @Accept       json
// @Produce      json
// @Param        id path int true "部门ID"
// @Success      200 {object} map[string]string "停用成功"
// @Failure      404 {object} map[string]string "部门不存在"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /department/{id}/deactivate [put]
func (h *DepartmentHandler) Deactivate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.service.Deactivate(uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrDepartmentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "部门不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "停用成功"})
}

// Activate 激活部门
// @Summary      激活部门
// @Description  激活已停用的部门
// @Tags         部门管理
// @Accept       json
// @Produce      json
// @Param        id path int true "部门ID"
// @Success      200 {object} map[string]string "激活成功"
// @Failure      404 {object} map[string]string "部门不存在"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /department/{id}/activate [put]
func (h *DepartmentHandler) Activate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.service.Activate(uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrDepartmentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "部门不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "激活成功"})
}

// RegisterDepartmentHandlers 注册部门路由
func RegisterDepartmentHandlers(rg *gin.RouterGroup, service *application.DepartmentService) {
	handler := NewDepartmentHandler(service)

	rg.POST("/department", handler.Create)
	rg.GET("/department/:id", handler.Get)
	rg.GET("/department", handler.List)
	rg.PUT("/department/:id", handler.Update)
	rg.PUT("/department/:id/deactivate", handler.Deactivate)
	rg.PUT("/department/:id/activate", handler.Activate)
}

