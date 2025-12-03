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

	process, err := h.processService.Create(&req)
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
	p, err := h.processService.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

// List godoc
// @Summary      获取工序列表
// @Description  获取所有工序列表
// @Tags         工序管理
// @Accept       json
// @Produce      json
// @Success      200 {array} Process "获取成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /process [get]
func (h *ProcessHandler) List(c *gin.Context) {
	list, err := h.processService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// Update godoc
// @Summary      更新工序信息
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

	if err := h.processService.Update(uint(id), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
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

	if err := h.processService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}