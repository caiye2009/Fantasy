package process

import (
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
)

type ProcessHandler struct {
	processService *ProcessService
}

func NewProcessHandler(processService *ProcessService) *ProcessHandler {
	return &ProcessHandler{processService: processService}
}

// Create godoc
// @Summary      创建工序
// @Description  创建新的工序信息
// @Tags         工序管理
// @Accept       json
// @Produce      json
// @Param        request body CreateProcessRequest true "工序信息"
// @Success      200 {object} Process "创建成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /process [post]
func (h *ProcessHandler) Create(c *gin.Context) {
	var req CreateProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	process, err := h.processService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, process)
}

// Get godoc
// @Summary      获取工序详情
// @Description  根据工序ID获取工序详细信息
// @Tags         工序管理
// @Accept       json
// @Produce      json
// @Param        id path int true "工序ID"
// @Success      200 {object} Process "获取成功"
// @Failure      404 {object} map[string]string "工序不存在"
// @Security     Bearer
// @Router       /process/{id} [get]
func (h *ProcessHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	process, err := h.processService.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "工序不存在"})
		return
	}

	c.JSON(http.StatusOK, process)
}

// List godoc
// @Summary      获取工序列表
// @Description  获取所有工序列表
// @Tags         工序管理
// @Accept       json
// @Produce      json
// @Param        limit query int false "每页数量" default(10)
// @Param        offset query int false "偏移量" default(0)
// @Success      200 {object} map[string]interface{} "获取成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /process [get]
func (h *ProcessHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	processes, err := h.processService.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     len(processes),
		"processes": processes,
	})
}

// Update godoc
// @Summary      更新工序
// @Description  根据工序ID更新工序信息
// @Tags         工序管理
// @Accept       json
// @Produce      json
// @Param        id path int true "工序ID"
// @Param        request body UpdateProcessRequest true "更新的工序信息"
// @Success      200 {object} map[string]string "更新成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /process/{id} [put]
func (h *ProcessHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	var req UpdateProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.processService.Update(c.Request.Context(), uint(id), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// Delete godoc
// @Summary      删除工序
// @Description  根据工序ID删除工序
// @Tags         工序管理
// @Accept       json
// @Produce      json
// @Param        id path int true "工序ID"
// @Success      200 {object} map[string]string "删除成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /process/{id} [delete]
func (h *ProcessHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	if err := h.processService.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}