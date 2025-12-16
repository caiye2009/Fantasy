package interfaces

import (
	"errors"
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
	
	"back/internal/material/application"
	"back/internal/material/domain"
)

// MaterialHandler 材料 Handler
type MaterialHandler struct {
	service *application.MaterialService
}

// NewMaterialHandler 创建 Handler
func NewMaterialHandler(service *application.MaterialService) *MaterialHandler {
	return &MaterialHandler{service: service}
}

// Create 创建材料
// @Summary      创建材料
// @Description  创建新的材料信息
// @Tags         材料管理
// @Accept       json
// @Produce      json
// @Param        request body application.CreateMaterialRequest true "材料信息"
// @Success      200 {object} application.MaterialResponse "创建成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /material [post]
func (h *MaterialHandler) Create(c *gin.Context) {
	var req application.CreateMaterialRequest
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

// Get 获取材料详情
// @Summary      获取材料详情
// @Description  根据材料ID获取材料详细信息
// @Tags         材料管理
// @Accept       json
// @Produce      json
// @Param        id path int true "材料ID"
// @Success      200 {object} application.MaterialResponse "获取成功"
// @Failure      404 {object} map[string]string "材料不存在"
// @Security     Bearer
// @Router       /material/{id} [get]
func (h *MaterialHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	resp, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrMaterialNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "材料不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}

// Update 更新材料信息
// @Summary      更新材料信息
// @Description  根据材料ID更新材料信息
// @Tags         材料管理
// @Accept       json
// @Produce      json
// @Param        id path int true "材料ID"
// @Param        request body application.UpdateMaterialRequest true "更新的材料信息"
// @Success      200 {object} map[string]string "更新成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /material/{id} [post]
func (h *MaterialHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	var req application.UpdateMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := h.service.Update(c.Request.Context(), uint(id), &req); err != nil {
		if errors.Is(err, domain.ErrMaterialNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "材料不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// Delete 删除材料
// @Summary      删除材料
// @Description  根据材料ID删除材料
// @Tags         材料管理
// @Accept       json
// @Produce      json
// @Param        id path int true "材料ID"
// @Success      200 {object} map[string]string "删除成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /material/{id} [delete]
func (h *MaterialHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrMaterialNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "材料不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// RegisterMaterialHandlers 注册路由
func RegisterMaterialHandlers(rg *gin.RouterGroup, service *application.MaterialService) {
	handler := NewMaterialHandler(service)

	rg.POST("/material", handler.Create)
	rg.GET("/material/:id", handler.Get)
	// List 接口已移除，使用 POST /search 替代
	rg.POST("/material/:id", handler.Update)
	rg.DELETE("/material/:id", handler.Delete)
}