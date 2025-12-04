package material

import (
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
)

type MaterialHandler struct {
	materialService *MaterialService
}

func NewMaterialHandler(materialService *MaterialService) *MaterialHandler {
	return &MaterialHandler{materialService: materialService}
}

// Create 创建材料
// @Summary      创建材料
// @Description  创建新的材料信息
// @Tags         材料管理
// @Accept       json
// @Produce      json
// @Param        request body CreateMaterialRequest true "材料信息"
// @Success      200 {object} Material "创建成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /material [post]
func (h *MaterialHandler) Create(c *gin.Context) {
	var req CreateMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	material, err := h.materialService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, material)
}

// Get 获取材料详情
// @Summary      获取材料详情
// @Description  根据材料ID获取材料详细信息
// @Tags         材料管理
// @Accept       json
// @Produce      json
// @Param        id path int true "材料ID"
// @Success      200 {object} Material "获取成功"
// @Failure      404 {object} map[string]string "材料不存在"
// @Security     Bearer
// @Router       /material/{id} [get]
func (h *MaterialHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	material, err := h.materialService.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "材料不存在"})
		return
	}

	c.JSON(http.StatusOK, material)
}

// List 获取材料列表
// @Summary      获取材料列表
// @Description  获取所有材料列表
// @Tags         材料管理
// @Accept       json
// @Produce      json
// @Param        limit query int false "每页数量" default(10)
// @Param        offset query int false "偏移量" default(0)
// @Success      200 {object} map[string]interface{} "获取成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /material [get]
func (h *MaterialHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	materials, err := h.materialService.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     len(materials),
		"materials": materials,
	})
}

// Update 更新材料信息
// @Summary      更新材料信息
// @Description  根据材料ID更新材料信息
// @Tags         材料管理
// @Accept       json
// @Produce      json
// @Param        id path int true "材料ID"
// @Param        request body UpdateMaterialRequest true "更新的材料信息"
// @Success      200 {object} map[string]string "更新成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /material/{id} [put]
func (h *MaterialHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	var req UpdateMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.materialService.Update(c.Request.Context(), uint(id), &req); err != nil {
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
	
	if err := h.materialService.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}