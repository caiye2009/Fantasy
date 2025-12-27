package interfaces

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"back/pkg/endpoint"
	"back/internal/plan/application"
	"back/internal/plan/domain"
)

// PlanHandler 计划 Handler
type PlanHandler struct {
	service *application.PlanService
}

// NewPlanHandler 创建 Handler
func NewPlanHandler(service *application.PlanService) *PlanHandler {
	return &PlanHandler{service: service}
}

// Create 创建计划
// @Summary      创建计划
// @Description  创建新的生产计划
// @Tags         计划管理
// @Accept       json
// @Produce      json
// @Param        request body application.CreatePlanRequest true "计划信息"
// @Success      200 {object} application.PlanResponse "创建成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /plan [post]
func (h *PlanHandler) Create(c *gin.Context) {
	var req application.CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	resp, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, domain.ErrPlanNoDuplicate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "计划编号已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}

// Get 获取计划详情
// @Summary      获取计划详情
// @Description  根据计划ID获取计划详细信息
// @Tags         计划管理
// @Accept       json
// @Produce      json
// @Param        id path int true "计划ID"
// @Success      200 {object} application.PlanResponse "获取成功"
// @Failure      404 {object} map[string]string "计划不存在"
// @Security     Bearer
// @Router       /plan/{id} [get]
func (h *PlanHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	resp, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrPlanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "计划不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}

// List 获取计划列表
// @Summary      获取计划列表
// @Description  获取所有生产计划列表
// @Tags         计划管理
// @Accept       json
// @Produce      json
// @Param        limit query int false "每页数量" default(10)
// @Param        offset query int false "偏移量" default(0)
// @Success      200 {object} application.PlanListResponse "获取成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /plan [get]
func (h *PlanHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	resp, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}

// Update 更新计划
// @Summary      更新计划
// @Description  根据计划ID更新计划信息
// @Tags         计划管理
// @Accept       json
// @Produce      json
// @Param        id path int true "计划ID"
// @Param        request body application.UpdatePlanRequest true "更新的计划信息"
// @Success      200 {object} map[string]string "更新成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /plan/{id} [post]
func (h *PlanHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	var req application.UpdatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := h.service.Update(c.Request.Context(), uint(id), &req); err != nil {
		if errors.Is(err, domain.ErrPlanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "计划不存在"})
			return
		}
		if errors.Is(err, domain.ErrCannotStartPlan) ||
			errors.Is(err, domain.ErrCannotCompletePlan) ||
			errors.Is(err, domain.ErrCannotUpdateCompletedPlan) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// Delete 删除计划
// @Summary      删除计划
// @Description  根据计划ID删除计划
// @Tags         计划管理
// @Accept       json
// @Produce      json
// @Param        id path int true "计划ID"
// @Success      200 {object} map[string]string "删除成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /plan/{id} [delete]
func (h *PlanHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrPlanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "计划不存在"})
			return
		}
		if errors.Is(err, domain.ErrCannotDeleteCompletedPlan) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "已完成的计划不能删除"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetRoutes 返回路由定义
func (h *PlanHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/plan", Handler: h.Create, Domain: "plan", Action: "create"},
		{Method: "GET", Path: "/plan/:id", Handler: h.Get, Domain: "", Action: ""},
		{Method: "GET", Path: "/plan", Handler: h.List, Domain: "", Action: ""},
		{Method: "PUT", Path: "/plan/:id", Handler: h.Update, Domain: "plan", Action: "update"},
		{Method: "DELETE", Path: "/plan/:id", Handler: h.Delete, Domain: "plan", Action: "delete"},
	}
}