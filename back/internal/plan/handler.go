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
// @Summary      创建生产计划
// @Description  创建新的生产计划
// @Tags         生产计划管理
// @Accept       json
// @Produce      json
// @Param        request body Plan true "生产计划信息"
// @Success      200 {object} Plan "创建成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /plan [post]
func (h *PlanHandler) Create(c *gin.Context) {
    var p Plan
    if err := c.ShouldBindJSON(&p); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.planService.Create(&p); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, p)
}

// Get godoc
// @Summary      获取生产计划详情
// @Description  根据生产计划ID获取详细信息
// @Tags         生产计划管理
// @Accept       json
// @Produce      json
// @Param        id path int true "生产计划ID"
// @Success      200 {object} Plan "获取成功"
// @Failure      404 {object} map[string]string "生产计划不存在"
// @Security     Bearer
// @Router       /plan/{id} [get]
func (h *PlanHandler) Get(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    p, err := h.planService.Get(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }
    c.JSON(http.StatusOK, p)
}

// List godoc
// @Summary      获取生产计划列表
// @Description  获取所有生产计划列表
// @Tags         生产计划管理
// @Accept       json
// @Produce      json
// @Success      200 {array} Plan "获取成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /plan [get]
func (h *PlanHandler) List(c *gin.Context) {
    list, err := h.planService.List()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, list)
}

// Update godoc
// @Summary      更新生产计划
// @Description  根据生产计划ID更新信息
// @Tags         生产计划管理
// @Accept       json
// @Produce      json
// @Param        id path int true "生产计划ID"
// @Param        request body map[string]interface{} true "更新的生产计划信息"
// @Success      200 {object} map[string]string "更新成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /plan/{id} [put]
func (h *PlanHandler) Update(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))

    data := make(map[string]interface{})
    if err := c.ShouldBindJSON(&data); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.planService.Update(uint(id), data); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// Delete godoc
// @Summary      删除生产计划
// @Description  根据生产计划ID删除
// @Tags         生产计划管理
// @Accept       json
// @Produce      json
// @Param        id path int true "生产计划ID"
// @Success      200 {object} map[string]string "删除成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /plan/{id} [delete]
func (h *PlanHandler) Delete(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))

    if err := h.planService.Delete(uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}