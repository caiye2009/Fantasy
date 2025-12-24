package interfaces

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"back/pkg/audit"
	"back/internal/process/application"
	"back/internal/process/domain"
)

// ProcessHandler 工序 Handler
type ProcessHandler struct {
	service *application.ProcessService
}

// NewProcessHandler 创建 Handler
func NewProcessHandler(service *application.ProcessService) *ProcessHandler {
	return &ProcessHandler{service: service}
}

// Create 创建工序
// @Summary      创建工序
// @Description  创建新的工序信息
// @Tags         工序管理
// @Accept       json
// @Produce      json
// @Param        request body application.CreateProcessRequest true "工序信息"
// @Success      200 {object} application.ProcessResponse "创建成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /process [post]
func (h *ProcessHandler) Create(c *gin.Context) {
	var req application.CreateProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	resp, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}

// Get 获取工序详情
// @Summary      获取工序详情
// @Description  根据工序ID获取工序详细信息
// @Tags         工序管理
// @Accept       json
// @Produce      json
// @Param        id path int true "工序ID"
// @Success      200 {object} application.ProcessResponse "获取成功"
// @Failure      404 {object} map[string]string "工序不存在"
// @Security     Bearer
// @Router       /process/{id} [get]
func (h *ProcessHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	resp, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrProcessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "工序不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}

// List 功能已迁移至 /search 接口
// 使用 POST /search 并指定 indices: ["processes"] 来获取工序列表

// Update 更新工序
// @Summary      更新工序
// @Description  根据工序ID更新工序信息
// @Tags         工序管理
// @Accept       json
// @Produce      json
// @Param        id path int true "工序ID"
// @Param        request body application.UpdateProcessRequest true "更新的工序信息"
// @Success      200 {object} map[string]string "更新成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /process/{id} [post]
func (h *ProcessHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	var req application.UpdateProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := h.service.Update(c.Request.Context(), uint(id), &req); err != nil {
		if errors.Is(err, domain.ErrProcessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "工序不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// Delete 删除工序
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
	
	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrProcessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "工序不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// RegisterProcessHandlers 注册路由
func RegisterProcessHandlers(rg *gin.RouterGroup, service *application.ProcessService) {
	handler := NewProcessHandler(service)

	rg.POST("/process", audit.Mark("process", "processCreation"), handler.Create)
	rg.GET("/process/:id", handler.Get)
	// List 接口已移除，使用 POST /search 替代
	rg.POST("/process/:id", audit.Mark("process", "processUpdate"), handler.Update)
	rg.DELETE("/process/:id", audit.Mark("process", "processDeletion"), handler.Delete)
}