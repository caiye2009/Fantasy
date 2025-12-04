package plan

import (
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
)

type PlanHandler struct {
	planService *PlanService
}

func NewPlanHandler(planService *PlanService) *PlanHandler {
	return &PlanHandler{planService: planService}
}

// Create godoc
// @Summary      创建计划
// @Description  创建新的生产计划
// @Tags         计划管理
// @Accept       json
// @Produce      json
// @Param        request body CreatePlanRequest true "计划信息"
// @Success      200 {object} Plan "创建成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /plan [post]
func (h *PlanHandler) Create(c *gin.Context) {
	var req CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan, err := h.planService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plan)
}

// Get godoc
// @Summary      获取计划详情
// @Description  根据计划ID获取计划详细信息
// @Tags         计划管理
// @Accept       json
// @Produce      json
// @Param        id path int true "计划ID"
// @Success      200 {object} Plan "获取成功"
// @Failure      404 {object} map[string]string "计划不存在"
// @Security     Bearer
// @Router       /plan/{id} [get]
func (h *PlanHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	plan, err := h.planService.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "计划不存在"})
		return
	}

	c.JSON(http.StatusOK, plan)
}

// List godoc
// @Summary      获取计划列表
// @Description  获取所有生产计划列表
// @Tags         计划管理
// @Accept       json
// @Produce      json
// @Param        limit query int false "每页数量" default(10)
// @Param        offset query int false "偏移量" default(0)
// @Success      200 {object} map[string]interface{} "获取成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /plan [get]
func (h *PlanHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	plans, err := h.planService.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": len(plans),
		"plans": plans,
	})
}

// Update godoc
// @Summary      更新计划
// @Description  根据计划ID更新计划信息
// @Tags         计划管理
// @Accept       json
// @Produce      json
// @Param        id path int true "计划ID"
// @Param        request body UpdatePlanRequest true "更新的计划信息"
// @Success      200 {object} map[string]string "更新成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /plan/{id} [put]
func (h *PlanHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	var req UpdatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.planService.Update(c.Request.Context(), uint(id), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// Delete godoc
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
	
	if err := h.planService.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}